package os

import (
	"context"
	"fmt"

	"github.com/NVIDIA/cloud-native-stack/pkg/collector/file"
	"github.com/NVIDIA/cloud-native-stack/pkg/measurement"
)

var (
	filePathGrub    = "/proc/cmdline"
	fileLineDelGrub = " "
	fileKVDelGrub   = "="

	// Keys to filter out from GRUB config for privacy/security
	filterOutGrubKeys = []string{
		"root",
	}
)

// collectGRUB retrieves the GRUB bootloader parameters from /proc/cmdline
// and returns them as a subtype with key-value pairs for each boot parameter.
func (c *Collector) collectGRUB(ctx context.Context) (*measurement.Subtype, error) {
	// Check if context is canceled
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	parser := file.NewParser(
		file.WithDelimiter(fileLineDelGrub),
		file.WithKVDelimiter(fileKVDelGrub),
	)

	params, err := parser.GetMap(filePathGrub)
	if err != nil {
		return nil, fmt.Errorf("failed to read GRUB params from %s: %w", filePathGrub, err)
	}

	props := make(map[string]measurement.Reading, 0)

	for k, v := range params {
		props[k] = measurement.Str(v)
	}

	res := &measurement.Subtype{
		Name: "grub",
		Data: measurement.FilterOut(props, filterOutGrubKeys),
	}

	return res, nil
}
