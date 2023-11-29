package scumm

import "strings"

// RoomID is the number of a room in the game.
type RoomID uint8

// RoomName is the name of a room in the game.
type RoomName [9]byte

func (name RoomName) String() string {
	return strings.Trim(string(name[:]), "\x00")
}
