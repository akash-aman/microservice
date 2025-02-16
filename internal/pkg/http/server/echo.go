package server

import (
	"context"
	"pkg/logger"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	ReadTimeout    = 10 * time.Second
	WriteTimeout   = 10 * time.Second
	MaxHeaderBytes = 1 << 20
)

type EchoConfig struct {
	Port      string `mapstructure:"port" validate:"required"`
	Host      string `mapstructure:"host"`
	BaseRoute string `mapstructure:"baseRoute" validate:"required"`
	DebugMode bool   `mapstructure:"debugMode" validate:"required"`
}

func NewEchoServer() *echo.Echo {
	e := echo.New()
	return e
}

func RunEchoServer(ctx context.Context, echo *echo.Echo, log logger.Zapper, cfg *EchoConfig) error {

	/**
	 * Configure the echo server.
	 */
	echo.Server.ReadTimeout = ReadTimeout
	echo.Server.WriteTimeout = WriteTimeout
	echo.Server.MaxHeaderBytes = MaxHeaderBytes

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Infof("Shutting down HTTP PORT: {%s}", cfg.Port)
				err := echo.Shutdown(ctx)
				if err != nil {
					log.Errorf("Error shutting down HTTP server {%v}", err)
					return
				}
				log.Info("HTTP server shutdown gracefully")
				return
			}
		}
	}()

	err := echo.Start(cfg.Port)
	return err
}
