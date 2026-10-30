package consolizer

import (
	"github.com/gdamore/tcell/v2"
	"sync"
	"time"
)

/*
keyboardMemoryType is a structure which holds information about the current keyboard state including buffer, key states,
and the last keystroke.
*/
type keyboardMemoryType struct {
	sync.Mutex
	entries [][]rune
	// Map to track the current state of each key (pressed or released)
	keyStates map[string]bool
	// Store the last keystroke for reference
	lastKeystroke []rune
}

/*
KeyboardMemory is a variable which acts as the global instance for managing the current keyboard state.
*/
var KeyboardMemory keyboardMemoryType

func init() {
	// Initialize the keyStates map
	KeyboardMemory.keyStates = make(map[string]bool)
}

/*
AddToBuffer is a method which allows you to add one or more keystrokes to the keyboard buffer. Each keystroke is
represented as a slice of runes.

Example:
    KeyboardMemory.AddToBuffer(rune{'a'}, rune{'b'})
*/
func (shared *keyboardMemoryType) AddToBuffer(keystroke ...[]rune) {
	shared.Lock()
	defer func() {
		shared.Unlock()
	}()
	for _, currentKeystroke := range keystroke {
		shared.entries = append(shared.entries, currentKeystroke)

		// Update the key state to pressed
		if len(currentKeystroke) > 0 {
			keyString := string(currentKeystroke)
			shared.keyStates[keyString] = true
			shared.lastKeystroke = currentKeystroke
		}
	}
}

/*
GetFromBuffer is a method which allows you to retrieve the next available keystroke from the keyboard buffer. The
keystroke is returned as a slice of runes.

Example:
    keystroke := KeyboardMemory.GetFromBuffer()
*/
func (shared *keyboardMemoryType) GetFromBuffer() []rune {
	shared.Lock()
	defer shared.Unlock()
	if len(shared.entries) == 0 {
		return nil
	}
	keystroke := shared.entries[0]
	shared.entries = shared.entries[1:]
	return keystroke
}

/*
LiveInkey is a method which allows you to retrieve the key currently being pressed. If no key is currently pressed, an
empty string is returned.

Example:
    key := KeyboardMemory.LiveInkey()
*/
func (shared *keyboardMemoryType) LiveInkey() string {
	shared.Lock()
	defer shared.Unlock()

	// Check if any key is currently pressed
	for key, pressed := range shared.keyStates {
		if pressed {
			return key
		}
	}

	return ""
}

/*
IsKeyPressed is a method which allows you to check if a specific key is currently being pressed.

Example:
    isPressed := KeyboardMemory.IsKeyPressed("a")
*/
func (shared *keyboardMemoryType) IsKeyPressed(key string) bool {
	shared.Lock()
	defer shared.Unlock()

	pressed, exists := shared.keyStates[key]
	return exists && pressed
}

/*
GetLastKeystroke is a method which allows you to retrieve the last keystroke that was processed by the keyboard buffer.

Example:
    lastKeystroke := KeyboardMemory.GetLastKeystroke()
*/
func (shared *keyboardMemoryType) GetLastKeystroke() []rune {
	shared.Lock()
	defer shared.Unlock()

	return shared.lastKeystroke
}

/*
ClearKeyStates is a method which allows you to mark all keys as released, effectively clearing the current state of all
keys.

Example:
    KeyboardMemory.ClearKeyStates()
*/
func (shared *keyboardMemoryType) ClearKeyStates() {
	shared.Lock()
	defer shared.Unlock()

	for key := range shared.keyStates {
		shared.keyStates[key] = false
	}
}

// orderedKeyboardInkeyTimeout is how recently the application must have called Inkey for ordered keyboard mode to
// hold key events back for it. Once Inkey has not been called for this long, the application is assumed to have
// stopped reading the keyboard, and the held key events are handed to the controls so they do not go unhandled.
const orderedKeyboardInkeyTimeout = 250 * time.Millisecond

/*
keyEventType is a structure which holds a single key event as it arrived from the terminal: the keystroke, in the same
form Inkey returns it, and the modifier keys held at the time, so that the event can be handled later exactly as it
would have been on arrival.

Example:

	keyEvent := keyEventType{keystroke: []rune("down"), modifiers: tcell.ModNone}
*/
type keyEventType struct {
	keystroke []rune
	modifiers tcell.ModMask
}

/*
orderedKeyboardType is a structure which holds the state of ordered keyboard mode: whether it is enabled, when the
application last called Inkey, and the key events held back until the application reads them. In addition, the
following should be noted:

  - Every field is guarded by mutex, a leaf lock which may be taken while uiStateMutex is held but never the other
    way round, and which is never held while a key event is handled.

Example:

	var orderedKeyboardState orderedKeyboardType
*/
type orderedKeyboardType struct {
	mutex            sync.Mutex
	isEnabled        bool
	lastInkeyTime    time.Time
	pendingKeyEvents []keyEventType
}

var orderedKeyboard orderedKeyboardType

