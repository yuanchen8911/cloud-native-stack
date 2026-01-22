package collector

import (
	"context"

	"github.com/NVIDIA/cloud-native-stack/pkg/measurement"
)

// Collector defines the interface for collecting system measurement data.
// Implementations gather data from various sources including system services,
// hardware components, OS configuration, and cluster state.
// All collectors must support context-based cancellation.
type Collector interface {
	Collect(ctx context.Context) (*measurement.Measurement, error)
}
