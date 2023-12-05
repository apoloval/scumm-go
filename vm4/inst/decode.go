package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

// Decode decodes an instruction from the bytecode reader.
func Decode(r *vm.BytecodeReader) (inst vm.Instruction, err error) {
	opcode := r.ReadOpCode()
	switch opcode {
	case 0x00:
		inst = new(StopObjectCode)
	case 0x04, 0x84:
		inst = new(IsGreaterEqual)
	case 0x08, 0x88:
		inst = new(IsNotEqual)
	case 0x0A, 0x2A, 0x4A, 0x6A, 0x8A, 0xAA, 0xCA, 0xEA:
		inst = new(StartScript)
	case 0x0C:
		return decodeResourceRoutine(opcode, r)
	case 0x18:
		inst = new(Goto)
	case 0x1A, 0x9A:
		inst = new(Move)
	case 0x26, 0xA6:
		inst = new(SetVarRange)
	case 0x27:
		return decodeStringOp(opcode, r)
	case 0x28:
		inst = new(IsEqualZero)
	case 0x2C:
		return decodeCursorCommand(opcode, r)
	case 0x33:
		return decodeRoomOp(opcode, r)
	case 0x38, 0xB8:
		inst = new(IsLessEqual)
	case 0x44, 0xC4:
		inst = new(IsLess)
	case 0x48, 0xC8:
		inst = new(IsEqual)
	case 0x5C:
		inst = new(RoomFade)
	case 0x68, 0xE8:
		inst = new(GetScriptRunning)
	case 0x72, 0xF2:
		inst = new(LoadRoom)
	case 0x78, 0xF8:
		inst = new(IsGreater)
	case 0x80:
		inst = new(BreakHere)
	case 0xCC:
		inst = new(PseudoRoom)
	case 0xA8:
		inst = new(IsNotEqualZero)
	case 0xAC:
		inst = new(Expression)
	default:
		return nil, fmt.Errorf("unknown opcode %02X", opcode)
	}
	err = vm.DecodeOperands(opcode, r, inst)
	return inst, err
}
