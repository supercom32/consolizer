package stringformat

import (
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/recast"
	"github.com/supercom32/filesystem"
	"time"
)

const maxLen = 4096

/*
IsRuneCharacterWide is a method which allows you to determine whether a rune occupies two terminal cells when
printed. It is a thin convenience wrapper defined entirely in terms of the width oracle, returning true exactly
when GetWidthOfRuneWhenPrinted reports a width of 2. It therefore never disagrees with what tcell draws and
carries no width rules of its own.

:param character: The rune whose printed width is being tested.
:return: True when the rune is wide (occupies two cells), false otherwise.

Example:
    isWide := IsRuneCharacterWide('读')
*/
func IsRuneCharacterWide(character rune) bool {
	return GetWidthOfRuneWhenPrinted(character) == 2
}

/*
GetWidthOfRunesWhenPrinted is a method which calculates the total width of an array of runes when printed.

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
GetWidthOfRunesWhenPrintedWithoutMarkup is a method which calculates the width of runes when printed, excluding any
markup characters.

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
GetTextWithoutMarkup is a method which removes markup tags (enclosed in {{ and }}) from a given string.

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
GetRunesFromString is a method which converts a string into an array of runes.

Example:
    runes := GetRunesFromString("test")
*/
func GetRunesFromString(stringToConvert string) []rune {
	var runes []rune
	runes = []rune(stringToConvert)
	return runes
}

/*
GetIntAsString is a method which converts a numeric variable to its string representation as an integer.

Example:
    strValue := GetIntAsString(123.45)
*/
func GetIntAsString(number interface{}) string {
	numberAsFloatint64 := recast.GetNumberAsInt64(number)
	return fmt.Sprintf("%d", numberAsFloatint64)
}

/*
GetFloatAsString is a method which converts a numeric variable to its string representation as a float.

Example:
    strValue := GetFloatAsString(123.45)
*/
func GetFloatAsString(number interface{}) string {
	numberAsFloat64 := recast.GetNumberAsFloat64(number)
	return fmt.Sprintf("%g", numberAsFloat64)
}

/*
GetSubString is a method which gets a substring from a given string.

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
GetStringAsBase64 is a method which encodes a string to base64.

Example:
    b64 := GetStringAsBase64("test")
*/
func GetStringAsBase64(inputString string) string {
	base64String := base64.StdEncoding.EncodeToString([]byte(inputString))
	return base64String
}

/*
GetStringFromBase64 is a method which decodes a base64 encoded string. In addition, the following should be noted:

- This method will panic if the input string is not a valid base64 encoded string.

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
GetColumnIndexBasedOnRuneIndex is a method which allows you to convert a rune index within an array into the
printed column offset at which that rune begins. The column offset is the sum of the printed widths of every rune
before the given index, so a leading run of narrow runes maps one to one while each preceding wide rune adds two.
The rune index is clamped to the closed range from zero to the array length, and passing the array length
returns the total printed width.

:param arrayOfRunes: The rune array being measured.
:param runeIndex: The rune index to locate, clamped to zero through len(arrayOfRunes).
:return: The zero-based printed column at which the rune at runeIndex begins.

Example:
    column := GetColumnIndexBasedOnRuneIndex([]rune("a中b"), 2)
*/
func GetColumnIndexBasedOnRuneIndex(arrayOfRunes []rune, runeIndex int) int {
	if runeIndex < 0 {
		runeIndex = 0
	}
	if runeIndex > len(arrayOfRunes) {
		runeIndex = len(arrayOfRunes)
	}
	column := 0
	for currentIndex := 0; currentIndex < runeIndex; currentIndex++ {
		column += GetWidthOfRuneWhenPrinted(arrayOfRunes[currentIndex])
	}
	return column
}

