package consolizer

import (
	"fmt"
	"github.com/supercom32/consolizer/memory"
	"github.com/supercom32/consolizer/stringformat"
	"slices"
	"sync/atomic"
	"time"

	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/types"
)

/*
SelectorInstanceType is a structure which represents an instance of a selector control.
*/
type SelectorInstanceType struct {
	BaseControlInstanceType
}

type selectorType struct{}

var Selector selectorType
var selectorRegistry = memory.NewControlMemoryManager[types.SelectorEntryType]()

/*
selectorClickType is a structure which records the most recent mouse click on a selector item, so that a following
click on the same item can be recognized as a double click.

Example:

	var lastClick selectorClickType
*/
type selectorClickType struct {
	layerAlias    string
	selectorAlias string
	itemIndex     int
	clickTime     time.Time
}

// selectorLastClick is only read and written by the event goroutine, so it needs no synchronization of its own.
var selectorLastClick selectorClickType

/*
selectorPressType is a structure which records the selector a held mouse button was first pressed on, and the item that
press last reported as selected, so that a press is only ever applied to the selector it began on.

Example:

	var press selectorPressType
*/
type selectorPressType struct {
	selectorEntry *types.SelectorEntryType
	itemIndex     int
}

// selectorPress is only read and written while the UI state lock is held, so it needs no synchronization of its own.
// A zero value means the held press, if any, did not begin on a selector item.
var selectorPress selectorPressType

// selectorDoubleClickInterval holds the double click interval in nanoseconds. It is atomic since the application
// goroutine may change it while the event goroutine reads it.
var selectorDoubleClickInterval atomic.Int64

func init() {
	selectorDoubleClickInterval.Store(int64(constants.DefaultDoubleClickInterval * time.Millisecond))
}

/*
SetDoubleClickInterval is a method which allows you to set the maximum time allowed between two mouse clicks on the
same selector item for the second click to be reported as a double click by GetSelectionSource. The interval applies to
every selector and defaults to constants.DefaultDoubleClickInterval milliseconds. An error is returned, and the current
interval is kept, if the interval given is not greater than zero.

Example:

	err := Selector.SetDoubleClickInterval(400 * time.Millisecond)
*/
func (shared *selectorType) SetDoubleClickInterval(interval time.Duration) error {
	if interval <= 0 {
		return fmt.Errorf("the double click interval must be greater than zero, but %v was given", interval)
	}
	selectorDoubleClickInterval.Store(int64(interval))
	return nil
}

/*
GetDoubleClickInterval is a method which allows you to obtain the maximum time allowed between two mouse clicks on the
same selector item for the second click to be reported as a double click.

Example:

	interval := Selector.GetDoubleClickInterval()
*/
func (shared *selectorType) GetDoubleClickInterval() time.Duration {
	return time.Duration(selectorDoubleClickInterval.Load())
}

/*
getMouseSelectionSource is a method which allows you to classify a mouse press on a selector item as a single or double
click, and records the press so the next one can be compared against it. A press is a double click when it lands on
the same item of the same selector as the previous recorded press, within the double click interval. In addition, the
following should be noted:

  - A double click consumes the recorded press, so a third click in quick succession starts a new pair and is
    reported as a single click rather than as another double click.

Example:

	selectionSource := Selector.getMouseSelectionSource("Layer1", "Selector1", 3, time.Now())
*/
func (shared *selectorType) getMouseSelectionSource(layerAlias string, selectorAlias string, itemIndex int, clickTime time.Time) int {
	previousClick := selectorLastClick
	isSameItem := previousClick.layerAlias == layerAlias && previousClick.selectorAlias == selectorAlias &&
		previousClick.itemIndex == itemIndex && !previousClick.clickTime.IsZero()
	elapsedTime := clickTime.Sub(previousClick.clickTime)
	if isSameItem && elapsedTime >= 0 && elapsedTime <= shared.GetDoubleClickInterval() {
		selectorLastClick = selectorClickType{}
		return constants.SelectionSourceDoubleClick
	}
	selectorLastClick = selectorClickType{layerAlias: layerAlias, selectorAlias: selectorAlias, itemIndex: itemIndex, clickTime: clickTime}
	return constants.SelectionSourceSingleClick
}

/*
updatePressOrigin is a method which allows you to record which selector the current mouse press began on. A fresh
press over a selector item records that selector, a fresh press anywhere else records that no selector owns the press,
and a release clears the record. Movement with the button still held leaves the record unchanged. In addition, the
following should be noted:

  - It must run for every mouse event, before anything that can stop the event from reaching the controls, such as a
    modal layer rejecting it, so that a press which was never delivered to a selector cannot be mistaken for one that
    was.

  - The selector is recorded by its entry rather than its alias, so a selector deleted during the press, and one
    added later under the same alias, never inherit the press.

Example:

	Selector.updatePressOrigin(characterEntry, buttonPressed, previousButtonPressed)
*/
func (shared *selectorType) updatePressOrigin(characterEntry types.CharacterEntryType, buttonPressed uint, previousButtonPressed uint) {
	if buttonPressed == 0 {
		selectorPress = selectorPressType{}
		return
	}
	if previousButtonPressed != 0 {
		return
	}
	selectorPress = selectorPressType{}
	if characterEntry.AttributeEntry.CellType != constants.CellTypeSelectorItem {
		return
	}
	if selectorEntry, isFound := selectorRegistry.Lookup(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias); isFound {
		selectorPress = selectorPressType{selectorEntry: selectorEntry, itemIndex: constants.NullItemSelection}
	}
}

/*
IsSelectorExists is a method which checks if a selector with the specified alias exists on a given text layer. In
addition, the following should be noted:

- Returns true if the selector exists, false otherwise.

- This method is useful for validating selector existence before performing operations on it.

Example:
    exists := IsSelectorExists("Layer1", "Selector1")
*/
func IsSelectorExists(layerAlias string, selectorAlias string) bool {
	// Use the generic memory manager to check existence
	return selectorRegistry.IsExists(layerAlias, selectorAlias)
}

/*
getSelector is a method which allows you to retrieve the entry holding a selector's state from a given text layer. If
the selector does not exist, a panic is generated to fail as fast as possible. The caller must hold uiStateMutex.

Example:

	selectorEntry := getSelector("Layer1", "Selector1")
*/
func getSelector(layerAlias string, selectorAlias string) *types.SelectorEntryType {
	// Use the generic memory manager to retrieve the selector entry
	selectorEntry := selectorRegistry.Get(layerAlias, selectorAlias)
	if selectorEntry == nil {
		panic(fmt.Sprintf("The selector '%s' under layer '%s' could not be obtained since it does not exist!", selectorAlias, layerAlias))
	}
	return selectorEntry
}

// ============================================================================
// REGULAR ENTRY
// ============================================================================

/*
Delete is a method which removes a selector from a text layer. In addition, the following should be noted:

- If you attempt to delete a selector which does not exist, then the request will simply be ignored.

- All memory associated with the selector will be freed.

Example:
    selector = selector.Delete()
*/
func (shared *SelectorInstanceType) Delete() *SelectorInstanceType {
	defer lockUIState()()
	shared.BaseControlInstanceType.delete()
	return nil
}

/*
IsNewItemSelected is a method which checks if a new item has been selected in the selector. Use GetSelectionSource to
tell whether the selection came from the keyboard, a single click, or a double click.

Example:
    if selector.IsNewItemSelected() { ... }
*/
func (shared *SelectorInstanceType) IsNewItemSelected() bool {
	defer lockUIState()()
	if selectorEntry, isFound := selectorRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		return selectorEntry.IsNewItemSelected
	}
	return false
}

/*
GetSelectionSource is a method which allows you to obtain how the selector's current selection was made, so that an
application can, for example, treat Enter or a double click as activating an item while a single click only selects
it. The value returned is one of constants.SelectionSourceKeyboard, constants.SelectionSourceSingleClick,
constants.SelectionSourceDoubleClick, or constants.SelectionSourceProgrammatic for a selection made with Select. If no
selection has been made yet, or the selector does not exist, constants.SelectionSourceNone is returned. In addition, the
following should be noted:

  - The source is not cleared by GetSelected, so it may be read before or after GetSelected for the same selection.

  - The first click of a double click is reported on its own as a single click, and the second click then reports
    the same item again as a double click. An application that acts on double clicks should therefore not treat a
    single click as final.

Example:

	if selector.IsNewItemSelected() && selector.GetSelectionSource() == constants.SelectionSourceDoubleClick {
		openItem(selector.GetSelected())
	}
*/
func (shared *SelectorInstanceType) GetSelectionSource() int {
	defer lockUIState()()
	if selectorEntry, isFound := selectorRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		return selectorEntry.SelectionSource
	}
	return constants.SelectionSourceNone
}

