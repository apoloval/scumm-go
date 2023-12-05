package vm

import (
	"strings"
)

// OpCode is an opcode of the bytecode scripting language.
type OpCode byte

// IsPointer returns true if the opcode expects a pointer parameter by the bit 7.
func (op OpCode) IsPointer(pos ParamPos) bool {
	return byte(op)&byte(pos) > 0
}

// ParamPos is the position of a instruction parameter. It is also a bitmask to figure out if the
// upcoming parameter is a pointer or a literal value.
type ParamPos byte

const (
	ParamPos1 ParamPos = 0x80
	ParamPos2 ParamPos = 0x40
	ParamPos3 ParamPos = 0x20
)

func ParseParamPos(str string) (ParamPos, bool) {
	switch str {
	case "1":
		return ParamPos1, true
	case "2":
		return ParamPos2, true
	case "3":
		return ParamPos3, true
	default:
		return 0, false
	}
}

// Param is a instruction parameter.
type Param interface {
	Display(st *SymbolTable) string
	Evaluate() int16
}

// Params is a sequence of instruction parameters.
type Params []Param

// Display displays the parameters.
func (p Params) Display(st *SymbolTable) string {
	var str strings.Builder
	for i, param := range p {
		if i > 0 {
			str.WriteString(", ")
		}
		str.WriteString(param.Display(st))
	}
	return str.String()
}

// Instruction is an instruction of the bytecode scripting language.
type Instruction interface{}
