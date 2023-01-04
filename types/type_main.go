package types

import (
	"github.com/supercom32/filesystem"
)

func LogInfo(info string) {
	var stringToAppend string
	// if len(parameters) != 0 {
	//	stringToAppend = fmt.Sprintf(info, parameters...)
	// } else {
	stringToAppend = info
	// }
	filesystem.AppendLineToFile("/tmp/debug.log", stringToAppend+"\n", 0)
}
