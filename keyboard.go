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

func (shared *keyboardMemoryType) AddKeystrokeToKeyboardBuffer(keystroke ...[]rune) {
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

func (shared *keyboardMemoryType) GetKeystrokeFromKeyboardBuffer() []rune {
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

// LiveInkey returns the key currently being pressed, or an empty string if no key is pressed
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

// IsKeyPressed checks if a specific key is currently pressed
func (shared *keyboardMemoryType) IsKeyPressed(key string) bool {
	shared.Lock()
	defer shared.Unlock()

	pressed, exists := shared.keyStates[key]
	return exists && pressed
}

// GetLastKeystroke returns the last keystroke that was processed
func (shared *keyboardMemoryType) GetLastKeystroke() []rune {
	shared.Lock()
	defer shared.Unlock()

	return shared.lastKeystroke
}

// ClearKeyStates marks all keys as released
func (shared *keyboardMemoryType) ClearKeyStates() {
	shared.Lock()
	defer shared.Unlock()

	for key := range shared.keyStates {
		shared.keyStates[key] = false
	}
}
