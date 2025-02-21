package memory

import (
	"fmt"
	"github.com/supercom32/consolizer/types"
)

var Dropdowns = NewControlMemoryManager[*types.DropdownEntryType]()

func AddDropdown(layerAlias string, dropdownAlias string, styleEntry types.TuiStyleEntryType, selectionEntry types.SelectionEntryType, xLocation int, yLocation int, itemWidth int, itemSelected int) {
	dropdownEntry := types.NewDropdownEntry()
	dropdownEntry.Alias = dropdownAlias
	dropdownEntry.StyleEntry = styleEntry
	dropdownEntry.SelectionEntry = selectionEntry
	dropdownEntry.XLocation = xLocation
	dropdownEntry.YLocation = yLocation
	dropdownEntry.ItemWidth = itemWidth
	dropdownEntry.ItemSelected = itemSelected

	// Use the ControlMemoryManager to add the dropdown entry
	Dropdowns.Add(layerAlias, dropdownAlias, &dropdownEntry)
}

func DeleteDropdown(layerAlias string, dropdownAlias string) {
	// Use ControlMemoryManager to remove the dropdown entry
	Dropdowns.Remove(layerAlias, dropdownAlias)
}

func IsDropdownExists(layerAlias string, dropdownAlias string) bool {
	// Use ControlMemoryManager to check if the dropdown exists
	return Dropdowns.Get(layerAlias, dropdownAlias) != nil
}

func GetDropdown(layerAlias string, dropdownAlias string) *types.DropdownEntryType {
	// Get the dropdown entry using ControlMemoryManager
	dropdownEntry := Dropdowns.Get(layerAlias, dropdownAlias)
	if dropdownEntry == nil {
		panic(fmt.Sprintf("The selector '%s' under layer '%s' could not be obtained since it does not exist!", dropdownAlias, layerAlias))
	}
	return dropdownEntry
}

func DeleteAllDropdownsFromLayer(layerAlias string) {
	// Get all dropdown entries from the layer
	dropdowns := Dropdowns.GetAllEntries(layerAlias)

	// Loop through all entries and delete them
	for _, dropdown := range dropdowns {
		Dropdowns.Remove(layerAlias, dropdown.Alias)
	}
}
