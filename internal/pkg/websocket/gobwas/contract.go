package gobwas

import (
	"context"

	"log"

	"github.com/gobwas/ws"
)

type WebSocketHandler interface {
	OnConnect(ctx context.Context, conn *Connection) error
	OnMessage(ctx context.Context, conn *Connection, msgType ws.OpCode, data []byte) error
	OnClose(ctx context.Context, conn *Connection)
}

type WebSocketHandlerImpl struct{}

func NewWebSocketHander() WebSocketHandler {
	return &WebSocketHandlerImpl{}
}

func (h *WebSocketHandlerImpl) OnConnect(ctx context.Context, conn *Connection) error {
	log.Println("Connected:", conn)
	return nil
}

func (h *WebSocketHandlerImpl) OnMessage(ctx context.Context, conn *Connection, msgType ws.OpCode, data []byte) error {
	log.Printf("Received message: %s", string(data))
	return nil
}

func (h *WebSocketHandlerImpl) OnClose(ctx context.Context, conn *Connection) {
	log.Println("Closed:", conn)
}
