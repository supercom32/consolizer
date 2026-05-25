package consolizer

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/stringformat"
	"github.com/supercom32/consolizer/types"
	"os"
	"os/signal"
	"math/rand"
	"runtime"
	"sync"
	"syscall"
	"time"
)

const (
	OS_WINDOWS = 1
	OS_LINUX   = 2
	OS_MAC     = 3
	OS_OTHER   = 4
)

/*
defaultValueType is a structure which holds common information about the current terminal session that needs to be
shared.
*/
type defaultValueType struct {
	screen               tcell.Screen
	terminalWidth        int
	terminalHeight       int
	screenLayer          types.LayerEntryType
	debugDirectory       string
	isDebugEnabled       bool
	displayUpdate        sync.Mutex
	updateDisplayChannel chan bool
}

/*
commonResource is a variable which holds shared data that is accessed by this package.

Example:

	commonResource
*/
var commonResource defaultValueType

/*
GetVersion is a method which allows you to obtain the current version of the consolizer library.

Example:

	version := GetVersion()
*/
func GetVersion() string {
	return "1"
}

/*
InitializeTerminal is a method which allows you to initialize consolizer for the first time. This method must be called
first before any operations take place. The parameters width and height represent the display size of the terminal
instance you wish to create. In addition, the following should be noted:

  - If you pass in a zero or negative value for either width or height a panic will be generated to fail as fast as
    possible.

Example:

	InitializeTerminal(80, 25)
*/
func InitializeTerminal(width int, height int) {
	InitializeTimerMemory()
	// Set the mouse location off screen so it won't trigger events at 0,0 which the user never moved to.
	SetMouseStatus(-1, -1, 0, "")
	var detectedWidth int
	var detectedHeight int
	if !commonResource.isDebugEnabled {
		screen, err := tcell.NewScreen()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		if err := screen.Init(); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		commonResource.screen = screen
		commonResource.screen.EnableMouse()
		commonResource.updateDisplayChannel = make(chan bool)
		setupCloseHandler()
		detectedWidth, detectedHeight = GetTerminalSize()
	}
	if width == 0 {
		commonResource.terminalWidth = detectedWidth
	} else {
		commonResource.terminalWidth = width
	}
	if height == 0 {
		commonResource.terminalHeight = detectedHeight
	} else {
		commonResource.terminalHeight = height
	}
	commonResource.debugDirectory = "/tmp/"
	validateTerminalWidthAndHeight(commonResource.terminalWidth, commonResource.terminalHeight)
	DeleteAllLayers()
	go setupEventUpdater()
	go setupPeriodicEventUpdater()
}

/*
setupPeriodicEventUpdater is a method which is a background method that updates periodic events.

Example:

	go setupPeriodicEventUpdater()
*/
func setupPeriodicEventUpdater() {
	for {
		UpdatePeriodicEvents()
		time.Sleep(10 * time.Millisecond)
	}
}

/*
setupEventUpdater is a method which is a background method that monitors all events coming into the terminal session.
When an event is detected, it is recorded and monitoring continues.

Example:

	go setupEventUpdater()
*/
func setupEventUpdater() {
	for {
		select {
		case <-commonResource.updateDisplayChannel:
			return
		default:
			UpdateEventQueues()
		}
	}
}

/*
setupCloseHandler is a method which enables the trapping of all unexpected system calls and shuts down the terminal
gracefully. This means all terminal settings should be reset back to normal if anything unexpected happens to the user
or if the process is killed.

Example:

	setupCloseHandler()
*/
func setupCloseHandler() {
	channel := make(chan os.Signal)
	signal.Notify(channel, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)
	go func() {
		<-channel
		commonResource.screen.Fini()
		os.Exit(1)
	}()

}

/*
RestoreTerminalSettings is a method which allows the user to gracefully return the terminal back to its normal settings.
This should be called once your application is finished using consolizer so that the users terminal environment is not
left in a bad state.

Example:

	RestoreTerminalSettings()
*/
func RestoreTerminalSettings() {
	commonResource.updateDisplayChannel <- true
	DeleteAllLayers()
	if commonResource.screen == nil {
		return
	}
	commonResource.screen.DisableMouse()
	commonResource.screen.Clear()
	commonResource.screen.Sync()
	commonResource.screen.Suspend()
	commonResource.screen.Fini()
	if getOsType() == OS_WINDOWS {
		fmt.Println("")
	}
}

/*
getOsType is a method which allows you to obtain the type of operating system currently running.

Example:

	osType := getOsType()
*/
func getOsType() int {
	os := runtime.GOOS
	switch os {
	case "windows":
		return OS_WINDOWS
	case "darwin":
		return OS_MAC
	case "linux":
		return OS_LINUX
	default:
		return OS_WINDOWS
	}
	return 0
}

/*
GetTerminalSize is a method which allows you to obtain width and height of the current terminal characters.

Example:

	width, height := GetTerminalSize()
*/
func GetTerminalSize() (int, int) {
	return commonResource.screen.Size()
}

