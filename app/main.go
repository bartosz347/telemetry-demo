package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"telemetry-demo/monitoring"
	"time"
)

var (
	instanceId     = uuid.New().String()
	serviceName    string
	targetServices []Service

	// Prometheus metrics
	promLatencyHistogram *prometheus.HistogramVec

	// OpenTelemetry metrics
	otelLatencyRecorder metric.Float64ValueRecorder
)

func main() {
	// Load app name
	var ok bool
	serviceName, ok = os.LookupEnv("APP_NAME")
	if !ok {
		log.Fatal("FATAL: APP_NAME not set!")
	}

	// Load services to call
	servicesStr, ok := os.LookupEnv("SERVICES_TO_CALL")
	if ok && servicesStr != "" {
		servicesStr = strings.ReplaceAll(servicesStr, " ", "")
		services := strings.Split(servicesStr, ",")

		for _, service := range services {
			s := strings.Split(service, ":")
			port, err := strconv.ParseUint(s[1], 10, 32)
			if err != nil {
				log.Fatalf("FATAL: Invalid port %s", s[1])
			}

			srv := Service{
				address: s[0],
				port:    uint(port),
			}
			targetServices = append(targetServices, srv)
		}
	}
	log.Printf("INFO: Loaded services: %s\n", targetServices)

	// Initialize Prometheus instrumentation
	log.Println("INFO: Initializing Prometheus")
	promLatencyHistogram = monitoring.InitPrometheus(serviceName)

	// Initialize OpenTelemetry instrumentation
	log.Println("INFO: Initializing OpenTelemetry")
	otelLatencyRecorder = monitoring.InitOpenTelemetry(serviceName, instanceId, os.Getenv("OTEL_AGENT"))

	// Start the HTTP server
	// Wrap `actionHandler` in OpenTelemetry middleware that provides tracing
	otelActionHandler := otelhttp.NewHandler(http.HandlerFunc(actionHandler), "/api/action")
	// Setup the endpoint
	http.Handle("/api/action", otelActionHandler)

	// Setup health check endpoint (without tracing)
	http.HandleFunc("/api/health", healthCheckHandler)

	log.Println("INFO: Starting HTTP server")
	log.Fatalf("ERROR: %s", http.ListenAndServe(":8080", nil))
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	fmt.Fprint(w, "OK")
}

func actionHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// call other services
	allRequestsSuccessful := true
	for _, service := range targetServices {
		if service.address != serviceName {
			err := service.call(r.Context(), r.URL.Query().Get("config"))
			if err != nil {
				allRequestsSuccessful = false
			}
		}
	}

	// do internal work
	trace.SpanFromContext(r.Context()).AddEvent("Starting internal processing")
	tr := otel.Tracer("tracer/internal")
	_, span := tr.Start(
		r.Context(),
		"internal-processing",
		trace.WithAttributes(semconv.PeerServiceKey.String("ExampleService")),
	)
	defer span.End()
	startInternal := time.Now()

	multiplier := GetMultiplier(r.URL.Query().Get("config"))
	runDummyLoop(multiplier)

	elapsedTimeInternal := time.Since(startInternal)
	promLatencyHistogram.WithLabelValues("internal-only", "OK").Observe(elapsedTimeInternal.Seconds())
	otelLatencyRecorder.Record(context.Background(), elapsedTimeInternal.Seconds(),
		attribute.String("instance", instanceId),
		attribute.Key("status").String("OK"),
		attribute.Key("type").String("internal-only"),
	)

	if allRequestsSuccessful {
		w.WriteHeader(200)
		fmt.Fprint(w, "OK")
	} else {
		w.WriteHeader(500)
		fmt.Fprint(w, "ERROR")
	}

	elapsedTime := time.Since(start)
	if allRequestsSuccessful {
		// Prometheus metrics
		promLatencyHistogram.WithLabelValues("total", "OK").Observe(elapsedTime.Seconds())
		// OpenTelemetry metrics
		otelLatencyRecorder.Record(context.Background(), elapsedTime.Seconds(),
			attribute.String("instance", instanceId),
			attribute.Key("status").String("OK"),
			attribute.Key("type").String("total"),
		)
	} else {
		promLatencyHistogram.WithLabelValues("total", "ERROR").Observe(elapsedTime.Seconds())
		otelLatencyRecorder.Record(context.Background(), elapsedTime.Seconds(),
			attribute.String("instance", instanceId),
			attribute.Key("status").String("ERROR"),
			attribute.Key("type").String("total"),
		)
	}
}
