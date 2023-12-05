package vm

import (
	"fmt"
	"unicode"
)

// NumberFormat are formatting options for displaying a parameter.
type NumberFormat int

const (
	// NumberFormatDecimal displays the parameter as a number.
	NumberFormatDecimal NumberFormat = iota

	// NumberFormatHex displays the parameter as a hexadecimal number.
	NumberFormatHex

	// NumberFormatChar displays the parameter as a character.
	NumberFormatChar

	// NumberFormatVarID displays the parameter as a variable resource ID.
	NumberFormatVarID

	// NumberFormatStringID displays the parameter as a string resource ID.
	NumberFormatStringID

	// NumberFormatCharsetID displays the parameter as a charset resource ID.
	NumberFormatCharsetID

	// NumberFormatSoundID displays the parameter as a sound resource ID.
	NumberFormatSoundID

	// NumberFormatRoomID displays the parameter as a room resource ID.
	NumberFormatRoomID

	// NumberFormatScriptID displays the parameter as a script resource ID.
	NumberFormatScriptID

	// NumberFormatCostumeID displays the parameter as a costume resource ID.
	NumberFormatCostumeID

	// NumberFormatAddress displays the parameter as a program address.
	NumberFormatAddress
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
		Format: NumberFormatDecimal,
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
	case NumberFormatChar:
		value := rune(c.Value & 0xFF)
		if unicode.IsGraphic(rune(value)) {
			str = fmt.Sprintf("'%c'", value)
		} else {
			str = fmt.Sprintf("'\\%02X'", value)
		}
	case NumberFormatHex:
		str = fmt.Sprintf("$%04X", uint16(c.Value))
	case NumberFormatStringID:
		str, _ = st.LookupSymbol(SymbolTypeString, uint16(c.Value), true)
	case NumberFormatCharsetID:
		str, _ = st.LookupSymbol(SymbolTypeCharset, uint16(c.Value), true)
	case NumberFormatSoundID:
		str, _ = st.LookupSymbol(SymbolTypeSound, uint16(c.Value), true)
	case NumberFormatRoomID:
		str, _ = st.LookupSymbol(SymbolTypeRoom, uint16(c.Value), true)
	case NumberFormatScriptID:
		str, _ = st.LookupSymbol(SymbolTypeScript, uint16(c.Value), true)
	case NumberFormatCostumeID:
		str, _ = st.LookupSymbol(SymbolTypeCostume, uint16(c.Value), true)
	case NumberFormatVarID:
		str, _ = st.LookupSymbol(SymbolTypeVar, uint16(c.Value), true)
	case NumberFormatAddress:
		str, _ = st.LookupSymbol(SymbolTypeLabel, uint16(c.Value), true)
	default:
		str = fmt.Sprintf("%d", int16(c.Value))
	}
	return
}
