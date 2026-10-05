package consolizer

import (
	"fmt"
	"sort"
	"sync"

	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/types"
)

// noExplicitTabIndex marks a tab stop that has no explicit order set through SetTabIndex, so it is ordered by when it
// was registered, after every stop that does have an explicit order.
const noExplicitTabIndex = -1

// focusChangeQueueCapacity is the most focus changes kept for GetFocusChange. When the queue is full, the oldest
// change is dropped to make room, so an application that never polls cannot grow it without bound.
const focusChangeQueueCapacity = 256

/*
FocusChangeType is a structure which describes a single change of keyboard focus, as returned by GetFocusChange. It
holds the control that had focus before the change and the control that has it after, each as a layer alias, control
alias, and cell type. An empty control alias with constants.NullControlType means no control.

Example:

	change, isFound := GetFocusChange()
*/
type FocusChangeType struct {
	PreviousLayerAlias   string
	PreviousControlAlias string
	PreviousControlType  int
	LayerAlias           string
	ControlAlias         string
	ControlType          int
}

/*
FocusableControlType is an interface which allows you to pass any control instance, such as a button, text field, or
selector, to functions like SetFocus. Every control instance type satisfies it through a pointer, since they all embed
BaseControlInstanceType.

Example:

	err := SetFocus(&okButton)
*/
type FocusableControlType interface {
	getControlIdentifier() controlIdentifierType
}

/*
focusStopType is a structure which describes one registered tab stop: the control, its explicit order if one was set
with SetTabIndex, and the sequence number of its registration, which orders stops without an explicit order.

Example:

	var stop focusStopType
*/
type focusStopType struct {
	control  controlIdentifierType
	tabIndex int
	sequence int
}

/*
focusScopeType is a structure which holds the tab order owned by one layer, along with the aliases of the layer's
default button, which Enter presses, and cancel button, which Esc presses.

Example:

	var scope focusScopeType
*/
type focusScopeType struct {
	stops         []focusStopType
	defaultButton string
	cancelButton  string
}

/*
modalEntryType is a structure which records an active modal layer and the control that had focus just before the
modal became active, so that focus can be returned to it when the modal closes.

Example:

	var modal modalEntryType
*/
type modalEntryType struct {
	layerAlias string
	savedFocus controlIdentifierType
}

/*
tabEntryType is a structure which describes one stop in the Tab sequence after radio groups have been folded together.
A plain control has a single member. A radio group has one member per registered radio button in the group, and the
whole group counts as one stop.

Example:

	var entry tabEntryType
*/
type tabEntryType struct {
	members      []controlIdentifierType
	isRadioGroup bool
	radioGroupId int
}

/*
focusManagerType is a structure which owns all keyboard focus state: the focused control, which is the single source of
truth for focus, every layer's tab order, the modal layers, and the queue of focus changes. In addition, the following
should be noted:

  - Every field is guarded by mutex, since the event goroutine, the periodic-event goroutine, and the application's
    own goroutine all use it. Functions whose names end in Locked must only be called with mutex held.

  - While mutex is held, only leaf locks may be taken: the control and layer memory managers, mouse memory, and the
    event state lock. commonResource.displayUpdate must never be taken while holding it, though it may already be
    held when mutex is taken, as it is while controls are drawn.

Example:

	focusManager.mutex.Lock()
*/
type focusManagerType struct {
	mutex                   sync.Mutex
	focusedControl          controlIdentifierType
	scopes                  map[string]*focusScopeType
	nextSequence            int
	modalLayers             map[string]int
	modalStack              []modalEntryType
	changeQueue             []FocusChangeType
	isFocusIndicatorVisible bool
}

var focusManager = focusManagerType{
	scopes:                  map[string]*focusScopeType{},
	modalLayers:             map[string]int{},
	isFocusIndicatorVisible: true,
}

/*
getFocusStopCellType is a method which allows you to obtain the cell type used to identify a control of the given
control type string, such as constants.TYPE_BUTTON, in focus and tab order state. It returns constants.NullControlType
for an unknown control type.

Example:

	cellType := getFocusStopCellType(constants.TYPE_BUTTON)
*/
func getFocusStopCellType(controlType string) int {
	switch controlType {
	case constants.TYPE_BUTTON:
		return constants.CellTypeButton
	case constants.TYPE_CHECKBOX:
		return constants.CellTypeCheckbox
	case constants.TYPE_DROPDOWN:
		return constants.CellTypeDropdown
	case constants.TYPE_LABEL:
		return constants.CellTypeLabel
	case constants.TYPE_PROGRESSBAR:
		return constants.CellTypeProgressBar
	case constants.TYPE_SCROLLBAR:
		return constants.CellTypeScrollbar
	case constants.TYPE_SELECTOR:
		return constants.CellTypeSelectorItem
	case constants.TYPE_TEXTBOX:
		return constants.CellTypeTextbox
	case constants.TYPE_TEXTFIELD:
		return constants.CellTypeTextField
	case constants.TYPE_TOOLTIP:
		return constants.CellTypeTooltip
	case constants.TYPE_RADIOBUTTON:
		return constants.CellTypeRadioButton
	case constants.TYPE_FILEMENU:
		return constants.CellTypeFileMenuHeading
	}
	return constants.NullControlType
}

/*
isTabStopCellType is a method which allows you to check whether controls of the given cell type are interactive and
can therefore be registered as tab stops and take keyboard focus.

Example:

	isAllowed := isTabStopCellType(constants.CellTypeButton)
*/
func isTabStopCellType(cellType int) bool {
	switch cellType {
	case constants.CellTypeButton, constants.CellTypeCheckbox, constants.CellTypeDropdown,
		constants.CellTypeRadioButton, constants.CellTypeScrollbar, constants.CellTypeSelectorItem,
		constants.CellTypeTextbox, constants.CellTypeTextField, constants.CellTypeFileMenuHeading:
		return true
	}
	return false
}

