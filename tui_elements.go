package consolizer

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/internal/memory"
	"github.com/supercom32/consolizer/internal/stringformat"
)

/*
connectorEntryType defines all the attributes that a connector should have.
A connector is a character that joins lines together.
*/
type connectorEntryType struct {
	connectorIndex    int
	isConnectedTop    bool
	isConnectedBottom bool
	isConnectedLeft   bool
	isConnectedRight  bool
}

/*
These constants represent the various types of connectors available.
*/
const (
	horizontalConnector = 0
	verticalConnector   = 1
	upperLeftConnector  = 2
	upperRightConnector = 3
	lowerLeftConnector  = 4
	lowerRightConnector = 5
	rightTConnector     = 6
	leftTConnector      = 7
	upTConnector        = 8
	downTConnector      = 9
	crossConnector      = 10
)

/*
connectorList is an array that contains definitions of what all connectors
should look like. This is used for determining when a connector character is
required and which character would be suitable to connect lines with.
*/
var connectorList = []connectorEntryType{
	{
		connectorIndex:    horizontalConnector,
		isConnectedTop:    false,
		isConnectedBottom: false,
		isConnectedLeft:   true,
		isConnectedRight:  true,
	},
	{
		connectorIndex:    verticalConnector,
		isConnectedTop:    true,
		isConnectedBottom: true,
		isConnectedLeft:   false,
		isConnectedRight:  false,
	},
	{
		connectorIndex:    upperLeftConnector,
		isConnectedTop:    false,
		isConnectedBottom: true,
		isConnectedLeft:   false,
		isConnectedRight:  true,
	},
	{
		connectorIndex:    upperRightConnector,
		isConnectedTop:    false,
		isConnectedBottom: true,
		isConnectedLeft:   true,
		isConnectedRight:  false,
	},
	{
		connectorIndex:    lowerLeftConnector,
		isConnectedTop:    true,
		isConnectedBottom: false,
		isConnectedLeft:   false,
		isConnectedRight:  true,
	},
	{
		connectorIndex:    lowerRightConnector,
		isConnectedTop:    true,
		isConnectedBottom: false,
		isConnectedLeft:   true,
		isConnectedRight:  false,
	},
	{
		connectorIndex:    rightTConnector,
		isConnectedTop:    true,
		isConnectedBottom: true,
		isConnectedLeft:   false,
		isConnectedRight:  true,
	},
	{
		connectorIndex:    leftTConnector,
		isConnectedTop:    true,
		isConnectedBottom: true,
		isConnectedLeft:   true,
		isConnectedRight:  false,
	},
	{
		connectorIndex:    upTConnector,
		isConnectedTop:    true,
		isConnectedBottom: false,
		isConnectedLeft:   true,
		isConnectedRight:  true,
	},
	{
		connectorIndex:    downTConnector,
		isConnectedTop:    false,
		isConnectedBottom: true,
		isConnectedLeft:   true,
		isConnectedRight:  true,
	},
	{
		connectorIndex:    crossConnector,
		isConnectedTop:    true,
		isConnectedBottom: true,
		isConnectedLeft:   true,
		isConnectedRight:  true,
	},
}

/*
getConnectorIndexByCharacter allows you to obtain a connector index ID based on
the actual connector rune provided. This is done by comparing the provided rune
with the passed in TUI style entry and determining what type of connector it
is. In addition, the following information should be noted:

- If the rune specified could not be found in your TUI style entry, then '-1'
will be returned instead.
*/
func getConnectorIndexByCharacter(sourceCharacter rune, styleEntry memory.TuiStyleEntryType) int {
	connectorIndex := -1
	if sourceCharacter == styleEntry.HorizontalLine {
		connectorIndex = horizontalConnector
	}
	if sourceCharacter == styleEntry.VerticalLine {
		connectorIndex = verticalConnector
	}
	if sourceCharacter == styleEntry.UpperLeftCorner {
		connectorIndex = upperLeftConnector
	}
	if sourceCharacter == styleEntry.UpperRightCorner {
		connectorIndex = upperRightConnector
	}
	if sourceCharacter == styleEntry.LowerLeftCorner {
		connectorIndex = lowerLeftConnector
	}
	if sourceCharacter == styleEntry.LowerRightCorner {
		connectorIndex = lowerRightConnector
	}
	if sourceCharacter == styleEntry.RightSideTConnector {
		connectorIndex = rightTConnector
	}
	if sourceCharacter == styleEntry.LeftSideTConnector {
		connectorIndex = leftTConnector
	}
	if sourceCharacter == styleEntry.UpSideTConnector {
		connectorIndex = upTConnector
	}
	if sourceCharacter == styleEntry.DownSideTConnector {
		connectorIndex = downTConnector
	}
	if sourceCharacter == styleEntry.CrossConnector {
		connectorIndex = crossConnector
	}
	return connectorIndex
}

