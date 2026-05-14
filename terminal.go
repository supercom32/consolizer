package consolizer

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/stringformat"
	"github.com/supercom32/consolizer/types"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
)

const (
	OS_WINDOWS = 1
	OS_LINUX   = 2
	OS_MAC     = 3
	OS_OTHER   = 4
)

/*
defaultValueType is a class which allows you to This class is a structure that holds common information about the
current terminal session that needs to be shared.
*/
type defaultValueType struct {
	screen               tcell.Screen
	layerInstance        *LayerInstanceType // What happens when last layer is deleted? This needs to be updated.
	terminalWidth        int
	terminalHeight       int
	screenLayer          types.LayerEntryType
	debugDirectory       string
	isDebugEnabled       bool
	displayUpdate        sync.Mutex
	updateDisplayChannel chan bool
}

/*
commonResource is a variable used to hold shared data that is accessed by this package.
*/
var commonResource defaultValueType

/*
GetVersion is a method which allows you to obtain the current version of the consolizer library.

:return: A string representing the version number.

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

- If you pass in a zero or negative value for either width or height a panic will be generated to fail as fast as.

:param width: The desired width of the terminal.
:param height: The desired height of the terminal.

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
setupPeriodicEventUpdater is a method which allows you to This method is a background method that updates periodic
events.

Example:

	go setupPeriodicEventUpdater()
*/
func setupPeriodicEventUpdater() {
	for {
		UpdatePeriodicEvents()
	}
}

/*
setupEventUpdater is a method which allows you to This method is a background method that monitors all events coming
into the terminal session. When an event is detected, it is recorded and monitoring continues.

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
setupCloseHandler is a method which allows you to This method enables the trapping of all unexpected system calls and
shuts down the terminal gracefully. This means all terminal settings should be reset back to normal if anything
unexpected happens to the user or if the process is killed.

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
RestoreTerminalSettings is a method which allows you to This method allows the user to gracefully return the terminal
back to its normal settings. This should be called once your application is finished using consolizer so that the users
terminal environment is not left in a bad state.

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

:return: An integer representing the OS type.

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

:return: The width and height of the terminal.

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

:return: A slice of runes representing the keystroke.

Example:

	keystroke := Inkey()
*/
func Inkey() []rune {
	return KeyboardMemory.GetFromBuffer()
}

/*
Layer is a method which allows you to specify a default layer alias that you wish to use when interacting with methods
which have a non-layer alias method signature. Non-layer alias method signatures can be identified by finding methods
which have both a layer and non-layer version. This makes interacting with methods faster, as the user does not need to
provide the layer alias context in which he is working on. In addition, the following should be noted:

- The provided layer instance must be valid and exist.

:param layerInstance: A pointer to the layer instance to set as default.

Example:

	consolizer.Layer(myLayerInstance)
*/
func Layer(layerInstance *LayerInstanceType) {
	validateLayer(layerInstance.layerAlias)
	commonResource.layerInstance = layerInstance
}

/*
AddTextStyle is a method which allows you to add a new printing style which can be used for writing dialog text with.
Dialog text printing differs from regular printing, since it allows for typewriter drawing effects as well as changing
text attributes on the fly (color, bold, etc). In addition, the following should be noted:

- This method expects you to pass in a textStyleEntry type obtained by calling consolizer.NewTextStyle.

:param textStyleAlias: The alias to assign to the text style.
:param textStyleEntry: The text style entry to add.

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

:param textStyleAlias: The alias of the text style to delete.

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

:return: A new text cell style entry.

Example:

	myTextStyleEntry := consolizer.NewTextStyle()
*/
func NewTextStyle() types.TextCellStyleEntryType {
	return types.NewTextCellStyleEntry()
}

/*
NewImageStyle is a constructor which allows you to obtain a new image style entry.

:return: A new image style entry.

Example:

	myImageStyle := NewImageStyle()
*/
func NewImageStyle() types.ImageStyleEntryType {
	return types.NewImageStyleEntry()
}

