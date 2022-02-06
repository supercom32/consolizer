package memory

import (
	"encoding/json"
)

type SelectorEntryType struct {
	StyleEntry TuiStyleEntryType
	SelectionEntry SelectionEntryType
	XLocation int
	YLocation  int
	ItemWidth        int
	NumberOfColumns  int
	ViewportPosition int
	ItemHighlighted  int
	ItemSelected int
}

func (shared SelectorEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		StyleEntry TuiStyleEntryType
		SelectionEntry SelectionEntryType
		XLocation int
		YLocation int
		ItemWidth int
		NumberOfColumns int
		ViewportPosition int
		ItemHighlighted int
		ItemSelected int
	}{
		StyleEntry: shared.StyleEntry,
		SelectionEntry: shared.SelectionEntry,
		XLocation: shared.XLocation,
		YLocation: shared.YLocation,
		ItemWidth: shared.ItemWidth,
		NumberOfColumns: shared.NumberOfColumns,
		ViewportPosition: shared.ViewportPosition,
		ItemHighlighted: shared.ItemHighlighted,
		ItemSelected: shared.ItemSelected,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (shared SelectorEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

func NewSelectorEntry(existingSelectorEntry ...*SelectorEntryType) SelectorEntryType {
	var menuBarEntry SelectorEntryType
	if existingSelectorEntry != nil {
		menuBarEntry.StyleEntry = existingSelectorEntry[0].StyleEntry
		menuBarEntry.SelectionEntry = existingSelectorEntry[0].SelectionEntry
		menuBarEntry.XLocation = existingSelectorEntry[0].XLocation
		menuBarEntry.YLocation = existingSelectorEntry[0].YLocation
		menuBarEntry.ItemWidth = existingSelectorEntry[0].ItemWidth
		menuBarEntry.NumberOfColumns = existingSelectorEntry[0].NumberOfColumns
		menuBarEntry.ViewportPosition = existingSelectorEntry[0].ViewportPosition
		menuBarEntry.ItemHighlighted = existingSelectorEntry[0].ItemHighlighted
		menuBarEntry.ItemSelected = existingSelectorEntry[0].ItemSelected
	}
	return menuBarEntry
}