package server

import (
	"context"
	"fmt"
	"pkg/discovery"
	"pkg/helper"
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
				log.Infof(ctx, "Shutting down HTTP PORT: {%s}", cfg.Port)
				err := echo.Shutdown(ctx)
				if err != nil {
					log.Errorf(ctx, "Error shutting down HTTP server {%v}", err)
					return
				}
				log.Info(ctx, "HTTP server shutdown gracefully")
				return
			}
		}
	}()

	if err := discovery.RegisterServiceWithConsul(
		ctx, "echo-http-service",
		fmt.Sprintf("echo-http-service-%s", helper.GetMachineID()),
		fmt.Sprintf("http://%s", "host.docker.internal"),
		helper.GetPort(cfg.Port),
		discovery.HTTPService,
		log,
	); err != nil {
		log.Errorf(ctx, "Error registering with Consul: %v", err)
	}

	err := echo.Start(cfg.Port)
	return err
}