/*
isControlFocusable is a method which allows you to check whether a control can currently take focus or be operated
from the keyboard. A control qualifies only when it is an interactive type, its layer exists and is visible, and the
control itself still exists, is visible, and is enabled. File menus, which have no visibility flag, are judged on the
remaining checks.

Example:

	isFocusable := isControlFocusable(getFocusedControl())
*/
func isControlFocusable(entry controlIdentifierType) bool {
	if Layers == nil || entry.controlAlias == "" {
		return false
	}
	layerEntry, isFound := Layers.Lookup(entry.layerAlias)
	if !isFound || !layerEntry.IsVisible {
		return false
	}
	var baseControl *types.BaseControlType
	switch entry.controlType {
	case constants.CellTypeButton:
		if controlEntry, isFound := Buttons.Lookup(entry.layerAlias, entry.controlAlias); isFound {
			baseControl = &controlEntry.BaseControlType
		}
	case constants.CellTypeCheckbox:
		if controlEntry, isFound := Checkboxes.Lookup(entry.layerAlias, entry.controlAlias); isFound {
			baseControl = &controlEntry.BaseControlType
		}
	case constants.CellTypeDropdown:
		if controlEntry, isFound := Dropdowns.Lookup(entry.layerAlias, entry.controlAlias); isFound {
			baseControl = &controlEntry.BaseControlType
		}
	case constants.CellTypeScrollbar:
		if controlEntry, isFound := ScrollBars.Lookup(entry.layerAlias, entry.controlAlias); isFound {
			baseControl = &controlEntry.BaseControlType
		}
	case constants.CellTypeSelectorItem:
		if controlEntry, isFound := Selectors.Lookup(entry.layerAlias, entry.controlAlias); isFound {
			baseControl = &controlEntry.BaseControlType
		}
	case constants.CellTypeTextField:
		if controlEntry, isFound := TextFields.Lookup(entry.layerAlias, entry.controlAlias); isFound {
			baseControl = &controlEntry.BaseControlType
		}
	case constants.CellTypeTextbox:
		if controlEntry, isFound := Textboxes.Lookup(entry.layerAlias, entry.controlAlias); isFound {
			baseControl = &controlEntry.BaseControlType
		}
	case constants.CellTypeRadioButton:
		if controlEntry, isFound := RadioButtons.Lookup(entry.layerAlias, entry.controlAlias); isFound {
			baseControl = &controlEntry.BaseControlType
		}
	case constants.CellTypeFileMenuHeading:
		fileMenuEntry, isFound := FileMenus.Lookup(entry.layerAlias, entry.controlAlias)
		return isFound && fileMenuEntry.IsEnabled
	}
	return baseControl != nil && baseControl.IsVisible && baseControl.IsEnabled
}

/*
isControlExists is a method which allows you to check whether the control a focus identifier refers to still exists,
regardless of whether it is enabled or visible.

Example:

	isFound := isControlExists(stop.control)
*/
func isControlExists(entry controlIdentifierType) bool {
	switch entry.controlType {
	case constants.CellTypeButton:
		return Buttons.IsExists(entry.layerAlias, entry.controlAlias)
	case constants.CellTypeCheckbox:
		return Checkboxes.IsExists(entry.layerAlias, entry.controlAlias)
	case constants.CellTypeDropdown:
		return Dropdowns.IsExists(entry.layerAlias, entry.controlAlias)
	case constants.CellTypeScrollbar:
		return ScrollBars.IsExists(entry.layerAlias, entry.controlAlias)
	case constants.CellTypeSelectorItem:
		return Selectors.IsExists(entry.layerAlias, entry.controlAlias)
	case constants.CellTypeTextField:
		return TextFields.IsExists(entry.layerAlias, entry.controlAlias)
	case constants.CellTypeTextbox:
		return Textboxes.IsExists(entry.layerAlias, entry.controlAlias)
	case constants.CellTypeRadioButton:
		return RadioButtons.IsExists(entry.layerAlias, entry.controlAlias)
	case constants.CellTypeFileMenuHeading:
		return FileMenus.IsExists(entry.layerAlias, entry.controlAlias)
	}
	return false
}

/*
getRadioGroupId is a method which allows you to obtain the group ID of a radio button, and whether the control is a
radio button that still exists.

Example:

	groupId, isRadioButton := getRadioGroupId(control)
*/
func getRadioGroupId(entry controlIdentifierType) (int, bool) {
	if entry.controlType != constants.CellTypeRadioButton {
		return 0, false
	}
	radioButtonEntry, isFound := RadioButtons.Lookup(entry.layerAlias, entry.controlAlias)
	if !isFound {
		return 0, false
	}
	return radioButtonEntry.GroupId, true
}

/*
isLayerInsideLayer is a method which allows you to check whether a layer is the given ancestor layer itself or one of
its descendants, by walking up the parent chain.

Example:

	isInside := isLayerInsideLayer("dialogButtons", "dialog")
*/
func isLayerInsideLayer(layerAlias string, ancestorAlias string) bool {
	if Layers == nil {
		return false
	}
	visitedAliases := map[string]bool{}
	for currentAlias := layerAlias; currentAlias != "" && !visitedAliases[currentAlias]; {
		if currentAlias == ancestorAlias {
			return true
		}
		visitedAliases[currentAlias] = true
		layerEntry, isFound := Layers.Lookup(currentAlias)
		if !isFound {
			return false
		}
		currentAlias = layerEntry.ParentAlias
	}
	return false
}

/*
getActiveModalAliasLocked is a method which allows you to obtain the alias of the modal layer that currently owns
keyboard focus, or an empty string when no modal is active.

Example:

	modalAlias := focusManager.getActiveModalAliasLocked()
*/
func (shared *focusManagerType) getActiveModalAliasLocked() string {
	if len(shared.modalStack) == 0 {
		return ""
	}
	return shared.modalStack[len(shared.modalStack)-1].layerAlias
}