/*
GetSelected is a method which allows you to retrieve the currently selected item of a selector, as its alias and its
zero based index. Reading the selection clears the flag reported by IsNewItemSelected. If nothing is selected, the
recorded selection is out of range, or the selector does not exist, an empty string and constants.SELECTED_NONE are
returned.

Example:

	alias, index := selector.GetSelected()
*/
func (shared *SelectorInstanceType) GetSelected() (string, int) {
	defer lockUIState()()
	if menuEntry, isFound := selectorRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		menuEntry.IsNewItemSelected = false
		value := menuEntry.ItemSelected
		if value >= 0 && value < Selector.getItemCount(menuEntry) {
			return menuEntry.SelectionEntry.SelectionAlias[value], value
		}
	}
	return "", constants.SELECTED_NONE
}

/*
GetAllItems is a method which allows you to retrieve every item of a selector, as one slice of aliases and one slice of
display values given in the order the items appear in the selector. If the selector does not exist, two empty slices
are returned. In addition, the following should be noted:

  - The slices returned are copies, so changing them does not change the selector. Use SetSelectionEntry, AddItem
    or DeleteItem to change its items.

Example:

	aliases, values := selector.GetAllItems()
*/
func (shared *SelectorInstanceType) GetAllItems() ([]string, []string) {
	defer lockUIState()()
	if menuEntry, isFound := selectorRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		selectionEntryCopy := Selector.getSelectionEntryCopy(menuEntry.SelectionEntry)
		return selectionEntryCopy.SelectionAlias, selectionEntryCopy.SelectionValue
	}
	return []string{}, []string{}
}

/*
Unselect is a method which clears the current selection for a selector. In addition, the following should be noted:

- The selector's selected item will be set to SELECTED_NONE.

- If the selector does not exist, no operation occurs.

Example:
    selector.Unselect()
*/
func (shared *SelectorInstanceType) Unselect() {
	defer lockUIState()()
	shared.unselect()
}

/*
unselect is a method which allows you to do exactly what Unselect does, for a caller that already holds uiStateMutex,
such as another entry point of the package that has taken the lock itself.

Example:

	selector.unselect()
*/
func (shared *SelectorInstanceType) unselect() {
	if selectorEntry, isFound := selectorRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		selectorEntry.ItemSelected = constants.SELECTED_NONE
	}
}

/*
Select is a method which allows you to select an item by its alias. In addition, the following should be noted:

- If an item with the matching alias is found, it will be set as the selected and highlighted item.

- If the selector or the item does not exist, no operation occurs.

Example:
    selector.Select("Option1")
*/
func (shared *SelectorInstanceType) Select(selectionAlias string) {
	defer lockUIState()()
	if selectorEntry, isFound := selectorRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {

		// Find the index of the item with the matching alias
		itemIndex := -1
		for i, alias := range selectorEntry.SelectionEntry.SelectionAlias {
			if alias == selectionAlias {
				itemIndex = i
				break
			}
		}

		// If the item exists, set it as selected.
		if itemIndex != -1 {
			selectorEntry.ItemSelected = itemIndex
			selectorEntry.ItemHighlighted = itemIndex
			selectorEntry.IsNewItemSelected = true
			selectorEntry.SelectionSource = constants.SelectionSourceProgrammatic
		}
	}
}

/*
FocusSelection is a method which allows you to scroll the item with a given alias into view, centring it in the
selector where the list allows and otherwise scrolling as far as the list goes, and moves the scroll bar to match. If
the selector or the item does not exist, no operation occurs. In addition, the following should be noted:

  - On a selector with more than one column the viewport is moved back to the start of its row, so items always
    stay in their own columns.

Example:

	selector.FocusSelection("Option5")
*/
func (shared *SelectorInstanceType) FocusSelection(selectionAlias string) {
	defer lockUIState()()
	if selectorEntry, isFound := selectorRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {

		// Find the index of the item with the matching alias
		itemIndex := -1
		for i, alias := range selectorEntry.SelectionEntry.SelectionAlias {
			if alias == selectionAlias {
				itemIndex = i
				break
			}
		}

		// If the item doesn't exist, do nothing
		if itemIndex == -1 {
			return
		}

		// Calculate the number of items visible in the viewport
		visibleRowCount := getClampedIndex(selectorEntry.Height, 1, selectorEntry.Height)
		visibleItems := visibleRowCount * Selector.getColumnCount(selectorEntry)

		// Center the item where possible. The viewport is kept in range and on a row boundary by normalizeIndexes.
		selectorEntry.ViewportPosition = itemIndex - (visibleItems / 2)
		Selector.normalizeIndexes(selectorEntry)
		Selector.updateScrollbar(shared.layerAlias, selectorEntry)
	}
}

/*
GetViewport is a method which allows you to obtain a selector's viewport position, which is the zero based index of
the item shown in the selector's top left corner. On a selector with more than one column this is always the first
item of a row. If the selector does not exist, 0 is returned.

Example:

	viewportPosition := selector.GetViewport()
*/
func (shared *SelectorInstanceType) GetViewport() int {
	defer lockUIState()()
	if selectorEntry, isFound := selectorRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		return selectorEntry.ViewportPosition
	}
	return 0
}

/*
SetViewport is a method which allows you to scroll a selector so that the item at a given zero based index is shown in
its top left corner. If the selector does not exist, no operation occurs. In addition, the following should be noted:

  - The position is clamped to the range the selector can scroll to and moved back to the start of its row, and the
    selector's scroll bar is moved to match.

  - The highlighted item is not moved, so it may be scrolled out of view.

Example:

	selector.SetViewport(10)
*/
func (shared *SelectorInstanceType) SetViewport(viewportPosition int) {
	defer lockUIState()()
	if menuEntry, isFound := selectorRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		menuEntry.ViewportPosition = viewportPosition
		Selector.normalizeIndexes(menuEntry)
		Selector.updateScrollbar(shared.layerAlias, menuEntry)
	}
}

/*
GetHighlighted is a method which allows you to obtain the zero based index of a selector's highlighted item, which is
the item the arrow keys move and Enter selects. If no item is highlighted, the recorded highlight is out of range for
the selector's items, or the selector does not exist, constants.SELECTED_NONE is returned. Reading the highlight does
not change any state of the selector.

Example:

	highlightedIndex := selector.GetHighlighted()
*/
func (shared *SelectorInstanceType) GetHighlighted() int {
	defer lockUIState()()
	if selectorEntry, isFound := selectorRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		itemHighlighted := selectorEntry.ItemHighlighted
		if itemHighlighted >= 0 && itemHighlighted < Selector.getItemCount(selectorEntry) {
			return itemHighlighted
		}
	}
	return constants.SELECTED_NONE
}

/*
SetHighlighted is a method which allows you to highlight the item at a given zero based index, as if the user had moved
to it with the arrow keys. The viewport scrolls by the fewest rows needed to show the item, and the scroll bar follows
it. A negative index clears the highlight, an index past the last item highlights the last item, and on a selector with
no items the highlight is cleared. If the selector does not exist, no operation occurs. In addition, the following
should be noted:

  - The selected item is not changed and no new selection is reported, so IsNewItemSelected is unaffected. Use
    Select to change the selection.

Example:

	selector.SetHighlighted(4)
*/
func (shared *SelectorInstanceType) SetHighlighted(itemIndex int) {
	defer lockUIState()()
	selectorEntry, isFound := selectorRegistry.Lookup(shared.layerAlias, shared.controlAlias)
	if !isFound {
		return
	}
	itemCount := Selector.getItemCount(selectorEntry)
	if itemIndex < 0 || itemCount == 0 {
		selectorEntry.ItemHighlighted = constants.NullItemSelection
	} else {
		selectorEntry.ItemHighlighted = getClampedIndex(itemIndex, 0, itemCount-1)
	}
	Selector.scrollHighlightIntoView(selectorEntry)
	Selector.normalizeIndexes(selectorEntry)
	Selector.updateScrollbar(shared.layerAlias, selectorEntry)
}

