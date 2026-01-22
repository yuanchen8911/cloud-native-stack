/*
Copyright © 2025 NVIDIA Corporation
SPDX-License-Identifier: Apache-2.0
*/

package umbrella

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/NVIDIA/cloud-native-stack/pkg/recipe"
)

func TestNewGenerator(t *testing.T) {
	g := NewGenerator()
	if g == nil {
		t.Fatal("NewGenerator returned nil")
	}
}

func TestGenerate_Success(t *testing.T) {
	g := NewGenerator()
	ctx := context.Background()
	outputDir := t.TempDir()

	input := &GeneratorInput{
		RecipeResult: createTestRecipeResult(),
		ComponentValues: map[string]map[string]interface{}{
			"cert-manager": {
				"installCRDs": true,
			},
			"gpu-operator": {
				"driver": map[string]interface{}{
					"enabled": true,
				},
			},
		},
		Version: "v1.0.0",
	}

	output, err := g.Generate(ctx, input, outputDir)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Verify output
	if len(output.Files) != 3 {
		t.Errorf("expected 3 files, got %d", len(output.Files))
	}

	// Check files exist
	expectedFiles := []string{"Chart.yaml", "values.yaml", "README.md"}
	for _, f := range expectedFiles {
		path := filepath.Join(outputDir, f)
		if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
			t.Errorf("expected file %s does not exist", f)
		}
	}

	// Verify Chart.yaml content
	chartContent, err := os.ReadFile(filepath.Join(outputDir, "Chart.yaml"))
	if err != nil {
		t.Fatalf("failed to read Chart.yaml: %v", err)
	}
	if !strings.Contains(string(chartContent), "cert-manager") {
		t.Error("Chart.yaml missing cert-manager dependency")
	}
	if !strings.Contains(string(chartContent), "gpu-operator") {
		t.Error("Chart.yaml missing gpu-operator dependency")
	}

	// Verify values.yaml content
	valuesContent, err := os.ReadFile(filepath.Join(outputDir, "values.yaml"))
	if err != nil {
		t.Fatalf("failed to read values.yaml: %v", err)
	}
	if !strings.Contains(string(valuesContent), "cert-manager") {
		t.Error("values.yaml missing cert-manager values")
	}
	if !strings.Contains(string(valuesContent), "gpu-operator") {
		t.Error("values.yaml missing gpu-operator values")
	}
	if !strings.Contains(string(valuesContent), "enabled: true") {
		t.Error("values.yaml missing enabled flag")
	}
}

func TestGenerate_NilInput(t *testing.T) {
	g := NewGenerator()
	ctx := context.Background()

	_, err := g.Generate(ctx, nil, t.TempDir())
	if err == nil {
		t.Error("expected error for nil input")
	}
}

func TestGenerate_NilRecipeResult(t *testing.T) {
	g := NewGenerator()
	ctx := context.Background()

	input := &GeneratorInput{
		RecipeResult: nil,
	}

	_, err := g.Generate(ctx, input, t.TempDir())
	if err == nil {
		t.Error("expected error for nil recipe result")
	}
}

func TestGenerate_ContextCancellation(t *testing.T) {
	g := NewGenerator()
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	input := &GeneratorInput{
		RecipeResult:    createEmptyRecipeResult(),
		ComponentValues: map[string]map[string]interface{}{},
		Version:         "v1.0.0",
	}

	_, err := g.Generate(ctx, input, t.TempDir())
	if err == nil {
		t.Error("expected error for cancelled context")
	}
}

func TestNormalizeVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"v1.0.0", "1.0.0"},
		{"1.0.0", "1.0.0"},
		{"v0.1.0-alpha", "0.1.0-alpha"},
		{"", "0.1.0"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := normalizeVersion(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeVersion(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSortComponentsByDeploymentOrder(t *testing.T) {
	components := []string{"gpu-operator", "cert-manager", "network-operator"}
	deploymentOrder := []string{"cert-manager", "gpu-operator", "network-operator"}

	sorted := SortComponentsByDeploymentOrder(components, deploymentOrder)

	if sorted[0] != "cert-manager" {
		t.Errorf("expected first component to be cert-manager, got %s", sorted[0])
	}
	if sorted[1] != "gpu-operator" {
		t.Errorf("expected second component to be gpu-operator, got %s", sorted[1])
	}
	if sorted[2] != "network-operator" {
		t.Errorf("expected third component to be network-operator, got %s", sorted[2])
	}
}

// Helper functions

func createTestRecipeResult() *recipe.RecipeResult {
	return &recipe.RecipeResult{
		Kind:       "RecipeResult",
		APIVersion: "cns.nvidia.com/v1alpha1",
		Metadata: struct {
			GeneratedAt        time.Time                  `json:"generatedAt" yaml:"generatedAt"`
			Version            string                     `json:"version,omitempty" yaml:"version,omitempty"`
			AppliedOverlays    []string                   `json:"appliedOverlays,omitempty" yaml:"appliedOverlays,omitempty"`
			ExcludedOverlays   []string                   `json:"excludedOverlays,omitempty" yaml:"excludedOverlays,omitempty"`
			ConstraintWarnings []recipe.ConstraintWarning `json:"constraintWarnings,omitempty" yaml:"constraintWarnings,omitempty"`
		}{
			GeneratedAt: time.Now(),
			Version:     "v0.1.0",
		},
		Criteria: &recipe.Criteria{
			Service:     "eks",
			Accelerator: "h100",
			Intent:      "training",
		},
		ComponentRefs: []recipe.ComponentRef{
			{
				Name:    "cert-manager",
				Version: "v1.17.2",
				Source:  "https://charts.jetstack.io",
			},
			{
				Name:    "gpu-operator",
				Version: "v25.3.3",
				Source:  "https://helm.ngc.nvidia.com/nvidia",
			},
		},
		DeploymentOrder: []string{"cert-manager", "gpu-operator"},
	}
}

func createEmptyRecipeResult() *recipe.RecipeResult {
	return &recipe.RecipeResult{
		Kind:       "RecipeResult",
		APIVersion: "cns.nvidia.com/v1alpha1",
		Metadata: struct {
			GeneratedAt        time.Time                  `json:"generatedAt" yaml:"generatedAt"`
			Version            string                     `json:"version,omitempty" yaml:"version,omitempty"`
			AppliedOverlays    []string                   `json:"appliedOverlays,omitempty" yaml:"appliedOverlays,omitempty"`
			ExcludedOverlays   []string                   `json:"excludedOverlays,omitempty" yaml:"excludedOverlays,omitempty"`
			ConstraintWarnings []recipe.ConstraintWarning `json:"constraintWarnings,omitempty" yaml:"constraintWarnings,omitempty"`
		}{
			GeneratedAt: time.Now(),
			Version:     "v0.1.0",
		},
		ComponentRefs:   []recipe.ComponentRef{},
		DeploymentOrder: []string{},
	}
}
