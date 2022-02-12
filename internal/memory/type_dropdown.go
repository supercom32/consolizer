package memory

import (
	"encoding/json"
)

type DropdownEntryType struct {
	StyleEntry TuiStyleEntryType
	SelectionEntry SelectionEntryType
	XLocation int
	YLocation  int
	SelectorHeight int
	ItemWidth        int
	ViewportPosition int
	ItemHighlighted  int
	ItemSelected int
	IsTrayOpen bool
}

func (shared DropdownEntryType) MarshalJSON() ([]byte, error) {
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
		IsTrayOpen bool
	}{
		StyleEntry: shared.StyleEntry,
		SelectionEntry: shared.SelectionEntry,
		XLocation: shared.XLocation,
		YLocation: shared.YLocation,
		ItemWidth: shared.ItemWidth,
		ViewportPosition: shared.ViewportPosition,
		ItemHighlighted: shared.ItemHighlighted,
		ItemSelected: shared.ItemSelected,
		IsTrayOpen: shared.IsTrayOpen,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (shared DropdownEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

func NewDropdownEntry(existingSelectorEntry ...*DropdownEntryType) DropdownEntryType {
	var menuBarEntry DropdownEntryType
	if existingSelectorEntry != nil {
		menuBarEntry.StyleEntry = existingSelectorEntry[0].StyleEntry
		menuBarEntry.SelectionEntry = existingSelectorEntry[0].SelectionEntry
		menuBarEntry.XLocation = existingSelectorEntry[0].XLocation
		menuBarEntry.YLocation = existingSelectorEntry[0].YLocation
		menuBarEntry.ItemWidth = existingSelectorEntry[0].ItemWidth
		menuBarEntry.ViewportPosition = existingSelectorEntry[0].ViewportPosition
		menuBarEntry.ItemHighlighted = existingSelectorEntry[0].ItemHighlighted
		menuBarEntry.ItemSelected = existingSelectorEntry[0].ItemSelected
		menuBarEntry.IsTrayOpen = existingSelectorEntry[0].IsTrayOpen
	}
	return menuBarEntry
}