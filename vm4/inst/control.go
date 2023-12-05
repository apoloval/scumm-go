package inst

import (
	"fmt"
	"strings"

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
	inst.Right = r.ReadWordParam(opcode, vm.ParamPos1, vm.NumberFormatDecimal)
	inst.Goto = r.ReadRelativeJump()
	return inst.base.Decode(opcode, r)

}

// StartScript is a instruction that starts a new script in a new thread.
type StartScript struct {
	base
	ScriptID vm.Param
	Args     vm.Params

	Recursive       bool
	FreezeResistant bool
}

func (inst StartScript) Mnemonic(st *vm.SymbolTable) string {
	var flags []string
	if inst.Recursive {
		flags = append(flags, "recursive")
	}
	if inst.FreezeResistant {
		flags = append(flags, "freeze-resistant")
	}
	return fmt.Sprintf("StartScript %s(%s) %s",
		inst.ScriptID.Display(st),
		inst.Params().Display(st),
		strings.Join(flags, ", "),
	)
}

func (inst *StartScript) Decode(opcode vm.OpCode, r *vm.BytecodeReader) error {
	inst.Recursive = opcode&0x40 > 0
	inst.FreezeResistant = opcode&0x20 > 0
	inst.ScriptID = r.ReadByteParam(opcode, vm.ParamPos1, vm.NumberFormatScriptID)
	inst.Args = r.ReadVarParams()
	return inst.base.Decode(opcode, r)
}
