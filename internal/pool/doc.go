// Package pool provides high-level abstractions for pool devices.
//
// # Overview
//
// This package provides a Bridge that manages communication with the Pentair
// gateway and exposes pool devices (switches and sensors) through a clean API.
//
// # Bridge
//
// The Bridge is the main entry point. It handles:
//   - Gateway discovery and connection
//   - Data caching with configurable update intervals
//   - Thread-safe access via sync.RWMutex
//   - Device state management
//
// # Devices
//
// Two device types are supported:
//
//   - Switch: Circuits that can be turned on/off (spa, jets, lights)
//   - Sensor: Read-only values (temperature, chemistry)
//
// Both implement the Device interface for uniform access.
//
// # Usage
//
//	// Create bridge (discovers gateway automatically)
//	bridge, err := pool.NewBridge("", 0, 30*time.Second)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Get pool status as JSON
//	json, _ := bridge.GetJSON()
//
//	// Control spa
//	bridge.SetCircuit(500, 1)  // Turn on
//	bridge.SetCircuit(500, 0)  // Turn off
//
//	// Query temperature
//	temp, _ := bridge.GetSpaTemperature()
//	fmt.Printf("Spa: %d°F\n", temp)
//
// # Circuit IDs
//
//	500 = Spa
//	501 = Cleaner
//	502 = Swim Jets
//	503 = Pool Light
//	504 = Spa Light
//	505 = Pool
package pool
