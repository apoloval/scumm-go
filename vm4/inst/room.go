package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

type RoomFade struct {
	Effect vm.Param `op:"p16" pos:"1"`
}

func (inst RoomFade) Acronym() string { return "ROFA" }

func (inst *RoomFade) DecodeOperands(opcode vm.OpCode, r *vm.BytecodeDecoder) error {
	sub := r.DecodeOpCode()
	if sub&0x1F == 3 {
		inst.Effect = r.DecodeWordParam(sub, vm.ParamPos1, vm.NumberFormatHex)
	}
	return nil
}

type PseudoRoom struct {
	Value       vm.Constant   `op:"8" fmt:"dec"`
	ResourceIDs []vm.Constant `op:"8" fmt:"dec"`
}

func (inst PseudoRoom) Acronym() string { return "PSRO" }

func (inst PseudoRoom) DisplayOperands(st *vm.SymbolTable) (ops []string) {
	ops = []string{inst.Value.Display(st)}
	for _, id := range inst.ResourceIDs {
		ops = append(ops, id.Display(st))
	}
	return ops
}

func (inst *PseudoRoom) DecodeOperands(opcode vm.OpCode, r *vm.BytecodeDecoder) error {
	inst.Value = r.DecodeByteConstant(vm.NumberFormatRoomID)
	inst.ResourceIDs = r.DecodeNullTerminatedBytes(vm.NumberFormatRoomID)
	params := []vm.Param{inst.Value}
	for _, id := range inst.ResourceIDs {
		params = append(params, id)
	}
	return nil
}

type RoomSetScrollLimits struct {
	MinX vm.Param `op:"p16" pos:"1"`
	MaxX vm.Param `op:"p16" pos:"2"`
}

func (inst RoomSetScrollLimits) Acronym() string { return "ROSL" }

type RoomColor struct {
	Color vm.Param `op:"p16" pos:"1" fmt:"hex"`
	Index vm.Param `op:"p16" pos:"2" fmt:"dec"`
}

func (inst RoomColor) Acronym() string { return "ROCO" }

type RoomInitScreen struct {
	B vm.Param `op:"p16" pos:"1" fmt:"dec"`
	H vm.Param `op:"p16" pos:"2" fmt:"dec"`
}

func (inst RoomInitScreen) Acronym() string { return "ROIS" }

type RoomShadowColor struct {
	Color vm.Param `op:"p16" pos:"1" fmt:"hex"`
	Index vm.Param `op:"p16" pos:"2" fmt:"dec"`
}

func (inst RoomShadowColor) Acronym() string { return "ROPAL" }

type RoomShakeOn struct{}

func (inst RoomShakeOn) Acronym() string { return "ROSHON" }

type RoomShakeOff struct{}

func (inst RoomShakeOff) Acronym() string { return "ROSHOFF" }

type RoomScale struct {
	Scale1 vm.Param `op:"p8" pos:"1" fmt:"dec"`
	Y1     vm.Param `op:"p8" pos:"2" fmt:"dec"`
	Scale2 vm.Param `op:"p8" pos:"1" fmt:"dec"` // pos relative to aux1
	Y2     vm.Param `op:"p8" pos:"2" fmt:"dec"` // pos relative to aux1
	Slot   vm.Param `op:"p8" pos:"1" fmt:"dec"` // pos relative to aux2
}

func (inst RoomScale) Acronym() string { return "ROSC" }

type RoomIntensity struct {
	Scale      vm.Param `op:"p8" pos:"1" fmt:"dec"`
	StartColor vm.Param `op:"p8" pos:"2" fmt:"hex"`
	EndColor   vm.Param `op:"p8" pos:"3" fmt:"hex"`
}

func (inst RoomIntensity) Acronym() string { return "ROINT" }

type RoomSaveGame struct {
	LoadFlag vm.Param `op:"p8" pos:"1" fmt:"dec"`
	LoadSlot vm.Param `op:"p8" pos:"2" fmt:"dec"`
}

func (inst RoomSaveGame) Acronym() string { return "ROSAVE" }

type RoomIntensityRGB struct {
	RedScale   vm.Param `op:"p16" pos:"1" fmt:"dec"`
	GreenScale vm.Param `op:"p16" pos:"2" fmt:"dec"`
	BlueScale  vm.Param `op:"p16" pos:"3" fmt:"dec"`

	StartColor vm.Param `op:"p8" pos:"1" fmt:"hex"` // pos relative to aux
	EndColor   vm.Param `op:"p8" pos:"2" fmt:"hex"` // pos relative to aux
}

