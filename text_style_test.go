package consolizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
TestCalculateStringLengthWithoutMarkupWithMultiByteRunes is a test which verifies that
CalculateStringLengthWithoutMarkup measures text in printed columns by rune, so multi-byte characters before, inside
or after tags are no longer counted once per byte, while ASCII text is measured as before.

Example:

	Expected Inputs:
	    "{{red}}Hello{{white}} World", "Efraín  {{warn}}Defector{{/}}  ○○○", "{{wärn}}x{{/}}", "·{{red}}○○○",
	    "í {{unclosed ○", "世界{{red}}中", "{{}}" and "".
	Expected Outputs:
	    11, 21, 1, 4, 14, 6, 0 and 0 respectively.
*/
func TestCalculateStringLengthWithoutMarkupWithMultiByteRunes(test *testing.T) {
	testCases := []struct {
		text           string
		expectedLength int
	}{
		{"{{red}}Hello{{white}} World", 11},
		{"Efraín  {{warn}}Defector{{/}}  ○○○", 21},
		{"{{wärn}}x{{/}}", 1},
		{"·{{red}}○○○", 4},
		{"í {{unclosed ○", 14},
		{"世界{{red}}中", 6},
		{"{{}}", 0},
		{"", 0},
	}

	for _, testCase := range testCases {
		assert.Equal(test, testCase.expectedLength, CalculateStringLengthWithoutMarkup(testCase.text), "wrong length for "+testCase.text)
	}
}
