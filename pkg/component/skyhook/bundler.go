package skyhook

import (
	"context"
	"log/slog"
	"path/filepath"

	"github.com/NVIDIA/cloud-native-stack/pkg/bundler/config"
	"github.com/NVIDIA/cloud-native-stack/pkg/bundler/result"
	"github.com/NVIDIA/cloud-native-stack/pkg/bundler/types"
	common "github.com/NVIDIA/cloud-native-stack/pkg/component/internal"
	"github.com/NVIDIA/cloud-native-stack/pkg/errors"
	"github.com/NVIDIA/cloud-native-stack/pkg/recipe"
)

const (
	Name = "skyhook-operator"
)

// componentConfig defines the Skyhook Operator bundler configuration.
var componentConfig = common.ComponentConfig{
	Name:              Name,
	DisplayName:       "skyhook",
	ValueOverrideKeys: []string{"skyhook"},
	AcceleratedNodeSelectorPaths: []string{
		"controllerManager.selectors",
	},
	AcceleratedTolerationPaths: []string{
		"controllerManager.tolerations",
	},
	DefaultHelmRepository: "https://helm.ngc.nvidia.com/nvidia",
	DefaultHelmChart:      "nvidia/skyhook-operator",
	TemplateGetter:        GetTemplate,
	CustomManifestFunc:    generateCustomManifests,
	MetadataFunc: func(configMap map[string]string) interface{} {
		return GenerateBundleMetadata(configMap)
	},
}

// Bundler creates Skyhook Operator application bundles based on recipes.
type Bundler struct {
	*common.BaseBundler
}

// NewBundler creates a new Skyhook bundler instance.
func NewBundler(conf *config.Config) *Bundler {
	return &Bundler{
		BaseBundler: common.NewBaseBundler(conf, types.BundleTypeSkyhook),
	}
}

// Make generates the Skyhook bundle based on the provided recipe.
func (b *Bundler) Make(ctx context.Context, input recipe.RecipeInput, dir string) (*result.Result, error) {
	return common.MakeBundle(ctx, b.BaseBundler, input, dir, componentConfig)
}

// generateCustomManifests generates Skyhook customization CR manifests.
func generateCustomManifests(ctx context.Context, b *common.BaseBundler, values map[string]interface{}, configMap map[string]string, dir string) ([]string, error) {
	// Check if customization is specified in values
	customizationName, ok := values["customization"].(string)
	if !ok || customizationName == "" {
		// No customization specified, nothing to generate
		return nil, nil
	}

	slog.Debug("generating Skyhook customization manifest",
		"customization", customizationName,
	)

	// Check if the customization template exists
	_, exists := GetCustomizationTemplate(customizationName)
	if !exists {
		availableCustomizations := ListCustomizations()
		return nil, errors.New(errors.ErrCodeInvalidRequest,
			"unknown Skyhook customization '"+customizationName+"'; available customizations: "+
				formatCustomizationList(availableCustomizations))
	}

	// Generate bundle metadata for manifest templates
	metadata := GenerateBundleMetadata(configMap)
	manifestData := map[string]interface{}{
		"Values": values,
		"Script": metadata,
	}

	// Add accelerated node tolerations if provided via CLI flags
	if tolerations := b.Config.AcceleratedNodeTolerations(); len(tolerations) > 0 {
		manifestData["Tolerations"] = common.TolerationsToPodSpec(tolerations)
	}

	// Add accelerated node selectors as matchExpressions if provided via CLI flags
	if nodeSelector := b.Config.AcceleratedNodeSelector(); len(nodeSelector) > 0 {
		manifestData["NodeSelectorExpressions"] = common.NodeSelectorToMatchExpressions(nodeSelector)
	}

	// Generate the customization manifest
	filePath := filepath.Join(dir, "manifests", customizationName+".yaml")
	if err := b.GenerateFileFromTemplate(ctx, GetCustomizationTemplate, customizationName,
		filePath, manifestData, 0644); err != nil {
		return nil, errors.Wrap(errors.ErrCodeInternal,
			"failed to generate customization manifest", err)
	}

	return []string{filePath}, nil
}

// formatCustomizationList formats a list of customization names for error messages.
func formatCustomizationList(names []string) string {
	if len(names) == 0 {
		return "(none available)"
	}
	result := ""
	for i, name := range names {
		if i > 0 {
			result += ", "
		}
		result += name
	}
	return result
}
