package stringformat

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supercom32/consolizer/constants"
)

/*
TestGetRunesThatFitInColumnCountFromStart is a test which verifies that GetRunesThatFitInColumnCountFromStart returns the longest prefix of "A中B" that fits a column
budget of 0 through 5 without ever splitting the wide rune, along with the exact column count that prefix occupies.

Example:
    Expected Inputs / Outputs for []rune("A中B") (widths 1, 2, 1):
        0 -> ("",    0) | 1 -> ("A",   1) | 2 -> ("A",    1) |
        3 -> ("A中", 3) | 4 -> ("A中B", 4) | 5 -> ("A中B", 4)
*/
func TestGetRunesThatFitInColumnCountFromStart(test *testing.T) {
	arrayOfRunes := []rune("A中B")
	tableOfCases := []struct {
		maxColumns     int
		expectedPrefix string
		expectedUsed   int
	}{
		{0, "", 0},
		{1, "A", 1},
		{2, "A", 1},
		{3, "A中", 3},
		{4, "A中B", 4},
		{5, "A中B", 4},
	}
	for _, currentCase := range tableOfCases {
		prefix, used := GetRunesThatFitInColumnCountFromStart(arrayOfRunes, currentCase.maxColumns)
		assert.Equalf(test, currentCase.expectedPrefix, string(prefix), "GetRunesThatFitInColumnCountFromStart prefix wrong for maxColumns %d.", currentCase.maxColumns)
		assert.Equalf(test, currentCase.expectedUsed, used, "GetRunesThatFitInColumnCountFromStart used count wrong for maxColumns %d.", currentCase.maxColumns)
	}
}

/*
TestGetRunesThatFitInColumnCountFromEnd is a test which verifies that GetRunesThatFitInColumnCountFromEnd returns the longest suffix of "A中B" that
fits a column budget of 0 through 5 without ever splitting the wide rune, along with the exact column count that
suffix occupies. It is the width mirror of TestGetRunesThatFitInColumnCountFromStart.

Example:
    Expected Inputs / Outputs for []rune("A中B") (widths 1, 2, 1):
        0 -> ("",    0) | 1 -> ("B",   1) | 2 -> ("B",    1) |
        3 -> ("中B", 3) | 4 -> ("A中B", 4) | 5 -> ("A中B", 4)
*/
func TestGetRunesThatFitInColumnCountFromEnd(test *testing.T) {
	arrayOfRunes := []rune("A中B")
	tableOfCases := []struct {
		maxColumns     int
		expectedSuffix string
		expectedUsed   int
	}{
		{0, "", 0},
		{1, "B", 1},
		{2, "B", 1},
		{3, "中B", 3},
		{4, "A中B", 4},
		{5, "A中B", 4},
	}
	for _, currentCase := range tableOfCases {
		suffix, used := GetRunesThatFitInColumnCountFromEnd(arrayOfRunes, currentCase.maxColumns)
		assert.Equalf(test, currentCase.expectedSuffix, string(suffix), "GetRunesThatFitInColumnCountFromEnd suffix wrong for maxColumns %d.", currentCase.maxColumns)
		assert.Equalf(test, currentCase.expectedUsed, used, "GetRunesThatFitInColumnCountFromEnd used count wrong for maxColumns %d.", currentCase.maxColumns)
	}
}

