package consolizer

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/memory"
	"github.com/supercom32/consolizer/types"
)

// Font is a global instance of fontType.
var Font fontType
var fonts = memory.NewMemoryManager[types.FontFamilyType]()

// Constants
const (
	numChars    = 94
	magicHeader = "\x13TheDraw FONTS file\x1a"
)

// The standard printable ASCII characters in TDF fonts
var charlist = []rune{
	'!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ':', ';', '<', '=', '>', '?',
	'@', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O',
	'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '[', '\\', ']', '^', '_',
	'`', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o',
	'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', '{', '|', '}', '~',
}

type fontInstanceType struct {
	fontAlias string
	fontIndex int
}

// Unload removes a font from memory.
func (shared *fontInstanceType) Unload() {
	if !fonts.IsExists(shared.fontAlias) {
		safeSttyPanic(fmt.Sprintf("Could not unload font with alias '%s' because it was not loaded.", shared.fontAlias))
	}
	fonts.Remove(shared.fontAlias)
}

// SwitchFont switches to a different font in the same file.
func (shared *fontInstanceType) SwitchFont(fontIndex int) {
	shared.fontIndex = fontIndex
}

// SetCharacterSpacing sets the character spacing for the font instance.
func (shared *fontInstanceType) SetCharacterSpacing(spacing int) {
	fontFamily := getFontFamilyFromMemory(shared.fontAlias)
	if shared.fontIndex >= len(fontFamily.Fonts) {
		safeSttyPanic(fmt.Sprintf("Font index %d not found in font alias '%s'.", shared.fontIndex, shared.fontAlias))
	}
	fontFamily.Fonts[shared.fontIndex].CharacterSpacing = spacing
}

// SetBlankSizeInCharacters sets the blank size for the font instance.
func (shared *fontInstanceType) SetBlankSizeInCharacters(size int) {
	fontFamily := getFontFamilyFromMemory(shared.fontAlias)
	if shared.fontIndex >= len(fontFamily.Fonts) {
		safeSttyPanic(fmt.Sprintf("Font index %d not found in font alias '%s'.", shared.fontIndex, shared.fontAlias))
	}
	fontFamily.Fonts[shared.fontIndex].BlankSizeInCharacters = size
}

// fontType struct acts as a manager/factory for fonts
type fontType struct{}

// Load loads a TDF font from the given file path.
// All fonts in the file are loaded and cached in memory.
// Returns a fontInstanceType with fontIndex set to 0 (the first font).
// Use SwitchFont to select a different font from the same file.
func (shared *fontType) Load(fontFile string) fontInstanceType {
	// If the font family is already loaded, just return a new instance.
	if !fonts.IsExists(fontFile) {
		data, err := getFileDataFromFileSystem(fontFile)
		if err != nil {
			safeSttyPanic(fmt.Sprintf("Could not load font file '%s': %s", fontFile, err.Error()))
		}

		if !strings.HasPrefix(string(data), magicHeader) {
			safeSttyPanic(fmt.Sprintf("The file '%s' is not a valid TDF font.", fontFile))
		}

		// Find all fonts in the file
		fontOffsets := shared.findFontOffsets(data)

		if len(fontOffsets) == 0 {
			safeSttyPanic(fmt.Sprintf("No fonts found in file '%s'.", fontFile))
		}

		fontFamily := &types.FontFamilyType{
			Fonts: make([]*types.FontEntryType, len(fontOffsets)),
		}

		// Load all fonts with unique aliases
		for i, offset := range fontOffsets {
			// Load the font at this offset
			fontEntry := shared.loadFontAtOffset(data, offset)
			fontFamily.Fonts[i] = fontEntry
		}
		fonts.Add(fontFile, fontFamily)
	}

	// Return instance for the first font (index 0)
	var fontInstance fontInstanceType
	fontInstance.fontAlias = fontFile
	fontInstance.fontIndex = 0
	return fontInstance
}

// findFontOffsets finds all font offsets in the file by reading each font's
// block size to determine the start of the next font.
func (shared *fontType) findFontOffsets(data []byte) []int {
	var offsets []int
	currentOffset := 20 // First font always starts at offset 20

	for currentOffset < len(data) {
		// Ensure we have enough data for the font header
		if currentOffset+25 > len(data) {
			break
		}

		// Add the current offset to our list of font offsets.
		offsets = append(offsets, currentOffset)

		// The block size is a 2-byte little-endian integer located 23 bytes
		// from the start of the font header.
		blockSizeOffset := currentOffset + 23
		blockSize := int(binary.LittleEndian.Uint16(data[blockSizeOffset : blockSizeOffset+2]))

		// The next font header begins immediately after the current font's
		// character data block. The character data block starts 213 bytes
		// after the font header.
		currentOffset += 213 + blockSize
	}

	return offsets
}

