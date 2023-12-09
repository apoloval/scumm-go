package inst

import "github.com/apoloval/scumm-go/vm"

type StartMusic struct {
	Music vm.Param `op:"p8" pos:"1" fmt:"id:music"`
}

func (inst StartMusic) Acronym() string { return "STMU" }

type StartSound struct {
	Sound vm.Param `op:"p8" pos:"1" fmt:"id:sound"`
}

func (inst StartSound) Acronym() string { return "STSN" }

type IsSoundRunning struct {
	Result vm.VarRef `op:"result"`
	Sound  vm.Param  `op:"p8" pos:"1" fmt:"id:sound"`
}

func (inst IsSoundRunning) Acronym() string { return "ISSNDRUN" }
