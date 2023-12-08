package inst

import "github.com/apoloval/scumm-go/vm"

type FreezeScripts struct {
	Flag vm.Param `op:"p8" pos:"1" fmt:"dec"`
}

func (inst FreezeScripts) Acronym() string { return "FREEZE" }

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

// ScriptRunning is a instruction that checks if a script is running. It is also known as
// IsScriptRunning in ScummVM.
type ScriptRunning struct {
	Result   vm.VarRef `op:"result"`
	ScriptID vm.Param  `op:"p8" pos:"1" fmt:"id:script"`
}

func (inst ScriptRunning) Acronym() string { return "SCRUN" }
