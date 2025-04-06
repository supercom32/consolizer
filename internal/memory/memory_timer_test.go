package memory

import (
	"github.com/stretchr/testify/assert"
	"supercom32.net/consolizer"
	"supercom32.net/consolizer/internal/recast"
	"testing"
)

func TestAddTimer(test *testing.T) {
	timerAlias := "timerAlias1"
	consolizer.InitializeTimerMemory()
	consolizer.AddTimer(timerAlias, 1000, true)
	obtainedResult := recast.GetArrayOfInterfaces(consolizer.Timers.Entries[timerAlias].TimerLength, consolizer.Timers.Entries[timerAlias].IsTimerEnabled)
	expectedResult := recast.GetArrayOfInterfaces(int64(1000), true)
	assert.Equalf(test, expectedResult, obtainedResult, "The added timer values do not match what was expected.")
}

func TestGetTimer(test *testing.T) {
	timerAlias := "timerAlias1"
	consolizer.InitializeTimerMemory()
	consolizer.AddTimer(timerAlias, 1234, true)
	timerEntry := consolizer.GetTimer(timerAlias)
	obtainedResult := recast.GetArrayOfInterfaces(timerEntry.TimerLength, timerEntry.IsTimerEnabled)
	expectedResult := recast.GetArrayOfInterfaces(int64(1234), true)
	assert.Equalf(test, expectedResult, obtainedResult, "The created dialog attribute style did not match what was supposed to be created!")

}

func TestDeleteTimer(test *testing.T) {
	timerAlias := "timerAlias1"
	consolizer.InitializeTimerMemory()
	consolizer.AddTimer(timerAlias, 1234, true)
	consolizer.DeleteTimer(timerAlias)
	obtainedResult := len(consolizer.Timers.Entries)
	expectedResult := 0
	assert.Equalf(test, expectedResult, obtainedResult, "The number of remaining timers does not match what was expected.")
}
