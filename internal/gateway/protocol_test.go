package gateway

import (
	"testing"
)

func TestMakeMessage(t *testing.T) {
	tests := []struct {
		name     string
		msgCode  uint16
		data     []byte
		wantLen  int
	}{
		{
			name:    "empty data",
			msgCode: VersionQuery,
			data:    nil,
			wantLen: HeaderSize,
		},
		{
			name:    "with data",
			msgCode: PoolStatusQuery,
			data:    []byte{0x00, 0x00, 0x00, 0x00},
			wantLen: HeaderSize + 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MakeMessage(tt.msgCode, tt.data)
			if len(msg) != tt.wantLen {
				t.Errorf("MakeMessage() len = %d, want %d", len(msg), tt.wantLen)
			}

			// Verify header structure
			// Bytes 0-1: MSG_CODE_1 (should be 0)
			if msg[0] != 0 || msg[1] != 0 {
				t.Errorf("MSG_CODE_1 should be 0, got %d %d", msg[0], msg[1])
			}

			// Bytes 2-3: MSG_CODE_2
			code := uint16(msg[2]) | uint16(msg[3])<<8
			if code != tt.msgCode {
				t.Errorf("MSG_CODE_2 = %d, want %d", code, tt.msgCode)
			}
		})
	}
}

func TestDecodeMessage(t *testing.T) {
	tests := []struct {
		name        string
		message     []byte
		wantCode    uint16
		wantDataLen int
		wantErr     bool
	}{
		{
			name:        "valid message",
			message:     []byte{0x00, 0x00, 0xb9, 0x1f, 0x04, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04},
			wantCode:    VersionAnswer,
			wantDataLen: 4,
			wantErr:     false,
		},
		{
			name:    "too short",
			message: []byte{0x00, 0x00, 0x00},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, data, err := DecodeMessage(tt.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if code != tt.wantCode {
					t.Errorf("DecodeMessage() code = %d, want %d", code, tt.wantCode)
				}
				if len(data) != tt.wantDataLen {
					t.Errorf("DecodeMessage() data len = %d, want %d", len(data), tt.wantDataLen)
				}
			}
		})
	}
}

func TestMakeMessageString(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		minLen int
	}{
		{
			name:   "short string",
			input:  "Android",
			minLen: 4 + 7, // length prefix + string
		},
		{
			name:   "empty string",
			input:  "",
			minLen: 4, // just length prefix
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MakeMessageString(tt.input)
			if len(result) < tt.minLen {
				t.Errorf("MakeMessageString() len = %d, want >= %d", len(result), tt.minLen)
			}

			// Should be padded to 4-byte boundary
			if len(result)%4 != 0 {
				t.Errorf("MakeMessageString() len = %d, not padded to 4 bytes", len(result))
			}
		})
	}
}

func TestGetMessageString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
	}{
		{name: "simple", input: "Android"},
		{name: "longer", input: "Pool Controller"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := MakeMessageString(tt.input)
			decoded, _ := GetMessageString(encoded)

			if decoded != tt.input {
				t.Errorf("GetMessageString() = %q, want %q", decoded, tt.input)
			}
		})
	}
}

func TestGetMessageStringEmpty(t *testing.T) {
	encoded := MakeMessageString("")
	decoded, consumed := GetMessageString(encoded)

	if decoded != "" {
		t.Errorf("GetMessageString() = %q, want empty", decoded)
	}
	// Empty string returns early after reading length prefix
	if consumed != 4 {
		t.Errorf("GetMessageString() consumed = %d, want 4 for empty string", consumed)
	}
}

func TestGetUint32(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}

	val, offset := GetUint32(data, 0)
	if val != 0x04030201 {
		t.Errorf("GetUint32() = 0x%x, want 0x04030201", val)
	}
	if offset != 4 {
		t.Errorf("GetUint32() offset = %d, want 4", offset)
	}

	val, offset = GetUint32(data, 4)
	if val != 0x08070605 {
		t.Errorf("GetUint32() = 0x%x, want 0x08070605", val)
	}
}

func TestGetByte(t *testing.T) {
	data := []byte{0xAB, 0xCD, 0xEF}

	val, offset := GetByte(data, 0)
	if val != 0xAB {
		t.Errorf("GetByte() = 0x%x, want 0xAB", val)
	}
	if offset != 1 {
		t.Errorf("GetByte() offset = %d, want 1", offset)
	}
}
