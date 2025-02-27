package gobwas

import (
	"context"
	"net"
	"pkg/logger"
	"time"

	"github.com/mailru/easygo/netpoll"
)

func closeConnection(conn net.Conn) {
	conn.Close()
}

func isConnectionClosed(ev netpoll.Event) bool {
	return ev&(netpoll.EventReadHup|netpoll.EventHup) != 0
}

func handleClose(ctx context.Context, s *WebSocketServer, desc *netpoll.Desc, wsConn *Connection, conn net.Conn) {
	s.Poller.Stop(desc)
	s.Handler.OnClose(ctx, wsConn)
	closeConnection(conn)
}

/**
 * The deadliner struct is a wrapper around net.Conn that ensures every Read() and Write() operation has a deadline set before execution.
 * This can be particularly useful in networking applications for the following reasons:
 * 1. Avoiding Infinite Blocking
 * 		Network operations (Read() and Write()) can potentially block forever if the remote connection becomes unresponsive.
 *		By setting a deadline, if the operation doesn’t complete within the specified duration (d.t), it will return a timeout error (i/o timeout).
 */
type deadliner struct {
	net.Conn
	t time.Duration
}

func (d *deadliner) Write(p []byte) (int, error) {
	if err := d.Conn.SetWriteDeadline(time.Now().Add(d.t)); err != nil {
		return 0, err
	}
	return d.Conn.Write(p)
}

func (d *deadliner) Read(p []byte) (int, error) {
	if err := d.Conn.SetReadDeadline(time.Now().Add(d.t)); err != nil {
		return 0, err
	}
	return d.Conn.Read(p)
}

func shouldRetryAfterCooldown(err error) bool {
	if err == ErrScheduleTimeout {
		return true
	}

	if ne, ok := err.(net.Error); ok && ne.Temporary() {
		return true
	}

	return false
}

func cooldownAndRetry(ctx context.Context, err error, log logger.Zapper) {
	delay := 5 * time.Millisecond
	log.Infof(ctx, "accept error: %v; retrying in %s", err, delay)
	time.Sleep(delay)
}
