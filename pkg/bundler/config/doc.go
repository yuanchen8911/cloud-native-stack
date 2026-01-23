// Package config provides configuration options for bundler implementations.
//
// This package defines the configuration structure and functional options pattern
// for customizing bundler behavior. All bundlers receive a Config instance that
// controls their output generation.
//
// # Configuration Options
//
// Config controls bundler behavior through various settings:
//   - Deployer: Deployment method (DeployerHelm or DeployerArgoCD)
//   - IncludeReadme: Generate deployment documentation
//   - IncludeChecksums: Generate SHA256 checksums.txt file
//   - Version: Bundler version string
//   - ValueOverrides: Per-bundler value overrides from CLI --set flags
//   - Verbose: Enable verbose output
//
// # Deployer Types
//
// The DeployerType constants define the supported deployment methods:
//   - DeployerHelm: Generates Helm umbrella charts (default)
//   - DeployerArgoCD: Generates ArgoCD App of Apps manifests
//
// Use ParseDeployerType() to parse user input and GetDeployerTypes() for CLI help.
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
//	    config.WithDeployer(config.DeployerHelm),
//	    config.WithIncludeChecksums(true),
//	    config.WithVersion("v1.0.0"),
//	)
//
// Access configuration:
//
//	if cfg.IncludeReadme() {
//	    // Generate README
//	}
//	version := cfg.Version()
//
// # Default Values
//
// The default configuration includes:
//   - Deployer: "helm"
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
//	    if b.cfg.IncludeReadme() {
//	        // Generate README
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
//	    if b.Config.IncludeChecksums() {
//	        // Generate checksums
//	    }
//	    // ...
//	}
package config
