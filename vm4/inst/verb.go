package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

type VerbImage struct {
	Object vm.Param `op:"p16" pos:"1" fmt:"id:object"`
}

func (inst VerbImage) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("IMG=%s", inst.Object.Display(st))
}

type VerbName struct {
	Name string `op:"string"`
}

func (inst VerbName) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("NAME=%q", inst.Name)
}

type VerbColor struct {
	Color vm.Param `op:"p8" pos:"1" fmt:"dec"`
}

func (inst VerbColor) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("COL=%s", inst.Color.Display(st))
}

type VerbHiColor struct {
	Color vm.Param `op:"p8" pos:"1" fmt:"dec"`
}

func (inst VerbHiColor) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("HICOL=%s", inst.Color.Display(st))
}

type VerbAt struct {
	Left vm.Param `op:"p16" pos:"1" fmt:"dec"`
	Top  vm.Param `op:"p16" pos:"2" fmt:"dec"`
}

func (inst VerbAt) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("AT=[%s,%s]", inst.Left.Display(st), inst.Top.Display(st))
}

type VerbOn struct{}

func (inst VerbOn) Display(st *vm.SymbolTable) string { return "ON" }

type VerbOff struct{}

func (inst VerbOff) Display(st *vm.SymbolTable) string { return "OFF" }

type VerbDelete struct{}

func (inst VerbDelete) Display(st *vm.SymbolTable) string { return "DEL" }

type VerbNew struct{}

func (inst VerbNew) Display(st *vm.SymbolTable) string { return "NEW" }

type VerbDimColor struct {
	Color vm.Param `op:"p8" pos:"1" fmt:"dec"`
}

func (inst VerbDimColor) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("DIMCOL=%s", inst.Color.Display(st))
}

type VerbDim struct{}

func (inst VerbDim) Display(st *vm.SymbolTable) string { return "DIM" }

type VerbKey struct {
	Key vm.Param `op:"p8" pos:"1" fmt:"dec"`
}

func (inst VerbKey) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("KEY=%s", inst.Key.Display(st))
}

type VerbCenter struct{}

func (inst VerbCenter) Display(st *vm.SymbolTable) string { return "CENT" }

type VerbNameStr struct {
	String vm.Param `op:"168" pos:"1" fmt:"id:string"`
}

func (inst VerbNameStr) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("NAMESTR=%s", inst.String.Display(st))
}

type VerbAssignObject struct {
	Object vm.Param `op:"p16" pos:"1" fmt:"id:object"`
	Room   vm.Param `op:"p8" pos:"2" fmt:"id:room"`
}

func (inst VerbAssignObject) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("ASSIGN=[%s,%s]", inst.Object.Display(st), inst.Room.Display(st))
}

type VerbSetBackColor struct {
	Color vm.Param `op:"p8" pos:"1" fmt:"dec"`
}

func (inst VerbSetBackColor) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("BACKCOL=%s", inst.Color.Display(st))
}

type Verb struct {
	Verb         vm.Param          `op:"p8" pos:"1" fmt:"id:verb"`
	Image        *VerbImage        // 0x01
	Name         *VerbName         // 0x02
	Color        *VerbColor        // 0x03
	HiColor      *VerbHiColor      // 0x04
	At           *VerbAt           // 0x05
	On           *VerbOn           // 0x06
	Off          *VerbOff          // 0x07
	Delete       *VerbDelete       // 0x08
	New          *VerbNew          // 0x09
	DimColor     *VerbDimColor     // 0x0A
	Dim          *VerbDim          // 0x0B
	Key          *VerbKey          // 0x0C
	Center       *VerbCenter       // 0x0D
	NameStr      *VerbNameStr      // 0x0E
	AssignObject *VerbAssignObject // 0x0F
	SetBackColor *VerbSetBackColor // 0x10

}

func (inst Verb) Acronym() string { return "VERB" }