/*
isInsideActiveModalLocked is a method which allows you to check whether a layer may receive input while the current
modal is active. Every layer qualifies when no modal is active; otherwise only the modal layer and its descendants do.

Example:

	isAllowed := focusManager.isInsideActiveModalLocked("dialog")
*/
func (shared *focusManagerType) isInsideActiveModalLocked(layerAlias string) bool {
	modalAlias := shared.getActiveModalAliasLocked()
	return modalAlias == "" || isLayerInsideLayer(layerAlias, modalAlias)
}

/*
beginChange is a method which allows you to take the focus manager lock for an operation that may change focus, and
returns the function that ends it. Ending the operation queues a single focus change if focus differs from what it was
when the operation began, then releases the lock. In addition, the following should be noted:

  - Changes are queued per operation rather than per step, so an operation that passes through intermediate focus
    values, such as reconciling after a modal closes, reports only where focus actually ended up.

Example:

	defer focusManager.beginChange()()
*/
func (shared *focusManagerType) beginChange() func() {
	shared.mutex.Lock()
	previousControl := shared.focusedControl
	return func() {
		shared.queueFocusChangeLocked(previousControl)
		shared.mutex.Unlock()
	}
}

/*
queueFocusChangeLocked is a method which allows you to add a focus change to the queue read by GetFocusChange, when the
focused control differs from the given previous control. When the queue is full, the oldest change is dropped.

Example:

	focusManager.queueFocusChangeLocked(previousControl)
*/
func (shared *focusManagerType) queueFocusChangeLocked(previousControl controlIdentifierType) {
	currentControl := shared.focusedControl
	if previousControl == currentControl {
		return
	}
	if len(shared.changeQueue) >= focusChangeQueueCapacity {
		shared.changeQueue = shared.changeQueue[1:]
	}
	shared.changeQueue = append(shared.changeQueue, FocusChangeType{
		PreviousLayerAlias:   previousControl.layerAlias,
		PreviousControlAlias: previousControl.controlAlias,
		PreviousControlType:  getNormalizedControlType(previousControl),
		LayerAlias:           currentControl.layerAlias,
		ControlAlias:         currentControl.controlAlias,
		ControlType:          getNormalizedControlType(currentControl),
	})
}

/*
setFocusLocked is a method which allows you to make a control the focused control. The change is queued for
GetFocusChange when the surrounding operation ends, as described for beginChange. In addition, the following should be
noted:

  - An identifier with an empty control alias clears focus, and is stored as the empty identifier so that "no focus"
    has exactly one representation.

  - When a selector gains focus with no valid item highlighted, its selected item, or its first item, is
    highlighted, so that its focus is visible and Enter has an item to pick. An empty selector is left with no
    highlight.

Example:

	okButton := controlIdentifierType{layerAlias: "main", controlAlias: "ok", controlType: constants.CellTypeButton}
	focusManager.setFocusLocked(okButton)
*/
func (shared *focusManagerType) setFocusLocked(control controlIdentifierType) {
	if control.controlAlias == "" {
		control = controlIdentifierType{}
	}
	if shared.focusedControl == control {
		return
	}
	shared.focusedControl = control
	if control.controlType == constants.CellTypeSelectorItem {
		if selectorEntry, isFound := Selectors.Lookup(control.layerAlias, control.controlAlias); isFound {
			Selector.highlightStartingItem(selectorEntry)
		}
	}
}

/*
getNormalizedControlType is a method which allows you to obtain the cell type reported for a focus identifier, which
is constants.NullControlType when the identifier refers to no control.

Example:

	controlType := getNormalizedControlType(control)
*/
func getNormalizedControlType(control controlIdentifierType) int {
	if control.controlAlias == "" {
		return constants.NullControlType
	}
	return control.controlType
}

/*
getOrCreateScopeLocked is a method which allows you to obtain the focus scope owned by a layer, creating an empty one
if the layer has none yet.

Example:

	scope := focusManager.getOrCreateScopeLocked("main")
*/
func (shared *focusManagerType) getOrCreateScopeLocked(layerAlias string) *focusScopeType {
	scope, isFound := shared.scopes[layerAlias]
	if !isFound {
		scope = &focusScopeType{}
		shared.scopes[layerAlias] = scope
	}
	return scope
}

/*
sortStops is a method which allows you to put a scope's tab stops into Tab order: stops with an explicit order from
SetTabIndex first, in ascending order, followed by every other stop in the order it was registered. Stops with the same
explicit order keep their registration order.

Example:

	sortStops(scope.stops)
*/
func sortStops(stops []focusStopType) {
	sort.SliceStable(stops, func(firstIndex int, secondIndex int) bool {
		first := stops[firstIndex]
		second := stops[secondIndex]
		isFirstExplicit := first.tabIndex != noExplicitTabIndex
		isSecondExplicit := second.tabIndex != noExplicitTabIndex
		if isFirstExplicit != isSecondExplicit {
			return isFirstExplicit
		}
		if isFirstExplicit && first.tabIndex != second.tabIndex {
			return first.tabIndex < second.tabIndex
		}
		return first.sequence < second.sequence
	})
}

/*
registerStopLocked is a method which allows you to add a control to its layer's tab order, or update its explicit order
if it is already registered. Passing noExplicitTabIndex keeps a newly registered stop in registration order and leaves
an existing stop's order unchanged.

Example:

	focusManager.registerStopLocked(control, noExplicitTabIndex)
*/
func (shared *focusManagerType) registerStopLocked(control controlIdentifierType, tabIndex int) {
	scope := shared.getOrCreateScopeLocked(control.layerAlias)
	for index := range scope.stops {
		if scope.stops[index].control == control {
			if tabIndex != noExplicitTabIndex {
				scope.stops[index].tabIndex = tabIndex
				sortStops(scope.stops)
			}
			return
		}
	}
	shared.nextSequence++
	scope.stops = append(scope.stops, focusStopType{control: control, tabIndex: tabIndex, sequence: shared.nextSequence})
	sortStops(scope.stops)
}

