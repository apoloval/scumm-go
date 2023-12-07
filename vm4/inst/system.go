package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

// Reset is a instruction to reset the game. It is also known as Restart in ScummVM.
type Reset struct{}

func (inst Reset) Acronym() string { return "RESET" }

// Pause is a instruction to pause the game.
type Pause struct{}

func (inst Pause) Acronym() string { return "PAUSE" }

// Quit is a instruction to quit the game.
type Quit struct{}

func (inst Quit) Acronym() string { return "QUIT" }

func decodeSystemOp(opcode vm.OpCode, r *vm.BytecodeDecoder) (inst vm.Instruction, err error) {
	sub := r.DecodeOpCode()
	switch sub & 0x0F {
	case 0x01:
		inst = new(Reset)
	case 0x02:
		inst = new(Pause)
	case 0x03:
		inst = new(Quit)
	default:
		err = fmt.Errorf("unknown sub-opcode %02X %02X for system operation", opcode, sub)
	}
	return
}
