package stringformat

import (
	"github.com/stretchr/testify/assert"
	"github.com/supercom32/consolizer/constants"
	"testing"
)

/*
TestIsRuneCharacterWide is a test which verifies that various rune characters are correctly identified as
wide or not.

Example:
    Expected Inputs:
        '读', 'ひ', '밥', 'W', '\u25B2', '\u25B3', '\u25B6', '\u25B7', '\u25BC', '\u25BD', '\u25C0', '\u25C1'
    Expected Outputs:
        true, true, true, false, false, false, false, false, false, false, false, false
*/
func TestIsRuneCharacterWide(test *testing.T) {
	wideChineseCharacter := '读'
	wideJapaneseCharacter := 'ひ'
	wideKoreanCharacter := '밥'
	englishCharacter := 'W'
	obtainedResult := IsRuneCharacterWide(wideChineseCharacter)
	assert.Equalf(test, true, obtainedResult, "The Chinese character specified is wide, but was not detected as such.")
	obtainedResult = IsRuneCharacterWide(wideJapaneseCharacter)
	assert.Equalf(test, true, obtainedResult, "The Japanese character specified is wide, but was not detected as such.")
	obtainedResult = IsRuneCharacterWide(wideKoreanCharacter)
	assert.Equalf(test, true, obtainedResult, "The Korean character specified is wide, but was not detected as such.")
	obtainedResult = IsRuneCharacterWide(englishCharacter)
	assert.Equalf(test, false, obtainedResult, "The English character specified is not wide, but was not detected as such.")

	// Test triangle characters (all directions)
	blackUpTriangle := '\u25B2'
	whiteUpTriangle := '\u25B3'
	blackRightTriangle := '\u25B6'
	whiteRightTriangle := '\u25B7'
	blackDownTriangle := '\u25BC'
	whiteDownTriangle := '\u25BD'
	blackLeftTriangle := '\u25C0'
	whiteLeftTriangle := '\u25C1'

	obtainedResult = IsRuneCharacterWide(blackUpTriangle)
	assert.Equalf(test, false, obtainedResult, "The BLACK UP-POINTING TRIANGLE character should not be wide.")
	obtainedResult = IsRuneCharacterWide(whiteUpTriangle)
	assert.Equalf(test, false, obtainedResult, "The WHITE UP-POINTING TRIANGLE character should not be wide.")
	obtainedResult = IsRuneCharacterWide(blackRightTriangle)
	assert.Equalf(test, false, obtainedResult, "The BLACK RIGHT-POINTING TRIANGLE character should not be wide.")
	obtainedResult = IsRuneCharacterWide(whiteRightTriangle)
	assert.Equalf(test, false, obtainedResult, "The WHITE RIGHT-POINTING TRIANGLE character should not be wide.")
	obtainedResult = IsRuneCharacterWide(blackDownTriangle)
	assert.Equalf(test, false, obtainedResult, "The BLACK DOWN-POINTING TRIANGLE character should not be wide.")
	obtainedResult = IsRuneCharacterWide(whiteDownTriangle)
	assert.Equalf(test, false, obtainedResult, "The WHITE DOWN-POINTING TRIANGLE character should not be wide.")
	obtainedResult = IsRuneCharacterWide(blackLeftTriangle)
	assert.Equalf(test, false, obtainedResult, "The BLACK LEFT-POINTING TRIANGLE character should not be wide.")
	obtainedResult = IsRuneCharacterWide(whiteLeftTriangle)
	assert.Equalf(test, false, obtainedResult, "The WHITE LEFT-POINTING TRIANGLE character should not be wide.")
}

/*
TestGetRunesFromString is a test which verifies that a string is correctly converted into an array of runes.

Example:
    Expected Inputs:
        "This is a test string to be converted into a rune array!"
    Expected Outputs:
        A rune array with a length of exactly 56.
*/
func TestGetRunesFromString(test *testing.T) {
	arrayOfRunes := GetRunesFromString("This is a test string to be converted into a rune array!")
	obtainedResult := len(arrayOfRunes)
	expectedResult := 56
	assert.Equalf(test, expectedResult, obtainedResult, "The string specified did not return a rune array of proper length!")
}

/*
TestGetIntAsString is a test which verifies that a numeric value is correctly converted to a string
representing an integer.

Example:
    Expected Inputs:
        Float value 123.456.
    Expected Outputs:
        String "123".
*/
func TestGetIntAsString(test *testing.T) {
	obtainedResult := GetIntAsString(123.456)
	expectedResult := "123"
	assert.Equalf(test, expectedResult, obtainedResult, "The number specified was not converted to a string correctly!")
}

/*
TestGetFloatAsString is a test which verifies that a numeric value is correctly converted to its string
representation as a float.

Example:
    Expected Inputs:
        Float value 123.456.
    Expected Outputs:
        String "123.456".
*/
func TestGetFloatAsString(test *testing.T) {
	obtainedResult := GetFloatAsString(123.456)
	expectedResult := "123.456"
	assert.Equalf(test, expectedResult, obtainedResult, "The number specified was not converted to a string correctly!")
}

/*
TestGetSubString is a test which verifies that a substring is correctly extracted from a given string.

Example:
    Expected Inputs:
        "This is a long string", start=5, length=2.
    Expected Outputs:
        String "is".
*/
func TestGetSubString(test *testing.T) {
	obtainedResult := GetSubString("This is a long string", 5, 2)
	expectedResult := "is"
	assert.Equalf(test, expectedResult, obtainedResult, "The substring requested was not correct!")
}

