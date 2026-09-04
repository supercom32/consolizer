package consolizer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/stringformat"
	"testing"
)

const TEXT_FIELD_TEST_SUITE_NAME = "text_field"

/*
TestTextFieldDefaultText is a test which verifies that a text field with default text is rendered correctly.

Example:
    Expected Inputs:
        A text field initialized with the default string "default".
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) showing the default text.
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
        A text field with password protection enabled and default text "default".
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) where "default" is replaced by masks.
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
        A text field containing text significantly longer than its 20-character display width.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) showing the scrolled viewport.
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
        A text field where the sequence "abcdefghijklmnopqrstuvwxyz" is programmatically typed.
    Expected Outputs:
        Screen content matches expected ANSI string (Base64 encoded) showing the end portion of the typed alphabet.
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
        A text field with full alphabet text followed by a "home" keystroke.
    Expected Outputs:
        The cursor is positioned at index 0 and the viewport scrolls back to the beginning.
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
        Alphabet string followed by moving cursor to index 5 and inserting "_INSERTED_".
    Expected Outputs:
        Screen content shows the merged string at the correct cursor position.
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
        Alphabet string followed by moving cursor to index 5 and performing 4 backspaces.
    Expected Outputs:
        The characters preceding index 5 are removed and the string is collapsed.
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
        Alphabet string followed by moving cursor to index 5 and performing 4 deletes.
    Expected Outputs:
        The characters at and after index 5 are removed as expected.
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
        Alphabet string with cursor moved to the final character followed by multiple delete commands.
    Expected Outputs:
        Delete commands at the end of the string have no visual or data effect.
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
        A text field with a 10-character limit where the full 26-character alphabet is typed.
    Expected Outputs:
        The field only contains the first 10 characters "abcdefghij".
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
        Cursor moved to index 4 followed by 10 backspace commands.
    Expected Outputs:
        The cursor remains at index 0 and no data corruption occurs.
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

/*
TestTextFieldAsciiCharacterization is a test which pins the pre-existing ASCII editing behaviour of the text field
so the wide-rune reconciliation work cannot silently regress plain-text typing, cursor movement, or viewport
scrolling.

Example:
    Expected Inputs:
        A width 6 text field, max length 20, into which "abcdef" is typed, then "home", "right" x3, "end",
        "backspace".
    Expected Outputs:
        After "abcdef": CurrentValue "abcdef " (runes a b c d e f space), CursorPosition 6, ViewportPosition 1.
        After "home": CursorPosition 0, ViewportPosition 0.
        After "right" x3: CursorPosition 3, ViewportPosition 0.
        After "end": CursorPosition 6, ViewportPosition 1.
        After "backspace": CurrentValue "abcde ", CursorPosition 5, ViewportPosition 1, and the rendered screen
        matches the committed master image.
*/
func TestTextFieldAsciiCharacterization(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	textFieldInstance := layer1.AddTextField(styleEntry, 2, 2, 6, 20, false, "", true)
	setFocusedControl(layer1.layerAlias, textFieldInstance.controlAlias, constants.CellTypeTextField)
	textFieldEntry := TextFields.Get(textFieldInstance.layerAlias, textFieldInstance.controlAlias)

	TextField.updateKeyboardEventTextboxWithString("abcdef")
	assert.Equal(test, "abcdef ", string(textFieldEntry.CurrentValue))
	assert.Equal(test, 6, textFieldEntry.CursorPosition)
	assert.Equal(test, 1, textFieldEntry.ViewportPosition)

	TextField.updateKeyboardEvent([]rune("home"))
	assert.Equal(test, 0, textFieldEntry.CursorPosition)
	assert.Equal(test, 0, textFieldEntry.ViewportPosition)

	TextField.updateKeyboardEventTextboxWithCommands("right", "right", "right")
	assert.Equal(test, 3, textFieldEntry.CursorPosition)
	assert.Equal(test, 0, textFieldEntry.ViewportPosition)

	TextField.updateKeyboardEvent([]rune("end"))
	assert.Equal(test, 6, textFieldEntry.CursorPosition)
	assert.Equal(test, 1, textFieldEntry.ViewportPosition)

	TextField.updateKeyboardEvent([]rune("backspace"))
	assert.Equal(test, "abcde ", string(textFieldEntry.CurrentValue))
	assert.Equal(test, 5, textFieldEntry.CursorPosition)
	assert.Equal(test, 1, textFieldEntry.ViewportPosition)

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TEXT_FIELD_TEST_SUITE_NAME, "TestTextFieldAsciiCharacterization", obtainedValue)
	expectedValue := LoadMasterImage(TEXT_FIELD_TEST_SUITE_NAME, "TestTextFieldAsciiCharacterization")
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", layerEntry.GetAnsiStringFromBase64(expectedValue))
		fmt.Println("Obtained:\n", layerEntry.GetAnsiStringFromBase64(obtainedValue))
	}
}

