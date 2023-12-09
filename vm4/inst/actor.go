package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

type ActorDummy struct {
	Arg1 vm.Param `op:"p8" pos:"1" fmt:"dec"`
}

func (inst ActorDummy) Display(st *vm.SymbolTable) string { return "DUMMY" }

type ActorCostume struct {
	Costume vm.Param `op:"p8" pos:"1" fmt:"id:costume"`
}

func (inst ActorCostume) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("CO=%s", inst.Costume.Display(st))
}

type ActorStepDist struct {
	XSpeed vm.Param `op:"p8" pos:"1" fmt:"dec"`
	YSpeed vm.Param `op:"p8" pos:"2" fmt:"dec"`
}

func (inst ActorStepDist) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("STEPDIST=%s,%s", inst.XSpeed.Display(st), inst.YSpeed.Display(st))
}

type ActorSound struct {
	Sound vm.Param `op:"p8" pos:"1" fmt:"id:sound"`
}

func (inst ActorSound) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("SOUND=%s", inst.Sound.Display(st))
}

type ActorWalkAnimation struct {
	WalkFrame vm.Param `op:"p8" pos:"1" fmt:"dec"`
}

func (inst ActorWalkAnimation) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("WALK=%s", inst.WalkFrame.Display(st))
}

type ActorTalkAnimation struct {
	StartTalk vm.Param `op:"p8" pos:"1" fmt:"dec"`
	EndTalk   vm.Param `op:"p8" pos:"2" fmt:"dec"`
}

func (inst ActorTalkAnimation) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("TALK=[%s-%s]", inst.StartTalk.Display(st), inst.EndTalk.Display(st))
}

type ActorStandAnimation struct {
	StandFrame vm.Param `op:"p8" pos:"1" fmt:"dec"`
}

func (inst ActorStandAnimation) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("STAND=%s", inst.StandFrame.Display(st))
}

type ActorAnimation struct {
	Arg1 vm.Param `op:"p8" pos:"1" fmt:"dec"`
	Arg2 vm.Param `op:"p8" pos:"2" fmt:"dec"`
	Arg3 vm.Param `op:"p8" pos:"3" fmt:"dec"`
}

func (inst ActorAnimation) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("ANIM=[%s,%s,%s]",
		inst.Arg1.Display(st), inst.Arg2.Display(st), inst.Arg3.Display(st))
}

type ActorDefault struct{}

func (inst ActorDefault) Display(st *vm.SymbolTable) string { return "DEF" }

type ActorElevation struct {
	Elevation vm.Param `op:"p16" pos:"1" fmt:"dec"`
}

func (inst ActorElevation) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("ELEV=%s", inst.Elevation.Display(st))
}

type ActorAnimationDefault struct{}

func (inst ActorAnimationDefault) Display(st *vm.SymbolTable) string { return "ANIMDEF" }

type ActorPalette struct {
	Index vm.Param `op:"p8" pos:"1" fmt:"dec"`
	Value vm.Param `op:"p8" pos:"2" fmt:"hex"`
}

func (inst ActorPalette) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("PAL[%s]=%s", inst.Index.Display(st), inst.Value.Display(st))
}

type ActorTalkColor struct {
	Color vm.Param `op:"p8" pos:"1" fmt:"hex"`
}

func (inst ActorTalkColor) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("TKCOL=%s", inst.Color.Display(st))
}

type ActorName struct {
	Name string `op:"str"`
}

func (inst ActorName) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("NAME=%q", inst.Name)
}

type ActorInitAnimation struct {
	InitFrame vm.Param `op:"p8" pos:"1" fmt:"dec"`
}

func (inst ActorInitAnimation) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("ANINIT[%s]", inst.InitFrame.Display(st))
}

type ActorWidth struct {
	Width vm.Param `op:"p8" pos:"1" fmt:"dec"`
}

func (inst ActorWidth) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("WIDTH=%s", inst.Width.Display(st))
}

type ActorScale struct {
	Scale vm.Param `op:"p8" pos:"1" fmt:"dec"`
}

func (inst ActorScale) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("SCALE=%s", inst.Scale.Display(st))
}

