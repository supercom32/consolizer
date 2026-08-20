package consolizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/types"
)

/*
makeBlockCharacterEntry is a test helper which builds a character entry for block compositor tests. It creates a cell
with the given glyph, foreground color, background color, and transparency flags, leaving all other attributes at their
defaults.
*/
func makeBlockCharacterEntry(character rune, foregroundColor constants.ColorType, backgroundColor constants.ColorType, isForegroundTransparent bool, isBackgroundTransparent bool) types.CharacterEntryType {
	attributeEntry := types.NewAttributeEntry()
	attributeEntry.ForegroundColor = foregroundColor
	attributeEntry.BackgroundColor = backgroundColor
	attributeEntry.IsForegroundTransparent = isForegroundTransparent
	attributeEntry.IsBackgroundTransparent = isBackgroundTransparent
	return types.CharacterEntryType{Character: character, AttributeEntry: attributeEntry}
}

/*
getCellPixelColors is a test helper which rasterizes a block element cell into its 64 per-pixel colors using the block
element mask table. It returns the pixel array indexed by bit position along with whether the glyph was a known block
element. Expected output for a lower half block with a green foreground and red background is red in pixels 0 through
31 and green in pixels 32 through 63.
*/
func getCellPixelColors(character rune, foregroundColor constants.ColorType, backgroundColor constants.ColorType) ([64]constants.ColorType, bool) {
	var pixelColors [64]constants.ColorType
	mask, isKnownBlockElement := blockElementMasks[character]
	if !isKnownBlockElement {
		return pixelColors, false
	}
	for pixelIndex := 0; pixelIndex < 64; pixelIndex++ {
		if mask&(uint64(1)<<uint(pixelIndex)) != 0 {
			pixelColors[pixelIndex] = foregroundColor
		} else {
			pixelColors[pixelIndex] = backgroundColor
		}
	}
	return pixelColors, true
}

/*
getPixelErrorAgainstExpected is a test helper which computes the total Manhattan color error between a rendered block
element cell and an expected array of 64 pixel colors. It is an independent per-pixel reference used to verify that the
best-fit search returns an error-minimizing result.
*/
func getPixelErrorAgainstExpected(character rune, foregroundColor constants.ColorType, backgroundColor constants.ColorType, expectedPixelColors [64]constants.ColorType) int64 {
	actualPixelColors, _ := getCellPixelColors(character, foregroundColor, backgroundColor)
	totalError := int64(0)
	for pixelIndex := 0; pixelIndex < 64; pixelIndex++ {
		expectedRed, expectedGreen, expectedBlue := GetRGBColorComponents(expectedPixelColors[pixelIndex])
		actualRed, actualGreen, actualBlue := GetRGBColorComponents(actualPixelColors[pixelIndex])
		totalError += colorComponentDistance(int64(expectedRed), int64(expectedGreen), int64(expectedBlue), int64(actualRed), int64(actualGreen), int64(actualBlue))
	}
	return totalError
}

/*
getMinimumAchievableError is a test helper which brute-forces the lowest possible per-pixel error any candidate glyph
could achieve against an expected array of 64 pixel colors. For every candidate it assigns the optimal average color to
each side of the glyph via a per-pixel loop, providing an independent reference for the popcount-based search.
*/
func getMinimumAchievableError(expectedPixelColors [64]constants.ColorType) int64 {
	minimumError := int64(1 << 62)
	for _, candidate := range blockElementCandidates {
		var setRed, setGreen, setBlue, setCount int64
		var unsetRed, unsetGreen, unsetBlue, unsetCount int64
		for pixelIndex := 0; pixelIndex < 64; pixelIndex++ {
			red, green, blue := GetRGBColorComponents(expectedPixelColors[pixelIndex])
			if candidate.mask&(uint64(1)<<uint(pixelIndex)) != 0 {
				setRed += int64(red)
				setGreen += int64(green)
				setBlue += int64(blue)
				setCount++
			} else {
				unsetRed += int64(red)
				unsetGreen += int64(green)
				unsetBlue += int64(blue)
				unsetCount++
			}
		}
		if setCount > 0 {
			setRed, setGreen, setBlue = setRed/setCount, setGreen/setCount, setBlue/setCount
		}
		if unsetCount > 0 {
			unsetRed, unsetGreen, unsetBlue = unsetRed/unsetCount, unsetGreen/unsetCount, unsetBlue/unsetCount
		}
		candidateError := int64(0)
		for pixelIndex := 0; pixelIndex < 64; pixelIndex++ {
			red, green, blue := GetRGBColorComponents(expectedPixelColors[pixelIndex])
			if candidate.mask&(uint64(1)<<uint(pixelIndex)) != 0 {
				candidateError += colorComponentDistance(int64(red), int64(green), int64(blue), setRed, setGreen, setBlue)
			} else {
				candidateError += colorComponentDistance(int64(red), int64(green), int64(blue), unsetRed, unsetGreen, unsetBlue)
			}
		}
		if candidateError < minimumError {
			minimumError = candidateError
		}
	}
	return minimumError
}

