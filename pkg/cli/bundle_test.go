/*
Copyright © 2025 NVIDIA Corporation
SPDX-License-Identifier: Apache-2.0
*/
package cli

import (
	"reflect"
	"testing"

	"github.com/NVIDIA/cloud-native-stack/pkg/bundler/config"
)

func TestParseSetFlags(t *testing.T) {
	tests := []struct {
		name     string
		setFlags []string
		want     map[string]map[string]string
		wantErr  bool
	}{
		{
			name:     "empty flags",
			setFlags: []string{},
			want:     map[string]map[string]string{},
			wantErr:  false,
		},
		{
			name:     "single flag",
			setFlags: []string{"gpuoperator:gds.enabled=true"},
			want: map[string]map[string]string{
				"gpuoperator": {
					"gds.enabled": "true",
				},
			},
			wantErr: false,
		},
		{
			name: "multiple flags same bundler",
			setFlags: []string{
				"gpuoperator:gds.enabled=true",
				"gpuoperator:driver.version=570.86.16",
			},
			want: map[string]map[string]string{
				"gpuoperator": {
					"gds.enabled":    "true",
					"driver.version": "570.86.16",
				},
			},
			wantErr: false,
		},
		{
			name: "multiple flags different bundlers",
			setFlags: []string{
				"gpuoperator:gds.enabled=true",
				"networkoperator:rdma.enabled=true",
			},
			want: map[string]map[string]string{
				"gpuoperator": {
					"gds.enabled": "true",
				},
				"networkoperator": {
					"rdma.enabled": "true",
				},
			},
			wantErr: false,
		},
		{
			name:     "value with equals sign",
			setFlags: []string{"gpuoperator:image.tag=v25.3.0=beta"},
			want: map[string]map[string]string{
				"gpuoperator": {
					"image.tag": "v25.3.0=beta",
				},
			},
			wantErr: false,
		},
		{
			name:     "value with spaces",
			setFlags: []string{"gpuoperator:custom.label=hello world"},
			want: map[string]map[string]string{
				"gpuoperator": {
					"custom.label": "hello world",
				},
			},
			wantErr: false,
		},
		{
			name:     "missing colon",
			setFlags: []string{"gpuoperatorgds.enabled=true"},
			wantErr:  true,
		},
		{
			name:     "missing equals sign",
			setFlags: []string{"gpuoperator:gds.enabledtrue"},
			wantErr:  true,
		},
		{
			name:     "empty path",
			setFlags: []string{"gpuoperator:=true"},
			wantErr:  true,
		},
		{
			name:     "empty value",
			setFlags: []string{"gpuoperator:gds.enabled="},
			wantErr:  true,
		},
		{
			name:     "only bundler name",
			setFlags: []string{"gpuoperator:"},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := config.ParseValueOverrides(tt.setFlags)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseValueOverrides() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseValueOverrides() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseOutputTarget(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantIsOCI bool
		wantReg   string
		wantRepo  string
		wantTag   string
		wantDir   string
		wantErr   bool
	}{
		{
			name:      "local directory relative",
			input:     "./bundle-out",
			wantIsOCI: false,
			wantDir:   "./bundle-out",
		},
		{
			name:      "local directory absolute",
			input:     "/tmp/bundles",
			wantIsOCI: false,
			wantDir:   "/tmp/bundles",
		},
		{
			name:      "local directory current",
			input:     ".",
			wantIsOCI: false,
			wantDir:   ".",
		},
		{
			name:      "OCI with tag",
			input:     "oci://ghcr.io/nvidia/bundle:v1.0.0",
			wantIsOCI: true,
			wantReg:   "ghcr.io",
			wantRepo:  "nvidia/bundle",
			wantTag:   "v1.0.0",
		},
		{
			name:      "OCI without tag defaults to latest",
			input:     "oci://ghcr.io/nvidia/bundle",
			wantIsOCI: true,
			wantReg:   "ghcr.io",
			wantRepo:  "nvidia/bundle",
			wantTag:   "latest",
		},
		{
			name:      "OCI with port and tag",
			input:     "oci://localhost:5000/test/bundle:v1",
			wantIsOCI: true,
			wantReg:   "localhost:5000",
			wantRepo:  "test/bundle",
			wantTag:   "v1",
		},
		{
			name:      "OCI with port no tag",
			input:     "oci://localhost:5000/test/bundle",
			wantIsOCI: true,
			wantReg:   "localhost:5000",
			wantRepo:  "test/bundle",
			wantTag:   "latest",
		},
		{
			name:      "OCI deeply nested repository",
			input:     "oci://ghcr.io/org/team/project/bundle:latest",
			wantIsOCI: true,
			wantReg:   "ghcr.io",
			wantRepo:  "org/team/project/bundle",
			wantTag:   "latest",
		},
		{
			name:    "OCI invalid reference",
			input:   "oci://",
			wantErr: true,
		},
		{
			name:    "OCI invalid characters",
			input:   "oci://ghcr.io/INVALID/Bundle:v1",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isOCI, reg, repo, tag, dir, err := parseOutputTarget(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("parseOutputTarget() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if isOCI != tt.wantIsOCI {
				t.Errorf("parseOutputTarget() isOCI = %v, want %v", isOCI, tt.wantIsOCI)
			}
			if reg != tt.wantReg {
				t.Errorf("parseOutputTarget() registry = %v, want %v", reg, tt.wantReg)
			}
			if repo != tt.wantRepo {
				t.Errorf("parseOutputTarget() repository = %v, want %v", repo, tt.wantRepo)
			}
			if tag != tt.wantTag {
				t.Errorf("parseOutputTarget() tag = %v, want %v", tag, tt.wantTag)
			}
			if dir != tt.wantDir {
				t.Errorf("parseOutputTarget() dir = %v, want %v", dir, tt.wantDir)
			}
		})
	}
}

func TestBundleCmd(t *testing.T) {
	cmd := bundleCmd()

	// Verify command configuration
	if cmd.Name != "bundle" {
		t.Errorf("expected command name 'bundle', got %q", cmd.Name)
	}

	// Verify required flags exist
	flagNames := make(map[string]bool)
	for _, flag := range cmd.Flags {
		names := flag.Names()
		for _, name := range names {
			flagNames[name] = true
		}
	}

	// Required flags for the new URI-based output approach
	requiredFlags := []string{"recipe", "r", "output", "o", "set", "plain-http", "insecure-tls"}
	for _, flag := range requiredFlags {
		if !flagNames[flag] {
			t.Errorf("expected flag %q to be defined", flag)
		}
	}

	// Verify node selector/toleration flags exist
	nodeFlags := []string{
		"system-node-selector",
		"system-node-toleration",
		"accelerated-node-selector",
		"accelerated-node-toleration",
	}
	for _, flag := range nodeFlags {
		if !flagNames[flag] {
			t.Errorf("expected flag %q to be defined", flag)
		}
	}

	// Verify removed flags don't exist (replaced by oci:// URI in --output)
	removedFlags := []string{"output-format", "registry", "repository", "tag", "push", "F"}
	for _, flag := range removedFlags {
		if flagNames[flag] {
			t.Errorf("flag %q should have been removed (use --output oci://... instead)", flag)
		}
	}
}
