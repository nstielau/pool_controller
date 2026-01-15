# Pool Controller - Development Notes

## Project Overview

Go-based controller for Pentair ScreenLogic pool systems. Runs on Raspberry Pi, provides REST API and Alexa skill integration.

## Quick Commands

```bash
# Build and test
make build          # Build for current platform
make build-arm      # Cross-compile for Pi (ARM64)
make test           # Run all tests
make lint           # Format and vet code

# Local development
make run            # Run locally (requires GATEWAY_IP in Makefile.local)

# Deployment
make deploy         # Build and deploy to Pi
make status         # Check service status on Pi
make logs           # Tail logs from Pi
```

## Local Configuration

Create `Makefile.local` from `Makefile.local.example` with your local settings:

```makefile
PI_HOST = pi@192.168.0.247      # Your Pi's IP
GATEWAY_IP = 192.168.0.225      # Pentair gateway IP
DEV_PORT = 8081
```

This file is gitignored to keep personal IPs out of version control.

## Architecture

```
cmd/pool-controller/   # Entry point
internal/
  gateway/             # Pentair protocol (TCP/UDP, binary)
  pool/                # Device abstractions (bridge, switches, sensors)
  api/                 # HTTP handlers (/pool endpoints)
  alexa/               # Alexa skill (signature verification, intents)
```

## Pentair Protocol Notes

- Discovery: UDP broadcast to 255.255.255.255:1444
- Connection: TCP with "CONNECTSERVERHOST\r\n\r\n" handshake
- Messages: 8-byte little-endian header (2 bytes padding, 2 bytes msg code, 4 bytes data length)
- Strings: 4-byte length prefix, padded to 4-byte boundary

## Circuit IDs

```go
CircuitSpa       = 500
CircuitCleaner   = 501
CircuitSwimJets  = 502
CircuitPoolLight = 503
CircuitSpaLight  = 504
CircuitPool      = 505
```

## Environment Variables

- `GATEWAY_IP` - Pentair gateway IP address (required)
- `TOKEN_REGEX` - Regex for API authentication (default: `.*` allows all)
- `ALEXA_SKIP_VERIFY` - Set to `true` to skip Alexa signature verification (dev only)

## Testing

```bash
# Run all tests
make test

# Run with coverage
make coverage

# Test specific package
go test -v ./internal/gateway/...
```

## Deployment

The service runs via systemd on the Pi:

```bash
# First-time setup
make setup-pi

# Subsequent deploys
make deploy

# Check status/logs
make status
make logs
```

## Alexa Intents

- `StartHotTubIntent` / `StopHotTubIntent` - Control spa (circuit 500)
- `StartSwimJetIntent` / `StopSwimJetIntent` - Control jets (circuit 502)
- `HotTubTempIntent` - Get spa temperature
