package logs

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"pkg/otel/conf"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
)

func ConfigOpenTelementryLogs(ctx context.Context, conf *conf.OtelConfig, res *resource.Resource) *sdklog.LoggerProvider {
	var secureOption otlploghttp.Option

	collectorURL := fmt.Sprintf("%s:%d", conf.Host, conf.Http)

	if !conf.Insecure {
		secureOption = otlploghttp.WithTLSClientConfig(&tls.Config{})
	} else {
		secureOption = otlploghttp.WithInsecure()
	}

	exp, err := otlploghttp.New(
		ctx,
		secureOption,
		otlploghttp.WithEndpoint(collectorURL),
		otlploghttp.WithCompression(otlploghttp.GzipCompression),
	)

	if err != nil {
		log.Fatalf("Could not set resources: %v", err)
	}

	bsp := sdklog.NewBatchProcessor(exp,
		sdklog.WithMaxQueueSize(10_000),
		sdklog.WithExportMaxBatchSize(10_000),
		sdklog.WithExportInterval(10*time.Second),
		sdklog.WithExportTimeout(10*time.Second),
	)

	provider := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(bsp),
	)

	global.SetLoggerProvider(provider)

	go func() {
		<-ctx.Done()
		err = provider.Shutdown(context.Background())
		if err != nil {
			log.Fatalf("Error exiting open-telemetry logs %v", err)
		} else {
			log.Printf("open-telemetry logs exited gracefully")
		}
	}()

	return provider
}
