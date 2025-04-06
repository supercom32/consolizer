package consolizer

import (
	"fmt"
	"github.com/google/uuid"
	"supercom32.net/consolizer/constants"
	"supercom32.net/consolizer/types"
)

type LayerInstanceType struct {
	layerAlias  string
	parentAlias string
	LayerWidth  int
	LayerHeight int
}

func getUUID() string {
	id := uuid.New()
	return id.String()
}

func (shared LayerInstanceType) Clear() {
	layerEntry := Screen.GetLayer(shared.layerAlias)
	localAttributeEntry := types.NewAttributeEntry()
	fillArea(layerEntry, localAttributeEntry, "", 0, 0, shared.LayerWidth, shared.LayerHeight, 0)
}

func (shared LayerInstanceType) DrawImage(fileName string, drawingStyle types.ImageStyleEntryType, xLocation int, yLocation int, widthInCharacters int, heightInCharacters int, blurSigma float64) error {
	var err error
	if !IsImageExists(fileName) {
		err = LoadImage(fileName)
		if err != nil {
			return err
		}
		defer func() {
			UnloadImage(fileName)
		}()
	}
	DrawImageToLayer(shared.layerAlias, fileName, drawingStyle, xLocation, yLocation, widthInCharacters, heightInCharacters, blurSigma)
	return err
}

func (shared LayerInstanceType) DrawComposedImage(imageComposeEntry ImageComposerEntryType, drawingStyle types.ImageStyleEntryType, xLocation int, yLocation int, widthInCharacters int, heightInCharacters int) error {
	var err error
	var imageLayer types.LayerEntryType
	baseImage := imageComposeEntry.RenderImage()
	if drawingStyle.DrawingStyle == constants.ImageStyleHighColor {
		imageLayer = getImageLayerAsHighColor(baseImage, drawingStyle, widthInCharacters, heightInCharacters, drawingStyle.BlurSigmaIntensity)
	} else {
		imageLayer = getImageLayerAsBraille(baseImage, drawingStyle, widthInCharacters, heightInCharacters, drawingStyle.BlurSigmaIntensity)
	}
	drawImageToLayer(shared.layerAlias, imageLayer, xLocation, yLocation)
	return err
}

func (shared LayerInstanceType) AddButton(buttonLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isEnabled bool) ButtonInstanceType {
	buttonAlias := getUUID()
	buttonInstance := Button.Add(shared.layerAlias, buttonAlias, buttonLabel, styleEntry, xLocation, yLocation, width, height, isEnabled)
	return buttonInstance
}

func (shared LayerInstanceType) AddCheckbox(checkboxLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, isSelected bool, isEnabled bool) CheckboxInstanceType {
	checkboxAlias := getUUID()
	checkboxInstance := Checkbox.Add(shared.layerAlias, checkboxAlias, checkboxLabel, styleEntry, xLocation, yLocation, isSelected, isEnabled)
	return checkboxInstance
}

func (shared LayerInstanceType) AddDropdown(styleEntry types.TuiStyleEntryType, selectionEntry types.SelectionEntryType, xLocation int, yLocation int, selectorHeight int, itemWidth int, defaultItemSelected int) DropdownInstanceType {
	dropdownAlias := getUUID()
	dropdownInstance := Dropdown.Add(shared.layerAlias, dropdownAlias, styleEntry, selectionEntry, xLocation, yLocation, selectorHeight, itemWidth, defaultItemSelected)
	return dropdownInstance
}

func (shared LayerInstanceType) AddLabel(labelValue string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int) LabelInstanceType {
	labelAlias := getUUID()
	labelInstance := Label.Add(shared.layerAlias, labelAlias, labelValue, styleEntry, xLocation, yLocation, width)
	return labelInstance
}

func (shared LayerInstanceType) AddProgressBar(progressBarLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, value int, maxValue int, isBackgroundTransparent bool) ProgressBarInstanceType {
	progressBarAlias := getUUID()
	progressBarInstance := ProgressBar.Add(shared.layerAlias, progressBarAlias, progressBarLabel, styleEntry, xLocation, yLocation, width, height, value, maxValue, isBackgroundTransparent)
	return progressBarInstance
}

