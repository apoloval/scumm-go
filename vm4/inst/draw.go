package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

type DrawBox struct {
	Left   vm.Param `op:"p16" pos:"1" fmt:"dec"`
	Top    vm.Param `op:"p16" pos:"2" fmt:"dec"`
	Right  vm.Param `op:"p16" pos:"1" fmt:"dec"`
	Bottom vm.Param `op:"p16" pos:"2" fmt:"dec"`
	Color  vm.Param `op:"p8" pos:"3" fmt:"dec"`
}

func (inst DrawBox) Acronym() string { return "DRAWBOX" }

func (inst *DrawBox) DecodeOperands(opcode vm.OpCode, r *vm.BytecodeDecoder) error {
	inst.Left = r.DecodeWordParam(opcode, vm.ParamPos1, vm.NumberFormatDecimal)
	inst.Top = r.DecodeWordParam(opcode, vm.ParamPos2, vm.NumberFormatDecimal)
	aux := r.DecodeOpCode()
	inst.Right = r.DecodeWordParam(aux, vm.ParamPos1, vm.NumberFormatDecimal)
	inst.Bottom = r.DecodeWordParam(aux, vm.ParamPos2, vm.NumberFormatDecimal)
	inst.Color = r.DecodeByteParam(aux, vm.ParamPos3, vm.NumberFormatDecimal)
	return nil
}

type DrawObjectAt struct {
	Object vm.Param `op:"p16" pos:"1" fmt:"id:object"`
	XPos   vm.Param `op:"p16" pos:"1" fmt:"dec"` // pos respect the sub-opcode
	YPos   vm.Param `op:"p16" pos:"2" fmt:"dec"` // pos respect the sub-opcode
}

type DrawObjectState struct {
	Object vm.Param `op:"p16" pos:"1" fmt:"id:object"`
	State  vm.Param `op:"p16" pos:"1" fmt:"dec"` // pos respect the sub-opcode
}

type DrawObject struct {
	Object vm.Param `op:"p16" pos:"1" fmt:"id:object"`
}

func decodeDrawObjectOp(opcode vm.OpCode, r *vm.BytecodeDecoder) (inst vm.Instruction, err error) {
	obj := r.DecodeWordParam(opcode, vm.ParamPos1, vm.NumberFormatDecimal)
	sub := r.DecodeOpCode()
	switch sub {
	case 0x01:
		inst = &DrawObjectAt{
			Object: obj,
			XPos:   r.DecodeWordParam(sub, vm.ParamPos1, vm.NumberFormatDecimal),
			YPos:   r.DecodeWordParam(sub, vm.ParamPos2, vm.NumberFormatDecimal),
		}
	case 0x02:
		inst = &DrawObjectState{
			Object: obj,
			State:  r.DecodeWordParam(sub, vm.ParamPos1, vm.NumberFormatDecimal),
		}
	case 0xFF:
		inst = &DrawObject{Object: obj}
	default:
		err = fmt.Errorf("invalid sub-opcode %02X %02X for draw object operation", opcode, sub)
	}
	return
}