/*
getConnectorIndexByConnections allows you to obtain a connector index ID based
on the connections specified. This is useful for figuring out what connector
rune should be used, since you only need to specify what existing connections
are present. In addition, the following information should be noted:

- If the connector described could not be found determined, then '-1'
will be returned instead.
*/
func getConnectorIndexByConnections(isConnectedTop bool, isConnectedBottom bool, isConnectedLeft bool, isConnectedRight bool) int {
	connectorIndex := -1
	if isConnectedTop && isConnectedBottom && isConnectedLeft && isConnectedRight {
		connectorIndex = crossConnector
	} else if isConnectedTop && isConnectedBottom == false && isConnectedLeft && isConnectedRight {
		connectorIndex = upTConnector
	} else if isConnectedTop == false && isConnectedBottom && isConnectedLeft && isConnectedRight {
		connectorIndex = downTConnector
	} else if isConnectedTop && isConnectedBottom && isConnectedLeft == false && isConnectedRight {
		connectorIndex = rightTConnector
	} else if isConnectedTop && isConnectedBottom && isConnectedLeft && isConnectedRight == false {
		connectorIndex = leftTConnector
	} else if isConnectedTop == false && isConnectedBottom == false && isConnectedLeft && isConnectedRight {
		connectorIndex = horizontalConnector
	} else if isConnectedTop && isConnectedBottom && isConnectedLeft == false && isConnectedRight == false {
		connectorIndex = verticalConnector
	} else if isConnectedTop == false && isConnectedBottom && isConnectedLeft == false && isConnectedRight {
		connectorIndex = upperLeftConnector
	} else if isConnectedTop == false && isConnectedBottom && isConnectedLeft && isConnectedRight == false {
		connectorIndex = upperRightConnector
	} else if isConnectedTop && isConnectedBottom == false && isConnectedLeft == false && isConnectedRight {
		connectorIndex = lowerLeftConnector
	} else if isConnectedTop && isConnectedBottom == false && isConnectedLeft && isConnectedRight == false {
		connectorIndex = lowerRightConnector
	}
	return connectorIndex
}

/*
getConnectorCharacterByIndex allows you to obtain a connector character based
on the connector index and the TUI style provided. In addition, the following
information should be noted:

- If the connector described could not be found determined, then a NullRune
will be returned instead.
*/
func getConnectorCharacterByIndex(connectorIndex int, styleEntry memory.TuiStyleEntryType) rune {
	characterToReturn := constants.NullRune
	if connectorIndex == horizontalConnector {
		characterToReturn = styleEntry.HorizontalLine
	}
	if connectorIndex == verticalConnector {
		characterToReturn = styleEntry.VerticalLine
	}
	if connectorIndex == upperLeftConnector {
		characterToReturn = styleEntry.UpperLeftCorner
	}
	if connectorIndex == upperRightConnector {
		characterToReturn = styleEntry.UpperRightCorner
	}
	if connectorIndex == lowerLeftConnector {
		characterToReturn = styleEntry.LowerLeftCorner
	}
	if connectorIndex == lowerRightConnector {
		characterToReturn = styleEntry.LowerRightCorner
	}
	if connectorIndex == rightTConnector {
		characterToReturn = styleEntry.RightSideTConnector
	}
	if connectorIndex == leftTConnector {
		characterToReturn = styleEntry.LeftSideTConnector
	}
	if connectorIndex == downTConnector {
		characterToReturn = styleEntry.DownSideTConnector
	}
	if connectorIndex == upTConnector {
		characterToReturn = styleEntry.UpSideTConnector
	}
	if connectorIndex == crossConnector {
		characterToReturn = styleEntry.CrossConnector
	}
	return characterToReturn
}