/*
SetSelectionEntry is a method which allows you to overwrite the current selection entry with a new one. The selected
item is cleared, the viewport returns to the first item, and the associated scroll bar is updated to fit the new list.
If the selector does not exist, no operation occurs. In addition, the following should be noted:

  - If the selector currently has keyboard focus and the new list is not empty, its first item is highlighted, as it
    would be when focus arrives, so the arrow keys and Enter act on a visible item. Otherwise the highlight is
    cleared.

  - The selector keeps its own copy of the item slices, so changing the slices of the entry passed in afterwards
    does not change the selector, and later changes to the selector never write into the caller's slices.

Example:

	selector.SetSelectionEntry(newSelection)
*/
func (shared *SelectorInstanceType) SetSelectionEntry(selectionEntry types.SelectionEntryType) {
	defer lockUIState()()
	isFocused := isControlCurrentlyFocused(shared.layerAlias, shared.controlAlias, constants.CellTypeSelectorItem)
	if selectorEntry, isFound := selectorRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		selectorEntry.SelectionEntry = Selector.getSelectionEntryCopy(selectionEntry)
		selectorEntry.ItemSelected = constants.SELECTED_NONE
		selectorEntry.ItemHighlighted = constants.NullItemSelection
		selectorEntry.ViewportPosition = 0
		if isFocused {
			Selector.highlightStartingItem(selectorEntry)
		}
		Selector.updateScrollbar(shared.layerAlias, selectorEntry)
	}
}

/*
SetHighlightOnClickOnly is a method which allows you to choose whether a selector highlights an item only when it is
clicked, or also when the mouse moves over it. If the selector does not exist, no operation occurs.

Example:

	selector.SetHighlightOnClickOnly(true)
*/
func (shared *SelectorInstanceType) SetHighlightOnClickOnly(isHighlightOnClickOnly bool) {
	defer lockUIState()()
	if selectorEntry, isFound := selectorRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		selectorEntry.HighlightOnClickOnly = isHighlightOnClickOnly
	}
}

/*
IsHighlightOnClickOnly is a method which allows you to check whether a selector highlights an item only when it is
clicked, rather than also when the mouse moves over it. If the selector does not exist, false is returned.

Example:

	isClickOnly := selector.IsHighlightOnClickOnly()
*/
func (shared *SelectorInstanceType) IsHighlightOnClickOnly() bool {
	defer lockUIState()()
	if selectorEntry, isFound := selectorRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		return selectorEntry.HighlightOnClickOnly
	}
	return false
}

/*
GetColumnCount is a method which allows you to obtain the number of columns a selector lays its items out in. If the
selector does not exist, 0 is returned.

Example:

	columnCount := selector.GetColumnCount()
*/
func (shared *SelectorInstanceType) GetColumnCount() int {
	defer lockUIState()()
	if selectorEntry, isFound := selectorRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		return Selector.getColumnCount(selectorEntry)
	}
	return 0
}

/*
SetColumnCount is a method which allows you to change the number of columns a selector lays its items out in. Each
column keeps its item width, so the selector becomes wider or narrower, and its border, scroll bar, tooltip and
viewport are updated to match. A count below one, or a selector that does not exist, leaves everything unchanged.

Example:

	selector.SetColumnCount(3)
*/
func (shared *SelectorInstanceType) SetColumnCount(columnCount int) {
	defer lockUIState()()
	selectorEntry, isFound := selectorRegistry.Lookup(shared.layerAlias, shared.controlAlias)
	if !isFound || columnCount < 1 {
		return
	}
	selectorEntry.NumberOfColumns = columnCount
	Selector.updateLayout(shared.layerAlias, selectorEntry)
}

/*
AddItem is a method which allows you to add a new item to the end of a selector's list of items. Both the alias and the
display value of the item must be provided. If the selector does not exist, no operation occurs. In addition, the
following should be noted:

  - The selector's scroll bar is enabled once the items overflow the visible area, and its range is extended to
    cover the new item.

Example:

	selector.AddItem("NewOpt", "New Option")
*/
func (shared *SelectorInstanceType) AddItem(selectionAlias string, selectionValue string) {
	defer lockUIState()()
	if menuEntry, isFound := selectorRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {
		menuEntry.SelectionEntry.Add(selectionAlias, selectionValue)
		Selector.normalizeIndexes(menuEntry)
		Selector.updateScrollbar(shared.layerAlias, menuEntry)
	}
}

/*
DeleteItem is a method which allows you to delete the item at a zero based index from a selector's list of items. If
the index is out of range, or the selector does not exist, no operation occurs. The highlighted and selected items are
moved back by one when they sit at or after the deleted item, so deleting the highlighted or selected item moves the
highlight or selection to the item before it. In addition, the following should be noted:

  - Deleting the last remaining item clears the highlight and the selection, and the viewport and scroll bar are
    pulled back into range when the list becomes shorter than they allow for.

  - The item slices are rebuilt rather than shifted in place, so a slice the application passed in through
    SetSelectionEntry or Add is never modified.

Example:

	selector.DeleteItem(2)
*/
func (shared *SelectorInstanceType) DeleteItem(index int) {
	defer lockUIState()()
	if menuEntry, isFound := selectorRegistry.Lookup(shared.layerAlias, shared.controlAlias); isFound {

		if index < 0 || index >= Selector.getItemCount(menuEntry) {
			return
		}

		menuEntry.SelectionEntry.SelectionAlias = slices.Concat(menuEntry.SelectionEntry.SelectionAlias[:index], menuEntry.SelectionEntry.SelectionAlias[index+1:])
		menuEntry.SelectionEntry.SelectionValue = slices.Concat(menuEntry.SelectionEntry.SelectionValue[:index], menuEntry.SelectionEntry.SelectionValue[index+1:])

		if menuEntry.ItemHighlighted >= index && menuEntry.ItemHighlighted > 0 {
			menuEntry.ItemHighlighted--
		}
		if menuEntry.ItemSelected >= index && menuEntry.ItemSelected > 0 {
			menuEntry.ItemSelected--
		}
		Selector.normalizeIndexes(menuEntry)
		Selector.updateScrollbar(shared.layerAlias, menuEntry)
	}
}

/*
Add is a method which allows you to add a selector to a given text layer. Once called, an instance of your control is
returned which will allow you to read or manipulate the properties for it. The style of the selector will be determined
by the style entry passed in. If you wish to remove a selector from a text layer, simply call DeleteSelector. In
addition, the following should be noted:

  - Selectors are not drawn physically to the text layer provided. Instead, they are rendered to the terminal at the
    same time when the text layer is rendered. This allows you to create selectors without actually overwriting the
    text layer data under it.

  - If the selector to be drawn falls outside the range of the provided layer, then only the visible portion of the
    selector will be drawn.

  - If the selector height is greater than the number of selections available, then no scroll bars are drawn.

  - A viewport position or highlighted item that does not fit the items given is corrected rather than stored: the
    viewport is clamped to the range the selector can scroll to, and an out of range highlight is cleared. The
    selector also keeps its own copy of the item slices.

Example:

	selectorInstance := Selector.Add("Layer1", "Selector1", style, selection, 0, 0, 10, 20, 1, 0, 0, false, true)
*/
func (shared *selectorType) Add(layerAlias string, selectorAlias string, styleEntry types.TuiStyleEntryType, selectionEntry types.SelectionEntryType, xLocation int, yLocation int, selectorHeight int, itemWidth int, numberOfColumns int, viewportPosition int, selectedItem int, highlightOnClickOnly bool, isBorderDrawn bool) SelectorInstanceType {
	defer lockUIState()()
	return shared.add(layerAlias, selectorAlias, styleEntry, selectionEntry, xLocation, yLocation, selectorHeight, itemWidth, numberOfColumns, viewportPosition, selectedItem, highlightOnClickOnly, isBorderDrawn)
}

