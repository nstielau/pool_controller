// Package gateway implements the Pentair ScreenLogic protocol.
// Protocol documentation: https://github.com/ceisenach/screenlogic_over_ip
package gateway

// Message codes for the Pentair protocol
const (
	MsgCode1 = 0

	ChallengeQuery    = 14
	ChallengeAnswer   = 15
	LocalLoginQuery   = 27
	LocalLoginAnswer  = 28
	VersionQuery      = 8120
	VersionAnswer     = 8121
	PoolStatusQuery   = 12526
	PoolStatusAnswer  = 12527
	ButtonPressQuery  = 12530
	ButtonPressAnswer = 12531
	CtrlConfigQuery   = 12532
	CtrlConfigAnswer  = 12533
	UnknownAnswer     = 13
)

// Circuit IDs
const (
	CircuitSpa       = 500
	CircuitCleaner   = 501
	CircuitSwimJets  = 502
	CircuitPoolLight = 503
	CircuitSpaLight  = 504
	CircuitPool      = 505
	CircuitAux5      = 506
	CircuitAux6      = 507
	CircuitAux7      = 508
)

// State mappings
var BodyType = []string{"Pool", "Spa"}
var HeatMode = []string{"Off", "Solar", "Solar Preferred", "Heat", "Don't Change"}
var OnOff = []string{"Off", "On", "Unknown"}
var ColorMode = []string{
	"Off", "On", "Set", "Sync", "Swim", "Party",
	"Romantic", "Caribbean", "American", "Sunset",
	"Royal", "Save", "Recall", "Blue", "Green",
	"Red", "White", "Magenta", "Thumper", "Next",
	"Reset", "Hold",
}

// HeaderSize is the size of the message header (8 bytes)
const HeaderSize = 8

// Discovery constants
const (
	DiscoveryPort      = 1444
	DiscoveryBroadcast = "255.255.255.255"
	ExpectedChecksum   = 2
)

// Login constants
const (
	LoginSchema         = 348
	LoginConnectionType = 0
	LoginPID            = 2
	LoginClientVersion  = "Android"
	LoginPassword       = "mypassword"
)

// ConnectString is sent to initiate connection with gateway
const ConnectString = "CONNECTSERVERHOST\r\n\r\n"
