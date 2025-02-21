package memory

import (
	"fmt"
	"github.com/supercom32/consolizer/types"
)

var Textboxes = NewControlMemoryManager[*types.TextboxEntryType]()

func AddTextbox(layerAlias string, textboxAlias string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isBorderDrawn bool) {
	textboxEntry := types.NewTexboxEntry()
	textboxEntry.Alias = textboxAlias
	textboxEntry.StyleEntry = styleEntry
	textboxEntry.XLocation = xLocation
	textboxEntry.YLocation = yLocation
	textboxEntry.Width = width
	textboxEntry.Height = height
	textboxEntry.IsBorderDrawn = isBorderDrawn

	// Use the generic memory manager to add the textbox entry
	Textboxes.Add(layerAlias, textboxAlias, &textboxEntry)
}

func GetTextbox(layerAlias string, textboxAlias string) *types.TextboxEntryType {
	// Use the generic memory manager to retrieve the textbox entry
	textboxEntry := Textboxes.Get(layerAlias, textboxAlias)
	if textboxEntry == nil {
		panic(fmt.Sprintf("The requested text with alias '%s' on layer '%s' could not be returned since it does not exist.", textboxAlias, layerAlias))
	}
	return textboxEntry
}

func IsTextboxExists(layerAlias string, textboxAlias string) bool {
	// Use the generic memory manager to check existence
	return Textboxes.Get(layerAlias, textboxAlias) != nil
}

func DeleteTextbox(layerAlias string, textboxAlias string) {
	// Use the generic memory manager to remove the textbox entry
	Textboxes.Remove(layerAlias, textboxAlias)
}

func DeleteAllTextboxesFromLayer(layerAlias string) {
	// Retrieve all textboxes in the specified layer
	textboxes := Textboxes.GetAllEntries(layerAlias)

	// Loop through all entries and delete them
	for _, textbox := range textboxes {
		Textboxes.Remove(layerAlias, textbox.Alias) // Assuming textbox.Alias contains the alias
	}
}
