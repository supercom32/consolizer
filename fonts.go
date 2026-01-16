package consolizer

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/types"
)

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

// tdfToRgbMap maps the 16-color TDF palette to 24-bit RGB ColorType values.
var tdfToRgbMap = []constants.ColorType{
	GetRGBColor(0, 0, 0),       // 0: Black
	GetRGBColor(0, 0, 170),     // 1: Blue
	GetRGBColor(0, 170, 0),     // 2: Green
	GetRGBColor(0, 170, 170),   // 3: Cyan
	GetRGBColor(170, 0, 0),     // 4: Red
	GetRGBColor(170, 0, 170),   // 5: Magenta
	GetRGBColor(170, 85, 0),    // 6: Brown
	GetRGBColor(170, 170, 170), // 7: Light Gray
	GetRGBColor(85, 85, 85),    // 8: Dark Gray
	GetRGBColor(85, 85, 255),   // 9: Light Blue
	GetRGBColor(85, 255, 85),   // 10: Light Green
	GetRGBColor(85, 255, 255),  // 11: Light Cyan
	GetRGBColor(255, 85, 85),   // 12: Light Red
	GetRGBColor(255, 85, 255),  // 13: Light Magenta
	GetRGBColor(255, 255, 85),  // 14: Yellow
	GetRGBColor(255, 255, 255), // 15: White
}

