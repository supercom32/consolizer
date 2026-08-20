package consolizer

import (
	"math"
	"math/bits"

	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/types"
)

// fullBlockElementMask represents a cell where every pixel in the 8x8 grid is covered by the foreground.
const fullBlockElementMask = ^uint64(0)

/*
blockPixelRegionType is a class which allows you to describe a group of pixels within a single character cell that all
share the same color. The mask field marks which pixels of the cell's 8x8 grid belong to the region, using the same bit
ordering as the masks found in constants.CharBlockBitmasks.
*/
type blockPixelRegionType struct {
	mask  uint64
	color constants.ColorType
	red   int32
	green int32
	blue  int32
}

/*
blockElementMatchType is a class which allows you to describe the block element glyph that exactly reproduces a given
pixel mask. The isForeground field records whether the mask corresponds to the glyph's foreground region or to its
background region, so callers know which color to assign to which attribute.
*/
type blockElementMatchType struct {
	character    rune
	isForeground bool
}

/*
blockElementCandidateType is a class which allows you to describe a single candidate glyph considered during best-fit
block element searches. Each candidate pairs a printable block element with the pixel mask of its foreground region.
*/
type blockElementCandidateType struct {
	character rune
	mask      uint64
}

var (
	blockElementMasks      map[rune]uint64
	blockElementExactMatch map[uint64]blockElementMatchType
	blockElementCandidates []blockElementCandidateType
)

func init() {
	initializeBlockElementTables()
}

/*
initializeBlockElementTables is a method which allows you to build the lookup tables used for sub-cell block element
compositing. It extends the base bitmask table from the constants package with the complementary glyphs that the image
renderers can emit, builds a reverse map for exact mask-to-glyph matching, and prepares the ordered candidate list used
by best-fit searches. In addition, the following should be noted:

- The candidate list is built in a fixed order so that tie-breaking during best-fit searches is deterministic across
  runs, since Go map iteration order is randomized.

- Complementary glyphs such as the upper half block are derived by inverting the bitmasks of their existing
  counterparts, which guarantees the table stays consistent with constants.CharBlockBitmasks.

Example:
    initializeBlockElementTables()
*/
func initializeBlockElementTables() {
	blockElementMasks = make(map[rune]uint64, len(constants.CharBlockBitmasks)+12)
	for character, bitmask := range constants.CharBlockBitmasks {
		blockElementMasks[character] = bitmask
	}
	blockElementMasks[constants.CharBlockFull] = fullBlockElementMask
	blockElementMasks[' '] = 0
	addComplementBlockElementMask(constants.CharBlockUpperHalf, constants.CharBlockLowerHalf)
	addComplementBlockElementMask(constants.CharBlockRightHalf, constants.CharBlockLeftHalf)
	addComplementBlockElementMask(constants.CharBlockUpperOneEighth, constants.CharBlockLowerSevenEighths)
	addComplementBlockElementMask(constants.CharBlockRightOneEighth, constants.CharBlockLeftSevenEighths)
	addComplementBlockElementMask(constants.CharBlockQuadrantUpperRightAndLowerLeft, constants.CharBlockQuadrantUpperLeftAndLowerRight)
	addComplementBlockElementMask(constants.CharBlockQuadrantUpperLeftAndLowerLeftAndLowerRight, constants.CharBlockQuadrantUpperRight)
	addComplementBlockElementMask(constants.CharBlockQuadrantUpperLeftAndUpperRightAndLowerLeft, constants.CharBlockQuadrantLowerRight)
	addComplementBlockElementMask(constants.CharBlockQuadrantUpperLeftAndUpperRightAndLowerRight, constants.CharBlockQuadrantLowerLeft)
	addComplementBlockElementMask(constants.CharBlockQuadrantUpperRightAndLowerLeftAndLowerRight, constants.CharBlockQuadrantUpperLeft)

	candidateOrder := []rune{
		constants.CharBlockLowerHalf,
		constants.CharBlockLeftHalf,
		constants.CharBlockQuadrantUpperLeft,
		constants.CharBlockQuadrantUpperRight,
		constants.CharBlockQuadrantLowerLeft,
		constants.CharBlockQuadrantLowerRight,
		constants.CharBlockQuadrantUpperLeftAndLowerRight,
		constants.CharBlockLowerOneQuarter,
		constants.CharBlockLowerThreeQuarters,
		constants.CharBlockLeftOneQuarter,
		constants.CharBlockLeftThreeQuarters,
		constants.CharBlockLowerOneEighth,
		constants.CharBlockLowerThreeEighths,
		constants.CharBlockLowerFiveEighths,
		constants.CharBlockLowerSevenEighths,
		constants.CharBlockLeftOneEighth,
		constants.CharBlockLeftThreeEighths,
		constants.CharBlockLeftFiveEighths,
		constants.CharBlockLeftSevenEighths,
		constants.CharBlockFull,
	}
	blockElementCandidates = make([]blockElementCandidateType, 0, len(candidateOrder))
	seenPartitions := make(map[uint64]bool, len(candidateOrder))
	for _, character := range candidateOrder {
		mask := blockElementMasks[character]
		normalizedPartition := mask
		if ^mask < mask {
			normalizedPartition = ^mask
		}
		if seenPartitions[normalizedPartition] {
			continue
		}
		seenPartitions[normalizedPartition] = true
		blockElementCandidates = append(blockElementCandidates, blockElementCandidateType{character: character, mask: mask})
	}

	exactMatchOrder := append([]rune{
		constants.CharBlockUpperHalf,
		constants.CharBlockRightHalf,
		constants.CharBlockUpperOneEighth,
		constants.CharBlockRightOneEighth,
		constants.CharBlockQuadrantUpperRightAndLowerLeft,
		constants.CharBlockQuadrantUpperLeftAndLowerLeftAndLowerRight,
		constants.CharBlockQuadrantUpperLeftAndUpperRightAndLowerLeft,
		constants.CharBlockQuadrantUpperLeftAndUpperRightAndLowerRight,
		constants.CharBlockQuadrantUpperRightAndLowerLeftAndLowerRight,
	}, candidateOrder...)
	blockElementExactMatch = make(map[uint64]blockElementMatchType, len(exactMatchOrder)*2)
	for _, character := range exactMatchOrder {
		mask := blockElementMasks[character]
		if mask == 0 || mask == fullBlockElementMask {
			continue
		}
		if _, isPresent := blockElementExactMatch[mask]; !isPresent {
			blockElementExactMatch[mask] = blockElementMatchType{character: character, isForeground: true}
		}
		if _, isPresent := blockElementExactMatch[^mask]; !isPresent {
			blockElementExactMatch[^mask] = blockElementMatchType{character: character, isForeground: false}
		}
	}
}

