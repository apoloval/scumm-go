package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

// LoadString is an instruction that loads a value into a string resource.
type LoadString struct {
	base
	StrID vm.Param
	Val   string
}

// Mnemonic implements the Instruction interface.
func (inst LoadString) Mnemonic(st *vm.SymbolTable) string {
	return fmt.Sprintf("LoadString %s, %q",
		inst.StrID.Display(st),
		inst.Val,
	)
}

// Decode implements the Instruction interface.
func (inst *LoadString) Decode(opcode vm.OpCode, r *vm.BytecodeReader) error {
	inst.StrID = r.ReadByteParam(opcode, vm.ParamPos1, vm.ParamFormatStringID)
	inst.Val = r.ReadString()
	return inst.base.Decode(opcode, r)
}

// WriteChar is an instruction that writes a character into a string resource.
type WriteChar struct {
	base
	StrID vm.Param
	Index vm.Param
	Val   vm.Param
}

// Mnemonic implements the Instruction interface.
func (inst WriteChar) Mnemonic(st *vm.SymbolTable) string {
	return fmt.Sprintf("WriteChar %s, %s, %s",
		inst.StrID.Display(st),
		inst.Index.Display(st),
		inst.Val.Display(st),
	)
}

// Decode implements the Instruction interface.
func (inst *WriteChar) Decode(opcode vm.OpCode, r *vm.BytecodeReader) error {
	inst.StrID = r.ReadByteParam(opcode, vm.ParamPos1, vm.ParamFormatStringID)
	inst.Index = r.ReadByteParam(opcode, vm.ParamPos2, vm.ParamFormatNumber)
	inst.Val = r.ReadByteParam(opcode, vm.ParamPos3, vm.ParamFormatChar)
	return inst.base.Decode(opcode, r)
}

// NewString is an instruction that allocates a new string resource.
type NewString struct {
	base
	StrID vm.Param
	Size  vm.Param
}

// Mnemonic implements the Instruction interface.
func (inst NewString) Mnemonic(st *vm.SymbolTable) string {
	return fmt.Sprintf("NewString %s, %s",
		inst.StrID.Display(st),
		inst.Size.Display(st),
	)
}

// Decode implements the Instruction interface.
func (inst *NewString) Decode(opcode vm.OpCode, r *vm.BytecodeReader) error {
	inst.StrID = r.ReadByteParam(opcode, vm.ParamPos1, vm.ParamFormatStringID)
	inst.Size = r.ReadByteParam(opcode, vm.ParamPos2, vm.ParamFormatNumber)
	return inst.base.Decode(opcode, r)
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
