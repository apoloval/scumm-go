package inst

import "github.com/apoloval/scumm-go/vm"

type RoomFade struct {
	base
	Effect vm.Param
}

func (inst *RoomFade) Decode(opcode vm.OpCode, r *vm.BytecodeReader) error {
	sub := r.ReadOpCode()
	inst.Effect = vm.Const(0)
	if sub&0x1F == 3 {
		inst.Effect = r.ReadWordParam(opcode, vm.ParamPos1, vm.NumberFormatHex)
	}
	return inst.base.decodeWithParams(r, inst.Effect)
}

type PseudoRoom struct {
	base
	Value       vm.Constant
	ResourceIDs []vm.Constant
}

func (inst *PseudoRoom) Decode(opcode vm.OpCode, r *vm.BytecodeReader) error {
	inst.Value = r.ReadByteConstant(vm.NumberFormatDecimal)
	inst.ResourceIDs = r.ReadNullTerminatedBytes()
	params := []vm.Param{inst.Value}
	for _, id := range inst.ResourceIDs {
		params = append(params, id)
	}
	return inst.base.decodeWithParams(r, params...)
}