/*
getConnectorCharacter allows you to select a source and target rune which
will then be combined into a new connector rune. The new connector rune
will have all lines joined together. For example, if a horizontal line
is combined with a vertical line, then a cross connector will be returned.
In addition, the following information should be noted:

- If the two runes could not be combined together, then the source rune
will be returned instead.
*/
func getConnectorCharacter(sourceCharacter rune, targetCharacter rune, styleEntry memory.TuiStyleEntryType) rune {
	connectorCharacter := sourceCharacter
	sourceCharacterIndex := getConnectorIndexByCharacter(sourceCharacter, styleEntry)
	targetCharacterIndex := getConnectorIndexByCharacter(targetCharacter, styleEntry)
	if sourceCharacterIndex != -1 && targetCharacterIndex != -1 {
		sourceConnectorEntry := connectorList[sourceCharacterIndex]
		targetConnectorEntry := connectorList[targetCharacterIndex]
		if targetConnectorEntry.isConnectedTop {
			sourceConnectorEntry.isConnectedTop = true
		}
		if targetConnectorEntry.isConnectedBottom {
			sourceConnectorEntry.isConnectedBottom = true
		}
		if targetConnectorEntry.isConnectedLeft {
			sourceConnectorEntry.isConnectedLeft = true
		}
		if targetConnectorEntry.isConnectedRight {
			sourceConnectorEntry.isConnectedRight = true
		}
		newConnectorIndex := getConnectorIndexByConnections(sourceConnectorEntry.isConnectedTop, sourceConnectorEntry.isConnectedBottom, sourceConnectorEntry.isConnectedLeft, sourceConnectorEntry.isConnectedRight)
		connectorCharacter = getConnectorCharacterByIndex(newConnectorIndex, styleEntry)
	}
	return connectorCharacter
}

/*
DrawVerticalLine allows you to draw a vertical line on a text layer. This
method also has the ability to draw connectors in case the line intersects
with other lines that have already been drawn. In addition, the following
information should be noted:

- If the the line to be drawn falls outside the area of the text layer
specified, then only the visible portion of the line will be drawn.
*/
func DrawVerticalLine(layerAlias string, styleEntry memory.TuiStyleEntryType, xLocation int, yLocation int, height int, isConnectorsDrawn bool) {
	layerEntry := memory.GetLayer(layerAlias)
	localAttributeEntry := memory.NewAttributeEntry()
	drawVerticalLine(layerEntry, styleEntry, localAttributeEntry, xLocation, yLocation, height, isConnectorsDrawn)
}

/*
DrawVerticalLine allows you to draw a vertical line on a text layer. This
method also has the ability to draw connectors in case the line intersects
with other lines that have already been drawn. In addition, the following
information should be noted:

- If the the line to be drawn falls outside the area of the text layer
specified, then only the visible portion of the line will be drawn.
*/
func drawVerticalLine(layerEntry *memory.LayerEntryType, styleEntry memory.TuiStyleEntryType, attributeEntry memory.AttributeEntryType, xLocation int, yLocation int, height int, isConnectorsDrawn bool) {
	localAttributeEntry := memory.NewAttributeEntry(&attributeEntry)
	localAttributeEntry.ForegroundColor = styleEntry.TextForegroundColor
	localAttributeEntry.BackgroundColor = styleEntry.TextBackgroundColor
	for currentRow := 0; currentRow < height; currentRow++ {
		if yLocation + currentRow >= 0 && xLocation >= 0 && yLocation+currentRow < layerEntry.Height && xLocation < layerEntry.Width {
			characterToDraw := constants.NullRune
			if isConnectorsDrawn && currentRow == 0 {
				characterToDraw = styleEntry.DownSideTConnector
			} else if isConnectorsDrawn && currentRow == height-1 {
				characterToDraw = styleEntry.UpSideTConnector
			} else {
				characterToDraw = styleEntry.VerticalLine
			}
			runeOnLayer := getRuneOnLayer(layerEntry, xLocation, yLocation+currentRow)
			stringToPrint := string(getConnectorCharacter(characterToDraw, runeOnLayer, styleEntry))
			if !isConnectorsDrawn { // If connectors are not drawn, set the target to empty so automatic drawing connectors does not happen
				stringToPrint = string(characterToDraw)
			}
			arrayOfRunes := stringformat.GetRunesFromString(stringToPrint)
			printLayer(layerEntry, localAttributeEntry, xLocation, yLocation+currentRow, arrayOfRunes)
		}
	}
}

