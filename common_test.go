package consolizer

import "github.com/supercom32/consolizer/types"

/*
CommonTestSetup is a method which allows you to initialize a standard testing environment with multiple layers and a
default TUI style.

Example:
    layer1, layer2, layer3, styleEntry := CommonTestSetup()
*/
func CommonTestSetup() (*LayerInstanceType, *LayerInstanceType, *LayerInstanceType, types.TuiStyleEntryType) {
	commonResource.isDebugEnabled = true
	layerWidth := 40
	layerHeight := 20
	styleEntry := types.NewTuiStyleEntry()
	styleEntry.Window.LineDrawingTextForegroundColor = GetRGBColor(255, 0, 255)
	styleEntry.Window.LineDrawingTextBackgroundColor = GetRGBColor(0, 0, 255)
	InitializeTerminal(layerWidth, layerHeight)
	layer1 := AddLayer(0, 0, layerWidth, layerHeight, 1, nil)
	layer2 := AddLayer(3, 10, layerWidth, layerHeight, 2, nil)
	layer3 := AddLayer(0, 0, layerWidth, layerHeight, 3, nil)
	layer1.Color(4, 6)
	layer1.FillLayer("a1a2a3a4a5")
	layer2.Color(11, 9)
	layer2.FillLayer("a1a2a3a4a5")
	return layer1, layer2, layer3, styleEntry
}

/*
CommonTestSetupImages is a method which allows you to initialize a standard testing environment for image-related tests.

Example:
    layer1, layer2, layer3, styleEntry := CommonTestSetupImages()
*/
func CommonTestSetupImages() (*LayerInstanceType, *LayerInstanceType, *LayerInstanceType, types.TuiStyleEntryType) {
	commonResource.isDebugEnabled = true
	layerWidth := 50
	layerHeight := 20
	styleEntry := types.NewTuiStyleEntry()
	styleEntry.Window.LineDrawingTextForegroundColor = GetRGBColor(255, 0, 255)
	styleEntry.Window.LineDrawingTextBackgroundColor = GetRGBColor(0, 0, 255)
	InitializeTerminal(layerWidth, layerHeight)
	layer1 := AddLayer(0, 0, layerWidth, layerHeight, 1, nil)
	layer2 := AddLayer(0, 0, layerWidth, layerHeight, 2, nil)
	layer3 := AddLayer(0, 0, layerWidth, layerHeight, 3, nil)
	return layer1, layer2, layer3, styleEntry
}
