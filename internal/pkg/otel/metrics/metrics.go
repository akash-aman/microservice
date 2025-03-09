package metrics

import (
	"context"
	"fmt"
	"log"
	"pkg/otel/conf"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc/credentials"
)

func ConfigOpenTelementryMeter(ctx context.Context, conf *conf.OtelConfig, res *resource.Resource) *metricsdk.MeterProvider {

	collectorURL := fmt.Sprintf("%s:%d", conf.Host, conf.Grpc)

	var secureOption otlpmetricgrpc.Option

	if !conf.Insecure {
		secureOption = otlpmetricgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	} else {
		secureOption = otlpmetricgrpc.WithInsecure()
	}

	exporter, err := otlpmetricgrpc.New(
		context.Background(),
		secureOption,
		otlpmetricgrpc.WithEndpoint(collectorURL),
	)

	if err != nil {
		log.Fatalf("Failed to create exporter: %v", err)
	}

	provider := metricsdk.NewMeterProvider(
		metricsdk.WithResource(res),
		metricsdk.WithReader(metricsdk.NewPeriodicReader(exporter,
			metricsdk.WithInterval(10*time.Second))),
	)

	go func() {
		<-ctx.Done()
		err = exporter.Shutdown(context.Background())
		if err != nil {
			log.Fatalf("Error exiting open-telemetry meter %v", err)
		} else {
			log.Print("open-telemetry meter exited gracefully")
		}
	}()

	return provider
}
