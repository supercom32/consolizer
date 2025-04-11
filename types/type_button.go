package types

import (
	"encoding/json"
)

// func DrawButton(LayerAlias string, ButtonLabel string, StyleEntry TuiStyleEntryType, IsPressed bool, IsSelected bool, XLocation int, YLocation int, Width int, Height int) {
type ButtonEntryType struct {
	BaseControlType
	IsPressed  bool
	IsSelected bool
}

/*
GetAlias allows you to retrieve the alias of a button control. In addition, the following
information should be noted:

- Returns the unique identifier for the button.
- This alias is used to reference the button in other operations.
- The alias is set when the button is created.
*/
func (shared ButtonEntryType) GetAlias() string {
	return shared.Alias
}

/*
MarshalJSON allows you to serialize a button control to JSON. In addition, the following
information should be noted:

- Converts the button's state to a JSON representation.
- Includes the base control properties and button-specific fields.
- Used for saving and loading button configurations.
*/
func (shared ButtonEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		BaseControlType
		IsPressed  bool
		IsSelected bool
	}{
		BaseControlType: shared.BaseControlType,
		IsPressed:       shared.IsPressed,
		IsSelected:      shared.IsSelected,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

/*
GetEntryAsJsonDump allows you to get a JSON string representation of a button control. In addition,
the following information should be noted:

- Returns a formatted JSON string of the button's state.
- Useful for debugging and logging purposes.
- Panics if JSON marshaling fails.
*/
func (shared ButtonEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

/*
NewButtonEntry allows you to create a new button control. In addition, the following
information should be noted:

- Initializes a button with default values.
- Can optionally copy properties from an existing button.
- Sets up the base control properties and button-specific fields.
*/
func NewButtonEntry(existingButtonEntry ...*ButtonEntryType) ButtonEntryType {
	var buttonEntry ButtonEntryType
	buttonEntry.BaseControlType = NewBaseControl()

	if existingButtonEntry != nil {
		buttonEntry.BaseControlType = existingButtonEntry[0].BaseControlType
		buttonEntry.IsPressed = existingButtonEntry[0].IsPressed
		buttonEntry.IsSelected = existingButtonEntry[0].IsSelected
	}
	return buttonEntry
}

/*
IsButtonEntryEqual allows you to compare two button controls for equality. In addition, the following
information should be noted:

- Compares all properties of both buttons.
- Returns true if all properties match, false otherwise.
- Used for change detection and state synchronization.
*/
func IsButtonEntryEqual(sourceButtonEntry *ButtonEntryType, targetButtonEntry *ButtonEntryType) bool {
	if sourceButtonEntry.BaseControlType == targetButtonEntry.BaseControlType &&
		sourceButtonEntry.IsPressed == targetButtonEntry.IsPressed &&
		sourceButtonEntry.IsSelected == targetButtonEntry.IsSelected {
		return true
	}
	return false
}

/*
GetButtonAlias allows you to retrieve the alias of a button control. In addition, the following
information should be noted:

- Returns the unique identifier for the button.
- This is a convenience method that delegates to GetAlias.
- The alias is used to reference the button in other operations.
*/
func GetButtonAlias(entry *ButtonEntryType) string {
	return entry.Alias
}