/*
NewTuiStyleEntry is a constructor which allows you to obtain a new style entry which can be used for specifying how TUI
controls and other TUI drawing operations should occur.

:return: A new TUI style entry.

Example:

	myTuiStyleEntry := consolizer.NewTuiStyleEntry()
*/
func NewTuiStyleEntry() types.TuiStyleEntryType {
	return types.NewTuiStyleEntry()
}

/*
NewSelectionEntry is a constructor which allows you to obtain an entry used for specifying what options you want to make
available for a given menu prompt.

:return: A new selection entry.

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

:return: A new asset list type.

Example:

	assetList := consolizer.NewAssetList()
*/
func NewAssetList() types.AssetListType {
	return types.NewAssetList()
}

/*
SetLayerZOrder is a method which allows you to set the z-order value for a given layer. In addition, the following
should be noted:

- The z-order priority controls which text layer should be drawn first and which text layer should be drawn last.

- Layers that have a higher priority will be drawn on top of layers that have a lower priority.

- In the event that two layers have the same priority, they will be drawn in random order.

:param layerInstance: A pointer to the layer instance.
:param zOrder: The z-order value to set.

Example:

	SetLayerZOrder(myLayer, 10)
*/
func SetLayerZOrder(layerInstance *LayerInstanceType, zOrder int) {
	layerEntry := Layers.Get(layerInstance.layerAlias)
	layerEntry.ZOrder = zOrder
}

/*
SetZOrder is a method which allows you to set the z-order value for the default currently selected layer. In addition,
the following should be noted:

- The z-order priority controls which text layer should be drawn first and which text layer should be drawn last.

- Layers that have a higher priority will be drawn on top of layers that have a lower priority.

:param zOrder: The z-order value to set.

Example:

	SetZOrder(10)
*/
func SetZOrder(zOrder int) {
	SetLayerZOrder(commonResource.layerInstance, zOrder)
}

/*
SetLayerAlpha is a method which allows you to set the alpha value for a given text layer. This lets you perform pseudo
transparencies by making the layer foreground and background colors blend with the layers underneath it to the degree
specified. In addition, the following should be noted:

- An alpha value of 1.0 is equal to 100% visible, while an alpha value of 0.0 is 0% visible.

- If the percent change specified is outside of the RGB color range, then the color will simply bottom or max out.

:param layerInstance: A pointer to the layer instance.
:param alphaValue: The alpha value to set.

Example:

	SetLayerAlpha(myLayer, 0.5)
*/
func SetLayerAlpha(layerInstance *LayerInstanceType, alphaValue float32) {
	setLayerAlpha(layerInstance, alphaValue)
}

/*
setLayerAlpha is a method which allows you to This method is an internal method that sets the alpha value for a given
layer instance.

:param layerInstance: A pointer to the layer instance.
:param alphaValue: The alpha value to set.

Example:

	setLayerAlpha(myLayer, 0.5)
*/
func setLayerAlpha(layerInstance *LayerInstanceType, alphaValue float32) {
	layerEntry := Layers.Get(layerInstance.layerAlias)
	layerEntry.DefaultAttribute.ForegroundAlphaValue = alphaValue
	layerEntry.DefaultAttribute.BackgroundAlphaValue = alphaValue
}

/*
SetLayer is a method which allows you to set the alpha value for the default currently selected layer.

:param alphaValue: The alpha value to set.

Example:

	SetLayer(0.5)
*/
func SetLayer(alphaValue float32) {
	setLayerAlpha(commonResource.layerInstance, alphaValue)
}

/*
GetColor is a method which allows you to obtain an RGB color from a predefined color palette list. This list corresponds
to the 16 color ANSI standard, where color 0 is Black and 15 is Bright White. In addition, the following should be
noted:

- If you specify a color index less than 0 or greater than 15 a panic will be generated to fail as fast as possible.

:param colorIndex: The index of the color in the ANSI palette.

:return: The color type corresponding to the index.

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

:param redColorIndex: The red channel value (0-255).
:param greenColorIndex: The green channel value (0-255).
:param blueColorIndex: The blue channel value (0-255).

:return: The resulting color type.

Example:

	color := GetRGBColor(255, 128, 0)
*/
func GetRGBColor(redColorIndex int32, greenColorIndex int32, blueColorIndex int32) constants.ColorType {
	validateRGBColorIndex(redColorIndex, greenColorIndex, blueColorIndex)
	return constants.ColorType(tcell.NewRGBColor(redColorIndex, greenColorIndex, blueColorIndex))
}

