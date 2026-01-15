// Package alexa implements Alexa skill request handling and verification.
package alexa

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Verifier handles Alexa request signature verification.
type Verifier struct {
	certCache map[string]*x509.Certificate
	cacheMu   sync.RWMutex
}

// NewVerifier creates a new Alexa request verifier.
func NewVerifier() *Verifier {
	return &Verifier{
		certCache: make(map[string]*x509.Certificate),
	}
}

// VerifyRequest validates an Alexa skill request.
// See: https://developer.amazon.com/docs/alexa/custom-skills/host-a-custom-skill-as-a-web-service.html
func (v *Verifier) VerifyRequest(r *http.Request, body []byte) error {
	// 1. Verify the URL for the signing certificate
	certURL := r.Header.Get("SignatureCertChainUrl")
	if certURL == "" {
		return fmt.Errorf("missing SignatureCertChainUrl header")
	}

	if err := v.validateCertURL(certURL); err != nil {
		return fmt.Errorf("invalid certificate URL: %w", err)
	}

	// 2. Get the signature
	signature := r.Header.Get("Signature-256")
	if signature == "" {
		// Fall back to SHA1 signature for older requests
		signature = r.Header.Get("Signature")
	}
	if signature == "" {
		return fmt.Errorf("missing Signature header")
	}

	// 3. Download and validate the certificate
	cert, err := v.getCertificate(certURL)
	if err != nil {
		return fmt.Errorf("failed to get certificate: %w", err)
	}

	// 4. Verify the signature
	if err := v.verifySignature(cert, signature, body); err != nil {
		return fmt.Errorf("signature verification failed: %w", err)
	}

	return nil
}

// validateCertURL checks that the certificate URL is valid for Alexa.
func (v *Verifier) validateCertURL(certURL string) error {
	u, err := url.Parse(certURL)
	if err != nil {
		return err
	}

	// Scheme must be HTTPS
	if strings.ToLower(u.Scheme) != "https" {
		return fmt.Errorf("certificate URL must use HTTPS")
	}

	// Host must be s3.amazonaws.com (case insensitive)
	// Use Hostname() to strip port if present
	hostname := strings.ToLower(u.Hostname())
	if hostname != "s3.amazonaws.com" &&
		!strings.HasSuffix(hostname, ".s3.amazonaws.com") {
		return fmt.Errorf("certificate URL host must be s3.amazonaws.com")
	}

	// Path must start with /echo.api/
	if !strings.HasPrefix(u.Path, "/echo.api/") {
		return fmt.Errorf("certificate URL path must start with /echo.api/")
	}

	// Port must be 443 or empty
	if u.Port() != "" && u.Port() != "443" {
		return fmt.Errorf("certificate URL port must be 443")
	}

	return nil
}

// getCertificate fetches and caches the signing certificate.
func (v *Verifier) getCertificate(certURL string) (*x509.Certificate, error) {
	// Check cache first
	v.cacheMu.RLock()
	if cert, ok := v.certCache[certURL]; ok {
		v.cacheMu.RUnlock()
		// Check if still valid
		if time.Now().Before(cert.NotAfter) {
			return cert, nil
		}
		// Expired, need to refresh
	} else {
		v.cacheMu.RUnlock()
	}

	// Download certificate
	resp, err := http.Get(certURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	certData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse PEM certificate
	block, _ := pem.Decode(certData)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	// Validate certificate
	if err := v.validateCertificate(cert); err != nil {
		return nil, err
	}

	// Cache the certificate
	v.cacheMu.Lock()
	v.certCache[certURL] = cert
	v.cacheMu.Unlock()

	return cert, nil
}

// validateCertificate checks that the certificate is valid for Alexa.
func (v *Verifier) validateCertificate(cert *x509.Certificate) error {
	// Check expiration
	now := time.Now()
	if now.Before(cert.NotBefore) || now.After(cert.NotAfter) {
		return fmt.Errorf("certificate is not valid at current time")
	}

	// Check that the domain echo-api.amazon.com is in the SAN
	found := false
	for _, name := range cert.DNSNames {
		if name == "echo-api.amazon.com" {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("certificate does not include echo-api.amazon.com in SAN")
	}

	return nil
}

// verifySignature checks the request signature.
func (v *Verifier) verifySignature(cert *x509.Certificate, signature string, body []byte) error {
	// Decode base64 signature
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("failed to decode signature: %w", err)
	}

	// Get public key
	pubKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("certificate does not contain RSA public key")
	}

	// Compute hash of body
	hash := sha1.Sum(body)

	// Verify signature
	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA1, hash[:], sig)
	if err != nil {
		return fmt.Errorf("RSA verification failed: %w", err)
	}

	return nil
}

// VerifyTimestamp checks that the request timestamp is within tolerance.
func VerifyTimestamp(timestamp string, tolerance time.Duration) error {
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return fmt.Errorf("failed to parse timestamp: %w", err)
	}

	diff := time.Since(t)
	if diff < 0 {
		diff = -diff
	}

	if diff > tolerance {
		return fmt.Errorf("request timestamp is outside tolerance: %v", diff)
	}

	return nil
}
