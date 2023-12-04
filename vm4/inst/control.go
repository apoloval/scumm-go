package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

// StopObjectCode is a stop instruction that stops the execution of the current script.
type StopObjectCode struct {
	base
}

// Mnemonic implements the Instruction interface.
func (inst StopObjectCode) Mnemonic(*vm.SymbolTable) string { return "StopObjectCode" }

// Goto is a goto instruction that jumps to the given address.
type Goto struct {
	base
	Goto vm.ProgramAddress
}

// Mnemonic implements the Instruction interface.
func (inst Goto) Mnemonic(st *vm.SymbolTable) string {
	return fmt.Sprintf("goto %s", inst.Goto.Display(st, vm.ParamFormatNumber))
}

// Decode implements the Instruction interface.
func (inst *Goto) Decode(opcode vm.OpCode, r *vm.BytecodeReader) error {
	inst.Goto = r.ReadRelativeJump()
	return inst.base.Decode(opcode, r)

}

// Branch is a base type to represent branching instructions.
type Branch struct {
	base
	Left  vm.WordPointer
	Right vm.Param
	Goto  vm.ProgramAddress
}

// IsEqual is a branching instruction that jumps to the given address unless the two operands are
// equal.
type IsEqual struct{ Branch }

func (inst IsEqual) Mnemonic(st *vm.SymbolTable) string {
	return fmt.Sprintf("unless (%s == %s) goto %s",
		inst.Left.Display(st, vm.ParamFormatNumber),
		inst.Right.Display(st, vm.ParamFormatNumber),
		inst.Goto.Display(st, vm.ParamFormatNumber),
	)
}

func (inst *IsEqual) Decode(opcode vm.OpCode, r *vm.BytecodeReader) error {
	inst.Left = r.ReadWordPointer()
	inst.Right = r.ReadWordParam(opcode, vm.ParamPos1)
	inst.Goto = r.ReadRelativeJump()
	return inst.base.Decode(opcode, r)

}
