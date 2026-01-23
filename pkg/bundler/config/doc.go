// Package config provides configuration options for bundler implementations.
//
// This package defines the configuration structure and functional options pattern
// for customizing bundler behavior. All bundlers receive a Config instance that
// controls their output generation.
//
// # Configuration Options
//
//   - Deployer: Deployment method (DeployerHelm or DeployerArgoCD)
//   - IncludeReadme: Generate deployment documentation
//   - IncludeChecksums: Generate SHA256 checksums.txt file
//   - Version: Bundler version string
//   - ValueOverrides: Per-bundler value overrides from CLI --set flags
//   - Verbose: Enable verbose output
//
// # Deployer Types
//
// DeployerType constants define supported deployment methods:
//   - DeployerHelm: Generates Helm umbrella charts (default)
//   - DeployerArgoCD: Generates ArgoCD App of Apps manifests
//
// Use ParseDeployerType() to parse user input and GetDeployerTypes() for CLI help.
//
// # Usage
//
//	cfg := config.NewConfig(
//	    config.WithDeployer(config.DeployerHelm),
//	    config.WithIncludeChecksums(true),
//	)
//
// # Defaults
//
//   - Deployer: DeployerHelm
//   - IncludeReadme: true
//   - IncludeChecksums: true
//   - Version: "dev"
//
// Config is immutable after creation, safe for concurrent use.
package config
