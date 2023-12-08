package inst

import "github.com/apoloval/scumm-go/vm"

type GetInventoryCount struct {
	Result vm.VarRef `op:"result"`
	Actor  vm.Param  `op:"p8" pos:"1" fmt:"id:actor"`
}

func (inst GetInventoryCount) Acronym() string { return "INVCNT" }
