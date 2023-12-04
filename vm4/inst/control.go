package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

// StopObjectCode is a stop instruction that stops the execution of the current script.
type StopObjectCode struct{ base }

// Goto is a goto instruction that jumps to the given address.
type Goto struct {
	base
	Goto vm.Constant
}

// Decode implements the Instruction interface.
func (inst *Goto) Decode(opcode vm.OpCode, r *vm.BytecodeReader) error {
	inst.Goto = r.ReadRelativeJump()
	return inst.decodeWithParams(r, inst.Goto)
}

// Branch is a base type to represent branching instructions.
type Branch struct {
	base
	Left  vm.WordPointer
	Right vm.Param
	Goto  vm.Constant
}

// IsEqual is a branching instruction that jumps to the given address unless the two operands are
// equal.
type IsEqual struct{ Branch }

func (inst IsEqual) Mnemonic(st *vm.SymbolTable) string {
	return fmt.Sprintf("Unless (%s == %s) Goto %s",
		inst.Left.Display(st),
		inst.Right.Display(st),
		inst.Goto.Display(st),
	)
}

func (inst *IsEqual) Decode(opcode vm.OpCode, r *vm.BytecodeReader) error {
	inst.Left = r.ReadWordPointer()
	inst.Right = r.ReadWordParam(opcode, vm.ParamPos1, vm.ParamFormatNumber)
	inst.Goto = r.ReadRelativeJump()
	return inst.base.Decode(opcode, r)

}