type RoomShadow struct {
	RedScale   vm.Param `op:"p16" pos:"1" fmt:"dec"`
	GreenScale vm.Param `op:"p16" pos:"2" fmt:"dec"`
	BlueScale  vm.Param `op:"p16" pos:"3" fmt:"dec"`

	StartColor vm.Param `op:"p8" pos:"1" fmt:"hex"` // pos relative to aux
	EndColor   vm.Param `op:"p8" pos:"2" fmt:"hex"` // pos relative to aux
}

func decodeRoomOp(opcode vm.OpCode, r *vm.BytecodeDecoder) (inst vm.Instruction, err error) {
	sub := r.DecodeOpCode()
	switch sub & 0x1F {
	case 0x01:
		inst = new(RoomSetScrollLimits)
	case 0x02:
		inst = new(RoomColor)
	case 0x03:
		inst = new(RoomInitScreen)
	case 0x04:
		inst = new(RoomShadowColor)
	case 0x05:
		inst = new(RoomShakeOn)
	case 0x06:
		inst = new(RoomShakeOff)
	case 0x07:
		scale1 := r.DecodeByteParam(sub, vm.ParamPos1, vm.NumberFormatDecimal)
		y1 := r.DecodeByteParam(sub, vm.ParamPos2, vm.NumberFormatDecimal)
		aux1 := r.DecodeOpCode()
		scale2 := r.DecodeByteParam(aux1, vm.ParamPos1, vm.NumberFormatDecimal)
		y2 := r.DecodeByteParam(aux1, vm.ParamPos2, vm.NumberFormatDecimal)
		aux2 := r.DecodeOpCode()
		slot := r.DecodeByteParam(aux2, vm.ParamPos1, vm.NumberFormatDecimal)
		inst = &RoomScale{
			Scale1: scale1,
			Y1:     y1,
			Scale2: scale2,
			Y2:     y2,
			Slot:   slot,
		}
		return
	case 0x08:
		inst = new(RoomIntensity)
	case 0x09:
		inst = new(RoomSaveGame)
	case 0x0A:
		inst = new(RoomFade)
	case 0x0B:
		red := r.DecodeWordParam(sub, vm.ParamPos1, vm.NumberFormatHex)
		green := r.DecodeWordParam(sub, vm.ParamPos2, vm.NumberFormatHex)
		blue := r.DecodeWordParam(sub, vm.ParamPos3, vm.NumberFormatHex)
		aux := r.DecodeOpCode()
		startColor := r.DecodeByteParam(aux, vm.ParamPos1, vm.NumberFormatHex)
		endColor := r.DecodeByteParam(aux, vm.ParamPos2, vm.NumberFormatHex)
		inst = &RoomIntensityRGB{
			RedScale:   red,
			GreenScale: green,
			BlueScale:  blue,
			StartColor: startColor,
			EndColor:   endColor,
		}
		return
	case 0x0C:
		red := r.DecodeWordParam(sub, vm.ParamPos1, vm.NumberFormatHex)
		green := r.DecodeWordParam(sub, vm.ParamPos2, vm.NumberFormatHex)
		blue := r.DecodeWordParam(sub, vm.ParamPos3, vm.NumberFormatHex)
		aux := r.DecodeOpCode()
		startColor := r.DecodeByteParam(aux, vm.ParamPos1, vm.NumberFormatHex)
		endColor := r.DecodeByteParam(aux, vm.ParamPos2, vm.NumberFormatHex)
		inst = &RoomShadow{
			RedScale:   red,
			GreenScale: green,
			BlueScale:  blue,
			StartColor: startColor,
			EndColor:   endColor,
		}
		return
	default:
		// I don't know exactly what other opcodes are supported in SCUMM v4. From ScummVM source
		// code, I can infer the 0x0D and 0x0E ops for string save and load, respectively, are not
		// used in any v4 game. Considering the sequence of op codes, I don't think the 0x0F and
		// 0x10 ops are used either. So they will not be implemented for now.
		return nil, fmt.Errorf("unknown opcode %02X %02X for room operation", opcode, sub)
	}
	err = vm.DecodeOperands(sub, r, inst)
	return inst, err
}