/*
DrawHorizontalLine allows you to draw a horizontal line on a text layer. This
method also has the ability to draw connectors in case the line intersects
with other lines that have already been drawn. In addition, the following
information should be noted:

- If the the line to be drawn falls outside the area of the text layer
specified, then only the visible portion of the line will be drawn.
*/
func DrawHorizontalLine(layerAlias string, styleEntry memory.TuiStyleEntryType, xLocation int, yLocation int, width int, isConnectorsDrawn bool) {
	layerEntry := memory.GetLayer(layerAlias)
	localAttributeEntry := memory.NewAttributeEntry()
	drawHorizontalLine(layerEntry, styleEntry, localAttributeEntry, xLocation, yLocation, width, isConnectorsDrawn)
}

/*
DrawHorizontalLine allows you to draw a horizontal line on a text layer. This
method also has the ability to draw connectors in case the line intersects
with other lines that have already been drawn. In addition, the following
information should be noted:

- If the the line to be drawn falls outside the area of the text layer
specified, then only the visible portion of the line will be drawn.
*/
func drawHorizontalLine(layerEntry *memory.LayerEntryType, styleEntry memory.TuiStyleEntryType, attributeEntry memory.AttributeEntryType, xLocation int, yLocation int, width int, isConnectorsDrawn bool) {
	localAttributeEntry := memory.NewAttributeEntry(&attributeEntry)
	localAttributeEntry.ForegroundColor = styleEntry.TextForegroundColor
	localAttributeEntry.BackgroundColor = styleEntry.TextBackgroundColor
	for currentCharacter := 0; currentCharacter < width; currentCharacter++ {
		if yLocation < layerEntry.Height && xLocation+currentCharacter < layerEntry.Width {
			if yLocation < 0 || xLocation+currentCharacter < 0 || yLocation > layerEntry.Height || xLocation+currentCharacter > layerEntry.Width {
				continue
			}
			characterToDraw := constants.NullRune
			if isConnectorsDrawn && currentCharacter == 0 {
				characterToDraw = styleEntry.RightSideTConnector
			} else if isConnectorsDrawn && currentCharacter == width-1 {
				characterToDraw = styleEntry.LeftSideTConnector
			} else {
				characterToDraw = styleEntry.HorizontalLine
			}
			runeOnLayer := getRuneOnLayer(layerEntry, xLocation+currentCharacter, yLocation)
			stringToPrint := string(getConnectorCharacter(characterToDraw, runeOnLayer, styleEntry))
			if !isConnectorsDrawn {
				stringToPrint = string(characterToDraw)
			}
			arrayOfRunes := stringformat.GetRunesFromString(stringToPrint)
			printLayer(layerEntry, localAttributeEntry, xLocation+currentCharacter, yLocation, arrayOfRunes)
		}
	}
}

/*
DrawBorder allows you to draw a border on a given text layer. Borders differ
from frames since they are flat shaded and do not have a raised or sunken
look to them. In addition, the following information should be noted:

- If the border to be drawn falls outside the range of the specified layer,
then only the visible portion of the border will be drawn.

- The 'isInteractive' option allows you to specify if the window should
interact with the layer being drawn on. For example, when enabled if the user
drags the window title bar, the whole layer will move to simulate movement of
the window itself.
*/
func DrawBorder(layerAlias string, styleEntry memory.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isInteractive bool) {
	layerEntry := memory.GetLayer(layerAlias)
	localAttributeEntry := memory.NewAttributeEntry()
	drawBorder(layerEntry, styleEntry, localAttributeEntry, xLocation, yLocation, width, height, isInteractive)
}

/*
DrawBorder allows you to draw a border on a given text layer. Borders differ
from frames since they are flat shaded and do not have a raised or sunken
look to them. In addition, the following information should be noted:

- If the border to be drawn falls outside the range of the specified layer,
then only the visible portion of the border will be drawn.

- The 'isInteractive' option allows you to specify if the window should
interact with the layer being drawn on. For example, when enabled if the user
drags the window title bar, the whole layer will move to simulate movement of
the window itself.
*/
func drawBorder(layerEntry *memory.LayerEntryType, styleEntry memory.TuiStyleEntryType, attributeEntry memory.AttributeEntryType, xLocation int, yLocation int, width int, height int, isInteractive bool) {
	localAttributeEntry := memory.NewAttributeEntry(&attributeEntry)
	drawFrame(layerEntry, styleEntry, localAttributeEntry, constants.FrameStyleNormal, xLocation, yLocation, width, height, isInteractive)
}