/*
SetOrderedKeyboard is a method which allows you to turn ordered keyboard mode on or off. It is off by default, in which
case a key that a control uses, such as Up or Down on a focused selector, is handled the moment it arrives, while a key
no control uses waits in the keyboard buffer for Inkey, so a fast typist's keys can reach the controls and the
application out of order. While the mode is on and the application has called Inkey within the last 250 milliseconds,
every key event is instead held back and handled by the Inkey call that reads it, one key per call, so the application
sees each key at the point it was typed. In addition, the following should be noted:

  - If the application stops calling Inkey for 250 milliseconds, or the mode is turned off, the held key events are
    handed to the controls in the order they were typed, and any key no control used is placed in the keyboard
    buffer for a later Inkey call.

  - Mouse events are never held back, so a click arriving while key events are held is handled before them.

Example:

	SetOrderedKeyboard(true)
*/
func SetOrderedKeyboard(isEnabled bool) {
	defer lockUIState()()
	orderedKeyboard.mutex.Lock()
	orderedKeyboard.isEnabled = isEnabled
	orderedKeyboard.mutex.Unlock()
	if !isEnabled && flushPendingKeyEvents() {
		updateDisplayLocked(false)
	}
}

/*
IsOrderedKeyboard is a method which allows you to check whether ordered keyboard mode, as described for
SetOrderedKeyboard, is turned on.

Example:

	isOrdered := IsOrderedKeyboard()
*/
func IsOrderedKeyboard() bool {
	orderedKeyboard.mutex.Lock()
	defer orderedKeyboard.mutex.Unlock()
	return orderedKeyboard.isEnabled
}

/*
recordInkeyCall is a method which allows you to record that the application has just called Inkey, which keeps ordered
keyboard mode holding key events back for another 250 milliseconds.

Example:

	orderedKeyboard.recordInkeyCall()
*/
func (shared *orderedKeyboardType) recordInkeyCall() {
	shared.mutex.Lock()
	defer shared.mutex.Unlock()
	shared.lastInkeyTime = time.Now()
}

/*
isHoldingLocked is a method which allows you to check whether key events are currently being held back, which is the
case while ordered keyboard mode is on and the application has called Inkey within the last 250 milliseconds. The
caller must hold mutex.

Example:

	isHolding := shared.isHoldingLocked()
*/
func (shared *orderedKeyboardType) isHoldingLocked() bool {
	return shared.isEnabled && !shared.lastInkeyTime.IsZero() && time.Since(shared.lastInkeyTime) < orderedKeyboardInkeyTimeout
}

/*
holdIfActive is a method which allows you to hold a key event back for Inkey if ordered keyboard mode is currently
holding key events, as described for SetOrderedKeyboard. It returns true if the event was held, in which case the
caller must not handle it, and false if the caller should handle it immediately.

Example:

	if orderedKeyboard.holdIfActive(keyEvent) {
		return
	}
*/
func (shared *orderedKeyboardType) holdIfActive(keyEvent keyEventType) bool {
	shared.mutex.Lock()
	defer shared.mutex.Unlock()
	if !shared.isHoldingLocked() {
		return false
	}
	shared.pendingKeyEvents = append(shared.pendingKeyEvents, keyEvent)
	return true
}

/*
hasPendingKeyEvents is a method which allows you to check whether any key event is being held back.

Example:

	isPending := orderedKeyboard.hasPendingKeyEvents()
*/
func (shared *orderedKeyboardType) hasPendingKeyEvents() bool {
	shared.mutex.Lock()
	defer shared.mutex.Unlock()
	return len(shared.pendingKeyEvents) > 0
}

/*
isStaleWithPendingKeyEvents is a method which allows you to check whether key events are being held back even though
they should no longer be, because the mode was turned off or the application has not called Inkey for 250
milliseconds, so that they can be handed to the controls.

Example:

	isStale := orderedKeyboard.isStaleWithPendingKeyEvents()
*/
func (shared *orderedKeyboardType) isStaleWithPendingKeyEvents() bool {
	shared.mutex.Lock()
	defer shared.mutex.Unlock()
	return len(shared.pendingKeyEvents) > 0 && !shared.isHoldingLocked()
}

/*
popPendingKeyEvent is a method which allows you to remove and return the oldest key event being held back. The second
return value is false, and the event empty, if none is being held.

Example:

	keyEvent, isFound := orderedKeyboard.popPendingKeyEvent()
*/
func (shared *orderedKeyboardType) popPendingKeyEvent() (keyEventType, bool) {
	shared.mutex.Lock()
	defer shared.mutex.Unlock()
	if len(shared.pendingKeyEvents) == 0 {
		return keyEventType{}, false
	}
	keyEvent := shared.pendingKeyEvents[0]
	shared.pendingKeyEvents[0] = keyEventType{}
	shared.pendingKeyEvents = shared.pendingKeyEvents[1:]
	if len(shared.pendingKeyEvents) == 0 {
		shared.pendingKeyEvents = nil
	}
	return keyEvent, true
}

/*
reset is a method which allows you to turn ordered keyboard mode off, forget when Inkey was last called, and discard
every key event being held back, without handling any of them.

Example:

	orderedKeyboard.reset()
*/
func (shared *orderedKeyboardType) reset() {
	shared.mutex.Lock()
	defer shared.mutex.Unlock()
	shared.isEnabled = false
	shared.lastInkeyTime = time.Time{}
	shared.pendingKeyEvents = nil
}
