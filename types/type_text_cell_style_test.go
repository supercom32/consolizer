package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTextStyleTypeCreation(test *testing.T) {
	firstAttributeEntry := NewTextCellStyleEntry()
	firstAttributeEntry.ForegroundColor = 1234
	firstAttributeEntry.BackgroundColor = 5678
	firstAttributeEntry.IsBlinking = false
	firstAttributeEntry.IsBold = true
	firstAttributeEntry.IsReversed = false
	firstAttributeEntry.IsUnderlined = true
	secondAttributeEntry := NewTextCellStyleEntry(&firstAttributeEntry)
	assert.Equalf(test, secondAttributeEntry, firstAttributeEntry, "The first text style entry is not the same as the second, even though it should be an identical clone.")

	firstAttributeEntry.ForegroundColor = 1234
	firstAttributeEntry.BackgroundColor = 5678
	firstAttributeEntry.IsBlinking = true
	firstAttributeEntry.IsBold = false
	firstAttributeEntry.IsReversed = true
	firstAttributeEntry.IsUnderlined = false
	assert.NotEqualf(test, secondAttributeEntry, firstAttributeEntry, "The second text style entry should not be the same as a the first, as manipulating it should only effect itself.")
}
