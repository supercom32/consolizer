package consolizer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/supercom32/consolizer/constants"
	"testing"
)

const TEXTBOX_TEST_SUITE_NAME = "textbox"

/*
TestTextboxMultiline is a test which verifies that a multiline textbox correctly renders several lines of text.

Example:
    Expected Inputs:
        A multiline textbox where multiple lines of text are programmatically inserted.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) showing multiple lines of rendered text.
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
TestTextboxLongLine is a test which verifies that a textbox correctly handles and renders a line of text that
exceeds its visible width.

Example:
    Expected Inputs:
        A textbox containing a line of text longer than 20 characters.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) showing horizontally scrolled or truncated text.
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
TestTextboxTallLine is a test which verifies that a textbox correctly handles and renders text content that
contains more lines than its visible height.

Example:
    Expected Inputs:
        A textbox containing more than 4 lines of text.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) showing vertically scrolled text.
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
TestTextboxLongAndTall is a test which verifies that a textbox correctly handles text content that exceeds both
its width and its height.

Example:
    Expected Inputs:
        A textbox with long lines and many rows of text.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) showing both horizontal and vertical scrolling.
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
TestTextboxLineBreak is a test which verifies that a textbox correctly handles explicit line break commands.

Example:
    Expected Inputs:
        A focused textbox where an "enter" keystroke is simulated.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) showing a new line inserted at the cursor.
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
TestTextboxDeleting is a test which verifies that deleting characters in a textbox correctly updates the text
data and rendering.

Example:
    Expected Inputs:
        A textbox where multiple "delete" commands are executed at specific cursor positions.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) with characters removed as expected.
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
TestTextboxHome is a test which verifies that the "home" key correctly moves the cursor to the beginning of a
line in the textbox.

Example:
    Expected Inputs:
        A focused textbox where the "home" command is simulated after moving to the end of a line.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) showing the cursor at column 0.
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
