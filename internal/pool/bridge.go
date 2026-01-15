package pool

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/nstielau/pool-controller/internal/gateway"
)

// Bridge is the main interface to the pool system.
type Bridge struct {
	mu             sync.RWMutex
	data           *gateway.PoolData
	devices        map[string]Device
	switches       map[int]*Switch
	conn           *gateway.Connection
	gatewayIP      string
	gatewayPort    int
	lastUpdate     time.Time
	updateInterval time.Duration
	timeout        time.Duration
}

// NewBridge creates a new Bridge, discovering the gateway if needed.
func NewBridge(gatewayIP string, gatewayPort int, updateInterval time.Duration) (*Bridge, error) {
	b := &Bridge{
		data:           gateway.NewPoolData(),
		devices:        make(map[string]Device),
		switches:       make(map[int]*Switch),
		updateInterval: updateInterval,
		timeout:        10 * time.Second,
	}

	// Discover gateway if not provided
	if gatewayIP == "" {
		info, err := gateway.DiscoverGateway(5 * time.Second)
		if err != nil {
			return nil, fmt.Errorf("gateway discovery failed: %w", err)
		}
		b.gatewayIP = info.IP
		b.gatewayPort = info.Port
	} else {
		b.gatewayIP = gatewayIP
		b.gatewayPort = gatewayPort
	}

	// Initial connection and data load
	err := b.loadInitialData()
	if err != nil {
		return nil, err
	}

	return b, nil
}

// loadInitialData connects to gateway and loads configuration and status.
func (b *Bridge) loadInitialData() error {
	conn := gateway.NewConnection(b.gatewayIP, b.gatewayPort)
	err := conn.Connect(b.timeout)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Query config first (needed for temperature unit)
	err = gateway.QueryConfig(conn, b.data, b.timeout)
	if err != nil {
		return fmt.Errorf("failed to query config: %w", err)
	}

	// Query status
	err = gateway.QueryStatus(conn, b.data, b.timeout)
	if err != nil {
		return fmt.Errorf("failed to query status: %w", err)
	}

	// Build device abstractions
	b.updateDevices()
	b.lastUpdate = time.Now()

	return nil
}

// updateDevices rebuilds the device map from raw data.
func (b *Bridge) updateDevices() {
	// Update switches from circuits
	for id, circuit := range b.data.Circuits {
		key := jsonName(circuit.Name)
		if sw, ok := b.switches[id]; ok {
			sw.Update(circuit)
		} else {
			sw := NewSwitch(circuit)
			b.switches[id] = sw
			b.devices[key] = sw
		}
	}

	// Update sensors
	for id, sensor := range b.data.Sensors {
		key := jsonName(id)
		if s, ok := b.devices[key].(*Sensor); ok {
			s.Update(sensor)
		} else {
			b.devices[key] = NewSensor(id, sensor)
		}
	}

	// Add body sensors
	unit := "°F"
	if b.data.Config.IsCelsius {
		unit = "°C"
	}

	for i, body := range b.data.Bodies {
		bodyName := "Pool"
		if body.BodyType == 1 {
			bodyName = "Spa"
		}

		// Current temperature
		tempKey := fmt.Sprintf("current_%s_temperature", strings.ToLower(bodyName))
		if s, ok := b.devices[tempKey].(*Sensor); ok {
			s.UpdateValue(body.CurrentTemperature)
		} else {
			b.devices[tempKey] = NewBodySensor(
				tempKey,
				fmt.Sprintf("Current %s Temperature", bodyName),
				body.CurrentTemperature,
				unit,
			)
		}

		// Heat status
		heatKey := fmt.Sprintf("%s_heater_%d", strings.ToLower(bodyName), i)
		if _, ok := b.devices[heatKey]; !ok {
			b.devices[heatKey] = &Sensor{
				id:       heatKey,
				name:     fmt.Sprintf("%s Heater", bodyName),
				state:    body.HeatStatus,
				hassType: "binary_sensor",
			}
		}
	}

	// Add chemistry sensors
	b.devices["ph"] = NewChemistrySensor("ph", "pH", b.data.Chemistry.PH, "")
	b.devices["orp"] = NewChemistrySensor("orp", "ORP", b.data.Chemistry.ORP, "")
	b.devices["saturation"] = NewChemistrySensor("saturation", "Saturation Index", b.data.Chemistry.Saturation, "")
	b.devices["salt_ppm"] = NewChemistrySensor("salt_ppm", "Salt", b.data.Chemistry.SaltPPM, "ppm")
}

