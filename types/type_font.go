package types

import (
	"encoding/json"
)

/*
FontFamilyType is a structure which represents a collection of fonts.

Example:

	var fontFamily types.FontFamilyType
*/
type FontFamilyType struct {
	Fonts []*FontEntryType
}

/*
FontEntryType is a structure which represents a single font and its properties.

Example:

	var fontEntry types.FontEntryType
*/
type FontEntryType struct {
	Name                  string
	FontType              byte
	Spacing               byte
	Height                int
	CharList              []uint16
	FontData              []byte
	Glyphs                []*Glyph
	CharacterSpacing      int
	BlankSizeInCharacters int
}

/*
Cell is a structure which represents a single character cell with a character and color.

Example:

	var cell types.Cell
*/
type Cell struct {
	Char  rune
	Color byte
}

/*
Glyph is a structure which represents a single character glyph with its width, height, and cells.

Example:

	var glyph types.Glyph
*/
type Glyph struct {
	Width, Height int
	Cells         []Cell
}

/*
MarshalJSON is a method which marshals the FontEntryType into a JSON byte array.

Example:

	instance.MarshalJSON()
*/
func (shared FontEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		Name                   string
		FontType               byte
		FontCharacterSpacing   byte
		Height                 int
		CharList               []uint16
		FontData               []byte
		Glyphs                 []*Glyph
		RenderCharacterSpacing int
		BlankSizeInCharacters  int
	}{
		Name:                   shared.Name,
		FontType:               shared.FontType,
		FontCharacterSpacing:   shared.Spacing,
		Height:                 shared.Height,
		CharList:               shared.CharList,
		FontData:               shared.FontData,
		Glyphs:                 shared.Glyphs,
		RenderCharacterSpacing: shared.CharacterSpacing,
		BlankSizeInCharacters:  shared.BlankSizeInCharacters,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

/*
GetEntryAsJsonDump is a method which returns a JSON string representation of the FontEntryType.

Example:

	instance.GetEntryAsJsonDump()
*/
func (shared FontEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}
