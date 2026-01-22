package snapshotter

import (
	"context"

	"github.com/NVIDIA/cloud-native-stack/pkg/header"
	"github.com/NVIDIA/cloud-native-stack/pkg/measurement"
)

// Snapshotter defines the interface for collecting system configuration snapshots.
// Implementations gather measurements from various system components and serialize
// the results for analysis or recommendation generation.
type Snapshotter interface {
	Measure(ctx context.Context) error
}

// NewSnapshot creates a new Snapshot instance with an initialized Measurements slice.
func NewSnapshot() *Snapshot {
	return &Snapshot{
		Measurements: make([]*measurement.Measurement, 0),
	}
}

// Snapshot represents a collected configuration snapshot from a system node.
// It contains metadata and measurements from various collectors including
// Kubernetes, GPU, OS configuration, and systemd services.
type Snapshot struct {
	header.Header `json:",inline" yaml:",inline"`

	// Measurements contains the collected measurements from various collectors.
	Measurements []*measurement.Measurement `json:"measurements" yaml:"measurements"`
}
