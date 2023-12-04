package vm

import (
	"fmt"
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

// ParamFlags is a bitmask of the parameter properties.
type ParamFlags int

const (
	// ParamFlagsNone indicates that the parameter has no special properties.
	ParamFlagsNone ParamFlags = 0x0000

	// ParamFlagsIsChar indicates that the parameter has to be represented as a character.
	ParamFlagsIsChar ParamFlags = 0x0001

	// ParamFlagsResString indicates that the parameter is a string resource.
	ParamFlagsResString ParamFlags = 0x0002
)

// Has returns true if the parameter has the given flags.
func (f ParamFlags) Has(other ParamFlags) bool {
	return f&other > 0
}

// Param is a instruction parameter.
type Param interface {
	Display(st *SymbolTable, flags ParamFlags) string
	Evaluate() int16
}

// Params is a sequence of instruction parameters.
type Params []Param

// Instruction is an instruction of the bytecode scripting language.
type Instruction interface {
	Decode(opcode OpCode, r *BytecodeReader) error
	Frame() BytecodeFrame
	Mnemonic(st *SymbolTable) string
	Params() Params
}

// ByteConstant is a constant byte value referenced from the bytecode.
type ByteConstant byte

// Evaluate implements the Param interface.
func (c ByteConstant) Evaluate() int16 {
	return int16(c)
}

// Display implements the Param interface.
func (c ByteConstant) Display(st *SymbolTable, flags ParamFlags) string {
	if flags.Has(ParamFlagsIsChar) {
		return fmt.Sprintf("'%c'", c)
	}
	if flags.Has(ParamFlagsResString) {
		sym, _ := st.LookupSymbol(SymbolTypeString, uint16(c), true)
		return sym
	}
	return fmt.Sprintf("%d", c)
}

// WordConstant is a constant word value referenced from the bytecode.
type WordConstant int16

// Evaluate implements the Param interface.
func (c WordConstant) Evaluate() int16 {
	return int16(c)
}

// Display implements the Param interface.
func (c WordConstant) Display(st *SymbolTable, flags ParamFlags) string {
	return fmt.Sprintf("%d", int16(c))
}

// Pointer is a pointer to a word, local or bit variable referenced from the bytecode.
type Pointer interface {
	Param
	Address() uint16
}

// WordPointer is a pointer to a word variable.
type WordPointer uint16

func (p WordPointer) Evaluate() int16 {
	panic("not implemented")
}

// Address returns the address of the pointer.
func (p WordPointer) Address() uint16 {
	return uint16(p) & 0x1FFF
}

// Display returns the symbol of the pointer.
func (p WordPointer) Display(st *SymbolTable, flags ParamFlags) string {
	sym, _ := st.LookupSymbol(SymbolTypeVar, uint16(p), true)
	return sym
}

// BitPointer is a pointer to a bit variable.
type BitPointer uint16

func (p BitPointer) Evaluate() int16 {
	panic("not implemented")
}

// Address returns the address of the pointer.
func (p BitPointer) Address() uint16 {
	return uint16(p) & 0x7FFF
}

// Display returns the symbol of the pointer.
func (p BitPointer) Display(st *SymbolTable, flags ParamFlags) string {
	sym, _ := st.LookupSymbol(SymbolTypeBit, uint16(p), true)
	return sym
}

// LocalPointer is a pointer to a local variable.
type LocalPointer uint8

func (p LocalPointer) Evaluate() int16 {
	panic("not implemented")
}

// Address returns the address of the pointer.
func (p LocalPointer) Address() uint16 {
	return uint16(p) & 0x0F
}

// Display returns the symbol of the pointer.
func (p LocalPointer) Display(st *SymbolTable, flags ParamFlags) string {
	sym, _ := st.LookupSymbol(SymbolTypeLocal, uint16(p), true)
	return sym
}

// ProgramAddress is a location in the program address space.
type ProgramAddress uint16

// Add returns the program address incremented by v.
func (p ProgramAddress) Add(v int16) ProgramAddress {
	return ProgramAddress(int16(p) + v)
}

// Display returns the symbol of the program address.
func (p ProgramAddress) Display(st *SymbolTable, flags ParamFlags) string {
	sym, _ := st.LookupSymbol(SymbolTypeLabel, uint16(p), true)
	return sym
}
