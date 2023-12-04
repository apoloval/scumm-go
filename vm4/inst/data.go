package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

// Move is a move instruction that puts the value from Src into Dest
type Move struct {
	base
	Dest vm.Pointer
	Src  vm.Param
}

// Mnemonic implements the Instruction interface.
func (inst Move) Mnemonic(st *vm.SymbolTable) string {
	return fmt.Sprintf("%s = %s",
		inst.Dest.Display(st),
		inst.Src.Display(st),
	)
}

// Decode implements the Instruction interface.
func (inst *Move) Decode(opcode vm.OpCode, r *vm.BytecodeReader) error {
	inst.Dest = r.ReadPointer()
	inst.Src = r.ReadWordParam(opcode, vm.ParamPos1, vm.ParamFormatNumber)
	return inst.base.Decode(opcode, r)
}

// SetVarRange is a instruction that sets a range of variables to the given values.
type SetVarRange struct {
	base
	Dest   vm.Pointer
	Count  vm.Constant
	Values []vm.Constant
}

// Mnemonic implements the Instruction interface.
func (inst SetVarRange) Mnemonic(st *vm.SymbolTable) string {
	return fmt.Sprintf("SetVarRange %s, %s, %v",
		inst.Dest.Display(st),
		inst.Count.Display(st),
		inst.Values,
	)
}

// Decode implements the Instruction interface.
func (inst *SetVarRange) Decode(opcode vm.OpCode, r *vm.BytecodeReader) error {
	inst.Dest = r.ReadPointer()
	inst.Count = r.ReadByteConstant(vm.ParamFormatNumber)
	for i := 0; i < int(inst.Count.Value); i++ {
		if opcode&0x80 > 0 {
			inst.Values = append(inst.Values, r.ReadWordConstant(vm.ParamFormatNumber))
		} else {
			inst.Values = append(inst.Values, r.ReadByteConstant(vm.ParamFormatNumber))
		}
	}
	return inst.base.Decode(opcode, r)
}