/*
Color is a method which allows you to set default colors on your text layer for printing with. The color index specified
corresponds to the 16 color ANSI standard, where color 0 is Black and 15 is Bright White. If you wish to change colors
settings for a text layer that is not currently set as your default, use ColorLayer instead.

:param foregroundColorIndex: The ANSI index for the foreground color.
:param backgroundColorIndex: The ANSI index for the background color.

Example:

	Color(15, 0)
*/
func Color(foregroundColorIndex int, backgroundColorIndex int) {
	validateDefaultLayerIsNotEmpty()
	ColorLayer(commonResource.layerInstance, foregroundColorIndex, backgroundColorIndex)
}

/*
ColorLayer is a method which allows you to set default colors on your specified text layer for printing with. The color
index specified corresponds to the 16 color ANSI standard, where color 0 is Black and 15 is Bright White.

:param layerInstance: A pointer to the layer instance.
:param foregroundColorIndex: The ANSI index for the foreground color.
:param backgroundColorIndex: The ANSI index for the background color.

Example:

	ColorLayer(myLayer, 15, 0)
*/
func ColorLayer(layerInstance *LayerInstanceType, foregroundColorIndex int, backgroundColorIndex int) {
	validateColorIndex(foregroundColorIndex)
	validateColorIndex(backgroundColorIndex)
	layerEntry := Layers.Get(layerInstance.layerAlias)
	layerEntry.DefaultAttribute.ForegroundColor = constants.AnsiColorByIndex[foregroundColorIndex]
	layerEntry.DefaultAttribute.BackgroundColor = constants.AnsiColorByIndex[backgroundColorIndex]
}

/*
ColorRGB is a method which allows you to set default colors on your text layer for printing with using RGB values. This
method allows you to specify colors using RGB color index values within the range of 0 to 255.

:param foregroundRedIndex: Red channel for foreground.
:param foregroundGreenIndex: Green channel for foreground.
:param foregroundBlueIndex: Blue channel for foreground.
:param backgroundRedIndex: Red channel for background.
:param backgroundGreenIndex: Green channel for background.
:param backgroundBlueIndex: Blue channel for background.

Example:

	ColorRGB(255, 255, 255, 0, 0, 0)
*/
func ColorRGB(foregroundRedIndex int32, foregroundGreenIndex int32, foregroundBlueIndex int32, backgroundRedIndex int32, backgroundGreenIndex int32, backgroundBlueIndex int32) {
	validateDefaultLayerIsNotEmpty()
	ColorLayerRGB(commonResource.layerInstance, foregroundRedIndex, foregroundGreenIndex, foregroundBlueIndex, backgroundRedIndex, backgroundGreenIndex, backgroundBlueIndex)
}

/*
ColorLayerRGB is a method which allows you to set default colors on your specified text layer for printing with using
RGB values. This method allows you to specify colors using RGB color index values within the range of 0 to 255.

:param layerInstance: A pointer to the layer instance.
:param foregroundRed: Red channel for foreground.
:param foregroundGreen: Green channel for foreground.
:param foregroundBlue: Blue channel for foreground.
:param backgroundRed: Red channel for background.
:param backgroundGreen: Green channel for background.
:param backgroundBlue: Blue channel for background.

Example:

	ColorLayerRGB(myLayer, 255, 255, 255, 0, 0, 0)
*/
func ColorLayerRGB(layerInstance *LayerInstanceType, foregroundRed int32, foregroundGreen int32, foregroundBlue int32, backgroundRed int32, backgroundGreen int32, backgroundBlue int32) {
	foregroundColor := GetRGBColor(foregroundRed, foregroundGreen, foregroundBlue)
	backgroundColor := GetRGBColor(backgroundRed, backgroundGreen, backgroundBlue)
	ColorLayer24Bit(layerInstance, foregroundColor, backgroundColor)
}

