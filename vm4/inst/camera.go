package inst

import "github.com/apoloval/scumm-go/vm"

type SetCameraAt struct {
	X vm.Param `op:"p16" pos:"1" fmt:"dec"`
}

func (inst SetCameraAt) MNemonic() string { return "SETCAMAT" }