/*
TestBlockElementMaskTableIntegrity is a test which verifies that the derived block element mask table stays consistent
with the base bitmask constants. It checks that every base bitmask is present, that derived complementary glyphs such
as the upper half block are exact inverses of their counterparts, and that the full block and space entries cover all
and none of the cell respectively.
*/
func TestBlockElementMaskTableIntegrity(test *testing.T) {
	for character, mask := range constants.CharBlockBitmasks {
		tableMask, isPresent := blockElementMasks[character]
		assert.True(test, isPresent, "base bitmask glyph missing from table: %c", character)
		assert.Equal(test, mask, tableMask, "base bitmask altered for glyph: %c", character)
	}
	assert.Equal(test, ^blockElementMasks[constants.CharBlockLowerHalf], blockElementMasks[constants.CharBlockUpperHalf])
	assert.Equal(test, ^blockElementMasks[constants.CharBlockLeftHalf], blockElementMasks[constants.CharBlockRightHalf])
	assert.Equal(test, ^blockElementMasks[constants.CharBlockQuadrantUpperRight], blockElementMasks[constants.CharBlockQuadrantUpperLeftAndLowerLeftAndLowerRight])
	assert.Equal(test, fullBlockElementMask, blockElementMasks[constants.CharBlockFull])
	assert.Equal(test, uint64(0), blockElementMasks[' '])
	assert.NotEmpty(test, blockElementCandidates)
}

/*
TestCompositeBlockElementCellsRevealsUpperHalfForeground is a test which verifies that a lower half block with a
transparent background reveals the true upper half content of the underlying cell. The source is a lower half block
with a green foreground over a target upper half block with a red foreground and blue background, so the expected
output pixels are red in the upper half and green in the lower half instead of the old behavior of blue in the upper
half.
*/
func TestCompositeBlockElementCellsRevealsUpperHalfForeground(test *testing.T) {
	colorRed := GetRGBColor(255, 0, 0)
	colorGreen := GetRGBColor(0, 255, 0)
	colorBlue := GetRGBColor(0, 0, 255)
	sourceEntry := makeBlockCharacterEntry(constants.CharBlockLowerHalf, colorGreen, colorBlue, false, true)
	targetEntry := makeBlockCharacterEntry(constants.CharBlockUpperHalf, colorRed, colorBlue, false, false)

	character, foregroundColor, backgroundColor, isComposited := compositeBlockElementCells(&sourceEntry, &targetEntry)
	assert.True(test, isComposited)

	actualPixelColors, isKnown := getCellPixelColors(character, foregroundColor, backgroundColor)
	assert.True(test, isKnown)
	upperHalfMask := blockElementMasks[constants.CharBlockUpperHalf]
	for pixelIndex := 0; pixelIndex < 64; pixelIndex++ {
		if upperHalfMask&(uint64(1)<<uint(pixelIndex)) != 0 {
			assert.Equal(test, colorRed, actualPixelColors[pixelIndex], "upper half pixel %d should reveal target foreground", pixelIndex)
		} else {
			assert.Equal(test, colorGreen, actualPixelColors[pixelIndex], "lower half pixel %d should keep source foreground", pixelIndex)
		}
	}
}

