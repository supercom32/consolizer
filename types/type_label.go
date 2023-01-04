package types

import (
	"encoding/json"
	"sync"
)

// func DrawButton(LayerAlias string, Value string, StyleEntry TuiStyleEntryType, IsPressed bool, IsSelected bool, XLocation int, YLocation int, Width int, Height int) {
type LabelEntryType struct {
	Mutex      sync.Mutex
	StyleEntry TuiStyleEntryType
	Alias      string
	Value      string
	XLocation  int
	YLocation  int
	Width      int
}

func (shared LabelEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		StyleEntry TuiStyleEntryType
		LabelAlias string
		LabelValue string
		XLocation  int
		YLocation  int
		Width      int
	}{
		StyleEntry: shared.StyleEntry,
		LabelAlias: shared.Alias,
		LabelValue: shared.Value,
		XLocation:  shared.XLocation,
		YLocation:  shared.YLocation,
		Width:      shared.Width,
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

func NewLabelEntry(existingButtonEntry ...*LabelEntryType) LabelEntryType {
	var labelEntry LabelEntryType
	if existingButtonEntry != nil {
		labelEntry.StyleEntry = NewTuiStyleEntry(&existingButtonEntry[0].StyleEntry)
		labelEntry.Alias = existingButtonEntry[0].Alias
		labelEntry.Value = existingButtonEntry[0].Value
		labelEntry.XLocation = existingButtonEntry[0].XLocation
		labelEntry.YLocation = existingButtonEntry[0].YLocation
		labelEntry.Width = existingButtonEntry[0].Width
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