// loadFontAtOffset loads a font from the given offset
func (shared *fontType) loadFontAtOffset(data []byte, offset int) *types.FontEntryType {
	fontEntry := &types.FontEntryType{}
	fontEntry.CharacterSpacing = -1 // Default character spacing
	fontEntry.BlankSizeInCharacters = 1

	// Read font name length
	nameLen := int(data[offset+4])

	// Read font name
	fontEntry.Name = string(data[offset+5 : offset+5+nameLen])

	// Read font type and spacing
	fontEntry.FontType = data[offset+21]
	fontEntry.Spacing = data[offset+22]

	// Calculate block size
	blockSize := int(binary.LittleEndian.Uint16(data[offset+23 : offset+25]))

	// Initialize glyph and character list arrays
	fontEntry.Glyphs = make([]*types.Glyph, numChars)
	fontEntry.CharList = make([]uint16, numChars)

	// Character list offset starts at offset+25
	charListOffset := offset + 25
	for i := 0; i < numChars; i++ {
		fontEntry.CharList[i] = binary.LittleEndian.Uint16(data[charListOffset+i*2 : charListOffset+i*2+2])
	}

	// Font data starts at offset+213
	fontDataOffset := offset + 213
	fontEntry.FontData = data[fontDataOffset : fontDataOffset+blockSize]

	// Calculate font height.
	for i := 0; i < numChars; i++ {
		charOffset := fontEntry.CharList[i]
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
	for i := 0; i < numChars; i++ {
		charIndex := shared.lookupChar(charlist[i])
		if charIndex != -1 {
			glyph, err := shared.readGlyph(fontEntry, i)
			if err != nil {
				// Skip this glyph if there's an error
				fontEntry.Glyphs[charIndex] = nil
			} else {
				fontEntry.Glyphs[charIndex] = glyph
			}
		} else {
			fontEntry.Glyphs[i] = nil
		}
	}

	return fontEntry
}

// GetAvailableFonts returns a list of all fonts in the specified file
func (shared *fontType) GetAvailableFonts(fontFile string) []string {
	data, err := getFileDataFromFileSystem(fontFile)
	if err != nil {
		safeSttyPanic(fmt.Sprintf("Could not load font file '%s': %s", fontFile, err.Error()))
	}

	if !strings.HasPrefix(string(data), magicHeader) {
		safeSttyPanic(fmt.Sprintf("The file '%s' is not a valid TDF font.", fontFile))
	}

	// Find all fonts in the file
	fontOffsets := shared.findFontOffsets(data)

	// Get the names of all fonts
	fontNames := make([]string, len(fontOffsets))
	for i, offset := range fontOffsets {
		nameLen := int(data[offset+4])
		fontNames[i] = string(data[offset+5 : offset+5+nameLen])
	}

	return fontNames
}

// readGlyph reads a single glyph from the font data.
func (shared *fontType) readGlyph(fontEntry *types.FontEntryType, index int) (*types.Glyph, error) {
	if fontEntry.CharList[index] == 0xffff {
		return nil, nil
	}

	offset := int(fontEntry.CharList[index])
	if offset+2 > len(fontEntry.FontData) {
		return nil, fmt.Errorf("offset beyond file")
	}

	width := int(fontEntry.FontData[offset])
	height := int(fontEntry.FontData[offset+1])
	if width <= 0 || height <= 0 {
		return nil, nil
	}
	offset += 2

	glyph := &types.Glyph{
		Width:  width,
		Height: height,
		Cells:  make([]types.Cell, width*fontEntry.Height),
	}

	for i := range glyph.Cells {
		glyph.Cells[i] = types.Cell{Char: ' ', Color: 0}
	}

	row, col := 0, 0
	for offset < len(fontEntry.FontData) {
		b := fontEntry.FontData[offset]
		offset++
		if b == 0 {
			break
		}

		if b == 13 {
			row++
			col = 0
			continue
		}

		if offset >= len(fontEntry.FontData) {
			break
		}
		color := fontEntry.FontData[offset]
		offset++

		if row < fontEntry.Height && col < width {
			cellIdx := row*width + col
			var char rune
			if int(b) < len(constants.CP437ToUnicode) {
				char = constants.CP437ToUnicode[b]
			} else {
				char = ' '
			}
			glyph.Cells[cellIdx] = types.Cell{Char: char, Color: color}
		}
		col++
	}

	return glyph, nil
}

// lookupChar finds the index of a character in the charlist.
func (shared *fontType) lookupChar(ch rune) int {
	for i, c := range charlist {
		if c == ch {
			return i
		}
	}
	return -1
}

// renderGlyph draws a single character and returns the width it occupied.
func (shared *fontType) renderGlyph(layerEntry *types.LayerEntryType, font *types.FontEntryType, ch rune, x, y int) int {
	if ch == ' ' {
		aIndex := shared.lookupChar('a')
		aWidth := 1
		if aIndex != -1 && font.Glyphs[aIndex] != nil {
			aWidth = font.Glyphs[aIndex].Width
		}
		return font.BlankSizeInCharacters * aWidth
	}
	idx := shared.lookupChar(ch)
	if idx == -1 {
		return 1 // Return a default width for unknown characters.
	}
	glyph := font.Glyphs[idx]
	if glyph == nil {
		return 1 // Return a default width for nil glyphs.
	}

	for row := 0; row < glyph.Height; row++ {
		for col := 0; col < glyph.Width; col++ {
			cell := glyph.Cells[row*glyph.Width+col]
			if cell.Char != ' ' { // Don't draw spaces
				attribute := types.NewAttributeEntry(&layerEntry.DefaultAttribute)
				if cell.Color != 0 {
					fg, bg := shared.convertTdfColorToRgb(cell.Color)
					attribute.ForegroundColor = fg
					attribute.BackgroundColor = bg
				}
				printLayer(layerEntry, attribute, x+col, y+row, []rune{cell.Char})
			}
		}
	}
	spacing := int(font.Spacing)
	if font.CharacterSpacing != -1 {
		spacing = font.CharacterSpacing
	}
	return glyph.Width + spacing
}

// PrintText renders a string onto a layer using the specified font.
func (shared *fontType) PrintText(layerEntry *types.LayerEntryType, fontInstance fontInstanceType, x, y int, str string) {
	fontFamily := getFontFamilyFromMemory(fontInstance.fontAlias)
	if fontInstance.fontIndex >= len(fontFamily.Fonts) {
		safeSttyPanic(fmt.Sprintf("Font index %d not found in font alias '%s'.", fontInstance.fontIndex, fontInstance.fontAlias))
	}
	font := fontFamily.Fonts[fontInstance.fontIndex]
	currentX := x
	for _, ch := range str {
		width := shared.renderGlyph(layerEntry, font, ch, currentX, y)
		currentX += width
	}
}

// PrintTextDialog renders a string onto a layer with a typewriter effect using the specified font.
// If widthOfLineInCharacters is greater than 0, text will wrap after that many characters.
func (shared *fontType) PrintTextDialog(layerEntry *types.LayerEntryType, fontInstance fontInstanceType, x, y, widthOfLineInCharacters, printDelayInMilliseconds int, isSkipable bool, str string) {
	fontFamily := getFontFamilyFromMemory(fontInstance.fontAlias)
	if fontInstance.fontIndex >= len(fontFamily.Fonts) {
		safeSttyPanic(fmt.Sprintf("Font index %d not found in font alias '%s'.", fontInstance.fontIndex, fontInstance.fontAlias))
	}
	font := fontFamily.Fonts[fontInstance.fontIndex]

	if printDelayInMilliseconds <= 0 {
		if widthOfLineInCharacters <= 0 {
			// Inline the PrintText functionality
			currentX := x
			for _, ch := range str {
				width := shared.renderGlyph(layerEntry, font, ch, currentX, y)
				currentX += width
			}
		} else {
			// Implement non-typewriter font printing with line wrapping
			currentX := x
			currentY := y
			charCount := 0

			for _, ch := range str {
				// Check if we need to wrap to the next line
				if charCount >= widthOfLineInCharacters {
					charCount = 0
					currentX = x
					currentY += font.Height + 1 // Move down by font height + 1
				}

				width := shared.renderGlyph(layerEntry, font, ch, currentX, currentY)
				currentX += width
				charCount++
			}
		}
		return
	}

	isPrintDelaySkipped := false
	currentX := x
	currentY := y
	charCount := 0

	for _, ch := range str {
		// Check if we need to wrap to the next line
		if widthOfLineInCharacters > 0 && charCount >= widthOfLineInCharacters {
			charCount = 0
			currentX = x
			currentY += font.Height + 1 // Move down by font height + 1
		}

		width := shared.renderGlyph(layerEntry, font, ch, currentX, currentY)
		currentX += width
		charCount++

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

// convertTdfColorToRgb converts a TDF color byte to RGBA ColorType values.
func (shared *fontType) convertTdfColorToRgb(c byte) (constants.ColorType, constants.ColorType) {
	fgIndex := int(c & 0x0F)
	bgIndex := int((c & 0xF0) >> 4)

	if fgIndex >= len(constants.TdfToRgbMap) {
		fgIndex = 0 // Default to black if out of bounds
	}
	if bgIndex >= len(constants.TdfToRgbMap) {
		bgIndex = 0 // Default to black if out of bounds
	}

	return constants.TdfToRgbMap[fgIndex], constants.TdfToRgbMap[bgIndex]
}

// LoadFont loads a TDF font from the given file path (kept for backward compatibility).
// If fontIndex is provided, it selects that specific font from the file.
func LoadFont(fontFile string) fontInstanceType {
	return Font.Load(fontFile)
}

func getFontFamilyFromMemory(alias string) *types.FontFamilyType {
	fontFamily := fonts.Get(alias)
	if fontFamily == nil {
		safeSttyPanic(fmt.Sprintf("font with alias '%s' not found", alias))
	}
	return fontFamily
}
