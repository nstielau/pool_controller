.PHONY: build build-arm deploy run test clean setup-pi logs status help fmt vet lint coverage

# Default configuration (override in Makefile.local)
PI_HOST ?= pi@raspberrypi.local
PI_PATH ?= /opt/pool-controller
BINARY_NAME ?= pool-controller
GATEWAY_IP ?=
DEV_PORT ?= 8081

# Include local overrides if they exist (not committed to git)
-include Makefile.local

# Default target
.DEFAULT_GOAL := help

## help: Show this help message
help:
	@echo "Pool Controller - Pentair ScreenLogic Control Server"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Build:"
	@echo "  build       Build binary for current platform"
	@echo "  build-arm   Build binary for Raspberry Pi (ARM64)"
	@echo "  clean       Remove build artifacts"
	@echo ""
	@echo "Development:"
	@echo "  run         Run locally (dev mode, port $(DEV_PORT))"
	@echo "  test        Run all tests"
	@echo "  coverage    Run tests with coverage report"
	@echo "  fmt         Format code with gofmt"
	@echo "  vet         Run go vet"
	@echo "  lint        Run all code quality checks"
	@echo ""
	@echo "Deployment (PI_HOST=$(PI_HOST)):"
	@echo "  setup-pi    First-time Pi setup (installs systemd service)"
	@echo "  deploy      Build and deploy to Pi"
	@echo "  status      Check service status on Pi"
	@echo "  logs        Tail logs from Pi"
	@echo ""
	@echo "Configuration (set in Makefile.local):"
	@echo "  PI_HOST     Raspberry Pi SSH target (current: $(PI_HOST))"
	@echo "  GATEWAY_IP  Pool gateway IP for local dev (current: $(GATEWAY_IP))"
	@echo ""
	@echo "Create Makefile.local from Makefile.local.example for local settings."

## build: Build for current platform
build:
	go build -o $(BINARY_NAME) ./cmd/pool-controller

## build-arm: Build for Raspberry Pi (ARM64)
build-arm:
	GOOS=linux GOARCH=arm64 go build -o $(BINARY_NAME)-arm64 ./cmd/pool-controller

## deploy: Deploy to Raspberry Pi
deploy: build-arm
	scp $(BINARY_NAME)-arm64 $(PI_HOST):/tmp/$(BINARY_NAME)
	ssh $(PI_HOST) "sudo mv /tmp/$(BINARY_NAME) $(PI_PATH)/$(BINARY_NAME) && sudo systemctl restart pool-controller"
	@echo "Deployed and restarted pool-controller on $(PI_HOST)"

## run: Run locally for development
run:
ifndef GATEWAY_IP
	$(error GATEWAY_IP is not set. Create Makefile.local or run: GATEWAY_IP=x.x.x.x make run)
endif
	ALEXA_SKIP_VERIFY=true GATEWAY_IP=$(GATEWAY_IP) go run ./cmd/pool-controller -port $(DEV_PORT)

## test: Run all tests
test:
	go test -v ./...

## coverage: Run tests with coverage report
coverage:
	go test -cover ./...
	@echo ""
	@echo "For detailed HTML report: go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out"

## fmt: Format code
fmt:
	go fmt ./...

## vet: Run go vet
vet:
	go vet ./...

## lint: Run all code quality checks
lint: fmt vet
	@echo "Code quality checks passed"

## clean: Remove build artifacts
clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME)-arm64 coverage.out

## setup-pi: Initial setup on Raspberry Pi (run once)
setup-pi: build-arm
	ssh $(PI_HOST) "sudo mkdir -p $(PI_PATH)"
	scp $(BINARY_NAME)-arm64 $(PI_HOST):/tmp/$(BINARY_NAME)
	scp pool-controller.service $(PI_HOST):/tmp/
	ssh $(PI_HOST) "sudo mv /tmp/$(BINARY_NAME) $(PI_PATH)/$(BINARY_NAME) && \
		sudo mv /tmp/pool-controller.service /etc/systemd/system/ && \
		sudo systemctl daemon-reload && \
		sudo systemctl enable pool-controller && \
		sudo systemctl start pool-controller"
	@echo "Pool controller installed and started on $(PI_HOST)"

## logs: View logs on Pi
logs:
	ssh $(PI_HOST) "sudo journalctl -u pool-controller -f"

## status: Check service status on Pi
status:
	ssh $(PI_HOST) "sudo systemctl status pool-controller"
