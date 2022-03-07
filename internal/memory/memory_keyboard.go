package memory

import (
	"sync"
)

type keyboardMemoryType struct {
	keyboardMemory [][]rune
	mutex          sync.Mutex
}

var KeyboardMemory keyboardMemoryType

func (shared *keyboardMemoryType) AddKeystrokeToKeyboardBuffer(keystroke ...[]rune) {
	shared.mutex.Lock()
	for _, currentKeystroke := range keystroke {
		shared.keyboardMemory = append(shared.keyboardMemory, currentKeystroke)
	}
	shared.mutex.Unlock()
}

func (shared *keyboardMemoryType) GetKeystrokeFromKeyboardBuffer() []rune {
	if shared.keyboardMemory == nil || len(shared.keyboardMemory) == 0 {
		return nil
	}
	var keystroke []rune
	shared.mutex.Lock()
	keystroke = shared.keyboardMemory[0]
	shared.keyboardMemory = shared.keyboardMemory[1:]
	shared.mutex.Unlock()
	return keystroke
}
