package scumm

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"io"
	"strings"

	"github.com/apoloval/scumm-go/ioutils"
)

// Charset is a charset resource of a SCUMM game.
type Charset struct {
	ColorMap     CharsetColorMap
	BitsPerPixel byte
	FontHeight   byte
	Characters   [256]*Character
}

// CharsetV4Magic is the magic number of a charset resource of a SCUMM v4 game.
const CharsetV4Magic = 0x0363

// DecodeCharsetV4 decodes a charset resource of a SCUMM v4 game.
func DecodeCharsetV4(r io.ReadSeeker) (Charset, error) {
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
		return Charset{}, err
	}

	if blockHeader.Magic != CharsetV4Magic {
		return Charset{}, fmt.Errorf("invalid magic number: %04X", blockHeader.Magic)
	}

	// For some unknown reason, the block size is 11 bytes less than the actual size of the block.
	// I discovered this in the source code of ScummVM.
	blockHeader.BlockSize += 11

	charset := Charset{
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
			return Charset{}, fmt.Errorf("invalid offset %d for character %d", offset, i)
		}

		if _, err := r.Seek(int64(offset), io.SeekStart); err != nil {
			return Charset{}, err
		}

		var charHeader struct {
			Width   uint8
			Height  uint8
			XOffset int8
			YOffset int8
		}
		if err := binary.Read(r, binary.LittleEndian, &charHeader); err != nil {
			return Charset{}, err
		}

		glyphPixels := int(charHeader.Width) * int(charHeader.Height)
		glyphBits := glyphPixels * int(blockHeader.BitsPerPixel)
		glyphSize := glyphBits / 8
		if glyphBits%8 != 0 {
			glyphSize++
		}

		glyph := make([]byte, glyphSize)
		if _, err := io.ReadFull(r, glyph); err != nil {
			return Charset{}, err
		}

		charset.Characters[i] = &Character{
			Width:   charHeader.Width,
			Height:  charHeader.Height,
			XOffset: charHeader.XOffset,
			YOffset: charHeader.YOffset,
			Glyph:   glyph,
		}
	}

	return charset, nil
}

// CharWidth returns the width of a character of a charset.
func (c Charset) CharWidth(ch rune) int {
	char := c.Characters[ch]
	if char == nil {
		return 0
	}
	return int(char.Width) + int(char.XOffset)
}

// PrintChar prints a character of a charset into an image.
func (c Charset) PrintChar(ch rune, img *image.Paletted, loc image.Point) int {
	char := c.Characters[ch]
	if char == nil {
		return 0
	}
	loc.X += int(char.XOffset)
	loc.Y += int(char.YOffset)
	bits := ioutils.NewBitsReader(bytes.NewReader(char.Glyph))
	for y := 0; y < int(char.Height); y++ {
		for x := 0; x < int(char.Width); x++ {
			color, err := bits.ReadBits(int(c.BitsPerPixel))
			if err != nil {
				panic(err)
			}
			if color == 0 {
				continue
			}
			img.SetColorIndex(loc.X+x, loc.Y+y, c.ColorMap[color-1])
		}
	}
	return int(char.Width) + int(char.XOffset)
}

// CharsetColorMap is the color map of a charset.
type CharsetColorMap [15]byte

// String implements the fmt.Stringer interface.
func (c CharsetColorMap) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, color := range c {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%02X", color))
	}
	sb.WriteString("]")
	return sb.String()
}

// Character is a character of a charset.
type Character struct {
	Width   uint8
	Height  uint8
	XOffset int8
	YOffset int8
	Glyph   []byte
}
