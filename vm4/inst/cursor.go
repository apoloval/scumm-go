package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

// CursorShow is a cursor command that shows the cursor.
type CursorShow struct{ base }

// CursorHide is a cursor command that hides the cursor.
type CursorHide struct{ base }

// UserputOn is a cursor command that enables user input.
type UserputOn struct{ base }

// UserputOff is a cursor command that disables user input.
type UserputOff struct{ base }

// CursorSoftOn is a cursor command that increments the cursor counter.
type CursorSoftOn struct{ base }

// CursorSoftOff is a cursor command that decrements the cursor counter.
type CursorSoftOff struct{ base }

// UserputSoftOn is a cursor command that increments the user input counter.
type UserputSoftOn struct{ base }

// UserputSoftOff is a cursor command that decrements the user input counter.
type UserputSoftOff struct{ base }

// SetCursorImg is a cursor command that sets the cursor image.
type SetCursorImg struct {
	base
	Cursor vm.Param
	Char   vm.Param
}

// SetCursorHotspot is a cursor command that sets the cursor hotspot.
type SetCursorHotspot struct {
	base
	Cursor vm.Param
	X      vm.Param
	Y      vm.Param
}

// InitCursor is a cursor command that initializes the cursor.
type InitCursor struct {
	base
	Cursor vm.Param
}

// InitCharset is a cursor command that initializes the charset.
type InitCharset struct {
	base
	Charset vm.Param
}

func (inst *SetCursorImg) Decode(opcode vm.OpCode, r *vm.BytecodeReader) error {
	inst.Cursor = r.ReadByteParam(opcode, vm.ParamPos1, vm.ParamFormatNumber)
	inst.Char = r.ReadByteParam(opcode, vm.ParamPos2, vm.ParamFormatChar)
	return inst.decodeWithParams(r, inst.Cursor, inst.Char)

}

func (inst *SetCursorHotspot) Decode(opcode vm.OpCode, r *vm.BytecodeReader) error {
	inst.Cursor = r.ReadByteParam(opcode, vm.ParamPos1, vm.ParamFormatNumber)
	inst.X = r.ReadByteParam(opcode, vm.ParamPos2, vm.ParamFormatNumber)
	inst.Y = r.ReadByteParam(opcode, vm.ParamPos3, vm.ParamFormatNumber)
	return inst.decodeWithParams(r, inst.Cursor, inst.X, inst.Y)
}

func (inst *InitCursor) Decode(opcode vm.OpCode, r *vm.BytecodeReader) error {
	inst.Cursor = r.ReadByteParam(opcode, vm.ParamPos1, vm.ParamFormatNumber)
	return inst.decodeWithParams(r, inst.Cursor)
}

func (inst *InitCharset) Decode(opcode vm.OpCode, r *vm.BytecodeReader) error {
	inst.Charset = r.ReadByteParam(opcode, vm.ParamPos1, vm.ParamFormatCharsetID)
	return inst.decodeWithParams(r, inst.Charset)
}

func decodeCursorCommand(opcode vm.OpCode, r *vm.BytecodeReader) (inst vm.Instruction, err error) {
	sub := r.ReadOpCode()

	switch sub & 0x1F {
	case 0x01:
		inst = &CursorShow{base: withName("CursorShow")}
	case 0x02:
		inst = &CursorHide{base: withName("CursorHide")}
	case 0x03:
		inst = &UserputOn{base: withName("UserputOn")}
	case 0x04:
		inst = &UserputOff{base: withName("UserputOff")}
	case 0x05:
		inst = &CursorSoftOn{base: withName("CursorSoftOn")}
	case 0x06:
		inst = &CursorSoftOff{base: withName("CursorSoftOff")}
	case 0x07:
		inst = &UserputSoftOn{base: withName("UserputSoftOn")}
	case 0x08:
		inst = &UserputSoftOff{base: withName("UserputSoftOff")}
	case 0x0A:
		inst = &SetCursorImg{base: withName("SetCursorImg")}
	case 0x0B:
		inst = &SetCursorHotspot{base: withName("SetCursorHotspot")}
	case 0x0C:
		inst = &InitCursor{base: withName("InitCursor")}
	case 0x0D:
		inst = &InitCharset{base: withName("InitCharset")}
	default:
		return nil, fmt.Errorf("unimplemented opcode %02X %02X for cursor command", opcode, sub)
	}

	err = inst.Decode(sub, r)
	return inst, err
}
