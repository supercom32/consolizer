package memory

import "fmt"

var DropdownMemory map[string]map[string]*DropdownEntryType

func InitializeDropdownMemory() {
	DropdownMemory = make(map[string]map[string]*DropdownEntryType)
}

func AddDropdown(layerAlias string, menuBarAlias string, styleEntry TuiStyleEntryType, selectionEntry SelectionEntryType, xLocation int, yLocation int, selectorHeight int, itemWidth int, numberOfColumns int, viewportPosition int, itemSelected int) {
	dropDownEntry := NewDropdownEntry()
	dropDownEntry.StyleEntry = styleEntry
	dropDownEntry.SelectionEntry = selectionEntry
	dropDownEntry.XLocation = xLocation
	dropDownEntry.YLocation = yLocation
	dropDownEntry.SelectorHeight = selectorHeight
	dropDownEntry.ItemWidth = itemWidth
	dropDownEntry.ViewportPosition = viewportPosition
	dropDownEntry.ItemHighlighted = itemSelected
	if DropdownMemory[layerAlias] == nil {
		DropdownMemory[layerAlias] = make(map[string]*DropdownEntryType)
	}
	DropdownMemory[layerAlias][menuBarAlias] = &dropDownEntry
}

func DeleteDropdown(layerAlias string, menuBarAlias string) {
	delete(DropdownMemory[layerAlias], menuBarAlias)
}

func IsDropdownExists(layerAlias string, menuBarAlias string) bool {
	if _, isExist := DropdownMemory[layerAlias][menuBarAlias]; isExist {
		return true
	}
	return false
}

func GetDropdown(layerAlias string, selectorAlias string) *DropdownEntryType {
	if !IsDropdownExists(layerAlias, selectorAlias) {
		panic(fmt.Sprintf("The selector '%s' under layer '%s' could not be obtained since it does not exist!", selectorAlias,  layerAlias))
	}
	return DropdownMemory[layerAlias][selectorAlias]
}