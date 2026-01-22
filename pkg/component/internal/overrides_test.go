package internal

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
)

// TestStruct is a test struct with various field types.
type TestStruct struct {
	// Simple fields
	Name    string
	Enabled string
	Count   int

	// Nested struct
	Driver struct {
		Version string
		Enabled string
	}

	// Acronym fields
	EnableGDS string
	MIG       struct {
		Strategy string
	}
	GPUOperator struct {
		Version string
	}

	// Complex nested
	DCGM struct {
		Exporter struct {
			Version string
			Enabled string
		}
	}
}

func TestApplyValueOverrides_SimpleFields(t *testing.T) {
	tests := []struct {
		name      string
		overrides map[string]string
		want      TestStruct
		wantErr   bool
	}{
		{
			name: "set string field",
			overrides: map[string]string{
				"name": "test-value",
			},
			want: TestStruct{
				Name: "test-value",
			},
		},
		{
			name: "set enabled field",
			overrides: map[string]string{
				"enabled": "true",
			},
			want: TestStruct{
				Enabled: "true",
			},
		},
		{
			name: "set multiple fields",
			overrides: map[string]string{
				"name":    "test",
				"enabled": "false",
			},
			want: TestStruct{
				Name:    "test",
				Enabled: "false",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TestStruct{}
			err := ApplyValueOverrides(&got, tt.overrides)

			if (err != nil) != tt.wantErr {
				t.Errorf("ApplyValueOverrides() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got.Name != tt.want.Name {
				t.Errorf("Name = %v, want %v", got.Name, tt.want.Name)
			}
			if got.Enabled != tt.want.Enabled {
				t.Errorf("Enabled = %v, want %v", got.Enabled, tt.want.Enabled)
			}
		})
	}
}

func TestApplyValueOverrides_NestedFields(t *testing.T) {
	tests := []struct {
		name      string
		overrides map[string]string
		want      TestStruct
		wantErr   bool
	}{
		{
			name: "set nested field",
			overrides: map[string]string{
				"driver.version": "550.127",
			},
			want: TestStruct{
				Driver: struct {
					Version string
					Enabled string
				}{
					Version: "550.127",
				},
			},
		},
		{
			name: "set multiple nested fields",
			overrides: map[string]string{
				"driver.version": "550.127",
				"driver.enabled": "true",
			},
			want: TestStruct{
				Driver: struct {
					Version string
					Enabled string
				}{
					Version: "550.127",
					Enabled: "true",
				},
			},
		},
		{
			name: "set deeply nested field",
			overrides: map[string]string{
				"dcgm.exporter.version": "3.3.11",
				"dcgm.exporter.enabled": "true",
			},
			want: TestStruct{
				DCGM: struct {
					Exporter struct {
						Version string
						Enabled string
					}
				}{
					Exporter: struct {
						Version string
						Enabled string
					}{
						Version: "3.3.11",
						Enabled: "true",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TestStruct{}
			err := ApplyValueOverrides(&got, tt.overrides)

			if (err != nil) != tt.wantErr {
				t.Errorf("ApplyValueOverrides() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want.Driver.Version != "" && got.Driver.Version != tt.want.Driver.Version {
				t.Errorf("Driver.Version = %v, want %v", got.Driver.Version, tt.want.Driver.Version)
			}
			if tt.want.Driver.Enabled != "" && got.Driver.Enabled != tt.want.Driver.Enabled {
				t.Errorf("Driver.Enabled = %v, want %v", got.Driver.Enabled, tt.want.Driver.Enabled)
			}
			if tt.want.DCGM.Exporter.Version != "" && got.DCGM.Exporter.Version != tt.want.DCGM.Exporter.Version {
				t.Errorf("DCGM.Exporter.Version = %v, want %v", got.DCGM.Exporter.Version, tt.want.DCGM.Exporter.Version)
			}
		})
	}
}

func TestApplyValueOverrides_AcronymFields(t *testing.T) {
	tests := []struct {
		name      string
		overrides map[string]string
		want      TestStruct
		wantErr   bool
	}{
		{
			name: "set MIG strategy",
			overrides: map[string]string{
				"mig.strategy": "mixed",
			},
			want: TestStruct{
				MIG: struct {
					Strategy string
				}{
					Strategy: "mixed",
				},
			},
		},
		{
			name: "set GPU operator version",
			overrides: map[string]string{
				"gpu-operator.version": "25.3.3",
			},
			want: TestStruct{
				GPUOperator: struct {
					Version string
				}{
					Version: "25.3.3",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TestStruct{}
			err := ApplyValueOverrides(&got, tt.overrides)

			if (err != nil) != tt.wantErr {
				t.Errorf("ApplyValueOverrides() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want.MIG.Strategy != "" && got.MIG.Strategy != tt.want.MIG.Strategy {
				t.Errorf("MIG.Strategy = %v, want %v", got.MIG.Strategy, tt.want.MIG.Strategy)
			}
			if tt.want.GPUOperator.Version != "" && got.GPUOperator.Version != tt.want.GPUOperator.Version {
				t.Errorf("GPUOperator.Version = %v, want %v", got.GPUOperator.Version, tt.want.GPUOperator.Version)
			}
		})
	}
}

func TestApplyValueOverrides_Errors(t *testing.T) {
	tests := []struct {
		name      string
		target    interface{}
		overrides map[string]string
		wantErr   bool
		errMsg    string
	}{
		{
			name:   "non-pointer target",
			target: TestStruct{},
			overrides: map[string]string{
				"name": "test",
			},
			wantErr: true,
			errMsg:  "must be a pointer",
		},
		{
			name:      "nil overrides",
			target:    &TestStruct{},
			overrides: nil,
			wantErr:   false,
		},
		{
			name:      "empty overrides",
			target:    &TestStruct{},
			overrides: map[string]string{},
			wantErr:   false,
		},
		{
			name:   "non-existent field",
			target: &TestStruct{},
			overrides: map[string]string{
				"nonexistent": "value",
			},
			wantErr: true,
			errMsg:  "field not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ApplyValueOverrides(tt.target, tt.overrides)

			if (err != nil) != tt.wantErr {
				t.Errorf("ApplyValueOverrides() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errMsg != "" {
				if err == nil || !containsSubstring(err.Error(), tt.errMsg) {
					t.Errorf("Expected error containing %q, got %v", tt.errMsg, err)
				}
			}
		})
	}
}

func TestApplyValueOverrides_CaseInsensitive(t *testing.T) {
	tests := []struct {
		name      string
		overrides map[string]string
		want      TestStruct
	}{
		{
			name: "lowercase field name",
			overrides: map[string]string{
				"name": "test",
			},
			want: TestStruct{
				Name: "test",
			},
		},
		{
			name: "uppercase field name",
			overrides: map[string]string{
				"NAME": "test",
			},
			want: TestStruct{
				Name: "test",
			},
		},
		{
			name: "mixed case field name",
			overrides: map[string]string{
				"NaMe": "test",
			},
			want: TestStruct{
				Name: "test",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TestStruct{}
			err := ApplyValueOverrides(&got, tt.overrides)

			if err != nil {
				t.Errorf("ApplyValueOverrides() unexpected error = %v", err)
				return
			}

			if got.Name != tt.want.Name {
				t.Errorf("Name = %v, want %v", got.Name, tt.want.Name)
			}
		})
	}
}

// Test with actual GPU Operator-like struct
type GPUOperatorValues struct {
	EnableDriver string
	Driver       struct {
		Version string
		Enabled string
	}
	EnableGDS string
	GDS       struct {
		Enabled string
	}
	GDRCopy struct {
		Enabled string
	}
	MIG struct {
		Strategy string
	}
	DCGM struct {
		Version string
	}
}

func TestApplyValueOverrides_GPUOperatorScenarios(t *testing.T) {
	tests := []struct {
		name      string
		overrides map[string]string
		verify    func(t *testing.T, values *GPUOperatorValues)
	}{
		{
			name: "gdrcopy enabled override",
			overrides: map[string]string{
				"gdrcopy.enabled": "false",
			},
			verify: func(t *testing.T, values *GPUOperatorValues) {
				if values.GDRCopy.Enabled != "false" {
					t.Errorf("GDRCopy.Enabled = %v, want false", values.GDRCopy.Enabled)
				}
			},
		},
		{
			name: "gds enabled override",
			overrides: map[string]string{
				"gds.enabled": "true",
			},
			verify: func(t *testing.T, values *GPUOperatorValues) {
				// Should match either EnableGDS or GDS.Enabled
				if values.EnableGDS != "true" && values.GDS.Enabled != "true" {
					t.Errorf("GDS not enabled: EnableGDS=%v, GDS.Enabled=%v", values.EnableGDS, values.GDS.Enabled)
				}
			},
		},
		{
			name: "driver version override",
			overrides: map[string]string{
				"driver.version": "570.86.16",
			},
			verify: func(t *testing.T, values *GPUOperatorValues) {
				if values.Driver.Version != "570.86.16" {
					t.Errorf("Driver.Version = %v, want 570.86.16", values.Driver.Version)
				}
			},
		},
		{
			name: "mig strategy override",
			overrides: map[string]string{
				"mig.strategy": "mixed",
			},
			verify: func(t *testing.T, values *GPUOperatorValues) {
				if values.MIG.Strategy != "mixed" {
					t.Errorf("MIG.Strategy = %v, want mixed", values.MIG.Strategy)
				}
			},
		},
		{
			name: "multiple overrides",
			overrides: map[string]string{
				"gdrcopy.enabled": "false",
				"gds.enabled":     "true",
				"driver.version":  "570.86.16",
				"mig.strategy":    "mixed",
			},
			verify: func(t *testing.T, values *GPUOperatorValues) {
				if values.GDRCopy.Enabled != "false" {
					t.Errorf("GDRCopy.Enabled = %v, want false", values.GDRCopy.Enabled)
				}
				if values.Driver.Version != "570.86.16" {
					t.Errorf("Driver.Version = %v, want 570.86.16", values.Driver.Version)
				}
				if values.MIG.Strategy != "mixed" {
					t.Errorf("MIG.Strategy = %v, want mixed", values.MIG.Strategy)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values := &GPUOperatorValues{}
			err := ApplyValueOverrides(values, tt.overrides)

			if err != nil {
				t.Fatalf("ApplyValueOverrides() unexpected error = %v", err)
			}

			tt.verify(t, values)
		})
	}
}

// containsSubstring checks if string contains substring (renamed to avoid collision)
func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsSubstringHelper(s, substr))
}

func containsSubstringHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestApplyNodeSelectorOverrides(t *testing.T) {
	tests := []struct {
		name         string
		values       map[string]interface{}
		nodeSelector map[string]string
		paths        []string
		verify       func(t *testing.T, values map[string]interface{})
	}{
		{
			name:   "applies to top-level nodeSelector",
			values: make(map[string]interface{}),
			nodeSelector: map[string]string{
				"nodeGroup": "system-cpu",
			},
			paths: []string{"nodeSelector"},
			verify: func(t *testing.T, values map[string]interface{}) {
				ns, ok := values["nodeSelector"].(map[string]interface{})
				if !ok {
					t.Fatal("nodeSelector not found or wrong type")
				}
				if ns["nodeGroup"] != "system-cpu" {
					t.Errorf("nodeSelector.nodeGroup = %v, want system-cpu", ns["nodeGroup"])
				}
			},
		},
		{
			name: "applies to nested paths",
			values: map[string]interface{}{
				"webhook": make(map[string]interface{}),
			},
			nodeSelector: map[string]string{
				"role": "control-plane",
			},
			paths: []string{"nodeSelector", "webhook.nodeSelector"},
			verify: func(t *testing.T, values map[string]interface{}) {
				// Check top-level
				ns, ok := values["nodeSelector"].(map[string]interface{})
				if !ok {
					t.Fatal("nodeSelector not found")
				}
				if ns["role"] != "control-plane" {
					t.Errorf("nodeSelector.role = %v, want control-plane", ns["role"])
				}
				// Check nested
				wh, ok := values["webhook"].(map[string]interface{})
				if !ok {
					t.Fatal("webhook not found")
				}
				whNs, ok := wh["nodeSelector"].(map[string]interface{})
				if !ok {
					t.Fatal("webhook.nodeSelector not found")
				}
				if whNs["role"] != "control-plane" {
					t.Errorf("webhook.nodeSelector.role = %v, want control-plane", whNs["role"])
				}
			},
		},
		{
			name:         "empty nodeSelector is no-op",
			values:       make(map[string]interface{}),
			nodeSelector: map[string]string{},
			paths:        []string{"nodeSelector"},
			verify: func(t *testing.T, values map[string]interface{}) {
				if _, ok := values["nodeSelector"]; ok {
					t.Error("nodeSelector should not be set for empty input")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ApplyNodeSelectorOverrides(tt.values, tt.nodeSelector, tt.paths...)
			tt.verify(t, tt.values)
		})
	}
}

func TestApplyTolerationsOverrides(t *testing.T) {
	tests := []struct {
		name        string
		values      map[string]interface{}
		tolerations []corev1.Toleration
		paths       []string
		verify      func(t *testing.T, values map[string]interface{})
	}{
		{
			name:   "applies single toleration",
			values: make(map[string]interface{}),
			tolerations: []corev1.Toleration{
				{
					Key:      "dedicated",
					Value:    "system-workload",
					Operator: corev1.TolerationOpEqual,
					Effect:   corev1.TaintEffectNoSchedule,
				},
			},
			paths: []string{"tolerations"},
			verify: func(t *testing.T, values map[string]interface{}) {
				tols, ok := values["tolerations"].([]interface{})
				if !ok {
					t.Fatal("tolerations not found or wrong type")
				}
				if len(tols) != 1 {
					t.Fatalf("expected 1 toleration, got %d", len(tols))
				}
				tol, ok := tols[0].(map[string]interface{})
				if !ok {
					t.Fatal("toleration entry wrong type")
				}
				if tol["key"] != "dedicated" {
					t.Errorf("key = %v, want dedicated", tol["key"])
				}
				if tol["value"] != "system-workload" {
					t.Errorf("value = %v, want system-workload", tol["value"])
				}
			},
		},
		{
			name: "applies to nested paths",
			values: map[string]interface{}{
				"webhook": make(map[string]interface{}),
			},
			tolerations: []corev1.Toleration{
				{Operator: corev1.TolerationOpExists},
			},
			paths: []string{"tolerations", "webhook.tolerations"},
			verify: func(t *testing.T, values map[string]interface{}) {
				// Check top-level
				tols, ok := values["tolerations"].([]interface{})
				if !ok {
					t.Fatal("tolerations not found")
				}
				if len(tols) != 1 {
					t.Fatalf("expected 1 toleration, got %d", len(tols))
				}
				// Check nested
				wh, ok := values["webhook"].(map[string]interface{})
				if !ok {
					t.Fatal("webhook not found")
				}
				whTols, ok := wh["tolerations"].([]interface{})
				if !ok {
					t.Fatal("webhook.tolerations not found")
				}
				if len(whTols) != 1 {
					t.Fatalf("expected 1 webhook toleration, got %d", len(whTols))
				}
			},
		},
		{
			name:        "empty tolerations is no-op",
			values:      make(map[string]interface{}),
			tolerations: []corev1.Toleration{},
			paths:       []string{"tolerations"},
			verify: func(t *testing.T, values map[string]interface{}) {
				if _, ok := values["tolerations"]; ok {
					t.Error("tolerations should not be set for empty input")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ApplyTolerationsOverrides(tt.values, tt.tolerations, tt.paths...)
			tt.verify(t, tt.values)
		})
	}
}

func TestTolerationsToPodSpec(t *testing.T) {
	tests := []struct {
		name        string
		tolerations []corev1.Toleration
		verify      func(t *testing.T, result []map[string]interface{})
	}{
		{
			name: "converts full toleration",
			tolerations: []corev1.Toleration{
				{
					Key:      "dedicated",
					Operator: corev1.TolerationOpEqual,
					Value:    "gpu",
					Effect:   corev1.TaintEffectNoSchedule,
				},
			},
			verify: func(t *testing.T, result []map[string]interface{}) {
				if len(result) != 1 {
					t.Fatalf("expected 1 result, got %d", len(result))
				}
				tol := result[0]
				if tol["key"] != "dedicated" {
					t.Errorf("key = %v, want dedicated", tol["key"])
				}
				if tol["operator"] != "Equal" {
					t.Errorf("operator = %v, want Equal", tol["operator"])
				}
				if tol["value"] != "gpu" {
					t.Errorf("value = %v, want gpu", tol["value"])
				}
				if tol["effect"] != "NoSchedule" {
					t.Errorf("effect = %v, want NoSchedule", tol["effect"])
				}
			},
		},
		{
			name: "omits empty fields",
			tolerations: []corev1.Toleration{
				{Operator: corev1.TolerationOpExists},
			},
			verify: func(t *testing.T, result []map[string]interface{}) {
				if len(result) != 1 {
					t.Fatalf("expected 1 result, got %d", len(result))
				}
				tol := result[0]
				if _, ok := tol["key"]; ok {
					t.Error("key should be omitted when empty")
				}
				if tol["operator"] != "Exists" {
					t.Errorf("operator = %v, want Exists", tol["operator"])
				}
				if _, ok := tol["value"]; ok {
					t.Error("value should be omitted when empty")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TolerationsToPodSpec(tt.tolerations)
			tt.verify(t, result)
		})
	}
}

func TestNodeSelectorToMatchExpressions(t *testing.T) {
	tests := []struct {
		name         string
		nodeSelector map[string]string
		verify       func(t *testing.T, result []map[string]interface{})
	}{
		{
			name: "converts single selector",
			nodeSelector: map[string]string{
				"nodeGroup": "gpu-nodes",
			},
			verify: func(t *testing.T, result []map[string]interface{}) {
				if len(result) != 1 {
					t.Fatalf("expected 1 expression, got %d", len(result))
				}
				expr := result[0]
				if expr["key"] != "nodeGroup" {
					t.Errorf("key = %v, want nodeGroup", expr["key"])
				}
				if expr["operator"] != "In" {
					t.Errorf("operator = %v, want In", expr["operator"])
				}
				values, ok := expr["values"].([]string)
				if !ok {
					t.Fatal("values not a []string")
				}
				if len(values) != 1 || values[0] != "gpu-nodes" {
					t.Errorf("values = %v, want [gpu-nodes]", values)
				}
			},
		},
		{
			name: "converts multiple selectors",
			nodeSelector: map[string]string{
				"nodeGroup":   "gpu-nodes",
				"accelerator": "nvidia-h100",
			},
			verify: func(t *testing.T, result []map[string]interface{}) {
				if len(result) != 2 {
					t.Fatalf("expected 2 expressions, got %d", len(result))
				}
				// Check both expressions exist (order may vary due to map iteration)
				foundNodeGroup := false
				foundAccelerator := false
				for _, expr := range result {
					if expr["key"] == "nodeGroup" {
						foundNodeGroup = true
						values := expr["values"].([]string)
						if values[0] != "gpu-nodes" {
							t.Errorf("nodeGroup values = %v, want [gpu-nodes]", values)
						}
					}
					if expr["key"] == "accelerator" {
						foundAccelerator = true
						values := expr["values"].([]string)
						if values[0] != "nvidia-h100" {
							t.Errorf("accelerator values = %v, want [nvidia-h100]", values)
						}
					}
				}
				if !foundNodeGroup {
					t.Error("missing nodeGroup expression")
				}
				if !foundAccelerator {
					t.Error("missing accelerator expression")
				}
			},
		},
		{
			name:         "returns nil for empty selector",
			nodeSelector: map[string]string{},
			verify: func(t *testing.T, result []map[string]interface{}) {
				if result != nil {
					t.Errorf("expected nil, got %v", result)
				}
			},
		},
		{
			name:         "returns nil for nil selector",
			nodeSelector: nil,
			verify: func(t *testing.T, result []map[string]interface{}) {
				if result != nil {
					t.Errorf("expected nil, got %v", result)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NodeSelectorToMatchExpressions(tt.nodeSelector)
			tt.verify(t, result)
		})
	}
}
