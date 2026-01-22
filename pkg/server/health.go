package server

import (
	"net/http"
	"time"

	"github.com/NVIDIA/cloud-native-stack/pkg/serializer"
)

// HealthResponse represents health check response
type HealthResponse struct {
	Status    string    `json:"status" yaml:"status"`
	Timestamp time.Time `json:"timestamp" yaml:"timestamp"`
	Reason    string    `json:"reason,omitempty" yaml:"reason,omitempty"`
}

// handleHealth handles GET /health
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	resp := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
	}

	serializer.RespondJSON(w, http.StatusOK, resp)
}

// handleReady handles GET /ready
func (s *Server) handleReady(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.mu.RLock()
	ready := s.ready
	s.mu.RUnlock()

	if !ready {
		resp := HealthResponse{
			Status:    "not_ready",
			Timestamp: time.Now(),
			Reason:    "service is initializing",
		}
		w.WriteHeader(http.StatusServiceUnavailable)
		serializer.RespondJSON(w, http.StatusServiceUnavailable, resp)
		return
	}

	resp := HealthResponse{
		Status:    "ready",
		Timestamp: time.Now(),
	}

	serializer.RespondJSON(w, http.StatusOK, resp)
}
