package types

import (
	"encoding/json"
)

/*
ButtonEntryType is a structure which represents a button control. In addition, the following should be noted:

- It includes base control properties and the pressed and selected states of the button.

Example:
    buttonEntry := ButtonEntryType{}
*/
type ButtonEntryType struct {
	BaseControlType
	IsPressed  bool
	IsSelected bool
}

/*
MarshalJSON is a method which serializes a button control to JSON. In addition, the following should be noted:

- It converts the button's state to a JSON representation.

- It includes the base control properties and button-specific fields.

- It is used for saving and loading button configurations.

Example:
    jsonData, err := buttonEntry.MarshalJSON()
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
NewButtonEntry is a constructor which creates a new button control. In addition, the following should be noted:

- It initializes a button with default values.

- It can optionally copy properties from an existing button.

- It sets up the base control properties and button-specific fields.

Example:
    buttonEntry := NewButtonEntry(&existingButtonEntry)
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
IsButtonEntryEqual is a method which compares two button controls for equality. In addition, the following should be noted:

- It compares all properties of both buttons.

- It returns true if all properties match, false otherwise.

- It is used for change detection and state synchronization.

Example:
    isEqual := IsButtonEntryEqual(&sourceButton, &targetButton)
*/
func IsButtonEntryEqual(sourceButtonEntry *ButtonEntryType, targetButtonEntry *ButtonEntryType) bool {
	return sourceButtonEntry.BaseControlType.IsEqual(&targetButtonEntry.BaseControlType) &&
		sourceButtonEntry.IsPressed == targetButtonEntry.IsPressed &&
		sourceButtonEntry.IsSelected == targetButtonEntry.IsSelected
}
