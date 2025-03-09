package gobwas

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"pkg/logger"
	"time"

	"github.com/failsafe-go/failsafe-go"
	"github.com/failsafe-go/failsafe-go/retrypolicy"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/mailru/easygo/netpoll"
	"github.com/panjf2000/ants/v2"
)

/**
 * Architecture Overview:
 *
 * ┌────────────────────┐     ┌────────────────────┐     ┌────────────────────┐
 * │                    │     │                    │     │                    │
 * │   Server           │     │   Connection       │     │   Handler          │
 * │                    │     │                    │     │                    │
 * │  - Manages TCP     │     │  - Manages WS      │     │  - Business Logic  │
 * │    listener        │◄───►│    protocol        │◄───►│  - Message         │
 * │  - Accepts         │     │  - Reads/writes    │     │    processing      │
 * │    connections     │     │    frames          │     │  - State           │
 * │  - Event polling   │     │  - Maintains       │     │    management      │
 * │                    │     │    connection      │     │                    │
 * └────────────────────┘     └────────────────────┘     └────────────────────┘
 *
 * Events flow:
 * 1. Server accepts TCP connection
 * 2. Connection handles WebSocket protocol details
 * 3. Handler receives parsed messages and implements business logic
 * 4. Handler sends messages back through Connection
 *
 */

// WebSocketConfig holds configuration for the WebSocket server
type WebSocketConfig struct {
	Host       string        `mapstructure:"host" validate:"required"`
	Port       int           `mapstructure:"port" validate:"required"`
	Workers    int           `mapstructure:"workers" validate:"required"`
	QueueSize  int           `mapstructure:"queueSize" validate:"required"`
	IOTimeout  time.Duration `mapstructure:"ioTimeout" validate:"required"`
	DebugPprof string        `mapstructure:"debugPprof"`
	MaxMsgSize int           `mapstructure:"maxMsgSize"`
}

type WebSocketServer struct {
	Handler WebSocketHandler
	Poller  netpoll.Poller
	Pool    *ants.Pool
	Config  *WebSocketConfig
}

func NewWebSocketServer(conf *WebSocketConfig, handler WebSocketHandler) *WebSocketServer {
	poller, err := netpoll.New(nil)
	if err != nil {
		return nil
	}

	// Create an ants worker pool
	pool, err := ants.NewPool(conf.Workers,
		ants.WithPreAlloc(true),
		ants.WithMaxBlockingTasks(conf.QueueSize),
	)
	if err != nil {
		return nil
	}

	return &WebSocketServer{
		Handler: handler,
		Poller:  poller,
		Config:  conf,
		Pool:    pool,
	}
}

/*
 * Start initializes and starts the WebSocket server. It sets up the pprof server if debugging is enabled,
 * creates a TCP listener on the configured address, and sets up the connection acceptor. The function
 * also handles graceful shutdown of the listener when the provided context is done.
 *
 * Parameters:
 *   - ctx: The context to control the server's lifecycle.
 *   - log: The logger instance for logging server events.
 *
 * Returns:
 *   - error: An error if the server fails to start or encounters an issue.
 */
func (s *WebSocketServer) Start(ctx context.Context, log logger.Zapper) error {
	if s.Config.DebugPprof != "" {
		log.Infof(ctx, "starting pprof server on %s", s.Config.DebugPprof)
		go func() {
			log.Infof(ctx, "pprof server error: %v", http.ListenAndServe(s.Config.DebugPprof, nil))
		}()
	}

	addr := fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to create listener: %w", err)
	}

	log.Infof(ctx, "websocket is listening on %s", ln.Addr().String())

	acceptDesc := netpoll.Must(netpoll.HandleListener(
		ln, netpoll.EventRead|netpoll.EventOneShot,
	))

	accept := make(chan error, 1)
	s.setupConnAcceptor(ctx, ln, acceptDesc, accept, log)

	go func() {
		<-ctx.Done()
		err = ln.Close()
		if err != nil {
			log.Fatalf(ctx, "Error closing WebSocket listener: %v", err)
		} else {
			log.Info(ctx, "WebSocket listener closed gracefully")
		}
		// Release the pool resources
		s.Pool.Release()
	}()

	return nil
}

/*
 * setupConnAcceptor sets up the connection acceptor for the WebSocket server. It starts the poller to listen for
 * incoming connections and schedules them to be handled by the worker pool. If an error occurs during acceptance,
 * it uses a retry policy to attempt reconnection.
 *
 * Parameters:
 *   - ctx: The context to control the server's lifecycle.
 *   - ln: The TCP listener for accepting connections.
 *   - acceptDesc: The netpoll descriptor for the listener.
 *   - accept: A channel to communicate acceptance errors.
 *   - log: The logger instance for logging server events.
 */
func (s *WebSocketServer) setupConnAcceptor(ctx context.Context, ln net.Listener, acceptDesc *netpoll.Desc, accept chan error, log logger.Zapper) {
	// Create retry policy outside the poller function to avoid recreation on each event
	retryPolicy := retrypolicy.Builder[any]().
		WithMaxRetries(3).
		WithDelay(100*time.Millisecond).
		WithJitter(time.Duration(0.2*float64(time.Second))).
		WithBackoff(2, 200*time.Millisecond).
		OnRetry(func(e failsafe.ExecutionEvent[any]) {
			log.Warnf(ctx, "retry attempt %d/%d after error: %v", e.Attempts(), 3, e.LastError())
		}).
		Build()

	s.Poller.Start(acceptDesc, func(e netpoll.Event) {
		// Accept new connection with retry capability
		err := acceptConnection(ctx, s, ln, accept, retryPolicy, log)

		if err != nil {
			log.Errorf(ctx, "fatal accept error: %v", err)
			return
		}

		// Resume accepting new connections
		s.Poller.Resume(acceptDesc)
	})
}