/*
Inkey is a method which allows you to read keyboard input from the user's terminal. This method returns the character
pressed or a keyword representing the special key pressed (for example: 'a', 'A', 'escape', 'f10', etc.). In addition,
the following should be noted:

- If more than one keystroke is recorded, it is stored sequentially in the input buffer and this method needs to be
called repeatedly in order to read them.

Example:

	keystroke := Inkey()
*/
func Inkey() []rune {
	return KeyboardMemory.GetFromBuffer()
}

/*
AddTextStyle is a method which allows you to add a new printing style which can be used for writing dialog text with.
Dialog text printing differs from regular printing, since it allows for typewriter drawing effects as well as changing
text attributes on the fly (color, bold, etc). In addition, the following should be noted:

- This method expects you to pass in a textStyleEntry type obtained by calling consolizer.NewTextStyle.

Example:

	AddTextStyle("MyStyle", myStyleEntry)
*/
func AddTextStyle(textStyleAlias string, textStyleEntry types.TextCellStyleEntryType) {
	TextStyles.Add(textStyleAlias, &textStyleEntry)
}

/*
DeleteTextStyle is a method which allows you to remove a text style that was added previously. In addition, the
following should be noted:

- If you attempt to delete an entry that does not exist, then no operation will be performed.

Example:

	DeleteTextStyle("MyStyle")
*/
func DeleteTextStyle(textStyleAlias string) {
	validateTextStyleExists(textStyleAlias)
	TextStyles.Remove(textStyleAlias)
}

/*
NewTextStyle is a constructor which allows you to obtain a new text style entry which can be used when printing dialog
text. By configuring attributes for your text style entry and adding your entry to consolizer, simple markup commands
can be used to switch between dialog printing styles automatically.

Example:

	myTextStyleEntry := consolizer.NewTextStyle()
*/
func NewTextStyle() types.TextCellStyleEntryType {
	return types.NewTextCellStyleEntry()
}

/*
NewImageStyle is a constructor which allows you to obtain a new image style entry.

Example:

	myImageStyle := NewImageStyle()
*/
func NewImageStyle() types.ImageStyleEntryType {
	return types.NewImageStyleEntry()
}

/*
NewTuiStyleEntry is a constructor which allows you to obtain a new style entry which can be used for specifying how TUI
controls and other TUI drawing operations should occur.

Example:

	myTuiStyleEntry := consolizer.NewTuiStyleEntry()
*/
func NewTuiStyleEntry() types.TuiStyleEntryType {
	return types.NewTuiStyleEntry()
}

/*
NewSelectionEntry is a constructor which allows you to obtain an entry used for specifying what options you want to make
available for a given menu prompt.

Example:

	selectionEntry := consolizer.NewSelectionEntry()
*/
func NewSelectionEntry() types.SelectionEntryType {
	return types.NewSelectionEntry()
}

/*
NewAssetList is a constructor which allows you to obtain a list for storing asset information with. This is useful for
loading assets in bulk since you can specify asset information as a collection instead of each one individually. In
addition, the following should be noted:

- An asset list can contain multiple asset types.

Example:

	assetList := consolizer.NewAssetList()
*/
func NewAssetList() types.AssetListType {
	return types.NewAssetList()
}

/*
setLayerZOrderInstance is a method which allows you to set the z-order value for a given layer. In addition, the following
should be noted:

- The z-order priority controls which text layer should be drawn first and which text layer should be drawn last.

- Layers that have a higher priority will be drawn on top of layers that have a lower priority.

- In the event that two layers have the same priority, they will be drawn in random order.

Example:

	setLayerZOrderInstance(myLayer, 10)
*/
func setLayerZOrderInstance(layerInstance *LayerInstanceType, zOrder int) {
	layerEntry := Layers.Get(layerInstance.layerAlias)
	layerEntry.ZOrder = zOrder
}

/*
setLayerAlphaInstance is a method which allows you to set the alpha value for a given text layer. This lets you perform pseudo
transparencies by making the layer foreground and background colors blend with the layers underneath it to the degree
specified. In addition, the following should be noted:

- An alpha value of 1.0 is equal to 100% visible, while an alpha value of 0.0 is 0% visible.

- If the percent change specified is outside of the RGB color range, then the color will simply bottom or max out.

Example:

	setLayerAlphaInstance(myLayer, 0.5)
*/
func setLayerAlphaInstance(layerInstance *LayerInstanceType, alphaValue float32) {
	layerEntry := Layers.Get(layerInstance.layerAlias)
	layerEntry.DefaultAttribute.ForegroundAlphaValue = alphaValue
	layerEntry.DefaultAttribute.BackgroundAlphaValue = alphaValue
}

/*
GetColor is a method which allows you to obtain an RGB color from a predefined color palette list. This list corresponds
to the 16 color ANSI standard, where color 0 is Black and 15 is Bright White. In addition, the following should be
noted:

- If you specify a color index less than 0 or greater than 15 a panic will be generated to fail as fast as possible.

Example:

	color := GetColor(1)
*/
func GetColor(colorIndex int) constants.ColorType {
	validateColorIndex(colorIndex)
	return constants.AnsiColorByIndex[colorIndex]
}