/*
DrawFrameLabel allows you to draw a label for a frame. The label will
be automatically enclosed by the characters "[" and "]" to blend in
with a border of a frame.

- If the frame label to be drawn falls outside the range of the
specified layer, then only the visible portion of the border will be
drawn.
*/
func DrawFrameLabel(layerAlias string, styleEntry memory.TuiStyleEntryType, label string, xLocation int, yLocation int) {
	layerEntry := memory.GetLayer(layerAlias)
	drawFrameLabel(layerEntry, styleEntry, label, xLocation, yLocation)
}

/*
DrawFrameLabel allows you to draw a label for a frame. The label will
be automatically enclosed by the characters "[" and "]" to blend in
with a border of a frame.

- If the frame label to be drawn falls outside the range of the
specified layer, then only the visible portion of the border will be
drawn.
*/
func drawFrameLabel(layerEntry *memory.LayerEntryType, styleEntry memory.TuiStyleEntryType, label string, xLocation int, yLocation int) {
	attributeEntry := memory.NewAttributeEntry()
	attributeEntry.ForegroundColor= styleEntry.TextForegroundColor
	attributeEntry.BackgroundColor= styleEntry.TextBackgroundColor
	printLayer(layerEntry, attributeEntry, xLocation, yLocation, []rune("[ "))
	printLayer(layerEntry, attributeEntry, xLocation + 2 + len(label), yLocation, []rune(" ]"))
	attributeEntry.ForegroundColor = styleEntry.TextLabelColor
	attributeEntry.BackgroundColor = styleEntry.TextBackgroundColor
	printLayer(layerEntry, attributeEntry, xLocation + 2, yLocation, []rune(label))
}

/*
DrawFrame allows you to draw a frame on a given text layer. Frames differ
from borders since borders are flat shaded and do not have a raised or
sunken look to them. In addition, the following information should be noted:

- If the frame to be drawn falls outside the range of the specified layer,
then only the visible portion of the frame will be drawn.

- The 'isInteractive' option allows you to specify if the window should
interact with the layer being drawn on. For example, when enabled if the user
drags the window title bar, the whole layer will move to simulate movement of
the window itself.
*/
func DrawFrame(layerAlias string, styleEntry memory.TuiStyleEntryType, isRaised bool, xLocation int, yLocation int, width int, height int, isInteractive bool) {
	layerEntry := memory.GetLayer(layerAlias)
	localAttributeEntry := memory.NewAttributeEntry()
	if isRaised {
		drawFrame(layerEntry, styleEntry, localAttributeEntry, constants.FrameStyleRaised, xLocation, yLocation, width, height, isInteractive)
	} else {
		drawFrame(layerEntry, styleEntry, localAttributeEntry, constants.FrameStyleSunken, xLocation, yLocation, width, height, isInteractive)
	}
}

