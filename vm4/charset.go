package vm4

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/apoloval/scumm-go/vm"
)

// CharsetMagic is the magic number of a charset resource of a SCUMM v4 game.
const CharsetMagic = 0x0363

// DecodeCharset decodes a charset resource of a SCUMM v4 game.
func DecodeCharset(r io.ReadSeeker) (vm.Charset, error) {
	r.Seek(0, io.SeekStart)
	var blockHeader struct {
		BlockSize       uint32
		Magic           uint16
		ColorMap        [15]byte
		BitsPerPixel    byte
		FontHeight      byte
		NumberOfChars   uint16
		CharDataOffsets [256]uint32
	}
	if err := binary.Read(r, binary.LittleEndian, &blockHeader); err != nil {
		return vm.Charset{}, err
	}

	if blockHeader.Magic != CharsetMagic {
		return vm.Charset{}, fmt.Errorf("invalid magic number: %04X", blockHeader.Magic)
	}

	// For some unknown reason, the block size is 11 bytes less than the actual size of the block.
	// I discovered this in the source code of ScummVM.
	blockHeader.BlockSize += 11

	charset := vm.Charset{
		ColorMap:     blockHeader.ColorMap,
		BitsPerPixel: blockHeader.BitsPerPixel,
		FontHeight:   blockHeader.FontHeight,
	}
	for i := 0; i < int(blockHeader.NumberOfChars); i++ {
		offset := blockHeader.CharDataOffsets[i]
		if offset == 0 {
			continue
		}

		// The offset is respect the end of the color map. Add the size of the color map, the magic
		// number and the block size to make it respect the start of the file.
		offset += 4 + 2 + 15
		if offset > blockHeader.BlockSize {
			return vm.Charset{}, fmt.Errorf("invalid offset %d for character %d", offset, i)
		}

		if _, err := r.Seek(int64(offset), io.SeekStart); err != nil {
			return vm.Charset{}, err
		}

		var charHeader struct {
			Width   uint8
			Height  uint8
			XOffset int8
			YOffset int8
		}
		if err := binary.Read(r, binary.LittleEndian, &charHeader); err != nil {
			return vm.Charset{}, err
		}

		glyphPixels := int(charHeader.Width) * int(charHeader.Height)
		glyphBits := glyphPixels * int(blockHeader.BitsPerPixel)
		glyphSize := glyphBits / 8
		if glyphBits%8 != 0 {
			glyphSize++
		}

		glyph := make([]byte, glyphSize)
		if _, err := io.ReadFull(r, glyph); err != nil {
			return vm.Charset{}, err
		}

		charset.Characters[i] = &vm.Character{
			Width:   charHeader.Width,
			Height:  charHeader.Height,
			XOffset: charHeader.XOffset,
			YOffset: charHeader.YOffset,
			Glyph:   glyph,
		}
	}

	return charset, nil
}
