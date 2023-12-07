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

func (inst ResourceLoadScript) Acronym() string { return "LDSC" }

type ResourceLoadSound struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:sound"`
}

func (inst ResourceLoadSound) Acronym() string { return "LDSN" }

type ResourceLoadCostume struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:costume"`
}

func (inst ResourceLoadCostume) Acronym() string { return "LDCO" }

type ResourceLoadRoom struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:room"`
}

func (inst ResourceLoadRoom) Acronym() string { return "LDRO" }

type ResourceLoadCharset struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:charset"`
}

func (inst ResourceLoadCharset) Acronym() string { return "LDCH" }

type ResourceNukeScript struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:script"`
}

func (inst ResourceNukeScript) Acronym() string { return "NKSC" }

type ResourceNukeSound struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:sound"`
}

func (inst ResourceNukeSound) Acronym() string { return "NKSN" }

type ResourceNukeCostume struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:costume"`
}

func (inst ResourceNukeCostume) Acronym() string { return "NKCO" }

type ResourceNukeRoom struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:room"`
}

func (inst ResourceNukeRoom) Acronym() string { return "NKRO" }

type ResourceNukeCharset struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:charset"`
}

func (inst ResourceNukeCharset) Acronym() string { return "NKCH" }

type ResourceLockScript struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:script"`
}

func (inst ResourceLockScript) Acronym() string { return "LKSC" }

type ResourceLockSound struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:sound"`
}

func (inst ResourceLockSound) Acronym() string { return "LKSN" }

type ResourceLockCostume struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:costume"`
}

func (inst ResourceLockCostume) Acronym() string { return "LKCO" }

type ResourceLockRoom struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:room"`
}

func (inst ResourceLockRoom) Acronym() string { return "LKRO" }

type ResourceUnlockScript struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:script"`
}

func (inst ResourceUnlockScript) Acronym() string { return "ULSC" }

type ResourceUnlockSound struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:sound"`
}

func (inst ResourceUnlockSound) Acronym() string { return "ULSN" }

type ResourceUnlockCostume struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:costume"`
}

func (inst ResourceUnlockCostume) Acronym() string { return "ULCO" }

type ResourceUnlockRoom struct {
	ResourceID vm.Param `op:"p8" pos:"1" fmt:"id:room"`
}

func (inst ResourceUnlockRoom) Acronym() string { return "ULRO" }

type ResourceClearHeap struct{}

func (inst ResourceClearHeap) Acronym() string { return "CLRH" }

type ResourceLoadObject struct {
	RoomID   vm.Param `type:"byte" pos:"1" fmt:"id:room"`
	ObjectID vm.Param `type:"word" pos:"2"`
}

func (inst ResourceLoadObject) Acronym() string { return "LDO" }

func decodeResourceRoutine(opcode vm.OpCode, r *vm.BytecodeDecoder) (inst vm.Instruction, err error) {
	sub := r.DecodeOpCode()
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
