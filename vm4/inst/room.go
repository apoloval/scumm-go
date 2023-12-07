package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

type RoomFade struct {
	Effect vm.Param `op:"p16" pos:"1"`
}

func (inst RoomFade) Acronym() string { return "ROFA" }

func (inst *RoomFade) DecodeOperands(opcode vm.OpCode, r *vm.BytecodeDecoder) error {
	sub := r.DecodeOpCode()
	if sub&0x1F == 3 {
		inst.Effect = r.DecodeWordParam(sub, vm.ParamPos1, vm.NumberFormatHex)
	}
	return nil
}

type PseudoRoom struct {
	Value       vm.Constant   `op:"8" fmt:"dec"`
	ResourceIDs []vm.Constant `op:"8" fmt:"dec"`
}

func (inst PseudoRoom) Acronym() string { return "PSRO" }

func (inst PseudoRoom) DisplayOperands(st *vm.SymbolTable) (ops []string) {
	ops = []string{inst.Value.Display(st)}
	for _, id := range inst.ResourceIDs {
		ops = append(ops, id.Display(st))
	}
	return ops
}

func (inst *PseudoRoom) DecodeOperands(opcode vm.OpCode, r *vm.BytecodeDecoder) error {
	inst.Value = r.DecodeByteConstant(vm.NumberFormatRoomID)
	inst.ResourceIDs = r.DecodeNullTerminatedBytes(vm.NumberFormatRoomID)
	params := []vm.Param{inst.Value}
	for _, id := range inst.ResourceIDs {
		params = append(params, id)
	}
	return nil
}

type RoomSetScrollLimits struct {
	MinX vm.Param `op:"p16" pos:"1"`
	MaxX vm.Param `op:"p16" pos:"2"`
}

func (inst RoomSetScrollLimits) Acronym() string { return "ROSL" }

type RoomInitScreen struct {
	B vm.Param `op:"p16" pos:"1" fmt:"dec"`
	H vm.Param `op:"p16" pos:"2" fmt:"dec"`
}

func (inst RoomInitScreen) Acronym() string { return "ROIS" }

func decodeRoomOp(opcode vm.OpCode, r *vm.BytecodeDecoder) (inst vm.Instruction, err error) {
	sub := r.DecodeOpCode()
	switch sub & 0x1F {
	case 0x01:
		inst = new(RoomSetScrollLimits)
	case 0x03:
		inst = new(RoomInitScreen)
	default:
		return nil, fmt.Errorf("unknown opcode %02X %02X for room operation", opcode, sub)
	}
	err = vm.DecodeOperands(sub, r, inst)
	return inst, err
}