/*
add is a method which allows you to do exactly what Add does, for a caller that already holds uiStateMutex, such as
another entry point of the package that has taken the lock itself.

Example:

	Selector.add(layerAlias, selectorAlias, styleEntry, selectionEntry, xLocation, yLocation, selectorHeight, itemWidth, numberOfColumns, viewportPosition, selectedItem, highlightOnClickOnly, isBorderDrawn)
*/
func (shared *selectorType) add(layerAlias string, selectorAlias string, styleEntry types.TuiStyleEntryType, selectionEntry types.SelectionEntryType, xLocation int, yLocation int, selectorHeight int, itemWidth int, numberOfColumns int, viewportPosition int, selectedItem int, highlightOnClickOnly bool, isBorderDrawn bool) SelectorInstanceType {
	newSelectorEntry := types.NewSelectorEntry()
	newSelectorEntry.Alias = selectorAlias
	newSelectorEntry.StyleEntry = styleEntry
	newSelectorEntry.SelectionEntry = shared.getSelectionEntryCopy(selectionEntry)
	newSelectorEntry.XLocation = xLocation
	newSelectorEntry.YLocation = yLocation
	newSelectorEntry.Height = selectorHeight
	newSelectorEntry.ItemWidth = itemWidth
	newSelectorEntry.NumberOfColumns = numberOfColumns
	newSelectorEntry.HighlightOnClickOnly = highlightOnClickOnly
	newSelectorEntry.ViewportPosition = viewportPosition
	newSelectorEntry.ItemHighlighted = selectedItem
	newSelectorEntry.IsBorderDrawn = isBorderDrawn
	newSelectorEntry.IsVisible = true
	shared.normalizeIndexes(&newSelectorEntry)
	selectorEntry := &newSelectorEntry
	selectorEntry.Width = shared.getSelectorWidth(selectorEntry)

	hotspotXLocation, hotspotYLocation, hotspotWidth, tooltipYLocation := shared.getTooltipGeometry(selectorEntry)
	tooltipInstance := Tooltip.add(layerAlias, selectorEntry.TooltipAlias, "", styleEntry,
		hotspotXLocation, hotspotYLocation, hotspotWidth, 1,
		hotspotXLocation, tooltipYLocation, hotspotWidth, 3,
		false, true, constants.DefaultTooltipHoverTime)
	tooltipInstance.setEnabled(false)
	tooltipInstance.setParentControlAlias(selectorAlias)
	selectorEntry.ScrollbarAlias = stringformat.GetLastSortedUUID()

	// The selector entry is complete before it is added to selectorRegistry, so nothing can find it half set up.
	// scrollbar.Add stores one less than the maximum value it is given, so one is added to the largest viewport position.
	scrollBarXLocation, scrollBarYLocation, scrollBarLength := shared.getScrollbarGeometry(selectorEntry)
	scrollbar.Add(layerAlias, selectorEntry.ScrollbarAlias, styleEntry, scrollBarXLocation, scrollBarYLocation, scrollBarLength,
		shared.getMaxViewportPosition(selectorEntry)+1, 0, shared.getColumnCount(selectorEntry), false)
	if scrollBarEntry, isFound := scrollbarRegistry.Lookup(layerAlias, selectorEntry.ScrollbarAlias); isFound {
		scrollBarEntry.ParentControlAlias = selectorAlias
		scrollBarEntry.ParentControlType = constants.CellTypeSelectorItem
	}
	shared.updateScrollbar(layerAlias, selectorEntry)
	selectorRegistry.Add(layerAlias, selectorAlias, selectorEntry)
	var selectorInstance SelectorInstanceType
	selectorInstance.layerAlias = layerAlias
	selectorInstance.controlAlias = selectorAlias
	selectorInstance.controlType = constants.TYPE_SELECTOR
	return selectorInstance
}

/*
Delete is a method which allows you to remove a selector from a text layer. In addition, the following should be noted:

- If you attempt to delete a selector which does not exist, then the request will simply be ignored.

- All memory associated with the selector will be freed.

Example:
    Selector.Delete("Layer1", "Selector1")
*/
func (shared *selectorType) Delete(layerAlias string, selectorAlias string) {
	defer lockUIState()()
	shared.delete(layerAlias, selectorAlias)
}

/*
delete is a method which allows you to do exactly what Delete does, for a caller that already holds uiStateMutex, such
as another entry point of the package that has taken the lock itself.

Example:

	Selector.delete(layerAlias, selectorAlias)
*/
func (shared *selectorType) delete(layerAlias string, selectorAlias string) {
	selectorRegistry.Remove(layerAlias, selectorAlias)
	clearStaleControlReferences(layerAlias, selectorAlias, constants.CellTypeSelectorItem)
}

/*
DeleteAll is a method which allows you to remove all selectors from a text layer. In addition, the following should be
noted:

- This operation cannot be undone.

- All memory associated with the selectors will be freed.

Example:
    Selector.DeleteAll("Layer1")
*/
func (shared *selectorType) DeleteAll(layerAlias string) {
	defer lockUIState()()
	removeAllControlsAndClearCells(selectorRegistry, layerAlias, constants.CellTypeSelectorItem, func(entry *types.SelectorEntryType) string { return entry.Alias })
}

/*
setupSelectorAttributes is a method which allows you to create and configure the standard and highlight attribute
entries for a selector based on the provided style entry. In addition, the following should be noted:

- Returns two attribute entries: one for normal menu items and one for highlighted items.

- The attributes are configured based on the colors specified in the style entry.

- These attributes control the visual appearance of selector items when drawn.

Example:
    attr, highAttr := Selector.setupSelectorAttributes(style)
*/
func (shared *selectorType) setupSelectorAttributes(styleEntry types.TuiStyleEntryType) (types.AttributeEntryType, types.AttributeEntryType) {
	menuAttributeEntry := types.NewAttributeEntry()
	menuAttributeEntry.ForegroundColor = styleEntry.Selector.ForegroundColor
	menuAttributeEntry.BackgroundColor = styleEntry.Selector.BackgroundColor

	highlightAttributeEntry := types.NewAttributeEntry()
	highlightAttributeEntry.ForegroundColor = styleEntry.Selector.HighlightForegroundColor
	highlightAttributeEntry.BackgroundColor = styleEntry.Selector.HighlightBackgroundColor

	return menuAttributeEntry, highlightAttributeEntry
}

/*
drawSelectorBorder is a method which allows you to draw a border around a selector on the specified text layer. The
border is drawn one cell outside the selector area, so it spans selectorWidth plus two columns and selectorHeight plus
two rows, and the area it encloses is first filled with spaces using the provided attribute entry. If IsShadowDrawn is
enabled, the border is drawn with a shadow using drawWindow, otherwise it is drawn using drawBorder. In addition, the
following should be noted:

  - The selectorWidth is the width of the whole selector, which is its item width multiplied by its number of
    columns, not the width of a single item.

Example:

	Selector.drawSelectorBorder(&layerEntry, style, attr, 0, 0, 20, 10)
*/
func (shared *selectorType) drawSelectorBorder(layerEntry *types.LayerEntryType, styleEntry types.TuiStyleEntryType,
	attributeEntry types.AttributeEntryType, xLocation int, yLocation int, selectorWidth int, selectorHeight int) {
	fillArea(layerEntry, attributeEntry, " ", xLocation-1, yLocation-1, selectorWidth+2, selectorHeight+2, constants.NullCellControlLocation)

	if styleEntry.Selector.IsShadowDrawn {
		drawWindow(layerEntry, styleEntry, attributeEntry, xLocation-1, yLocation-1, selectorWidth+2, selectorHeight+2, false)
	} else {
		drawBorder(layerEntry, styleEntry, attributeEntry, xLocation-1, yLocation-1, selectorWidth+2, selectorHeight+2, false)
	}
}

/*
formatSelectorItemText is a method which allows you to format the text for a selector item based on the provided style
and whether the item is highlighted. In addition, the following should be noted:

- Handles text alignment according to the style entry's text alignment setting.

- For centered text, special handling is applied to ensure proper centering with markup.

- For highlighted items, markup tags are stripped to ensure consistent highlighting.

- Returns a string formatted to the specified width with appropriate padding.

Example:
    formattedText := Selector.formatSelectorItemText("Option 1", 20, style, false)
*/
func (shared *selectorType) formatSelectorItemText(menuItemText string, itemWidth int, styleEntry types.TuiStyleEntryType, isHighlighted bool) string {
	var menuItemName string

	if styleEntry.Selector.IsSelectionCentered {
		// If centered, we need to handle markup specially
		if isHighlighted {
			// For highlighted items, strip markup tags to ensure they don't apply
			menuItemText = GetNonMarkupText(menuItemText)
			menuItemName = stringformat.GetFormattedString(menuItemText, itemWidth, constants.AlignmentCenter)
		} else {
			// For non-highlighted items, calculate centering that accounts for markup
			textLength := CalculateStringLengthWithoutMarkup(menuItemText)
			padding := (itemWidth - textLength) / 2
			if padding < 0 {
				padding = 0
			}

			// Create padding
			leftPadding := stringformat.GetFilledString(padding, " ")
			rightPadding := stringformat.GetFilledString(itemWidth-textLength-padding, " ")

			// Combine with original text (preserving markup)
			menuItemName = leftPadding + menuItemText
			if len(menuItemName) < itemWidth {
				menuItemName += rightPadding
			}
		}
	} else {
		// Use standard formatting for non-centered items
		if isHighlighted {
			// For highlighted items, strip markup tags
			menuItemText = GetNonMarkupText(menuItemText)
		}
		menuItemName = stringformat.GetFormattedString(menuItemText, itemWidth, styleEntry.Selector.TextAlignment)
	}

	return menuItemName
}

