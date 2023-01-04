package memory

import (
	"github.com/supercom32/consolizer/types"
	"testing"
)

func TestCreateDeleteButton(test *testing.T) {
	InitializeButtonMemory()
	styleEntry := types.NewTuiStyleEntry()
	AddButton("MyLayer", "MyButtonAlias", "Label", styleEntry, 0, 0, 10, 11)
	if Button.Entries["MyLayer"]["MyButtonAlias"] == nil {
		test.Errorf("A button was requested to be created, but could not be found in memory!")
	}
	DeleteButton("MyLayer", "MyButtonAlias")
	if Button.Entries["MyLayer"]["MyButtonAlias"] != nil {
		test.Errorf("A button was requested to be delete, but it could still be found in memory!")
	}
}
