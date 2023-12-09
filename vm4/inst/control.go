package inst

import (
	"github.com/apoloval/scumm-go/vm"
)

// StopObjectCode is a stop instruction that stops the execution of the current script after reaching the
// end of the code. This is also known as StopObjectCode in ScummVM.
type StopObjectCode struct{}

func (inst StopObjectCode) Acronym() string { return "SOC" }

// Jump is a instruction that jumps to the given address. This is also known as JumpRelative in
// ScummVM.
type Jump struct {
	Target vm.Constant `op:"reljmp" fmt:"addr"`
}

func (inst Jump) Acronym() string { return "JMP" }

type UnaryBranch struct {
	Var    vm.VarRef   `op:"var"`
	Target vm.Constant `op:"reljmp" fmt:"addr"`
}

type BinaryBranch struct {
	Var    vm.VarRef   `op:"var"`
	Value  vm.Param    `op:"p16" pos:"1" fmt:"dec"`
	Target vm.Constant `op:"reljmp" fmt:"addr"`
}

type BranchUnlessEqual BinaryBranch

func (inst BranchUnlessEqual) Acronym() string { return "BREQ" }

type BranchUnlessNotEqual BinaryBranch

func (inst BranchUnlessNotEqual) Acronym() string { return "BRNE" }

type BranchUnlessLess BinaryBranch

func (inst BranchUnlessLess) Acronym() string { return "BRLT" }

type BranchUnlessLessEqual BinaryBranch

func (inst BranchUnlessLessEqual) Acronym() string { return "BRLE" }

type BranchUnlessGreater BinaryBranch

func (inst BranchUnlessGreater) Acronym() string { return "BRGT" }

type BranchUnlessGreaterEqual BinaryBranch

func (inst BranchUnlessGreaterEqual) Acronym() string { return "BRGE" }

type BranchUnlessZero UnaryBranch

func (inst BranchUnlessZero) Acronym() string { return "BRZE" }

type BranchUnlessNotZero UnaryBranch

func (inst BranchUnlessNotZero) Acronym() string { return "BRNZ" }

// StartObject is a instruction that starts a object script.
type StartObject struct {
	Object vm.Param  `op:"p16" pos:"1" fmt:"id:object"`
	Script vm.Param  `op:"p8" pos:"2" fmt:"id:script"`
	Args   vm.Params `op:"v16"`
}

func (inst StartObject) Acronym() string { return "STOB" }

type BreakHere struct{}

func (inst BreakHere) Acronym() string { return "BREAK" }

// LoadRoom is a instruction that loads a new room.
type LoadRoom struct {
	RoomID vm.Param `op:"p8" pos:"1" fmt:"id:room"`
}

func (inst LoadRoom) Acronym() string { return "LDRO" }

type BranchUnlessState struct {
	Object vm.Param    `op:"p16" pos:"1" fmt:"id:object"`
	State  vm.Param    `op:"p8" pos:"2" fmt:"dec"`
	Target vm.Constant `op:"reljmp" fmt:"addr"`
}

func (inst BranchUnlessState) Acronym() string { return "BRST" }

type BranchUnlessNotState struct {
	Object vm.Param    `op:"p16" pos:"1" fmt:"id:object"`
	State  vm.Param    `op:"p8" pos:"2" fmt:"dec"`
	Target vm.Constant `op:"reljmp" fmt:"addr"`
}

func (inst BranchUnlessNotState) Acronym() string { return "BRNST" }

type BranchUnlessActorInBox struct {
	Actor  vm.Param    `op:"p8" pos:"1" fmt:"dec"`
	Box    vm.Param    `op:"p8" pos:"2" fmt:"dec"`
	Target vm.Constant `op:"reljmp" fmt:"addr"`
}

func (inst BranchUnlessActorInBox) Acronym() string { return "BRAB" }

type BranchUnlessClass struct {
	Object  vm.Param    `op:"p16" pos:"1" fmt:"id:object"`
	Classes vm.Params   `op:"v16"`
	Target  vm.Constant `op:"reljmp" fmt:"addr"`
}

func (inst BranchUnlessClass) Acronym() string { return "BRCL" }

type Delay struct {
	Param vm.Param `op:"24" fmt:"dec"`
}

func (inst Delay) Acronym() string { return "DELAY" }

type DelayVar struct {
	Var vm.VarRef `op:"var"`
}

func (inst DelayVar) Acronym() string { return "DELAYVAR" }

type Debug struct {
	Param vm.Param `op:"p16" pos:"1" fmt:"dec"`
}

func (inst Debug) Acronym() string { return "DEBUG" }
