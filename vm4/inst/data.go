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
		inst.Dest.Display(st, vm.ParamFormatNumber),
		inst.Src.Display(st, vm.ParamFormatNumber),
	)
}

// Decode implements the Instruction interface.
func (inst *Move) Decode(opcode vm.OpCode, r *vm.BytecodeReader) error {
	inst.Dest = r.ReadPointer()
	inst.Src = r.ReadWordParam(opcode, vm.ParamPos1)
	return inst.base.Decode(opcode, r)
}

// SetVarRange is a instruction that sets a range of variables to the given values.
type SetVarRange struct {
	base
	Dest   vm.Pointer
	Count  vm.ByteConstant
	Values []vm.WordConstant
}

// Mnemonic implements the Instruction interface.
func (inst SetVarRange) Mnemonic(st *vm.SymbolTable) string {
	return fmt.Sprintf("SetVarRange %s, %s, %v",
		inst.Dest.Display(st, vm.ParamFormatVarID),
		inst.Count.Display(st, vm.ParamFormatNumber),
		inst.Values,
	)
}

// Decode implements the Instruction interface.
func (inst *SetVarRange) Decode(opcode vm.OpCode, r *vm.BytecodeReader) error {
	inst.Dest = r.ReadPointer()
	inst.Count = r.ReadByteConstant()
	for i := 0; i < int(inst.Count); i++ {
		if opcode&0x80 > 0 {
			inst.Values = append(inst.Values, r.ReadWordConstant())
		} else {
			inst.Values = append(inst.Values, vm.WordConstant(r.ReadByteConstant()))
		}
	}
	return inst.base.Decode(opcode, r)
}
