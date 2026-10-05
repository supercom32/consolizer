package consolizer

import (
	"github.com/gdamore/tcell/v2"
	"github.com/supercom32/consolizer/types"
	"testing"
)

/*
addTestFileMenu is a method which allows you to add a file menu with the single heading "File" and the items
"Open" and "Save" at location (2, 2).

Example:

	fileMenu := addTestFileMenu(layerAlias, styleEntry)
*/
func addTestFileMenu(layerAlias string, styleEntry types.TuiStyleEntryType) FileMenuInstanceType {
	selectionEntry := types.NewSelectionEntry()
	selectionEntry.Add("open", "Open")
	selectionEntry.Add("save", "Save")
	return FileMenu.Add(layerAlias, "menu", styleEntry, []string{"File"}, []types.SelectionEntryType{selectionEntry}, 2, 2, true)
}

/*
TestFileMenuEscClosesOpenMenu is a test which allows you to verify that Esc closes an open file menu and is consumed.
In addition, the following should be noted:

  - Before the fix, the file menu only reacted to "escape", which tcell never reports (it reports "Esc"), so Esc
    could never close an open file menu.

Example:

	Expected Inputs:
		File menu with heading "File", opened on heading 0, then the key Esc.

	Expected Outputs:
		IsOpen is false, ActiveHeadingIndex is -1, and the keyboard buffer is empty.
*/
func TestFileMenuEscClosesOpenMenu(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	fileMenu := addTestFileMenu(layerAlias, styleEntry)
	fileMenuEntry := FileMenus.Get(layerAlias, "menu")
	fileMenuEntry.ActiveHeadingIndex = 0
	fileMenuEntry.IsSubmenuOpen = true
	simScreen := startInputSimulation(test)

	pressKey(simScreen, tcell.KeyEscape, 0, tcell.ModNone)

	if fileMenu.IsOpen() || fileMenuEntry.ActiveHeadingIndex != -1 {
		test.Fatalf("expected Esc to close the menu, got open %v with active heading %d", fileMenu.IsOpen(), fileMenuEntry.ActiveHeadingIndex)
	}
	assertKeyboardBuffer(test, "after Esc on an open file menu")
}

/*
TestFileMenuEscWithNoOpenMenuIsNotConsumed is a test which allows you to verify that Esc reaches the keyboard buffer
when no file menu is open, so the application can still treat it as Back.

Example:

	Expected Inputs:
		Closed file menu with heading "File", then the key Esc.

	Expected Outputs:
		IsOpen is false and the keyboard buffer is ["esc"].
*/
func TestFileMenuEscWithNoOpenMenuIsNotConsumed(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	fileMenu := addTestFileMenu(layerAlias, styleEntry)
	simScreen := startInputSimulation(test)

	pressKey(simScreen, tcell.KeyEscape, 0, tcell.ModNone)

	if fileMenu.IsOpen() {
		test.Fatalf("expected the menu to stay closed")
	}
	assertKeyboardBuffer(test, "after Esc with no open file menu", "esc")
}
