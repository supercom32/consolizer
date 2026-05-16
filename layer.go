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

type layerType struct{}

var layer layerType
var Layers *memory.MemoryManager[types.LayerEntryType]

type LayerInstanceType struct {
	layerAlias  string
	parentAlias string
}
type layerAliasZOrderPair struct {
	Key   string
	Value int
}
type LayerAliasZOrderPairList []layerAliasZOrderPair

func init() {
	layer.ReInitializeScreenMemory()
}

func clearLayerInstance(layerInstance *LayerInstanceType) {
	layerEntry := Layers.Get(layerInstance.layerAlias)
	layer.clear(layerEntry)
}

/*
clear is a method which allows you to empty the specified text layer of all its contents. This is useful for
internal methods that want to clear a text layer directly.

:param layerEntry: A pointer to the layer entry structure that you wish to clear.

Example:

	layer.clear(layerEntry)
*/
func (shared *layerType) clear(layerEntry *types.LayerEntryType) {
	types.InitializeCharacterMemory(layerEntry)
}

/*
processMarkupTag is a method which allows you to process a markup tag in the text and update the attribute entry
accordingly. It returns the updated attribute entry and the new character index after the tag. In addition, the
following should be noted:

- If no valid closing tag is found, the opening tag is treated as regular text.

- Special tag "/" resets to the default attribute entry.

:param textToPrint: A slice of runes representing the text being processed.
:param currentIndex: The current index in the rune slice where the tag starts.
:param textString: The string representation of the text being processed.
:param defaultAttributeEntry: The default attribute entry to use as a fallback.

:return: The updated attribute entry and the new character index.

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
handleWordWrap is a method which allows you to manage word wrapping logic when a space character is encountered. It
returns the updated cursor positions after applying word wrap if needed. In addition, the following should be noted:

- Returns early if word wrapping is disabled (wordWrapWidth <= 0).

- Wraps to the next line if the word would exceed the word wrap width or layer width.

:param cursorX: The current X position of the cursor.
:param cursorY: The current Y position of the cursor.
:param xLocation: The starting X coordinate for the text.
:param wordWidth: The width of the word being processed.
:param wordWrapWidth: The width at which word wrapping should occur.
:param layerWidth: The total width of the layer.
:param layerHeight: The total height of the layer.

:return: The updated X and Y cursor positions.

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
skipped. This is typically used with word wrapping to avoid leading spaces after a wrap. In addition, the following
should be noted:

- Only applies when word wrapping is enabled and the character is a space at the start of a line.

:param wordWrapWidth: The width at which word wrapping is occurring.
:param character: The rune representing the character to check.
:param cursorX: The current X position of the cursor.
:param xLocation: The starting X coordinate for the text.

:return: A boolean indicating whether the leading space should be skipped.

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

:param y: The Y coordinate to check.
:param height: The height of the layer.

:return: A boolean indicating whether the position is within vertical bounds.

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

:param x: The X coordinate to check.
:param width: The width of the layer.

:return: A boolean indicating whether the position is within horizontal bounds.

Example:

	isValid := layer.isWithinHorizontalBounds(10, 80)
*/
func (shared *layerType) isWithinHorizontalBounds(x, width int) bool {
	return x >= 0 && x < width
}

/*
renderCharacter is a method which allows you to render a character at the specified position with the given attributes.
It handles wide characters and background transparency. In addition, the following should be noted:

- For wide characters, it occupies two character cells.

- Preserves the original background color when transparency is enabled.

:param characterMemory: The character memory structure of the layer.
:param cursorX: The X position to render the character at.
:param cursorY: The Y position to render the character at.
:param character: The rune representing the character to render.
:param attributeEntry: The attribute entry containing styling information.

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
advanceCursor is a method which allows you to move the cursor position after rendering a character. It handles line
wrapping when the cursor reaches the end of a line. In addition, the following should be noted:

- When word wrapping is enabled, wraps to the next line when reaching the layer width.

- When word wrapping is disabled, stops at the layer width.

:param cursorX: The current X position of the cursor.
:param cursorY: The current Y position of the cursor.
:param xLocation: The starting X coordinate for the text.
:param layerWidth: The total width of the layer.
:param wordWrapWidth: The width at which word wrapping should occur.

:return: The updated X and Y cursor positions.

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
print is a method which allows you to handle all types of printing with configurable options. It supports word wrapping
and markup/styling based on the provided options. In addition, the following should be noted:

- Handles boundary checking to ensure text stays within the layer.

- Processes markup tags when useMarkup is true.

- Supports word wrapping when wordWrapWidth > 0.

:param layerEntry: A pointer to the layer entry structure where text will be printed.
:param attributeEntry: The attribute entry containing styling information.
:param xLocation: The starting X coordinate for the text.
:param yLocation: The starting Y coordinate for the text.
:param textToPrint: A slice of runes representing the text to print.
:param wordWrapWidth: The width at which word wrapping should occur. Set to 0 to disable.
:param useMarkup: A boolean indicating whether to process markup tags.

:return: The final cursor X position relative to the starting position.

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

:param layerEntry: A pointer to the layer entry structure where text will be printed.
:param attributeEntry: The attribute entry containing styling information.
:param xLocation: The starting X coordinate for the text.
:param yLocation: The starting Y coordinate for the text.
:param width: The width at which word wrapping should occur.
:param textToPrint: A slice of runes representing the text to print.

:return: The final cursor X position relative to the starting position.

Example:

	finalX := layer.printWithWordWrap(layerEntry, attr, 0, 0, 20, rune("Hello"))
*/
func (shared *layerType) printWithWordWrap(layerEntry *types.LayerEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int, width int, textToPrint []rune) int {
	return shared.print(layerEntry, attributeEntry, xLocation, yLocation, textToPrint, width, false)
}

