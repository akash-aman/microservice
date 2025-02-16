package commands

import (
	"context"
	"pkg/logger"
	dtos_v1 "products/app/apis/products/get_by_id/v1/dtos"
	model_v1 "products/app/apis/products/get_by_id/v1/model"

	"github.com/labstack/echo/v4"
)

type GetProductByIdHandler struct {
	log  logger.Zapper
	echo *echo.Echo
	ctx  context.Context
}

func RegisterProductHandler(echo *echo.Echo, log logger.Zapper, ctx context.Context) *GetProductByIdHandler {
	return &GetProductByIdHandler{
		log:  log,
		echo: echo,
		ctx:  ctx,
	}
}

func (c *GetProductByIdHandler) Handle(ctx context.Context, cmd *model_v1.GetProductById) (*dtos_v1.ProductResponseDto, error) {

	data := &dtos_v1.ProductResponseDto{
		ID:          "0",
		Name:        "Sample Product",
		Description: "This is a sample product description.",
		Price:       0.0,
	}
	return data, nil
}
