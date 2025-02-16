package inits

import (
	"context"
	"pkg/logger"
	registerGetProductById_v1 "products/app/apis/products/get_by_id/v1/endpoints"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

func ConfigEndpoints(validate *validator.Validate, log logger.Zapper, echo *echo.Echo, ctx context.Context) {
	//register_user_v1.MapRoute(validate, log, echo, ctx)
	registerGetProductById_v1.MapRoute(validate, log, echo, ctx)
	// Add graphql endpoint.

}
