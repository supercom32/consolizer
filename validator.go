package consolizer

import (
	"fmt"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/types"
)

/*
validateTextFieldWidth is a method which allows you to validate that a specified text field width is greater than zero.

:param width: The width value to be validated.

Example:

	validateTextFieldWidth(20)
*/
func validateTextFieldWidth(width int) {
	if width <= 0 {
		safeSttyPanic(fmt.Sprintf("The specified text field width '%d' is invalid.", width))
	}
}

/*
validateLayerLocationByLayerAlias is a method which allows you to validate that a specific coordinate location is within
the bounds of a layer, identified by its alias.

:param layerAlias: The alias of the layer to check the location against.
:param xLocation: The X coordinate to validate.
:param yLocation: The Y coordinate to validate.

Example:

	validateLayerLocationByLayerAlias("mainLayer", 10, 5)
*/
func validateLayerLocationByLayerAlias(layerAlias string, xLocation int, yLocation int) {
	validateLayer(layerAlias)
	layerEntry := Layers.Get(layerAlias)
	validateLayerLocationByLayerEntry(layerEntry, xLocation, yLocation)
}

/*
validateSelectionEntry is a method which allows you to validate that a selection entry is not empty.

:param selectionEntry: The selection entry to be validated.

Example:

	validateSelectionEntry(mySelection)
*/
func validateSelectionEntry(selectionEntry types.SelectionEntryType) {
	if len(selectionEntry.SelectionValue) == 0 {
		safeSttyPanic(fmt.Sprintf("The selection entry passed was empty."))
	}
}

/*
validateLayerLocationByLayerEntry is a method which allows you to validate that a specific coordinate location is within
the bounds of a provided layer entry.

:param layerEntry: The layer entry to check the location against.
:param xLocation: The X coordinate to validate.
:param yLocation: The Y coordinate to validate.

Example:

	validateLayerLocationByLayerEntry(myLayer, 10, 5)
*/
func validateLayerLocationByLayerEntry(layerEntry *types.LayerEntryType, xLocation int, yLocation int) {
	if xLocation < 0 || yLocation < 0 ||
		xLocation >= layerEntry.Width || yLocation >= layerEntry.Height {
		safeSttyPanic(fmt.Sprintf("The specified location (%d, %d) is out of bounds for layer with a size of (%d, %d).", xLocation, yLocation, layerEntry.Width, layerEntry.Height))
	}
}

/*
validateRGBColorIndex is a method which allows you to validate that the provided RGB color components are within the
valid range of 0 to 255.

:param redColorIndex: The red component value.
:param greenColorIndex: The green component value.
:param blueColorIndex: The blue component value.

Example:

	validateRGBColorIndex(255, 128, 0)
*/
func validateRGBColorIndex(redColorIndex int32, greenColorIndex int32, blueColorIndex int32) {
	if redColorIndex < 0 || redColorIndex > 255 || greenColorIndex < 0 || greenColorIndex > 255 ||
		blueColorIndex < 0 || blueColorIndex > 255 {
		safeSttyPanic(fmt.Sprintf("The specified RGB color index '%d, %d, %d' is invalid!", redColorIndex, greenColorIndex, blueColorIndex))
	}
}

/*
validateColorIndex is a method which allows you to validate that a color index is within the range of defined ANSI
colors.

:param colorIndex: The color index to be validated.

Example:

	validateColorIndex(7)
*/
func validateColorIndex(colorIndex int) {
	if colorIndex < 0 || colorIndex > len(constants.AnsiColorByIndex) {
		safeSttyPanic(fmt.Sprintf("The specified color index '%d' is invalid!", colorIndex))
	}
}

/*
validateTextStyleExists is a method which allows you to validate that a text style with the specified alias exists.

:param textStyleAlias: The alias of the text style to check for.

Example:

	validateTextStyleExists("boldStyle")
*/
func validateTextStyleExists(textStyleAlias string) {
	if !IsTextStyleExists(textStyleAlias) {
		safeSttyPanic(fmt.Sprintf("The specified text style '%s' does not exist.", textStyleAlias))
	}
}

/*
validateDefaultLayerIsNotEmpty is a method which allows you to validate that a default layer has been established and is
not empty.

Example:

	validateDefaultLayerIsNotEmpty()
*/
func validateDefaultLayerIsNotEmpty() {
	if commonResource.layerInstance.layerAlias == "" {
		safeSttyPanic(fmt.Sprintf("The action could not be completed since no default text layer exists!"))
	}
}

/*
validateTerminalWidthAndHeight is a method which allows you to validate that the specified terminal width and height are
greater than zero.

:param width: The terminal width to validate.
:param height: The terminal height to validate.

Example:

	validateTerminalWidthAndHeight(80, 24)
*/
func validateTerminalWidthAndHeight(width int, height int) {
	if width <= 0 || height <= 0 {
		safeSttyPanic(fmt.Sprintf("The specified terminal width and height of '%d, %d' is invalid!", width, height))
	}
}

/*
validateLayer is a method which allows you to validate that a layer with the specified alias exists.

:param layerAlias: The alias of the layer to check for.

Example:

	validateLayer("backgroundLayer")
*/
func validateLayer(layerAlias string) {
	if !Layers.IsExists(layerAlias) {
		safeSttyPanic(fmt.Sprintf("The specified layer '%s' does not exist.", layerAlias))
	}
}

/*
validatorTextField is a method which allows you to validate that a text field with the specified alias exists under the
given layer alias.

:param layerAlias: The alias of the layer the text field belongs to.
:param textFieldAlias: The alias of the text field to check for.

Example:

	validatorTextField("loginLayer", "usernameField")
*/
func validatorTextField(layerAlias string, textFieldAlias string) {
	if !(TextFields.IsExists(layerAlias, textFieldAlias)) {
		safeSttyPanic(fmt.Sprintf("The text field '%s' under layer '%s' could not be obtained since it does not exist!", textFieldAlias, layerAlias))
	}
}

/*
validatorMenu is a method which allows you to validate that a menu selector with the specified alias exists under the
given layer alias.

:param layerAlias: The alias of the layer the menu belongs to.
:param menuAlias: The alias of the menu to check for.

Example:

	validatorMenu("mainLayer", "fileMenu")
*/
func validatorMenu(layerAlias string, menuAlias string) {
	if !(Selectors.IsExists(layerAlias, menuAlias)) {
		safeSttyPanic(fmt.Sprintf("The menu '%s' under layer '%s' could not be obtained since it does not exist!", menuAlias, layerAlias))
	}
}

/*
safeSttyPanic is a method which allows you to restore original terminal settings and then generate a panic with the
provided message. In addition, the following should be noted:

- This is used to ensure the terminal is left in a usable state when a fatal error occurs.

:param panicMessage: The message to be included in the panic.

Example:

	safeSttyPanic("A fatal error has occurred.")
*/
func safeSttyPanic(panicMessage interface{}) {
	RestoreTerminalSettings()
	panic(panicMessage)
}
