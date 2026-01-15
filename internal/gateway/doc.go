// Package gateway implements the Pentair ScreenLogic binary protocol.
//
// # Protocol Overview
//
// The Pentair ScreenLogic system uses a proprietary binary protocol over TCP.
// Communication follows a query-response pattern with 8-byte message headers.
//
// # Message Format
//
// Every message consists of an 8-byte header followed by optional data:
//
//	Offset  Size  Field       Description
//	0       2     MSG_CODE_1  Always 0
//	2       2     MSG_CODE_2  Message type identifier
//	4       4     Data Size   Length of message data
//	8       N     Data        Message payload (little-endian)
//
// # Discovery
//
// Gateways are discovered via UDP broadcast to 255.255.255.255:1444.
// The gateway responds with its IP address, port, and name.
//
// # Connection Flow
//
//  1. TCP connect to gateway IP:port
//  2. Send "CONNECTSERVERHOST\r\n\r\n" (no response)
//  3. Challenge exchange (gateway returns MAC address)
//  4. Login with credentials
//  5. Query config and status as needed
//
// # Usage
//
//	// Discover gateway
//	info, _ := gateway.DiscoverGateway(5 * time.Second)
//
//	// Connect
//	conn := gateway.NewConnection(info.IP, info.Port)
//	conn.Connect(10 * time.Second)
//	defer conn.Close()
//
//	// Query status
//	data := gateway.NewPoolData()
//	gateway.QueryConfig(conn, data, 10*time.Second)
//	gateway.QueryStatus(conn, data, 10*time.Second)
//
//	// Control circuit
//	gateway.SetCircuit(conn, gateway.CircuitSpa, 1, 10*time.Second)
package gateway
