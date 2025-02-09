package commands

import (
	"context"
	"pkg/logger"

	"github.com/labstack/echo/v4"
)

type GetProductByIdHandler struct {
	log  logger.ILogger
	echo *echo.Echo
	ctx  context.Context
}

func RegisterProductHandler(e *echo.Echo, log logger.ILogger, ctx context.Context) *GetProductByIdHandler {
	return &GetProductByIdHandler{
		log:  log,
		echo: e,
		ctx:  ctx,
	}
}

func (c *GetProductByIdHandler) Handle(ctx echo.Context) error {
	// Implement the handler logic here
	return nil
}
