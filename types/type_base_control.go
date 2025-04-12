package types

import (
	"sync"
)

type BaseControlType struct {
	Mutex            *sync.Mutex
	StyleEntry       TuiStyleEntryType
	Alias            string
	XLocation        int
	YLocation        int
	Width            int
	Height           int
	IsEnabled        bool
	IsVisible        bool
	TabIndex         int
	Label            string
	IsBorderDrawn    bool
	TooltipAlias     string
	IsTooltipEnabled bool
}

/*
NewBaseControl allows you to create a new base control with default values. In addition, the following
information should be noted:

- Initializes all fields with sensible defaults.
- Creates a new mutex for thread safety.
- Sets up a default style entry.
- All numeric fields are initialized to 0.
- Boolean fields are initialized to true where appropriate.
*/
func NewBaseControl() BaseControlType {
	var baseControl BaseControlType
	baseControl.Mutex = &sync.Mutex{}
	baseControl.StyleEntry = NewTuiStyleEntry()
	baseControl.Alias = ""
	baseControl.XLocation = 0
	baseControl.YLocation = 0
	baseControl.Width = 0
	baseControl.Height = 0
	baseControl.IsEnabled = true
	baseControl.IsVisible = true
	baseControl.TabIndex = 0
	baseControl.Label = ""
	baseControl.IsBorderDrawn = false
	baseControl.TooltipAlias = ""
	baseControl.IsTooltipEnabled = false
	return baseControl
}

/*
GetBounds allows you to retrieve the position and size of a control. In addition, the following
information should be noted:

- Returns the X and Y location coordinates.
- Returns the width and height of the control.
- These values can be used for collision detection and layout calculations.
*/
func (shared *BaseControlType) GetBounds() (int, int, int, int) {
	return shared.XLocation, shared.YLocation, shared.Width, shared.Height
}

/*
SetPosition allows you to set the position of a control. In addition, the following
information should be noted:

- Updates the X and Y location coordinates.
- These values determine where the control is drawn on the screen.
- The position is relative to the parent layer's origin.
*/
func (shared *BaseControlType) SetPosition(x, y int) {
	shared.XLocation = x
	shared.YLocation = y
}

/*
SetSize allows you to set the dimensions of a control. In addition, the following
information should be noted:

- Updates the width and height of the control.
- These values determine the control's visible area.
- The size affects how the control is drawn and how it responds to input.
*/
func (shared *BaseControlType) SetSize(width, height int) {
	shared.Width = width
	shared.Height = height
}

/*
SetEnabled allows you to enable or disable a control. In addition, the following
information should be noted:

- When disabled, the control will not respond to user input.
- The visual appearance may change to indicate the disabled state.
- This state can be used to implement control dependencies.
*/
func (shared *BaseControlType) SetEnabled(enabled bool) {
	shared.IsEnabled = enabled
}

/*
SetVisible allows you to show or hide a control. In addition, the following
information should be noted:

- When hidden, the control is not drawn and does not respond to input.
- The control's state is preserved while hidden.
- This can be used to implement dynamic interfaces.
*/
func (shared *BaseControlType) SetVisible(visible bool) {
	shared.IsVisible = visible
}

/*
SetStyle allows you to change the visual appearance of a control. In addition, the following
information should be noted:

- Updates the control's style entry with new visual properties.
- Affects colors, borders, and other visual attributes.
- The style can be changed dynamically at runtime.
*/
func (shared *BaseControlType) SetStyle(style TuiStyleEntryType) {
	shared.StyleEntry = style
}

/*
SetTabIndex allows you to set the tab order of a control. In addition, the following
information should be noted:

- Determines the order in which controls receive focus when tabbing.
- Lower values receive focus before higher values.
- A value of 0 means the control will not receive focus via tabbing.
*/
func (shared *BaseControlType) SetTabIndex(index int) {
	shared.TabIndex = index
}

/*
Lock allows you to acquire the control's mutex. In addition, the following
information should be noted:

- Ensures thread-safe access to the control's properties.
- Blocks if the mutex is already locked by another goroutine.
- Should be paired with Unlock in a defer statement.
*/
func (shared *BaseControlType) Lock() {
	if shared.Mutex != nil {
		shared.Mutex.Lock()
	}
}

/*
Unlock allows you to release the control's mutex. In addition, the following
information should be noted:

- Releases the lock acquired by Lock.
- Allows other goroutines to access the control's properties.
- Should be called after Lock when the critical section is complete.
*/
func (shared *BaseControlType) Unlock() {
	if shared.Mutex != nil {
		shared.Mutex.Unlock()
	}
}