/*
drawFrame allows you to draw a frame on a given text layer. Frames differ
from borders since borders are flat shaded and do not have a raised or
sunken look to them. In addition, the following information should be noted:

- If the frame to be drawn falls outside the range of the specified layer,
then only the visible portion of the frame will be drawn.

- The 'isInteractive' option allows you to specify if the window should
interact with the layer being drawn on. For example, when enabled if the user
drags the window title bar, the whole layer will move to simulate movement of
the window itself.
*/
func drawFrame(layerEntry *memory.LayerEntryType, styleEntry memory.TuiStyleEntryType, attributeEntry memory.AttributeEntryType, frameStyle int, xLocation int, yLocation int, width int, height int, isInteractive bool) {
	localAttributeEntry := memory.NewAttributeEntry(&attributeEntry)
	localAttributeEntry.ForegroundColor = styleEntry.TextForegroundColor
	localAttributeEntry.BackgroundColor = styleEntry.TextBackgroundColor
	for currentRow := 0; currentRow < height; currentRow++ {
		for currentCharacter := 0; currentCharacter < width; currentCharacter++ {
			if yLocation+currentRow >= 0 || xLocation+currentCharacter >= 0 ||
				yLocation+currentRow < layerEntry.Height && xLocation+currentCharacter < layerEntry.Width {
				currentAttributeEntry := memory.NewAttributeEntry(&localAttributeEntry)
				characterToDraw := constants.NullRune
				if currentRow == 0 {
					if isInteractive {
						// Set cell type so the frame can be tied to mouse dragging and layer interactivity.
						currentAttributeEntry.CellType = constants.CellTypeFrameTop
					}
					if currentCharacter == 0 {
						characterToDraw = styleEntry.UpperLeftCorner
						if frameStyle == constants.FrameStyleSunken {
							currentAttributeEntry.ForegroundColor = constants.AnsiColorByIndex[0]
						}
					} else if currentCharacter == width-1 {
						characterToDraw = styleEntry.UpperRightCorner
						if frameStyle == constants.FrameStyleRaised {
							currentAttributeEntry.ForegroundColor = constants.AnsiColorByIndex[0]
						}
					} else {
						characterToDraw = styleEntry.HorizontalLine
						if frameStyle == constants.FrameStyleSunken {
							currentAttributeEntry.ForegroundColor = constants.AnsiColorByIndex[0]
						}
					}
				} else if currentRow == height-1 {
					if isInteractive {
						// Reset cell attribute if it was set to something layer interactive.
						currentAttributeEntry.CellType = 0
					}
					if currentCharacter == 0 {
						characterToDraw = styleEntry.LowerLeftCorner
						if frameStyle == constants.FrameStyleSunken {
							currentAttributeEntry.ForegroundColor = constants.AnsiColorByIndex[0]
						}
					} else if currentCharacter == width-1 {
						characterToDraw = styleEntry.LowerRightCorner
						if frameStyle == constants.FrameStyleRaised {
							currentAttributeEntry.ForegroundColor = constants.AnsiColorByIndex[0]
						}
					} else {
						characterToDraw = styleEntry.HorizontalLine
						if frameStyle == constants.FrameStyleRaised {
							currentAttributeEntry.ForegroundColor = constants.AnsiColorByIndex[0]
						}
					}
				} else {
					if currentCharacter == 0 {
						characterToDraw = styleEntry.VerticalLine
						if frameStyle == constants.FrameStyleSunken {
							currentAttributeEntry.ForegroundColor = constants.AnsiColorByIndex[0]
						}
					}
					if currentCharacter == width-1 {
						characterToDraw = styleEntry.VerticalLine
						if frameStyle == constants.FrameStyleRaised {
							currentAttributeEntry.ForegroundColor = constants.AnsiColorByIndex[0]
						}
					}
				}
				if characterToDraw != constants.NullRune {
					runeOnLayer := getRuneOnLayer(layerEntry, xLocation+currentCharacter, yLocation+currentRow)
					stringToPrint := string(getConnectorCharacter(characterToDraw, runeOnLayer, styleEntry))
					arrayOfRunes := stringformat.GetRunesFromString(stringToPrint)
					printLayer(layerEntry, currentAttributeEntry, xLocation+currentCharacter, yLocation+currentRow, arrayOfRunes)
				}
			}
		}
	}
}

/*
DrawWindow allows you to draw a window on a given text layer. Windows differ
from borders since the entire area the window surrounds gets filled with
a solid background color. In addition, the following information should be noted:

- If the window to be drawn falls outside the range of the specified layer,
then only the visible portion of the window will be drawn.

- The 'isInteractive' option allows you to specify if the window should
interact with the layer being drawn on. For example, when enabled if the user
drags the window title bar, the whole layer will move to simulate movement of
the window itself.
*/
func DrawWindow(layerAlias string, styleEntry memory.TuiStyleEntryType, xLocation int, yLocation int, width int, height int, isInteractive bool) {
	layerEntry := memory.GetLayer(layerAlias)
	localAttributeEntry := memory.NewAttributeEntry()
	drawWindow(layerEntry, styleEntry, localAttributeEntry, xLocation, yLocation, width, height, isInteractive)
}

