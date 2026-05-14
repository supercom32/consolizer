package consolizer

import (
	"github.com/supercom32/consolizer/memory"
	"github.com/supercom32/consolizer/stringformat"
	"strings"

	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/types"
)

type LabelInstanceType struct {
	BaseControlInstanceType
}

type labelType struct{}

var Label labelType
var Labels = memory.NewControlMemoryManager[types.LabelEntryType]()

// ============================================================================
// REGULAR ENTRY
// ============================================================================

/*
AddToTabIndex is a method which allows you to add the label to the tab index of its parent layer.

Example:

	labelInstance.AddToTabIndex()
*/
func (shared *LabelInstanceType) AddToTabIndex() {
	addTabIndex(shared.layerAlias, shared.controlAlias, constants.CellTypeLabel)
}

/*
SetIsTooltipEnabled is a method which allows you to enable or disable the tooltip associated with the label.

:param isEnabled: A boolean indicating whether the tooltip should be enabled.

Example:

	labelInstance.SetIsTooltipEnabled(true)
*/
func (shared *LabelInstanceType) SetIsTooltipEnabled(isEnabled bool) {
	labelEntry := Labels.Get(shared.layerAlias, shared.controlAlias)
	if labelEntry != nil && labelEntry.TooltipAlias != "" {
		tooltipEntry := Tooltips.Get(shared.layerAlias, labelEntry.TooltipAlias)
		tooltipEntry.IsEnabled = isEnabled
	}
}

/*
SetTooltipText is a method which allows you to update the text displayed in the label's tooltip.

:param text: The new text for the tooltip.

Example:

	labelInstance.SetTooltipText("New tooltip text")
*/
func (shared *LabelInstanceType) SetTooltipText(text string) {
	labelEntry := Labels.Get(shared.layerAlias, shared.controlAlias)
	tooltipEntry := Tooltips.Get(shared.layerAlias, labelEntry.TooltipAlias)
	tooltipEntry.Label = text
}

/*
Delete is a method which allows you to remove the label instance and its associated resources.

:return: A nil pointer of type LabelInstanceType.

Example:

	labelInstance = labelInstance.Delete()
*/
func (shared *LabelInstanceType) Delete() *LabelInstanceType {
	shared.BaseControlInstanceType.Delete()
	return nil
}

/*
Add is a method which allows you to create and add a new label to a specified layer.

:param layerAlias: The alias of the layer to which the label will be added.
:param labelAlias: A unique alias for the new label.
:param labelValue: The text content of the label.
:param styleEntry: The visual style to apply to the label.
:param xLocation: The X-coordinate for the label's position.
:param yLocation: The Y-coordinate for the label's position.
:param width: The width of the label in characters.

:return: A LabelInstanceType representing the newly created label.

Example:

	labelInstance := Label.Add("mainLayer", "myLabel", "Hello World", style, 10, 5, 20)
*/
func (shared *labelType) Add(layerAlias string, labelAlias string, labelValue string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int) LabelInstanceType {
	labelEntry := types.NewLabelEntry()
	labelEntry.StyleEntry = styleEntry
	labelEntry.Alias = labelAlias
	labelEntry.Label = labelValue
	labelEntry.XLocation = xLocation
	labelEntry.YLocation = yLocation
	labelEntry.Width = width
	labelEntry.TooltipAlias = stringformat.GetLastSortedUUID()
	// Use the ControlMemoryManager to add the label entry
	Labels.Add(layerAlias, labelAlias, &labelEntry)

	tooltipWidth := len(stringformat.GetRunesFromString(labelValue)) + 13 // Add padding

	tooltipInstance := Tooltip.Add(layerAlias, labelEntry.TooltipAlias, labelValue, styleEntry,
		labelEntry.XLocation, labelEntry.YLocation,
		labelEntry.Width+2, 1,
		labelEntry.XLocation, labelEntry.YLocation+1,
		tooltipWidth, 1,
		false, true, constants.DefaultTooltipHoverTime)
	tooltipInstance.SetEnabled(false)
	tooltipInstance.setParentControlAlias(labelAlias)
	var labelInstance LabelInstanceType
	labelInstance.layerAlias = layerAlias
	labelInstance.controlAlias = labelAlias
	labelInstance.controlType = constants.TYPE_LABEL
	return labelInstance
}

/*
Delete is a method which allows you to remove a label from a specific layer using its alias.

:param layerAlias: The alias of the layer containing the label.
:param labelAlias: The alias of the label to be deleted.

Example:

	Label.Delete("mainLayer", "myLabel")
*/
func (shared *labelType) Delete(layerAlias string, labelAlias string) {
	Labels.Remove(layerAlias, labelAlias)
}

/*
DeleteAll is a method which allows you to remove all labels associated with a specific layer.

:param layerAlias: The alias of the layer from which to remove all labels.

Example:

	Label.DeleteAll("mainLayer")
*/
func (shared *labelType) DeleteAll(layerAlias string) {
	Labels.RemoveAll(layerAlias)
}

/*
drawOnLayer is a method which allows you to draw all labels associated with a given layer entry.

:param layerEntry: The layer entry structure containing the labels to be drawn.

Example:

	Label.drawOnLayer(layerEntry)
*/
func (shared *labelType) drawOnLayer(layerEntry types.LayerEntryType) {
	layerAlias := layerEntry.LayerAlias
	for _, currentLabelEntry := range Labels.GetAllEntries(layerAlias) {
		labelEntry := currentLabelEntry
		drawLabel(&layerEntry, labelEntry.Alias, labelEntry.Label, labelEntry.StyleEntry, labelEntry.XLocation, labelEntry.YLocation, labelEntry.Width)
	}
}

/*
drawLabel is a method which allows you to draw a specific label onto a layer entry.

:param layerEntry: A pointer to the layer entry where the label will be drawn.
:param labelAlias: The alias of the label being drawn.
:param labelValue: The text content of the label.
:param styleEntry: The visual style to apply to the label.
:param xLocation: The X-coordinate for the label's position.
:param yLocation: The Y-coordinate for the label's position.
:param width: The width allocated for the label.

Example:

	drawLabel(&layerEntry, "myLabel", "Hello", style, 10, 5, 20)
*/
func drawLabel(layerEntry *types.LayerEntryType, labelAlias string, labelValue string, styleEntry types.TuiStyleEntryType, xLocation int, yLocation int, width int) {
	attributeEntry := types.NewAttributeEntry()
	attributeEntry.ForegroundColor = styleEntry.Label.ForegroundColor
	attributeEntry.BackgroundColor = styleEntry.Label.BackgroundColor
	attributeEntry.CellType = constants.CellTypeLabel
	attributeEntry.CellControlAlias = labelAlias
	emptyString := strings.Repeat(" ", width)
	layer.printLayer(layerEntry, attributeEntry, xLocation, yLocation, stringformat.GetRunesFromString(emptyString))
	if len(labelValue) > width {
		labelValue = string([]rune(labelValue)[:width-3])
		labelValue = labelValue + "..."
	}
	arrayOfRunes := stringformat.GetRunesFromString(labelValue)
	layer.printLayer(layerEntry, attributeEntry, xLocation, yLocation, arrayOfRunes)
}
