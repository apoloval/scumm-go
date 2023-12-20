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

	// ReadWord reads the value of a word variable.
	ReadWord(idx uint16) int

	// WriteWord writes the value of a word variable.
	WriteWord(idx uint16, value int)

	// ReadBit reads the value of a bit variable.
	ReadBit(idx uint16) bool

	// WriteBit writes the value of a bit variable.
	WriteBit(idx uint16, value bool)

	// ReadLocal reads the value of a local variable.
	ReadLocal(idx uint16) int

	// WriteLocal writes the value of a local variable.
	WriteLocal(idx uint16, value int)
}

func ExecContextFrom(e *Engine, t *Thread) ExecutionContext {
	return &executionContext{e, t}
}

type executionContext struct {
	*Engine
	*Thread
}