/*
GetRGBColor is a method which allows you to obtain a specific RGB color based on the red, green, and blue index values
provided. In addition, the following should be noted:

- If you specify a color channel index less than 0 or greater than 255 a panic will be generated to fail as fast as
possible.

Example:

	color := GetRGBColor(255, 128, 0)
*/
func GetRGBColor(redColorIndex int32, greenColorIndex int32, blueColorIndex int32) constants.ColorType {
	validateRGBColorIndex(redColorIndex, greenColorIndex, blueColorIndex)
	return constants.ColorType(tcell.NewRGBColor(redColorIndex, greenColorIndex, blueColorIndex))
}

/*
colorLayerInstance is a method which allows you to set default colors on your specified text layer for printing with. The color
index specified corresponds to the 16 color ANSI standard, where color 0 is Black and 15 is Bright White.

Example:

	colorLayerInstance(myLayer, 15, 0)
*/
func colorLayerInstance(layerInstance *LayerInstanceType, foregroundColorIndex int, backgroundColorIndex int) {
	validateColorIndex(foregroundColorIndex)
	validateColorIndex(backgroundColorIndex)
	layerEntry := Layers.Get(layerInstance.layerAlias)
	layerEntry.DefaultAttribute.ForegroundColor = constants.AnsiColorByIndex[foregroundColorIndex]
	layerEntry.DefaultAttribute.BackgroundColor = constants.AnsiColorByIndex[backgroundColorIndex]
}

/*
colorLayerRGBInstance is a method which allows you to set default colors on your specified text layer for printing with using
RGB values. This method allows you to specify colors using RGB color index values within the range of 0 to 255.

Example:

	colorLayerRGBInstance(myLayer, 255, 255, 255, 0, 0, 0)
*/
func colorLayerRGBInstance(layerInstance *LayerInstanceType, foregroundRed int32, foregroundGreen int32, foregroundBlue int32, backgroundRed int32, backgroundGreen int32, backgroundBlue int32) {
	foregroundColor := GetRGBColor(foregroundRed, foregroundGreen, foregroundBlue)
	backgroundColor := GetRGBColor(backgroundRed, backgroundGreen, backgroundBlue)
	colorLayer24BitInstance(layerInstance, foregroundColor, backgroundColor)
}

/*
colorLayer24BitInstance is a method which allows you to color a specified layer using a 24-bit color expressed as an int32.

Example:

	colorLayer24BitInstance(myLayer, fgColor, bgColor)
*/
func colorLayer24BitInstance(layerInstance *LayerInstanceType, foregroundColor constants.ColorType, backgroundColor constants.ColorType) {
	layerEntry := Layers.Get(layerInstance.layerAlias)
	layerEntry.DefaultAttribute.ForegroundColor = foregroundColor
	layerEntry.DefaultAttribute.BackgroundColor = backgroundColor
}

/*
locateLayerInstance is a method which allows you to set the default cursor location on your specified text layer for printing
with. In addition, the following should be noted:

- If you pass in a location value that falls outside the dimensions of the specified text layer, a panic will be
generated.

Example:

	locateLayerInstance(myLayer, 10, 5)
*/
func locateLayerInstance(layerInstance *LayerInstanceType, xLocation int, yLocation int) {
	validateLayer(layerInstance.layerAlias)
	layerEntry := Layers.Get(layerInstance.layerAlias)
	validateLayerLocationByLayerEntry(layerEntry, xLocation, yLocation)
	layerEntry.CursorXLocation = xLocation
	layerEntry.CursorYLocation = yLocation
}

/*
printLayerInstance is a method which allows you to write text to a specified text layer. In addition, the following should be
noted:

- When text is written to the text layer, the cursor position is also updated to reflect its new location.

  - If the string to print ends up being too long to fit at its current location, then only the visible portion of your
    text will be printed.

Example:

	printLayerInstance(myLayer, "Hello World")
*/
func printLayerInstance(layerInstance *LayerInstanceType, textToPrint string) {
	layerEntry := Layers.Get(layerInstance.layerAlias)
	if layerEntry.CursorYLocation >= layerEntry.Height {
		layerEntry.CursorYLocation = layerEntry.Height - 1
		layerEntry.CharacterMemory = scrollCharacterMemory(layerEntry)
	}
	arrayOfRunes := stringformat.GetRunesFromString(textToPrint)
	printLayer(layerEntry, layerEntry.DefaultAttribute, layerEntry.CursorXLocation, layerEntry.CursorYLocation, arrayOfRunes)
	layerEntry.CursorXLocation = 0
	layerEntry.CursorYLocation = layerEntry.CursorYLocation + 1
}

