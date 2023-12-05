package inst

import (
	"fmt"
	"strings"

	"github.com/apoloval/scumm-go/vm"
)

// StopObjectCode is a stop instruction that stops the execution of the current script.
type StopObjectCode struct{ base }

// Goto is a goto instruction that jumps to the given address.
type Goto struct{ branch }

type IsEqual struct{ binaryBranch }
type IsNotEqual struct{ binaryBranch }
type IsLess struct{ binaryBranch }
type IsLessEqual struct{ binaryBranch }
type IsGreater struct{ binaryBranch }
type IsGreaterEqual struct{ binaryBranch }

type IsEqualZero struct{ unaryBranch }
type IsNotEqualZero struct{ unaryBranch }

type branch struct {
	base
	Goto vm.Constant
}

func withBranch(name string) branch {
	return branch{base: withName(name)}
}

func (b *branch) Decode(opcode vm.OpCode, r *vm.BytecodeReader) error {
	return b.decodeWithParams(r)
}

func (b *branch) decodeWithParams(r *vm.BytecodeReader, params ...vm.Param) error {
	b.Goto = r.ReadRelativeJump()
	return b.base.decodeWithParams(r, append([]vm.Param{b.Goto}, params...)...)
}

type unaryBranch struct {
	branch
	Var vm.Pointer
}

func withUnaryBranch(name string) unaryBranch {
	return unaryBranch{branch: withBranch(name)}
}

func (inst unaryBranch) Mnemonic(st *vm.SymbolTable) string {
	return fmt.Sprintf("Unless (%s %s) Goto %s",
		inst.Var.Display(st),
		inst.name,
		inst.Goto.Display(st),
	)
}

func (b *unaryBranch) Decode(opcode vm.OpCode, r *vm.BytecodeReader) error {
	return b.decodeWithParams(r)
}

func (b *unaryBranch) decodeWithParams(r *vm.BytecodeReader, params ...vm.Param) error {
	b.Var = r.ReadPointer()
	return b.branch.decodeWithParams(r, append([]vm.Param{b.Var}, params...)...)
}

type binaryBranch struct {
	branch
	Var   vm.Pointer
	Value vm.Param
}

func withBinaryBranch(name string) binaryBranch {
	return binaryBranch{branch: withBranch(name)}
}

func (inst binaryBranch) Mnemonic(st *vm.SymbolTable) string {
	return fmt.Sprintf("Unless (%s %s %s) Goto %s",
		inst.Value.Display(st),
		inst.name,
		inst.Var.Display(st),
		inst.Goto.Display(st),
	)
}

func (b *binaryBranch) Decode(opcode vm.OpCode, r *vm.BytecodeReader) error {
	return b.decodeWithParams(opcode, r)
}

func (b *binaryBranch) decodeWithParams(
	opcode vm.OpCode,
	r *vm.BytecodeReader,
	params ...vm.Param,
) error {
	b.Var = r.ReadPointer()
	b.Value = r.ReadWordParam(opcode, vm.ParamPos1, vm.NumberFormatDecimal)
	return b.branch.decodeWithParams(r, append([]vm.Param{b.Var, b.Value}, params...)...)
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
