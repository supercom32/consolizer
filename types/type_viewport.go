package types

import (
	"encoding/json"
	"reflect"
)

// ViewportEntryType represents a read-only text viewport control.
// It supports scrollback history, text wrapping, and markup codes for text colorization.
type ViewportEntryType struct {
	BaseControlType
	HorizontalScrollbarAlias string
	VerticalScrollbarAlias   string
	TextData                 [][]rune
	ViewportXLocation        int
	ViewportYLocation        int
	MaxHistoryLines          int // Maximum number of lines to keep in history
	IsHistoryEnabled         bool // Whether to keep history beyond what's visible
	IsTransparent            bool // Whether the viewport background is transparent
	IsLinesWrapped           bool // Whether text is wrapped to fit the viewport width
}

// GetAlias allows you to retrieve the alias of a viewport control.
func (shared ViewportEntryType) GetAlias() string {
	return shared.Alias
}

// MarshalJSON allows you to serialize a viewport control to JSON.
func (shared ViewportEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		BaseControlType
		HorizontalScrollbarAlias string
		VerticalScrollbarAlias   string
		TextData                 [][]rune
		ViewportX                int
		ViewportY                int
		MaxHistoryLines          int
		IsHistoryEnabled         bool
		IsTransparent            bool
		IsLinesWrapped           bool
	}{
		BaseControlType:          shared.BaseControlType,
		HorizontalScrollbarAlias: shared.HorizontalScrollbarAlias,
		VerticalScrollbarAlias:   shared.VerticalScrollbarAlias,
		TextData:                 shared.TextData,
		ViewportX:                shared.ViewportXLocation,
		ViewportY:                shared.ViewportYLocation,
		MaxHistoryLines:          shared.MaxHistoryLines,
		IsHistoryEnabled:         shared.IsHistoryEnabled,
		IsTransparent:            shared.IsTransparent,
		IsLinesWrapped:           shared.IsLinesWrapped,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

// GetEntryAsJsonDump allows you to get a JSON string representation of a viewport control.
func (shared ViewportEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

// NewViewportEntry allows you to create a new viewport control.
func NewViewportEntry(existingViewportEntry ...*ViewportEntryType) ViewportEntryType {
	var viewportEntry ViewportEntryType
	viewportEntry.BaseControlType = NewBaseControl()
	viewportEntry.MaxHistoryLines = 1000 // Default max history lines
	viewportEntry.IsHistoryEnabled = true // Default to enabled history
	viewportEntry.IsTransparent = false // Default to non-transparent background
	viewportEntry.IsLinesWrapped = true // Default to wrapped text

	if existingViewportEntry != nil {
		viewportEntry.BaseControlType = existingViewportEntry[0].BaseControlType
		viewportEntry.HorizontalScrollbarAlias = existingViewportEntry[0].HorizontalScrollbarAlias
		viewportEntry.VerticalScrollbarAlias = existingViewportEntry[0].VerticalScrollbarAlias
		viewportEntry.TextData = existingViewportEntry[0].TextData
		viewportEntry.ViewportXLocation = existingViewportEntry[0].ViewportXLocation
		viewportEntry.ViewportYLocation = existingViewportEntry[0].ViewportYLocation
		viewportEntry.MaxHistoryLines = existingViewportEntry[0].MaxHistoryLines
		viewportEntry.IsHistoryEnabled = existingViewportEntry[0].IsHistoryEnabled
		viewportEntry.IsTransparent = existingViewportEntry[0].IsTransparent
		viewportEntry.IsLinesWrapped = existingViewportEntry[0].IsLinesWrapped
	}
	return viewportEntry
}

// IsViewportEntryEqual allows you to compare two viewport controls for equality.
func IsViewportEntryEqual(sourceViewportEntry *ViewportEntryType, targetViewportEntry *ViewportEntryType) bool {
	return sourceViewportEntry.BaseControlType.IsEqual(&targetViewportEntry.BaseControlType) &&
		sourceViewportEntry.HorizontalScrollbarAlias == targetViewportEntry.HorizontalScrollbarAlias &&
		sourceViewportEntry.VerticalScrollbarAlias == targetViewportEntry.VerticalScrollbarAlias &&
		reflect.DeepEqual(sourceViewportEntry.TextData, targetViewportEntry.TextData) &&
		sourceViewportEntry.ViewportXLocation == targetViewportEntry.ViewportXLocation &&
		sourceViewportEntry.ViewportYLocation == targetViewportEntry.ViewportYLocation &&
		sourceViewportEntry.MaxHistoryLines == targetViewportEntry.MaxHistoryLines &&
		sourceViewportEntry.IsHistoryEnabled == targetViewportEntry.IsHistoryEnabled &&
		sourceViewportEntry.IsTransparent == targetViewportEntry.IsTransparent &&
		sourceViewportEntry.IsLinesWrapped == targetViewportEntry.IsLinesWrapped
}

// GetViewportAlias allows you to retrieve the alias of a viewport control.
func GetViewportAlias(entry *ViewportEntryType) string {
	return entry.Alias
}