func (shared LayerInstanceType) AddRadioButton(radioButtonLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, groupId int, isSelected bool) RadioButtonInstanceType {
	radioButtonAlias := getUUID()
	radioButtonInstance := radioButton.Add(shared.layerAlias, radioButtonAlias, radioButtonLabel, styleEntry, xLocation, yLocation, groupId, isSelected)
	return radioButtonInstance
}

func (shared LayerInstanceType) AddScrollbar(styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, length int, maxScrollValue int, scrollValue int, scrollIncrement int, isHorizontal bool) ScrollbarInstanceType {
	scrollbarAlias := getUUID()
	scrollbarInstance := scrollbar.Add(shared.layerAlias, scrollbarAlias, styleEntry, xLocation, yLocation, length, maxScrollValue, scrollValue, scrollIncrement, isHorizontal)
	return scrollbarInstance
}

func (shared LayerInstanceType) AddSelector(styleEntry types.TuiStyleEntryType, selectionEntry types.SelectionEntryType, xLocation int, yLocation int, selectorHeight int, itemWidth int, numberOfColumns int, viewportPosition int, selectedItem int, isBorderDrawn bool) selectorInstanceType {
	selectorAlias := getUUID()
	selectorInstance := Selector.Add(shared.layerAlias, selectorAlias, styleEntry, selectionEntry, xLocation, yLocation, selectorHeight, itemWidth, numberOfColumns, viewportPosition, selectedItem, isBorderDrawn)
	return selectorInstance
}

func (shared LayerInstanceType) AddTextField(styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, maxLengthAllowed int, isPasswordProtected bool, defaultValue string, isEnabled bool) textFieldInstanceType {
	textFieldAlias := getUUID()
	textFieldInstance := TextField.Add(shared.layerAlias, textFieldAlias, styleEntry, xLocation, yLocation, width, maxLengthAllowed, isPasswordProtected, defaultValue, isEnabled)
	return textFieldInstance
}

func (shared LayerInstanceType) AddTextbox(styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isBorderDrawn bool) TextboxInstanceType {
	textBoxAlias := getUUID()
	textBoxInstance := textbox.AddTextbox(shared.layerAlias, textBoxAlias, styleEntry, xLocation, yLocation, width, height, isBorderDrawn)
	return textBoxInstance
}

func (shared LayerInstanceType) AddTooltip(tooltipValue string, styleEntry types.TuiStyleEntryType, hotspotXLocation int, hotspotYLocation int, hotspotWidth int, hotspotHeight int, tooltipXLocation int, tooltipYLocation int, tooltipWidth int, tooltipHeight int, isLocationAbsolute bool, isBorderDrawn bool, hoverTime int) TooltipInstanceType {
	tooltipAlias := getUUID()
	tooltipInstance := Tooltip.Add(shared.layerAlias, tooltipAlias, tooltipValue, styleEntry, hotspotXLocation, hotspotYLocation, hotspotWidth, hotspotHeight, tooltipXLocation, tooltipYLocation, tooltipWidth, tooltipHeight, isLocationAbsolute, isBorderDrawn, hoverTime)
	return tooltipInstance
}

/*
DrawVerticalLine allows you to draw a vertical line on a text layer. This
method also has the ability to draw connectors in case the line intersects
with other lines that have already been drawn. In addition, the following
information should be noted:

- If the the line to be drawn falls outside the area of the text layer
specified, then only the visible portion of the line will be drawn.
*/
func (shared LayerInstanceType) DrawVerticalLine(styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, height int, isConnectorsDrawn bool) {
	layerEntry := Screen.GetLayer(shared.layerAlias)
	localAttributeEntry := types.NewAttributeEntry()
	drawVerticalLine(layerEntry, styleEntry, localAttributeEntry, xLocation, yLocation, height, isConnectorsDrawn)
}

/*
DrawHorizontalLine allows you to draw a horizontal line on a text layer. This
method also has the ability to draw connectors in case the line intersects
with other lines that have already been drawn. In addition, the following
information should be noted:

- If the the line to be drawn falls outside the area of the text layer
specified, then only the visible portion of the line will be drawn.
*/
func (shared LayerInstanceType) DrawHorizontalLine(styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, isConnectorsDrawn bool) {
	layerEntry := Screen.GetLayer(shared.layerAlias)
	localAttributeEntry := types.NewAttributeEntry()
	drawHorizontalLine(layerEntry, styleEntry, localAttributeEntry, xLocation, yLocation, width, isConnectorsDrawn)
}

