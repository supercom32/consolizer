package memory

import (
	"github.com/supercom32/consolizer/types"
)

// FileMenus is a control memory manager for file menu entries
var FileMenus = NewControlMemoryManager[types.FileMenuEntryType]()
