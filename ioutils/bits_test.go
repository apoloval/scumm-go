package ioutils_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/apoloval/scumm-go/ioutils"
	"github.com/stretchr/testify/assert"
)

func TestBitsReader_1bit(t *testing.T) {
	data := []byte{0b00011011, 0b11100100}
	r := ioutils.NewBitsReader(bytes.NewReader(data))
	v0, _ := r.ReadBits(1)
	v1, _ := r.ReadBits(1)
	v2, _ := r.ReadBits(1)
	v3, _ := r.ReadBits(1)
	v4, _ := r.ReadBits(1)
	v5, _ := r.ReadBits(1)
	v6, _ := r.ReadBits(1)
	v7, _ := r.ReadBits(1)
	v8, _ := r.ReadBits(1)
	v9, _ := r.ReadBits(1)
	v10, _ := r.ReadBits(1)
	v11, _ := r.ReadBits(1)
	v12, _ := r.ReadBits(1)
	v13, _ := r.ReadBits(1)
	v14, _ := r.ReadBits(1)
	v15, _ := r.ReadBits(1)

	assert.Equal(t, byte(0b0), v0)
	assert.Equal(t, byte(0b0), v1)
	assert.Equal(t, byte(0b0), v2)
	assert.Equal(t, byte(0b1), v3)
	assert.Equal(t, byte(0b1), v4)
	assert.Equal(t, byte(0b0), v5)
	assert.Equal(t, byte(0b1), v6)
	assert.Equal(t, byte(0b1), v7)
	assert.Equal(t, byte(0b1), v8)
	assert.Equal(t, byte(0b1), v9)
	assert.Equal(t, byte(0b1), v10)
	assert.Equal(t, byte(0b0), v11)
	assert.Equal(t, byte(0b0), v12)
	assert.Equal(t, byte(0b1), v13)
	assert.Equal(t, byte(0b0), v14)
	assert.Equal(t, byte(0b0), v15)

	_, err := r.ReadBits(1)
	assert.Equal(t, io.EOF, err)
}

func TestBitsReader_2bits(t *testing.T) {
	data := []byte{0b00_01_10_11, 0b11_10_01_00}
	r := ioutils.NewBitsReader(bytes.NewReader(data))
	v0, _ := r.ReadBits(2)
	v1, _ := r.ReadBits(2)
	v2, _ := r.ReadBits(2)
	v3, _ := r.ReadBits(2)
	v4, _ := r.ReadBits(2)
	v5, _ := r.ReadBits(2)
	v6, _ := r.ReadBits(2)
	v7, _ := r.ReadBits(2)

	assert.Equal(t, byte(0b00), v0)
	assert.Equal(t, byte(0b01), v1)
	assert.Equal(t, byte(0b10), v2)
	assert.Equal(t, byte(0b11), v3)
	assert.Equal(t, byte(0b11), v4)
	assert.Equal(t, byte(0b10), v5)
	assert.Equal(t, byte(0b01), v6)
	assert.Equal(t, byte(0b00), v7)

	_, err := r.ReadBits(2)
	assert.Equal(t, io.EOF, err)
}

func TestBitsReader_4bits(t *testing.T) {
	data := []byte{0b0001_1011, 0b1110_0100}
	r := ioutils.NewBitsReader(bytes.NewReader(data))
	v0, _ := r.ReadBits(4)
	v1, _ := r.ReadBits(4)
	v2, _ := r.ReadBits(4)
	v3, _ := r.ReadBits(4)

	assert.Equal(t, byte(0b0001), v0)
	assert.Equal(t, byte(0b01011), v1)
	assert.Equal(t, byte(0b1110), v2)
	assert.Equal(t, byte(0b0100), v3)

	_, err := r.ReadBits(4)
	assert.Equal(t, io.EOF, err)
}
