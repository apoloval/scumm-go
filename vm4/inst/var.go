package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

// Move is a move instruction that puts the value from Src into Dest
type Move struct {
	Result vm.VarRef `op:"result"`
	Src    vm.Param  `op:"p16" pos:"1" fmt:"dec"`
}

func (inst Move) Acronym() string { return "MOVE" }

func (inst Move) Execute(ctx vm.ExecutionContext) {
	inst.Result.Write(ctx, inst.Src.Evaluate(ctx))
}

// SetVarRange is a instruction that sets a range of variables to the given values.
type SetVarRange struct {
	Result vm.VarRef     `op:"result"`
	Values []vm.Constant `op:"16"`
}

func (inst SetVarRange) Acronym() string { return "SETVR" }

func (inst SetVarRange) DisplayOperands(st *vm.SymbolTable) []string {
	ops := []string{
		inst.Result.Display(st),
	}
	for _, val := range inst.Values {
		ops = append(ops, val.Display(st))
	}
	return ops
}

// Decode implements the Instruction interface.
func (inst *SetVarRange) DecodeOperands(opcode vm.OpCode, r *vm.BytecodeDecoder) error {
	inst.Result = r.DecodeVarRef()
	count := r.DecodeByte()
	for i := 0; i < int(count); i++ {
		if opcode&0x80 > 0 {
			inst.Values = append(inst.Values, r.DecodeWordConstant(vm.NumberFormatDecimal))
		} else {
			inst.Values = append(inst.Values, r.DecodeByteConstant(vm.NumberFormatDecimal))
		}
	}
	return nil
}

type SaveVarsWriteVars struct {
	ResultA vm.VarRef `op:"result"`
	ResultB vm.VarRef `op:"result"`
}

func (inst SaveVarsWriteVars) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("VARS=[%s, %s]", inst.ResultA.Display(st), inst.ResultB.Display(st))
}

type SaveVarsWriteStrings struct {
	Arg1 vm.Param `op:"p8" pos:"1" fmt:"id:string"`
	Arg2 vm.Param `op:"p8" pos:"2" fmt:"id:string"`
}

func (inst SaveVarsWriteStrings) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("STRINGS=[%s, %s]", inst.Arg1.Display(st), inst.Arg2.Display(st))
}

type SaveVarsOpenFile struct {
	Filename string `op:"string"`
}

func (inst SaveVarsOpenFile) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("OPEN(%s)", inst.Filename)
}

type SaveVarsDummy struct{}

func (inst SaveVarsDummy) Display(st *vm.SymbolTable) string {
	return "DUMMY"
}

type SaveVarsCloseFile struct{}

func (inst SaveVarsCloseFile) Display(st *vm.SymbolTable) string {
	return "CLOSE"
}

type SaveVars struct {
	WriteVars    *SaveVarsWriteVars    // 0x01
	WriteStrings *SaveVarsWriteStrings // 0x02
	OpenFile     *SaveVarsOpenFile     // 0x03
	Dummy        *SaveVarsDummy        // 0x04
	CloseFile    *SaveVarsCloseFile    // 0x1F
}

func (inst SaveVars) Acronym() string { return "SAVEVARS" }

func (inst SaveVars) DisplayOperands(st *vm.SymbolTable) []string {
	var props []string
	if inst.WriteVars != nil {
		props = append(props, inst.WriteVars.Display(st))
	}
	if inst.WriteStrings != nil {
		props = append(props, inst.WriteStrings.Display(st))
	}
	if inst.OpenFile != nil {
		props = append(props, inst.OpenFile.Display(st))
	}
	if inst.Dummy != nil {
		props = append(props, inst.Dummy.Display(st))
	}
	if inst.CloseFile != nil {
		props = append(props, inst.CloseFile.Display(st))
	}
	return props
}

func (inst *SaveVars) DecodeOperands(opcode vm.OpCode, r *vm.BytecodeDecoder) error {
	for {
		sub := r.DecodeOpCode()
		switch sub {
		case 0x00:
			return nil
		case 0x01:
			inst.WriteVars = &SaveVarsWriteVars{
				ResultA: r.DecodeVarRef(),
				ResultB: r.DecodeVarRef(),
			}
		case 0x02:
			inst.WriteStrings = &SaveVarsWriteStrings{
				Arg1: r.DecodeByteParam(sub, vm.ParamPos1, vm.NumberFormatStringID),
				Arg2: r.DecodeByteParam(sub, vm.ParamPos2, vm.NumberFormatStringID),
			}
		case 0x03:
			inst.OpenFile = &SaveVarsOpenFile{
				Filename: r.DecodeString(),
			}
		case 0x04:
			inst.Dummy = new(SaveVarsDummy)
		case 0x1F:
			inst.CloseFile = new(SaveVarsCloseFile)
		default:
			return fmt.Errorf("unknown opcode %02X %02X in save vars op", sub, opcode)
		}
		return nil
	}
}

