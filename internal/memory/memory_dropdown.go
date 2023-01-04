package memory

import (
	"fmt"
	"github.com/supercom32/consolizer/types"
	"sync"
)

type dropdownMemoryType struct {
	sync.Mutex
	Entries map[string]map[string]*types.DropdownEntryType
}

var Dropdown dropdownMemoryType

func InitializeDropdownMemory() {
	Dropdown.Entries = make(map[string]map[string]*types.DropdownEntryType)
}

func AddDropdown(layerAlias string, dropdownAlias string, styleEntry types.TuiStyleEntryType, selectionEntry types.SelectionEntryType, xLocation int, yLocation int, itemWidth int, itemSelected int) {
	Dropdown.Lock()
	defer func() {
		Dropdown.Unlock()
	}()
	dropDownEntry := types.NewDropdownEntry()
	dropDownEntry.StyleEntry = styleEntry
	dropDownEntry.SelectionEntry = selectionEntry
	dropDownEntry.XLocation = xLocation
	dropDownEntry.YLocation = yLocation
	dropDownEntry.ItemWidth = itemWidth
	dropDownEntry.ItemSelected = itemSelected
	if Dropdown.Entries[layerAlias] == nil {
		Dropdown.Entries[layerAlias] = make(map[string]*types.DropdownEntryType)
	}
	Dropdown.Entries[layerAlias][dropdownAlias] = &dropDownEntry
}

func DeleteDropdown(layerAlias string, dropdownAlias string) {
	Dropdown.Lock()
	defer func() {
		Dropdown.Unlock()
	}()
	delete(Dropdown.Entries[layerAlias], dropdownAlias)
}

func IsDropdownExists(layerAlias string, dropdownAlias string) bool {
	Dropdown.Lock()
	defer func() {
		Dropdown.Unlock()
	}()
	if _, isExist := Dropdown.Entries[layerAlias][dropdownAlias]; isExist {
		return true
	}
	return false
}

func GetDropdown(layerAlias string, dropdownAlias string) *types.DropdownEntryType {
	Dropdown.Lock()
	defer func() {
		Dropdown.Unlock()
	}()
	if _, isExist := Dropdown.Entries[layerAlias][dropdownAlias]; !isExist {
		panic(fmt.Sprintf("The selector '%s' under layer '%s' could not be obtained since it does not exist!", dropdownAlias, layerAlias))
	}
	return Dropdown.Entries[layerAlias][dropdownAlias]
}

func DeleteAllDropdownsFromLayer(layerAlias string) {
	Dropdown.Lock()
	defer func() {
		Dropdown.Unlock()
	}()
	for entryToRemove := range Dropdown.Entries[layerAlias] {
		delete(Dropdown.Entries[layerAlias], entryToRemove)
	}
}
