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

func (inst CursorShow) Mnemonic(*vm.SymbolTable) string     { return "CursorShow" }
func (inst CursorHide) Mnemonic(*vm.SymbolTable) string     { return "CursorHide" }
func (inst UserputOn) Mnemonic(*vm.SymbolTable) string      { return "UserputOn" }
func (inst UserputOff) Mnemonic(*vm.SymbolTable) string     { return "UserputOff" }
func (inst CursorSoftOn) Mnemonic(*vm.SymbolTable) string   { return "CursorSoftOn" }
func (inst CursorSoftOff) Mnemonic(*vm.SymbolTable) string  { return "CursorSoftOff" }
func (inst UserputSoftOn) Mnemonic(*vm.SymbolTable) string  { return "UserputSoftOn" }
func (inst UserputSoftOff) Mnemonic(*vm.SymbolTable) string { return "UserputSoftOff" }

func (inst SetCursorImg) Mnemonic(st *vm.SymbolTable) string {
	return fmt.Sprintf("SetCursorImg %s, %s",
		inst.Cursor.Display(st),
		inst.Char.Display(st),
	)
}
func (inst SetCursorHotspot) Mnemonic(st *vm.SymbolTable) string {
	return fmt.Sprintf("SetCursorHotspot %s, %s, %s",
		inst.Cursor.Display(st),
		inst.X.Display(st),
		inst.Y.Display(st),
	)
}
func (inst InitCursor) Mnemonic(st *vm.SymbolTable) string {
	return fmt.Sprintf("InitCursor %s", inst.Cursor.Display(st))
}
func (inst InitCharset) Mnemonic(st *vm.SymbolTable) string {
	return fmt.Sprintf("InitCharset %s",
		inst.Charset.Display(st),
	)
}

func (inst *SetCursorImg) Decode(opcode vm.OpCode, r *vm.BytecodeReader) error {
	inst.Cursor = r.ReadByteParam(opcode, vm.ParamPos1, vm.ParamFormatNumber)
	inst.Char = r.ReadByteParam(opcode, vm.ParamPos2, vm.ParamFormatChar)
	return inst.base.Decode(opcode, r)

}

func (inst *SetCursorHotspot) Decode(opcode vm.OpCode, r *vm.BytecodeReader) error {
	inst.Cursor = r.ReadByteParam(opcode, vm.ParamPos1, vm.ParamFormatNumber)
	inst.X = r.ReadByteParam(opcode, vm.ParamPos2, vm.ParamFormatNumber)
	inst.Y = r.ReadByteParam(opcode, vm.ParamPos3, vm.ParamFormatNumber)
	return inst.base.Decode(opcode, r)

}

func (inst *InitCursor) Decode(opcode vm.OpCode, r *vm.BytecodeReader) error {
	inst.Cursor = r.ReadByteParam(opcode, vm.ParamPos1, vm.ParamFormatNumber)
	return inst.base.Decode(opcode, r)

}

func (inst *InitCharset) Decode(opcode vm.OpCode, r *vm.BytecodeReader) error {
	inst.Charset = r.ReadByteParam(opcode, vm.ParamPos1, vm.ParamFormatCharsetID)
	return inst.base.Decode(opcode, r)

}

func decodeCursorCommand(opcode vm.OpCode, r *vm.BytecodeReader) (inst vm.Instruction, err error) {
	sub := r.ReadOpCode()

	switch sub & 0x1F {
	case 0x01:
		inst = &CursorShow{}
	case 0x02:
		inst = &CursorHide{}
	case 0x03:
		inst = &UserputOn{}
	case 0x04:
		inst = &UserputOff{}
	case 0x05:
		inst = &CursorSoftOn{}
	case 0x06:
		inst = &CursorSoftOff{}
	case 0x07:
		inst = &UserputSoftOn{}
	case 0x08:
		inst = &UserputSoftOff{}
	case 0x0A:
		inst = &SetCursorImg{}
	case 0x0B:
		inst = &SetCursorHotspot{}
	case 0x0C:
		inst = &InitCursor{}
	case 0x0D:
		inst = &InitCharset{}
	default:
		return nil, fmt.Errorf("unimplemented opcode %02X %02X for cursor command", opcode, sub)
	}

	err = inst.Decode(sub, r)
	return inst, err
}
