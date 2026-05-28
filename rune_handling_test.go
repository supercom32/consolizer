package consolizer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
TestRuneAndMarkupHandling is a test which verifies that methods correctly handle multi-byte UTF-8 characters 
and markup tags without boundary errors or incorrect length calculations.

Example:
    Expected Inputs:
        Various strings containing emojis and Chinese characters mixed with {{style}} tags.
    Expected Outputs:
        Correct non-markup text and word lengths.
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
		input string
		start int
		expected int
	}{
		{"Next {{red}}Word{{/}}", 5, 4}, // "Word"
		{"こんにちは {{blue}}世界{{/}} test", 6, 2}, // "世界" (indices in runes)
		{"🌟 {{green}}Sparkle{{/}}", 2, 7}, // "Sparkle"
	}

	for _, tc := range wordTestCases {
		result := getLengthOfNextWord(tc.input, tc.start)
		assert.Equal(test, tc.expected, result, "getLengthOfNextWord failed for: "+tc.input)
	}

	// 3. Test getAttributeTag with multi-byte characters
	tagTestCases := []struct {
		input string
		start int
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

/*
TestCalculateWordWidthWithRunes is a test which verifies that calculateWordWidth handles rune arrays 
containing multi-byte characters and markup correctly.

Example:
    Expected Inputs:
        Rune arrays with Chinese characters and markup tags.
    Expected Outputs:
        Correct word width excluding markup.
*/
func TestCalculateWordWidthWithRunes(test *testing.T) {
	input := []rune(" こんにちは {{red}}世界{{/}} test")
	
	// Test width of "世界" starting from the space at index 6
	width := calculateWordWidth(input, 6, true)
	assert.Equal(test, 2, width, "calculateWordWidth failed to calculate correct width for '世界'")
	
	// Test width without markup
	inputNoMarkup := []rune(" Hello World")
	widthNoMarkup := calculateWordWidth(inputNoMarkup, 6, false)
	assert.Equal(test, 5, widthNoMarkup, "calculateWordWidth failed for standard ASCII")
}
