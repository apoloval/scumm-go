package inst

import "github.com/apoloval/scumm-go/vm"

// ActorPut is a instruction to put an actor at a given position.
type ActorPut struct {
	Actor vm.Param `op:"p8" pos:"1" fmt:"dec"`
	X     vm.Param `op:"p16" pos:"2" fmt:"dec"`
	Y     vm.Param `op:"p16" pos:"3" fmt:"dec"`
}

func (inst ActorPut) Acronym() string { return "ACPUT" }

type PutActorInRoom struct {
	Actor vm.Param `op:"p8" pos:"1" fmt:"id:actor"`
	Room  vm.Param `op:"p8" pos:"2" fmt:"id:room"`
}

func (inst PutActorInRoom) Acronym() string { return "PAIR" }

type AnimateActor struct {
	Actor     vm.Param `op:"p8" pos:"1" fmt:"id:actor"`
	Animation vm.Param `op:"p8" pos:"2" fmt:"dec"`
}

func (inst AnimateActor) Acronym() string { return "ANIM" }

// PickUpObject is a instruction for the ego actor to pick up an object.
type PickUpObject struct {
	Object vm.Param `op:"p16" pos:"1" fmt:"dec"`
}

func (inst PickUpObject) Acronym() string { return "PICK" }

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

type PanCameraTo struct {
	X vm.Param `op:"p16" pos:"1" fmt:"dec"`
}

func (inst PanCameraTo) Acronym() string { return "PANC" }

type ActorFollowCamera struct {
	Actor vm.Param `op:"p8" pos:"1" fmt:"id:actor"`
}

func (inst ActorFollowCamera) Acronym() string { return "AFC" }

type DoSentence struct {
	Verb vm.Param `op:"p8" pos:"1" fmt:"dec"`
	Obj1 vm.Param `op:"p16" pos:"2" fmt:"id:object"`
	Obj2 vm.Param `op:"p16" pos:"3" fmt:"id:object"`
}

func (inst DoSentence) Acronym() string { return "DOSENT" }

type WalkActorToObject struct {
	Actor  vm.Param `op:"p8" pos:"1" fmt:"id:actor"`
	Object vm.Param `op:"p16" pos:"2" fmt:"id:object"`
}

func (inst WalkActorToObject) Acronym() string { return "WALKO" }
