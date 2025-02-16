package endpoints

import (
	"context"
	"net/http"

	"pkg/logger"
	ot "pkg/otel"

	dtos_v1 "products/app/apis/products/get_by_id/v1/dtos"
	model_v1 "products/app/apis/products/get_by_id/v1/model"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

func MapRoute(validator *validator.Validate, log logger.Zapper, echo *echo.Echo, ctx context.Context) {
	group := echo.Group("/api/v1/products")
	group.GET("", getProductById(validator, log, ctx))
}

func getProductById(validator *validator.Validate, log logger.Zapper, _ context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {

		/**
		 * 
		 */
		tracer, ctx := ot.NewTracer(c.Request().Context(), "GetProductById Controller")
		defer tracer.End()

		request := &dtos_v1.GetProductByIdRequestDto{}

		/**
		 * 
		 */
		if err := c.Bind(request); err != nil {
			errMsg := "[getProductById_handler.Bind] error in the binding request"
			log.Errorf("Error while binding request: %v", err)
			tracer.RecordError(err, "[getProductById_handler.Bind] error in the binding request")
			return echo.NewHTTPError(http.StatusBadRequest, errMsg)
		}

		tracer.AddAttributes(attribute.String("productId", request.ProductID))

		command := model_v1.NewGetProductById(request.ProductID)

		/**
		 * 
		 */
		if err := validator.StructCtx(ctx, command); err != nil {
			errMsg := "[getProductById_handler.StructCtx] command validation failed"
			log.Errorf("Validation error: %v", zap.Error(err))
			tracer.RecordError(err, "[getProductById_handler.StructCtx] command validation failed")
			return echo.NewHTTPError(http.StatusBadRequest, errMsg)
		}

		/**
		 * 
		 */
		result, err := mediatr.Send[*model_v1.GetProductById, *dtos_v1.GetProductByIdResponseDto](ctx, command)

		
		/**
		 * 
		 */
		if err != nil {
			log.Errorf("Error processing request: %v", err)
			tracer.RecordError(err, "Error processing request")
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		/**
		 * 
		 */
		log.Infof("Product retrieved successfully, id: %d", result.ID)

		return c.JSON(http.StatusOK, result)
	}
}
