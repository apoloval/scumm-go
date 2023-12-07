package vm

import (
	"encoding/binary"
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

// BytecodeDecoder is a reader for reading bytecode elements.
type BytecodeDecoder struct {
	r     io.ReadSeeker
	frame BytecodeFrame
	err   error
}

// NewBytecodeDecoder creates a new bytecode reader.
func NewBytecodeDecoder(r io.ReadSeeker) *BytecodeDecoder {
	return &BytecodeDecoder{r: r}
}

// BeginFrame starts a new bytecode frame.
func (d *BytecodeDecoder) BeginFrame() {
	addr, _ := d.r.Seek(0, io.SeekCurrent)
	d.frame = NewBytecodeFrame(uint16(addr))
	d.err = nil
}

// EndFrame ends the current bytecode frame and returns its content.
func (d *BytecodeDecoder) EndFrame() (BytecodeFrame, error) {
	c := d.frame
	e := d.err
	d.BeginFrame()
	return c, e
}

// DecodeByte decodes a byte.
func (d *BytecodeDecoder) DecodeByte() byte {
	var b [1]byte
	d.readBytes(b[:])
	return b[0]
}

// DecodeWord decodes a word.
func (d *BytecodeDecoder) DecodeWord() uint16 {
	var b [2]byte
	d.readBytes(b[:])
	return binary.LittleEndian.Uint16(b[:])
}

// DecodeOpCode decodes an opcode.
func (d *BytecodeDecoder) DecodeOpCode() OpCode {
	return OpCode(d.DecodeByte())
}

// DecodeByteConstant decodes a byte constant.
func (d *BytecodeDecoder) DecodeByteConstant(format NumberFormat) Constant {
	return Constant{
		Value:  int16(d.DecodeByte()),
		Format: format,
	}
}

// DecodeWordConstant decodes a word constant.
func (d *BytecodeDecoder) DecodeWordConstant(format NumberFormat) Constant {
	return Constant{
		Value:  int16(d.DecodeWord()),
		Format: format,
	}
}

// DecodeVarRef decodes a variable reference.
func (d *BytecodeDecoder) DecodeVarRef() (ref VarRef) {
	ref.VarID = d.DecodeWord()
	if ref.VarID&0x2000 != 0 {
		ref.Offset = d.DecodeWord()
	}
	return
}

// ReadParam8 decodes a parameter of 8 bits.
func (d *BytecodeDecoder) DecodeByteParam(opcode OpCode, pos ParamPos, format NumberFormat) Param {
	if opcode.IsPointer(pos) {
		return d.DecodeVarRef()
	}
	return d.DecodeByteConstant(format)
}

// ReadParam16 decodes a parameter of 16 bits.
func (d *BytecodeDecoder) DecodeWordParam(opcode OpCode, pos ParamPos, format NumberFormat) Param {
	if opcode.IsPointer(pos) {
		return d.DecodeVarRef()
	}
	return d.DecodeWordConstant(format)
}

// DecodeVarParams decodes a variable number of parameters.
func (d *BytecodeDecoder) DecodeVarParams() (params Params) {
	for {
		b := d.DecodeOpCode()
		if b == 0xFF {
			return
		}
		params = append(params, d.DecodeWordParam(b, ParamPos1, NumberFormatDecimal))
	}
}

// DecodeNullTerminatedBytes decodes a sequence of bytes terminated by a null byte.
func (d *BytecodeDecoder) DecodeNullTerminatedBytes(fmt NumberFormat) (bytes []Constant) {
	for {
		b := d.DecodeByte()
		if b == 0 {
			return
		}
		bytes = append(bytes, Constant{
			Value:  int16(b),
			Format: fmt,
		})
	}
}

// DecodeRelativeJump decodes a program address.
func (d *BytecodeDecoder) DecodeRelativeJump() Constant {
	rel := int16(d.DecodeWord())
	pos := d.currentPos()
	return pos.Add(rel)
}

// DecodeString decodes a null-terminated string from the bytecode.
func (d *BytecodeDecoder) DecodeString() string {
	var s string
	for {
		b := d.DecodeByte()
		if b == 0 {
			return s
		}
		s += string(b)
	}
}

func (d *BytecodeDecoder) readBytes(b []byte) {
	if d.err == nil {
		_, d.err = d.r.Read(b[:])
		d.frame.Append(b[:])
	}
}

func (d *BytecodeDecoder) currentPos() Constant {
	addr, _ := d.r.Seek(0, io.SeekCurrent)
	return Constant{
		Value:  int16(addr),
		Format: NumberFormatAddress,
	}
}
