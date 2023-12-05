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

	// Frames is the list of bytecode frames. Only available if the script was decoded.
	Frames []vm.BytecodeFrame
}

// Decode decodes the script bytecode using the given instruction decoder.
func (s *Script) Decode(dec vm.InstructionDecoder) (err error) {
	r := vm.NewBytecodeReader(bytes.NewReader(s.Bytecode))
	for {
		r.BeginFrame()
		inst, err := dec(r)
		if err != nil {
			return err
		}
		frame, err := r.EndFrame()
		s.Code = append(s.Code, inst)
		s.Frames = append(s.Frames, frame)
		if err == io.EOF {
			return nil
		}
	}
}

// Listing prints the script listing to the given writer.
func (s Script) Listing(st *vm.SymbolTable, w io.Writer) error {
	var text bytes.Buffer
	for i, inst := range s.Code {
		frame := s.Frames[i]
		label, ok := st.LookupSymbol(vm.SymbolTypeLabel, frame.StartAddress, false)
		if ok {
			label += ":"
		}
		line := fmt.Sprintf("%- 12s%s", label, vm.DisplayInstruction(st, inst))
		if err := frame.Display(&text, line); err != nil {
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
