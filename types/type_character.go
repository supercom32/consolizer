package types

import (
	"encoding/json"
)

type CharacterEntryType struct {
	Character      rune
	AttributeEntry AttributeEntryType
	LayerAlias     string
	ParentAlias    string
}

/*
MarshalJSON allows you to convert a character entry to JSON format. In addition, the following
information should be noted:

- Implements the json.Marshaler interface for CharacterEntryType.
- Converts the character entry to a JSON string representation.
- Used for serializing character entries when saving state or transmitting data.
*/
func (shared CharacterEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		Character      rune
		AttributeEntry AttributeEntryType
		LayerAlias     string
		ParentAlias    string
	}{
		Character:      shared.Character,
		AttributeEntry: shared.AttributeEntry,
		LayerAlias:     shared.LayerAlias,
		ParentAlias:    shared.ParentAlias,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

/*
GetEntryAsJsonDump allows you to get a JSON string representation of the character entry. In addition, the following
information should be noted:

- Returns a formatted JSON string of the character entry.
- Used for debugging and logging purposes.
- The output is human-readable and includes all character properties.
*/
func (shared CharacterEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

/*
NewCharacterEntry allows you to create a new character entry. In addition, the following
information should be noted:

- Initializes a character entry with default values.
- Sets up a character with specified properties like foreground color, background color, and character value.
- The character can be used to render text in the terminal interface.
*/
func NewCharacterEntry(existingCharacterEntry ...*CharacterEntryType) CharacterEntryType {
	var characterEntry CharacterEntryType
	if existingCharacterEntry != nil {
		characterEntry.Character = existingCharacterEntry[0].Character
		characterEntry.AttributeEntry = NewAttributeEntry(&existingCharacterEntry[0].AttributeEntry)
		characterEntry.LayerAlias = existingCharacterEntry[0].LayerAlias
		characterEntry.ParentAlias = existingCharacterEntry[0].ParentAlias
	} else {
		characterEntry.AttributeEntry = NewAttributeEntry()
	}
	return characterEntry
}
