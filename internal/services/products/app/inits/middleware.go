package inits

import (
	"pkg/otel"
	"products/app/core/constants"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	otelMiddleware "go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

func ConfigMiddlewares(e *echo.Echo, otew *otel.OtelConfig) {

	e.HideBanner = false

	e.Use(middleware.Logger())

	e.Use(otelMiddleware.Middleware(otew.Service))

	e.Use(middleware.RequestID())

	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: constants.GzipLevel,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))

	e.Use(middleware.BodyLimit(constants.BodyLimit))
}