/*
TestCompositeBlockElementCellsRevealsFullBlockForeground is a test which verifies that transparency over a full block
target reveals the target's foreground color rather than its background color. The source is a lower half block with a
green foreground over a target full block with a red foreground and blue background, so all revealed upper half pixels
must be red even though the old attribute-level substitution would have produced blue.
*/
func TestCompositeBlockElementCellsRevealsFullBlockForeground(test *testing.T) {
	colorRed := GetRGBColor(255, 0, 0)
	colorGreen := GetRGBColor(0, 255, 0)
	colorBlue := GetRGBColor(0, 0, 255)
	sourceEntry := makeBlockCharacterEntry(constants.CharBlockLowerHalf, colorGreen, colorBlue, false, true)
	targetEntry := makeBlockCharacterEntry(constants.CharBlockFull, colorRed, colorBlue, false, false)

	character, foregroundColor, backgroundColor, isComposited := compositeBlockElementCells(&sourceEntry, &targetEntry)
	assert.True(test, isComposited)

	actualPixelColors, isKnown := getCellPixelColors(character, foregroundColor, backgroundColor)
	assert.True(test, isKnown)
	upperHalfMask := blockElementMasks[constants.CharBlockUpperHalf]
	for pixelIndex := 0; pixelIndex < 64; pixelIndex++ {
		if upperHalfMask&(uint64(1)<<uint(pixelIndex)) != 0 {
			assert.Equal(test, colorRed, actualPixelColors[pixelIndex], "revealed pixel %d should show full block foreground", pixelIndex)
		} else {
			assert.Equal(test, colorGreen, actualPixelColors[pixelIndex], "opaque pixel %d should keep source foreground", pixelIndex)
		}
	}
}

/*
TestCompositeBlockElementCellsBestFitIsErrorMinimizing is a test which verifies that a three-color composite resolves
to the glyph and color pair with the lowest achievable per-pixel error. The source is a left half block with a green
foreground and transparent background over a target upper half block with a red foreground and blue background, so the
expected pixel content is green on the left, red in the upper right quadrant, and blue in the lower right quadrant,
which no single glyph can reproduce exactly.
*/
func TestCompositeBlockElementCellsBestFitIsErrorMinimizing(test *testing.T) {
	colorRed := GetRGBColor(255, 0, 0)
	colorGreen := GetRGBColor(0, 255, 0)
	colorBlue := GetRGBColor(0, 0, 255)
	sourceEntry := makeBlockCharacterEntry(constants.CharBlockLeftHalf, colorGreen, colorBlue, false, true)
	targetEntry := makeBlockCharacterEntry(constants.CharBlockUpperHalf, colorRed, colorBlue, false, false)

	var expectedPixelColors [64]constants.ColorType
	leftHalfMask := blockElementMasks[constants.CharBlockLeftHalf]
	upperHalfMask := blockElementMasks[constants.CharBlockUpperHalf]
	for pixelIndex := 0; pixelIndex < 64; pixelIndex++ {
		pixelBit := uint64(1) << uint(pixelIndex)
		switch {
		case leftHalfMask&pixelBit != 0:
			expectedPixelColors[pixelIndex] = colorGreen
		case upperHalfMask&pixelBit != 0:
			expectedPixelColors[pixelIndex] = colorRed
		default:
			expectedPixelColors[pixelIndex] = colorBlue
		}
	}

	character, foregroundColor, backgroundColor, isComposited := compositeBlockElementCells(&sourceEntry, &targetEntry)
	assert.True(test, isComposited)

	actualError := getPixelErrorAgainstExpected(character, foregroundColor, backgroundColor, expectedPixelColors)
	minimumError := getMinimumAchievableError(expectedPixelColors)
	assert.Equal(test, minimumError, actualError, "best-fit result should achieve the minimum possible pixel error")
}

/*
TestCompositeBlockElementCellsTextTargetFallsBackToBackground is a test which verifies that transparency over a
non-block target reveals the target's background color, matching the previous attribute-level behavior. The source is
a lower half block with a green foreground over a target letter A with a white foreground and black background, so the
revealed upper half pixels must be black.
*/
func TestCompositeBlockElementCellsTextTargetFallsBackToBackground(test *testing.T) {
	colorGreen := GetRGBColor(0, 255, 0)
	colorWhite := GetRGBColor(255, 255, 255)
	colorBlack := GetRGBColor(0, 0, 0)
	sourceEntry := makeBlockCharacterEntry(constants.CharBlockLowerHalf, colorGreen, colorBlack, false, true)
	targetEntry := makeBlockCharacterEntry('A', colorWhite, colorBlack, false, false)

	character, foregroundColor, backgroundColor, isComposited := compositeBlockElementCells(&sourceEntry, &targetEntry)
	assert.True(test, isComposited)

	actualPixelColors, isKnown := getCellPixelColors(character, foregroundColor, backgroundColor)
	assert.True(test, isKnown)
	lowerHalfMask := blockElementMasks[constants.CharBlockLowerHalf]
	for pixelIndex := 0; pixelIndex < 64; pixelIndex++ {
		if lowerHalfMask&(uint64(1)<<uint(pixelIndex)) != 0 {
			assert.Equal(test, colorGreen, actualPixelColors[pixelIndex])
		} else {
			assert.Equal(test, colorBlack, actualPixelColors[pixelIndex])
		}
	}
}

