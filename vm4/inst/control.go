package inst

import (
	"github.com/apoloval/scumm-go/vm"
)

// StopObjectCode is a stop instruction that stops the execution of the current script after reaching the
// end of the code. This is also known as StopObjectCode in ScummVM.
type StopObjectCode struct{}

func (inst StopObjectCode) Acronym() string { return "SOC" }

// Jump is a instruction that jumps to the given address. This is also known as JumpRelative in
// ScummVM.
type Jump struct {
	Target vm.Constant `op:"reljmp" fmt:"addr"`
}

func (inst Jump) Acronym() string { return "JMP" }

type UnaryBranch struct {
	Var    vm.VarRef   `op:"var"`
	Target vm.Constant `op:"reljmp" fmt:"addr"`
}

type BinaryBranch struct {
	Var    vm.VarRef   `op:"var"`
	Value  vm.Param    `op:"p16" pos:"1" fmt:"dec"`
	Target vm.Constant `op:"reljmp" fmt:"addr"`
}

type BranchUnlessEqual BinaryBranch

func (inst BranchUnlessEqual) Acronym() string { return "BREQ" }

type BranchUnlessNotEqual BinaryBranch

func (inst BranchUnlessNotEqual) Acronym() string { return "BRNE" }

type BranchUnlessLess BinaryBranch

func (inst BranchUnlessLess) Acronym() string { return "BRLT" }

type BranchUnlessLessEqual BinaryBranch

func (inst BranchUnlessLessEqual) Acronym() string { return "BRLE" }

type BranchUnlessGreater BinaryBranch

func (inst BranchUnlessGreater) Acronym() string { return "BRGT" }

type BranchUnlessGreaterEqual BinaryBranch

func (inst BranchUnlessGreaterEqual) Acronym() string { return "BRGE" }

type BranchUnlessZero UnaryBranch

func (inst BranchUnlessZero) Acronym() string { return "BRZE" }

type BranchUnlessNotZero UnaryBranch

func (inst BranchUnlessNotZero) Acronym() string { return "BRNZ" }

// StartScript is a instruction that starts a new script in a new thread.
type StartScript struct {
	ScriptID vm.Param  `op:"p8"`
	Args     vm.Params `op:"v16"`

	Recursive       bool
	FreezeResistant bool
}

func (inst *StartScript) DecodeOperands(opcode vm.OpCode, r *vm.BytecodeDecoder) error {
	inst.ScriptID = r.DecodeByteParam(opcode, vm.ParamPos1, vm.NumberFormatScriptID)
	inst.Args = r.DecodeVarParams()
	inst.Recursive = opcode&0x40 > 0
	inst.FreezeResistant = opcode&0x20 > 0
	return nil
}

func (inst StartScript) Acronym() string { return "STRSC" }

func (inst StartScript) DisplayOperands(st *vm.SymbolTable) (ops []string) {
	var flags string
	if inst.Recursive {
		flags += "R"
	}
	if inst.FreezeResistant {
		flags += "F"
	}
	ops = []string{
		inst.ScriptID.Display(st),
		inst.Args.Display(st),
	}
	if flags != "" {
		ops = append(ops, flags)
	}
	return
}

type StopScript struct {
	Script vm.Param `op:"p8" pos:"1" fmt:"id:script"`
}

func (inst StopScript) Acronym() string { return "STPSC" }

type ChainStript struct {
	Script vm.Param  `op:"p8" pos:"1" fmt:"id:script"`
	Args   vm.Params `op:"v16"`
}

func (inst ChainStript) Acronym() string { return "CHNSC" }

// StartObject is a instruction that starts a object script.
type StartObject struct {
	Object vm.Param  `op:"p16" pos:"1" fmt:"id:object"`
	Script vm.Param  `op:"p8" pos:"2" fmt:"id:script"`
	Args   vm.Params `op:"v16"`
}

func (inst StartObject) Acronym() string { return "STOB" }

type BreakHere struct{}

func (inst BreakHere) Acronym() string { return "BREAK" }

// ScriptRunning is a instruction that checks if a script is running. It is also known as
// IsScriptRunning in ScummVM.
type ScriptRunning struct {
	Result   vm.VarRef `op:"result"`
	ScriptID vm.Param  `op:"p8" pos:"1" fmt:"id:script"`
}

func (inst ScriptRunning) Acronym() string { return "SCRUN" }

// LoadRoom is a instruction that loads a new room.
type LoadRoom struct {
	RoomID vm.Param `op:"p8" pos:"1" fmt:"id:room"`
}

func (inst LoadRoom) Acronym() string { return "LDRO" }

type BranchUnlessState struct {
	Object vm.Param    `op:"p16" pos:"1" fmt:"id:object"`
	State  vm.Param    `op:"p8" pos:"2" fmt:"id:state"`
	Target vm.Constant `op:"reljmp" fmt:"addr"`
}

func (inst BranchUnlessState) Acronym() string { return "BRST" }

type BranchUnlessActorInBox struct {
	Actor  vm.Param    `op:"p8" pos:"1" fmt:"dec"`
	Box    vm.Param    `op:"p8" pos:"2" fmt:"dec"`
	Target vm.Constant `op:"reljmp" fmt:"addr"`
}

func (inst BranchUnlessActorInBox) Acronym() string { return "BRAB" }

type BranchUnlessClass struct {
	Object  vm.Param    `op:"p16" pos:"1" fmt:"id:object"`
	Classes vm.Params   `op:"v16"`
	Target  vm.Constant `op:"reljmp" fmt:"addr"`
}

func (inst BranchUnlessClass) Acronym() string { return "BRCL" }
