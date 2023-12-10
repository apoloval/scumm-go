package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

type BeginOverride struct {
	Target vm.Constant `op:"reljmp" fmt:"addr"`
}

func (inst BeginOverride) Acronym() string { return "BEGOVER" }

type EndOverride struct{}

func (inst EndOverride) Acronym() string { return "ENDOVER" }

func decodeOverrideOp(opcode vm.OpCode, r *vm.BytecodeDecoder) (inst vm.Instruction, err error) {
	sub := r.DecodeOpCode()
	switch sub & 0x1F {
	case 0x00:
		inst = new(EndOverride)
	case 0x01:
		sub2 := r.DecodeOpCode()
		if sub2 != 0x18 {
			err = fmt.Errorf("unknown override opcode: %02X %02X %02X", opcode, sub, sub2)
			return
		}
		inst = &BeginOverride{
			Target: r.DecodeRelativeJump(),
		}
	default:
		err = fmt.Errorf("unknown override opcode: %02X %02X", opcode, sub)
	}
	return
}
