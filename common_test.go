package consolizer

import (
	"encoding/base64"
	"fmt"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/types"
	"os"
	"testing"
)

/*
UpdateMasterImages is a method which allows you to update the master regression files for a test. In addition, the
following should be noted:

- If the environment variable UPDATE_MASTER_IMAGES is set to true, it will always perform the update.

Example:
    isUpdated := UpdateMasterImages(false, "TestName", "base64data")
*/
func UpdateMasterImages(isUpdateRequested bool, testSuiteName string, testCaseName string, ansiBase64 string) bool {
	if os.Getenv("UPDATE_MASTER_IMAGES") == "true" || isUpdateRequested {
		fullPath := constants.MasterImagesPath + testSuiteName + "/"
		os.MkdirAll(fullPath, 0755)
		os.WriteFile(fullPath+testCaseName+".base64", []byte(ansiBase64), 0644)
		ansiData, _ := base64.StdEncoding.DecodeString(ansiBase64)
		os.WriteFile(fullPath+testCaseName+".ansi", ansiData, 0644)
		fmt.Println("Updated master image for: " + testSuiteName + "/" + testCaseName)
		return true
	}
	return false
}

/*
LoadMasterImage is a method which allows you to load a master regression file for a test.

Example:
    expectedValue := LoadMasterImage("TestSuite", "TestCase")
*/
func LoadMasterImage(testSuiteName string, testCaseName string) string {
	expectedValueBytes, _ := os.ReadFile(constants.MasterImagesPath + testSuiteName + "/" + testCaseName + ".base64")
	return string(expectedValueBytes)
}

/*
CommonTestSetup is a method which initializes a standard testing environment with multiple layers and a
default TUI style. In addition, the following should be noted:

- Registers RestoreTerminalSettings as a cleanup on test, so the terminal session this call starts is always torn
  down when the calling test finishes, even on a failure or panic. Without this, the background goroutines
  InitializeTerminal starts keep running into later tests and racing their own terminal state.

Example:
    layer1, layer2, layer3, styleEntry := CommonTestSetup(test)
*/
func CommonTestSetup(test testing.TB) (*LayerInstanceType, *LayerInstanceType, *LayerInstanceType, types.TuiStyleEntryType) {
	test.Helper()
	commonResource.isDebugEnabled = true
	layerWidth := 40
	layerHeight := 20
	styleEntry := types.NewTuiStyleEntry()
	styleEntry.Window.LineDrawingTextForegroundColor = GetRGBColor(255, 0, 255)
	styleEntry.Window.LineDrawingTextBackgroundColor = GetRGBColor(0, 0, 255)
	InitializeTerminal(layerWidth, layerHeight)
	test.Cleanup(RestoreTerminalSettings)
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
CommonTestSetupImages is a method which initializes a standard testing environment for image-related tests. In
addition, the following should be noted:

- Registers RestoreTerminalSettings as a cleanup on test, so the terminal session this call starts is always torn
  down when the calling test finishes, even on a failure or panic. Without this, the background goroutines
  InitializeTerminal starts keep running into later tests and racing their own terminal state.

Example:
    layer1, layer2, layer3, styleEntry := CommonTestSetupImages(test)
*/
func CommonTestSetupImages(test testing.TB) (*LayerInstanceType, *LayerInstanceType, *LayerInstanceType, types.TuiStyleEntryType) {
	test.Helper()
	commonResource.isDebugEnabled = true
	layerWidth := 50
	layerHeight := 20
	styleEntry := types.NewTuiStyleEntry()
	styleEntry.Window.LineDrawingTextForegroundColor = GetRGBColor(255, 0, 255)
	styleEntry.Window.LineDrawingTextBackgroundColor = GetRGBColor(0, 0, 255)
	InitializeTerminal(layerWidth, layerHeight)
	test.Cleanup(RestoreTerminalSettings)
	layer1 := AddLayer(0, 0, layerWidth, layerHeight, 1, nil)
	layer2 := AddLayer(0, 0, layerWidth, layerHeight, 2, nil)
	layer3 := AddLayer(0, 0, layerWidth, layerHeight, 3, nil)
	return layer1, layer2, layer3, styleEntry
}

/*
CommonTestSetupHighResolutionImages is a method which initializes a standard testing environment for high resolution
image tests. In addition, the following should be noted:

- Registers RestoreTerminalSettings as a cleanup on test, so the terminal session this call starts is always torn
  down when the calling test finishes, even on a failure or panic. Without this, the background goroutines
  InitializeTerminal starts keep running into later tests and racing their own terminal state.

Example:
    layer1, layer2, layer3, styleEntry := CommonTestSetupHighResolutionImages(test)
*/
func CommonTestSetupHighResolutionImages(test testing.TB) (*LayerInstanceType, *LayerInstanceType, *LayerInstanceType, types.TuiStyleEntryType) {
	test.Helper()
	commonResource.isDebugEnabled = true
	layerWidth := 140
	layerHeight := 50
	styleEntry := types.NewTuiStyleEntry()
	styleEntry.Window.LineDrawingTextForegroundColor = GetRGBColor(255, 0, 255)
	styleEntry.Window.LineDrawingTextBackgroundColor = GetRGBColor(0, 0, 255)
	InitializeTerminal(layerWidth, layerHeight)
	test.Cleanup(RestoreTerminalSettings)
	layer1 := AddLayer(0, 0, layerWidth, layerHeight, 1, nil)
	layer2 := AddLayer(0, 0, layerWidth, layerHeight, 2, nil)
	layer3 := AddLayer(0, 0, layerWidth, layerHeight, 3, nil)
	return layer1, layer2, layer3, styleEntry
}
