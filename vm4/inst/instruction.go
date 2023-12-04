package inst

import "github.com/apoloval/scumm-go/vm"

type base struct {
	frame vm.BytecodeFrame

	operand []vm.Param
	format  []vm.NumberFormat
}

func (b *base) Decode(opcode vm.OpCode, r *vm.BytecodeReader) (err error) {
	b.frame, err = r.EndFrame()
	return
}

func (b base) Frame() vm.BytecodeFrame {
	return b.frame
}

func (b base) Params() vm.Params { return nil }
