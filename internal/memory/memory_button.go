package memory

import (
	"fmt"
	"github.com/supercom32/consolizer/types"
)

var Buttons = NewControlMemoryManager[*types.ButtonEntryType]()

func AddButton(layerAlias string, buttonAlias string, buttonLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int) {
	buttonEntry := types.NewButtonEntry()
	buttonEntry.StyleEntry = styleEntry
	buttonEntry.ButtonAlias = buttonAlias
	buttonEntry.ButtonLabel = buttonLabel
	buttonEntry.XLocation = xLocation
	buttonEntry.YLocation = yLocation
	buttonEntry.IsEnabled = true
	buttonEntry.Width = width
	buttonEntry.Height = height

	// Use the ControlMemoryManager to handle button entries
	Buttons.Add(layerAlias, buttonAlias, &buttonEntry)
}

func GetButton(layerAlias string, buttonAlias string) *types.ButtonEntryType {
	// Get the button entry using the ControlMemoryManager
	buttonEntry := Buttons.Get(layerAlias, buttonAlias)
	if buttonEntry == nil {
		panic(fmt.Sprintf("The requested button with alias '%s' on layer '%s' could not be returned since it does not exist.", buttonAlias, layerAlias))
	}
	return buttonEntry
}

func IsButtonExists(layerAlias string, buttonAlias string) bool {
	// Check existence using ControlMemoryManager
	return Buttons.Get(layerAlias, buttonAlias) != nil
}

func DeleteButton(layerAlias string, buttonAlias string) {
	// Use the ControlMemoryManager to delete the button entry
	Buttons.Remove(layerAlias, buttonAlias)
}

func DeleteAllButtonsFromLayer(layerAlias string) {
	// Retrieve all buttons from the layer and delete them
	buttons := Buttons.GetAllEntries(layerAlias)
	for _, button := range buttons {
		// Assume the button has an alias function or equivalent field to remove
		Buttons.Remove(layerAlias, button.ButtonAlias)
	}
}
