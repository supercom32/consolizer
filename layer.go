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
	LayerWidth  int
	LayerHeight int
}

type layerAliasZOrderPair struct {
	Key   string
	Value int
}
type LayerAliasZOrderPairList []layerAliasZOrderPair

func init() {
	layer.ReInitializeScreenMemory()
}

/*
ClearLayer allows you to empty the specified text layer of all its contents. If you
do not wish to specify a text layer, you can use the method 'Clear' which will
simply clear the default text layer previously set.
*/
func ClearLayer(layerInstance *LayerInstanceType) {
	layerEntry := Layers.Get(layerInstance.layerAlias)
	layer.clearLayer(layerEntry)
}

/*
clearLayer allows you to empty the specified text layer of all its contents.
This is useful for internal methods that want to clear a text layer directly.
*/
func (shared *layerType) clearLayer(layerEntry *types.LayerEntryType) {
	types.InitializeCharacterMemory(layerEntry)
}

/*
processMarkupTag processes a markup tag in the text and updates the attribute entry accordingly.
It returns the updated attribute entry and the new character index after the tag.
In addition, the following information should be noted:

- If no valid closing tag is found, the opening tag is treated as regular text.
- Special tag "/" resets to the default attribute entry.
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
handleWordWrap manages word wrapping logic when a space character is encountered.
It returns the updated cursor positions after applying word wrap if needed.
In addition, the following information should be noted:

- Returns early if word wrapping is disabled (wordWrapWidth <= 0)
- Wraps to the next line if the word would exceed the word wrap width or layer width
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
shouldSkipLeadingSpace determines if a space character at the start of a line should be skipped.
This is typically used with word wrapping to avoid leading spaces after a wrap.
In addition, the following information should be noted:

- Only applies when word wrapping is enabled and the character is a space at the start of a line
*/
func (shared *layerType) shouldSkipLeadingSpace(wordWrapWidth int, character rune, cursorX, xLocation int) bool {
	return wordWrapWidth > 0 && character == ' ' && cursorX == xLocation
}

/*
isWithinVerticalBounds checks if the current cursor Y position is within the layer's height bounds.
In addition, the following information should be noted:

- A position is valid if it's greater than or equal to 0 and less than the layer height
*/
func (shared *layerType) isWithinVerticalBounds(y, height int) bool {
	return y >= 0 && y < height
}

/*
isWithinHorizontalBounds checks if the current cursor X position is within the layer's width bounds.
In addition, the following information should be noted:

- A position is valid if it's greater than or equal to 0 and less than the layer width
*/
func (shared *layerType) isWithinHorizontalBounds(x, width int) bool {
	return x >= 0 && x < width
}

/*
renderCharacter renders a character at the specified position with the given attributes.
It handles wide characters and background transparency.
In addition, the following information should be noted:

- For wide characters, it occupies two character cells
- Preserves the original background color when transparency is enabled
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
advanceCursor moves the cursor position after rendering a character.
It handles line wrapping when the cursor reaches the end of a line.
Returns the new cursor positions.
In addition, the following information should be noted:

- When word wrapping is enabled, wraps to the next line when reaching the layer width
- When word wrapping is disabled, stops at the layer width
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
print is a base method that handles all types of printing with configurable options.
It supports word wrapping and markup/styling based on the provided options.

Options:
- wordWrapWidth: Width for word wrapping. Set to 0 to disable word wrapping.
- useMarkup: Whether to process markup tags for styling.

Returns the final cursor X position relative to the starting position.
In addition, the following information should be noted:

- Handles boundary checking to ensure text stays within the layer
- Processes markup tags when useMarkup is true
- Supports word wrapping when wordWrapWidth > 0
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

func (shared *layerType) printLayerWithWordWrap(layerEntry *types.LayerEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int, width int, textToPrint []rune) int {
	return shared.print(layerEntry, attributeEntry, xLocation, yLocation, textToPrint, width, false)
}

/*
calculateWordWidth calculates the width of a word from the given position.
The first position is always assumed to be ' ' and is skipped.
In addition, the following information should be noted:

- When markup is enabled, it processes the text to exclude markup tags from the width calculation
- Returns the number of characters until the next space or end of text
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
printLayer allows you to write text to a text layer. This is useful
for internal methods that want to write text to a text layer directly, without
effecting user settings (such as current cursor location, etc). In addition,
the following information should be noted:

- If the location to print falls outside the range of the text layer,
then only the visible portion of your text will be printed.
*/
func (shared *layerType) printLayer(layerEntry *types.LayerEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int, textToPrint []rune) int {
	return shared.print(layerEntry, attributeEntry, xLocation, yLocation, textToPrint, 0, false)
}