/*
drawWindow allows you to draw a window on a given text layer. Windows differ
from borders since the entire area the window surrounds gets filled with
a solid background color. In addition, the following information should be noted:

- If the window to be drawn falls outside the range of the specified layer,
then only the visible portion of the window will be drawn.

- The 'isInteractive' option allows you to specify if the window should
interact with the layer being drawn on. For example, when enabled if the user
drags the window title bar, the whole layer will move to simulate movement of
the window itself.
*/
func drawWindow(layerEntry *memory.LayerEntryType, styleEntry memory.TuiStyleEntryType, attributeEntry memory.AttributeEntryType, xLocation int, yLocation int, width int, height int, isInteractive bool) {
	localAttributeEntry := memory.NewAttributeEntry(&attributeEntry)
	localAttributeEntry.ForegroundColor = styleEntry.TextForegroundColor
	localAttributeEntry.BackgroundColor = styleEntry.TextBackgroundColor
	if styleEntry.IsSquareFont {
		drawShadow(layerEntry, localAttributeEntry, xLocation+1, yLocation+1, width, height, 0.5)
	} else {
		drawShadow(layerEntry, localAttributeEntry, xLocation+2, yLocation+1, width, height, 0.5)
	}
	fillArea(layerEntry, localAttributeEntry, " ", xLocation, yLocation, width, height, constants.NullCellControlLocation)
	drawBorder(layerEntry, styleEntry, localAttributeEntry, xLocation, yLocation, width, height, isInteractive)
	if styleEntry.IsWindowHeaderDrawn {
		drawHorizontalLine(layerEntry, styleEntry, localAttributeEntry, xLocation, yLocation+2, width, true)
	}
	if styleEntry.IsWindowFooterDrawn {
		drawHorizontalLine(layerEntry, styleEntry, localAttributeEntry, xLocation, yLocation+height-3, width, true)
	}
}

/*
DrawShadow allows you to draw shadows on a given text layer. Shadows are simply
transparent areas which darken whatever text layers are underneath it by a
specified degree. In addition, the following information should be noted:

- The alpha value can range from 0.0 (no shadow) to 1.0 (totally black).
*/
func DrawShadow(layerAlias string, xLocation int, yLocation int, width int, height int, alphaValue float32) {
	layerEntry := memory.GetLayer(layerAlias)
	localAttributeEntry := memory.NewAttributeEntry()
	drawShadow(layerEntry, localAttributeEntry, xLocation, yLocation, width, height, alphaValue)
}

/*
drawShadow allows you to draw shadows on a given text layer. Shadows are simply
transparent areas which darken whatever text layers are underneath it by a
specified degree. In addition, the following information should be noted:

- The alpha value can range from 0.0 (no shadow) to 1.0 (totally black).
*/
func drawShadow(layerEntry *memory.LayerEntryType, attributeEntry memory.AttributeEntryType, xLocation int, yLocation int, width int, height int, alphaValue float32) {
	localAttributeEntry := memory.NewAttributeEntry(&attributeEntry)
	localAttributeEntry.ForegroundTransformValue = alphaValue
	localAttributeEntry.BackgroundTransformValue = alphaValue
	fillArea(layerEntry, localAttributeEntry, "", xLocation, yLocation, width, height, constants.NullCellControlLocation)
}

/*
FillArea allows you to fill an area of a given text layer with characters of
your choice. If you wish to fill the area with repeating text, simply provide
the string you wish to repeat. In addition, the following information should be
noted:

- If the area to fill falls outside the range of the specified layer, then only
the visible portion of the fill will be drawn.
 */
func FillArea(layerAlias string, fillCharacters string, xLocation int, yLocation int, width int, height int) {
	layerEntry := memory.GetLayer(layerAlias)
	attributeEntry := layerEntry.DefaultAttribute
	fillArea(layerEntry, attributeEntry, fillCharacters, xLocation, yLocation, width, height, constants.NullCellControlLocation)
}

