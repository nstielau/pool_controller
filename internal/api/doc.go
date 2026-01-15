// Package api provides HTTP handlers for the pool controller REST API.
//
// # Endpoints
//
// The API provides the following endpoints:
//
//   - GET /        Health check, returns "hello"
//   - GET /pool    Returns full pool status as JSON (requires auth)
//   - GET /pool/{attr}  Returns specific attribute (requires auth)
//   - POST /       Alexa skill endpoint (uses Alexa verification)
//
// # Authentication
//
// The /pool endpoints use Bearer token authentication. Tokens are validated
// against the TOKEN_REGEX environment variable (default: ".*" accepts any token).
//
// Example request:
//
//	curl -H "Authorization: Bearer mytoken" http://localhost/pool
//
// # Usage
//
// Create a router with NewRouter and pass it to http.ListenAndServe:
//
//	bridge, _ := pool.NewBridge("", 0, 30*time.Second)
//	alexaHandler := alexa.NewHandler(bridge)
//	router := api.NewRouter(bridge, alexaHandler)
//	http.ListenAndServe(":80", router.Handler())
package api
