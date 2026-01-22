package snapshotter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Snapshot collection metrics
	snapshotCollectionDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "cns_snapshot_collection_duration_seconds",
			Help:    "Time taken to collect a complete node snapshot",
			Buckets: []float64{1, 5, 10, 30, 60, 120, 300},
		},
	)

	snapshotCollectionTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cns_snapshot_collection_total",
			Help: "Total number of snapshot collection attempts",
		},
		[]string{"status"}, // success or error
	)

	snapshotCollectorDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cns_snapshot_collector_duration_seconds",
			Help:    "Time taken by individual collectors",
			Buckets: []float64{0.1, 0.5, 1, 5, 10, 30},
		},
		[]string{"collector"}, // image, k8s, kmod, systemd, grub, sysctl, smi, metadata
	)

	snapshotMeasurementCount = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "cns_snapshot_measurements",
			Help: "Number of measurements in the last collected snapshot",
		},
	)
)