/*
printLayer is a method which allows you to write text to a text layer directly. This is useful for internal methods that
want to write text to a text layer directly, without affecting user settings. In addition, the following should be
noted:

  - If the location to print falls outside the range of the text layer, then only the visible portion of your text will
    be rendered.

Example:

	printLayer(layerEntry, attr, 0, 0, runes)
*/
func printLayer(layerEntry *types.LayerEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int, textToPrint []rune) int {
	layerWidth := layerEntry.Width
	layerHeight := layerEntry.Height
	cursorXLocation := xLocation
	cursorYLocation := yLocation
	characterMemory := layerEntry.CharacterMemory
	for _, currentCharacter := range textToPrint {
		if cursorXLocation >= 0 && cursorXLocation < layerWidth && cursorYLocation >= 0 && cursorYLocation < layerHeight {
			originalBackgroundColor := characterMemory[cursorYLocation][cursorXLocation].AttributeEntry.BackgroundColor
			characterMemory[cursorYLocation][cursorXLocation].AttributeEntry = types.NewAttributeEntry(&attributeEntry)
			characterMemory[cursorYLocation][cursorXLocation].Character = currentCharacter
			if stringformat.IsRuneCharacterWide(currentCharacter) {
				cursorXLocation++
				if cursorXLocation >= layerWidth {
					return cursorXLocation - xLocation
				}
				characterMemory[cursorYLocation][cursorXLocation].AttributeEntry = types.NewAttributeEntry(&attributeEntry)
				characterMemory[cursorYLocation][cursorXLocation].Character = ' '
			}
			if characterMemory[cursorYLocation][cursorXLocation].AttributeEntry.IsBackgroundTransparent {
				characterMemory[cursorYLocation][cursorXLocation].AttributeEntry.BackgroundColor = originalBackgroundColor
			}
		}
		cursorXLocation++
		if cursorXLocation >= layerWidth {
			return cursorXLocation - xLocation
		}
	}
	return cursorXLocation - xLocation
}

/*
clearLayer is a method which allows you to empty the specified text layer of all its contents.

Example:

	clearLayer(layerEntry)
*/
func clearLayer(layerEntry *types.LayerEntryType) {
	types.InitializeCharacterMemory(layerEntry)
}

/*
GetCharacterOnScreen is a method which allows you to obtain the character currently being displayed on the screen at a
specific location.

Example:

	char := GetCharacterOnScreen(10, 5)
*/
func GetCharacterOnScreen(xLocation int, yLocation int) rune {
	layerEntry := commonResource.screenLayer
	validateLayerLocationByLayerEntry(&layerEntry, xLocation, yLocation)
	return layerEntry.CharacterMemory[xLocation][yLocation].Character
}

/*
scrollCharacterMemory is a method which allows you to advance the specified text layer up by one row. In addition, the
following should be noted:

- The first row is discarded and all subsequent rows are moved up by one position.

- The new row created at the bottom of the text layer will be filled with spaces.

Example:

	scrollCharacterMemory(layerEntry)
*/
func scrollCharacterMemory(layerEntry *types.LayerEntryType) [][]types.CharacterEntryType {
	layerWidth := layerEntry.Width
	characterMemory := layerEntry.CharacterMemory
	characterMemory = characterMemory[1:]
	characterObjectArray := make([]types.CharacterEntryType, layerWidth)
	for currentCharacterCell := 0; currentCharacterCell < layerWidth; currentCharacterCell++ {
		characterEntry := types.NewCharacterEntry()
		characterEntry.AttributeEntry = types.NewAttributeEntry(&layerEntry.DefaultAttribute)
		characterEntry.Character = ' '
		characterObjectArray[currentCharacterCell] = characterEntry
	}
	characterMemory = append(characterMemory, characterObjectArray)
	layerEntry.CharacterMemory = characterMemory
	return characterMemory
}

/*
getRuneOnLayer is a method which allows you to obtain a specific rune at the location specified on the given text layer.
In addition, the following should be noted:

- If the location specified is outside the valid range, 0 is returned.

Example:

	char := getRuneOnLayer(layerEntry, 10, 5)
*/
func getRuneOnLayer(layerEntry *types.LayerEntryType, xLocation int, yLocation int) rune {
	// validateLayerLocationByLayerEntry(layerEntry, xLocation, yLocation)
	// We don't do validation here because if we draw outside the layer intentionally, we don't want to panic.
	// For example, if rendering a control that pops off-screen.
	if xLocation < 0 || yLocation < 0 ||
		xLocation >= layerEntry.Width || yLocation >= layerEntry.Height {
		return 0
	}
	characterMemory := layerEntry.CharacterMemory
	return characterMemory[yLocation][xLocation].Character
}

/*
GetCellIdUnderMouseLocation is a method which allows you to obtain the cell ID for the text directly under your mouse
cursor. In addition, the following should be noted:

- If multiple text layers are being displayed, the cell ID returned will be from the top-most visible text cell.

- The cell ID returned will only reflect what is currently being displayed on the terminal display.

Example:

	cellId := GetCellIdUnderMouseLocation()
*/
func GetCellIdUnderMouseLocation() int {
	mouseXLocation, mouseYLocation, _, _ := GetMouseStatus()
	return getCellIdByLayerEntry(&commonResource.screenLayer, mouseXLocation, mouseYLocation)
}

/*
getCellIdByLayerAlias is a method which allows you to obtain a cell ID from a given text layer by layer alias.

Example:

	cellId := getCellIdByLayerAlias("MyLayer", 10, 5)
*/
func getCellIdByLayerAlias(layerAlias string, mouseXLocation int, mouseYLocation int) int {
	validateLayer(layerAlias)
	layerEntry := Layers.Get(layerAlias)
	return getCellIdByLayerEntry(layerEntry, mouseXLocation, mouseYLocation)
}

