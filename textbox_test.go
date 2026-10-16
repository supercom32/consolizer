package consolizer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/stringformat"
	"github.com/supercom32/consolizer/types"
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
	layer1, _, _, styleEntry := CommonTestSetup(test)
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
	layer1, _, _, styleEntry := CommonTestSetup(test)
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
	layer1, _, _, styleEntry := CommonTestSetup(test)
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
	layer1, _, _, styleEntry := CommonTestSetup(test)
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
	layer1, _, _, styleEntry := CommonTestSetup(test)
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
	layer1, _, _, styleEntry := CommonTestSetup(test)
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
	layer1, _, _, styleEntry := CommonTestSetup(test)
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

/*
TestTextboxAsciiCharacterization is a test which pins the pre-existing ASCII behaviour of the textbox (typing,
horizontal viewport scroll, line split, end key) so the wide-rune reconciliation work cannot silently regress
plain-text editing.

Example:
    Expected Inputs:
        A width 10, height 3 textbox into which "abcdefghijkl" is typed, then "home", "enter", "end".
    Expected Outputs:
        After "abcdefghijkl": single line [a..l space], CursorXLocation 12, CursorYLocation 0,
        ViewportXLocation 3, ViewportYLocation 0.
        After "home": CursorXLocation 0, ViewportXLocation 0.
        After "enter": two lines, CursorYLocation 1, CursorXLocation 0.
        After "end": CursorXLocation 12, ViewportXLocation 3, and the rendered screen matches the master image.
*/
func TestTextboxAsciiCharacterization(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	textboxInstance := layer1.AddTextbox(styleEntry, 2, 2, 10, 3, false)
	setFocusedControl(layer1.layerAlias, textboxInstance.controlAlias, constants.CellTypeTextbox)
	textboxEntry := Textboxes.Get(textboxInstance.layerAlias, textboxInstance.controlAlias)

	textbox.UpdateKeyboardEventTextboxWithString("abcdefghijkl")
	assert.Equal(test, "abcdefghijkl ", string(textboxEntry.TextData[0]))
	assert.Equal(test, 12, textboxEntry.CursorXLocation)
	assert.Equal(test, 0, textboxEntry.CursorYLocation)
	assert.Equal(test, 3, textboxEntry.ViewportXLocation)
	assert.Equal(test, 0, textboxEntry.ViewportYLocation)

	textbox.UpdateKeyboardEvent([]rune("home"))
	assert.Equal(test, 0, textboxEntry.CursorXLocation)
	assert.Equal(test, 0, textboxEntry.ViewportXLocation)

	textbox.UpdateKeyboardEvent([]rune("enter"))
	assert.Equal(test, 2, len(textboxEntry.TextData))
	assert.Equal(test, 1, textboxEntry.CursorYLocation)
	assert.Equal(test, 0, textboxEntry.CursorXLocation)

	textbox.UpdateKeyboardEvent([]rune("end"))
	assert.Equal(test, 12, textboxEntry.CursorXLocation)
	assert.Equal(test, 3, textboxEntry.ViewportXLocation)

	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, TEXTBOX_TEST_SUITE_NAME, "TestTextboxAsciiCharacterization", obtainedValue)
	expectedValue := LoadMasterImage(TEXTBOX_TEST_SUITE_NAME, "TestTextboxAsciiCharacterization")
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", layerEntry.GetAnsiStringFromBase64(expectedValue))
		fmt.Println("Obtained:\n", layerEntry.GetAnsiStringFromBase64(obtainedValue))
	}
}

/*
TestTextboxCjkInsertDeleteBackspace is a test which verifies that inserting, backspacing, and forward deleting
wide CJK runes in a multiline textbox operates one rune at a time and never leaves a half rune behind.

Example:
    Expected Inputs:
        A width 12 textbox: type "中文", "enter", "字符", then "backspace", then "home" "delete".
    Expected Outputs:
        After typing: TextData [[中 文 space] [字 符 space]].
        After "backspace": the trailing 符 on line 2 is removed, leaving [字 space] with CursorXLocation 1,
        CursorYLocation 1.
        After "home" then "delete": the leading 字 on line 2 is removed, leaving [中 文 space] and [space].
*/
func TestTextboxCjkInsertDeleteBackspace(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	textboxInstance := layer1.AddTextbox(styleEntry, 2, 2, 12, 4, false)
	setFocusedControl(layer1.layerAlias, textboxInstance.controlAlias, constants.CellTypeTextbox)
	textboxEntry := Textboxes.Get(textboxInstance.layerAlias, textboxInstance.controlAlias)

	textbox.UpdateKeyboardEventTextboxWithString("中文")
	textbox.UpdateKeyboardEvent([]rune("enter"))
	textbox.UpdateKeyboardEventTextboxWithString("字符")
	assert.Equal(test, [][]rune{[]rune("中文 "), []rune("字符 ")}, textboxEntry.TextData)

	textbox.UpdateKeyboardEvent([]rune("backspace"))
	assert.Equal(test, [][]rune{[]rune("中文 "), []rune("字 ")}, textboxEntry.TextData)
	assert.Equal(test, 1, textboxEntry.CursorXLocation)
	assert.Equal(test, 1, textboxEntry.CursorYLocation)

	textbox.UpdateKeyboardEvent([]rune("home"))
	textbox.UpdateKeyboardEvent([]rune("delete"))
	assert.Equal(test, [][]rune{[]rune("中文 "), []rune(" ")}, textboxEntry.TextData)
}

