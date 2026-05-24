package consolizer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/supercom32/consolizer/constants"
	"testing"
)

const TEXT_FIELD_TEST_SUITE_NAME = "text_field"

/*
TestTextFieldDefaultText is a test which verifies that a text field with default text is rendered correctly.

Example:

	Expected Inputs:
	    A text field with default text.

	Expected Outputs:
	    A rendered text field with the default text.
*/
func TestTextFieldDefaultText(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	textFieldInstance := layer1.AddTextField(styleEntry, 2, 2, 20, 10, false, "default", true)
	setFocusedControl(layer1.layerAlias, textFieldInstance.controlAlias, constants.CellTypeTextField)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TEXT_FIELD_TEST_SUITE_NAME, "TestTextFieldDefaultText", obtainedValue)
	expectedValue := LoadMasterImage(TEXT_FIELD_TEST_SUITE_NAME, "TestTextFieldDefaultText")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestTextFieldPasswordText is a test which verifies that a text field in password mode masks its text.

Example:

	Expected Inputs:
	    A text field in password mode with some text.

	Expected Outputs:
	    A rendered text field where the text is replaced by masks.
*/
func TestTextFieldPasswordText(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	textFieldInstance := layer1.AddTextField(styleEntry, 2, 2, 20, 10, true, "default", true)
	setFocusedControl(layer1.layerAlias, textFieldInstance.controlAlias, constants.CellTypeTextField)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TEXT_FIELD_TEST_SUITE_NAME, "TestTextFieldPasswordText", obtainedValue)
	expectedValue := LoadMasterImage(TEXT_FIELD_TEST_SUITE_NAME, "TestTextFieldPasswordText")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestTextFieldLongText is a test which verifies that a text field handles scrolling for long text correctly.

Example:

	Expected Inputs:
	    A text field with text longer than its display width.

	Expected Outputs:
	    A rendered text field showing the scrolled portion of the text.
*/
func TestTextFieldLongText(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	textFieldInstance := layer1.AddTextField(styleEntry, 2, 2, 20, 10, false, "this is a long string of text which i know is long.", true)
	setFocusedControl(layer1.layerAlias, textFieldInstance.controlAlias, constants.CellTypeTextField)
	TextField.updateKeyboardEvent([]rune("end"))
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TEXT_FIELD_TEST_SUITE_NAME, "TestTextFieldLongText", obtainedValue)
	expectedValue := LoadMasterImage(TEXT_FIELD_TEST_SUITE_NAME, "TestTextFieldLongText")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestTextFieldLongTypedText is a test which verifies that a text field correctly displays text typed by the user.

Example:

	Expected Inputs:
	    A text field where a long string is typed.

	Expected Outputs:
	    A rendered text field showing the end portion of the typed text.
*/
func TestTextFieldLongTypedText(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	textFieldInstance := layer1.AddTextField(styleEntry, 2, 2, 20, 30, false, "", true)
	setFocusedControl(layer1.layerAlias, textFieldInstance.controlAlias, constants.CellTypeTextField)
	TextField.updateKeyboardEventTextboxWithString("abcdefghijklmnopqrstuvwxyz")
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TEXT_FIELD_TEST_SUITE_NAME, "TestTextFieldLongTypedText", obtainedValue)
	expectedValue := LoadMasterImage(TEXT_FIELD_TEST_SUITE_NAME, "TestTextFieldLongTypedText")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestTextFieldHomeKey is a test which verifies that the Home key correctly moves the cursor to the beginning of the text field.

Example:

	Expected Inputs:
	    A text field with typed text, followed by a Home key press.

	Expected Outputs:
	    A rendered text field with the cursor at the first character.
*/
func TestTextFieldHomeKey(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	textFieldInstance := layer1.AddTextField(styleEntry, 2, 2, 20, 30, false, "", true)
	setFocusedControl(layer1.layerAlias, textFieldInstance.controlAlias, constants.CellTypeTextField)
	TextField.updateKeyboardEventTextboxWithString("abcdefghijklmnopqrstuvwxyz")
	TextField.updateKeyboardEvent([]rune("home"))
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TEXT_FIELD_TEST_SUITE_NAME, "TestTextFieldHomeKey", obtainedValue)
	expectedValue := LoadMasterImage(TEXT_FIELD_TEST_SUITE_NAME, "TestTextFieldHomeKey")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestTextFieldInsert is a test which verifies that text can be correctly inserted into an existing string.

Example:

	Expected Inputs:
	    A text field with text, followed by inserting a new string at a specific position.

	Expected Outputs:
	    A rendered text field containing both the original and inserted text.
*/
func TestTextFieldInsert(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	textFieldInstance := layer1.AddTextField(styleEntry, 2, 2, 20, 50, false, "", true)
	setFocusedControl(layer1.layerAlias, textFieldInstance.controlAlias, constants.CellTypeTextField)
	TextField.updateKeyboardEventTextboxWithString("abcdefghijklmnopqrstuvwxyz")
	TextField.updateKeyboardEvent([]rune("home"))
	TextField.updateKeyboardEventTextboxWithCommands("right", "right", "right", "right", "right")
	TextField.updateKeyboardEventTextboxWithString("_INSERTED_")
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TEXT_FIELD_TEST_SUITE_NAME, "TestTextFieldInsert", obtainedValue)
	expectedValue := LoadMasterImage(TEXT_FIELD_TEST_SUITE_NAME, "TestTextFieldInsert")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestTextFieldBackspace is a test which verifies that the Backspace key correctly removes characters from the text field.

