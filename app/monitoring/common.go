package monitoring

import "github.com/prometheus/client_golang/prometheus"

var (
	bucketsConfig = prometheus.ExponentialBuckets(0.01, 1.8, 20)
)
