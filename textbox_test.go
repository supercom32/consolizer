package consolizer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/supercom32/consolizer/constants"
	"testing"
)

const TEXTBOX_TEST_SUITE_NAME = "textbox"

/*
TestTextboxMultiline is a test which verifies that a multiline textbox is rendered correctly.

Example:

	Expected Inputs:
	    A multiline textbox with several lines of text.

	Expected Outputs:
	    A rendered multiline textbox.
*/
func TestTextboxMultiline(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	textboxInstance := layer1.AddTextbox(styleEntry, 2, 2, 20, 4, false)
	setFocusedControl(layer1.layerAlias, textboxInstance.controlAlias, constants.CellTypeTextbox)
	textbox.UpdateKeyboardEventTextboxWithString("abcdefghijklmnop")
	textbox.UpdateKeyboardEvent([]rune("enter"))
	textbox.UpdateKeyboardEventTextboxWithString("abcdefghijklmnop")
	textbox.UpdateKeyboardEvent([]rune("enter"))
	textbox.UpdateKeyboardEventTextboxWithString("abcdefghijklmnop")
	textbox.UpdateKeyboardEvent([]rune("enter"))
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TEXTBOX_TEST_SUITE_NAME, "TestTextboxMultiline", obtainedValue)
	expectedValue := LoadMasterImage(TEXTBOX_TEST_SUITE_NAME, "TestTextboxMultiline")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestTextboxLongLine is a test which verifies that a textbox with a long line is rendered correctly.

Example:

	Expected Inputs:
	    A textbox with a line of text exceeding its width.

	Expected Outputs:
	    A rendered textbox with horizontally scrolled text.
*/
func TestTextboxLongLine(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	textboxInstance := layer1.AddTextbox(styleEntry, 2, 2, 20, 4, false)
	setFocusedControl(layer1.layerAlias, textboxInstance.controlAlias, constants.CellTypeTextbox)
	textbox.UpdateKeyboardEventTextboxWithString("abcdefghijklmnop")
	textbox.UpdateKeyboardEvent([]rune("enter"))
	textbox.UpdateKeyboardEventTextboxWithString("abcdefghijklmnopqrstuvwxyz")
	textbox.UpdateKeyboardEvent([]rune("enter"))
	textbox.UpdateKeyboardEventTextboxWithString("abcdefghijklmnop")
	textbox.UpdateKeyboardEvent([]rune("enter"))
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TEXTBOX_TEST_SUITE_NAME, "TestTextboxLongLine", obtainedValue)
	expectedValue := LoadMasterImage(TEXTBOX_TEST_SUITE_NAME, "TestTextboxLongLine")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestTextboxTallLine is a test which verifies that a textbox with many lines is rendered correctly.

Example:

	Expected Inputs:
	    A textbox with more lines of text than its height.

	Expected Outputs:
	    A rendered textbox with vertically scrolled text.
*/
func TestTextboxTallLine(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	textboxInstance := layer1.AddTextbox(styleEntry, 2, 2, 20, 4, false)
	setFocusedControl(layer1.layerAlias, textboxInstance.controlAlias, constants.CellTypeTextbox)
	textbox.UpdateKeyboardEventTextboxWithString("abcdefghijklmnop")
	textbox.UpdateKeyboardEvent([]rune("enter"))
	textbox.UpdateKeyboardEventTextboxWithString("abcdefghijklmnop")
	textbox.UpdateKeyboardEvent([]rune("enter"))
	textbox.UpdateKeyboardEventTextboxWithString("abcd")
	textbox.UpdateKeyboardEvent([]rune("enter"))
	textbox.UpdateKeyboardEventTextboxWithString("abcdefghijklmnop")
	textbox.UpdateKeyboardEvent([]rune("enter"))
	textbox.UpdateKeyboardEventTextboxWithString("END")
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TEXTBOX_TEST_SUITE_NAME, "TestTextboxTallLine", obtainedValue)
	expectedValue := LoadMasterImage(TEXTBOX_TEST_SUITE_NAME, "TestTextboxTallLine")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestTextboxLongAndTall is a test which verifies that a textbox with both long and many lines is rendered correctly.

Example:

	Expected Inputs:
	    A textbox with lines exceeding both width and height.

	Expected Outputs:
	    A rendered textbox with both horizontally and vertically scrolled text.
