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

// Param is a instruction parameter.
type Param interface {
	Evaluate() int16
}

// Params is a sequence of instruction parameters.
type Params []Param

// Instruction is an instruction of the bytecode scripting language.
type Instruction interface {
	Bytecode() []byte
	Mnemonic() string
	Params() Params
}

// ByteConstant is a constant byte value referenced from the bytecode.
type ByteConstant byte

func (c ByteConstant) Evaluate() int16 {
	return int16(c)
}

// String returns the string representation of the constant.
func (c ByteConstant) String() string {
	return fmt.Sprintf("%d", byte(c))
}

// WordConstant is a constant word value referenced from the bytecode.
type WordConstant int16

func (c WordConstant) Evaluate() int16 {
	return int16(c)
}

// String returns the string representation of the constant.
func (c WordConstant) String() string {
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

// String implements the Stringer interface.
func (p WordPointer) String() string {
	return fmt.Sprintf("word:%04X", p.Address())
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

// String returns the string representation of the pointer.
func (p BitPointer) String() string {
	return fmt.Sprintf("bit:%04X", p.Address())
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

func (p LocalPointer) String() string {
	return fmt.Sprintf("loc:%01X", p.Address())
}
