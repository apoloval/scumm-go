package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

// Decode decodes an instruction from the bytecode reader.
func Decode(r *vm.BytecodeReader) (vm.Instruction, error) {
	r.BeginFrame()
	opcode := r.ReadOpCode()

	switch opcode {
	case 0x2C:
		return decodeCursorCommand(opcode, r)
	default:
		return nil, fmt.Errorf("unknown opcode %02X", opcode)
	}
}