/*
TestTextboxCjkHorizontalWindow is a test which verifies that with word wrap off ViewportXLocation is a printed
COLUMN offset: the drawn window is Width columns wide, always begins on a rune boundary, never splits a wide
rune, and a column offset that lands on the trailing half of a wide rune advances to the next whole rune rather
than clipping it. The horizontal scrollbar ceiling is the exact printed width of the widest line.

Example:

	Expected Inputs:
	    A width 6 textbox whose single line is "中文中文中文" (six wide runes, twelve printed columns), rendered
	    with ViewportXLocation set to 0, 1, 2, 3, 4, then 6.
	Expected Outputs:
	    getMaxHorizontalTextValue is 12. The six drawn cells are, per column offset:
	        0 -> 中,blank,文,blank,中,blank
	        1 -> 文,blank,中,blank,文,blank   (offset 1 is the trailing half of 中, rounds up to 文)
	        2 -> 文,blank,中,blank,文,blank
	        3 -> 中,blank,文,blank,中,blank   (offset 3 is the trailing half of 文, rounds up to 中)
	        4 -> 中,blank,文,blank,中,blank
	        6 -> 文,blank,中,blank,文,blank
	    The last drawn cell is never a wide-rune lead.
*/
func TestTextboxCjkHorizontalWindow(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	textboxInstance := layer1.AddTextbox(styleEntry, 2, 2, 6, 4, false)
	setFocusedControl(layer1.layerAlias, textboxInstance.controlAlias, constants.CellTypeTextbox)
	textboxEntry := Textboxes.Get(textboxInstance.layerAlias, textboxInstance.controlAlias)
	textboxEntry.TextData = [][]rune{[]rune("中文中文中文")}
	textboxEntry.CursorXLocation = 0
	textboxEntry.CursorYLocation = 0

	assert.Equal(test, 12, textbox.getMaxHorizontalTextValue(textboxInstance.layerAlias, textboxInstance.controlAlias))

	expectedWindows := map[int][]rune{
		0: []rune("中 文 中 "),
		1: []rune("文 中 文 "),
		2: []rune("文 中 文 "),
		3: []rune("中 文 中 "),
		4: []rune("中 文 中 "),
		6: []rune("文 中 文 "),
	}
	for _, viewportXLocation := range []int{0, 1, 2, 3, 4, 6} {
		textboxEntry.ViewportXLocation = viewportXLocation
		UpdateDisplay(false)
		layerEntry := commonResource.screenLayer
		window := make([]rune, 0, textboxEntry.Width)
		for column := 0; column < textboxEntry.Width; column++ {
			window = append(window, layerEntry.CharacterMemory[2][2+column].Character)
		}
		assert.Equalf(test, expectedWindows[viewportXLocation], window, "window mismatch at ViewportXLocation %d", viewportXLocation)
		assert.Falsef(test, stringformat.IsRuneCharacterWide(layerEntry.CharacterMemory[2][2+textboxEntry.Width-1].Character), "wide rune split at right edge for ViewportXLocation %d", viewportXLocation)
	}
}

/*
TestTextboxCjkWordWrap is a test which verifies that with word wrap on every produced line fits the column budget
and no wrapped line ends on a half rune.

Example:
    Expected Inputs:
        A width 10 textbox with word wrap enabled whose single line is "中文 中文中文 中文中文中文" (spaces between
        groups), which is 22 printed columns.
    Expected Outputs:
        wrapTextToWidth returns more than one line, each of printed width at most 10, and concatenating every
        wrapped line's non-space runes preserves the original non-space runes in order.
*/
func TestTextboxCjkWordWrap(test *testing.T) {
	_, _, _, _ = CommonTestSetup(test)
	source := [][]rune{[]rune("中文 中文中文 中文中文中文")}
	wrapped := textbox.wrapTextToWidth(source, 10)
	assert.Greater(test, len(wrapped), 1)
	for _, line := range wrapped {
		assert.LessOrEqual(test, stringformat.GetWidthOfRunesWhenPrinted(line), 10)
	}
	var rebuilt []rune
	for _, line := range wrapped {
		for _, currentRune := range line {
			if currentRune != ' ' {
				rebuilt = append(rebuilt, currentRune)
			}
		}
	}
	var original []rune
	for _, currentRune := range source[0] {
		if currentRune != ' ' {
			original = append(original, currentRune)
		}
	}
	assert.Equal(test, original, rebuilt)
}

/*
TestTextboxCjkClickCoordinates is a test which verifies that a mouse click on the trailing half of a wide rune
resolves to that rune's index and that a click past the end of a line clamps to the last rune index.

Example:
    Expected Inputs:
        A width 10 textbox at layer position (2,2) whose single line is the rune slice for "中文字 "; a click at
        screen column 3 (the trailing half of the leading 中) and a click at screen column 10 (past the text).
    Expected Outputs:
        Screen cells for 中 (columns 2 and 3) both report CellControlId 0; cells for 文 (columns 4 and 5) both
        report CellControlId 1. The column 3 click sets CursorXLocation 0, CursorYLocation 0. The column 10
        click clamps CursorXLocation to 3 (len(line) minus one).
*/
func TestTextboxCjkClickCoordinates(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	textboxInstance := layer1.AddTextbox(styleEntry, 2, 2, 10, 4, false)
	setFocusedControl(layer1.layerAlias, textboxInstance.controlAlias, constants.CellTypeTextbox)
	textboxEntry := Textboxes.Get(textboxInstance.layerAlias, textboxInstance.controlAlias)
	textboxEntry.TextData = [][]rune{[]rune("中文字 ")}
	textboxEntry.CursorXLocation = 0
	textboxEntry.CursorYLocation = 0
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer

	assert.Equal(test, 0, layerEntry.CharacterMemory[2][2].AttributeEntry.CellControlId)
	assert.Equal(test, 0, layerEntry.CharacterMemory[2][3].AttributeEntry.CellControlId)
	assert.Equal(test, 1, layerEntry.CharacterMemory[2][4].AttributeEntry.CellControlId)
	assert.Equal(test, 1, layerEntry.CharacterMemory[2][5].AttributeEntry.CellControlId)

	setEventStateId(0)
	SetMouseStatus(3, 2, 1, "")
	textbox.updateMouseEvent()
	assert.Equal(test, 0, textboxEntry.CursorXLocation)
	assert.Equal(test, 0, textboxEntry.CursorYLocation)

	setEventStateId(0)
	SetMouseStatus(10, 2, 1, "")
	textbox.updateMouseEvent()
	assert.Equal(test, 3, textboxEntry.CursorXLocation)
}

/*
TestTextboxCjkVerticalScrollUnaffected is a test which verifies that wide runes on a line do not change the
vertical scrollbar ceiling, which depends only on the number of lines.

Example:
    Expected Inputs:
        A width 6, height 4 textbox with eight lines, each the rune slice for "中文中".
    Expected Outputs:
        The vertical scrollbar MaxScrollValue equals len(TextData) minus Height, which is 4.
*/
func TestTextboxCjkVerticalScrollUnaffected(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	textboxInstance := layer1.AddTextbox(styleEntry, 2, 2, 6, 4, false)
	setFocusedControl(layer1.layerAlias, textboxInstance.controlAlias, constants.CellTypeTextbox)
	textboxEntry := Textboxes.Get(textboxInstance.layerAlias, textboxInstance.controlAlias)
	textboxEntry.TextData = nil
	for lineIndex := 0; lineIndex < 8; lineIndex++ {
		textboxEntry.TextData = append(textboxEntry.TextData, []rune("中文中"))
	}
	textbox.setTextboxMaxScrollBarValues(textboxInstance.layerAlias, textboxInstance.controlAlias)

	verticalScrollbarEntry := ScrollBars.Get(textboxInstance.layerAlias, textboxEntry.VerticalScrollbarAlias)
	assert.Equal(test, len(textboxEntry.TextData)-textboxEntry.Height, verticalScrollbarEntry.MaxScrollValue)
	assert.Equal(test, 4, verticalScrollbarEntry.MaxScrollValue)
}

