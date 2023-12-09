package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

type SetBoxFlags struct {
	Box   vm.Param `op:"p8" pos:"1" fmt:"dec"`
	Value vm.Param `op:"p8" pos:"2" fmt:"dec"`
}

func (inst SetBoxFlags) Acronym() string { return "SETBOXF" }

type SetBoxScale struct {
	Box   vm.Param `op:"p8" pos:"1" fmt:"dec"`
	Value vm.Param `op:"p8" pos:"2" fmt:"dec"`
}

func (inst SetBoxScale) Acronym() string { return "SETBOXSCL" }

type SetBoxScaleInv struct {
	Box   vm.Param `op:"p8" pos:"1" fmt:"dec"`
	Value vm.Param `op:"p8" pos:"2" fmt:"dec"`
}

func (inst SetBoxScaleInv) Acronym() string { return "SETBOXSCLI" }

type CreateBoxMatrix struct{}

func (inst CreateBoxMatrix) Acronym() string { return "INITBOXMTX" }

func decodeBoxOp(opcode vm.OpCode, r *vm.BytecodeDecoder) (inst vm.Instruction, err error) {
	sub := r.DecodeOpCode()
	switch sub & 0x1F {
	case 0x01:
		inst = &SetBoxFlags{
			Box:   r.DecodeByteParam(sub, vm.ParamPos1, vm.NumberFormatDecimal),
			Value: r.DecodeByteParam(sub, vm.ParamPos2, vm.NumberFormatDecimal),
		}
	case 0x02:
		inst = &SetBoxScale{
			Box:   r.DecodeByteParam(sub, vm.ParamPos1, vm.NumberFormatDecimal),
			Value: r.DecodeByteParam(sub, vm.ParamPos2, vm.NumberFormatDecimal),
		}
	case 0x03:
		inst = &SetBoxScaleInv{
			Box:   r.DecodeByteParam(sub, vm.ParamPos1, vm.NumberFormatDecimal),
			Value: r.DecodeByteParam(sub, vm.ParamPos2, vm.NumberFormatDecimal),
		}
	case 0x04:
		inst = new(CreateBoxMatrix)
	default:
		err = fmt.Errorf("unknown box opcode: %02X %02X", opcode, sub)
	}
	return
}
