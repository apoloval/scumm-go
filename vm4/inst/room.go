package inst

import (
	"fmt"
	"strings"

	"github.com/apoloval/scumm-go/vm"
)

type RoomFade struct {
	Effect vm.Param `op:"p16" pos:"1"`
}

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

func (inst *PseudoRoom) DecodeOperands(opcode vm.OpCode, r *vm.BytecodeDecoder) error {
	inst.Value = r.DecodeByteConstant(vm.NumberFormatDecimal)
	inst.ResourceIDs = r.DecodeNullTerminatedBytes()
	params := []vm.Param{inst.Value}
	for _, id := range inst.ResourceIDs {
		params = append(params, id)
	}
	return nil
}

func (inst PseudoRoom) Display(st *vm.SymbolTable) string {
	ids := make([]string, len(inst.ResourceIDs))
	for i, val := range inst.ResourceIDs {
		ids[i] = val.Display(st)
	}
	return fmt.Sprintf("PseudoRoom %s: [ %s ]",
		inst.Value.Display(st),
		strings.Join(ids, ", "),
	)
}

type RoomSetScrollLimits struct {
	MinX vm.Param `op:"p16" pos:"1"`
	MaxX vm.Param `op:"p16" pos:"2"`
}

type RoomInitScreen struct {
	B vm.Param `op:"p16" pos:"1"`
	H vm.Param `op:"p16" pos:"2"`
}

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
	return inst, nil
}