/*
GetRuneIndexBasedOnColumnIndex is a method which allows you to convert a printed column offset into the index of
the rune that occupies that column. A column that lands on the trailing half of a wide rune resolves to that
rune's index, and a column greater than or equal to the total printed width resolves to the array length. A
negative column resolves to zero. It is the inverse of GetColumnIndexBasedOnRuneIndex for arrays that contain no
zero-width runes.

:param arrayOfRunes: The rune array being indexed.
:param columnIndex: The printed column offset to resolve.
:return: The rune index covering the given column, or len(arrayOfRunes) when the column is past the end.

Example:
    runeIndex := GetRuneIndexBasedOnColumnIndex([]rune("a中b"), 2)
*/
func GetRuneIndexBasedOnColumnIndex(arrayOfRunes []rune, columnIndex int) int {
	if columnIndex < 0 {
		columnIndex = 0
	}
	consumedColumns := 0
	for currentIndex := 0; currentIndex < len(arrayOfRunes); currentIndex++ {
		widthOfRune := GetWidthOfRuneWhenPrinted(arrayOfRunes[currentIndex])
		if columnIndex < consumedColumns+widthOfRune {
			return currentIndex
		}
		consumedColumns += widthOfRune
	}
	return len(arrayOfRunes)
}

/*
GetRunesThatFitInColumnCountFromStart is a method which allows you to take the longest prefix of a rune array
whose printed width does not exceed a column budget. It never splits a wide rune, so when only a single column of
budget remains and the next rune is wide it stops without consuming it and without emitting any padding. Padding
a straddling boundary is the caller's responsibility via GetRunesPaddedToColumnWidth.

:param arrayOfRunes: The rune array to take a prefix from.
:param maxColumns: The maximum printed width the returned prefix may occupy.
:return: The prefix rune slice and the exact number of columns it occupies.

Example:
    prefix, columnsUsed := GetRunesThatFitInColumnCountFromStart([]rune("A中B"), 2)
*/
func GetRunesThatFitInColumnCountFromStart(arrayOfRunes []rune, maxColumns int) ([]rune, int) {
	columnsUsed := 0
	for currentIndex := 0; currentIndex < len(arrayOfRunes); currentIndex++ {
		widthOfRune := GetWidthOfRuneWhenPrinted(arrayOfRunes[currentIndex])
		if columnsUsed+widthOfRune > maxColumns {
			return GetRuneArrayCopy(arrayOfRunes[:currentIndex]), columnsUsed
		}
		columnsUsed += widthOfRune
	}
	return GetRuneArrayCopy(arrayOfRunes), columnsUsed
}

/*
GetRunesThatFitInColumnCountFromEnd is a method which allows you to take the longest suffix of a rune array whose
printed width does not exceed a column budget. It is the width inverse of GetRunesThatFitInColumnCountFromStart:
it never splits a wide rune, so when only a single column of budget remains and the preceding rune is wide it
stops without consuming it.

:param arrayOfRunes: The rune array to take a suffix from.
:param maxColumns: The maximum printed width the returned suffix may occupy.
:return: The suffix rune slice and the exact number of columns it occupies.

Example:
    suffix, columnsUsed := GetRunesThatFitInColumnCountFromEnd([]rune("A中B"), 2)
*/
func GetRunesThatFitInColumnCountFromEnd(arrayOfRunes []rune, maxColumns int) ([]rune, int) {
	columnsUsed := 0
	for currentIndex := len(arrayOfRunes) - 1; currentIndex >= 0; currentIndex-- {
		widthOfRune := GetWidthOfRuneWhenPrinted(arrayOfRunes[currentIndex])
		if columnsUsed+widthOfRune > maxColumns {
			return GetRuneArrayCopy(arrayOfRunes[currentIndex+1:]), columnsUsed
		}
		columnsUsed += widthOfRune
	}
	return GetRuneArrayCopy(arrayOfRunes), columnsUsed
}