/*
fillArea allows you to fill an area of a given text layer with characters of
your choice. If you wish to fill the area with repeating text, simply provide
the string you wish to repeat. In addition, the following information should be
noted:

- If the area to fill falls outside the range of the specified layer, then only
the visible portion of the fill will be drawn.
*/
func fillArea(layerEntry *memory.LayerEntryType, attributeEntry memory.AttributeEntryType, fillCharacters string, xLocation int, yLocation int, width int, height int, startingControlLocation int) {
	currentFillCharacterIndex := 0
	arrayOfRunes := stringformat.GetRunesFromString(fillCharacters)
	for currentRow := 0; currentRow < height; currentRow++ {
		// We make the cell location equal to the current location being printed.
		attributeEntry.CellControlLocation = startingControlLocation + currentRow
		for currentCharacter := 0; currentCharacter < width; currentCharacter++ {
			if yLocation >=0 && yLocation < layerEntry.Height && xLocation+currentCharacter >= 0 && xLocation+currentCharacter < layerEntry.Width {
				if len(arrayOfRunes) == 0 {
					printLayer(layerEntry, attributeEntry, xLocation+currentCharacter, yLocation+currentRow, []rune{0})
				} else {
					printLayer(layerEntry, attributeEntry, xLocation+currentCharacter, yLocation+currentRow, []rune{arrayOfRunes[currentFillCharacterIndex]})
					// Double Width characters advance by 2 spaces. But what happens to characters between? Get lost?
					if stringformat.IsRuneCharacterWide(arrayOfRunes[currentFillCharacterIndex]) {
						currentCharacter++
					}
				}
				currentFillCharacterIndex++
				if currentFillCharacterIndex >= len(arrayOfRunes) {
					currentFillCharacterIndex = 0
				}
			}
		}
	}
}

/*
FillLayer allows you to fill an entire layer with characters of your choice.
If you wish to fill the layer with repeating text, simply provide the string
you wish to repeat.
*/
func FillLayer(layerAlias string, fillCharacters string) {
	layerEntry := memory.GetLayer(layerAlias)
	attributeEntry := layerEntry.DefaultAttribute
	fillLayer(layerEntry, attributeEntry, fillCharacters)
}

/*
fillLayer allows you to fill an entire layer with characters of your choice.
If you wish to fill the layer with repeating text, simply provide the string
you wish to repeat.
*/
func fillLayer(layerEntry *memory.LayerEntryType, attributeEntry memory.AttributeEntryType, fillCharacters string) {
	fillArea(layerEntry, attributeEntry, fillCharacters, 0, 0, layerEntry.Width, layerEntry.Height, constants.NullCellControlLocation)
}

/*
DrawBar allows you to draw a horizontal bar on a given text layer row. This is
useful for drawing application headers or status bar footers.
*/
func DrawBar(layerAlias string, xLocation int, yLocation int, barLength int, fillCharacters string) {
	layerEntry := memory.GetLayer(layerAlias)
	attributeEntry := layerEntry.DefaultAttribute
	fillArea(layerEntry, attributeEntry, fillCharacters, xLocation, yLocation, barLength, 1, constants.NullCellControlLocation)
}

/*
fillLayer allows you to fill an entire layer with characters of your choice.
If you wish to fill the layer with repeating text, simply provide the string
you wish to repeat.
*/
func GetDarkenedCharacterEntry(characterEntry *memory.CharacterEntryType, alphaValue float32) memory.CharacterEntryType {
	var newCharacterEntry = memory.NewCharacterEntry(characterEntry)
	foregroundColor := newCharacterEntry.AttributeEntry.ForegroundColor
	backgroundColor := newCharacterEntry.AttributeEntry.BackgroundColor
	newCharacterEntry.AttributeEntry.ForegroundColor = GetDarkenedColor(foregroundColor, alphaValue)
	newCharacterEntry.AttributeEntry.BackgroundColor = GetDarkenedColor(backgroundColor, alphaValue)
	return newCharacterEntry
}


/*
GetDarkenedColor allows you to obtain a color that has been darkened
uniformly by a specific amount. In addition, the following information
should be noted:

- The percent change can range from 0.0 (totally dark) to 1.0
(no difference).

- If you pass in a percent change of less than 0.0 or greater
than 1.0, a panic will be generated to fail as fast as possible.
*/
func GetDarkenedColor(color constants.ColorType, percentChange float32) constants.ColorType {
	var redColorIndex int32
	var greenColorIndex int32
	var blueColorIndex int32
	if percentChange < 0 || percentChange > 1 {
		panic(fmt.Sprintf("The specified brightness percent value of '%f' is invalid!", percentChange))
	}
	redColorIndex, greenColorIndex, blueColorIndex = GetRGBColorComponents(color)
	redColorIndex = int32(float32(redColorIndex) * percentChange)
	greenColorIndex = int32(float32(greenColorIndex) * percentChange)
	blueColorIndex = int32(float32(blueColorIndex) * percentChange)
	newColor := tcell.NewRGBColor(redColorIndex, greenColorIndex, blueColorIndex)
	return constants.ColorType(newColor)
}