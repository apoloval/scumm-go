package vm

import (
	"fmt"
	"unicode"
)

// NumberFormat are formatting options for displaying a parameter.
type NumberFormat int

const (
	// ParamFormatNumber displays the parameter as a number.
	ParamFormatNumber NumberFormat = iota

	// ParamFormatChar displays the parameter as a character.
	ParamFormatChar

	// ParamFormatVarID displays the parameter as a variable resource ID.
	ParamFormatVarID

	// ParamFormatStringID displays the parameter as a string resource ID.
	ParamFormatStringID

	// ParamFormatCharsetID displays the parameter as a charset resource ID.
	ParamFormatCharsetID

	// ParamFormatSoundID displays the parameter as a sound resource ID.
	ParamFormatSoundID

	// ParamFormatRoomID displays the parameter as a room resource ID.
	ParamFormatRoomID

	// ParamFormatScriptID displays the parameter as a script resource ID.
	ParamFormatScriptID

	// ParamFormatCostumeID displays the parameter as a costume resource ID.
	ParamFormatCostumeID

	// ParamFormatProgramAddress displays the parameter as a program address.
	ParamFormatProgramAddress
)

// Constant is a constant value referenced from the bytecode.
type Constant struct {
	Value  int16
	Format NumberFormat
}

// Const creates a new word constant with number format.
func Const(v int16) Constant {
	return Constant{
		Value:  v,
		Format: ParamFormatNumber,
	}
}

// Add adds a value to the constant.
func (c Constant) Add(v int16) Constant {
	return Constant{
		Value:  c.Value + v,
		Format: c.Format,
	}
}

// Evaluate implements the Param interface.
func (c Constant) Evaluate() int16 {
	return int16(c.Value)
}

// Display implements the Param interface.
func (c Constant) Display(st *SymbolTable) (str string) {
	switch c.Format {
	case ParamFormatChar:
		value := rune(c.Value & 0xFF)
		if unicode.IsGraphic(rune(value)) {
			str = fmt.Sprintf("'%c'", value)
		} else {
			str = fmt.Sprintf("'\\%02X'", value)
		}
	case ParamFormatStringID:
		str, _ = st.LookupSymbol(SymbolTypeString, uint16(c.Value), true)
	case ParamFormatCharsetID:
		str, _ = st.LookupSymbol(SymbolTypeCharset, uint16(c.Value), true)
	case ParamFormatSoundID:
		str, _ = st.LookupSymbol(SymbolTypeSound, uint16(c.Value), true)
	case ParamFormatRoomID:
		str, _ = st.LookupSymbol(SymbolTypeRoom, uint16(c.Value), true)
	case ParamFormatScriptID:
		str, _ = st.LookupSymbol(SymbolTypeScript, uint16(c.Value), true)
	case ParamFormatCostumeID:
		str, _ = st.LookupSymbol(SymbolTypeCostume, uint16(c.Value), true)
	case ParamFormatVarID:
		str, _ = st.LookupSymbol(SymbolTypeVar, uint16(c.Value), true)
	case ParamFormatProgramAddress:
		str, _ = st.LookupSymbol(SymbolTypeLabel, uint16(c.Value), true)
	default:
		str = fmt.Sprintf("%d", int16(c.Value))
	}
	return
}
