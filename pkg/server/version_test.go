package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNegotiateAPIVersion(t *testing.T) {
	tests := []struct {
		name   string
		accept string
		want   string
	}{
		{"empty accept defaults", "", DefaultAPIVersion},
		{"non-vendor accept defaults", "application/json", DefaultAPIVersion},
		{"vendor v1", "application/vnd.nvidia.cns.v1+json", "v1"},
		{"vendor v2 unsupported defaults", "application/vnd.nvidia.cns.v2+json", DefaultAPIVersion},
		{"vendor malformed defaults", "application/vnd.nvidia.cns.vBAD+json", DefaultAPIVersion},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.accept != "" {
				req.Header.Set("Accept", tt.accept)
			}
			if got := negotiateAPIVersion(req); got != tt.want {
				t.Fatalf("negotiateAPIVersion(Accept=%q) = %q, want %q", tt.accept, got, tt.want)
			}
		})
	}
}

func TestIsValidAPIVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    bool
	}{
		{"v1 valid", "v1", true},
		{"v2 invalid", "v2", false},
		{"empty invalid", "", false},
		{"random invalid", "nope", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidAPIVersion(tt.version); got != tt.want {
				t.Fatalf("isValidAPIVersion(%q) = %v, want %v", tt.version, got, tt.want)
			}
		})
	}
}
