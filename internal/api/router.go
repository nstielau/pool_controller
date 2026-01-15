package api

import (
	"net/http"

	"github.com/nstielau/pool-controller/internal/pool"
)

// Router sets up the HTTP routes for the pool controller.
type Router struct {
	mux          *http.ServeMux
	poolHandler  *PoolHandler
	alexaHandler http.Handler
}

// NewRouter creates a new Router with all routes configured.
func NewRouter(bridge *pool.Bridge, alexaHandler http.Handler) *Router {
	r := &Router{
		mux:          http.NewServeMux(),
		poolHandler:  NewPoolHandler(bridge),
		alexaHandler: alexaHandler,
	}

	r.setupRoutes()
	return r
}

// setupRoutes configures all HTTP routes.
func (r *Router) setupRoutes() {
	// Health check (no auth)
	r.mux.HandleFunc("GET /", r.poolHandler.HandleIndex)

	// Alexa skill endpoint (POST /, no auth - Alexa verifies itself)
	if r.alexaHandler != nil {
		r.mux.Handle("POST /", r.alexaHandler)
	}

	// Pool endpoints with authentication
	authPool := AuthMiddleware(http.HandlerFunc(r.poolHandler.HandlePool))
	r.mux.Handle("GET /pool", authPool)

	// Pool attribute endpoint
	authPoolAttr := AuthMiddleware(http.HandlerFunc(r.poolHandler.HandlePoolAttribute))
	r.mux.Handle("GET /pool/", authPoolAttr)
}

// ServeHTTP implements the http.Handler interface.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

// Handler returns the HTTP handler for the router.
func (r *Router) Handler() http.Handler {
	return r.mux
}