*/
func TestTextboxLongAndTall(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	textboxInstance := layer1.AddTextbox(styleEntry, 2, 2, 20, 4, false)
	setFocusedControl(layer1.layerAlias, textboxInstance.controlAlias, constants.CellTypeTextbox)
	textbox.UpdateKeyboardEventTextboxWithString("abcdefghijklmnop")
	textbox.UpdateKeyboardEvent([]rune("enter"))
	textbox.UpdateKeyboardEventTextboxWithString("abcdefghijklmnopqrstuvwxyz")
	textbox.UpdateKeyboardEvent([]rune("enter"))
	textbox.UpdateKeyboardEventTextboxWithString("abcd")
	textbox.UpdateKeyboardEvent([]rune("enter"))
	textbox.UpdateKeyboardEventTextboxWithString("abcdefghijklmnop")
	textbox.UpdateKeyboardEvent([]rune("enter"))
	textbox.UpdateKeyboardEventTextboxWithString("END")
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TEXTBOX_TEST_SUITE_NAME, "TestTextboxLongAndTall", obtainedValue)
	expectedValue := LoadMasterImage(TEXTBOX_TEST_SUITE_NAME, "TestTextboxLongAndTall")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestTextboxLineBreak is a test which verifies that a textbox handles line breaks correctly.

Example:

	Expected Inputs:
	    A textbox where an enter key is pressed to create a new line.

	Expected Outputs:
	    A rendered textbox with a new line inserted.
*/
func TestTextboxLineBreak(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	textboxInstance := layer1.AddTextbox(styleEntry, 2, 2, 20, 4, false)
	setFocusedControl(layer1.layerAlias, textboxInstance.controlAlias, constants.CellTypeTextbox)
	textbox.UpdateKeyboardEventTextboxWithString("abcdefghijklmnop")
	textbox.UpdateKeyboardEvent([]rune("enter"))
	textbox.UpdateKeyboardEventTextboxWithString("abcdefghijklmnopqrstuvwxyz")
	textbox.UpdateKeyboardEvent([]rune("enter"))
	textbox.UpdateKeyboardEventTextboxWithString("abcd")
	textbox.UpdateKeyboardEvent([]rune("enter"))
	textbox.UpdateKeyboardEventTextboxWithString("abcdefghijklmnop")
	textbox.UpdateKeyboardEvent([]rune("enter"))
	textbox.UpdateKeyboardEventTextboxWithString("END")
	textbox.UpdateKeyboardEvent([]rune("up"))
	textbox.UpdateKeyboardEvent([]rune("enter"))
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TEXTBOX_TEST_SUITE_NAME, "TestTextboxLineBreak", obtainedValue)
	expectedValue := LoadMasterImage(TEXTBOX_TEST_SUITE_NAME, "TestTextboxLineBreak")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestTextboxDeleting is a test which verifies that deleting text in a textbox works correctly.

Example:

	Expected Inputs:
	    A textbox where the delete key is pressed to remove characters.

	Expected Outputs:
	    A rendered textbox with the specified characters removed.
*/
func TestTextboxDeleting(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	textboxInstance := layer1.AddTextbox(styleEntry, 2, 2, 20, 4, false)
	setFocusedControl(layer1.layerAlias, textboxInstance.controlAlias, constants.CellTypeTextbox)
	textbox.UpdateKeyboardEventTextboxWithString("abcdefghijklmnop")
	textbox.UpdateKeyboardEvent([]rune("enter"))
	textbox.UpdateKeyboardEventTextboxWithString("abcdefghijklmnopqrstuvwxyz")
	textbox.UpdateKeyboardEvent([]rune("enter"))
	textbox.UpdateKeyboardEventTextboxWithString("abcd")
	textbox.UpdateKeyboardEventTextboxWithCommands("up", "end", "left", "left", "left", "delete", "delete", "delete", "delete")
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TEXTBOX_TEST_SUITE_NAME, "TestTextboxDeleting", obtainedValue)
	expectedValue := LoadMasterImage(TEXTBOX_TEST_SUITE_NAME, "TestTextboxDeleting")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestTextboxHome is a test which verifies that the home key in a textbox works correctly.

Example:

	Expected Inputs:
	    A textbox where the home key is pressed to move the cursor to the beginning of the line.

	Expected Outputs:
	    A rendered textbox with the cursor at the beginning of the line.
*/
func TestTextboxHome(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	textboxInstance := layer1.AddTextbox(styleEntry, 2, 2, 20, 4, false)
	setFocusedControl(layer1.layerAlias, textboxInstance.controlAlias, constants.CellTypeTextbox)
	textbox.UpdateKeyboardEventTextboxWithString("abcdefghijklmnop")
	textbox.UpdateKeyboardEvent([]rune("enter"))
	textbox.UpdateKeyboardEventTextboxWithString("abcdefghijklmnopqrstuvwxyz")
	textbox.UpdateKeyboardEvent([]rune("enter"))
	textbox.UpdateKeyboardEventTextboxWithString("abcd")
	textbox.UpdateKeyboardEventTextboxWithCommands("up", "end", "left", "left", "left", "delete", "delete", "delete", "delete", "home")
	textbox.updateMouseEvent()
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TEXTBOX_TEST_SUITE_NAME, "TestTextboxHome", obtainedValue)
	expectedValue := LoadMasterImage(TEXTBOX_TEST_SUITE_NAME, "TestTextboxHome")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}
