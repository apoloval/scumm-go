package inst

import "github.com/apoloval/scumm-go/vm"

type LoadString struct {
	instruction
	Dest vm.Param
	Val  []byte
}
