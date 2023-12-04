package inst

import (
	"strings"

	"github.com/apoloval/scumm-go/vm"
)

type base struct {
	frame vm.BytecodeFrame

	name   string
	params vm.Params
}

func withName(name string) base {
	return base{name: name}
}

func (b *base) Decode(opcode vm.OpCode, r *vm.BytecodeReader) error {
	return b.decodeFrame(r)
}

func (b *base) Mnemonic(st *vm.SymbolTable) string {
	params := make([]string, 0, len(b.params))
	for _, p := range b.params {
		params = append(params, p.Display(st))
	}
	return b.name + " " + strings.Join(params, ", ")
}

func (b *base) decodeWithParams(r *vm.BytecodeReader, params ...vm.Param) error {
	b.params = vm.Params(params)
	return b.decodeFrame(r)
}

func (b *base) decodeFrame(r *vm.BytecodeReader) (err error) {
	b.frame, err = r.EndFrame()
	return
}

func (b base) Frame() vm.BytecodeFrame {
	return b.frame
}

func (b base) Params() vm.Params { return b.params }