/*
removeStopLocked is a method which allows you to remove a control from its layer's tab order, along with any default or
cancel button role it has. Removing a control that is not registered does nothing.

Example:

	focusManager.removeStopLocked(control)
*/
func (shared *focusManagerType) removeStopLocked(control controlIdentifierType) {
	scope, isFound := shared.scopes[control.layerAlias]
	if !isFound {
		return
	}
	remainingStops := scope.stops[:0]
	for _, stop := range scope.stops {
		if stop.control != control {
			remainingStops = append(remainingStops, stop)
		}
	}
	scope.stops = remainingStops
	if control.controlType == constants.CellTypeButton {
		if scope.defaultButton == control.controlAlias {
			scope.defaultButton = ""
		}
		if scope.cancelButton == control.controlAlias {
			scope.cancelButton = ""
		}
	}
}

/*
buildTabEntriesLocked is a method which allows you to obtain a layer's tab order with radio groups folded together, so
that each radio group appears once, at the position of its first registered member, and counts as a single stop.

Example:

	entries := focusManager.buildTabEntriesLocked("main")
*/
func (shared *focusManagerType) buildTabEntriesLocked(layerAlias string) []tabEntryType {
	scope, isFound := shared.scopes[layerAlias]
	if !isFound {
		return nil
	}
	var entries []tabEntryType
	radioGroupEntryIndex := map[int]int{}
	for _, stop := range scope.stops {
		if groupId, isRadioButton := getRadioGroupId(stop.control); isRadioButton {
			if entryIndex, isGroupFound := radioGroupEntryIndex[groupId]; isGroupFound {
				entries[entryIndex].members = append(entries[entryIndex].members, stop.control)
				continue
			}
			radioGroupEntryIndex[groupId] = len(entries)
			entries = append(entries, tabEntryType{members: []controlIdentifierType{stop.control}, isRadioGroup: true, radioGroupId: groupId})
			continue
		}
		entries = append(entries, tabEntryType{members: []controlIdentifierType{stop.control}})
	}
	return entries
}

/*
getEntryTarget is a method which allows you to obtain the control that should receive focus when Tab lands on an
entry, and whether the entry can take focus at all. A plain entry targets its control. A radio group targets its
selected radio button when that one can take focus, and otherwise its first member that can.

Example:

	target, isFocusable := getEntryTarget(entry)
*/
func getEntryTarget(entry tabEntryType) (controlIdentifierType, bool) {
	if !entry.isRadioGroup {
		return entry.members[0], isControlFocusable(entry.members[0])
	}
	if len(entry.members) > 0 {
		selectedAlias := getSelectedRadioButton(entry.members[0].layerAlias, entry.members[0].controlAlias)
		selectedControl := controlIdentifierType{layerAlias: entry.members[0].layerAlias, controlAlias: selectedAlias, controlType: constants.CellTypeRadioButton}
		if selectedAlias != "" && isControlFocusable(selectedControl) {
			return selectedControl, true
		}
	}
	for _, member := range entry.members {
		if isControlFocusable(member) {
			return member, true
		}
	}
	return controlIdentifierType{}, false
}

/*
findEntryIndex is a method which allows you to locate the entry a control belongs to within a folded tab order. A radio
button matches its group's entry even when it is not itself registered, as long as it is in the same group on the same
layer. It returns -1 when the control has no entry.

Example:

	index := findEntryIndex(entries, focusedControl)
*/
func findEntryIndex(entries []tabEntryType, control controlIdentifierType) int {
	groupId, isRadioButton := getRadioGroupId(control)
	for entryIndex, entry := range entries {
		if entry.isRadioGroup && isRadioButton && entry.radioGroupId == groupId && entry.members[0].layerAlias == control.layerAlias {
			return entryIndex
		}
		for _, member := range entry.members {
			if member == control {
				return entryIndex
			}
		}
	}
	return -1
}

/*
getActiveScopeAliasLocked is a method which allows you to obtain the alias of the layer whose tab order Tab and
Shift+Tab currently move through, and whose default and cancel buttons Enter and Esc press. This is the active modal
layer when there is one, otherwise the layer of the focused control, and otherwise the topmost visible layer that has a
stop able to take focus. An empty string means there is no active scope.

Example:

	scopeAlias := focusManager.getActiveScopeAliasLocked()
*/
func (shared *focusManagerType) getActiveScopeAliasLocked() string {
	if modalAlias := shared.getActiveModalAliasLocked(); modalAlias != "" {
		return modalAlias
	}
	if shared.focusedControl.controlAlias != "" {
		return shared.focusedControl.layerAlias
	}
	if Layers == nil {
		return ""
	}
	sortedLayers := layer.GetSortedLayerMemoryAliasSlice()
	for index := len(sortedLayers) - 1; index >= 0; index-- {
		layerAlias := sortedLayers[index].Key
		layerEntry, isFound := Layers.Lookup(layerAlias)
		if !isFound || !layerEntry.IsVisible {
			continue
		}
		for _, entry := range shared.buildTabEntriesLocked(layerAlias) {
			if _, isFocusable := getEntryTarget(entry); isFocusable {
				return layerAlias
			}
		}
	}
	return ""
}