/*
DrawBorder allows you to draw a border on a given text layer. Borders differ
from frames since they are flat shaded and do not have a raised or sunken
look to them. In addition, the following information should be noted:

- If the border to be drawn falls outside the range of the specified layer,
then only the visible portion of the border will be drawn.

- The 'isInteractive' option allows you to specify if the window should
interact with the layer being drawn on. For example, when enabled if the user
drags the window title bar, the whole layer will move to simulate movement of
the window itself.
*/
func (shared LayerInstanceType) DrawBorder(styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isInteractive bool) {
	layerEntry := Screen.GetLayer(shared.layerAlias)
	localAttributeEntry := types.NewAttributeEntry()
	drawBorder(layerEntry, styleEntry, localAttributeEntry, xLocation, yLocation, width, height, isInteractive)
}

/*
DrawFrameLabel allows you to draw a label for a frame. The label will
be automatically enclosed by the characters "[" and "]" to blend in
with a border of a frame.

- If the frame label to be drawn falls outside the range of the
specified layer, then only the visible portion of the border will be
drawn.
*/
func (shared LayerInstanceType) DrawFrameLabel(styleEntry types.TuiStyleEntryType, label string, xLocation int, yLocation int) {
	layerEntry := Screen.GetLayer(shared.layerAlias)
	drawFrameLabel(layerEntry, styleEntry, label, xLocation, yLocation)
}

/*
DrawFrame allows you to draw a frame on a given text layer. Frames differ
from borders since borders are flat shaded and do not have a raised or
sunken look to them. In addition, the following information should be noted:

- If the frame to be drawn falls outside the range of the specified layer,
then only the visible portion of the frame will be drawn.

- The 'isInteractive' option allows you to specify if the window should
interact with the layer being drawn on. For example, when enabled if the user
drags the window title bar, the whole layer will move to simulate movement of
the window itself.
*/
func (shared LayerInstanceType) DrawFrame(styleEntry types.TuiStyleEntryType, isRaised bool, xLocation int, yLocation int, width int, height int, isInteractive bool) {
	layerEntry := Screen.GetLayer(shared.layerAlias)
	localAttributeEntry := types.NewAttributeEntry()
	if isRaised {
		drawFrame(layerEntry, styleEntry, localAttributeEntry, constants.FrameStyleRaised, xLocation, yLocation, width, height, isInteractive)
	} else {
		drawFrame(layerEntry, styleEntry, localAttributeEntry, constants.FrameStyleSunken, xLocation, yLocation, width, height, isInteractive)
	}
}

/*
DrawWindow allows you to draw a window on a given text layer. Windows differ
from borders since the entire area the window surrounds gets filled with
a solid background color. In addition, the following information should be noted:

- If the window to be drawn falls outside the range of the specified layer,
then only the visible portion of the window will be drawn.

- The 'isInteractive' option allows you to specify if the window should
interact with the layer being drawn on. For example, when enabled if the user
drags the window title bar, the whole layer will move to simulate movement of
the window itself.
*/
func (shared LayerInstanceType) DrawWindow(styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isInteractive bool) {
	layerEntry := Screen.GetLayer(shared.layerAlias)
	localAttributeEntry := types.NewAttributeEntry()
	drawWindow(layerEntry, styleEntry, localAttributeEntry, xLocation, yLocation, width, height, isInteractive)
}

/*
DrawShadow allows you to draw shadows on a given text layer. Shadows are simply
transparent areas which darken whatever text layers are underneath it by a
specified degree. In addition, the following information should be noted:

- The alpha value can range from 0.0 (no shadow) to 1.0 (totally black).
*/
func (shared LayerInstanceType) DrawShadow(xLocation int, yLocation int, width int, height int, alphaValue float32) {
	layerEntry := Screen.GetLayer(shared.layerAlias)
	localAttributeEntry := types.NewAttributeEntry()
	drawShadow(layerEntry, localAttributeEntry, xLocation, yLocation, width, height, alphaValue)
}

