package inst

import "github.com/apoloval/scumm-go/vm"

type PutActor struct {
	Actor vm.Param `op:"p8" pos:"1" fmt:"dec"`
	X     vm.Param `op:"p16" pos:"2" fmt:"dec"`
	Y     vm.Param `op:"p16" pos:"3" fmt:"dec"`
}
