package vm

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
func (p WordPointer) Display(st *SymbolTable) string {
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
func (p BitPointer) Display(st *SymbolTable) string {
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
func (p LocalPointer) Display(st *SymbolTable) string {
	sym, _ := st.LookupSymbol(SymbolTypeLocal, uint16(p), true)
	return sym
}
