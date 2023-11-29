package scumm

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

// Charset is a charset resource of a SCUMM game.
type Charset struct {
	ColorMap     CharsetColorMap
	BitsPerPixel byte
	FontHeight   byte
	Characters   []*Character
}

// CharsetColorMap is the color map of a charset.
type CharsetColorMap [15]byte

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

// CharsetV4Magic is the magic number of a charset resource of a SCUMM v4 game.
const CharsetV4Magic = 0x0363

// DecodeCharsetV4 decodes a charset resource of a SCUMM v4 game.
func DecodeCharsetV4(r io.ReadSeeker) (Charset, error) {
	r.Seek(0, io.SeekStart)
	var header struct {
		BlockSize       uint32
		Magic           uint16
		ColorMap        [15]byte
		BitsPerPixel    byte
		FontHeight      byte
		NumberOfChars   uint16
		CharDataOffsets [256]uint32
	}
	if err := binary.Read(r, binary.LittleEndian, &header); err != nil {
		return Charset{}, err
	}

	if header.Magic != CharsetV4Magic {
		return Charset{}, fmt.Errorf("invalid magic number: %04X", header.Magic)
	}

	charsets := make([]*Character, header.NumberOfChars)
	for i, offset := range header.CharDataOffsets {
		if offset == 0 {
			fmt.Printf("WARN: char %d has no data\n", i)
			continue
		}

		// The offset is respect the end of the color map.
		if _, err := r.Seek(int64(offset+4+2+15), io.SeekStart); err != nil {
			return Charset{}, err
		}

		var charData struct {
			Width   uint8
			Height  uint8
			XOffset int8
			YOffset int8
		}
		if err := binary.Read(r, binary.LittleEndian, &charData); err != nil {
			return Charset{}, err
		}

		glypthSize := int(charData.Width) * int(charData.Height) * int(header.BitsPerPixel) / 8
		glyph := make([]byte, glypthSize)
		if _, err := io.ReadFull(r, glyph); err != nil {
			return Charset{}, err
		}

		charsets[i] = &Character{
			Width:   charData.Width,
			Height:  charData.Height,
			XOffset: charData.XOffset,
			YOffset: charData.YOffset,
			Glyph:   glyph,
		}
	}

	return Charset{
		ColorMap:     header.ColorMap,
		BitsPerPixel: header.BitsPerPixel,
		FontHeight:   header.FontHeight,
		Characters:   charsets,
	}, nil
}
