package inst

import (
	"fmt"
	"io"

	"github.com/apoloval/scumm-go/vm"
)

// DecodeAll decodes all instructions from the bytecode reader.
func DecodeAll(r *vm.BytecodeReader) (code []vm.Instruction, err error) {
	for {
		var i vm.Instruction
		i, err = Decode(r)
		if err == io.EOF {
			return code, nil
		}
		if err != nil {
			return code, err
		}
		code = append(code, i)
	}
}

// Decode decodes an instruction from the bytecode reader.
func Decode(r *vm.BytecodeReader) (vm.Instruction, error) {
	r.BeginFrame()
	opcode := r.ReadOpCode()

	var inst vm.Instruction
	switch opcode {
	case 0x00:
		inst = &StopObjectCode{base: withName("StopObjectCode")}
	case 0x04, 0x84:
		inst = &IsGreaterEqual{binaryBranch: withBinaryBranch(">=")}
	case 0x08, 0x88:
		inst = &IsNotEqual{binaryBranch: withBinaryBranch("!=")}
	case 0x0A, 0x2A, 0x4A, 0x6A, 0x8A, 0xAA, 0xCA, 0xEA:
		inst = &StartScript{base: withName("StartScript")}
	case 0x0C:
		return decodeResourceRoutine(opcode, r)
	case 0x18:
		inst = &Goto{branch: withBranch("Goto")}
	case 0x1A, 0x9A:
		inst = &Move{base: withName("Move")}
	case 0x26, 0xA6:
		inst = &SetVarRange{base: withName("SetVarRange")}
	case 0x27:
		return decodeStringOp(opcode, r)
	case 0x28:
		inst = &IsEqualZero{unaryBranch: withUnaryBranch("== 0")}
	case 0x2C:
		return decodeCursorCommand(opcode, r)
	case 0x38, 0xB8:
		inst = &IsLessEqual{binaryBranch: withBinaryBranch("<=")}
	case 0x44, 0xC4:
		inst = &IsLess{binaryBranch: withBinaryBranch("<")}
	case 0x48, 0xC8:
		inst = &IsEqual{binaryBranch: withBinaryBranch("==")}
	case 0x78, 0xF8:
		inst = &IsGreater{binaryBranch: withBinaryBranch(">")}
	case 0x5C:
		inst = &RoomFade{base: withName("RoomFade")}
	case 0xCC:
		inst = &PseudoRoom{base: withName("PseudoRoom")}
	case 0xA8:
		inst = &IsNotEqualZero{unaryBranch: withUnaryBranch("!= 0")}
	case 0xAC:
		inst = &Expression{}
	default:
		return nil, fmt.Errorf("unknown opcode %02X", opcode)
	}
	if err := inst.Decode(opcode, r); err != nil {
		return nil, err
	}
	return inst, nil
}
