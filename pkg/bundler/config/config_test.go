package config

import "testing"

func TestConfigImmutability(t *testing.T) {
	cfg := NewConfig()

	// Verify getters return expected default values
	outputFormat := cfg.OutputFormat()
	if outputFormat != "yaml" {
		t.Errorf("OutputFormat() = %s, want yaml", outputFormat)
	}

	// Verify getters return expected default values
	if !cfg.IncludeScripts() {
		t.Error("IncludeScripts() = false, want true")
	}

	if !cfg.IncludeReadme() {
		t.Error("IncludeReadme() = false, want true")
	}

	if !cfg.IncludeChecksums() {
		t.Error("IncludeChecksums() = false, want true")
	}

	if cfg.Compression() {
		t.Error("Compression() = true, want false")
	}

	if cfg.Verbose() {
		t.Error("Verbose() = true, want false")
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name:    "valid default config",
			config:  NewConfig(),
			wantErr: false,
		},
		{
			name:    "valid yaml format",
			config:  NewConfig(WithOutputFormat("yaml")),
			wantErr: false,
		},
		{
			name:    "valid json format",
			config:  NewConfig(WithOutputFormat("json")),
			wantErr: false,
		},
		{
			name:    "valid helm format",
			config:  NewConfig(WithOutputFormat("helm")),
			wantErr: false,
		},
		{
			name: "invalid output format",
			config: &Config{
				outputFormat: "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewConfigWithOptions(t *testing.T) {
	cfg := NewConfig(
		WithOutputFormat("json"),
		WithCompression(true),
		WithIncludeScripts(false),
		WithIncludeReadme(false),
		WithIncludeChecksums(false),
		WithVerbose(true),
	)

	// Verify all options were applied
	if cfg.OutputFormat() != "json" {
		t.Errorf("OutputFormat() = %s, want json", cfg.OutputFormat())
	}
	if !cfg.Compression() {
		t.Error("Compression() = false, want true")
	}
	if cfg.IncludeScripts() {
		t.Error("IncludeScripts() = true, want false")
	}
	if cfg.IncludeReadme() {
		t.Error("IncludeReadme() = true, want false")
	}
	if cfg.IncludeChecksums() {
		t.Error("IncludeChecksums() = true, want false")
	}
	if !cfg.Verbose() {
		t.Error("Verbose() = false, want true")
	}
}

func TestAllGetters(t *testing.T) {
	cfg := NewConfig(
		WithOutputFormat("helm"),
		WithCompression(true),
		WithIncludeScripts(false),
		WithIncludeReadme(true),
		WithIncludeChecksums(false),
		WithVerbose(true),
	)

	tests := []struct {
		name     string
		got      interface{}
		want     interface{}
		getterFn string
	}{
		{"OutputFormat", cfg.OutputFormat(), "helm", "OutputFormat()"},
		{"Compression", cfg.Compression(), true, "Compression()"},
		{"IncludeScripts", cfg.IncludeScripts(), false, "IncludeScripts()"},
		{"IncludeReadme", cfg.IncludeReadme(), true, "IncludeReadme()"},
		{"IncludeChecksums", cfg.IncludeChecksums(), false, "IncludeChecksums()"},
		{"Verbose", cfg.Verbose(), true, "Verbose()"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("%s = %v, want %v", tt.getterFn, tt.got, tt.want)
			}
		})
	}
}

func TestVersionOption(t *testing.T) {
	// Test WithVersion sets the version
	cfg := NewConfig(WithVersion("v1.2.3"))
	if cfg.Version() != "v1.2.3" {
		t.Errorf("Version() = %s, want v1.2.3", cfg.Version())
	}

	// Test default version
	cfgDefault := NewConfig()
	if cfgDefault.Version() != "dev" {
		t.Errorf("default Version() = %s, want dev", cfgDefault.Version())
	}
}

func TestValueOverridesOption(t *testing.T) {
	overrides := map[string]map[string]string{
		"gpuoperator": {
			"gds.enabled":    "true",
			"driver.version": "570.86.16",
		},
		"networkoperator": {
			"rdma.enabled": "true",
		},
	}

	cfg := NewConfig(WithValueOverrides(overrides))

	// Verify overrides were set
	got := cfg.ValueOverrides()
	if got == nil {
		t.Fatal("ValueOverrides() returned nil")
	}

	// Verify gpuoperator overrides
	if got["gpuoperator"]["gds.enabled"] != "true" {
		t.Errorf("gpuoperator gds.enabled = %s, want true", got["gpuoperator"]["gds.enabled"])
	}
	if got["gpuoperator"]["driver.version"] != "570.86.16" {
		t.Errorf("gpuoperator driver.version = %s, want 570.86.16", got["gpuoperator"]["driver.version"])
	}

	// Verify networkoperator overrides
	if got["networkoperator"]["rdma.enabled"] != "true" {
		t.Errorf("networkoperator rdma.enabled = %s, want true", got["networkoperator"]["rdma.enabled"])
	}
}

func TestValueOverridesImmutability(t *testing.T) {
	overrides := map[string]map[string]string{
		"gpuoperator": {"key": "value"},
	}

	cfg := NewConfig(WithValueOverrides(overrides))

	// Get and modify returned map
	got := cfg.ValueOverrides()
	got["gpuoperator"]["key"] = "modified"
	got["gpuoperator"]["new"] = "added"

	// Verify original config unchanged
	fresh := cfg.ValueOverrides()
	if fresh["gpuoperator"]["key"] != "value" {
		t.Error("modifying returned map affected config - not immutable")
	}
	if _, exists := fresh["gpuoperator"]["new"]; exists {
		t.Error("adding key to returned map affected config - not immutable")
	}
}

func TestValueOverridesNil(t *testing.T) {
	// WithValueOverrides with nil should not panic
	cfg := NewConfig(WithValueOverrides(nil))

	// ValueOverrides on empty config should return nil
	got := cfg.ValueOverrides()
	if len(got) > 0 {
		t.Errorf("ValueOverrides() = %v, want nil or empty", got)
	}
}

func TestValidateErrorMessages(t *testing.T) {
	tests := []struct {
		name            string
		config          *Config
		wantErrContains string
	}{
		{
			name: "invalid format error message",
			config: &Config{
				outputFormat: "xml",
			},
			wantErrContains: "invalid output format: xml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if err == nil {
				t.Fatal("Validate() error = nil, want error")
			}
			if tt.wantErrContains != "" {
				errMsg := err.Error()
				if len(errMsg) < len(tt.wantErrContains) || errMsg[:len(tt.wantErrContains)] != tt.wantErrContains {
					t.Errorf("Validate() error = %q, want error containing %q", errMsg, tt.wantErrContains)
				}
			}
		})
	}
}
