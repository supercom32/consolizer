package memory

import (
	"fmt"
	"github.com/supercom32/consolizer/types"
)

var Selectors = NewControlMemoryManager[*types.SelectorEntryType]()

func AddSelector(layerAlias string, selectorAlias string, styleEntry types.TuiStyleEntryType, selectionEntry types.SelectionEntryType, xLocation int, yLocation int, selectorHeight int, itemWidth int, numberOfColumns int, viewportPosition int, itemSelected int, isBorderDrawn bool) {
	selectorEntry := types.NewSelectorEntry()
	selectorEntry.SelectorAlias = selectorAlias
	selectorEntry.StyleEntry = styleEntry
	selectorEntry.SelectionEntry = selectionEntry
	selectorEntry.XLocation = xLocation
	selectorEntry.YLocation = yLocation
	selectorEntry.SelectorHeight = selectorHeight
	selectorEntry.ItemWidth = itemWidth
	selectorEntry.NumberOfColumns = numberOfColumns
	selectorEntry.ViewportPosition = viewportPosition
	selectorEntry.ItemHighlighted = itemSelected
	selectorEntry.IsBorderDrawn = isBorderDrawn
	selectorEntry.IsVisible = true

	// Use the generic memory manager to add the selector entry
	Selectors.Add(layerAlias, selectorAlias, &selectorEntry)
}

func DeleteSelector(layerAlias string, selectorAlias string) {
	// Use the generic memory manager to remove the selector entry
	Selectors.Remove(layerAlias, selectorAlias)
}

func DeleteAllSelectorsFromLayer(layerAlias string) {
	// Retrieve all selectors in the specified layer
	selectors := Selectors.GetAllEntries(layerAlias)

	// Loop through all entries and delete them
	for _, selector := range selectors {
		Selectors.Remove(layerAlias, selector.SelectorAlias) // Assuming selector.Alias contains the alias
	}
}

func IsSelectorExists(layerAlias string, selectorAlias string) bool {
	// Use the generic memory manager to check existence
	return Selectors.Get(layerAlias, selectorAlias) != nil
}

func GetSelector(layerAlias string, selectorAlias string) *types.SelectorEntryType {
	// Use the generic memory manager to retrieve the selector entry
	selectorEntry := Selectors.Get(layerAlias, selectorAlias)
	if selectorEntry == nil {
		panic(fmt.Sprintf("The selector '%s' under layer '%s' could not be obtained since it does not exist!", selectorAlias, layerAlias))
	}
	return selectorEntry
}
