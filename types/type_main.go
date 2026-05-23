package types

import (
	"github.com/supercom32/filesystem"
)

/*
LogInfo is a method which appends information to a log file.

Example:
    LogInfo(info)
*/
func LogInfo(info string) {
	var stringToAppend string
	// if len(parameters) != 0 {
	//	stringToAppend = fmt.Sprintf(info, parameters...)
	// } else {
	stringToAppend = info
	// }
	filesystem.AppendLineToFile("/tmp/debug.log", stringToAppend+"\n", 0)
}
