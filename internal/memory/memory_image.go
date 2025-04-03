package memory

import (
	"fmt"
	"supercom32.net/consolizer/types"
	"sync"
)

type imageMemoryType struct {
	sync.Mutex
	Entries map[string]*types.ImageEntryType
}

var Image imageMemoryType

func InitializeImageMemory() {
	Image.Entries = make(map[string]*types.ImageEntryType)
}

func AddImage(imageAlias string, imageEntry types.ImageEntryType) {
	Image.Lock()
	defer func() {
		Image.Unlock()
	}()
	// verify if any errors occurred?
	Image.Entries[imageAlias] = &imageEntry
}

func GetImage(imageAlias string) *types.ImageEntryType {
	Image.Lock()
	defer func() {
		Image.Unlock()
	}()
	if Image.Entries[imageAlias] == nil {
		panic(fmt.Sprintf("The requested Image with alias '%s' could not be returned since it does not exist.", imageAlias))
	}
	return Image.Entries[imageAlias]
}
func DeleteImage(imageAlias string) {
	Image.Lock()
	defer func() {
		Image.Unlock()
	}()
	delete(Image.Entries, imageAlias)
}
func IsImageExists(imageAlias string) bool {
	Image.Lock()
	defer func() {
		Image.Unlock()
	}()
	if Image.Entries[imageAlias] == nil {
		return false
	}
	return true
}
