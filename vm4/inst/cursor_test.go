package inst_test

import (
	"bytes"
	"testing"

	"github.com/apoloval/scumm-go/vm"
	"github.com/apoloval/scumm-go/vm4/inst"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeInstV4(t *testing.T) {
	for _, testCase := range []struct {
		bytecode []byte
		expected string
	}{
		{
			bytecode: []byte{0x2C, 0x01},
			expected: "CursorShow",
		},
		{
			bytecode: []byte{0x2C, 0x02},
			expected: "CursorHide",
		},
		{
			bytecode: []byte{0x2C, 0x03},
			expected: "UserputOn",
		},
		{
			bytecode: []byte{0x2C, 0x04},
			expected: "UserputOff",
		},
		{
			bytecode: []byte{0x2C, 0x05},
			expected: "CursorSoftOn",
		},
		{
			bytecode: []byte{0x2C, 0x06},
			expected: "CursorSoftOff",
		},
		{
			bytecode: []byte{0x2C, 0x07},
			expected: "UserputSoftOn",
		},
		{
			bytecode: []byte{0x2C, 0x08},
			expected: "UserputSoftOff",
		},
		{
			bytecode: []byte{0x2C, 0x0A, 0x01, 0x41},
			expected: "SetCursorImg 1, 'A'",
		},
		{
			bytecode: []byte{0x2C, 0x0B, 0x01, 0x02, 0x03},
			expected: "SetCursorHotspot 1, 2, 3",
		},
		{
			bytecode: []byte{0x2C, 0x0C, 0x01},
			expected: "InitCursor 1",
		},
		{
			bytecode: []byte{0x2C, 0x0D, 0x01},
			expected: "InitCharset CHARSET_0001",
		},
	} {
		t.Run(testCase.expected, func(t *testing.T) {
			st := vm.NewSymbolTable()
			r := vm.NewBytecodeDecoder(bytes.NewReader(testCase.bytecode))
			inst, err := inst.Decode(r)
			require.NoError(t, err)
			assert.Equal(t, testCase.expected, vm.DisplayInstruction(st, inst))
		})
	}
}
