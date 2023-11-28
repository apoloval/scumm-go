package ioutils

import "io"

// XorReader is a reader that XORs the bytes read with a given key.
type XorReader struct {
	r   io.Reader
	key byte
}

// NewXorReader returns a new XorReader that reads from r.
func NewXorReader(r io.Reader, key byte) *XorReader {
	return &XorReader{r: r, key: key}
}

// Read implements the io.Reader interface.
func (r *XorReader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	for i := 0; i < n; i++ {
		p[i] ^= r.key
	}
	return
}
