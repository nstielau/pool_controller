package pool

import (
	"testing"

	"github.com/nstielau/pool-controller/internal/gateway"
)

func TestNewSwitch(t *testing.T) {
	circuit := &gateway.Circuit{
		ID:    500,
		Name:  "Spa",
		State: 1,
	}

	sw := NewSwitch(circuit)

	if sw.IntID() != 500 {
		t.Errorf("ID() = %d, want 500", sw.IntID())
	}
	if sw.Name() != "Spa" {
		t.Errorf("Name() = %s, want Spa", sw.Name())
	}
	if sw.IntState() != 1 {
		t.Errorf("IntState() = %d, want 1", sw.IntState())
	}
}

func TestSwitchFriendlyState(t *testing.T) {
	tests := []struct {
		name  string
		state int
		want  string
	}{
		{name: "off", state: 0, want: "Off"},
		{name: "on", state: 1, want: "On"},
		{name: "high value", state: 5, want: "On"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sw := &Switch{state: tt.state}
			if got := sw.FriendlyState(); got != tt.want {
				t.Errorf("FriendlyState() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestSwitchIsOn(t *testing.T) {
	tests := []struct {
		state int
		want  bool
	}{
		{state: 0, want: false},
		{state: 1, want: true},
		{state: 2, want: true},
	}

	for _, tt := range tests {
		sw := &Switch{state: tt.state}
		if got := sw.IsOn(); got != tt.want {
			t.Errorf("IsOn() with state %d = %v, want %v", tt.state, got, tt.want)
		}
	}
}

func TestSwitchHassType(t *testing.T) {
	sw := &Switch{}
	if got := sw.HassType(); got != "switch" {
		t.Errorf("HassType() = %s, want switch", got)
	}
}

func TestSwitchUpdate(t *testing.T) {
	sw := NewSwitch(&gateway.Circuit{ID: 500, Name: "Spa", State: 0})

	if sw.IsOn() {
		t.Error("Switch should be off initially")
	}

	sw.Update(&gateway.Circuit{ID: 500, Name: "Spa", State: 1})

	if !sw.IsOn() {
		t.Error("Switch should be on after update")
	}
}
