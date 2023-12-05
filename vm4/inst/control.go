package inst

import (
	"fmt"
	"strings"

	"github.com/apoloval/scumm-go/vm"
)

// StopObjectCode is a stop instruction that stops the execution of the current script.
type StopObjectCode struct{}

// Goto is a goto instruction that jumps to the given address.
type Goto struct {
	Target vm.Constant `op:"reljmp" fmt:"addr"`
}

type UnaryBranch struct {
	Var    vm.Pointer  `op:"var"`
	Target vm.Constant `op:"reljmp" fmt:"addr"`
}

type BinaryBranch struct {
	Var    vm.Pointer  `op:"var"`
	Value  vm.Param    `op:"p16" pos:"1"`
	Target vm.Constant `op:"reljmp" fmt:"addr"`
}

type IsEqual BinaryBranch
type IsNotEqual BinaryBranch
type IsLess BinaryBranch
type IsLessEqual BinaryBranch
type IsGreater BinaryBranch
type IsGreaterEqual BinaryBranch

type IsEqualZero UnaryBranch
type IsNotEqualZero UnaryBranch

func (inst IsEqual) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("Unless (%s == %s) Goto %s",
		inst.Var.Display(st),
		inst.Value.Display(st),
		inst.Target.Display(st),
	)
}

// StartScript is a instruction that starts a new script in a new thread.
type StartScript struct {
	ScriptID vm.Param  `op:"p8"`
	Args     vm.Params `op:"v16"`

	Recursive       bool
	FreezeResistant bool
}

func (inst *StartScript) DecodeOperands(opcode vm.OpCode, r *vm.BytecodeReader) error {
	inst.ScriptID = r.ReadByteParam(opcode, vm.ParamPos1, vm.NumberFormatScriptID)
	inst.Args = r.ReadVarParams()
	inst.Recursive = opcode&0x40 > 0
	inst.FreezeResistant = opcode&0x20 > 0
	return nil
}

func (inst StartScript) Display(st *vm.SymbolTable) string {
	var flags []string
	if inst.Recursive {
		flags = append(flags, "recursive")
	}
	if inst.FreezeResistant {
		flags = append(flags, "freeze-resistant")
	}
	return fmt.Sprintf("StartScript %s(%s) %s",
		inst.ScriptID.Display(st),
		inst.Args.Display(st),
		strings.Join(flags, ", "),
	)
}