/*
GetRunesPaddedToColumnWidth is a method which allows you to right pad a rune array with blank spaces until its
printed width reaches a minimum column count. When the array already meets or exceeds the requested width it is
returned unchanged, so this method never truncates. In addition, the following should be noted:

- The returned slice is always a fresh copy; the input array is not mutated.

:param arrayOfRunes: The rune array to pad.
:param columns: The minimum printed width the result must occupy.
:return: A new rune slice whose printed width is at least the requested column count.

Example:
    padded := GetRunesPaddedToColumnWidth([]rune("中"), 4)
*/
func GetRunesPaddedToColumnWidth(arrayOfRunes []rune, columns int) []rune {
	result := GetRuneArrayCopy(arrayOfRunes)
	paddingColumns := columns - GetWidthOfRunesWhenPrinted(arrayOfRunes)
	for currentIndex := 0; currentIndex < paddingColumns; currentIndex++ {
		result = append(result, ' ')
	}
	return result
}

/*
GetNumberOfWideCharacters is a method which gets the number of wide characters in an array of runes.

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
GetMaxCharactersThatFitInStringSize is a method which gets the longest prefix of a rune array that fits inside a
given column budget, accounting for markup. It walks the array keeping a running column count from
GetWidthOfRuneWhenPrinted; complete markup tags are copied verbatim and count zero columns, and a wide rune that
would straddle the end of the budget is dropped and replaced by a single blank so callers that measure the
result by rune count still see the slice filled to the requested width.

:param arrayOfRunes: The rune array to take a prefix from, markup included.
:param maxLengthOfString: The column budget the returned prefix must fit within.
:return: The prefix rune array, its printed width at most maxLengthOfString.

Example:
    runes := GetMaxCharactersThatFitInStringSize([]rune("{{red}}test{{/}}"), 2)
*/
func GetMaxCharactersThatFitInStringSize(arrayOfRunes []rune, maxLengthOfString int) []rune {
	numberOfColumnsUsed := 0
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
		if GetWidthOfRuneWhenPrinted(currentRune) == 2 {
			numberOfColumnsUsed = numberOfColumnsUsed + 2
			if numberOfColumnsUsed > maxLengthOfString {
				// A wide rune straddling the end is dropped and padded with one blank so the slice still
				// measures the requested column width by rune count.
				formattedArray = append(formattedArray, ' ')
				return formattedArray
			}
		} else {
			numberOfColumnsUsed++
		}

		formattedArray = append(formattedArray, currentRune)
		if numberOfColumnsUsed == maxLengthOfString {
			return formattedArray
		}

		i++
	}

	return formattedArray
}

/*
GetMaxCharactersThatFitInStringSizeReverse is a method which calculates how many runes from the end of an array
fit within a given column budget once markup has been stripped. It is the exact width inverse of the forward fit:
the count is the length of the suffix GetRunesThatFitInColumnCountFromEnd would return, so a trailing wide rune is
only counted when its full two columns fit and is never partially admitted.

:param arrayOfRunes: The rune array to measure from its end, markup included.
:param maxLengthOfString: The column budget the counted suffix must fit within.
:return: The number of trailing runes that fit within the column budget.

Example:
    count := GetMaxCharactersThatFitInStringSizeReverse([]rune("a中b"), 3)
*/
func GetMaxCharactersThatFitInStringSizeReverse(arrayOfRunes []rune, maxLengthOfString int) int {
	runesWithoutMarkup := []rune(GetTextWithoutMarkup(string(arrayOfRunes)))
	suffix, _ := GetRunesThatFitInColumnCountFromEnd(runesWithoutMarkup, maxLengthOfString)
	return len(suffix)
}

/*
GetRuneArrayCopy is a method which creates a copy of a rune array.

Example:
    copy := GetRuneArrayCopy(original)
*/
func GetRuneArrayCopy(sourceRuneArray []rune) []rune {
	copyOfRuneArray := make([]rune, len(sourceRuneArray))
	copy(copyOfRuneArray, sourceRuneArray)
	return copyOfRuneArray
}

/*
logInfo is a method which logs information to a debug file.

Example:
    logInfo("debug message")
*/
func logInfo(info string) {
	filesystem.AppendLineToFile("/tmp/debug.log", info+"\n", 0)
}

/*
GetFormattedString is a method which gets a formatted string based on specified length and alignment.

Example:
    fmtStr := GetFormattedString("test", 10, constants.AlignmentCenter)
*/
func GetFormattedString(stringToFormat string, lengthOfString int, position int) string {
	arrayOfRunes := GetRunesFromString(stringToFormat)
	return string(GetFormattedRuneArray(arrayOfRunes, lengthOfString, position))
}

