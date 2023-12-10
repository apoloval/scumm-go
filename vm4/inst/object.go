package inst

import "github.com/apoloval/scumm-go/vm"

// PickUpObject is a instruction for the ego actor to pick up an object.
type PickUpObject struct {
	Object vm.Param `op:"p16" pos:"1" fmt:"dec"`
}

func (inst PickUpObject) Acronym() string { return "PICK" }

type FindObject struct {
	Result vm.VarRef `op:"result"`
	X      vm.Param  `op:"p8" pos:"1" fmt:"dec"`
	Y      vm.Param  `op:"p8" pos:"2" fmt:"dec"`
}

func (inst FindObject) Acronym() string { return "FINDOBJ" }

type SetClass struct {
	Object  vm.Param  `op:"p16" pos:"1" fmt:"id:object"`
	Classes vm.Params `op:"v16"`
}

func (inst SetClass) Acronym() string { return "SOCL" }

type SetObjectName struct {
	Object vm.Param `op:"p16" pos:"1" fmt:"id:object"`
	Name   string   `op:"string"`
}

func (inst SetObjectName) Acronym() string { return "SONM" }

type GetObjectOwner struct {
	Result vm.VarRef `op:"result"`
	Object vm.Param  `op:"p16" pos:"1" fmt:"id:object"`
}

func (inst GetObjectOwner) Acronym() string { return "GOOW" }

type SetObjectOwner struct {
	Object vm.Param `op:"p16" pos:"1" fmt:"id:object"`
	Owner  vm.Param `op:"p8" pos:"2" fmt:"id:actor"`
}

func (inst SetObjectOwner) Acronym() string { return "SOOW" }

type SetObjectState struct {
	Object vm.Param `op:"p16" pos:"1" fmt:"id:object"`
	State  vm.Param `op:"p8" pos:"2" fmt:"hex"`
}

func (inst SetObjectState) Acronym() string { return "SOST" }

type GetDistance struct {
	Result vm.VarRef `op:"result"`
	Obj1   vm.Param  `op:"p16" pos:"1" fmt:"id:object"`
	Obj2   vm.Param  `op:"p16" pos:"2" fmt:"id:object"`
}

func (inst GetDistance) Acronym() string { return "DIST" }
