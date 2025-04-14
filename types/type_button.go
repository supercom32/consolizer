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
	return sourceButtonEntry.BaseControlType.IsEqual(&targetButtonEntry.BaseControlType) &&
		sourceButtonEntry.IsPressed == targetButtonEntry.IsPressed &&
		sourceButtonEntry.IsSelected == targetButtonEntry.IsSelected
}
