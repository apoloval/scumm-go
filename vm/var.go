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

func (r VarRef) Read(ctx ExecutionContext) int {
	switch {
	case r.IsWordVar():
		return ctx.ReadWord(r.VarID)
	case r.IsBitVar():
		if ctx.ReadBit(r.VarID) {
			return 1
		}
		return 0
	case r.IsLocalVar():
		return ctx.ReadLocal(r.VarID)
	case r.IsIndirectWord():
		return ctx.ReadWord(r.VarID + uint16(r.Offset))
	case r.IsIndirectDerefWord():
		return ctx.ReadWord(r.VarID + uint16(ctx.ReadWord(r.Offset)))
	default:
		panic("unknown variable reference")
	}
}

func (r VarRef) Write(ctx ExecutionContext, value int) {
	switch {
	case r.IsWordVar():
		ctx.WriteWord(r.VarID, value)
	case r.IsBitVar():
		ctx.WriteBit(r.VarID, value != 0)
	case r.IsLocalVar():
		ctx.WriteLocal(r.VarID, value)
	case r.IsIndirectWord():
		ctx.WriteWord(r.VarID+uint16(r.Offset), value)
	case r.IsIndirectDerefWord():
		ctx.WriteWord(r.VarID+uint16(ctx.ReadWord(r.Offset)), value)
	default:
		panic("unknown variable reference")
	}
}

// Evaluate implements the Param interface.
func (r VarRef) Evaluate(ctx ExecutionContext) int {
	return r.Read(ctx)
}
