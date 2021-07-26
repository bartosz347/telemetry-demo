package monitoring

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/unit"
	"google.golang.org/grpc"
	"log"
	"strings"
	"time"
)

func InitOpenTelemetry(serviceName string, instanceId string, otelAgentAddr string) metric.Float64ValueRecorder {
	ctx := context.Background()
	exp, err := otlp.NewExporter(ctx, otlpgrpc.NewDriver(
		otlpgrpc.WithInsecure(),
		otlpgrpc.WithEndpoint(otelAgentAddr),
		otlpgrpc.WithDialOption(grpc.WithBlock()), // block until connected, useful for testing
	))

	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}

	// For the demonstration, use sdktrace.AlwaysSample sampler to sample all traces.
	// In a production application, use sdktrace.ProbabilitySampler with a desired probability.
	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()), // All traces are sampled
		trace.WithBatcher(exp),
		trace.WithResource(resource.NewWithAttributes(semconv.ServiceNameKey.String(serviceName))),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// Metrics
	histogramBoundaries := histogram.WithExplicitBoundaries(
		bucketsConfig,
	)

	pusher := controller.New(
		processor.New(
			simple.NewWithHistogramDistribution(histogramBoundaries),
			exp,
		),
		controller.WithExporter(exp),
		controller.WithCollectPeriod(5*time.Second),
	)
	err = pusher.Start(ctx)
	if err != nil {
		log.Fatalf("failed to initialize metric controller: %v", err)
	}

	global.SetMeterProvider(pusher.MeterProvider())

	meter := global.GetMeterProvider().Meter("app")
	return metric.Must(meter).NewFloat64ValueRecorder(
		strings.ReplaceAll(serviceName, " ", "_")+"_operation_latency",
		metric.WithDescription(fmt.Sprintf("Processing time for %s (native OpenTelemetry metric).", serviceName)),
		metric.WithUnit(unit.Milliseconds),
	)
}
