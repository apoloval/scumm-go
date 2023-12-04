package scumm

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/apoloval/scumm-go/vm"
)

// ScriptID is the ID of a script.
type ScriptID int

// ParseScriptID parses a string into a script ID.
func ParseScriptID(s string) (ScriptID, error) {
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid script ID: %w", err)
	}
	return ScriptID(id), nil
}

// Script is a piece of SCUMM bytecode that can be executed by the game engine.
type Script struct {
	// ID is the script ID.
	ID ScriptID

	// Bytecode is the raw bytecode.
	Bytecode []byte

	// Code is the decoded bytecode. Only available if the script was decoded.
	Code []vm.Instruction
}

// Listing prints the script listing to the given writer.
func (s Script) Listing(st *vm.SymbolTable, w io.Writer) error {
	var text bytes.Buffer
	for _, i := range s.Code {
		label, ok := st.LookupSymbol(vm.SymbolTypeLabel, i.Frame().StartAddress, false)
		if ok {
			label += ":"
		}
		line := fmt.Sprintf("%- 12s%s", label, i.Mnemonic(st))
		if err := i.Frame().Display(&text, line); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintf(w, "Script %d: %d bytes\n", s.ID, len(s.Bytecode)); err != nil {
		return err
	}

	if err := st.Listing(w); err != nil {
		return err
	}

	fmt.Fprintf(w, "\nCode text:\n")
	_, err := io.Copy(w, &text)
	return err
}
