package memory

import "fmt"

var TextFieldMemory map[string]map[string]*TextFieldEntryType

func InitializeTextFieldMemory() {
	TextFieldMemory = make(map[string]map[string]*TextFieldEntryType)
}

func AddTextField(layerAlias string, textFieldAlias string, styleEntry TuiStyleEntryType, xLocation int, yLocation int, width int, maxLengthAllowed int, IsPasswordProtected bool, defaultValue string) {
	textFieldEntry := NewTextFieldEntry()
	textFieldEntry.StyleEntry = styleEntry
	textFieldEntry.XLocation = xLocation
	textFieldEntry.YLocation = yLocation
	textFieldEntry.Width = width
	textFieldEntry.MaxLengthAllowed = maxLengthAllowed
	textFieldEntry.IsPasswordProtected = IsPasswordProtected
	textFieldEntry.CurrentValue = []rune(defaultValue)
	textFieldEntry.DefaultValue = defaultValue
	if TextFieldMemory[layerAlias] == nil {
		TextFieldMemory[layerAlias] = make(map[string]*TextFieldEntryType)
	}
	TextFieldMemory[layerAlias][textFieldAlias] = &textFieldEntry
}

func DeleteTextField(layerAlias string, textFieldAlias string) {
	delete(TextFieldMemory[layerAlias], textFieldAlias)
}

func IsTextFieldExists(layerAlias string, textFieldAlias string) bool {
	if _, isExist := TextFieldMemory[layerAlias][textFieldAlias]; isExist {
		return true
	}
	return false
}

func GetTextField(layerAlias string, textFieldAlias string) *TextFieldEntryType{
	if !IsTextFieldExists(layerAlias, textFieldAlias) {
		panic(fmt.Sprintf("The text field '%s' under layer '%s' could not be obtained since it does not exist!", textFieldAlias,  layerAlias))
	}
	return TextFieldMemory[layerAlias][textFieldAlias]
}