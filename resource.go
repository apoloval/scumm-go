package scumm

import (
	"fmt"
	"os"
	"path"
	"strconv"
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

	// GetScript returns a script from its ID.
	GetScript(id ScriptID) (*Script, error)
}

// FromIndex creates a resource manager from an index file.
func FromIndexFile(f string) (ResourceManager, error) {
	indexFile, err := os.Open(f)
	if err != nil {
		return nil, fmt.Errorf("failed to open index file: %w", err)
	}

	rt := DetectResourceFile(indexFile)
	switch rt {
	case ResourceFileIndexV4:
		index, err := DecodeIndexV4(indexFile)
		if err != nil {
			return nil, err
		}
		return NewResourceManagerV4(path.Dir(f), index), nil
	default:
		return nil, fmt.Errorf("invalid index resource: unexpected %s", rt)
	}
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
