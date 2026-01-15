package alexa

// Response is the Alexa skill response structure.
type Response struct {
	Version           string                 `json:"version"`
	SessionAttributes map[string]interface{} `json:"sessionAttributes,omitempty"`
	Response          ResponseBody           `json:"response"`
}

// ResponseBody contains the actual response content.
type ResponseBody struct {
	OutputSpeech     *OutputSpeech `json:"outputSpeech,omitempty"`
	Card             *Card         `json:"card,omitempty"`
	ShouldEndSession bool          `json:"shouldEndSession"`
}

// OutputSpeech defines speech output.
type OutputSpeech struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
	SSML string `json:"ssml,omitempty"`
}

// Card defines a visual card.
type Card struct {
	Type    string `json:"type"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
}

// SpeakResponse creates a simple speech response.
func SpeakResponse(text string, endSession bool) *Response {
	return &Response{
		Version: "1.0",
		Response: ResponseBody{
			OutputSpeech: &OutputSpeech{
				Type: "PlainText",
				Text: text,
			},
			Card: &Card{
				Type:    "Simple",
				Title:   text,
				Content: text,
			},
			ShouldEndSession: endSession,
		},
	}
}

// SpeakWithReprompt creates a speech response with a reprompt.
func SpeakWithReprompt(text string, reprompt string) *Response {
	return &Response{
		Version: "1.0",
		Response: ResponseBody{
			OutputSpeech: &OutputSpeech{
				Type: "PlainText",
				Text: text,
			},
			ShouldEndSession: false,
		},
	}
}
