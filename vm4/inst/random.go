package inst

import "github.com/apoloval/scumm-go/vm"

type GetRandomNumber struct {
	Result vm.VarRef `op:"result"`
	Seed   vm.Param  `op:"p8" pos:"1" fmt:"hex"`
}

func (inst GetRandomNumber) Acronym() string { return "RAND" }