/*
addComplementBlockElementMask is a method which allows you to register a block element glyph whose pixel mask is the
exact inverse of an already registered counterpart glyph. This keeps derived masks consistent with the base bitmask
table without duplicating any literal bit patterns.

:param character: The glyph to register in the block element mask table.
:param complementCharacter: The already registered glyph whose inverted mask defines the new entry.

Example:
    addComplementBlockElementMask(constants.CharBlockUpperHalf, constants.CharBlockLowerHalf)
*/
func addComplementBlockElementMask(character rune, complementCharacter rune) {
	blockElementMasks[character] = ^blockElementMasks[complementCharacter]
}

/*
appendBlockPixelRegion is a method which allows you to add a pixel region to a fixed-size region list while merging
regions that share the same color. Merging keeps the region count as small as possible, which lets many composites
resolve through the cheap uniform and exact-match paths instead of a best-fit search. In addition, the following should
be noted:

- The regions array is mutated in place and the region count is updated through the provided pointer.

- The RGB components of new regions are computed once here so that best-fit searches never repeat the conversion.

:param pixelRegions: A pointer to the fixed-size array holding the accumulated regions.
:param regionCount: A pointer to the number of regions currently stored in the array.
:param mask: The pixel mask describing which cell pixels belong to the region.
:param color: The color shared by every pixel in the region.

Example:
    appendBlockPixelRegion(&pixelRegions, &regionCount, 0x0F0F0F0F0F0F0F0F, foregroundColor)
*/
func appendBlockPixelRegion(pixelRegions *[4]blockPixelRegionType, regionCount *int, mask uint64, color constants.ColorType) {
	for regionIndex := 0; regionIndex < *regionCount; regionIndex++ {
		if pixelRegions[regionIndex].color == color {
			pixelRegions[regionIndex].mask |= mask
			return
		}
	}
	red, green, blue := GetRGBColorComponents(color)
	pixelRegions[*regionCount] = blockPixelRegionType{mask: mask, color: color, red: red, green: green, blue: blue}
	*regionCount++
}

