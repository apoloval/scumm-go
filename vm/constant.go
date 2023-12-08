package vm

import (
	"fmt"
	"unicode"
)

// NumberFormat are formatting options for displaying a parameter.
type NumberFormat string

const (
	// NumberFormatDecimal displays the parameter as a number.
	NumberFormatDecimal NumberFormat = "dec"

	// NumberFormatHex displays the parameter as a hexadecimal number.
	NumberFormatHex NumberFormat = "hex"

	// NumberFormatChar displays the parameter as a character.
	NumberFormatChar NumberFormat = "char"

	// NumberFormatAddress displays the parameter as a program address.
	NumberFormatAddress NumberFormat = "addr"

	// NumberFormatVarID displays the parameter as a variable resource ID.
	NumberFormatVarID NumberFormat = "id:var"

	// NumberFormatStringID displays the parameter as a string resource ID.
	NumberFormatStringID NumberFormat = "id:string"

	// NumberFormatCharsetID displays the parameter as a charset resource ID.
	NumberFormatCharsetID NumberFormat = "id:charset"

	// NumberFormatSoundID displays the parameter as a sound resource ID.
	NumberFormatSoundID NumberFormat = "id:sound"

	// NumberFormatRoomID displays the parameter as a room resource ID.
	NumberFormatRoomID NumberFormat = "id:room"

	// NumberFormatScriptID displays the parameter as a script resource ID.
	NumberFormatScriptID NumberFormat = "id:script"

	// NumberFormatCostumeID displays the parameter as a costume resource ID.
	NumberFormatCostumeID NumberFormat = "id:costume"

	// NumberFormatActorID displays the parameter as a actor resource ID.
	NumberFormatActorID NumberFormat = "id:actor"

	// NumberFormatObjectID displays the parameter as a object resource ID.
	NumberFormatObjectID NumberFormat = "id:object"

	// NumberFormatClassID displays the parameter as a class ID.
	NumberFormatClassID NumberFormat = "id:class"

	// NumberFormatMusicID displays the parameter as a music resource ID.
	NumberFormatMusicID NumberFormat = "id:music"
)

// Constant is a constant value referenced from the bytecode.
type Constant struct {
	Value  int
	Format NumberFormat
}

// Const16 creates a new word constant with number format.
func Const16(v int16) Constant {
	return Constant{
		Value:  int(v),
		Format: NumberFormatDecimal,
	}
}

// Add adds a value to the constant.
func (c Constant) Add(v int16) Constant {
	return Constant{
		Value:  c.Value + int(v),
		Format: c.Format,
	}
}

// Evaluate implements the Param interface.
func (c Constant) Evaluate() int {
	return int(c.Value)
}

// Display implements the Param interface.
func (c Constant) Display(st *SymbolTable) (str string) {
	switch c.Format {
	case NumberFormatDecimal:
		str = fmt.Sprintf("%d", int16(c.Value))
	case NumberFormatHex:
		str = fmt.Sprintf("$%04X", uint16(c.Value))
	case NumberFormatChar:
		value := rune(c.Value & 0xFF)
		if unicode.IsGraphic(rune(value)) {
			str = fmt.Sprintf("'%c'", value)
		} else {
			str = fmt.Sprintf("'\\%02X'", value)
		}
	case NumberFormatAddress:
		str, _ = st.LookupSymbol(SymbolTypeLabel, uint16(c.Value), true)
	case NumberFormatVarID:
		str, _ = st.LookupSymbol(SymbolTypeVar, uint16(c.Value), true)
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
	case NumberFormatActorID:
		str, _ = st.LookupSymbol(SymbolTypeActor, uint16(c.Value), true)
	case NumberFormatObjectID:
		str, _ = st.LookupSymbol(SymbolTypeObject, uint16(c.Value), true)
	case NumberFormatClassID:
		str, _ = st.LookupSymbol(SymbolTypeClass, uint16(c.Value), true)
	case NumberFormatMusicID:
		str, _ = st.LookupSymbol(SymbolTypeMusic, uint16(c.Value), true)
	default:
		panic("invalid number format")
	}
	return
}
