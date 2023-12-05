package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

type ResourceRoutine struct {
	ResourceID vm.Param `op:"p8" pos:"1"`
}

type LoadScript struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:script"`
}
type LoadSound struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:sound"`
}
type LoadCostume struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:costume"`
}
type LoadRoom struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:room"`
}
type LoadCharset struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:charset"`
}

type NukeScript struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:script"`
}
type NukeSound struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:sound"`
}
type NukeCostume struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:costume"`
}
type NukeRoom struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:room"`
}
type NukeCharset struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:charset"`
}

type LockScript struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:script"`
}
type LockSound struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:sound"`
}
type LockCostume struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:costume"`
}
type LockRoom struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:room"`
}

type UnlockScript struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:script"`
}
type UnlockSound struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:sound"`
}
type UnlockCostume struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:costume"`
}
type UnlockRoom struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:room"`
}

type ClearHeap struct{}

type LoadObject struct {
	RoomID   vm.Param `type:"byte" pos:"1" fmt:"id:room"`
	ObjectID vm.Param `type:"word" pos:"2"`
}

func decodeResourceRoutine(opcode vm.OpCode, r *vm.BytecodeReader) (inst vm.Instruction, err error) {
	sub := r.ReadOpCode()
	switch sub & 0x1F {
	case 0x01:
		inst = new(LoadScript)
	case 0x02:
		inst = new(LoadSound)
	case 0x03:
		inst = new(LoadCostume)
	case 0x04:
		inst = new(LoadRoom)
	case 0x05:
		inst = new(NukeScript)
	case 0x06:
		inst = new(NukeSound)
	case 0x07:
		inst = new(NukeCostume)
	case 0x08:
		inst = new(NukeRoom)
	case 0x09:
		inst = new(LockScript)
	case 0x0A:
		inst = new(LockSound)
	case 0x0B:
		inst = new(LockCostume)
	case 0x0C:
		inst = new(LockRoom)
	case 0x0D:
		inst = new(UnlockScript)
	case 0x0E:
		inst = new(UnlockSound)
	case 0x0F:
		inst = new(UnlockCostume)
	case 0x10:
		inst = new(UnlockRoom)
	case 0x11:
		inst = new(ClearHeap)
	case 0x12:
		inst = new(LoadCharset)
	case 0x13:
		inst = new(NukeCharset)
	case 0x14:
		inst = new(LoadObject)
	default:
		return nil, fmt.Errorf("unknown opcode %02X %02X for resource routine", opcode, sub)
	}
	err = vm.DecodeOperands(sub, r, inst)
	return
}