/*
moveFocusLocked is a method which allows you to move focus to the next stop, when direction is 1, or the previous
stop, when direction is -1, in the active scope's tab order, wrapping around at either end and skipping stops that
cannot take focus. It returns true if a stop took focus. In addition, the following should be noted:

  - The search starts from the focused control's stop. When the focused control has no stop in the scope, it starts
    just before the first stop when moving forward, or just after the last stop when moving backward.

  - The focused control's own stop is checked last, so when it is the only stop able to take focus, it keeps focus.

  - A radio group is a single stop, so moving away from any of its members skips the rest of the group.

Example:

	isMoved := focusManager.moveFocusLocked(1)
*/
func (shared *focusManagerType) moveFocusLocked(direction int) bool {
	scopeAlias := shared.getActiveScopeAliasLocked()
	if scopeAlias == "" {
		return false
	}
	entries := shared.buildTabEntriesLocked(scopeAlias)
	entryCount := len(entries)
	if entryCount == 0 {
		return false
	}
	startIndex := -1
	if shared.focusedControl.layerAlias == scopeAlias {
		startIndex = findEntryIndex(entries, shared.focusedControl)
	}
	if startIndex == -1 {
		if direction > 0 {
			startIndex = entryCount - 1
		} else {
			startIndex = 0
		}
	}
	for step := 1; step <= entryCount; step++ {
		candidateIndex := ((startIndex+direction*step)%entryCount + entryCount) % entryCount
		if target, isFocusable := getEntryTarget(entries[candidateIndex]); isFocusable {
			shared.setFocusLocked(target)
			return true
		}
	}
	return false
}

/*
reconcileLocked is a method which allows you to bring focus back into a valid state after anything that may have
invalidated it, such as a control being disabled, hidden, or deleted, a layer being hidden or deleted, or a modal
opening or closing. In addition, the following should be noted:

  - Modal layers that were hidden, deleted, or unmarked are closed. When the topmost one closes, focus returns to the
    control that had it before that modal opened, if it can still take focus.

  - Modal layers that are marked and visible but not yet active are opened in the order they were marked, each
    saving the current focus for later and moving focus to its own first stop.

  - A focused control that can no longer take focus moves to the next stop in its scope, or loses focus if none can
    take it. A focused control outside the active modal moves to the modal's first stop.

  - Tab stops, scopes, and default or cancel buttons that refer to controls or layers that no longer exist are
    removed last, after the focus decisions above have had a chance to use their positions.

Example:

	focusManager.reconcileLocked()
*/
func (shared *focusManagerType) reconcileLocked() {
	if Layers == nil {
		return
	}
	for layerAlias := range shared.modalLayers {
		if !isLayerExists(layerAlias) {
			delete(shared.modalLayers, layerAlias)
		}
	}
	isModalActive := func(layerAlias string) bool {
		_, isMarked := shared.modalLayers[layerAlias]
		layerEntry, isFound := Layers.Lookup(layerAlias)
		return isMarked && isFound && layerEntry.IsVisible
	}
	var focusToRestore *controlIdentifierType
	for len(shared.modalStack) > 0 && !isModalActive(shared.getActiveModalAliasLocked()) {
		savedFocus := shared.modalStack[len(shared.modalStack)-1].savedFocus
		focusToRestore = &savedFocus
		shared.modalStack = shared.modalStack[:len(shared.modalStack)-1]
	}
	remainingModals := shared.modalStack[:0]
	for _, modalEntry := range shared.modalStack {
		if isModalActive(modalEntry.layerAlias) {
			remainingModals = append(remainingModals, modalEntry)
		}
	}
	shared.modalStack = remainingModals
	if focusToRestore != nil {
		shared.setFocusLocked(*focusToRestore)
		shared.moveLostFocusLocked()
	}

	var pendingModals []string
	for layerAlias := range shared.modalLayers {
		isOnStack := false
		for _, modalEntry := range shared.modalStack {
			if modalEntry.layerAlias == layerAlias {
				isOnStack = true
				break
			}
		}
		if !isOnStack && isModalActive(layerAlias) {
			pendingModals = append(pendingModals, layerAlias)
		}
	}
	sort.Slice(pendingModals, func(firstIndex int, secondIndex int) bool {
		return shared.modalLayers[pendingModals[firstIndex]] < shared.modalLayers[pendingModals[secondIndex]]
	})
	for _, layerAlias := range pendingModals {
		shared.modalStack = append(shared.modalStack, modalEntryType{layerAlias: layerAlias, savedFocus: shared.focusedControl})
		if !isLayerInsideLayer(shared.focusedControl.layerAlias, layerAlias) {
			shared.setFocusLocked(controlIdentifierType{})
			shared.moveFocusLocked(1)
		}
	}

	shared.moveLostFocusLocked()

	for layerAlias, scope := range shared.scopes {
		if !isLayerExists(layerAlias) {
			delete(shared.scopes, layerAlias)
			continue
		}
		remainingStops := scope.stops[:0]
		for _, stop := range scope.stops {
			if isControlExists(stop.control) {
				remainingStops = append(remainingStops, stop)
			}
		}
		scope.stops = remainingStops
		if scope.defaultButton != "" && !Buttons.IsExists(layerAlias, scope.defaultButton) {
			scope.defaultButton = ""
		}
		if scope.cancelButton != "" && !Buttons.IsExists(layerAlias, scope.cancelButton) {
			scope.cancelButton = ""
		}
	}
}

/*
moveLostFocusLocked is a method which allows you to move focus away from a focused control that can no longer have it.
A control that cannot take focus moves to the next stop in its own scope, and focus is cleared when no stop there can
take it, rather than jumping to an unrelated layer. A control outside the active modal moves to the modal's first stop.
Focus is left alone when the focused control is still valid.

Example:

	focusManager.moveLostFocusLocked()
*/
func (shared *focusManagerType) moveLostFocusLocked() {
	focusedControl := shared.focusedControl
	if focusedControl.controlAlias == "" {
		return
	}
	if !shared.isInsideActiveModalLocked(focusedControl.layerAlias) {
		shared.setFocusLocked(controlIdentifierType{})
		shared.moveFocusLocked(1)
		return
	}
	if isControlFocusable(focusedControl) {
		return
	}
	if !shared.moveFocusLocked(1) {
		shared.setFocusLocked(controlIdentifierType{})
	}
}

/*
reconcileFocus is a method which allows you to bring focus back into a valid state, as described for reconcileLocked,
taking the focus manager lock itself. It is called after anything that can invalidate focus, and before every event and
redraw as a safety net.

Example:

	reconcileFocus()
*/
func reconcileFocus() {
	defer focusManager.beginChange()()
	focusManager.reconcileLocked()
}

