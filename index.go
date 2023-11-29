package scumm

import (
	"encoding/binary"
	"fmt"
	"io"
	"sort"
	"strings"
)

// RoomNumber is the number of a room in the game.
type RoomNumber uint8

// RoomName is the name of a room in the game.
type RoomName [9]byte

func (name RoomName) String() string {
	return strings.Trim(string(name[:]), "\x00")
}

// RoomIndex is the index of the room resources.
type RoomIndex struct {
	// Name is the room name, if given.
	Name RoomName

	// FileNumber is the number of the disk data file. Or 0 if not used.
	FileNumber uint8

	// FileOffset is the offset respect the beginning of the disk data file where room data is
	// located.
	FileOffset uint32

	// ScriptOffsets are the offsets of the global scripts from the beginning of the room (LF) data
	// section.
	ScriptOffsets []uint32

	// SoundOffsets are the offsets of the global sounds from the beginning of the room (LF) data
	// section.
	SoundOffsets []uint32

	// CostumeOffsets are the offsets of the room costumes from the beginning of the room (LF) data
	// section.
	CostumeOffsets []uint32
}

type ObjectClass uint32

const (
	ObjectClassNone        ObjectClass = 0
	ObjectClassYFlip       ObjectClass = 18
	ObjectClassXFlip       ObjectClass = 19
	ObjectClassNeverClip   ObjectClass = 20
	ObjectClassAlwaysClip  ObjectClass = 21
	ObjectClassIgnoreBoxes ObjectClass = 22
	ObjectClassPlayer      ObjectClass = 23 // Actor is controlled by the player
	ObjectClassUntouchable ObjectClass = 24
)

type ObjectIndex struct {
	Class ObjectClass
	Owner byte
	State byte
}

type Index struct {
	// Rooms is the indexed data of rooms
	Rooms map[RoomNumber]RoomIndex

	Objects []ObjectIndex
}

func DecodeIndexFile(r io.Reader) (index Index, err error) {
	// TODO: detect SCUMM version from the index data. For now, assume v4.

	index.Rooms = make(map[RoomNumber]RoomIndex)
	for {

		var blockSize uint32
		var blockName [2]byte
		if err := binary.Read(r, binary.LittleEndian, &blockSize); err != nil {
			if err == io.EOF {
				err = nil
			}
			return index, err
		}
		if err := binary.Read(r, binary.LittleEndian, &blockName); err != nil {
			return index, err
		}

		blockSize -= 6 // ignore block header size
		switch string(blockName[:]) {
		case "RN":
			if err := index.decodeRoomNames(r, int(blockSize)); err != nil {
				return index, err
			}
		case "0R":
			if err := index.decodeDirectoryOfRooms(r, int(blockSize)); err != nil {
				return index, err
			}
		case "0S":
			if err := index.decodeDirectoryOfScripts(r, int(blockSize)); err != nil {
				return index, err
			}
		case "0N":
			if err := index.decodeDirectoryOfSounds(r, int(blockSize)); err != nil {
				return index, err
			}
		case "0C":
			if err := index.decodeDirectoryOfCostumes(r, int(blockSize)); err != nil {
				return index, err
			}
		case "0O":
			if err := index.decodeDirectoryOfObjects(r, int(blockSize)); err != nil {
				return index, err
			}
		default:
			return index, fmt.Errorf("unknown index block type: %s", blockName)
		}
	}
}

// VisitRooms iterates over all rooms in the index and calls the given function for each of them.
func (index *Index) VisitRooms(fn func(RoomNumber, *RoomIndex)) {
	keys := make([]int, 0, len(index.Rooms))
	for key := range index.Rooms {
		keys = append(keys, int(key))
	}
	sort.Ints(keys)
	for _, key := range keys {
		num := RoomNumber(key)
		room := index.Rooms[num]
		fn(num, &room)
	}
}

func (index *Index) decodeRoomNames(r io.Reader, size int) (err error) {
	var nread int
	for {
		var number RoomNumber
		var name RoomName
		if err := binary.Read(r, binary.LittleEndian, &number); err != nil {
			return err
		}
		nread += 1

		if number == 0x00 {
			if nread != size {
				return fmt.Errorf(
					"invalid room names block size: %d expected, %d read", size, nread)
			}
			return nil
		}

		if err := binary.Read(r, binary.LittleEndian, &name); err != nil {
			return err
		}
		nread += len(name)
		for i := 0; i < len(name); i++ {
			name[i] = name[i] ^ 0xFF
		}

		index.updateRoom(number, func(room *RoomIndex) {
			room.Name = name
		})
	}
}

