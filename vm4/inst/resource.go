package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

type resourceRoutine struct {
	base
	ResourceID     vm.Param
	resourceFormat vm.NumberFormat
}

func withResourceRoutine(name string, format vm.NumberFormat) resourceRoutine {
	return resourceRoutine{base: withName(name), resourceFormat: format}
}

type LoadScript struct{ resourceRoutine }
type LoadSound struct{ resourceRoutine }
type LoadCostume struct{ resourceRoutine }
type LoadRoom struct{ resourceRoutine }
type LoadCharset struct{ resourceRoutine }

type NukeScript struct{ resourceRoutine }
type NukeSound struct{ resourceRoutine }
type NukeCostume struct{ resourceRoutine }
type NukeRoom struct{ resourceRoutine }
type NukeCharset struct{ resourceRoutine }

type LockSound struct{ resourceRoutine }
type LockScript struct{ resourceRoutine }
type LockCostume struct{ resourceRoutine }
type LockRoom struct{ resourceRoutine }

type UnlockSound struct{ resourceRoutine }
type UnlockScript struct{ resourceRoutine }
type UnlockCostume struct{ resourceRoutine }
type UnlockRoom struct{ resourceRoutine }

type ClearHeap struct{ base }

type LoadObject struct {
	base
	RoomID   vm.Param
	ObjectID vm.Param
}

func (inst *resourceRoutine) Decode(opcode vm.OpCode, r *vm.BytecodeReader) error {
	inst.ResourceID = r.ReadByteParam(opcode, vm.ParamPos1, inst.resourceFormat)
	return inst.base.decodeWithParams(r, inst.ResourceID)
}

func (inst *LoadObject) Decode(opcode vm.OpCode, r *vm.BytecodeReader) error {
	inst.RoomID = r.ReadByteParam(opcode, vm.ParamPos1, vm.ParamFormatRoomID)
	inst.ObjectID = r.ReadWordParam(opcode, vm.ParamPos2, vm.ParamFormatNumber)
	return inst.base.decodeWithParams(r, inst.RoomID, inst.ObjectID)
}

func decodeResourceRoutine(opcode vm.OpCode, r *vm.BytecodeReader) (inst vm.Instruction, err error) {
	sub := r.ReadOpCode()
	switch sub & 0x1F {
	case 0x01:
		inst = &LoadScript{withResourceRoutine("LoadScript", vm.ParamFormatScriptID)}
	case 0x02:
		inst = &LoadSound{withResourceRoutine("LoadSound", vm.ParamFormatSoundID)}
	case 0x03:
		inst = &LoadCostume{withResourceRoutine("LoadCostume", vm.ParamFormatCostumeID)}
	case 0x04:
		inst = &LoadRoom{withResourceRoutine("LoadRoom", vm.ParamFormatRoomID)}
	case 0x05:
		inst = &NukeScript{withResourceRoutine("NukeScript", vm.ParamFormatScriptID)}
	case 0x06:
		inst = &NukeSound{withResourceRoutine("NukeSound", vm.ParamFormatSoundID)}
	case 0x07:
		inst = &NukeCostume{withResourceRoutine("NukeCostume", vm.ParamFormatCostumeID)}
	case 0x08:
		inst = &NukeRoom{withResourceRoutine("NukeRoom", vm.ParamFormatRoomID)}
	case 0x09:
		inst = &LockScript{withResourceRoutine("LockScript", vm.ParamFormatScriptID)}
	case 0x0A:
		inst = &LockSound{withResourceRoutine("LockSound", vm.ParamFormatSoundID)}
	case 0x0B:
		inst = &LockCostume{withResourceRoutine("LockCostume", vm.ParamFormatCostumeID)}
	case 0x0C:
		inst = &LockRoom{withResourceRoutine("LockRoom", vm.ParamFormatRoomID)}
	case 0x0D:
		inst = &UnlockScript{withResourceRoutine("UnlockScript", vm.ParamFormatScriptID)}
	case 0x0E:
		inst = &UnlockSound{withResourceRoutine("UnlockSound", vm.ParamFormatSoundID)}
	case 0x0F:
		inst = &UnlockCostume{withResourceRoutine("UnlockCostume", vm.ParamFormatCostumeID)}
	case 0x10:
		inst = &UnlockRoom{withResourceRoutine("UnlockRoom", vm.ParamFormatRoomID)}
	case 0x11:
		inst = &ClearHeap{withName("ClearHeap")}
	case 0x12:
		inst = &LoadCharset{withResourceRoutine("LoadCharset", vm.ParamFormatCharsetID)}
	case 0x13:
		inst = &NukeCharset{withResourceRoutine("NukeCharset", vm.ParamFormatCharsetID)}
	case 0x14:
		inst = &LoadObject{}
	default:
		return nil, fmt.Errorf("unknown opcode %02X %02X for resource routine", opcode, sub)
	}
	if err := inst.Decode(opcode, r); err != nil {
		return nil, err
	}
	return inst, nil
}
