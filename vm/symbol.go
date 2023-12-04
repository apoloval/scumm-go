package vm

import (
	"fmt"
	"io"

	"github.com/apoloval/scumm-go/collections"
)

type SymbolTable struct {
	wordVariables  map[string]uint16
	bitVariables   map[string]uint16
	localVariables map[string]uint8
	progAddress    map[string]uint16

	wordVariablesRev map[uint16]string
	progAddressRev   map[uint16]string
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		wordVariables:  make(map[string]uint16),
		bitVariables:   make(map[string]uint16),
		localVariables: make(map[string]uint8),
		progAddress:    make(map[string]uint16),

		progAddressRev:   make(map[uint16]string),
		wordVariablesRev: make(map[uint16]string),
	}
}

func (st *SymbolTable) DeclWordVariable(name string, addr uint16) *SymbolTable {
	st.wordVariables[name] = addr
	st.wordVariablesRev[addr] = name
	return st
}

func (st *SymbolTable) DeclBitVariable(name string, addr uint16) *SymbolTable {
	st.bitVariables[name] = addr
	return st
}

func (st *SymbolTable) DeclLocalVariable(name string, addr uint8) *SymbolTable {
	st.localVariables[name] = addr
	return st
}

func (st *SymbolTable) DeclLabel(name string, addr uint16) *SymbolTable {
	st.progAddress[name] = addr
	st.progAddressRev[addr] = name
	return st
}

func (st SymbolTable) Listing(w io.Writer) error {
	if len(st.wordVariables) > 0 {
		fmt.Fprintf(w, "Word variables:\n")
		collections.VisitMap(st.wordVariablesRev, func(addr uint16, name string) {
			fmt.Fprintf(w, "%04X: \t%s\n", addr, name)
		})
	}
	if len(st.bitVariables) > 0 {
		fmt.Fprintf(w, "Bit variables:\n")
		for name, addr := range st.bitVariables {
			if _, err := fmt.Fprintf(w, "\t%s at %d\n", name, addr); err != nil {
				return err
			}
		}
	}
	if len(st.localVariables) > 0 {
		fmt.Fprintf(w, "Local variables:\n")
		for name, addr := range st.localVariables {
			if _, err := fmt.Fprintf(w, "\t%s at %d\n", name, addr); err != nil {
				return err
			}
		}
	}

	return nil
}

func (st *SymbolTable) WordVariableAt(addr uint16, create bool) string {
	name, ok := st.wordVariablesRev[addr]
	if ok {
		return name
	}
	if create {
		name := fmt.Sprintf("VAR_%d", addr)
		st.wordVariables[name] = addr
		st.wordVariablesRev[addr] = name
		return name
	}
	return ""
}

func (st *SymbolTable) BitVariableAt(addr uint16, create bool) string {
	for name, a := range st.bitVariables {
		if a == addr {
			return name
		}
	}
	if create {
		name := fmt.Sprintf("BIT_%d", addr)
		st.bitVariables[name] = addr
		return name
	}
	return ""
}

func (st *SymbolTable) LocalVariableAt(addr uint8, create bool) string {
	for name, a := range st.localVariables {
		if a == addr {
			return name
		}
	}
	if create {
		name := fmt.Sprintf("LOCAL_%d", addr)
		st.localVariables[name] = addr
		return name
	}
	return ""
}

func (st *SymbolTable) LabelAt(addr uint16, create bool) string {
	name, ok := st.progAddressRev[addr]
	if ok {
		return name
	}
	if create {
		name = fmt.Sprintf("LABEL_%04X", addr)
		st.progAddress[name] = addr
		st.progAddressRev[addr] = name
		return name
	}
	return ""
}
