package consolizer

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/memory"
	"github.com/supercom32/consolizer/types"
)

/*
Font is a structure which represents a global instance of fontType.

Example:

	Font
*/
var Font fontType
var fonts = memory.NewMemoryManager[types.FontFamilyType]()

// Constants
const (
	numberOfCharacters = 94
	magicHeader        = "\x13TheDraw FONTS file\x1a"
)

// The standard printable ASCII characters in TDF fonts
var characterList = []rune{
	'!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ':', ';', '<', '=', '>', '?',
	'@', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O',
	'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '[', '\\', ']', '^', '_',
	'`', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o',
	'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', '{', '|', '}', '~',
}

/*
fontInstanceType is a structure which represents a font instance that can be used for rendering.
*/
type fontInstanceType struct {
	fontAlias string
	fontIndex int
}

/*
Unload is a method which allows you to remove a font from memory. In addition, the following should be noted:

- This method is used to free up memory when a font is no longer needed.

- If the font alias does not exist, the application will panic.

Example:

	font.Unload()
*/
func (instance *fontInstanceType) Unload() {
	if !fonts.IsExists(instance.fontAlias) {
		safeSttyPanic(fmt.Sprintf("Could not unload font with alias '%s' because it was not loaded.", instance.fontAlias))
	}
	fonts.Remove(instance.fontAlias)
}

/*
SwitchFont is a method which allows you to switch to a different font in the same file.

Example:

	font.SwitchFont(1)
*/
func (instance *fontInstanceType) SwitchFont(fontIndex int) {
	instance.fontIndex = fontIndex
}

/*
SetCharacterSpacing is a method which allows you to set the character spacing for the font instance.

Example:

	font.SetCharacterSpacing(1)
*/
func (instance *fontInstanceType) SetCharacterSpacing(characterSpacing int) {
	fontFamily := getFontFamilyFromMemory(instance.fontAlias)
	if instance.fontIndex >= len(fontFamily.Fonts) {
		safeSttyPanic(fmt.Sprintf("Font index %d not found in font alias '%s'.", instance.fontIndex, instance.fontAlias))
	}
	fontFamily.Fonts[instance.fontIndex].CharacterSpacing = characterSpacing
}

/*
SetBlankSizeInCharacters is a method which allows you to set the blank size for the font instance.

Example:

	font.SetBlankSizeInCharacters(2)
*/
func (instance *fontInstanceType) SetBlankSizeInCharacters(blankSize int) {
	fontFamily := getFontFamilyFromMemory(instance.fontAlias)
	if instance.fontIndex >= len(fontFamily.Fonts) {
		safeSttyPanic(fmt.Sprintf("Font index %d not found in font alias '%s'.", instance.fontIndex, instance.fontAlias))
	}
	fontFamily.Fonts[instance.fontIndex].BlankSizeInCharacters = blankSize
}

/*
fontType is a structure which acts as a manager/factory for fonts.
*/
type fontType struct{}

/*
Load is a method which allows you to load a TDF font from the given file path. In addition, the following should be
noted:

- All fonts in the file are loaded and cached in memory.

- Use SwitchFont to select a different font from the same file.

Example:

	font := Font.Load("myfont.tdf")
*/
func (manager *fontType) Load(fontFile string) fontInstanceType {
	// If the font family is already loaded, just return a new instance.
	if !fonts.IsExists(fontFile) {
		fileData, err := getFileDataFromFileSystem(fontFile)
		if err != nil {
			safeSttyPanic(fmt.Sprintf("Could not load font file '%s': %s", fontFile, err.Error()))
		}

		if !strings.HasPrefix(string(fileData), magicHeader) {
			safeSttyPanic(fmt.Sprintf("The file '%s' is not a valid TDF font.", fontFile))
		}

		// Find all fonts in the file
		fontOffsets := manager.findFontOffsets(fileData)

		if len(fontOffsets) == 0 {
			safeSttyPanic(fmt.Sprintf("No fonts found in file '%s'.", fontFile))
		}

		fontFamily := &types.FontFamilyType{
			Fonts: make([]*types.FontEntryType, len(fontOffsets)),
		}

		// Load all fonts with unique aliases
		for fontIndex, fontOffset := range fontOffsets {
			// Load the font at this offset
			fontEntry := manager.loadFontAtOffset(fileData, fontOffset)
			fontFamily.Fonts[fontIndex] = fontEntry
		}
		fonts.Add(fontFile, fontFamily)
	}

	// Return instance for the first font (index 0)
	var fontInstance fontInstanceType
	fontInstance.fontAlias = fontFile
	fontInstance.fontIndex = 0
	return fontInstance
}