/*
TestGetStringAsBase64 is a test which verifies that a string is correctly encoded to base64.

Example:
    Expected Inputs:
        "This is base64 encoded string".
    Expected Outputs:
        The resulting base64 string "VGhpcyBpcyBiYXNlNjQgZW5jb2RlZCBzdHJpbmc=".
*/
func TestGetStringAsBase64(test *testing.T) {
	obtainedResult := GetStringAsBase64("This is base64 encoded string")
	expectedResult := "VGhpcyBpcyBiYXNlNjQgZW5jb2RlZCBzdHJpbmc="
	assert.Equalf(test, expectedResult, obtainedResult, "The base64 encoded string requested is incorrect!")
}

/*
TestGetStringFromBase64 is a test which verifies that a base64 encoded string is correctly decoded.

Example:
    Expected Inputs:
        Base64 string "VGhpcyBpcyBiYXNlNjQgZW5jb2RlZCBzdHJpbmc=".
    Expected Outputs:
        The decoded plaintext string "This is base64 encoded string".
*/
func TestGetStringFromBase64(test *testing.T) {
	obtainedResult := GetStringFromBase64("VGhpcyBpcyBiYXNlNjQgZW5jb2RlZCBzdHJpbmc=")
	expectedResult := "This is base64 encoded string"
	assert.Equalf(test, expectedResult, obtainedResult, "The converted base64 string did not return the result expected!")
}

/*
TestGetNumberOfWideCharacters is a test which verifies the counting of wide characters in a rune array. The U+2600
symbol block ballot box glyph is counted as narrow because the width oracle now matches what tcell draws, so only
the CJK ideographs, kana, and hangul in the sample are wide.

Example:
    Expected Inputs:
        "AL 读写汉字 ひらがな コンピュータワンワンローソク 보리밥보리밥☑  EX"
    Expected Outputs:
        Correct count of 28 wide characters detected (4 CJK + 4 kana + 14 katakana + 6 hangul; the ☑ is narrow).
*/
func TestGetNumberOfWideCharacters(test *testing.T) {
	arrayOfRunes := GetRunesFromString("AL 读写汉字 ひらがな コンピュータワンワンローソク 보리밥보리밥☑  EX")
	obtainedResult := GetNumberOfWideCharacters(arrayOfRunes)
	expectedResult := 28
	assert.Equalf(test, expectedResult, obtainedResult, "The number of wide characters detected did not match what was expected!")
}

/*
TestGetFormattedString is a test which verifies that a string is correctly formatted based on specified
length and alignment. The source contains one hangul syllable (wide) and one ballot box glyph from the U+2600
block which the corrected width oracle treats as narrow, so the printed content is 20 columns wide and the pad
counts are one larger than under the old assume-wide oracle.

Example:
    Expected Inputs:
        "Formatted 밥☑ String" (printed width 20), 40, Left/Right/Center alignment.
    Expected Outputs:
        Left:   "Formatted 밥☑ String" followed by 20 spaces        (39 runes, printed width 40)
        Right:  20 spaces followed by "Formatted 밥☑ String"        (39 runes, printed width 40)
        Center: 10 spaces, "Formatted 밥☑ String", 10 spaces        (39 runes, printed width 40)
*/
func TestGetFormattedString(test *testing.T) {
	obtainedResult := GetFormattedString("Formatted 밥☑ String", 40, constants.AlignmentLeft)
	expectedResult := "Formatted 밥☑ String                    "
	assert.Equalf(test, expectedResult, obtainedResult, "The formatted string obtained was not left aligned as expected.")
	obtainedSize := len(GetRunesFromString(obtainedResult))
	expectedSize := 39
	assert.Equalf(test, expectedSize, obtainedSize, "The formatted string obtained was not the right size as expected.")

	obtainedResult = GetFormattedString("Formatted 밥☑ String", 40, constants.AlignmentRight)
	expectedResult = "                    Formatted 밥☑ String"
	assert.Equalf(test, expectedResult, obtainedResult, "The formatted string obtained was not right aligned as expected.")
	obtainedSize = len(GetRunesFromString(obtainedResult))
	expectedSize = 39
	assert.Equalf(test, expectedSize, obtainedSize, "The formatted string obtained was not the right size as expected.")

	obtainedResult = GetFormattedString("Formatted 밥☑ String", 40, constants.AlignmentCenter)
	expectedResult = "          Formatted 밥☑ String          "
	assert.Equalf(test, expectedResult, obtainedResult, "The formatted string obtained was not center aligned as expected.")
	obtainedSize = len(GetRunesFromString(obtainedResult))
	expectedSize = 39
	assert.Equalf(test, expectedSize, obtainedSize, "The formatted string obtained was not the right size as expected.")
}

/*
TestGetFormattedStringWithShortAsciiInput is a test which verifies that a plain ASCII string shorter than the
desired length is padded correctly for left, right, and no-padding alignment.

Example:
    Expected Inputs:
        "Hi", 5, Left/Right/NoPadding alignment.
    Expected Outputs:
        "Hi   ", "   Hi", " Hi ".
*/
func TestGetFormattedStringWithShortAsciiInput(test *testing.T) {
	obtainedResult := GetFormattedString("Hi", 5, constants.AlignmentLeft)
	assert.Equalf(test, "Hi   ", obtainedResult, "The formatted string obtained was not left aligned as expected.")

	obtainedResult = GetFormattedString("Hi", 5, constants.AlignmentRight)
	assert.Equalf(test, "   Hi", obtainedResult, "The formatted string obtained was not right aligned as expected.")

	obtainedResult = GetFormattedString("Hi", 5, constants.AlignmentNoPadding)
	assert.Equalf(test, " Hi ", obtainedResult, "The formatted string obtained was not formatted with no padding as expected.")
}
