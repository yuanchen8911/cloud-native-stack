package k8s

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/NVIDIA/cloud-native-stack/pkg/measurement"
)

// Collect retrieves Kubernetes cluster version information from the API server.
// This provides cluster version details for comparison across environments.
func (k *Collector) collectServer(ctx context.Context) (map[string]measurement.Reading, error) {
	// Check if context is canceled
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Server Version
	serverVersion, err := k.ClientSet.Discovery().ServerVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get kubernetes version: %w", err)
	}

	slog.Debug("collected kubernetes version", slog.String("version", serverVersion.GitVersion))

	versionInfo := map[string]measurement.Reading{
		measurement.KeyVersion: measurement.Str(serverVersion.GitVersion),
		"platform":             measurement.Str(serverVersion.Platform),
		"goVersion":            measurement.Str(serverVersion.GoVersion),
	}

	return versionInfo, nil
}