type LoadVarsReadVars struct {
	ResultA vm.VarRef `op:"result"`
	ResultB vm.VarRef `op:"result"`
}

func (inst LoadVarsReadVars) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("VARS=[%s, %s]", inst.ResultA.Display(st), inst.ResultB.Display(st))
}

type LoadVarsReadStrings struct {
	Arg1 vm.Param `op:"p8" pos:"1" fmt:"id:string"`
	Arg2 vm.Param `op:"p8" pos:"2" fmt:"id:string"`
}

func (inst LoadVarsReadStrings) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("STRINGS=[%s, %s]", inst.Arg1.Display(st), inst.Arg2.Display(st))
}

type LoadVarsOpenFile struct {
	Filename string `op:"string"`
}

func (inst LoadVarsOpenFile) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("OPEN(%s)", inst.Filename)
}

type LoadVarsDummy struct{}

func (inst LoadVarsDummy) Display(st *vm.SymbolTable) string {
	return "DUMMY"
}

type LoadVarsCloseFile struct{}

func (inst LoadVarsCloseFile) Display(st *vm.SymbolTable) string {
	return "CLOSE"
}

type LoadVars struct {
	ReadVars    *LoadVarsReadVars    // 0x01
	ReadStrings *LoadVarsReadStrings // 0x02
	OpenFile    *LoadVarsOpenFile    // 0x03
	Dummy       *LoadVarsDummy       // 0x04
	CloseFile   *LoadVarsCloseFile   // 0x1F
}

func (inst LoadVars) Acronym() string { return "LOADVARS" }

func (inst LoadVars) DisplayOperands(st *vm.SymbolTable) []string {
	var props []string
	if inst.ReadVars != nil {
		props = append(props, inst.ReadVars.Display(st))
	}
	if inst.ReadStrings != nil {
		props = append(props, inst.ReadStrings.Display(st))
	}
	if inst.OpenFile != nil {
		props = append(props, inst.OpenFile.Display(st))
	}
	if inst.Dummy != nil {
		props = append(props, inst.Dummy.Display(st))
	}
	if inst.CloseFile != nil {
		props = append(props, inst.CloseFile.Display(st))
	}
	return props
}

func (inst *LoadVars) DecodeOperands(opcode vm.OpCode, r *vm.BytecodeDecoder) error {
	for {
		sub := r.DecodeOpCode()
		switch sub {
		case 0x00:
			return nil
		case 0x01:
			inst.ReadVars = &LoadVarsReadVars{
				ResultA: r.DecodeVarRef(),
				ResultB: r.DecodeVarRef(),
			}
		case 0x02:
			inst.ReadStrings = &LoadVarsReadStrings{
				Arg1: r.DecodeByteParam(sub, vm.ParamPos1, vm.NumberFormatStringID),
				Arg2: r.DecodeByteParam(sub, vm.ParamPos2, vm.NumberFormatStringID),
			}
		case 0x03:
			inst.OpenFile = &LoadVarsOpenFile{
				Filename: r.DecodeString(),
			}
		case 0x04:
			inst.Dummy = new(LoadVarsDummy)
		case 0x1F:
			inst.CloseFile = new(LoadVarsCloseFile)
		default:
			return fmt.Errorf("unknown opcode %02X %02X in load vars op", sub, opcode)
		}
		return nil
	}
}

func decodeSaveLoadVarsOp(opcode vm.OpCode, r *vm.BytecodeDecoder) (inst vm.Instruction, err error) {
	sub := r.DecodeOpCode()
	switch sub {
	case 0x01:
		inst = new(SaveVars)
	case 0x02:
		inst = new(LoadVars)
	default:
		return nil, fmt.Errorf("unknown opcode %02X %02X in save/load vars op", sub, opcode)
	}
	err = vm.DecodeOperands(opcode, r, inst)
	return
}