/*
compositeBlockElementCells is a method which allows you to composite a partially transparent block element cell over a
target cell with sub-cell accuracy. Instead of substituting a single underlying color for the whole transparent region,
it decomposes both cells into per-pixel color regions using their block element bitmasks, reveals the target's true
visual content beneath the transparent portion of the source, and re-encodes the combined result as the best matching
block element glyph with new foreground and background colors. In addition, the following should be noted:

- The method only applies when the source character is a known block element and at least one of its transparency
  flags actually exposes pixels, otherwise it reports failure so callers can fall back to attribute-level compositing.

- Target cells whose characters are not block elements are treated as a uniform surface of their background color,
  which matches the previous attribute-level behavior for text beneath images.

- When the combined regions collapse to two colors whose layout exactly matches a known glyph, the result is resolved
  through a single map lookup with no search being performed. The source glyph is preferred whenever it already
  represents the resulting partition, which keeps output glyphs stable across composites.

:param sourceEntry: The source character cell being drawn, expected to contain a block element glyph.
:param targetEntry: The already composited target cell underneath the source cell.
:return: The resulting glyph, its foreground color, its background color, and whether compositing was performed.

Example:
    character, foregroundColor, backgroundColor, isComposited := compositeBlockElementCells(&sourceCell, &targetCell)
*/
func compositeBlockElementCells(sourceEntry *types.CharacterEntryType, targetEntry *types.CharacterEntryType) (rune, constants.ColorType, constants.ColorType, bool) {
	sourceAttributes := sourceEntry.AttributeEntry
	sourceMask, isKnownBlockElement := blockElementMasks[sourceEntry.Character]
	if !isKnownBlockElement {
		return 0, 0, 0, false
	}
	transparentMask := uint64(0)
	if sourceAttributes.IsForegroundTransparent {
		transparentMask |= sourceMask
	}
	if sourceAttributes.IsBackgroundTransparent {
		transparentMask |= ^sourceMask
	}
	if transparentMask == 0 {
		return 0, 0, 0, false
	}

	targetAttributes := targetEntry.AttributeEntry
	targetMask, isTargetBlockElement := blockElementMasks[targetEntry.Character]
	if !isTargetBlockElement {
		targetMask = 0
	}

	var pixelRegions [4]blockPixelRegionType
	regionCount := 0
	if !sourceAttributes.IsForegroundTransparent && sourceMask != 0 {
		appendBlockPixelRegion(&pixelRegions, &regionCount, sourceMask, sourceAttributes.ForegroundColor)
	}
	if !sourceAttributes.IsBackgroundTransparent && ^sourceMask != 0 {
		appendBlockPixelRegion(&pixelRegions, &regionCount, ^sourceMask, sourceAttributes.BackgroundColor)
	}
	if revealedForegroundMask := transparentMask & targetMask; revealedForegroundMask != 0 {
		appendBlockPixelRegion(&pixelRegions, &regionCount, revealedForegroundMask, targetAttributes.ForegroundColor)
	}
	if revealedBackgroundMask := transparentMask &^ targetMask; revealedBackgroundMask != 0 {
		appendBlockPixelRegion(&pixelRegions, &regionCount, revealedBackgroundMask, targetAttributes.BackgroundColor)
	}

	if regionCount == 1 {
		uniformColor := pixelRegions[0].color
		return constants.CharBlockFull, uniformColor, uniformColor, true
	}
	if regionCount == 2 {
		// Prefer keeping the source glyph when it already represents the partition, so glyphs stay stable.
		if pixelRegions[0].mask == sourceMask {
			return sourceEntry.Character, pixelRegions[0].color, pixelRegions[1].color, true
		}
		if pixelRegions[1].mask == sourceMask {
			return sourceEntry.Character, pixelRegions[1].color, pixelRegions[0].color, true
		}
		if match, isPresent := blockElementExactMatch[pixelRegions[0].mask]; isPresent {
			if match.isForeground {
				return match.character, pixelRegions[0].color, pixelRegions[1].color, true
			}
			return match.character, pixelRegions[1].color, pixelRegions[0].color, true
		}
	}
	character, foregroundColor, backgroundColor := findBestBlockElementForPixelRegions(pixelRegions[:regionCount])
	return character, foregroundColor, backgroundColor, true
}

