package vm

import "fmt"

type Engine struct {
	rm    ResourceManager
	props map[Property]int
}

func NewEngine(rm ResourceManager) *Engine {
	return &Engine{
		rm:    rm,
		props: make(map[Property]int),
	}
}

func (e *Engine) GetProperty(prop Property) int {
	return e.props[prop]
}

func (e *Engine) SetProperty(prop Property, value int) {
	e.props[prop] = value
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
