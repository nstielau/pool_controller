package gateway

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

// GatewayInfo contains information about a discovered gateway.
type GatewayInfo struct {
	IP      string
	Port    int
	Type    byte
	Subtype byte
	Name    string
}

// DiscoverGateway broadcasts to find a Pentair gateway on the local network.
// Returns gateway information or an error if not found.
func DiscoverGateway(timeout time.Duration) (*GatewayInfo, error) {
	// Create UDP socket
	conn, err := net.ListenUDP("udp4", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UDP socket: %w", err)
	}
	defer conn.Close()

	// Set read timeout
	conn.SetReadDeadline(time.Now().Add(timeout))

	// Prepare broadcast address
	bcastAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", DiscoveryBroadcast, DiscoveryPort))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve broadcast address: %w", err)
	}

	// Send discovery packet: 8 bytes [1,0,0,0,0,0,0,0]
	packet := []byte{1, 0, 0, 0, 0, 0, 0, 0}
	_, err = conn.WriteToUDP(packet, bcastAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to send discovery broadcast: %w", err)
	}

	// Listen for response
	buf := make([]byte, 4096)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		return nil, fmt.Errorf("no gateway response: %w", err)
	}

	if n < 12 {
		return nil, fmt.Errorf("response too short: %d bytes", n)
	}

	// Parse response:
	// Offset 0-3: Checksum (should be 2)
	// Offset 4-7: IP address (4 bytes)
	// Offset 8-9: Port (uint16 LE)
	// Offset 10: Gateway type
	// Offset 11: Gateway subtype
	// Offset 12+: Gateway name (null-terminated string)

	checksum := binary.LittleEndian.Uint32(buf[0:4])
	if checksum != ExpectedChecksum {
		return nil, fmt.Errorf("invalid checksum: expected %d, got %d", ExpectedChecksum, checksum)
	}

	ip := fmt.Sprintf("%d.%d.%d.%d", buf[4], buf[5], buf[6], buf[7])
	port := int(binary.LittleEndian.Uint16(buf[8:10]))
	gwType := buf[10]
	gwSubtype := buf[11]

	// Extract name (null-terminated)
	name := ""
	if n > 12 {
		nameBytes := buf[12:n]
		for i, b := range nameBytes {
			if b == 0 {
				name = string(nameBytes[:i])
				break
			}
		}
		if name == "" {
			name = string(nameBytes)
		}
	}

	return &GatewayInfo{
		IP:      ip,
		Port:    port,
		Type:    gwType,
		Subtype: gwSubtype,
		Name:    name,
	}, nil
}
