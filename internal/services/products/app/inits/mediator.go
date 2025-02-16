package inits

import (
	"context"
	"pkg/logger"

	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"

	getProductById_dto_v1 "products/app/apis/products/get_by_id/v1/dtos"
	getProductById_handler_v1 "products/app/apis/products/get_by_id/v1/handler"
	getProductById_model_v1 "products/app/apis/products/get_by_id/v1/model"
)

func InitMediator(log logger.Zapper, echo *echo.Echo, ctx context.Context) error {

	err := mediatr.RegisterRequestHandler[*getProductById_model_v1.GetProductById, *getProductById_dto_v1.ProductResponseDto](getProductById_handler_v1.RegisterProductHandler(echo, log, ctx))
	if err != nil {
		return err
	}

	return nil
}
