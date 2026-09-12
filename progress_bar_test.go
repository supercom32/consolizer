package consolizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
TestProgressBarTransparentLabelBackgroundFollowsFill is a test which verifies that with isBackgroundTransparent
set the label cells take the filled background colour where the fill has reached and the unfilled background
colour beyond it, across both halves of a wide rune, rather than a single fixed label background.

Example:
    Expected Inputs:
        A width 20, value 10 of 20 horizontal progress bar (fill boundary at column 10) at layer position
        (2,2), label "确定进度AB" (ten printed columns, centred so it starts at screen column 7), drawn with
        isBackgroundTransparent true. Filled background is a green, unfilled a grey, fixed text background a blue.
    Expected Outputs:
        The lead and placeholder cells of the leading wide rune (screen columns 7 and 8, inside the fill) both
        carry the green filled background and are not flagged transparent; the trailing "A" and "B" cells
        (screen columns 15 and 16, past the fill) carry the grey unfilled background. None carry the blue fixed
        text background.
*/
func TestProgressBarTransparentLabelBackgroundFollowsFill(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	filledBackgroundColor := GetRGBColor(0, 200, 0)
	unfilledBackgroundColor := GetRGBColor(90, 90, 90)
	fixedTextBackgroundColor := GetRGBColor(0, 0, 200)
	styleEntry.ProgressBar.FilledBackgroundColor = filledBackgroundColor
	styleEntry.ProgressBar.UnfilledBackgroundColor = unfilledBackgroundColor
	styleEntry.ProgressBar.TextBackgroundColor = fixedTextBackgroundColor
	styleEntry.ProgressBar.TextForegroundColor = GetRGBColor(255, 255, 255)

	layer1.AddProgressBar("确定进度AB", styleEntry, 2, 2, 20, 1, false, 10, 20, true)
	UpdateDisplay(false)
	characterMemory := commonResource.screenLayer.CharacterMemory

	assert.Equal(test, '确', characterMemory[2][7].Character)
	assert.Equal(test, filledBackgroundColor, characterMemory[2][7].AttributeEntry.BackgroundColor, "wide rune lead cell inside the fill")
	assert.Equal(test, filledBackgroundColor, characterMemory[2][8].AttributeEntry.BackgroundColor, "wide rune placeholder cell inside the fill")
	assert.False(test, characterMemory[2][7].AttributeEntry.IsBackgroundTransparent, "the label cell background is set explicitly")

	assert.Equal(test, 'A', characterMemory[2][15].Character)
	assert.Equal(test, unfilledBackgroundColor, characterMemory[2][15].AttributeEntry.BackgroundColor, "narrow rune past the fill")
	assert.Equal(test, unfilledBackgroundColor, characterMemory[2][16].AttributeEntry.BackgroundColor, "narrow rune past the fill")

	assert.NotEqual(test, fixedTextBackgroundColor, characterMemory[2][7].AttributeEntry.BackgroundColor)
	assert.NotEqual(test, fixedTextBackgroundColor, characterMemory[2][15].AttributeEntry.BackgroundColor)
}

/*
TestProgressBarOpaqueLabelBackgroundUnchanged is a test which guards that with isBackgroundTransparent NOT set
every label cell still carries the single fixed text background colour, so the previous behaviour is untouched.

Example:
    Expected Inputs:
        The same width 20, value 10 of 20 progress bar and "确定进度AB" label as the transparent case, but drawn
        with isBackgroundTransparent false.
    Expected Outputs:
        Both a label cell inside the fill (screen column 7) and one past it (screen column 15) carry the blue
        fixed text background colour, not the filled or unfilled bar colours.
*/
func TestProgressBarOpaqueLabelBackgroundUnchanged(test *testing.T) {
	layer1, _, _, styleEntry := CommonTestSetup()
	filledBackgroundColor := GetRGBColor(0, 200, 0)
	unfilledBackgroundColor := GetRGBColor(90, 90, 90)
	fixedTextBackgroundColor := GetRGBColor(0, 0, 200)
	styleEntry.ProgressBar.FilledBackgroundColor = filledBackgroundColor
	styleEntry.ProgressBar.UnfilledBackgroundColor = unfilledBackgroundColor
	styleEntry.ProgressBar.TextBackgroundColor = fixedTextBackgroundColor
	styleEntry.ProgressBar.TextForegroundColor = GetRGBColor(255, 255, 255)

	layer1.AddProgressBar("确定进度AB", styleEntry, 2, 2, 20, 1, false, 10, 20, false)
	UpdateDisplay(false)
	characterMemory := commonResource.screenLayer.CharacterMemory

	assert.Equal(test, fixedTextBackgroundColor, characterMemory[2][7].AttributeEntry.BackgroundColor, "label cell inside the fill keeps the fixed text background")
	assert.Equal(test, fixedTextBackgroundColor, characterMemory[2][15].AttributeEntry.BackgroundColor, "label cell past the fill keeps the fixed text background")
}
