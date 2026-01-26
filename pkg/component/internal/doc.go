// Package internal provides the generic bundler framework and shared utilities for component implementations.
//
// This package is internal to the bundler framework and should not be imported
// by external packages. It contains reusable components that reduce boilerplate
// and ensure consistency across bundler implementations.
//
// # Generic Bundler Framework
//
// The framework provides a declarative approach to bundle generation using:
//
// ComponentConfig: Defines all component-specific settings in one struct:
//   - Name and DisplayName for identification
//   - ValueOverrideKeys for CLI --set flag mapping
//   - Node selector and toleration paths for workload placement
//   - DefaultHelmRepository, DefaultHelmChart, DefaultHelmChartVersion for Helm deployment
//   - TemplateGetter function for embedded templates
//   - Optional CustomManifestFunc for generating additional manifests
//   - Optional MetadataExtensions map for custom README template data (preferred over MetadataFunc)
//
// MakeBundle: Generic function that handles all common bundling steps:
//   - Extracting component values from recipe input
//   - Applying user value overrides from CLI flags
//   - Applying node selectors and tolerations to Helm paths
//   - Creating directory structure
//   - Writing values.yaml with proper YAML headers
//   - Calling optional CustomManifestFunc for additional files
//   - Generating README from templates
//   - Computing checksums
//
// # Minimal Bundler Example
//
// Most bundlers can be implemented in ~50 lines using the framework:
//
//	var componentConfig = internal.ComponentConfig{
//	    Name:                  "my-operator",
//	    DisplayName:           "My Operator",
//	    ValueOverrideKeys:     []string{"myoperator"},
//	    DefaultHelmRepository: "https://charts.example.com",
//	    DefaultHelmChart:      "example/my-operator",
//	    TemplateGetter:        GetTemplate,
//	}
//
//	type Bundler struct {
//	    *internal.BaseBundler
//	}
//
//	func NewBundler(cfg *config.Config) *Bundler {
//	    return &Bundler{
//	        BaseBundler: internal.NewBaseBundler(cfg, types.BundleTypeMyOperator),
//	    }
//	}
//
//	func (b *Bundler) Make(ctx context.Context, input recipe.RecipeInput, dir string) (*result.Result, error) {
//	    return internal.MakeBundle(ctx, b.BaseBundler, input, dir, componentConfig)
//	}
//
// # Custom Metadata
//
// Components that need additional template data beyond the default BundleMetadata
// can provide a MetadataExtensions map in ComponentConfig:
//
//	var componentConfig = internal.ComponentConfig{
//	    // ... other fields ...
//	    MetadataExtensions: map[string]interface{}{
//	        "InstallCRDs":   true,
//	        "CustomField":   "custom-value",
//	    },
//	}
//
// These extensions are merged into the BundleMetadata.Extensions map and can be
// accessed in templates via {{ .Script.Extensions.InstallCRDs }}.
//
// For more complex metadata requirements, MetadataFunc is still supported but
// MetadataExtensions is preferred for simple key-value additions.
//
// # Custom Manifest Generation
//
// Components that need to generate additional manifests can provide a CustomManifestFunc:
//
//	var componentConfig = internal.ComponentConfig{
//	    // ... other fields ...
//	    CustomManifestFunc: func(ctx context.Context, b *internal.BaseBundler,
//	        values map[string]interface{}, configMap map[string]string, dir string) ([]string, error) {
//	        // Generate manifests using b.WriteFile() or b.GenerateFileFromTemplate()
//	        return []string{"manifests/custom.yaml"}, nil
//	    },
//	}
//
// # BaseBundler Helper Methods
//
// BaseBundler provides common functionality for file operations:
//
//   - CreateBundleDir: Creates directory structure with proper permissions
//   - WriteFile: Writes content with automatic directory creation
//   - WriteFileString: Convenience wrapper for string content
//   - RenderTemplate: Renders Go templates with error handling
//   - GenerateFileFromTemplate: One-step template rendering and file writing
//   - GenerateChecksums: Creates checksums.txt with SHA256 hashes
//   - CheckContext: Periodic context cancellation checking
//   - Finalize: Records timing and result metadata
//   - BuildConfigMapFromInput: Creates baseline config map from recipe input
//
// # Helper Functions
//
// Utility functions for common operations:
//
//   - GetConfigValue: Safely extracts config map values with defaults
//   - GetBundlerVersion: Returns bundler version from config
//   - GetRecipeBundlerVersion: Returns recipe version from config
//   - MarshalYAMLWithHeader: Serializes values with component header
//   - ApplyMapOverrides: Applies dot-notation overrides to nested maps
//   - ApplyNodeSelectorOverrides: Applies node selectors to Helm paths
//   - ApplyTolerationsOverrides: Applies tolerations to Helm paths
//   - GenerateDefaultBundleMetadata: Creates default BundleMetadata struct
//
// # Default BundleMetadata
//
// Components using the default metadata get:
//
//   - Namespace, HelmRepository, HelmChart, HelmChartVersion
//   - Version (bundler version), RecipeVersion
//
// Access in templates via {{ .Script.Namespace }}, {{ .Script.Version }}, etc.
//
// # TestHarness
//
// TestHarness simplifies bundler testing by providing common setup and assertions:
//
//	func TestMyBundler_Make(t *testing.T) {
//	    h := internal.NewTestHarness(t, "my-bundler")
//	    bundler := NewMyBundler(h.Config())
//	    h.TestMake(bundler)
//	}
//
// The harness automatically creates temporary directories, generates test recipes,
// validates output files, and cleans up resources.
package internal
