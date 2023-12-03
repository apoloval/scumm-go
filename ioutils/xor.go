package ioutils

import "io"

// XorReader is a reader that XORs the bytes read with a given key.
type XorReader struct {
	r   io.ReadSeeker
	key byte
}

// NewXorReader returns a new XorReader that reads from r.
func NewXorReader(r io.ReadSeeker, key byte) *XorReader {
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

// Seek implements the io.Seeker interface.
func (r *XorReader) Seek(offset int64, whence int) (int64, error) {
	return r.r.Seek(offset, whence)
}
