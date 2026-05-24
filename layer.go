package consolizer

import (
	"fmt"
	"github.com/supercom32/consolizer/memory"
	"github.com/supercom32/consolizer/stringformat"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/types"
)

/*
layerType is a structure which provides the internal namespace for layer management operations.
*/
type layerType struct{}

/*
layer is a variable which is the global instance for managing layer operations.
*/
var layer layerType

/*
Layers is a variable which is the global instance for managing the memory of all layer entries.
*/
var Layers *memory.MemoryManager[types.LayerEntryType]

/*
LayerInstanceType is a structure which represents a handle to a layer instance.
*/
type LayerInstanceType struct {
	layerAlias  string
	parentAlias string
}

/*
layerAliasZOrderPair is a structure which represents a pairing between a layer alias and its z-order.
*/
type layerAliasZOrderPair struct {
	Key   string
	Value int
}

/*
LayerAliasZOrderPairList is a structure which represents a list of layer alias and z-order pairings.
*/
type LayerAliasZOrderPairList []layerAliasZOrderPair

/*
init is a method which initializes the layer management system by re-initializing the screen memory.

Example:

	init()
*/
func init() {
	layer.ReInitializeScreenMemory()
}

/*
clearLayerInstance is a method which allows you to clear the contents of a layer instance.

Example:

	clearLayerInstance(instance)
*/
func clearLayerInstance(layerInstance *LayerInstanceType) {
	layerEntry := Layers.Get(layerInstance.layerAlias)
	layer.clear(layerEntry)
}

/*
clear is a method which allows you to empty the specified text layer of all its contents. In addition, the following should be noted:

- This is useful for internal methods that want to clear a text layer directly.

Example:

	layer.clear(layerEntry)
*/
func (shared *layerType) clear(layerEntry *types.LayerEntryType) {
	types.InitializeCharacterMemory(layerEntry)
}

/*
processMarkupTag is a method which allows you to process a markup tag in the text and update the attribute entry
accordingly. In addition, the following should be noted:

- It returns the updated attribute entry and the new character index after the tag.

- If no valid closing tag is found, the opening tag is treated as regular text.

- Special tag "/" resets to the default attribute entry.

Example:

	attr, nextIndex := layer.processMarkupTag(runes, 5, "Hello {{red}}World", defaultAttr)
*/
func (shared *layerType) processMarkupTag(textToPrint []rune, currentIndex int, textString string, defaultAttributeEntry types.AttributeEntryType) (types.AttributeEntryType, int) {
	tagStartIndex := currentIndex + 2
	tagEndRelIndex := strings.Index(textString[tagStartIndex:], "}}")
	if tagEndRelIndex != -1 {
		tagEndIndex := tagStartIndex + tagEndRelIndex
		tagContent := textString[tagStartIndex:tagEndIndex]
		return getDialogAttributeEntry(tagContent, defaultAttributeEntry), tagEndIndex + 1
	}
	return defaultAttributeEntry, currentIndex
}

/*
handleWordWrap is a method which allows you to manage word wrapping logic when a space character is encountered. In addition, the following should be noted:

- It returns the updated cursor positions after applying word wrap if needed.

- Returns early if word wrapping is disabled (wordWrapWidth <= 0).

- Wraps to the next line if the word would exceed the word wrap width or layer width.

Example:

	newX, newY := layer.handleWordWrap(10, 5, 0, 5, 20, 80, 25)
*/
func (shared *layerType) handleWordWrap(cursorX, cursorY, xLocation, wordWidth, wordWrapWidth, layerWidth, layerHeight int) (int, int) {
	if wordWrapWidth <= 0 {
		return cursorX, cursorY
	}

	if cursorX+wordWidth >= xLocation+wordWrapWidth || cursorX+wordWidth >= layerWidth {
		cursorX = xLocation
		cursorY++
	}

	return cursorX, cursorY
}

/*
shouldSkipLeadingSpace is a method which allows you to determine if a space character at the start of a line should be
skipped. In addition, the following should be noted:

- This is typically used with word wrapping to avoid leading spaces after a wrap.

- Only applies when word wrapping is enabled and the character is a space at the start of a line.

Example:

	skip := layer.shouldSkipLeadingSpace(20, ' ', 0, 0)
*/
func (shared *layerType) shouldSkipLeadingSpace(wordWrapWidth int, character rune, cursorX, xLocation int) bool {
	return wordWrapWidth > 0 && character == ' ' && cursorX == xLocation
}

/*
isWithinVerticalBounds is a method which allows you to check if the current cursor Y position is within the layer's
height bounds. In addition, the following should be noted:

- A position is valid if it's greater than or equal to 0 and less than the layer height.

Example:

	isValid := layer.isWithinVerticalBounds(10, 25)
*/
func (shared *layerType) isWithinVerticalBounds(y, height int) bool {
	return y >= 0 && y < height
}

/*
isWithinHorizontalBounds is a method which allows you to check if the current cursor X position is within the layer's
width bounds. In addition, the following should be noted:

- A position is valid if it's greater than or equal to 0 and less than the layer width.

Example:

	isValid := layer.isWithinHorizontalBounds(10, 80)
*/
func (shared *layerType) isWithinHorizontalBounds(x, width int) bool {
	return x >= 0 && x < width
}

/*
renderCharacter is a method which allows you to render a character at the specified position with the given attributes. In addition, the following should be noted:

- It handles wide characters and background transparency.

- For wide characters, it occupies two character cells.

- Preserves the original background color when transparency is enabled.

Example:

	layer.renderCharacter(memory, 10, 5, 'A', attr)
*/
func (shared *layerType) renderCharacter(characterMemory [][]types.CharacterEntryType, cursorX, cursorY int, character rune, attributeEntry types.AttributeEntryType) {
	originalBackgroundColor := characterMemory[cursorY][cursorX].AttributeEntry.BackgroundColor

	characterMemory[cursorY][cursorX].AttributeEntry = types.NewAttributeEntry(&attributeEntry)
	characterMemory[cursorY][cursorX].Character = character

	// Handle wide characters
	if stringformat.IsRuneCharacterWide(character) {
		cursorX++
		if cursorX < len(characterMemory[0]) {
			characterMemory[cursorY][cursorX].AttributeEntry = types.NewAttributeEntry(&attributeEntry)
			characterMemory[cursorY][cursorX].Character = ' '
		}
	}

	if characterMemory[cursorY][cursorX].AttributeEntry.IsBackgroundTransparent {
		characterMemory[cursorY][cursorX].AttributeEntry.BackgroundColor = originalBackgroundColor
	}
}

/*
advanceCursor is a method which allows you to move the cursor position after rendering a character.

In addition, the following should be noted:

- It handles line wrapping when the cursor reaches the end of a line.

- When word wrapping is enabled, wraps to the next line when reaching the layer width.

- When word wrapping is disabled, stops at the layer width.

Example:

	newX, newY := layer.advanceCursor(10, 5, 0, 80, 20)
*/
func (shared *layerType) advanceCursor(cursorX, cursorY, xLocation, layerWidth, wordWrapWidth int) (int, int) {
	cursorX++

	if cursorX >= layerWidth {
		if wordWrapWidth > 0 {
			cursorX = xLocation
			cursorY++
		}
	}

	return cursorX, cursorY
}

/*
print is a method which allows you to handle all types of printing with configurable options.

In addition, the following should be noted:

- It supports word wrapping and markup/styling based on the provided options.

- Handles boundary checking to ensure text stays within the layer.

- Processes markup tags when useMarkup is true.

- Supports word wrapping when wordWrapWidth > 0.

Example:

	finalX := layer.print(layerEntry, attr, 0, 0, rune("Hello"), 20, true)
*/
func (shared *layerType) print(layerEntry *types.LayerEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int, textToPrint []rune, wordWrapWidth int, useMarkup bool) int {
	layerWidth := layerEntry.Width
	layerHeight := layerEntry.Height
	cursorXLocation := xLocation
	cursorYLocation := yLocation
	characterMemory := layerEntry.CharacterMemory

	// For markup
	currentAttributeEntry := attributeEntry

	// Convert runes to string if markup is needed
	var textString string
	if useMarkup {
		textString = string(textToPrint)
	}

	for currentCharacterIndex := 0; currentCharacterIndex < len(textToPrint); currentCharacterIndex++ {
		currentCharacter := textToPrint[currentCharacterIndex]

		// Word wrap logic
		if wordWrapWidth > 0 && currentCharacter == ' ' {
			wordWidth := calculateWordWidth(textToPrint, currentCharacterIndex, useMarkup)
			cursorXLocation, cursorYLocation = shared.handleWordWrap(cursorXLocation, cursorYLocation, xLocation, wordWidth, wordWrapWidth, layerWidth, layerHeight)
			if !shared.isWithinVerticalBounds(cursorYLocation, layerHeight) {
				return cursorXLocation - xLocation
			}
		}

		// Skip space at start of line with word wrap
		if shared.shouldSkipLeadingSpace(wordWrapWidth, currentCharacter, cursorXLocation, xLocation) {
			continue
		}

		// Handle markup
		if useMarkup && currentCharacter == '{' && currentCharacterIndex+1 < len(textToPrint) && textToPrint[currentCharacterIndex+1] == '{' {
			savedCharacterIndex := currentCharacterIndex
			currentAttributeEntry, currentCharacterIndex = shared.processMarkupTag(textToPrint, currentCharacterIndex, textString, attributeEntry)
			if savedCharacterIndex != currentCharacterIndex {
				continue
			}
		}

		// Skip if character is off-screen (vertically)
		if !shared.isWithinVerticalBounds(cursorYLocation, layerHeight) {
			cursorXLocation, cursorYLocation = shared.advanceCursor(cursorXLocation, cursorYLocation, xLocation, layerWidth, wordWrapWidth)
			if !shared.isWithinVerticalBounds(cursorYLocation, layerHeight) {
				return cursorXLocation - xLocation
			}
			continue
		}

		// Render character if it's within horizontal bounds
		if shared.isWithinHorizontalBounds(cursorXLocation, layerWidth) {
			attrToUse := currentAttributeEntry
			if !useMarkup {
				attrToUse = attributeEntry
			}
			shared.renderCharacter(characterMemory, cursorXLocation, cursorYLocation, currentCharacter, attrToUse)
		}

		// Advance cursor
		cursorXLocation, cursorYLocation = shared.advanceCursor(cursorXLocation, cursorYLocation, xLocation, layerWidth, wordWrapWidth)
		if !shared.isWithinVerticalBounds(cursorYLocation, layerHeight) {
			return cursorXLocation - xLocation
		}
	}
	return cursorXLocation - xLocation
}

