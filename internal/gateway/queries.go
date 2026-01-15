package gateway

import (
	"encoding/binary"
	"fmt"
	"time"
)

// PoolData contains all pool information from the gateway.
type PoolData struct {
	Config    ConfigData
	Circuits  map[int]*Circuit
	Bodies    map[int]*Body
	Sensors   map[string]*Sensor
	Chemistry ChemistryData
}

// ConfigData contains pool configuration.
type ConfigData struct {
	ControllerID      uint32
	MinSetPoint       [2]int
	MaxSetPoint       [2]int
	IsCelsius         bool
	ControllerType    byte
	HardwareType      byte
	EquipmentFlags    int32
	CircuitCount      int
	Colors            []Color
	Pumps             map[int]byte
	InterfaceTabFlags uint32
	ShowAlarms        uint32
}

// Circuit represents a pool circuit (switch).
type Circuit struct {
	ID            int
	Name          string
	State         int
	Function      byte
	Interface     byte
	Flags         byte
	ColorSet      byte
	ColorPosition byte
	ColorStagger  byte
	DeviceID      byte
	DefaultRT     uint16
}

// Body represents a body of water (pool or spa).
type Body struct {
	BodyType           int
	CurrentTemperature int
	HeatStatus         int
	HeatSetPoint       int
	CoolSetPoint       int
	HeatMode           int
}

// Sensor represents a sensor reading.
type Sensor struct {
	Name     string
	State    interface{}
	Unit     string
	HassType string
}

// ChemistryData contains chemistry readings.
type ChemistryData struct {
	PH           float64
	ORP          int
	Saturation   float64
	SaltPPM      int
	PHTankLevel  int
	ORPTankLevel int
	Alarms       int
}

// Color represents a light color.
type Color struct {
	Name string
	R    uint32
	G    uint32
	B    uint32
}

// NewPoolData creates an empty PoolData structure.
func NewPoolData() *PoolData {
	return &PoolData{
		Circuits: make(map[int]*Circuit),
		Bodies:   make(map[int]*Body),
		Sensors:  make(map[string]*Sensor),
		Config: ConfigData{
			Pumps: make(map[int]byte),
		},
	}
}

// QueryVersion queries the gateway version.
func QueryVersion(conn *Connection, timeout time.Duration) (string, error) {
	resp, err := conn.Send(VersionQuery, nil, timeout)
	if err != nil {
		return "", err
	}

	code, data, err := DecodeMessage(resp)
	if err != nil {
		return "", err
	}
	if code != VersionAnswer {
		return "", fmt.Errorf("unexpected version response code: %d", code)
	}

	version, _ := GetMessageString(data)
	return version, nil
}

// QueryConfig queries the pool configuration.
func QueryConfig(conn *Connection, data *PoolData, timeout time.Duration) error {
	// Send config query with two zeros
	payload := make([]byte, 8)
	binary.LittleEndian.PutUint32(payload[0:4], 0)
	binary.LittleEndian.PutUint32(payload[4:8], 0)

	resp, err := conn.Send(CtrlConfigQuery, payload, timeout)
	if err != nil {
		return err
	}

	code, buf, err := DecodeMessage(resp)
	if err != nil {
		return err
	}
	if code != CtrlConfigAnswer {
		return fmt.Errorf("unexpected config response code: %d", code)
	}

	return decodeConfigAnswer(buf, data)
}

// QueryStatus queries the current pool status.
func QueryStatus(conn *Connection, data *PoolData, timeout time.Duration) error {
	// Send status query with one zero
	payload := make([]byte, 4)
	binary.LittleEndian.PutUint32(payload[0:4], 0)

	resp, err := conn.Send(PoolStatusQuery, payload, timeout)
	if err != nil {
		return err
	}

	code, buf, err := DecodeMessage(resp)
	if err != nil {
		return err
	}
	if code != PoolStatusAnswer {
		return fmt.Errorf("unexpected status response code: %d", code)
	}

	return decodeStatusAnswer(buf, data)
}

