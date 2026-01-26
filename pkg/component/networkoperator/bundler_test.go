package networkoperator

import (
	"testing"

	"github.com/NVIDIA/cloud-native-stack/pkg/bundler/config"
	"github.com/NVIDIA/cloud-native-stack/pkg/component/internal"
)

func TestBundler(t *testing.T) {
	internal.RunStandardBundlerTests(t, internal.StandardBundlerTestConfig{
		ComponentName:     Name,
		NewBundler:        func(cfg *config.Config) internal.BundlerInterface { return NewBundler(cfg) },
		GetTemplate:       GetTemplate,
		ExpectedTemplates: []string{"README.md"},
		ExpectedFiles:     []string{"values.yaml", "README.md", "checksums.txt"},
	})
}
