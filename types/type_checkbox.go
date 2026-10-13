package types

import (
	"encoding/json"
	"strconv"
)

/*
CheckboxEntryType is a structure which represents a checkbox control.
*/
type CheckboxEntryType struct {
	BaseControlType
	IsSelected   bool
	IsLabelFirst bool
}

/*
MarshalJSON is a method which serializes a checkbox control to JSON. The result is a single flat JSON object holding
every base control property alongside the checkbox's selection state and whether its label is drawn first. In
addition, the following should be noted:

  - The base control properties are produced by BaseControlType's own MarshalJSON and then merged with the checkbox
    fields. Embedding BaseControlType in an anonymous struct instead would promote that MarshalJSON to the anonymous
    struct, which would silently drop every checkbox field from the output.

  - Because the fields are merged through a map, the keys in the output are in sorted order.

Example:

	jsonData, err := checkboxEntry.MarshalJSON()
*/
func (shared CheckboxEntryType) MarshalJSON() ([]byte, error) {
	baseJson, err := json.Marshal(shared.BaseControlType)
	if err != nil {
		return nil, err
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(baseJson, &fields); err != nil {
		return nil, err
	}
	fields["IsSelected"] = json.RawMessage(strconv.FormatBool(shared.IsSelected))
	fields["IsLabelFirst"] = json.RawMessage(strconv.FormatBool(shared.IsLabelFirst))
	return json.Marshal(fields)
}

/*
GetEntryAsJsonDump is a method which retrieves a JSON string representation of a checkbox control. In addition, the
following should be noted:

- It returns a formatted JSON string of the checkbox's state.

Example:
    instance.GetEntryAsJsonDump()
*/
func (shared CheckboxEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

/*
NewCheckboxEntry is a constructor which creates a new checkbox control. In addition, the following should be noted:

- It initializes a checkbox with default values.

- It can optionally copy properties from an existing checkbox.

- It uses Unicode characters U+2610 and U+2611 for rendering states.

Example:
    NewCheckboxEntry(existingCheckboxEntry)
*/
func NewCheckboxEntry(existingCheckboxEntry ...*CheckboxEntryType) CheckboxEntryType {
	var checkboxEntry CheckboxEntryType
	checkboxEntry.BaseControlType = NewBaseControl()
	checkboxEntry.IsSelected = false
	checkboxEntry.IsLabelFirst = false

	if existingCheckboxEntry != nil {
		checkboxEntry.BaseControlType = existingCheckboxEntry[0].BaseControlType
		checkboxEntry.IsSelected = existingCheckboxEntry[0].IsSelected
		checkboxEntry.IsLabelFirst = existingCheckboxEntry[0].IsLabelFirst
	}
	return checkboxEntry
}

/*
IsCheckboxEqual is a method which compares two checkbox controls for equality. In addition, the following should be
noted:

- It compares the base control properties, the selection state, and whether the label is drawn first.

Example:
    IsCheckboxEqual(sourceCheckboxEntry, targetCheckboxEntry)
*/
func IsCheckboxEqual(sourceCheckboxEntry *CheckboxEntryType, targetCheckboxEntry *CheckboxEntryType) bool {
	return sourceCheckboxEntry.BaseControlType.IsEqual(&targetCheckboxEntry.BaseControlType) &&
		sourceCheckboxEntry.IsSelected == targetCheckboxEntry.IsSelected &&
		sourceCheckboxEntry.IsLabelFirst == targetCheckboxEntry.IsLabelFirst
}