/*
TestGetRunesThatFitInColumnCountNeverSplitsWideRuneFuzz is a test which drives GetRunesThatFitInColumnCountFromStart and GetRunesThatFitInColumnCountFromEnd with random mixed
ASCII and CJK arrays and random budgets and checks that neither ever splits a wide rune, that the reported column
count matches the returned slice's printed width, that the count never exceeds the budget, and that the prefix and
suffix widths for the same budget are consistent with one another.

Example:
    For a random []rune drawn from "abcXYZ中文字あ" and a random budget 0..20, both helpers return a slice whose
    printed width equals the reported used count and is at most the budget.
*/
func TestGetRunesThatFitInColumnCountNeverSplitsWideRuneFuzz(test *testing.T) {
	pool := []rune("abcXYZ中文字あ")
	randomSource := rand.New(rand.NewSource(2))
	for iteration := 0; iteration < 3000; iteration++ {
		arrayOfRunes := make([]rune, randomSource.Intn(12))
		for currentIndex := range arrayOfRunes {
			arrayOfRunes[currentIndex] = pool[randomSource.Intn(len(pool))]
		}
		maxColumns := randomSource.Intn(21)

		prefix, prefixUsed := GetRunesThatFitInColumnCountFromStart(arrayOfRunes, maxColumns)
		suffix, suffixUsed := GetRunesThatFitInColumnCountFromEnd(arrayOfRunes, maxColumns)

		assert.Equal(test, GetWidthOfRunesWhenPrinted(prefix), prefixUsed, "GetRunesThatFitInColumnCountFromStart used count must equal the prefix printed width.")
		assert.Equal(test, GetWidthOfRunesWhenPrinted(suffix), suffixUsed, "GetRunesThatFitInColumnCountFromEnd used count must equal the suffix printed width.")
		assert.LessOrEqual(test, prefixUsed, maxColumns, "GetRunesThatFitInColumnCountFromStart must not exceed the column budget.")
		assert.LessOrEqual(test, suffixUsed, maxColumns, "GetRunesThatFitInColumnCountFromEnd must not exceed the column budget.")
		assert.True(test, string(prefix) == string(arrayOfRunes[:len(prefix)]), "GetRunesThatFitInColumnCountFromStart must return an unbroken prefix.")
		assert.True(test, string(suffix) == string(arrayOfRunes[len(arrayOfRunes)-len(suffix):]), "GetRunesThatFitInColumnCountFromEnd must return an unbroken suffix.")
		if len(arrayOfRunes) > 0 && maxColumns >= GetWidthOfRunesWhenPrinted(arrayOfRunes) {
			assert.Equal(test, string(arrayOfRunes), string(prefix), "A budget covering the whole array must take all of it.")
			assert.Equal(test, string(arrayOfRunes), string(suffix), "A budget covering the whole array must take all of it.")
		}
	}
}

/*
TestGetColumnIndexBasedOnRuneIndex is a test which verifies that GetColumnIndexBasedOnRuneIndex maps a rune index in "a中b" to the printed
column where that rune begins, counting two columns for the wide rune, and that the array length maps to the total
printed width.

Example:
    Expected Inputs / Outputs for []rune("a中b"):
        0 -> 0 | 1 -> 1 | 2 -> 3 | 3 -> 4
*/
func TestGetColumnIndexBasedOnRuneIndex(test *testing.T) {
	arrayOfRunes := []rune("a中b")
	for runeIndex, expectedColumn := range []int{0, 1, 3, 4} {
		assert.Equalf(test, expectedColumn, GetColumnIndexBasedOnRuneIndex(arrayOfRunes, runeIndex), "GetColumnIndexBasedOnRuneIndex wrong for rune index %d.", runeIndex)
	}
}

/*
TestGetRuneIndexBasedOnColumnIndex is a test which verifies that GetRuneIndexBasedOnColumnIndex maps a printed column in "a中b" to the rune
covering it, that a column on the trailing half of the wide rune resolves to that rune's index, and that a column
at or past the printed width resolves to the array length.

Example:
    Expected Inputs / Outputs for []rune("a中b") (printed width 4):
        0 -> 0 | 1 -> 1 | 2 -> 1 | 3 -> 2 | 4 -> 3 | 5 -> 3
    Also GetRuneIndexBasedOnColumnIndex([]rune("中"), 1) -> 0
*/
func TestGetRuneIndexBasedOnColumnIndex(test *testing.T) {
	arrayOfRunes := []rune("a中b")
	for column, expectedRuneIndex := range []int{0, 1, 1, 2, 3, 3} {
		assert.Equalf(test, expectedRuneIndex, GetRuneIndexBasedOnColumnIndex(arrayOfRunes, column), "GetRuneIndexBasedOnColumnIndex wrong for column %d.", column)
	}
	assert.Equal(test, 0, GetRuneIndexBasedOnColumnIndex([]rune("中"), 1), "A column on the trailing half of a wide rune must resolve to that rune's index.")
}