// CP437 → Unicode mapping
var cp437ToUnicode = func() []rune {
	arr := make([]rune, 256)
	for i := 0; i < 256; i++ {
		arr[i] = rune(i) // default fallback
	}

	// Complete mapping from CP437 to Unicode
	arr[0] = 0x0000
	arr[1] = 0x263A
	arr[2] = 0x263B
	arr[3] = 0x2665
	arr[4] = 0x2666
	arr[5] = 0x2663
	arr[6] = 0x2660
	arr[7] = 0x2022
	arr[8] = 0x25D8
	arr[9] = 0x25CB
	arr[10] = 0x25D9
	arr[11] = 0x2642
	arr[12] = 0x2640
	arr[13] = 0x266A
	arr[14] = 0x266B
	arr[15] = 0x263C
	arr[16] = 0x25BA
	arr[17] = 0x25C4
	arr[18] = 0x2195
	arr[19] = 0x203C
	arr[20] = 0x00B6
	arr[21] = 0x00A7
	arr[22] = 0x25AC
	arr[23] = 0x21A8
	arr[24] = 0x2191
	arr[25] = 0x2193
	arr[26] = 0x2192
	arr[27] = 0x2190
	arr[28] = 0x221F
	arr[29] = 0x2194
	arr[30] = 0x25B2
	arr[31] = 0x25BC
	arr[32] = 0x0020
	arr[33] = 0x0021
	arr[34] = 0x0022
	arr[35] = 0x0023
	arr[36] = 0x0024
	arr[37] = 0x0025
	arr[38] = 0x0026
	arr[39] = 0x0027
	arr[40] = 0x0028
	arr[41] = 0x0029
	arr[42] = 0x002A
	arr[43] = 0x002B
	arr[44] = 0x002C
	arr[45] = 0x002D
	arr[46] = 0x002E
	arr[47] = 0x002F
	arr[48] = 0x0030
	arr[49] = 0x0031
	arr[50] = 0x0032
	arr[51] = 0x0033
	arr[52] = 0x0034
	arr[53] = 0x0035
	arr[54] = 0x0036
	arr[55] = 0x0037
	arr[56] = 0x0038
	arr[57] = 0x0039
	arr[58] = 0x003A
	arr[59] = 0x003B
	arr[60] = 0x003C
	arr[61] = 0x003D
	arr[62] = 0x003E
	arr[63] = 0x003F
	arr[64] = 0x0040
	arr[65] = 0x0041
	arr[66] = 0x0042
	arr[67] = 0x0043
	arr[68] = 0x0044
	arr[69] = 0x0045
	arr[70] = 0x0046
	arr[71] = 0x0047
	arr[72] = 0x0048
	arr[73] = 0x0049
	arr[74] = 0x004A
	arr[75] = 0x004B
	arr[76] = 0x004C
	arr[77] = 0x004D
	arr[78] = 0x004E
	arr[79] = 0x004F
	arr[80] = 0x0050
	arr[81] = 0x0051
	arr[82] = 0x0052
	arr[83] = 0x0053
	arr[84] = 0x0054
	arr[85] = 0x0055
	arr[86] = 0x0056
	arr[87] = 0x0057
	arr[88] = 0x0058
	arr[89] = 0x0059
	arr[90] = 0x005A
	arr[91] = 0x005B
	arr[92] = 0x005C
	arr[93] = 0x005D
	arr[94] = 0x005E
	arr[95] = 0x005F
	arr[96] = 0x0060
	arr[97] = 0x0061
	arr[98] = 0x0062
	arr[99] = 0x0063
	arr[100] = 0x0064
	arr[101] = 0x0065
	arr[102] = 0x0066
	arr[103] = 0x0067
	arr[104] = 0x0068
	arr[105] = 0x0069
	arr[106] = 0x006A
	arr[107] = 0x006B
	arr[108] = 0x006C
	arr[109] = 0x006D
	arr[110] = 0x006E
	arr[111] = 0x006F
	arr[112] = 0x0070
	arr[113] = 0x0071
	arr[114] = 0x0072
	arr[115] = 0x0073
	arr[116] = 0x0074
	arr[117] = 0x0075
	arr[118] = 0x0076
	arr[119] = 0x0077
	arr[120] = 0x0078
	arr[121] = 0x0079
	arr[122] = 0x007A
	arr[123] = 0x007B
	arr[124] = 0x007C
	arr[125] = 0x007D
	arr[126] = 0x007E
	arr[127] = 0x2302
	arr[128] = 0x00C7
	arr[129] = 0x00FC
	arr[130] = 0x00E9
	arr[131] = 0x00E2
	arr[132] = 0x00E4
	arr[133] = 0x00E0
	arr[134] = 0x00E5
	arr[135] = 0x00E7
	arr[136] = 0x00EA
	arr[137] = 0x00EB
	arr[138] = 0x00E8
	arr[139] = 0x00EF
	arr[140] = 0x00EE
	arr[141] = 0x00EC
	arr[142] = 0x00C4
	arr[143] = 0x00C5
	arr[144] = 0x00C9
	arr[145] = 0x00E6
	arr[146] = 0x00C6
	arr[147] = 0x00F4
	arr[148] = 0x00F6
	arr[149] = 0x00F2
	arr[150] = 0x00FB
	arr[151] = 0x00F9
	arr[152] = 0x00FF
	arr[153] = 0x00D6
	arr[154] = 0x00DC
	arr[155] = 0x00A2
	arr[156] = 0x00A3
	arr[157] = 0x00A5
	arr[158] = 0x20A7
	arr[159] = 0x0192
	arr[160] = 0x00E1
	arr[161] = 0x00ED
	arr[162] = 0x00F3
	arr[163] = 0x00FA
	arr[164] = 0x00F1
	arr[165] = 0x00D1
	arr[166] = 0x00AA
	arr[167] = 0x00BA
	arr[168] = 0x00BF
	arr[169] = 0x2310
	arr[170] = 0x00AC
	arr[171] = 0x00BD
	arr[172] = 0x00BC
	arr[173] = 0x00A1
	arr[174] = 0x00AB
	arr[175] = 0x00BB
	arr[176] = 0x2591
	arr[177] = 0x2592
	arr[178] = 0x2593
	arr[179] = 0x2502
	arr[180] = 0x2524
	arr[181] = 0x2561
	arr[182] = 0x2562
	arr[183] = 0x2556
	arr[184] = 0x2555
	arr[185] = 0x2563
	arr[186] = 0x2551
	arr[187] = 0x2557
	arr[188] = 0x255D
	arr[189] = 0x255C
	arr[190] = 0x255B
	arr[191] = 0x2510
	arr[192] = 0x2514
	arr[193] = 0x2534
	arr[194] = 0x252C
	arr[195] = 0x251C
	arr[196] = 0x2500
	arr[197] = 0x253C
	arr[198] = 0x255E
	arr[199] = 0x255F
	arr[200] = 0x255A
	arr[201] = 0x2554
	arr[202] = 0x2569
	arr[203] = 0x2566
	arr[204] = 0x2560
	arr[205] = 0x2550
	arr[206] = 0x256C
	arr[207] = 0x2567
	arr[208] = 0x2568
	arr[209] = 0x2564
	arr[210] = 0x2565
	arr[211] = 0x2559
	arr[212] = 0x2558
	arr[213] = 0x2552
	arr[214] = 0x2553
	arr[215] = 0x256B
	arr[216] = 0x256A
	arr[217] = 0x2518
	arr[218] = 0x250C
	arr[219] = 0x2588
	arr[220] = 0x2584
	arr[221] = 0x258C
	arr[222] = 0x2590
	arr[223] = 0x2580
	arr[224] = 0x03B1
	arr[225] = 0x00DF
	arr[226] = 0x0393
	arr[227] = 0x03C0
	arr[228] = 0x03A3
	arr[229] = 0x03C3
	arr[230] = 0x00B5
	arr[231] = 0x03C4
	arr[232] = 0x03A6
	arr[233] = 0x0398
	arr[234] = 0x03A9
	arr[235] = 0x03B4
	arr[236] = 0x221E
	arr[237] = 0x03C6
	arr[238] = 0x03B5
	arr[239] = 0x2229
	arr[240] = 0x2261
	arr[241] = 0x00B1
	arr[242] = 0x2265
	arr[243] = 0x2264
	arr[244] = 0x2320
	arr[245] = 0x2321
	arr[246] = 0x00F7
	arr[247] = 0x2248
	arr[248] = 0x00B0
	arr[249] = 0x2219
	arr[250] = 0x00B7
	arr[251] = 0x221A
	arr[252] = 0x207F
	arr[253] = 0x00B2
	arr[254] = 0x25A0
	arr[255] = 0x00A0

	return arr
}()