/*
TestTextFieldCjkInsert is a test which verifies that typing wide CJK runes into a text field stores one rune per
ideograph and leaves the cursor on a rune boundary.

Example:
    Expected Inputs:
        A width 6 text field into which the string "中文" is typed.
    Expected Outputs:
        CurrentValue is the rune slice [中 文 space], CursorPosition 2, ViewportPosition 0, and the rendered
        screen matches the committed master image (the two ideographs occupy four columns plus the cursor cell).
*/
func TestTextFieldCjkInsert(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	textFieldInstance := layer1.AddTextField(styleEntry, 2, 2, 6, 20, false, "", true)
	setFocusedControl(layer1.layerAlias, textFieldInstance.controlAlias, constants.CellTypeTextField)
	textFieldEntry := TextFields.Get(textFieldInstance.layerAlias, textFieldInstance.controlAlias)

	TextField.updateKeyboardEventTextboxWithString("中文")
	assert.Equal(test, []rune("中文 "), textFieldEntry.CurrentValue)
	assert.Equal(test, 2, textFieldEntry.CursorPosition)
	assert.Equal(test, 0, textFieldEntry.ViewportPosition)

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TEXT_FIELD_TEST_SUITE_NAME, "TestTextFieldCjkInsert", obtainedValue)
	expectedValue := LoadMasterImage(TEXT_FIELD_TEST_SUITE_NAME, "TestTextFieldCjkInsert")
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", layerEntry.GetAnsiStringFromBase64(expectedValue))
		fmt.Println("Obtained:\n", layerEntry.GetAnsiStringFromBase64(obtainedValue))
	}
}

/*
TestTextFieldCjkBackspace is a test which verifies that backspacing across a wide CJK rune removes exactly one
rune and keeps the trailing sentinel blank.

Example:
    Expected Inputs:
        A width 6 text field where "中文" is typed and then one "backspace" is issued.
    Expected Outputs:
        CurrentValue is the rune slice [中 space], CursorPosition 1, ViewportPosition 0.
*/
func TestTextFieldCjkBackspace(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	textFieldInstance := layer1.AddTextField(styleEntry, 2, 2, 6, 20, false, "", true)
	setFocusedControl(layer1.layerAlias, textFieldInstance.controlAlias, constants.CellTypeTextField)
	textFieldEntry := TextFields.Get(textFieldInstance.layerAlias, textFieldInstance.controlAlias)

	TextField.updateKeyboardEventTextboxWithString("中文")
	TextField.updateKeyboardEvent([]rune("backspace"))
	assert.Equal(test, []rune("中 "), textFieldEntry.CurrentValue)
	assert.Equal(test, 1, textFieldEntry.CursorPosition)
	assert.Equal(test, 0, textFieldEntry.ViewportPosition)
	assert.Equal(test, ' ', textFieldEntry.CurrentValue[len(textFieldEntry.CurrentValue)-1])
}

