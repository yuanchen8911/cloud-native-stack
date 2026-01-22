package os

import (
	"context"
	"fmt"
	"strings"

	"github.com/NVIDIA/cloud-native-stack/pkg/collector/file"
	"github.com/NVIDIA/cloud-native-stack/pkg/measurement"
)

var (
	filePathKMod = "/proc/modules"
)

// collectKMod retrieves the list of loaded kernel modules from /proc/modules
// and returns them as a subtype with module names as keys.
func (c *Collector) collectKMod(ctx context.Context) (*measurement.Subtype, error) {
	// Check if context is canceled
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	parser := file.NewParser()

	lines, err := parser.GetLines(filePathKMod)
	if err != nil {
		return nil, fmt.Errorf("failed to read kernel modules from %s: %w", filePathKMod, err)
	}

	readings := make(map[string]measurement.Reading)

	for _, line := range lines {
		// Module name is the first field (space-separated)
		fields := strings.Fields(line)
		if len(fields) > 0 {
			readings[fields[0]] = measurement.Bool(true)
		}
	}

	res := &measurement.Subtype{
		Name: "kmod",
		Data: readings,
	}

	return res, nil
}
