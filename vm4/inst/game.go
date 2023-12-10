package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

type Game struct {
	Result vm.VarRef `op:"result"`
	Arg    vm.Param  `op:"p8" pos:"1" fmt:"decimal"`
}

func (inst Game) Mnemonic() string { return "GAME" }

func (inst Game) DisplayOperands(st *vm.SymbolTable) (ops []string) {
	// It is bizarre, but by design this instruction could have the action encoded in a variable.
	// Probably it is never used like that, but we have to handle it.
	if c, ok := inst.Arg.(vm.Constant); ok {
		switch c.Value & 0xE0 {
		case 0x00:
			ops = []string{"SLOTS"}
		case 0x20:
			ops = []string{"DRIVE"}
		case 0x40:
			ops = []string{"LOAD", fmt.Sprintf("%d", c.Value&0x1F)}
		case 0x80:
			ops = []string{"SAVE", fmt.Sprintf("%d", c.Value&0x1F)}
		case 0xC0:
			ops = []string{"TEST", fmt.Sprintf("%d", c.Value&0x1F)}
		}
		return
	}
	ops = []string{"?", inst.Arg.Display(st)}
	return
}
