package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

type ResourceRoutine struct {
	ResourceID vm.Param `op:"p8" pos:"1"`
}

type ResourceLoadScript struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:script"`
}
type ResourceLoadSound struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:sound"`
}
type ResourceLoadCostume struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:costume"`
}
type ResourceLoadRoom struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:room"`
}
type ResourceLoadCharset struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:charset"`
}

type ResourceNukeScript struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:script"`
}
type ResourceNukeSound struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:sound"`
}
type ResourceNukeCostume struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:costume"`
}
type ResourceNukeRoom struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:room"`
}
type ResourceNukeCharset struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:charset"`
}

type ResourceLockScript struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:script"`
}
type ResourceLockSound struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:sound"`
}
type ResourceLockCostume struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:costume"`
}
type ResourceLockRoom struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:room"`
}

type ResourceUnlockScript struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:script"`
}
type ResourceUnlockSound struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:sound"`
}
type ResourceUnlockCostume struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:costume"`
}
type ResourceUnlockRoom struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:room"`
}

type ResourceClearHeap struct{}

type ResourceLoadObject struct {
	RoomID   vm.Param `type:"byte" pos:"1" fmt:"id:room"`
	ObjectID vm.Param `type:"word" pos:"2"`
}

func decodeResourceRoutine(opcode vm.OpCode, r *vm.BytecodeReader) (inst vm.Instruction, err error) {
	sub := r.ReadOpCode()
	switch sub & 0x1F {
	case 0x01:
		inst = new(ResourceLoadScript)
	case 0x02:
		inst = new(ResourceLoadSound)
	case 0x03:
		inst = new(ResourceLoadCostume)
	case 0x04:
		inst = new(ResourceLoadRoom)
	case 0x05:
		inst = new(ResourceNukeScript)
	case 0x06:
		inst = new(ResourceNukeSound)
	case 0x07:
		inst = new(ResourceNukeCostume)
	case 0x08:
		inst = new(ResourceNukeRoom)
	case 0x09:
		inst = new(ResourceLockScript)
	case 0x0A:
		inst = new(ResourceLockSound)
	case 0x0B:
		inst = new(ResourceLockCostume)
	case 0x0C:
		inst = new(ResourceLockRoom)
	case 0x0D:
		inst = new(ResourceUnlockScript)
	case 0x0E:
		inst = new(ResourceUnlockSound)
	case 0x0F:
		inst = new(ResourceUnlockCostume)
	case 0x10:
		inst = new(ResourceUnlockRoom)
	case 0x11:
		inst = new(ResourceClearHeap)
	case 0x12:
		inst = new(ResourceLoadCharset)
	case 0x13:
		inst = new(ResourceNukeCharset)
	case 0x14:
		inst = new(ResourceLoadObject)
	default:
		return nil, fmt.Errorf("unknown opcode %02X %02X for resource routine", opcode, sub)
	}
	err = vm.DecodeOperands(sub, r, inst)
	return
}
