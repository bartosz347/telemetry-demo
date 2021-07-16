package monitoring

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"strings"
)

func InitPrometheus(serviceName string) *prometheus.HistogramVec {
	processingDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    strings.ReplaceAll(serviceName, " ", "_") + "_operation_latency",
		Help:    fmt.Sprintf("Processing time for %s (native Prometheus metric).", serviceName),
		Buckets: bucketsConfig,
	}, []string{"type", "status"})

	prometheus.MustRegister(processingDuration)
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(":9000", nil)
		if err != nil {
			log.Printf("WARNING: Prometheus metric initialization failed: %s", err)
		}
	}()

	return processingDuration
}
