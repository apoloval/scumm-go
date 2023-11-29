package scumm

import "fmt"

// BlockOffset is the offset of a block respect its parent in a resource file.
type BlockOffset uint32

// String returns the string representation of the block offset.
func (offset BlockOffset) String() string {
	return fmt.Sprintf("$%08x", uint32(offset))
}
