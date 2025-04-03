package types

import (
	"encoding/json"
)

// func DrawButton(LayerAlias string, Label string, StyleEntry TuiStyleEntryType, IsPressed bool, IsSelected bool, XLocation int, YLocation int, Width int, Height int) {
type ProgressBarEntryType struct {
	BaseControlType
	Label                   string
	Value                   int
	MaxValue                int
	IsBackgroundTransparent bool
	Length                  int
	CurrentValue            int
	IsHorizontal            bool
}

func (shared ProgressBarEntryType) GetAlias() string {
	return shared.Alias
}

func (shared ProgressBarEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		BaseControlType
		Label                   string
		Value                   int
		MaxValue                int
		IsBackgroundTransparent bool
		Length                  int
		CurrentValue            int
		IsHorizontal            bool
	}{
		BaseControlType:         shared.BaseControlType,
		Label:                   shared.Label,
		Value:                   shared.Value,
		MaxValue:                shared.MaxValue,
		IsBackgroundTransparent: shared.IsBackgroundTransparent,
		Length:                  shared.Length,
		CurrentValue:            shared.CurrentValue,
		IsHorizontal:            shared.IsHorizontal,
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
		progressBarEntry.BaseControlType = ExistingProgressBarEntry[0].BaseControlType
		progressBarEntry.Label = ExistingProgressBarEntry[0].Label
		progressBarEntry.Value = ExistingProgressBarEntry[0].Value
		progressBarEntry.MaxValue = ExistingProgressBarEntry[0].MaxValue
		progressBarEntry.IsBackgroundTransparent = ExistingProgressBarEntry[0].IsBackgroundTransparent
		progressBarEntry.Length = ExistingProgressBarEntry[0].Length
		progressBarEntry.CurrentValue = ExistingProgressBarEntry[0].CurrentValue
		progressBarEntry.IsHorizontal = ExistingProgressBarEntry[0].IsHorizontal
	}
	return progressBarEntry
}

func IsProgressBarEntryEqual(sourceProgressBarEntry *ProgressBarEntryType, targetProgressBarEntry *ProgressBarEntryType) bool {
	if sourceProgressBarEntry.BaseControlType == targetProgressBarEntry.BaseControlType &&
		sourceProgressBarEntry.Label == targetProgressBarEntry.Label &&
		sourceProgressBarEntry.Value == targetProgressBarEntry.Value &&
		sourceProgressBarEntry.MaxValue == targetProgressBarEntry.MaxValue &&
		sourceProgressBarEntry.IsBackgroundTransparent == targetProgressBarEntry.IsBackgroundTransparent &&
		sourceProgressBarEntry.Length == targetProgressBarEntry.Length &&
		sourceProgressBarEntry.CurrentValue == targetProgressBarEntry.CurrentValue &&
		sourceProgressBarEntry.IsHorizontal == targetProgressBarEntry.IsHorizontal {
		return true
	}
	return false
}

func GetProgressBarAlias(entry *ProgressBarEntryType) string {
	return entry.Alias
}
