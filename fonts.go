package consolizer

import (
	"encoding/binary"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/memory"
	"github.com/supercom32/consolizer/types"
)

var fontManager = memory.NewMemoryManager[types.FontEntryType]()

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

func UnloadFont(fontAlias string) {
	if !fontManager.IsExists(fontAlias) {
		safeSttyPanic(fmt.Sprintf("Could not unload font with alias '%s' because it was not loaded.", fontAlias))
	}
	fontManager.Remove(fontAlias)
}

func GetFontFromMemory(alias string) (*types.FontEntryType, error) {
	font := fontManager.Get(alias)
	if font == nil {
		return nil, fmt.Errorf("font with alias '%s' not found", alias)
	}
	return font, nil
}

// LoadFont loads a TDF font from the given file path.
func LoadFont(fontFile string) string {
	if filepath.Ext(fontFile) == "" {
		fontFile += ".tdf"
	}

	data, err := getFileDataFromFileSystem(fontFile)
	if err != nil {
		safeSttyPanic(fmt.Sprintf("Could not load font file '%s': %s", fontFile, err.Error()))
	}

	if !strings.HasPrefix(string(data), magicHeader) {
		safeSttyPanic(fmt.Sprintf("The file '%s' is not a valid TDF font.", fontFile))
	}

	fontEntry := &types.FontEntryType{}
	nameLen := int(data[24])
	fontEntry.Name = string(data[25 : 25+nameLen])
	fontEntry.FontType = data[41]
	fontEntry.Spacing = data[42]
	fontEntry.Glyphs = make([]*types.Glyph, numChars)
	fontEntry.CharList = make([]uint16, numChars)

	// Character list offset starts at 45
	charListOffset := 45
	for i := 0; i < numChars; i++ {
		fontEntry.CharList[i] = binary.LittleEndian.Uint16(data[charListOffset+i*2 : charListOffset+i*2+2])
	}

	// Font data starts at 233
	fontEntry.FontData = data[233:]

	// Calculate fontEntry height and load glyphs
	for i := 0; i < numChars; i++ {
		charIndex := lookupChar(charlist[i])
		if charIndex != -1 {
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
	}

	for i := 0; i < numChars; i++ {
		charIndex := lookupChar(charlist[i])
		if charIndex != -1 {
			glyph, err := readGlyph(fontEntry, charIndex)
			if err != nil {
				safeSttyPanic(fmt.Sprintf("Failed to read glyph %d from font '%s': %v", i, fontFile, err))
			}
			fontEntry.Glyphs[i] = glyph
		} else {
			fontEntry.Glyphs[i] = nil
		}
	}
	fontManager.Add(fontFile, fontEntry)
	return fontFile
}

// readGlyph reads a single glyph from the font data.
func readGlyph(fontEntry *types.FontEntryType, index int) (*types.Glyph, error) {
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
func lookupChar(ch rune) int {
	for i, c := range charlist {
		if c == ch {
			return i
		}
	}
	return -1
}

// renderGlyph draws a single character and returns the width it occupied.
func renderGlyph(layerEntry *types.LayerEntryType, font *types.FontEntryType, ch rune, x, y int) int {
	idx := lookupChar(ch)
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
					fg, bg := convertTdfColorToRgb(cell.Color)
					attribute.ForegroundColor = fg
					attribute.BackgroundColor = bg
				}
				printLayer(layerEntry, attribute, x+col, y+row, []rune{cell.Char})
			}
		}
	}
	return glyph.Width + int(font.Spacing)
}

// printFont renders a string onto a layer using the specified font.
func printFont(layerEntry *types.LayerEntryType, x, y int, str string, fontAlias string) {
	font, err := GetFontFromMemory(fontAlias)
	if err != nil {
		safeSttyPanic(fmt.Sprintf("Font with alias '%s' not found.", fontAlias))
	}
	currentX := x
	for _, ch := range str {
		width := renderGlyph(layerEntry, font, ch, currentX, y)
		currentX += width
	}
}

// printFontDialog renders a string onto a layer with a typewriter effect using the specified font.
func printFontDialog(layerEntry *types.LayerEntryType, x, y, printDelayInMilliseconds int, isSkipable bool, str string, fontAlias string) {
	font, err := GetFontFromMemory(fontAlias)
	if err != nil {
		safeSttyPanic(fmt.Sprintf("Font with alias '%s' not found.", fontAlias))
	}

	if printDelayInMilliseconds <= 0 {
		printFont(layerEntry, x, y, str, fontAlias)
		return
	}

	isPrintDelaySkipped := false
	currentX := x
	for _, ch := range str {
		width := renderGlyph(layerEntry, font, ch, currentX, y)
		currentX += width

		// Check for skip input
		if isSkipable {
			_, _, mouseButtonPressed, _ := GetMouseStatus()
			keyPressed := Inkey()
			if mouseButtonPressed != 0 || string(keyPressed) == "enter" {
				isPrintDelaySkipped = true
			}
		}

		// Apply delay unless skipped
		if !isPrintDelaySkipped && printDelayInMilliseconds > 0 {
			SleepInMilliseconds(uint(printDelayInMilliseconds))
			UpdateDisplay(false)
		}
	}
	UpdateDisplay(false)
}

// convertTdfColorToRgb converts a TDF color byte to RGBA ColorType values.
func convertTdfColorToRgb(c byte) (constants.ColorType, constants.ColorType) {
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
