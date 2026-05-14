package types

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/klauspost/compress/zstd"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/stringformat"
	"github.com/supercom32/filesystem"
	"io"
	"os"
)

const (
	LayerMagicHeader = "CONS"
	layerVersion     = 1
)

const (
	flagFgTransparent = 1 << 0
	flagBgTransparent = 1 << 1
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

/*
MarshalJSON is a method which allows you to marshaljson.

Example:

	instance.MarshalJSON()
*/
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

/*
GetBasicAnsiString is a method which allows you to getbasicansistring.

:return: string

Example:

	instance.GetBasicAnsiString()
*/
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

/*
GetBasicAnsiStringAsBase64 is a method which allows you to getbasicansistringasbase64.

:return: string

Example:

	instance.GetBasicAnsiStringAsBase64()
*/
func (shared LayerEntryType) GetBasicAnsiStringAsBase64() string {
	ansiString := shared.GetBasicAnsiString()
	err := filesystem.WriteBytesToFile("/tmp/test_output/output.ans", []byte(ansiString), 0)
	if err != nil {
		panic(err)
	}
	return stringformat.GetStringAsBase64(ansiString)
}

/*
GetBasicAnsiStringAsBase642 is a method which allows you to getbasicansistringasbase642.

:return: string

Example:

	instance.GetBasicAnsiStringAsBase642()
*/
func (shared LayerEntryType) GetBasicAnsiStringAsBase642() string {
	ansiString := shared.GetBasicAnsiString()
	err := filesystem.WriteBytesToFile("/tmp/output.ans", []byte(ansiString), 0)
	if err != nil {
		panic(err)
	}
	return stringformat.GetStringAsBase64(ansiString)
}

/*
GetAnsiStringFromBase64 is a method which allows you to getansistringfrombase64.

:param base64String: The base64String parameter.

:return: string

Example:

	instance.GetAnsiStringFromBase64(base64String)
*/
func (shared LayerEntryType) GetAnsiStringFromBase64(base64String string) string {
	return stringformat.GetStringFromBase64(base64String)
}

/*
WriteAnsiStringFromBase64 is a method which allows you to decodes a base64 string to ANSI and writes it to the specified
file. In addition, the following information should be noted:

- This is useful for comparing expected and actual values when tests fail.

:param base64String: The base64String parameter.

:return: error

Example:

	WriteAnsiStringFromBase64(base64String)
*/
func WriteAnsiStringFromBase64(base64String string) error {
	ansiString := stringformat.GetStringFromBase64(base64String)
	return filesystem.WriteBytesToFile("/tmp/test_output/expected.ans", []byte(ansiString), 0)
}

/*
GetAnsiForegroundColorString is a method which allows you to getansiforegroundcolorstring.

:param color: The color parameter.

:return: string

Example:

	instance.GetAnsiForegroundColorString(color)
*/
func (shared LayerEntryType) GetAnsiForegroundColorString(color constants.ColorType) string {
	var ansiString string
	redIndex, greenIndex, blueIndex := shared.GetRGBColorComponents(color)
	ansiString = "\u001b[38;2;" + stringformat.GetIntAsString(redIndex) + ";" + stringformat.GetIntAsString(greenIndex) + ";" + stringformat.GetIntAsString(blueIndex) + "m"
	return ansiString
}

/*
GetAnsiBackgroundColorString is a method which allows you to getansibackgroundcolorstring.

:param color: The color parameter.

:return: string

Example:

	instance.GetAnsiBackgroundColorString(color)
*/
func (shared LayerEntryType) GetAnsiBackgroundColorString(color constants.ColorType) string {
	var ansiString string
	redIndex, greenIndex, blueIndex := shared.GetRGBColorComponents(color)
	ansiString = "\u001b[48;2;" + stringformat.GetIntAsString(redIndex) + ";" + stringformat.GetIntAsString(greenIndex) + ";" + stringformat.GetIntAsString(blueIndex) + "m"
	return ansiString
}

/*
GetAnsiLocateString is a method which allows you to getansilocatestring.

:param xLocation: The xLocation parameter.
:param yLocation: The yLocation parameter.

:return: string

Example:

	instance.GetAnsiLocateString(xLocation, yLocation)
*/
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

/*
GetEntryAsJsonDump is a method which allows you to getentryasjsondump.

:return: string

Example:

	instance.GetEntryAsJsonDump()
*/
func (shared LayerEntryType) GetEntryAsJsonDump() string {
	j, err := json.Marshal(shared)
	if err != nil {
		panic(err)
	}
	return string(j)
}

/*
GetRGBColorComponents is a method which allows you to getrgbcolorcomponents.

:param color: The color parameter.

Example:

	instance.GetRGBColorComponents(color)
*/
func (shared LayerEntryType) GetRGBColorComponents(color constants.ColorType) (int32, int32, int32) {
	var redColorIndex int32
	var greenColorIndex int32
	var blueColorIndex int32
	redColorIndex, greenColorIndex, blueColorIndex = tcell.Color.RGB(tcell.Color(color))
	return redColorIndex, greenColorIndex, blueColorIndex
}

/*
NewLayerEntry is a constructor which allows you to newlayerentry.

:param layerAlias: The layerAlias parameter.
:param parentAlias: The parentAlias parameter.
:param width: The width parameter.
:param height: The height parameter.
:param existingLayerEntry: The existingLayerEntry parameter.

:return: LayerEntryType

Example:

	NewLayerEntry(layerAlias, parentAlias, width, height, existingLayerEntry)
*/
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

/*
InitializeCharacterMemory is a method which allows you to initializecharactermemory.

:param layerEntry: The layerEntry parameter.

Example:

	InitializeCharacterMemory(layerEntry)
*/
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

/*
SaveLayer is a method which allows you to writes the layer to a file with zstd compression.

:param path: The path parameter.

:return: error

Example:

	instance.SaveLayer(path)
*/
func (shared *LayerEntryType) SaveLayer(path string) error {
	// Open file
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Create zstd writer
	zstdWriter, err := zstd.NewWriter(file)
	if err != nil {
		return fmt.Errorf("failed to create zstd writer: %w", err)
	}
	defer zstdWriter.Close()

	// Buffered writer
	writer := bufio.NewWriter(zstdWriter)

	// --- Header ---
	if _, err := writer.Write([]byte(LayerMagicHeader)); err != nil {
		return fmt.Errorf("failed to write magic header: %w", err)
	}
	if err := binary.Write(writer, binary.LittleEndian, uint16(layerVersion)); err != nil {
		return fmt.Errorf("failed to write version: %w", err)
	}

	height := uint16(len(shared.CharacterMemory))
	var width uint16
	if height > 0 {
		width = uint16(len(shared.CharacterMemory[0]))
	}
	if err := binary.Write(writer, binary.LittleEndian, width); err != nil {
		return fmt.Errorf("failed to write width: %w", err)
	}
	if err := binary.Write(writer, binary.LittleEndian, height); err != nil {
		return fmt.Errorf("failed to write height: %w", err)
	}

	// --- Layer Data ---
	for y := 0; y < int(height); y++ {
		for x := 0; x < int(width); x++ {
			entry := shared.CharacterMemory[y][x]

			// Rune
			if err := binary.Write(writer, binary.LittleEndian, int32(entry.Character)); err != nil {
				return fmt.Errorf("failed to write character at (%d,%d): %w", x, y, err)
			}

			// Colors (write full uint64 exactly)
			if err := binary.Write(writer, binary.LittleEndian, entry.AttributeEntry.ForegroundColor); err != nil {
				return fmt.Errorf("failed to write foreground color at (%d,%d): %w", x, y, err)
			}
			if err := binary.Write(writer, binary.LittleEndian, entry.AttributeEntry.BackgroundColor); err != nil {
				return fmt.Errorf("failed to write background color at (%d,%d): %w", x, y, err)
			}

			// Flags
			var flags byte
			if entry.AttributeEntry.IsForegroundTransparent {
				flags |= flagFgTransparent
			}
			if entry.AttributeEntry.IsBackgroundTransparent {
				flags |= flagBgTransparent
			}
			if err := writer.WriteByte(flags); err != nil {
				return fmt.Errorf("failed to write flags at (%d,%d): %w", x, y, err)
			}
		}
	}

	// Flush before closing zstd
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	return nil
}

/*
LoadLayer is a method which allows you to reads a layer from a file.

:param path: The path parameter.

:return: error

Example:

	instance.LoadLayer(path)
*/
func (shared *LayerEntryType) LoadLayer(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	return shared.LoadLayerFromBytes(data)
}

/*
LoadLayerFromBytes is a method which allows you to reads a layer from a byte slice.

:param data: The data parameter.

:return: error

Example:

	instance.LoadLayerFromBytes(data)
*/
func (shared *LayerEntryType) LoadLayerFromBytes(data []byte) error {
	reader := bytes.NewReader(data)

	// Create zstd reader
	zstdReader, err := zstd.NewReader(reader)
	if err != nil {
		return fmt.Errorf("failed to create zstd reader: %w", err)
	}
	defer zstdReader.Close()

	buffReader := bufio.NewReader(zstdReader)

	// --- Header ---
	magicHeader := make([]byte, len(LayerMagicHeader))
	if _, err := io.ReadFull(buffReader, magicHeader); err != nil {
		return fmt.Errorf("failed to read magic header: %w", err)
	}
	if string(magicHeader) != LayerMagicHeader {
		return fmt.Errorf("not a valid layer file")
	}

	var fileVersion uint16
	if err := binary.Read(buffReader, binary.LittleEndian, &fileVersion); err != nil {
		return fmt.Errorf("failed to read version: %w", err)
	}
	if fileVersion != layerVersion {
		return fmt.Errorf("unsupported version %d", fileVersion)
	}

	var width, height uint16
	if err := binary.Read(buffReader, binary.LittleEndian, &width); err != nil {
		return fmt.Errorf("failed to read width: %w", err)
	}
	if err := binary.Read(buffReader, binary.LittleEndian, &height); err != nil {
		return fmt.Errorf("failed to read height: %w", err)
	}

	shared.Width = int(width)
	shared.Height = int(height)
	characterMemory := make([][]CharacterEntryType, height)

	// --- Layer Data ---
	for y := 0; y < int(height); y++ {
		characterMemory[y] = make([]CharacterEntryType, width)
		for x := 0; x < int(width); x++ {
			var char int32
			if err := binary.Read(buffReader, binary.LittleEndian, &char); err != nil {
				return fmt.Errorf("failed to read character at (%d,%d): %w", x, y, err)
			}

			var fgColor, bgColor uint64
			if err := binary.Read(buffReader, binary.LittleEndian, &fgColor); err != nil {
				return fmt.Errorf("failed to read foreground color at (%d,%d): %w", x, y, err)
			}
			if err := binary.Read(buffReader, binary.LittleEndian, &bgColor); err != nil {
				return fmt.Errorf("failed to read background color at (%d,%d): %w", x, y, err)
			}

			flags, err := buffReader.ReadByte()
			if err != nil {
				return fmt.Errorf("failed to read flags at (%d,%d): %w", x, y, err)
			}

			characterMemory[y][x] = CharacterEntryType{
				Character: char,
				AttributeEntry: AttributeEntryType{
					ForegroundColor:         constants.ColorType(fgColor),
					BackgroundColor:         constants.ColorType(bgColor),
					IsForegroundTransparent: flags&flagFgTransparent != 0,
					IsBackgroundTransparent: flags&flagBgTransparent != 0,
				},
			}
		}
	}

	shared.CharacterMemory = characterMemory
	return nil
}
