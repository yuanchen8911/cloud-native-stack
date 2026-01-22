/*
Copyright © 2025 NVIDIA Corporation
SPDX-License-Identifier: Apache-2.0
*/

// Package validator provides recipe constraint validation against system snapshots.
//
// # Overview
//
// The validator package evaluates recipe constraints against actual system measurements
// captured in snapshots. It supports version comparison operators and exact string matching
// to determine if a cluster meets the requirements specified in a recipe.
//
// # Constraint Format
//
// Constraints use fully qualified measurement paths in the format: {Type}.{Subtype}.{Key}
//
// Examples:
//
//	K8s.server.version         -> Kubernetes server version
//	OS.release.ID              -> Operating system identifier (e.g., "ubuntu")
//	OS.release.VERSION_ID      -> OS version (e.g., "24.04")
//	OS.sysctl./proc/sys/kernel/osrelease -> Kernel version
//
// # Supported Operators
//
// The following comparison operators are supported in constraint values:
//   - ">=" - Greater than or equal (version comparison)
//   - "<=" - Less than or equal (version comparison)
//   - ">"  - Greater than (version comparison)
//   - "<"  - Less than (version comparison)
//   - "==" - Exact match (string or version)
//   - "!=" - Not equal (string or version)
//   - (no operator) - Exact string match
//
// # Usage
//
// Basic validation:
//
//	v := validator.New()
//	result, err := v.Validate(ctx, recipe, snapshot)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Status: %s\n", result.Summary.Status)
//	for _, r := range result.Results {
//	    fmt.Printf("  %s: expected %q, got %q - %v\n",
//	        r.Name, r.Expected, r.Actual, r.Status)
//	}
//
// # Result Structure
//
// ValidationResult contains:
//   - Summary: Overall pass/fail counts and status
//   - Results: Per-constraint validation results with expected/actual values
//
// # Error Handling
//
// Constraints that cannot be evaluated (e.g., path not found in snapshot) are
// marked as "skipped" with appropriate warning messages, allowing partial
// validation results to be returned.
package validator
