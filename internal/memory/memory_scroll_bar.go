package memory

import (
	"fmt"
	"github.com/supercom32/consolizer/types"
)

var ScrollBars = NewControlMemoryManager[*types.ScrollbarEntryType]()

func AddScrollbar(layerAlias string, scrollBarAlias string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, length int, maxScrollValue int, scrollValue int, scrollIncrement int, isHorizontal bool) {
	scrollbarEntry := types.NewScrollbarEntry()
	scrollbarEntry.ScrollBarAlias = scrollBarAlias
	scrollbarEntry.StyleEntry = styleEntry
	scrollbarEntry.XLocation = xLocation
	scrollbarEntry.YLocation = yLocation
	scrollbarEntry.Length = length
	scrollbarEntry.MaxScrollValue = maxScrollValue - 1 // Adjusted for 0-based indexing
	scrollbarEntry.ScrollValue = scrollValue
	scrollbarEntry.IsVisible = true
	scrollbarEntry.IsEnabled = true
	scrollbarEntry.IsHorizontal = isHorizontal
	scrollbarEntry.ScrollIncrement = scrollIncrement

	// Use the generic memory manager to add the scrollbar entry
	ScrollBars.Add(layerAlias, scrollBarAlias, &scrollbarEntry)
}

func GetScrollbar(layerAlias string, scrollBarAlias string) *types.ScrollbarEntryType {
	// Use the generic memory manager to retrieve the scrollbar entry
	scrollbarEntry := ScrollBars.Get(layerAlias, scrollBarAlias)
	if scrollbarEntry == nil {
		panic(fmt.Sprintf("The requested scroll bar with alias '%s' on layer '%s' could not be returned since it does not exist.", scrollBarAlias, layerAlias))
	}
	return scrollbarEntry
}

func IsScrollbarExists(layerAlias string, scrollBarAlias string) bool {
	// Check existence using the generic memory manager
	return ScrollBars.Get(layerAlias, scrollBarAlias) != nil
}

func DeleteScrollbar(layerAlias string, scrollBarAlias string) {
	// Use the generic memory manager to delete the scrollbar entry
	ScrollBars.Remove(layerAlias, scrollBarAlias)
}

func DeleteAllScrollbarsFromLayer(layerAlias string) {
	// Retrieve all scrollbars in the specified layer
	scrollbars := ScrollBars.GetAllEntries(layerAlias)

	// Loop through all entries and delete them
	for _, scrollbar := range scrollbars {
		ScrollBars.Remove(layerAlias, scrollbar.ScrollBarAlias) // Assuming scrollbar.Alias contains the alias
	}
}

func compareScrollbarAlias(a, b *types.ScrollbarEntryType) bool {
	return a.ScrollBarAlias < b.ScrollBarAlias // Compare lexicographically in ascending order
}
