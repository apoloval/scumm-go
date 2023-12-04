package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

// StopObjectCode is a stop instruction that stops the execution of the current script.
type StopObjectCode struct {
	instruction
}

// Mnemonic implements the Instruction interface.
func (inst StopObjectCode) Mnemonic(*vm.SymbolTable) string { return "StopObjectCode" }

// Decode implements the Instruction interface.
func (inst *StopObjectCode) Decode(_ vm.OpCode, r *vm.BytecodeReader) (err error) {
	inst.frame, err = r.EndFrame()
	return
}

type Branch struct {
	instruction
	Left  vm.WordPointer
	Right vm.Param
	Goto  vm.ProgramAddress
}

type IsEqual struct{ Branch }

func (inst IsEqual) Mnemonic(st *vm.SymbolTable) string {
	return fmt.Sprintf("unless (%s == %s) goto %s",
		inst.Left.Represent(st, vm.ParamFlagsNone),
		inst.Right.Represent(st, vm.ParamFlagsNone),
		inst.Goto.Represent(st, vm.ParamFlagsNone),
	)
}

func (inst *IsEqual) Decode(opcode vm.OpCode, r *vm.BytecodeReader) (err error) {
	inst.Left = r.ReadWordPointer()
	inst.Right = r.ReadWordParam(opcode, vm.ParamPos1)
	inst.Goto = r.ReadRelativeJump()
	inst.frame, err = r.EndFrame()
	return
}