/*
findFontOffsets is a method which allows you to find all font offsets in the file by reading each font's block size to
determine the start of the next font.

Example:

	offsets := font.findFontOffsets(data)
*/
func (manager *fontType) findFontOffsets(fileData []byte) []int {
	var fontOffsets []int
	currentOffset := 20 // First font always starts at offset 20

	for currentOffset < len(fileData) {
		// Ensure we have enough data for the font header
		if currentOffset+25 > len(fileData) {
			break
		}

		// Add the current offset to our list of font offsets.
		fontOffsets = append(fontOffsets, currentOffset)

		// The block size is a 2-byte little-endian integer located 23 bytes
		// from the start of the font header.
		blockSizeOffset := currentOffset + 23
		blockSize := int(binary.LittleEndian.Uint16(fileData[blockSizeOffset : blockSizeOffset+2]))

		// The next font header begins immediately after the current font's
		// character data block. The character data block starts 213 bytes
		// after the font header.
		currentOffset += 213 + blockSize
	}

	return fontOffsets
}

/*
loadFontAtOffset is a method which allows you to load a font from the given offset.

Example:

	fontEntry := font.loadFontAtOffset(data, offset)
*/
func (manager *fontType) loadFontAtOffset(fileData []byte, fontOffset int) *types.FontEntryType {
	fontEntry := &types.FontEntryType{}
	fontEntry.CharacterSpacing = -1 // Default character spacing
	fontEntry.BlankSizeInCharacters = 1

	// Read font name length
	nameLength := int(fileData[fontOffset+4])

	// Read font name
	fontEntry.Name = string(fileData[fontOffset+5 : fontOffset+5+nameLength])

	// Read font type and spacing
	fontEntry.FontType = fileData[fontOffset+21]
	fontEntry.Spacing = fileData[fontOffset+22]

	// Calculate block size
	blockSize := int(binary.LittleEndian.Uint16(fileData[fontOffset+23 : fontOffset+25]))

	// Initialize glyph and character list arrays
	fontEntry.Glyphs = make([]*types.Glyph, numberOfCharacters)
	fontEntry.CharList = make([]uint16, numberOfCharacters)

	// Character list offset starts at offset+25
	characterListOffset := fontOffset + 25
	for charIndex := 0; charIndex < numberOfCharacters; charIndex++ {
		fontEntry.CharList[charIndex] = binary.LittleEndian.Uint16(fileData[characterListOffset+charIndex*2 : characterListOffset+charIndex*2+2])
	}

	// Font data starts at offset+213
	fontDataOffset := fontOffset + 213
	fontEntry.FontData = fileData[fontDataOffset : fontDataOffset+blockSize]

	// Calculate font height.
	for charIndex := 0; charIndex < numberOfCharacters; charIndex++ {
		charOffset := fontEntry.CharList[charIndex]
		if charOffset != 0xffff {
			heightOffset := charOffset + 1
			if int(heightOffset) < len(fontEntry.FontData) {
				height := int(fontEntry.FontData[heightOffset])
				if height > fontEntry.Height {
					fontEntry.Height = height
				}
			}
		}
	}

	// Load glyphs
	for glyphIndex := 0; glyphIndex < numberOfCharacters; glyphIndex++ {
		charIndex := manager.lookupChar(characterList[glyphIndex])
		if charIndex != -1 {
			glyph, err := manager.readGlyph(fontEntry, glyphIndex)
			if err != nil {
				// Skip this glyph if there's an error
				fontEntry.Glyphs[charIndex] = nil
			} else {
				fontEntry.Glyphs[charIndex] = glyph
			}
		} else {
			fontEntry.Glyphs[glyphIndex] = nil
		}
	}

	return fontEntry
}

/*
GetAvailableFonts is a method which allows you to retrieve a list of all font names in the specified file.

Example:

	fonts := Font.GetAvailableFonts("myfont.tdf")
*/
func (manager *fontType) GetAvailableFonts(fontFile string) []string {
	fileData, err := getFileDataFromFileSystem(fontFile)
	if err != nil {
		safeSttyPanic(fmt.Sprintf("Could not load font file '%s': %s", fontFile, err.Error()))
	}

	if !strings.HasPrefix(string(fileData), magicHeader) {
		safeSttyPanic(fmt.Sprintf("The file '%s' is not a valid TDF font.", fontFile))
	}

	// Find all fonts in the file
	fontOffsets := manager.findFontOffsets(fileData)

	// Get the names of all fonts
	fontNames := make([]string, len(fontOffsets))
	for fontIndex, fontOffset := range fontOffsets {
		nameLength := int(fileData[fontOffset+4])
		fontNames[fontIndex] = string(fileData[fontOffset+5 : fontOffset+5+nameLength])
	}

	return fontNames
}

