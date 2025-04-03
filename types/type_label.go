package types

import (
	"encoding/json"
)

// func DrawButton(LayerAlias string, Value string, StyleEntry TuiStyleEntryType, IsPressed bool, IsSelected bool, XLocation int, YLocation int, Width int, Height int) {
type LabelEntryType struct {
	BaseControlType
	Label                   string
	Value                   string
	IsBackgroundTransparent bool
}

func (shared LabelEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		BaseControlType
		LabelValue string
	}{
		BaseControlType: shared.BaseControlType,
		LabelValue:      shared.Value,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (shared LabelEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

func NewLabelEntry(existingLabelEntry ...*LabelEntryType) LabelEntryType {
	var labelEntry LabelEntryType
	labelEntry.BaseControlType = NewBaseControl()

	if existingLabelEntry != nil {
		labelEntry.StyleEntry = NewTuiStyleEntry(&existingLabelEntry[0].StyleEntry)
		labelEntry.Alias = existingLabelEntry[0].Alias
		labelEntry.Value = existingLabelEntry[0].Value
		labelEntry.XLocation = existingLabelEntry[0].XLocation
		labelEntry.YLocation = existingLabelEntry[0].YLocation
		labelEntry.Width = existingLabelEntry[0].Width
		labelEntry.IsEnabled = existingLabelEntry[0].IsEnabled
		labelEntry.IsVisible = existingLabelEntry[0].IsVisible
		labelEntry.TabIndex = existingLabelEntry[0].TabIndex
	}
	return labelEntry
}

func IsLabelEntryEqual(sourceButtonEntry *LabelEntryType, targetButtonEntry *LabelEntryType) bool {
	if sourceButtonEntry.StyleEntry == targetButtonEntry.StyleEntry &&
		sourceButtonEntry.Alias == targetButtonEntry.Alias &&
		sourceButtonEntry.Value == targetButtonEntry.Value &&
		sourceButtonEntry.XLocation == targetButtonEntry.XLocation &&
		sourceButtonEntry.YLocation == targetButtonEntry.YLocation &&
		sourceButtonEntry.Width == targetButtonEntry.Width {
		return true
	}
	return false
}