type ActorIgnoreBoxes struct{}

func (inst ActorIgnoreBoxes) Display(st *vm.SymbolTable) string { return "IGNBOX" }

type Actor struct {
	Actor            vm.Param               `op:"p8" pos:"1" fmt:"id:actor"`
	Dummy            *ActorDummy            // 0x00
	SetCostume       *ActorCostume          // 0x01
	StepDist         *ActorStepDist         // 0x04
	Sound            *ActorSound            // 0x05
	WalkAnimation    *ActorWalkAnimation    // 0x06
	TalkAnimation    *ActorTalkAnimation    // 0x07
	StandAnimation   *ActorStandAnimation   // 0x08
	Animation        *ActorAnimation        // 0x09
	Default          *ActorDefault          // 0x0A
	Elevation        *ActorElevation        // 0x0B
	AnimationDefault *ActorAnimationDefault // 0x0C
	Palette          *ActorPalette          // 0x0D
	TalkColor        *ActorTalkColor        // 0x0E
	ActorName        *ActorName             // 0x0F
	InitAnimation    *ActorInitAnimation    // 0x10
	ActorWidth       *ActorWidth            // 0x12
	ActorScale       *ActorScale            // 0x13
	IgnoreBoxes      *ActorIgnoreBoxes      // 0x14
}

func (inst Actor) Acronym() string { return "ACTOR" }

