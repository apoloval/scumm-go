package vm

import "fmt"

type Thread struct {
	script  *Script
	ip      int
	local   []int
	symbols *SymbolTable
}

func NewThread(script *Script) *Thread {
	return NewThreadOn(script, 0)
}

func NewThreadOn(script *Script, ip int) *Thread {
	return &Thread{
		script:  script,
		ip:      ip,
		local:   make([]int, MaxLocals),
		symbols: NewSymbolTable(),
	}
}

func (t *Thread) ReadLocal(idx uint16) int {
	return t.local[idx]
}

func (t *Thread) WriteLocal(idx uint16, value int) {
	t.local[idx] = value
}

func (t *Thread) Run(eng *Engine) error {
	ctx := ExecContextFrom(eng, t)
	for {
		inst := t.script.Code[t.ip]
		disp := DisplayInstruction(t.symbols, inst)
		if exec, ok := inst.(hasExecute); ok {
			exec.Execute(ctx)
			fmt.Printf("%04X: %s\n", t.ip, disp)
		} else {
			return fmt.Errorf("instruction does not implement execute: %s", disp)
		}
		t.ip++
	}
}

type hasExecute interface {
	Execute(ExecutionContext)
}
