# Pool Party!

A Go-based controller for Pentair ScreenLogic pool systems. Provides REST API and Alexa skill integration for controlling your pool and spa.

## Features

- **REST API** - Get pool status, control circuits via HTTP
- **Alexa Skill** - Voice control for spa, swim jets, and temperature queries
- **Auto-discovery** - Automatically finds your Pentair gateway on the network
- **Cross-platform** - Builds for Raspberry Pi, Linux, macOS

## Quick Start

### Build

```bash
# Build for current platform
make build

# Build for Raspberry Pi (ARM64)
make build-arm

# Run locally (development)
make run
```

### Deploy to Raspberry Pi

First-time setup:
```bash
make setup-pi PI_HOST=pi@raspberrypi.local
```

Subsequent deploys:
```bash
make deploy PI_HOST=pi@raspberrypi.local
```

### View Logs

```bash
make logs PI_HOST=pi@raspberrypi.local
```

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | GET | Health check, returns "hello" |
| `/pool` | GET | Returns full pool status as JSON |
| `/pool/{attr}` | GET | Returns specific attribute |
| `/` | POST | Alexa skill endpoint |

### Authentication

The `/pool` endpoints require a Bearer token. Set `TOKEN_REGEX` environment variable to control token validation (default: `.*` accepts any token).

```bash
curl -H "Authorization: Bearer mytoken" http://localhost:8080/pool
```

## Alexa Intents

| Intent | Example Phrase | Action |
|--------|---------------|--------|
| `StartHotTubIntent` | "Turn on the hot tub" | Turns on spa circuit |
| `StopHotTubIntent` | "Turn off the hot tub" | Turns off spa circuit |
| `StartSwimJetIntent` | "Turn on the swim jets" | Turns on swim jets |
| `StopSwimJetIntent` | "Turn off the swim jets" | Turns off swim jets |
| `HotTubTempIntent` | "What's the hot tub temperature?" | Reports spa temperature |

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `80` | HTTP server port |
| `GATEWAY_IP` | (auto-discover) | Pentair gateway IP address |
| `TOKEN_REGEX` | `.*` | Regex pattern for token validation |
| `ALEXA_SKIP_VERIFY` | `false` | Skip Alexa signature verification (dev only) |

### Command Line Flags

```bash
pool-controller -port 8080 -gateway-ip 192.168.1.100 -update-interval 30s
```

## Circuit IDs

| ID | Name |
|----|------|
| 500 | Spa |
| 501 | Cleaner |
| 502 | Swim Jets |
| 503 | Pool Light |
| 504 | Spa Light |
| 505 | Pool |

## Development

### Project Structure

```
pool-controller/
├── cmd/pool-controller/     # Main entry point
├── internal/
│   ├── gateway/             # Pentair protocol implementation
│   ├── pool/                # Device abstractions and bridge
│   ├── api/                 # HTTP handlers
│   └── alexa/               # Alexa skill handlers
├── Makefile                 # Build/deploy commands
├── pool-controller.service  # systemd unit file
└── .github/workflows/       # CI/CD
```

### Running Tests

```bash
make test
```

### Local Development

Run with Alexa verification disabled:
```bash
ALEXA_SKIP_VERIFY=true go run ./cmd/pool-controller -port 8080
```

## Legacy Python Implementation

The original Python implementation files (`echoserver.py`, `pool_controller.py`, `gateway/`, `screenlogic/`) are kept for reference.

## Attribution

Pool controller via https://github.com/dieselrabbit/screenlogicpy_example and
http://github.com/keithpjolley/soipip. A million thanks!

## License

MIT
