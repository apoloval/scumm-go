package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

type ActorOpSetCostume struct {
	Costume vm.Param `op:"p8" pos:"1" fmt:"id:costume"`
}

func (inst ActorOpSetCostume) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("CO=%s", inst.Costume.Display(st))
}

type ActorOpDefault struct{}

func (inst ActorOpDefault) Display(st *vm.SymbolTable) string { return "DEF" }

type ActorOpTalkColor struct {
	Color vm.Param `op:"p8" pos:"1" fmt:"hex"`
}

func (inst ActorOpTalkColor) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("TKCOL=%s", inst.Color.Display(st))
}

type ActorOpName struct {
	Name string `op:"str"`
}

func (inst ActorOpName) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("NAME=%q", inst.Name)
}

type ActorOpInitAnimation struct {
	InitFrame vm.Param `op:"p8" pos:"1" fmt:"dec"`
}

func (inst ActorOpInitAnimation) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("ANINIT[%s]", inst.InitFrame.Display(st))
}

type ActorOps struct {
	Actor            vm.Param              `op:"p8" pos:"1" fmt:"id:actor"`
	SetCostume       *ActorOpSetCostume    // 0x01
	StepDist         any                   // 0x04
	Sound            any                   // 0x05
	WalkAnimation    any                   // 0x06
	TalkAnimation    any                   // 0x07
	StandAnimation   any                   // 0x08
	Animation        any                   // 0x09
	Default          *ActorOpDefault       // 0x0A
	Elevation        any                   // 0x0B
	AnimationDefault any                   // 0x0C
	Palette          any                   // 0x0D
	TalkColor        *ActorOpTalkColor     // 0x0E
	ActorName        *ActorOpName          // 0x0F
	InitAnimation    *ActorOpInitAnimation // 0x10
	ActorWidth       any                   // 0x12
	ActorScale       any                   // 0x13
	IgnoreBoxes      any                   // 0x14
}

func (inst ActorOps) Acronym() string { return "ACTOR" }

func (inst ActorOps) DisplayOperands(st *vm.SymbolTable) []string {
	var props []string
	if inst.SetCostume != nil {
		props = append(props, inst.SetCostume.Display(st))
	}
	if inst.Default != nil {
		props = append(props, inst.Default.Display(st))
	}
	if inst.TalkColor != nil {
		props = append(props, inst.TalkColor.Display(st))
	}
	if inst.ActorName != nil {
		props = append(props, inst.ActorName.Display(st))
	}
	if inst.InitAnimation != nil {
		props = append(props, inst.InitAnimation.Display(st))
	}
	return append([]string{inst.Actor.Display(st)}, props...)
}

func (inst *ActorOps) DecodeOperands(opcode vm.OpCode, r *vm.BytecodeDecoder) error {
	inst.Actor = r.DecodeByteParam(opcode, vm.ParamPos1, vm.NumberFormatActorID)
	for {
		sub := r.DecodeOpCode()
		if sub == 0xFF {
			return nil
		}
		switch sub & 0x0F {
		case 0x01:
			inst.SetCostume = &ActorOpSetCostume{
				Costume: r.DecodeByteParam(sub, vm.ParamPos1, vm.NumberFormatCostumeID),
			}
		case 0x0A:
			inst.Default = new(ActorOpDefault)
		case 0x0E:
			inst.TalkColor = &ActorOpTalkColor{
				Color: r.DecodeByteParam(sub, vm.ParamPos1, vm.NumberFormatHex),
			}
		case 0x0F:
			inst.ActorName = &ActorOpName{
				Name: r.DecodeString(),
			}
		case 0x10:
			inst.InitAnimation = &ActorOpInitAnimation{
				InitFrame: r.DecodeByteParam(sub, vm.ParamPos1, vm.NumberFormatDecimal),
			}
		default:
			return fmt.Errorf("unknown sub-opcode %02X for actor ops instructions", sub)
		}
	}
}

type ActorFromPos struct {
	Result vm.VarRef `op:"result"`
	X      vm.Param  `op:"p16" pos:"1" fmt:"dec"`
	Y      vm.Param  `op:"p16" pos:"2" fmt:"dec"`
}

func (inst ActorFromPos) Acronym() string { return "ACTORAT" }