/*
printWithWordWrap is a method which allows you to write text to a text layer with word wrapping enabled.

Example:

	finalX := layer.printWithWordWrap(layerEntry, attr, 0, 0, 20, rune("Hello"))
*/
func (shared *layerType) printWithWordWrap(layerEntry *types.LayerEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int, width int, textToPrint []rune) int {
	return shared.print(layerEntry, attributeEntry, xLocation, yLocation, textToPrint, width, false)
}

/*
calculateWordWidth is a method which allows you to calculate the width of a word from a given position.

In addition, the following should be noted:

- The first position is always assumed to be ' ' and is skipped.

- When markup is enabled, it processes the text to exclude markup tags from the width calculation.

- Returns the number of characters until the next space or end of text.

Example:

	width := calculateWordWidth(rune(" Hello"), 0, false)
*/
func calculateWordWidth(textToPrint []rune, start int, useMarkup bool) int {
	// Calculate the width of a word from the given position. The first position is
	// always assumed to be ' ' and is skipped.

	// If markup is enabled, use the string version and handle markup tags
	if useMarkup {
		textString := string(textToPrint)
		// Use the substring from the starting index
		substring := textString[start+1:]
		// Get the text without markup
		textWithoutMarkup := GetNonMarkupText(substring)
		// Calculate the length of the next word
		var wordWidth int
		for i := 0; i < len(textWithoutMarkup); i++ {
			if stringformat.GetSubString(textWithoutMarkup, i, 1) == " " {
				return wordWidth
			}
			wordWidth++
		}
		return wordWidth
	}

	// Standard case without markup
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
printLayer is a method which allows you to write text to a text layer.

In addition, the following should be noted:

- This is useful for internal methods that want to write text to a text layer directly, without affecting user settings.

- If the location to print falls outside the range of the text layer, then only the visible portion of your text will.

Example:

	finalX := layer.printLayer(layerEntry, attr, 0, 0, rune("Hello"))
*/
func (shared *layerType) printLayer(layerEntry *types.LayerEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int, textToPrint []rune) int {
	return shared.print(layerEntry, attributeEntry, xLocation, yLocation, textToPrint, 0, false)
}

/*
printMarkup is a method which allows you to write text to the terminal screen with word wrapping and attribute tags.

In addition, the following should be noted:

- This is similar to printDialog but without the typewriter effect and printing delay.

Example:

	layer.printMarkup(layerEntry, attr, 0, 0, 20, "Hello {{red}}World")
*/
func (shared *layerType) printMarkup(layerEntry *types.LayerEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int, widthOfLineInCharacters int, stringToPrint string) {
	arrayOfRunes := stringformat.GetRunesFromString(stringToPrint)
	shared.print(layerEntry, attributeEntry, xLocation, yLocation, arrayOfRunes, widthOfLineInCharacters, true)
}

/*
ReInitializeScreenMemory is a method which allows you to reset and reinitialize the screen memory, effectively clearing
all managed layers.

Example:

	layer.ReInitializeScreenMemory()
*/
func (shared *layerType) ReInitializeScreenMemory() {
	Layers = memory.NewMemoryManager[types.LayerEntryType]() // Initialize MemoryManager
}

/*
Add is a method which allows you to create and add a new text layer to the layer management system.

Example:

	layer.Add("main", 0, 0, 80, 25, 1, "")
*/
func (shared *layerType) Add(layerAlias string, xLocation int, yLocation int, width int, height int, zOrderPriority int, parentAlias string) {
	if width <= 0 {
		panic(fmt.Sprintf("The layer '%s' could not be created since a HotspotWidth of '%d' was specified!", layerAlias, width))
	}
	if height <= 0 {
		panic(fmt.Sprintf("The layer '%s' could not be created since a Length of '%d' was specified!", layerAlias, height))
	}
	layerEntry := types.NewLayerEntry(layerAlias, parentAlias, width, height)
	layerEntry.LayerAlias = layerAlias
	layerEntry.ScreenXLocation = xLocation
	layerEntry.ScreenYLocation = yLocation
	if zOrderPriority < 0 {
		layerEntry.ZOrder = shared.getHighestZOrderNumber(parentAlias) + 1
		layerEntry.IsTopmost = true
	} else {
		layerEntry.ZOrder = zOrderPriority
	}
	layerEntry.ParentAlias = parentAlias

	if parentAlias != "" {
		parentEntry := Layers.Get(parentAlias)
		if parentEntry != nil {
			parentEntry.IsParent = true
		} else {
			panic(fmt.Sprintf("The layer '%s' could not be created since the parent alias '%s' does not exist!", layerAlias, parentAlias))
		}
	}
	Layers.Add(layerAlias, &layerEntry)
}

/*
GetNextAlias is a method which allows you to retrieve the alias of the first available layer in the system.

In addition, the following should be noted:

- Returns an empty string if no layers exist.

Example:

	nextAlias := layer.GetNextAlias()
*/
func (shared *layerType) GetNextAlias() string {
	for _, currentEntry := range Layers.GetAllEntries() {
		return currentEntry.LayerAlias
	}
	return ""
}

/*
Delete is a method which allows you to remove a layer and all its associated child layers and controls from the system.

Example:

	layer.Delete("myLayer")
*/
func (shared *layerType) Delete(layerAlias string) {
	screenEntry := Layers.Get(layerAlias)
	if screenEntry == nil {
		panic(fmt.Sprintf("The layer '%s' could not be deleted since it does not exist!", layerAlias))
	}
	layerEntry := Layers.Get(layerAlias)
	parentAlias := layerEntry.ParentAlias
	isParent := layerEntry.IsParent

	// If this layer is a parent, recursively delete all children first
	if isParent {
		shared.deleteAllChildrenOfParent(layerAlias)
	}

	// Delete all controls on this layer
	Labels.RemoveAll(layerAlias)
	Buttons.RemoveAll(layerAlias)
	Checkboxes.RemoveAll(layerAlias)
	Dropdowns.RemoveAll(layerAlias)
	ProgressBars.RemoveAll(layerAlias)
	RadioButtons.RemoveAll(layerAlias)
	ScrollBars.RemoveAll(layerAlias)
	Selectors.RemoveAll(layerAlias)
	Textboxes.RemoveAll(layerAlias)
	TextFields.RemoveAll(layerAlias)
	Tooltips.RemoveAll(layerAlias)
	Viewports.RemoveAll(layerAlias)
	// Remove the layer itself
	Layers.Remove(layerAlias)

	// Update parent's IsParent status if needed
	if parentAlias != "" {
		parentEntry := Layers.Get(parentAlias)
		if parentEntry != nil {
			if !shared.IsAParent(parentAlias) {
				layerEntry = Layers.Get(parentAlias)
				layerEntry.IsParent = false
			}
		}
	}
}

/*
deleteAllChildrenOfParent is a method which allows you to recursively delete all child layers associated with a given
parent layer.

Example:

	layer.deleteAllChildrenOfParent("parentLayer")
*/
func (shared *layerType) deleteAllChildrenOfParent(parentAlias string) {
	// Get all entries first to avoid modification during iteration
	entries := make([]string, 0)
	for _, currentValue := range Layers.GetAllEntries() {
		if currentValue.ParentAlias == parentAlias {
			entries = append(entries, currentValue.LayerAlias)
		}
	}

	// Delete each child layer
	for _, childAlias := range entries {
		shared.Delete(childAlias)
	}
}

/*
IsAParent is a method which allows you to determine if a given layer acts as a parent to any other layers.

Example:

	isParent := layer.IsAParent("myLayer")
*/
func (shared *layerType) IsAParent(parentAlias string) bool {
	isParent := false
	for _, currentValue := range Layers.GetAllEntries() {
		if currentValue.ParentAlias == parentAlias {
			isParent = true
		}
	}
	return isParent
}

/*
GetSortedLayerMemoryAliasSlice is a method which allows you to retrieve a sorted list of layer aliases based on their
z-order.

In addition, the following should be noted:

- Returns a list of layer-alias and z-order pairs, sorted by z-order.

Example:

	sortedLayers := layer.GetSortedLayerMemoryAliasSlice()
*/
func (shared *layerType) GetSortedLayerMemoryAliasSlice() LayerAliasZOrderPairList {
	pairList := make(LayerAliasZOrderPairList, len(Layers.GetAllEntries()))
	currentEntry := 0
	for currentKey, currentValue := range Layers.GetAllEntriesWithKeys() {
		pairList[currentEntry].Key = currentKey
		pairList[currentEntry].Value = currentValue.ZOrder
		currentEntry++
	}
	sort.SliceStable(pairList, func(firstIndex, secondIndex int) bool {
		return pairList[firstIndex].Value < pairList[secondIndex].Value
	})
	return pairList
}

/*
SetHighestZOrderNumber is a method which allows you to set the highest z-order number for a specific layer among its
siblings.

Example:

	layer.SetHighestZOrderNumber("topLayer", "parentLayer")
*/
func (shared *layerType) SetHighestZOrderNumber(layerAlias string, parentAlias string) {
	if Layers.IsExists(layerAlias) {
		highestZOrderNumber := shared.getHighestZOrderNumber(parentAlias)
		for _, currentValue := range Layers.GetAllEntries() {
			if currentValue.ParentAlias == parentAlias && currentValue.ZOrder == highestZOrderNumber {
				currentValue.ZOrder = highestZOrderNumber - 1
				currentValue.IsTopmost = false
			}
		}
		Layers.Get(layerAlias).ZOrder = highestZOrderNumber
		Layers.Get(layerAlias).IsTopmost = true
	}
}

/*
getHighestZOrderNumber is a method which allows you to retrieve the highest z-order value currently used among siblings
under a given parent.

Example:

	highest := layer.getHighestZOrderNumber("parentLayer")
*/
func (shared *layerType) getHighestZOrderNumber(parentAlias string) int {
	highestZOrderNumber := 0
	for _, currentValue := range Layers.GetAllEntries() {
		if currentValue.ParentAlias == parentAlias && currentValue.ZOrder > highestZOrderNumber {
			highestZOrderNumber = currentValue.ZOrder
		}
	}
	return highestZOrderNumber
}

/*
SetTopmostLayer is a method which allows you to set the specified layer to be the topmost layer among its siblings.

Example:

	layer.SetTopmostLayer("myLayer")
*/
func (shared *layerType) SetTopmostLayer(layerAlias string) {
	if !Layers.IsExists(layerAlias) {
		return
	}

	targetLayer := Layers.Get(layerAlias)

	// Reset the current topmost layer's flag.
	for _, layerEntry := range Layers.GetAllEntries() {
		if layerEntry.ParentAlias == targetLayer.ParentAlias && layerEntry.IsTopmost {
			layerEntry.IsTopmost = false
		}
	}

	// Find the highest Z-order and set the target layer's Z-order higher.
	highestZOrder := 0
	for _, layerEntry := range Layers.GetAllEntries() {
		if layerEntry.ParentAlias == targetLayer.ParentAlias && layerEntry.ZOrder > highestZOrder {
			highestZOrder = layerEntry.ZOrder
		}
	}

	targetLayer.ZOrder = highestZOrder + 1
	targetLayer.IsTopmost = true
}

/*
GetRootParentAlias is a method which allows you to recursively find the root parent layer alias for a given layer.

In addition, the following should be noted:

- It returns the alias of the root parent layer and the alias of its immediate child in the path.

Example:

	root, child := layer.GetRootParentAlias("grandchild", "")
*/
func (shared *layerType) GetRootParentAlias(layerAlias string, previousChildAlias string) (string, string) {
	layerEntry := Layers.Get(layerAlias)
	if layerEntry.ParentAlias != "" {
		childToTrack := previousChildAlias
		if childToTrack == "" {
			childToTrack = layerAlias
		}
		return shared.GetRootParentAlias(layerEntry.ParentAlias, childToTrack)
	}
	return layerAlias, previousChildAlias
}

// ============================================================================
// REGULAR ENTRY
// ============================================================================

/*
getUUID is a method which allows you to generate a unique version 4 UUID string.

Example:

	id := getUUID()
*/
func getUUID() string {
	id := uuid.New()
	return id.String()
}

/*
Clear is a method which allows you to empty the current text layer of all its contents.

Example:

	layerInstance.Clear()
*/
func (shared *LayerInstanceType) Clear() {
	layerEntry := Layers.Get(shared.layerAlias)
	localAttributeEntry := types.NewAttributeEntry()
	fillArea(layerEntry, localAttributeEntry, "", 0, 0, layerEntry.Width, layerEntry.Height, 0)
}

/*
GetAlias is a method which allows you to retrieve the unique alias associated with the current layer instance.

In addition, the following should be noted:

- Returns the string alias of the layer.

Example:

	alias := layerInstance.GetAlias()
*/
func (shared *LayerInstanceType) GetAlias() string {
	return shared.layerAlias
}

/*
DrawImage is a method which allows you to draw an image on a given text layer. This method supports various image
formats and drawing styles, allowing for flexible rendering of images as text-based art.

In addition, the following should be noted:

- If the image to be drawn is not already loaded in memory, it will be loaded automatically and then unloaded after the.

- When drawing images with transparencies, the transparent edges are only computed once against the layer. Moving the.

Example:

	err := layerInstance.DrawImage("photo.png", style, 0, 0, 40, 20, 0.5)
*/
func (shared *LayerInstanceType) DrawImage(fileName string, drawingStyle types.ImageStyleEntryType, xLocation int, yLocation int, widthInCharacters int, heightInCharacters int, blurSigma float64) error {
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
	imageEntryType := getImage(fileName)
	imageLayer := imageEntryType.LayerEntry
	var currentLayer *types.LayerEntryType
	currentLayer = Layers.Get(shared.layerAlias)
	if imageEntryType.ImageData != nil {
		imageData := imageEntryType.ImageData
		// Get the current layer to pass for transparency handling
		imageLayer = getImageLayer(fileName, imageData, drawingStyle, widthInCharacters, heightInCharacters, blurSigma)
	}
	drawImageToLayer(currentLayer, imageLayer, xLocation, yLocation)
	return err
}

/*
DrawComposedImage is a method which allows you to draw a composed image on a text layer using a specific drawing style.
In addition, the following should be noted:

  - The image will be rendered according to the settings in the provided `drawingStyle`, including any blurring or
    dithering effects.

Example:

	err := layerInstance.DrawComposedImage(composedImage, style, 10, 5, 20, 10)
*/
func (shared *LayerInstanceType) DrawComposedImage(imageComposeEntry ImageComposerEntryType, drawingStyle types.ImageStyleEntryType, xLocation int, yLocation int, widthInCharacters int, heightInCharacters int) error {
	var err error
	baseImage := imageComposeEntry.RenderImage()

	// Get the current layer to pass for transparency handling
	var currentLayer *types.LayerEntryType
	currentLayer = Layers.Get(shared.layerAlias)

	imageLayer := getImageLayer("", baseImage, drawingStyle, widthInCharacters, heightInCharacters, drawingStyle.BlurSigmaIntensity)
	drawImageToLayer(currentLayer, imageLayer, xLocation, yLocation)
	return err
}

/*
AddButton is a method which allows you to add a new button control to the current layer.

Example:

	btn := layerInstance.AddButton("Click Me", style, 10, 5, 12, 1, true)
*/
func (shared *LayerInstanceType) AddButton(buttonLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isEnabled bool) ButtonInstanceType {
	buttonAlias := getUUID()
	buttonInstance := Button.Add(shared.layerAlias, buttonAlias, buttonLabel, styleEntry, xLocation, yLocation, width, height, isEnabled)
	return buttonInstance
}

/*
AddCheckbox is a method which allows you to add a new checkbox control to the current layer.

Example:

	cb := layerInstance.AddCheckbox("Option A", style, 10, 5, false, true)
*/
func (shared *LayerInstanceType) AddCheckbox(checkboxLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, isSelected bool, isEnabled bool) CheckboxInstanceType {
	checkboxAlias := getUUID()
	checkboxInstance := Checkbox.Add(shared.layerAlias, checkboxAlias, checkboxLabel, styleEntry, xLocation, yLocation, isSelected, isEnabled)
	return checkboxInstance
}

/*
AddDropdown is a method which allows you to add a new dropdown control to the current layer.

Example:

	dd := layerInstance.AddDropdown(style, selections, 10, 5, 5, 15, 0)
*/
func (shared *LayerInstanceType) AddDropdown(styleEntry types.TuiStyleEntryType, selectionEntry types.SelectionEntryType, xLocation int, yLocation int, selectorHeight int, itemWidth int, defaultItemSelected int) DropdownInstanceType {
	dropdownAlias := getUUID()
	dropdownInstance := Dropdown.Add(shared.layerAlias, dropdownAlias, styleEntry, selectionEntry, xLocation, yLocation, selectorHeight, itemWidth, defaultItemSelected)
	return dropdownInstance
}

/*
AddLabel is a method which allows you to add a new label control to the current layer.

Example:

	lbl := layerInstance.AddLabel("Username:", style, 5, 5, 10)
*/
func (shared *LayerInstanceType) AddLabel(labelValue string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int) LabelInstanceType {
	labelAlias := getUUID()
	labelInstance := Label.Add(shared.layerAlias, labelAlias, labelValue, styleEntry, xLocation, yLocation, width)
	return labelInstance
}

/*
AddProgressBar is a method which allows you to add a new progress bar control to the current layer.

Example:

	pb := layerInstance.AddProgressBar("Loading", style, 10, 5, 20, 1, false, 50, 100, false)
*/
func (shared *LayerInstanceType) AddProgressBar(progressBarLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isVertical bool, value int, maxValue int, isBackgroundTransparent bool) ProgressBarInstanceType {
	progressBarAlias := getUUID()
	progressBarInstance := ProgressBar.Add(shared.layerAlias, progressBarAlias, progressBarLabel, styleEntry, xLocation, yLocation, width, height, isVertical, value, maxValue, isBackgroundTransparent)
	return progressBarInstance
}

/*
AddRadioButton is a method which allows you to add a new radio button control to the current layer.

Example:

	rb := layerInstance.AddRadioButton("Option 1", style, 10, 5, 1, false)
*/
func (shared *LayerInstanceType) AddRadioButton(radioButtonLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, groupId int, isSelected bool) RadioButtonInstanceType {
	radioButtonAlias := getUUID()
	radioButtonInstance := radioButton.Add(shared.layerAlias, radioButtonAlias, radioButtonLabel, styleEntry, xLocation, yLocation, groupId, isSelected)
	return radioButtonInstance
}

/*
AddScrollbar is a method which allows you to add a new scrollbar control to the current layer.

Example:

	sb := layerInstance.AddScrollbar(style, 10, 5, 20, 100, 0, 1, false)
*/
func (shared *LayerInstanceType) AddScrollbar(styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, length int, maxScrollValue int, scrollValue int, scrollIncrement int, isHorizontal bool) ScrollbarInstanceType {
	scrollbarAlias := getUUID()
	scrollbarInstance := scrollbar.Add(shared.layerAlias, scrollbarAlias, styleEntry, xLocation, yLocation, length, maxScrollValue, scrollValue, scrollIncrement, isHorizontal)
	return scrollbarInstance
}

/*
AddSelector is a method which allows you to add a new selector control to the current layer.

Example:

	sel := layerInstance.AddSelector(style, selections, 10, 5, 10, 20, 1, 0, 0, false, true)
*/
func (shared *LayerInstanceType) AddSelector(styleEntry types.TuiStyleEntryType, selectionEntry types.SelectionEntryType, xLocation int, yLocation int, selectorHeight int, itemWidth int, numberOfColumns int, viewportPosition int, selectedItem int, highlightOnClickOnly bool, isBorderDrawn bool) SelectorInstanceType {
	selectorAlias := getUUID()
	selectorInstance := Selector.Add(shared.layerAlias, selectorAlias, styleEntry, selectionEntry, xLocation, yLocation, selectorHeight, itemWidth, numberOfColumns, viewportPosition, selectedItem, highlightOnClickOnly, isBorderDrawn)
	return selectorInstance
}

/*
AddTextField is a method which allows you to add a new text field control to the current layer.

Example:

	tf := layerInstance.AddTextField(style, 10, 5, 20, 50, false, "", true)
*/
func (shared *LayerInstanceType) AddTextField(styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, maxLengthAllowed int, isPasswordProtected bool, defaultValue string, isEnabled bool) TextFieldInstanceType {
	textFieldAlias := getUUID()
	textFieldInstance := TextField.Add(shared.layerAlias, textFieldAlias, styleEntry, xLocation, yLocation, width, maxLengthAllowed, isPasswordProtected, defaultValue, isEnabled)
	return textFieldInstance
}

/*
AddTextbox is a method which allows you to add a new multi-line textbox control to the current layer.

Example:

	tb := layerInstance.AddTextbox(style, 10, 5, 30, 10, true)
*/
func (shared *LayerInstanceType) AddTextbox(styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isBorderDrawn bool) TextboxInstanceType {
	textBoxAlias := getUUID()
	textBoxInstance := textbox.Add(shared.layerAlias, textBoxAlias, styleEntry, xLocation, yLocation, width, height, isBorderDrawn)
	return textBoxInstance
}

/*
AddTooltip is a method which allows you to add a new tooltip control to the current layer.

Example:

	tt := layerInstance.AddTooltip("Help text", style, 10, 5, 5, 1, 10, 6, 20, 3, false, true, 500)
*/
func (shared *LayerInstanceType) AddTooltip(tooltipValue string, styleEntry types.TuiStyleEntryType, hotspotXLocation int, hotspotYLocation int, hotspotWidth int, hotspotHeight int, tooltipXLocation int, tooltipYLocation int, tooltipWidth int, tooltipHeight int, isLocationAbsolute bool, isBorderDrawn bool, hoverTime int) TooltipInstanceType {
	tooltipAlias := getUUID()
	tooltipInstance := Tooltip.Add(shared.layerAlias, tooltipAlias, tooltipValue, styleEntry, hotspotXLocation, hotspotYLocation, hotspotWidth, hotspotHeight, tooltipXLocation, tooltipYLocation, tooltipWidth, tooltipHeight, isLocationAbsolute, isBorderDrawn, hoverTime)
	return tooltipInstance
}

/*
AddViewport is a method which allows you to add a viewport to a given text layer. A viewport is a read-only text display
control that can show text with markup codes for colorization. It supports scrollback history and text wrapping.

In addition, the following should be noted:

  - If vertical scrollbars are enabled, the viewport will maintain scrollback history up to the specified
    maxHistoryLines.

  - If vertical scrollbars are not enabled, then no history is needed and only memory for the visible display is
    required.

  - If isLinesWrapped is enabled, text printed to screen will wrap text cleanly like dialog.go's printMarkup
    method and no horizontal scrollbars will be shown.

  - If isLinesWrapped is disabled, lines will remain on the same line and horizontal scrollbars will be rendered
    if needed.

  - Text can be added to the viewport using the Print method.

Example:

	vp := layerInstance.AddViewport(style, 0, 0, 40, 10, true, true, 100)
*/
func (shared *LayerInstanceType) AddViewport(styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isLinesWrapped bool, isBorderDrawn bool, maxHistoryLines int) ViewportInstanceType {
	viewportAlias := getUUID()
	viewportInstance := viewport.Add(shared.layerAlias, viewportAlias, styleEntry, xLocation, yLocation, width, height, isLinesWrapped, isBorderDrawn, maxHistoryLines)
	return viewportInstance
}

/*
AddFileMenu is a method which allows you to add a file menu to a layer.

In addition, the following should be noted:

- The file menu will be drawn at the specified location with the given style.

- Each heading in the menu can have its own dropdown with selectable items.

- The top level headings widths are always dynamic based on how large the heading is.

- The file menu reuses existing selectors for dropdown functionality.

Example:

	fm := layerInstance.AddFileMenu(style, string{"File", "Edit"}, selections, 0, 0, true)
*/
func (shared *LayerInstanceType) AddFileMenu(styleEntry types.TuiStyleEntryType, menuHeadings []string, menuSelections []types.SelectionEntryType, xLocation int, yLocation int, isEnabled bool) FileMenuInstanceType {
	menuAlias := getUUID()
	fileMenuInstance := FileMenu.Add(shared.layerAlias, menuAlias, styleEntry, menuHeadings, menuSelections, xLocation, yLocation, isEnabled)
	return fileMenuInstance
}

/*
DeleteAllButtons is a method which allows you to remove all buttons from the current layer.

Example:

	layerInstance.DeleteAllButtons()
*/
func (shared *LayerInstanceType) DeleteAllButtons() {
	Buttons.RemoveAll(shared.layerAlias)
}

/*
DeleteAllCheckboxes is a method which allows you to remove all checkboxes from the current layer.

Example:

	layerInstance.DeleteAllCheckboxes()
*/
func (shared *LayerInstanceType) DeleteAllCheckboxes() {
	Checkboxes.RemoveAll(shared.layerAlias)
}

/*
DeleteAllDropdowns is a method which allows you to remove all dropdowns from the current layer.

Example:

	layerInstance.DeleteAllDropdowns()
*/
func (shared *LayerInstanceType) DeleteAllDropdowns() {
	Dropdowns.RemoveAll(shared.layerAlias)
}

/*
DeleteAllLabels is a method which allows you to remove all labels from the current layer.

Example:

	layerInstance.DeleteAllLabels()
*/
func (shared *LayerInstanceType) DeleteAllLabels() {
	Labels.RemoveAll(shared.layerAlias)
}

/*
DeleteAllProgressBars is a method which allows you to remove all progress bars from the current layer.

Example:

	layerInstance.DeleteAllProgressBars()
*/
func (shared *LayerInstanceType) DeleteAllProgressBars() {
	ProgressBars.RemoveAll(shared.layerAlias)
}

/*
DeleteAllRadioButtons is a method which allows you to remove all radio buttons from the current layer.

Example:

	layerInstance.DeleteAllRadioButtons()
*/
func (shared *LayerInstanceType) DeleteAllRadioButtons() {
	RadioButtons.RemoveAll(shared.layerAlias)
}

/*
DeleteAllScrollbars is a method which allows you to remove all scroll bars from the current layer.

Example:

	layerInstance.DeleteAllScrollbars()
*/
func (shared *LayerInstanceType) DeleteAllScrollbars() {
	ScrollBars.RemoveAll(shared.layerAlias)
}

/*
DeleteAllSelectors is a method which allows you to remove all selectors from the current layer.

Example:

	layerInstance.DeleteAllSelectors()
*/
func (shared *LayerInstanceType) DeleteAllSelectors() {
	Selectors.RemoveAll(shared.layerAlias)
}

/*
DeleteAllTextFields is a method which allows you to remove all text fields from the current layer.

Example:

	layerInstance.DeleteAllTextFields()
*/
func (shared *LayerInstanceType) DeleteAllTextFields() {
	TextFields.RemoveAll(shared.layerAlias)
}

/*
DeleteAllTextboxes is a method which allows you to remove all textboxes from the current layer.

Example:

	layerInstance.DeleteAllTextboxes()
*/
func (shared *LayerInstanceType) DeleteAllTextboxes() {
	Textboxes.RemoveAll(shared.layerAlias)
}

/*
DeleteAllTooltips is a method which allows you to remove all tooltips from the current layer.

Example:

	layerInstance.DeleteAllTooltips()
*/
func (shared *LayerInstanceType) DeleteAllTooltips() {
	Tooltips.RemoveAll(shared.layerAlias)
}

/*
DeleteAllViewports is a method which allows you to remove all viewports from the current layer.

Example:

	layerInstance.DeleteAllViewports()
*/
func (shared *LayerInstanceType) DeleteAllViewports() {
	Viewports.RemoveAll(shared.layerAlias)
}

/*
DeleteAllFileMenus is a method which allows you to remove all file menus from the current layer.

Example:

	layerInstance.DeleteAllFileMenus()
*/
func (shared *LayerInstanceType) DeleteAllFileMenus() {
	FileMenus.RemoveAll(shared.layerAlias)
}

/*
Print is a method which allows you to write text to the current layer.

In addition, the following should be noted:

  - When text is written to the text layer, the cursor position is also updated to reflect its new location.

  - If the string to print ends up being too long to fit at its current location, then only the visible portion of
    your text will be printed.

Example:

	layerInstance.Print("Hello World")
*/
func (shared *LayerInstanceType) Print(textToPrint string) {
	printLayerInstance(shared, textToPrint)
}

/*
Locate is a method which allows you to set the default cursor location on your text layer for printing with.

In addition, the following should be noted:

  - If you pass in a location value that falls outside the dimensions of the default text layer, a panic will be
    generated.

  - Valid text layer locations start at position (0,0) for the upper left corner.

Example:

	layerInstance.Locate(10, 5)
*/
func (shared *LayerInstanceType) Locate(xLocation int, yLocation int) {
	locateLayerInstance(shared, xLocation, yLocation)
}

/*
Color is a method which allows you to set default colors on your text layer for printing with.

In addition, the following should be noted:

- The color index specified corresponds to the 16 color ANSI standard, where color 0 is Black and 15 is Bright White.

Example:

	layerInstance.Color(15, 0)
*/
func (shared *LayerInstanceType) Color(foregroundColorIndex int, backgroundColorIndex int) {
	colorLayerInstance(shared, foregroundColorIndex, backgroundColorIndex)
}

/*
ColorRGB is a method which allows you to set default colors on your text layer for printing with using RGB values.

In addition, the following should be noted:

- This method allows you to specify colors using RGB color index values within the range of 0 to 255.

Example:

	layerInstance.ColorRGB(255, 255, 255, 0, 0, 0)
*/
func (shared *LayerInstanceType) ColorRGB(foregroundRedIndex int32, foregroundGreenIndex int32, foregroundBlueIndex int32, backgroundRedIndex int32, backgroundGreenIndex int32, backgroundBlueIndex int32) {
	colorLayerRGBInstance(shared, foregroundRedIndex, foregroundGreenIndex, foregroundBlueIndex, backgroundRedIndex, backgroundGreenIndex, backgroundBlueIndex)
}

/*
Color24Bit is a method which allows you to color the current layer using a 24-bit color expressed as an int32.

Example:

	layerInstance.Color24Bit(fgColor, bgColor)
*/
func (shared *LayerInstanceType) Color24Bit(foregroundColor constants.ColorType, backgroundColor constants.ColorType) {
	colorLayer24BitInstance(shared, foregroundColor, backgroundColor)
}

/*
SetAlpha is a method which allows you to set the alpha value for the current layer. This lets you perform pseudo
transparencies by making the layer foreground and background colors blend with the layers underneath it to the degree
specified.

In addition, the following should be noted:

- An alpha value of 1.0 is equal to 100% visible, while an alpha value of 0.0 is 0% visible.

- If the percent change specified is outside of the RGB color range, then the color will simply bottom or max out.

Example:

	layerInstance.SetAlpha(0.5)
*/
func (shared *LayerInstanceType) SetAlpha(alphaValue float32) {
	setLayerAlphaInstance(shared, alphaValue)
}

/*
SetZOrder is a method which allows you to set the z-order priority for the current layer.

Example:

	layerInstance.SetZOrder(10)
*/
func (shared *LayerInstanceType) SetZOrder(zOrder int) {
	setLayerZOrderInstance(shared, zOrder)
}

/*
DrawVerticalLine is a method which allows you to draw a vertical line on a text layer.

In addition, the following should be noted:

  - This method also has the ability to draw connectors in case the line intersects with other lines that have
    already been drawn.

  - If the line to be drawn falls outside the area of the text layer specified, then only the visible portion of the
    line will be drawn.

Example:

	layerInstance.DrawVerticalLine(style, 10, 5, 10, true)
*/
func (shared *LayerInstanceType) DrawVerticalLine(styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, height int, isConnectorsDrawn bool) {
	layerEntry := Layers.Get(shared.layerAlias)
	localAttributeEntry := types.NewAttributeEntry()
	drawVerticalLine(layerEntry, styleEntry, localAttributeEntry, xLocation, yLocation, height, isConnectorsDrawn)
}

/*
DrawHorizontalLine is a method which allows you to draw a horizontal line on a text layer.

In addition, the following should be noted:

  - This method also has the ability to draw connectors in case the line intersects with other lines that have
    already been drawn.

  - If the line to be drawn falls outside the area of the text layer specified, then only the visible portion of the
    line will be drawn.

Example:

	layerInstance.DrawHorizontalLine(style, 10, 5, 20, true)
*/
func (shared *LayerInstanceType) DrawHorizontalLine(styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, isConnectorsDrawn bool) {
	layerEntry := Layers.Get(shared.layerAlias)
	localAttributeEntry := types.NewAttributeEntry()
	drawHorizontalLine(layerEntry, styleEntry, localAttributeEntry, xLocation, yLocation, width, isConnectorsDrawn)
}

/*
DrawBorder is a method which allows you to draw a border on a given text layer.

In addition, the following should be noted:

  - Borders differ from frames since they are flat shaded and do not have a raised or sunken look to them.

  - If the border to be drawn falls outside the range of the specified layer, then only the visible portion of the
    border will be drawn.

  - The 'isInteractive' option allows you to specify if the window should interact with the layer being drawn on.

Example:

	layerInstance.DrawBorder(style, 5, 5, 40, 10, true)
*/
func (shared *LayerInstanceType) DrawBorder(styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isInteractive bool) {
	layerEntry := Layers.Get(shared.layerAlias)
	localAttributeEntry := types.NewAttributeEntry()
	drawBorder(layerEntry, styleEntry, localAttributeEntry, xLocation, yLocation, width, height, isInteractive)
}

/*
DrawFrameLabel is a method which allows you to draw a label for a frame.

In addition, the following should be noted:

  - The label will be automatically enclosed by characters to blend in with a border of a frame.

  - If the frame label to be drawn falls outside the range of the specified layer, then only the visible portion of
    the label will be drawn.

Example:

	layerInstance.DrawFrameLabel(style, "Settings", 7, 5)
*/
func (shared *LayerInstanceType) DrawFrameLabel(styleEntry types.TuiStyleEntryType, label string, xLocation int, yLocation int) {
	layerEntry := Layers.Get(shared.layerAlias)
	drawFrameLabel(layerEntry, styleEntry, label, xLocation, yLocation)
}

/*
DrawFrame is a method which allows you to draw a frame on a given text layer.

In addition, the following should be noted:

  - Frames differ from borders since borders are flat shaded and do not have a raised or sunken look to them.

  - If the frame to be drawn falls outside the range of the specified layer, then only the visible portion of the
    frame will be drawn.

  - The 'isInteractive' option allows you to specify if the window should interact with the layer being drawn on.

Example:

	layerInstance.DrawFrame(style, true, 5, 5, 40, 10, true)
*/
func (shared *LayerInstanceType) DrawFrame(styleEntry types.TuiStyleEntryType, isRaised bool, xLocation int, yLocation int, width int, height int, isInteractive bool) {
	layerEntry := Layers.Get(shared.layerAlias)
	localAttributeEntry := types.NewAttributeEntry()
	if isRaised {
		drawFrame(layerEntry, styleEntry, localAttributeEntry, constants.FrameStyleRaised, xLocation, yLocation, width, height, isInteractive)
	} else {
		drawFrame(layerEntry, styleEntry, localAttributeEntry, constants.FrameStyleSunken, xLocation, yLocation, width, height, isInteractive)
	}
}

/*
DrawWindow is a method which allows you to draw a window on a given text layer.

In addition, the following should be noted:

  - Windows differ from borders since the entire area the window surrounds gets filled with a solid background
    color.

  - If the window to be drawn falls outside the range of the specified layer, then only the visible portion of the
    window will be drawn.

  - The 'isInteractive' option allows you to specify if the window should interact with the layer being drawn on.

Example:

	layerInstance.DrawWindow(style, 5, 5, 40, 10, true)
*/
func (shared *LayerInstanceType) DrawWindow(styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isInteractive bool) {
	layerEntry := Layers.Get(shared.layerAlias)
	localAttributeEntry := types.NewAttributeEntry()
	drawWindow(layerEntry, styleEntry, localAttributeEntry, xLocation, yLocation, width, height, isInteractive)
}

/*
DrawShadow is a method which allows you to draw shadows on a given text layer.

In addition, the following should be noted:

- Shadows are simply transparent areas which darken whatever text layers are underneath it by a specified degree.

- The alpha value can range from 0.0 (no shadow) to 1.0 (totally black).

Example:

	layerInstance.DrawShadow(7, 7, 40, 10, 0.5)
*/
func (shared *LayerInstanceType) DrawShadow(xLocation int, yLocation int, width int, height int, alphaValue float32) {
	layerEntry := Layers.Get(shared.layerAlias)
	localAttributeEntry := types.NewAttributeEntry()
	drawShadow(layerEntry, localAttributeEntry, xLocation, yLocation, width, height, alphaValue)
}

/*
FillArea is a method which allows you to fill an area of a given text layer with characters of your choice.

In addition, the following should be noted:

  - If you wish to fill the area with repeating text, simply provide the string you wish to repeat.

  - If the area to fill falls outside the range of the specified layer, then only the visible portion of the fill
    will be drawn.

Example:

	layerInstance.FillArea("*", 0, 0, 80, 25)
*/
func (shared *LayerInstanceType) FillArea(fillCharacters string, xLocation int, yLocation int, width int, height int) {
	layerEntry := Layers.Get(shared.layerAlias)
	attributeEntry := layerEntry.DefaultAttribute
	fillArea(layerEntry, attributeEntry, fillCharacters, xLocation, yLocation, width, height, constants.NullCellControlLocation)
}

/*
FillLayer is a method which allows you to fill an entire layer with characters of your choice.

In addition, the following should be noted:

- If you wish to fill the layer with repeating text, simply provide the string you wish to repeat.

Example:

	layerInstance.FillLayer(".")
*/
func (shared *LayerInstanceType) FillLayer(fillCharacters string) {
	layerEntry := Layers.Get(shared.layerAlias)
	attributeEntry := layerEntry.DefaultAttribute
	fillLayer(layerEntry, attributeEntry, fillCharacters)
}

/*
DrawBar is a method which allows you to draw a horizontal bar on a given text layer row.

In addition, the following should be noted:

- This is useful for drawing application headers or status bar footers.

Example:

	layerInstance.DrawBar(style, 0, 0, 80, "=")
*/
func (shared *LayerInstanceType) DrawBar(styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, barLength int, fillCharacters string) {
	layerEntry := Layers.Get(shared.layerAlias)
	attributeEntry := types.NewAttributeEntry()
	attributeEntry.ForegroundColor = styleEntry.Bar.ForegroundColor
	attributeEntry.BackgroundColor = styleEntry.Bar.BackgroundColor
	fillArea(layerEntry, attributeEntry, fillCharacters, xLocation, yLocation, barLength, 1, constants.NullCellControlLocation)
}

/*
MoveLayerByAbsoluteValue is a method which allows you to move a text layer by an absolute value.

In addition, the following should be noted:

- This is useful if you know exactly what position you wish to move your text layer to.

- If you move your layer outside the visible terminal display, only the visible display area will be rendered.

Example:

	layerInstance.MoveLayerByAbsoluteValue(10, 5)
*/
func (shared *LayerInstanceType) MoveLayerByAbsoluteValue(xLocation int, yLocation int) {
	moveLayerByAbsoluteValue(shared.layerAlias, xLocation, yLocation)
}

/*
MoveLayerByRelativeValue is a method which allows you to move a text layer by a relative value.

In addition, the following should be noted:

  - This is useful for windows, foregrounds, backgrounds, or any kind of animations or movement you may wish
    to do in increments.

  - If you move your layer outside the visible terminal display, only the visible display area will be rendered.

  - If your text layer is a child of a parent layer, then only the visible display area will be rendered on
    the parent.

Example:

	layerInstance.MoveLayerByRelativeValue(-1, 2)
*/
func (shared *LayerInstanceType) MoveLayerByRelativeValue(xLocation int, yLocation int) {
	moveLayerByRelativeValue(shared.layerAlias, xLocation, yLocation)
}

/*
Delete is a method which allows you to remove a text layer.

In addition, the following should be noted:

  - If you wish to reuse a text layer for a future purpose, you may also consider making the layer invisible
    instead of deleting it.

  - When a text layer is deleted, all child text layers are recursively deleted as well.

  - If any dynamically drawn TUI controls reference the deleted layer, they will still be present.

  - If you attempt to delete a text layer which is currently set as your default text layer, then a panic will
    be generated.

  - If you attempt to delete a text layer that does not exist, then the operation will be ignored.

Example:

	layerInstance.Delete()
*/
func (shared *LayerInstanceType) Delete() {
	deleteLayer(shared.layerAlias)
	shared.layerAlias = ""
}

/*
IsExists is a method which allows you to check if the current layer instance still exists in the layer management
system.

Example:

	exists := layerInstance.IsExists()
*/
func (shared *LayerInstanceType) IsExists() bool {
	if shared.layerAlias != "" {
		return true
	}
	return false
}

/*
SetIsVisible is a method which allows you to set the visibility of the current layer.

Example:

	layerInstance.SetIsVisible(false)
*/
func (shared *LayerInstanceType) SetIsVisible(isVisible bool) {
	validateLayer(shared.layerAlias)
	setLayerIsVisible(shared.layerAlias, isVisible)
}

/*
SetAlphaValue is a method which defines a global transparency level for the entire layer by providing a percentage
between zero and one. This value acts as a multiplier for all individual cell alpha values during the layer composition
phase to determine the final visual opacity of the layer's contents.

Example:
    instance.SetAlphaValue(0.5)
*/
func (shared *LayerInstanceType) SetAlphaValue(alphaValue float32) {
	validateLayer(shared.layerAlias)
	setLayerAlphaValue(shared.layerAlias, alphaValue)
}

/*
GetAlphaValue is a method which retrieves the current global transparency level of the layer. It returns a float32
where 0.0 represents a fully transparent layer and 1.0 represents a fully opaque one.

Example:
    alpha := instance.GetAlphaValue()
*/
func (shared *LayerInstanceType) GetAlphaValue() float32 {
	validateLayer(shared.layerAlias)
	return getLayerAlphaValue(shared.layerAlias)
}

/*
SetTransparencyStrategy is a method which configures the specific algorithm used to handle character-level visibility
for the layer. It allows you to toggle between standard color blending and various dithering patterns that determine
whether a cell's character should be rendered or if the layer beneath it should show through.

Example:
    instance.SetTransparencyStrategy(constants.TransparencyStrategyBayer8x8)
*/
func (shared *LayerInstanceType) SetTransparencyStrategy(strategy constants.TransparencyStrategy) {
	validateLayer(shared.layerAlias)
	setLayerTransparencyStrategy(shared.layerAlias, strategy)
}

/*
GetTransparencyStrategy is a method which obtains the current rendering algorithm used for cell-level transparency.
It returns the transparency strategy constant currently assigned to the layer entry.

Example:
    strategy := instance.GetTransparencyStrategy()
*/
func (shared *LayerInstanceType) GetTransparencyStrategy() constants.TransparencyStrategy {
	validateLayer(shared.layerAlias)
	return getLayerTransparencyStrategy(shared.layerAlias)
}

/*
SetTopmost is a method which allows you to set the current layer to be the topmost layer within its parent hierarchy.

Example:

	layerInstance.SetTopmost()
*/
func (shared *LayerInstanceType) SetTopmost() {
	validateLayer(shared.layerAlias)
	layer.SetTopmostLayer(shared.layerAlias)
}

/*
GetLocation is a method which allows you to retrieve the current screen X and Y coordinates of the layer.

Example:

	x, y := layerInstance.GetLocation()
*/
func (shared *LayerInstanceType) GetLocation() (int, int) {
	validateLayer(shared.layerAlias)
	layerEntry := Layers.Get(shared.layerAlias)
	return layerEntry.ScreenXLocation, layerEntry.ScreenYLocation
}

/*
GetSize is a method which allows you to retrieve the current width and height of the layer in characters.

Example:

	w, h := layerInstance.GetSize()
*/
func (shared *LayerInstanceType) GetSize() (int, int) {
	validateLayer(shared.layerAlias)
	layerEntry := Layers.Get(shared.layerAlias)
	return layerEntry.Width, layerEntry.Height
}

/*
GetAlpha is a method which allows you to retrieve the current alpha transparency level of the layer.

Example:

	alpha := layerInstance.GetAlpha()
*/
func (shared *LayerInstanceType) GetAlpha() float32 {
	validateLayer(shared.layerAlias)
	layerEntry := Layers.Get(shared.layerAlias)
	return layerEntry.DefaultAttribute.ForegroundAlphaValue
}

/*
LoadLayer is a method which allows you to load a pre-rendered layer from disk and add it to the layer system.

In addition, the following should be noted:

- This is useful for quickly loading complex layers that were previously saved, such as image layers.

- The layer is loaded from a compressed format that was created by SaveLayer.

- The file extension ".clayer" is automatically appended to the filename if not provided.

- If the file cannot be read or is not a valid layer file, an error is returned.

- The loaded layer is added to the layer system with the specified alias, position, and z-order.

- The function returns a LayerInstanceType that can be used to manipulate the loaded layer.

Example:

	err := layerInstance.LoadLayer(0, 0, 1, "image.clayer")
*/
func (shared *LayerInstanceType) LoadLayer(xLocation int, yLocation int, zOrderPriority int, filePath string) error {
	// Get the file data using the virtual file system
	fileData, err := getFileDataFromFileSystem(filePath)
	if err != nil {
		return err
	}
	layerEntry := Layers.Get(shared.layerAlias)
	layerEntry.LoadLayerFromBytes(fileData)
	return nil
}

/*
LoadPreRenderedLayerImage is a method which allows you to load a pre-rendered layer image directly into image memory.

In addition, the following should be noted:

  - This is different from loading an image and pre-rendering it afterwards, as it directly loads a layer
    that has already been rendered.

  - The file extension ".clayer" is automatically appended to the filename if not provided.

  - If the file cannot be read or is not a valid layer file, an error is returned.

  - The loaded layer is added to the image system with the specified alias.

Example:

	err := layerInstance.LoadPreRenderedLayerImage("pre.clayer", "myImage")
*/
func (shared *LayerInstanceType) LoadPreRenderedLayerImage(filePath string, imageAlias string) error {
	// Load the pre-rendered layer image using the image.go function
	return LoadPreRenderedLayerImage(filePath, imageAlias)
}

/*
SaveLayer is a method which allows you to save the current layer to disk.

In addition, the following should be noted:

- This is useful for caching complex layers that take time to render, such as image layers.

- The file extension ".clayer" is automatically appended to the filename if not provided.

- The layer is saved using gzip compression to minimize disk space.

- If the file cannot be written, an error is returned.

Example:

	err := layerInstance.SaveLayer("saved.clayer")
*/
func (shared *LayerInstanceType) SaveLayer(filePath string) error {
	validateLayer(shared.layerAlias)
	layerEntry := Layers.Get(shared.layerAlias)
	return layerEntry.SaveLayer(filePath)
}

/*
ColorStyle is a method which allows you to apply a TUI style entry to the current layer, setting its default foreground
and background colors.

Example:

	layerInstance.ColorStyle(myStyle)
*/
func (shared *LayerInstanceType) ColorStyle(styleEntry types.TuiStyleEntryType) {
	colorLayerInstance(shared, int(styleEntry.Text.ForegroundColor), int(styleEntry.Text.BackgroundColor))
}

/*
Resize is a method which allows you to change the width and height of a layer.

In addition, the following should be noted:

  - If you pass in a zero or negative value for either width or height a panic will be generated to fail as fast
    as possible.

Example:

	layerInstance.Resize(100, 50)
*/
func (shared *LayerInstanceType) Resize(width int, height int) {
	validateLayer(shared.layerAlias)
	validateLayerSize(shared.layerAlias, width, height)
	layerEntry := Layers.Get(shared.layerAlias)

	// Create a new character memory with the new dimensions
	newCharacterMemory := make([][]types.CharacterEntryType, height)
	for i := range newCharacterMemory {
		newCharacterMemory[i] = make([]types.CharacterEntryType, width)
		for j := range newCharacterMemory[i] {
			newCharacterMemory[i][j] = types.NewCharacterEntry()
			newCharacterMemory[i][j].AttributeEntry = layerEntry.DefaultAttribute
			newCharacterMemory[i][j].LayerAlias = layerEntry.LayerAlias
			newCharacterMemory[i][j].ParentAlias = layerEntry.ParentAlias
		}
	}

	// Copy the existing character memory to the new one
	copyCharacterMemory(layerEntry.CharacterMemory, newCharacterMemory, 0, 0, layerEntry.Width, layerEntry.Height)

	layerEntry.Width = width
	layerEntry.Height = height
	layerEntry.CharacterMemory = newCharacterMemory
}

/*
PrintDialog is a method which allows you to write text immediately to the terminal screen via a typewriter effect.

In addition, the following should be noted:

  - This is useful for video games or other applications that may require printing text in a dialog box.

  - If you specify a print location outside the range of your specified text layer, a panic will be generated
    to fail as fast as possible.

  - If printing has reached the last line of your text layer, printing will not advance to the next line.
    Instead, it will resume and overwrite what was already printed on the same line.

  - Specifying the width of your text line allows you to control when text wrapping occurs.

  - When a word is too long to be printed on a text layer line, or the width of your line has already exceed
    its allowed maximum, the word will be pushed to the line directly under it.

  - When specifying a printing delay, the amount of time to wait is inserted between each character printed.

  - If the dialog being printed is flagged as skipable, the user can speed up printing by pressing the 'enter'
    key or right mouse button.

  - This method supports the use of text styles during printing to add color or styles to specific words in
    your string. All text styles must be enclosed around the "{" and "}" characters.

Example:

	layerInstance.PrintDialog(0, 0, 30, 50, true, "Hello {red}World{}")
*/
func (shared *LayerInstanceType) PrintDialog(xLocation int, yLocation int, widthOfLineInCharacters int, printDelayInMilliseconds int, isSkipable bool, stringToPrint string) {
	formattedTextToPrint := fmt.Sprint(stringToPrint)
	layerEntry := Layers.Get(shared.layerAlias)
	if xLocation < 0 || xLocation > layerEntry.Width || yLocation < 0 || yLocation > layerEntry.Height {
		panic(fmt.Sprintf("The specified location (%d, %d) is out of bounds for layer '%s' with a size of (%d, %d).", xLocation, yLocation, layerEntry.LayerAlias, layerEntry.Width, layerEntry.Height))
	}
	printDialog(layerEntry, layerEntry.DefaultAttribute, xLocation, yLocation, widthOfLineInCharacters, printDelayInMilliseconds, isSkipable, formattedTextToPrint)
}

/*
PrintMarkup is a method which allows you to write text immediately to the terminal screen with word wrapping and
attribute tags.

In addition, the following should be noted:

  - This is similar to PrintDialog but without the typewriter effect and printing delay.

  - If you specify a print location outside the range of your specified text layer, a panic will be generated
    to fail as fast as possible.

  - If printing has reached the last line of your text layer, printing will not advance to the next line.
    Instead, it will resume and overwrite what was already printed on the same line.

  - Specifying the width of your text line allows you to control when text wrapping occurs.

  - When a word is too long to be printed on a text layer line, or the width of your line has already exceeded
    its allowed maximum, the word will be pushed to the line directly under it.

  - This method supports the use of text styles during printing to add color or styles to specific words in
    your string. All text styles must be enclosed around the "{{" and "}}" characters.

Example:

	layerInstance.PrintMarkup(0, 0, 30, "This is {red}red{} text.")
*/
func (shared *LayerInstanceType) PrintMarkup(xLocation int, yLocation int, widthOfLineInCharacters int, stringToPrint string) {
	formattedTextToPrint := fmt.Sprint(stringToPrint)
	layerEntry := Layers.Get(shared.layerAlias)
	if xLocation < 0 || xLocation > layerEntry.Width || yLocation < 0 || yLocation > layerEntry.Height {
		panic(fmt.Sprintf("The specified location (%d, %d) is out of bounds for layer '%s' with a size of (%d, %d).", xLocation, yLocation, layerEntry.LayerAlias, layerEntry.Width, layerEntry.Height))
	}
	layer.printMarkup(layerEntry, layerEntry.DefaultAttribute, xLocation, yLocation, widthOfLineInCharacters, formattedTextToPrint)
}

/*
PrintFont is a method which allows you to render a string onto a layer using the specified font.

Example:

	layerInstance.PrintFont(myFont, 10, 5, "Big Text")
*/
func (shared *LayerInstanceType) PrintFont(fontInstance fontInstanceType, xLocation int, yLocation int, stringToPrint string) {
	layerEntry := Layers.Get(shared.layerAlias)
	if layerEntry == nil {
		panic(fmt.Sprintf("Layer with alias '%s' not found.", shared.layerAlias))
	}
	Font.PrintText(layerEntry, fontInstance, xLocation, yLocation, stringToPrint)
}

/*
PrintFontDialog is a method which allows you to write text to the terminal screen with a typewriter effect using a
specified font.

In addition, the following should be noted:

  - This is useful for creating animated text sequences with custom fonts.

  - If you specify a print location outside the range of your specified text layer, a panic will be generated.

  - When specifying a printing delay, the amount of time to wait is inserted between each character printed.

  - If the dialog being printed is flagged as skippable, the user can speed up printing by pressing the 'enter'
    key or right mouse button.

  - Specifying the width of your text line in characters allows you to control when text wrapping occurs.

Example:

	layerInstance.PrintFontDialog(myFont, 0, 0, 30, 50, true, "Animated font text")
*/
func (shared *LayerInstanceType) PrintFontDialog(fontInstance fontInstanceType, xLocation int, yLocation int, widthOfLineInCharacters int, printDelayInMilliseconds int, isSkipable bool, stringToPrint string) {
	formattedTextToPrint := fmt.Sprint(stringToPrint)
	layerEntry := Layers.Get(shared.layerAlias)
	if xLocation < 0 || xLocation > layerEntry.Width || yLocation < 0 || yLocation > layerEntry.Height {
		panic(fmt.Sprintf("The specified location (%d, %d) is out of bounds for layer '%s' with a size of (%d, %d).", xLocation, yLocation, layerEntry.LayerAlias, layerEntry.Width, layerEntry.Height))
	}
	Font.PrintTextDialog(layerEntry, fontInstance, xLocation, yLocation, widthOfLineInCharacters, printDelayInMilliseconds, isSkipable, formattedTextToPrint)
}

/*
AddLayer is a method which allows you to add a text layer to the current terminal display.

In addition, the following should be noted:

  - You can add as many layers as you wish to suit your applications needs.

  - Text layers are useful for setting up windows, modal dialogs, viewports, game foregrounds and backgrounds,
    and even effects like parallax scrolling.

  - If you specify a location for your layer that is outside the visible terminal display, then only the
    visible portion will be rendered.

  - If you pass in a zero or negative value for either width or height a panic will be generated to fail as fast
    as possible.

  - The z-order priority controls which text layer should be drawn first and which text layer should be drawn
    last.

  - The parent layer instance specifies which text layer is the parent of the one being created.

  - When adding a new text layer, it will become the default working text layer automatically.

Example:

	layerInstance := AddLayer(0, 0, 80, 25, 1, nil)
*/
func AddLayer(xLocation int, yLocation int, width int, height int, zOrderPriority int, parentLayerInstance *LayerInstanceType) *LayerInstanceType {
	layerAlias := getUUID()
	validateTerminalWidthAndHeight(width, height)
	var parentAlias string
	if parentLayerInstance != nil {
		parentAlias = parentLayerInstance.layerAlias
	}
	layer.Add(layerAlias, xLocation, yLocation, width, height, zOrderPriority, parentAlias)
	layerInstance := LayerInstanceType{layerAlias: layerAlias, parentAlias: parentAlias}
	return &layerInstance
}

/*
deleteLayer is a method which allows you to remove a text layer.

In addition, the following should be noted:

  - If any dynamically drawn TUI controls reference the deleted layer, they will still be present but no
    longer rendered.

  - If you attempt to delete a text layer which is currently set as your default text layer, then a panic will
    be generated.

  - If you attempt to delete a text layer that does not exist, then the operation will be ignored.

Example:

	deleteLayer("myLayer")
*/
func deleteLayer(layerAlias string) {
	validateLayer(layerAlias)
	layer.Delete(layerAlias)
}

/*
moveLayerByAbsoluteValue is a method which allows you to move a layer to an absolute screen position.

Example:

	moveLayerByAbsoluteValue("myLayer", 10, 5)
*/
func moveLayerByAbsoluteValue(layerAlias string, xLocation int, yLocation int) {
	validateLayer(layerAlias)
	layerEntry := Layers.Get(layerAlias)
	layerEntry.ScreenXLocation = xLocation
	layerEntry.ScreenYLocation = yLocation
}

/*
moveLayerByRelativeValue is a method which allows you to move a layer relative to its current screen position.

Example:

	moveLayerByRelativeValue("myLayer", 1, -1)
*/
func moveLayerByRelativeValue(layerAlias string, xLocation int, yLocation int) {
	validateLayer(layerAlias)
	layerEntry := Layers.Get(layerAlias)
	layerEntry.ScreenXLocation += xLocation
	layerEntry.ScreenYLocation += yLocation
}

/*
DeleteAllLayers is a method which allows you to remove all layers from memory and reinitialize screen memory.

Example:

	DeleteAllLayers()
*/
func DeleteAllLayers() {
	for _, entryToRemove := range Layers.GetAllEntries() {
		if !Layers.IsExists(entryToRemove.LayerAlias) {
			continue
		}
		layer.Delete(entryToRemove.LayerAlias)
	}
	layer.ReInitializeScreenMemory()
}

/*
SetTopmostLayer is a method which allows you to set the specified layer to be the topmost layer among its siblings.

Example:

	SetTopmostLayer(layerInstance)
*/
func SetTopmostLayer(layerInstance *LayerInstanceType) {
	layer.SetTopmostLayer(layerInstance.layerAlias)
}

/*
isLayerExists is a method which allows you to check if a layer with the given alias exists.

Example:

	exists := isLayerExists("myLayer")
*/
func isLayerExists(layerAlias string) bool {
	if Layers.IsExists(layerAlias) {
		return true
	}
	return false
}

/*
setLayerIsVisible is a method which allows you to set the visibility state of a layer.

Example:

	setLayerIsVisible("myLayer", true)
*/
func setLayerIsVisible(layerAlias string, isVisible bool) {
	validateLayer(layerAlias)
	layerEntry := Layers.Get(layerAlias)
	layerEntry.IsVisible = isVisible
}

/*
setLayerAlphaValue is a method which allows you to set the alpha value for a layer. In addition, the following should
be noted:

- This internal function updates the underlying layer entry memory.

Example:
    setLayerAlphaValue("myLayer", 0.5)
*/
/*
setLayerAlphaValue is a method which updates the underlying transparency property of a layer entry in memory. It
directly modifies the alpha value of the layer entry identified by the provided alias, affecting all subsequent
rendering operations for that layer. In addition, the following should be noted:

- This function triggers a full layer validation check before performing the update.

Example:
    setLayerAlphaValue("myLayer", 0.5)
*/
func setLayerAlphaValue(layerAlias string, alphaValue float32) {
	validateLayer(layerAlias)
	layerEntry := Layers.Get(layerAlias)
	layerEntry.AlphaValue = alphaValue
}

/*
getLayerAlphaValue is a method which retrieves the transparency multiplier directly from the layer entry memory. It
accesses the layer entry via its alias and returns the current float32 alpha value. In addition, the following should
be noted:

- This function performs a layer validation check to ensure the entry exists before access.

Example:
    alpha := getLayerAlphaValue("myLayer")
*/
func getLayerAlphaValue(layerAlias string) float32 {
	validateLayer(layerAlias)
	layerEntry := Layers.Get(layerAlias)
	return layerEntry.AlphaValue
}

/*
setLayerTransparencyStrategy is a method which updates the dithering or blending algorithm assigned to a layer entry.
It changes how individual cells are evaluated for visibility during composition, allowing for dynamic transitions
between opaque and transparent states. In addition, the following should be noted:

- This function performs a layer validation check to ensure the entry exists before access.

Example:
    setLayerTransparencyStrategy("myLayer", constants.TransparencyStrategyBayer)
*/
func setLayerTransparencyStrategy(layerAlias string, strategy constants.TransparencyStrategy) {
	validateLayer(layerAlias)
	layerEntry := Layers.Get(layerAlias)
	layerEntry.TransparencyStrategy = strategy
}

/*
getLayerTransparencyStrategy is a method which retrieves the active transparency algorithm from a layer entry's memory.
It returns the TransparencyStrategy constant currently determining the layer's character-level rendering behavior. In
addition, the following should be noted:

- This function performs a layer validation check to ensure the entry exists before access.

Example:
    strategy := getLayerTransparencyStrategy("myLayer")
*/
func getLayerTransparencyStrategy(layerAlias string) constants.TransparencyStrategy {
	validateLayer(layerAlias)
	layerEntry := Layers.Get(layerAlias)
	return layerEntry.TransparencyStrategy
}

/*
validateLayerSize is a method which allows you to check if the given width and height are valid for a layer.

Example:

	validateLayerSize("myLayer", 80, 25)
*/
func validateLayerSize(layerAlias string, width int, height int) {
	if width <= 0 {
		panic(fmt.Sprintf("The layer '%s' could not be resized since a width of '%d' was specified!", layerAlias, width))
	}
	if height <= 0 {
		panic(fmt.Sprintf("The layer '%s' could not be resized since a height of '%d' was specified!", layerAlias, height))
	}
}
