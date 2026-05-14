package consolizer

import (
	"fmt"
	"github.com/supercom32/consolizer/stringformat"
	"github.com/supercom32/filesystem"
	"os"
	"time"
)

/*
dumpScreenComparisons is a method which allows you to dump base64 encoded screen data to local files for comparison.

:param originalScreenAsBase64: The base64 string of the original screen.
:param expectedScreenAsBase64: The base64 string of the expected screen.

Example:

	dumpScreenComparisons(originalB64, expectedB64)
*/
func dumpScreenComparisons(originalScreenAsBase64 string, expectedScreenAsBase64 string) {
	originalScreen := stringformat.GetStringFromBase64(originalScreenAsBase64)
	os.WriteFile("/tmp/test_output/original.txt", []byte(originalScreen), 0644)
	expectedScreen := stringformat.GetStringFromBase64(expectedScreenAsBase64)
	os.WriteFile("/tmp/test_output/expected.txt", []byte(expectedScreen), 0644)
}

/*
LogInfo is a method which allows you to log information to a debug file with a timestamp.

:param info: The information string to log (supports format specifiers).
:param parameters: Optional parameters for the format specifiers.

Example:

	LogInfo("Value is %d", 10)
*/
func LogInfo(info string, parameters ...any) {
	var stringToAppend string
	// if len(parameters) != 0 {
	//	stringToAppend = fmt.Sprintf(info, parameters...)
	// } else {
	stringToAppend = fmt.Sprintf(info, parameters...)
	stringToAppend = time.Now().String() + " - " + stringToAppend
	// }
	filesystem.AppendLineToFile("/tmp/test_output/debug.log", stringToAppend+"\n", 0)
}
