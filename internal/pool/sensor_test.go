package pool

import (
	"testing"

	"github.com/nstielau/pool-controller/internal/gateway"
)

func TestNewSensor(t *testing.T) {
	data := &gateway.Sensor{
		Name:     "Air Temperature",
		State:    72,
		Unit:     "°F",
		HassType: "sensor",
	}

	sensor := NewSensor("air_temperature", data)

	if sensor.ID() != "air_temperature" {
		t.Errorf("ID() = %v, want air_temperature", sensor.ID())
	}
	if sensor.Name() != "Air Temperature" {
		t.Errorf("Name() = %s, want Air Temperature", sensor.Name())
	}
	if sensor.State() != 72 {
		t.Errorf("State() = %v, want 72", sensor.State())
	}
	if sensor.Unit() != "°F" {
		t.Errorf("Unit() = %s, want °F", sensor.Unit())
	}
	if sensor.HassType() != "sensor" {
		t.Errorf("HassType() = %s, want sensor", sensor.HassType())
	}
}

func TestSensorFriendlyState(t *testing.T) {
	tests := []struct {
		name     string
		state    interface{}
		unit     string
		hassType string
		want     string
	}{
		{
			name:     "temperature",
			state:    72,
			unit:     "°F",
			hassType: "sensor",
			want:     "72 °F",
		},
		{
			name:     "no unit",
			state:    7.4,
			unit:     "",
			hassType: "sensor",
			want:     "7.4",
		},
		{
			name:     "binary on",
			state:    1,
			unit:     "",
			hassType: "binary_sensor",
			want:     "On",
		},
		{
			name:     "binary off",
			state:    0,
			unit:     "",
			hassType: "binary_sensor",
			want:     "Off",
		},
		{
			name:     "salt ppm",
			state:    3200,
			unit:     "ppm",
			hassType: "sensor",
			want:     "3200 ppm",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sensor := &Sensor{
				state:    tt.state,
				unit:     tt.unit,
				hassType: tt.hassType,
			}
			if got := sensor.FriendlyState(); got != tt.want {
				t.Errorf("FriendlyState() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestNewBodySensor(t *testing.T) {
	sensor := NewBodySensor("current_spa_temperature", "Current Spa Temperature", 102, "°F")

	if sensor.Name() != "Current Spa Temperature" {
		t.Errorf("Name() = %s, want Current Spa Temperature", sensor.Name())
	}
	if sensor.State() != 102 {
		t.Errorf("State() = %v, want 102", sensor.State())
	}
	if sensor.FriendlyState() != "102 °F" {
		t.Errorf("FriendlyState() = %s, want 102 °F", sensor.FriendlyState())
	}
}

func TestNewChemistrySensor(t *testing.T) {
	sensor := NewChemistrySensor("ph", "pH", 7.4, "")

	if sensor.Name() != "pH" {
		t.Errorf("Name() = %s, want pH", sensor.Name())
	}
	if sensor.FriendlyState() != "7.4" {
		t.Errorf("FriendlyState() = %s, want 7.4", sensor.FriendlyState())
	}
}

func TestSensorUpdateValue(t *testing.T) {
	sensor := NewBodySensor("temp", "Temperature", 70, "°F")

	sensor.UpdateValue(75)

	if sensor.State() != 75 {
		t.Errorf("State() = %v, want 75", sensor.State())
	}
	if sensor.FriendlyState() != "75 °F" {
		t.Errorf("FriendlyState() = %s, want 75 °F", sensor.FriendlyState())
	}
}
