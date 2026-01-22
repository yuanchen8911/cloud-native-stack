// Package types defines the type system for bundler implementations.
//
// This package provides a type-safe way to identify and work with different
// bundler types throughout the framework. It ensures compile-time safety when
// working with bundler identifiers.
//
// # Core Type
//
// BundleType: String-based type identifier for bundlers
//
//	type BundleType string
//
// # Predefined Types
//
// The package defines constants for supported bundlers:
//
//	const (
//	    BundleTypeGpuOperator     BundleType = "gpu-operator"
//	    BundleTypeNetworkOperator BundleType = "network-operator"
//	)
//
// # Usage
//
// Use constants for type-safe bundler references:
//
//	bundlerType := types.BundleTypeGpuOperator
//	fmt.Println(bundlerType.String()) // Output: gpu-operator
//
// Parse string input with validation:
//
//	bundlerType, err := types.ParseType("gpu-operator")
//	if err != nil {
//	    // Handle invalid type
//	}
//
// Get list of all supported types:
//
//	allTypes := types.SupportedTypes()
//	for _, t := range allTypes {
//	    fmt.Println(t)
//	}
//
// Convert to string slice:
//
//	typeNames := types.SupportedBundleTypesAsStrings()
//	// ["gpu-operator", "network-operator"]
//
// # Type Validation
//
// ParseType validates input and returns descriptive errors:
//
//	_, err := types.ParseType("invalid-type")
//	// Error: invalid bundler type: "invalid-type", supported types: gpu-operator, network-operator
//
// # Type Comparison
//
// Types can be compared directly:
//
//	if bundlerType == types.BundleTypeGpuOperator {
//	    // Handle GPU Operator
//	}
//
// # Map Keys
//
// BundleType can be used as map keys:
//
//	bundlers := map[types.BundleType]Bundler{
//	    types.BundleTypeGpuOperator: gpuBundler,
//	    types.BundleTypeNetworkOperator: networkBundler,
//	}
//
// # Adding New Types
//
// To add a new bundler type:
//
// 1. Define the constant:
//
//	const BundleTypeMyOperator BundleType = "my-operator"
//
// 2. Add to SupportedTypes():
//
//	func SupportedTypes() []BundleType {
//	    return []BundleType{
//	        BundleTypeGpuOperator,
//	        BundleTypeNetworkOperator,
//	        BundleTypeMyOperator,  // Add here
//	    }
//	}
//
// 3. Register the bundler implementation in its package init():
//
//	func init() {
//	    registry.MustRegister(types.BundleTypeMyOperator, NewMyBundler)
//	}
//
// # Zero Value
//
// The zero value of BundleType is an empty string "", which is not a valid
// bundler type. Always use the predefined constants or ParseType.
package types