// Update refreshes data from the gateway if the update interval has elapsed.
func (b *Bridge) Update() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if time.Since(b.lastUpdate) < b.updateInterval {
		return nil
	}

	conn := gateway.NewConnection(b.gatewayIP, b.gatewayPort)
	err := conn.Connect(b.timeout)
	if err != nil {
		return err
	}
	defer conn.Close()

	err = gateway.QueryStatus(conn, b.data, b.timeout)
	if err != nil {
		return err
	}

	b.updateDevices()
	b.lastUpdate = time.Now()

	return nil
}

// GetJSON returns all devices as a JSON string.
func (b *Bridge) GetJSON() (string, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	out := make(map[string]interface{})

	for key, dev := range b.devices {
		switch d := dev.(type) {
		case *Switch:
			out[key] = map[string]interface{}{
				"id":            d.IntID(),
				"name":          d.Name(),
				"friendlyState": strings.ToLower(d.FriendlyState()),
				"state":         d.IntState(),
			}
		case *Sensor:
			out[key] = map[string]interface{}{
				"name":  d.Name(),
				"state": d.FriendlyState(),
			}
		}
	}

	data, err := json.Marshal(out)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// GetCircuit returns the friendly state of a circuit.
func (b *Bridge) GetCircuit(circuitID int) string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if sw, ok := b.switches[circuitID]; ok {
		return sw.FriendlyState()
	}
	return "error"
}

// GetCircuitState returns the raw state of a circuit (0 or 1).
func (b *Bridge) GetCircuitState(circuitID int) int {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if sw, ok := b.switches[circuitID]; ok {
		return sw.IntState()
	}
	return -1
}

// SetCircuit changes a circuit's state.
func (b *Bridge) SetCircuit(circuitID, state int) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	conn := gateway.NewConnection(b.gatewayIP, b.gatewayPort)
	err := conn.Connect(b.timeout)
	if err != nil {
		return err
	}
	defer conn.Close()

	err = gateway.SetCircuit(conn, circuitID, state, b.timeout)
	if err != nil {
		return err
	}

	// Refresh status after change
	err = gateway.QueryStatus(conn, b.data, b.timeout)
	if err != nil {
		return err
	}

	b.updateDevices()
	b.lastUpdate = time.Now()

	return nil
}

// GetBodyTemperature returns the current temperature for a body (0=Pool, 1=Spa).
func (b *Bridge) GetBodyTemperature(bodyIndex int) (int, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if body, ok := b.data.Bodies[bodyIndex]; ok {
		return body.CurrentTemperature, nil
	}
	return 0, fmt.Errorf("body %d not found", bodyIndex)
}

// GetSpaTemperature returns the current spa temperature.
func (b *Bridge) GetSpaTemperature() (int, error) {
	return b.GetBodyTemperature(1)
}

// IsSpaOn returns true if the spa circuit is on.
func (b *Bridge) IsSpaOn() bool {
	return b.GetCircuitState(gateway.CircuitSpa) > 0
}

// GetDevice returns a device by its JSON key name.
func (b *Bridge) GetDevice(key string) (Device, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	dev, ok := b.devices[key]
	return dev, ok
}

// GetAttribute returns a specific attribute from the pool data.
func (b *Bridge) GetAttribute(attr string) (interface{}, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	dev, ok := b.devices[attr]
	if !ok {
		return nil, false
	}

	switch d := dev.(type) {
	case *Switch:
		return map[string]interface{}{
			"id":            d.IntID(),
			"name":          d.Name(),
			"friendlyState": strings.ToLower(d.FriendlyState()),
			"state":         d.IntState(),
		}, true
	case *Sensor:
		return map[string]interface{}{
			"name":  d.Name(),
			"state": d.FriendlyState(),
		}, true
	}

	return nil, false
}

// TemperatureUnit returns the temperature unit (°F or °C).
func (b *Bridge) TemperatureUnit() string {
	if b.data.Config.IsCelsius {
		return "°C"
	}
	return "°F"
}

// jsonName converts a name to JSON-friendly format (lowercase, underscores).
func jsonName(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "_"))
}
