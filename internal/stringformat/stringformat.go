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
	setUnicodeRangeWidth('\u2600', '\u26FF', true) // Misc Symbols
	setUnicodeGeometricShapeWidth()
}

func setUnicodeRangeWidth(startingIndex int, endingIndex int, isWide bool) {
	offset := startingIndex
	for currentIndex := 0; currentIndex < endingIndex - startingIndex; currentIndex++ {
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
	// If Asian font which is detected as wide, return true.
	properties := width.LookupRune(character)
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

func GetMaxCharactersThatFitInStringSize(arrayOfRunes []rune, lengthOfString int) []rune {
	numberOfCharactersUsed := 0
	for currentIndex, currentRune := range arrayOfRunes {
		if IsRuneCharacterWide(currentRune) {
			numberOfCharactersUsed = numberOfCharactersUsed + 2
		} else {
			numberOfCharactersUsed++
		}
		if numberOfCharactersUsed == lengthOfString {
			// If your last character printed fills your string exactly, include it to your string.
			formattedArray := arrayOfRunes[:currentIndex+1]
			return formattedArray
		} else if numberOfCharactersUsed > lengthOfString {
			formattedArray := arrayOfRunes[:currentIndex]
			// If you just printed a double width character that exceeds the printing limit, just add a blank space
			// padding, since you only have 1 space left.
			formattedArray = append(formattedArray, ' ')
			return formattedArray
		}
	}
	return arrayOfRunes
}

func logInfo(info string) {
	filesystem.WriteBytesToFile("/tmp/debug.log", []byte(info), 666)
}

func GetFormattedString(stringToFormat string, lengthOfString int, position int) string {
	arrayOfRunes := GetRunesFromString(stringToFormat)
	if len(stringToFormat) == 0 {
		return GetFilledString(lengthOfString, " ")
	}
	widthOfRunesWhenPrinted := GetWidthOfRunesWhenPrinted(arrayOfRunes)
	paddingSize := lengthOfString - widthOfRunesWhenPrinted
	if paddingSize <= 0 {
		paddingSize = 0
		return string(GetMaxCharactersThatFitInStringSize(GetRunesFromString(stringToFormat), lengthOfString))
	}
	stringPaddingInRunes := GetRunesFromString(GetFilledString(paddingSize, " "))
	formattedArrayOfRunes := []rune{}
	if position == constants.AlignmentRight {
		formattedArrayOfRunes = append(stringPaddingInRunes, arrayOfRunes...)
	} else if position == constants.AlignmentCenter {
		formattedArrayOfRunes = append(stringPaddingInRunes, arrayOfRunes...)
		formattedArrayOfRunes = append(formattedArrayOfRunes, stringPaddingInRunes...)
		if len(formattedArrayOfRunes) < lengthOfString {
			formattedArrayOfRunes = append(formattedArrayOfRunes, ' ')
		}
	} else if position == constants.AlignmentNoPadding {
		formattedArrayOfRunes = append(formattedArrayOfRunes, ' ')
		formattedArrayOfRunes = append(formattedArrayOfRunes, arrayOfRunes...)
		formattedArrayOfRunes = append(formattedArrayOfRunes, ' ')
	} else {
		formattedArrayOfRunes = append(arrayOfRunes, stringPaddingInRunes...)
	}
	return string(formattedArrayOfRunes)
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
	time := string(time.Now().Unix())
	return "zzzzzzz" + time + id.String()
}