/*
getCellIdByLayerEntry is a method which allows you to obtain a cell ID from a given text layer by layer entry. In
addition, the following should be noted:

- If the location specified is outside the valid range of the text layer, then a value of -1 is returned.

Example:

	cellId := getCellIdByLayerEntry(layerEntry, 10, 5)
*/
func getCellIdByLayerEntry(layerEntry *types.LayerEntryType, xLocation int, yLocation int) int {
	returnValue := -1
	if xLocation < 0 || xLocation >= layerEntry.Width || yLocation < 0 || yLocation >= layerEntry.Height {
		return returnValue
	}
	if yLocation-layerEntry.ScreenYLocation >= 0 && xLocation-layerEntry.ScreenXLocation >= 0 &&
		yLocation-layerEntry.ScreenYLocation < len(layerEntry.CharacterMemory) && xLocation-layerEntry.ScreenXLocation < len(layerEntry.CharacterMemory[0]) {
		characterEntry := layerEntry.CharacterMemory[yLocation-layerEntry.ScreenYLocation][xLocation-layerEntry.ScreenXLocation]
		returnValue = characterEntry.AttributeEntry.CellUserId
	}
	return returnValue
}

/*
UpdateDisplay is a method which allows you to synchronize the terminal's visible display area with your current changes.
In addition, the following should be noted:

- All text layers are sorted from lowest to highest z-order priority level.

- Layers with the same z-order priority will appear in random display order.

Example:

	UpdateDisplay(false)
*/
func UpdateDisplay(isRefreshForced bool) {
	commonResource.displayUpdate.Lock()
	defer func() {
		commonResource.displayUpdate.Unlock()
	}()
	sortedLayerAliasSlice := layer.GetSortedLayerMemoryAliasSlice()
	baseLayerEntry := types.NewLayerEntry("", "", commonResource.terminalWidth, commonResource.terminalHeight)
	baseLayerEntry = renderLayers(&baseLayerEntry, sortedLayerAliasSlice)
	Tooltip.renderAll(baseLayerEntry)
	DrawLayerToScreen(&baseLayerEntry, isRefreshForced)
	commonResource.screenLayer = baseLayerEntry
}

/*
RefreshDisplay is a method which allows you to sync the terminal screen.

Example:

	RefreshDisplay()
*/
func RefreshDisplay() {
	commonResource.screen.Sync()
}

/*
renderLayers is a method which allows you to render a list of text layers to the specified root text layer. In addition,
the following should be noted:

- If a text layer being rendered is a parent, then all child text layers will be rendered recursively.

- Any text layer which is marked as not visible will be ignored.

Example:

	renderLayers(&rootLayer, aliases)
*/
func renderLayers(rootLayerEntry *types.LayerEntryType, sortedLayerAliasSlice LayerAliasZOrderPairList) types.LayerEntryType {
	baseLayerEntry := types.NewLayerEntry("", "", 0, 0, rootLayerEntry)
	isOpaque := true
	for currentListIndex := 0; currentListIndex < len(sortedLayerAliasSlice); currentListIndex++ {
		if !Layers.IsExists(sortedLayerAliasSlice[currentListIndex].Key) {
			continue
		}
		currentLayerEntry := types.NewLayerEntry("", "", 0, 0, Layers.Get(sortedLayerAliasSlice[currentListIndex].Key))
		if currentLayerEntry.IsVisible {
			renderControls(currentLayerEntry)
			if currentLayerEntry.IsParent && (currentLayerEntry.LayerAlias != baseLayerEntry.LayerAlias && currentLayerEntry.ParentAlias == baseLayerEntry.LayerAlias) {
				renderedLayer := renderLayers(&currentLayerEntry, sortedLayerAliasSlice)
				overlayLayers(&renderedLayer, &baseLayerEntry, isOpaque)
			} else {
				if currentLayerEntry.ParentAlias == baseLayerEntry.LayerAlias {
					overlayLayers(&currentLayerEntry, &baseLayerEntry, isOpaque)
				}
			}
		}
		// After the first layer is rendered to the base, allow for transparencies.
		if isOpaque {
			isOpaque = false
		}
	}
	return baseLayerEntry
}

/*
renderControls is a method which allows you to draw various control elements on the specified layer. The order of
drawing matters, as complex controls are drawn first above basic controls. In addition, the following should be noted:

- Tooltip hotspot zones must be drawn before FileMenu to prevent them from capturing clicks intended for file menu.

Example:

	renderControls(layerEntry)
*/
func renderControls(currentLayerEntry types.LayerEntryType) {
	Button.drawOnLayer(currentLayerEntry)
	TextField.drawOnLayer(currentLayerEntry)
	Checkbox.drawOnLayer(currentLayerEntry)
	radioButton.drawOnLayer(currentLayerEntry)
	ProgressBar.drawOnLayer(currentLayerEntry)
	Label.drawOnLayer(currentLayerEntry)
	scrollbar.drawOnLayer(currentLayerEntry)

	textbox.drawOnLayer(currentLayerEntry)
	Tooltip.drawHotspotZonesOnLayer(currentLayerEntry)
	viewport.drawOnLayer(currentLayerEntry)
	FileMenu.drawOnLayer(currentLayerEntry) // File menu must appear before selector or selectors won't render when menu is open.
	Dropdown.drawOnLayer(currentLayerEntry) // Dropdowns must come before selectors, or it won't show on top.
	Selector.drawSelectorsOnLayer(currentLayerEntry)

}

