.PHONY: build build-arm deploy run test clean setup-pi

# Configuration
PI_HOST ?= pi@raspberrypi.local
PI_PATH ?= /opt/pool-controller
BINARY_NAME ?= pool-controller

# Build for current platform
build:
	go build -o $(BINARY_NAME) ./cmd/pool-controller

# Build for Raspberry Pi (ARM64)
build-arm:
	GOOS=linux GOARCH=arm64 go build -o $(BINARY_NAME)-arm64 ./cmd/pool-controller

# Deploy to Raspberry Pi
deploy: build-arm
	scp $(BINARY_NAME)-arm64 $(PI_HOST):$(PI_PATH)/$(BINARY_NAME)
	ssh $(PI_HOST) "sudo systemctl restart pool-controller"
	@echo "Deployed and restarted pool-controller on $(PI_HOST)"

# Run locally (for development)
run:
	ALEXA_SKIP_VERIFY=true GATEWAY_IP=192.168.0.225 go run ./cmd/pool-controller -port 8081

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME)-arm64

# Initial setup on Raspberry Pi (run once)
setup-pi: build-arm
	ssh $(PI_HOST) "sudo mkdir -p $(PI_PATH)"
	scp $(BINARY_NAME)-arm64 $(PI_HOST):$(PI_PATH)/$(BINARY_NAME)
	scp pool-controller.service $(PI_HOST):/tmp/
	ssh $(PI_HOST) "sudo mv /tmp/pool-controller.service /etc/systemd/system/ && \
		sudo systemctl daemon-reload && \
		sudo systemctl enable pool-controller && \
		sudo systemctl start pool-controller"
	@echo "Pool controller installed and started on $(PI_HOST)"

# View logs on Pi
logs:
	ssh $(PI_HOST) "sudo journalctl -u pool-controller -f"

# Check status on Pi
status:
	ssh $(PI_HOST) "sudo systemctl status pool-controller"
