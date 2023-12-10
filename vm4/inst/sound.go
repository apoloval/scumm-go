package inst

import "github.com/apoloval/scumm-go/vm"

type StartMusic struct {
	Music vm.Param `op:"p8" pos:"1" fmt:"id:music"`
}

func (inst StartMusic) Acronym() string { return "STARTMUS" }

type StopMusic struct{}

func (inst StopMusic) Acronym() string { return "STOPMUS" }

type StartSound struct {
	Sound vm.Param `op:"p8" pos:"1" fmt:"id:sound"`
}

func (inst StartSound) Acronym() string { return "STARTSND" }

type StopSound struct {
	Sound vm.Param `op:"p8" pos:"1" fmt:"id:sound"`
}

type IsSoundRunning struct {
	Result vm.VarRef `op:"result"`
	Sound  vm.Param  `op:"p8" pos:"1" fmt:"id:sound"`
}

func (inst IsSoundRunning) Acronym() string { return "SNDRUN" }
