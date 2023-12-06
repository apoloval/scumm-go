package inst

import (
	"fmt"
	"strings"

	"github.com/apoloval/scumm-go/vm"
)

// Move is a move instruction that puts the value from Src into Dest
type Move struct {
	Result vm.VarRef `op:"result"`
	Src    vm.Param  `op:"p16" pos:"1"`
}

// Mnemonic implements the Instruction interface.
func (inst Move) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("%s = %s", inst.Result.Display(st), inst.Src.Display(st))
}

// SetVarRange is a instruction that sets a range of variables to the given values.
type SetVarRange struct {
	Result vm.VarRef     `op:"result"`
	Count  vm.Constant   `op:"8"`
	Values []vm.Constant `op:"16"`
}

// Decode implements the Instruction interface.
func (inst *SetVarRange) DecodeOperands(opcode vm.OpCode, r *vm.BytecodeDecoder) error {
	inst.Result = r.DecodeVarRef()
	inst.Count = r.DecodeByteConstant(vm.NumberFormatDecimal)
	for i := 0; i < int(inst.Count.Value); i++ {
		if opcode&0x80 > 0 {
			inst.Values = append(inst.Values, r.DecodeWordConstant(vm.NumberFormatDecimal))
		} else {
			inst.Values = append(inst.Values, r.DecodeByteConstant(vm.NumberFormatDecimal))
		}
	}
	return nil
}

func (inst SetVarRange) Display(st *vm.SymbolTable) string {
	vals := make([]string, len(inst.Values))
	for i, val := range inst.Values {
		vals[i] = val.Display(st)
	}
	return fmt.Sprintf("SetVarRange %s: [ %s ]",
		inst.Result.Display(st),
		strings.Join(vals, ", "),
	)
}
