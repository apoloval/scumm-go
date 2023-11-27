package scumm

import (
	"encoding/binary"
	"fmt"
	"io"
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

	// ScriptOffsets are the offsets respect the beginning of the room data where scripts are found.
	ScriptOffsets []uint32

	// SoundOffsets are the offsets respect the beginning of the room data where sounds are found.
	SoundOffsets []uint32

	// CostumeOffsets are the offsets respect the beginning of the room data where costumes are
	// found.
	CostumeOffsets []uint32
}

type ObjectClass byte

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

	GlobalScripts  []uint32
	GlobalSounds   []uint32
	GlobalCostumes []uint32
}

func DecodeIndexFile(r io.Reader) (index Index, err error) {
	// TODO: detect SCUMM version from the index data. For now, assume v3.

	index.Rooms = make(map[RoomNumber]RoomIndex)
	for {

		var blockSize uint32
		var blockName [2]byte
		if err := binary.Read(r, binary.LittleEndian, &blockSize); err != nil {
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
			return index, nil
			//return index, fmt.Errorf("unknown index block type: %s", blockName)
		}
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
	var nread int
	var numberOfItems uint16
	if err := binary.Read(r, binary.LittleEndian, &numberOfItems); err != nil {
		return err
	}
	nread += 2

	for i := 1; i <= int(numberOfItems); i++ {
		var fileNumber uint8
		var offset uint32
		if err := binary.Read(r, binary.LittleEndian, &fileNumber); err != nil {
			return err
		}
		if err := binary.Read(r, binary.LittleEndian, &offset); err != nil {
			return err
		}
		nread += 5
		index.updateRoom(RoomNumber(i), func(room *RoomIndex) {
			room.FileNumber = fileNumber
			room.FileOffset = offset
		})
	}
	if nread != size {
		return fmt.Errorf(
			"invalid directory of rooms block size: %d expected, %d read", size, nread)
	}
	return nil
}

func (index *Index) decodeDirectoryOfScripts(r io.Reader, size int) (err error) {
	var nread int
	var numberOfItems uint16
	if err := binary.Read(r, binary.LittleEndian, &numberOfItems); err != nil {
		return err
	}
	nread += 2

	for i := 1; i <= int(numberOfItems); i++ {
		var roomNumber RoomNumber
		var scriptOffset uint32
		if err := binary.Read(r, binary.LittleEndian, &roomNumber); err != nil {
			return err
		}
		if err := binary.Read(r, binary.LittleEndian, &scriptOffset); err != nil {
			return err
		}
		nread += 5

		if roomNumber == 0x00 {
			index.GlobalScripts = append(index.GlobalScripts, scriptOffset)
		} else {
			index.updateRoom(roomNumber, func(room *RoomIndex) {
				room.ScriptOffsets = append(room.ScriptOffsets, scriptOffset)
			})
		}
	}
	if nread != size {
		return fmt.Errorf(
			"invalid directory of scripts block size: %d expected, %d read", size, nread)
	}
	return nil
}

func (index *Index) decodeDirectoryOfSounds(r io.Reader, size int) (err error) {
	var nread int
	var numberOfItems uint16
	if err := binary.Read(r, binary.LittleEndian, &numberOfItems); err != nil {
		return err
	}
	nread += 2

	for i := 1; i <= int(numberOfItems); i++ {
		var roomNumber RoomNumber
		var soundOffset uint32
		if err := binary.Read(r, binary.LittleEndian, &roomNumber); err != nil {
			return err
		}
		if err := binary.Read(r, binary.LittleEndian, &soundOffset); err != nil {
			return err
		}
		nread += 5

		if roomNumber == 0x00 {
			index.GlobalSounds = append(index.GlobalSounds, soundOffset)
		} else {
			index.updateRoom(roomNumber, func(room *RoomIndex) {
				room.SoundOffsets = append(room.SoundOffsets, soundOffset)
			})
		}
	}
	if nread != size {
		return fmt.Errorf(
			"invalid directory of sounds block size: %d expected, %d read", size, nread)
	}
	return nil
}

func (index *Index) decodeDirectoryOfCostumes(r io.Reader, size int) (err error) {
	var nread int
	var numberOfItems uint16
	if err := binary.Read(r, binary.LittleEndian, &numberOfItems); err != nil {
		return err
	}
	nread += 2

	for i := 1; i <= int(numberOfItems); i++ {
		var roomNumber RoomNumber
		var costumeOffset uint32
		if err := binary.Read(r, binary.LittleEndian, &roomNumber); err != nil {
			return err
		}
		if err := binary.Read(r, binary.LittleEndian, &costumeOffset); err != nil {
			return err
		}
		nread += 5

		if roomNumber == 0x00 {
			index.GlobalCostumes = append(index.GlobalCostumes, costumeOffset)
		} else {
			index.updateRoom(roomNumber, func(room *RoomIndex) {
				room.CostumeOffsets = append(room.CostumeOffsets, costumeOffset)
			})
		}
	}
	if nread != size {
		return fmt.Errorf(
			"invalid directory of sounds block size: %d expected, %d read", size, nread)
	}
	return nil
}

func (index *Index) decodeDirectoryOfObjects(r io.Reader, size int) (err error) {
	var nread int
	var numberOfItems uint16
	if err := binary.Read(r, binary.LittleEndian, &numberOfItems); err != nil {
		return err
	}
	nread += 2

	for i := 1; i <= int(numberOfItems); i++ {
		var class ObjectClass
		var ownerState byte
		if err := binary.Read(r, binary.LittleEndian, &class); err != nil {
			return err
		}
		if err := binary.Read(r, binary.LittleEndian, &ownerState); err != nil {
			return err
		}
		nread += 4

		object := ObjectIndex{
			Class: class,
			Owner: (ownerState & 0xF0) >> 4,
			State: ownerState & 0x0F,
		}
		index.Objects = append(index.Objects, object)
	}
	if nread != size {
		return fmt.Errorf(
			"invalid directory of objects block size: %d expected, %d read", size, nread)
	}
	return nil
}

func (index *Index) updateRoom(roomNumber RoomNumber, update func(*RoomIndex)) {
	room := index.Rooms[roomNumber]
	update(&room)
	index.Rooms[roomNumber] = room
}
