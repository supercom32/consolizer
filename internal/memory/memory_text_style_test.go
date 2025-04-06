package memory

import (
	"github.com/stretchr/testify/assert"
	"supercom32.net/consolizer"
	"supercom32.net/consolizer/constants"
	"supercom32.net/consolizer/internal/recast"
	"supercom32.net/consolizer/types"
	"testing"
)

func TestAddTextStyle(test *testing.T) {
	expectedAlias := "MyCustomAttribute"
	expectedForegroundColor := constants.AnsiColorByIndex[constants.ColorRed]
	expectedBackgroundColor := constants.AnsiColorByIndex[constants.ColorBrightGreen]
	expectedIsBlinking := true
	InitializeTextStyleMemory()
	attributeEntry := types.NewTextCellStyleEntry()
	attributeEntry.ForegroundColor = expectedForegroundColor
	attributeEntry.BackgroundColor = expectedBackgroundColor
	attributeEntry.IsBlinking = expectedIsBlinking
	consolizer.AddTextStyle(expectedAlias, attributeEntry)
	expectedResult := recast.GetArrayOfInterfaces(expectedForegroundColor, expectedBackgroundColor, expectedIsBlinking)
	obtainedAttributeEntry := consolizer.GetTextStyle(expectedAlias)
	obtainedResult := recast.GetArrayOfInterfaces(obtainedAttributeEntry.ForegroundColor, obtainedAttributeEntry.BackgroundColor, obtainedAttributeEntry.IsBlinking)
	assert.Equalf(test, expectedResult, obtainedResult, "The created dialog attribute style did not match what was supposed to be created!")
}

func TestDeleteTextStyle(test *testing.T) {
	expectedAlias := "MyCustomAttribute"
	expectedForegroundColor := constants.AnsiColorByIndex[constants.ColorRed]
	expectedBackgroundColor := constants.AnsiColorByIndex[constants.ColorBrightGreen]
	expectedIsBlinking := true
	InitializeTextStyleMemory()
	attributeEntry := types.NewTextCellStyleEntry()
	attributeEntry.ForegroundColor = expectedForegroundColor
	attributeEntry.BackgroundColor = expectedBackgroundColor
	attributeEntry.IsBlinking = expectedIsBlinking
	consolizer.AddTextStyle(expectedAlias, attributeEntry)
	consolizer.DeleteTextStyle(expectedAlias)
	assert.Panics(test, func() { consolizer.GetTextStyle("expectedAlias") }, "The created dialog attribute style did not match what was supposed to be created!")
	// assert.Equalf(test, (*AttributeEntryType)(nil), TextStyleMemory[expectedAlias], "The created dialog attribute style did not match what was supposed to be created!")
}
