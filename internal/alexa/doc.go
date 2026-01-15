// Package alexa implements Alexa skill request handling and verification.
//
// # Overview
//
// This package provides handlers for Alexa Smart Home and Custom Skill requests,
// including signature verification as required by Amazon's security guidelines.
//
// # Supported Intents
//
//   - LaunchRequest         Skill invocation ("Alexa, open pool party")
//   - StartSwimJetIntent    Turn on swim jets
//   - StopSwimJetIntent     Turn off swim jets
//   - StartHotTubIntent     Turn on spa/hot tub
//   - StopHotTubIntent      Turn off spa/hot tub
//   - HotTubTempIntent      Query spa temperature
//   - AMAZON.CancelIntent   Cancel/stop skill
//   - AMAZON.StopIntent     Stop skill
//
// # Security
//
// All incoming requests are verified using Amazon's signature verification:
//   - Certificate URL validation (must be from s3.amazonaws.com/echo.api/)
//   - Certificate download and caching
//   - RSA signature verification
//   - Timestamp validation (within 150 seconds)
//
// Set ALEXA_SKIP_VERIFY=true to disable verification during development.
//
// # Usage
//
//	bridge, _ := pool.NewBridge("", 0, 30*time.Second)
//	handler := alexa.NewHandler(bridge)
//	http.Handle("/", handler)
package alexa
