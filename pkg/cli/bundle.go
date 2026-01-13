/*
Copyright © 2025 NVIDIA Corporation
SPDX-License-Identifier: Apache-2.0
*/

package cli

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/NVIDIA/cloud-native-stack/pkg/bundler"
	"github.com/NVIDIA/cloud-native-stack/pkg/bundler/config"
	"github.com/NVIDIA/cloud-native-stack/pkg/bundler/result"
	"github.com/NVIDIA/cloud-native-stack/pkg/oci"
	"github.com/NVIDIA/cloud-native-stack/pkg/recipe"
	"github.com/NVIDIA/cloud-native-stack/pkg/serializer"
	"github.com/NVIDIA/cloud-native-stack/pkg/snapshotter"
	"github.com/urfave/cli/v3"
)

// deployerArgoCD is the ArgoCD deployer type.
const deployerArgoCD = "argocd"

func bundleCmd() *cli.Command {
	return &cli.Command{
		Name:                  "bundle",
		EnableShellCompletion: true,
		Usage:                 "Generate deployment bundle from a given recipe",
		Description: `Generates a deployment bundle from a given recipe. By default, this creates
a Helm umbrella chart. Use --deployer argocd to generate ArgoCD Applications.

# Default Output (Helm Umbrella Chart)

  - Chart.yaml: Helm chart metadata with component dependencies
  - values.yaml: Combined values for all components
  - README.md: Deployment instructions
  - recipe.yaml: Copy of the input recipe for reference

# ArgoCD Output (--deployer argocd)

  - app-of-apps.yaml: Parent ArgoCD Application
  - <component>/application.yaml: ArgoCD Application per component
  - <component>/values.yaml: Values for each component
  - README.md: Deployment instructions

# Examples

Generate Helm umbrella chart (default):
  cnsctl bundle --recipe recipe.yaml --output ./my-bundle

Generate ArgoCD App of Apps:
  cnsctl bundle --recipe recipe.yaml --output ./my-bundle --deployer argocd

Override values in generated bundle:
  cnsctl bundle --recipe recipe.yaml --set gpuoperator:driver.version=570.133.20

Set node selectors for GPU workloads:
  cnsctl bundle --recipe recipe.yaml \
    --accelerated-node-selector nodeGroup=gpu-nodes \
    --accelerated-node-toleration nvidia.com/gpu=present:NoSchedule

Push bundle to OCI registry:
  cnsctl bundle --recipe recipe.yaml --output ./my-bundle \
    --push --registry ghcr.io --repository nvidia/cns-bundle --tag v1.0.0

# Deployment (Helm)

After generating the Helm bundle, deploy using:
  cd my-bundle
  helm dependency update
  helm install cns-stack . -n cns-system --create-namespace

# Deployment (ArgoCD)

After generating the ArgoCD bundle:
  1. Push the generated files to your GitOps repository
  2. Apply the app-of-apps.yaml to your ArgoCD cluster:
     kubectl apply -f app-of-apps.yaml`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "recipe",
				Aliases:  []string{"r"},
				Required: true,
				Usage: `Path/URI to previously generated recipe from which to build the bundle.
	Supports: file paths, HTTP/HTTPS URLs, or ConfigMap URIs (cm://namespace/name).`,
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Value:   ".",
				Usage:   "Output directory path for the generated Helm umbrella chart",
			},
			&cli.StringSliceFlag{
				Name:  "set",
				Usage: "Override values in generated bundle files (format: bundler:path.to.field=value, e.g., --set gpuoperator:gds.enabled=true)",
			},
			&cli.StringSliceFlag{
				Name:  "system-node-selector",
				Usage: "Node selector for system components (format: key=value, can be repeated)",
			},
			&cli.StringSliceFlag{
				Name:  "system-node-toleration",
				Usage: "Toleration for system components (format: key=value:effect, can be repeated)",
			},
			&cli.StringSliceFlag{
				Name:  "accelerated-node-selector",
				Usage: "Node selector for accelerated/GPU nodes (format: key=value, can be repeated)",
			},
			&cli.StringSliceFlag{
				Name:  "accelerated-node-toleration",
				Usage: "Toleration for accelerated/GPU nodes (format: key=value:effect, can be repeated)",
			},
			&cli.StringFlag{
				Name:    "deployer",
				Aliases: []string{"d"},
				Value:   "",
				Usage:   "Deployment method: '' (default Helm umbrella chart) or 'argocd' (App of Apps pattern)",
			},
			&cli.StringFlag{
				Name:  "repo",
				Value: "",
				Usage: "Git repository URL for ArgoCD applications (only used with --deployer argocd)",
			},
			// OCI push flags
			&cli.BoolFlag{
				Name:  "push",
				Usage: "Push generated bundle as OCI artifact to registry",
			},
			&cli.StringFlag{
				Name:  "registry",
				Usage: "OCI registry host (e.g., ghcr.io, localhost:5000)",
			},
			&cli.StringFlag{
				Name:  "repository",
				Usage: "OCI repository path (e.g., nvidia/cns-bundle)",
			},
			&cli.StringFlag{
				Name:  "tag",
				Usage: "OCI image tag (default: latest)",
			},
			&cli.BoolFlag{
				Name:  "insecure-tls",
				Usage: "Skip TLS certificate verification for OCI registry",
			},
			&cli.BoolFlag{
				Name:  "plain-http",
				Usage: "Use HTTP instead of HTTPS for OCI registry (for local development)",
			},
			kubeconfigFlag,
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			recipeFilePath := cmd.String("recipe")
			outputDir := cmd.String("output")
			kubeconfig := cmd.String("kubeconfig")
			deployer := cmd.String("deployer")
			repoURL := cmd.String("repo")

			// Validate deployer flag
			if deployer != "" && deployer != deployerArgoCD {
				return fmt.Errorf("invalid --deployer value: %q (must be '' or 'argocd')", deployer)
			}

			// OCI push options
			pushEnabled := cmd.Bool("push")
			registryHost := cmd.String("registry")
			repository := cmd.String("repository")
			tag := cmd.String("tag")
			insecureTLS := cmd.Bool("insecure-tls")
			plainHTTP := cmd.Bool("plain-http")

			// Validate push flags
			if pushEnabled {
				if registryHost == "" {
					return fmt.Errorf("--registry is required when --push is enabled")
				}
				if repository == "" {
					return fmt.Errorf("--repository is required when --push is enabled")
				}
				// Validate registry and repository format
				if err := oci.ValidateRegistryReference(registryHost, repository); err != nil {
					return fmt.Errorf("invalid OCI reference: %w", err)
				}
			}

			// Parse value overrides from --set flags
			valueOverrides, err := config.ParseValueOverrides(cmd.StringSlice("set"))
			if err != nil {
				return fmt.Errorf("invalid --set flag: %w", err)
			}

			// Parse node selectors
			systemNodeSelector, err := snapshotter.ParseNodeSelectors(cmd.StringSlice("system-node-selector"))
			if err != nil {
				return fmt.Errorf("invalid --system-node-selector: %w", err)
			}
			acceleratedNodeSelector, err := snapshotter.ParseNodeSelectors(cmd.StringSlice("accelerated-node-selector"))
			if err != nil {
				return fmt.Errorf("invalid --accelerated-node-selector: %w", err)
			}

			// Parse tolerations
			systemNodeTolerations, err := snapshotter.ParseTolerations(cmd.StringSlice("system-node-toleration"))
			if err != nil {
				return fmt.Errorf("invalid --system-node-toleration: %w", err)
			}
			acceleratedNodeTolerations, err := snapshotter.ParseTolerations(cmd.StringSlice("accelerated-node-toleration"))
			if err != nil {
				return fmt.Errorf("invalid --accelerated-node-toleration: %w", err)
			}

			outputType := "umbrella chart"
			if deployer == deployerArgoCD {
				outputType = "ArgoCD applications"
			}
			slog.Info("generating bundle",
				slog.String("type", outputType),
				slog.String("recipe", recipeFilePath),
				slog.String("output", outputDir),
			)

			// Load recipe from file/URL/ConfigMap
			rec, err := serializer.FromFileWithKubeconfig[recipe.RecipeResult](recipeFilePath, kubeconfig)
			if err != nil {
				slog.Error("failed to load recipe file", "error", err, "path", recipeFilePath)
				return err
			}

			// Create bundler with config
			cfg := config.NewConfig(
				config.WithVersion(version),
				config.WithDeployer(deployer),
				config.WithRepoURL(repoURL),
				config.WithValueOverrides(valueOverrides),
				config.WithSystemNodeSelector(systemNodeSelector),
				config.WithSystemNodeTolerations(systemNodeTolerations),
				config.WithAcceleratedNodeSelector(acceleratedNodeSelector),
				config.WithAcceleratedNodeTolerations(acceleratedNodeTolerations),
			)

			b, err := bundler.NewWithConfig(cfg)
			if err != nil {
				slog.Error("failed to create bundler", "error", err)
				return err
			}

			// Generate bundle
			out, err := b.Make(ctx, rec, outputDir)
			if err != nil {
				slog.Error("bundle generation failed", "error", err)
				return err
			}

			slog.Info("bundle generated",
				"type", outputType,
				"files", out.TotalFiles,
				"size_bytes", out.TotalSize,
				"duration_sec", out.TotalDuration.Seconds(),
				"output_dir", out.OutputDir,
			)

			// Print deployment instructions
			printBundleDeploymentInstructions(deployer, repoURL, out)

			// Push to OCI registry if enabled
			if pushEnabled {
				absOutputDir, err := filepath.Abs(outputDir)
				if err != nil {
					return fmt.Errorf("failed to resolve output directory: %w", err)
				}

				// Default tag if not provided
				imageTag := tag
				if imageTag == "" {
					imageTag = "latest"
				}

				slog.Info("pushing bundle to OCI registry",
					"registry", registryHost,
					"repository", repository,
					"tag", imageTag,
				)

				// Package locally first
				packageResult, err := oci.Package(ctx, oci.PackageOptions{
					SourceDir:  absOutputDir,
					OutputDir:  absOutputDir,
					Registry:   registryHost,
					Repository: repository,
					Tag:        imageTag,
				})
				if err != nil {
					return fmt.Errorf("failed to package OCI artifact: %w", err)
				}

				slog.Info("OCI artifact packaged locally",
					"reference", packageResult.Reference,
					"digest", packageResult.Digest,
					"store_path", packageResult.StorePath,
				)

				// Push to remote registry
				pushResult, err := oci.PushFromStore(ctx, packageResult.StorePath, oci.PushOptions{
					Registry:    registryHost,
					Repository:  repository,
					Tag:         imageTag,
					PlainHTTP:   plainHTTP,
					InsecureTLS: insecureTLS,
				})
				if err != nil {
					return fmt.Errorf("failed to push OCI artifact to registry: %w", err)
				}

				slog.Info("OCI artifact pushed successfully",
					"reference", pushResult.Reference,
					"digest", pushResult.Digest,
				)
			}

			return nil
		},
	}
}

