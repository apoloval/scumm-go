package ioutils_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/apoloval/scumm-go/ioutils"
	"github.com/stretchr/testify/assert"
)

func TestXorReader(t *testing.T) {
	input := []byte{0x01, 0x02, 0x03, 0x04}
	r := ioutils.NewXorReader(bytes.NewReader(input), 0x69)
	output, err := io.ReadAll(r)

	assert.NoError(t, err)
	assert.Equal(t, []byte{0x68, 0x6b, 0x6a, 0x6d}, output)
}
