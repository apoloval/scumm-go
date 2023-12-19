package vm

type Engine struct {
	props map[Property]int
}

func NewEngine() *Engine {
	return &Engine{
		props: make(map[Property]int),
	}
}

func (e *Engine) GetProperty(prop Property) int {
	return e.props[prop]
}

func (e *Engine) SetProperty(prop Property, value int) {
	e.props[prop] = value
}
