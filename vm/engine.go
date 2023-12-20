package vm

import "fmt"

const (
	MaxWords  = 8192
	MaxBits   = 32768
	MaxLocals = 16
)

type Engine struct {
	rm    ResourceManager
	props map[Property]int
	words []int
	bits  []byte
}

func NewEngine(rm ResourceManager) *Engine {
	return &Engine{
		rm:    rm,
		props: make(map[Property]int),
		words: make([]int, MaxWords),
		bits:  make([]byte, MaxBits/8),
	}
}

func (e *Engine) GetProperty(prop Property) int {
	return e.props[prop]
}

func (e *Engine) SetProperty(prop Property, value int) {
	e.props[prop] = value
}

func (e *Engine) ReadWord(idx uint16) int {
	return e.words[idx]
}

func (e *Engine) WriteWord(idx uint16, value int) {
	e.words[idx] = value
}

func (e *Engine) ReadBit(idx uint16) bool {
	return e.bits[idx/8]&(1<<(idx%8)) != 0
}

func (e *Engine) WriteBit(idx uint16, value bool) {
	if value {
		e.bits[idx/8] |= 1 << (idx % 8)
	} else {
		e.bits[idx/8] &^= 1 << (idx % 8)
	}
}

func (e *Engine) Run() error {
	bootscript, err := e.rm.GetScript(1, true)
	if err != nil {
		return fmt.Errorf("could not load bootscript: %v", err)
	}

	// TODO: determine what goes next after running bootscript
	th := NewThread(bootscript)
	return th.Run(e)
}