/*
overlayLayersByLayerAlias is a method which allows you to overlay a text layer by its layer alias. In addition, the
following should be noted:

- This is useful when you do not have actual layer data and only know the alias.

Example:

	overlayLayersByLayerAlias("SrcLayer", &targetLayer)
*/
func overlayLayersByLayerAlias(sourceLayerAlias string, targetLayerEntry *types.LayerEntryType) {
	validateLayer(sourceLayerAlias)
	layerEntry := Layers.Get(sourceLayerAlias)
	overlayLayers(layerEntry, targetLayerEntry, false)
}

/*
copyCharacterMemory is a method which allows you to copy a portion of a source character memory to a target character
memory. In addition, the following should be noted:

- If the source character memory to be drawn is outside the target, then only the visible portion will be rendered.

Example:

	copyCharacterMemory(srcMem, targetMem, 0, 0, 10, 10)
*/
func copyCharacterMemory(sourceCharacterMemory [][]types.CharacterEntryType, targetCharacterMemory [][]types.CharacterEntryType, xLocation, yLocation, width, height int) {
	sourceHeight := len(sourceCharacterMemory)
	if sourceHeight == 0 {
		return
	}
	sourceWidth := len(sourceCharacterMemory[0])
	if sourceWidth == 0 {
		return
	}
	targetHeight := len(targetCharacterMemory)
	if targetHeight == 0 {
		return
	}
	targetWidth := len(targetCharacterMemory[0])
	if targetWidth == 0 {
		return
	}
	for currentRow := 0; currentRow < height; currentRow++ {
		for currentColumn := 0; currentColumn < width; currentColumn++ {
			targetColumn := xLocation + currentColumn
			targetRow := yLocation + currentRow
			if currentColumn < sourceWidth && currentRow < sourceHeight && targetColumn < targetWidth && targetRow < targetHeight {
				targetCharacterMemory[targetRow][targetColumn] = sourceCharacterMemory[currentRow][currentColumn]
			}
		}
	}
}

/*
getEffectiveAlpha is a method which calculates the final transparency value of a cell by multiplying the layer's
global alpha with the cell-specific or default attribute alpha. It ensures that transparency is applied
hierarchically so that a semi-transparent layer correctly attenuates the opacity of its individual character cells.

Example:
    alpha := getEffectiveAlpha(0.5, 1.0, 0.8)
*/
func getEffectiveAlpha(layerAlpha float32, defaultAlpha float32, cellAlpha float32) float32 {
	effectiveAlpha := layerAlpha
	if cellAlpha < 1 {
		effectiveAlpha *= cellAlpha
	} else if defaultAlpha < 1 {
		effectiveAlpha *= defaultAlpha
	}
	return effectiveAlpha
}

/*
shouldShowSource is a method which determines if a character from the source layer should be drawn over the target
layer based on its alpha value and a dithering strategy. It compares the effective alpha against either a random float
for stochastic transparency or a pre-defined threshold from a Bayer matrix for patterned dithering, returning true if
the character is visible. In addition, the following should be noted:

- The stochastic strategy relies on the global random number generator which can affect reproducibility if not seeded.

Example:
    show := shouldShowSource(0.5, constants.TransparencyStrategyStochastic, 10, 5)
*/
func shouldShowSource(effectiveAlpha float32, strategy constants.TransparencyStrategy, x, y int) bool {
	if effectiveAlpha >= 1 {
		return true
	}
	// Add a small epsilon to ensure that values very close to 0 are treated as 0
	if effectiveAlpha <= 0.0001 {
		return false
	}
	switch strategy {
	case constants.TransparencyStrategyStochastic:
		return rand.Float32() < effectiveAlpha
	case constants.TransparencyStrategy2x2Bayer:
		return effectiveAlpha > constants.BayerMatrix2x2[y%2][x%2]
	case constants.TransparencyStrategy4x4Bayer:
		return effectiveAlpha > constants.BayerMatrix4x4[y%4][x%4]
	case constants.TransparencyStrategy8x8Bayer:
		return effectiveAlpha > constants.BayerMatrix8x8[y%8][x%8]
	case constants.TransparencyStrategyDissolve:
		return true
	default:
		return true
	}
}

