package inst

import "github.com/apoloval/scumm-go/vm"

type StartMusic struct {
	Music vm.Param `op:"p8" pos:"1" fmt:"id:music"`
}

func (inst StartMusic) Acronym() string { return "STMU" }
