package memory

import (
	"supercom32.net/consolizer"
	"testing"
)

func TestKeyboardMemory(test *testing.T) {
	consolizer.KeyboardMemory.AddKeystrokeToKeyboardBuffer([]rune{'a'})
	consolizer.KeyboardMemory.AddKeystrokeToKeyboardBuffer([]rune{'b'})
	consolizer.KeyboardMemory.AddKeystrokeToKeyboardBuffer([]rune{'c'})
	if consolizer.KeyboardMemory.GetKeystrokeFromKeyboardBuffer()[0] != 'a' {
		test.Errorf("The first keyboard Character was not returned when it should be next in queue.")
	}
	if consolizer.KeyboardMemory.GetKeystrokeFromKeyboardBuffer()[0] != 'b' {
		test.Errorf("The second keyboard Character was not returned when it should be next in queue.")
	}
	if consolizer.KeyboardMemory.GetKeystrokeFromKeyboardBuffer()[0] != 'c' {
		test.Errorf("The second keyboard Character was not returned when it should be next in queue.")
	}
	if consolizer.KeyboardMemory.GetKeystrokeFromKeyboardBuffer() != nil {
		test.Errorf("No keyboard keystrokes should have been returned, but one was given.")
	}
}
