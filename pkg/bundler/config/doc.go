// Package config provides configuration options for bundler implementations.
//
// This package defines the configuration structure and functional options pattern
// for customizing bundler behavior. All bundlers receive a Config instance that
// controls their output generation.
//
// # Configuration Options
//
// Config controls bundler behavior through various settings:
//   - OutputFormat: Output format (yaml, json, helm)
//   - Compression: Enable gzip compression
//   - IncludeScripts: Generate installation/uninstallation scripts
//   - IncludeReadme: Generate deployment documentation
//   - IncludeChecksums: Generate SHA256 checksums.txt file
//   - Version: Bundler version string
//   - ValueOverrides: Per-bundler value overrides from CLI --set flags
//   - Verbose: Enable verbose output
//
// # Usage
//
// Create with defaults:
//
//	cfg := config.NewConfig()
//
// Customize with functional options:
//
//	cfg := config.NewConfig(
//	    config.WithOutputFormat("yaml"),
//	    config.WithIncludeScripts(true),
//	    config.WithIncludeChecksums(true),
//	    config.WithVersion("v1.0.0"),
//	)
//
// Access configuration:
//
//	if cfg.IncludeScripts() {
//	    // Generate scripts
//	}
//	version := cfg.Version()
//
// # Default Values
//
// The default configuration includes:
//   - OutputFormat: "yaml"
//   - IncludeScripts: true
//   - IncludeReadme: true
//   - IncludeChecksums: true
//   - Version: "dev"
//
// # Thread Safety
//
// Config is immutable after creation, making it safe for concurrent use by
// multiple bundlers executing in parallel.
//
// # Integration with Bundlers
//
// Bundlers receive Config through their constructor:
//
//	type MyBundler struct {
//	    cfg *config.Config
//	}
//
//	func NewMyBundler(cfg *config.Config) *MyBundler {
//	    return &MyBundler{cfg: cfg}
//	}
//
//	func (b *MyBundler) Make(ctx context.Context, r *recipe.Recipe, outputDir string) (*result.Result, error) {
//	    if b.cfg.IncludeScripts() {
//	        // Generate scripts
//	    }
//	    // ...
//	}
//
// Or use BaseBundler which embeds Config:
//
//	type MyBundler struct {
//	    *internal.BaseBundler
//	}
//
//	func (b *MyBundler) Make(ctx context.Context, r *recipe.Recipe, outputDir string) (*result.Result, error) {
//	    if b.Config.IncludeScripts() {
//	        // Generate scripts
//	    }
//	    // ...
//	}
package config
