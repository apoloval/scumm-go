package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

type PanCameraTo struct {
	X vm.Param `op:"p16" pos:"1" fmt:"dec"`
}

func (inst PanCameraTo) Acronym() string { return "PANC" }

type DoSentence struct {
	Verb vm.Param `op:"p8" pos:"1" fmt:"id:verb"`
	Obj1 vm.Param `op:"p16" pos:"2" fmt:"id:object"`
	Obj2 vm.Param `op:"p16" pos:"3" fmt:"id:object"`
}

func (inst DoSentence) Acronym() string { return "DOSENT" }

func (inst *DoSentence) DecodeOperands(opcode vm.OpCode, r *vm.BytecodeDecoder) error {
	inst.Verb = r.DecodeByteParam(opcode, vm.ParamPos1, vm.NumberFormatVerbID)
	if verb, ok := inst.Verb.(vm.Constant); ok && verb.Value == 0xFE {
		// Special case to break the sentence, not reading any further operand.
		return nil
	}
	inst.Obj1 = r.DecodeWordParam(opcode, vm.ParamPos2, vm.NumberFormatObjectID)
	inst.Obj2 = r.DecodeWordParam(opcode, vm.ParamPos3, vm.NumberFormatObjectID)
	return nil
}

type WaitForActor struct {
	Actor vm.Param `op:"p8" pos:"1" fmt:"id:actor"`
}

func (inst WaitForActor) Acronym() string { return "WAITA" }

type WaitForMessage struct{}

func (inst WaitForMessage) Acronym() string { return "WAITM" }

type WaitForCamera struct{}

func (inst WaitForCamera) Acronym() string { return "WAITC" }

type WaitForSentence struct{}

func (inst WaitForSentence) Acronym() string { return "WAITS" }

func decodeWaitOp(opcode vm.OpCode, r *vm.BytecodeDecoder) (inst vm.Instruction, err error) {
	sub := r.DecodeOpCode()
	switch sub & 0x1F {
	case 0x01:
		inst = &WaitForActor{Actor: r.DecodeByteParam(sub, vm.ParamPos1, vm.NumberFormatActorID)}
	case 0x02:
		inst = new(WaitForMessage)
	case 0x03:
		inst = new(WaitForCamera)
	case 0x04:
		inst = new(WaitForSentence)
	default:
		return nil, fmt.Errorf("unknown opcode %02X %02X for wait operation", opcode, sub)
	}
	return
}

type CutScene struct {
	Args vm.Params `op:"v16"`
}

func (inst CutScene) Acronym() string { return "CUTSCE" }

type EndCutScene struct{}

func (inst EndCutScene) Acronym() string { return "ENDCUT" }

type Lights struct {
	Arg1 vm.Param    `op:"p8" pos:"1" fmt:"dec"`
	Arg2 vm.Constant `op:"8" fmt:"dec"`
	Arg3 vm.Constant `op:"8" fmt:"dec"`
}

func (inst Lights) Acronym() string { return "LIGHTS" }
