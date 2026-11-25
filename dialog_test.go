package consolizer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/types"
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
	defer RestoreTerminalSettings()
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
TestRuneAndMarkupHandling is a test which verifies that GetNonMarkupText removes markup tags from text holding
multi-byte UTF-8 characters without boundary errors, and leaves an unclosed tag as regular text.

Example:

	Expected Inputs:
	    "Hello {{red}}World{{/}}", "こんにちは {{blue}}世界{{/}}", "Emoji Test 🌟 {{green}}Sparkles{{/}} 🌟" and
	    "Unclosed {{tag and multi-byte 汉字".
	Expected Outputs:
	    "Hello World", "こんにちは 世界", "Emoji Test 🌟 Sparkles 🌟" and "Unclosed {{tag and multi-byte 汉字".
*/
func TestRuneAndMarkupHandling(test *testing.T) {
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
}

/*
TestGetDialogAttributeEntryKeepsCellFields is a test which verifies that a markup tag naming a text style takes only
the style's visual fields, keeping the default entry's cell and control fields, while the closing tag, an empty tag
and an unknown tag return the default entry unchanged.

Example:

	Expected Inputs:
	    A style "cellFieldsWarn" with a red foreground, a blue background, bold, underline, reverse, blink and italic.
	    A default entry with foreground (1, 2, 3), background (4, 5, 6), alphas 0.5, both transparency flags set,
	    CellType CellTypeSelectorItem, CellControlAlias "list", CellControlId 3, CellControlLocation 7, CellUserId 9
	    and CellUserAlias "user". The tags "cellFieldsWarn", "/", "" and "cellFieldsMissing".
	Expected Outputs:
	    "cellFieldsWarn" returns the style's colours, flags, alphas of 1 and cleared transparency flags, with
	    CellType CellTypeSelectorItem, CellControlAlias "list", CellControlId 3, CellControlLocation 7, CellUserId 9
	    and CellUserAlias "user". Every other tag returns the default entry exactly.
*/
func TestGetDialogAttributeEntryKeepsCellFields(test *testing.T) {
	textStyleAlias := "cellFieldsWarn"
	textStyleEntry := NewTextStyle()
	textStyleEntry.ForegroundColor = GetRGBColor(255, 0, 0)
	textStyleEntry.BackgroundColor = GetRGBColor(0, 0, 255)
	textStyleEntry.IsBold = true
	textStyleEntry.IsUnderlined = true
	textStyleEntry.IsReversed = true
	textStyleEntry.IsBlinking = true
	textStyleEntry.IsItalic = true
	AddTextStyle(textStyleAlias, textStyleEntry)
	test.Cleanup(func() { DeleteTextStyle(textStyleAlias) })

	defaultAttributeEntry := types.NewAttributeEntry()
	defaultAttributeEntry.ForegroundColor = GetRGBColor(1, 2, 3)
	defaultAttributeEntry.BackgroundColor = GetRGBColor(4, 5, 6)
	defaultAttributeEntry.ForegroundAlphaValue = 0.5
	defaultAttributeEntry.BackgroundAlphaValue = 0.5
	defaultAttributeEntry.IsForegroundTransparent = true
	defaultAttributeEntry.IsBackgroundTransparent = true
	defaultAttributeEntry.CellType = constants.CellTypeSelectorItem
	defaultAttributeEntry.CellControlAlias = "list"
	defaultAttributeEntry.CellControlId = 3
	defaultAttributeEntry.CellControlLocation = 7
	defaultAttributeEntry.CellUserId = 9
	defaultAttributeEntry.CellUserAlias = "user"

	expectedAttributeEntry := defaultAttributeEntry
	expectedAttributeEntry.ForegroundColor = textStyleEntry.ForegroundColor
	expectedAttributeEntry.BackgroundColor = textStyleEntry.BackgroundColor
	expectedAttributeEntry.IsBold = true
	expectedAttributeEntry.IsUnderlined = true
	expectedAttributeEntry.IsReversed = true
	expectedAttributeEntry.IsBlinking = true
	expectedAttributeEntry.IsItalic = true
	expectedAttributeEntry.ForegroundAlphaValue = 1
	expectedAttributeEntry.BackgroundAlphaValue = 1
	expectedAttributeEntry.IsForegroundTransparent = false
	expectedAttributeEntry.IsBackgroundTransparent = false

	assert.Equal(test, expectedAttributeEntry, getDialogAttributeEntry(textStyleAlias, defaultAttributeEntry), "a style tag should apply only the style's visual fields")
	for _, attributeTag := range []string{"/", "", "cellFieldsMissing"} {
		assert.Equal(test, defaultAttributeEntry, getDialogAttributeEntry(attributeTag, defaultAttributeEntry), "tag \""+attributeTag+"\" should return the default entry")
	}
}

