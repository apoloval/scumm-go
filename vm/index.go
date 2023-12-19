package vm

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
	FileOffset ChunkOffset

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
	Offset ChunkOffset
}

// IndexedSound is the a sound indexed in the directory of sounds.
type IndexedSound struct {
	ID     SoundID
	Room   RoomID
	Offset ChunkOffset
}

// IndexedCostume is the a costume indexed in the directory of costumes.
type IndexedCostume struct {
	ID     CostumeID
	Room   RoomID
	Offset ChunkOffset
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
