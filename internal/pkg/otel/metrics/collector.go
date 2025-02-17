package metrics

/**
 * Ref: https://github.com/SigNoz/sample-golang-app/blob/master/metrics/metrics.go
 *
 * Note: Modification needed
 */
import (
	"context"
	"math/rand"
	"pkg/logger"
	"time"

	"go.opentelemetry.io/otel/attribute"
	api "go.opentelemetry.io/otel/metric"
)

// Example use cases for sync counter
// - count the number of bytes received
// - count the number of requests completed
// - count the number of accounts created
// - count the number of checkpoints run
// - count the number of HTTP 5xx errors
//
// The increments should be non-negative.

func exceptionsCounter(ctx context.Context, meter api.Meter, log logger.Zapper) {
	counter, err := meter.Int64Counter("exceptions", api.WithUnit("1"),
		api.WithDescription("Counts exceptions since start"),
	)
	if err != nil {
		log.Fatalf(ctx, "Error in creating exceptions counter: ", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Info(ctx, "Terminating exception counter gracefully")
			return
		default:
			// Increment the counter by 1.
			// The attributes describe the exception.
			counter.Add(ctx, 1, api.WithAttributes(attribute.KeyValue{
				Key: attribute.Key("exception_type"), Value: attribute.StringValue("NullPointerException"),
			}))
			time.Sleep(time.Duration(rand.Int63n(5)) * time.Millisecond)
		}
	}
}

// Example use cases for async counter
//   - count the number of page faults
//   - CPU time, which could be reported for each thread, each process or the
//     entire system. For example "the CPU time for process
//     A running in user mode, measured in seconds".
//
// Basically, any value that is monotonically increasing and happens in the background.
// The increments should be non-negative.
func pageFaultsCounter(ctx context.Context, meter api.Meter, log logger.Zapper) {
	counter, err := meter.Int64ObservableCounter(
		"page_faults",
		api.WithUnit("1"),
		api.WithDescription("Counts page faults since start"),
	)
	if err != nil {
		log.Fatalf(ctx, "Error in creating page faults counter: ", err)
	}

	_, err = meter.RegisterCallback(
		func(_ context.Context, o api.Observer) error {
			attrSet := attribute.NewSet(attribute.String("process", "boo"))
			withAttrSet := api.WithAttributeSet(attrSet)
			o.ObserveInt64(counter, rand.Int63n(100), withAttrSet)
			return nil
		},
		counter,
	)
	if err != nil {
		log.Fatalf(ctx, "Error in registering callback: ", err)
	}
}

// Example use cases for Histogram
// - the request duration
// - the size of the response payload
func requestDurationHistogram(ctx context.Context, meter api.Meter, log logger.Zapper) {
	histogram, err := meter.Int64Histogram(
		"http_request_duration",
		api.WithUnit("ms"),
		api.WithDescription("The HTTP request duration in milliseconds"),
	)
	if err != nil {
		log.Fatalf(ctx, "Error in creating request duration histogram: ", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Info(ctx, "Terminating request duration counter gracefully")
			return
		default:
			histogram.Record(context.Background(), rand.Int63n(1000), api.WithAttributes(attribute.String("path", "/api/boo")))
			time.Sleep(time.Duration(rand.Int63n(5)) * time.Millisecond)
		}
	}

}

// Asynchronous Gauge is an Instrument
// which reports non-additive value(s)
// (e.g. the room temperature - it makes no sense to report the
// temperature value from multiple rooms and sum them up) when the
// instrument is being observed.

// Example use cases for Async Gauge
// - the current room temperature
// - the CPU fan speed

// Note: if the values are additive (e.g. the process heap size -
// it makes sense to report the heap size from multiple processes and sum them up,
// so we get the total heap usage),
// use Asynchronous Counter or Asynchronous UpDownCounter.

func roomTemperatureGauge(ctx context.Context, meter api.Meter, log logger.Zapper) {
	gauge, err := meter.Float64ObservableGauge(
		"room_temperature",
		api.WithUnit("1"),
		api.WithDescription("The room temperature in celsius"),
	)
	if err != nil {
		log.Fatalf(ctx, "Error in creating room temperature gauge: ", err)
	}

	_, err = meter.RegisterCallback(
		func(_ context.Context, o api.Observer) error {
			attrSet := attribute.NewSet(attribute.String("process", "boo"))
			withAttrSet := api.WithAttributeSet(attrSet)
			o.ObserveFloat64(gauge, rand.Float64()*100, withAttrSet)
			return nil
		},
		gauge,
	)
	if err != nil {
		log.Fatalf(ctx, "Error in registering callback: ", err)
	}
}

// UpDownCounter is an Instrument which supports increments and decrements.
// if the value is monotonically increasing, use Counter instead.
// Example use cases for UpDownCounter
// - the number of active requests
// - the number of items in a queue

func itemsInQueueUpDownCounter(ctx context.Context, meter api.Meter, log logger.Zapper) {
	counter, err := meter.Int64UpDownCounter(
		"items_in_queue",
		api.WithUnit("1"),
		api.WithDescription("The number of items in the queue"),
	)
	if err != nil {
		log.Fatalf(ctx, "Error in creating items in queue up down counter: ", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Info(ctx, "Terminating queue counter gracefully")
			return
		default:
			counter.Add(context.Background(), rand.Int63n(100), api.WithAttributes(attribute.String("queue", "A")))
			time.Sleep(time.Duration(rand.Int63n(5)) * time.Millisecond)
		}
	}
}

// Asynchronous UpDownCounter is an asynchronous Instrument
// which reports additive value(s)
// (e.g. the process heap size - it makes sense to report the heap size
// from multiple processes and sum them up, so we get the total heap usage)
// when the instrument is being observed.
//
// Example use cases for Asynchronous UpDownCounter
// - the process heap size
// - the approximate number of items in a lock-free circular buffer

func processHeapSizeUpDownCounter(ctx context.Context, meter api.Meter, log logger.Zapper) {
	counter, err := meter.Float64ObservableUpDownCounter(
		"process_heap_size",
		api.WithUnit("1"),
		api.WithDescription("The process heap size"),
	)
	if err != nil {
		log.Fatalf(ctx, "Error in creating process heap size up down counter: ", err)
	}

	_, err = meter.RegisterCallback(
		func(_ context.Context, o api.Observer) error {
			attrSet := attribute.NewSet(attribute.String("process", "boo"))
			withAttrSet := api.WithAttributeSet(attrSet)
			o.ObserveFloat64(counter, rand.Float64()*100, withAttrSet)
			return nil
		},
		counter,
	)
	if err != nil {
		log.Fatalf(ctx, "Error in registering callback: ", err)
	}
}

func GenerateMetrics(ctx context.Context, meter api.Meter, log logger.Zapper) {
	go exceptionsCounter(ctx, meter, log)
	go pageFaultsCounter(ctx, meter, log)
	go requestDurationHistogram(ctx, meter, log)
	go roomTemperatureGauge(ctx, meter, log)
	go itemsInQueueUpDownCounter(ctx, meter, log)
	go processHeapSizeUpDownCounter(ctx, meter, log)
}