/*
findBestBlockElementForPixelRegions is a method which allows you to find the block element glyph and color pair that
most closely reproduces a set of colored pixel regions within a single character cell. For every candidate glyph it
derives the optimal foreground and background colors as the pixel-weighted averages of the regions covered by each side
of the glyph, then scores the candidate by the total per-pixel color error and keeps the best one. In addition, the
following should be noted:

- All pixel counting is performed with bitwise mask intersections and population counts rather than per-pixel loops,
  which keeps the cost of a full search in the sub-microsecond range.

- Ties are broken in favor of the earliest candidate in the fixed candidate order, keeping results deterministic.

:param pixelRegions: The colored pixel regions to reproduce, whose masks must partition the full 8x8 cell grid.
:return: The best matching glyph along with its derived foreground and background colors.

Example:
    character, foregroundColor, backgroundColor := findBestBlockElementForPixelRegions(pixelRegions)
*/
func findBestBlockElementForPixelRegions(pixelRegions []blockPixelRegionType) (rune, constants.ColorType, constants.ColorType) {
	bestError := int64(math.MaxInt64)
	bestCharacter := constants.CharBlockFull
	var bestSetRed, bestSetGreen, bestSetBlue int64
	var bestUnsetRed, bestUnsetGreen, bestUnsetBlue int64
	var bestSetCount, bestUnsetCount int64

	for _, candidate := range blockElementCandidates {
		var setPixelCounts, unsetPixelCounts [4]int64
		var setCount, unsetCount int64
		var setRed, setGreen, setBlue int64
		var unsetRed, unsetGreen, unsetBlue int64
		for regionIndex := range pixelRegions {
			region := &pixelRegions[regionIndex]
			setPixels := int64(bits.OnesCount64(region.mask & candidate.mask))
			unsetPixels := int64(bits.OnesCount64(region.mask &^ candidate.mask))
			setPixelCounts[regionIndex] = setPixels
			unsetPixelCounts[regionIndex] = unsetPixels
			setCount += setPixels
			unsetCount += unsetPixels
			setRed += setPixels * int64(region.red)
			setGreen += setPixels * int64(region.green)
			setBlue += setPixels * int64(region.blue)
			unsetRed += unsetPixels * int64(region.red)
			unsetGreen += unsetPixels * int64(region.green)
			unsetBlue += unsetPixels * int64(region.blue)
		}
		var setAverageRed, setAverageGreen, setAverageBlue int64
		if setCount > 0 {
			setAverageRed = setRed / setCount
			setAverageGreen = setGreen / setCount
			setAverageBlue = setBlue / setCount
		}
		var unsetAverageRed, unsetAverageGreen, unsetAverageBlue int64
		if unsetCount > 0 {
			unsetAverageRed = unsetRed / unsetCount
			unsetAverageGreen = unsetGreen / unsetCount
			unsetAverageBlue = unsetBlue / unsetCount
		}
		candidateError := int64(0)
		for regionIndex := range pixelRegions {
			region := &pixelRegions[regionIndex]
			if setPixelCounts[regionIndex] > 0 {
				candidateError += setPixelCounts[regionIndex] * colorComponentDistance(int64(region.red), int64(region.green), int64(region.blue), setAverageRed, setAverageGreen, setAverageBlue)
			}
			if unsetPixelCounts[regionIndex] > 0 {
				candidateError += unsetPixelCounts[regionIndex] * colorComponentDistance(int64(region.red), int64(region.green), int64(region.blue), unsetAverageRed, unsetAverageGreen, unsetAverageBlue)
			}
		}
		if candidateError < bestError {
			bestError = candidateError
			bestCharacter = candidate.character
			bestSetRed, bestSetGreen, bestSetBlue = setAverageRed, setAverageGreen, setAverageBlue
			bestUnsetRed, bestUnsetGreen, bestUnsetBlue = unsetAverageRed, unsetAverageGreen, unsetAverageBlue
			bestSetCount, bestUnsetCount = setCount, unsetCount
		}
	}

	if bestSetCount == 0 {
		bestSetRed, bestSetGreen, bestSetBlue = bestUnsetRed, bestUnsetGreen, bestUnsetBlue
	}
	if bestUnsetCount == 0 {
		bestUnsetRed, bestUnsetGreen, bestUnsetBlue = bestSetRed, bestSetGreen, bestSetBlue
	}
	foregroundColor := GetRGBColor(int32(bestSetRed), int32(bestSetGreen), int32(bestSetBlue))
	backgroundColor := GetRGBColor(int32(bestUnsetRed), int32(bestUnsetGreen), int32(bestUnsetBlue))
	return bestCharacter, foregroundColor, backgroundColor
}

/*
colorComponentDistance is a method which allows you to measure the difference between two colors expressed as separate
RGB components. It returns the Manhattan distance across the three channels, which is inexpensive to compute and is the
same style of metric used by the block element renderer during image loading.

:param firstRed: The red component of the first color.
:param firstGreen: The green component of the first color.
:param firstBlue: The blue component of the first color.
:param secondRed: The red component of the second color.
:param secondGreen: The green component of the second color.
:param secondBlue: The blue component of the second color.
:return: The sum of the absolute differences between each pair of color components.

Example:
    distance := colorComponentDistance(255, 0, 0, 0, 0, 255)
*/
func colorComponentDistance(firstRed, firstGreen, firstBlue, secondRed, secondGreen, secondBlue int64) int64 {
	distance := firstRed - secondRed
	if distance < 0 {
		distance = -distance
	}
	greenDistance := firstGreen - secondGreen
	if greenDistance < 0 {
		greenDistance = -greenDistance
	}
	blueDistance := firstBlue - secondBlue
	if blueDistance < 0 {
		blueDistance = -blueDistance
	}
	return distance + greenDistance + blueDistance
}
