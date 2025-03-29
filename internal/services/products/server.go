package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

// Configuration for the load test
type Config struct {
	Host          string
	Port          string
	Path          string
	TotalConns    int
	ConnRate      int
	MessageRate   int
	MessageSize   int
	RunDuration   time.Duration
	ConnTimeout   time.Duration
	SuccessLogDir string
}

// Statistics to track during the test
type Stats struct {
	ConnectionAttempts int64
	ConnectionSuccess  int64
	ConnectionErrors   int64
	MessagesSent       int64
	MessagesReceived   int64
	ResponseErrors     int64
	CurrentConnections int64
}

func main() {
	// Parse command line flags
	host := flag.String("host", "localhost", "WebSocket server host")
	port := flag.String("port", "5006", "WebSocket server port")
	path := flag.String("path", "/", "WebSocket endpoint path")
	totalConns := flag.Int("connections", 50000, "Total number of connections to establish")
	connRate := flag.Int("conn-rate", 100, "Connection establishment rate per second")
	messageRate := flag.Int("msg-rate", 10, "Messages per second per connection")
	messageSize := flag.Int("msg-size", 64, "Size of each message in bytes")
	runDuration := flag.Duration("duration", 5*time.Minute, "Total test duration")
	connTimeout := flag.Duration("timeout", 10*time.Second, "Connection timeout")
	logDir := flag.String("log-dir", "./logs", "Directory for success logs")
	verbose := flag.Bool("v", false, "Verbose logging")
	flag.Parse()

	// Create configuration
	config := &Config{
		Host:          *host,
		Port:          *port,
		Path:          *path,
		TotalConns:    *totalConns,
		ConnRate:      *connRate,
		MessageRate:   *messageRate,
		MessageSize:   *messageSize,
		RunDuration:   *runDuration,
		ConnTimeout:   *connTimeout,
		SuccessLogDir: *logDir,
	}

	// Ensure log directory exists
	if err := os.MkdirAll(config.SuccessLogDir, 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	// Create success log file
	logFile, err := os.Create(fmt.Sprintf("%s/load_test_%s.log", config.SuccessLogDir, time.Now().Format("20060102_150405")))
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}
	defer logFile.Close()
	logger := log.New(logFile, "", log.LstdFlags)

	// Initialize stats
	var stats Stats

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("Received shutdown signal, stopping test...")
		cancel()
	}()

	// Set test end timer
	go func() {
		select {
		case <-time.After(config.RunDuration):
			log.Printf("Test duration of %s reached, stopping...", config.RunDuration)
			cancel()
		case <-ctx.Done():
			return
		}
	}()

	// WebSocket URL
	u := url.URL{
		Scheme: "ws",
		Host:   fmt.Sprintf("%s:%s", config.Host, config.Port),
		Path:   config.Path,
	}
	serverAddr := u.String()
	log.Printf("Connecting to WebSocket server: %s", serverAddr)
	log.Printf("Test configuration: %d connections at %d conn/s, %d msg/s per conn, %d bytes per msg",
		config.TotalConns, config.ConnRate, config.MessageRate, config.MessageSize)

	// Generate test message
	testMessage := make([]byte, config.MessageSize)
	for i := range testMessage {
		testMessage[i] = 'A' + byte(i%26)
	}

	// WaitGroup to track active goroutines
	var wg sync.WaitGroup

	// Start connection rate limiter
	connTicker := time.NewTicker(time.Second / time.Duration(config.ConnRate))
	defer connTicker.Stop()

	// Launch statistics reporter
	go reportStats(ctx, &stats, *verbose)

	// Launch connections
	activeConns := make(map[int]*websocket.Conn)
	var connMapMutex sync.Mutex

	// Start creating connections
	for i := 0; i < config.TotalConns; i++ {
		select {
		case <-ctx.Done():
			log.Println("Test stopped before all connections were created")
			break
		case <-connTicker.C:
			wg.Add(1)
			connID := i
			go func() {
				defer wg.Done()
				atomic.AddInt64(&stats.ConnectionAttempts, 1)

				// Create connection with timeout
				connCtx, connCancel := context.WithTimeout(ctx, config.ConnTimeout)
				defer connCancel()

				// Connect to WebSocket server
				dialer := websocket.Dialer{
					HandshakeTimeout: config.ConnTimeout,
				}
				conn, _, err := dialer.DialContext(connCtx, serverAddr, nil)
				if err != nil {
					atomic.AddInt64(&stats.ConnectionErrors, 1)
					log.Printf("Connection error (ID %d): %v", connID, err)
					return
				}

				// Track successful connection
				atomic.AddInt64(&stats.ConnectionSuccess, 1)
				atomic.AddInt64(&stats.CurrentConnections, 1)
				defer func() {
					conn.Close()
					atomic.AddInt64(&stats.CurrentConnections, -1)
				}()

				// Add connection to map
				connMapMutex.Lock()
				activeConns[connID] = conn
				connMapMutex.Unlock()

				// Log successful connection
				logger.Printf("Connection %d established to %s", connID, serverAddr)

				// Create message rate limiter for this connection
				msgTicker := time.NewTicker(time.Second / time.Duration(config.MessageRate))
				defer msgTicker.Stop()

				// Set read handler
				go func() {
					for {
						_, message, err := conn.ReadMessage()
						if err != nil {
							if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
								log.Printf("Read error (ID %d): %v", connID, err)
								atomic.AddInt64(&stats.ResponseErrors, 1)
							}
							return
						}
						atomic.AddInt64(&stats.MessagesReceived, 1)
						if *verbose {
							log.Printf("Received from %d: %s", connID, message)
						}
						logger.Printf("Conn %d received: %s", connID, message)
					}
				}()

				// Send messages periodically
				for {
					select {
					case <-ctx.Done():
						return
					case <-msgTicker.C:
						err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Conn %d: %s", connID, testMessage)))
						if err != nil {
							log.Printf("Write error (ID %d): %v", connID, err)
							return
						}
						atomic.AddInt64(&stats.MessagesSent, 1)
						if *verbose {
							log.Printf("Sent message from connection %d", connID)
						}
					}
				}
			}()
		}
	}

	// Wait for test completion
	<-ctx.Done()
	log.Println("Closing all connections...")
	
	// Close all connections
	connMapMutex.Lock()
	for id, conn := range activeConns {
		if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
			log.Printf("Error closing connection %d: %v", id, err)
		}
		conn.Close()
	}
	connMapMutex.Unlock()
	
	// Wait for all goroutines to finish
	wg.Wait()
	
	// Final report
	log.Printf("Test completed. Final statistics:")
	log.Printf("- Connection attempts: %d", stats.ConnectionAttempts)
	log.Printf("- Successful connections: %d", stats.ConnectionSuccess)
	log.Printf("- Connection errors: %d", stats.ConnectionErrors)
	log.Printf("- Messages sent: %d", stats.MessagesSent)
	log.Printf("- Messages received: %d", stats.MessagesReceived)
	log.Printf("- Response errors: %d", stats.ResponseErrors)
	
	// Write summary to log file
	logger.Printf("TEST SUMMARY:")
	logger.Printf("- Connection attempts: %d", stats.ConnectionAttempts)
	logger.Printf("- Successful connections: %d", stats.ConnectionSuccess)
	logger.Printf("- Connection errors: %d", stats.ConnectionErrors)
	logger.Printf("- Messages sent: %d", stats.MessagesSent)
	logger.Printf("- Messages received: %d", stats.MessagesReceived)
	logger.Printf("- Response errors: %d", stats.ResponseErrors)
}

// reportStats periodically reports test statistics
func reportStats(ctx context.Context, stats *Stats, verbose bool) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			log.Printf("STATS - Connections: %d/%d (errors: %d), Messages: sent=%d received=%d (errors: %d)",
				stats.ConnectionSuccess, stats.ConnectionAttempts, stats.ConnectionErrors,
				stats.MessagesSent, stats.MessagesReceived, stats.ResponseErrors)
			
			if verbose {
				log.Printf("Current active connections: %d", stats.CurrentConnections)
			}
		}
	}
}