package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

// Decode decodes an instruction from the bytecode reader.
func Decode(r *vm.BytecodeDecoder) (inst vm.Instruction, err error) {
	opcode := r.DecodeOpCode()
	switch opcode {
	case 0x00, 0xA0:
		inst = new(StopObjectCode)
	case 0x01, 0x21, 0x41, 0x61, 0x81, 0xA1, 0xC1, 0xE1:
		inst = new(ActorPut)
	case 0x02:
		inst = new(StartMusic)
	case 0x03, 0x83:
		inst = new(GetActorRoom)
	case 0x04, 0x84:
		inst = new(BranchUnlessGreaterEqual)
	case 0x05, 0x85:
		inst = new(DrawObject)
	case 0x06, 0x86:
		inst = new(GetActorElevation)
	case 0x07, 0x47, 0x87, 0xC7:
		inst = new(SetObjectState)
	case 0x08, 0x88:
		inst = new(BranchUnlessNotEqual)
	case 0x09, 0x49, 0x89, 0xC9:
		inst = new(FaceActor)
	case 0x0A, 0x2A, 0x4A, 0x6A, 0x8A, 0xAA, 0xCA, 0xEA:
		inst = new(StartScript)
	case 0x0B, 0x4B, 0x8B, 0xCB:
		inst = new(GetVerbEntrypoint)
	case 0x0C:
		return decodeResourceRoutine(opcode, r)
	case 0x0D, 0x4D, 0x8D, 0xCD:
		inst = new(WalkActorToActor)
	case 0x0E, 0x4E, 0x8E, 0xCE:
		inst = new(PutActorAtObject)
	case 0x0F, 0x4F, 0x8F, 0xCF:
		inst = new(BranchUnlessState)
	case 0x10, 0x90:
		inst = new(GetObjectOwner)
	case 0x11, 0x51, 0x91, 0xD1:
		inst = new(AnimateActor)
	case 0x12, 0x92:
		inst = new(PanCameraTo)
	case 0x13, 0x53, 0x93, 0xD3:
		inst = new(Actor)
	case 0x14, 0x94:
		inst = new(Print)
	case 0x15, 0x55, 0x95, 0xD5:
		inst = new(ActorFromPos)
	case 0x16, 0x96:
		inst = new(GetRandomNumber)
	case 0x17, 0x97:
		inst = new(And)
	case 0x18:
		inst = new(Jump)
	case 0x19, 0x39, 0x59, 0x79, 0x99, 0xB9, 0xD9, 0xF9:
		inst = new(DoSentence)
	case 0x1A, 0x9A:
		inst = new(Move)
	case 0x1B, 0x9B:
		inst = new(Mult)
	case 0x1C, 0x9C:
		inst = new(StartSound)
	case 0x1D, 0x9D:
		inst = new(BranchUnlessClass)
	case 0x1E, 0x3E, 0x5E, 0x7E, 0x9E, 0xBE, 0xDE, 0xFE:
		inst = new(WalkActorTo)
	case 0x1F, 0x5F, 0x9F, 0xDF:
		inst = new(BranchUnlessActorInBox)
	case 0x20:
		inst = new(StopMusic)
	case 0x22, 0xA2:
		inst = new(Game)
	case 0x23, 0xA3:
		inst = new(GetActorY)
	case 0x24, 0x64, 0xA4, 0xE4:
		inst = new(LoadRoomWithEgo)
	case 0x26, 0xA6:
		inst = new(SetVarRange)
	case 0x27:
		return decodeStringOp(opcode, r)
	case 0x28:
		inst = new(BranchUnlessZero)
	case 0x29, 0x69, 0xA9, 0xE9:
		inst = new(SetObjectOwner)
	case 0x2B:
		inst = new(DelayVar)
	case 0x2C:
		return decodeCursorCommand(opcode, r)
	case 0x2D, 0x6D, 0xAD, 0xED:
		inst = new(PutActorInRoom)
	case 0x2E:
		inst = new(Delay)
	case 0x2F:
		inst = new(BranchUnlessNotState)
	case 0x30, 0xB0:
		return decodeBoxOp(opcode, r)
	case 0x31, 0xB1:
		inst = new(GetInventoryCount)
	case 0x32:
		inst = new(SetCameraAt)
	case 0x33:
		return decodeRoomOp(opcode, r)
	case 0x34, 0x74, 0xB4, 0xF4:
		inst = new(GetDistance)
	case 0x35, 0x75, 0xB5, 0xF5:
		inst = new(FindObject)
	case 0x36, 0x76, 0xB6, 0xF6:
		inst = new(WalkActorToObject)
	case 0x37, 0x77, 0xB7, 0xF7:
		inst = new(StartObject)
	case 0x38, 0xB8:
		inst = new(BranchUnlessLessEqual)
	case 0x3A, 0xBA:
		inst = new(Sub)
	case 0x3B, 0xBB:
		inst = new(GetActorScale)
	case 0x3C, 0xBC:
		inst = new(StopSound)
	case 0x3D, 0x7D, 0xBD, 0xFD:
		inst = new(FindInventory)
	case 0x3F, 0x7F, 0xBF, 0xFF:
		inst = new(DrawBox)
	case 0x40:
		inst = new(CutScene)
	case 0x42, 0xC2:
		inst = new(ChainStript)
	case 0x43, 0xC3:
		inst = new(GetActorX)
	case 0x44, 0xC4:
		inst = new(BranchUnlessLess)
	case 0x46:
		inst = new(Increment)
	case 0x48, 0xC8:
		inst = new(BranchUnlessEqual)
	case 0x50, 0xD0:
		inst = new(PickUpObject)
	case 0x52, 0xD2:
		inst = new(ActorFollowCamera)
	case 0x54, 0xD4:
		inst = new(SetObjectName)
	case 0x56, 0xD6:
		inst = new(GetActorMoving)
	case 0x57, 0xD7:
		inst = new(Or)
	case 0x58:
		return decodeOverrideOp(opcode, r)
	case 0x5A, 0xDA:
		inst = new(And)
	case 0x5B, 0xDB:
		inst = new(Div)
	case 0x5C:
		inst = new(RoomFade)
	case 0x5D, 0xDD:
		inst = new(SetClass)
	case 0x60, 0xE0:
		inst = new(FreezeScripts)
	case 0x62, 0xE2:
		inst = new(StopScript)
	case 0x63, 0xE3:
		inst = new(GetActorFacing)
	case 0x66, 0xE6:
		inst = new(GetActorClosestObject)
	case 0x67, 0xE7:
		inst = new(StringWidth)
	case 0x68, 0xE8:
		inst = new(ScriptRunning)
	case 0x6B, 0xEB:
		inst = new(Debug)
	case 0x6C, 0xEC:
		inst = new(GetActorWidth)
	case 0x6E, 0xEE:
		inst = new(StopObjectScript)
	case 0x70, 0xF0:
		inst = new(Lights)
	case 0x71, 0xF1:
		inst = new(GetActorCostume)
	case 0x72, 0xF2:
		inst = new(LoadRoom)
	case 0x78, 0xF8:
		inst = new(BranchUnlessGreater)
	case 0x7A, 0xFA:
		inst = new(Verb)
	case 0x7B, 0xFB:
		inst = new(GetActorWalkBox)
	case 0x7C, 0xFC:
		inst = new(IsSoundRunning)
	case 0x80:
		inst = new(BreakHere)
	case 0x98:
		return decodeSystemOp(opcode, r)
	case 0xC0:
		inst = new(EndCutScene)
	case 0xC6:
		inst = new(Decrement)
	case 0xCC:
		inst = new(PseudoRoom)
	case 0xA7:
		return decodeSaveLoadVarsOp(opcode, r)
	case 0xA8:
		inst = new(BranchUnlessNotZero)
	case 0xAB:
		return decodeSaveRestoreDeleteVerbs(opcode, r)
	case 0xAC:
		inst = new(Expression)
	case 0xAE:
		return decodeWaitOp(opcode, r)
	case 0xD8:
		inst = new(PrintEgo)
	default:
		return nil, fmt.Errorf("unknown opcode %02X", opcode)
	}
	err = vm.DecodeOperands(opcode, r, inst)
	return inst, err
}
