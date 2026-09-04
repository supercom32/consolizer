package consolizer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/stringformat"
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
	layer1, _, _, styleEntry := CommonTestSetup()
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
	layer1, _, _, styleEntry := CommonTestSetup()
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
	layer1, _, _, styleEntry := CommonTestSetup()
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
	_, _, _, _ = CommonTestSetup()
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
	layer1, _, _, styleEntry := CommonTestSetup()
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

	eventStateMemory.stateId = 0
	SetMouseStatus(3, 2, 1, "")
	textbox.updateMouseEvent()
	assert.Equal(test, 0, textboxEntry.CursorXLocation)
	assert.Equal(test, 0, textboxEntry.CursorYLocation)

	eventStateMemory.stateId = 0
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
	layer1, _, _, styleEntry := CommonTestSetup()
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
	layer1, _, _, styleEntry := CommonTestSetup()
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
	_, _, _, _ = CommonTestSetup()
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
	layer1, _, _, styleEntry := CommonTestSetup()
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

	eventStateMemory.stateId = 0
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
	layer1, _, _, styleEntry := CommonTestSetup()
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
	layer1, _, _, styleEntry := CommonTestSetup()
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
	layer1, _, _, styleEntry := CommonTestSetup()
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
	_, _, _, _ = CommonTestSetup()
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
	layer1, _, _, styleEntry := CommonTestSetup()
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
	layer1, _, _, styleEntry := CommonTestSetup()
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
TestTextboxSetTextLinesEndWithSentinel is a test which verifies that SetText stores each line with the trailing
blank sentinel rune, the same shape typed lines have, so a caret can rest past the last character.

Example:

	Expected Inputs:
	    A textbox on which SetText("abc\nde") is called.
	Expected Outputs:
	    The two non-seed lines are the rune slices [a b c space] and [d e space].
*/
func TestTextboxSetTextLinesEndWithSentinel(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	textboxInstance := layer1.AddTextbox(styleEntry, 2, 2, 20, 4, false)
	textboxInstance.SetText("abc\nde")
	textboxEntry := Textboxes.Get(textboxInstance.layerAlias, textboxInstance.controlAlias)

	lines := textboxEntry.TextData
	assert.Equal(test, []rune("abc "), lines[len(lines)-2])
	assert.Equal(test, []rune("de "), lines[len(lines)-1])
}

/*
TestTextboxDeleteLastCharacterThenType is a test which verifies that forward delete removes the final character
of a line set with SetText, leaving the caret over the empty line's sentinel, and that typing then inserts a
character and advances the caret.

Example:

	Expected Inputs:
	    A width 20 textbox whose only content line is set to "X"; the caret is put on the 'X', "delete" is
	    pressed, then the character 'Y' is typed.
	Expected Outputs:
	    After "delete": the line is the single rune slice [space], CursorXLocation 0, CursorYLocation the last
	    line index.
	    After typing 'Y': the line is [Y space] and CursorXLocation is 1.
*/
func TestTextboxDeleteLastCharacterThenType(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	textboxInstance := layer1.AddTextbox(styleEntry, 2, 2, 20, 4, false)
	setFocusedControl(layer1.layerAlias, textboxInstance.controlAlias, constants.CellTypeTextbox)
	textboxInstance.SetText("X")
	textboxEntry := Textboxes.Get(textboxInstance.layerAlias, textboxInstance.controlAlias)
	lastLineIndex := len(textboxEntry.TextData) - 1
	textboxEntry.CursorYLocation = lastLineIndex
	textboxEntry.CursorXLocation = 0

	textbox.UpdateKeyboardEventManually(textboxInstance.layerAlias, textboxInstance.controlAlias, []rune("delete"))
	assert.Equal(test, []rune(" "), textboxEntry.TextData[lastLineIndex])
	assert.Equal(test, 0, textboxEntry.CursorXLocation)
	assert.Equal(test, lastLineIndex, textboxEntry.CursorYLocation)

	textbox.UpdateKeyboardEventManually(textboxInstance.layerAlias, textboxInstance.controlAlias, []rune("Y"))
	assert.Equal(test, []rune("Y "), textboxEntry.TextData[lastLineIndex])
	assert.Equal(test, 1, textboxEntry.CursorXLocation)
}
