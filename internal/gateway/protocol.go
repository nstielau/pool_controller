package gateway

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// MakeMessage creates a complete protocol message with header.
// Header format (little-endian):
//   - 2 bytes: MSG_CODE_1 (always 0)
//   - 2 bytes: MSG_CODE_2 (message type)
//   - 4 bytes: Data size
//   - N bytes: Message data
func MakeMessage(msgCode2 uint16, data []byte) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, uint16(MsgCode1))
	binary.Write(buf, binary.LittleEndian, msgCode2)
	binary.Write(buf, binary.LittleEndian, uint32(len(data)))
	buf.Write(data)
	return buf.Bytes()
}

// DecodeMessage extracts message code and data from a raw message.
// Returns the message code (MSG_CODE_2) and the message data.
func DecodeMessage(message []byte) (uint16, []byte, error) {
	if len(message) < HeaderSize {
		return 0, nil, fmt.Errorf("message too short: %d bytes", len(message))
	}

	msgCode2 := binary.LittleEndian.Uint16(message[2:4])
	dataLen := binary.LittleEndian.Uint32(message[4:8])

	if msgCode2 == UnknownAnswer {
		return msgCode2, nil, fmt.Errorf("received UNKNOWN_ANSWER")
	}

	data := message[HeaderSize:]
	if uint32(len(data)) < dataLen {
		// Partial data is OK, just return what we have
	}

	return msgCode2, data, nil
}

// MakeMessageString encodes a string for the protocol.
// Format: 4-byte length prefix + string data + padding to 4-byte boundary
func MakeMessageString(s string) []byte {
	data := []byte(s)
	length := len(data)
	pad := 4 - (length % 4)
	if pad == 4 {
		pad = 0
	}
	// Always pad to 4-byte boundary (protocol requires this even for exact lengths)
	if pad == 0 {
		pad = 4
	}

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, uint32(length))
	buf.Write(data)
	buf.Write(make([]byte, pad))
	return buf.Bytes()
}

// GetMessageString decodes a length-prefixed string from the buffer.
// Returns the string and the number of bytes consumed.
func GetMessageString(data []byte) (string, int) {
	if len(data) < 4 {
		return "", 0
	}

	length := binary.LittleEndian.Uint32(data[0:4])
	if length == 0 {
		return "", 4
	}

	// Calculate padding
	pad := 4 - (length % 4)
	if pad == 4 {
		pad = 0
	}

	end := 4 + length
	if int(end) > len(data) {
		end = uint32(len(data))
	}

	str := string(data[4:end])
	consumed := 4 + int(length) + int(pad)

	return str, consumed
}

// GetUint32 reads a little-endian uint32 from buffer at offset.
func GetUint32(data []byte, offset int) (uint32, int) {
	if offset+4 > len(data) {
		return 0, offset
	}
	val := binary.LittleEndian.Uint32(data[offset : offset+4])
	return val, offset + 4
}

// GetInt32 reads a little-endian int32 from buffer at offset.
func GetInt32(data []byte, offset int) (int32, int) {
	if offset+4 > len(data) {
		return 0, offset
	}
	val := int32(binary.LittleEndian.Uint32(data[offset : offset+4]))
	return val, offset + 4
}

// GetUint16 reads a little-endian uint16 from buffer at offset.
func GetUint16(data []byte, offset int) (uint16, int) {
	if offset+2 > len(data) {
		return 0, offset
	}
	val := binary.LittleEndian.Uint16(data[offset : offset+2])
	return val, offset + 2
}

// GetByte reads a single byte from buffer at offset.
func GetByte(data []byte, offset int) (byte, int) {
	if offset >= len(data) {
		return 0, offset
	}
	return data[offset], offset + 1
}

// GetString reads a length-prefixed string from buffer at offset.
func GetString(data []byte, offset int) (string, int) {
	if offset+4 > len(data) {
		return "", offset
	}

	length := binary.LittleEndian.Uint32(data[offset : offset+4])
	offset += 4

	if length == 0 {
		return "", offset
	}

	// Calculate padded length
	paddedLen := length
	if length%4 != 0 {
		paddedLen += 4 - (length % 4)
	}

	end := offset + int(length)
	if end > len(data) {
		end = len(data)
	}

	str := string(data[offset:end])

	// Trim null bytes
	for len(str) > 0 && str[len(str)-1] == 0 {
		str = str[:len(str)-1]
	}

	return str, offset + int(paddedLen)
}
