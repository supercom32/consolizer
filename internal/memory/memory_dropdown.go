package memory

import "fmt"

var DropdownMemory map[string]map[string]*DropdownEntryType

func InitializeDropdownMemory() {
	DropdownMemory = make(map[string]map[string]*DropdownEntryType)
}

func AddDropdown(layerAlias string, dropdownAlias string, styleEntry TuiStyleEntryType, selectionEntry SelectionEntryType, xLocation int, yLocation int, itemWidth int, itemSelected int) {
	dropDownEntry := NewDropdownEntry()
	dropDownEntry.StyleEntry = styleEntry
	dropDownEntry.SelectionEntry = selectionEntry
	dropDownEntry.XLocation = xLocation
	dropDownEntry.YLocation = yLocation
	dropDownEntry.ItemWidth = itemWidth
	dropDownEntry.ItemSelected = itemSelected
	if DropdownMemory[layerAlias] == nil {
		DropdownMemory[layerAlias] = make(map[string]*DropdownEntryType)
	}
	DropdownMemory[layerAlias][dropdownAlias] = &dropDownEntry
}

func DeleteDropdown(layerAlias string, dropdownAlias string) {
	delete(DropdownMemory[layerAlias], dropdownAlias)
}

func IsDropdownExists(layerAlias string, dropdownAlias string) bool {
	if _, isExist := DropdownMemory[layerAlias][dropdownAlias]; isExist {
		return true
	}
	return false
}

func GetDropdown(layerAlias string, dropdownAlias string) *DropdownEntryType {
	if !IsDropdownExists(layerAlias, dropdownAlias) {
		panic(fmt.Sprintf("The selector '%s' under layer '%s' could not be obtained since it does not exist!", dropdownAlias,  layerAlias))
	}
	return DropdownMemory[layerAlias][dropdownAlias]
}