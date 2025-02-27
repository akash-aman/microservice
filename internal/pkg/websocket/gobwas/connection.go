package gobwas

import (
	"io"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type Connection struct {
	Hub    WebSocketHandler
	Conn   io.ReadWriteCloser
	Mutex  *sync.RWMutex
	ConnID string
	Closed bool
}

func NewConnection(hub WebSocketHandler, conn io.ReadWriteCloser) *Connection {
	return &Connection{
		Hub:    hub,
		Conn:   conn,
		Closed: false,
	}
}

func (c *Connection) WritePong(msg []byte) error {
	return wsutil.WriteServerMessage(c.Conn, ws.OpPong, msg)
}