/*
drawSelectorItem is a method which allows you to draw a single selector item on the specified text layer. In addition,
the following should be noted:

- For non-highlighted items, markup processing is applied to support styled text.

- For highlighted items, standard printing is used with the highlight attributes.

- The cell control ID and alias are set to enable mouse and keyboard interaction.

- Returns the width of the drawn item, which may vary based on the content.

Example:
    width := Selector.drawSelectorItem(&layerEntry, attr, "Option 1", 0, 0, 0, false)
*/
func (shared *selectorType) drawSelectorItem(layerEntry *types.LayerEntryType, attributeEntry types.AttributeEntryType,
	menuItemName string, xLocation int, currentXOffset int, currentYLocation int, isHighlighted bool) int {

	arrayOfRunes := stringformat.GetRunesFromString(menuItemName)

	// Use printMarkup for non-highlighted items, otherwise use printLayer
	if !isHighlighted {
		layer.printMarkup(layerEntry, attributeEntry, xLocation+(currentXOffset), currentYLocation, 0, menuItemName)
	} else {
		layer.printLayer(layerEntry, attributeEntry, xLocation+(currentXOffset), currentYLocation, arrayOfRunes)
	}

	return stringformat.GetWidthOfRunesWhenPrinted(arrayOfRunes)
}

/*
drawSelector is a method which allows you to draw a selector on a given text layer, from the state held in its entry.
The style of the selector is taken from its entry. Selectors are not drawn physically to the text layer provided.
Instead, they are rendered to the terminal at the same time when the text layer is rendered, and if the selector falls
outside the range of the layer, only its visible portion is drawn. A focused selector shows its focus on its border,
drawn in the focused colours, when isFocusShown is true. In addition, the following should be noted:

  - Drawing never indexes outside the item list, whatever state it is given: a negative viewport position is drawn
    from the first item, and a highlighted item that is out of range is drawn as no highlight. Since rendering runs
    on the event goroutine, a panic here could not be recovered by the application.

  - The caller must hold uiStateMutex.

Example:

	Selector.drawSelector(&layerEntry, selectorEntry, false)
*/
func (shared *selectorType) drawSelector(layerEntry *types.LayerEntryType, selectorEntry *types.SelectorEntryType, isFocusShown bool) {
	if !selectorEntry.IsVisible {
		return
	}
	selectorAlias := selectorEntry.Alias
	styleEntry := selectorEntry.StyleEntry
	selectionEntry := selectorEntry.SelectionEntry
	xLocation := selectorEntry.XLocation
	yLocation := selectorEntry.YLocation
	selectorHeight := selectorEntry.Height
	itemWidth := selectorEntry.ItemWidth
	columnCount := shared.getColumnCount(selectorEntry)
	itemHighlighted := selectorEntry.ItemHighlighted

	menuAttributeEntry, highlightAttributeEntry := shared.setupSelectorAttributes(styleEntry)

	if selectorEntry.IsBorderDrawn {
		borderStyleEntry := styleEntry
		// A focused selector shows its focus on its border, drawn in the focused colours while keyboard focus is shown.
		if isFocusShown {
			borderStyleEntry.Window.LineDrawingTextForegroundColor, borderStyleEntry.Window.LineDrawingTextBackgroundColor = getFocusedColors(
				styleEntry.Window.LineDrawingTextForegroundColor, styleEntry.Window.LineDrawingTextBackgroundColor,
				styleEntry.Selector.FocusedForegroundColor, styleEntry.Selector.FocusedBackgroundColor)
		}
		// Border cells, and the empty cells inside it, are tagged with the selector so the mouse wheel can scroll it there.
		borderAttributeEntry := menuAttributeEntry
		borderAttributeEntry.CellType = constants.CellTypeSelectorBorder
		borderAttributeEntry.CellControlAlias = selectorAlias
		shared.drawSelectorBorder(layerEntry, borderStyleEntry, borderAttributeEntry, xLocation, yLocation, shared.getSelectorWidth(selectorEntry), selectorHeight)
	}

	currentYLocation := yLocation
	currentMenuItemIndex := getClampedIndex(selectorEntry.ViewportPosition, 0, selectorEntry.ViewportPosition)
	currentXOffset := 0
	currentColumn := 0
	currentRow := 0
	for currentMenuItemIndex < len(selectionEntry.SelectionValue) && currentRow < selectorHeight {
		isHighlighted := currentMenuItemIndex == itemHighlighted
		attributeEntry := menuAttributeEntry
		if isHighlighted {
			attributeEntry = highlightAttributeEntry
		}
		menuItemName := shared.formatSelectorItemText(selectionEntry.SelectionValue[currentMenuItemIndex], itemWidth, styleEntry, isHighlighted)
		attributeEntry.CellControlId = currentMenuItemIndex
		attributeEntry.CellControlAlias = selectorAlias
		attributeEntry.CellType = constants.CellTypeSelectorItem
		printedItemWidth := shared.drawSelectorItem(layerEntry, attributeEntry, menuItemName, xLocation, currentXOffset, currentYLocation, isHighlighted)
		currentMenuItemIndex++
		currentXOffset = currentXOffset + printedItemWidth
		currentColumn++
		if currentColumn >= columnCount {
			currentXOffset = 0
			currentColumn = 0
			currentYLocation++
			currentRow++
		}
	}
	scrollbar.drawOnLayerByAlias(layerEntry, selectorEntry.ScrollbarAlias)
}

/*
drawSelectorsOnLayer is a method which allows you to draw all selectors on a given text layer. In addition, the
following should be noted:

  - Selectors are drawn in alphabetical order by their alias, which keeps the order the same across frames, and
    internally generated selectors, such as dropdown trays, are drawn last.

  - The caller must hold uiStateMutex, as UpdateDisplay does for the whole redraw.

Example:

	Selector.drawSelectorsOnLayer(layerEntry)
*/
func (shared *selectorType) drawSelectorsOnLayer(layerEntry types.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	compareByAlias := func(a, b *types.SelectorEntryType) bool {
		return a.Alias < b.Alias
	}
	// Sort array so internally generated selectors appear last (Since sorted by name, and
	// UUID generates "zzz" prefixes). This prevents Dropdown selectors from appearing under
	// user created selectors, when they should always be on top.
	for _, currentSelectorEntry := range selectorRegistry.SortEntries(layerAlias, true, compareByAlias) {
		selectorEntry := currentSelectorEntry
		isFocusShown := isFocusIndicatorShown(layerAlias, selectorEntry.Alias, constants.CellTypeSelectorItem)
		shared.drawSelector(&layerEntry, selectorEntry, isFocusShown)
	}
}

/*
getSelectionEntryCopy is a method which allows you to obtain a copy of a selection entry whose alias and value slices
share no memory with the original. Selectors keep such a copy so that the application's slices and the selector's own
slices can each be changed without affecting the other.

Example:

	selectionEntryCopy := Selector.getSelectionEntryCopy(selectionEntry)
*/
func (shared *selectorType) getSelectionEntryCopy(selectionEntry types.SelectionEntryType) types.SelectionEntryType {
	return types.SelectionEntryType{
		SelectionAlias: slices.Clone(selectionEntry.SelectionAlias),
		SelectionValue: slices.Clone(selectionEntry.SelectionValue),
	}
}

/*
getItemCount is a method which allows you to obtain the number of items of a selector that can be safely indexed. It is
the length of the shorter of the alias and value slices, so any index below it is valid for both even if an
application has edited them to different lengths.

Example:

	itemCount := Selector.getItemCount(selectorEntry)
*/
func (shared *selectorType) getItemCount(selectorEntry *types.SelectorEntryType) int {
	return getClampedIndex(len(selectorEntry.SelectionEntry.SelectionAlias), 0, len(selectorEntry.SelectionEntry.SelectionValue))
}

