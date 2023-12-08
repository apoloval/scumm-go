package inst

import "github.com/apoloval/scumm-go/vm"

type GetInventoryCount struct {
	Result vm.VarRef `op:"result"`
	Actor  vm.Param  `op:"p8" pos:"1" fmt:"id:actor"`
}

func (inst GetInventoryCount) Acronym() string { return "INVCNT" }

type FindInventory struct {
	Result vm.VarRef `op:"result"`
	Owner  vm.Param  `op:"p8" pos:"1" fmt:"id:actor"`
	Index  vm.Param  `op:"p8" pos:"2" fmt:"dec"`
}

func (inst FindInventory) Acronym() string { return "FINDINV" }
