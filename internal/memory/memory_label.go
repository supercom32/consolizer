package memory

import (
	"fmt"
	"supercom32.net/consolizer/types"
)

var Labels = NewControlMemoryManager[types.LabelEntryType]()

func AddLabel(layerAlias string, labelAlias string, labelValue string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int) {
	labelEntry := types.NewLabelEntry()
	labelEntry.StyleEntry = styleEntry
	labelEntry.Alias = labelAlias
	labelEntry.Value = labelValue
	labelEntry.XLocation = xLocation
	labelEntry.YLocation = yLocation
	labelEntry.Width = width

	// Use the ControlMemoryManager to add the label entry
	Labels.Add(layerAlias, labelAlias, &labelEntry)
}

func GetLabel(layerAlias string, labelAlias string) *types.LabelEntryType {
	// Get the label entry using ControlMemoryManager
	labelEntry := Labels.Get(layerAlias, labelAlias)
	if labelEntry == nil {
		panic(fmt.Sprintf("The requested label with alias '%s' on layer '%s' could not be returned since it does not exist.", labelAlias, layerAlias))
	}
	return labelEntry
}

func IsLabelExists(layerAlias string, labelAlias string) bool {
	// Use ControlMemoryManager to check if the label exists
	return Labels.Get(layerAlias, labelAlias) != nil
}

func DeleteLabel(layerAlias string, labelAlias string) {
	// Use ControlMemoryManager to remove the label entry
	Labels.Remove(layerAlias, labelAlias)
}

func DeleteAllLabelsFromLayer(layerAlias string) {
	// Get all label entries from the layer
	labels := Labels.GetAllEntries(layerAlias)

	// Loop through all entries and delete them
	for _, label := range labels {
		Labels.Remove(layerAlias, label.Alias) // Assuming label.Alias is used as the alias
	}
}