// SetCircuit sends a button press to change circuit state.
func SetCircuit(conn *Connection, circuitID, state int, timeout time.Duration) error {
	// Payload: padding (4 bytes), circuit ID (4 bytes), state (4 bytes)
	payload := make([]byte, 12)
	binary.LittleEndian.PutUint32(payload[0:4], 0)
	binary.LittleEndian.PutUint32(payload[4:8], uint32(circuitID))
	binary.LittleEndian.PutUint32(payload[8:12], uint32(state))

	resp, err := conn.Send(ButtonPressQuery, payload, timeout)
	if err != nil {
		return err
	}

	code, _, err := DecodeMessage(resp)
	if err != nil {
		return err
	}
	if code != ButtonPressAnswer {
		return fmt.Errorf("unexpected button press response code: %d", code)
	}

	return nil
}

// decodeConfigAnswer parses the configuration response.
func decodeConfigAnswer(buf []byte, data *PoolData) error {
	offset := 0

	// Controller ID
	data.Config.ControllerID, offset = GetUint32(buf, offset)

	// Min/Max set points
	var b byte
	b, offset = GetByte(buf, offset)
	data.Config.MinSetPoint[0] = int(b)
	b, offset = GetByte(buf, offset)
	data.Config.MaxSetPoint[0] = int(b)
	b, offset = GetByte(buf, offset)
	data.Config.MinSetPoint[1] = int(b)
	b, offset = GetByte(buf, offset)
	data.Config.MaxSetPoint[1] = int(b)

	// Is Celsius
	b, offset = GetByte(buf, offset)
	data.Config.IsCelsius = b != 0

	// Controller type, hardware type, buffer
	data.Config.ControllerType, offset = GetByte(buf, offset)
	data.Config.HardwareType, offset = GetByte(buf, offset)
	_, offset = GetByte(buf, offset) // controller buffer

	// Equipment flags
	data.Config.EquipmentFlags, offset = GetInt32(buf, offset)

	// Generic circuit name (skip)
	_, offset = GetString(buf, offset)

	// Circuit count
	var circuitCount uint32
	circuitCount, offset = GetUint32(buf, offset)
	data.Config.CircuitCount = int(circuitCount)

	// Circuits
	for i := 0; i < int(circuitCount); i++ {
		circuit := &Circuit{}

		var id int32
		id, offset = GetInt32(buf, offset)
		circuit.ID = int(id)

		circuit.Name, offset = GetString(buf, offset)

		circuit.Function, offset = GetByte(buf, offset) // nameIndex actually, reusing
		circuit.Function, offset = GetByte(buf, offset)
		circuit.Interface, offset = GetByte(buf, offset)
		circuit.Flags, offset = GetByte(buf, offset)
		circuit.ColorSet, offset = GetByte(buf, offset)
		circuit.ColorPosition, offset = GetByte(buf, offset)
		circuit.ColorStagger, offset = GetByte(buf, offset)
		circuit.DeviceID, offset = GetByte(buf, offset)
		circuit.DefaultRT, offset = GetUint16(buf, offset)

		// Skip 2 padding bytes
		offset += 2

		data.Circuits[circuit.ID] = circuit
	}

	// Color count
	var colorCount uint32
	colorCount, offset = GetUint32(buf, offset)

	data.Config.Colors = make([]Color, colorCount)
	for i := 0; i < int(colorCount); i++ {
		data.Config.Colors[i].Name, offset = GetString(buf, offset)
		data.Config.Colors[i].R, offset = GetUint32(buf, offset)
		data.Config.Colors[i].G, offset = GetUint32(buf, offset)
		data.Config.Colors[i].B, offset = GetUint32(buf, offset)
	}

	// Pump data (8 entries)
	for i := 0; i < 8; i++ {
		data.Config.Pumps[i], offset = GetByte(buf, offset)
	}

	// Interface tab flags and show alarms
	data.Config.InterfaceTabFlags, offset = GetUint32(buf, offset)
	data.Config.ShowAlarms, offset = GetUint32(buf, offset)

	return nil
}