/*
TestTextboxCjkCursorRendering is a test which pins the policy that a cursor sitting on a wide rune paints BOTH of
that rune's columns with the cursor colours.

Example:
    Expected Inputs:
        A width 10 focused textbox whose single line is the rune slice for "中文字 " with CursorXLocation 1
        (on the second ideograph).
    Expected Outputs:
        The rendered lead cell (column 4) and its placeholder cell (column 5) both carry the style's textbox
        cursor foreground and background colours.
*/
func TestTextboxCjkCursorRendering(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	textboxInstance := layer1.AddTextbox(styleEntry, 2, 2, 10, 4, false)
	setFocusedControl(layer1.layerAlias, textboxInstance.controlAlias, constants.CellTypeTextbox)
	textboxEntry := Textboxes.Get(textboxInstance.layerAlias, textboxInstance.controlAlias)
	textboxEntry.TextData = [][]rune{[]rune("中文字 ")}
	textboxEntry.CursorXLocation = 1
	textboxEntry.CursorYLocation = 0
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer

	leadCell := layerEntry.CharacterMemory[2][4]
	placeholderCell := layerEntry.CharacterMemory[2][5]
	assert.Equal(test, styleEntry.Textbox.CursorForegroundColor, leadCell.AttributeEntry.ForegroundColor)
	assert.Equal(test, styleEntry.Textbox.CursorBackgroundColor, leadCell.AttributeEntry.BackgroundColor)
	assert.Equal(test, styleEntry.Textbox.CursorForegroundColor, placeholderCell.AttributeEntry.ForegroundColor)
	assert.Equal(test, styleEntry.Textbox.CursorBackgroundColor, placeholderCell.AttributeEntry.BackgroundColor)
}

/*
TestTextboxGetWrappedRows is a test which pins the wrapped-row layout the textbox draws from, including the
logical line index and start rune index each display row maps back to.

Example:

	Expected Inputs:
	    Two logical lines, "abcdefghij klmnop qrstuv wxyz 0123456789" and "short", wrapped at width 10.
	Expected Outputs:
	    Five display rows: {line 0, start 0, "abcdefghij"}, {line 0, start 11, "klmnop qrs"},
	    {line 0, start 21, "tuv wxyz "}, {line 0, start 30, "0123456789"}, {line 1, start 0, "short"}.
*/
func TestTextboxGetWrappedRows(test *testing.T) {
	_, _, _, _ = CommonTestSetup(test)
	textData := [][]rune{[]rune("abcdefghij klmnop qrstuv wxyz 0123456789"), []rune("short")}
	wrappedRows := textbox.getWrappedRows(textData, 10, true)

	type expectedRow struct {
		logicalLineIndex int
		startRuneIndex   int
		text             string
	}
	expected := []expectedRow{
		{0, 0, "abcdefghij"},
		{0, 11, "klmnop qrs"},
		{0, 21, "tuv wxyz "},
		{0, 30, "0123456789"},
		{1, 0, "short"},
	}
	assert.Equal(test, len(expected), len(wrappedRows))
	for currentIndex, wantRow := range expected {
		assert.Equalf(test, wantRow.logicalLineIndex, wrappedRows[currentIndex].logicalLineIndex, "row %d logical line", currentIndex)
		assert.Equalf(test, wantRow.startRuneIndex, wrappedRows[currentIndex].startRuneIndex, "row %d start rune", currentIndex)
		assert.Equalf(test, wantRow.text, string(wrappedRows[currentIndex].runes), "row %d text", currentIndex)
	}

	// Round trip: every logical position maps to a display position and back to itself.
	line := textData[0]
	for runeIndex := 0; runeIndex <= len(line); runeIndex++ {
		displayRow, displayColumn := textbox.getCursorDisplayCoordinates(wrappedRows, 0, runeIndex)
		gotLine, gotRune := textbox.getLogicalCoordinatesFromDisplay(wrappedRows, displayRow, displayColumn)
		if runeIndex < len(line) && line[runeIndex] != ' ' {
			assert.Equalf(test, 0, gotLine, "round trip line for rune %d", runeIndex)
			assert.Equalf(test, runeIndex, gotRune, "round trip rune for rune %d", runeIndex)
		}
	}
}

/*
TestTextboxWordWrapClickPlacesCursor is a test which verifies that clicking a cell on a wrapped display row moves
the cursor to the matching LOGICAL line and rune index, not to a wrapped-segment offset.

Example:

	Expected Inputs:
	    A width 10, height 4 textbox at layer position (2,2) with word wrap on and the single logical line
	    "abcdefghij klmnop qrstuv wxyz 0123456789"; a click at screen column 4, row 4 (the third character of the
	    third wrapped row, which begins at logical rune 21).
	Expected Outputs:
	    CursorYLocation 0 and CursorXLocation 23, the logical index of 'v'.
*/
func TestTextboxWordWrapClickPlacesCursor(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	textboxInstance := layer1.AddTextbox(styleEntry, 2, 2, 10, 4, false)
	setFocusedControl(layer1.layerAlias, textboxInstance.controlAlias, constants.CellTypeTextbox)
	textboxEntry := Textboxes.Get(textboxInstance.layerAlias, textboxInstance.controlAlias)
	textboxInstance.SetWordWrap(true)
	textboxEntry.TextData = [][]rune{[]rune("abcdefghij klmnop qrstuv wxyz 0123456789")}
	textboxEntry.CursorXLocation = 0
	textboxEntry.CursorYLocation = 0
	textboxEntry.ViewportXLocation = 0
	textboxEntry.ViewportYLocation = 0
	UpdateDisplay(false)

	assert.Equal(test, 'v', commonResource.screenLayer.CharacterMemory[4][4].Character)

	setEventStateId(0)
	SetMouseStatus(4, 4, 1, "")
	textbox.updateMouseEvent()
	assert.Equal(test, 0, textboxEntry.CursorYLocation)
	assert.Equal(test, 23, textboxEntry.CursorXLocation)
}

