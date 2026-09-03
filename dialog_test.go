package consolizer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const DIALOG_TEST_SUITE_NAME = "dialog"

/*
TestPrintDialogWithTextStyles is a test which allows you to verify that the printDialog function correctly applies text
styles from markup tags and handles word wrapping.

Example:

	Expected Inputs:
	    Input: "This is a {{redColor}}sample{{/}} line of text. This line will {{redColor}}automatically{{/}} wrap without cutting words."

	Expected Outputs:
	    Output: Screen content matches expected ANSI string (Base64 encoded).
*/
func TestPrintDialogWithTextStyles(test *testing.T) {
	commonResource.isDebugEnabled = true
	layerWidth := 40
	layerHeight := 10
	textStyleAlias := "redColor"
	InitializeTerminal(layerWidth, layerHeight)
	layerAlias1 := AddLayer(0, 0, layerWidth, layerHeight, 1, nil)
	layerAlias1.Color(4, 6)
	layerAlias1.FillLayer("a1a2a3a4a5")
	layerAlias1.Color(8, 9)
	attributeEntry := NewTextStyle()
	attributeEntry.ForegroundColor = GetRGBColor(255, 0, 0)
	AddTextStyle(textStyleAlias, attributeEntry)
	stringToPrint := "This is a {{redColor}}sample{{/}} line of text. This line will {{redColor}}automatically{{/}} wrap without cutting words."
	layerAlias1.PrintDialog(2, 2, 20, 0, false, stringToPrint)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, DIALOG_TEST_SUITE_NAME, "TestPrintDialogWithTextStyles", obtainedValue)
	expectedValue := LoadMasterImage(DIALOG_TEST_SUITE_NAME, "TestPrintDialogWithTextStyles")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestRuneAndMarkupHandling is a test which verifies that methods correctly handle multi-byte UTF-8 characters
and markup tags without boundary errors or incorrect length calculations.

Example:

	Expected Inputs:
	    Various strings containing emojis and Japanese/Chinese characters mixed with {{style}} tags.
	Expected Outputs:
	    Correct non-markup text, correct word lengths (counting runes), and correctly identified attribute tags.
*/
func TestRuneAndMarkupHandling(test *testing.T) {
	// 1. Test GetNonMarkupText with multi-byte characters
	testCases := []struct {
		input    string
		expected string
	}{
		{"Hello {{red}}World{{/}}", "Hello World"},
		{"こんにちは {{blue}}世界{{/}}", "こんにちは 世界"},
		{"Emoji Test 🌟 {{green}}Sparkles{{/}} 🌟", "Emoji Test 🌟 Sparkles 🌟"},
		{"Unclosed {{tag and multi-byte 汉字", "Unclosed {{tag and multi-byte 汉字"},
	}

	for _, tc := range testCases {
		result := GetNonMarkupText(tc.input)
		assert.Equal(test, tc.expected, result, "GetNonMarkupText failed for: "+tc.input)
	}

	// 2. Test getLengthOfNextWord with multi-byte characters
	wordTestCases := []struct {
		input    string
		start    int
		expected int
	}{
		{"Next {{red}}Word{{/}}", 5, 4},      // "Word"
		{"こんにちは {{blue}}世界{{/}} test", 6, 2}, // "世界" (indices in runes)
		{"🌟 {{green}}Sparkle{{/}}", 2, 7},    // "Sparkle"
	}

	for _, tc := range wordTestCases {
		result := getLengthOfNextWord(tc.input, tc.start)
		assert.Equal(test, tc.expected, result, "getLengthOfNextWord failed for: "+tc.input)
	}

	// 3. Test getAttributeTag with multi-byte characters
	tagTestCases := []struct {
		input    string
		start    int
		expected string
	}{
		{"Text {{red}} More", 5, "{{red}}"},
		{"汉字 {{blue}} 更多", 3, "{{blue}}"},
		{"Emoji 🌟 {{green}} More", 8, "{{green}}"},
	}

	for _, tc := range tagTestCases {
		result := getAttributeTag(tc.input, tc.start)
		assert.Equal(test, tc.expected, result, "getAttributeTag failed for: "+tc.input)
	}
}
