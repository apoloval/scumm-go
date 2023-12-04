package inst

import (
	"fmt"

	"github.com/apoloval/scumm-go/vm"
)

// LoadCharset is an instruction that loads a charset resource.
type LoadCharset struct {
	instruction
	CharsetID vm.Param
}

func (inst LoadCharset) Mnemonic(st *vm.SymbolTable) string {
	return fmt.Sprintf("LoadCharset %s",
		inst.CharsetID.Display(st, vm.ParamFormatCharsetID),
	)
}

func (inst *LoadCharset) Decode(opcode vm.OpCode, r *vm.BytecodeReader) (err error) {
	inst.CharsetID = r.ReadByteParam(opcode, vm.ParamPos1)
	inst.frame, err = r.EndFrame()
	return
}

// LoadSound is an instruction that loads a sound resource.
type LoadSound struct {
	instruction
	SoundID vm.Param
}

func (inst LoadSound) Mnemonic(st *vm.SymbolTable) string {
	return fmt.Sprintf("LoadSound %s",
		inst.SoundID.Display(st, vm.ParamFormatSoundID),
	)
}

func (inst *LoadSound) Decode(opcode vm.OpCode, r *vm.BytecodeReader) (err error) {
	inst.SoundID = r.ReadByteParam(opcode, vm.ParamPos1)
	inst.frame, err = r.EndFrame()
	return
}

// LockSound is an instruction that locks the sound.
type LockSound struct {
	instruction
	SoundID vm.Param
}

func (inst LockSound) Mnemonic(st *vm.SymbolTable) string {
	return fmt.Sprintf("LockSound %s",
		inst.SoundID.Display(st, vm.ParamFormatSoundID),
	)
}

func (inst *LockSound) Decode(opcode vm.OpCode, r *vm.BytecodeReader) (err error) {
	inst.SoundID = r.ReadByteParam(opcode, vm.ParamPos1)
	inst.frame, err = r.EndFrame()
	return
}

func decodeResourceRoutine(opcode vm.OpCode, r *vm.BytecodeReader) (inst vm.Instruction, err error) {
	sub := r.ReadOpCode()
	switch sub & 0x1F {
	case 0x02:
		inst = &LoadSound{}
	case 0x0A:
		inst = &LockSound{}
	case 0x12:
		inst = &LoadCharset{}
	default:
		return nil, fmt.Errorf("unknown opcode %02X %02X for resource routine", opcode, sub)
	}
	if err := inst.Decode(opcode, r); err != nil {
		return nil, err
	}
	return inst, nil
}
