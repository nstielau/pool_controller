package gateway

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

// Connection manages a TCP connection to the Pentair gateway.
type Connection struct {
	conn net.Conn
	ip   string
	port int
}

// NewConnection creates a new gateway connection.
func NewConnection(ip string, port int) *Connection {
	return &Connection{
		ip:   ip,
		port: port,
	}
}

// Connect establishes a connection to the gateway and performs login.
func (c *Connection) Connect(timeout time.Duration) error {
	// Establish TCP connection
	addr := fmt.Sprintf("%s:%d", c.ip, c.port)
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return fmt.Errorf("failed to connect to gateway: %w", err)
	}
	c.conn = conn

	// Set read/write deadline
	c.conn.SetDeadline(time.Now().Add(timeout))

	// Send connect string (no response expected)
	_, err = c.conn.Write([]byte(ConnectString))
	if err != nil {
		c.Close()
		return fmt.Errorf("failed to send connect string: %w", err)
	}

	// Challenge exchange
	err = c.sendChallenge()
	if err != nil {
		c.Close()
		return fmt.Errorf("challenge failed: %w", err)
	}

	// Login
	err = c.sendLogin()
	if err != nil {
		c.Close()
		return fmt.Errorf("login failed: %w", err)
	}

	// Clear deadline for future operations
	c.conn.SetDeadline(time.Time{})

	return nil
}

// sendChallenge performs the challenge exchange with the gateway.
func (c *Connection) sendChallenge() error {
	msg := MakeMessage(ChallengeQuery, nil)
	_, err := c.conn.Write(msg)
	if err != nil {
		return err
	}

	resp := make([]byte, 256)
	n, err := c.conn.Read(resp)
	if err != nil {
		return err
	}
	if n < HeaderSize {
		return fmt.Errorf("challenge response too short")
	}

	code, _, err := DecodeMessage(resp[:n])
	if err != nil {
		return err
	}
	if code != ChallengeAnswer {
		return fmt.Errorf("unexpected challenge response code: %d", code)
	}

	return nil
}

// sendLogin sends login credentials to the gateway.
func (c *Connection) sendLogin() error {
	// Build login message:
	// - uint32: schema (348)
	// - uint32: connectionType (0)
	// - string: clientVersion ("Android")
	// - string: password
	// - byte: padding (0)
	// - uint32: pid (2)

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, uint32(LoginSchema))
	binary.Write(buf, binary.LittleEndian, uint32(LoginConnectionType))
	buf.Write(MakeMessageString(LoginClientVersion))
	buf.Write(MakeMessageString(LoginPassword))
	buf.WriteByte(0) // padding
	binary.Write(buf, binary.LittleEndian, uint32(LoginPID))

	msg := MakeMessage(LocalLoginQuery, buf.Bytes())
	_, err := c.conn.Write(msg)
	if err != nil {
		return err
	}

	resp := make([]byte, 256)
	n, err := c.conn.Read(resp)
	if err != nil {
		return err
	}
	if n < HeaderSize {
		return fmt.Errorf("login response too short")
	}

	code, _, err := DecodeMessage(resp[:n])
	if err != nil {
		return err
	}
	if code != LocalLoginAnswer {
		return fmt.Errorf("unexpected login response code: %d", code)
	}

	return nil
}

// Close closes the gateway connection.
func (c *Connection) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Send sends a message and returns the response.
func (c *Connection) Send(msgCode uint16, data []byte, timeout time.Duration) ([]byte, error) {
	if c.conn == nil {
		return nil, fmt.Errorf("not connected")
	}

	c.conn.SetDeadline(time.Now().Add(timeout))
	defer c.conn.SetDeadline(time.Time{})

	msg := MakeMessage(msgCode, data)
	_, err := c.conn.Write(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	// Read response - may need multiple reads for large responses
	resp := make([]byte, 2048)
	n, err := c.conn.Read(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return resp[:n], nil
}

// IsConnected returns true if the connection is established.
func (c *Connection) IsConnected() bool {
	return c.conn != nil
}
