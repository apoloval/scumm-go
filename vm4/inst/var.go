package inst

import (
	"github.com/apoloval/scumm-go/vm"
)

// Move is a move instruction that puts the value from Src into Dest
type Move struct {
	Result vm.VarRef `op:"result"`
	Src    vm.Param  `op:"p16" pos:"1" fmt:"dec"`
}

func (inst Move) Acronym() string { return "MOVE" }

// SetVarRange is a instruction that sets a range of variables to the given values.
type SetVarRange struct {
	Result vm.VarRef     `op:"result"`
	Values []vm.Constant `op:"16"`
}

func (inst SetVarRange) Acronym() string { return "SETVR" }

func (inst SetVarRange) DisplayOperands(st *vm.SymbolTable) []string {
	ops := []string{
		inst.Result.Display(st),
	}
	for _, val := range inst.Values {
		ops = append(ops, val.Display(st))
	}
	return ops
}

// Decode implements the Instruction interface.
func (inst *SetVarRange) DecodeOperands(opcode vm.OpCode, r *vm.BytecodeDecoder) error {
	inst.Result = r.DecodeVarRef()
	count := r.DecodeByte()
	for i := 0; i < int(count); i++ {
		if opcode&0x80 > 0 {
			inst.Values = append(inst.Values, r.DecodeWordConstant(vm.NumberFormatDecimal))
		} else {
			inst.Values = append(inst.Values, r.DecodeByteConstant(vm.NumberFormatDecimal))
		}
	}
	return nil
}
