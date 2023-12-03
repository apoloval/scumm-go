package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

type CursorShow struct{ instruction }
type CursorHide struct{ instruction }
type UserputOn struct{ instruction }
type UserputOff struct{ instruction }
type CursorSoftOn struct{ instruction }
type CursorSoftOff struct{ instruction }
type UserputSoftOn struct{ instruction }
type UserputSoftOff struct{ instruction }
type SetCursorImg struct {
	instruction
	Cursor vm.Param
	Char   vm.Param
}
type SetCursorHotspot struct {
	instruction
	Cursor vm.Param
	X      vm.Param
	Y      vm.Param
}
type InitCursor struct {
	instruction
	Cursor vm.Param
}
type InitCharset struct {
	instruction
	Charset vm.Param
}

func (inst CursorShow) Mnemonic() string     { return "CursorShow" }
func (inst CursorHide) Mnemonic() string     { return "CursorHide" }
func (inst UserputOn) Mnemonic() string      { return "UserputOn" }
func (inst UserputOff) Mnemonic() string     { return "UserputOff" }
func (inst CursorSoftOn) Mnemonic() string   { return "CursorSoftOn" }
func (inst CursorSoftOff) Mnemonic() string  { return "CursorSoftOff" }
func (inst UserputSoftOn) Mnemonic() string  { return "UserputSoftOn" }
func (inst UserputSoftOff) Mnemonic() string { return "UserputSoftOff" }
func (inst SetCursorImg) Mnemonic() string {
	return fmt.Sprintf("SetCursorImg %s, %s", inst.Cursor, inst.Char)
}
func (inst SetCursorHotspot) Mnemonic() string {
	return fmt.Sprintf("SetCursorHotspot %s, %s, %s", inst.Cursor, inst.X, inst.Y)
}
func (inst InitCursor) Mnemonic() string {
	return fmt.Sprintf("InitCursor %s", inst.Cursor)
}
func (inst InitCharset) Mnemonic() string {
	return fmt.Sprintf("InitCharset %s", inst.Charset)
}

func decodeCursorCommand(opcode vm.OpCode, r *vm.BytecodeReader) (vm.Instruction, error) {
	sub := r.ReadOpCode() & 0x1F

	var err error
	switch sub {
	case 0x01:
		var inst CursorShow
		inst.bytecode, err = r.EndFrame()
		return inst, err
	case 0x02:
		var inst CursorHide
		inst.bytecode, err = r.EndFrame()
		return inst, err
	case 0x03:
		var inst UserputOn
		inst.bytecode, err = r.EndFrame()
		return inst, err
	case 0x04:
		var inst UserputOff
		inst.bytecode, err = r.EndFrame()
		return inst, err
	case 0x05:
		var inst CursorSoftOn
		inst.bytecode, err = r.EndFrame()
		return inst, err
	case 0x06:
		var inst CursorSoftOff
		inst.bytecode, err = r.EndFrame()
		return inst, err
	case 0x07:
		var inst UserputSoftOn
		inst.bytecode, err = r.EndFrame()
		return inst, err
	case 0x08:
		var inst UserputSoftOff
		inst.bytecode, err = r.EndFrame()
		return inst, err
	case 0x0A:
		var inst SetCursorImg
		inst.Cursor = r.ReadByteParam(sub, vm.ParamPos1)
		inst.Char = r.ReadByteParam(sub, vm.ParamPos2)
		inst.bytecode, err = r.EndFrame()
		return inst, err
	case 0x0B:
		var inst SetCursorHotspot
		inst.Cursor = r.ReadByteParam(sub, vm.ParamPos1)
		inst.X = r.ReadByteParam(sub, vm.ParamPos2)
		inst.Y = r.ReadByteParam(sub, vm.ParamPos3)
		inst.bytecode, err = r.EndFrame()
		return inst, err
	case 0x0C:
		var inst InitCursor
		inst.Cursor = r.ReadByteParam(sub, vm.ParamPos1)
		inst.bytecode, err = r.EndFrame()
		return inst, err
	case 0x0D:
		var inst InitCharset
		inst.Charset = r.ReadByteParam(sub, vm.ParamPos1)
		inst.bytecode, err = r.EndFrame()
		return inst, err

	default:
		return nil, fmt.Errorf("unimplemented sub-opcode %02X for cursor command", sub)
	}
}