/*
compositeCell is a method which blends the contents and attributes of a source cell with a target cell to produce a
single rendered character entry. It calculates color transitions for foregrounds and backgrounds, handles special cell
types like shadows and tooltips, and processes transparency flags to ensure layers are combined with visual accuracy.

Example:
    result := compositeCell(&source, &target, 0.5, 1.0, 1.0, false)
*/
func compositeCell(sourceEntry *types.CharacterEntryType, targetEntry *types.CharacterEntryType, layerAlpha float32, defaultFgAlpha float32, defaultBgAlpha float32, isOpaque bool, isBinaryAlpha bool, strategy constants.TransparencyStrategy) types.CharacterEntryType {
	sourceAttributeEntry := sourceEntry.AttributeEntry
	targetAttributeEntry := targetEntry.AttributeEntry

	// For dithering strategies, we treat the layer as fully opaque if it's shown at all
	blendAlpha := layerAlpha
	if isBinaryAlpha {
		blendAlpha = 1.0
	}

	// Handle NullRune (transparent) cells
	if sourceEntry.Character == constants.NullRune {
		resultEntry := *targetEntry
		if sourceAttributeEntry.CellType == constants.CellTypeShadow {
			if sourceAttributeEntry.ForegroundAlphaValue < 1 {
				resultEntry.AttributeEntry.ForegroundColor = GetTransitionedColor(targetAttributeEntry.ForegroundColor, GetRGBColor(0, 0, 0), sourceAttributeEntry.ForegroundAlphaValue*blendAlpha)
			}
			if sourceAttributeEntry.BackgroundAlphaValue < 1 {
				resultEntry.AttributeEntry.BackgroundColor = GetTransitionedColor(targetAttributeEntry.BackgroundColor, GetRGBColor(0, 0, 0), sourceAttributeEntry.BackgroundAlphaValue*blendAlpha)
			}
			resultEntry.AttributeEntry.CellType = constants.CellTypeShadow
		}
		if sourceAttributeEntry.CellType == constants.CellTypeTooltip {
			resultEntry.AttributeEntry.CellType = constants.CellTypeTooltip
			resultEntry.AttributeEntry.CellControlAlias = sourceAttributeEntry.CellControlAlias
			resultEntry.LayerAlias = sourceEntry.LayerAlias
		}
		return resultEntry
	}

	resultEntry := *sourceEntry
	newAttributeEntry := types.NewAttributeEntry(&sourceAttributeEntry)

	// Copy layer and parent aliases if they exist
	if sourceEntry.LayerAlias != "" {
		resultEntry.LayerAlias = sourceEntry.LayerAlias
	} else {
		resultEntry.LayerAlias = targetEntry.LayerAlias
	}
	if sourceEntry.ParentAlias != "" {
		resultEntry.ParentAlias = sourceEntry.ParentAlias
	} else {
		resultEntry.ParentAlias = targetEntry.ParentAlias
	}

	// --- Efficient transparency handling ---
	if !isOpaque {
		if sourceAttributeEntry.IsForegroundTransparent {
			newAttributeEntry.ForegroundColor = targetAttributeEntry.ForegroundColor
			newAttributeEntry.IsForegroundTransparent = true
		}
		if sourceAttributeEntry.IsBackgroundTransparent {
			newAttributeEntry.BackgroundColor = targetAttributeEntry.BackgroundColor
			newAttributeEntry.IsBackgroundTransparent = true
		}
	}

	// Apply color transformations
	effectiveFgAlpha := getEffectiveAlpha(blendAlpha, defaultFgAlpha, sourceAttributeEntry.ForegroundAlphaValue)
	effectiveBgAlpha := getEffectiveAlpha(blendAlpha, defaultBgAlpha, sourceAttributeEntry.BackgroundAlphaValue)

	if strategy == constants.TransparencyStrategyDissolve {
		// For Dissolve, background always transitions smoothly
		newAttributeEntry.BackgroundColor = GetTransitionedColor(targetAttributeEntry.BackgroundColor, sourceAttributeEntry.BackgroundColor, effectiveBgAlpha)
		blendedBG := newAttributeEntry.BackgroundColor

		if effectiveFgAlpha > 0.5 {
			// Alpha 1.0 -> 0.5: Transition from source FG to blended BG
			normalizedAlpha := (effectiveFgAlpha - 0.5) * 2
			newAttributeEntry.ForegroundColor = GetTransitionedColor(blendedBG, sourceAttributeEntry.ForegroundColor, normalizedAlpha)
			resultEntry.Character = sourceEntry.Character
		} else {
			// Alpha 0.5 -> 0.0: Transition from blended BG to target FG
			normalizedAlpha := effectiveFgAlpha * 2
			newAttributeEntry.ForegroundColor = GetTransitionedColor(targetAttributeEntry.ForegroundColor, blendedBG, normalizedAlpha)
			resultEntry.Character = targetEntry.Character
		}
	} else {
		if effectiveFgAlpha < 1 {
			newAttributeEntry.ForegroundColor = GetTransitionedColor(targetAttributeEntry.ForegroundColor, sourceAttributeEntry.ForegroundColor, effectiveFgAlpha)
		}
		if effectiveBgAlpha < 1 {
			newAttributeEntry.BackgroundColor = GetTransitionedColor(targetAttributeEntry.BackgroundColor, sourceAttributeEntry.BackgroundColor, effectiveBgAlpha)
		}
	}

	resultEntry.AttributeEntry = newAttributeEntry
	return resultEntry
}

