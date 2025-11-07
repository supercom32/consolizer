package stringformat

import (
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/recast"
	"github.com/supercom32/filesystem"
	"golang.org/x/text/width"
	"time"
)

var isUnicodeWide map[int]bool

const maxLen = 4096
const nullRune = '\x00'
const leftAligned = 0

func InitializeUnicodeWidthMemory() {
	isUnicodeWide = make(map[int]bool)
	setUnicodeRangeWidth('\u2500', '\u257F', false) // Box Drawing
	setUnicodeRangeWidth('\u2580', '\u259F', false) // Block Elements
	setUnicodeRangeWidth('\u2600', '\u26FF', true)  // Misc Symbols
	setUnicodeGeometricShapeWidth()
}

func setUnicodeRangeWidth(startingIndex int, endingIndex int, isWide bool) {
	offset := startingIndex
	for currentIndex := 0; currentIndex <= endingIndex-startingIndex; currentIndex++ {
		isUnicodeWide[offset+currentIndex] = isWide
	}
}

func setUnicodeGeometricShapeWidth() {
	setUnicodeRangeWidth(9632, 9727, true) // Symbols

	// Small arrows
	isUnicodeWide['\u25B4'] = false
	isUnicodeWide['\u25B5'] = false
	isUnicodeWide['\u25B8'] = false
	isUnicodeWide['\u25B9'] = false
	isUnicodeWide['\u25BE'] = false
	isUnicodeWide['\u25BF'] = false
	isUnicodeWide['\u25C2'] = false
	isUnicodeWide['\u25C3'] = false

	// Triangles (all directions)
	isUnicodeWide['\u25B2'] = false // BLACK UP-POINTING TRIANGLE
	isUnicodeWide['\u25B3'] = false // WHITE UP-POINTING TRIANGLE
	isUnicodeWide['\u25B6'] = false // BLACK RIGHT-POINTING TRIANGLE
	isUnicodeWide['\u25B7'] = false // WHITE RIGHT-POINTING TRIANGLE
	isUnicodeWide['\u25BC'] = false // BLACK DOWN-POINTING TRIANGLE
	isUnicodeWide['\u25BD'] = false // WHITE DOWN-POINTING TRIANGLE
	isUnicodeWide['\u25C0'] = false // BLACK LEFT-POINTING TRIANGLE
	isUnicodeWide['\u25C1'] = false // WHITE LEFT-POINTING TRIANGLE

	isUnicodeWide['\u25C4'] = false // BLACK LEFT-POINTING TRIANGLE (bold variant)
	isUnicodeWide['\u25C5'] = false // WHITE LEFT-POINTING TRIANGLE (bold variant)
	isUnicodeWide['\u25BA'] = false // BLACK RIGHT-POINTING POINTER (bold variant)
	isUnicodeWide['\u25BB'] = false // WHITE RIGHT-POINTING POINTER (bold variant)

	// Circle pieces
	isUnicodeWide['\u25DC'] = false
	isUnicodeWide['\u25DD'] = false
	isUnicodeWide['\u25DE'] = false
	isUnicodeWide['\u25DF'] = false

	// Additional circles
	isUnicodeWide['\u25CF'] = false // BLACK CIRCLE
	isUnicodeWide['\u25CB'] = false // WHITE CIRCLE
	isUnicodeWide['\u25C9'] = false // FISHEYE
	isUnicodeWide['\u25CE'] = false // BULLSEYE
	isUnicodeWide['\u25EF'] = false // LARGE CIRCLE
	isUnicodeWide['\u25CD'] = false // CIRCLE WITH VERTICAL FILL
	isUnicodeWide['\u25CC'] = false // DOTTED CIRCLE
	isUnicodeWide['\u25D0'] = false // CIRCLE WITH LEFT HALF BLACK
	isUnicodeWide['\u25D1'] = false // CIRCLE WITH RIGHT HALF BLACK
	isUnicodeWide['\u25D2'] = false // CIRCLE WITH LOWER HALF BLACK
	isUnicodeWide['\u25D3'] = false // CIRCLE WITH UPPER HALF BLACK

	// Extra medium circles from Misc Symbols block
	isUnicodeWide['\u26AA'] = false // MEDIUM WHITE CIRCLE
	isUnicodeWide['\u26AB'] = false // MEDIUM BLACK CIRCLE
}

func IsRuneCharacterWide(character rune) bool {
	if isUnicodeWide == nil {
		InitializeUnicodeWidthMemory()
	}

	properties := width.LookupRune(character)
	// If Asian font which is detected as wide, return true.
	if properties.Kind() == width.EastAsianWide || properties.Kind() == width.EastAsianFullwidth {
		return true
	}
	// If not multi-byte, then return false.
	_, numberOfBytes := width.LookupString(string(character))
	if numberOfBytes == 1 {
		return false
	}
	// If a specific override value is found in table memory, return that value
	if isWide, exists := isUnicodeWide[int(character)]; exists {
		return isWide
	}
	// Otherwise, by default assume character is wide.
	return true
}