/*
FillArea allows you to fill an area of a given text layer with characters of
your choice. If you wish to fill the area with repeating text, simply provide
the string you wish to repeat. In addition, the following information should be
noted:

- If the area to fill falls outside the range of the specified layer, then only
the visible portion of the fill will be drawn.
*/
func (shared LayerInstanceType) FillArea(fillCharacters string, xLocation int, yLocation int, width int, height int) {
	layerEntry := Screen.GetLayer(shared.layerAlias)
	attributeEntry := layerEntry.DefaultAttribute
	fillArea(layerEntry, attributeEntry, fillCharacters, xLocation, yLocation, width, height, constants.NullCellControlLocation)
}

/*
FillLayer allows you to fill an entire layer with characters of your choice.
If you wish to fill the layer with repeating text, simply provide the string
you wish to repeat.
*/
func (shared LayerInstanceType) FillLayer(fillCharacters string) {
	layerEntry := Screen.GetLayer(shared.layerAlias)
	attributeEntry := layerEntry.DefaultAttribute
	fillLayer(layerEntry, attributeEntry, fillCharacters)
}

/*
DrawBar allows you to draw a horizontal bar on a given text layer row. This is
useful for drawing application headers or status bar footers.
*/
func (shared LayerInstanceType) DrawBar(xLocation int, yLocation int, barLength int, fillCharacters string) {
	layerEntry := Screen.GetLayer(shared.layerAlias)
	attributeEntry := layerEntry.DefaultAttribute
	fillArea(layerEntry, attributeEntry, fillCharacters, xLocation, yLocation, barLength, 1, constants.NullCellControlLocation)
}

/*
MoveLayerByAbsoluteValue allows you to move a text layer by an absolute value.
This is useful if you know exactly what position you wish to move your text
layer to. In addition, the following information should be noted:

- If you move your layer outside the visible terminal display, only the visible
display area will be rendered. Likewise, if your text layer is a child of
a parent layer, then only the visible display area will be rendered on the
parent.
*/
func (shared LayerInstanceType) MoveLayerByAbsoluteValue(xLocation int, yLocation int) {
	validateLayer(shared.layerAlias)
	layerEntry := Screen.GetLayer(shared.layerAlias)
	layerEntry.ScreenXLocation = xLocation
	layerEntry.ScreenYLocation = yLocation
}

/*
MoveLayerByRelativeValue allows you to move a text layer by a relative value.
This is useful for windows, foregrounds, backgrounds, or any kind of
animations or movement you may wish to do in increments. For example:

	// Move the text layer with the alias "ForegroundLayer" one character to
	// the left and two characters down from its current location.
	consolizer.MoveLayerByRelativeValue("ForegroundLayer", -1, 2)

In addition, the following information should be noted:

- If you move your layer outside the visible terminal display, only the visible
display area will be rendered. Likewise, if your text layer is a child of
a parent layer, then only the visible display area will be rendered on the
parent.
*/
func (shared LayerInstanceType) MoveLayerByRelativeValue(xLocation int, yLocation int) {
	validateLayer(shared.layerAlias)
	layerEntry := Screen.GetLayer(shared.layerAlias)
	layerEntry.ScreenXLocation += xLocation
	layerEntry.ScreenYLocation += yLocation
}

/*
DeleteLayer allows you to remove a text layer. If you wish to reuse a text
layer for a future purpose, you may also consider making the layer invisible
instead of deleting it. In addition, the following information should be noted:

- When a text layer is deleted, all child text layers are recursively deleted
as well.

- If any dynamically drawn TUI controls reference the deleted layer, they will
still be present. However, because the layer they were created for no longer
exists, they will never be rendered. Consider removing any TUI controls before
deleting the layer they reference. If you delete a layer that is referenced
by dynamic TUI controls, creating a new layer with the same layer alias will
allow them to be rendered again.

- If you attempt to delete a text layer which is currently set as your default
text layer, then a panic will be generated in order to fail as fast as
possible.

- If you attempt to delete a text layer that does not exist, then the operation
will be ignored.
*/
func (shared LayerInstanceType) DeleteLayer() {
	validateLayer(shared.layerAlias)
	Screen.DeleteLayer(shared.layerAlias)
	if commonResource.layerInstance.layerAlias == shared.layerAlias {
		nextLayerAlias := Screen.GetNextLayerAlias()
		nextLayerInstance := Screen.GetLayer(nextLayerAlias)
		commonResource.layerInstance = LayerInstanceType{layerAlias: nextLayerAlias, parentAlias: nextLayerInstance.ParentAlias, LayerWidth: nextLayerInstance.Width, LayerHeight: nextLayerInstance.Height}
	}
	shared.layerAlias = ""
}

