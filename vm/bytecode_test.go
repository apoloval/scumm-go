package vm_test

import (
	"bytes"
	"testing"

	"github.com/apoloval/scumm-go/vm"
	"github.com/stretchr/testify/assert"
)

func TestByteCodeReader_ReadPointer(t *testing.T) {
	for _, test := range []struct {
		name     string
		bytecode []byte
		expected vm.Pointer
	}{
		{"word pointer", []byte{0b01010101, 0b000_10101}, vm.WordPointer(0b1010101010101)},
		{"bit pointer", []byte{0b01010101, 0b1_1010101}, vm.BitPointer(0b1010101_01010101)},
		{"local pointer #1", []byte{0b0000_1010, 0b01_0_00000}, vm.LocalPointer(0b1010)},
		{"local pointer #2", []byte{0b0000_1010, 0b01_1_00000}, vm.LocalPointer(0b1010)},
	} {
		t.Run(test.name, func(t *testing.T) {
			r := vm.NewBytecodeReader(bytes.NewReader(test.bytecode))
			r.BeginFrame()
			assert.Equal(t, test.expected, r.ReadPointer())
			b, err := r.EndFrame()
			assert.NoError(t, err)
			assert.Equal(t, test.bytecode, b)
		})
	}
}

func TestByteCodeReader_ReadByteParam(t *testing.T) {
	for _, test := range []struct {
		name     string
		opcode   vm.OpCode
		pos      vm.ParamPos
		bytecode []byte
		expected vm.Param
	}{
		{"const pos 1", vm.OpCode(0x00), vm.ParamPos1, []byte{0x42}, vm.Const(0x42)},
		{"const pos 2", vm.OpCode(0x00), vm.ParamPos2, []byte{0x42}, vm.Const(0x42)},
		{"const pos 3", vm.OpCode(0x00), vm.ParamPos3, []byte{0x42}, vm.Const(0x42)},
		{"wptr pos 1", vm.OpCode(0x80), vm.ParamPos1, []byte{0x42, 0x0A}, vm.WordPointer(0xA42)},
		{"wptr pos 2", vm.OpCode(0x40), vm.ParamPos2, []byte{0x42, 0x0A}, vm.WordPointer(0xA42)},
		{"wptr pos 3", vm.OpCode(0x20), vm.ParamPos3, []byte{0x42, 0x0A}, vm.WordPointer(0xA42)},
		{"bptr pos 1", vm.OpCode(0x80), vm.ParamPos1, []byte{0xCD, 0xAB}, vm.BitPointer(0x2BCD)},
		{"bptr pos 2", vm.OpCode(0x40), vm.ParamPos2, []byte{0xCD, 0xAB}, vm.BitPointer(0x2BCD)},
		{"bptr pos 3", vm.OpCode(0x20), vm.ParamPos3, []byte{0xCD, 0xAB}, vm.BitPointer(0x2BCD)},
		{"lptr pos 1", vm.OpCode(0x80), vm.ParamPos1, []byte{0x0A, 0x40}, vm.LocalPointer(0xA)},
		{"lptr pos 2", vm.OpCode(0x40), vm.ParamPos2, []byte{0x0A, 0x40}, vm.LocalPointer(0xA)},
		{"lptr pos 3", vm.OpCode(0x20), vm.ParamPos3, []byte{0x0A, 0x40}, vm.LocalPointer(0xA)},
	} {
		t.Run(test.name, func(t *testing.T) {
			r := vm.NewBytecodeReader(bytes.NewReader(test.bytecode))
			r.BeginFrame()
			assert.Equal(t, test.expected, r.ReadByteParam(test.opcode, test.pos, vm.ParamFormatNumber))
			b, err := r.EndFrame()
			assert.NoError(t, err)
			assert.Equal(t, test.bytecode, b)
		})
	}
}

func TestByteCodeReader_ReadWordParam(t *testing.T) {
	for _, test := range []struct {
		name     string
		opcode   vm.OpCode
		pos      vm.ParamPos
		bytecode []byte
		expected vm.Param
	}{
		{"const pos 1", vm.OpCode(0x00), vm.ParamPos1, []byte{0x34, 0x12}, vm.Const(0x1234)},
		{"const pos 2", vm.OpCode(0x00), vm.ParamPos2, []byte{0x34, 0x12}, vm.Const(0x1234)},
		{"const pos 3", vm.OpCode(0x00), vm.ParamPos3, []byte{0x34, 0x12}, vm.Const(0x1234)},
		{"wptr pos 1", vm.OpCode(0x80), vm.ParamPos1, []byte{0x42, 0x0A}, vm.WordPointer(0xA42)},
		{"wptr pos 2", vm.OpCode(0x40), vm.ParamPos2, []byte{0x42, 0x0A}, vm.WordPointer(0xA42)},
		{"wptr pos 3", vm.OpCode(0x20), vm.ParamPos3, []byte{0x42, 0x0A}, vm.WordPointer(0xA42)},
		{"bptr pos 1", vm.OpCode(0x80), vm.ParamPos1, []byte{0xCD, 0xAB}, vm.BitPointer(0x2BCD)},
		{"bptr pos 2", vm.OpCode(0x40), vm.ParamPos2, []byte{0xCD, 0xAB}, vm.BitPointer(0x2BCD)},
		{"bptr pos 3", vm.OpCode(0x20), vm.ParamPos3, []byte{0xCD, 0xAB}, vm.BitPointer(0x2BCD)},
		{"lptr pos 1", vm.OpCode(0x80), vm.ParamPos1, []byte{0x0A, 0x40}, vm.LocalPointer(0xA)},
		{"lptr pos 2", vm.OpCode(0x40), vm.ParamPos2, []byte{0x0A, 0x40}, vm.LocalPointer(0xA)},
		{"lptr pos 3", vm.OpCode(0x20), vm.ParamPos3, []byte{0x0A, 0x40}, vm.LocalPointer(0xA)},
	} {
		t.Run(test.name, func(t *testing.T) {
			r := vm.NewBytecodeReader(bytes.NewReader(test.bytecode))
			r.BeginFrame()
			assert.Equal(t, test.expected, r.ReadWordParam(test.opcode, test.pos, vm.ParamFormatNumber))
			b, err := r.EndFrame()
			assert.NoError(t, err)
			assert.Equal(t, test.bytecode, b)
		})
	}
}
