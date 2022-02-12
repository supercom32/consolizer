package memory

import (
	"encoding/json"
	"sync"
)

// func DrawButton(LayerAlias string, ButtonLabel string, StyleEntry TuiStyleEntryType, IsPressed bool, IsSelected bool, XLocation int, YLocation int, Width int, Length int) {
type ScrollBarEntryType struct {
	Mutex                   sync.Mutex
	StyleEntry              TuiStyleEntryType
	ScrollBarAlias          string
	XLocation               int
	YLocation               int
	Length         int
	MaxScrollValue int
	ScrollValue    int
	HandlePosition int
	IsVisible bool
	IsHorizontal            bool
}

func (shared ScrollBarEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		StyleEntry  TuiStyleEntryType
		ScrollBarAlias string
		XLocation   int
		YLocation   int
		Length      int
		MaxScrollValue   int
		ScrollValue int
		ScrollPosition int
		IsVisible bool
		IsHorizontal bool
	}{
		StyleEntry: shared.StyleEntry,
		ScrollBarAlias: shared.ScrollBarAlias,
		XLocation: shared.XLocation,
		YLocation: shared.YLocation,
		Length: shared.Length,
		MaxScrollValue: shared.MaxScrollValue,
		ScrollPosition: shared.HandlePosition,
		IsVisible: shared.IsVisible,
		IsHorizontal: shared.IsHorizontal,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (shared ScrollBarEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

func NewScrollBarEntry(existingScrollBarEntry ...*ScrollBarEntryType) ScrollBarEntryType {
	var scrollBarEntry ScrollBarEntryType
	if existingScrollBarEntry != nil {
		scrollBarEntry.StyleEntry = NewTuiStyleEntry(&existingScrollBarEntry[0].StyleEntry)
		scrollBarEntry.XLocation = existingScrollBarEntry[0].XLocation
		scrollBarEntry.YLocation = existingScrollBarEntry[0].YLocation
		scrollBarEntry.Length = existingScrollBarEntry[0].Length
		scrollBarEntry.MaxScrollValue = existingScrollBarEntry[0].MaxScrollValue
		scrollBarEntry.ScrollValue = existingScrollBarEntry[0].ScrollValue
		scrollBarEntry.HandlePosition = existingScrollBarEntry[0].HandlePosition
		scrollBarEntry.IsVisible = existingScrollBarEntry[0].IsVisible
		scrollBarEntry.IsHorizontal = existingScrollBarEntry[0].IsHorizontal
	}
	return scrollBarEntry
}

func IsScrollBarEntryEqual(sourceScrollBarEntry *ScrollBarEntryType, targetScrollBarEntry *ScrollBarEntryType) bool {
	if sourceScrollBarEntry.StyleEntry == targetScrollBarEntry.StyleEntry &&
		sourceScrollBarEntry.XLocation == targetScrollBarEntry.XLocation &&
		sourceScrollBarEntry.YLocation == targetScrollBarEntry.YLocation &&
		sourceScrollBarEntry.Length == targetScrollBarEntry.Length &&
		sourceScrollBarEntry.MaxScrollValue == targetScrollBarEntry.MaxScrollValue &&
		sourceScrollBarEntry.ScrollValue == targetScrollBarEntry.ScrollValue &&
		sourceScrollBarEntry.HandlePosition == targetScrollBarEntry.HandlePosition &&
		sourceScrollBarEntry.IsVisible == targetScrollBarEntry.IsVisible &&
		sourceScrollBarEntry.IsHorizontal == targetScrollBarEntry.IsHorizontal {
		return true
	}
	return false
}
