package vm

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// BytecodeFrame is a bytecode frame.
type BytecodeFrame struct {
	StartAddress uint16
	Bytes        []byte
}

// NewBytecodeFrame creates a new bytecode frame.
func NewBytecodeFrame(addr uint16) BytecodeFrame {
	return BytecodeFrame{StartAddress: addr, Bytes: make([]byte, 0, 16)}
}

// Append appends bytes to the bytecode frame.
func (f *BytecodeFrame) Append(bytes []byte) {
	f.Bytes = append(f.Bytes, bytes...)
}

// String implements the Stringer interface.
func (f BytecodeFrame) Display(w io.Writer, line string) error {
	base := f.StartAddress
	data := f.Bytes
	for {
		rem := min(len(data), 8)

		if _, err := fmt.Fprintf(w, "%04X: %- 24X\t%s\n", base, data[:rem], line); err != nil {
			return err
		}
		if rem < 8 {
			return nil
		}
		line = ""
		base += 8
		data = data[8:]
	}
}

// BytecodeReader is a reader for reading bytecode elements.
type BytecodeReader struct {
	r     io.ReadSeeker
	frame BytecodeFrame
	err   error
}

// NewBytecodeReader creates a new bytecode reader.
func NewBytecodeReader(r io.ReadSeeker) *BytecodeReader {
	return &BytecodeReader{r: r}
}

// BeginFrame starts a new bytecode frame.
func (r *BytecodeReader) BeginFrame() {
	addr, _ := r.r.Seek(0, io.SeekCurrent)
	r.frame = NewBytecodeFrame(uint16(addr))
	r.err = nil
}

// EndFrame ends the current bytecode frame and returns its content.
func (r *BytecodeReader) EndFrame() (BytecodeFrame, error) {
	c := r.frame
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
func (r *BytecodeReader) ReadByteConstant(format NumberFormat) Constant {
	return Constant{
		Value:  int16(r.ReadByte()),
		Format: format,
	}
}

// ReadWordConstant reads a word constant.
func (r *BytecodeReader) ReadWordConstant(format NumberFormat) Constant {
	return Constant{
		Value:  int16(r.ReadWord()),
		Format: format,
	}
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

// ReadWordPointer reads a word pointer.
func (r *BytecodeReader) ReadWordPointer() WordPointer {
	if p, ok := r.ReadPointer().(WordPointer); ok {
		return p
	}
	r.err = errors.New("invalid pointer type")
	return 0
}

// ReadParam8 reads a parameter of 8 bits.
func (r *BytecodeReader) ReadByteParam(opcode OpCode, pos ParamPos, format NumberFormat) Param {
	if opcode.IsPointer(pos) {
		return r.ReadPointer()
	}
	return r.ReadByteConstant(format)
}

// ReadByteParams reads n parameters of 8 bits. n must be 1, 2 or 3.
func (r *BytecodeReader) ReadByteParams(opcode OpCode, n uint, format NumberFormat) []Param {
	if n < 1 || n > 3 {
		panic("invalid parameters")
	}
	pos := []ParamPos{ParamPos1, ParamPos2, ParamPos3}
	params := make([]Param, n)
	for i := uint(0); i < n; i++ {
		params[i] = r.ReadByteParam(opcode, pos[i], format)
	}
	return params
}

// ReadParam16 reads a parameter of 16 bits.
func (r *BytecodeReader) ReadWordParam(opcode OpCode, pos ParamPos, format NumberFormat) Param {
	if opcode.IsPointer(pos) {
		return r.ReadPointer()
	}
	return r.ReadWordConstant(format)
}

// ReadRelativeJump reads a program address.
func (r *BytecodeReader) ReadRelativeJump() Constant {
	rel := int16(r.ReadWord())
	pos := r.currentPos()
	return pos.Add(rel)
}

// ReadString reads a null-terminated string from the bytecode.
func (r *BytecodeReader) ReadString() string {
	var s string
	for {
		b := r.ReadByte()
		if b == 0 {
			return s
		}
		s += string(b)
	}
}

func (r *BytecodeReader) readBytes(b []byte) {
	if r.err == nil {
		_, r.err = r.r.Read(b[:])
		r.frame.Append(b[:])
	}
}

func (r *BytecodeReader) currentPos() Constant {
	addr, _ := r.r.Seek(0, io.SeekCurrent)
	return Constant{
		Value:  int16(addr),
		Format: ParamFormatProgramAddress,
	}
}
