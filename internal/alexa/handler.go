package alexa

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/nstielau/pool-controller/internal/gateway"
	"github.com/nstielau/pool-controller/internal/pool"
)

// Request is the Alexa skill request structure.
type Request struct {
	Version string         `json:"version"`
	Session SessionRequest `json:"session"`
	Request RequestBody    `json:"request"`
}

// SessionRequest contains session information.
type SessionRequest struct {
	SessionID   string                 `json:"sessionId"`
	Application ApplicationRequest     `json:"application"`
	Attributes  map[string]interface{} `json:"attributes"`
	User        UserRequest            `json:"user"`
	New         bool                   `json:"new"`
}

// ApplicationRequest contains skill application info.
type ApplicationRequest struct {
	ApplicationID string `json:"applicationId"`
}

// UserRequest contains user information.
type UserRequest struct {
	UserID string `json:"userId"`
}

// RequestBody contains the actual request content.
type RequestBody struct {
	Type      string `json:"type"`
	RequestID string `json:"requestId"`
	Timestamp string `json:"timestamp"`
	Locale    string `json:"locale"`
	Intent    Intent `json:"intent,omitempty"`
}

// Intent contains intent information.
type Intent struct {
	Name               string                 `json:"name"`
	ConfirmationStatus string                 `json:"confirmationStatus"`
	Slots              map[string]interface{} `json:"slots"`
}

// Handler handles Alexa skill requests.
type Handler struct {
	bridge     *pool.Bridge
	verifier   *Verifier
	logger     *log.Logger
	skipVerify bool
}

// NewHandler creates a new Alexa skill handler.
func NewHandler(bridge *pool.Bridge) *Handler {
	skipVerify := os.Getenv("ALEXA_SKIP_VERIFY") == "true"
	return &Handler{
		bridge:     bridge,
		verifier:   NewVerifier(),
		logger:     log.New(os.Stdout, "[alexa] ", log.LstdFlags),
		skipVerify: skipVerify,
	}
}

// ServeHTTP handles Alexa skill HTTP requests.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Printf("Failed to read request body: %v", err)
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	// Verify request signature (skip in development)
	if !h.skipVerify {
		if err := h.verifier.VerifyRequest(r, body); err != nil {
			h.logger.Printf("Request verification failed: %v", err)
			http.Error(w, "verification failed", http.StatusUnauthorized)
			return
		}
	}

	// Parse request
	var req Request
	if err := json.Unmarshal(body, &req); err != nil {
		h.logger.Printf("Failed to parse request: %v", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Verify timestamp (within 150 seconds)
	if !h.skipVerify {
		if err := VerifyTimestamp(req.Request.Timestamp, 150*time.Second); err != nil {
			h.logger.Printf("Timestamp verification failed: %v", err)
			http.Error(w, "timestamp verification failed", http.StatusUnauthorized)
			return
		}
	}

	// Log request details
	sessionID := req.Session.SessionID
	if len(sessionID) > 20 {
		sessionID = sessionID[:20] + "..."
	}
	h.logger.Printf("Request: type=%s intent=%s session=%s new=%v locale=%s",
		req.Request.Type,
		req.Request.Intent.Name,
		sessionID,
		req.Session.New,
		req.Request.Locale)

	// Handle request
	startTime := time.Now()
	var response *Response
	switch req.Request.Type {
	case "LaunchRequest":
		response = h.handleLaunchRequest()
	case "IntentRequest":
		response = h.handleIntent(req.Request.Intent.Name)
	case "SessionEndedRequest":
		response = SpeakResponse("Goodbye!", true)
	default:
		response = SpeakResponse("I don't know how to handle that.", true)
	}

	// Log response
	responseText := ""
	if response.Response.OutputSpeech != nil {
		responseText = response.Response.OutputSpeech.Text
	}
	h.logger.Printf("Response: text=%q endSession=%v duration=%v",
		responseText,
		response.Response.ShouldEndSession,
		time.Since(startTime))

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleLaunchRequest handles skill launch.
func (h *Handler) handleLaunchRequest() *Response {
	// Keep session open (false) so user can follow up with commands
	return SpeakResponse("Pool party time. Do you want to check the hot tub temp or turn on the pool jets?", false)
}

// handleIntent routes to the appropriate intent handler.
func (h *Handler) handleIntent(intentName string) *Response {
	switch intentName {
	case "StartSwimJetIntent":
		return h.handleStartSwimJet()
	case "StopSwimJetIntent":
		return h.handleStopSwimJet()
	case "StartHotTubIntent":
		return h.handleStartHotTub()
	case "StopHotTubIntent":
		return h.handleStopHotTub()
	case "HotTubTempIntent":
		return h.handleHotTubTemp()
	case "AMAZON.CancelIntent", "AMAZON.StopIntent":
		return SpeakResponse("Party on!", true)
	case "AMAZON.HelpIntent":
		return SpeakResponse("You can ask me to turn on the hot tub, turn on the swim jets, or get the hot tub temperature.", false)
	default:
		return SpeakResponse("I don't know how to do that.", true)
	}
}

// handleStartSwimJet turns on the swim jets.
func (h *Handler) handleStartSwimJet() *Response {
	err := h.bridge.SetCircuit(gateway.CircuitSwimJets, 1)
	if err != nil {
		h.logger.Printf("Failed to start swim jet: %v", err)
		return SpeakResponse("Sorry, I couldn't start the swim jet.", true)
	}
	return SpeakResponse("Pool jet started", true)
}

// handleStopSwimJet turns off the swim jets.
func (h *Handler) handleStopSwimJet() *Response {
	err := h.bridge.SetCircuit(gateway.CircuitSwimJets, 0)
	if err != nil {
		h.logger.Printf("Failed to stop swim jet: %v", err)
		return SpeakResponse("Sorry, I couldn't stop the swim jet.", true)
	}
	return SpeakResponse("Pool jet stopped", true)
}

// handleStartHotTub turns on the spa.
func (h *Handler) handleStartHotTub() *Response {
	err := h.bridge.SetCircuit(gateway.CircuitSpa, 1)
	if err != nil {
		h.logger.Printf("Failed to start hot tub: %v", err)
		return SpeakResponse("Sorry, I couldn't start the hot tub.", true)
	}
	return SpeakResponse("Hot Tub started", true)
}

// handleStopHotTub turns off the spa.
func (h *Handler) handleStopHotTub() *Response {
	err := h.bridge.SetCircuit(gateway.CircuitSpa, 0)
	if err != nil {
		h.logger.Printf("Failed to stop hot tub: %v", err)
		return SpeakResponse("Sorry, I couldn't stop the hot tub.", true)
	}
	return SpeakResponse("Hot Tub stopped", true)
}

// handleHotTubTemp returns the current spa temperature.
func (h *Handler) handleHotTubTemp() *Response {
	// Refresh data
	h.bridge.Update()

	// Check if spa is on
	if !h.bridge.IsSpaOn() {
		return SpeakResponse("Hot tub is off", true)
	}

	// Get temperature
	temp, err := h.bridge.GetSpaTemperature()
	if err != nil {
		h.logger.Printf("Failed to get spa temperature: %v", err)
		return SpeakResponse("Sorry, I couldn't get the hot tub temperature.", true)
	}

	unit := h.bridge.TemperatureUnit()
	text := fmt.Sprintf("Hot Tub is %d %s", temp, unit)
	return SpeakResponse(text, true)
}
