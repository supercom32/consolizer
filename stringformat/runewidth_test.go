package stringformat

import (
	"os"
	"testing"

	runewidth "github.com/mattn/go-runewidth"
	"github.com/stretchr/testify/assert"

	// Blank-imported so tcell's real init() runs in this test binary. tcell
	// mutates the same global runewidth.DefaultCondition that stringformat's
	// own init() touches, so this proves the two converge in a real build.
	_ "github.com/gdamore/tcell/v2"
)

/*
TestGetWidthOfRuneWhenPrintedClassification is a test which verifies that GetWidthOfRuneWhenPrinted reports the correct terminal cell width for a
representative set of runes covering ASCII, CJK ideographs, kana, hangul, fullwidth forms, the narrow multibyte
scripts the old oracle wrongly assumed wide, box-drawing, block elements, and the U+2600 symbol block.

Example:
    Expected Inputs / Outputs:
        'A' -> 1 | ' ' -> 1 | '中' -> 2 | '漢' -> 2 | 'あ' -> 2 | 'カ' -> 2 | '가' -> 2 |
        'Ａ' (fullwidth) -> 2 | '　' (ideographic space) -> 2 |
        'é' -> 1 | 'Ж' -> 1 | 'ก' -> 1 | '─' -> 1 | '█' -> 1 | '❤' -> 1
*/
func TestGetWidthOfRuneWhenPrintedClassification(test *testing.T) {
	tableOfCases := []struct {
		character     rune
		expectedWidth int
		description   string
	}{
		{'A', 1, "printable ASCII letter"},
		{' ', 1, "ASCII space"},
		{'中', 2, "CJK ideograph"},
		{'漢', 2, "CJK ideograph"},
		{'あ', 2, "hiragana"},
		{'カ', 2, "katakana"},
		{'가', 2, "hangul syllable"},
		{'Ａ', 2, "fullwidth Latin A"},
		{'　', 2, "ideographic space"},
		{'é', 1, "precomposed accented Latin (regression vs old assume-wide)"},
		{'Ж', 1, "Cyrillic (regression)"},
		{'ก', 1, "Thai (regression)"},
		{'─', 1, "box drawing"},
		{'█', 1, "block element"},
		{'❤', 1, "U+2764 heart; must be 1 to match tcell"},
	}
	for _, currentCase := range tableOfCases {
		obtainedWidth := GetWidthOfRuneWhenPrinted(currentCase.character)
		assert.Equalf(test, currentCase.expectedWidth, obtainedWidth,
			"GetWidthOfRuneWhenPrinted(%q) [%s] returned %d, expected %d.", currentCase.character, currentCase.description, obtainedWidth, currentCase.expectedWidth)
	}
}

/*
TestGetWidthOfRuneWhenPrintedAsciiFastPath is a test which verifies that every printable ASCII rune in the range 0x20 through 0x7E
resolves to width 1 via the fast path, while control runes in 0x00 through 0x1F and the 0x7F delete rune are left
to the runewidth backend rather than forced to a fixed value.

Example:
    GetWidthOfRuneWhenPrinted(0x41) == 1
    GetWidthOfRuneWhenPrinted(0x07) == runewidth.RuneWidth(0x07)
*/
func TestGetWidthOfRuneWhenPrintedAsciiFastPath(test *testing.T) {
	for currentRune := rune(0x20); currentRune <= 0x7E; currentRune++ {
		assert.Equalf(test, 1, GetWidthOfRuneWhenPrinted(currentRune), "Printable ASCII rune 0x%02X should have width 1.", currentRune)
	}
	for currentRune := rune(0x00); currentRune <= 0x1F; currentRune++ {
		assert.Equalf(test, runewidth.RuneWidth(currentRune), GetWidthOfRuneWhenPrinted(currentRune),
			"Control rune 0x%02X should defer to the runewidth backend.", currentRune)
	}
	assert.Equalf(test, runewidth.RuneWidth(0x7F), GetWidthOfRuneWhenPrinted(0x7F), "The 0x7F delete rune should defer to the runewidth backend.")
}

/*
TestOracleMatchesTcellBackend is a test which verifies that GetWidthOfRuneWhenPrinted never disagrees with the package-level
runewidth.RuneWidth that tcell calls in cell.go. It checks the classification table plus a 500-rune stride sample
across the Basic Multilingual Plane. Because this test binary blank-imports tcell, tcell's real init() has run and
both packages share one global condition.

Example:
    For every sampled rune r: GetWidthOfRuneWhenPrinted(r) == runewidth.RuneWidth(r)
*/
func TestOracleMatchesTcellBackend(test *testing.T) {
	assertOracleMatchesBackend(test)
}

/*
TestOracleMatchesTcellBackendEastAsian is a test which verifies the same backend agreement as
TestOracleMatchesTcellBackend but under RUNEWIDTH_EASTASIAN=1. The runewidth backend reads that environment
variable once at init time, so this test skips unless the variable is already set in the process environment; run
it with a dedicated invocation such as RUNEWIDTH_EASTASIAN=1 go test -run TestOracleMatchesTcellBackendEastAsian.

Example:
    RUNEWIDTH_EASTASIAN=1 go test -run TestOracleMatchesTcellBackendEastAsian ./stringformat/...
*/
func TestOracleMatchesTcellBackendEastAsian(test *testing.T) {
	if os.Getenv("RUNEWIDTH_EASTASIAN") != "1" {
		test.Skip("RUNEWIDTH_EASTASIAN is not set to 1; run this test in a dedicated invocation with that environment variable.")
	}
	assertOracleMatchesBackend(test)
}

/*
assertOracleMatchesBackend is a method which asserts that stringformat.GetWidthOfRuneWhenPrinted equals the package-level
runewidth.RuneWidth for the classification table and a 500-rune stride sample across the Basic Multilingual Plane.

Example:
    assertOracleMatchesBackend(test)
*/
func assertOracleMatchesBackend(test *testing.T) {
	sampledRunes := []rune{'A', ' ', '中', '漢', 'あ', 'カ', '가', 'Ａ', '　', 'é', 'Ж', 'ก', '─', '█', '❤'}
	const strideCount = 500
	stride := rune(0xFFFF / strideCount)
	for currentRune := rune(0); currentRune <= 0xFFFF; currentRune += stride {
		sampledRunes = append(sampledRunes, currentRune)
	}
	for _, currentRune := range sampledRunes {
		assert.Equalf(test, runewidth.RuneWidth(currentRune), GetWidthOfRuneWhenPrinted(currentRune),
			"GetWidthOfRuneWhenPrinted(0x%04X) disagrees with the runewidth backend tcell uses.", currentRune)
	}
}

/*
TestPrintedWidthMixed is a test which verifies that GetWidthOfRunesWhenPrinted, which is built on GetWidthOfRuneWhenPrinted, sums
mixed narrow and wide runes in columns rather than in rune count.

Example:
    GetWidthOfRunesWhenPrinted([]rune("a中b文")) == 6
    GetWidthOfRunesWhenPrinted([]rune("café")) == 4
*/
func TestPrintedWidthMixed(test *testing.T) {
	assert.Equal(test, 6, GetWidthOfRunesWhenPrinted([]rune("a中b文")), "Mixed ASCII and CJK width should be counted in columns.")
	assert.Equal(test, 4, GetWidthOfRunesWhenPrinted([]rune("café")), "Accented Latin must not be double counted.")
	assert.Equal(test, 10, GetWidthOfRunesWhenPrinted([]rune("中文字五个")), "Five CJK ideographs occupy ten columns.")
}
