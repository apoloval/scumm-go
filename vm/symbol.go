package vm

import (
	"fmt"
	"io"

	"github.com/apoloval/scumm-go/collections"
)

type SymbolType string

const (
	SymbolTypeVar     SymbolType = "VAR"
	SymbolTypeBit     SymbolType = "BIT"
	SymbolTypeLocal   SymbolType = "LOCAL"
	SymbolTypeLabel   SymbolType = "LABEL"
	SymbolTypeString  SymbolType = "STRING"
	SymbolTypeCharset SymbolType = "CHARSET"
	SymbolTypeSound   SymbolType = "SOUND"
	SymbolTypeRoom    SymbolType = "ROOM"
	SymbolTypeScript  SymbolType = "SCRIPT"
	SymbolTypeCostume SymbolType = "COSTUME"
	SymbolTypeActor   SymbolType = "ACTOR"
	SymbolTypeObject  SymbolType = "OBJECT"
	SymbolTypeClass   SymbolType = "CLASS"
	SymbolTypeMusic   SymbolType = "MUSIC"
	SymbolTypeVerb    SymbolType = "VERB"
)

type SymbolTable struct {
	values  map[SymbolType]map[string]uint16
	symbols map[SymbolType]map[uint16]string
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		values: map[SymbolType]map[string]uint16{
			SymbolTypeVar:     make(map[string]uint16),
			SymbolTypeBit:     make(map[string]uint16),
			SymbolTypeLocal:   make(map[string]uint16),
			SymbolTypeLabel:   make(map[string]uint16),
			SymbolTypeString:  make(map[string]uint16),
			SymbolTypeCharset: make(map[string]uint16),
			SymbolTypeSound:   make(map[string]uint16),
			SymbolTypeRoom:    make(map[string]uint16),
			SymbolTypeScript:  make(map[string]uint16),
			SymbolTypeCostume: make(map[string]uint16),
			SymbolTypeActor:   make(map[string]uint16),
			SymbolTypeObject:  make(map[string]uint16),
			SymbolTypeClass:   make(map[string]uint16),
			SymbolTypeMusic:   make(map[string]uint16),
			SymbolTypeVerb:    make(map[string]uint16),
		},
		symbols: map[SymbolType]map[uint16]string{
			SymbolTypeVar:     make(map[uint16]string),
			SymbolTypeBit:     make(map[uint16]string),
			SymbolTypeLocal:   make(map[uint16]string),
			SymbolTypeLabel:   make(map[uint16]string),
			SymbolTypeString:  make(map[uint16]string),
			SymbolTypeCharset: make(map[uint16]string),
			SymbolTypeSound:   make(map[uint16]string),
			SymbolTypeRoom:    make(map[uint16]string),
			SymbolTypeScript:  make(map[uint16]string),
			SymbolTypeCostume: make(map[uint16]string),
			SymbolTypeActor:   make(map[uint16]string),
			SymbolTypeObject:  make(map[uint16]string),
			SymbolTypeClass:   make(map[uint16]string),
			SymbolTypeMusic:   make(map[uint16]string),
			SymbolTypeVerb:    make(map[uint16]string),
		},
	}
}

func (st *SymbolTable) Declare(t SymbolType, name string, value uint16) *SymbolTable {
	st.symbols[t][value] = name
	st.values[t][name] = value
	return st
}

func (st *SymbolTable) LookupValue(t SymbolType, name string) (uint16, bool) {
	v, ok := st.values[t][name]
	return v, ok
}

func (st *SymbolTable) LookupSymbol(t SymbolType, value uint16, create bool) (string, bool) {
	sym, ok := st.symbols[t][value]
	if !ok && create {
		sym = fmt.Sprintf(symbolTypeFormats[t], value)
		st.Declare(t, sym, value)
	}
	return sym, ok
}

func (st *SymbolTable) SymbolsOf(t SymbolType) map[string]uint16 {
	return st.values[t]
}

func (st SymbolTable) Listing(w io.Writer) error {
	tables := []struct {
		name   string
		values map[uint16]string
	}{
		{"Word variables", st.symbols[SymbolTypeVar]},
		{"Bit variables", st.symbols[SymbolTypeBit]},
		{"Local variables", st.symbols[SymbolTypeLocal]},
		{"Labels", st.symbols[SymbolTypeLabel]},
		{"Strings", st.symbols[SymbolTypeString]},
		{"Charsets", st.symbols[SymbolTypeCharset]},
		{"Sounds", st.symbols[SymbolTypeSound]},
		{"Rooms", st.symbols[SymbolTypeRoom]},
		{"Scripts", st.symbols[SymbolTypeScript]},
		{"Costumes", st.symbols[SymbolTypeCostume]},
		{"Actors", st.symbols[SymbolTypeActor]},
		{"Objects", st.symbols[SymbolTypeObject]},
		{"Classes", st.symbols[SymbolTypeClass]},
		{"Music", st.symbols[SymbolTypeMusic]},
		{"Verbs", st.symbols[SymbolTypeVerb]},
	}
	for _, table := range tables {
		if len(table.values) > 0 {
			fmt.Fprintf(w, "\n%s:\n", table.name)
			collections.VisitMap(table.values, func(addr uint16, name string) {
				fmt.Fprintf(w, "%04X: \t%s\n", addr, name)
			})
		}
	}
	return nil
}

var symbolTypeFormats = map[SymbolType]string{
	SymbolTypeVar:     "VAR_%d",
	SymbolTypeBit:     "BIT_%d",
	SymbolTypeLocal:   "LOCAL_%d",
	SymbolTypeLabel:   "LABEL_%04X",
	SymbolTypeString:  "STRING_%d",
	SymbolTypeCharset: "CHARSET_%d",
	SymbolTypeSound:   "SOUND_%d",
	SymbolTypeRoom:    "ROOM_%d",
	SymbolTypeScript:  "SCRIPT_%d",
	SymbolTypeCostume: "COSTUME_%d",
	SymbolTypeActor:   "ACTOR_%d",
	SymbolTypeObject:  "OBJECT_%d",
	SymbolTypeClass:   "CLASS_%d",
	SymbolTypeMusic:   "MUSIC_%d",
	SymbolTypeVerb:    "VERB_%d",
}
