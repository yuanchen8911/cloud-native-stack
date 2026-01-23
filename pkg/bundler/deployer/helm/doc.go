/*
Copyright © 2025 NVIDIA Corporation
SPDX-License-Identifier: Apache-2.0
*/

// Package helm generates Helm umbrella charts from recipe results.
//
// An umbrella chart is a Helm chart that uses dependencies to deploy multiple
// sub-charts in a single release. This approach provides:
//
//   - Single deployment point with `helm install`
//   - Automatic dependency resolution with `helm dependency update`
//   - Shared configuration through values.yaml
//   - Consistent versioning across all components
//
// Output Structure:
//
//	output-dir/
//	├── Chart.yaml    # Chart metadata with dependencies
//	├── values.yaml   # Combined values for all components
//	├── README.md     # Deployment instructions
//	└── checksums.txt # SHA256 checksums (optional)
//
// Usage:
//
//	generator := helm.NewGenerator()
//	input := &helm.GeneratorInput{
//	    RecipeResult:    recipeResult,
//	    ComponentValues: componentValues,
//	    Version:         "1.0.0",
//	}
//	output, err := generator.Generate(ctx, input, "/path/to/output")
package helm