/*
TestTextFieldCjkDeleteMidString is a test which verifies that a forward delete in the middle of a string mixing
wide and narrow runes removes only the targeted rune.

Example:
    Expected Inputs:
        A width 6 text field where "中A文" is typed, then "home", "right", "delete".
    Expected Outputs:
        The narrow 'A' at rune index 1 is removed, leaving CurrentValue [中 文 space] with CursorPosition 1.
*/
func TestTextFieldCjkDeleteMidString(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	textFieldInstance := layer1.AddTextField(styleEntry, 2, 2, 6, 20, false, "", true)
	setFocusedControl(layer1.layerAlias, textFieldInstance.controlAlias, constants.CellTypeTextField)
	textFieldEntry := TextFields.Get(textFieldInstance.layerAlias, textFieldInstance.controlAlias)

	TextField.updateKeyboardEventTextboxWithString("中A文")
	TextField.updateKeyboardEvent([]rune("home"))
	TextField.updateKeyboardEvent([]rune("right"))
	TextField.updateKeyboardEvent([]rune("delete"))
	assert.Equal(test, []rune("中文 "), textFieldEntry.CurrentValue)
	assert.Equal(test, 1, textFieldEntry.CursorPosition)
}

/*
TestTextFieldCjkViewportScrollRight is a test which verifies that typing past the visible width with wide runes
scrolls the viewport in columns so the caret stays visible and no wide rune is split at either edge.

Example:
    Expected Inputs:
        A width 6 text field into which "中文中文中" (ten printed columns) is typed.
    Expected Outputs:
        CurrentValue [中 文 中 文 中 space], CursorPosition 5, ViewportPosition 3. The visible prefix printed
        width is at most 6 columns and the caret column relative to the viewport is strictly less than 6.
*/
func TestTextFieldCjkViewportScrollRight(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	textFieldInstance := layer1.AddTextField(styleEntry, 2, 2, 6, 20, false, "", true)
	setFocusedControl(layer1.layerAlias, textFieldInstance.controlAlias, constants.CellTypeTextField)
	textFieldEntry := TextFields.Get(textFieldInstance.layerAlias, textFieldInstance.controlAlias)

	TextField.updateKeyboardEventTextboxWithString("中文中文中")
	assert.Equal(test, []rune("中文中文中 "), textFieldEntry.CurrentValue)
	assert.Equal(test, 5, textFieldEntry.CursorPosition)
	assert.Equal(test, 3, textFieldEntry.ViewportPosition)

	visibleRunes := textFieldEntry.CurrentValue[textFieldEntry.ViewportPosition:]
	visiblePrefix, _ := stringformat.GetRunesThatFitInColumnCountFromStart(visibleRunes, textFieldEntry.Width)
	assert.LessOrEqual(test, stringformat.GetWidthOfRunesWhenPrinted(visiblePrefix), textFieldEntry.Width)
	caretColumn := stringformat.GetColumnIndexBasedOnRuneIndex(visibleRunes, textFieldEntry.CursorPosition-textFieldEntry.ViewportPosition)
	assert.Less(test, caretColumn, textFieldEntry.Width)
}

/*
TestTextFieldCjkViewportScrollLeftHome is a test which verifies that pressing Home from a horizontally scrolled
state snaps the viewport back to the first rune and renders the first wide rune at the left edge.

Example:
    Expected Inputs:
        A width 6 text field scrolled by typing "中文中文中", then a "home" keystroke.
    Expected Outputs:
        CursorPosition 0, ViewportPosition 0. On screen the first ideograph occupies columns 0 and 1 (the lead
        cell holds 中 and the placeholder cell holds a blank).
*/
func TestTextFieldCjkViewportScrollLeftHome(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	textFieldInstance := layer1.AddTextField(styleEntry, 2, 2, 6, 20, false, "", true)
	setFocusedControl(layer1.layerAlias, textFieldInstance.controlAlias, constants.CellTypeTextField)
	textFieldEntry := TextFields.Get(textFieldInstance.layerAlias, textFieldInstance.controlAlias)

	TextField.updateKeyboardEventTextboxWithString("中文中文中")
	TextField.updateKeyboardEvent([]rune("home"))
	assert.Equal(test, 0, textFieldEntry.CursorPosition)
	assert.Equal(test, 0, textFieldEntry.ViewportPosition)

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	assert.Equal(test, '中', layerEntry.CharacterMemory[2][2].Character)
	assert.Equal(test, ' ', layerEntry.CharacterMemory[2][3].Character)
}

