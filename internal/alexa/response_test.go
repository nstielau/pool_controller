package alexa

import (
	"encoding/json"
	"testing"
)

func TestSpeakResponse(t *testing.T) {
	tests := []struct {
		name       string
		text       string
		endSession bool
	}{
		{
			name:       "end session",
			text:       "Pool jet started",
			endSession: true,
		},
		{
			name:       "keep session",
			text:       "What would you like to do?",
			endSession: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := SpeakResponse(tt.text, tt.endSession)

			if resp.Version != "1.0" {
				t.Errorf("Version = %s, want 1.0", resp.Version)
			}
			if resp.Response.OutputSpeech == nil {
				t.Fatal("OutputSpeech should not be nil")
			}
			if resp.Response.OutputSpeech.Type != "PlainText" {
				t.Errorf("OutputSpeech.Type = %s, want PlainText", resp.Response.OutputSpeech.Type)
			}
			if resp.Response.OutputSpeech.Text != tt.text {
				t.Errorf("OutputSpeech.Text = %s, want %s", resp.Response.OutputSpeech.Text, tt.text)
			}
			if resp.Response.ShouldEndSession != tt.endSession {
				t.Errorf("ShouldEndSession = %v, want %v", resp.Response.ShouldEndSession, tt.endSession)
			}
			if resp.Response.Card == nil {
				t.Fatal("Card should not be nil")
			}
			if resp.Response.Card.Type != "Simple" {
				t.Errorf("Card.Type = %s, want Simple", resp.Response.Card.Type)
			}
		})
	}
}

func TestSpeakResponseJSON(t *testing.T) {
	resp := SpeakResponse("Hot Tub is 102 °F", true)

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	// Verify it can be unmarshaled back
	var parsed Response
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if parsed.Response.OutputSpeech.Text != "Hot Tub is 102 °F" {
		t.Errorf("Round-trip text = %s, want Hot Tub is 102 °F", parsed.Response.OutputSpeech.Text)
	}
}

func TestSpeakWithReprompt(t *testing.T) {
	resp := SpeakWithReprompt("What would you like?", "Please say a command")

	if resp.Response.ShouldEndSession != false {
		t.Error("ShouldEndSession should be false for reprompt")
	}
	if resp.Response.OutputSpeech.Text != "What would you like?" {
		t.Errorf("Text = %s, want What would you like?", resp.Response.OutputSpeech.Text)
	}
}