/*
readGlyph is a method which allows you to read a single glyph from the font data.

Example:

	glyph, err := font.readGlyph(entry, 0)
*/
func (manager *fontType) readGlyph(fontEntry *types.FontEntryType, glyphIndex int) (*types.Glyph, error) {
	if fontEntry.CharList[glyphIndex] == 0xffff {
		return nil, nil
	}

	dataOffset := int(fontEntry.CharList[glyphIndex])
	if dataOffset+2 > len(fontEntry.FontData) {
		return nil, fmt.Errorf("offset beyond file")
	}

	glyphWidth := int(fontEntry.FontData[dataOffset])
	glyphHeight := int(fontEntry.FontData[dataOffset+1])
	if glyphWidth <= 0 || glyphHeight <= 0 {
		return nil, nil
	}
	dataOffset += 2

	glyph := &types.Glyph{
		Width:  glyphWidth,
		Height: glyphHeight,
		Cells:  make([]types.Cell, glyphWidth*fontEntry.Height),
	}

	for cellIndex := range glyph.Cells {
		glyph.Cells[cellIndex] = types.Cell{Char: ' ', Color: 0}
	}

	rowIndex, columnIndex := 0, 0
	for dataOffset < len(fontEntry.FontData) {
		byteValue := fontEntry.FontData[dataOffset]
		dataOffset++
		if byteValue == 0 {
			break
		}

		if byteValue == 13 {
			rowIndex++
			columnIndex = 0
			continue
		}

		if dataOffset >= len(fontEntry.FontData) {
			break
		}
		colorValue := fontEntry.FontData[dataOffset]
		dataOffset++

		if rowIndex < fontEntry.Height && columnIndex < glyphWidth {
			cellIndex := rowIndex*glyphWidth + columnIndex
			var character rune
			if int(byteValue) < len(constants.CP437ToUnicode) {
				character = constants.CP437ToUnicode[byteValue]
			} else {
				character = ' '
			}
			glyph.Cells[cellIndex] = types.Cell{Char: character, Color: colorValue}
		}
		columnIndex++
	}

	return glyph, nil
}

/*
lookupChar is a method which allows you to find the index of a character in the character list.

Example:

	index := font.lookupChar('A')
*/
func (manager *fontType) lookupChar(characterToFind rune) int {
	for index, character := range characterList {
		if character == characterToFind {
			return index
		}
	}
	return -1
}

/*
renderGlyph is a method which allows you to draw a single character and returns the width it occupied.

Example:

	width := font.renderGlyph(layer, fontEntry, 'A', 0, 0)
*/
func (manager *fontType) renderGlyph(layerEntry *types.LayerEntryType, font *types.FontEntryType, character rune, xLocation, yLocation int) int {
	if character == ' ' {
		characterAIndex := manager.lookupChar('a')
		characterAWidth := 1
		if characterAIndex != -1 && font.Glyphs[characterAIndex] != nil {
			characterAWidth = font.Glyphs[characterAIndex].Width
		}
		return font.BlankSizeInCharacters * characterAWidth
	}
	characterIndex := manager.lookupChar(character)
	if characterIndex == -1 {
		return 1 // Return a default width for unknown characters.
	}
	glyph := font.Glyphs[characterIndex]
	if glyph == nil {
		return 1 // Return a default width for nil glyphs.
	}

	for rowIndex := 0; rowIndex < glyph.Height; rowIndex++ {
		for columnIndex := 0; columnIndex < glyph.Width; columnIndex++ {
			cell := glyph.Cells[rowIndex*glyph.Width+columnIndex]
			if cell.Char != ' ' { // Don't draw spaces
				attribute := types.NewAttributeEntry(&layerEntry.DefaultAttribute)
				if cell.Color != 0 {
					foregroundColor, backgroundColor := manager.convertTdfColorToRgb(cell.Color)
					attribute.ForegroundColor = foregroundColor
					attribute.BackgroundColor = backgroundColor
				}
				printLayer(layerEntry, attribute, xLocation+columnIndex, yLocation+rowIndex, []rune{cell.Char})
			}
		}
	}
	spacing := int(font.Spacing)
	if font.CharacterSpacing != -1 {
		spacing = font.CharacterSpacing
	}
	return glyph.Width + spacing
}

/*
PrintText is a method which allows you to render a string onto a layer using the specified font.

Example:

	Font.PrintText(layer, fontInstance, 0, 0, "Hello")
*/
func (manager *fontType) PrintText(layerEntry *types.LayerEntryType, fontInstance fontInstanceType, xLocation, yLocation int, textToPrint string) {
	fontFamily := getFontFamilyFromMemory(fontInstance.fontAlias)
	if fontInstance.fontIndex >= len(fontFamily.Fonts) {
		safeSttyPanic(fmt.Sprintf("Font index %d not found in font alias '%s'.", fontInstance.fontIndex, fontInstance.fontAlias))
	}
	font := fontFamily.Fonts[fontInstance.fontIndex]
	currentXLocation := xLocation
	for _, character := range textToPrint {
		characterWidth := manager.renderGlyph(layerEntry, font, character, currentXLocation, yLocation)
		currentXLocation += characterWidth
	}
}

