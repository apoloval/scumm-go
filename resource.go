package scumm

import (
	"encoding/binary"
	"fmt"
	"io"
)

// BlockOffset is the offset of a block respect its parent in a resource file.
type BlockOffset uint32

// String returns the string representation of the block offset.
func (offset BlockOffset) String() string {
	return fmt.Sprintf("$%08x", uint32(offset))
}

// ResourceFileType is the type of a resource file.
type ResourceFileType string

const (
	// ResourceFileIndexV4 is a LFL index resource file for SCUMM v4.
	ResourceFileIndexV4 ResourceFileType = "SCUMM v4 LFL index file"

	// ResourceFileUknown is an unknown resource file.
	ResourceFileUknown ResourceFileType = "unknown resource file"
)

// DetectResourceFile detects the type of a resource file.
func DetectResourceFile(r io.ReadSeeker) ResourceFileType {
	defer r.Seek(0, io.SeekStart)

	if isFileIndexv4(r) {
		return ResourceFileIndexV4
	}
	return ResourceFileUknown
}

func isFileIndexv4(r io.ReadSeeker) bool {
	r.Seek(4, io.SeekStart)
	var blockName [2]byte
	if err := binary.Read(r, binary.LittleEndian, &blockName); err != nil {
		return false
	}
	return string(blockName[:]) == "RN" || string(blockName[:]) == "0R"
}
