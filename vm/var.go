package vm

import "fmt"

// VarRef is a reference to a VM variable.
type VarRef struct {
	VarID  uint16
	Offset uint16
}

// IsWordVar returns true if r references a word variable.
func (r VarRef) IsWordVar() bool {
	return r.VarID&0xE000 == 0
}

// IsLocalVar returns true if r references a bit variable.
func (r VarRef) IsBitVar() bool {
	return r.VarID&0x8000 != 0
}

// IsLocalVar returns true if r references a local variable.
func (r VarRef) IsLocalVar() bool {
	return r.VarID&0xFFF0 == 0x4000 || r.VarID&0xFFF0 == 0x6000
}

// IsIndirectWord returns true if r references a word variable indirectly.
func (r VarRef) IsIndirectWord() bool {
	return r.VarID&0x2000 != 0 && r.Offset&0x2000 == 0
}

// IsIndirectDerefWord returns true if r references a word variable indirectly with dereference.
func (r VarRef) IsIndirectDerefWord() bool {
	return r.VarID&0x2000 != 0 && r.Offset&0x2000 != 0
}

// Display implements the Param interface.
func (r VarRef) Display(st *SymbolTable) string {
	switch {
	case r.IsWordVar():
		sym, _ := st.LookupSymbol(SymbolTypeVar, r.VarID&0x1FFF, true)
		return sym
	case r.IsBitVar():
		sym, _ := st.LookupSymbol(SymbolTypeBit, r.VarID&0x7FFF, true)
		return sym
	case r.IsLocalVar():
		sym, _ := st.LookupSymbol(SymbolTypeLocal, r.VarID&0x000F, true)
		return sym
	case r.IsIndirectWord():
		sym, _ := st.LookupSymbol(SymbolTypeVar, r.VarID&0x1FFF, true)
		return fmt.Sprintf("%s[%d]", sym, r.Offset&0xFFF)
	case r.IsIndirectDerefWord():
		sym, _ := st.LookupSymbol(SymbolTypeVar, r.VarID&0x1FFF, true)
		ind := VarRef{VarID: r.Offset & 0xDFFF}
		return fmt.Sprintf("%s[%s]", sym, ind.Display(st))
	default:
		panic(fmt.Errorf("unknown variable reference %04X", r.VarID))
	}
}

// Evaluate implements the Param interface.
func (r VarRef) Evaluate(_ ExecutionContext) int {
	panic("not implemented")
}
