// Package defaults provides centralized configuration constants for the CNS system.
//
// This package defines timeout values, retry parameters, and other configuration
// defaults used across the codebase. Centralizing these values ensures consistency
// and makes tuning easier.
//
// # Timeout Categories
//
// Timeouts are organized by component:
//
//   - Collector timeouts: For system data collection operations
//   - Handler timeouts: For HTTP request processing
//   - Server timeouts: For HTTP server configuration
//   - Kubernetes timeouts: For K8s API operations
//   - HTTP client timeouts: For outbound HTTP requests
//
// # Usage
//
// Import and use constants directly:
//
//	import "github.com/NVIDIA/cloud-native-stack/pkg/defaults"
//
//	ctx, cancel := context.WithTimeout(ctx, defaults.CollectorTimeout)
//	defer cancel()
//
// # Timeout Guidelines
//
// When choosing timeout values:
//
//   - Collectors: 10s default, respects parent context deadline
//   - HTTP handlers: 30s for recipes, 60s for bundles
//   - K8s operations: 30s for API calls, 5m for job completion
//   - Server shutdown: 30s for graceful shutdown
package defaults
