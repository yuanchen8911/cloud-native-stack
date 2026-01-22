// Package internal provides shared utilities and helpers for bundler implementations.
//
// This package is internal to the bundler framework and should not be imported
// by external packages. It contains reusable components that reduce boilerplate
// and ensure consistency across bundler implementations.
//
// # Core Components
//
// BaseBundler: Helper that reduces boilerplate in bundler implementations by providing:
//   - File and directory management with proper error handling
//   - Template rendering with embedded template support
//   - Checksum generation with SHA256
//   - Context cancellation checking
//   - Result tracking with automatic metadata
//   - Prometheus metrics recording
//
// ContextChecker: Periodic context cancellation checking for long-running operations
//
// TestHarness: Shared test utilities for bundler testing with:
//   - Automatic recipe generation
//   - Temporary directory management
//   - Common assertions (files, checksums, metadata)
//   - Test isolation with cleanup
//
// # BaseBundler Usage
//
// Bundler implementations should embed BaseBundler to access helper methods:
//
//	type MyBundler struct {
//	    *internal.BaseBundler
//	}
//
//	func NewMyBundler(cfg *config.Config) *MyBundler {
//	    return &MyBundler{
//	        BaseBundler: internal.NewBaseBundler(cfg, types.BundleTypeMyOperator),
//	    }
//	}
//
// Then use helper methods in Make():
//
//	func (b *MyBundler) Make(ctx context.Context, r *recipe.Recipe, outputDir string) (*result.Result, error) {
//	    // Create directory structure
//	    dirs, err := b.CreateBundleDir(outputDir, "my-operator")
//	    if err != nil {
//	        return nil, err
//	    }
//
//	    // Render template
//	    content, err := b.RenderTemplate(myTemplate, "values.yaml", data)
//	    if err != nil {
//	        return nil, err
//	    }
//
//	    // Write file
//	    path := filepath.Join(dirs.Root, "values.yaml")
//	    if err := b.WriteFileString(path, content, 0644); err != nil {
//	        return nil, err
//	    }
//
//	    // Generate checksums (if enabled)
//	    if b.Config.IncludeChecksums() {
//	        if err := b.GenerateChecksums(ctx, dirs.Root); err != nil {
//	            return nil, err
//	        }
//	    }
//
//	    return b.Result, nil
//	}
//
// # Template Helpers
//
// RenderTemplate: Renders Go templates with automatic error handling
//
//	content, err := b.RenderTemplate(templateFS, "values.yaml", data)
//
// Supports go:embed templates:
//
//	//go:embed templates/*.tmpl
//	var templates embed.FS
//
// # File Helpers
//
// WriteFile and WriteFileString: Atomic file writes with directory creation:
//
//	err := b.WriteFile(path, []byte("content"), 0644)
//	err := b.WriteFileString(path, "content", 0644)
//
// Both methods automatically:
//   - Create parent directories
//   - Track files in result.Files
//   - Handle errors consistently
//
// # Checksum Generation
//
// GenerateChecksums: Creates checksums.txt with SHA256 hashes:
//
//	err := b.GenerateChecksums(ctx, bundleDir)
//
// Automatically includes all files in result.Files and respects context cancellation.
//
// # Context Checking
//
// CheckContext: Periodically check for context cancellation in loops:
//
//	for _, item := range items {
//	    if err := b.CheckContext(ctx); err != nil {
//	        return err
//	    }
//	    // Process item
//	}
//
// # Test Harness
//
// TestHarness simplifies bundler testing by providing common setup and assertions:
//
//	func TestMyBundler_Make(t *testing.T) {
//	    h := internal.NewTestHarness(t, "my-bundler")
//	    bundler := NewMyBundler(h.Config())
//
//	    h.TestMake(bundler)
//	}
//
// The harness automatically:
//   - Creates temporary directories
//   - Generates test recipes
//   - Validates output files
//   - Checks checksums
//   - Verifies metadata
//   - Cleans up resources
//
// # Helper Functions
//
// BuildBaseConfigMap: Creates baseline config map from recipe measurements
//
//	configMap, err := internal.BuildBaseConfigMap(recipe)
//
// GenerateFileFromTemplate: One-step template rendering and file writing
//
//	err := internal.GenerateFileFromTemplate(templateFS, "config.yaml", outputPath, data, 0644)
//
// ExtractK8sImageSubtype: Safely extracts K8s image subtype from recipe
//
//	subtype, err := internal.ExtractK8sImageSubtype(recipe)
package internal
