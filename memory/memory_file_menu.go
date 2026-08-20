package memory

import (
	"github.com/supercom32/consolizer/types"
)

/*
FileMenus is a variable which contains a control memory manager for file menu entries.

Example:
    memory.FileMenus.Add(layerAlias, controlAlias, entry)
*/
var FileMenus = NewControlMemoryManager[types.FileMenuEntryType]()
