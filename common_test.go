package consolizer

import "supercom32.net/consolizer/types"

func CommonTestSetup() (*LayerInstanceType, *LayerInstanceType, *LayerInstanceType, types.TuiStyleEntryType) {
	commonResource.isDebugEnabled = true
	layerWidth := 40
	layerHeight := 20
	styleEntry := types.NewTuiStyleEntry()
	styleEntry.LineDrawingTextForegroundColor = GetRGBColor(255, 0, 255)
	styleEntry.LineDrawingTextBackgroundColor = GetRGBColor(0, 0, 255)
	InitializeTerminal(layerWidth, layerHeight)
	layer1 := AddLayer(0, 0, layerWidth, layerHeight, 1, nil)
	layer2 := AddLayer(3, 10, layerWidth, layerHeight, 2, nil)
	layer3 := AddLayer(0, 0, layerWidth, layerHeight, 3, nil)
	Layer(layer1)
	Color(4, 6)
	layer1.FillLayer("a1a2a3a4a5")
	Layer(layer2)
	Color(11, 9)
	layer2.FillLayer("a1a2a3a4a5")
	return &layer1, &layer2, &layer3, styleEntry
}