/*
TestRuneIndexColumnRoundTrip is a test which verifies that for random mixed ASCII and CJK arrays the identity
GetRuneIndexBasedOnColumnIndex(runes, GetColumnIndexBasedOnRuneIndex(runes, i)) == i holds for every rune index i from zero to the array
length.

Example:
    For []rune("a中b") and i in 0..3, GetColumnIndexBasedOnRuneIndex gives 0,1,3,4 and GetRuneIndexBasedOnColumnIndex maps each back to
    0,1,2,3.
*/
func TestRuneIndexColumnRoundTrip(test *testing.T) {
	pool := []rune("abcXYZ中文字あ")
	randomSource := rand.New(rand.NewSource(3))
	for iteration := 0; iteration < 2000; iteration++ {
		arrayOfRunes := make([]rune, randomSource.Intn(12))
		for currentIndex := range arrayOfRunes {
			arrayOfRunes[currentIndex] = pool[randomSource.Intn(len(pool))]
		}
		for runeIndex := 0; runeIndex <= len(arrayOfRunes); runeIndex++ {
			column := GetColumnIndexBasedOnRuneIndex(arrayOfRunes, runeIndex)
			assert.Equalf(test, runeIndex, GetRuneIndexBasedOnColumnIndex(arrayOfRunes, column), "Round trip failed at rune index %d for %q.", runeIndex, string(arrayOfRunes))
		}
	}
}

/*
TestFitFunctionsAsciiCharacterization is a test which pins the current pure-ASCII outputs of
GetMaxCharactersThatFitInStringSize and GetMaxCharactersThatFitInStringSizeReverse as literals, so Phase 2's
column accounting cannot silently regress ASCII behaviour.

Example:
    GetMaxCharactersThatFitInStringSize("hello world", 5)      -> "hello"
    GetMaxCharactersThatFitInStringSize("{{red}}test{{/}}", 2) -> "{{red}}te"
    GetMaxCharactersThatFitInStringSizeReverse("hello world", 5) -> 5
    GetMaxCharactersThatFitInStringSizeReverse("{{red}}test{{/}}", 2) -> 2
*/
func TestFitFunctionsAsciiCharacterization(test *testing.T) {
	forwardCases := []struct {
		input          string
		maxColumns     int
		expectedResult string
	}{
		{"hello world", 5, "hello"},
		{"hello", 10, "hello"},
		{"hello", 5, "hello"},
		{"hello", 3, "hel"},
		{"hello", 1, "h"},
		{"", 4, ""},
		{"{{red}}test{{/}}", 2, "{{red}}te"},
		{"{{red}}test{{/}}", 4, "{{red}}test"},
		{"ab{{x}}cd", 3, "ab{{x}}c"},
	}
	for _, currentCase := range forwardCases {
		obtained := string(GetMaxCharactersThatFitInStringSize([]rune(currentCase.input), currentCase.maxColumns))
		assert.Equalf(test, currentCase.expectedResult, obtained, "GetMaxCharactersThatFitInStringSize(%q, %d) regressed.", currentCase.input, currentCase.maxColumns)
	}

	reverseCases := []struct {
		input         string
		maxColumns    int
		expectedCount int
	}{
		{"hello world", 5, 5},
		{"hello", 10, 5},
		{"hello", 5, 5},
		{"hello", 3, 3},
		{"hello", 1, 1},
		{"", 4, 0},
		{"{{red}}test{{/}}", 2, 2},
		{"{{red}}test{{/}}", 4, 4},
		{"ab{{x}}cd", 3, 3},
	}
	for _, currentCase := range reverseCases {
		obtained := GetMaxCharactersThatFitInStringSizeReverse([]rune(currentCase.input), currentCase.maxColumns)
		assert.Equalf(test, currentCase.expectedCount, obtained, "GetMaxCharactersThatFitInStringSizeReverse(%q, %d) regressed.", currentCase.input, currentCase.maxColumns)
	}
}

/*
TestFormattedRuneArrayAsciiCharacterization is a test which pins the current pure-ASCII outputs of
GetFormattedString for left, right, centre and no-padding alignment as literals, locking Rule 2 no-regression
across Phase 2's move to column-based padding.

Example:
    GetFormattedString("Hi", 5, AlignmentCenter)  -> " Hi  "
    GetFormattedString("test", 5, AlignmentLeft)  -> "test "
    GetFormattedString("test", 10, AlignmentCenter) -> "   test   "
*/
func TestFormattedRuneArrayAsciiCharacterization(test *testing.T) {
	tableOfCases := []struct {
		input          string
		width          int
		alignment      int
		expectedResult string
	}{
		{"Hi", 5, constants.AlignmentLeft, "Hi   "},
		{"Hi", 5, constants.AlignmentRight, "   Hi"},
		{"Hi", 5, constants.AlignmentCenter, " Hi  "},
		{"Hi", 5, constants.AlignmentNoPadding, " Hi "},
		{"Hi", 10, constants.AlignmentCenter, "    Hi    "},
		{"test", 5, constants.AlignmentLeft, "test "},
		{"test", 5, constants.AlignmentCenter, "test "},
		{"test", 10, constants.AlignmentCenter, "   test   "},
		{"test", 10, constants.AlignmentRight, "      test"},
		{"hello", 5, constants.AlignmentCenter, "hello"},
		{"hello", 1, constants.AlignmentLeft, "h"},
	}
	for _, currentCase := range tableOfCases {
		obtained := GetFormattedString(currentCase.input, currentCase.width, currentCase.alignment)
		assert.Equalf(test, currentCase.expectedResult, obtained, "GetFormattedString(%q, %d, %d) regressed.", currentCase.input, currentCase.width, currentCase.alignment)
	}
}

