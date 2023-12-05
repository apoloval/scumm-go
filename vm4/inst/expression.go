package inst

import (
	"fmt"
	"strings"

	"github.com/apoloval/scumm-go/vm"
)

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
	base
	Dest   vm.Pointer
	Values []vm.Param
	Ops    []ExpressionOp
}

func (inst Expression) Mnemonic(st *vm.SymbolTable) string {
	var str strings.Builder

	fmt.Fprintf(&str, "%s = %s", inst.Dest.Display(st), inst.Values[0].Display(st))
	for i, op := range inst.Ops {
		fmt.Fprintf(&str, " %s %s", op, inst.Values[i+1].Display(st))
	}

	return str.String()
}

func (inst *Expression) Decode(opcode vm.OpCode, r *vm.BytecodeReader) error {
	inst.Dest = r.ReadPointer()
	for {
		sub := r.ReadOpCode()
		if sub == 0xFF {
			return inst.base.Decode(opcode, r)
		}
		switch sub & 0x1F {
		case 0x01:
			inst.Values = append(inst.Values, r.ReadWordParam(sub, vm.ParamPos1, vm.ParamFormatNumber))
		case 0x02, 0x03, 0x04, 0x05:
			inst.Ops = append(inst.Ops, ExpressionOp(sub))
		default:
			return fmt.Errorf("unknown sub-opcode %02X decoding expression", sub)
		}
	}
}
