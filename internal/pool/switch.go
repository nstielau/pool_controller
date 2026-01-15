package pool

import "github.com/nstielau/pool-controller/internal/gateway"

// Switch represents a toggleable circuit (on/off).
type Switch struct {
	id    int
	name  string
	state int
}

// NewSwitch creates a new Switch from circuit data.
func NewSwitch(circuit *gateway.Circuit) *Switch {
	return &Switch{
		id:    circuit.ID,
		name:  circuit.Name,
		state: circuit.State,
	}
}

// ID returns the circuit ID.
func (s *Switch) ID() interface{} {
	return s.id
}

// IntID returns the circuit ID as an int.
func (s *Switch) IntID() int {
	return s.id
}

// Name returns the circuit name.
func (s *Switch) Name() string {
	return s.name
}

// State returns the current state (0 or 1).
func (s *Switch) State() interface{} {
	return s.state
}

// IntState returns the state as an int.
func (s *Switch) IntState() int {
	return s.state
}

// HassType returns the Home Assistant device type.
func (s *Switch) HassType() string {
	return "switch"
}

// FriendlyState returns "On" or "Off".
func (s *Switch) FriendlyState() string {
	if s.state > 0 {
		return "On"
	}
	return "Off"
}

// IsOn returns true if the switch is on.
func (s *Switch) IsOn() bool {
	return s.state > 0
}

// Update updates the switch state from new circuit data.
func (s *Switch) Update(circuit *gateway.Circuit) {
	s.state = circuit.State
}