// LoadFont loads a TDF font from the given file path.
func LoadFont(fontFile, fontDir string) (*types.FontEntryType, error) {
	if filepath.Ext(fontFile) == "" {
		fontFile += ".tdf"
	}

	fullpath := filepath.Join(fontDir, fontFile)
	data, err := os.ReadFile(fullpath)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(string(data), magicHeader) {
		return nil, fmt.Errorf("invalid TDF font: %s", fullpath)
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
				return nil, fmt.Errorf("failed reading glyph %d: %v", i, err)
			}
			fontEntry.Glyphs[i] = glyph
		} else {
			fontEntry.Glyphs[i] = nil
		}
	}

	return fontEntry, nil
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
			if int(b) < len(cp437ToUnicode) {
				char = cp437ToUnicode[b]
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

// RenderOnLayer renders a string onto a layer using the specified font.
func RenderOnLayer(layerEntry *types.LayerEntryType, x, y int, str string, font *types.FontEntryType) {
	currentX := x
	for _, ch := range str {
		idx := lookupChar(ch)
		if idx == -1 {
			currentX += 1 // Advance for unknown characters
			continue
		}
		glyph := font.Glyphs[idx]
		if glyph == nil {
			currentX += 1 // Advance for nil glyphs
			continue
		}

		for row := 0; row < glyph.Height; row++ {
			for col := 0; col < glyph.Width; col++ {
				cell := glyph.Cells[row*glyph.Width+col]
				if cell.Char != ' ' { // Don't draw spaces
					// Start with a copy of the layer's default attribute.
					attribute := types.NewAttributeEntry(&layerEntry.DefaultAttribute)
					// If the font cell has a non-zero color attribute, it overrides the default.
					if cell.Color != 0 {
						fg, bg := convertTdfColorToRgb(cell.Color)
						attribute.ForegroundColor = fg
						attribute.BackgroundColor = bg
					}
					printLayer(layerEntry, attribute, currentX+col, y+row, []rune{cell.Char})
				}
			}
		}
		currentX += glyph.Width + int(font.Spacing)
	}
}

// convertTdfColorToRgb converts a TDF color byte to RGBA ColorType values.
func convertTdfColorToRgb(c byte) (constants.ColorType, constants.ColorType) {
	fgIndex := int(c & 0x0F)
	bgIndex := int((c & 0xF0) >> 4)

	if fgIndex >= len(tdfToRgbMap) {
		fgIndex = 0 // Default to black if out of bounds
	}
	if bgIndex >= len(tdfToRgbMap) {
		bgIndex = 0 // Default to black if out of bounds
	}

	return tdfToRgbMap[fgIndex], tdfToRgbMap[bgIndex]
}
