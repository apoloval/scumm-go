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
