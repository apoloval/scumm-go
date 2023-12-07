package inst

import "github.com/apoloval/scumm-go/vm"

// ActorPut is a instruction to put an actor at a given position.
type ActorPut struct {
	Actor vm.Param `op:"p8" pos:"1" fmt:"dec"`
	X     vm.Param `op:"p16" pos:"2" fmt:"dec"`
	Y     vm.Param `op:"p16" pos:"3" fmt:"dec"`
}

func (inst ActorPut) Acronym() string { return "ACPUT" }

// PickUpObject is a instruction for the ego actor to pick up an object.
type PickUpObject struct {
	Object vm.Param `op:"p16" pos:"1" fmt:"dec"`
}

func (inst PickUpObject) Acronym() string { return "PICK" }

type SetClass struct {
	Object  vm.Param  `op:"p16" pos:"1" fmt:"id:object"`
	Classes vm.Params `op:"v16"`
}

func (inst SetClass) Acronym() string { return "CLASS" }

type SetObjectName struct {
	Object vm.Param `op:"p16" pos:"1" fmt:"id:object"`
	Name   string   `op:"string"`
}

func (inst SetObjectName) Acronym() string { return "OBJN" }

type GetObjectOwner struct {
	Result vm.VarRef `op:"result"`
	Object vm.Param  `op:"p16" pos:"1" fmt:"id:object"`
}

func (inst GetObjectOwner) Acronym() string { return "OBJO" }

type SetObjectOwner struct {
	Object vm.Param `op:"p16" pos:"1" fmt:"id:object"`
	Owner  vm.Param `op:"p8" pos:"2" fmt:"id:actor"`
}

func (inst SetObjectOwner) Acronym() string { return "OWN" }
