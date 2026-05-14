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

/*
InitializeUnicodeWidthMemory is a method which allows you to initialize the unicode width memory table with default
values and overrides.

Example:

	InitializeUnicodeWidthMemory()
*/
func InitializeUnicodeWidthMemory() {
	isUnicodeWide = make(map[int]bool)
	setUnicodeRangeWidth('\u2500', '\u257F', false) // Box Drawing
	setUnicodeRangeWidth('\u2580', '\u259F', false) // Block Elements
	setUnicodeRangeWidth('\u2600', '\u26FF', true)  // Misc Symbols
	setUnicodeGeometricShapeWidth()
}

/*
setUnicodeRangeWidth is a method which allows you to set the width property for a range of unicode characters.

:param startingIndex: The starting unicode index.
:param endingIndex: The ending unicode index.
:param isWide: Whether the characters in the range are wide.

Example:

	setUnicodeRangeWidth(0x2500, 0x257F, false)
*/
func setUnicodeRangeWidth(startingIndex int, endingIndex int, isWide bool) {
	offset := startingIndex
	for currentIndex := 0; currentIndex <= endingIndex-startingIndex; currentIndex++ {
		isUnicodeWide[offset+currentIndex] = isWide
	}
}

/*
setUnicodeGeometricShapeWidth is a method which allows you to set the width property for specific unicode geometric
shapes.

Example:

	setUnicodeGeometricShapeWidth()
*/
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

/*
IsRuneCharacterWide is a method which allows you to determine if a rune character is wide.

:param character: The rune character to check.

:return: True if the character is wide, false otherwise.

Example:

	isWide := IsRuneCharacterWide('读')
*/
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

/*
GetWidthOfRunesWhenPrinted is a method which allows you to calculate the total width of an array of runes when printed.

:param arrayOfRunes: The array of runes to calculate the width for.

:return: The total width of the runes.

Example:

	width := GetWidthOfRunesWhenPrinted(rune("test"))
*/
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

/*
GetWidthOfRunesWhenPrintedWithoutMarkup is a method which allows you to calculate the width of runes when printed,
excluding any markup characters.

:param arrayOfRunes: The array of runes to calculate the width for.

:return: The width of the runes excluding markup.

Example:

	width := GetWidthOfRunesWhenPrintedWithoutMarkup(rune("{{red}}test{{/}}"))
*/
func GetWidthOfRunesWhenPrintedWithoutMarkup(arrayOfRunes []rune) int {
	// Convert runes to string for easier markup detection
	textString := string(arrayOfRunes)

	// Get text without markup
	textWithoutMarkup := GetTextWithoutMarkup(textString)

	// Calculate width of the text without markup
	return GetWidthOfRunesWhenPrinted([]rune(textWithoutMarkup))
}

/*
GetTextWithoutMarkup is a method which allows you to remove markup tags (enclosed in {{ and }}) from a given string.

:param textString: The string containing markup tags.

:return: The string with markup tags removed.

Example:

	plainText := GetTextWithoutMarkup("{{red}}test{{/}}")
*/
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

/*
GetRunesFromString is a method which allows you to convert a string into an array of runes.

:param stringToConvert: The string to convert.

:return: An array of runes.

Example:

	runes := GetRunesFromString("test")
*/
func GetRunesFromString(stringToConvert string) []rune {
	var runes []rune
	runes = []rune(stringToConvert)
	return runes
}

/*
GetIntAsString is a method which allows you to convert a numeric variable to its string representation as an integer.

:param number: The number to convert.

:return: The string representation of the number as an integer.

Example:

	strValue := GetIntAsString(123.45)
*/
func GetIntAsString(number interface{}) string {
	numberAsFloatint64 := recast.GetNumberAsInt64(number)
	return fmt.Sprintf("%d", numberAsFloatint64)
}

/*
GetFloatAsString is a method which allows you to convert a numeric variable to its string representation as a float.

:param number: The number to convert.

:return: The string representation of the number as a float.

Example:

	strValue := GetFloatAsString(123.45)
*/
func GetFloatAsString(number interface{}) string {
	numberAsFloat64 := recast.GetNumberAsFloat64(number)
	return fmt.Sprintf("%g", numberAsFloat64)
}

/*
GetSubString is a method which allows you to get a substring from a given string.

:param input: The input string.
:param start: The starting index.
:param length: The length of the substring.

:return: The substring.

Example:

	sub := GetSubString("hello world", 0, 5)
*/
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

/*
GetStringAsBase64 is a method which allows you to encode a string to base64.

:param inputString: The string to encode.

:return: The base64 encoded string.

Example:

	b64 := GetStringAsBase64("test")
*/
func GetStringAsBase64(inputString string) string {
	base64String := base64.StdEncoding.EncodeToString([]byte(inputString))
	return base64String
}