func (index *Index) decodeDirectoryOfRooms(r io.Reader, size int) (err error) {
	return index.decodeDirectoryOfResources(r, size, func(idx int, p1 uint8, p2 uint32) {
		// For unknown reasons, this directory usually has a fixed size of 100 entries. No matter if
		// the game doesn't use them all. The remaning entries are zero-filled. Thus, we ignore any
		// entry whose disk ID is zero.
		if p1 != 0 {
			index.updateRoom(RoomNumber(idx), func(room *RoomIndex) {
				room.FileNumber = p1
				room.FileOffset = p2
			})
		}
	})
}

func (index *Index) decodeDirectoryOfScripts(r io.Reader, size int) (err error) {
	return index.decodeDirectoryOfResources(r, size, func(idx int, p1 uint8, p2 uint32) {
		// For unknown reasons, this directory usually has a fixed size of 200 entries. No matter if
		// the game doesn't use them all. The remaning entries are zero-filled. Thus, we ignore any
		// entry whose room ID is zero.
		if p1 != 0 {
			index.updateRoom(RoomNumber(p1), func(room *RoomIndex) {
				room.ScriptOffsets = append(room.ScriptOffsets, p2)
			})
		}
	})
}

func (index *Index) decodeDirectoryOfSounds(r io.Reader, size int) (err error) {
	return index.decodeDirectoryOfResources(r, size, func(idx int, p1 uint8, p2 uint32) {
		// For unknown reasons, this directory usually has a fixed size of 200 entries. No matter if
		// the game doesn't use them all. The remaning entries are zero-filled. Thus, we ignore any
		// entry whose room ID is zero.
		if p1 != 0 {
			index.updateRoom(RoomNumber(p1), func(room *RoomIndex) {
				room.SoundOffsets = append(room.SoundOffsets, p2)
			})
		}
	})
}

func (index *Index) decodeDirectoryOfCostumes(r io.Reader, size int) (err error) {
	return index.decodeDirectoryOfResources(r, size, func(idx int, p1 uint8, p2 uint32) {
		// For unknown reasons, this directory usually has a fixed size of 200 entries. No matter if
		// the game doesn't use them all. The remaning entries are zero-filled. Thus, we ignore any
		// entry whose room ID is zero.
		if p1 != 0 {
			index.updateRoom(RoomNumber(p1), func(room *RoomIndex) {
				room.CostumeOffsets = append(room.CostumeOffsets, p2)
			})
		}
	})
}

func (index *Index) decodeDirectoryOfObjects(r io.Reader, size int) (err error) {
	var nread int
	var numberOfItems uint16
	if err := binary.Read(r, binary.LittleEndian, &numberOfItems); err != nil {
		return err
	}
	nread += 2

	for i := 1; i <= int(numberOfItems); i++ {
		var entry struct {
			Class      [3]byte
			OwnerState byte
		}
		if err := binary.Read(r, binary.LittleEndian, &entry); err != nil {
			return err
		}
		nread += 4

		object := ObjectIndex{
			Class: ObjectClass(
				uint32(entry.Class[0]) | uint32(entry.Class[1])<<8 | uint32(entry.Class[2])<<16,
			),
			Owner: (entry.OwnerState & 0xF0) >> 4,
			State: entry.OwnerState & 0x0F,
		}
		index.Objects = append(index.Objects, object)
	}
	if nread != size {
		return fmt.Errorf(
			"invalid directory of objects block size: %d expected, %d read", size, nread)
	}
	return nil
}

func (index *Index) decodeDirectoryOfResources(
	r io.Reader,
	size int,
	fn func(idx int, p1 uint8, p2 uint32),
) error {
	var nread int
	var numberOfItems uint16
	if err := binary.Read(r, binary.LittleEndian, &numberOfItems); err != nil {
		return err
	}
	nread += 2

	for i := 0; i < int(numberOfItems); i++ {
		var entry struct {
			P1 uint8
			P2 uint32
		}
		if err := binary.Read(r, binary.LittleEndian, &entry); err != nil {
			return err
		}
		nread += 5

		fn(i, entry.P1, entry.P2)
	}
	if nread != size {
		return fmt.Errorf(
			"invalid directory block size: %d expected, %d read", size, nread)
	}
	return nil
}

func (index *Index) updateRoom(roomNumber RoomNumber, update func(*RoomIndex)) {
	room := index.Rooms[roomNumber]
	update(&room)
	index.Rooms[roomNumber] = room
}
