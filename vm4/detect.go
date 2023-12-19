package vm4

import (
	"encoding/binary"
	"io"

	"github.com/apoloval/scumm-go/ioutils"
	"github.com/apoloval/scumm-go/vm"
)

const (
	// ResourceFileIndex is a LFL index resource file for SCUMM v4.
	ResourceFileIndex vm.ResourceFileType = "SCUMM v4 LFL index file"

	// ResourceFileCharset is a charset resource file for SCUMM v4.
	ResourceFileCharset vm.ResourceFileType = "SCUMM v4 LFL charset file"

	// ResourceFileBundle is a data resource file for SCUMM v4.
	ResourceFileBundle vm.ResourceFileType = "SCUMM v4 LEC data file"
)

func IsFileIndex(r io.ReadSeeker) bool {
	r.Seek(4, io.SeekStart)
	var blockName [2]byte
	if err := binary.Read(r, binary.LittleEndian, &blockName); err != nil {
		return false
	}
	return string(blockName[:]) == "RN" || string(blockName[:]) == "0R"
}

func IsCharset(r io.ReadSeeker) bool {
	r.Seek(4, io.SeekStart)
	var magic uint16
	if err := binary.Read(r, binary.LittleEndian, &magic); err != nil {
		return false
	}
	return magic == CharsetMagic
}

func IsResourceBundle(r io.ReadSeeker) bool {
	r.Seek(4, io.SeekStart)
	xor := ioutils.NewXorReader(r, ResourceBundleKey)
	var blockType [2]byte

	if err := binary.Read(xor, binary.LittleEndian, &blockType); err != nil {
		return false
	}
	return blockType == ChunkTypeLE
}