/*
TestTextboxWordWrapCursorStaysVisible is a test which verifies that after the cursor moves to a wrapped display
row below the visible window the viewport scrolls so the cursor row is inside it and the cursor cell is actually
painted.

Example:

	Expected Inputs:
	    A width 10, height 2 textbox with word wrap on and one logical line that wraps to four display rows; the
	    cursor is placed at the end of the logical line and updateViewport is run.
	Expected Outputs:
	    ViewportYLocation is 2 so the cursor's display row 3 sits in the window [2,4); a cell carrying the textbox
	    cursor background colour is present in the drawn area.
*/
func TestTextboxWordWrapCursorStaysVisible(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	textboxInstance := layer1.AddTextbox(styleEntry, 2, 2, 10, 2, false)
	setFocusedControl(layer1.layerAlias, textboxInstance.controlAlias, constants.CellTypeTextbox)
	textboxEntry := Textboxes.Get(textboxInstance.layerAlias, textboxInstance.controlAlias)
	textboxInstance.SetWordWrap(true)
	textboxEntry.TextData = [][]rune{[]rune("abcdefghij klmnop qrstuv wxyz 0123456789")}
	textboxEntry.CursorYLocation = 0
	textboxEntry.CursorXLocation = len(textboxEntry.TextData[0]) - 1
	textbox.updateViewport(textboxEntry)

	wrappedRows := textbox.getWrappedRows(textboxEntry.TextData, textboxEntry.Width, true)
	cursorDisplayRow, _ := textbox.getCursorDisplayCoordinates(wrappedRows, textboxEntry.CursorYLocation, textboxEntry.CursorXLocation)
	assert.Equal(test, 2, textboxEntry.ViewportYLocation)
	assert.GreaterOrEqual(test, cursorDisplayRow, textboxEntry.ViewportYLocation)
	assert.Less(test, cursorDisplayRow, textboxEntry.ViewportYLocation+textboxEntry.Height)

	UpdateDisplay(false)
	cursorCellFound := false
	for rowIndex := 2; rowIndex < 4; rowIndex++ {
		for columnIndex := 2; columnIndex < 12; columnIndex++ {
			if commonResource.screenLayer.CharacterMemory[rowIndex][columnIndex].AttributeEntry.BackgroundColor == styleEntry.Textbox.CursorBackgroundColor {
				cursorCellFound = true
			}
		}
	}
	assert.True(test, cursorCellFound)
}

/*
TestTextboxWordWrapScrollbarUnits is a test which verifies that with word wrap on the vertical scrollbar ceiling
counts wrapped display rows, not logical lines, and the horizontal scrollbar is disabled.

Example:

	Expected Inputs:
	    A width 10, height 2 textbox with word wrap on whose single logical line wraps to four display rows.
	Expected Outputs:
	    The vertical scrollbar MaxScrollValue is 2 (four rows minus height two); the horizontal scrollbar is not
	    enabled.
*/
func TestTextboxWordWrapScrollbarUnits(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	textboxInstance := layer1.AddTextbox(styleEntry, 2, 2, 10, 2, false)
	setFocusedControl(layer1.layerAlias, textboxInstance.controlAlias, constants.CellTypeTextbox)
	textboxEntry := Textboxes.Get(textboxInstance.layerAlias, textboxInstance.controlAlias)
	textboxInstance.SetWordWrap(true)
	textboxEntry.TextData = [][]rune{[]rune("abcdefghij klmnop qrstuv wxyz 0123456789")}
	textbox.setTextboxMaxScrollBarValues(textboxInstance.layerAlias, textboxInstance.controlAlias)

	verticalScrollbarEntry := ScrollBars.Get(textboxInstance.layerAlias, textboxEntry.VerticalScrollbarAlias)
	horizontalScrollbarEntry := ScrollBars.Get(textboxInstance.layerAlias, textboxEntry.HorizontalScrollbarAlias)
	assert.Equal(test, 4-textboxEntry.Height, verticalScrollbarEntry.MaxScrollValue)
	assert.Equal(test, 2, verticalScrollbarEntry.MaxScrollValue)
	assert.False(test, horizontalScrollbarEntry.IsEnabled)
}

/*
TestTextboxWordWrapInsertKeepsCursorOnRune is a test which verifies that inserting a character in the middle of a
wrapped line puts it at the correct logical index, advances the cursor by one, and leaves the cursor inside the
viewport.

Example:

	Expected Inputs:
	    A width 10, height 3 textbox with word wrap on and the logical line "abcdefghij klmnop qrstuv"; the cursor
	    is set to logical rune 13 and the character 'X' is typed.
	Expected Outputs:
	    TextData[0] rune 13 is 'X', CursorXLocation is 14, and the cursor's display row is within the viewport.
*/
func TestTextboxWordWrapInsertKeepsCursorOnRune(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	textboxInstance := layer1.AddTextbox(styleEntry, 2, 2, 10, 3, false)
	setFocusedControl(layer1.layerAlias, textboxInstance.controlAlias, constants.CellTypeTextbox)
	textboxEntry := Textboxes.Get(textboxInstance.layerAlias, textboxInstance.controlAlias)
	textboxInstance.SetWordWrap(true)
	textboxEntry.TextData = [][]rune{[]rune("abcdefghij klmnop qrstuv")}
	textboxEntry.CursorYLocation = 0
	textboxEntry.CursorXLocation = 13
	textboxEntry.ViewportYLocation = 0

	textbox.UpdateKeyboardEventManually(textboxInstance.layerAlias, textboxInstance.controlAlias, []rune("X"))

	assert.Equal(test, 'X', textboxEntry.TextData[0][13])
	assert.Equal(test, 14, textboxEntry.CursorXLocation)
	wrappedRows := textbox.getWrappedRows(textboxEntry.TextData, textboxEntry.Width, true)
	cursorDisplayRow, _ := textbox.getCursorDisplayCoordinates(wrappedRows, textboxEntry.CursorYLocation, textboxEntry.CursorXLocation)
	assert.GreaterOrEqual(test, cursorDisplayRow, textboxEntry.ViewportYLocation)
	assert.Less(test, cursorDisplayRow, textboxEntry.ViewportYLocation+textboxEntry.Height)
}

/*
TestTextboxWrapDisabledCoordinatesUnchanged is a test which guards that with word wrap OFF getWrappedRows returns
one row per logical line with no rune offset, so display coordinates equal logical coordinates and the existing
non-wrapped behaviour is untouched.

Example:

	Expected Inputs:
	    Three logical lines of ASCII text with word wrap disabled.
	Expected Outputs:
	    Three display rows, each with startRuneIndex 0 and logicalLineIndex equal to its position; a click on any
	    cell yields the same line and rune index it was stamped with.
*/
func TestTextboxWrapDisabledCoordinatesUnchanged(test *testing.T) {
	_, _, _, _ = CommonTestSetup(test)
	textData := [][]rune{[]rune("first line here"), []rune("second"), []rune("third line of text")}
	wrappedRows := textbox.getWrappedRows(textData, 8, false)

	assert.Equal(test, 3, len(wrappedRows))
	for currentIndex := range wrappedRows {
		assert.Equal(test, currentIndex, wrappedRows[currentIndex].logicalLineIndex)
		assert.Equal(test, 0, wrappedRows[currentIndex].startRuneIndex)
		assert.Equal(test, string(textData[currentIndex]), string(wrappedRows[currentIndex].runes))
	}
	gotLine, gotRune := textbox.getLogicalCoordinatesFromDisplay(wrappedRows, 2, 6)
	assert.Equal(test, 2, gotLine)
	assert.Equal(test, 6, gotRune)
}

