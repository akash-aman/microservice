package otel

/**
 *  * https://github.com/uptrace/opentelemetry-go-extra/tree/main/otelzap
 *  * https://uptrace.dev/get/opentelemetry-go#exporting-logs
 */
import (
	"context"
	"fmt"
	"log"
	"os"
	"pkg/otel/conf"
	"pkg/otel/logs"
	"pkg/otel/metrics"
	"pkg/otel/tracer"

	"go.opentelemetry.io/otel/attribute"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	collectorURL = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
)

func InitOpentelemetry(ctx context.Context, conf *conf.OtelConfig) (*sdklog.LoggerProvider, *sdktrace.TracerProvider, *metric.MeterProvider) {
	if collectorURL == "" {
		collectorURL = fmt.Sprintf("%s%s", conf.Host, conf.Port)
	}

	resources, err := resource.New(ctx,
		//resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			attribute.String("service.name", conf.Service),
			attribute.String("service.version", "1.0.0"),
		))

	if err != nil {
		log.Fatalf("Could not set resources: %v", err)
		return nil, nil, nil
	}

	logProvider := logs.ConfigOpenTelementryLogs(ctx, conf, resources)
	traceProvider := tracer.ConfigOpenTelementryTracer(ctx, conf, resources)
	metricProvider := metrics.ConfigOpenTelementryMeter(ctx, conf, resources)

	return logProvider, traceProvider, metricProvider
}
