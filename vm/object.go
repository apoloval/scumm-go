package vm

import "fmt"

// ObjectID is the ID of an object.
type ObjectID int

// ObjectClass is the class of an object.
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

func (class ObjectClass) String() string {
	return fmt.Sprintf("$%06x", uint32(class))
}

type ObjectOwner byte

func (owner ObjectOwner) String() string {
	return fmt.Sprintf("$%02x", byte(owner))
}

type ObjectState byte

func (state ObjectState) String() string {
	return fmt.Sprintf("$%02x", byte(state))
}
