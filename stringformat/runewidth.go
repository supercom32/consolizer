package stringformat

import (
	"os"

	runewidth "github.com/mattn/go-runewidth"
)

/*
init is a method which allows you to align this package's width backend with the one tcell uses at render time. It
mirrors the mutation performed by github.com/gdamore/tcell/v2 cell.go so that GetWidthOfRuneWhenPrinted stays
correct even in builds or unit tests that never import tcell. In addition, the following should be noted:

- The mutation is idempotent, so the order in which this init and tcell's own init run does not matter.

- The RUNEWIDTH_EASTASIAN environment variable is honoured exactly as tcell honours it: when set, the ambiguous
  East Asian width behaviour of the backend is left untouched.
*/
func init() {
	if os.Getenv("RUNEWIDTH_EASTASIAN") == "" {
		runewidth.DefaultCondition.EastAsianWidth = false
	}
}

/*
GetWidthOfRuneWhenPrinted is a method which allows you to obtain the number of terminal cells a single rune
occupies when printed. The result is 0 for a zero-width rune, 1 for a narrow rune, and 2 for a wide (double-cell)
rune such as a CJK ideograph, kana, hangul, or a fullwidth form. GetWidthOfRuneWhenPrinted is DEFINED as
"whatever tcell's runewidth backend returns" and must never acquire independent width logic of its own. It
delegates straight to the package-level runewidth.RuneWidth, the exact function tcell calls in cell.go, so it can
never disagree with what tcell draws. In addition, the following should be noted:

- Printable ASCII (0x20 through 0x7E) takes a fast path that is provably identical to the backend result.

- Control runes below 0x20 and 0x7F are deliberately left to the backend rather than forced to a fixed width.

:param character: The rune whose printed cell width is required.
:return: The number of terminal cells the rune occupies, one of 0, 1, or 2.

Example:
    width := GetWidthOfRuneWhenPrinted('中')
*/
func GetWidthOfRuneWhenPrinted(character rune) int {
	if character >= 0x20 && character < 0x7f {
		return 1
	}
	return runewidth.RuneWidth(character)
}