/*
TestGetAttributeTagContent is a test which verifies that the style name is taken from inside a full markup tag, and
that anything not wrapped in "{{" and "}}" yields an empty string.

Example:

	Expected Inputs:
	    "{{red}}", "{{/}}", "{{汉字}}", "{{}}", "", "red", "{{red" and "red}}".
	Expected Outputs:
	    "red", "/", "汉字", and an empty string for each of the rest.
*/
func TestGetAttributeTagContent(test *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"{{red}}", "red"},
		{"{{/}}", "/"},
		{"{{汉字}}", "汉字"},
		{"{{}}", ""},
		{"", ""},
		{"red", ""},
		{"{{red", ""},
		{"red}}", ""},
	}
	for _, testCase := range testCases {
		assert.Equal(test, testCase.expected, getAttributeTagContent(testCase.input), "getAttributeTagContent failed for: "+testCase.input)
	}
}

/*
TestPrintMarkupAppliesStylesAndKeepsCellFields is a test which verifies that the package-level printMarkup colours
text inside a style tag and keeps the cell fields of the attribute entry it was given. Before the fix the tag was
looked up with its braces, as "{{dialogMarkupWarn}}", so no style ever matched and the text stayed uncoloured.

Example:

	Expected Inputs:
	    A style "dialogMarkupWarn" (red) and the layer's default attribute entry with CellType CellTypeSelectorItem and
	    CellControlAlias "printed". printMarkup writes "Go {{dialogMarkupWarn}}Red{{/}} end" at (1, 1) with a line
	    width of 30.
	Expected Outputs:
	    Row 1 of the layer shows "Go Red end" from column 1. All ten cells have CellType CellTypeSelectorItem and
	    CellControlAlias "printed", and only the three cells of "Red" are red.
*/
func TestPrintMarkupAppliesStylesAndKeepsCellFields(test *testing.T) {
	layerAlias, _ := setupInputTest(test)
	warnColor := GetRGBColor(255, 0, 0)
	addTestTextStyle(test, "dialogMarkupWarn", warnColor)
	layerEntry := layerRegistry.Get(layerAlias)
	attributeEntry := layerEntry.DefaultAttribute
	attributeEntry.CellType = constants.CellTypeSelectorItem
	attributeEntry.CellControlAlias = "printed"

	printMarkup(layerEntry, attributeEntry, 1, 1, 30, "Go {{dialogMarkupWarn}}Red{{/}} end")

	getCell := func(column int) types.CharacterEntryType { return layerEntry.CharacterMemory[1][column] }
	assertMarkupRowCells(test, "printMarkup", getCell, 1, constants.CellTypeSelectorItem, "printed", warnColor)
}

/*
TestPrintMarkupWithoutWrapAndLeadingTags is a test which verifies that the package-level printMarkup prints a whole
line when word wrapping is off, never prints a tag that opens the text or follows another tag, and gives a wide
character both of its columns.

Example:

	Expected Inputs:
	    Styles "dialogLeadA" (green) and "dialogLeadB" (red). printMarkup with a line width of 0 writes
	    "{{dialogLeadA}}{{dialogLeadB}}Hi 汉字 ok" at (1, 1) on the layer.
	Expected Outputs:
	    Row 1 shows "Hi 汉字 ok" from column 1 on that single row, the wide characters filling columns 4 to 7, and
	    "H" is red. Row 2 is untouched.
*/
func TestPrintMarkupWithoutWrapAndLeadingTags(test *testing.T) {
	layerAlias, _ := setupInputTest(test)
	addTestTextStyle(test, "dialogLeadA", GetRGBColor(0, 255, 0))
	redColor := GetRGBColor(255, 0, 0)
	addTestTextStyle(test, "dialogLeadB", redColor)
	layerEntry := layerRegistry.Get(layerAlias)
	untouchedCharacter := layerEntry.CharacterMemory[2][1].Character

	printMarkup(layerEntry, layerEntry.DefaultAttribute, 1, 1, 0, "{{dialogLeadA}}{{dialogLeadB}}Hi 汉字 ok")

	getCell := func(column int) types.CharacterEntryType { return layerEntry.CharacterMemory[1][column] }
	assert.Equal(test, 11, assertCellsShowText(test, "printMarkup without wrapping", getCell, 1, "Hi 汉字 ok"), "the text should end at column 11")
	assert.Equal(test, redColor, getCell(1).AttributeEntry.ForegroundColor, "\"H\" should use the last tag's style")
	assert.Equal(test, untouchedCharacter, layerEntry.CharacterMemory[2][1].Character, "nothing should wrap onto row 2")
}
