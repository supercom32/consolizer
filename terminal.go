package consolizer

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/internal/math"
	"github.com/supercom32/consolizer/internal/memory"
	"github.com/supercom32/consolizer/internal/stringformat"
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
defaultValueType is a structure that holds common information about the
current terminal session that needs to be shared.
*/
type defaultValueType struct {
	screen               tcell.Screen
	layerInstance        LayerInstanceType // What happens when last layer is deleted? This needs to be updated.
	terminalWidth        int
	terminalHeight       int
	screenLayer          types.LayerEntryType
	debugDirectory       string
	isDebugEnabled       bool
	displayUpdate        sync.Mutex
	updateDisplayChannel chan bool
}

/*
commonResource is a variable used to hold shared data that is accessed
by this package.
*/
var commonResource defaultValueType

func Test() {

}

/*
InitializeTerminal allows you to initialize consolizer for the first time.
This method must be called first before any operations take place. The
parameters 'width' and 'height' represent the display size of the
terminal instance you wish to create. In addition, the following
information should be noted:

- If you pass in a zero or negative value for ether width or height a panic
will be generated to fail as fast as possible.
*/
func InitializeTerminal(width int, height int) {
	memory.InitializeScreenMemory()
	memory.InitializeImageMemory()
	memory.InitializeTimerMemory()
	// Set the mouse location off screen so it won't trigger events at 0,0 which the user never moved to.
	memory.SetMouseStatus(-1, -1, 0, "")
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
		go setupEventUpdater()
		go setupPeriodicEventUpdater()
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
}

func setupPeriodicEventUpdater() {
	for {
		UpdatePeriodicEvents()
	}
}

/*
setupEventUpdater is a background method that monitors all events coming
into the terminal session. When an event is detected, it is recorded and
monitoring continues.
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
setupCloseHandler enables the trapping of all unexpected system calls and shuts
down the terminal gracefully. This means all terminal settings should be reset
back to normal if anything unexpected happens to the user or if the process is
killed.
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
RestoreTerminalSettings allows the user to gracefully return the terminal
back to its normal settings. This should be called once your application
is finished using consolizer so that the users terminal environment is not
left in a bad state.
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
GetTerminalSize allows you to obtain width and height of the current terminal
characters.
*/
func GetTerminalSize() (int, int) {
	return commonResource.screen.Size()
}

/*
Inkey allows you to read keyboard input from the user's terminal. This
method returns the character pressed or a keyword representing the
special key pressed (For example: 'a', 'A', 'escape', 'f10', etc.).
In addition, the following information should be noted:

- If more than one keystroke is recorded, it is stored sequentially
in the input buffer and this method needs to be called repeatedly in
order to read them.
*/
func Inkey() []rune {
	return memory.KeyboardMemory.GetKeystrokeFromKeyboardBuffer()
}

/*
Layer allows you to specify a default layer alias that you wish to use when
interacting with methods which have a non-layer alias method signature.
Non-layer alias method signatures can be identified by finding methods which
have both a layer and non-layer version. This makes interacting with methods
faster, as the user does not need to provide the layer alias context in which
he is working on. For example:

	// On the layer with a layer alias of "MyForegroundLayer" print a string.
	consolizer.PrintLayer("MyForegroundLayer", "Hello World")

	// Set the default layer alias to be "MyForegroundLayer".
	consolizer.Layer("MyForegroundLayer")

	// Since we set the default layer, we don't need to call the method
	// PrintLayer anymore. Instead, we can use the shorter method Print.
	consolizer.Print("Hello World")
*/
func Layer(layerInstance LayerInstanceType) {
	validateLayer(layerInstance.layerAlias)
	commonResource.layerInstance = layerInstance
}

/*
AddTextStyle allows you to add a new printing style which can be used for
writing dialog text with. Dialog text printing differs from regular
printing, since it allows for "typewriter" drawing effects as well as
changing text attributes on the fly (color, isBold, etc). In addition,
the following information should be noted:

- This method expects you to pass in a 'textStyleEntry' type obtained by
calling 'consolizer.NewTextStyle()'. This entry type contains all the
options available for your text style and can be configured easily
by setting each attribute accordingly.
*/
func AddTextStyle(textStyleAlias string, textStyleEntry types.TextCellStyleEntryType) {
	memory.AddTextStyle(textStyleAlias, textStyleEntry)
}