/*
PrintTextDialog is a method which allows you to render a string onto a layer with a typewriter effect using the
specified font. In addition, the following should be noted:

- If widthOfLineInCharacters is greater than 0, text will wrap after that many characters.

Example:

	Font.PrintTextDialog(layer, fontInstance, 0, 0, 40, 100, true, "Typewriter text")
*/
func (manager *fontType) PrintTextDialog(layerEntry *types.LayerEntryType, fontInstance fontInstanceType, xLocation, yLocation, widthOfLineInCharacters, printDelayInMilliseconds int, isSkipable bool, textToPrint string) {
	fontFamily := getFontFamilyFromMemory(fontInstance.fontAlias)
	if fontInstance.fontIndex >= len(fontFamily.Fonts) {
		safeSttyPanic(fmt.Sprintf("Font index %d not found in font alias '%s'.", fontInstance.fontIndex, fontInstance.fontAlias))
	}
	font := fontFamily.Fonts[fontInstance.fontIndex]

	if printDelayInMilliseconds <= 0 {
		if widthOfLineInCharacters <= 0 {
			// Inline the PrintText functionality
			currentXLocation := xLocation
			for _, character := range textToPrint {
				characterWidth := manager.renderGlyph(layerEntry, font, character, currentXLocation, yLocation)
				currentXLocation += characterWidth
			}
		} else {
			// Implement non-typewriter font printing with line wrapping
			currentXLocation := xLocation
			currentYLocation := yLocation
			characterCount := 0

			for _, character := range textToPrint {
				// Check if we need to wrap to the next line
				if characterCount >= widthOfLineInCharacters {
					characterCount = 0
					currentXLocation = xLocation
					currentYLocation += font.Height + 1 // Move down by font height + 1
				}

				characterWidth := manager.renderGlyph(layerEntry, font, character, currentXLocation, currentYLocation)
				currentXLocation += characterWidth
				characterCount++
			}
		}
		return
	}

	isPrintDelaySkipped := false
	currentXLocation := xLocation
	currentYLocation := yLocation
	characterCount := 0

	for _, character := range textToPrint {
		// Check if we need to wrap to the next line
		if widthOfLineInCharacters > 0 && characterCount >= widthOfLineInCharacters {
			characterCount = 0
			currentXLocation = xLocation
			currentYLocation += font.Height + 1 // Move down by font height + 1
		}

		characterWidth := manager.renderGlyph(layerEntry, font, character, currentXLocation, currentYLocation)
		currentXLocation += characterWidth
		characterCount++

		// Check for skip input
		if isSkipable {
			_, _, mouseButtonPressed, _ := GetMouseStatus()
			keyPressed := Inkey()
			if mouseButtonPressed != 0 || string(keyPressed) == "enter" {
				isPrintDelaySkipped = true
			}
		}

		// Apply delay unless skipped
		if !isPrintDelaySkipped {
			SleepInMilliseconds(uint(printDelayInMilliseconds))
			UpdateDisplay(false)
		}
	}
	UpdateDisplay(false)
}

/*
convertTdfColorToRgb is a method which allows you to convert a TDF color byte to RGBA ColorType values.

Example:

	fg, bg := font.convertTdfColorToRgb(0x1F)
*/
func (manager *fontType) convertTdfColorToRgb(tdfColor byte) (constants.ColorType, constants.ColorType) {
	foregroundIndex := int(tdfColor & 0x0F)
	backgroundIndex := int((tdfColor & 0xF0) >> 4)

	if foregroundIndex >= len(constants.TdfToRgbMap) {
		foregroundIndex = 0 // Default to black if out of bounds
	}
	if backgroundIndex >= len(constants.TdfToRgbMap) {
		backgroundIndex = 0 // Default to black if out of bounds
	}

	return constants.TdfToRgbMap[foregroundIndex], constants.TdfToRgbMap[backgroundIndex]
}

/*
LoadFont is a method which allows you to load a TDF font from the given file path. In addition, the following should be
noted:

- This is kept for backward compatibility.

Example:

	font := LoadFont("myfont.tdf")
*/
func LoadFont(fontFile string) fontInstanceType {
	return Font.Load(fontFile)
}

/*
getFontFamilyFromMemory is a method which allows you to retrieve a font family from memory by its alias.

Example:

	family := getFontFamilyFromMemory("myfont.tdf")
*/
func getFontFamilyFromMemory(fontAlias string) *types.FontFamilyType {
	fontFamily := fonts.Get(fontAlias)
	if fontFamily == nil {
		safeSttyPanic(fmt.Sprintf("font with alias '%s' not found", fontAlias))
	}
	return fontFamily
}