/*
calculateWordWidth is a method which allows you to calculate the width of a word from a given position. In addition, the
following should be noted:

- The first position is always assumed to be ' ' and is skipped.

- When markup is enabled, it processes the text to exclude markup tags from the width calculation.

- Returns the number of characters until the next space or end of text.

:param textToPrint: A slice of runes representing the text being processed.
:param start: The starting index in the rune slice.
:param useMarkup: A boolean indicating whether to account for markup tags.

:return: The calculated width of the word.

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
printLayer is a method which allows you to write text to a text layer. This is useful for internal methods that want to
write text to a text layer directly, without affecting user settings. In addition, the following should be noted:

- If the location to print falls outside the range of the text layer, then only the visible portion of your text will.

:param layerEntry: A pointer to the layer entry structure where text will be printed.
:param attributeEntry: The attribute entry containing styling information.
:param xLocation: The starting X coordinate for the text.
:param yLocation: The starting Y coordinate for the text.
:param textToPrint: A slice of runes representing the text to print.

:return: The final cursor X position relative to the starting position.

Example:

	finalX := layer.printLayer(layerEntry, attr, 0, 0, rune("Hello"))
*/
func (shared *layerType) printLayer(layerEntry *types.LayerEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int, textToPrint []rune) int {
	return shared.print(layerEntry, attributeEntry, xLocation, yLocation, textToPrint, 0, false)
}

/*
printMarkup is a method which allows you to write text to the terminal screen with word wrapping and attribute tags.
This is similar to printDialog but without the typewriter effect and printing delay.

:param layerEntry: A pointer to the layer entry structure where text will be printed.
:param attributeEntry: The attribute entry containing styling information.
:param xLocation: The starting X coordinate for the text.
:param yLocation: The starting Y coordinate for the text.
:param widthOfLineInCharacters: The width at which word wrapping should occur.
:param stringToPrint: The string content to print, potentially containing markup tags.

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

:param layerAlias: A unique alias for the new layer.
:param xLocation: The starting X coordinate for the layer's screen position.
:param yLocation: The starting Y coordinate for the layer's screen position.
:param width: The width of the layer in characters.
:param height: The height of the layer in characters.
:param zOrderPriority: The rendering priority. Higher values are rendered on top.
:param parentAlias: The alias of the parent layer, or an empty string if it has no parent.

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

:return: The alias of the next available layer, or an empty string if no layers exist.

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

:param layerAlias: The alias of the layer to be deleted.

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

:param parentAlias: The alias of the parent layer whose children will be deleted.

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

:param parentAlias: The alias of the layer to check.

:return: A boolean indicating whether the layer is a parent.

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

:return: A list of layer-alias and z-order pairs, sorted by z-order.

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

:param layerAlias: The alias of the layer to update.
:param parentAlias: The alias of the parent layer to scope the z-order update.

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

:param parentAlias: The alias of the parent layer.

:return: The highest z-order value found.

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

:param layerAlias: The alias of the layer to be set as topmost.

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

:param layerAlias: The alias of the layer to start the search from.
:param previousChildAlias: The alias of the previous child in the hierarchy.

:return: The alias of the root parent layer and the alias of its immediate child in the path.

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

:return: The string alias of the layer.

Example:

	alias := layerInstance.GetAlias()
*/
func (shared *LayerInstanceType) GetAlias() string {
	return shared.layerAlias
}

/*
DrawImage is a method which allows you to draw an image on a given text layer. This method supports various image
formats and drawing styles, allowing for flexible rendering of images as text-based art. In addition, the following
should be noted:

- If the image to be drawn is not already loaded in memory, it will be loaded automatically and then unloaded after the.

- When drawing images with transparencies, the transparent edges are only computed once against the layer. Moving the.

:param fileName: The path or filename of the image to draw.
:param drawingStyle: The style entry defining how the image should be rendered.
:param xLocation: The X-coordinate where the image should be drawn.
:param yLocation: The Y-coordinate where the image should be drawn.
:param widthInCharacters: The target width of the rendered image in characters.
:param heightInCharacters: The target height of the rendered image in characters.
:param blurSigma: The sigma value for Gaussian blur preprocessing.

:return: An error if the image could not be loaded or drawn.

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
		imageLayer = getImageLayer(imageData, drawingStyle, widthInCharacters, heightInCharacters, blurSigma)
	}
	drawImageToLayer(currentLayer, imageLayer, xLocation, yLocation)
	return err
}

/*
DrawComposedImage is a method which allows you to draw a composed image on a text layer using a specific drawing style.

:param imageComposeEntry: The composed image entry containing the image data to render.
:param drawingStyle: The visual style to apply during rendering.
:param xLocation: The X-coordinate for drawing.
:param yLocation: The Y-coordinate for drawing.
:param widthInCharacters: The target width in characters.
:param heightInCharacters: The target height in characters.

:return: An error if the rendering fails.

Example:

	err := layerInstance.DrawComposedImage(composedImage, style, 10, 5, 20, 10)
*/
func (shared *LayerInstanceType) DrawComposedImage(imageComposeEntry ImageComposerEntryType, drawingStyle types.ImageStyleEntryType, xLocation int, yLocation int, widthInCharacters int, heightInCharacters int) error {
	var err error
	var imageLayer types.LayerEntryType
	baseImage := imageComposeEntry.RenderImage()

	// Get the current layer to pass for transparency handling
	var currentLayer *types.LayerEntryType
	currentLayer = Layers.Get(shared.layerAlias)

	if drawingStyle.DrawingStyle == constants.ImageStyleHalfBlock {
		imageLayer = getImageLayerAsHalfBlock(baseImage, drawingStyle, widthInCharacters, heightInCharacters, drawingStyle.BlurSigmaIntensity)
	} else if drawingStyle.DrawingStyle == constants.ImageStyleCharacters {
		imageLayer = GetImageLayerAsAsciiColorArt(baseImage, drawingStyle, widthInCharacters, heightInCharacters, drawingStyle.BlurSigmaIntensity)
	} else if drawingStyle.DrawingStyle == constants.ImageStyleBlockElementsAccurate {
		imageLayer = getImageLayerAsBlockElementsAccurate(baseImage, drawingStyle, widthInCharacters, heightInCharacters, drawingStyle.BlurSigmaIntensity)
	} else if drawingStyle.DrawingStyle == constants.ImageStyleBlockElementsFast {
		imageLayer = getImageLayerAsBlockElementsFast(baseImage, drawingStyle, widthInCharacters, heightInCharacters, drawingStyle.BlurSigmaIntensity)
	} else if drawingStyle.DrawingStyle == constants.ImageStyleBraille {
		imageLayer = getImageLayerAsBraille(baseImage, drawingStyle, widthInCharacters, heightInCharacters, drawingStyle.BlurSigmaIntensity)
	} else if drawingStyle.DrawingStyle == constants.ImageStyleFullBlock {
		imageLayer = getImageLayerAsFullBlock(baseImage, drawingStyle, widthInCharacters, heightInCharacters, drawingStyle.BlurSigmaIntensity)
	} else {
		safeSttyPanic("Invalid image style rendering type!")
	}
	drawImageToLayer(currentLayer, imageLayer, xLocation, yLocation)
	return err
}

