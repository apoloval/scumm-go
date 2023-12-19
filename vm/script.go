package vm

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
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
	Code []Instruction

	// Frames is the list of bytecode frames. Only available if the script was decoded.
	Frames []BytecodeFrame
}

// Decode decodes the script bytecode using the given instruction decoder.
func (s *Script) Decode(dec InstructionDecoder) (err error) {
	r := NewBytecodeDecoder(bytes.NewReader(s.Bytecode))
	for {
		r.BeginFrame()
		inst, err := dec(r)
		if err != nil {
			frame, _ := r.EndFrame()
			return fmt.Errorf("error decoding bytecode frame [% 3X ]: %w", frame.Bytes, err)
		}
		frame, err := r.EndFrame()
		s.Code = append(s.Code, inst)
		s.Frames = append(s.Frames, frame)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("error decoding bytecode frame [% 3X ]: %w", frame.Bytes, err)
		}
	}
}

// Listing prints the script listing to the given writer.
func (s Script) Listing(st *SymbolTable, w io.Writer) error {
	var text bytes.Buffer
	instructions := make([]string, len(s.Code))
	for i, inst := range s.Code {
		instructions[i] = DisplayInstruction(st, inst)
	}

	if err := s.checkBranchConsistency(st); err != nil {
		return err
	}

	for i, inst := range instructions {
		frame := s.Frames[i]
		label, ok := st.LookupSymbol(SymbolTypeLabel, frame.StartAddress, false)
		if ok {
			label += ":"
		}
		line := fmt.Sprintf("%- 12s%s", label, inst)
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

func (s Script) checkBranchConsistency(st *SymbolTable) error {
	for sym, addr := range st.SymbolsOf(SymbolTypeLabel) {
		if addr := s.instructionOnAddress(addr); addr == nil {
			return fmt.Errorf(
				"Branch consistency check failed: label %s points to invalid address", sym)
		}
	}
	return nil
}

func (s Script) instructionOnAddress(addr uint16) Instruction {
	for i, frame := range s.Frames {
		if frame.StartAddress == addr {
			return s.Code[i]
		}
	}
	return nil
}
