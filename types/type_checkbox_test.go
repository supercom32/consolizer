package types

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
TestCheckboxTypeLabelFirst is a test which verifies that the IsLabelFirst field defaults to false, is copied by
NewCheckboxEntry, takes part in IsCheckboxEqual, and is included in the JSON produced by MarshalJSON.

Example:

	Expected Inputs:
		A new checkbox entry with Alias "MyCheckbox", IsSelected true, and IsLabelFirst true.

	Expected Outputs:
		A fresh entry has IsLabelFirst false. A clone made with NewCheckboxEntry equals the original. Clearing
		IsLabelFirst on the clone makes IsCheckboxEqual report false. The JSON of the original decodes to an object
		containing "IsLabelFirst": true and "IsSelected": true.
*/
func TestCheckboxTypeLabelFirst(test *testing.T) {
	assert.False(test, NewCheckboxEntry().IsLabelFirst, "IsLabelFirst should default to false")

	firstCheckboxEntry := NewCheckboxEntry()
	firstCheckboxEntry.Alias = "MyCheckbox"
	firstCheckboxEntry.IsSelected = true
	firstCheckboxEntry.IsLabelFirst = true

	secondCheckboxEntry := NewCheckboxEntry(&firstCheckboxEntry)
	assert.Equal(test, firstCheckboxEntry, secondCheckboxEntry, "The clone should be identical to the original")
	assert.True(test, IsCheckboxEqual(&firstCheckboxEntry, &secondCheckboxEntry))

	secondCheckboxEntry.IsLabelFirst = false
	assert.False(test, IsCheckboxEqual(&firstCheckboxEntry, &secondCheckboxEntry), "IsLabelFirst should affect equality")

	jsonBytes, err := firstCheckboxEntry.MarshalJSON()
	assert.NoError(test, err)
	var decoded map[string]interface{}
	assert.NoError(test, json.Unmarshal(jsonBytes, &decoded))
	assert.Equal(test, true, decoded["IsLabelFirst"])
	assert.Equal(test, true, decoded["IsSelected"])
}

/*
TestCheckboxTypeMarshalJSONKeepsBaseFields is a test which verifies that MarshalJSON produces one flat JSON object that
still holds the base control properties next to the checkbox fields, rather than only one or the other.

Example:

	Expected Inputs:
		A new checkbox entry with Alias "MyCheckbox", Label "Enable", XLocation 3, and IsSelected false.

	Expected Outputs:
		The decoded object has "Alias": "MyCheckbox", "Label": "Enable", "XLocation": 3, "IsSelected": false, and
		"IsLabelFirst": false, and has no "Mutex" key.
*/
func TestCheckboxTypeMarshalJSONKeepsBaseFields(test *testing.T) {
	checkboxEntry := NewCheckboxEntry()
	checkboxEntry.Alias = "MyCheckbox"
	checkboxEntry.Label = "Enable"
	checkboxEntry.XLocation = 3

	jsonBytes, err := json.Marshal(checkboxEntry)
	assert.NoError(test, err)
	var decoded map[string]interface{}
	assert.NoError(test, json.Unmarshal(jsonBytes, &decoded))
	assert.Equal(test, "MyCheckbox", decoded["Alias"])
	assert.Equal(test, "Enable", decoded["Label"])
	assert.Equal(test, float64(3), decoded["XLocation"])
	assert.Equal(test, false, decoded["IsSelected"])
	assert.Equal(test, false, decoded["IsLabelFirst"])
	assert.NotContains(test, decoded, "Mutex")
}
