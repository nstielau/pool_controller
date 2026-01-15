package alexa

import (
	"testing"
	"time"
)

func TestValidateCertURL(t *testing.T) {
	v := NewVerifier()

	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid URL",
			url:     "https://s3.amazonaws.com/echo.api/echo-api-cert.pem",
			wantErr: false,
		},
		{
			name:    "valid URL with subdomain",
			url:     "https://s3.amazonaws.com/echo.api/echo-api-cert-4.pem",
			wantErr: false,
		},
		{
			name:    "HTTP not allowed",
			url:     "http://s3.amazonaws.com/echo.api/echo-api-cert.pem",
			wantErr: true,
		},
		{
			name:    "wrong host",
			url:     "https://example.com/echo.api/echo-api-cert.pem",
			wantErr: true,
		},
		{
			name:    "wrong path",
			url:     "https://s3.amazonaws.com/other/echo-api-cert.pem",
			wantErr: true,
		},
		{
			name:    "wrong port",
			url:     "https://s3.amazonaws.com:8080/echo.api/echo-api-cert.pem",
			wantErr: true,
		},
		{
			name:    "port 443 allowed",
			url:     "https://s3.amazonaws.com:443/echo.api/echo-api-cert.pem",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.validateCertURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateCertURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVerifyTimestamp(t *testing.T) {
	tests := []struct {
		name      string
		timestamp string
		tolerance time.Duration
		wantErr   bool
	}{
		{
			name:      "current time",
			timestamp: time.Now().UTC().Format(time.RFC3339),
			tolerance: 150 * time.Second,
			wantErr:   false,
		},
		{
			name:      "within tolerance",
			timestamp: time.Now().Add(-60 * time.Second).UTC().Format(time.RFC3339),
			tolerance: 150 * time.Second,
			wantErr:   false,
		},
		{
			name:      "outside tolerance",
			timestamp: time.Now().Add(-200 * time.Second).UTC().Format(time.RFC3339),
			tolerance: 150 * time.Second,
			wantErr:   true,
		},
		{
			name:      "future within tolerance",
			timestamp: time.Now().Add(60 * time.Second).UTC().Format(time.RFC3339),
			tolerance: 150 * time.Second,
			wantErr:   false,
		},
		{
			name:      "invalid format",
			timestamp: "not-a-timestamp",
			tolerance: 150 * time.Second,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := VerifyTimestamp(tt.timestamp, tt.tolerance)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyTimestamp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewVerifier(t *testing.T) {
	v := NewVerifier()
	if v == nil {
		t.Error("NewVerifier() returned nil")
	}
	if v.certCache == nil {
		t.Error("certCache should be initialized")
	}
}