/*
handleControlDeleted is a method which allows you to update focus state for a control that has just been deleted. If
it had focus, focus moves to the next stop in its scope. Its tab stop, and any default or cancel button role it had,
are then removed right away, so that a new control later created with the same alias does not inherit them. In
addition, the following should be noted:

  - This must be called after the control has been removed from its memory manager, so that it is no longer counted
    as able to take focus.

Example:

	handleControlDeleted("main", "ok", constants.CellTypeButton)
*/
func handleControlDeleted(layerAlias string, controlAlias string, cellType int) {
	if controlAlias == "" {
		return
	}
	defer focusManager.beginChange()()
	deletedControl := controlIdentifierType{layerAlias: layerAlias, controlAlias: controlAlias, controlType: cellType}
	if focusManager.focusedControl == deletedControl {
		focusManager.moveLostFocusLocked()
	}
	focusManager.removeStopLocked(deletedControl)
	focusManager.reconcileLocked()
}

/*
focusControlFromClick is a method which allows you to give focus to the control under a mouse press. A press on an item
in a dropdown's tray focuses the dropdown, and a press on an item of a file menu's tray leaves focus unchanged. A press
on anything that cannot take focus, such as a label, a disabled control, empty layer space, or a layer outside the
active modal, also leaves focus unchanged. In addition, the following should be noted:

  - Scroll bars are not handled here, since the scroll bar's own mouse handling takes focus when it starts a drag.

Example:

	focusControlFromClick(characterEntry)
*/
func focusControlFromClick(characterEntry types.CharacterEntryType) {
	layerAlias := characterEntry.LayerAlias
	controlAlias := characterEntry.AttributeEntry.CellControlAlias
	cellType := characterEntry.AttributeEntry.CellType
	if controlAlias == "" || cellType == constants.CellTypeScrollbar || !isTabStopCellType(cellType) {
		return
	}
	target := controlIdentifierType{layerAlias: layerAlias, controlAlias: controlAlias, controlType: cellType}
	if cellType == constants.CellTypeSelectorItem {
		for _, dropdownEntry := range Dropdowns.GetAllEntries(layerAlias) {
			if dropdownEntry.SelectorAlias == controlAlias {
				target = controlIdentifierType{layerAlias: layerAlias, controlAlias: dropdownEntry.Alias, controlType: constants.CellTypeDropdown}
			}
		}
		for _, fileMenuEntry := range FileMenus.GetAllEntries(layerAlias) {
			for _, selectorAlias := range fileMenuEntry.SelectorAliases {
				if selectorAlias == controlAlias {
					return
				}
			}
		}
	}
	defer focusManager.beginChange()()
	if isControlFocusable(target) && focusManager.isInsideActiveModalLocked(target.layerAlias) {
		focusManager.setFocusLocked(target)
	}
}

/*
isLayerAcceptingInput is a method which allows you to check whether a layer may currently receive mouse input. Every
layer may when no modal is active; otherwise only the active modal layer and its descendants may.

Example:

	isAllowed := isLayerAcceptingInput("main")
*/
func isLayerAcceptingInput(layerAlias string) bool {
	focusManager.mutex.Lock()
	defer focusManager.mutex.Unlock()
	return focusManager.isInsideActiveModalLocked(layerAlias)
}

/*
getFocusedControl is a method which allows you to obtain the control that currently has focus, read as a single
consistent snapshot. Callers that need more than one of its fields must read them all from one snapshot rather than
calling this repeatedly, since focus may change between calls.

Example:

	focusedControl := getFocusedControl()
*/
func getFocusedControl() controlIdentifierType {
	focusManager.mutex.Lock()
	defer focusManager.mutex.Unlock()
	return focusManager.focusedControl
}

/*
setFocusedControl is a method which allows you to set which control has focus from inside the package, without the
checks SetFocus applies. It is used by control handlers that already know the control is valid, such as a click on an
enabled text field. Passing an empty control alias clears focus.

Example:

	setFocusedControl("layer1", "textfield1", constants.CellTypeTextField)
*/
func setFocusedControl(layerAlias string, controlAlias string, controlType int) {
	defer focusManager.beginChange()()
	focusManager.setFocusLocked(controlIdentifierType{layerAlias: layerAlias, controlAlias: controlAlias, controlType: controlType})
}

/*
isControlCurrentlyFocused is a method which checks if a specific control is currently the focused control in the
application.

Example:

	isControlCurrentlyFocused("layer1", "textfield1", constants.CellTypeTextField)
*/
func isControlCurrentlyFocused(layerAlias string, controlAlias string, cellType int) bool {
	if controlAlias == "" {
		return false
	}
	focusManager.mutex.Lock()
	defer focusManager.mutex.Unlock()
	focusedControl := focusManager.focusedControl
	return focusedControl.layerAlias == layerAlias && focusedControl.controlAlias == controlAlias && focusedControl.controlType == cellType
}

/*
setFocusIndicatorVisible is a method which allows you to record whether the focused control should currently be drawn
with its focus indicator, the style's focused colours. The indicator follows the input method last used, as the
:focus-visible rule does in web browsers: a keystroke shows it, so a keyboard user can always see which control will
receive their keys, and a mouse button press hides it, since a mouse user already knows what they clicked. It returns
true when the setting changed while a control has focus, which means the screen needs to be redrawn. In addition, the
following should be noted:

  - Only how focus is drawn depends on this setting. Which control has focus, and so where keystrokes go, is the
    same either way, and a control focused by a click still shows its indicator once the user presses a key.

  - Focus moved by the application with SetFocus is drawn according to the input method the user last used, so it
    shows the indicator until the user first presses a mouse button.

Example:

	isRedrawRequired := setFocusIndicatorVisible(false)
*/
func setFocusIndicatorVisible(isVisible bool) bool {
	focusManager.mutex.Lock()
	defer focusManager.mutex.Unlock()
	if focusManager.isFocusIndicatorVisible == isVisible {
		return false
	}
	focusManager.isFocusIndicatorVisible = isVisible
	return focusManager.focusedControl.controlAlias != ""
}

