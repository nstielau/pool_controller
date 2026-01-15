package api

import (
	"encoding/json"
	"net/http"

	"github.com/nstielau/pool-controller/internal/pool"
)

// PoolHandler handles pool-related HTTP requests.
type PoolHandler struct {
	bridge *pool.Bridge
}

// NewPoolHandler creates a new PoolHandler.
func NewPoolHandler(bridge *pool.Bridge) *PoolHandler {
	return &PoolHandler{bridge: bridge}
}

// HandleIndex is the health check endpoint (GET /).
func (h *PoolHandler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("hello"))
}

// HandlePool returns the full pool status as JSON (GET /pool).
func (h *PoolHandler) HandlePool(w http.ResponseWriter, r *http.Request) {
	// Refresh data if needed
	h.bridge.Update()

	jsonData, err := h.bridge.GetJSON()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonData))
}

// HandlePoolAttribute returns a specific pool attribute (GET /pool/{attribute}).
func (h *PoolHandler) HandlePoolAttribute(w http.ResponseWriter, r *http.Request) {
	// Extract attribute from URL path
	// Path is like /pool/spa or /pool/current_spa_temperature
	path := r.URL.Path
	if len(path) <= 6 { // "/pool/"
		http.Error(w, "attribute required", http.StatusBadRequest)
		return
	}
	attribute := path[6:] // Strip "/pool/"

	// Refresh data if needed
	h.bridge.Update()

	// Get the attribute
	data, ok := h.bridge.GetAttribute(attribute)
	if !ok {
		http.Error(w, "attribute not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