/*
DeleteTextStyle allows you to remove a text style that was added previously.
In addition, the following information should be noted:

- If you attempt to delete an entry that does not exist, then no operation
will be performed.
*/
func DeleteTextStyle(textStyleAlias string) {
	validateTextStyleExists(textStyleAlias)
	memory.DeleteTextStyle(textStyleAlias)
}

/*
NewTextStyle allows you to obtain a new text style entry which can be
used when printing dialog text. By configuring attributes for your text
style entry and adding your entry to consolizer, simple markup commands can
be used to switch between dialog printing styles automatically.
For example:

	// Create a new text style entry to configure.
	myTextStyleEntry := consolizer.NewTextStyle()
	// Configure your text style so that the foreground color is red.
	myTextStyleEntry.ForegroundColor = consolizer.GetRGBColor(255, 0, 0)
	// Add a new text style called "RedColor".
	consolizer.AddTextStyle("RedColor", myTextStyleEntry)
*/
func NewTextStyle() types.TextCellStyleEntryType {
	return types.NewTextCellStyleEntry()
}

func NewImageStyle() types.ImageStyleEntryType {
	return types.NewImageStyleEntry()
}

/*
NewTuiStyleEntry allows you to obtain a new style entry which can be used
for specifying how TUI controls and other TUI drawing operations should
occur. For example:

	// Create a new TUI style entry to configure.
	myTuiStyleEntry := consolizer.NewTuiStyleEntry()
	// Configure the style entry so that the upper left corner character
	// for drawing a window is the '╔' character. Here we use '\u' to
	// denote the specific byte value in our code since not all editors may
	// know how to display the actual character.
	myTuiStyleEntry.UpperLeftCorner = '\u2554'
	// Draw a window on the text layer with the alias of "ForegroundLayer",
	// using the TUI style entry "myTuiStyleEntry", at layer location (0, 0),
	// with a width and height of 10x10 characters.
	consolizer.DrawWindow("ForegroundLayer", myTuiStyleEntry, 0, 0, 10, 10)
*/
func NewTuiStyleEntry() types.TuiStyleEntryType {
	return types.NewTuiStyleEntry()
}

/*
NewSelectionEntry allows you to obtain an entry used for specifying what
options you want to make available for a given menu prompt. For example:

	// Create a new TUI style entry with default settings.
	tuiStyleEntry := consolizer.NewTuiStyleEntry()
	// Create a new selection entry to populate our menu entries with.
	selectionEntry := consolizer.NewSelectionEntry()
	// Add a selection with the alias "Opt1" and a display value of "OK".
	selectionEntry.Add("Opt1", "OK")
	// Add a selection with the alias "Opt2" with the display value of "CANCEL".
	selectionEntry.Add("Opt2", "CANCEL")
	// Prompt the user with a vertical selection menu, on the text layer
	// with the alias "ForegroundLayer", using a default TUI style entry,
	// a selection entry with two options, at the layer location (0, 0),
	// with a menu width and height of 15x15 characters.
	selectionMade := consolizer.GetSelectionFromVerticalMenu ("ForegroundLayer", tuiStyleEntry, selectionEntry, 0, 0, 15, 15)
*/
func NewSelectionEntry() types.SelectionEntryType {
	return types.NewSelectionEntry()
}

/*
NewAssetList allows you to obtain a list for storing asset information with.
This is useful for loading assets in bulk since you can specify asset
information as a collection instead of each one individually. An example
use of this method is as follows:

	// Create a new asset list.
	assetList := consolizer.NewAssetList()
	// Add an image file with the image filename 'MyImageFile', and image
	// alias of 'MyImageAlias'.
	assetList.AddImage("MyImageFile", "MyImageAlias")
	// Load the list of images into memory.
	err := loadImagesInBulk(assetList)

In addition, the following information should be noted:

- An asset list can contain multiple asset types. This allows the same
asset list to be shared by multiple methods that load different kinds of
assets.
*/
func NewAssetList() types.AssetListType {
	return types.NewAssetList()
}

/*
SetLayerZOrder allows you to set the zOrder value for a given layer. In
addition, the following information should be noted:

- The z order priority controls which text layer should be drawn first and
which text layer should be drawn last. Layers that have a higher priority
will be drawn on top of layers that have a lower priority. In the event
that two layers have the same priority, they will be drawn in random order.
This is to ensure that programmers do not attempt to rely on any specific
behavior that might be a coincidental side effect.
*/
func SetLayerZOrder(layerInstance LayerInstanceType, zOrder int) {
	layerEntry := memory.GetLayer(layerInstance.layerAlias)
	layerEntry.ZOrder = zOrder
}

