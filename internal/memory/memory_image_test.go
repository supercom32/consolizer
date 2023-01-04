package memory

import (
	"github.com/supercom32/consolizer/types"
	"testing"
)

func TestCreateDeleteImage(test *testing.T) {
	InitializeImageMemory()
	imageEntry := types.NewImageEntry()
	AddImage("MyImageAlias", imageEntry)
	if Image.Entries["MyImageAlias"] == nil {
		test.Errorf("An image entry was requested to be created, but could not be found in memory!")
	}
	DeleteImage("MyImageAlias")
	if Image.Entries["MyImageAlias"] != nil {
		test.Errorf("An image was requested to be delete, but it could still be found in memory!")
	}
}