/*
printMarkup allows you to write text to the terminal screen with word wrapping
and attribute tags. This is similar to printDialog but without the typewriter
effect and printing delay.
*/
func (shared *layerType) printMarkup(layerEntry *types.LayerEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int, widthOfLineInCharacters int, stringToPrint string) {
	arrayOfRunes := stringformat.GetRunesFromString(stringToPrint)
	shared.print(layerEntry, attributeEntry, xLocation, yLocation, arrayOfRunes, widthOfLineInCharacters, true)
}

func (shared *layerType) ReInitializeScreenMemory() {
	Layers = memory.NewMemoryManager[types.LayerEntryType]() // Initialize MemoryManager
}

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

// GetNextLayerAlias retrieves the next available layer alias.
func (shared *layerType) GetNextLayerAlias() string {
	for _, currentEntry := range Layers.GetAllEntries() {
		return currentEntry.LayerAlias
	}
	return ""
}

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

func (shared *layerType) IsAParent(parentAlias string) bool {
	isParent := false
	for _, currentValue := range Layers.GetAllEntries() {
		if currentValue.ParentAlias == parentAlias {
			isParent = true
		}
	}
	return isParent
}

// GetSortedLayerMemoryAliasSlice returns a sorted list of layer aliases based on z-order.
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

// SetHighestZOrderNumber sets the highest z-order number for the given layer.
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

func (shared *layerType) getHighestZOrderNumber(parentAlias string) int {
	highestZOrderNumber := 0
	for _, currentValue := range Layers.GetAllEntries() {
		if currentValue.ParentAlias == parentAlias && currentValue.ZOrder > highestZOrderNumber {
			highestZOrderNumber = currentValue.ZOrder
		}
	}
	return highestZOrderNumber
}

// SetTopmostLayer sets the specified layer to be the topmost layer.
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