/*
isFocusIndicatorShown is a method which allows you to check whether a control should be drawn with its focus
indicator, which is the case only while it is the focused control and the user's last input came from the keyboard,
as described for setFocusIndicatorVisible. Control drawing uses this rather than isControlCurrentlyFocused, so that
clicking a control focuses it without flashing the keyboard focus colours over it.

Example:

	if isFocusIndicatorShown("layer1", "ok", constants.CellTypeButton) {
		foregroundColor, backgroundColor = getFocusedColors(foregroundColor, backgroundColor, focusedForeground, focusedBackground)
	}
*/
func isFocusIndicatorShown(layerAlias string, controlAlias string, cellType int) bool {
	if controlAlias == "" {
		return false
	}
	focusManager.mutex.Lock()
	defer focusManager.mutex.Unlock()
	focusedControl := focusManager.focusedControl
	return focusManager.isFocusIndicatorVisible && focusedControl.layerAlias == layerAlias &&
		focusedControl.controlAlias == controlAlias && focusedControl.controlType == cellType
}

/*
getDefaultButtonForKey is a method which allows you to obtain the button that a keystroke nothing else consumed should
press: the active scope's default button for Enter, or its cancel button for Esc. It returns an empty identifier when
the keystroke is neither, the scope has no such button, or the button cannot currently take focus.

Example:

	button := getDefaultButtonForKey("enter")
*/
func getDefaultButtonForKey(keystroke string) controlIdentifierType {
	if keystroke != "enter" && keystroke != "esc" {
		return controlIdentifierType{}
	}
	defer focusManager.beginChange()()
	focusManager.reconcileLocked()
	scopeAlias := focusManager.getActiveScopeAliasLocked()
	scope, isFound := focusManager.scopes[scopeAlias]
	if !isFound {
		return controlIdentifierType{}
	}
	buttonAlias := scope.defaultButton
	if keystroke == "esc" {
		buttonAlias = scope.cancelButton
	}
	button := controlIdentifierType{layerAlias: scopeAlias, controlAlias: buttonAlias, controlType: constants.CellTypeButton}
	if buttonAlias == "" || !isControlFocusable(button) {
		return controlIdentifierType{}
	}
	return button
}

/*
setScopeButton is a method which allows you to make a button the default or cancel button of the layer it belongs to.
Passing a nil button clears the role. An error is returned if the button does not exist. In addition, the following
should be noted:

  - The button does not need to be registered in the tab order to be a default or cancel button.

Example:

	err := setScopeButton("main", &okButton, true)
*/
func setScopeButton(layerAlias string, button *ButtonInstanceType, isDefault bool) error {
	buttonAlias := ""
	if button != nil {
		if button.layerAlias != layerAlias {
			return fmt.Errorf("the button '%s' belongs to layer '%s', so it cannot be a default or cancel button of layer '%s'", button.controlAlias, button.layerAlias, layerAlias)
		}
		if !Buttons.IsExists(button.layerAlias, button.controlAlias) {
			return fmt.Errorf("the button '%s' on layer '%s' does not exist", button.controlAlias, button.layerAlias)
		}
		buttonAlias = button.controlAlias
	}
	focusManager.mutex.Lock()
	defer focusManager.mutex.Unlock()
	scope := focusManager.getOrCreateScopeLocked(layerAlias)
	if isDefault {
		scope.defaultButton = buttonAlias
	} else {
		scope.cancelButton = buttonAlias
	}
	return nil
}

/*
setLayerModal is a method which allows you to mark or unmark a layer as modal, then bring focus into line with the
change. While a modal layer is visible it is the active focus scope: Tab cannot leave it, and mouse input to other
layers is ignored. When it is unmarked, hidden, or deleted, focus returns to the control that had it before the modal
opened.

Example:

	setLayerModal("dialog", true)
*/
func setLayerModal(layerAlias string, isModal bool) {
	defer focusManager.beginChange()()
	if isModal {
		if _, isMarked := focusManager.modalLayers[layerAlias]; !isMarked {
			focusManager.nextSequence++
			focusManager.modalLayers[layerAlias] = focusManager.nextSequence
		}
	} else {
		delete(focusManager.modalLayers, layerAlias)
	}
	focusManager.reconcileLocked()
}

/*
isLayerModal is a method which allows you to check whether a layer is marked as modal.

Example:

	isModal := isLayerModal("dialog")
*/
func isLayerModal(layerAlias string) bool {
	focusManager.mutex.Lock()
	defer focusManager.mutex.Unlock()
	_, isMarked := focusManager.modalLayers[layerAlias]
	return isMarked
}

/*
clearLayerTabIndex is a method which allows you to remove every tab stop of one layer, along with its default and
cancel buttons. Focus is left where it is.

Example:

	clearLayerTabIndex("main")
*/
func clearLayerTabIndex(layerAlias string) {
	focusManager.mutex.Lock()
	defer focusManager.mutex.Unlock()
	delete(focusManager.scopes, layerAlias)
}

/*
ClearTabIndex is a method which allows you to reset keyboard focus completely, for example when switching to a new
screen. Every layer's tab order, default button, and cancel button is removed, focus is cleared, and any pending focus
changes are discarded, so nothing carries over to the next screen. In addition, the following should be noted:

  - Layers marked as modal stay marked, since that is a property of the layer rather than of the tab order. A modal
    that is still visible becomes active again with nothing saved to restore.

Example:

	ClearTabIndex()
*/
func ClearTabIndex() {
	focusManager.mutex.Lock()
	defer focusManager.mutex.Unlock()
	focusManager.scopes = map[string]*focusScopeType{}
	focusManager.modalStack = nil
	focusManager.focusedControl = controlIdentifierType{}
	focusManager.changeQueue = nil
	focusManager.reconcileLocked()
	focusManager.changeQueue = nil
}