func (shared LayerInstanceType) IsLayerExists() bool {
	if shared.layerAlias != "" {
		return true
	}
	return false
}

func (shared LayerInstanceType) SetIsVisible(isVisible bool) {
	validateLayer(shared.layerAlias)
	setLayerIsVisible(shared.layerAlias, isVisible)
}

/*
Color24Bit allows you to color a layer using a 24-bit color expressed as
an int32. This is useful for when you have colors which are already defined.
*/

func (shared LayerInstanceType) Color24Bit(foregroundColor constants.ColorType, backgroundColor constants.ColorType) {
	ColorLayer24Bit(shared, foregroundColor, backgroundColor)
}

/*
Locate allows you to set the default cursor location on your specified text
layer for printing with. This is useful for when you wish to print text
at different locations of your text layer at any given time. If you wish to
change the cursor location for a text layer that is not currently set as your
default, use 'LocateLayer' instead. In addition, the following information
should be noted:

- If you pass in a location value that falls outside the dimensions of the
default text layer, a panic will be generated to fail as fast as possible.

- Valid text layer locations start at position (0,0) for the upper left corner.
Since location values do not start at (1,1), valid end positions for the bottom
right corner will be one less than the text layer width and height. For
example:

	// Create a new text layer with the alias "ForegroundLayer", at location
	// (0,0), with a width and height of 15x15, a z order priority of 1,
	// and no parent layer associated with it.
	consolizer.AddLayer("ForegroundLayer", 0, 0, 15, 15, 1, "")
	// Set the text layer with the alias "ForegroundLayer" as our default.
	consolizer.Layer("ForegroundLayer")
	// Move our cursor location to the bottom right corner of our text layer.
	consolizer.Locate(14, 14)
*/
func (shared LayerInstanceType) Locate(xLocation int, yLocation int) {
	validateDefaultLayerIsNotEmpty()
	LocateLayer(shared, xLocation, yLocation)
}

/*
Print allows you to write text to the default text layer. If you wish to
print to a text layer that is not currently set as the default, use
'PrintLayer' instead. In addition, the following information should be noted:

- When text is written to the text layer, the cursor position is also updated
to reflect its new location. Like a typewriter, the cursor position moves to
the start of the next line after each print statement.

- If the string to print ends up being too long to fit at its current location,
then only the visible portion of your string will be printed.

- If printing has not yet finished and there are no available lines left, then
all remaining characters will be discarded and printing will stop.
*/
func (shared LayerInstanceType) Print(textToPrint string) {
	validateDefaultLayerIsNotEmpty()
	PrintLayer(shared, textToPrint)
}