func GetWidthOfRunesWhenPrinted(arrayOfRunes []rune) int {
	widthOfString := 0
	for _, currentCharacter := range arrayOfRunes {
		if IsRuneCharacterWide(currentCharacter) {
			widthOfString = widthOfString + 2
		} else {
			widthOfString++
		}
	}
	return widthOfString
}

// GetWidthOfRunesWhenPrintedWithoutMarkup calculates the width of runes when printed,
// excluding any markup characters (enclosed in {{ and }}).
func GetWidthOfRunesWhenPrintedWithoutMarkup(arrayOfRunes []rune) int {
	// Convert runes to string for easier markup detection
	textString := string(arrayOfRunes)

	// Get text without markup
	textWithoutMarkup := GetTextWithoutMarkup(textString)

	// Calculate width of the text without markup
	return GetWidthOfRunesWhenPrinted([]rune(textWithoutMarkup))
}

// GetTextWithoutMarkup removes markup tags (enclosed in {{ and }}) from the given text.
func GetTextWithoutMarkup(textString string) string {
	var result []rune
	runes := []rune(textString) // convert string to runes
	for i := 0; i < len(runes); i++ {
		if i+1 < len(runes) && runes[i] == '{' && runes[i+1] == '{' {
			// look for closing tag
			foundClosing := false
			for j := i + 2; j < len(runes)-1; j++ {
				if runes[j] == '}' && runes[j+1] == '}' {
					i = j + 1 // skip the tag
					foundClosing = true
					break
				}
			}
			if !foundClosing {
				result = append(result, '{', '{')
				i++ // skip the first {
			}
		} else {
			result = append(result, runes[i])
		}
	}
	return string(result)
}

func GetRunesFromString(stringToConvert string) []rune {
	var runes []rune
	runes = []rune(stringToConvert)
	return runes
}

func GetIntAsString(number interface{}) string {
	numberAsFloatint64 := recast.GetNumberAsInt64(number)
	return fmt.Sprintf("%d", numberAsFloatint64)
}
func GetFloatAsString(number interface{}) string {
	numberAsFloat64 := recast.GetNumberAsFloat64(number)
	return fmt.Sprintf("%g", numberAsFloat64)
}

func GetSubString(input string, start int, length int) string {
	asRunes := []rune(input)
	if start >= len(asRunes) {
		return ""
	}
	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}
	return string(asRunes[start : start+length])
}

func GetStringAsBase64(inputString string) string {
	base64String := base64.StdEncoding.EncodeToString([]byte(inputString))
	return base64String
}

func GetStringFromBase64(inputString string) string {
	decodedString, err := base64.StdEncoding.DecodeString(inputString)
	if err != nil {
		panic(err)
	}
	return string(decodedString)
}

func GetNumberOfWideCharacters(arrayOfRunes []rune) int {
	numberOfWideCharacters := 0
	for _, currentRune := range arrayOfRunes {
		if IsRuneCharacterWide((currentRune)) {
			numberOfWideCharacters++
		}
	}
	return numberOfWideCharacters
}

// This returns an array of all the characters that will fit inside the length specified.
// It accounts for markup characters by not counting them towards the length.
func GetMaxCharactersThatFitInStringSize(arrayOfRunes []rune, maxLengthOfString int) []rune {
	numberOfCharactersUsed := 0
	formattedArray := []rune{}

	// Process the runes, handling markup tags
	i := 0
	for i < len(arrayOfRunes) {
		// Check for markup opening
		if i+1 < len(arrayOfRunes) && arrayOfRunes[i] == '{' && arrayOfRunes[i+1] == '{' {
			// Look for closing tag
			markupStart := i
			markupEnd := -1

			// Find the closing tag
			for j := i + 2; j < len(arrayOfRunes)-1; j++ {
				if arrayOfRunes[j] == '}' && arrayOfRunes[j+1] == '}' {
					markupEnd = j + 1 // Position of the last '}'
					break
				}
			}

			if markupEnd != -1 {
				// Found a complete markup tag, add it without counting towards length
				for k := markupStart; k <= markupEnd; k++ {
					formattedArray = append(formattedArray, arrayOfRunes[k])
				}
				i = markupEnd + 1
				continue
			}
		}

		// Regular character (or incomplete markup)
		currentRune := arrayOfRunes[i]
		if IsRuneCharacterWide(currentRune) {
			numberOfCharactersUsed = numberOfCharactersUsed + 2
			if numberOfCharactersUsed > maxLengthOfString {
				// If you added a wide character and it won't fit (needs two free spaces),
				// we just add a blank space to pad it out.
				formattedArray = append(formattedArray, ' ')
				return formattedArray
			}
		} else {
			numberOfCharactersUsed++
		}

		formattedArray = append(formattedArray, currentRune)
		if numberOfCharactersUsed == maxLengthOfString {
			return formattedArray
		}

		i++
	}

	return formattedArray
}

