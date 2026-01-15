package gateway

import "testing"

func TestConstants(t *testing.T) {
	// Verify message code pairs (query + 1 = answer)
	if ChallengeAnswer != ChallengeQuery+1 {
		t.Error("ChallengeAnswer should be ChallengeQuery + 1")
	}
	if LocalLoginAnswer != LocalLoginQuery+1 {
		t.Error("LocalLoginAnswer should be LocalLoginQuery + 1")
	}
	if VersionAnswer != VersionQuery+1 {
		t.Error("VersionAnswer should be VersionQuery + 1")
	}
	if PoolStatusAnswer != PoolStatusQuery+1 {
		t.Error("PoolStatusAnswer should be PoolStatusQuery + 1")
	}
	if ButtonPressAnswer != ButtonPressQuery+1 {
		t.Error("ButtonPressAnswer should be ButtonPressQuery + 1")
	}
	if CtrlConfigAnswer != CtrlConfigQuery+1 {
		t.Error("CtrlConfigAnswer should be CtrlConfigQuery + 1")
	}
}

func TestCircuitIDs(t *testing.T) {
	// Verify circuit IDs are in expected range
	circuits := []int{CircuitSpa, CircuitCleaner, CircuitSwimJets, CircuitPoolLight, CircuitSpaLight, CircuitPool}
	for _, id := range circuits {
		if id < 500 || id > 600 {
			t.Errorf("Circuit ID %d outside expected range 500-600", id)
		}
	}
}

func TestMappings(t *testing.T) {
	// Verify mappings have expected entries
	if len(BodyType) != 2 {
		t.Errorf("BodyType should have 2 entries, got %d", len(BodyType))
	}
	if BodyType[0] != "Pool" || BodyType[1] != "Spa" {
		t.Error("BodyType should be [Pool, Spa]")
	}

	if len(OnOff) < 2 {
		t.Errorf("OnOff should have at least 2 entries, got %d", len(OnOff))
	}
	if OnOff[0] != "Off" || OnOff[1] != "On" {
		t.Error("OnOff should start with [Off, On]")
	}
}
