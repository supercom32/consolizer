package memory

import (
	"fmt"
	"supercom32.net/consolizer/types"
)

var TextFields = NewControlMemoryManager[types.TextFieldEntryType]()

func AddTextField(layerAlias string, textFieldAlias string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, maxLengthAllowed int, IsPasswordProtected bool, defaultValue string, isEnabled bool) {
	textFieldEntry := types.NewTextFieldEntry()
	textFieldEntry.Alias = textFieldAlias
	textFieldEntry.StyleEntry = styleEntry
	textFieldEntry.XLocation = xLocation
	textFieldEntry.YLocation = yLocation
	textFieldEntry.Width = width
	textFieldEntry.MaxLengthAllowed = maxLengthAllowed
	textFieldEntry.IsPasswordProtected = IsPasswordProtected
	textFieldEntry.CurrentValue = []rune(defaultValue)
	textFieldEntry.DefaultValue = defaultValue
	textFieldEntry.IsEnabled = isEnabled

	// Use the generic memory manager to add the text field entry
	TextFields.Add(layerAlias, textFieldAlias, &textFieldEntry)
}

func DeleteTextField(layerAlias string, textFieldAlias string) {
	// Use the generic memory manager to remove the text field entry
	TextFields.Remove(layerAlias, textFieldAlias)
}

func DeleteAllTextFieldsFromLayer(layerAlias string) {
	// Retrieve all text fields in the specified layer
	textFields := TextFields.GetAllEntries(layerAlias)

	// Loop through all entries and delete them
	for _, textField := range textFields {
		TextFields.Remove(layerAlias, textField.Alias) // Assuming textField.Alias contains the alias
	}
}

func IsTextFieldExists(layerAlias string, textFieldAlias string) bool {
	// Use the generic memory manager to check existence
	return TextFields.Get(layerAlias, textFieldAlias) != nil
}

func GetTextField(layerAlias string, textFieldAlias string) *types.TextFieldEntryType {
	// Use the generic memory manager to retrieve the text field entry
	textFieldEntry := TextFields.Get(layerAlias, textFieldAlias)
	if textFieldEntry == nil {
		panic(fmt.Sprintf("The text field '%s' under layer '%s' could not be obtained since it does not exist!", textFieldAlias, layerAlias))
	}
	return textFieldEntry
}