/*
overlayLayers is a method which allows you to overlay one text layer on top of another text layer. In addition, the
following should be noted:

- If the source rune to be drawn is null, then it will be considered transparent.

- If a transparent rune has a foreground or background alpha value set, then it will be drawn as a shadow.

Example:

	overlayLayers(&srcLayer, &targetLayer, false)
*/
func overlayLayers(sourceLayerEntry *types.LayerEntryType, targetLayerEntry *types.LayerEntryType, isOpaque bool) {
	// 1. Simplified Clipping Logic (Integer Math)
	sourceStartX := 0
	if sourceLayerEntry.ScreenXLocation < 0 {
		sourceStartX = -sourceLayerEntry.ScreenXLocation
	}

	sourceStartY := 0
	if sourceLayerEntry.ScreenYLocation < 0 {
		sourceStartY = -sourceLayerEntry.ScreenYLocation
	}

	targetStartX := 0
	if sourceLayerEntry.ScreenXLocation > 0 {
		targetStartX = sourceLayerEntry.ScreenXLocation
	}

	targetStartY := 0
	if sourceLayerEntry.ScreenYLocation > 0 {
		targetStartY = sourceLayerEntry.ScreenYLocation
	}

	// Calculate the width of the overlapping area.
	widthToCopy := sourceLayerEntry.Width - sourceStartX
	if targetLayerEntry.Width-targetStartX < widthToCopy {
		widthToCopy = targetLayerEntry.Width - targetStartX
	}

	// Calculate the height of the overlapping area.
	heightToCopy := sourceLayerEntry.Height - sourceStartY
	if targetLayerEntry.Height-targetStartY < heightToCopy {
		heightToCopy = targetLayerEntry.Height - targetStartY
	}

	// If there's no overlapping area, we have nothing to do.
	if widthToCopy <= 0 || heightToCopy <= 0 {
		return
	}

	// 2. Cache Lookups
	sourceCharacterMemory := sourceLayerEntry.CharacterMemory
	targetCharacterMemory := targetLayerEntry.CharacterMemory
	defaultFgAlpha := sourceLayerEntry.DefaultAttribute.ForegroundAlphaValue
	defaultBgAlpha := sourceLayerEntry.DefaultAttribute.BackgroundAlphaValue
	layerAlpha := sourceLayerEntry.AlphaValue
	strategy := sourceLayerEntry.TransparencyStrategy
	isBinaryAlpha := strategy != constants.TransparencyStrategyColorOnly && strategy != constants.TransparencyStrategyDissolve

	// 3. Parallel Processing with Goroutines
	var wg sync.WaitGroup
	wg.Add(heightToCopy)

	for currentRow := 0; currentRow < heightToCopy; currentRow++ {
		go func(row int) {
			defer wg.Done()
			sourceRow := row + sourceStartY
			targetRow := row + targetStartY

			for currentColumn := 0; currentColumn < widthToCopy; currentColumn++ {
				targetCol := currentColumn + targetStartX
				sourceCol := currentColumn + sourceStartX

				sourceEntry := &sourceCharacterMemory[sourceRow][sourceCol]
				targetEntry := &targetCharacterMemory[targetRow][targetCol]

				effectiveAlpha := getEffectiveAlpha(layerAlpha, defaultFgAlpha, sourceEntry.AttributeEntry.ForegroundAlphaValue)

				if !shouldShowSource(effectiveAlpha, strategy, targetCol, targetRow) {
					targetCharacterMemory[targetRow][targetCol] = *targetEntry
					continue
				}

				targetCharacterMemory[targetRow][targetCol] = compositeCell(sourceEntry, targetEntry, layerAlpha, defaultFgAlpha, defaultBgAlpha, isOpaque, isBinaryAlpha, strategy)
			}
		}(currentRow)
	}
	wg.Wait()
}

/*
DrawLayerToScreen is a method which allows you to render a text layer to the visible terminal screen. In addition, the
following should be noted:

- If debug is enabled, this method does nothing since the terminal is virtual.

Example:

	DrawLayerToScreen(layerEntry, false)
*/
func DrawLayerToScreen(layerEntry *types.LayerEntryType, isForcedRefreshRequired bool) {
	if !commonResource.isDebugEnabled {
		width := layerEntry.Width
		height := layerEntry.Height
		for currentRow := 0; currentRow < height; currentRow++ {
			for currentCharacter := 0; currentCharacter < width; currentCharacter++ {
				style := tcell.StyleDefault
				attributeEntry := layerEntry.CharacterMemory[currentRow][currentCharacter].AttributeEntry
				style = style.Foreground(tcell.Color(attributeEntry.ForegroundColor))
				style = style.Background(tcell.Color(attributeEntry.BackgroundColor))
				style = style.Blink(attributeEntry.IsBlinking)
				style = style.Bold(attributeEntry.IsBold)
				style = style.Reverse(attributeEntry.IsReversed)
				style = style.Underline(attributeEntry.IsUnderlined)
				var character = layerEntry.CharacterMemory[currentRow][currentCharacter].Character
				r2 := []rune("")
				commonResource.screen.SetContent(currentCharacter, currentRow, character, r2, style)
			}
		}
		if isForcedRefreshRequired {
			commonResource.screen.Sync()
		}
		commonResource.screen.Show()
	}
}

/*
GetOs is a method which allows you to obtain the name of the operating system currently running.

Example:

	osName := GetOs()
*/
func GetOs() string {
	switch runtime.GOOS {
	case "windows":
		return constants.OS_WINDOWS
	case "linux":
		return constants.OS_LINUX
	case "darwin":
		return constants.OS_MAC
	default:
		return constants.OS_OTHER
	}
}