/*
TestTextboxWordWrapArrowDownStepsWrappedRows is a test which verifies that the down arrow walks a word-wrapped
logical line one VISIBLE row at a time and pulls the viewport with it, instead of getting stuck because there is
only one logical line.

Example:

	Expected Inputs:
	    A width 10, height 2 textbox with word wrap on and one logical line that wraps to four display rows; the
	    cursor starts at the top and "down" is pressed three times, then a fourth time.
	Expected Outputs:
	    The cursor's display row goes 0 -> 1 -> 2 -> 3 and then stays at 3; ViewportYLocation ends at 2 so the
	    cursor row is inside the two-row window.
*/
func TestTextboxWordWrapArrowDownStepsWrappedRows(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	textboxInstance := layer1.AddTextbox(styleEntry, 2, 2, 10, 2, false)
	setFocusedControl(layer1.layerAlias, textboxInstance.controlAlias, constants.CellTypeTextbox)
	textboxEntry := Textboxes.Get(textboxInstance.layerAlias, textboxInstance.controlAlias)
	textboxInstance.SetWordWrap(true)
	textboxEntry.TextData = [][]rune{[]rune("abcdefghij klmnop qrstuv wxyz 0123456789")}
	textboxEntry.CursorYLocation = 0
	textboxEntry.CursorXLocation = 0
	textboxEntry.ViewportYLocation = 0

	displayRowOf := func() int {
		wrappedRows := textbox.getWrappedRows(textboxEntry.TextData, textboxEntry.Width, true)
		displayRow, _ := textbox.getCursorDisplayCoordinates(wrappedRows, textboxEntry.CursorYLocation, textboxEntry.CursorXLocation)
		return displayRow
	}

	assert.Equal(test, 0, displayRowOf())
	for expectedRow := 1; expectedRow <= 3; expectedRow++ {
		textbox.UpdateKeyboardEventManually(textboxInstance.layerAlias, textboxInstance.controlAlias, []rune("down"))
		assert.Equalf(test, expectedRow, displayRowOf(), "after %d down presses", expectedRow)
	}
	textbox.UpdateKeyboardEventManually(textboxInstance.layerAlias, textboxInstance.controlAlias, []rune("down"))
	assert.Equal(test, 3, displayRowOf())
	assert.Equal(test, 2, textboxEntry.ViewportYLocation)
}

/*
TestTextboxArrowDownWrapOffMovesLogicalLine is a test which guards that with word wrap OFF the down and up arrows
still move exactly one logical line, keeping the printed column where the destination line is long enough.

Example:

	Expected Inputs:
	    A width 20 textbox with word wrap off holding three logical lines; the cursor is on line 0 at rune 4 and
	    "down" then "down" then "up" are pressed.
	Expected Outputs:
	    After the first "down" CursorYLocation is 1, after the second it is 2, and after "up" it is 1; the cursor
	    stays at rune 4 on lines that are long enough.
*/
func TestTextboxArrowDownWrapOffMovesLogicalLine(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	textboxInstance := layer1.AddTextbox(styleEntry, 2, 2, 20, 4, false)
	setFocusedControl(layer1.layerAlias, textboxInstance.controlAlias, constants.CellTypeTextbox)
	textboxEntry := Textboxes.Get(textboxInstance.layerAlias, textboxInstance.controlAlias)
	textboxEntry.TextData = [][]rune{[]rune("first line of text "), []rune("second line here "), []rune("third line as well ")}
	textboxEntry.CursorYLocation = 0
	textboxEntry.CursorXLocation = 4

	textbox.UpdateKeyboardEventManually(textboxInstance.layerAlias, textboxInstance.controlAlias, []rune("down"))
	assert.Equal(test, 1, textboxEntry.CursorYLocation)
	assert.Equal(test, 4, textboxEntry.CursorXLocation)
	textbox.UpdateKeyboardEventManually(textboxInstance.layerAlias, textboxInstance.controlAlias, []rune("down"))
	assert.Equal(test, 2, textboxEntry.CursorYLocation)
	textbox.UpdateKeyboardEventManually(textboxInstance.layerAlias, textboxInstance.controlAlias, []rune("up"))
	assert.Equal(test, 1, textboxEntry.CursorYLocation)
}

/*
TestTextboxSetTextLinesEndWithSentinel is a test which verifies that SetText replaces the blank line a new textbox
starts with, and stores each line with the trailing blank sentinel rune, the same shape typed lines have, so a caret
can rest past the last character.

Example:

	Expected Inputs:
	    A new textbox on which SetText("abc\nde") is called.
	Expected Outputs:
	    The text is exactly the two rune slices [a b c space] and [d e space], with no blank line before them.
*/
func TestTextboxSetTextLinesEndWithSentinel(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	textboxInstance := layer1.AddTextbox(styleEntry, 2, 2, 20, 4, false)
	textboxInstance.SetText("abc\nde")
	textboxEntry := Textboxes.Get(textboxInstance.layerAlias, textboxInstance.controlAlias)

	assert.Equal(test, [][]rune{[]rune("abc "), []rune("de ")}, textboxEntry.TextData)
}

/*
TestTextboxDeleteLastCharacterThenType is a test which verifies that forward delete removes the final character
of a line set with SetText, leaving the caret over the empty line's sentinel, and that typing then inserts a
character and advances the caret.

Example:

	Expected Inputs:
	    A width 20 textbox whose text is set to "X", which leaves the caret on the 'X'. "delete" is pressed, then
	    the character 'Y' is typed.
	Expected Outputs:
	    After "delete": the text is the single line [space], with the cursor at line 0, column 0.
	    After typing 'Y': the text is the single line [Y space], CursorXLocation is 1, and GetText returns "Y".
*/
func TestTextboxDeleteLastCharacterThenType(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup(test)
	textboxInstance := layer1.AddTextbox(styleEntry, 2, 2, 20, 4, false)
	setFocusedControl(layer1.layerAlias, textboxInstance.controlAlias, constants.CellTypeTextbox)
	textboxInstance.SetText("X")
	textboxEntry := Textboxes.Get(textboxInstance.layerAlias, textboxInstance.controlAlias)

	textbox.UpdateKeyboardEventManually(textboxInstance.layerAlias, textboxInstance.controlAlias, []rune("delete"))
	assert.Equal(test, [][]rune{[]rune(" ")}, textboxEntry.TextData)
	assert.Equal(test, 0, textboxEntry.CursorXLocation)
	assert.Equal(test, 0, textboxEntry.CursorYLocation)

	textbox.UpdateKeyboardEventManually(textboxInstance.layerAlias, textboxInstance.controlAlias, []rune("Y"))
	assert.Equal(test, [][]rune{[]rune("Y ")}, textboxEntry.TextData)
	assert.Equal(test, 1, textboxEntry.CursorXLocation)
	assert.Equal(test, "Y", textboxInstance.GetText())
}

