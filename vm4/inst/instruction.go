package inst

import "github.com/apoloval/scumm-go/vm"

type instruction struct {
	frame vm.BytecodeFrame
}

func (inst instruction) Frame() vm.BytecodeFrame {
	return inst.frame
}

func (inst instruction) Params() vm.Params { return nil }