/*
AddButton is a method which allows you to add a new button control to the current layer.

:param buttonLabel: The text label to display on the button.
:param styleEntry: The visual style to apply to the button.
:param xLocation: The X-coordinate for the button's position.
:param yLocation: The Y-coordinate for the button's position.
:param width: The width of the button in characters.
:param height: The height of the button in characters.
:param isEnabled: A boolean indicating whether the button starts enabled.

:return: A ButtonInstanceType representing the new button.

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

:param checkboxLabel: The text label to display next to the checkbox.
:param styleEntry: The visual style to apply to the checkbox.
:param xLocation: The X-coordinate for the checkbox's position.
:param yLocation: The Y-coordinate for the checkbox's position.
:param isSelected: A boolean indicating whether the checkbox starts selected.
:param isEnabled: A boolean indicating whether the checkbox starts enabled.

:return: A CheckboxInstanceType representing the new checkbox.

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

:param styleEntry: The visual style to apply to the dropdown.
:param selectionEntry: The selection options for the dropdown.
:param xLocation: The X-coordinate for the dropdown's position.
:param yLocation: The Y-coordinate for the dropdown's position.
:param selectorHeight: The height of the expanded dropdown selector.
:param itemWidth: The width of each item in the dropdown.
:param defaultItemSelected: The index of the item selected by default.

:return: A DropdownInstanceType representing the new dropdown.

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

:param labelValue: The text content of the label.
:param styleEntry: The visual style to apply to the label.
:param xLocation: The X-coordinate for the label's position.
:param yLocation: The Y-coordinate for the label's position.
:param width: The width allocated for the label.

:return: A LabelInstanceType representing the new label.

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

:param progressBarLabel: The text label to display on or near the progress bar.
:param styleEntry: The visual style to apply to the progress bar.
:param xLocation: The X-coordinate for position.
:param yLocation: The Y-coordinate for position.
:param width: The width of the progress bar in characters.
:param height: The height of the progress bar in characters.
:param isVertical: A boolean indicating if the progress bar should be vertical.
:param value: The current value of the progress bar.
:param maxValue: The maximum possible value for the progress bar.
:param isBackgroundTransparent: A boolean indicating if the background should be transparent.

:return: A ProgressBarInstanceType representing the new progress bar.

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

:param radioButtonLabel: The text label for the radio button.
:param styleEntry: The visual style to apply.
:param xLocation: The X-coordinate for position.
:param yLocation: The Y-coordinate for position.
:param groupId: The ID of the group this radio button belongs to.
:param isSelected: A boolean indicating if the button starts selected.

:return: A RadioButtonInstanceType representing the new radio button.

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

:param styleEntry: The visual style to apply to the scrollbar.
:param xLocation: The X-coordinate for position.
:param yLocation: The Y-coordinate for position.
:param length: The length of the scrollbar in characters.
:param maxScrollValue: The maximum scrollable value.
:param scrollValue: The initial scroll value.
:param scrollIncrement: The amount to scroll per step.
:param isHorizontal: A boolean indicating if the scrollbar should be horizontal.

:return: A ScrollbarInstanceType representing the new scrollbar.

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

:param styleEntry: The visual style to apply to the selector.
:param selectionEntry: The selection options for the selector.
:param xLocation: The X-coordinate for position.
:param yLocation: The Y-coordinate for position.
:param selectorHeight: The height of the selector control.
:param itemWidth: The width of each selectable item.
:param numberOfColumns: The number of columns in the selector.
:param viewportPosition: The initial viewport offset.
:param selectedItem: The index of the initially selected item.
:param highlightOnClickOnly: A boolean for highlighting behavior.
:param isBorderDrawn: A boolean indicating if a border should be drawn.

:return: A SelectorInstanceType representing the new selector.

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

:param styleEntry: The visual style to apply to the text field.
:param xLocation: The X-coordinate for position.
:param yLocation: The Y-coordinate for position.
:param width: The width of the text field in characters.
:param maxLengthAllowed: The maximum character length allowed in the field.
:param isPasswordProtected: A boolean indicating if characters should be masked.
:param defaultValue: The initial text value of the field.
:param isEnabled: A boolean indicating if the field starts enabled.

:return: A TextFieldInstanceType representing the new text field.

Example:

	tf := layerInstance.AddTextField(style, 10, 5, 20, 50, false, "", true)
*/
func (shared *LayerInstanceType) AddTextField(styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, maxLengthAllowed int, isPasswordProtected bool, defaultValue string, isEnabled bool) TextFieldInstanceType {
	textFieldAlias := getUUID()
	textFieldInstance := TextField.Add(shared.layerAlias, textFieldAlias, styleEntry, xLocation, yLocation, width, maxLengthAllowed, isPasswordProtected, defaultValue, isEnabled)
	return textFieldInstance
}

