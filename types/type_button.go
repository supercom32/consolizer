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

func (shared ButtonEntryType) GetAlias() string {
	return shared.Alias
}

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

func (shared ButtonEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

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

func IsButtonEntryEqual(sourceButtonEntry *ButtonEntryType, targetButtonEntry *ButtonEntryType) bool {
	if sourceButtonEntry.BaseControlType == targetButtonEntry.BaseControlType &&
		sourceButtonEntry.IsPressed == targetButtonEntry.IsPressed &&
		sourceButtonEntry.IsSelected == targetButtonEntry.IsSelected {
		return true
	}
	return false
}

func GetButtonAlias(entry *ButtonEntryType) string {
	return entry.Alias
}