/*
Color24Bit is a method which allows you to color the default layer using a 24-bit color expressed as an int32.

:param foregroundColor: The 24-bit foreground color.
:param backgroundColor: The 24-bit background color.

Example:

	Color24Bit(fgColor, bgColor)
*/
func Color24Bit(foregroundColor constants.ColorType, backgroundColor constants.ColorType) {
	ColorLayer24Bit(commonResource.layerInstance, foregroundColor, backgroundColor)
}

/*
ColorLayer24Bit is a method which allows you to color a specified layer using a 24-bit color expressed as an int32.

:param layerInstance: A pointer to the layer instance.
:param foregroundColor: The 24-bit foreground color.
:param backgroundColor: The 24-bit background color.

Example:

	ColorLayer24Bit(myLayer, fgColor, bgColor)
*/
func ColorLayer24Bit(layerInstance *LayerInstanceType, foregroundColor constants.ColorType, backgroundColor constants.ColorType) {
	layerEntry := Layers.Get(layerInstance.layerAlias)
	layerEntry.DefaultAttribute.ForegroundColor = foregroundColor
	layerEntry.DefaultAttribute.BackgroundColor = backgroundColor
}

/*
Locate is a method which allows you to set the default cursor location on your specified text layer for printing with.
In addition, the following should be noted:

  - If you pass in a location value that falls outside the dimensions of the default text layer, a panic will be
    generated.

- Valid text layer locations start at position (0,0) for the upper left corner.

:param xLocation: The x-axis location for the cursor.
:param yLocation: The y-axis location for the cursor.

Example:

	Locate(10, 5)
*/
func Locate(xLocation int, yLocation int) {
	validateDefaultLayerIsNotEmpty()
	LocateLayer(commonResource.layerInstance, xLocation, yLocation)
}

/*
LocateLayer is a method which allows you to set the default cursor location on your specified text layer for printing
with. In addition, the following should be noted:

  - If you pass in a location value that falls outside the dimensions of the specified text layer, a panic will be
    generated.

:param layerInstance: A pointer to the layer instance.
:param xLocation: The x-axis location for the cursor.
:param yLocation: The y-axis location for the cursor.

Example:

	LocateLayer(myLayer, 10, 5)
*/
func LocateLayer(layerInstance *LayerInstanceType, xLocation int, yLocation int) {
	validateLayer(layerInstance.layerAlias)
	layerEntry := Layers.Get(layerInstance.layerAlias)
	validateLayerLocationByLayerEntry(layerEntry, xLocation, yLocation)
	layerEntry.CursorXLocation = xLocation
	layerEntry.CursorYLocation = yLocation
}

/*
Print is a method which allows you to write text to the default text layer. In addition, the following should be noted:

- When text is written to the text layer, the cursor position is also updated to reflect its new location.

- If the string to print ends up being too long to fit at its current location, then only the visible portion of your.

- If printing has not yet finished and there are no available lines left, then all remaining characters will be.

:param textToPrint: The string of text to print.

Example:

	Print("Hello World")
*/
func Print(textToPrint string) {
	validateDefaultLayerIsNotEmpty()
	PrintLayer(commonResource.layerInstance, textToPrint)
}

