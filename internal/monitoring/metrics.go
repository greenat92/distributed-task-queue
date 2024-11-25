package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
)

// define metrics
var (
	TasksProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tasks_processed_total",
			Help: "Total number of tasks processed by the worker",
		},
		[]string{"status"},
	)

	TaskRetries = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "task_retires_total",
			Help: "Total number of task retries",
		},
	)

	TaskProcessingTime = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "task_processing_time_seconds",
			Help:    "Histogram of task processing times",
			Buckets: prometheus.DefBuckets,
		},
	)
)

// Register metrics
func InitMetrics() {
	prometheus.MustRegister(TasksProcessed)
	prometheus.MustRegister(TaskRetries)
	prometheus.MustRegister(TaskProcessingTime)
}