/*
TestTextFieldCjkClickRightHalf is a test which verifies that both cells of a wide rune carry the lead rune's
CellControlId so a click on either half resolves to the same cursor position.

Example:
    Expected Inputs:
        A width 6 focused text field with default value "中文" rendered to the screen; a simulated click on the
        trailing-half column of each ideograph.
    Expected Outputs:
        Screen cells for 中 (columns 2 and 3) both report CellControlId 0; cells for 文 (columns 4 and 5) both
        report CellControlId 1. A click at column 3 sets CursorPosition 0; a click at column 5 sets
        CursorPosition 1.
*/
func TestTextFieldCjkClickRightHalf(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	textFieldInstance := layer1.AddTextField(styleEntry, 2, 2, 6, 20, false, "中文", true)
	setFocusedControl(layer1.layerAlias, textFieldInstance.controlAlias, constants.CellTypeTextField)
	textFieldEntry := TextFields.Get(textFieldInstance.layerAlias, textFieldInstance.controlAlias)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer

	assert.Equal(test, 0, layerEntry.CharacterMemory[2][2].AttributeEntry.CellControlId)
	assert.Equal(test, 0, layerEntry.CharacterMemory[2][3].AttributeEntry.CellControlId)
	assert.Equal(test, 1, layerEntry.CharacterMemory[2][4].AttributeEntry.CellControlId)
	assert.Equal(test, 1, layerEntry.CharacterMemory[2][5].AttributeEntry.CellControlId)

	eventStateMemory.stateId = 0
	SetMouseStatus(3, 2, 1, "")
	TextField.updateMouseEvent()
	assert.Equal(test, 0, textFieldEntry.CursorPosition)

	eventStateMemory.stateId = 0
	SetMouseStatus(5, 2, 1, "")
	TextField.updateMouseEvent()
	assert.Equal(test, 1, textFieldEntry.CursorPosition)
}

/*
TestTextFieldCjkHighlightDrag is a test which verifies that dragging a selection across wide runes records
highlight bounds as rune indices.

Example:
    Expected Inputs:
        A width 6 focused text field with default value "中文"; a click on the first ideograph followed by a
        drag whose mouse position lands on the trailing half of the second ideograph.
    Expected Outputs:
        HighlightStart 0, HighlightEnd 1, IsHighlightActive true.
*/
func TestTextFieldCjkHighlightDrag(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	textFieldInstance := layer1.AddTextField(styleEntry, 2, 2, 6, 20, false, "中文", true)
	setFocusedControl(layer1.layerAlias, textFieldInstance.controlAlias, constants.CellTypeTextField)
	textFieldEntry := TextFields.Get(textFieldInstance.layerAlias, textFieldInstance.controlAlias)
	UpdateDisplay(false)

	eventStateMemory.stateId = 0
	SetMouseStatus(2, 2, 1, "")
	TextField.updateMouseEvent()

	eventStateMemory.stateId = constants.EventStateDragAndDrop
	SetMouseStatus(5, 2, 1, "")
	TextField.updateMouseEvent()
	eventStateMemory.stateId = 0

	assert.True(test, textFieldEntry.IsHighlightActive)
	assert.Equal(test, 0, textFieldEntry.HighlightStart)
	assert.Equal(test, 1, textFieldEntry.HighlightEnd)
}

