package inst

import (
	"fmt"
	"io"

	"github.com/apoloval/scumm-go/vm"
)

// DecodeAll decodes all instructions from the bytecode reader.
func DecodeAll(r *vm.BytecodeReader) (code []vm.Instruction, err error) {
	for {
		var i vm.Instruction
		i, err = Decode(r)
		if err == io.EOF {
			return code, nil
		}
		if err != nil {
			return code, err
		}
		code = append(code, i)
	}
}

// Decode decodes an instruction from the bytecode reader.
func Decode(r *vm.BytecodeReader) (vm.Instruction, error) {
	r.BeginFrame()
	opcode := r.ReadOpCode()

	var inst vm.Instruction
	switch opcode {
	case 0x00:
		inst = &StopObjectCode{}
	case 0x1A:
		inst = &Move{}
	case 0x2C:
		return decodeCursorCommand(opcode, r)
	case 0x48:
		inst = &IsEqual{}
	default:
		return nil, fmt.Errorf("unknown opcode %02X", opcode)
	}
	if err := inst.Decode(opcode, r); err != nil {
		return nil, err
	}
	return inst, nil
}