/*
TestTextboxStaleCursorIsClamped is a test which allows you to verify that a textbox cannot index outside its text
through a cursor left out of range, whether set through the exported entry or left behind by deleting a highlighted
whole line, each of which used to panic or leave the cursor at minus one.

Example:

	Expected Inputs:
		Case one: a textbox holding "ab" whose cursor is set to line 9, column 40, then Right. Case two: a textbox
		holding "ab" with its whole line, sentinel included, highlighted from column 0 to column 2 and deleted.

	Expected Outputs:
		Case one: no panic, and the cursor ends inside the text. Case two: the line is restored to a single blank
		sentinel and the cursor is at line 0, column 0.
*/
func TestTextboxStaleCursorIsClamped(test *testing.T) {
	test.Run("cursor set out of range", func(test *testing.T) {
		layerAlias, styleEntry := setupInputTest(test)
		textboxInstance := textbox.Add(layerAlias, "box", styleEntry, 2, 2, 20, 5, false)
		textboxInstance.SetText("ab")
		textboxEntry := GetTextbox(layerAlias, "box")
		textboxEntry.CursorYLocation = 9
		textboxEntry.CursorXLocation = 40
		textbox.UpdateKeyboardEventManually(layerAlias, "box", []rune("right"))
		UpdateDisplay(false)
		if textboxEntry.CursorYLocation >= len(textboxEntry.TextData) ||
			textboxEntry.CursorXLocation >= len(textboxEntry.TextData[textboxEntry.CursorYLocation]) {
			test.Fatalf("expected the cursor inside the text, got line %d, column %d", textboxEntry.CursorYLocation, textboxEntry.CursorXLocation)
		}
	})
	test.Run("whole line highlighted and deleted", func(test *testing.T) {
		layerAlias, styleEntry := setupInputTest(test)
		textbox.Add(layerAlias, "box", styleEntry, 2, 2, 20, 5, false)
		textboxEntry := GetTextbox(layerAlias, "box")
		textboxEntry.TextData = [][]rune{[]rune("ab ")}
		textboxEntry.IsHighlightActive = true
		textboxEntry.HighlightStartX, textboxEntry.HighlightStartY = 0, 0
		textboxEntry.CursorXLocation, textboxEntry.CursorYLocation = 2, 0
		textbox.deleteHighlightedText(textboxEntry)
		if textboxEntry.CursorXLocation != 0 || textboxEntry.CursorYLocation != 0 || string(textboxEntry.TextData[0]) != " " {
			test.Fatalf("expected an empty line with the cursor at 0, 0, got line %q with the cursor at %d, %d",
				string(textboxEntry.TextData[0]), textboxEntry.CursorYLocation, textboxEntry.CursorXLocation)
		}
	})
}

/*
assertTextboxScrollbarLayout is a method which allows you to fail the test unless a textbox's horizontal and vertical
scroll bars sit at the expected columns and rows with the expected lengths, given as horizontal column, row and
length followed by vertical column, row and length.

Example:

	assertTextboxScrollbarLayout(test, "after Add", layerAlias, textboxEntry, 1, 6, 12, 13, 1, 5)
*/
func assertTextboxScrollbarLayout(test *testing.T, context string, layerAlias string, textboxEntry *types.TextboxEntryType,
	hXLocation int, hYLocation int, hLength int, vXLocation int, vYLocation int, vLength int) {
	test.Helper()
	hEntry := ScrollBars.Get(layerAlias, textboxEntry.HorizontalScrollbarAlias)
	vEntry := ScrollBars.Get(layerAlias, textboxEntry.VerticalScrollbarAlias)
	if hEntry.XLocation != hXLocation || hEntry.YLocation != hYLocation || hEntry.Length != hLength {
		test.Fatalf("%s: expected the horizontal scroll bar at (%d, %d) with length %d, got (%d, %d) with length %d", context,
			hXLocation, hYLocation, hLength, hEntry.XLocation, hEntry.YLocation, hEntry.Length)
	}
	if vEntry.XLocation != vXLocation || vEntry.YLocation != vYLocation || vEntry.Length != vLength {
		test.Fatalf("%s: expected the vertical scroll bar at (%d, %d) with length %d, got (%d, %d) with length %d", context,
			vXLocation, vYLocation, vLength, vEntry.XLocation, vEntry.YLocation, vEntry.Length)
	}
}

/*
TestTextboxSetSizePositionAndBorderMoveScrollbars is a test which allows you to verify that moving a textbox,
resizing it, and turning its border off each move its scroll bars and tooltip with it, and that sizes below one cell
are ignored.

Example:

	Expected Inputs:
		Bordered textbox at (2, 2), 10 wide and 3 tall. SetPosition(4, 3), then SetSize(15, 4), then SetSize(0, 4)
		and SetSize(15, 0), then SetBorderDrawn(false).

	Expected Outputs:
		Horizontal scroll bar (column, row, length) and vertical scroll bar after each step:
		after Add (1, 6, 12) and (13, 1, 5); after SetPosition (3, 7, 12) and (15, 2, 5); after SetSize (3, 8, 17)
		and (20, 2, 6), with the tooltip hotspot at (4, 3) sized 15 by 4; the invalid sizes leave the size (15, 4);
		without the border (4, 7, 15) and (19, 3, 4).
*/
func TestTextboxSetSizePositionAndBorderMoveScrollbars(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	textboxInstance := textbox.Add(layerAlias, "textbox", styleEntry, 2, 2, 10, 3, true)
	textboxEntry := Textboxes.Get(layerAlias, "textbox")
	assertTextboxScrollbarLayout(test, "after Add", layerAlias, textboxEntry, 1, 6, 12, 13, 1, 5)

	textboxInstance.SetPosition(4, 3)
	assertTextboxScrollbarLayout(test, "after SetPosition", layerAlias, textboxEntry, 3, 7, 12, 15, 2, 5)

	textboxInstance.SetSize(15, 4)
	assertTextboxScrollbarLayout(test, "after SetSize", layerAlias, textboxEntry, 3, 8, 17, 20, 2, 6)
	tooltipEntry := Tooltips.Get(layerAlias, textboxEntry.TooltipAlias)
	if tooltipEntry.HotspotXLocation != 4 || tooltipEntry.HotspotYLocation != 3 || tooltipEntry.HotspotWidth != 15 || tooltipEntry.HotspotHeight != 4 {
		test.Fatalf("expected the tooltip hotspot at (4, 3) sized 15 by 4, got (%d, %d) sized %d by %d", tooltipEntry.HotspotXLocation,
			tooltipEntry.HotspotYLocation, tooltipEntry.HotspotWidth, tooltipEntry.HotspotHeight)
	}

	textboxInstance.SetSize(0, 4)
	textboxInstance.SetSize(15, 0)
	if width, height := textboxInstance.GetSize(); width != 15 || height != 4 {
		test.Fatalf("expected invalid sizes to be ignored, got (%d, %d)", width, height)
	}

	textboxInstance.SetBorderDrawn(false)
	assertTextboxScrollbarLayout(test, "without border", layerAlias, textboxEntry, 4, 7, 15, 19, 3, 4)
}

