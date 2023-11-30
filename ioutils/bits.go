package ioutils

import (
	"encoding/binary"
	"fmt"
	"io"
)

// BitsReader reads arbitrary bits from an io.Reader.
type BitsReader struct {
	r io.Reader
	b byte
	n int
}

// NewBitsReader creates a new BitsReader.
func NewBitsReader(r io.Reader) *BitsReader {
	return &BitsReader{r: r}
}

// ReadBits reads w bits from the underlying io.Reader.
func (r *BitsReader) ReadBits(w int) (byte, error) {
	if w > 8 {
		return 0, fmt.Errorf("invalid width: %d", w)
	}
	if r.n == 0 {
		if err := binary.Read(r.r, binary.LittleEndian, &r.b); err != nil {
			return 0, err
		}
		r.n = 8
	}
	if w > r.n {
		return 0, io.ErrUnexpectedEOF
	}
	v := r.b >> (8 - w)
	r.b <<= w
	r.n -= w
	return v, nil
}