/*
GetFocusedControl is a method which allows you to obtain the control that currently has keyboard focus, so that an
application can act on it or restore focus to it later. It returns the alias of the layer the control belongs to, the
alias of the control, and its cell type as one of the constants.CellType values. When no control has focus, empty
aliases and constants.NullControlType are returned. In addition, the following should be noted:

  - Focus is brought into a valid state first, so a control that was disabled, hidden, or deleted since the last
    event is never reported.

  - While a dropdown's tray is open, the dropdown itself is reported rather than the internal selector that draws
    its tray.

  - The value is a consistent snapshot, but a later event or call may change focus as soon as this method returns.

Example:

	layerAlias, controlAlias, controlType := GetFocusedControl()
*/
func GetFocusedControl() (string, string, int) {
	defer focusManager.beginChange()()
	focusManager.reconcileLocked()
	focusedControl := focusManager.focusedControl
	return focusedControl.layerAlias, focusedControl.controlAlias, getNormalizedControlType(focusedControl)
}

/*
SetFocus is a method which allows you to give keyboard focus to a control, the same as clicking it. An error is
returned, and focus is left unchanged, if the control is not an interactive type, does not exist, is disabled or hidden,
is on a hidden layer, or is outside the active modal. In addition, the following should be noted:

  - The control does not need to be registered in the tab order. Tab then continues from the start of its layer's
    order.

Example:

	err := SetFocus(&nameTextField)
*/
func SetFocus(control FocusableControlType) error {
	if control == nil {
		return fmt.Errorf("a control must be given to receive focus")
	}
	target := control.getControlIdentifier()
	if !isTabStopCellType(target.controlType) {
		return fmt.Errorf("the control '%s' on layer '%s' cannot take focus since it is not an interactive control", target.controlAlias, target.layerAlias)
	}
	if !isControlFocusable(target) {
		return fmt.Errorf("the control '%s' on layer '%s' cannot take focus since it does not exist, is disabled, or is hidden", target.controlAlias, target.layerAlias)
	}
	defer focusManager.beginChange()()
	focusManager.reconcileLocked()
	if !focusManager.isInsideActiveModalLocked(target.layerAlias) {
		return fmt.Errorf("the control '%s' on layer '%s' cannot take focus while the modal layer '%s' is active", target.controlAlias, target.layerAlias, focusManager.getActiveModalAliasLocked())
	}
	focusManager.setFocusLocked(target)
	return nil
}

/*
FocusNext is a method which allows you to move focus to the next stop in the active scope's tab order, exactly as
pressing Tab does, wrapping from the last stop to the first and skipping stops that cannot take focus. It returns true
if a stop took focus.

Example:

	isMoved := FocusNext()
*/
func FocusNext() bool {
	defer focusManager.beginChange()()
	focusManager.reconcileLocked()
	return focusManager.moveFocusLocked(1)
}

/*
FocusPrevious is a method which allows you to move focus to the previous stop in the active scope's tab order, exactly
as pressing Shift+Tab does, wrapping from the first stop to the last and skipping stops that cannot take focus. It
returns true if a stop took focus.

Example:

	isMoved := FocusPrevious()
*/
func FocusPrevious() bool {
	defer focusManager.beginChange()()
	focusManager.reconcileLocked()
	return focusManager.moveFocusLocked(-1)
}

/*
GetFocusChange is a method which allows you to obtain the oldest focus change that has not been read yet, removing it
from the queue. The second value is false when there are no unread changes. Every change of focus is queued, whether
it came from Tab, a click, SetFocus, a modal opening or closing, or the focused control being disabled, hidden, or
deleted. In addition, the following should be noted:

  - Nothing is ever called back on the event goroutine. The application polls this queue from its own goroutine,
    typically once per pass of its main loop, and it is safe to call from any goroutine.

  - At most 256 changes are kept. When more arrive before they are read, the oldest are dropped, so the most recent
    change always reflects the current focus.

Example:

	for change, isFound := GetFocusChange(); isFound; change, isFound = GetFocusChange() {
		updateStatusBar(change.ControlAlias)
	}
*/
func GetFocusChange() (FocusChangeType, bool) {
	focusManager.mutex.Lock()
	defer focusManager.mutex.Unlock()
	if len(focusManager.changeQueue) == 0 {
		return FocusChangeType{}, false
	}
	change := focusManager.changeQueue[0]
	focusManager.changeQueue = focusManager.changeQueue[1:]
	return change, true
}

/*
getFocusedColors is a method which allows you to obtain the foreground and background colours to draw a focused control
with, given its normal colours and the focused colours from its style, so that focus is always visible (WCAG 2.4.7
Focus Visible). In addition, the following should be noted:

  - A focused colour left unset (zero) in the style falls back to the corresponding normal colour.

  - If the result would look exactly like the unfocused control, which happens when a style sets neither focused
    colour or sets them to the normal ones, the normal foreground and background are swapped instead.

Example:

	foregroundColor, backgroundColor := getFocusedColors(style.Button.ForegroundColor, style.Button.BackgroundColor,
		style.Button.FocusedForegroundColor, style.Button.FocusedBackgroundColor)
*/
func getFocusedColors(foregroundColor constants.ColorType, backgroundColor constants.ColorType, focusedForegroundColor constants.ColorType, focusedBackgroundColor constants.ColorType) (constants.ColorType, constants.ColorType) {
	resolvedForegroundColor := focusedForegroundColor
	if resolvedForegroundColor == 0 {
		resolvedForegroundColor = foregroundColor
	}
	resolvedBackgroundColor := focusedBackgroundColor
	if resolvedBackgroundColor == 0 {
		resolvedBackgroundColor = backgroundColor
	}
	if resolvedForegroundColor == foregroundColor && resolvedBackgroundColor == backgroundColor {
		return backgroundColor, foregroundColor
	}
	return resolvedForegroundColor, resolvedBackgroundColor
}
