package pool

import (
	"fmt"

	"github.com/nstielau/pool-controller/internal/gateway"
)

// Sensor represents a measurement sensor (temperature, chemistry, etc.).
type Sensor struct {
	id       string
	name     string
	state    interface{}
	unit     string
	hassType string
}

// NewSensor creates a new Sensor from gateway sensor data.
func NewSensor(id string, data *gateway.Sensor) *Sensor {
	return &Sensor{
		id:       id,
		name:     data.Name,
		state:    data.State,
		unit:     data.Unit,
		hassType: data.HassType,
	}
}

// NewBodySensor creates a sensor from body temperature data.
func NewBodySensor(id string, name string, value int, unit string) *Sensor {
	return &Sensor{
		id:       id,
		name:     name,
		state:    value,
		unit:     unit,
		hassType: "sensor",
	}
}

// NewChemistrySensor creates a sensor for chemistry data.
func NewChemistrySensor(id string, name string, value interface{}, unit string) *Sensor {
	return &Sensor{
		id:       id,
		name:     name,
		state:    value,
		unit:     unit,
		hassType: "sensor",
	}
}

// ID returns the sensor ID.
func (s *Sensor) ID() interface{} {
	return s.id
}

// Name returns the sensor name.
func (s *Sensor) Name() string {
	return s.name
}

// State returns the current state value.
func (s *Sensor) State() interface{} {
	return s.state
}

// Unit returns the measurement unit.
func (s *Sensor) Unit() string {
	return s.unit
}

// HassType returns the Home Assistant device type.
func (s *Sensor) HassType() string {
	return s.hassType
}

// FriendlyState returns a formatted state string with unit.
func (s *Sensor) FriendlyState() string {
	if s.hassType == "binary_sensor" {
		if val, ok := s.state.(int); ok {
			if val > 0 {
				return "On"
			}
			return "Off"
		}
	}

	if s.unit != "" {
		return fmt.Sprintf("%v %s", s.state, s.unit)
	}
	return fmt.Sprintf("%v", s.state)
}

// Update updates the sensor from new data.
func (s *Sensor) Update(data *gateway.Sensor) {
	s.state = data.State
}

// UpdateValue updates just the state value.
func (s *Sensor) UpdateValue(value interface{}) {
	s.state = value
}
