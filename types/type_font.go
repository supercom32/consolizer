
package types

import (
	"encoding/json"
)

type FontEntryType struct {
	Name     string
	FontType byte
	Spacing  byte
	Height   int
	CharList []uint16
	FontData []byte
	Glyphs   []*Glyph
}

type Cell struct {
	Char  rune
	Color byte
}

type Glyph struct {
	Width, Height int
	Cells         []Cell
}

func (shared FontEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		Name     string
		FontType byte
		Spacing  byte
		Height   int
		CharList []uint16
		FontData []byte
		Glyphs   []*Glyph
	}{
		Name:     shared.Name,
		FontType: shared.FontType,
		Spacing:  shared.Spacing,
		Height:   shared.Height,
		CharList: shared.CharList,
		FontData: shared.FontData,
		Glyphs:   shared.Glyphs,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (shared FontEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

func NewFontEntry(existingFontEntry ...*FontEntryType) FontEntryType {
	var fontEntry FontEntryType
	if existingFontEntry != nil {
		fontEntry.Name = existingFontEntry[0].Name
		fontEntry.FontType = existingFontEntry[0].FontType
		fontEntry.Spacing = existingFontEntry[0].Spacing
		fontEntry.Height = existingFontEntry[0].Height
		fontEntry.CharList = existingFontEntry[0].CharList
		fontEntry.FontData = existingFontEntry[0].FontData
		fontEntry.Glyphs = existingFontEntry[0].Glyphs
	}
	return fontEntry
}