/*
PrintDialog allows you to write text immediately to the terminal screen via a
typewriter effect. This is useful for video games or other applications that
may require printing text in a dialog box. In addition, the following
information should be noted:

- If you specify a print location outside the range of your specified text
layer, a panic will be generated to fail as fast as possible.

- If printing has reached the last line of your text layer, printing will
not advance to the next line. Instead, it will resume and overwrite
what was already printed on the same line.

- Specifying the width of your text line allows you to control when text
wrapping occurs. For example, if printing starts at location (2, 2) and you set
a line width of 10 characters, text wrapping will occur when the printing
exceeds the text layer location (12, 2). When this happens, text will continue
to print underneath the previous line at (2, 3).

- When a word is too long to be printed on a text layer line, or the width
of your line has already exceed its allowed maximum, the word will be pushed
to the line directly under it. This prevents words from being split across
two lines.

- When specifying a printing delay, the amount of time to wait is inserted
between each character printed and does not reflect the overall time to
print your specified text.

- If the dialog being printed is flagged as skipable, the user can speed up
printing by pressing the 'enter' key or right mouse button. Otherwise, they
must wait for the animation to completely finish before execution continues.

- This method supports the use of text styles during printing to add color
or styles to specific words in your string. All text styles must be enclosed
around the "{" and "}" characters. If you wish to use the default text
style, simply omit specifying any text style between your enclosing braces.
For example:

	// AddLayer a text layer with the alias "ForegroundLayer", at location (0, 0),
	// with a width and height of 80x20 characters, z order priority of 1,
	// with no parent layer.
	dosktop.AddLayer("ForegroundLayer", 0, 0, 80, 20, 1, "")
	// Obtain a new text style entry.
	redTextStyle := dosktop.GetTextStyle()
	// Change the default foreground color of our text style to be red.
	redTextStyle.ForegroundColor = dosktop.GetRGBColor(255,0,0)
	// Register our new text style so Dosktop can use it.
	dosktop.AddTextStyle("red", redTextStyle)
	// Print some dialog text on the text layer "ForegroundLayer", at location
	// (0, 0), with a text wrapping location at 30 characters, a 10 millisecond
	// delay between each character printed, and mark the dialog as skipable.
	// Inside our string to print, we add the "{red}" tag to switch printing
	// styles on the fly to "red" and change back to the default style using
	// "{}".
	dosktop.PrintDialog("ForegroundLayer", 0, 0, 30, 10, true, "This is some dialog text in {red}red color{}. Only the words 'red color' should be colored.")
*/
func (shared LayerInstanceType) PrintDialog(xLocation int, yLocation int, widthOfLineInCharacters int, printDelayInMilliseconds int, isSkipable bool, stringToPrint string) {
	layerEntry := Screen.GetLayer(shared.layerAlias)
	if xLocation < 0 || xLocation > layerEntry.Width || yLocation < 0 || yLocation > layerEntry.Height {
		panic(fmt.Sprintf("The specified location (%d, %d) is out of bounds for layer '%s' with a size of (%d, %d).", xLocation, yLocation, layerEntry.LayerAlias, layerEntry.Width, layerEntry.Height))
	}
	printDialog(layerEntry, layerEntry.DefaultAttribute, xLocation, yLocation, widthOfLineInCharacters, printDelayInMilliseconds, isSkipable, stringToPrint)
}

/*
AddLayer allows you to add a text layer to the current terminal display. You
can add as many layers as you wish to suite your applications needs. Text
layers are useful for setting up windows, modal dialogs, viewports, game
foregrounds and backgrounds, and even effects like parallax scrolling. In
addition, the following information should be noted:

- If you specify location for your layer that is outside the visible
terminal display, then only the visible portion of your text layer will be
rendered. Likewise, if your text layer is larger than the visible area of your
terminal display, then only the visible portion of it will be displayed.

- If you pass in a zero or negative value for ether width or height a panic
will be generated to fail as fast as possible.

- The z order priority controls which text layer should be drawn first and
which text layer should be drawn last. Layers that have a higher priority
will be drawn on top of layers that have a lower priority. In the event
that two layers have the same priority, they will be drawn in random order.
This is to ensure that programmers do not attempt to rely on any specific
behavior that might be a coincidental side effect.

- The parent alias specifies which text layer is the parent of the one being
created. Having a parent layer means that the child layer will only render
on the parent and not the main terminal. This allows you to have text layers
within text layers that can be moved or manipulated relative to the parent.
If you pass in a value of "" for the parent alias, then no parent is used
and the layer is rendered directly to the terminal display. This feature
is useful for creating 'Window' effects where content is contained within
something else.

- When adding a new text layer, it will become the default
working text layer automatically. If you wish to set another text layer
as your default, use 'Layer' to explicitly set it.
*/
func AddLayer(xLocation int, yLocation int, width int, height int, zOrderPriority int, parentLayerInstance *LayerInstanceType) LayerInstanceType {
	layerAlias := getUUID()
	validateTerminalWidthAndHeight(width, height)
	if parentLayerInstance == nil {
		Screen.AddLayer(layerAlias, xLocation, yLocation, width, height, zOrderPriority, "")
		layerInstance := LayerInstanceType{layerAlias: layerAlias, parentAlias: "", LayerWidth: width, LayerHeight: height}
		commonResource.layerInstance = layerInstance
		return layerInstance
	} else {
		Screen.AddLayer(layerAlias, xLocation, yLocation, width, height, zOrderPriority, parentLayerInstance.layerAlias)
		layerInstance := LayerInstanceType{layerAlias: layerAlias, parentAlias: "", LayerWidth: width, LayerHeight: height}
		commonResource.layerInstance = layerInstance
		return layerInstance
	}
}