// GetMaxCharactersThatFitInStringSizeReverse returns the number of characters from the end of the array
// that will fit within the specified length, accounting for markup characters.
func GetMaxCharactersThatFitInStringSizeReverse(arrayOfRunes []rune, maxLengthOfString int) int {
	// Convert to string for easier markup handling
	textString := string(arrayOfRunes)

	// Remove markup tags
	textWithoutMarkup := GetTextWithoutMarkup(textString)
	runesWithoutMarkup := []rune(textWithoutMarkup)

	// Calculate how many characters from the end will fit
	numberOfCharactersUsed := 0
	charactersToInclude := 0

	for i := len(runesWithoutMarkup) - 1; i >= 0; i-- {
		currentRune := runesWithoutMarkup[i]

		if IsRuneCharacterWide(currentRune) {
			numberOfCharactersUsed += 2
		} else {
			numberOfCharactersUsed++
		}

		charactersToInclude++

		if numberOfCharactersUsed >= maxLengthOfString {
			break
		}
	}

	// Now we need to map this back to the original string with markup
	// This is a simplified approach - we'll just return the number of characters
	// that would fit if we were to process from the end
	return charactersToInclude
}

func GetRuneArrayCopy(sourceRuneArray []rune) []rune {
	copyOfRuneArray := make([]rune, len(sourceRuneArray))
	copy(copyOfRuneArray, sourceRuneArray)
	return copyOfRuneArray
}

func logInfo(info string) {
	filesystem.AppendLineToFile("/tmp/debug.log", info+"\n", 0)
}

func GetFormattedString(stringToFormat string, lengthOfString int, position int) string {
	arrayOfRunes := GetRunesFromString(stringToFormat)
	return string(GetFormattedRuneArray(arrayOfRunes, lengthOfString, position))
}
func GetFormattedRuneArray(arrayOfRunes []rune, desiredLengthOfArray int, textAlignment int) []rune {
	if len(arrayOfRunes) == 0 {
		return GetRunesFromString(GetFilledString(desiredLengthOfArray, " "))
	}
	// Use GetWidthOfRunesWhenPrintedWithoutMarkup to exclude markup characters from the width calculation
	widthOfRunesWhenPrinted := GetWidthOfRunesWhenPrintedWithoutMarkup(arrayOfRunes)
	paddingSize := desiredLengthOfArray - widthOfRunesWhenPrinted
	if paddingSize <= 0 {
		paddingSize = 0
		return GetMaxCharactersThatFitInStringSize(arrayOfRunes, desiredLengthOfArray)
	}

	// If you're viewing the end of a long string (so you need padding) and some characters are wide,
	// you need to add padding to compensate for the missing width.
	// paddingSize = paddingSize + GetNumberOfWideCharacters(arrayOfRunes)

	// stringPaddingInRunes := GetRunesFromString(GetFilledString(paddingSize, " "))
	fullStringPadding := GetFilledRuneArray(paddingSize, ' ')
	halfStringPadding := GetFilledRuneArray(paddingSize/2, ' ')

	formattedArrayOfRunes := []rune{}
	if textAlignment == constants.AlignmentRight {
		formattedArrayOfRunes = append(GetMaxCharactersThatFitInStringSize(arrayOfRunes, desiredLengthOfArray))
		formattedArrayOfRunes = append(fullStringPadding, formattedArrayOfRunes...)
	} else if textAlignment == constants.AlignmentCenter {
		formattedArrayOfRunes = append(halfStringPadding, arrayOfRunes...)
		formattedArrayOfRunes = append(formattedArrayOfRunes, halfStringPadding...)
		if len(formattedArrayOfRunes) < desiredLengthOfArray {
			formattedArrayOfRunes = append(formattedArrayOfRunes, ' ')
		}
	} else if textAlignment == constants.AlignmentNoPadding {
		formattedArrayOfRunes = append(formattedArrayOfRunes, ' ')
		formattedArrayOfRunes = append(formattedArrayOfRunes, arrayOfRunes...)
		formattedArrayOfRunes = append(formattedArrayOfRunes, ' ')
	} else {
		formattedArrayOfRunes = append(GetMaxCharactersThatFitInStringSize(arrayOfRunes, desiredLengthOfArray))
		formattedArrayOfRunes = append(formattedArrayOfRunes, fullStringPadding...)
	}
	return formattedArrayOfRunes
}
func GetFilledRuneArray(lengthOfString int, character rune) []rune {
	result := GetFilledString(lengthOfString, string(character))
	return GetRunesFromString((result))
}

func GetFilledString(lengthOfString int, character string) string {
	newString := ""
	for currentIndex := 0; currentIndex < lengthOfString; currentIndex++ {
		newString += character
	}
	return newString
}

func GetLastSortedUUID() string {
	id := uuid.New()
	time := fmt.Sprint(time.Now().Unix())
	return "zzzzzzz" + time + id.String()
}
