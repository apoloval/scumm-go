package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

// Decode decodes an instruction from the bytecode reader.
func Decode(r *vm.BytecodeDecoder) (inst vm.Instruction, err error) {
	opcode := r.DecodeOpCode()
	switch opcode {
	case 0x00:
		inst = new(EndOfCode)
	case 0x01, 0x21, 0x41, 0x61, 0x81, 0xA1, 0xC1, 0xE1:
		inst = new(ActorPut)
	case 0x04, 0x84:
		inst = new(BranchUnlessGreaterEqual)
	case 0x07, 0x47, 0x87, 0xC7:
		inst = new(SetObjectState)
	case 0x08, 0x88:
		inst = new(BranchUnlessNotEqual)
	case 0x0A, 0x2A, 0x4A, 0x6A, 0x8A, 0xAA, 0xCA, 0xEA:
		inst = new(StartScript)
	case 0x0C:
		return decodeResourceRoutine(opcode, r)
	case 0x0F, 0x4F, 0x8F, 0xCF:
		inst = new(BranchUnlessState)
	case 0x10, 0x90:
		inst = new(GetObjectOwner)
	case 0x11, 0x51, 0x91, 0xD1:
		inst = new(AnimateActor)
	case 0x13, 0x53, 0x93, 0xD3:
		inst = new(ActorOps)
	case 0x14, 0x94:
		inst = new(Print)
	case 0x17, 0x97:
		inst = new(And)
	case 0x18:
		inst = new(Jump)
	case 0x1A, 0x9A:
		inst = new(Move)
	case 0x1F, 0x3F, 0x5F, 0x7F, 0x9F, 0xBF, 0xDF, 0xFF:
		inst = new(BranchUnlessActorInBox)
	case 0x26, 0xA6:
		inst = new(SetVarRange)
	case 0x27:
		return decodeStringOp(opcode, r)
	case 0x28:
		inst = new(BranchUnlessZero)
	case 0x29, 0x69, 0xA9, 0xE9:
		inst = new(SetObjectOwner)
	case 0x2C:
		return decodeCursorCommand(opcode, r)
	case 0x2D, 0x6D, 0xAD, 0xED:
		inst = new(PutActorInRoom)
	case 0x33:
		return decodeRoomOp(opcode, r)
	case 0x37, 0x77, 0xB7, 0xF7:
		inst = new(StartObject)
	case 0x38, 0xB8:
		inst = new(BranchUnlessLessEqual)
	case 0x44, 0xC4:
		inst = new(BranchUnlessLess)
	case 0x48, 0xC8:
		inst = new(BranchUnlessEqual)
	case 0x50, 0xD0:
		inst = new(PickUpObject)
	case 0x54, 0xD4:
		inst = new(SetObjectName)
	case 0x5C:
		inst = new(RoomFade)
	case 0x5D, 0xDD:
		inst = new(SetClass)
	case 0x68, 0xE8:
		inst = new(ScriptRunning)
	case 0x72, 0xF2:
		inst = new(LoadRoom)
	case 0x78, 0xF8:
		inst = new(BranchUnlessGreater)
	case 0x80:
		inst = new(BreakHere)
	case 0x98:
		return decodeSystemOp(opcode, r)
	case 0xCC:
		inst = new(PseudoRoom)
	case 0xA8:
		inst = new(BranchUnlessNotZero)
	case 0xAC:
		inst = new(Expression)
	default:
		return nil, fmt.Errorf("unknown opcode %02X", opcode)
	}
	err = vm.DecodeOperands(opcode, r, inst)
	return inst, err
}