/*
GetStringFromBase64 is a method which allows you to decode a base64 encoded string.

:param inputString: The base64 encoded string to decode.

:return: The decoded string.

Example:

	str := GetStringFromBase64("dGVzdA==")
*/
func GetStringFromBase64(inputString string) string {
	decodedString, err := base64.StdEncoding.DecodeString(inputString)
	if err != nil {
		panic(err)
	}
	return string(decodedString)
}

/*
GetNumberOfWideCharacters is a method which allows you to get the number of wide characters in an array of runes.

:param arrayOfRunes: The array of runes to check.

:return: The number of wide characters.

Example:

	count := GetNumberOfWideCharacters(rune("读test"))
*/
func GetNumberOfWideCharacters(arrayOfRunes []rune) int {
	numberOfWideCharacters := 0
	for _, currentRune := range arrayOfRunes {
		if IsRuneCharacterWide((currentRune)) {
			numberOfWideCharacters++
		}
	}
	return numberOfWideCharacters
}

/*
GetMaxCharactersThatFitInStringSize is a method which allows you to get an array of runes that will fit inside a
specified string size, accounting for markup.

:param arrayOfRunes: The array of runes to process.
:param maxLengthOfString: The maximum length of the resulting string.

:return: An array of runes that fits the specified size.

Example:

	runes := GetMaxCharactersThatFitInStringSize(rune("{{red}}test{{/}}"), 2)
*/
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

/*
GetMaxCharactersThatFitInStringSizeReverse is a method which allows you to calculate the number of characters from the
end of an array that fit within a specified length.

:param arrayOfRunes: The array of runes to process.
:param maxLengthOfString: The maximum length.

:return: The number of characters that fit from the end.

Example:

	count := GetMaxCharactersThatFitInStringSizeReverse(rune("test"), 2)
*/
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

/*
GetRuneArrayCopy is a method which allows you to create a copy of a rune array.

:param sourceRuneArray: The source rune array.

:return: A copy of the rune array.

Example:

	copy := GetRuneArrayCopy(original)
*/
func GetRuneArrayCopy(sourceRuneArray []rune) []rune {
	copyOfRuneArray := make([]rune, len(sourceRuneArray))
	copy(copyOfRuneArray, sourceRuneArray)
	return copyOfRuneArray
}

/*
logInfo is a method which allows you to log information to a debug file.

:param info: The information string to log.

Example:

	logInfo("debug message")
*/
func logInfo(info string) {
	filesystem.AppendLineToFile("/tmp/debug.log", info+"\n", 0)
}

/*
GetFormattedString is a method which allows you to get a formatted string based on specified length and alignment.

:param stringToFormat: The string to format.
:param lengthOfString: The desired length of the resulting string.
:param position: The alignment (e.g., left, center, right).

:return: The formatted string.

Example:

	fmtStr := GetFormattedString("test", 10, constants.AlignmentCenter)
*/
func GetFormattedString(stringToFormat string, lengthOfString int, position int) string {
	arrayOfRunes := GetRunesFromString(stringToFormat)
	return string(GetFormattedRuneArray(arrayOfRunes, lengthOfString, position))
}

/*
GetFormattedRuneArray is a method which allows you to get a formatted rune array based on desired length and alignment.

:param arrayOfRunes: The array of runes to format.
:param desiredLengthOfArray: The desired length of the resulting array.
:param textAlignment: The alignment (e.g., left, center, right).

:return: The formatted array of runes.

Example:

	fmtRunes := GetFormattedRuneArray(rune("test"), 10, constants.AlignmentLeft)
*/
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

/*
GetFilledRuneArray is a method which allows you to create a rune array of a specified length filled with a given
character.

:param lengthOfString: The length of the array.
:param character: The character to fill the array with.

:return: The filled array of runes.

Example:

	runes := GetFilledRuneArray(5, ' ')
*/
func GetFilledRuneArray(lengthOfString int, character rune) []rune {
	result := GetFilledString(lengthOfString, string(character))
	return GetRunesFromString((result))
}

/*
GetFilledString is a method which allows you to create a string of a specified length filled with a given character.

:param lengthOfString: The length of the string.
:param character: The character string to fill with.

:return: The filled string.

Example:

	str := GetFilledString(5, " ")
*/
func GetFilledString(lengthOfString int, character string) string {
	newString := ""
	for currentIndex := 0; currentIndex < lengthOfString; currentIndex++ {
		newString += character
	}
	return newString
}

/*
GetLastSortedUUID is a method which allows you to generate a UUID that is prefixed to ensure it sorts correctly based on
time.

:return: A time-sorted UUID string.

Example:

	uuid := GetLastSortedUUID()
*/
func GetLastSortedUUID() string {
	id := uuid.New()
	time := fmt.Sprint(time.Now().Unix())
	return "zzzzzzz" + time + id.String()
}
