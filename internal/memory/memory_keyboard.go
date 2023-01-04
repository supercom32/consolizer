package memory

import (
	"sync"
)

type keyboardMemoryType struct {
	sync.Mutex
	entries [][]rune
}

var KeyboardMemory keyboardMemoryType

func (shared *keyboardMemoryType) AddKeystrokeToKeyboardBuffer(keystroke ...[]rune) {
	shared.Lock()
	defer func() {
		shared.Unlock()
	}()
	for _, currentKeystroke := range keystroke {
		shared.entries = append(shared.entries, currentKeystroke)
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
