package memory

import (
	"github.com/stretchr/testify/assert"
	"github.com/supercom32/consolizer/internal/recast"
	"testing"
)

func TestAddTimer(test *testing.T) {
	timerAlias := "timerAlias1"
	InitializeTimerMemory()
	AddTimer(timerAlias, 1000, true)
	obtainedResult := recast.GetArrayOfInterfaces(Timer.Entries[timerAlias].TimerLength, Timer.Entries[timerAlias].IsTimerEnabled)
	expectedResult := recast.GetArrayOfInterfaces(int64(1000), true)
	assert.Equalf(test, expectedResult, obtainedResult, "The added timer values do not match what was expected.")
}

func TestGetTimer(test *testing.T) {
	timerAlias := "timerAlias1"
	InitializeTimerMemory()
	AddTimer(timerAlias, 1234, true)
	timerEntry := GetTimer(timerAlias)
	obtainedResult := recast.GetArrayOfInterfaces(timerEntry.TimerLength, timerEntry.IsTimerEnabled)
	expectedResult := recast.GetArrayOfInterfaces(int64(1234), true)
	assert.Equalf(test, expectedResult, obtainedResult, "The created dialog attribute style did not match what was supposed to be created!")

}

func TestDeleteTimer(test *testing.T) {
	timerAlias := "timerAlias1"
	InitializeTimerMemory()
	AddTimer(timerAlias, 1234, true)
	DeleteTimer(timerAlias)
	obtainedResult := len(Timer.Entries)
	expectedResult := 0
	assert.Equalf(test, expectedResult, obtainedResult, "The number of remaining timers does not match what was expected.")
}
