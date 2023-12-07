package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

type PrintPos struct {
	XPos vm.Param
	YPos vm.Param
}

func (inst PrintPos) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("POS=[%s,%s]", inst.XPos.Display(st), inst.YPos.Display(st))
}

type PrintColor struct {
	Color vm.Param
}

func (inst PrintColor) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("COL=%s", inst.Color.Display(st))
}

type PrintClipped struct {
	Right vm.Param
}

func (inst PrintClipped) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("CLI=%s", inst.Right.Display(st))
}

type PrintErase struct {
	Width  vm.Param
	Height vm.Param
}

func (inst PrintErase) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("ERA=[%s,%s]", inst.Width.Display(st), inst.Height.Display(st))
}

type PrintCenter struct{}

func (inst PrintCenter) Display(st *vm.SymbolTable) string { return "CENT" }

type PrintLeft struct{}

func (inst PrintLeft) Display(st *vm.SymbolTable) string { return "LEFT" }

type PrintOverhead struct{}

func (inst PrintOverhead) Display(st *vm.SymbolTable) string { return "OVER" }

type PrintText struct {
	Text string
}

func (inst PrintText) Display(st *vm.SymbolTable) string {
	return fmt.Sprintf("text=%q", inst.Text)
}

type Print struct {
	Actor    vm.Param
	Pos      *PrintPos
	Color    *PrintColor
	Clipped  *PrintClipped
	Erase    *PrintErase
	Center   *PrintCenter
	Left     *PrintLeft
	Overhead *PrintOverhead
	Text     *PrintText
}

func (inst Print) Acronym() string { return "PRINT" }

func (inst Print) DisplayOperands(st *vm.SymbolTable) []string {
	var props []string
	if inst.Pos != nil {
		props = append(props, inst.Pos.Display(st))
	}
	if inst.Color != nil {
		props = append(props, inst.Color.Display(st))
	}
	if inst.Clipped != nil {
		props = append(props, inst.Clipped.Display(st))
	}
	if inst.Erase != nil {
		props = append(props, inst.Erase.Display(st))
	}
	if inst.Center != nil {
		props = append(props, inst.Center.Display(st))
	}
	if inst.Left != nil {
		props = append(props, inst.Left.Display(st))
	}
	if inst.Overhead != nil {
		props = append(props, inst.Overhead.Display(st))
	}
	if inst.Text != nil {
		props = append(props, inst.Text.Display(st))
	}

	return append([]string{inst.Actor.Display(st)}, props...)
}

func (inst *Print) DecodeOperands(opcode vm.OpCode, r *vm.BytecodeDecoder) error {
	inst.Actor = r.DecodeByteParam(opcode, vm.ParamPos1, vm.NumberFormatActorID)
	for {
		sub := r.DecodeOpCode()
		if sub == 0xFF {
			return nil
		}
		switch sub & 0x0F {
		case 0x00:
			inst.Pos = &PrintPos{
				XPos: r.DecodeWordParam(sub, vm.ParamPos1, vm.NumberFormatDecimal),
				YPos: r.DecodeWordParam(sub, vm.ParamPos2, vm.NumberFormatDecimal),
			}
		case 0x01:
			inst.Color = &PrintColor{
				Color: r.DecodeByteParam(sub, vm.ParamPos1, vm.NumberFormatDecimal),
			}
		case 0x02:
			inst.Clipped = &PrintClipped{
				Right: r.DecodeWordParam(sub, vm.ParamPos1, vm.NumberFormatDecimal),
			}
		case 0x03:
			inst.Erase = &PrintErase{
				Width:  r.DecodeWordParam(sub, vm.ParamPos1, vm.NumberFormatDecimal),
				Height: r.DecodeWordParam(sub, vm.ParamPos2, vm.NumberFormatDecimal),
			}
		case 0x04:
			inst.Center = &PrintCenter{}
		case 0x06:
			inst.Left = &PrintLeft{}
		case 0x07:
			inst.Overhead = &PrintOverhead{}
		case 0x0F:
			// TODO: check if this is null terminated of 0xFF terminated
			inst.Text = &PrintText{Text: r.DecodeString()}
		default:
			return fmt.Errorf("unknown sub-opcode %02X for print op operation", sub)
		}
	}
}
