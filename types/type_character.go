package types

import (
	"encoding/json"
)

/*
CharacterEntryType is a structure which represents a single character cell in the terminal. In addition, the following should be noted:

- It contains the character rune, its visual attributes, and its layer and parent relationships.

Example:

	var character types.CharacterEntryType
*/
type CharacterEntryType struct {
	Character      rune
	AttributeEntry AttributeEntryType
	LayerAlias     string
	ParentAlias    string
}

/*
MarshalJSON is a method which serializes a character entry to JSON. In addition, the following should be noted:

- It implements the json.Marshaler interface for CharacterEntryType.

- It converts the character entry to a JSON string representation.

- It is used for serializing character entries when saving state or transmitting data.

Example:

	instance.MarshalJSON()
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
GetEntryAsJsonDump is a method which retrieves a JSON string representation of the character entry. In addition, the following should be noted:

- It returns a formatted JSON string of the character entry.

- It is used for debugging and logging purposes.

- The output is human-readable and includes all character properties.

Example:

	instance.GetEntryAsJsonDump()
*/
func (shared CharacterEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

/*
NewCharacterEntry is a constructor which creates a new character entry. In addition, the following should be noted:

- It initializes a character entry with default values.

- It sets up a character with specified properties like foreground color, background color, and character value.

- The character can be used to render text in the terminal interface.

Example:

	NewCharacterEntry(existingCharacterEntry)
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
