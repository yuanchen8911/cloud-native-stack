/*
Copyright © 2025 NVIDIA Corporation
SPDX-License-Identifier: Apache-2.0
*/

package validator

import (
	"fmt"
	"strings"

	"github.com/NVIDIA/cloud-native-stack/pkg/measurement"
	"github.com/NVIDIA/cloud-native-stack/pkg/snapshotter"
)

// ConstraintPath represents a parsed fully qualified constraint path.
// Format: {Type}.{Subtype}.{Key}
// Example: "K8s.server.version" -> Type="K8s", Subtype="server", Key="version"
type ConstraintPath struct {
	Type    measurement.Type
	Subtype string
	Key     string
}

// ParseConstraintPath parses a fully qualified constraint path.
// The path format is: {Type}.{Subtype}.{Key}
// The key portion may contain dots (e.g., "/proc/sys/kernel/osrelease").
func ParseConstraintPath(path string) (*ConstraintPath, error) {
	if path == "" {
		return nil, fmt.Errorf("constraint path cannot be empty")
	}

	// Split by dots, but we need at least 3 parts
	parts := strings.SplitN(path, ".", 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid constraint path %q: expected format {Type}.{Subtype}.{Key}", path)
	}

	// Parse and validate the measurement type
	measurementType, valid := measurement.ParseType(parts[0])
	if !valid {
		return nil, fmt.Errorf("invalid measurement type %q in constraint path %q: valid types are %v",
			parts[0], path, measurement.Types)
	}

	return &ConstraintPath{
		Type:    measurementType,
		Subtype: parts[1],
		Key:     parts[2], // Key is everything after the second dot (preserves dots in key)
	}, nil
}

// String returns the fully qualified path string.
func (cp *ConstraintPath) String() string {
	return fmt.Sprintf("%s.%s.%s", cp.Type, cp.Subtype, cp.Key)
}

// ExtractValue extracts the value at this path from a snapshot.
// Returns the value as a string, or an error if the path doesn't exist.
func (cp *ConstraintPath) ExtractValue(snap *snapshotter.Snapshot) (string, error) {
	if snap == nil {
		return "", fmt.Errorf("snapshot is nil")
	}

	// Find the measurement with matching type
	var targetMeasurement *measurement.Measurement
	for _, m := range snap.Measurements {
		if m.Type == cp.Type {
			targetMeasurement = m
			break
		}
	}

	if targetMeasurement == nil {
		return "", fmt.Errorf("measurement type %q not found in snapshot", cp.Type)
	}

	// Find the subtype
	var targetSubtype *measurement.Subtype
	for i := range targetMeasurement.Subtypes {
		if targetMeasurement.Subtypes[i].Name == cp.Subtype {
			targetSubtype = &targetMeasurement.Subtypes[i]
			break
		}
	}

	if targetSubtype == nil {
		return "", fmt.Errorf("subtype %q not found in measurement type %q", cp.Subtype, cp.Type)
	}

	// Find the key in data
	reading, exists := targetSubtype.Data[cp.Key]
	if !exists {
		return "", fmt.Errorf("key %q not found in subtype %q of measurement type %q", cp.Key, cp.Subtype, cp.Type)
	}

	// Convert reading to string
	return reading.String(), nil
}
