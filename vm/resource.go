package vm

import (
	"fmt"
	"strconv"
)

// ResourceFileType is the type of a resource file.
type ResourceFileType string

const (
	// ResourceFileUknown is an unknown resource file.
	ResourceFileUknown ResourceFileType = "unknown resource file"
)

// ChunkOffset is the offset of a block respect its parent in a resource file.
type ChunkOffset uint32

// String returns the string representation of the block offset.
func (offset ChunkOffset) String() string {
	return fmt.Sprintf("$%08x", uint32(offset))
}

// ResourceManager is a manager for SCUMM resources.
type ResourceManager interface {
	// GetRoom returns a room from its ID.
	GetRoom(id RoomID) (*Room, error)

	// GetRoomByName returns a room from its name.
	GetRoomByName(name RoomName) (*Room, error)

	// GetScript returns a script from its ID. If decode is true, the script bytecode is decoded.
	GetScript(id ScriptID, decode bool) (*Script, error)
}

// GetRoomFromRef returns a room from a reference in a string form that can be either a room ID or a
// room name.
func GetRoomFromRef(man ResourceManager, ref string) (*Room, error) {
	id, err := strconv.Atoi(ref)
	if err == nil {
		return man.GetRoom(RoomID(id))
	}

	name, err := ParseRoomName(ref)
	if err != nil {
		return nil, err
	}
	return man.GetRoomByName(name)
}