/*
TestCompositeBlockElementCellsRejectsNonBlockSource is a test which verifies that a source cell whose character is not
a known block element is rejected so callers fall back to attribute-level compositing. The source is the letter A with
a transparent background over an upper half block target, and the expected result is a failed composite.
*/
func TestCompositeBlockElementCellsRejectsNonBlockSource(test *testing.T) {
	colorRed := GetRGBColor(255, 0, 0)
	colorBlue := GetRGBColor(0, 0, 255)
	sourceEntry := makeBlockCharacterEntry('A', colorRed, colorBlue, false, true)
	targetEntry := makeBlockCharacterEntry(constants.CharBlockUpperHalf, colorRed, colorBlue, false, false)

	_, _, _, isComposited := compositeBlockElementCells(&sourceEntry, &targetEntry)
	assert.False(test, isComposited)
}

/*
TestCompositeBlockElementCellsRejectsOpaqueSource is a test which verifies that a block element source without any
effective transparency is rejected. A full block with only the background transparency flag set exposes no pixels
because the glyph covers the whole cell, so the composite must report failure and leave the fast path untouched.
*/
func TestCompositeBlockElementCellsRejectsOpaqueSource(test *testing.T) {
	colorRed := GetRGBColor(255, 0, 0)
	colorBlue := GetRGBColor(0, 0, 255)
	sourceEntry := makeBlockCharacterEntry(constants.CharBlockFull, colorRed, colorBlue, false, true)
	targetEntry := makeBlockCharacterEntry(constants.CharBlockUpperHalf, colorRed, colorBlue, false, false)

	_, _, _, isComposited := compositeBlockElementCells(&sourceEntry, &targetEntry)
	assert.False(test, isComposited)

	sourceEntry = makeBlockCharacterEntry(constants.CharBlockLowerHalf, colorRed, colorBlue, false, false)
	_, _, _, isComposited = compositeBlockElementCells(&sourceEntry, &targetEntry)
	assert.False(test, isComposited)
}

/*
TestCompositeBlockElementCellsUniformResult is a test which verifies that a composite collapsing to a single color is
returned as a full block in that color. The source is a lower half block with a red foreground and transparent
background over a target whose background is the same red, so every output pixel must be red.
*/
func TestCompositeBlockElementCellsUniformResult(test *testing.T) {
	colorRed := GetRGBColor(255, 0, 0)
	colorWhite := GetRGBColor(255, 255, 255)
	sourceEntry := makeBlockCharacterEntry(constants.CharBlockLowerHalf, colorRed, colorRed, false, true)
	targetEntry := makeBlockCharacterEntry('A', colorWhite, colorRed, false, false)

	character, foregroundColor, backgroundColor, isComposited := compositeBlockElementCells(&sourceEntry, &targetEntry)
	assert.True(test, isComposited)
	assert.Equal(test, constants.CharBlockFull, character)
	assert.Equal(test, colorRed, foregroundColor)
	assert.Equal(test, colorRed, backgroundColor)
}

/*
TestCompositeBlockElementCellsFullyTransparentSpacePassthrough is a test which verifies that a space source cell with
both transparency flags set reproduces the target cell exactly. The target is an upper half block with a red foreground
and blue background, and the expected output pixels are red in the upper half and blue in the lower half.
*/
func TestCompositeBlockElementCellsFullyTransparentSpacePassthrough(test *testing.T) {
	colorRed := GetRGBColor(255, 0, 0)
	colorGreen := GetRGBColor(0, 255, 0)
	colorBlue := GetRGBColor(0, 0, 255)
	sourceEntry := makeBlockCharacterEntry(' ', colorGreen, colorGreen, true, true)
	targetEntry := makeBlockCharacterEntry(constants.CharBlockUpperHalf, colorRed, colorBlue, false, false)

	character, foregroundColor, backgroundColor, isComposited := compositeBlockElementCells(&sourceEntry, &targetEntry)
	assert.True(test, isComposited)

	actualPixelColors, isKnown := getCellPixelColors(character, foregroundColor, backgroundColor)
	assert.True(test, isKnown)
	expectedPixelColors, _ := getCellPixelColors(constants.CharBlockUpperHalf, colorRed, colorBlue)
	assert.Equal(test, expectedPixelColors, actualPixelColors)
}

