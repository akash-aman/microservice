package endpoints

import (
	"context"
	"net/http"

	"pkg/logger"

	dtos_v1 "products/app/apis/products/get_by_id/v1/dtos"
	model_v1 "products/app/apis/products/get_by_id/v1/model"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/pkg/errors"
)

func MapRoute(validator *validator.Validate, log logger.ILogger, echo *echo.Echo, ctx context.Context) {
	group := echo.Group("/api/v1/products")
	group.GET("", getProductById(validator, log, ctx))
}

func getProductById(validator *validator.Validate, log logger.ILogger, ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := &dtos_v1.ProductRequestDto{}

		/**
		 * Data Mapping.
		 */
		if err := c.Bind(request); err != nil {
			badRequestErr := errors.Wrap(err, "[registerBulkOrderTask_handler.Bind] error in the binding request")
			log.Error("Error while binding request : %v", badRequestErr)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		command := model_v1.NewGetProductById(request.ProductID)

		/**
		 * Data Validation.
		 */
		if err := validator.StructCtx(ctx, command); err != nil {
			validationErr := errors.Wrap(err, "[registerBulkOrderTask_handler.StructCtx] command validation failed")
			log.Error("Validation error : %v", validationErr)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		/**
		 * Call Business Logic Handler.
		 */

		result, err := mediatr.Send[*model_v1.GetProductById, *dtos_v1.ProductResponseDto](ctx, command)

		/**
		 * Response Validation & Mapping.
		 */
		if err != nil {
			log.Error("Error enqueuing the bulk order task : %v", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		log.Infof("Bulk order task enqueued, id: %d", result.ID)
		return c.JSON(http.StatusCreated, result)
	}
}