/*
 * acceptConnection accepts a new WebSocket connection with retry capability. It uses a retry policy to handle
 * transient errors during the acceptance process. The accepted connection is then handled by the worker pool.
 *
 * Parameters:
 *   - ctx: The context to control the server's lifecycle.
 *   - s: The WebSocket server instance.
 *   - ln: The TCP listener for accepting connections.
 *   - accept: A channel to communicate acceptance errors.
 *   - retryPolicy: The retry policy to use for accepting connections.
 *   - log: The logger instance for logging server events.
 *
 * Returns:
 *   - error: An error if the connection acceptance fails after retries.
 */
func acceptConnection(ctx context.Context, s *WebSocketServer, ln net.Listener, accept chan error,
	retryPolicy retrypolicy.RetryPolicy[any], log logger.Zapper) error {

	_, err := failsafe.Get(

		func() (interface{}, error) {

			return nil, s.Pool.Submit(func() {

				conn, err := ln.Accept()

				if err != nil {
					accept <- err
					return
				}

				accept <- nil

				s.handleConnection(ctx, conn, log)
			})
		},
		retryPolicy,
	)

	if err == nil {
		// Wait for the result of connection acceptance
		err = <-accept
	}

	return err
}

/*
 * handleConnection handles an incoming WebSocket connection. It upgrades the connection to a WebSocket,
 * logs the connection establishment, and invokes the handler's OnConnect method. If the connection is
 * accepted, it starts the connection poller to listen for incoming messages.
 *
 * Parameters:
 *   - ctx: The context to control the server's lifecycle.
 *   - conn: The net.Conn representing the incoming connection.
 *   - log: The logger instance for logging connection events.
 */
func (s *WebSocketServer) handleConnection(ctx context.Context, conn net.Conn, log logger.Zapper) {
	safeConn := &deadliner{conn, s.Config.IOTimeout}

	hs, err := ws.Upgrade(safeConn)
	if err != nil {
		log.Errorf(ctx, "%s: upgrade error: %v", conn.RemoteAddr().String(), err)
		closeConnection(conn)
		return
	}

	log.Infof(ctx, "%s: established websocket connection: %+v", conn.RemoteAddr().String(), hs)

	wsConn := NewConnection(s.Handler, safeConn)
	if err := s.Handler.OnConnect(ctx, wsConn); err != nil {
		log.Errorf(ctx, "handler rejected connection: %v", err)
		closeConnection(conn)
		return
	}

	desc := netpoll.Must(netpoll.HandleRead(conn))
	s.startConnectionPoller(ctx, desc, wsConn, conn, log)
}

/*
 * startConnectionPoller starts the poller for the WebSocket connection. It listens for events on the connection
 * and schedules message reading tasks to be handled by the worker pool. If the connection is closed, it handles
 * the cleanup process.
 *
 * Parameters:
 *   - ctx: The context to control the server's lifecycle.
 *   - desc: The netpoll descriptor for the connection.
 *   - wsConn: The WebSocket connection wrapper.
 *   - conn: The net.Conn representing the connection.
 *   - log: The logger instance for logging connection events.
 */
func (s *WebSocketServer) startConnectionPoller(ctx context.Context, desc *netpoll.Desc, wsConn *Connection, conn net.Conn, log logger.Zapper) {
	s.Poller.Start(desc, func(ev netpoll.Event) {
		if isConnectionClosed(ev) {
			log.Infof(ctx, "connection closed: %s", conn.RemoteAddr().String())
			handleClose(ctx, s, desc, wsConn, conn)
			return
		}

		// Use ants pool for message handling
		err := s.Pool.Submit(func() {
			if err := s.readMessage(ctx, wsConn); err != nil {
				log.Errorf(ctx, "error reading message: %v", err)
				handleClose(ctx, s, desc, wsConn, conn)
			}
		})

		if err != nil {
			log.Errorf(ctx, "failed to schedule message reading: %v", err)
			handleClose(ctx, s, desc, wsConn, conn)
		}
	})
}

/*
 * readMessage reads a message from the WebSocket connection. It handles different WebSocket opcodes
 * such as close, ping, and pong, and delegates message handling to the WebSocketHandler.
 *
 * Parameters:
 *   - ctx: The context to control the server's lifecycle.
 *   - conn: The WebSocket connection wrapper.
 *
 * Returns:
 *   - error: An error if reading the message fails or if the connection is closed by the client.
 */
func (s *WebSocketServer) readMessage(ctx context.Context, conn *Connection) error {
	msg, opCode, err := wsutil.ReadClientData(conn.Conn)

	if err != nil {
		return fmt.Errorf("read message error: %w", err)
	}

	switch opCode {
	case ws.OpClose:
		return fmt.Errorf("connection closed by client")
	case ws.OpPing:
		return conn.WritePong(msg)
	case ws.OpPong:
		return nil
	}

	return s.Handler.OnMessage(ctx, conn, opCode, msg)
}
