package stringformat

import (
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/internal/recast"
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
	for currentIndex := 0; currentIndex < endingIndex-startingIndex; currentIndex++ {
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
	// Circle pieces
	isUnicodeWide['\u25DC'] = false
	isUnicodeWide['\u25DD'] = false
	isUnicodeWide['\u25DE'] = false
	isUnicodeWide['\u25DF'] = false
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

// This returns an array pf all the characters that will fit inside the length specified.
func GetMaxCharactersThatFitInStringSize(arrayOfRunes []rune, maxLengthOfString int) []rune {
	numberOfCharactersUsed := 0
	formattedArray := []rune{}
	for _, currentRune := range arrayOfRunes {
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
	}
	return formattedArray
}

func GetMaxCharactersThatFitInStringSizeReverse(arrayOfRunes []rune, maxLengthOfString int) int {
	numberOfCharactersUsed := 0
	formattedArray := []rune{}
	for currentRuneIndex := range arrayOfRunes {
		currentRune := arrayOfRunes[len(arrayOfRunes)-1-currentRuneIndex]
		if IsRuneCharacterWide(currentRune) {
			numberOfCharactersUsed = numberOfCharactersUsed + 2
			if numberOfCharactersUsed > maxLengthOfString {
				// If you added a wide character and it won't fit (needs two free spaces),
				// we just add a blank space to pad it out.
				formattedArray = append(formattedArray, ' ')
				return len(formattedArray)
			}
		} else {
			numberOfCharactersUsed++
		}
		formattedArray = append(formattedArray, currentRune)
		if numberOfCharactersUsed == maxLengthOfString {
			return len(formattedArray)
		}
	}
	return len(formattedArray)
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
	widthOfRunesWhenPrinted := GetWidthOfRunesWhenPrinted(arrayOfRunes)
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
