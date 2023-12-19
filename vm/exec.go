package vm

// Property is a VM property.
type Property string

const (
	PropUICursorCurrent  Property = "ui.cursor.current"
	PropUICursorVisible  Property = "ui.cursor.visible"
	PropUIUserputEnabled Property = "ui.userput.enabled"
)

type ExecutionContext interface {
	// GetProperty returns the value of a property.
	GetProperty(prop Property) int

	// SetProperty sets the value of a property.
	SetProperty(prop Property, value int)
}

func ExecContextFrom(e *Engine, t *Thread) ExecutionContext {
	return &executionContext{e, t}
}

type executionContext struct {
	engine *Engine
	thread *Thread
}

func (c *executionContext) GetProperty(prop Property) int {
	return c.engine.GetProperty(prop)
}

func (c *executionContext) SetProperty(prop Property, value int) {
	c.engine.SetProperty(prop, value)
}