/*
Add is a method which allows you to add a new multi-line textbox control to the current layer.

:param styleEntry: The visual style to apply to the textbox.
:param xLocation: The X-coordinate for position.
:param yLocation: The Y-coordinate for position.
:param width: The width of the textbox.
:param height: The height of the textbox.
:param isBorderDrawn: A boolean indicating if a border should be drawn.

:return: A TextboxInstanceType representing the new textbox.

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

:param tooltipValue: The text content of the tooltip.
:param styleEntry: The visual style to apply.
:param hotspotXLocation: The X-coordinate of the interactive hotspot.
:param hotspotYLocation: The Y-coordinate of the interactive hotspot.
:param hotspotWidth: The width of the hotspot area.
:param hotspotHeight: The height of the hotspot area.
:param tooltipXLocation: The X-coordinate where the tooltip appears.
:param tooltipYLocation: The Y-coordinate where the tooltip appears.
:param tooltipWidth: The width of the tooltip display.
:param tooltipHeight: The height of the tooltip display.
:param isLocationAbsolute: A boolean for coordinate mode.
:param isBorderDrawn: A boolean indicating if a border should be drawn.
:param hoverTime: The time in milliseconds the user must hover to show the tooltip.

:return: A TooltipInstanceType representing the new tooltip.

Example:

	tt := layerInstance.AddTooltip("Help text", style, 10, 5, 5, 1, 10, 6, 20, 3, false, true, 500)
*/
func (shared *LayerInstanceType) AddTooltip(tooltipValue string, styleEntry types.TuiStyleEntryType, hotspotXLocation int, hotspotYLocation int, hotspotWidth int, hotspotHeight int, tooltipXLocation int, tooltipYLocation int, tooltipWidth int, tooltipHeight int, isLocationAbsolute bool, isBorderDrawn bool, hoverTime int) TooltipInstanceType {
	tooltipAlias := getUUID()
	tooltipInstance := Tooltip.Add(shared.layerAlias, tooltipAlias, tooltipValue, styleEntry, hotspotXLocation, hotspotYLocation, hotspotWidth, hotspotHeight, tooltipXLocation, tooltipYLocation, tooltipWidth, tooltipHeight, isLocationAbsolute, isBorderDrawn, hoverTime)
	return tooltipInstance
}

/*
Add is a method which allows you to add a viewport to a given text layer. A viewport is a read-only text display
control that can show text with markup codes for colorization. It supports scrollback history and text wrapping. In
addition, the following should be noted:

- If vertical scrollbars are enabled, the viewport will maintain scrollback history up to the specified maxHistoryLines.

- If vertical scrollbars are not enabled, then no history is needed and only memory for the visible display is required.

  - If isLinesWrapped is enabled, text printed to screen will wrap text cleanly like dialog.go's printMarkup method and no
    horizontal scrollbars will be shown.

  - If isLinesWrapped is disabled, lines will remain on the same line and horizontal scrollbars will be rendered if
    needed.

- Text can be added to the viewport using the Print method.

:param styleEntry: The visual style to apply to the viewport.
:param xLocation: The X-coordinate for position.
:param yLocation: The Y-coordinate for position.
:param width: The width of the viewport.
:param height: The height of the viewport.
:param isLinesWrapped: A boolean indicating if lines should wrap.
:param isBorderDrawn: A boolean indicating if a border should be drawn.
:param maxHistoryLines: The maximum number of history lines to maintain.

:return: A ViewportInstanceType representing the new viewport.

Example:

	vp := layerInstance.AddViewport(style, 0, 0, 40, 10, true, true, 100)
*/
func (shared *LayerInstanceType) AddViewport(styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isLinesWrapped bool, isBorderDrawn bool, maxHistoryLines int) ViewportInstanceType {
	viewportAlias := getUUID()
	viewportInstance := viewport.Add(shared.layerAlias, viewportAlias, styleEntry, xLocation, yLocation, width, height, isLinesWrapped, isBorderDrawn, maxHistoryLines)
	return viewportInstance
}

/*
AddFileMenu is a method which allows you to add a file menu to a layer. In addition, the following should be noted:

- The file menu will be drawn at the specified location with the given style.

- Each heading in the menu can have its own dropdown with selectable items.

- The top level headings widths are always dynamic based on how large the heading is.

- The file menu reuses existing selectors for dropdown functionality.

:param styleEntry: The visual style to apply to the file menu.
:param menuHeadings: A list of string headings for the top-level menu.
:param menuSelections: A list of selection entries for each menu heading's dropdown.
:param xLocation: The X-coordinate for position.
:param yLocation: The Y-coordinate for position.
:param isEnabled: A boolean indicating if the menu starts enabled.

:return: A FileMenuInstanceType representing the new file menu.

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
Print is a method which allows you to write text to the current layer. In addition, the following should be noted:

- When text is written to the text layer, the cursor position is also updated to reflect its new location.

- If the string to print ends up being too long to fit at its current location, then only the visible portion of your.

:param textToPrint: The string of text to print.

Example:

	layerInstance.Print("Hello World")
*/
func (shared *LayerInstanceType) Print(textToPrint string) {
	printLayerInstance(shared, textToPrint)
}

