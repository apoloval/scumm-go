package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

// CursorShow is a cursor command that shows the cursor.
type CursorShow struct{}

// CursorHide is a cursor command that hides the cursor.
type CursorHide struct{}

// UserputOn is a cursor command that enables user input.
type UserputOn struct{}

// UserputOff is a cursor command that disables user input.
type UserputOff struct{}

// CursorSoftOn is a cursor command that increments the cursor counter.
type CursorSoftOn struct{}

// CursorSoftOff is a cursor command that decrements the cursor counter.
type CursorSoftOff struct{}

// UserputSoftOn is a cursor command that increments the user input counter.
type UserputSoftOn struct{}

// UserputSoftOff is a cursor command that decrements the user input counter.
type UserputSoftOff struct{}

// SetCursorImg is a cursor command that sets the cursor image.
type SetCursorImg struct {
	Cursor vm.Param `op:"p8" pos:"1"`
	Char   vm.Param `op:"p8" pos:"2"`
}

// SetCursorHotspot is a cursor command that sets the cursor hotspot.
type SetCursorHotspot struct {
	Cursor vm.Param `op:"p8" pos:"1"`
	X      vm.Param `op:"p8" pos:"2"`
	Y      vm.Param `op:"p8" pos:"3"`
}

// CursorSet is a cursor command that initializes the cursor.
type CursorSet struct {
	Cursor vm.Param `op:"p8" pos:"1"`
}

// CharsetSet is a cursor command that initializes the charset.
type CharsetSet struct {
	Charset vm.Param `op:"p8" pos:"1"`
}

func decodeCursorCommand(opcode vm.OpCode, r *vm.BytecodeDecoder) (inst vm.Instruction, err error) {
	sub := r.DecodeOpCode()

	switch sub & 0x1F {
	case 0x01:
		inst = new(CursorShow)
	case 0x02:
		inst = new(CursorHide)
	case 0x03:
		inst = new(UserputOn)
	case 0x04:
		inst = new(UserputOff)
	case 0x05:
		inst = new(CursorSoftOn)
	case 0x06:
		inst = new(CursorSoftOff)
	case 0x07:
		inst = new(UserputSoftOn)
	case 0x08:
		inst = new(UserputSoftOff)
	case 0x0A:
		inst = new(SetCursorImg)
	case 0x0B:
		inst = new(SetCursorHotspot)
	case 0x0C:
		inst = new(CursorSet)
	case 0x0D:
		inst = new(CharsetSet)
	default:
		return nil, fmt.Errorf("unimplemented opcode %02X %02X for cursor command", opcode, sub)
	}

	err = vm.DecodeOperands(sub, r, inst)
	return inst, err
}
