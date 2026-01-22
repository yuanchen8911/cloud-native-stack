/*
Copyright © 2025 NVIDIA Corporation
SPDX-License-Identifier: Apache-2.0
*/

// Package umbrella generates Helm umbrella charts from recipe results.
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
//	└── README.md     # Deployment instructions
//
// Usage:
//
//	generator := umbrella.NewGenerator()
//	input := &umbrella.GeneratorInput{
//	    RecipeResult:    recipeResult,
//	    ComponentValues: map[string]map[string]interface{}{...},
//	    Version:         "v1.0.0",
//	}
//	output, err := generator.Generate(ctx, input, "./output")
package umbrella
