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

	Input: "This is a {{redColor}}sample{{/}} line of text. This line will {{redColor}}automatically{{/}} wrap without cutting words."
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
