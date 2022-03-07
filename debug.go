package consolizer

import (
	"github.com/supercom32/consolizer/internal/stringformat"
	"github.com/supercom32/filesystem"
	"os"
)

func dumpScreenComparisons(originalScreenAsBase64 string, expectedScreenAsBase64 string) {
	originalScreen := stringformat.GetStringFromBase64(originalScreenAsBase64)
	os.WriteFile("/tmp/original.txt", []byte(originalScreen), 0644)
	expectedScreen := stringformat.GetStringFromBase64(expectedScreenAsBase64)
	os.WriteFile("/tmp/expected.txt", []byte(expectedScreen), 0644)
}

func logInfo(info string) {
	filesystem.AppendLineToFile("/tmp/debug.log", info + "\n", 0)
}