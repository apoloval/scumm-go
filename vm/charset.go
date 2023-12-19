package vm

import (
	"bytes"
	"fmt"
	"image"
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