/*
getColumnCount is a method which allows you to obtain the number of columns a selector lays its items out in, treating
a column count below one as a single column so that row and column arithmetic never divides by zero.

Example:

	columnCount := Selector.getColumnCount(selectorEntry)
*/
func (shared *selectorType) getColumnCount(selectorEntry *types.SelectorEntryType) int {
	return getClampedIndex(selectorEntry.NumberOfColumns, 1, selectorEntry.NumberOfColumns)
}

/*
getSelectorWidth is a method which allows you to obtain the width of a selector's item area, which is its item width
multiplied by its number of columns, treating a column count below one as a single column. The border, when drawn, and
a scroll bar placed outside it are not included, in the same way that a selector's height counts only its item rows.

Example:

	selectorWidth := Selector.getSelectorWidth(selectorEntry)
*/
func (shared *selectorType) getSelectorWidth(selectorEntry *types.SelectorEntryType) int {
	return selectorEntry.ItemWidth * shared.getColumnCount(selectorEntry)
}

/*
getScrollbarGeometry is a method which allows you to obtain where a selector's scroll bar belongs, as its column, its
top row and its length. Without a border, the scroll bar covers the last column of the item area and is as tall as the
item rows. With a border, it sits one column to the right of the right border, starting on the top border row and
spanning the full height of the border.

Example:

	xLocation, yLocation, length := Selector.getScrollbarGeometry(selectorEntry)
*/
func (shared *selectorType) getScrollbarGeometry(selectorEntry *types.SelectorEntryType) (int, int, int) {
	selectorWidth := shared.getSelectorWidth(selectorEntry)
	if selectorEntry.IsBorderDrawn {
		return selectorEntry.XLocation + selectorWidth + 1, selectorEntry.YLocation - 1, selectorEntry.Height + 2
	}
	return selectorEntry.XLocation + selectorWidth - 1, selectorEntry.YLocation, selectorEntry.Height
}

/*
getTooltipGeometry is a method which allows you to obtain where a selector's tooltip belongs, as the column and row of
its one row hotspot, the width shared by the hotspot and the tooltip, and the row the three row tooltip is drawn
from, which is the row below the hotspot. Both are as wide as the item area plus two columns.

Example:

	hotspotXLocation, hotspotYLocation, width, tooltipYLocation := Selector.getTooltipGeometry(selectorEntry)
*/
func (shared *selectorType) getTooltipGeometry(selectorEntry *types.SelectorEntryType) (int, int, int, int) {
	return selectorEntry.XLocation, selectorEntry.YLocation, shared.getSelectorWidth(selectorEntry) + 2, selectorEntry.YLocation + 1
}

/*
updateLayout is a method which allows you to bring everything derived from a selector's position, size, column count
and border back in line with them, and is called whenever any of those change after the selector is added. The
selector's width is recomputed, its scroll bar and tooltip are moved and resized to match, its viewport and item
indexes are brought back into range for the number of rows now visible, and its scroll bar is updated to match.
The caller must hold uiStateMutex.

Example:

	Selector.updateLayout("Layer1", selectorEntry)
*/
func (shared *selectorType) updateLayout(layerAlias string, selectorEntry *types.SelectorEntryType) {
	selectorEntry.Width = shared.getSelectorWidth(selectorEntry)
	if scrollBarEntry, isFound := scrollbarRegistry.Lookup(layerAlias, selectorEntry.ScrollbarAlias); isFound {
		scrollBarEntry.XLocation, scrollBarEntry.YLocation, scrollBarEntry.Length = shared.getScrollbarGeometry(selectorEntry)
	}
	if tooltipEntry, isFound := tooltipRegistry.Lookup(layerAlias, selectorEntry.TooltipAlias); isFound {
		hotspotXLocation, hotspotYLocation, width, tooltipYLocation := shared.getTooltipGeometry(selectorEntry)
		tooltipEntry.HotspotXLocation, tooltipEntry.HotspotYLocation, tooltipEntry.HotspotWidth = hotspotXLocation, hotspotYLocation, width
		tooltipEntry.TooltipXLocation, tooltipEntry.TooltipYLocation, tooltipEntry.TooltipWidth = hotspotXLocation, tooltipYLocation, width
	}
	shared.normalizeIndexes(selectorEntry)
	shared.updateScrollbar(layerAlias, selectorEntry)
}

/*
setSize is a method which allows you to resize a selector's item area to a given width and height, where the height is
the number of item rows and the width is shared equally between its columns. The item width becomes the width divided
by the number of columns, rounded down, so the selector's width is the largest multiple of its column count that does
not exceed the width given. The scroll bar, tooltip and viewport are then updated to match. A height below one, a
width too narrow to give every column at least one cell, or a selector that does not exist leaves everything
unchanged. The caller must hold uiStateMutex.

Example:

	Selector.setSize("Layer1", "Selector1", 30, 5)
*/
func (shared *selectorType) setSize(layerAlias string, selectorAlias string, width int, height int) {
	selectorEntry, isFound := selectorRegistry.Lookup(layerAlias, selectorAlias)
	if !isFound {
		return
	}
	columnCount := shared.getColumnCount(selectorEntry)
	if height < 1 || width < columnCount {
		return
	}
	selectorEntry.ItemWidth = width / columnCount
	selectorEntry.Height = height
	shared.updateLayout(layerAlias, selectorEntry)
}

/*
getMaxViewportPosition is a method which allows you to obtain the largest viewport position a selector can scroll to.
This is the index of the first item on the row that puts the last row of items at the bottom of the selector, so it is
always the start of a row, and it is zero when every item fits. It is also the maximum value of the selector's scroll
bar.

Example:

	maxViewportPosition := Selector.getMaxViewportPosition(selectorEntry)
*/
func (shared *selectorType) getMaxViewportPosition(selectorEntry *types.SelectorEntryType) int {
	columnCount := shared.getColumnCount(selectorEntry)
	rowCount := (shared.getItemCount(selectorEntry) + columnCount - 1) / columnCount
	visibleRowCount := getClampedIndex(selectorEntry.Height, 1, selectorEntry.Height)
	maxViewportPosition := (rowCount - visibleRowCount) * columnCount
	return getClampedIndex(maxViewportPosition, 0, maxViewportPosition)
}

/*
normalizeIndexes is a method which allows you to bring a selector's highlighted item, selected item and viewport
position back in line with its current list of items, so that none of them can be used to index outside the list. It
is called by everything that changes the list or those indexes. A highlighted or selected item that is out of range is
cleared, and the viewport position is clamped between zero and getMaxViewportPosition and moved back to the start of
its row. The caller must hold uiStateMutex.

Example:

	Selector.normalizeIndexes(selectorEntry)
*/
func (shared *selectorType) normalizeIndexes(selectorEntry *types.SelectorEntryType) {
	itemCount := shared.getItemCount(selectorEntry)
	if selectorEntry.ItemHighlighted < 0 || selectorEntry.ItemHighlighted >= itemCount {
		selectorEntry.ItemHighlighted = constants.NullItemSelection
	}
	if selectorEntry.ItemSelected < 0 || selectorEntry.ItemSelected >= itemCount {
		selectorEntry.ItemSelected = constants.SELECTED_NONE
	}
	columnCount := shared.getColumnCount(selectorEntry)
	viewportPosition := getClampedIndex(selectorEntry.ViewportPosition, 0, shared.getMaxViewportPosition(selectorEntry))
	selectorEntry.ViewportPosition = viewportPosition - viewportPosition%columnCount
}

/*
highlightStartingItem is a method which allows you to make sure a selector with items has one of them highlighted, as
happens when the selector gains keyboard focus. If the highlighted item is missing or out of range, the selected item
is highlighted when it is valid, and otherwise the first item is. An empty selector is left with no highlight, since
there is nothing to highlight. The caller must hold uiStateMutex.

Example:

	Selector.highlightStartingItem(selectorEntry)
*/
func (shared *selectorType) highlightStartingItem(selectorEntry *types.SelectorEntryType) {
	shared.normalizeIndexes(selectorEntry)
	if selectorEntry.ItemHighlighted != constants.NullItemSelection || shared.getItemCount(selectorEntry) == 0 {
		return
	}
	selectorEntry.ItemHighlighted = 0
	if selectorEntry.ItemSelected != constants.SELECTED_NONE {
		selectorEntry.ItemHighlighted = selectorEntry.ItemSelected
	}
}

