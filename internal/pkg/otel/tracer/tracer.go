package tracer

/**
 * Reference: https://signoz.io/docs/instrumentation/opentelemetry-golang/
 */

import (
	"context"
	"fmt"
	"pkg/otel/conf"

	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

	"google.golang.org/grpc/credentials"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func ConfigOpenTelementryTracer(ctx context.Context, conf *conf.OtelConfig, res *resource.Resource) {

	collectorURL := fmt.Sprintf("%s%s", conf.Host, conf.Grpc)

	var secureOption otlptracegrpc.Option

	if !conf.Insecure {
		secureOption = otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	} else {
		secureOption = otlptracegrpc.WithInsecure()
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(collectorURL),
		),
	)

	if err != nil {
		log.Fatalf("Failed to create exporter: %v", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(res),
		),
	)

	go func() {
		<-ctx.Done()
		err = exporter.Shutdown(context.Background())
		if err != nil {
			log.Fatalf("Error exiting open-telemetry tracer %v", err)
		} else {
			log.Print("open-telemetry tracer exited gracefully")
		}
	}()
}
