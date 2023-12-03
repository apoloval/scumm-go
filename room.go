package scumm

import (
	"fmt"
	"strings"
)

// RoomID is the number of a room in the game.
type RoomID uint8

// RoomName is the name of a room in the game.
type RoomName [9]byte

// ParseRoomName parses a string into a room name.
func ParseRoomName(str string) (RoomName, error) {
	if len(str) > 8 {
		return RoomName{}, fmt.Errorf("invalid room name: %s", str)
	}
	var name RoomName
	for i := 0; i < len(name); i++ {
		name[i] = 0
	}
	for i := 0; i < len(str); i++ {
		name[i] = str[i]
	}
	return name, nil
}

func (name RoomName) String() string {
	return strings.Trim(string(name[:]), "\x00")
}

// Room is a room in the game.
type Room struct {
	ID                   RoomID
	Name                 RoomName
	Width                uint16
	Height               uint16
	NumberOfObjects      uint16
	NumberOfLocalScripts uint8
	LocalScripts         []Script
}
