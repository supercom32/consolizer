package memory

import (
	"supercom32.net/consolizer/types"
	"testing"
)

func TestCreateDeleteButton(test *testing.T) {
	// No need to call InitializeButtonMemory() anymore as it's handled by ControlMemoryManager initialization.
	styleEntry := types.NewTuiStyleEntry()

	// Add a button using the new memory manager
	AddButton("MyLayer", "MyButtonAlias", "Label", styleEntry, 0, 0, 10, 11)

	// Check if the button was added
	button := GetButton("MyLayer", "MyButtonAlias")
	if button == nil {
		test.Errorf("A button was requested to be created, but could not be found in memory!")
	}

	// Delete the button using the new memory manager
	DeleteButton("MyLayer", "MyButtonAlias")

	// Check if the button was deleted
	if IsButtonExists("MyLayer", "MyButtonAlias") {
		test.Errorf("A button was requested to be deleted, but it could still be found in memory!")
	}
}