/*
MoveLayerByAbsoluteValue allows you to move a text layer by an absolute value.
This is useful if you know exactly what position you wish to move your text
layer to. In addition, the following information should be noted:

- If you move your layer outside the visible terminal display, only the visible
display area will be rendered. Likewise, if your text layer is a child of
a parent layer, then only the visible display area will be rendered on the
parent.
*/
func MoveLayerByAbsoluteValue(layerAlias string, xLocation int, yLocation int) {
	validateLayer(layerAlias)
	layerEntry := Screen.GetLayer(layerAlias)
	layerEntry.ScreenXLocation = xLocation
	layerEntry.ScreenYLocation = yLocation
}

/*
MoveLayerByRelativeValue allows you to move a text layer by a relative value.
This is useful for windows, foregrounds, backgrounds, or any kind of
animations or movement you may wish to do in increments. For example:

	// Move the text layer with the alias "ForegroundLayer" one character to
	// the left and two characters down from its current location.
	consolizer.MoveLayerByRelativeValue("ForegroundLayer", -1, 2)

In addition, the following information should be noted:

- If you move your layer outside the visible terminal display, only the visible
display area will be rendered. Likewise, if your text layer is a child of
a parent layer, then only the visible display area will be rendered on the
parent.
*/
func MoveLayerByRelativeValue(layerAlias string, xLocation int, yLocation int) {
	validateLayer(layerAlias)
	layerEntry := Screen.GetLayer(layerAlias)
	layerEntry.ScreenXLocation += xLocation
	layerEntry.ScreenYLocation += yLocation
}

/*
DeleteLayer allows you to remove a text layer. If you wish to reuse a text
layer for a future purpose, you may also consider making the layer invisible
instead of deleting it. In addition, the following information should be noted:

- When a text layer is deleted, all child text layers are recursively deleted
as well.

- If any dynamically drawn TUI controls reference the deleted layer, they will
still be present. However, because the layer they were created for no longer
exists, they will never be rendered. Consider removing any TUI controls before
deleting the layer they reference. If you delete a layer that is referenced
by dynamic TUI controls, creating a new layer with the same layer alias will
allow them to be rendered again.

- If you attempt to delete a text layer which is currently set as your default
text layer, then a panic will be generated in order to fail as fast as
possible.

- If you attempt to delete a text layer that does not exist, then the operation
will be ignored.
*/
func deleteLayer(layerAlias string) {
	validateLayer(layerAlias)
	Screen.DeleteLayer(layerAlias)
	if commonResource.layerInstance.layerAlias == layerAlias {
		nextLayerAlias := Screen.GetNextLayerAlias()
		// If last entry and no more layers, just return. Do not set anything.
		if nextLayerAlias == "" {
			commonResource.layerInstance = LayerInstanceType{layerAlias: "", parentAlias: "", LayerWidth: 0, LayerHeight: 0}
			return
		}
		nextLayerInstance := Screen.GetLayer(nextLayerAlias)
		commonResource.layerInstance = LayerInstanceType{layerAlias: nextLayerAlias, parentAlias: nextLayerInstance.ParentAlias, LayerWidth: nextLayerInstance.Width, LayerHeight: nextLayerInstance.Height}
	}
}

func DeleteLayer(layerInstance LayerInstanceType) {
	Screen.DeleteLayer(layerInstance.layerAlias)
	if commonResource.layerInstance.layerAlias == layerInstance.layerAlias {
		nextLayerAlias := Screen.GetNextLayerAlias()
		nextLayerInstance := Screen.GetLayer(nextLayerAlias)
		commonResource.layerInstance = LayerInstanceType{layerAlias: nextLayerAlias, parentAlias: nextLayerInstance.ParentAlias, LayerWidth: nextLayerInstance.Width, LayerHeight: nextLayerInstance.Height}
	}
}

/*
DeleteAllLayers allows you to remove all layers from memory.
*/
func DeleteAllLayers() {
	for _, entryToRemove := range Screen.ScreenMemory.GetAllEntries() {
		deleteLayer(entryToRemove.LayerAlias)
	}
	Screen.ReInitializeScreenMemory()
}

func isLayerExists(layerAlias string) bool {
	if Screen.IsLayerExists(layerAlias) {
		return true
	}
	return false
}

func setLayerIsVisible(layerAlias string, isVisible bool) {
	validateLayer(layerAlias)
	layerEntry := Screen.GetLayer(layerAlias)
	layerEntry.IsVisible = isVisible
}
