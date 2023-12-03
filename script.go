package scumm

import (
	"fmt"
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
	ID       ScriptID
	Bytecode []byte
}
