package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

// Move is a move instruction that puts the value from Src into Dest
type Move struct {
	instruction
	Dest vm.Pointer
	Src  vm.Param
}

// Mnemonic implements the Instruction interface.
func (inst Move) Mnemonic(st *vm.SymbolTable) string {
	return fmt.Sprintf("%s = %s",
		inst.Dest.Represent(st, vm.ParamFlagsNone),
		inst.Src.Represent(st, vm.ParamFlagsNone),
	)
}

// Decode implements the Instruction interface.
func (inst *Move) Decode(opcode vm.OpCode, r *vm.BytecodeReader) (err error) {
	inst.Dest = r.ReadPointer()
	inst.Src = r.ReadWordParam(opcode, vm.ParamPos1)
	inst.frame, err = r.EndFrame()
	return
}
