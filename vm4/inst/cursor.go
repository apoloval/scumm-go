package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

// CursorShow is a cursor command that shows the cursor.
type CursorShow struct{ instruction }

// CursorHide is a cursor command that hides the cursor.
type CursorHide struct{ instruction }

// UserputOn is a cursor command that enables user input.
type UserputOn struct{ instruction }

// UserputOff is a cursor command that disables user input.
type UserputOff struct{ instruction }

// CursorSoftOn is a cursor command that increments the cursor counter.
type CursorSoftOn struct{ instruction }

// CursorSoftOff is a cursor command that decrements the cursor counter.
type CursorSoftOff struct{ instruction }

// UserputSoftOn is a cursor command that increments the user input counter.
type UserputSoftOn struct{ instruction }

// UserputSoftOff is a cursor command that decrements the user input counter.
type UserputSoftOff struct{ instruction }

// SetCursorImg is a cursor command that sets the cursor image.
type SetCursorImg struct {
	instruction
	Cursor vm.Param
	Char   vm.Param
}

// SetCursorHotspot is a cursor command that sets the cursor hotspot.
type SetCursorHotspot struct {
	instruction
	Cursor vm.Param
	X      vm.Param
	Y      vm.Param
}

// InitCursor is a cursor command that initializes the cursor.
type InitCursor struct {
	instruction
	Cursor vm.Param
}

// InitCharset is a cursor command that initializes the charset.
type InitCharset struct {
	instruction
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
func (inst SetCursorImg) Mnemonic(*vm.SymbolTable) string {
	return fmt.Sprintf("SetCursorImg %s, %s", inst.Cursor, inst.Char)
}
func (inst SetCursorHotspot) Mnemonic(*vm.SymbolTable) string {
	return fmt.Sprintf("SetCursorHotspot %s, %s, %s", inst.Cursor, inst.X, inst.Y)
}
func (inst InitCursor) Mnemonic(*vm.SymbolTable) string {
	return fmt.Sprintf("InitCursor %s", inst.Cursor)
}
func (inst InitCharset) Mnemonic(st *vm.SymbolTable) string {
	return fmt.Sprintf("InitCharset %s",
		inst.Charset.Display(st, vm.ParamFormatCharsetID))
}

func (inst *CursorShow) Decode(opcode vm.OpCode, r *vm.BytecodeReader) (err error) {
	inst.frame, err = r.EndFrame()
	return
}

func (inst *CursorHide) Decode(opcode vm.OpCode, r *vm.BytecodeReader) (err error) {
	inst.frame, err = r.EndFrame()
	return
}

func (inst *UserputOn) Decode(opcode vm.OpCode, r *vm.BytecodeReader) (err error) {
	inst.frame, err = r.EndFrame()
	return
}

func (inst *UserputOff) Decode(opcode vm.OpCode, r *vm.BytecodeReader) (err error) {
	inst.frame, err = r.EndFrame()
	return
}

func (inst *CursorSoftOn) Decode(opcode vm.OpCode, r *vm.BytecodeReader) (err error) {
	inst.frame, err = r.EndFrame()
	return
}

func (inst *CursorSoftOff) Decode(opcode vm.OpCode, r *vm.BytecodeReader) (err error) {
	inst.frame, err = r.EndFrame()
	return
}

func (inst *UserputSoftOn) Decode(opcode vm.OpCode, r *vm.BytecodeReader) (err error) {
	inst.frame, err = r.EndFrame()
	return
}

func (inst *UserputSoftOff) Decode(opcode vm.OpCode, r *vm.BytecodeReader) (err error) {
	inst.frame, err = r.EndFrame()
	return
}

func (inst *SetCursorImg) Decode(opcode vm.OpCode, r *vm.BytecodeReader) (err error) {
	inst.Cursor = r.ReadByteParam(opcode, vm.ParamPos1)
	inst.Char = r.ReadByteParam(opcode, vm.ParamPos2)
	inst.frame, err = r.EndFrame()
	return
}

func (inst *SetCursorHotspot) Decode(opcode vm.OpCode, r *vm.BytecodeReader) (err error) {
	inst.Cursor = r.ReadByteParam(opcode, vm.ParamPos1)
	inst.X = r.ReadByteParam(opcode, vm.ParamPos2)
	inst.Y = r.ReadByteParam(opcode, vm.ParamPos3)
	inst.frame, err = r.EndFrame()
	return
}

func (inst *InitCursor) Decode(opcode vm.OpCode, r *vm.BytecodeReader) (err error) {
	inst.Cursor = r.ReadByteParam(opcode, vm.ParamPos1)
	inst.frame, err = r.EndFrame()
	return
}

func (inst *InitCharset) Decode(opcode vm.OpCode, r *vm.BytecodeReader) (err error) {
	inst.Charset = r.ReadByteParam(opcode, vm.ParamPos1)
	inst.frame, err = r.EndFrame()
	return
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