// decodeStatusAnswer parses the status response.
func decodeStatusAnswer(buf []byte, data *PoolData) error {
	offset := 0

	// OK flag
	_, offset = GetUint32(buf, offset)

	// Freeze mode, remotes, delays
	_, offset = GetByte(buf, offset) // freezeMode
	_, offset = GetByte(buf, offset) // remotes
	_, offset = GetByte(buf, offset) // poolDelay
	_, offset = GetByte(buf, offset) // spaDelay
	_, offset = GetByte(buf, offset) // cleanerDelay
	_, offset = GetByte(buf, offset) // ff1
	_, offset = GetByte(buf, offset) // ff2
	_, offset = GetByte(buf, offset) // ff3

	// Temperature unit
	unit := "°F"
	if data.Config.IsCelsius {
		unit = "°C"
	}

	// Air temperature
	var airTemp int32
	airTemp, offset = GetInt32(buf, offset)
	data.Sensors["air_temperature"] = &Sensor{
		Name:     "Air Temperature",
		State:    int(airTemp),
		Unit:     unit,
		HassType: "sensor",
	}

	// Bodies count
	var bodiesCount uint32
	bodiesCount, offset = GetUint32(buf, offset)
	if bodiesCount > 2 {
		bodiesCount = 2
	}

	for i := 0; i < int(bodiesCount); i++ {
		body := &Body{}

		var bodyType uint32
		bodyType, offset = GetUint32(buf, offset)
		if bodyType > 1 {
			bodyType = 0
		}
		body.BodyType = int(bodyType)

		var temp int32
		temp, offset = GetInt32(buf, offset)
		body.CurrentTemperature = int(temp)

		temp, offset = GetInt32(buf, offset)
		body.HeatStatus = int(temp)

		temp, offset = GetInt32(buf, offset)
		body.HeatSetPoint = int(temp)

		temp, offset = GetInt32(buf, offset)
		body.CoolSetPoint = int(temp)

		temp, offset = GetInt32(buf, offset)
		body.HeatMode = int(temp)

		data.Bodies[i] = body

		// Also add as sensors for compatibility
		bodyName := "Pool"
		if body.BodyType == 1 {
			bodyName = "Spa"
		}

		data.Sensors[fmt.Sprintf("current_%s_temperature", bodyName)] = &Sensor{
			Name:     fmt.Sprintf("Current %s Temperature", bodyName),
			State:    body.CurrentTemperature,
			Unit:     unit,
			HassType: "sensor",
		}
	}

	// Circuit count and states
	var circuitCount uint32
	circuitCount, offset = GetUint32(buf, offset)

	for i := 0; i < int(circuitCount); i++ {
		var circuitID, circuitState uint32
		circuitID, offset = GetUint32(buf, offset)
		circuitState, offset = GetUint32(buf, offset)

		if circuit, ok := data.Circuits[int(circuitID)]; ok {
			circuit.State = int(circuitState)
		}

		// Skip color bytes
		_, offset = GetByte(buf, offset) // colorSet
		_, offset = GetByte(buf, offset) // colorPos
		_, offset = GetByte(buf, offset) // colorStagger
		_, offset = GetByte(buf, offset) // delay
	}

	// Chemistry data
	var val int32
	val, offset = GetInt32(buf, offset)
	data.Chemistry.PH = float64(val) / 100.0

	val, offset = GetInt32(buf, offset)
	data.Chemistry.ORP = int(val)

	val, offset = GetInt32(buf, offset)
	data.Chemistry.Saturation = float64(val) / 100.0

	val, offset = GetInt32(buf, offset)
	data.Chemistry.SaltPPM = int(val)

	val, offset = GetInt32(buf, offset)
	data.Chemistry.PHTankLevel = int(val)

	val, offset = GetInt32(buf, offset)
	data.Chemistry.ORPTankLevel = int(val)

	val, offset = GetInt32(buf, offset)
	data.Chemistry.Alarms = int(val)

	return nil
}