/*
SetZOrder allows you to set the zOrder value for default currently selected
layer. In addition, the following information should be noted:

- The z order priority controls which text layer should be drawn first and
which text layer should be drawn last. Layers that have a higher priority
will be drawn on top of layers that have a lower priority. In the event
that two layers have the same priority, they will be drawn in random order.
This is to ensure that programmers do not attempt to rely on any specific
behavior that might be a coincidental side effect.
*/
func SetZOrder(zOrder int) {
	SetLayerZOrder(commonResource.layerInstance, zOrder)
}

/*
SetLayerAlpha allows you to set the alpha value for a given text layer. This lets
you perform pseudo transparencies by making the layer foreground and background
colors blend with the layers underneath it to the degree specified. In
addition, the following information should be noted:

- An alpha value of 1.0 is equal to 100% visible, while an alpha value of
0.0 is 0% visible. Specifying a value outside this range indicates that
you want to over amplify or under amplify the color transparency effect.

- If the percent change specified is outside of the RGB color range (for
example, if you specified 200%), then the color will simply bottom or max
out at RGB(0, 0, 0) or RGB(255, 255, 255) respectively.
*/
func SetLayerAlpha(layerInstance LayerInstanceType, alphaValue float32) {
	setLayerAlpha(layerInstance, alphaValue)
}

func setLayerAlpha(layerInstance LayerInstanceType, alphaValue float32) {
	layerEntry := memory.GetLayer(layerInstance.layerAlias)
	layerEntry.DefaultAttribute.ForegroundTransformValue = alphaValue
	layerEntry.DefaultAttribute.BackgroundTransformValue = alphaValue
}

func SetLayer(alphaValue float32) {
	setLayerAlpha(commonResource.layerInstance, alphaValue)
}

/*
GetColor allows you to obtain an RGB color from a predefined color palette
list. This list corresponds to the 16 color ANSI standard, where color
0 is Black and 15 is Bright White.  In addition, the following information
should be noted:

- If you specify a color index less than 0 or greater than 15 a panic
will be generated to fail as fast as possible.
*/
func GetColor(colorIndex int) constants.ColorType {
	validateColorIndex(colorIndex)
	return constants.AnsiColorByIndex[colorIndex]
}

/*
GetRGBColor allows you to obtain a specific RGB color based on the red, green, and
blue index values provided. In addition, the following information should be noted:

- If you specify a color channel index less than 0 or greater than 255 a panic
will be generated to fail as fast as possible.
*/
func GetRGBColor(redColorIndex int32, greenColorIndex int32, blueColorIndex int32) constants.ColorType {
	validateRGBColorIndex(redColorIndex, greenColorIndex, blueColorIndex)
	return constants.ColorType(tcell.NewRGBColor(redColorIndex, greenColorIndex, blueColorIndex))
}

/*
Color allows you to set default colors on your text layer for printing with.
The color index specified corresponds to the 16 color ANSI standard, where
color 0 is Black and 15 is Bright White. If you wish to change colors settings
for a text layer that is not currently set as your default, use 'ColorLayer'
instead.
*/
func Color(foregroundColorIndex int, backgroundColorIndex int) {
	validateDefaultLayerIsNotEmpty()
	ColorLayer(commonResource.layerInstance, foregroundColorIndex, backgroundColorIndex)
}

/*
ColorLayer allows you to set default colors on your specified text layer for
printing with. The color index specified corresponds to the 16 color ANSI
standard, where color 0 is Black and 15 is Bright White. If you do not wish
to specify a text layer, you can use the method 'Color' which will simply
change the color for the default text layer previously set.
*/
func ColorLayer(layerInstance LayerInstanceType, foregroundColorIndex int, backgroundColorIndex int) {
	validateColorIndex(foregroundColorIndex)
	validateColorIndex(backgroundColorIndex)
	layerEntry := memory.GetLayer(layerInstance.layerAlias)
	layerEntry.DefaultAttribute.ForegroundColor = constants.AnsiColorByIndex[foregroundColorIndex]
	layerEntry.DefaultAttribute.BackgroundColor = constants.AnsiColorByIndex[backgroundColorIndex]
}

/*
ColorRGB allows you to set default colors on your text layer for printing with.
This method allows you to specify colors using RGB color index values within
the range of 0 to 255. If you wish to change colors settings for a text layer
that is not currently set as your default, use 'ColorLayerRGB' instead.
*/
func ColorRGB(foregroundRedIndex int32, foregroundGreenIndex int32, foregroundBlueIndex int32, backgroundRedIndex int32, backgroundGreenIndex int32, backgroundBlueIndex int32) {
	validateDefaultLayerIsNotEmpty()
	ColorLayerRGB(commonResource.layerInstance, foregroundRedIndex, foregroundGreenIndex, foregroundBlueIndex, backgroundRedIndex, backgroundGreenIndex, backgroundBlueIndex)
}

