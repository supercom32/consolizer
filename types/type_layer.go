package types

import (
	"encoding/json"
	"github.com/gdamore/tcell/v2"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/internal/stringformat"
	"github.com/supercom32/filesystem"
)

type LayerEntryType struct {
	Width            int
	Height           int
	ScreenXLocation  int
	ScreenYLocation  int
	CursorXLocation  int
	CursorYLocation  int
	ZOrder           int
	IsTopmost        bool
	IsFocusable      bool
	IsVisible        bool
	LayerAlias       string
	ParentAlias      string
	IsParent         bool
	DefaultAttribute AttributeEntryType
	CharacterMemory  [][]CharacterEntryType
}

func (shared LayerEntryType) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		Width            int
		Height           int
		ScreenXLocation  int
		ScreenYLocation  int
		CursorXLocation  int
		CursorYLocation  int
		ZOrder           int
		IsTopmost        bool
		IsFocusable      bool
		IsVisible        bool
		LayerAlias       string
		ParentAlias      string
		IsParent         bool
		DefaultAttribute AttributeEntryType
		CharacterMemory  [][]CharacterEntryType
	}{
		Width:            shared.Width,
		Height:           shared.Height,
		ScreenXLocation:  shared.ScreenXLocation,
		ScreenYLocation:  shared.ScreenYLocation,
		CursorXLocation:  shared.CursorXLocation,
		CursorYLocation:  shared.CursorYLocation,
		ZOrder:           shared.ZOrder,
		IsTopmost:        shared.IsTopmost,
		IsFocusable:      shared.IsFocusable,
		IsVisible:        shared.IsVisible,
		LayerAlias:       shared.LayerAlias,
		ParentAlias:      shared.ParentAlias,
		IsParent:         shared.IsParent,
		DefaultAttribute: shared.DefaultAttribute,
		CharacterMemory:  shared.CharacterMemory,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (shared LayerEntryType) GetBasicAnsiString() string {
	var ansiString string
	var currentForegroundColor constants.ColorType
	var currentBackgroundColor constants.ColorType
	for currentRow := 0; currentRow < shared.Height; currentRow++ {
		for currentCharacter := 0; currentCharacter < shared.Width; currentCharacter++ {
			if shared.CharacterMemory[currentRow][currentCharacter].AttributeEntry.ForegroundColor != currentForegroundColor {
				ansiString += shared.GetAnsiForegroundColorString(shared.CharacterMemory[currentRow][currentCharacter].AttributeEntry.ForegroundColor)
				currentForegroundColor = shared.CharacterMemory[currentRow][currentCharacter].AttributeEntry.ForegroundColor
			}
			if shared.CharacterMemory[currentRow][currentCharacter].AttributeEntry.BackgroundColor != currentBackgroundColor {
				ansiString += shared.GetAnsiBackgroundColorString(shared.CharacterMemory[currentRow][currentCharacter].AttributeEntry.BackgroundColor)
				currentBackgroundColor = shared.CharacterMemory[currentRow][currentCharacter].AttributeEntry.BackgroundColor
			}
			if shared.CharacterMemory[currentRow][currentCharacter].Character == constants.NullRune {
				ansiString += " "
			} else {
				ansiString += string(shared.CharacterMemory[currentRow][currentCharacter].Character)
			}
		}
		ansiString += shared.GetAnsiForegroundColorString(constants.AnsiColorByIndex[constants.ColorBlack])
		ansiString += shared.GetAnsiBackgroundColorString(constants.AnsiColorByIndex[constants.ColorBlack])
		currentForegroundColor = constants.AnsiColorByIndex[constants.ColorBlack]
		currentBackgroundColor = constants.AnsiColorByIndex[constants.ColorBlack]
		ansiString += "\n"
	}
	return ansiString
}

func (shared LayerEntryType) GetBasicAnsiStringAsBase64() string {
	ansiString := shared.GetBasicAnsiString()
	err := filesystem.WriteBytesToFile("/tmp/output.ans", []byte(ansiString), 0)
	if err != nil {
		panic(err)
	}
	return stringformat.GetStringAsBase64(ansiString)
}

func (shared LayerEntryType) GetBasicAnsiStringAsBase642() string {
	ansiString := shared.GetBasicAnsiString()
	err := filesystem.WriteBytesToFile("/tmp/output.ans", []byte(ansiString), 0)
	if err != nil {
		panic(err)
	}
	return stringformat.GetStringAsBase64(ansiString)
}

func (shared LayerEntryType) GetAnsiForegroundColorString(color constants.ColorType) string {
	var ansiString string
	redIndex, greenIndex, blueIndex := shared.GetRGBColorComponents(color)
	ansiString = "\u001b[38;2;" + stringformat.GetIntAsString(redIndex) + ";" + stringformat.GetIntAsString(greenIndex) + ";" + stringformat.GetIntAsString(blueIndex) + "m"
	return ansiString
}

func (shared LayerEntryType) GetAnsiBackgroundColorString(color constants.ColorType) string {
	var ansiString string
	redIndex, greenIndex, blueIndex := shared.GetRGBColorComponents(color)
	ansiString = "\u001b[48;2;" + stringformat.GetIntAsString(redIndex) + ";" + stringformat.GetIntAsString(greenIndex) + ";" + stringformat.GetIntAsString(blueIndex) + "m"
	return ansiString
}

func (shared LayerEntryType) GetAnsiLocateString(xLocation int, yLocation int) string {
	var ansiString string
	ansiString += "\033[99999A"
	ansiString += "\033[99999D"
	if yLocation != 0 {
		ansiString += "\033[" + stringformat.GetIntAsString(yLocation) + "B"
	}
	if xLocation != 0 {
		ansiString += "\033[" + stringformat.GetIntAsString(xLocation) + "C"
	}
	return ansiString
}

func (shared LayerEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

func (shared LayerEntryType) GetRGBColorComponents(color constants.ColorType) (int32, int32, int32) {
	var redColorIndex int32
	var greenColorIndex int32
	var blueColorIndex int32
	redColorIndex, greenColorIndex, blueColorIndex = tcell.Color.RGB(tcell.Color(color))
	return redColorIndex, greenColorIndex, blueColorIndex
}

func NewLayerEntry(layerAlias string, parentAlias string, width int, height int, existingLayerEntry ...*LayerEntryType) LayerEntryType {
	var layerEntry LayerEntryType
	if existingLayerEntry != nil {
		layerEntry.Width = existingLayerEntry[0].Width
		layerEntry.Height = existingLayerEntry[0].Height
		layerEntry.LayerAlias = existingLayerEntry[0].LayerAlias
		layerEntry.ScreenXLocation = existingLayerEntry[0].ScreenXLocation
		layerEntry.ScreenYLocation = existingLayerEntry[0].ScreenYLocation
		layerEntry.CursorXLocation = existingLayerEntry[0].CursorXLocation
		layerEntry.CursorYLocation = existingLayerEntry[0].CursorYLocation
		layerEntry.ZOrder = existingLayerEntry[0].ZOrder
		layerEntry.IsVisible = existingLayerEntry[0].IsVisible
		layerEntry.IsTopmost = existingLayerEntry[0].IsTopmost
		layerEntry.IsFocusable = existingLayerEntry[0].IsFocusable
		layerEntry.LayerAlias = existingLayerEntry[0].LayerAlias
		layerEntry.ParentAlias = existingLayerEntry[0].ParentAlias
		layerEntry.IsParent = existingLayerEntry[0].IsParent
		layerEntry.DefaultAttribute = existingLayerEntry[0].DefaultAttribute
		for currentRow := 0; currentRow < existingLayerEntry[0].Height; currentRow++ {
			var characterObjectArray = make([]CharacterEntryType, existingLayerEntry[0].Width)
			for currentCharacter := 0; currentCharacter < existingLayerEntry[0].Width; currentCharacter++ {
				characterObjectArray[currentCharacter] = NewCharacterEntry()
				characterObjectArray[currentCharacter].LayerAlias = layerAlias
				characterObjectArray[currentCharacter].ParentAlias = parentAlias
				characterObjectArray[currentCharacter] = existingLayerEntry[0].CharacterMemory[currentRow][currentCharacter]
			}
			layerEntry.CharacterMemory = append(layerEntry.CharacterMemory, characterObjectArray)
		}
	} else {
		layerEntry.Width = width
		layerEntry.Height = height
		layerEntry.IsVisible = true
		layerEntry.DefaultAttribute = NewAttributeEntry()
		for currentRow := 0; currentRow < height; currentRow++ {
			var characterObjectArray = make([]CharacterEntryType, width)
			for currentCharacter := 0; currentCharacter < width; currentCharacter++ {
				characterObjectArray[currentCharacter] = NewCharacterEntry()
				characterObjectArray[currentCharacter].LayerAlias = layerAlias
				characterObjectArray[currentCharacter].ParentAlias = parentAlias
			}
			layerEntry.CharacterMemory = append(layerEntry.CharacterMemory, characterObjectArray)
		}
	}
	return layerEntry
}

func InitializeCharacterMemory(layerEntry *LayerEntryType) {
	// This is used exclusively for clearing layer data.
	layerEntry.CharacterMemory = [][]CharacterEntryType{}
	for currentRow := 0; currentRow < layerEntry.Height; currentRow++ {
		var characterObjectArray = make([]CharacterEntryType, layerEntry.Width)
		for currentCharacter := 0; currentCharacter < layerEntry.Width; currentCharacter++ {
			characterObjectArray[currentCharacter] = NewCharacterEntry()
			characterObjectArray[currentCharacter].LayerAlias = layerEntry.LayerAlias
			characterObjectArray[currentCharacter].ParentAlias = layerEntry.ParentAlias
		}
		layerEntry.CharacterMemory = append(layerEntry.CharacterMemory, characterObjectArray)
	}
}