/*
TestTextFieldCjkPasswordAlignment is a test which verifies that a password field masks each wide source rune
with a single asterisk plus a blank pad cell so the masked field keeps the same column layout as the unmasked
one.

Example:
    Expected Inputs:
        A width 6 password protected text field with default value "中文" rendered to the screen.
    Expected Outputs:
        Screen columns 2 and 4 hold '*'; columns 3 and 5 hold a blank pad; the rendered screen matches the
        committed master image.
*/
func TestTextFieldCjkPasswordAlignment(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	textFieldInstance := layer1.AddTextField(styleEntry, 2, 2, 6, 20, true, "中文", true)
	setFocusedControl(layer1.layerAlias, textFieldInstance.controlAlias, constants.CellTypeTextField)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer

	assert.Equal(test, '*', layerEntry.CharacterMemory[2][2].Character)
	assert.Equal(test, ' ', layerEntry.CharacterMemory[2][3].Character)
	assert.Equal(test, '*', layerEntry.CharacterMemory[2][4].Character)
	assert.Equal(test, ' ', layerEntry.CharacterMemory[2][5].Character)

	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TEXT_FIELD_TEST_SUITE_NAME, "TestTextFieldCjkPasswordAlignment", obtainedValue)
	expectedValue := LoadMasterImage(TEXT_FIELD_TEST_SUITE_NAME, "TestTextFieldCjkPasswordAlignment")
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", layerEntry.GetAnsiStringFromBase64(expectedValue))
		fmt.Println("Obtained:\n", layerEntry.GetAnsiStringFromBase64(obtainedValue))
	}
}

/*
TestTextFieldMaxLengthIsRunes is a test which pins the contract that the text field maximum length counts RUNES,
not printed columns, so a limit of three admits three wide ideographs even though they print as six columns.

Example:
    Expected Inputs:
        A width 10 text field with max length 3 into which "中文字字" is typed.
    Expected Outputs:
        Only the first three ideographs are stored: CurrentValue is the rune slice [中 文 字 space] of length 4.
*/
func TestTextFieldMaxLengthIsRunes(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	textFieldInstance := layer1.AddTextField(styleEntry, 2, 2, 10, 3, false, "", true)
	setFocusedControl(layer1.layerAlias, textFieldInstance.controlAlias, constants.CellTypeTextField)
	textFieldEntry := TextFields.Get(textFieldInstance.layerAlias, textFieldInstance.controlAlias)

	TextField.updateKeyboardEventTextboxWithString("中文字字")
	assert.Equal(test, []rune("中文字 "), textFieldEntry.CurrentValue)
	assert.Equal(test, 4, len(textFieldEntry.CurrentValue))
}

/*
TestTextFieldSentinelPreserved is a test which verifies that the trailing sentinel blank in CurrentValue survives
an arbitrary mix of insert, delete, backspace, and navigation edits over wide and narrow runes.

Example:
    Expected Inputs:
        A width 6 text field driven through: type "中a文b", home, "delete", end, "backspace", type "字".
    Expected Outputs:
        The last rune of CurrentValue is always a blank space after every edit.
*/
func TestTextFieldSentinelPreserved(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	textFieldInstance := layer1.AddTextField(styleEntry, 2, 2, 6, 20, false, "", true)
	setFocusedControl(layer1.layerAlias, textFieldInstance.controlAlias, constants.CellTypeTextField)
	textFieldEntry := TextFields.Get(textFieldInstance.layerAlias, textFieldInstance.controlAlias)
	lastRuneIsBlank := func() bool {
		return textFieldEntry.CurrentValue[len(textFieldEntry.CurrentValue)-1] == ' '
	}

	TextField.updateKeyboardEventTextboxWithString("中a文b")
	assert.True(test, lastRuneIsBlank())
	TextField.updateKeyboardEvent([]rune("home"))
	TextField.updateKeyboardEvent([]rune("delete"))
	assert.True(test, lastRuneIsBlank())
	TextField.updateKeyboardEvent([]rune("end"))
	TextField.updateKeyboardEvent([]rune("backspace"))
	assert.True(test, lastRuneIsBlank())
	TextField.updateKeyboardEventTextboxWithString("字")
	assert.True(test, lastRuneIsBlank())
}
