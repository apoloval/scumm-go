package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

// CursorShow is a cursor command that shows the cursor.
type CursorShow struct{}

func (inst CursorShow) Acronym() string { return "CRS" }

// CursorHide is a cursor command that hides the cursor.
type CursorHide struct{}

func (inst CursorHide) Acronym() string { return "CRH" }

// CursorInc is a cursor command that increments the cursor counter. Also known as CursorSoftOn in
// ScummVM.
type CursorInc struct{}

func (inst CursorInc) Acronym() string { return "CRINC" }

// CursorDec is a cursor command that decrements the cursor counter. Also known as CursorDec in
// ScummVM.
type CursorDec struct{}

func (inst CursorDec) Acronym() string { return "CRDEC" }

// UserputEnable is a cursor command that enables user input. Also known as UserputOn in ScummVM.
type UserputEnable struct{}

func (inst UserputEnable) Acronym() string { return "UPE" }

// UserputDisable is a cursor command that disables user input. Also known as UserputDisable in
// ScummVM.
type UserputDisable struct{}

func (inst UserputDisable) Acronym() string { return "UPD" }

// UserputInc is a cursor command that increments the user input counter. Also known as
// UserputSoftOn in ScummVM.
type UserputInc struct{}

func (inst UserputInc) Acronym() string { return "UPINC" }

// UserputDec is a cursor command that decrements the user input counter. Also known as
// UserputSoftOff in ScummVM.
type UserputDec struct{}

func (inst UserputDec) Acronym() string { return "UPDEC" }

// SetCursorImg is a cursor command that sets the cursor image.
type SetCursorImg struct {
	Cursor vm.Param `op:"p8" pos:"1"`
	Char   vm.Param `op:"p8" pos:"2"`
}

func (inst SetCursorImg) Acronym() string { return "CRIMG" }

// SetCursorHotspot is a cursor command that sets the cursor hotspot.
type SetCursorHotspot struct {
	Cursor vm.Param `op:"p8" pos:"1"`
	X      vm.Param `op:"p8" pos:"2"`
	Y      vm.Param `op:"p8" pos:"3"`
}

func (inst SetCursorHotspot) Acronym() string { return "CRHOT" }

// CursorSelect is a cursor command to select the cursor.
type CursorSelect struct {
	Cursor vm.Param `op:"p8" pos:"1"`
}

func (inst CursorSelect) Acronym() string { return "CRSEL" }

// CharsetSelect is a cursor command to select the charset.
type CharsetSelect struct {
	Charset vm.Param `op:"p8" pos:"1" fmt:"id:charset"`
}

func (inst CharsetSelect) Acronym() string { return "CHSEL" }

func decodeCursorCommand(opcode vm.OpCode, r *vm.BytecodeDecoder) (inst vm.Instruction, err error) {
	sub := r.DecodeOpCode()

	switch sub & 0x1F {
	case 0x01:
		inst = new(CursorShow)
	case 0x02:
		inst = new(CursorHide)
	case 0x03:
		inst = new(UserputEnable)
	case 0x04:
		inst = new(UserputDisable)
	case 0x05:
		inst = new(CursorInc)
	case 0x06:
		inst = new(CursorDec)
	case 0x07:
		inst = new(UserputInc)
	case 0x08:
		inst = new(UserputDec)
	case 0x0A:
		inst = new(SetCursorImg)
	case 0x0B:
		inst = new(SetCursorHotspot)
	case 0x0C:
		inst = new(CursorSelect)
	case 0x0D:
		inst = new(CharsetSelect)
	default:
		return nil, fmt.Errorf("unimplemented opcode %02X %02X for cursor command", opcode, sub)
	}

	err = vm.DecodeOperands(sub, r, inst)
	return inst, err
}
