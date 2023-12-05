package vm

import (
	"fmt"
	"reflect"
	"strings"
)

type Displayer interface {
	Display(st *SymbolTable) string
}

func DisplayInstruction(st *SymbolTable, inst Instruction) string {
	if disp, ok := inst.(Displayer); ok {
		return disp.Display(st)
	}
	return displayInstruction(st, inst)
}

func displayInstruction(st *SymbolTable, inst Instruction) string {
	elem := reflect.ValueOf(inst)
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}

	var str strings.Builder
	fmt.Fprintf(&str, "%s ", elem.Type().Name())
	displayOperands(st, elem, &str)
	return str.String()
}

func displayOperands(st *SymbolTable, elem reflect.Value, str *strings.Builder) {
	var ops []string
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		fieldType := elem.Type().Field(i)

		if field.Kind() == reflect.Struct && fieldType.Anonymous {
			displayOperands(st, field, str)
			continue
		}

		if !fieldType.IsExported() {
			continue
		}

		if tagOp := fieldType.Tag.Get("op"); tagOp == "" {
			continue
		}

		switch value := field.Interface().(type) {
		case Param:
			ops = append(ops, value.Display(st))
		case string:
			ops = append(ops, fmt.Sprintf("%q", value))
		default:
			ops = append(ops, fmt.Sprintf("%v", field.Interface()))
		}
	}
	str.WriteString(strings.Join(ops, ", "))
}
