package inits

import (
	"context"
	"pkg/logger"
	"pkg/websocket/gobwas"
)

func InitWebsocket(ctx context.Context, server *gobwas.WebSocketServer, log logger.Zapper) {
	if err := server.Start(ctx, log); err != nil {
		log.Errorf(ctx, "failed to start websocket server %v", err)
	}
}