/*
GetFormattedRuneArray is a method which gets a formatted rune array padded to a desired printed width according to
the requested alignment. The pad amount is the difference between the desired width and the array's printed width
with markup excluded, and it is always distributed in COLUMNS rather than rune counts so wide runes stay aligned.
When the content is already at least the desired width it is truncated to fit on a rune boundary and returned
without padding. In addition, the following should be noted:

- Centre alignment places the odd leftover column on the RIGHT side of the content.

:param arrayOfRunes: The content runes to format, markup included.
:param desiredLengthOfArray: The printed width the result should occupy.
:param textAlignment: One of the constants.Alignment values selecting how the pad is distributed.
:return: A new rune array padded or truncated to the desired printed width.

Example:
    fmtRunes := GetFormattedRuneArray([]rune("test"), 10, constants.AlignmentLeft)
*/
func GetFormattedRuneArray(arrayOfRunes []rune, desiredLengthOfArray int, textAlignment int) []rune {
	if len(arrayOfRunes) == 0 {
		return GetRunesFromString(GetFilledString(desiredLengthOfArray, " "))
	}
	// Exclude markup characters from the width calculation so pad is measured against printed columns.
	widthOfRunesWhenPrinted := GetWidthOfRunesWhenPrintedWithoutMarkup(arrayOfRunes)
	paddingSize := desiredLengthOfArray - widthOfRunesWhenPrinted
	if paddingSize <= 0 {
		return GetMaxCharactersThatFitInStringSize(arrayOfRunes, desiredLengthOfArray)
	}

	fullStringPadding := GetFilledRuneArray(paddingSize, ' ')

	formattedArrayOfRunes := []rune{}
	if textAlignment == constants.AlignmentRight {
		formattedArrayOfRunes = GetMaxCharactersThatFitInStringSize(arrayOfRunes, desiredLengthOfArray)
		formattedArrayOfRunes = append(fullStringPadding, formattedArrayOfRunes...)
	} else if textAlignment == constants.AlignmentCenter {
		leftPaddingSize := paddingSize / 2
		rightPaddingSize := paddingSize - leftPaddingSize
		formattedArrayOfRunes = append(formattedArrayOfRunes, GetFilledRuneArray(leftPaddingSize, ' ')...)
		formattedArrayOfRunes = append(formattedArrayOfRunes, arrayOfRunes...)
		formattedArrayOfRunes = append(formattedArrayOfRunes, GetFilledRuneArray(rightPaddingSize, ' ')...)
	} else if textAlignment == constants.AlignmentNoPadding {
		formattedArrayOfRunes = append(formattedArrayOfRunes, ' ')
		formattedArrayOfRunes = append(formattedArrayOfRunes, arrayOfRunes...)
		formattedArrayOfRunes = append(formattedArrayOfRunes, ' ')
	} else {
		formattedArrayOfRunes = GetMaxCharactersThatFitInStringSize(arrayOfRunes, desiredLengthOfArray)
		formattedArrayOfRunes = append(formattedArrayOfRunes, fullStringPadding...)
	}
	return formattedArrayOfRunes
}

/*
GetFilledRuneArray is a method which creates a rune array of a specified length filled with a given character.

Example:
    runes := GetFilledRuneArray(5, ' ')
*/
func GetFilledRuneArray(lengthOfString int, character rune) []rune {
	result := GetFilledString(lengthOfString, string(character))
	return GetRunesFromString((result))
}

/*
GetFilledString is a method which creates a string of a specified length filled with a given character.

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
GetLastSortedUUID is a method which generates a UUID that is prefixed to ensure it sorts correctly based on time.

Example:
    uuid := GetLastSortedUUID()
*/
func GetLastSortedUUID() string {
	id := uuid.New()
	time := fmt.Sprint(time.Now().Unix())
	return "zzzzzzz" + time + id.String()
}
