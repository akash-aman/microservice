package otel

/**
 * Reference: https://signoz.io/docs/instrumentation/opentelemetry-golang/
 */

import (
	"context"
	"fmt"
	"os"
	"pkg/logger"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"google.golang.org/grpc/credentials"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	collectorURL = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
)

type OtelConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Service  string `mapstructure:"service"`
	Insecure bool   `mapstructure:"insecure"`
}

type OtelCleanUp func(context.Context) error

func InitTracer(ctx context.Context, conf *OtelConfig, log logger.ILogger) OtelCleanUp {

	if collectorURL == "" {
		collectorURL = fmt.Sprintf("%s:%s", conf.Host, conf.Port)
	}

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
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", conf.Service),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		log.Fatalf("Could not set resources: %v", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)

	go func() {
		<-ctx.Done()
		err = exporter.Shutdown(context.Background())
		if err != nil {
			log.Errorf("Error exiting open-telemetry %v", err)
		} else {
			log.Info("open-telemetry exited gracefully")
		}
	}()

	return exporter.Shutdown
}