Example:

	Expected Inputs:
	    A text field with text, followed by multiple Backspace key presses.

	Expected Outputs:
	    A rendered text field with the corresponding characters removed.
*/
func TestTextFieldBackspace(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	textFieldInstance := layer1.AddTextField(styleEntry, 2, 2, 20, 90, false, "", true)
	setFocusedControl(layer1.layerAlias, textFieldInstance.controlAlias, constants.CellTypeTextField)
	TextField.updateKeyboardEventTextboxWithString("abcdefghijklmnopqrstuvwxyz")
	TextField.updateKeyboardEvent([]rune("home"))
	TextField.updateKeyboardEventTextboxWithCommands("right", "right", "right", "right", "right", "backspace", "backspace", "backspace", "backspace")
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TEXT_FIELD_TEST_SUITE_NAME, "TestTextFieldBackspace", obtainedValue)
	expectedValue := LoadMasterImage(TEXT_FIELD_TEST_SUITE_NAME, "TestTextFieldBackspace")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestTextFieldDelete is a test which verifies that the Delete key correctly removes characters following the cursor.

Example:

	Expected Inputs:
	    A text field with text, followed by cursor movement and multiple Delete key presses.

	Expected Outputs:
	    A rendered text field with the corresponding characters removed.
*/
func TestTextFieldDelete(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	textFieldInstance := layer1.AddTextField(styleEntry, 2, 2, 20, 90, false, "", true)
	setFocusedControl(layer1.layerAlias, textFieldInstance.controlAlias, constants.CellTypeTextField)
	TextField.updateKeyboardEventTextboxWithString("abcdefghijklmnopqrstuvwxyz")
	TextField.updateKeyboardEvent([]rune("home"))
	TextField.updateKeyboardEventTextboxWithCommands("right", "right", "right", "right", "right", "delete", "delete", "delete", "delete")
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TEXT_FIELD_TEST_SUITE_NAME, "TestTextFieldDelete", obtainedValue)
	expectedValue := LoadMasterImage(TEXT_FIELD_TEST_SUITE_NAME, "TestTextFieldDelete")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestTextFieldDeleteAtEnd is a test which verifies that the Delete key behaves correctly at the end of the text.

Example:

	Expected Inputs:
	    A text field with text, followed by cursor movement to the end and multiple Delete key presses.

	Expected Outputs:
	    A rendered text field where Delete key presses at the end have no effect.
*/
func TestTextFieldDeleteAtEnd(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	textFieldInstance := layer1.AddTextField(styleEntry, 2, 2, 20, 90, false, "", true)
	setFocusedControl(layer1.layerAlias, textFieldInstance.controlAlias, constants.CellTypeTextField)
	TextField.updateKeyboardEventTextboxWithString("abcdefghijklmnopqrstuvwxyz")
	TextField.updateKeyboardEvent([]rune("end"))
	TextField.updateKeyboardEventTextboxWithCommands("left", "left", "left", "left", "delete", "delete", "delete", "delete", "delete", "delete")
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TEXT_FIELD_TEST_SUITE_NAME, "TestTextFieldDeleteAtEnd", obtainedValue)
	expectedValue := LoadMasterImage(TEXT_FIELD_TEST_SUITE_NAME, "TestTextFieldDeleteAtEnd")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestTextFieldMaxFieldLimit is a test which verifies that the text field respects the maximum character limit.

Example:

	Expected Inputs:
	    A text field with a set maximum character limit, followed by typing text exceeding that limit.

	Expected Outputs:
	    A rendered text field containing only the text up to the maximum limit.
*/
func TestTextFieldMaxFieldLimit(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	textFieldInstance := layer1.AddTextField(styleEntry, 2, 2, 20, 10, false, "", true)
	setFocusedControl(layer1.layerAlias, textFieldInstance.controlAlias, constants.CellTypeTextField)
	TextField.updateKeyboardEventTextboxWithString("abcdefghijklmnopqrstuvwxyz")
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TEXT_FIELD_TEST_SUITE_NAME, "TestTextFieldMaxFieldLimit", obtainedValue)
	expectedValue := LoadMasterImage(TEXT_FIELD_TEST_SUITE_NAME, "TestTextFieldMaxFieldLimit")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestTextFieldBackspaceStop is a test which verifies that the Backspace key correctly stops at the beginning of the text field.

Example:

	Expected Inputs:
	    A text field with text, followed by cursor movement and more Backspace key presses than characters.

	Expected Outputs:
	    A rendered text field where excess Backspace key presses have no effect.
*/
func TestTextFieldBackspaceStop(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	textFieldInstance := layer1.AddTextField(styleEntry, 2, 2, 20, 70, false, "", true)
	setFocusedControl(layer1.layerAlias, textFieldInstance.controlAlias, constants.CellTypeTextField)
	TextField.updateKeyboardEventTextboxWithString("abcdefghijklmnopqrstuvwxyz")
	TextField.updateKeyboardEvent([]rune("home"))
	TextField.updateKeyboardEventTextboxWithCommands("right", "right", "right", "right", "backspace", "backspace", "backspace", "backspace", "backspace", "backspace")
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TEXT_FIELD_TEST_SUITE_NAME, "TestTextFieldBackspaceStop", obtainedValue)
	expectedValue := LoadMasterImage(TEXT_FIELD_TEST_SUITE_NAME, "TestTextFieldBackspaceStop")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}
