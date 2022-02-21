package memory

import (
	"encoding/json"
)

type TextboxEntryType struct {
	HorizontalScrollbarAlias string
	VerticalScrollbarAlias string
	StyleEntry  TuiStyleEntryType
	TextData []string
	XLocation   int
	YLocation   int
	Width int
	Height            int
	ViewportXLocation int
	ViewportYLocation int
	IsEnabled         bool
	IsVisible bool
	IsBorderDrawn bool
}

func (shared TextboxEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		HorizontalScrollbarAlias string
		VerticalScrollbarAlias string
		StyleEntry TuiStyleEntryType
		TextData   []string
		XLocation  int
		YLocation   int
		Width int
		Height int
		ViewportXLocation int
		ViewportYLocation int
		IsEnabled bool
		IsVisible bool
		IsBorderDrawn bool
	}{
		HorizontalScrollbarAlias: shared.HorizontalScrollbarAlias,
		VerticalScrollbarAlias: shared.VerticalScrollbarAlias,
		StyleEntry: shared.StyleEntry,
		TextData:   shared.TextData,
		XLocation:  shared.XLocation,
		YLocation:  shared.YLocation,
		Width: shared.Width,
		Height: shared.Height,
		ViewportXLocation: shared.ViewportXLocation,
		ViewportYLocation: shared.ViewportYLocation,
		IsEnabled:  shared.IsEnabled,
		IsVisible:  shared.IsVisible,
		IsBorderDrawn: shared.IsBorderDrawn,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (shared TextboxEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

func NewTexboxEntry(existingButtonEntry ...*TextboxEntryType) TextboxEntryType {
	var textboxEntry TextboxEntryType
	if existingButtonEntry != nil {
		textboxEntry.HorizontalScrollbarAlias = existingButtonEntry[0].HorizontalScrollbarAlias
		textboxEntry.VerticalScrollbarAlias = existingButtonEntry[0].VerticalScrollbarAlias
		textboxEntry.StyleEntry = NewTuiStyleEntry(&existingButtonEntry[0].StyleEntry)
		textboxEntry.TextData = existingButtonEntry[0].TextData
		textboxEntry.XLocation = existingButtonEntry[0].XLocation
		textboxEntry.YLocation = existingButtonEntry[0].YLocation
		textboxEntry.Width = existingButtonEntry[0].Width
		textboxEntry.Height = existingButtonEntry[0].Height
		textboxEntry.ViewportXLocation = existingButtonEntry[0].ViewportXLocation
		textboxEntry.ViewportYLocation = existingButtonEntry[0].ViewportYLocation
		textboxEntry.IsEnabled = existingButtonEntry[0].IsEnabled
		textboxEntry.IsVisible = existingButtonEntry[0].IsVisible
		textboxEntry.IsBorderDrawn = existingButtonEntry[0].IsBorderDrawn
	}
	return textboxEntry
}

func IsTexboxEqual(sourceCheckboxEntry *TextboxEntryType, targetCheckboxEntry *TextboxEntryType) bool {
	if sourceCheckboxEntry.StyleEntry == targetCheckboxEntry.StyleEntry &&
		sourceCheckboxEntry.XLocation == targetCheckboxEntry.XLocation &&
		sourceCheckboxEntry.YLocation == targetCheckboxEntry.YLocation &&
		sourceCheckboxEntry.Width == targetCheckboxEntry.Width &&
		sourceCheckboxEntry.Height == targetCheckboxEntry.Height &&
		sourceCheckboxEntry.IsEnabled == targetCheckboxEntry.IsEnabled &&
		sourceCheckboxEntry.IsVisible == targetCheckboxEntry.IsVisible &&
		sourceCheckboxEntry.IsBorderDrawn == targetCheckboxEntry.IsBorderDrawn {
		return true
	}
	return false
}
