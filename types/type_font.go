package types

import (
	"encoding/json"
)

type FontFamilyType struct {
	Fonts []*FontEntryType
}

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

type Cell struct {
	Char  rune
	Color byte
}

type Glyph struct {
	Width, Height int
	Cells         []Cell
}

/*
MarshalJSON is a method which allows you to marshaljson.

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
GetEntryAsJsonDump is a method which allows you to getentryasjsondump.

:return: string

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
