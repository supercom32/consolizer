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

func LogInfo(info string) {
	var stringToAppend string
	// if len(parameters) != 0 {
	//	stringToAppend = fmt.Sprintf(info, parameters...)
	// } else {
	stringToAppend = info
	// }
	filesystem.AppendLineToFile("/tmp/debug.log", stringToAppend+"\n", 0)
}
