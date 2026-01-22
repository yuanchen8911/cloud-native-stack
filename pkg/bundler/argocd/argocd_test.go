/*
Copyright © 2025 NVIDIA Corporation
SPDX-License-Identifier: Apache-2.0
*/

package argocd

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/NVIDIA/cloud-native-stack/pkg/recipe"
)

const (
	testCertManager     = "cert-manager"
	testGPUOperator     = "gpu-operator"
	testNetworkOperator = "network-operator"
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
			testCertManager: {
				"installCRDs": true,
			},
			testGPUOperator: {
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

	// Verify output has expected number of files:
	// - app-of-apps.yaml (1)
	// - README.md (1)
	// - cert-manager/application.yaml, cert-manager/values.yaml (2)
	// - gpu-operator/application.yaml, gpu-operator/values.yaml (2)
	// Total: 6 files
	if len(output.Files) != 6 {
		t.Errorf("expected 6 files, got %d: %v", len(output.Files), output.Files)
	}

	// Check root files exist
	rootFiles := []string{"app-of-apps.yaml", "README.md"}
	for _, f := range rootFiles {
		path := filepath.Join(outputDir, f)
		if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
			t.Errorf("expected file %s does not exist", f)
		}
	}

	// Check component directories and files exist
	components := []string{testCertManager, testGPUOperator}
	for _, comp := range components {
		compDir := filepath.Join(outputDir, comp)
		if _, statErr := os.Stat(compDir); os.IsNotExist(statErr) {
			t.Errorf("expected directory %s does not exist", comp)
		}

		for _, f := range []string{"application.yaml", "values.yaml"} {
			path := filepath.Join(compDir, f)
			if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
				t.Errorf("expected file %s/%s does not exist", comp, f)
			}
		}
	}

	// Verify app-of-apps.yaml content
	appOfAppsContent, err := os.ReadFile(filepath.Join(outputDir, "app-of-apps.yaml"))
	if err != nil {
		t.Fatalf("failed to read app-of-apps.yaml: %v", err)
	}
	if !strings.Contains(string(appOfAppsContent), "kind: Application") {
		t.Error("app-of-apps.yaml missing 'kind: Application'")
	}
	if !strings.Contains(string(appOfAppsContent), "cns-stack") {
		t.Error("app-of-apps.yaml missing 'cns-stack' name")
	}

	// Verify component application.yaml has sync-wave
	gpuAppContent, err := os.ReadFile(filepath.Join(outputDir, testGPUOperator, "application.yaml"))
	if err != nil {
		t.Fatalf("failed to read gpu-operator/application.yaml: %v", err)
	}
	if !strings.Contains(string(gpuAppContent), "argocd.argoproj.io/sync-wave") {
		t.Error("gpu-operator/application.yaml missing sync-wave annotation")
	}

	// Verify README content
	readmeContent, err := os.ReadFile(filepath.Join(outputDir, "README.md"))
	if err != nil {
		t.Fatalf("failed to read README.md: %v", err)
	}
	if !strings.Contains(string(readmeContent), "ArgoCD") {
		t.Error("README.md missing 'ArgoCD'")
	}
	if !strings.Contains(string(readmeContent), testCertManager) {
		t.Error("README.md missing cert-manager component")
	}
	if !strings.Contains(string(readmeContent), testGPUOperator) {
		t.Error("README.md missing gpu-operator component")
	}
	if !strings.Contains(string(readmeContent), "Recipe Version:") {
		t.Error("README.md missing Recipe Version")
	}
	if !strings.Contains(string(readmeContent), "Bundler Version:") {
		t.Error("README.md missing Bundler Version")
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

func TestGenerate_WithRepoURL(t *testing.T) {
	g := NewGenerator()
	ctx := context.Background()
	outputDir := t.TempDir()

	customRepoURL := "https://github.com/my-org/my-gitops-repo.git"
	input := &GeneratorInput{
		RecipeResult: createTestRecipeResult(),
		ComponentValues: map[string]map[string]interface{}{
			testCertManager: {"installCRDs": true},
			testGPUOperator: {"driver": map[string]interface{}{"enabled": true}},
		},
		Version: "v1.0.0",
		RepoURL: customRepoURL,
	}

	output, err := g.Generate(ctx, input, outputDir)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if len(output.Files) == 0 {
		t.Fatal("expected files to be generated")
	}

	// Verify app-of-apps.yaml contains the custom repo URL
	appOfAppsContent, err := os.ReadFile(filepath.Join(outputDir, "app-of-apps.yaml"))
	if err != nil {
		t.Fatalf("failed to read app-of-apps.yaml: %v", err)
	}

	if !strings.Contains(string(appOfAppsContent), customRepoURL) {
		t.Errorf("app-of-apps.yaml should contain custom repo URL %q", customRepoURL)
	}

	// Verify it does NOT contain the placeholder
	if strings.Contains(string(appOfAppsContent), "YOUR-ORG") {
		t.Error("app-of-apps.yaml should not contain placeholder 'YOUR-ORG' when custom repo URL is provided")
	}
}

func TestGenerate_WithoutRepoURL(t *testing.T) {
	g := NewGenerator()
	ctx := context.Background()
	outputDir := t.TempDir()

	input := &GeneratorInput{
		RecipeResult: createTestRecipeResult(),
		ComponentValues: map[string]map[string]interface{}{
			testCertManager: {"installCRDs": true},
			testGPUOperator: {"driver": map[string]interface{}{"enabled": true}},
		},
		Version: "v1.0.0",
		RepoURL: "", // Empty - should use placeholder
	}

	output, err := g.Generate(ctx, input, outputDir)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if len(output.Files) == 0 {
		t.Fatal("expected files to be generated")
	}

	// Verify app-of-apps.yaml contains the placeholder
	appOfAppsContent, err := os.ReadFile(filepath.Join(outputDir, "app-of-apps.yaml"))
	if err != nil {
		t.Fatalf("failed to read app-of-apps.yaml: %v", err)
	}

	if !strings.Contains(string(appOfAppsContent), "YOUR-ORG") {
		t.Error("app-of-apps.yaml should contain placeholder 'YOUR-ORG' when no custom repo URL is provided")
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
		{"v25.3.3", "25.3.3"},
		{"", ""},
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

func TestGetNamespace(t *testing.T) {
	tests := []struct {
		name     string
		compName string
		expected string
	}{
		{testGPUOperator, testGPUOperator, testGPUOperator},
		{testNetworkOperator, testNetworkOperator, "nvidia-network-operator"},
		{testCertManager, testCertManager, testCertManager},
		{"unknown", "unknown-component", defaultNamespace},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comp := recipe.ComponentRef{Name: tt.compName}
			result := getNamespace(comp)
			if result != tt.expected {
				t.Errorf("getNamespace(%q) = %q, want %q", tt.compName, result, tt.expected)
			}
		})
	}
}

func TestSortComponentsByDeploymentOrder(t *testing.T) {
	components := []recipe.ComponentRef{
		{Name: testGPUOperator, Version: "v25.3.3"},
		{Name: testCertManager, Version: "v1.17.2"},
		{Name: testNetworkOperator, Version: "v25.4.0"},
	}
	deploymentOrder := []string{testCertManager, testGPUOperator, testNetworkOperator}

	sorted := sortComponentsByDeploymentOrder(components, deploymentOrder)

	if sorted[0].Name != testCertManager {
		t.Errorf("expected first component to be cert-manager, got %s", sorted[0].Name)
	}
	if sorted[1].Name != testGPUOperator {
		t.Errorf("expected second component to be gpu-operator, got %s", sorted[1].Name)
	}
	if sorted[2].Name != testNetworkOperator {
		t.Errorf("expected third component to be network-operator, got %s", sorted[2].Name)
	}
}

func TestSortComponentsByDeploymentOrder_EmptyOrder(t *testing.T) {
	components := []recipe.ComponentRef{
		{Name: testGPUOperator},
		{Name: testCertManager},
	}

	sorted := sortComponentsByDeploymentOrder(components, nil)

	// Should return original order when no deployment order specified
	if len(sorted) != 2 {
		t.Errorf("expected 2 components, got %d", len(sorted))
	}
	if sorted[0].Name != testGPUOperator {
		t.Errorf("expected first component to be gpu-operator, got %s", sorted[0].Name)
	}
}

func TestSortComponentsByDeploymentOrder_PartialOrder(t *testing.T) {
	components := []recipe.ComponentRef{
		{Name: "unknown-component"},
		{Name: testGPUOperator},
		{Name: testCertManager},
	}
	// Only cert-manager and gpu-operator in order, unknown-component should be last
	deploymentOrder := []string{testCertManager, testGPUOperator}

	sorted := sortComponentsByDeploymentOrder(components, deploymentOrder)

	if sorted[0].Name != testCertManager {
		t.Errorf("expected first component to be cert-manager, got %s", sorted[0].Name)
	}
	if sorted[1].Name != testGPUOperator {
		t.Errorf("expected second component to be gpu-operator, got %s", sorted[1].Name)
	}
	// unknown-component should be sorted to the end
	if sorted[2].Name != "unknown-component" {
		t.Errorf("expected third component to be unknown-component, got %s", sorted[2].Name)
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
				Name:    testCertManager,
				Version: "v1.17.2",
				Source:  "https://charts.jetstack.io",
			},
			{
				Name:    testGPUOperator,
				Version: "v25.3.3",
				Source:  "https://helm.ngc.nvidia.com/nvidia",
			},
		},
		DeploymentOrder: []string{testCertManager, testGPUOperator},
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
