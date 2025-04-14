package types

import (
	"encoding/json"
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
GetAlias allows you to retrieve the alias of a control. In addition, the following
information should be noted:

- Returns the unique identifier for the control.
- This alias is used to reference the control in other operations.
- The alias is set when the control is created.
*/
func (shared *BaseControlType) GetAlias() string {
	return shared.Alias
}

/*
MarshalJSON allows you to serialize a control to JSON. In addition, the following
information should be noted:

- Converts the control's state to a JSON representation.
- Includes all base control properties.
- Used for saving and loading control configurations.
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
GetEntryAsJsonDump allows you to get a JSON string representation of a control. In addition,
the following information should be noted:

- Returns a formatted JSON string of the control's state.
- Useful for debugging and logging purposes.
- Panics if JSON marshaling fails.
*/
func (shared BaseControlType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

/*
IsEqual allows you to compare two controls for equality. In addition, the following
information should be noted:

- Compares all base control properties.
- Returns true if all properties match, false otherwise.
- Used for change detection and state synchronization.
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
