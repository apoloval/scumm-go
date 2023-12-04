package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

// LoadString is an instruction that loads a value into a string resource.
type LoadString struct {
	instruction
	StrID vm.Param
	Val   string
}

// Mnemonic implements the Instruction interface.
func (inst LoadString) Mnemonic(st *vm.SymbolTable) string {
	return fmt.Sprintf("LoadString %s, %q",
		inst.StrID.Display(st, vm.ParamFormatStringID),
		inst.Val,
	)
}

// Decode implements the Instruction interface.
func (inst *LoadString) Decode(opcode vm.OpCode, r *vm.BytecodeReader) (err error) {
	inst.StrID = r.ReadByteParam(opcode, vm.ParamPos1)
	inst.Val = r.ReadString()
	inst.frame, err = r.EndFrame()
	return
}

// WriteChar is an instruction that writes a character into a string resource.
type WriteChar struct {
	instruction
	StrID vm.Param
	Index vm.Param
	Val   vm.Param
}

// Mnemonic implements the Instruction interface.
func (inst WriteChar) Mnemonic(st *vm.SymbolTable) string {
	return fmt.Sprintf("WriteChar %s, %s, %s",
		inst.StrID.Display(st, vm.ParamFormatStringID),
		inst.Index.Display(st, vm.ParamFormatNumber),
		inst.Val.Display(st, vm.ParamFormatChar),
	)
}

// Decode implements the Instruction interface.
func (inst *WriteChar) Decode(opcode vm.OpCode, r *vm.BytecodeReader) (err error) {
	inst.StrID = r.ReadByteParam(opcode, vm.ParamPos1)
	inst.Index = r.ReadByteParam(opcode, vm.ParamPos2)
	inst.Val = r.ReadByteParam(opcode, vm.ParamPos3)
	inst.frame, err = r.EndFrame()
	return
}

// NewString is an instruction that allocates a new string resource.
type NewString struct {
	instruction
	StrID vm.Param
	Size  vm.Param
}

// Mnemonic implements the Instruction interface.
func (inst NewString) Mnemonic(st *vm.SymbolTable) string {
	return fmt.Sprintf("NewString %s, %s",
		inst.StrID.Display(st, vm.ParamFormatStringID),
		inst.Size.Display(st, vm.ParamFormatNumber),
	)
}

// Decode implements the Instruction interface.
func (inst *NewString) Decode(opcode vm.OpCode, r *vm.BytecodeReader) (err error) {
	inst.StrID = r.ReadByteParam(opcode, vm.ParamPos1)
	inst.Size = r.ReadByteParam(opcode, vm.ParamPos2)
	inst.frame, err = r.EndFrame()
	return
}

func decodeStringOp(opcode vm.OpCode, r *vm.BytecodeReader) (inst vm.Instruction, err error) {
	sub := r.ReadOpCode()

	switch sub & 0x1F {
	case 0x01:
		inst = &LoadString{}
	case 0x03:
		inst = &WriteChar{}
	case 0x05:
		inst = &NewString{}
	default:
		return nil, fmt.Errorf("unknown opcode %02X %02X for string operation", opcode, sub)
	}
	err = inst.Decode(sub, r)
	return inst, nil
}