/*
TestTextboxSetSizePullsViewportBack is a test which allows you to verify that resizing a scrolled textbox keeps the
viewport where it is while it is still in range, pulls it back when the larger text area would scroll past the end of
the text, and turns the scroll bars off once all the text fits.

Example:

	Expected Inputs:
		Textbox at (2, 2), 10 wide and 3 tall, holding "line0" to "line9" and a line of 25 x characters, so 11
		rows whose widest is 26 columns with its sentinel. The viewport is set to (12, 8), the furthest down it can
		scroll. Then SetSize(10, 8), then SetSize(30, 12).

	Expected Outputs:
		After SetSize(10, 8): viewport (12, 3), with both scroll bars enabled and their values matching. After
		SetSize(30, 12): viewport (0, 0) with both scroll bars disabled.
*/
func TestTextboxSetSizePullsViewportBack(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	textboxInstance := textbox.Add(layerAlias, "textbox", styleEntry, 2, 2, 10, 3, false)
	textboxInstance.SetText("line0\nline1\nline2\nline3\nline4\nline5\nline6\nline7\nline8\nline9\nxxxxxxxxxxxxxxxxxxxxxxxxx")
	textboxInstance.SetViewport(12, 8)
	textboxEntry := Textboxes.Get(layerAlias, "textbox")
	hEntry := ScrollBars.Get(layerAlias, textboxEntry.HorizontalScrollbarAlias)
	vEntry := ScrollBars.Get(layerAlias, textboxEntry.VerticalScrollbarAlias)

	textboxInstance.SetSize(10, 8)
	if textboxEntry.ViewportXLocation != 12 || textboxEntry.ViewportYLocation != 3 {
		test.Fatalf("expected viewport (12, 3), got (%d, %d)", textboxEntry.ViewportXLocation, textboxEntry.ViewportYLocation)
	}
	if !hEntry.IsEnabled || !vEntry.IsEnabled || hEntry.ScrollValue != 12 || vEntry.ScrollValue != 3 {
		test.Fatalf("expected both scroll bars enabled with values 12 and 3, got %t %t with values %d and %d",
			hEntry.IsEnabled, vEntry.IsEnabled, hEntry.ScrollValue, vEntry.ScrollValue)
	}

	textboxInstance.SetSize(30, 12)
	if textboxEntry.ViewportXLocation != 0 || textboxEntry.ViewportYLocation != 0 {
		test.Fatalf("expected viewport (0, 0), got (%d, %d)", textboxEntry.ViewportXLocation, textboxEntry.ViewportYLocation)
	}
	if hEntry.IsEnabled || vEntry.IsEnabled {
		test.Fatalf("expected both scroll bars disabled, got %t and %t", hEntry.IsEnabled, vEntry.IsEnabled)
	}
}

/*
TestTextboxSetTextReplacesTextAndResetsState is a test which allows you to verify that SetText replaces any text
already in a textbox rather than adding to it, and moves the cursor, viewport and highlight back to the start.

Example:

	Expected Inputs:
		Textbox 10 wide and 2 tall holding "line0" to "line5", with its cursor at (3, 5), its viewport at (2, 4)
		and a highlight from (1, 1) to (3, 5) that is active and toggled on. Then SetText("new"), then SetText("").

	Expected Outputs:
		After SetText("new"): the text is the single line [n e w space], the cursor and viewport are at (0, 0),
		the highlight is off with every highlight position 0, the vertical scroll bar is disabled, and GetText
		returns "new". After SetText(""): the text is the single line [space] and GetText returns "".
*/
func TestTextboxSetTextReplacesTextAndResetsState(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	textboxInstance := textbox.Add(layerAlias, "textbox", styleEntry, 2, 2, 10, 2, false)
	textboxInstance.SetText("line0\nline1\nline2\nline3\nline4\nline5")
	textboxEntry := Textboxes.Get(layerAlias, "textbox")
	textboxEntry.CursorXLocation, textboxEntry.CursorYLocation = 3, 5
	textboxEntry.ViewportXLocation, textboxEntry.ViewportYLocation = 2, 4
	textboxEntry.HighlightStartX, textboxEntry.HighlightStartY = 1, 1
	textboxEntry.HighlightEndX, textboxEntry.HighlightEndY = 3, 5
	textboxEntry.IsHighlightActive, textboxEntry.IsHighlightModeToggled = true, true

	textboxInstance.SetText("new")
	assert.Equal(test, [][]rune{[]rune("new ")}, textboxEntry.TextData)
	if textboxEntry.CursorXLocation != 0 || textboxEntry.CursorYLocation != 0 ||
		textboxEntry.ViewportXLocation != 0 || textboxEntry.ViewportYLocation != 0 {
		test.Fatalf("expected the cursor and viewport at (0, 0), got (%d, %d) and (%d, %d)", textboxEntry.CursorXLocation,
			textboxEntry.CursorYLocation, textboxEntry.ViewportXLocation, textboxEntry.ViewportYLocation)
	}
	if textboxEntry.IsHighlightActive || textboxEntry.IsHighlightModeToggled || textboxEntry.HighlightStartX != 0 ||
		textboxEntry.HighlightStartY != 0 || textboxEntry.HighlightEndX != 0 || textboxEntry.HighlightEndY != 0 {
		test.Fatalf("expected the highlight to be cleared, got active %t, toggled %t, from (%d, %d) to (%d, %d)",
			textboxEntry.IsHighlightActive, textboxEntry.IsHighlightModeToggled, textboxEntry.HighlightStartX,
			textboxEntry.HighlightStartY, textboxEntry.HighlightEndX, textboxEntry.HighlightEndY)
	}
	if ScrollBars.Get(layerAlias, textboxEntry.VerticalScrollbarAlias).IsEnabled {
		test.Fatalf("expected the vertical scroll bar to be disabled once the text fits")
	}
	assert.Equal(test, "new", textboxInstance.GetText())

	textboxInstance.SetText("")
	assert.Equal(test, [][]rune{[]rune(" ")}, textboxEntry.TextData)
	assert.Equal(test, "", textboxInstance.GetText())
}

