package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	// Handler that returns 200 OK if auth passes
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	tests := []struct {
		name       string
		tokenRegex string
		authHeader string
		wantBody   string
	}{
		{
			name:       "default regex accepts any token",
			tokenRegex: "",
			authHeader: "Bearer anytoken",
			wantBody:   "success",
		},
		{
			name:       "exact match success",
			tokenRegex: "^secret123$",
			authHeader: "Bearer secret123",
			wantBody:   "success",
		},
		{
			name:       "exact match failure",
			tokenRegex: "^secret123$",
			authHeader: "Bearer wrongtoken",
			wantBody:   "Unauthed",
		},
		{
			name:       "no auth header with permissive regex",
			tokenRegex: ".*",
			authHeader: "",
			wantBody:   "success",
		},
		{
			name:       "prefix match",
			tokenRegex: "^pool-",
			authHeader: "Bearer pool-abc123",
			wantBody:   "success",
		},
		{
			name:       "prefix match failure",
			tokenRegex: "^pool-",
			authHeader: "Bearer notpool-abc",
			wantBody:   "Unauthed",
		},
		{
			name:       "legacy Authentication header",
			tokenRegex: "^test$",
			authHeader: "Bearer test",
			wantBody:   "success",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set TOKEN_REGEX env var
			if tt.tokenRegex != "" {
				os.Setenv("TOKEN_REGEX", tt.tokenRegex)
			} else {
				os.Unsetenv("TOKEN_REGEX")
			}

			// Create middleware
			handler := AuthMiddleware(nextHandler)

			// Create request
			req := httptest.NewRequest("GET", "/pool", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Record response
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Body.String() != tt.wantBody {
				t.Errorf("body = %q, want %q", rr.Body.String(), tt.wantBody)
			}
		})
	}

	// Clean up
	os.Unsetenv("TOKEN_REGEX")
}

func TestAuthMiddlewareInvalidRegex(t *testing.T) {
	os.Setenv("TOKEN_REGEX", "[invalid")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("success"))
	})

	handler := AuthMiddleware(nextHandler)

	req := httptest.NewRequest("GET", "/pool", nil)
	req.Header.Set("Authorization", "Bearer anytoken")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Invalid regex should fall back to permissive pattern
	if rr.Body.String() != "success" {
		t.Errorf("invalid regex should fall back to permissive, got %q", rr.Body.String())
	}

	os.Unsetenv("TOKEN_REGEX")
}
