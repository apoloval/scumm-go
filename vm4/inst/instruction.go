package inst

import "github.com/apoloval/scumm-go/vm"

type instruction struct {
	bytecode []byte
}

func (inst instruction) Bytecode() []byte {
	return inst.bytecode
}

func (inst instruction) Params() vm.Params { return nil }