/*
TestTextboxGetTextReturnsSetText is a test which allows you to verify that GetText returns the text given to SetText,
without the blank sentinel each line is stored with, including real trailing and leading spaces, empty lines and
wide characters, and that carriage return and newline pairs come back as single newlines.

Example:

	Expected Inputs:
		SetText with each of "", "abc", "abc\nde", "trailing space \n  indented", "\n\nmiddle\n\n", "中文\nEnglish",
		and "one\r\ntwo".

	Expected Outputs:
		GetText returns each input unchanged, except "one\r\ntwo", which comes back as "one\ntwo".
*/
func TestTextboxGetTextReturnsSetText(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	textboxInstance := textbox.Add(layerAlias, "textbox", styleEntry, 2, 2, 20, 4, false)
	for _, testCase := range []struct {
		text         string
		expectedText string
	}{
		{"", ""},
		{"abc", "abc"},
		{"abc\nde", "abc\nde"},
		{"trailing space \n  indented", "trailing space \n  indented"},
		{"\n\nmiddle\n\n", "\n\nmiddle\n\n"},
		{"中文\nEnglish", "中文\nEnglish"},
		{"one\r\ntwo", "one\ntwo"},
	} {
		textboxInstance.SetText(testCase.text)
		if text := textboxInstance.GetText(); text != testCase.expectedText {
			test.Fatalf("SetText(%q): expected GetText to return %q, got %q", testCase.text, testCase.expectedText, text)
		}
	}
}

/*
TestTextboxDeleteHighlightKeepsRealTrailingSpace is a test which allows you to verify that deleting a highlight that
runs to the end of a line puts the line's sentinel back even when the text left behind ends in a real space, so the
space is not mistaken for the sentinel and lost.

Example:

	Expected Inputs:
		Textbox holding "ab c", stored as [a b space c space]. A highlight from column 3, the 'c', to column 4,
		the sentinel, is deleted.

	Expected Outputs:
		The line is [a b space space], the cursor is at line 0, column 3, and GetText returns "ab ".
*/
func TestTextboxDeleteHighlightKeepsRealTrailingSpace(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	textboxInstance := textbox.Add(layerAlias, "textbox", styleEntry, 2, 2, 20, 4, false)
	textboxInstance.SetText("ab c")
	textboxEntry := Textboxes.Get(layerAlias, "textbox")
	textboxEntry.IsHighlightActive = true
	textboxEntry.HighlightStartX, textboxEntry.HighlightStartY = 3, 0
	textboxEntry.CursorXLocation, textboxEntry.CursorYLocation = 4, 0
	textbox.deleteHighlightedText(textboxEntry)
	assert.Equal(test, [][]rune{[]rune("ab  ")}, textboxEntry.TextData)
	if textboxEntry.CursorXLocation != 3 || textboxEntry.CursorYLocation != 0 {
		test.Fatalf("expected the cursor at line 0, column 3, got line %d, column %d", textboxEntry.CursorYLocation, textboxEntry.CursorXLocation)
	}
	assert.Equal(test, "ab ", textboxInstance.GetText())
}

/*
TestTextboxDeleteMultiLineHighlightJoinsLines is a test which allows you to verify that deleting a highlight spanning
several lines joins the text before it on the first line to the text after it on the last line, removes every line in
between, keeps a single sentinel on the joined line, and moves the cursor to where the highlight began.

Example:

	Expected Inputs:
		Each case sets the text, then the highlight start and the cursor, given as (column, line), and deletes:
		"abc\ndef" from (1, 0) to (1, 1); "abc\nmid\ndef" from (1, 0) to (0, 2); "abc\ndef" with the cursor at
		(1, 0) before the highlight start at (1, 1); "abc\ndef" from (3, 0), the first line's sentinel, to (0, 1);
		"abc\ndef\nzzz" from (0, 0) to (3, 1), the second line's sentinel; and "abc\ndef" from (-5, -1) to (40, 9),
		both out of range.

	Expected Outputs:
		Text and cursor (column, line) after each: ["af "] at (1, 0); ["aef "] at (1, 0); ["af "] at (1, 0);
		["abcef "] at (3, 0); [" ", "zzz "] at (0, 0); and [" "] at (0, 0), with GetText returning "af", "aef",
		"af", "abcef", "\nzzz" and "" respectively.
*/
func TestTextboxDeleteMultiLineHighlightJoinsLines(test *testing.T) {
	testCases := []struct {
		name              string
		text              string
		highlightStartX   int
		highlightStartY   int
		cursorXLocation   int
		cursorYLocation   int
		expectedTextData  [][]rune
		expectedCursorX   int
		expectedCursorY   int
		expectedFinalText string
	}{
		{"two lines", "abc\ndef", 1, 0, 1, 1, [][]rune{[]rune("af ")}, 1, 0, "af"},
		{"three lines", "abc\nmid\ndef", 1, 0, 0, 2, [][]rune{[]rune("aef ")}, 1, 0, "aef"},
		{"cursor before highlight start", "abc\ndef", 1, 1, 1, 0, [][]rune{[]rune("af ")}, 1, 0, "af"},
		{"starting on the first line's sentinel", "abc\ndef", 3, 0, 0, 1, [][]rune{[]rune("abcef ")}, 3, 0, "abcef"},
		{"through the last line's sentinel", "abc\ndef\nzzz", 0, 0, 3, 1, [][]rune{[]rune(" "), []rune("zzz ")}, 0, 0, "\nzzz"},
		{"out of range positions", "abc\ndef", -5, -1, 40, 9, [][]rune{[]rune(" ")}, 0, 0, ""},
	}
	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			layerAlias, styleEntry := setupInputTest(test)
			textboxInstance := textbox.Add(layerAlias, "textbox", styleEntry, 2, 2, 20, 4, false)
			textboxInstance.SetText(testCase.text)
			textboxEntry := Textboxes.Get(layerAlias, "textbox")
			textboxEntry.IsHighlightActive = true
			textboxEntry.HighlightStartX, textboxEntry.HighlightStartY = testCase.highlightStartX, testCase.highlightStartY
			textboxEntry.CursorXLocation, textboxEntry.CursorYLocation = testCase.cursorXLocation, testCase.cursorYLocation
			textbox.deleteHighlightedText(textboxEntry)
			assert.Equal(test, testCase.expectedTextData, textboxEntry.TextData)
			if textboxEntry.CursorXLocation != testCase.expectedCursorX || textboxEntry.CursorYLocation != testCase.expectedCursorY {
				test.Fatalf("expected the cursor at column %d, line %d, got column %d, line %d", testCase.expectedCursorX,
					testCase.expectedCursorY, textboxEntry.CursorXLocation, textboxEntry.CursorYLocation)
			}
			assert.Equal(test, testCase.expectedFinalText, textboxInstance.GetText())
			if textboxEntry.IsHighlightActive {
				test.Fatalf("expected the highlight to be turned off")
			}
		})
	}
}