/*
scrollHighlightIntoView is a method which allows you to move a selector's viewport by the fewest rows needed for its
highlighted item to be visible. Nothing is moved when no item is highlighted. The caller must hold uiStateMutex.

Example:

	Selector.scrollHighlightIntoView(selectorEntry)
*/
func (shared *selectorType) scrollHighlightIntoView(selectorEntry *types.SelectorEntryType) {
	if selectorEntry.ItemHighlighted < 0 {
		return
	}
	columnCount := shared.getColumnCount(selectorEntry)
	visibleRowCount := getClampedIndex(selectorEntry.Height, 1, selectorEntry.Height)
	highlightedRowStart := selectorEntry.ItemHighlighted - selectorEntry.ItemHighlighted%columnCount
	if highlightedRowStart < selectorEntry.ViewportPosition {
		selectorEntry.ViewportPosition = highlightedRowStart
	} else if highlightedRowStart >= selectorEntry.ViewportPosition+visibleRowCount*columnCount {
		selectorEntry.ViewportPosition = highlightedRowStart - (visibleRowCount-1)*columnCount
	}
}

/*
updateScrollbar is a method which allows you to bring a selector's scroll bar in line with its items and viewport. The
scroll bar is enabled only while the items overflow the selector, is shown only while it is enabled and the selector
itself is visible, ranges from zero to getMaxViewportPosition, and has its value and handle moved to the selector's
viewport position. Nothing happens if the selector has no scroll bar. The caller must hold uiStateMutex.

Example:

	Selector.updateScrollbar("Layer1", selectorEntry)
*/
func (shared *selectorType) updateScrollbar(layerAlias string, selectorEntry *types.SelectorEntryType) {
	scrollBarEntry, isFound := scrollbarRegistry.Lookup(layerAlias, selectorEntry.ScrollbarAlias)
	if !isFound {
		return
	}
	visibleRowCount := getClampedIndex(selectorEntry.Height, 1, selectorEntry.Height)
	isOverflowing := shared.getItemCount(selectorEntry) > visibleRowCount*shared.getColumnCount(selectorEntry) &&
		selectorEntry.StyleEntry.Selector.TextAlignment != constants.AlignmentNoPadding
	scrollBarEntry.IsEnabled = isOverflowing
	scrollBarEntry.IsVisible = isOverflowing && selectorEntry.IsVisible
	scrollBarEntry.MaxScrollValue = shared.getMaxViewportPosition(selectorEntry)
	scrollBarEntry.ScrollValue = selectorEntry.ViewportPosition
	scrollbar.computeHandlePositionByScrollValue(layerAlias, selectorEntry.ScrollbarAlias)
}

/*
updateKeyboardEventForSelector is a method which allows you to process a keystroke for a specific selector. Up and down
move the highlight by one row, left and right move it by one column within its row, and Enter selects the highlighted
item. The viewport scrolls to keep the highlighted item visible and the scroll bar follows it. The first result reports
whether the screen needs updating, and the second whether the keystroke was consumed. In addition, the following should
be noted:

  - On a selector with no items, nothing is changed and Enter selects nothing. The arrow keys are still consumed, so
    they do not reach the application's keyboard buffer as if no control had focus, while Enter is left unconsumed
    so the layer's default button can still act on it.

  - If no valid item is highlighted when a key arrives, for example because the list was refilled, the selected item
    is highlighted first, or the first item when nothing is selected, and the key is then applied to it.

  - Left on the first column and right on the last column are not consumed, so they reach the keyboard buffer as
    before. On a single column selector this is true of every left and right press.

Example:

	isUpdateRequired, isConsumed := Selector.updateKeyboardEventForSelector("Layer1", "Sel1", []rune("down"))
*/
func (shared *selectorType) updateKeyboardEventForSelector(layerAlias string, selectorAlias string, keystroke []rune) (bool, bool) {
	keystrokeAsString := string(keystroke)
	isArrowKey := keystrokeAsString == "up" || keystrokeAsString == "down" || keystrokeAsString == "left" || keystrokeAsString == "right"
	if !isArrowKey && keystrokeAsString != "enter" {
		return false, false
	}
	selectorEntry, isFound := selectorRegistry.Lookup(layerAlias, selectorAlias)
	if !isFound {
		return false, false
	}
	itemCount := shared.getItemCount(selectorEntry)
	if itemCount == 0 {
		return false, isArrowKey
	}
	shared.highlightStartingItem(selectorEntry)
	columnCount := shared.getColumnCount(selectorEntry)
	isKeystrokeConsumed := false
	switch keystrokeAsString {
	case "down":
		if selectorEntry.ItemHighlighted+columnCount < itemCount {
			selectorEntry.ItemHighlighted += columnCount
		}
		isKeystrokeConsumed = true
	case "up":
		if selectorEntry.ItemHighlighted-columnCount >= 0 {
			selectorEntry.ItemHighlighted -= columnCount
		}
		isKeystrokeConsumed = true
	case "left":
		if selectorEntry.ItemHighlighted%columnCount != 0 {
			selectorEntry.ItemHighlighted--
			isKeystrokeConsumed = true
		}
	case "right":
		if selectorEntry.ItemHighlighted%columnCount != columnCount-1 {
			if selectorEntry.ItemHighlighted+1 < itemCount {
				selectorEntry.ItemHighlighted++
			}
			isKeystrokeConsumed = true
		}
	case "enter":
		selectorEntry.ItemSelected = selectorEntry.ItemHighlighted
		selectorEntry.IsNewItemSelected = true
		selectorEntry.SelectionSource = constants.SelectionSourceKeyboard
		isKeystrokeConsumed = true
	}
	shared.scrollHighlightIntoView(selectorEntry)
	shared.normalizeIndexes(selectorEntry)
	shared.updateScrollbar(layerAlias, selectorEntry)
	return true, isKeystrokeConsumed
}

/*
updateKeyboardEvent is a method which allows you to update the state of all selectors according to the current keystroke
event. In addition, the following should be noted:

- Handles navigation keys (up, down, left, right) to move between items.

- Enter key selects the currently highlighted item.

- Returns true if the screen needs to be updated due to state changes.

Example:
    updateRequired, consumed := Selector.updateKeyboardEvent(rune("up"))
*/
func (shared *selectorType) updateKeyboardEvent(keystroke []rune) (bool, bool) {
	isScreenUpdateRequired := false
	isKeystrokeConsumed := false
	focusedControl := getFocusedControl()
	if focusedControl.controlType != constants.CellTypeSelectorItem || !selectorRegistry.IsExists(focusedControl.layerAlias, focusedControl.controlAlias) {
		return isScreenUpdateRequired, isKeystrokeConsumed
	}
	return shared.updateKeyboardEventForSelector(focusedControl.layerAlias, focusedControl.controlAlias, keystroke)
}