func (inst Verb) DisplayOperands(st *vm.SymbolTable) []string {
	var props []string
	if inst.Image != nil {
		props = append(props, inst.Image.Display(st))
	}
	if inst.Name != nil {
		props = append(props, inst.Name.Display(st))
	}
	if inst.Color != nil {
		props = append(props, inst.Color.Display(st))
	}
	if inst.HiColor != nil {
		props = append(props, inst.HiColor.Display(st))
	}
	if inst.At != nil {
		props = append(props, inst.At.Display(st))
	}
	if inst.On != nil {
		props = append(props, inst.On.Display(st))
	}
	if inst.Off != nil {
		props = append(props, inst.Off.Display(st))
	}
	if inst.Delete != nil {
		props = append(props, inst.Delete.Display(st))
	}
	if inst.New != nil {
		props = append(props, inst.New.Display(st))
	}
	if inst.DimColor != nil {
		props = append(props, inst.DimColor.Display(st))
	}
	if inst.Dim != nil {
		props = append(props, inst.Dim.Display(st))
	}
	if inst.Key != nil {
		props = append(props, inst.Key.Display(st))
	}
	if inst.Center != nil {
		props = append(props, inst.Center.Display(st))
	}
	if inst.NameStr != nil {
		props = append(props, inst.NameStr.Display(st))
	}
	if inst.AssignObject != nil {
		props = append(props, inst.AssignObject.Display(st))
	}
	if inst.SetBackColor != nil {
		props = append(props, inst.SetBackColor.Display(st))
	}
	return append([]string{inst.Verb.Display(st)}, props...)
}

func (inst *Verb) DecodeOperands(opcode vm.OpCode, r *vm.BytecodeDecoder) error {
	inst.Verb = r.DecodeByteParam(opcode, vm.ParamPos1, vm.NumberFormatVerbID)
	for {
		sub := r.DecodeOpCode()
		if sub == 0xFF {
			return nil
		}
		switch sub & 0x1F {
		case 0x01:
			inst.Image = &VerbImage{
				Object: r.DecodeByteParam(opcode, vm.ParamPos1, vm.NumberFormatObjectID),
			}
		case 0x02:
			inst.Name = &VerbName{Name: r.DecodeString()}
		case 0x03:
			inst.Color = &VerbColor{
				Color: r.DecodeByteParam(opcode, vm.ParamPos1, vm.NumberFormatDecimal),
			}
		case 0x04:
			inst.HiColor = &VerbHiColor{
				Color: r.DecodeByteParam(opcode, vm.ParamPos1, vm.NumberFormatDecimal),
			}
		case 0x05:
			inst.At = &VerbAt{
				Left: r.DecodeByteParam(opcode, vm.ParamPos1, vm.NumberFormatDecimal),
				Top:  r.DecodeByteParam(opcode, vm.ParamPos2, vm.NumberFormatDecimal),
			}
		case 0x06:
			inst.On = &VerbOn{}
		case 0x07:
			inst.Off = &VerbOff{}
		case 0x08:
			inst.Delete = &VerbDelete{}
		case 0x09:
			inst.New = &VerbNew{}
		case 0x0A:
			inst.DimColor = &VerbDimColor{
				Color: r.DecodeByteParam(opcode, vm.ParamPos1, vm.NumberFormatDecimal),
			}
		case 0x0B:
			inst.Dim = &VerbDim{}
		case 0x0C:
			inst.Key = &VerbKey{
				Key: r.DecodeByteParam(opcode, vm.ParamPos1, vm.NumberFormatDecimal),
			}
		case 0x0D:
			inst.Center = &VerbCenter{}
		case 0x0E:
			inst.NameStr = &VerbNameStr{
				String: r.DecodeByteParam(opcode, vm.ParamPos1, vm.NumberFormatStringID),
			}
		case 0x0F:
			inst.AssignObject = &VerbAssignObject{
				Object: r.DecodeByteParam(opcode, vm.ParamPos1, vm.NumberFormatObjectID),
				Room:   r.DecodeByteParam(opcode, vm.ParamPos2, vm.NumberFormatRoomID),
			}
		case 0x10:
			inst.SetBackColor = &VerbSetBackColor{
				Color: r.DecodeByteParam(opcode, vm.ParamPos1, vm.NumberFormatDecimal),
			}
		default:
			return fmt.Errorf("unknown opcode %02X %02X for verb operation", opcode, sub)
		}
	}
}