/*
PrintLayer is a method which allows you to write text to a specified text layer. In addition, the following should be
noted:

- When text is written to the text layer, the cursor position is also updated to reflect its new location.

- If the string to print ends up being too long to fit at its current location, then only the visible portion of your.

:param layerInstance: A pointer to the layer instance.
:param textToPrint: The string of text to print.

Example:

	PrintLayer(myLayer, "Hello World")
*/
func PrintLayer(layerInstance *LayerInstanceType, textToPrint string) {
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

- If the location to print falls outside the range of the text layer, then only the visible portion of your text will.

:param layerEntry: A pointer to the layer entry.
:param attributeEntry: The attribute entry for the text.
:param xLocation: The x-axis location.
:param yLocation: The y-axis location.
:param textToPrint: A slice of runes to print.

:return: The final x-axis location after printing.

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
Clear is a method which allows you to empty the default text layer of all its contents.

Example:

	Clear()
*/
func Clear() {
	validateDefaultLayerIsNotEmpty()
	ClearLayer(commonResource.layerInstance)
}

/*
clear is a method which allows you to empty the specified text layer of all its contents.

:param layerEntry: A pointer to the layer entry to clear.

Example:

	clearLayer(layerEntry)
*/
func clearLayer(layerEntry *types.LayerEntryType) {
	types.InitializeCharacterMemory(layerEntry)
}

/*
GetCharacterOnScreen is a method which allows you to obtain the character currently being displayed on the screen at a
specific location.

:param xLocation: The x-axis location.
:param yLocation: The y-axis location.

:return: The rune at the specified location.

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

:param layerEntry: A pointer to the layer entry to scroll.

:return: The updated character memory.

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

:param layerEntry: A pointer to the layer entry.
:param xLocation: The x-axis location.
:param yLocation: The y-axis location.

:return: The rune at the specified location.

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

:return: The cell ID under the mouse cursor.

Example:

	cellId := GetCellIdUnderMouseLocation()
*/
func GetCellIdUnderMouseLocation() int {
	mouseXLocation, mouseYLocation, _, _ := GetMouseStatus()
	return getCellIdByLayerEntry(&commonResource.screenLayer, mouseXLocation, mouseYLocation)
}

/*
getCellIdByLayerAlias is a method which allows you to obtain a cell ID from a given text layer by layer alias.

:param layerAlias: The alias of the layer.
:param mouseXLocation: The x-axis location.
:param mouseYLocation: The y-axis location.

:return: The cell ID at the specified location.

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

:param layerEntry: A pointer to the layer entry.
:param xLocation: The x-axis location.
:param yLocation: The y-axis location.

:return: The cell ID at the specified location.

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

:param isRefreshForced: Whether to force a full refresh of the screen.

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

:param rootLayerEntry: A pointer to the root layer entry.
:param sortedLayerAliasSlice: A list of layer aliases sorted by z-order.

:return: The rendered root layer entry.

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

:param currentLayerEntry: The layer entry to render controls for.

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

:param sourceLayerAlias: The alias of the source layer.
:param targetLayerEntry: A pointer to the target layer entry.

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

:param sourceCharacterMemory: The source memory to copy from.
:param targetCharacterMemory: The target memory to copy to.
:param xLocation: The x-axis location in the target.
:param yLocation: The y-axis location in the target.
:param width: The width to copy.
:param height: The height to copy.

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
	var waitGroup sync.WaitGroup
	waitGroup.Add(height)
	for currentRow := 0; currentRow < height; currentRow++ {
		go func(row int) {
			defer waitGroup.Done()
			for currentColumn := 0; currentColumn < width; currentColumn++ {
				targetColumn := xLocation + currentColumn
				targetRow := yLocation + row
				if currentColumn < sourceWidth && row < sourceHeight && targetColumn < targetWidth && targetRow < targetHeight {
					targetCharacterMemory[targetRow][targetColumn] = sourceCharacterMemory[row][currentColumn]
				}
			}
		}(currentRow)
	}
	waitGroup.Wait()
}

/*
overlayLayers is a method which allows you to overlay one text layer on top of another text layer. In addition, the
following should be noted:

- If the source rune to be drawn is null, then it will be considered transparent.

- If a transparent rune has a foreground or background alpha value set, then it will be drawn as a shadow.

:param sourceLayerEntry: A pointer to the source layer entry.
:param targetLayerEntry: A pointer to the target layer entry.
:param isOpaque: Whether the source layer should be treated as opaque.

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

	// 3. Parallel Processing with Goroutines
	var wg sync.WaitGroup
	wg.Add(heightToCopy)

	for currentRow := 0; currentRow < heightToCopy; currentRow++ {
		go func(row int) {
			defer wg.Done()
			sourceRow := row + sourceStartY
			targetRow := row + targetStartY

			for currentColumn := 0; currentColumn < widthToCopy; currentColumn++ {
				sourceCol := currentColumn + sourceStartX
				targetCol := currentColumn + targetStartX

				sourceCharacterEntry := &sourceCharacterMemory[sourceRow][sourceCol]
				targetCharacterEntry := &targetCharacterMemory[targetRow][targetCol]
				sourceAttributeEntry := sourceCharacterEntry.AttributeEntry
				targetAttributeEntry := targetCharacterEntry.AttributeEntry

				// Handle NullRune (transparent) cells
				if sourceCharacterEntry.Character == constants.NullRune {
					if sourceAttributeEntry.CellType == constants.CellTypeShadow && sourceAttributeEntry.ForegroundAlphaValue < 1 {
						targetAttributeEntry.ForegroundColor = GetTransitionedColor(targetAttributeEntry.ForegroundColor, GetRGBColor(0, 0, 0), sourceAttributeEntry.ForegroundAlphaValue)
					}
					if sourceAttributeEntry.CellType == constants.CellTypeShadow && sourceAttributeEntry.BackgroundAlphaValue < 1 {
						targetAttributeEntry.BackgroundColor = GetTransitionedColor(targetAttributeEntry.BackgroundColor, GetRGBColor(0, 0, 0), sourceAttributeEntry.BackgroundAlphaValue)
					}
					targetCharacterEntry.AttributeEntry = targetAttributeEntry

					if sourceAttributeEntry.CellType == constants.CellTypeTooltip {
						targetCharacterEntry.AttributeEntry.CellType = constants.CellTypeTooltip
						targetCharacterEntry.AttributeEntry.CellControlAlias = sourceAttributeEntry.CellControlAlias
						targetCharacterEntry.LayerAlias = sourceCharacterEntry.LayerAlias
					}
					if sourceAttributeEntry.CellType == constants.CellTypeShadow {
						targetCharacterEntry.AttributeEntry.CellType = constants.CellTypeShadow
					}
					targetCharacterMemory[targetRow][targetCol] = *targetCharacterEntry
					continue
				}

				// Copy layer and parent aliases
				if sourceCharacterEntry.LayerAlias != "" {
					targetCharacterEntry.LayerAlias = sourceCharacterEntry.LayerAlias
				}
				if sourceCharacterEntry.ParentAlias != "" {
					targetCharacterEntry.ParentAlias = sourceCharacterEntry.ParentAlias
				}

				newAttributeEntry := types.NewAttributeEntry(&sourceAttributeEntry)
				targetCharacterEntry.Character = sourceCharacterEntry.Character

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
				if sourceAttributeEntry.ForegroundAlphaValue < 1 {
					newAttributeEntry.ForegroundColor = GetTransitionedColor(targetAttributeEntry.ForegroundColor, sourceAttributeEntry.ForegroundColor, sourceAttributeEntry.ForegroundAlphaValue)
				} else if defaultFgAlpha < 1 {
					newAttributeEntry.ForegroundColor = GetTransitionedColor(targetAttributeEntry.ForegroundColor, sourceAttributeEntry.ForegroundColor, defaultFgAlpha)
				}

				if sourceAttributeEntry.BackgroundAlphaValue < 1 {
					newAttributeEntry.BackgroundColor = GetTransitionedColor(targetAttributeEntry.BackgroundColor, sourceAttributeEntry.BackgroundColor, sourceAttributeEntry.BackgroundAlphaValue)
				} else if defaultBgAlpha < 1 {
					newAttributeEntry.BackgroundColor = GetTransitionedColor(targetAttributeEntry.BackgroundColor, sourceAttributeEntry.BackgroundColor, defaultBgAlpha)
				}

				targetCharacterEntry.AttributeEntry = newAttributeEntry
				targetCharacterMemory[targetRow][targetCol] = *targetCharacterEntry
			}
		}(currentRow)
	}
	wg.Wait()
}

/*
DrawLayerToScreen is a method which allows you to render a text layer to the visible terminal screen. In addition, the
following should be noted:

- If debug is enabled, this method does nothing since the terminal is virtual.

:param layerEntry: A pointer to the layer entry to draw.
:param isForcedRefreshRequired: Whether to force a full refresh of the screen.

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

:return: A string representing the OS.

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
