package consolizer

import (
	"sync"
)

type keyboardMemoryType struct {
	sync.Mutex
	entries [][]rune
	// Map to track the current state of each key (pressed or released)
	keyStates map[string]bool
	// Store the last keystroke for reference
	lastKeystroke []rune
}

var KeyboardMemory keyboardMemoryType

func init() {
	// Initialize the keyStates map
	KeyboardMemory.keyStates = make(map[string]bool)
}

/*
AddToBuffer is a method which allows you to add one or more keystrokes to the keyboard buffer. Each
keystroke is represented as a slice of runes.

:param keystroke: One or more keystrokes to be added to the buffer.

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
GetFromBuffer is a method which allows you to retrieve the next available keystroke from the keyboard
buffer. The keystroke is returned as a slice of runes.

:return: The next keystroke in the buffer, or nil if the buffer is empty.

Example:

	keystroke := KeyboardMemory.GetFromBuffer()
*/
func (shared *keyboardMemoryType) GetFromBuffer() []rune {
	if shared.entries == nil || len(shared.entries) == 0 {
		return nil
	}
	var keystroke []rune
	shared.Lock()
	defer func() {
		shared.Unlock()
	}()
	keystroke = shared.entries[0]
	shared.entries = shared.entries[1:]
	return keystroke
}

/*
LiveInkey is a method which allows you to retrieve the key currently being pressed. If no key is currently pressed, an
empty string is returned.

:return: The string representation of the key currently being pressed.

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

:param key: The string representation of the key to check.

:return: A boolean indicating whether the specified key is currently pressed.

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

:return: The last keystroke as a slice of runes.

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
