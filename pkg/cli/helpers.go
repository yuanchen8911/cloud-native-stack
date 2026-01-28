/*
Copyright © 2025 NVIDIA Corporation
SPDX-License-Identifier: Apache-2.0
*/
package cli

import (
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/NVIDIA/cloud-native-stack/pkg/serializer"
)

// parseOutputFormat extracts and validates the output format from CLI flags.
// Returns the validated format or an error if the format is unknown.
func parseOutputFormat(cmd *cli.Command) (serializer.Format, error) {
	outFormat := serializer.Format(cmd.String("format"))
	if outFormat.IsUnknown() {
		return "", fmt.Errorf("unknown output format: %q, valid formats are: yaml, json, table", outFormat)
	}
	return outFormat, nil
}