/*
updateWheelEvent is a method which allows you to scroll a visible selector by one row for each notch of the mouse
wheel, when the cursor is over one of its items, its border, the empty rows inside its border, or its scroll bar. This
is what scrolls an open dropdown tray, since the tray is a selector. Wheel up scrolls toward the first item and wheel down toward the last, the viewport stays
within the range the selector can scroll to, and the scroll bar follows the viewport. It returns true only if the
viewport moved. In addition, the following should be noted:

  - When the cursor is over an item and the selector highlights on hover, the highlight moves to the item that has
    scrolled under the cursor, so it does not stay on an item that may have scrolled out of view.

  - Nothing is scrolled while a scroll bar is being dragged, so the wheel cannot fight the drag over the viewport.

Example:

	isUpdateRequired := Selector.updateWheelEvent(characterEntry, "Down")
*/
func (shared *selectorType) updateWheelEvent(characterEntry types.CharacterEntryType, wheelState string) bool {
	if (wheelState != "Up" && wheelState != "Down") || getEventStateId() != constants.EventStateNone {
		return false
	}
	layerAlias := characterEntry.LayerAlias
	selectorAlias := characterEntry.AttributeEntry.CellControlAlias
	isOverItem := characterEntry.AttributeEntry.CellType == constants.CellTypeSelectorItem
	switch characterEntry.AttributeEntry.CellType {
	case constants.CellTypeSelectorItem, constants.CellTypeSelectorBorder:
	case constants.CellTypeScrollbar:
		scrollBarEntry, isFound := scrollbarRegistry.Lookup(layerAlias, selectorAlias)
		if !isFound || scrollBarEntry.ParentControlType != constants.CellTypeSelectorItem {
			return false
		}
		selectorAlias = scrollBarEntry.ParentControlAlias
	default:
		return false
	}
	selectorEntry, isFound := selectorRegistry.Lookup(layerAlias, selectorAlias)
	if !isFound || !selectorEntry.IsVisible {
		return false
	}
	columnCount := shared.getColumnCount(selectorEntry)
	previousViewportPosition := selectorEntry.ViewportPosition
	if wheelState == "Up" {
		selectorEntry.ViewportPosition -= columnCount
	} else {
		selectorEntry.ViewportPosition += columnCount
	}
	shared.normalizeIndexes(selectorEntry)
	viewportShift := selectorEntry.ViewportPosition - previousViewportPosition
	if viewportShift == 0 {
		return false
	}
	if isOverItem && !selectorEntry.HighlightOnClickOnly {
		hoveredItemIndex := characterEntry.AttributeEntry.CellControlId + viewportShift
		if hoveredItemIndex >= 0 && hoveredItemIndex < shared.getItemCount(selectorEntry) {
			selectorEntry.ItemHighlighted = hoveredItemIndex
		}
	}
	shared.updateScrollbar(layerAlias, selectorEntry)
	return true
}

/*
updateMouseEvent is a method which allows you to update the state of all selectors according to the current mouse event
state. Pressing on an item selects and highlights it, hovering highlights it unless the selector only highlights on a
click, the mouse wheel scrolls the selector under the cursor, and dragging a scroll bar moves the viewport of the
selector it belongs to. It returns true if the screen needs to be updated. In addition, the following should be noted:

  - A press only selects items in the selector it began on, as recorded by updatePressOrigin. A held button dragged
    onto a selector from anywhere else changes neither its selection nor its highlight.

  - While the button is held, a new selection is only reported for the press itself and for each move onto a
    different item, so further events over the same item never report the same single or double click again.

  - The item under the cursor is taken from the last drawn screen, so a click or hover on an item that no longer
    exists, because the list was shortened and not yet redrawn, is ignored rather than stored.

  - A viewport position copied from a scroll bar is clamped into the range the selector can scroll to.

Example:

	isUpdateRequired := Selector.updateMouseEvent()
*/
func (shared *selectorType) updateMouseEvent() bool {
	isScreenUpdateRequired := false
	focusedLayerAlias := getFocusedControl().layerAlias
	var characterEntry types.CharacterEntryType
	mouseXLocation, mouseYLocation, buttonPressed, wheelState := GetMouseStatus()
	characterEntry = getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
	if shared.updateWheelEvent(characterEntry, wheelState) {
		return true
	}
	selectorEntry, isFound := selectorRegistry.Lookup(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias)
	isOverSelectorItem := isFound && characterEntry.AttributeEntry.CellType == constants.CellTypeSelectorItem &&
		getEventStateId() == constants.EventStateNone
	// A press belongs to the control it began on, so dragging it over this selector neither selects nor highlights.
	isPressFromElsewhere := buttonPressed != 0 && selectorPress.selectorEntry != selectorEntry
	if isOverSelectorItem && !isPressFromElsewhere {
		itemIndex := characterEntry.AttributeEntry.CellControlId
		// The cell may predate a change to the item list that has not been redrawn yet, so its item may be gone.
		isItemIndexValid := itemIndex >= 0 && itemIndex < shared.getItemCount(selectorEntry)
		if buttonPressed != 0 && isItemIndexValid {
			_, _, previousButtonPressed, _ := GetPreviousMouseStatus()
			isNewSelection := true
			if previousButtonPressed == 0 {
				selectorEntry.SelectionSource = shared.getMouseSelectionSource(characterEntry.LayerAlias, selectorEntry.Alias, itemIndex, time.Now())
			} else if selectorPress.itemIndex != itemIndex {
				// Dragging with the button held onto another item selects it, but only as a single click.
				selectorLastClick = selectorClickType{}
				selectorEntry.SelectionSource = constants.SelectionSourceSingleClick
			} else {
				// The held press already reported this item, so reporting it again would repeat the same selection.
				isNewSelection = false
			}
			if isNewSelection {
				selectorPress.itemIndex = itemIndex
				selectorEntry.ItemHighlighted = itemIndex
				selectorEntry.ItemSelected = itemIndex
				selectorEntry.IsNewItemSelected = true
			}
		} else if buttonPressed == 0 && isItemIndexValid && !selectorEntry.HighlightOnClickOnly {
			selectorEntry.ItemHighlighted = itemIndex
		}
		// Check if this selector belongs to a dropdown
		for _, currentDropdownEntry := range dropdownRegistry.GetAllEntries(characterEntry.LayerAlias) {
			dropdownEntry := currentDropdownEntry
			if dropdownEntry.SelectorAlias == characterEntry.AttributeEntry.CellControlAlias {
				// Focus itself is given by the press, through focusControlFromClick, which maps a tray item to its
				// dropdown. Hovering only tracks the highlight.
				setPreviouslyHighlightedControl(characterEntry.LayerAlias, dropdownEntry.Alias, constants.CellTypeDropdown)
				isScreenUpdateRequired = true
				return isScreenUpdateRequired
			}
		}
		// Focus itself is given by the press, through focusControlFromClick. Hovering only tracks the highlight.
		setPreviouslyHighlightedControl(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias, constants.CellTypeSelectorItem)
		isScreenUpdateRequired = true
	} else if !isOverSelectorItem {
		previouslyHighlightedControl := getPreviouslyHighlightedControl()
		if selectorEntry, isFound := selectorRegistry.Lookup(previouslyHighlightedControl.layerAlias, previouslyHighlightedControl.controlAlias); previouslyHighlightedControl.controlType == constants.CellTypeSelectorItem &&
			isFound && characterEntry.AttributeEntry.CellType == constants.CellTypeSelectorItem &&
			selectorRegistry.IsExists(characterEntry.LayerAlias, characterEntry.AttributeEntry.CellControlAlias) {
			// Only clear highlighting if HighlightOnClickOnly is false
			if !selectorEntry.HighlightOnClickOnly {
				selectorEntry.ItemHighlighted = constants.NullItemSelection
			}
			setPreviouslyHighlightedControl("", "", constants.NullControlType)
			isScreenUpdateRequired = true
		}
	}

	// --- SCROLL BAR SYNC CODE ---
	layerAlias := characterEntry.LayerAlias

	// If a buttonType is pressed AND (you are in a drag and drop event OR the cell type is scroll bar), then
	// sync all Dropdown selectors with their appropriate scroll bars. If the control under focus
	// matches a control that belongs to a Dropdown list, then stop processing (Do not attempt to close Dropdown).
	if buttonPressed != 0 && (getEventStateId() == constants.EventStateDragAndDropScrollbar ||
		characterEntry.AttributeEntry.CellType == constants.CellTypeScrollbar) {
		for _, currentSelectorEntry := range selectorRegistry.GetAllEntries(focusedLayerAlias) {
			selectorEntry := currentSelectorEntry
			scrollBarEntry, isFound := scrollbarRegistry.Lookup(focusedLayerAlias, selectorEntry.ScrollbarAlias)
			if !isFound {
				continue
			}
			if selectorEntry.ViewportPosition != scrollBarEntry.ScrollValue {
				selectorEntry.ViewportPosition = scrollBarEntry.ScrollValue
				shared.normalizeIndexes(selectorEntry)
				isScreenUpdateRequired = true
			}
		}
	}
	// If a Selector is no longer visible, then make the scroll bars associated with it invisible as well.
	for _, currentSelectorEntry := range selectorRegistry.GetAllEntries(layerAlias) {
		selectorEntry := currentSelectorEntry
		scrollBarEntry, isFound := scrollbarRegistry.Lookup(layerAlias, selectorEntry.ScrollbarAlias)
		if !isFound {
			continue
		}
		if !selectorEntry.IsVisible {
			scrollBarEntry.IsVisible = false
		} else {
			if scrollBarEntry.IsEnabled {
				scrollBarEntry.IsVisible = true
			}
		}
	}
	return isScreenUpdateRequired
}