// printBundleDeploymentInstructions prints user-friendly deployment instructions.
func printBundleDeploymentInstructions(deployer, repoURL string, out *result.Output) {
	if deployer == deployerArgoCD {
		fmt.Printf("\nArgoCD applications generated successfully!\n")
		fmt.Printf("Output directory: %s\n", out.OutputDir)
		fmt.Printf("Files generated: %d\n", out.TotalFiles)
		fmt.Printf("\nTo deploy:\n")
		fmt.Printf("  1. Push the generated files to your GitOps repository\n")
		if repoURL == "" {
			fmt.Printf("  2. Update app-of-apps.yaml with your repository URL\n")
			fmt.Printf("  3. Apply to your cluster:\n")
		} else {
			fmt.Printf("  2. Apply to your cluster:\n")
		}
		fmt.Printf("     kubectl apply -f %s/app-of-apps.yaml\n", out.OutputDir)
	} else {
		fmt.Printf("\nUmbrella chart generated successfully!\n")
		fmt.Printf("Output directory: %s\n", out.OutputDir)
		fmt.Printf("Files generated: %d\n", out.TotalFiles)
		fmt.Printf("\nTo deploy:\n")
		fmt.Printf("  cd %s\n", out.OutputDir)
		fmt.Printf("  helm dependency update\n")
		fmt.Printf("  helm install cns-stack . -n cns-system --create-namespace\n")
	}
}
