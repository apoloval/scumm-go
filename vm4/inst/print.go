package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

type PrintColorSet struct {
	Actor vm.Param `op:"p8" pos:"1" fmt:"dec"`
	Color vm.Param `op:"p8" pos:"2" fmt:"dec"`
}

func decodePrintOp(opcode vm.OpCode, r *vm.BytecodeReader) (inst vm.Instruction, err error) {
	actor := r.ReadByteParam(opcode, vm.ParamPos1, vm.NumberFormatDecimal)
	sub := r.ReadOpCode()
	switch sub & 0x1F {
	case 0x01:
		inst = &PrintColorSet{
			Actor: actor,
			Color: r.ReadByteParam(sub, vm.ParamPos1, vm.NumberFormatDecimal),
		}
	default:
		return nil, fmt.Errorf("unknown opcode %02X %02X for print op operation", opcode, sub)
	}
	err = vm.DecodeOperands(sub, r, inst)
	return inst, nil
}
