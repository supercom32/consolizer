package consolizer

import (
	"fmt"
	"supercom32.net/consolizer/constants"
	"supercom32.net/consolizer/types"
)

func validateTextFieldWidth(width int) {
	if width <= 0 {
		safeSttyPanic(fmt.Sprintf("The specified text field width '%d' is invalid.", width))
	}
}

func validateLayerLocationByLayerAlias(layerAlias string, xLocation int, yLocation int) {
	validateLayer(layerAlias)
	layerEntry := Layers.Get(layerAlias)
	validateLayerLocationByLayerEntry(layerEntry, xLocation, yLocation)
}

func validateSelectionEntry(selectionEntry types.SelectionEntryType) {
	if len(selectionEntry.SelectionValue) == 0 {
		safeSttyPanic(fmt.Sprintf("The selection entry passed was empty."))
	}
}
func validateLayerLocationByLayerEntry(layerEntry *types.LayerEntryType, xLocation int, yLocation int) {
	if xLocation < 0 || yLocation < 0 ||
		xLocation >= layerEntry.Width || yLocation >= layerEntry.Height {
		safeSttyPanic(fmt.Sprintf("The specified location (%d, %d) is out of bounds for layer with a size of (%d, %d).", xLocation, yLocation, layerEntry.Width, layerEntry.Height))
	}
}

func validateRGBColorIndex(redColorIndex int32, greenColorIndex int32, blueColorIndex int32) {
	if redColorIndex < 0 || redColorIndex > 255 || greenColorIndex < 0 || greenColorIndex > 255 ||
		blueColorIndex < 0 || blueColorIndex > 255 {
		safeSttyPanic(fmt.Sprintf("The specified RGB color index '%d, %d, %d' is invalid!", redColorIndex, greenColorIndex, blueColorIndex))
	}
}

func validateColorIndex(colorIndex int) {
	if colorIndex < 0 || colorIndex > len(constants.AnsiColorByIndex) {
		safeSttyPanic(fmt.Sprintf("The specified color index '%d' is invalid!", colorIndex))
	}
}

func validateTextStyleExists(textStyleAlias string) {
	if !IsTextStyleExists(textStyleAlias) {
		safeSttyPanic(fmt.Sprintf("The specified text style '%s' does not exist.", textStyleAlias))
	}
}

func validateDefaultLayerIsNotEmpty() {
	if commonResource.layerInstance.layerAlias == "" {
		safeSttyPanic(fmt.Sprintf("The action could not be completed since no default text layer exists!"))
	}
}

func validateTerminalWidthAndHeight(width int, height int) {
	if width <= 0 || height <= 0 {
		safeSttyPanic(fmt.Sprintf("The specified terminal width and height of '%d, %d' is invalid!", width, height))
	}
}

func validateLayer(layerAlias string) {
	if !Layers.IsExists(layerAlias) {
		safeSttyPanic(fmt.Sprintf("The specified layer '%s' does not exist.", layerAlias))
	}
}

func validatorTextField(layerAlias string, textFieldAlias string) {
	if !(TextFields.IsExists(layerAlias, textFieldAlias)) {
		safeSttyPanic(fmt.Sprintf("The text field '%s' under layer '%s' could not be obtained since it does not exist!", textFieldAlias, layerAlias))
	}
}

func validatorMenu(layerAlias string, menuAlias string) {
	if !(Selectors.IsExists(layerAlias, menuAlias)) {
		safeSttyPanic(fmt.Sprintf("The menu '%s' under layer '%s' could not be obtained since it does not exist!", menuAlias, layerAlias))
	}
}

func safeSttyPanic(panicMessage string) {
	RestoreTerminalSettings()
	panic(panicMessage)
}
