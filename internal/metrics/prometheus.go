package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	TasksSubmitted = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "task_submitted_total",
			Help: "Total number of submitted tasks",
		},
		[]string{"priority"},
	)

	TasksProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "task_processed_total",
			Help: "Total number of processed tasks",
		},
		[]string{"status"},
	)

	TasksInQueue = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "task_queue_length",
			Help: "Number of tasks currently in the queue",
		},
	)

	TaskDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "task_processing_seconds",
			Help:    "Duration in seconds of task processing",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"priority"},
	)
)

// Init registers all custom metrics
func Init() {
	prometheus.MustRegister(
		TasksSubmitted,
		TasksProcessed,
		TasksInQueue,
		TaskDuration,
	)
}