/*
ColorLayerRGB allows you to set default colors on your specified text layer
for printing with. This method allows you to specify colors using RGB color
index values within the range of 0 to 255. If you do not wish to specify a
text layer, you can use the method 'ColorRGB' which will simply change the
color for the default text layer previously set.
*/
func ColorLayerRGB(layerInstance LayerInstanceType, foregroundRed int32, foregroundGreen int32, foregroundBlue int32, backgroundRed int32, backgroundGreen int32, backgroundBlue int32) {
	foregroundColor := GetRGBColor(foregroundRed, foregroundGreen, foregroundBlue)
	backgroundColor := GetRGBColor(backgroundRed, backgroundGreen, backgroundBlue)
	ColorLayer24Bit(layerInstance, foregroundColor, backgroundColor)
}

/*
Color24Bit allows you to color a layer using a 24-bit color expressed as
an int32. This is useful for when you have colors which are already defined.
*/

func Color24Bit(foregroundColor constants.ColorType, backgroundColor constants.ColorType) {
	ColorLayer24Bit(commonResource.layerInstance, foregroundColor, backgroundColor)
}

/*
ColorLayer24Bit allows you to color a layer using a 24-bit color expressed as
an int32. This is useful for internal methods that already have a 24-bit color
and do not require to compute it again.
*/
func ColorLayer24Bit(layerInstance LayerInstanceType, foregroundColor constants.ColorType, backgroundColor constants.ColorType) {
	layerEntry := memory.GetLayer(layerInstance.layerAlias)
	layerEntry.DefaultAttribute.ForegroundColor = foregroundColor
	layerEntry.DefaultAttribute.BackgroundColor = backgroundColor
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
func Locate(xLocation int, yLocation int) {
	validateDefaultLayerIsNotEmpty()
	LocateLayer(commonResource.layerInstance, xLocation, yLocation)
}

/*
LocateLayer allows you to set the default cursor location on your specified text
layer for printing with. This is useful for when you wish to print text
at different locations of your text layer at any given time. If you do not
wish to specify a text layer, you can use the method 'Locate' which will
simply change the cursor location for the default text layer previously set.
In addition, the following information should be noted:

- If you pass in a location value that falls outside the dimensions of the
specified text layer, a panic will be generated to fail as fast as possible.

- Valid text layer locations start at position (0,0) for the upper left corner.
Since location values do not start at (1,1), valid end positions for the bottom
right corner will be one less than the text layer width and height. For
example:

	// Create a new text layer with the alias "ForegroundLayer", at location
	// (0,0), with a width and height of 15x15, a z order priority of 1,
	// and no parent layer associated with it.
	consolizer.AddLayer("ForegroundLayer", 0, 0, 15, 15, 1, "")
	// Move our cursor location to the bottom right corner of our text layer.
	consolizer.LocateLayer(14, 14)
*/
func LocateLayer(layerInstance LayerInstanceType, xLocation int, yLocation int) {
	validateLayer(layerInstance.layerAlias)
	layerEntry := memory.GetLayer(layerInstance.layerAlias)
	validateLayerLocationByLayerEntry(layerEntry, xLocation, yLocation)
	layerEntry.CursorXLocation = xLocation
	layerEntry.CursorYLocation = yLocation
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
func Print(textToPrint string) {
	validateDefaultLayerIsNotEmpty()
	PrintLayer(commonResource.layerInstance, textToPrint)
}

/*
PrintLayer allows you to write text to a specified text layer. If you do not
wish to specify a text layer, you can use the method 'Print' which will
simply print to the default text layer previously set. In addition, the
following information should be noted:

- When text is written to the text layer, the cursor position is also updated
to reflect its new location. Like a typewriter, the cursor position moves to
the start of the next line after each print statement.

- If the string to print ends up being too long to fit at its current location,
then only the visible portion of your string will be printed.

- If printing has not yet finished and there are no available lines left, then
all remaining characters will be discarded and printing will stop.
*/
func PrintLayer(layerInstance LayerInstanceType, textToPrint string) {
	layerEntry := memory.GetLayer(layerInstance.layerAlias)
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
printLayer allows you to write text to a text layer. This is useful
for internal methods that want to write text to a text layer directly, without
effecting user settings (such as current cursor location, etc). In addition,
the following information should be noted:

- If the location to print falls outside the range of the text layer,
then only the visible portion of your text will be printed.
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

func printLayerWithWordWrap(layerEntry *types.LayerEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int, width int, textToPrint []rune) int {
	layerWidth := layerEntry.Width
	layerHeight := layerEntry.Height
	cursorXLocation := xLocation
	cursorYLocation := yLocation
	characterMemory := layerEntry.CharacterMemory
	for currentCharacterIndex, currentCharacter := range textToPrint {
		if currentCharacter == ' ' {
			// Check if the word fits within the remaining space on the current line.
			wordWidth := calculateWordWidth(textToPrint, currentCharacterIndex)
			if cursorXLocation+wordWidth >= xLocation+width {
				// Word doesn't fit, move to the next line.
				cursorXLocation = xLocation
				cursorYLocation++
			}
		}
		if currentCharacter == ' ' && cursorXLocation == xLocation {
			// Skip the first blank space at the start of a line if one exists.
			continue
		}
		if cursorXLocation >= 0 && cursorXLocation < layerWidth && cursorYLocation >= 0 && cursorYLocation < layerHeight {
			originalBackgroundColor := characterMemory[cursorYLocation][cursorXLocation].AttributeEntry.BackgroundColor
			characterMemory[cursorYLocation][cursorXLocation].AttributeEntry = types.NewAttributeEntry(&attributeEntry)
			characterMemory[cursorYLocation][cursorXLocation].Character = currentCharacter
			if stringformat.IsRuneCharacterWide(currentCharacter) {
				cursorXLocation++
				if cursorXLocation >= layerWidth {
					continue
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
			continue
		}
	}
	return cursorXLocation - xLocation
}

func calculateWordWidth(textToPrint []rune, start int) int {
	// Calculate the width of a word from the given position. The first position is
	// always assumed to be ' ' and is skipped.
	wordWidth := 0
	for i := start + 1; i < len(textToPrint); i++ {
		if textToPrint[i] == ' ' {
			break
		}
		wordWidth++
	}
	return wordWidth
}

/*
Clear allows you to empty the default text layer of all its contents. If you
wish to clear a text layer that is not currently set as the default, use
'ClearLayer' instead.
*/
func Clear() {
	validateDefaultLayerIsNotEmpty()
	ClearLayer(commonResource.layerInstance)
}

/*
ClearLayer allows you to empty the specified text layer of all its contents. If you
do not wish to specify a text layer, you can use the method 'Clear' which will
simply clear the default text layer previously set.
*/
func ClearLayer(layerInstance LayerInstanceType) {
	layerEntry := memory.GetLayer(layerInstance.layerAlias)
	clearLayer(layerEntry)
}

/*
clearLayer allows you to empty the specified text layer of all its contents.
This is useful for internal methods that want to clear a text layer directly.
*/
func clearLayer(layerEntry *types.LayerEntryType) {
	types.InitializeCharacterMemory(layerEntry)
}

func GetCharacterOnScreen(xLocation int, yLocation int) rune {
	layerEntry := commonResource.screenLayer
	validateLayerLocationByLayerEntry(&layerEntry, xLocation, yLocation)
	return layerEntry.CharacterMemory[xLocation][yLocation].Character
}

/*
scrollCharacterMemory allows you to advance the specified text layer up by one
row. This means that the first row is discarded and all subsequent rows are
moved up by one position. The new row created at the bottom of the text layer
will be filled with spaces (" ") colored with the layers default attributes.
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
getRuneOnLayer allows you to obtain a specific rune at the location specified
on the given text layer. In addition, the following information should be
noted:

- If the location specified is outside the valid range of the text layer, then
a panic will be thrown to fail as fast as possible.
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
GetCellIdUnderMouseLocation allows you to obtain the cell ID for the text
directly under your mouse cursor. This is useful for tracking
elements on a screen, creating "hot spots", or interactive zones which you
want the user to interact with. In addition, the following information should be
noted:

- If multiple text layers are being displayed, the cell ID returned will be
from the top-most visible text cell.

- The cell ID returned will only reflect what is currently being displayed
on the terminal display. If you wish for any new changes to take effect,
call 'UpdateDisplay' to refresh the visible display area first.
*/
func GetCellIdUnderMouseLocation() int {
	mouseXLocation, mouseYLocation, _, _ := memory.GetMouseStatus()
	return getCellIdByLayerEntry(&commonResource.screenLayer, mouseXLocation, mouseYLocation)
}

/*
getCellIdByLayerAlias allows you to obtain a cell ID from a given text layer
by layer alias. This is simply a wrapper method that converts the text
layer alias into a layer entry and calls 'getCellIdByLayerEntry'.
*/
func getCellIdByLayerAlias(layerAlias string, mouseXLocation int, mouseYLocation int) int {
	validateLayer(layerAlias)
	layerEntry := memory.GetLayer(layerAlias)
	return getCellIdByLayerEntry(layerEntry, mouseXLocation, mouseYLocation)
}

/*
getCellIdByLayerEntry allows you to obtain a cell ID from a given text layer
by layer entry. In addition, the following information should be noted:

- If the location specified is outside the valid range of the text layer, then
a value of '-1' is returned instead.
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
UpdateDisplay allows you to synchronize the terminals visible display area with
your current changes. In addition, the following information should be noted:

- All text layers are sorted from lowest to highest z order priority level.

- Layers with the same z order priority will appear in random display order.
This is to ensure that programmers do not attempt to rely on any specific
behavior that might be a coincidental side effect.
*/
func UpdateDisplay(isRefreshForced bool) {
	commonResource.displayUpdate.Lock()
	defer func() {
		commonResource.displayUpdate.Unlock()
	}()
	sortedLayerAliasSlice := memory.GetSortedLayerMemoryAliasSlice()
	baseLayerEntry := types.NewLayerEntry("", "", commonResource.terminalWidth, commonResource.terminalHeight)
	baseLayerEntry = renderLayers(&baseLayerEntry, sortedLayerAliasSlice)
	DrawLayerToScreen(&baseLayerEntry, isRefreshForced)
	commonResource.screenLayer = baseLayerEntry
}

func RefreshDisplay() {
	commonResource.screen.Sync()
}

/*
renderLayers allows you to render a list of text layers to the specified root
text layer. In addition, the following information should be noted:

- The root layer entry is considered the parent entry. Only text layers under
the specified parent will be rendered on it.

- If a text layer being rendered is a parent, then all child text layers will
be rendered on the parent before the parent is drawn. This is done by making
a recursive call to 'renderLayers' with the new parent layer.

- The list of text layers alias provided should be sorted so that layers with
a lower z order priority are rendered first.

- Any text layer which is marked as not visible will be ignored.

- All rendering occurs on a temporary text layer until it is ready to be
overlaid on the final (terminal ready) text layer. Button and other special
TUI controls are also dynamically rendered at this time so that the original
text layer data underneath them is preserved.
*/
func renderLayers(rootLayerEntry *types.LayerEntryType, sortedLayerAliasSlice memory.LayerAliasZOrderPairList) types.LayerEntryType {
	baseLayerEntry := types.NewLayerEntry("", "", 0, 0, rootLayerEntry)
	for currentListIndex := 0; currentListIndex < len(sortedLayerAliasSlice); currentListIndex++ {
		if !memory.IsLayerExists(sortedLayerAliasSlice[currentListIndex].Key) {
			continue
		}
		currentLayerEntry := types.NewLayerEntry("", "", 0, 0, memory.GetLayer(sortedLayerAliasSlice[currentListIndex].Key))
		if currentLayerEntry.IsVisible {
			renderControls(currentLayerEntry)
			if currentLayerEntry.IsParent && (currentLayerEntry.LayerAlias != baseLayerEntry.LayerAlias && currentLayerEntry.ParentAlias == baseLayerEntry.LayerAlias) {
				renderedLayer := renderLayers(&currentLayerEntry, sortedLayerAliasSlice)
				overlayLayers(&renderedLayer, &baseLayerEntry)
			} else {
				if currentLayerEntry.ParentAlias == baseLayerEntry.LayerAlias {
					overlayLayers(&currentLayerEntry, &baseLayerEntry)
				}
			}
		}
	}
	return baseLayerEntry
}

/*
renderControls allows you to draw various control elements on the specified layer.
The order of drawing matters, as complex controls are drawn first above basic controls
to ensure that any pop-up controls appear over complex controls.
*/
func renderControls(currentLayerEntry types.LayerEntryType) {
	Button.drawButtonsOnLayer(currentLayerEntry)
	TextField.drawTextFieldOnLayer(currentLayerEntry)
	Checkbox.drawCheckboxesOnLayer(currentLayerEntry)
	Dropdown.drawDropdownsOnLayer(currentLayerEntry)
	Selector.drawSelectorsOnLayer(currentLayerEntry)
	scrollbar.drawScrollbarsOnLayer(currentLayerEntry)
	textbox.drawTextboxesOnLayer(currentLayerEntry)
	radioButton.drawRadioButtonsOnLayer(currentLayerEntry)
	ProgressBar.drawProgressBarsOnLayer(currentLayerEntry)
	Label.drawLabelsOnLayer(currentLayerEntry)
	Tooltip.drawTooltipsOnLayer(currentLayerEntry)
}

/*
overlayLayersByLayerAlias allows you to overlay a text layer by its layer
alias. This is useful when you do not have actual layer data and only
know the alias of the layer you wish to overlay.
*/
func overlayLayersByLayerAlias(sourceLayerAlias string, targetLayerEntry *types.LayerEntryType) {
	validateLayer(sourceLayerAlias)
	layerEntry := memory.GetLayer(sourceLayerAlias)
	overlayLayers(layerEntry, targetLayerEntry)
}

/*
overlayLayers allows you to overlay one text layer on top of another text
layer. In addition, the following information should be noted:

- If the source text layer is set to be drawn outside the target layer,
then only the visible portion of the source text layer will be rendered.

- If the source text layer is found to be completely outside the range
of the target layer, then no rendering will occur.

- If the source rune to be drawn is null, then it will be considered
transparent.

- If a transparent rune has a foreground or background alpha value set,
then it will be drawn as a shadow with the color and intensity matching
the rune underneath it.
*/
func overlayLayers(sourceLayerEntry *types.LayerEntryType, targetLayerEntry *types.LayerEntryType) {
	sourceCharacterMemory := sourceLayerEntry.CharacterMemory
	targetCharacterMemory := targetLayerEntry.CharacterMemory
	sourceWidthToCopy := sourceLayerEntry.Width
	sourceHeightToCopy := sourceLayerEntry.Height
	// Calculate how much of the source Width to copy.
	sourceWidthToCopy = sourceLayerEntry.Width - int(math.GetAbsoluteValueAsFloat64(sourceLayerEntry.ScreenXLocation))
	if sourceLayerEntry.ScreenXLocation < 0 {
		if sourceWidthToCopy > targetLayerEntry.Width {
			sourceWidthToCopy = targetLayerEntry.Width
		}
	} else {
		if sourceWidthToCopy < targetLayerEntry.Width {
			sourceWidthToCopy = sourceLayerEntry.Width
		}
	}
	// Calculate how much of the source Length to copy.
	sourceHeightToCopy = sourceLayerEntry.Height - int(math.GetAbsoluteValueAsFloat64(sourceLayerEntry.ScreenYLocation))
	if sourceLayerEntry.ScreenYLocation < 0 {
		if sourceHeightToCopy > targetLayerEntry.Height {
			sourceHeightToCopy = targetLayerEntry.Height
		}
	} else {
		if sourceHeightToCopy < targetLayerEntry.Height {
			sourceHeightToCopy = sourceLayerEntry.Height
		}
	}
	// Adjust where rendering on the layer should start.
	startingSourceXLocation := 0
	startingSourceYLocation := 0
	startingTargetXLocation := 0
	startingTargetYLocation := 0
	if sourceLayerEntry.ScreenXLocation < 0 {
		startingSourceXLocation = int(math.GetAbsoluteValueAsFloat64(sourceLayerEntry.ScreenXLocation))
	} else {
		startingTargetXLocation = int(math.GetAbsoluteValueAsFloat64(sourceLayerEntry.ScreenXLocation))
	}
	if sourceLayerEntry.ScreenYLocation < 0 {
		startingSourceYLocation = int(math.GetAbsoluteValueAsFloat64(sourceLayerEntry.ScreenYLocation))
	} else {
		startingTargetYLocation = int(math.GetAbsoluteValueAsFloat64(sourceLayerEntry.ScreenYLocation))
	}
	if sourceWidthToCopy+startingTargetXLocation > targetLayerEntry.Width {
		sourceWidthToCopy = targetLayerEntry.Width - startingTargetXLocation
	}
	if sourceHeightToCopy+startingTargetYLocation > targetLayerEntry.Height {
		sourceHeightToCopy = targetLayerEntry.Height - startingTargetYLocation
	}
	// If the layer is totally off screen, don't bother to render it.
	if startingSourceXLocation+sourceWidthToCopy < 0 || sourceLayerEntry.ScreenXLocation+sourceWidthToCopy > targetLayerEntry.Width ||
		startingSourceYLocation+sourceHeightToCopy < 0 || sourceLayerEntry.ScreenYLocation+sourceHeightToCopy > targetLayerEntry.Height {
		return
	}
	// Perform the actual copy using the starting offsets previously calculated.
	for currentRow := 0; currentRow < sourceHeightToCopy; currentRow++ {
		for currentColumn := 0; currentColumn < sourceWidthToCopy; currentColumn++ {
			sourceCharacterEntry := &sourceCharacterMemory[currentRow+startingSourceYLocation][currentColumn+startingSourceXLocation]
			targetCharacterEntry := &targetCharacterMemory[currentRow+startingTargetYLocation][currentColumn+startingTargetXLocation]
			sourceAttributeEntry := sourceCharacterEntry.AttributeEntry
			targetAttributeEntry := targetCharacterEntry.AttributeEntry

			// Handle transformations
			if sourceCharacterEntry.Character == constants.NullRune {
				if sourceAttributeEntry.ForegroundTransformValue < 1 {
					targetAttributeEntry.ForegroundColor = GetTransitionedColor(targetAttributeEntry.ForegroundColor, GetRGBColor(0, 0, 0), sourceAttributeEntry.ForegroundTransformValue)
				}
				if sourceAttributeEntry.BackgroundTransformValue < 1 {
					targetAttributeEntry.BackgroundColor = GetTransitionedColor(targetAttributeEntry.BackgroundColor, GetRGBColor(0, 0, 0), sourceAttributeEntry.BackgroundTransformValue)
				}
				targetCharacterEntry.AttributeEntry = targetAttributeEntry
				// Here if we detect that our transparent rune is a tooltip, we propagate it to the cell underneath it.
				if sourceCharacterEntry.AttributeEntry.CellType == constants.CellTypeTooltip {
					targetCharacterEntry.AttributeEntry.CellType = constants.CellTypeTooltip
					targetCharacterEntry.AttributeEntry.CellControlAlias = sourceCharacterEntry.AttributeEntry.CellControlAlias
					targetCharacterEntry.LayerAlias = sourceCharacterEntry.LayerAlias
				}
				targetCharacterMemory[currentRow+startingTargetYLocation][currentColumn+startingTargetXLocation] = *targetCharacterEntry
			} else {
				// If your copying a composed image, don't clobber the target layer alias
				if sourceCharacterEntry.LayerAlias != "" {
					targetCharacterEntry.LayerAlias = sourceCharacterEntry.LayerAlias
				}
				if sourceCharacterEntry.ParentAlias != "" {
					targetCharacterEntry.ParentAlias = sourceCharacterEntry.ParentAlias
				}
				targetCharacterEntry.AttributeEntry = types.NewAttributeEntry(&sourceAttributeEntry)
				targetCharacterEntry.Character = sourceCharacterEntry.Character
				// If there is no local color transforming being done on cells
				if sourceAttributeEntry.ForegroundTransformValue != 1 || sourceAttributeEntry.BackgroundTransformValue != 1 {
					if sourceAttributeEntry.ForegroundTransformValue < 1 {
						targetCharacterEntry.AttributeEntry.ForegroundColor = GetTransitionedColor(targetAttributeEntry.ForegroundColor, sourceAttributeEntry.ForegroundColor, sourceAttributeEntry.ForegroundTransformValue)
					}
					if sourceAttributeEntry.BackgroundTransformValue < 1 {
						targetCharacterEntry.AttributeEntry.BackgroundColor = GetTransitionedColor(targetAttributeEntry.BackgroundColor, sourceAttributeEntry.BackgroundColor, sourceAttributeEntry.BackgroundTransformValue)
					}
				} else {
					if sourceLayerEntry.DefaultAttribute.ForegroundTransformValue < 1 {
						targetCharacterEntry.AttributeEntry.ForegroundColor = GetTransitionedColor(targetAttributeEntry.ForegroundColor, sourceAttributeEntry.ForegroundColor, sourceLayerEntry.DefaultAttribute.ForegroundTransformValue)
					}
					if sourceLayerEntry.DefaultAttribute.BackgroundTransformValue < 1 {
						targetCharacterEntry.AttributeEntry.BackgroundColor = GetTransitionedColor(targetAttributeEntry.BackgroundColor, sourceAttributeEntry.BackgroundColor, sourceLayerEntry.DefaultAttribute.BackgroundTransformValue)
					}
				}
			}
		}
	}
}

/*
DrawLayerToScreen allows you to render a text layer to the visible terminal
screen. If debug is enabled, this method does nothing since the terminal
is virtual.
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
