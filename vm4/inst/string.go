package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

// LoadString is an instruction that loads a value into a string resource.
type LoadString struct {
	StrID vm.Param `op:"p8" pos:"1" fmt:"id:string"`
	Val   string   `op:"str"`
}

// WriteChar is an instruction that writes a character into a string resource.
type WriteChar struct {
	StrID vm.Param    `op:"p8" pos:"1" fmt:"id:string"`
	Index vm.Param    `op:"p8" pos:"2" fmt:"dec"`
	Val   vm.Constant `op:"c" fmt:"char"`
}

// NewString is an instruction that allocates a new string resource.
type NewString struct {
	StrID vm.Param `op:"p8" pos:"1" fmt:"id:string"`
	Size  vm.Param `op:"p8" pos:"2" fmt:"dec"`
}

func decodeStringOp(opcode vm.OpCode, r *vm.BytecodeDecoder) (inst vm.Instruction, err error) {
	sub := r.DecodeOpCode()
	switch sub & 0x1F {
	case 0x01:
		inst = new(LoadString)
	case 0x03:
		inst = new(WriteChar)
	case 0x05:
		inst = new(NewString)
	default:
		return nil, fmt.Errorf("unknown opcode %02X %02X for string operation", opcode, sub)
	}
	err = vm.DecodeOperands(sub, r, inst)
	return inst, nil
}
