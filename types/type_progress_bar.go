package types

import (
	"encoding/json"
	"sync"
)

// func DrawButton(LayerAlias string, Label string, StyleEntry TuiStyleEntryType, IsPressed bool, IsSelected bool, XLocation int, YLocation int, Width int, Height int) {
type ProgressBarEntryType struct {
	Mutex                   sync.Mutex
	StyleEntry              TuiStyleEntryType
	Alias                   string
	Label                   string
	Value                   int
	MaxValue                int
	IsBackgroundTransparent bool
	XLocation               int
	YLocation               int
	Width                   int
	Height                  int
	TabIndex                int
}

func (shared ProgressBarEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		StyleEntry              TuiStyleEntryType
		Alias                   string
		Label                   string
		Value                   int
		MaxValue                int
		IsBackgroundTransparent bool
		XLocation               int
		YLocation               int
		Width                   int
		Height                  int
		TabIndex                int
	}{
		StyleEntry:              shared.StyleEntry,
		Alias:                   shared.Alias,
		Label:                   shared.Label,
		Value:                   shared.Value,
		MaxValue:                shared.MaxValue,
		IsBackgroundTransparent: shared.IsBackgroundTransparent,
		XLocation:               shared.XLocation,
		YLocation:               shared.YLocation,
		Width:                   shared.Width,
		Height:                  shared.Height,
		TabIndex:                shared.TabIndex,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (shared ProgressBarEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

func NewProgressBarEntry(ExistingProgressBarEntry ...*ProgressBarEntryType) ProgressBarEntryType {
	var progressBarEntry ProgressBarEntryType
	if ExistingProgressBarEntry != nil {
		progressBarEntry.StyleEntry = NewTuiStyleEntry(&ExistingProgressBarEntry[0].StyleEntry)
		progressBarEntry.Alias = ExistingProgressBarEntry[0].Alias
		progressBarEntry.Label = ExistingProgressBarEntry[0].Label
		progressBarEntry.Value = ExistingProgressBarEntry[0].Value
		progressBarEntry.MaxValue = ExistingProgressBarEntry[0].MaxValue
		progressBarEntry.IsBackgroundTransparent = ExistingProgressBarEntry[0].IsBackgroundTransparent
		progressBarEntry.XLocation = ExistingProgressBarEntry[0].XLocation
		progressBarEntry.YLocation = ExistingProgressBarEntry[0].YLocation
		progressBarEntry.Width = ExistingProgressBarEntry[0].Width
		progressBarEntry.Height = ExistingProgressBarEntry[0].Height
		progressBarEntry.TabIndex = ExistingProgressBarEntry[0].TabIndex
	}
	return progressBarEntry
}

func IsProgressBarEntryEqual(sourceProgressBarEntry *ProgressBarEntryType, targetProgressBarEntry *ProgressBarEntryType) bool {
	if sourceProgressBarEntry.StyleEntry == targetProgressBarEntry.StyleEntry &&
		sourceProgressBarEntry.Alias == targetProgressBarEntry.Alias &&
		sourceProgressBarEntry.Label == targetProgressBarEntry.Label &&
		sourceProgressBarEntry.Value == targetProgressBarEntry.Value &&
		sourceProgressBarEntry.MaxValue == targetProgressBarEntry.MaxValue &&
		sourceProgressBarEntry.IsBackgroundTransparent == targetProgressBarEntry.IsBackgroundTransparent &&
		sourceProgressBarEntry.XLocation == targetProgressBarEntry.XLocation &&
		sourceProgressBarEntry.YLocation == targetProgressBarEntry.YLocation &&
		sourceProgressBarEntry.Width == targetProgressBarEntry.Width &&
		sourceProgressBarEntry.Height == targetProgressBarEntry.Height &&
		sourceProgressBarEntry.TabIndex == targetProgressBarEntry.TabIndex {
		return true
	}
	return false
}
