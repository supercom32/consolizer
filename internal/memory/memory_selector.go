package memory

import "fmt"

var SelectorMemory map[string]map[string]*SelectorEntryType

func InitializeSelectorMemory() {
	SelectorMemory = make(map[string]map[string]*SelectorEntryType)
}

func AddSelector(layerAlias string, selectorAlias string, styleEntry TuiStyleEntryType, selectionEntry SelectionEntryType, xLocation int, yLocation int, selectorHeight int, itemWidth int, numberOfColumns int, viewportPosition int, itemSelected int, isBorderDrawn bool) {
	selectorEntry := NewSelectorEntry()
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
	if SelectorMemory[layerAlias] == nil {
		SelectorMemory[layerAlias] = make(map[string]*SelectorEntryType)
	}
	SelectorMemory[layerAlias][selectorAlias] = &selectorEntry
}

func DeleteSelector(layerAlias string, selectorAlias string) {
	delete(SelectorMemory[layerAlias], selectorAlias)
}

func IsSelectorExists(layerAlias string, selectorAlias string) bool {
	if _, isExist := SelectorMemory[layerAlias][selectorAlias]; isExist {
		return true
	}
	return false
}

func GetSelector(layerAlias string, selectorAlias string) *SelectorEntryType {
	if !IsSelectorExists(layerAlias, selectorAlias) {
		panic(fmt.Sprintf("The selector '%s' under layer '%s' could not be obtained since it does not exist!", selectorAlias,  layerAlias))
	}
	return SelectorMemory[layerAlias][selectorAlias]
}