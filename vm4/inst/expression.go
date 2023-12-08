package inst

import (
	"fmt"
	"strings"

	"github.com/apoloval/scumm-go/vm"
)

type ExpressionVal struct {
	Value vm.Param
	Inst  vm.Instruction
}

func (val ExpressionVal) Display(st *vm.SymbolTable) string {
	if val.Inst != nil {
		return fmt.Sprintf("(%s)",
			strings.TrimSpace(vm.DisplayInstruction(st, val.Inst)))
	}
	return val.Value.Display(st)
}

type ExpressionOp vm.OpCode

const (
	ExpressionOpAdd ExpressionOp = 0x02
	ExpressionOpSub ExpressionOp = 0x03
	ExpressionOpMul ExpressionOp = 0x04
	ExpressionOpDiv ExpressionOp = 0x05
)

func (op ExpressionOp) String() string {
	switch op {
	case ExpressionOpAdd:
		return "+"
	case ExpressionOpSub:
		return "-"
	case ExpressionOpMul:
		return "*"
	case ExpressionOpDiv:
		return "/"
	default:
		panic("unknown expression op")
	}
}

type Expression struct {
	Result vm.VarRef `op:"result"`
	Values []ExpressionVal
	Ops    []ExpressionOp
}

func (inst Expression) Acronym() string { return "EXPR" }

func (inst Expression) DisplayOperands(st *vm.SymbolTable) []string {
	ops := []string{inst.Values[0].Display(st)}
	for i, op := range inst.Ops {
		ops = append(ops, fmt.Sprintf("%s %s", op, inst.Values[i+1].Display(st)))
	}
	return []string{
		inst.Result.Display(st),
		strings.Join(ops, " "),
	}
}

func (inst *Expression) DecodeOperands(opcode vm.OpCode, r *vm.BytecodeDecoder) error {
	inst.Result = r.DecodeVarRef()
	for {
		sub := r.DecodeOpCode()
		if sub == 0xFF {
			return nil
		}
		switch sub & 0x1F {
		case 0x01:
			val := ExpressionVal{
				Value: r.DecodeWordParam(sub, vm.ParamPos1, vm.NumberFormatDecimal),
			}
			inst.Values = append(inst.Values, val)
		case 0x02, 0x03, 0x04, 0x05:
			inst.Ops = append(inst.Ops, ExpressionOp(sub))
		case 0x06:
			nested, err := Decode(r)
			if err != nil {
				return err
			}
			val := ExpressionVal{Inst: nested}
			inst.Values = append(inst.Values, val)
		default:
			return fmt.Errorf("unknown sub-opcode %02X decoding expression", sub)
		}
	}
}

type Add struct {
	Result vm.VarRef `op:"result"`
	Value  vm.Param  `op:"p16" pos:"1" fmt:"dec"`
}

func (inst Add) Acronym() string { return "ADD" }

type Sub struct {
	Result vm.VarRef `op:"result"`
	Value  vm.Param  `op:"p16" pos:"1" fmt:"dec"`
}

func (inst Sub) Acronym() string { return "SUB" }

type Mult struct {
	Result vm.VarRef `op:"result"`
	Value  vm.Param  `op:"p16" pos:"1" fmt:"dec"`
}

func (inst Mult) Acronym() string { return "MULT" }

type Div struct {
	Result vm.VarRef `op:"result"`
	Value  vm.Param  `op:"p16" pos:"1" fmt:"dec"`
}

func (inst Div) Acronym() string { return "DIV" }

type And struct {
	Result vm.VarRef `op:"result"`
	Value  vm.Param  `op:"p16" pos:"1" fmt:"dec"`
}

func (inst And) Acronym() string { return "AND" }

type Or struct {
	Result vm.VarRef `op:"result"`
	Value  vm.Param  `op:"p16" pos:"1" fmt:"dec"`
}

func (inst Or) Acronym() string { return "OR" }