/*
TestGetFormattedRuneArrayCjk is a test which verifies that GetFormattedString distributes its pad in columns for
CJK content: left and right alignment place the whole pad on one side, centre splits it evenly with any odd column
on the right, and content wider than the target is truncated on a rune boundary rather than through a wide rune.

Example:
    GetFormattedString("中文", 10, AlignmentLeft)   -> "中文" + 6 spaces        (printed width 10)
    GetFormattedString("中文", 10, AlignmentRight)  -> 6 spaces + "中文"        (printed width 10)
    GetFormattedString("中文", 10, AlignmentCenter) -> 3 spaces + "中文" + 3 spaces (printed width 10)
    GetFormattedString("中文", 11, AlignmentCenter) -> 3 spaces + "中文" + 4 spaces (odd column on the right)
    GetFormattedString("中文字", 4, AlignmentLeft)  -> "中文"                    (rune-boundary truncation)
*/
func TestGetFormattedRuneArrayCjk(test *testing.T) {
	assert.Equal(test, "中文      ", GetFormattedString("中文", 10, constants.AlignmentLeft), "Left CJK pad wrong.")
	assert.Equal(test, "      中文", GetFormattedString("中文", 10, constants.AlignmentRight), "Right CJK pad wrong.")
	assert.Equal(test, "   中文   ", GetFormattedString("中文", 10, constants.AlignmentCenter), "Centre CJK pad wrong.")
	assert.Equal(test, "   中文    ", GetFormattedString("中文", 11, constants.AlignmentCenter), "Centre CJK odd column must fall on the right.")
	assert.Equal(test, "中文", GetFormattedString("中文字", 4, constants.AlignmentLeft), "CJK truncation must land on a rune boundary.")

	for _, alignment := range []int{constants.AlignmentLeft, constants.AlignmentRight, constants.AlignmentCenter} {
		obtained := GetFormattedString("中文", 10, alignment)
		assert.Equalf(test, 10, GetWidthOfRunesWhenPrinted([]rune(obtained)), "CJK formatted printed width wrong for alignment %d.", alignment)
	}
}

/*
TestFitReverseIsInverseOfFit is a test which verifies across random mixed ASCII and CJK strings and random column
budgets that GetMaxCharactersThatFitInStringSizeReverse counts exactly the runes of the suffix GetRunesThatFitInColumnCountFromEnd
would keep, that the counted suffix never exceeds the budget in printed width, and that it never begins on the
trailing half of a wide rune.

Example:
    For "a中b" and budget 3, the fitting suffix is "中b" (3 columns) so the reverse count is 2.
*/
func TestFitReverseIsInverseOfFit(test *testing.T) {
	pool := []rune("abcXYZ 中文字あ")
	randomSource := rand.New(rand.NewSource(4))
	for iteration := 0; iteration < 3000; iteration++ {
		arrayOfRunes := make([]rune, randomSource.Intn(16))
		for currentIndex := range arrayOfRunes {
			arrayOfRunes[currentIndex] = pool[randomSource.Intn(len(pool))]
		}
		maxColumns := randomSource.Intn(21)

		reverseCount := GetMaxCharactersThatFitInStringSizeReverse(arrayOfRunes, maxColumns)
		suffix, suffixUsed := GetRunesThatFitInColumnCountFromEnd(arrayOfRunes, maxColumns)

		assert.Equal(test, len(suffix), reverseCount, "The reverse count must equal the fitting suffix rune length.")
		assert.LessOrEqual(test, suffixUsed, maxColumns, "The fitting suffix must not exceed the budget.")
		assert.Equal(test, GetWidthOfRunesWhenPrinted(suffix), suffixUsed, "The fitting suffix printed width must match its used count.")
		if len(suffix) > 0 {
			assert.Equal(test, string(arrayOfRunes[len(arrayOfRunes)-len(suffix):]), string(suffix), "The suffix must be an unbroken tail of the array.")
		}
	}
}
