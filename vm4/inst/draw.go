package inst

import "github.com/apoloval/scumm-go/vm"

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
