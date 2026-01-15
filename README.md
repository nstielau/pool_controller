# Pool Party!

A Go-based controller for Pentair ScreenLogic pool systems. Provides REST API and Alexa skill integration for controlling your pool and spa.

## Features

- **REST API** - Get pool status, control circuits via HTTP
- **Alexa Skill** - Voice control for spa, swim jets, and temperature queries
- **Auto-discovery** - Automatically finds your Pentair gateway on the network
- **Cross-platform** - Builds for Raspberry Pi, Linux, macOS
- **Simple deployment** - Single binary, systemd service included

## Prerequisites

- **Go 1.22+** (for building from source)
- **Pentair ScreenLogic** system on your local network
- **Raspberry Pi** (recommended) or any Linux server for deployment
- **SSH access** to your Pi for deployment

## Quick Start

```bash
# Show all available commands
make help

# Build for Raspberry Pi
make build-arm

# First-time setup on Pi
make setup-pi

# Deploy updates
make deploy
```

## API Endpoints

| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| `/` | GET | No | Health check |
| `/pool` | GET | Yes | Full pool status as JSON |
| `/pool/{attr}` | GET | Yes | Specific attribute |
| `/` | POST | Alexa | Alexa skill endpoint |

### Example Requests

```bash
# Health check
curl http://192.168.0.247/
# Response: hello

# Get full pool status
curl -H "Authorization: Bearer mytoken" http://192.168.0.247/pool
# Response: {"spa":{"id":500,"name":"Spa","friendlyState":"off","state":0},...}

# Get specific attribute
curl -H "Authorization: Bearer mytoken" http://192.168.0.247/pool/spa
# Response: {"id":500,"name":"Spa","friendlyState":"off","state":0}

# Get temperature
curl -H "Authorization: Bearer mytoken" http://192.168.0.247/pool/current_spa_temperature
# Response: {"name":"Current Spa Temperature","state":"102 °F"}
```

### Authentication

The `/pool` endpoints require a Bearer token validated against the `TOKEN_REGEX` environment variable.

| TOKEN_REGEX | Effect |
|-------------|--------|
| `.*` (default) | Accepts any token |
| `^mysecret$` | Requires exact match |
| `^pool-.*` | Requires prefix |

## Alexa Voice Commands

| Say | Action |
|-----|--------|
| "Alexa, turn on the hot tub" | Turns on spa circuit |
| "Alexa, turn off the hot tub" | Turns off spa circuit |
| "Alexa, turn on the swim jets" | Turns on swim jets |
| "Alexa, turn off the swim jets" | Turns off swim jets |
| "Alexa, what's the hot tub temperature?" | Reports spa temperature |

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `80` | HTTP server port |
| `GATEWAY_IP` | (auto-discover) | Pentair gateway IP (skip discovery) |
| `TOKEN_REGEX` | `.*` | Regex for token validation |
| `ALEXA_SKIP_VERIFY` | `false` | Skip Alexa signature verification (dev only) |

### Command Line Flags

```bash
pool-controller -port 8080 -gateway-ip 192.168.1.100 -update-interval 30s
```

### Circuit IDs

| ID | Name |
|----|------|
| 500 | Spa |
| 501 | Cleaner |
| 502 | Swim Jets |
| 503 | Pool Light |
| 504 | Spa Light |
| 505 | Pool |

## Deployment

### First-Time Setup

```bash
# Build and install on Pi (creates systemd service)
make setup-pi
```

This will:
1. Build ARM64 binary
2. Copy to `/opt/pool-controller/` on Pi
3. Install systemd service
4. Enable and start the service

### Updating

```bash
# Rebuild and deploy
make deploy
```

### Monitoring

```bash
# Check service status
make status

# Tail logs
make logs
```

## Development

### Project Structure

```
pool-controller/
├── cmd/pool-controller/     # Main entry point
├── internal/
│   ├── gateway/             # Pentair protocol (discovery, connection, queries)
│   ├── pool/                # Device abstractions (bridge, switch, sensor)
│   ├── api/                 # HTTP handlers and auth middleware
│   └── alexa/               # Alexa skill handlers and verification
├── Makefile                 # Build, test, deploy commands
├── pool-controller.service  # systemd unit file
└── .github/workflows/       # CI/CD pipeline
```

### Make Targets

```bash
make help      # Show all commands
make test      # Run tests
make coverage  # Run tests with coverage
make fmt       # Format code
make vet       # Run go vet
make lint      # Run all quality checks
```

### Running Locally

```bash
# Run with your gateway IP
GATEWAY_IP=192.168.0.225 ALEXA_SKIP_VERIFY=true go run ./cmd/pool-controller -port 8081

# Or use make (uses configured GATEWAY_IP)
make run
```

## Troubleshooting

### Gateway not found

If auto-discovery fails, set `GATEWAY_IP` explicitly:
```bash
GATEWAY_IP=192.168.1.100 ./pool-controller
```

Or in the systemd service file:
```ini
Environment=GATEWAY_IP=192.168.1.100
```

### Permission denied on deploy

Ensure SSH key authentication is set up:
```bash
ssh-copy-id pi@192.168.0.247
```

### Service won't start

Check logs for errors:
```bash
make logs
# or
ssh pi@192.168.0.247 "sudo journalctl -u pool-controller -n 50"
```

### API returns "Unauthed"

Ensure you're sending the Authorization header:
```bash
curl -H "Authorization: Bearer anytoken" http://192.168.0.247/pool
```

## Attribution

Protocol implementation based on:
- https://github.com/dieselrabbit/screenlogicpy_example
- https://github.com/keithpjolley/soipip

A million thanks!

## License

MIT - see [LICENSE](LICENSE)
