package memory

import "testing"

func TestKeyboardMemory(test *testing.T) {
	KeyboardMemory.AddKeystrokeToKeyboardBuffer([]rune{'a'})
	KeyboardMemory.AddKeystrokeToKeyboardBuffer([]rune{'b'})
	KeyboardMemory.AddKeystrokeToKeyboardBuffer([]rune{'c'})
	if KeyboardMemory.GetKeystrokeFromKeyboardBuffer()[0] != 'a' {
		test.Errorf("The first keyboard Character was not returned when it should be next in queue.")
	}
	if KeyboardMemory.GetKeystrokeFromKeyboardBuffer()[0] != 'b' {
		test.Errorf("The second keyboard Character was not returned when it should be next in queue.")
	}
	if KeyboardMemory.GetKeystrokeFromKeyboardBuffer()[0] != 'c' {
		test.Errorf("The second keyboard Character was not returned when it should be next in queue.")
	}
	if KeyboardMemory.GetKeystrokeFromKeyboardBuffer() != nil {
		test.Errorf("No keyboard keystrokes should have been returned, but one was given.")
	}
}