/*
Locate is a method which allows you to set the default cursor location on your text layer for printing with. In
addition, the following should be noted:

  - If you pass in a location value that falls outside the dimensions of the default text layer, a panic will be
    generated.

- Valid text layer locations start at position (0,0) for the upper left corner.

:param xLocation: The x-axis location for the cursor.
:param yLocation: The y-axis location for the cursor.

Example:

	layerInstance.Locate(10, 5)
*/
func (shared *LayerInstanceType) Locate(xLocation int, yLocation int) {
	locateLayerInstance(shared, xLocation, yLocation)
}

/*
Color is a method which allows you to set default colors on your text layer for printing with. The color index specified
corresponds to the 16 color ANSI standard, where color 0 is Black and 15 is Bright White.

:param foregroundColorIndex: The ANSI index for the foreground color.
:param backgroundColorIndex: The ANSI index for the background color.

Example:

	layerInstance.Color(15, 0)
*/
func (shared *LayerInstanceType) Color(foregroundColorIndex int, backgroundColorIndex int) {
	colorLayerInstance(shared, foregroundColorIndex, backgroundColorIndex)
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

	layerInstance.ColorRGB(255, 255, 255, 0, 0, 0)
*/
func (shared *LayerInstanceType) ColorRGB(foregroundRedIndex int32, foregroundGreenIndex int32, foregroundBlueIndex int32, backgroundRedIndex int32, backgroundGreenIndex int32, backgroundBlueIndex int32) {
	colorLayerRGBInstance(shared, foregroundRedIndex, foregroundGreenIndex, foregroundBlueIndex, backgroundRedIndex, backgroundGreenIndex, backgroundBlueIndex)
}

/*
Color24Bit is a method which allows you to color the current layer using a 24-bit color expressed as an int32.

:param foregroundColor: The 24-bit foreground color.
:param backgroundColor: The 24-bit background color.

Example:

	layerInstance.Color24Bit(fgColor, bgColor)
*/
func (shared *LayerInstanceType) Color24Bit(foregroundColor constants.ColorType, backgroundColor constants.ColorType) {
	colorLayer24BitInstance(shared, foregroundColor, backgroundColor)
}

/*
SetAlpha is a method which allows you to set the alpha value for the current layer. This lets you perform pseudo
transparencies by making the layer foreground and background colors blend with the layers underneath it to the degree
specified. In addition, the following should be noted:

- An alpha value of 1.0 is equal to 100% visible, while an alpha value of 0.0 is 0% visible.

- If the percent change specified is outside of the RGB color range, then the color will simply bottom or max out.

:param alphaValue: The alpha value to set.

Example:

	layerInstance.SetAlpha(0.5)
*/
func (shared *LayerInstanceType) SetAlpha(alphaValue float32) {
	setLayerAlphaInstance(shared, alphaValue)
}

/*
SetZOrder is a method which allows you to set the z-order priority for the current layer.

:param zOrder: The z-order value to set.

Example:

	layerInstance.SetZOrder(10)
*/
func (shared *LayerInstanceType) SetZOrder(zOrder int) {
	setLayerZOrderInstance(shared, zOrder)
}

/*
DrawVerticalLine is a method which allows you to draw a vertical line on a text layer. This method also has the ability
to draw connectors in case the line intersects with other lines that have already been drawn. In addition, the following
should be noted:

  - If the line to be drawn falls outside the area of the text layer specified, then only the visible portion of the line
    will be drawn.

:param styleEntry: The visual style to apply to the line.
:param xLocation: The X-coordinate where the line starts.
:param yLocation: The Y-coordinate where the line starts.
:param height: The height of the vertical line.
:param isConnectorsDrawn: A boolean indicating if intersection connectors should be drawn.

Example:

	layerInstance.DrawVerticalLine(style, 10, 5, 10, true)
*/
func (shared *LayerInstanceType) DrawVerticalLine(styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, height int, isConnectorsDrawn bool) {
	layerEntry := Layers.Get(shared.layerAlias)
	localAttributeEntry := types.NewAttributeEntry()
	drawVerticalLine(layerEntry, styleEntry, localAttributeEntry, xLocation, yLocation, height, isConnectorsDrawn)
}

/*
DrawHorizontalLine is a method which allows you to draw a horizontal line on a text layer. This method also has the
ability to draw connectors in case the line intersects with other lines that have already been drawn. In addition, the
following should be noted:

  - If the line to be drawn falls outside the area of the text layer specified, then only the visible portion of the line
    will be drawn.

:param styleEntry: The visual style to apply to the line.
:param xLocation: The X-coordinate where the line starts.
:param yLocation: The Y-coordinate where the line starts.
:param width: The width of the horizontal line.
:param isConnectorsDrawn: A boolean indicating if intersection connectors should be drawn.

Example:

	layerInstance.DrawHorizontalLine(style, 10, 5, 20, true)
*/
func (shared *LayerInstanceType) DrawHorizontalLine(styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, isConnectorsDrawn bool) {
	layerEntry := Layers.Get(shared.layerAlias)
	localAttributeEntry := types.NewAttributeEntry()
	drawHorizontalLine(layerEntry, styleEntry, localAttributeEntry, xLocation, yLocation, width, isConnectorsDrawn)
}

/*
DrawBorder is a method which allows you to draw a border on a given text layer. Borders differ from frames since they
are flat shaded and do not have a raised or sunken look to them. In addition, the following should be noted:

- If the border to be drawn falls outside the range of the specified layer, then only the visible portion of the border.

- The 'isInteractive' option allows you to specify if the window should interact with the layer being drawn on. For.

:param styleEntry: The visual style to apply to the border.
:param xLocation: The X-coordinate for position.
:param yLocation: The Y-coordinate for position.
:param width: The width of the bordered area.
:param height: The height of the bordered area.
:param isInteractive: A boolean indicating if the border is interactive.

Example:

	layerInstance.DrawBorder(style, 5, 5, 40, 10, true)
*/
func (shared *LayerInstanceType) DrawBorder(styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isInteractive bool) {
	layerEntry := Layers.Get(shared.layerAlias)
	localAttributeEntry := types.NewAttributeEntry()
	drawBorder(layerEntry, styleEntry, localAttributeEntry, xLocation, yLocation, width, height, isInteractive)
}

/*
DrawFrameLabel is a method which allows you to draw a label for a frame. The label will be automatically enclosed by the
characters "" and "" to blend in with a border of a frame. In addition, the following should be noted:

- If the frame label to be drawn falls outside the range of the specified layer, then only the visible portion of the.

:param styleEntry: The visual style to apply to the label.
:param label: The text content of the label.
:param xLocation: The X-coordinate for position.
:param yLocation: The Y-coordinate for position.

Example:

	layerInstance.DrawFrameLabel(style, "Settings", 7, 5)
*/
func (shared *LayerInstanceType) DrawFrameLabel(styleEntry types.TuiStyleEntryType, label string, xLocation int, yLocation int) {
	layerEntry := Layers.Get(shared.layerAlias)
	drawFrameLabel(layerEntry, styleEntry, label, xLocation, yLocation)
}

/*
DrawFrame is a method which allows you to draw a frame on a given text layer. Frames differ from borders since borders
are flat shaded and do not have a raised or sunken look to them. In addition, the following should be noted:

- If the frame to be drawn falls outside the range of the specified layer, then only the visible portion of the frame.

- The 'isInteractive' option allows you to specify if the window should interact with the layer being drawn on. For.

:param styleEntry: The visual style to apply to the frame.
:param isRaised: A boolean indicating if the frame should look raised (true) or sunken (false).
:param xLocation: The X-coordinate for position.
:param yLocation: The Y-coordinate for position.
:param width: The width of the framed area.
:param height: The height of the framed area.
:param isInteractive: A boolean indicating if the frame is interactive.

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
DrawWindow is a method which allows you to draw a window on a given text layer. Windows differ from borders since the
entire area the window surrounds gets filled with a solid background color. In addition, the following should be noted:

- If the window to be drawn falls outside the range of the specified layer, then only the visible portion of the window.

- The 'isInteractive' option allows you to specify if the window should interact with the layer being drawn on. For.

:param styleEntry: The visual style to apply to the window.
:param xLocation: The X-coordinate for position.
:param yLocation: The Y-coordinate for position.
:param width: The width of the window.
:param height: The height of the window.
:param isInteractive: A boolean indicating if the window is interactive.

Example:

	layerInstance.DrawWindow(style, 5, 5, 40, 10, true)
*/
func (shared *LayerInstanceType) DrawWindow(styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isInteractive bool) {
	layerEntry := Layers.Get(shared.layerAlias)
	localAttributeEntry := types.NewAttributeEntry()
	drawWindow(layerEntry, styleEntry, localAttributeEntry, xLocation, yLocation, width, height, isInteractive)
}

/*
DrawShadow is a method which allows you to draw shadows on a given text layer. Shadows are simply transparent areas
which darken whatever text layers are underneath it by a specified degree. In addition, the following should be noted:

- The alpha value can range from 0.0 (no shadow) to 1.0 (totally black).

:param xLocation: The X-coordinate for position.
:param yLocation: The Y-coordinate for position.
:param width: The width of the shadow area.
:param height: The height of the shadow area.
:param alphaValue: The transparency level of the shadow.

Example:

	layerInstance.DrawShadow(7, 7, 40, 10, 0.5)
*/
func (shared *LayerInstanceType) DrawShadow(xLocation int, yLocation int, width int, height int, alphaValue float32) {
	layerEntry := Layers.Get(shared.layerAlias)
	localAttributeEntry := types.NewAttributeEntry()
	drawShadow(layerEntry, localAttributeEntry, xLocation, yLocation, width, height, alphaValue)
}

/*
FillArea is a method which allows you to fill an area of a given text layer with characters of your choice. If you wish
to fill the area with repeating text, simply provide the string you wish to repeat. In addition, the following should be
noted:

- If the area to fill falls outside the range of the specified layer, then only the visible portion of the fill will be.

:param fillCharacters: The string of characters to use for filling.
:param xLocation: The starting X-coordinate.
:param yLocation: The starting Y-coordinate.
:param width: The width of the area to fill.
:param height: The height of the area to fill.

Example:

	layerInstance.FillArea("*", 0, 0, 80, 25)
*/
func (shared *LayerInstanceType) FillArea(fillCharacters string, xLocation int, yLocation int, width int, height int) {
	layerEntry := Layers.Get(shared.layerAlias)
	attributeEntry := layerEntry.DefaultAttribute
	fillArea(layerEntry, attributeEntry, fillCharacters, xLocation, yLocation, width, height, constants.NullCellControlLocation)
}

/*
FillLayer is a method which allows you to fill an entire layer with characters of your choice. If you wish to fill the
layer with repeating text, simply provide the string you wish to repeat.

:param fillCharacters: The string of characters to use for filling.

Example:

	layerInstance.FillLayer(".")
*/
func (shared *LayerInstanceType) FillLayer(fillCharacters string) {
	layerEntry := Layers.Get(shared.layerAlias)
	attributeEntry := layerEntry.DefaultAttribute
	fillLayer(layerEntry, attributeEntry, fillCharacters)
}

/*
DrawBar is a method which allows you to draw a horizontal bar on a given text layer row. This is useful for drawing
application headers or status bar footers.

:param styleEntry: The visual style to apply to the bar.
:param xLocation: The starting X-coordinate.
:param yLocation: The starting Y-coordinate.
:param barLength: The length of the bar in characters.
:param fillCharacters: The characters to use for filling the bar.

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
MoveLayerByAbsoluteValue is a method which allows you to move a text layer by an absolute value. This is useful if you
know exactly what position you wish to move your text layer to. In addition, the following should be noted:

  - If you move your layer outside the visible terminal display, only the visible display area will be rendered.
    Likewise,.

:param xLocation: The new absolute X-coordinate for the layer.
:param yLocation: The new absolute Y-coordinate for the layer.

Example:

	layerInstance.MoveLayerByAbsoluteValue(10, 5)
*/
func (shared *LayerInstanceType) MoveLayerByAbsoluteValue(xLocation int, yLocation int) {
	moveLayerByAbsoluteValue(shared.layerAlias, xLocation, yLocation)
}

/*
MoveLayerByRelativeValue is a method which allows you to move a text layer by a relative value. This is useful for
windows, foregrounds, backgrounds, or any kind of animations or movement you may wish to do in increments. In addition,
the following should be noted:

  - If you move your layer outside the visible terminal display, only the visible display area will be rendered. Likewise,
    if your text layer is a child of a parent layer, then only the visible display area will be rendered on the parent.

:param xLocation: The relative X-coordinate offset.
:param yLocation: The relative Y-coordinate offset.

Example:

	layerInstance.MoveLayerByRelativeValue(-1, 2)
*/
func (shared *LayerInstanceType) MoveLayerByRelativeValue(xLocation int, yLocation int) {
	moveLayerByRelativeValue(shared.layerAlias, xLocation, yLocation)
}

/*
Delete is a method which allows you to remove a text layer. If you wish to reuse a text layer for a future purpose,
you may also consider making the layer invisible instead of deleting it. In addition, the following should be noted:

- When a text layer is deleted, all child text layers are recursively deleted as well.

- If any dynamically drawn TUI controls reference the deleted layer, they will still be present. However, because the.

  - If you attempt to delete a text layer which is currently set as your default text layer, then a panic will be
    generated.

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

:return: A boolean indicating whether the layer exists.

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

:param isVisible: A boolean indicating whether the layer should be visible.

Example:

	layerInstance.SetIsVisible(false)
*/
func (shared *LayerInstanceType) SetIsVisible(isVisible bool) {
	validateLayer(shared.layerAlias)
	setLayerIsVisible(shared.layerAlias, isVisible)
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

:return: The current X and Y coordinates of the layer.

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

:return: The width and height of the layer.

Example:
Example:

	alpha := layerInstance.GetAlpha()
*/
func (shared *LayerInstanceType) GetAlpha() float32 {
	validateLayer(shared.layerAlias)
	layerEntry := Layers.Get(shared.layerAlias)
	return layerEntry.DefaultAttribute.ForegroundAlphaValue
}

/*
LoadLayer is a method which allows you to load a pre-rendered layer from disk and add it to the layer system. This is
useful for quickly loading complex layers that were previously saved, such as image layers. The layer is loaded from a
compressed format that was created by SaveLayer. In addition, the following should be noted:

- The file extension ".clayer" is automatically appended to the filename if not provided.

- If the file cannot be read or is not a valid layer file, an error is returned.

- The loaded layer is added to the layer system with the specified alias, position, and z-order.

- The function returns a LayerInstanceType that can be used to manipulate the loaded layer.

:param xLocation: The X-coordinate for the layer's position.
:param yLocation: The Y-coordinate for the layer's position.
:param zOrderPriority: The rendering priority.
:param filePath: The path to the layer file.

:return: An error if the layer could not be loaded.

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
This is different from loading an image and pre-rendering it afterwards, as it directly loads a layer that has already
been rendered. In addition, the following should be noted:

- The file extension ".clayer" is automatically appended to the filename if not provided.

- If the file cannot be read or is not a valid layer file, an error is returned.

- The loaded layer is added to the image system with the specified alias.

:param filePath: The path to the layer file.
:param imageAlias: The alias to assign to the loaded image.

:return: An error if the image could not be loaded.

Example:

	err := layerInstance.LoadPreRenderedLayerImage("pre.clayer", "myImage")
*/
func (shared *LayerInstanceType) LoadPreRenderedLayerImage(filePath string, imageAlias string) error {
	// Load the pre-rendered layer image using the image.go function
	return LoadPreRenderedLayerImage(filePath, imageAlias)
}

/*
SaveLayer is a method which allows you to save the current layer to disk. This is useful for caching complex layers that
take time to render, such as image layers. In addition, the following should be noted:

- The file extension ".clayer" is automatically appended to the filename if not provided.

- The layer is saved using gzip compression to minimize disk space.

- If the file cannot be written, an error is returned.

:param filePath: The path where the layer should be saved.

:return: An error if the layer could not be saved.

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

:param styleEntry: The visual style to apply.

Example:

	layerInstance.ColorStyle(myStyle)
*/
func (shared *LayerInstanceType) ColorStyle(styleEntry types.TuiStyleEntryType) {
	colorLayerInstance(shared, int(styleEntry.Text.ForegroundColor), int(styleEntry.Text.BackgroundColor))
}

/*
Resize is a method which allows you to change the width and height of a layer. In addition, the following should be
noted:

- If you pass in a zero or negative value for either width or height a panic will be generated to fail as fast as.

:param width: The new width for the layer.
:param height: The new height for the layer.

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
PrintDialog is a method which allows you to write text immediately to the terminal screen via a typewriter effect. This
is useful for video games or other applications that may require printing text in a dialog box. In addition, the
following should be noted:

  - If you specify a print location outside the range of your specified text layer, a panic will be generated to fail as
    fast as possible.

  - If printing has reached the last line of your text layer, printing will not advance to the next line. Instead, it will
    resume and overwrite what was already printed on the same line.

- Specifying the width of your text line allows you to control when text wrapping occurs.

  - When a word is too long to be printed on a text layer line, or the width of your line has already exceed its allowed
    maximum, the word will be pushed to the line directly under it.

- When specifying a printing delay, the amount of time to wait is inserted between each character printed.

  - If the dialog being printed is flagged as skipable, the user can speed up printing by pressing the 'enter' key or
    right mouse button.

  - This method supports the use of text styles during printing to add color or styles to specific words in your string.
    All text styles must be enclosed around the "{" and "}" characters.

:param xLocation: The starting X-coordinate for the dialog.
:param yLocation: The starting Y-coordinate for the dialog.
:param widthOfLineInCharacters: The width at which word wrapping should occur.
:param printDelayInMilliseconds: The delay in milliseconds between each character.
:param isSkipable: A boolean indicating if the animation can be skipped.
:param stringToPrint: The text content to print, potentially containing style tags.

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
attribute tags. This is similar to PrintDialog but without the typewriter effect and printing delay. In addition, the
following should be noted:

  - If you specify a print location outside the range of your specified text layer, a panic will be generated to fail as
    fast as possible.

  - If printing has reached the last line of your text layer, printing will not advance to the next line. Instead, it will
    resume and overwrite what was already printed on the same line.

- Specifying the width of your text line allows you to control when text wrapping occurs.

  - When a word is too long to be printed on a text layer line, or the width of your line has already exceeded its allowed
    maximum, the word will be pushed to the line directly under it.

  - This method supports the use of text styles during printing to add color or styles to specific words in your string.
    All text styles must be enclosed around the "{" and "}" characters.

:param xLocation: The starting X-coordinate.
:param yLocation: The starting Y-coordinate.
:param widthOfLineInCharacters: The width at which word wrapping should occur.
:param stringToPrint: The text content to print, potentially containing style tags.

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

:param fontInstance: The font instance to use for rendering.
:param xLocation: The X-coordinate for position.
:param yLocation: The Y-coordinate for position.
:param stringToPrint: The text content to render.

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
specified font. This is useful for creating animated text sequences with custom fonts. In addition, the following should
be noted:

- If you specify a print location outside the range of your specified text layer, a panic will be generated.

- When specifying a printing delay, the amount of time to wait is inserted between each character printed.

- If the dialog being printed is flagged as skippable, the user can speed up printing by pressing the 'enter' key or.

- Specifying the width of your text line in characters allows you to control when text wrapping occurs.

:param fontInstance: The font instance to use for rendering.
:param xLocation: The starting X-coordinate.
:param yLocation: The starting Y-coordinate.
:param widthOfLineInCharacters: The width at which word wrapping should occur.
:param printDelayInMilliseconds: The delay in milliseconds between each character.
:param isSkipable: A boolean indicating if the animation can be skipped.
:param stringToPrint: The text content to print.

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
AddLayer is a method which allows you to add a text layer to the current terminal display. You can add as many layers as
you wish to suit your applications needs. Text layers are useful for setting up windows, modal dialogs, viewports, game
foregrounds and backgrounds, and even effects like parallax scrolling. In addition, the following should be noted:

  - If you specify a location for your layer that is outside the visible terminal display, then only the visible portion
    of.

- If you pass in a zero or negative value for either width or height a panic will be generated to fail as fast as.

  - The z-order priority controls which text layer should be drawn first and which text layer should be drawn last.
    Layers.

- The parent layer instance specifies which text layer is the parent of the one being created. Having a parent layer.

- When adding a new text layer, it will become the default working text layer automatically.

:param xLocation: The X-coordinate for the layer's screen position.
:param yLocation: The Y-coordinate for the layer's screen position.
:param width: The width of the layer in characters.
:param height: The height of the layer in characters.
:param zOrderPriority: The rendering priority. Higher values are rendered on top.
:param parentLayerInstance: A pointer to the parent layer instance, or nil if it has no parent.

:return: A pointer to the newly created LayerInstanceType.

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
MoveLayerByAbsoluteValue is a method which allows you to move a text layer by an absolute value. This is useful if you

- If any dynamically drawn TUI controls reference the deleted layer, they will still be present but no longer rendered.

  - If you attempt to delete a text layer which is currently set as your default text layer, then a panic will be
    generated.

- If you attempt to delete a text layer that does not exist, then the operation will be ignored.

:param layerAlias: The alias of the layer to be deleted.

Example:

	deleteLayer("myLayer")
*/
func deleteLayer(layerAlias string) {
	validateLayer(layerAlias)
	layer.Delete(layerAlias)
}

func moveLayerByAbsoluteValue(layerAlias string, xLocation int, yLocation int) {
	validateLayer(layerAlias)
	layerEntry := Layers.Get(layerAlias)
	layerEntry.ScreenXLocation = xLocation
	layerEntry.ScreenYLocation = yLocation
}

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

:param layerInstance: A pointer to the layer instance to be set as topmost.

Example:

	SetTopmostLayer(layerInstance)
*/
func SetTopmostLayer(layerInstance *LayerInstanceType) {
	layer.SetTopmostLayer(layerInstance.layerAlias)
}

func isLayerExists(layerAlias string) bool {
	if Layers.IsExists(layerAlias) {
		return true
	}
	return false
}

func setLayerIsVisible(layerAlias string, isVisible bool) {
	validateLayer(layerAlias)
	layerEntry := Layers.Get(layerAlias)
	layerEntry.IsVisible = isVisible
}

/*
validateLayerSize is a method which allows you to check if the given width and height are valid for a layer.

:param layerAlias: The alias of the layer being validated.
:param width: The width value to check.
:param height: The height value to check.

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
