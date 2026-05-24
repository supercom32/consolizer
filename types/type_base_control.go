package types

import (
	"encoding/json"
	"sync"
)

/*
BaseControlType is a structure which represents the base properties for all UI controls. In addition, the following should be noted:

- It includes common attributes like position, size, visibility, and styling.

- It uses a mutex to ensure thread-safe access to its properties.

Example:

	baseControl := BaseControlType{}
*/
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
	Label            string
	IsBorderDrawn    bool
	TooltipAlias     string
	IsTooltipEnabled bool
}

/*
NewBaseControl is a constructor which creates a new base control with default values. In addition, the following should be noted:

- It initializes all fields with sensible defaults.

- It creates a new mutex for thread safety.

- It sets up a default style entry.

- All numeric fields are initialized to 0.

- Boolean fields are initialized to true where appropriate.

Example:

	baseControl := NewBaseControl()
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
	baseControl.Label = ""
	baseControl.IsBorderDrawn = false
	baseControl.TooltipAlias = ""
	baseControl.IsTooltipEnabled = false
	return baseControl
}

/*
GetBounds is a method which retrieves the position and size of a control. In addition, the following should be noted:

- It returns the X and Y location coordinates.

- It returns the width and height of the control.

- These values can be used for collision detection and layout calculations.

Example:

	x, y, w, h := baseControl.GetBounds()
*/
func (shared *BaseControlType) GetBounds() (int, int, int, int) {
	return shared.XLocation, shared.YLocation, shared.Width, shared.Height
}

/*
SetPosition is a method which sets the position of a control. In addition, the following should be noted:

- It updates the X and Y location coordinates.

- These values determine where the control is drawn on the screen.

- The position is relative to the parent layer's origin.

Example:

	baseControl.SetPosition(10, 20)
*/
func (shared *BaseControlType) SetPosition(x, y int) {
	shared.XLocation = x
	shared.YLocation = y
}

/*
SetSize is a method which sets the dimensions of a control. In addition, the following should be noted:

- It updates the width and height of the control.

- These values determine the control's visible area.

- The size affects how the control is drawn and how it responds to input.

Example:

	baseControl.SetSize(80, 24)
*/
func (shared *BaseControlType) SetSize(width, height int) {
	shared.Width = width
	shared.Height = height
}

/*
SetEnabled is a method which enables or disables a control. In addition, the following should be noted:

- When disabled, the control will not respond to user input.

- The visual appearance may change to indicate the disabled state.

- This state can be used to implement control dependencies.

Example:

	baseControl.SetEnabled(false)
*/
func (shared *BaseControlType) SetEnabled(enabled bool) {
	shared.IsEnabled = enabled
}

/*
SetVisible is a method which shows or hides a control. In addition, the following should be noted:

- When hidden, the control is not drawn and does not respond to input.

- The control's state is preserved while hidden.

- This can be used to implement dynamic interfaces.

Example:

	baseControl.SetVisible(false)
*/
func (shared *BaseControlType) SetVisible(visible bool) {
	shared.IsVisible = visible
}

/*
SetStyle is a method which changes the visual appearance of a control. In addition, the following should be noted:

- It updates the control's style entry with new visual properties.

- It affects colors, borders, and other visual attributes.

- The style can be changed dynamically at runtime.

Example:

	baseControl.SetStyle(style)
*/
func (shared *BaseControlType) SetStyle(style TuiStyleEntryType) {
	shared.StyleEntry = style
}

/*
GetAlias is a method which retrieves the alias of a control. In addition, the following should be noted:

- It returns the unique identifier for the control.

- This alias is used to reference the control in other operations.

- The alias is set when the control is created.

Example:

	alias := baseControl.GetAlias()
*/
func (shared *BaseControlType) GetAlias() string {
	return shared.Alias
}

/*
MarshalJSON is a method which serializes a control to JSON. In addition, the following should be noted:

- It converts the control's state to a JSON representation.

- It includes all base control properties.

- It is used for saving and loading control configurations.

Example:

	jsonData, err := baseControl.MarshalJSON()
*/
func (shared BaseControlType) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		StyleEntry       TuiStyleEntryType
		Alias            string
		XLocation        int
		YLocation        int
		Width            int
		Height           int
		IsEnabled        bool
		IsVisible        bool
		Label            string
		IsBorderDrawn    bool
		TooltipAlias     string
		IsTooltipEnabled bool
	}{
		StyleEntry:       shared.StyleEntry,
		Alias:            shared.Alias,
		XLocation:        shared.XLocation,
		YLocation:        shared.YLocation,
		Width:            shared.Width,
		Height:           shared.Height,
		IsEnabled:        shared.IsEnabled,
		IsVisible:        shared.IsVisible,
		Label:            shared.Label,
		IsBorderDrawn:    shared.IsBorderDrawn,
		TooltipAlias:     shared.TooltipAlias,
		IsTooltipEnabled: shared.IsTooltipEnabled,
	})
}

/*
GetEntryAsJsonDump is a method which returns a JSON string representation of a control. In addition, the following should be noted:

- It returns a formatted JSON string of the control's state.

- It is useful for debugging and logging purposes.

- It panics if JSON marshaling fails.

Example:

	jsonString := baseControl.GetEntryAsJsonDump()
*/
func (shared BaseControlType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

/*
IsEqual is a method which compares two controls for equality. In addition, the following should be noted:

- It compares all base control properties.

- It returns true if all properties match, false otherwise.

- It is used for change detection and state synchronization.

Example:

	isEqual := baseControl.IsEqual(otherControl)
*/
func (shared *BaseControlType) IsEqual(other *BaseControlType) bool {
	// Compare all fields except the mutex
	return shared.Alias == other.Alias &&
		shared.XLocation == other.XLocation &&
		shared.YLocation == other.YLocation &&
		shared.Width == other.Width &&
		shared.Height == other.Height &&
		shared.IsEnabled == other.IsEnabled &&
		shared.IsVisible == other.IsVisible &&
		shared.Label == other.Label &&
		shared.IsBorderDrawn == other.IsBorderDrawn &&
		shared.TooltipAlias == other.TooltipAlias &&
		shared.IsTooltipEnabled == other.IsTooltipEnabled
}
