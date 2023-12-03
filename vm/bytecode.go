package vm

import (
	"encoding/binary"
	"errors"
	"io"
)

// BytecodeReader is a reader for reading bytecode elements.
type BytecodeReader struct {
	r     io.ReadSeeker
	chunk []byte
	err   error
}

// NewBytecodeReader creates a new bytecode reader.
func NewBytecodeReader(r io.ReadSeeker) *BytecodeReader {
	return &BytecodeReader{r: r}
}

// BeginFrame starts a new bytecode frame.
func (r *BytecodeReader) BeginFrame() {
	r.chunk = make([]byte, 0, 16)
	r.err = nil
}

// EndFrame ends the current bytecode frame and returns its content.
func (r *BytecodeReader) EndFrame() ([]byte, error) {
	c := r.chunk
	e := r.err
	r.BeginFrame()
	return c, e
}

// ReadByte reads a byte.
func (r *BytecodeReader) ReadByte() byte {
	var b [1]byte
	r.readBytes(b[:])
	return b[0]
}

// ReadWord reads a word.
func (r *BytecodeReader) ReadWord() uint16 {
	var b [2]byte
	r.readBytes(b[:])
	return binary.LittleEndian.Uint16(b[:])
}

// ReadOpCode reads an opcode.
func (r *BytecodeReader) ReadOpCode() OpCode {
	return OpCode(r.ReadByte())
}

// ReadByteConstant reads a byte constant.
func (r *BytecodeReader) ReadByteConstant() ByteConstant {
	return ByteConstant(r.ReadByte())
}

// ReadWordConstant reads a word constant.
func (r *BytecodeReader) ReadWordConstant() WordConstant {
	return WordConstant(r.ReadWord())
}

// ReadPointer reads a word address.
func (r *BytecodeReader) ReadPointer() Pointer {
	word := r.ReadWord()
	if word&0xE000 == 0 {
		return WordPointer(word & 0x1FFF)
	}
	if word&0x8000 > 0 {
		return BitPointer(word & 0x7FFF)
	}
	if word&0xFFF0 == 0x4000 || word&0xFFF0 == 0x6000 {
		return LocalPointer(word & 0x000F)
	}
	r.err = errors.New("invalid pointer")
	return WordPointer(0)
}

// ReadParam8 reads a parameter of 8 bits.
func (r *BytecodeReader) ReadByteParam(opcode OpCode, pos ParamPos) Param {
	if opcode.IsPointer(pos) {
		return r.ReadPointer()
	}
	return r.ReadByteConstant()
}

// ReadByteParams reads n parameters of 8 bits. n must be 1, 2 or 3.
func (r *BytecodeReader) ReadByteParams(opcode OpCode, n uint) []Param {
	if n < 1 || n > 3 {
		panic("invalid parameters")
	}
	pos := []ParamPos{ParamPos1, ParamPos2, ParamPos3}
	params := make([]Param, n)
	for i := uint(0); i < n; i++ {
		params[i] = r.ReadByteParam(opcode, pos[i])
	}
	return params
}

// ReadParam16 reads a parameter of 16 bits.
func (r *BytecodeReader) ReadWordParam(opcode OpCode, pos ParamPos) Param {
	if opcode.IsPointer(pos) {
		return r.ReadPointer()
	}
	return r.ReadWordConstant()
}

func (r *BytecodeReader) readBytes(b []byte) {
	if r.err == nil {
		_, r.err = r.r.Read(b[:])
		r.chunk = append(r.chunk, b[:]...)
	}
}
