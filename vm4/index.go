package vm4

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/apoloval/scumm-go/vm"
)

func DecodeIndex(r io.Reader) (index vm.Index, err error) {
	index.Rooms = make(map[vm.RoomID]vm.IndexedRoom)
	index.Scripts = make(map[vm.ScriptID]vm.IndexedScript)
	index.Sounds = make(map[vm.SoundID]vm.IndexedSound)
	index.Costumes = make(map[vm.CostumeID]vm.IndexedCostume)
	index.Objects = make(map[vm.ObjectID]vm.IndexedObject)
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
			if err := decodeRoomNames(&index, r, int(blockSize)); err != nil {
				return index, err
			}
		case "0R":
			if err := decodeDirectoryOfRooms(&index, r, int(blockSize)); err != nil {
				return index, err
			}
		case "0S":
			if err := decodeDirectoryOfScripts(&index, r, int(blockSize)); err != nil {
				return index, err
			}
		case "0N":
			if err := decodeDirectoryOfSounds(&index, r, int(blockSize)); err != nil {
				return index, err
			}
		case "0C":
			if err := decodeDirectoryOfCostumes(&index, r, int(blockSize)); err != nil {
				return index, err
			}
		case "0O":
			if err := decodeDirectoryOfObjects(&index, r, int(blockSize)); err != nil {
				return index, err
			}
		default:
			return index, fmt.Errorf("unknown index block type: %s", blockName)
		}
	}
}

func decodeRoomNames(index *vm.Index, r io.Reader, size int) (err error) {
	var nread int
	for {
		var number vm.RoomID
		var name vm.RoomName
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

		updateRoom(index, number, func(room *vm.IndexedRoom) {
			room.Name = name
		})
	}
}

func decodeDirectoryOfRooms(index *vm.Index, r io.Reader, size int) (err error) {
	return decodeDirectoryOfResources(index, r, size, func(idx int, p1 uint8, p2 uint32) {
		// For unknown reasons, this directory usually has a fixed size of 100 entries. No matter if
		// the game doesn't use them all. The remaning entries are zero-filled. Thus, we ignore any
		// entry whose disk ID is zero.
		if p1 != 0 {
			updateRoom(index, vm.RoomID(idx), func(room *vm.IndexedRoom) {
				room.ID = vm.RoomID(idx)
				room.FileNumber = p1
				room.FileOffset = vm.ChunkOffset(p2)
			})
		}
	})
}

func decodeDirectoryOfScripts(index *vm.Index, r io.Reader, size int) (err error) {
	return decodeDirectoryOfResources(index, r, size, func(idx int, p1 uint8, p2 uint32) {
		// For unknown reasons, this directory usually has a fixed size of 200 entries. No matter if
		// the game doesn't use them all. The remaning entries are zero-filled. Thus, we ignore any
		// entry whose room ID is zero.
		if p1 != 0 {
			script := vm.IndexedScript{
				ID:     vm.ScriptID(idx),
				Room:   vm.RoomID(p1),
				Offset: vm.ChunkOffset(p2),
			}
			index.Scripts[script.ID] = script
			updateRoom(index, script.Room, func(room *vm.IndexedRoom) {
				room.Scripts = append(room.Scripts, script)
			})
		}
	})
}

func decodeDirectoryOfSounds(index *vm.Index, r io.Reader, size int) (err error) {
	return decodeDirectoryOfResources(index, r, size, func(idx int, p1 uint8, p2 uint32) {
		// For unknown reasons, this directory usually has a fixed size of 200 entries. No matter if
		// the game doesn't use them all. The remaning entries are zero-filled. Thus, we ignore any
		// entry whose room ID is zero.
		if p1 != 0 {
			sound := vm.IndexedSound{
				ID:     vm.SoundID(idx),
				Room:   vm.RoomID(p1),
				Offset: vm.ChunkOffset(p2),
			}
			index.Sounds[sound.ID] = sound
			updateRoom(index, sound.Room, func(room *vm.IndexedRoom) {
				room.Sounds = append(room.Sounds, sound)
			})
		}
	})
}

func decodeDirectoryOfCostumes(index *vm.Index, r io.Reader, size int) (err error) {
	return decodeDirectoryOfResources(index, r, size, func(idx int, p1 uint8, p2 uint32) {
		// For unknown reasons, this directory usually has a fixed size of 200 entries. No matter if
		// the game doesn't use them all. The remaning entries are zero-filled. Thus, we ignore any
		// entry whose room ID is zero.
		if p1 != 0 {
			costume := vm.IndexedCostume{
				ID:     vm.CostumeID(idx),
				Room:   vm.RoomID(p1),
				Offset: vm.ChunkOffset(p2),
			}
			index.Costumes[costume.ID] = costume
			updateRoom(index, costume.Room, func(room *vm.IndexedRoom) {
				room.Costumes = append(room.Costumes, costume)
			})
		}
	})
}

func decodeDirectoryOfObjects(index *vm.Index, r io.Reader, size int) (err error) {
	var nread int
	var numberOfItems uint16
	if err := binary.Read(r, binary.LittleEndian, &numberOfItems); err != nil {
		return err
	}
	nread += 2

	for i := 0; i < int(numberOfItems); i++ {
		var entry struct {
			Class      [3]byte
			OwnerState byte
		}
		if err := binary.Read(r, binary.LittleEndian, &entry); err != nil {
			return err
		}
		nread += 4

		object := vm.IndexedObject{
			ID: vm.ObjectID(i),
			Class: vm.ObjectClass(
				uint32(entry.Class[0]) | uint32(entry.Class[1])<<8 | uint32(entry.Class[2])<<16,
			),
			Owner: vm.ObjectOwner(entry.OwnerState & 0x0F),
			State: vm.ObjectState((entry.OwnerState & 0xF0) >> 4),
		}
		index.Objects[object.ID] = object
	}
	if nread != size {
		return fmt.Errorf(
			"invalid directory of objects block size: %d expected, %d read", size, nread)
	}
	return nil
}

func decodeDirectoryOfResources(
	index *vm.Index,
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

func updateRoom(index *vm.Index, roomNumber vm.RoomID, update func(*vm.IndexedRoom)) {
	room := index.Rooms[roomNumber]
	update(&room)
	index.Rooms[roomNumber] = room
}
