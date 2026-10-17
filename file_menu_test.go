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
	fileMenuEntry := fileMenuRegistry.Get(layerAlias, "menu")
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

/*
TestFileMenuGetSelectedItemIgnoresStaleSelection is a test which allows you to verify that GetSelectedItem does not
index a submenu with a recorded selection that no longer fits it, which used to panic.

Example:

	Expected Inputs:
		A file menu with the heading "File" whose submenu has two items, with the submenu's selection set to five
		through the exported selector entry.

	Expected Outputs:
		GetSelectedItem reports nothing selected: minus one, minus one, and two empty strings.
*/
func TestFileMenuGetSelectedItemIgnoresStaleSelection(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	menuSelection := types.NewSelectionEntry()
	menuSelection.Add("open", "Open")
	menuSelection.Add("exit", "Exit")
	fileMenu := FileMenu.Add(layerAlias, "menu", styleEntry, []string{"File"}, []types.SelectionEntryType{menuSelection}, 0, 0, true)
	getSelector(layerAlias, fileMenuRegistry.Get(layerAlias, "menu").SelectorAliases[0]).ItemSelected = 5
	headingIndex, itemIndex, itemAlias, itemValue := fileMenu.GetSelectedItem()
	if headingIndex != -1 || itemIndex != -1 || itemAlias != "" || itemValue != "" {
		test.Fatalf("expected nothing selected, got (%d, %d, %q, %q)", headingIndex, itemIndex, itemAlias, itemValue)
	}
}

/*
TestFileMenuOpenCloseEnableAndSetMenus is a test which allows you to verify that a file menu can be opened, closed,
enabled, disabled and given new headings through its methods, that invalid requests are ignored or refused, and that
replacing the headings replaces their submenus.

Example:

	Expected Inputs:
		The test file menu with the heading "File". Open(0), Open(5), Close, SetEnabled(false), Open(0),
		SetEnabled(true), SetMenus with two headings and one selection, then SetMenus with "File" and "Edit" and a
		selection for each, then Open(1).

	Expected Outputs:
		IsOpen reads true after Open(0), true after the ignored Open(5), false after Close, and false after Open(0) on
		the disabled menu, which reads IsEnabled false. The mismatched SetMenus returns an error and keeps one
		submenu. The second SetMenus returns no error, removes the old submenu's selector, creates two new ones, and
		Open(1) opens the second heading.
*/
func TestFileMenuOpenCloseEnableAndSetMenus(test *testing.T) {
	layerAlias, styleEntry := setupInputTest(test)
	fileMenu := addTestFileMenu(layerAlias, styleEntry)
	if !fileMenu.IsEnabled() {
		test.Fatalf("expected the file menu to start enabled")
	}
	fileMenu.Open(0)
	fileMenu.Open(5)
	if !fileMenu.IsOpen() {
		test.Fatalf("expected the file menu to be open after Open(0), and an invalid heading to be ignored")
	}
	fileMenu.Close()
	if fileMenu.IsOpen() {
		test.Fatalf("expected the file menu to be closed after Close")
	}
	fileMenu.SetEnabled(false)
	fileMenu.Open(0)
	if fileMenu.IsEnabled() || fileMenu.IsOpen() {
		test.Fatalf("expected a disabled file menu to stay closed, got enabled %t and open %t", fileMenu.IsEnabled(), fileMenu.IsOpen())
	}
	fileMenu.SetEnabled(true)

	getSelectorAliases := func() []string {
		defer lockUIState()()
		return append([]string(nil), fileMenuRegistry.Get(layerAlias, "menu").SelectorAliases...)
	}
	oldSelectorAliases := getSelectorAliases()
	if err := fileMenu.SetMenus([]string{"File", "Edit"}, []types.SelectionEntryType{getTestSelectionEntry(2)}); err == nil {
		test.Fatalf("expected an error for two headings and one selection")
	}
	if selectorAliases := getSelectorAliases(); len(selectorAliases) != 1 {
		test.Fatalf("expected the refused SetMenus to keep one submenu, got %d", len(selectorAliases))
	}
	if err := fileMenu.SetMenus([]string{"File", "Edit"}, []types.SelectionEntryType{getTestSelectionEntry(2), getTestSelectionEntry(3)}); err != nil {
		test.Fatalf("expected no error, got %v", err)
	}
	newSelectorAliases := getSelectorAliases()
	if len(newSelectorAliases) != 2 || IsSelectorExists(layerAlias, oldSelectorAliases[0]) {
		test.Fatalf("expected two new submenus and the old one removed, got %d submenus and old one existing %t",
			len(newSelectorAliases), IsSelectorExists(layerAlias, oldSelectorAliases[0]))
	}
	fileMenu.Open(1)
	if !fileMenu.IsOpen() {
		test.Fatalf("expected Open(1) to open the second heading")
	}
}
