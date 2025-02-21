package memory

import (
	"fmt"
	"github.com/supercom32/consolizer/types"
)

var ProgressBars = NewControlMemoryManager[*types.ProgressBarEntryType]()

func AddProgressBar(layerAlias string, progressBarAlias string, progressBarLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, value int, maxValue int, isBackgroundTransparent bool) {
	progressBarEntry := types.NewProgressBarEntry()
	progressBarEntry.StyleEntry = styleEntry
	progressBarEntry.Alias = progressBarAlias
	progressBarEntry.Label = progressBarLabel
	progressBarEntry.Value = value
	progressBarEntry.MaxValue = maxValue
	progressBarEntry.IsBackgroundTransparent = isBackgroundTransparent
	progressBarEntry.XLocation = xLocation
	progressBarEntry.YLocation = yLocation
	progressBarEntry.Width = width
	progressBarEntry.Height = height

	// Use the ControlMemoryManager to add the progress bar entry
	ProgressBars.Add(layerAlias, progressBarAlias, &progressBarEntry)
}

func GetProgressBar(layerAlias string, progressBarAlias string) *types.ProgressBarEntryType {
	// Get the progress bar entry using ControlMemoryManager
	progressBarEntry := ProgressBars.Get(layerAlias, progressBarAlias)
	if progressBarEntry == nil {
		panic(fmt.Sprintf("The requested progress bar with alias '%s' on layer '%s' could not be returned since it does not exist.", progressBarAlias, layerAlias))
	}
	return progressBarEntry
}

func IsProgressBarExists(layerAlias string, progressBarAlias string) bool {
	// Use ControlMemoryManager to check if the progress bar exists
	return ProgressBars.Get(layerAlias, progressBarAlias) != nil
}

func DeleteProgressBar(layerAlias string, progressBarAlias string) {
	// Use ControlMemoryManager to remove the progress bar entry
	ProgressBars.Remove(layerAlias, progressBarAlias)
}

func DeleteAllProgressBarsFromLayer(layerAlias string) {
	// Get all progress bar entries from the layer
	progressBars := ProgressBars.GetAllEntries(layerAlias)

	// Loop through all entries and delete them
	for _, progressBar := range progressBars {
		ProgressBars.Remove(layerAlias, progressBar.Alias) // Assuming progressBar.Alias is used as the alias
	}
}
