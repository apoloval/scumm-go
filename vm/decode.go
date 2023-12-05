package vm

import (
	"errors"
	"fmt"
	"reflect"
)

// InstructionDecoder is a function that decodes an instruction from a bytecode reader.
type InstructionDecoder func(r *BytecodeReader) (Instruction, error)

type OperandDecoder interface {
	DecodeOperands(opcode OpCode, r *BytecodeReader) error
}

func DecodeOperands(opcode OpCode, r *BytecodeReader, inst Instruction) error {
	if dec, ok := inst.(OperandDecoder); ok {
		return dec.DecodeOperands(opcode, r)
	}

	elem := reflect.ValueOf(inst).Elem()
	if !elem.CanAddr() {
		return errors.New("cannot decode operands of non-addressable instruction")
	}
	return decodeOperands(opcode, r, elem)
}

func decodeOperands(opcode OpCode, r *BytecodeReader, elem reflect.Value) error {
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		fieldType := elem.Type().Field(i)
		fieldName := fmt.Sprintf("%s.%s", elem.Type().Name(), fieldType.Name)

		if field.Kind() == reflect.Struct && fieldType.Anonymous {
			if err := decodeOperands(opcode, r, field); err != nil {
				return err
			}
			continue
		}

		if !fieldType.IsExported() {
			continue
		}

		tagPos := fieldType.Tag.Get("pos")
		tagOp := fieldType.Tag.Get("op")
		tagFmt := fieldType.Tag.Get("fmt")

		if tagOp == "" {
			continue
		}

		var value any
		switch tagOp {
		case "result", "var":
			value = r.ReadPointer()
		case "byte", "8", "c":
			value = r.ReadByteConstant(NumberFormat(tagFmt))
		case "word", "16":
			value = r.ReadWordConstant(NumberFormat(tagFmt))
		case "param8", "p8":
			pos, ok := ParseParamPos(tagPos)
			if !ok {
				return fmt.Errorf("invalid param position in %s: %s", fieldName, tagPos)
			}
			value = r.ReadByteParam(opcode, pos, NumberFormat(tagFmt))
		case "param16", "p16":
			pos, ok := ParseParamPos(tagPos)
			if !ok {
				return fmt.Errorf("invalid param position in %s: %s", fieldName, tagPos)
			}
			value = r.ReadWordParam(opcode, pos, NumberFormat(tagFmt))
		case "str":
			value = r.ReadString()
		case "varargs", "v16":
			value = r.ReadVarParams()
		case "reljmp", "jmp":
			value = r.ReadRelativeJump()
		default:
			return fmt.Errorf("invalid operand type in %s: %s", fieldName, tagOp)
		}

		if !reflect.TypeOf(value).AssignableTo(fieldType.Type) {
			return fmt.Errorf("cannot assign %T to %s", value, fieldName)
		}
		field.Set(reflect.ValueOf(value))
	}
	return nil
}