func (inst Actor) DisplayOperands(st *vm.SymbolTable) []string {
	var props []string
	if inst.Dummy != nil {
		props = append(props, inst.Dummy.Display(st))
	}
	if inst.SetCostume != nil {
		props = append(props, inst.SetCostume.Display(st))
	}
	if inst.StepDist != nil {
		props = append(props, inst.StepDist.Display(st))
	}
	if inst.Sound != nil {
		props = append(props, inst.Sound.Display(st))
	}
	if inst.WalkAnimation != nil {
		props = append(props, inst.WalkAnimation.Display(st))
	}
	if inst.TalkAnimation != nil {
		props = append(props, inst.TalkAnimation.Display(st))
	}
	if inst.StandAnimation != nil {
		props = append(props, inst.StandAnimation.Display(st))
	}
	if inst.Animation != nil {
		props = append(props, inst.Animation.Display(st))
	}
	if inst.Default != nil {
		props = append(props, inst.Default.Display(st))
	}
	if inst.Elevation != nil {
		props = append(props, inst.Elevation.Display(st))
	}
	if inst.AnimationDefault != nil {
		props = append(props, inst.AnimationDefault.Display(st))
	}
	if inst.Palette != nil {
		props = append(props, inst.Palette.Display(st))
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
	if inst.ActorWidth != nil {
		props = append(props, inst.ActorWidth.Display(st))
	}
	if inst.ActorScale != nil {
		props = append(props, inst.ActorScale.Display(st))
	}
	if inst.IgnoreBoxes != nil {
		props = append(props, inst.IgnoreBoxes.Display(st))
	}
	return append([]string{inst.Actor.Display(st)}, props...)
}

func (inst *Actor) DecodeOperands(opcode vm.OpCode, r *vm.BytecodeDecoder) error {
	inst.Actor = r.DecodeByteParam(opcode, vm.ParamPos1, vm.NumberFormatActorID)
	for {
		sub := r.DecodeOpCode()
		if sub == 0xFF {
			return nil
		}
		switch sub & 0x1F {
		case 0x00:
			inst.Dummy = &ActorDummy{
				Arg1: r.DecodeByteParam(sub, vm.ParamPos1, vm.NumberFormatDecimal),
			}
		case 0x01:
			inst.SetCostume = &ActorCostume{
				Costume: r.DecodeByteParam(sub, vm.ParamPos1, vm.NumberFormatCostumeID),
			}
		case 0x04:
			inst.StepDist = &ActorStepDist{
				XSpeed: r.DecodeByteParam(sub, vm.ParamPos1, vm.NumberFormatDecimal),
				YSpeed: r.DecodeByteParam(sub, vm.ParamPos2, vm.NumberFormatDecimal),
			}
		case 0x05:
			inst.Sound = &ActorSound{
				Sound: r.DecodeByteParam(sub, vm.ParamPos1, vm.NumberFormatSoundID),
			}
		case 0x06:
			inst.WalkAnimation = &ActorWalkAnimation{
				WalkFrame: r.DecodeByteParam(sub, vm.ParamPos1, vm.NumberFormatDecimal),
			}
		case 0x07:
			inst.TalkAnimation = &ActorTalkAnimation{
				StartTalk: r.DecodeByteParam(sub, vm.ParamPos1, vm.NumberFormatDecimal),
				EndTalk:   r.DecodeByteParam(sub, vm.ParamPos2, vm.NumberFormatDecimal),
			}
		case 0x08:
			inst.StandAnimation = &ActorStandAnimation{
				StandFrame: r.DecodeByteParam(sub, vm.ParamPos1, vm.NumberFormatDecimal),
			}
		case 0x09:
			inst.Animation = &ActorAnimation{
				Arg1: r.DecodeByteParam(sub, vm.ParamPos1, vm.NumberFormatDecimal),
				Arg2: r.DecodeByteParam(sub, vm.ParamPos2, vm.NumberFormatDecimal),
				Arg3: r.DecodeByteParam(sub, vm.ParamPos3, vm.NumberFormatDecimal),
			}
		case 0x0A:
			inst.Default = new(ActorDefault)
		case 0x0B:
			inst.Elevation = &ActorElevation{
				Elevation: r.DecodeWordParam(sub, vm.ParamPos1, vm.NumberFormatDecimal),
			}
		case 0x0C:
			inst.AnimationDefault = new(ActorAnimationDefault)
		case 0x0D:
			inst.Palette = &ActorPalette{
				Index: r.DecodeByteParam(sub, vm.ParamPos1, vm.NumberFormatDecimal),
				Value: r.DecodeByteParam(sub, vm.ParamPos2, vm.NumberFormatHex),
			}
		case 0x0E:
			inst.TalkColor = &ActorTalkColor{
				Color: r.DecodeByteParam(sub, vm.ParamPos1, vm.NumberFormatHex),
			}
		case 0x0F:
			inst.ActorName = &ActorName{
				Name: r.DecodeString(),
			}
		case 0x10:
			inst.InitAnimation = &ActorInitAnimation{
				InitFrame: r.DecodeByteParam(sub, vm.ParamPos1, vm.NumberFormatDecimal),
			}
		case 0x12:
			inst.ActorWidth = &ActorWidth{
				Width: r.DecodeByteParam(sub, vm.ParamPos1, vm.NumberFormatDecimal),
			}
		case 0x13:
			inst.ActorScale = &ActorScale{
				Scale: r.DecodeByteParam(sub, vm.ParamPos1, vm.NumberFormatDecimal),
			}
		case 0x14:
			inst.IgnoreBoxes = new(ActorIgnoreBoxes)
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

type GetActorX struct {
	Result vm.VarRef `op:"result"`
	Actor  vm.Param  `op:"p16" pos:"1" fmt:"id:actor"`
}

func (inst GetActorX) Acronym() string { return "ACTORX" }

type GetActorY struct {
	Result vm.VarRef `op:"result"`
	Actor  vm.Param  `op:"p16" pos:"1" fmt:"id:actor"`
}

func (inst GetActorY) Acronym() string { return "ACTORY" }

type GetActorWidth struct {
	Result vm.VarRef `op:"result"`
	Actor  vm.Param  `op:"p8" pos:"1" fmt:"id:actor"`
}

func (inst GetActorWidth) Acronym() string { return "ACTORW" }

type GetActorScale struct {
	Result vm.VarRef `op:"result"`
	Actor  vm.Param  `op:"p8" pos:"1" fmt:"id:actor"`
}

func (inst GetActorScale) Acronym() string { return "ACTORSC" }

type GetActorWalkBox struct {
	Result vm.VarRef `op:"result"`
	Actor  vm.Param  `op:"p8" pos:"1" fmt:"id:actor"`
}

func (inst GetActorWalkBox) Acronym() string { return "ACTORWB" }

type GetActorFacing struct {
	Result vm.VarRef `op:"result"`
	Actor  vm.Param  `op:"p8" pos:"1" fmt:"id:actor"`
}

func (inst GetActorFacing) Acronym() string { return "ACTORFA" }

type GetActorElevation struct {
	Result vm.VarRef `op:"result"`
	Actor  vm.Param  `op:"p8" pos:"1" fmt:"id:actor"`
}

func (inst GetActorElevation) Acronym() string { return "ACTOREL" }

type GetActorMoving struct {
	Result vm.VarRef `op:"result"`
	Actor  vm.Param  `op:"p8" pos:"1" fmt:"id:actor"`
}

func (inst GetActorMoving) Acronym() string { return "ACTORMOV" }

type GetActorRoom struct {
	Result vm.VarRef `op:"result"`
	Actor  vm.Param  `op:"p8" pos:"1" fmt:"id:actor"`
}

func (inst GetActorRoom) Acronym() string { return "ACTRO" }

type GetActorCostume struct {
	Result vm.VarRef `op:"result"`
	Actor  vm.Param  `op:"p8" pos:"1" fmt:"id:actor"`
}

func (inst GetActorCostume) Acronym() string { return "ACTORCO" }

type GetActorAnimCounter struct {
	Result vm.VarRef `op:"result"`
	Actor  vm.Param  `op:"p8" pos:"1" fmt:"id:actor"`
}

func (inst GetActorAnimCounter) Acronym() string { return "ACTORAC" }

type GetActorClosestObject struct {
	Result vm.VarRef `op:"result"`
	Actor  vm.Param  `op:"p16" pos:"1" fmt:"id:actor"`
}

func (inst GetActorClosestObject) Acronym() string { return "ACTORCLOBJ" }

type FaceActor struct {
	Actor  vm.Param `op:"p8" pos:"1" fmt:"id:actor"`
	Object vm.Param `op:"p16" pos:"2" fmt:"id:object"`
}

func (inst FaceActor) Acronym() string { return "FACEA" }

type WalkActorTo struct {
	Actor vm.Param `op:"p8" pos:"1" fmt:"id:actor"`
	X     vm.Param `op:"p16" pos:"2" fmt:"dec"`
	Y     vm.Param `op:"p16" pos:"3" fmt:"dec"`
}

func (inst WalkActorTo) Acronym() string { return "WALKT" }

type WalkActorToObject struct {
	Actor  vm.Param `op:"p8" pos:"1" fmt:"id:actor"`
	Object vm.Param `op:"p16" pos:"2" fmt:"id:object"`
}

func (inst WalkActorToObject) Acronym() string { return "WALKO" }

type WalkActorToActor struct {
	Walker   vm.Param    `op:"p8" pos:"1" fmt:"id:actor"`
	Walkee   vm.Param    `op:"p8" pos:"2" fmt:"id:actor"`
	Distance vm.Constant `op:"8" fmt:"dec"`
}

func (inst WalkActorToActor) Acronym() string { return "WALKA" }

type ActorFollowCamera struct {
	Actor vm.Param `op:"p8" pos:"1" fmt:"id:actor"`
}

func (inst ActorFollowCamera) Acronym() string { return "AFC" }

// ActorPut is a instruction to put an actor at a given position.
type ActorPut struct {
	Actor vm.Param `op:"p8" pos:"1" fmt:"dec"`
	X     vm.Param `op:"p16" pos:"2" fmt:"dec"`
	Y     vm.Param `op:"p16" pos:"3" fmt:"dec"`
}

func (inst ActorPut) Acronym() string { return "ACPUT" }

type PutActorInRoom struct {
	Actor vm.Param `op:"p8" pos:"1" fmt:"id:actor"`
	Room  vm.Param `op:"p8" pos:"2" fmt:"id:room"`
}

func (inst PutActorInRoom) Acronym() string { return "PAIR" }

type AnimateActor struct {
	Actor     vm.Param `op:"p8" pos:"1" fmt:"id:actor"`
	Animation vm.Param `op:"p8" pos:"2" fmt:"dec"`
}

func (inst AnimateActor) Acronym() string { return "ANIM" }
