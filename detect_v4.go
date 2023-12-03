package scumm

import (
	"encoding/binary"
	"io"

	"github.com/apoloval/scumm-go/ioutils"
)

const (
	// ResourceFileIndexV4 is a LFL index resource file for SCUMM v4.
	ResourceFileIndexV4 ResourceFileType = "SCUMM v4 LFL index file"

	// ResourceFileCharsetV4 is a charset resource file for SCUMM v4.
	ResourceFileCharsetV4 ResourceFileType = "SCUMM v4 LFL charset file"

	// ResourceFileBundleV4 is a data resource file for SCUMM v4.
	ResourceFileBundleV4 ResourceFileType = "SCUMM v4 LEC data file"
)

func isFileIndexv4(r io.ReadSeeker) bool {
	r.Seek(4, io.SeekStart)
	var blockName [2]byte
	if err := binary.Read(r, binary.LittleEndian, &blockName); err != nil {
		return false
	}
	return string(blockName[:]) == "RN" || string(blockName[:]) == "0R"
}

func isCharsetV4(r io.ReadSeeker) bool {
	r.Seek(4, io.SeekStart)
	var magic uint16
	if err := binary.Read(r, binary.LittleEndian, &magic); err != nil {
		return false
	}
	return magic == CharsetV4Magic
}

func isResourceBundleV4(r io.ReadSeeker) bool {
	r.Seek(4, io.SeekStart)
	xor := ioutils.NewXorReader(r, ResourceBundleV4Key)
	var blockType [2]byte

	if err := binary.Read(xor, binary.LittleEndian, &blockType); err != nil {
		return false
	}
	return blockType == ChunkTypeV4LE
}
