package scumm

import (
	"encoding/binary"
	"fmt"
	"io"
)

// IndexedRoom is the a room indexed in the directory of rooms.
type IndexedRoom struct {
	// ID is the room ID.
	ID RoomID

	// Name is the room name, if given.
	Name RoomName

	// FileNumber is the number of the disk data file. Or 0 if not used.
	FileNumber uint8

	// FileOffset is the offset respect the beginning of the disk data file where room data is
	// located.
	FileOffset uint32

	// Scripts are the indexed scripts that belong to this room.
	Scripts []IndexedScript

	// Sounds are the indexed sounds that belong to this room.
	Sounds []IndexedSound

	// Costumes are the indexed costumes that belong to this room.
	Costumes []IndexedCostume
}

// IndexedScript is the a script indexed in the directory of scripts.
type IndexedScript struct {
	ID     ScriptID
	Room   RoomID
	Offset BlockOffset
}

// IndexedSound is the a sound indexed in the directory of sounds.
type IndexedSound struct {
	ID     SoundID
	Room   RoomID
	Offset BlockOffset
}

// IndexedCostume is the a costume indexed in the directory of costumes.
type IndexedCostume struct {
	ID     CostumeID
	Room   RoomID
	Offset BlockOffset
}

// IndexedObject is the an object indexed in the directory of objects.
type IndexedObject struct {
	ID    ObjectID
	Class ObjectClass
	Owner ObjectOwner
	State ObjectState
}

type Index struct {
	// Rooms is the indexed rooms
	Rooms map[RoomID]IndexedRoom

	// Scripts is the indexed scripts
	Scripts map[ScriptID]IndexedScript

	// Sounds is the indexed sounds
	Sounds map[SoundID]IndexedSound

	// Costumes is the indexed costumes
	Costumes map[CostumeID]IndexedCostume

	// Objects is the indexed objects
	Objects map[ObjectID]IndexedObject
}

func DecodeIndexV4(r io.Reader) (index Index, err error) {
	index.Rooms = make(map[RoomID]IndexedRoom)
	index.Scripts = make(map[ScriptID]IndexedScript)
	index.Sounds = make(map[SoundID]IndexedSound)
	index.Costumes = make(map[CostumeID]IndexedCostume)
	index.Objects = make(map[ObjectID]IndexedObject)
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

func (index *Index) decodeRoomNames(r io.Reader, size int) (err error) {
	var nread int
	for {
		var number RoomID
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

		index.updateRoom(number, func(room *IndexedRoom) {
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
			index.updateRoom(RoomID(idx), func(room *IndexedRoom) {
				room.ID = RoomID(idx)
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
			script := IndexedScript{
				ID:     ScriptID(idx),
				Room:   RoomID(p1),
				Offset: BlockOffset(p2),
			}
			index.Scripts[script.ID] = script
			index.updateRoom(script.Room, func(room *IndexedRoom) {
				room.Scripts = append(room.Scripts, script)
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
			sound := IndexedSound{
				ID:     SoundID(idx),
				Room:   RoomID(p1),
				Offset: BlockOffset(p2),
			}
			index.Sounds[sound.ID] = sound
			index.updateRoom(sound.Room, func(room *IndexedRoom) {
				room.Sounds = append(room.Sounds, sound)
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
			costume := IndexedCostume{
				ID:     CostumeID(idx),
				Room:   RoomID(p1),
				Offset: BlockOffset(p2),
			}
			index.Costumes[costume.ID] = costume
			index.updateRoom(costume.Room, func(room *IndexedRoom) {
				room.Costumes = append(room.Costumes, costume)
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

	for i := 0; i < int(numberOfItems); i++ {
		var entry struct {
			Class      [3]byte
			OwnerState byte
		}
		if err := binary.Read(r, binary.LittleEndian, &entry); err != nil {
			return err
		}
		nread += 4

		object := IndexedObject{
			ID: ObjectID(i),
			Class: ObjectClass(
				uint32(entry.Class[0]) | uint32(entry.Class[1])<<8 | uint32(entry.Class[2])<<16,
			),
			Owner: ObjectOwner(entry.OwnerState & 0x0F),
			State: ObjectState((entry.OwnerState & 0xF0) >> 4),
		}
		index.Objects[object.ID] = object
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

func (index *Index) updateRoom(roomNumber RoomID, update func(*IndexedRoom)) {
	room := index.Rooms[roomNumber]
	update(&room)
	index.Rooms[roomNumber] = room
}