func (shared *layerType) GetRootParentLayerAlias(layerAlias string, previousChildAlias string) (string, string) {
	layerEntry := Layers.Get(layerAlias)
	if layerEntry.ParentAlias != "" {
		childToTrack := previousChildAlias
		if childToTrack == "" {
			childToTrack = layerAlias
		}
		return shared.GetRootParentLayerAlias(layerEntry.ParentAlias, childToTrack)
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

func (shared *LayerInstanceType) Clear() {
	layerEntry := Layers.Get(shared.layerAlias)
	localAttributeEntry := types.NewAttributeEntry()
	fillArea(layerEntry, localAttributeEntry, "", 0, 0, shared.LayerWidth, shared.LayerHeight, 0)
}

/*
DrawImage allows you to draw an image on a given text layer. This method
supports various image formats and drawing styles, allowing for flexible
rendering of images as text-based art. In addition, the following information
should be noted:

- If the image to be drawn is not already loaded in memory, it will be
loaded automatically and then unloaded after the drawing operation is complete.

- When drawing images with transparencies, the transparent edges are only
computed once against the layer. Moving the layer around or changing the
underlying layers will not make the image update.
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

func (shared *LayerInstanceType) AddButton(buttonLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isEnabled bool) buttonInstanceType {
	buttonAlias := getUUID()
	buttonInstance := Button.Add(shared.layerAlias, buttonAlias, buttonLabel, styleEntry, xLocation, yLocation, width, height, isEnabled)
	return buttonInstance
}

func (shared *LayerInstanceType) AddCheckbox(checkboxLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, isSelected bool, isEnabled bool) checkboxInstanceType {
	checkboxAlias := getUUID()
	checkboxInstance := Checkbox.Add(shared.layerAlias, checkboxAlias, checkboxLabel, styleEntry, xLocation, yLocation, isSelected, isEnabled)
	return checkboxInstance
}

func (shared *LayerInstanceType) AddDropdown(styleEntry types.TuiStyleEntryType, selectionEntry types.SelectionEntryType, xLocation int, yLocation int, selectorHeight int, itemWidth int, defaultItemSelected int) DropdownInstanceType {
	dropdownAlias := getUUID()
	dropdownInstance := Dropdown.Add(shared.layerAlias, dropdownAlias, styleEntry, selectionEntry, xLocation, yLocation, selectorHeight, itemWidth, defaultItemSelected)
	return dropdownInstance
}

func (shared *LayerInstanceType) AddLabel(labelValue string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int) LabelInstanceType {
	labelAlias := getUUID()
	labelInstance := Label.Add(shared.layerAlias, labelAlias, labelValue, styleEntry, xLocation, yLocation, width)
	return labelInstance
}

func (shared *LayerInstanceType) AddProgressBar(progressBarLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isVertical bool, value int, maxValue int, isBackgroundTransparent bool) ProgressBarInstanceType {
	progressBarAlias := getUUID()
	progressBarInstance := ProgressBar.Add(shared.layerAlias, progressBarAlias, progressBarLabel, styleEntry, xLocation, yLocation, width, height, isVertical, value, maxValue, isBackgroundTransparent)
	return progressBarInstance
}

func (shared *LayerInstanceType) AddRadioButton(radioButtonLabel string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, groupId int, isSelected bool) RadioButtonInstanceType {
	radioButtonAlias := getUUID()
	radioButtonInstance := radioButton.Add(shared.layerAlias, radioButtonAlias, radioButtonLabel, styleEntry, xLocation, yLocation, groupId, isSelected)
	return radioButtonInstance
}

func (shared *LayerInstanceType) AddScrollbar(styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, length int, maxScrollValue int, scrollValue int, scrollIncrement int, isHorizontal bool) ScrollbarInstanceType {
	scrollbarAlias := getUUID()
	scrollbarInstance := scrollbar.Add(shared.layerAlias, scrollbarAlias, styleEntry, xLocation, yLocation, length, maxScrollValue, scrollValue, scrollIncrement, isHorizontal)
	return scrollbarInstance
}

func (shared *LayerInstanceType) AddSelector(styleEntry types.TuiStyleEntryType, selectionEntry types.SelectionEntryType, xLocation int, yLocation int, selectorHeight int, itemWidth int, numberOfColumns int, viewportPosition int, selectedItem int, highlightOnClickOnly bool, isBorderDrawn bool) SelectorInstanceType {
	selectorAlias := getUUID()
	selectorInstance := Selector.Add(shared.layerAlias, selectorAlias, styleEntry, selectionEntry, xLocation, yLocation, selectorHeight, itemWidth, numberOfColumns, viewportPosition, selectedItem, highlightOnClickOnly, isBorderDrawn)
	return selectorInstance
}

func (shared *LayerInstanceType) AddTextField(styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, maxLengthAllowed int, isPasswordProtected bool, defaultValue string, isEnabled bool) textFieldInstanceType {
	textFieldAlias := getUUID()
	textFieldInstance := TextField.Add(shared.layerAlias, textFieldAlias, styleEntry, xLocation, yLocation, width, maxLengthAllowed, isPasswordProtected, defaultValue, isEnabled)
	return textFieldInstance
}

func (shared *LayerInstanceType) AddTextbox(styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isBorderDrawn bool) TextboxInstanceType {
	textBoxAlias := getUUID()
	textBoxInstance := textbox.AddTextbox(shared.layerAlias, textBoxAlias, styleEntry, xLocation, yLocation, width, height, isBorderDrawn)
	return textBoxInstance
}

func (shared *LayerInstanceType) AddTooltip(tooltipValue string, styleEntry types.TuiStyleEntryType, hotspotXLocation int, hotspotYLocation int, hotspotWidth int, hotspotHeight int, tooltipXLocation int, tooltipYLocation int, tooltipWidth int, tooltipHeight int, isLocationAbsolute bool, isBorderDrawn bool, hoverTime int) TooltipInstanceType {
	tooltipAlias := getUUID()
	tooltipInstance := Tooltip.Add(shared.layerAlias, tooltipAlias, tooltipValue, styleEntry, hotspotXLocation, hotspotYLocation, hotspotWidth, hotspotHeight, tooltipXLocation, tooltipYLocation, tooltipWidth, tooltipHeight, isLocationAbsolute, isBorderDrawn, hoverTime)
	return tooltipInstance
}

/*
AddViewport allows you to add a viewport to a given text layer. A viewport is a read-only
text display control that can show text with markup codes for colorization. It supports
scrollback history and text wrapping. In addition, the following information should be noted:

- If vertical scrollbars are enabled, the viewport will maintain scrollback history up to
the specified maxHistoryLines.
- If vertical scrollbars are not enabled, then no history is needed and only memory for
the visible display is required.
- If isLinesWrapped is enabled, text printed to screen will wrap text cleanly
like dialog.go's printMarkup method and no horizontal scrollbars will be shown.
- If isLinesWrapped is disabled, lines will remain on the same line and horizontal
scrollbars will be rendered if needed.
- Text can be added to the viewport using the Print method.
*/
func (shared *LayerInstanceType) AddViewport(styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isLinesWrapped bool, isBorderDrawn bool, maxHistoryLines int) ViewportInstanceType {
	viewportAlias := getUUID()
	viewportInstance := AddViewport(shared.layerAlias, viewportAlias, styleEntry, xLocation, yLocation, width, height, isLinesWrapped, isBorderDrawn, maxHistoryLines)
	return viewportInstance
}

/*
AddFileMenu allows you to add a file menu to a layer. In addition, the following
information should be noted:

- The file menu will be drawn at the specified location with the given style.
- Each heading in the menu can have its own dropdown with selectable items.
- The top level headings widths are always dynamic based on how large the heading is.
- The file menu reuses existing selectors for dropdown functionality.
*/
func (shared *LayerInstanceType) AddFileMenu(styleEntry types.TuiStyleEntryType, menuHeadings []string, menuSelections []types.SelectionEntryType, xLocation int, yLocation int, isEnabled bool) FileMenuInstanceType {
	menuAlias := getUUID()
	fileMenuInstance := FileMenu.Add(shared.layerAlias, menuAlias, styleEntry, menuHeadings, menuSelections, xLocation, yLocation, isEnabled)
	return fileMenuInstance
}

/*
DrawVerticalLine allows you to draw a vertical line on a text layer. This
method also has the ability to draw connectors in case the line intersects
with other lines that have already been drawn. In addition, the following
information should be noted:

- If the the line to be drawn falls outside the area of the text layer
specified, then only the visible portion of the line will be drawn.
*/
func (shared *LayerInstanceType) DrawVerticalLine(styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, height int, isConnectorsDrawn bool) {
	layerEntry := Layers.Get(shared.layerAlias)
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
func (shared *LayerInstanceType) DrawHorizontalLine(styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, isConnectorsDrawn bool) {
	layerEntry := Layers.Get(shared.layerAlias)
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
func (shared *LayerInstanceType) DrawBorder(styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isInteractive bool) {
	layerEntry := Layers.Get(shared.layerAlias)
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
func (shared *LayerInstanceType) DrawFrameLabel(styleEntry types.TuiStyleEntryType, label string, xLocation int, yLocation int) {
	layerEntry := Layers.Get(shared.layerAlias)
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
func (shared *LayerInstanceType) DrawWindow(styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isInteractive bool) {
	layerEntry := Layers.Get(shared.layerAlias)
	localAttributeEntry := types.NewAttributeEntry()
	drawWindow(layerEntry, styleEntry, localAttributeEntry, xLocation, yLocation, width, height, isInteractive)
}

/*
DrawShadow allows you to draw shadows on a given text layer. Shadows are simply
transparent areas which darken whatever text layers are underneath it by a
specified degree. In addition, the following information should be noted:

- The alpha value can range from 0.0 (no shadow) to 1.0 (totally black).
*/
func (shared *LayerInstanceType) DrawShadow(xLocation int, yLocation int, width int, height int, alphaValue float32) {
	layerEntry := Layers.Get(shared.layerAlias)
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
func (shared *LayerInstanceType) FillArea(fillCharacters string, xLocation int, yLocation int, width int, height int) {
	layerEntry := Layers.Get(shared.layerAlias)
	attributeEntry := layerEntry.DefaultAttribute
	fillArea(layerEntry, attributeEntry, fillCharacters, xLocation, yLocation, width, height, constants.NullCellControlLocation)
}

/*
FillLayer allows you to fill an entire layer with characters of your choice.
If you wish to fill the layer with repeating text, simply provide the string
you wish to repeat.
*/
func (shared *LayerInstanceType) FillLayer(fillCharacters string) {
	layerEntry := Layers.Get(shared.layerAlias)
	attributeEntry := layerEntry.DefaultAttribute
	fillLayer(layerEntry, attributeEntry, fillCharacters)
}

/*
DrawBar allows you to draw a horizontal bar on a given text layer row. This is
useful for drawing application headers or status bar footers.
*/
func (shared *LayerInstanceType) DrawBar(styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, barLength int, fillCharacters string) {
	layerEntry := Layers.Get(shared.layerAlias)
	attributeEntry := types.NewAttributeEntry()
	attributeEntry.ForegroundColor = styleEntry.Bar.ForegroundColor
	attributeEntry.BackgroundColor = styleEntry.Bar.BackgroundColor
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
func (shared *LayerInstanceType) MoveLayerByAbsoluteValue(xLocation int, yLocation int) {
	validateLayer(shared.layerAlias)
	layerEntry := Layers.Get(shared.layerAlias)
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
func (shared *LayerInstanceType) MoveLayerByRelativeValue(xLocation int, yLocation int) {
	validateLayer(shared.layerAlias)
	layerEntry := Layers.Get(shared.layerAlias)
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
func (shared *LayerInstanceType) DeleteLayer() {
	validateLayer(shared.layerAlias)
	layer.Delete(shared.layerAlias)
	if commonResource.layerInstance.layerAlias == shared.layerAlias {
		nextLayerAlias := layer.GetNextLayerAlias()
		var nextLayerInstance *types.LayerEntryType
		if nextLayerAlias != "" {
			nextLayerInstance = Layers.Get(nextLayerAlias)
			commonResource.layerInstance = &LayerInstanceType{layerAlias: nextLayerAlias, parentAlias: nextLayerInstance.ParentAlias, LayerWidth: nextLayerInstance.Width, LayerHeight: nextLayerInstance.Height}
		}
	}
	shared.layerAlias = ""
}

func (shared *LayerInstanceType) IsLayerExists() bool {
	if shared.layerAlias != "" {
		return true
	}
	return false
}

func (shared *LayerInstanceType) SetIsVisible(isVisible bool) {
	validateLayer(shared.layerAlias)
	setLayerIsVisible(shared.layerAlias, isVisible)
}

/*
SetTopmost sets the current layer to be the topmost layer within its parent.
This ensures that it is rendered on top of all other sibling layers.
*/
func (shared *LayerInstanceType) SetTopmost() {
	validateLayer(shared.layerAlias)
	layer.SetTopmostLayer(shared.layerAlias)
}

/*
GetLocation returns the current x and y location of the layer. This is useful
when you need to determine the exact position of a layer on the screen. The
function returns two integer values representing the x and y coordinates.
*/
func (shared *LayerInstanceType) GetLocation() (int, int) {
	validateLayer(shared.layerAlias)
	layerEntry := Layers.Get(shared.layerAlias)
	return layerEntry.ScreenXLocation, layerEntry.ScreenYLocation
}

/*
GetLayerSize returns the current width and height of the layer. This is useful
when you need to determine the exact dimensions of a layer. The function
returns two integer values representing the width and height.
*/
func (shared *LayerInstanceType) GetLayerSize() (int, int) {
	validateLayer(shared.layerAlias)
	layerEntry := Layers.Get(shared.layerAlias)
	return layerEntry.Width, layerEntry.Height
}

/*
SetAlpha allows you to set the alpha value for a given text layer. This lets
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
func (shared *LayerInstanceType) SetAlpha(alphaValue float32) {
	validateLayer(shared.layerAlias)
	layerEntry := Layers.Get(shared.layerAlias)
	layerEntry.DefaultAttribute.ForegroundAlphaValue = alphaValue
	layerEntry.DefaultAttribute.BackgroundAlphaValue = alphaValue
}

/*
GetAlpha allows you to retrieve the alpha value for a given text layer.
This value represents the layer's transparency level.
*/
func (shared *LayerInstanceType) GetAlpha() float32 {
	validateLayer(shared.layerAlias)
	layerEntry := Layers.Get(shared.layerAlias)
	return layerEntry.DefaultAttribute.ForegroundAlphaValue
}

/*
LoadLayer allows you to load a pre-rendered layer from disk and add it to the layer system.
This is useful for quickly loading complex layers that were previously saved, such as image layers.
The layer is loaded from a compressed format that was created by SaveLayer. In addition, the following
information should be noted:

- The file extension ".clayer" is automatically appended to the filename if not provided.
- If the file cannot be read or is not a valid layer file, an error is returned.
- The loaded layer is added to the layer system with the specified alias, position, and z-order.
- The function returns a LayerInstanceType that can be used to manipulate the loaded layer.
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
LoadPreRenderedLayerImage allows you to load a pre-rendered layer image directly into image memory.
This is different from loading an image and pre-rendering it afterwards, as it directly loads
a layer that has already been rendered. This is useful for quickly loading complex images
that have been pre-processed and saved as layers. In addition, the following information
should be noted:

- The file extension ".clayer" is automatically appended to the filename if not provided.
- If the file cannot be read or is not a valid layer file, an error is returned.
- The loaded layer is added to the image system with the specified alias.
*/
func (shared *LayerInstanceType) LoadPreRenderedLayerImage(filePath string, imageAlias string) error {
	// Load the pre-rendered layer image using the image.go function
	return LoadPreRenderedLayerImage(filePath, imageAlias)
}

/*
SaveLayer allows you to save the current layer to disk. This is useful for caching
complex layers that take time to render, such as image layers. The layer is saved in a
compressed format to minimize disk space usage. In addition, the following information
should be noted:

- The file extension ".clayer" is automatically appended to the filename if not provided.
- The layer is saved using gzip compression to minimize disk space.
- If the file cannot be written, an error is returned.
*/
func (shared *LayerInstanceType) SaveLayer(filePath string) error {
	validateLayer(shared.layerAlias)
	layerEntry := Layers.Get(shared.layerAlias)
	return layerEntry.SaveLayer(filePath)
}

func (shared *LayerInstanceType) ColorStyle(styleEntry types.TuiStyleEntryType) {
	shared.colorLayer(*shared, styleEntry.Text.ForegroundColor, styleEntry.Text.BackgroundColor)
}

/*
Color allows you to set default colors on your text layer for printing with.
The color index specified corresponds to the 16 color ANSI standard, where
color 0 is Black and 15 is Bright White.
*/
func (shared *LayerInstanceType) Color(foregroundColorIndex int, backgroundColorIndex int) {
	validateColorIndex(foregroundColorIndex)
	validateColorIndex(backgroundColorIndex)
	shared.colorLayer(*shared, constants.AnsiColorByIndex[foregroundColorIndex], constants.AnsiColorByIndex[backgroundColorIndex])
}

/*
colorLayer allows you to set default colors on your specified text layer for
printing with. The color index specified corresponds to the 16 color ANSI
standard, where color 0 is Black and 15 is Bright White.
*/
func (shared *LayerInstanceType) colorLayer(layerInstance LayerInstanceType, foregroundColor constants.ColorType, backgroundColor constants.ColorType) {
	layerEntry := Layers.Get(layerInstance.layerAlias)
	layerEntry.DefaultAttribute.ForegroundColor = foregroundColor
	layerEntry.DefaultAttribute.BackgroundColor = backgroundColor
}

/*
ColorRGB allows you to set default colors on your text layer for printing with.
This method allows you to specify colors using RGB color index values within
the range of 0 to 255.
*/
func (shared *LayerInstanceType) ColorRGB(foregroundRedIndex int32, foregroundGreenIndex int32, foregroundBlueIndex int32, backgroundRedIndex int32, backgroundGreenIndex int32, backgroundBlueIndex int32) {
	shared.colorLayer(*shared, GetRGBColor(foregroundRedIndex, foregroundGreenIndex, foregroundBlueIndex), GetRGBColor(backgroundRedIndex, backgroundGreenIndex, backgroundBlueIndex))
}

/*
Color24Bit allows you to color a layer using a 24-bit color expressed as
an int32. This is useful for when you have colors which are already defined.
*/

func (shared *LayerInstanceType) Color24Bit(foregroundColor constants.ColorType, backgroundColor constants.ColorType) {
	shared.colorLayer(*shared, foregroundColor, backgroundColor)
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
func (shared *LayerInstanceType) Locate(xLocation int, yLocation int) {
	validateDefaultLayerIsNotEmpty()
	LocateLayer(shared, xLocation, yLocation)
}

/*
printLayer allows you to write text to a text layer. This is useful
for internal methods that want to write text to a text layer directly, without
effecting user settings (such as current cursor location, etc). In addition,
the following information should be noted:

- If the location to print falls outside the range of the text layer,
then only the visible portion of your text will be printed.
*/
func (shared *LayerInstanceType) printLayer(layerEntry *types.LayerEntryType, attributeEntry types.AttributeEntryType, xLocation int, yLocation int, textToPrint []rune) int {
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
func (shared *LayerInstanceType) Print(textToPrint string) {
	validateDefaultLayerIsNotEmpty()
	formattedTextToPrint := fmt.Sprint(textToPrint)
	layerEntry := Layers.Get(shared.layerAlias)
	if layerEntry.CursorYLocation >= layerEntry.Height {
		layerEntry.CursorYLocation = layerEntry.Height - 1
		layerEntry.CharacterMemory = scrollCharacterMemory(layerEntry)
	}
	arrayOfRunes := stringformat.GetRunesFromString(formattedTextToPrint)
	shared.printLayer(layerEntry, layerEntry.DefaultAttribute, layerEntry.CursorXLocation, layerEntry.CursorYLocation, arrayOfRunes)
	layerEntry.CursorXLocation = 0
	layerEntry.CursorYLocation = layerEntry.CursorYLocation + 1
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
func (shared *LayerInstanceType) PrintDialog(xLocation int, yLocation int, widthOfLineInCharacters int, printDelayInMilliseconds int, isSkipable bool, stringToPrint string) {
	formattedTextToPrint := fmt.Sprint(stringToPrint)
	layerEntry := Layers.Get(shared.layerAlias)
	if xLocation < 0 || xLocation > layerEntry.Width || yLocation < 0 || yLocation > layerEntry.Height {
		panic(fmt.Sprintf("The specified location (%d, %d) is out of bounds for layer '%s' with a size of (%d, %d).", xLocation, yLocation, layerEntry.LayerAlias, layerEntry.Width, layerEntry.Height))
	}
	printDialog(layerEntry, layerEntry.DefaultAttribute, xLocation, yLocation, widthOfLineInCharacters, printDelayInMilliseconds, isSkipable, formattedTextToPrint)
}

/*
PrintMarkup allows you to write text immediately to the terminal screen with word wrapping
and attribute tags. This is similar to PrintDialog but without the typewriter
effect and printing delay. In addition, the following information should be noted:

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
of your line has already exceeded its allowed maximum, the word will be pushed
to the line directly under it. This prevents words from being split across
two lines.

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
	// Print some text on the text layer "ForegroundLayer", at location
	// (0, 0), with a text wrapping location at 30 characters.
	// Inside our string to print, we add the "{red}" tag to switch printing
	// styles on the fly to "red" and change back to the default style using
	// "{}".
	dosktop.PrintMarkup("ForegroundLayer", 0, 0, 30, "This is some text with {red}red color{}. Only the words 'red color' should be colored.")
*/
func (shared *LayerInstanceType) PrintMarkup(xLocation int, yLocation int, widthOfLineInCharacters int, stringToPrint string) {
	formattedTextToPrint := fmt.Sprint(stringToPrint)
	layerEntry := Layers.Get(shared.layerAlias)
	if xLocation < 0 || xLocation > layerEntry.Width || yLocation < 0 || yLocation > layerEntry.Height {
		panic(fmt.Sprintf("The specified location (%d, %d) is out of bounds for layer '%s' with a size of (%d, %d).", xLocation, yLocation, layerEntry.LayerAlias, layerEntry.Width, layerEntry.Height))
	}
	layer.printMarkup(layerEntry, layerEntry.DefaultAttribute, xLocation, yLocation, widthOfLineInCharacters, formattedTextToPrint)
}

// PrintFont renders a string onto a layer using the specified font.
func (shared *LayerInstanceType) PrintFont(fontInstance fontInstanceType, xLocation int, yLocation int, stringToPrint string) {
	layerEntry := Layers.Get(shared.layerAlias)
	if layerEntry == nil {
		panic(fmt.Sprintf("Layer with alias '%s' not found.", shared.layerAlias))
	}
	Font.PrintText(layerEntry, fontInstance, xLocation, yLocation, stringToPrint)
}

/*
PrintFontDialog allows you to write text to the terminal screen with a typewriter effect using a specified font.
This is useful for creating animated text sequences with custom fonts. In addition, the following information should be noted:

  - If you specify a print location outside the range of your specified text layer, a panic will be generated.
  - When specifying a printing delay, the amount of time to wait is inserted between each character printed.
  - If the dialog being printed is flagged as skippable, the user can speed up printing by pressing the 'enter' key or right mouse button.
  - Specifying the width of your text line in characters allows you to control when text wrapping occurs. For example, if printing starts at location (2, 2) and you set
    a line width of 10 characters, text wrapping will occur after 10 characters have been printed. When this happens, text will continue
    to print underneath the previous line at a distance of font height + 1.
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
func AddLayer(xLocation int, yLocation int, width int, height int, zOrderPriority int, parentLayerInstance *LayerInstanceType) *LayerInstanceType {
	layerAlias := getUUID()
	validateTerminalWidthAndHeight(width, height)
	if parentLayerInstance == nil {
		layer.Add(layerAlias, xLocation, yLocation, width, height, zOrderPriority, "")
		layerInstance := LayerInstanceType{layerAlias: layerAlias, parentAlias: "", LayerWidth: width, LayerHeight: height}
		commonResource.layerInstance = &layerInstance
		return &layerInstance
	} else {
		layer.Add(layerAlias, xLocation, yLocation, width, height, zOrderPriority, parentLayerInstance.layerAlias)
		layerInstance := LayerInstanceType{layerAlias: layerAlias, parentAlias: "", LayerWidth: width, LayerHeight: height}
		commonResource.layerInstance = &layerInstance
		return &layerInstance
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
	layerEntry := Layers.Get(layerAlias)
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
	layerEntry := Layers.Get(layerAlias)
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
	layer.Delete(layerAlias)
	if commonResource.layerInstance.layerAlias == layerAlias {
		nextLayerAlias := layer.GetNextLayerAlias()
		// If last entry and no more layers, just return. Do not set anything.
		if nextLayerAlias == "" {
			commonResource.layerInstance = &LayerInstanceType{layerAlias: "", parentAlias: "", LayerWidth: 0, LayerHeight: 0}
			return
		}
		nextLayerInstance := Layers.Get(nextLayerAlias)
		commonResource.layerInstance = &LayerInstanceType{layerAlias: nextLayerAlias, parentAlias: nextLayerInstance.ParentAlias, LayerWidth: nextLayerInstance.Width, LayerHeight: nextLayerInstance.Height}
	}
}

func DeleteLayer(layerInstance *LayerInstanceType) *LayerInstanceType {
	deleteLayer(layerInstance.layerAlias)
	if commonResource.layerInstance.layerAlias == layerInstance.layerAlias {
		nextLayerAlias := layer.GetNextLayerAlias()
		nextLayerInstance := Layers.Get(nextLayerAlias)
		commonResource.layerInstance = &LayerInstanceType{layerAlias: nextLayerAlias, parentAlias: nextLayerInstance.ParentAlias, LayerWidth: nextLayerInstance.Width, LayerHeight: nextLayerInstance.Height}
	}
	return nil
}

/*
DeleteAllLayers allows you to remove all layers from memory.
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
SetTopmostLayer sets the specified layer to be the topmost layer.
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