/*
TestCompositeCellResolvesSubCellTransparency is a test which verifies the integration of sub-cell compositing into the
cell compositing pipeline. A lower half block source with a transparent background composited at full opacity in a
final composite over an upper half block target must produce a cell whose upper half shows the target's red foreground
and whose transparency flags are consumed, while a partial blend alpha must fall back to the previous attribute-level
substitution behavior.
*/
func TestCompositeCellResolvesSubCellTransparency(test *testing.T) {
	colorRed := GetRGBColor(255, 0, 0)
	colorGreen := GetRGBColor(0, 255, 0)
	colorBlue := GetRGBColor(0, 0, 255)
	sourceEntry := makeBlockCharacterEntry(constants.CharBlockLowerHalf, colorGreen, colorBlue, false, true)
	targetEntry := makeBlockCharacterEntry(constants.CharBlockUpperHalf, colorRed, colorBlue, false, false)

	resultEntry := compositeCell(&sourceEntry, &targetEntry, 1.0, 1.0, false, false, constants.TransparencyStrategyNone, true)
	actualPixelColors, isKnown := getCellPixelColors(resultEntry.Character, resultEntry.AttributeEntry.ForegroundColor, resultEntry.AttributeEntry.BackgroundColor)
	assert.True(test, isKnown)
	upperHalfMask := blockElementMasks[constants.CharBlockUpperHalf]
	for pixelIndex := 0; pixelIndex < 64; pixelIndex++ {
		if upperHalfMask&(uint64(1)<<uint(pixelIndex)) != 0 {
			assert.Equal(test, colorRed, actualPixelColors[pixelIndex])
		} else {
			assert.Equal(test, colorGreen, actualPixelColors[pixelIndex])
		}
	}
	assert.False(test, resultEntry.AttributeEntry.IsForegroundTransparent, "sub-cell resolution should consume the foreground transparency flag")
	assert.False(test, resultEntry.AttributeEntry.IsBackgroundTransparent, "sub-cell resolution should consume the background transparency flag")

	partialAlphaResult := compositeCell(&sourceEntry, &targetEntry, 0.5, 0.5, false, false, constants.TransparencyStrategyNone, true)
	assert.Equal(test, sourceEntry.Character, partialAlphaResult.Character, "partial alpha blending should not use sub-cell resolution")
}

/*
TestCompositeCellIntermediateCompositeKeepsFlags is a test which verifies that a non-final composite never re-encodes
a partially transparent block element cell. A lower half block source with a transparent background composited over an
upper half block target in an intermediate composite must keep its original glyph, keep its transparency flags set for
later compositing passes, and receive the target's background color through plain attribute substitution.
*/
func TestCompositeCellIntermediateCompositeKeepsFlags(test *testing.T) {
	colorRed := GetRGBColor(255, 0, 0)
	colorGreen := GetRGBColor(0, 255, 0)
	colorBlue := GetRGBColor(0, 0, 255)
	colorBlack := GetRGBColor(0, 0, 0)
	sourceEntry := makeBlockCharacterEntry(constants.CharBlockLowerHalf, colorGreen, colorBlack, false, true)
	targetEntry := makeBlockCharacterEntry(constants.CharBlockUpperHalf, colorRed, colorBlue, false, false)

	resultEntry := compositeCell(&sourceEntry, &targetEntry, 1.0, 1.0, false, false, constants.TransparencyStrategyNone, false)
	assert.Equal(test, sourceEntry.Character, resultEntry.Character, "intermediate composites should keep the source glyph")
	assert.True(test, resultEntry.AttributeEntry.IsBackgroundTransparent, "intermediate composites should keep transparency flags")
	assert.Equal(test, colorGreen, resultEntry.AttributeEntry.ForegroundColor)
	assert.Equal(test, colorBlue, resultEntry.AttributeEntry.BackgroundColor, "intermediate composites should substitute the target background color")
}

/*
BenchmarkCompositeBlockElementCells is a test which measures the cost of a worst-case sub-cell composite that cannot be
resolved through the uniform or exact-match fast paths and must run the full best-fit search across all candidates.
*/
func BenchmarkCompositeBlockElementCells(benchmark *testing.B) {
	colorRed := GetRGBColor(255, 0, 0)
	colorGreen := GetRGBColor(0, 255, 0)
	colorBlue := GetRGBColor(0, 0, 255)
	sourceEntry := makeBlockCharacterEntry(constants.CharBlockLeftHalf, colorGreen, colorBlue, false, true)
	targetEntry := makeBlockCharacterEntry(constants.CharBlockUpperHalf, colorRed, colorBlue, false, false)
	benchmark.ResetTimer()
	for iteration := 0; iteration < benchmark.N; iteration++ {
		compositeBlockElementCells(&sourceEntry, &targetEntry)
	}
}
