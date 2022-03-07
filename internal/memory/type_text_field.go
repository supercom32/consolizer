package memory

import "encoding/json"

type TextFieldEntryType struct {
	StyleEntry TuiStyleEntryType
	XLocation int
	YLocation int
	Width int
	MaxLengthAllowed int
	DefaultValue string
	CursorPosition int
	ViewportPosition int
	IsPasswordProtected bool
	CurrentValue []rune
}

func (shared TextFieldEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		StyleEntry TuiStyleEntryType
		XLocation int
		YLocation int
		Width int
		MaxLengthAllowed int
		IsPasswordProtected bool
		CurrentValue []rune
		DefaultValue string
		StringPosition int
		CursorPosition int
		ViewportPosition int
	}{
		StyleEntry: shared.StyleEntry,
		XLocation: shared.XLocation,
		YLocation: shared.YLocation,
		Width: shared.Width,
		MaxLengthAllowed: shared.MaxLengthAllowed,
		IsPasswordProtected: shared.IsPasswordProtected,
		CurrentValue: shared.CurrentValue,
		DefaultValue: shared.DefaultValue,
		CursorPosition: shared.CursorPosition,
		ViewportPosition: shared.ViewportPosition,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (shared TextFieldEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

func NewTextFieldEntry(existingTextFieldEntry ...*TextFieldEntryType) TextFieldEntryType {
	var textFieldEntry TextFieldEntryType
	if existingTextFieldEntry != nil {
		textFieldEntry.XLocation = existingTextFieldEntry[0].XLocation
		textFieldEntry.YLocation = existingTextFieldEntry[0].YLocation
		textFieldEntry.Width = existingTextFieldEntry[0].Width
		textFieldEntry.MaxLengthAllowed = existingTextFieldEntry[0].MaxLengthAllowed
		textFieldEntry.IsPasswordProtected = existingTextFieldEntry[0].IsPasswordProtected
		textFieldEntry.CurrentValue = existingTextFieldEntry[0].CurrentValue
		textFieldEntry.DefaultValue = existingTextFieldEntry[0].DefaultValue
		textFieldEntry.CursorPosition = existingTextFieldEntry[0].CursorPosition
	}
	textFieldEntry.CurrentValue = []rune{' '}
	return textFieldEntry
}
