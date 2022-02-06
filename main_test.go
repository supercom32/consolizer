package consolizer

import (
	"github.com/supercom32/consolizer/internal/memory"
	"strconv"
	"testing"
	"time"
)

func TestMainStub(test *testing.T) {
	testSelector()
	//testTextField()
	//testWindowMovement()
	//testButtonPressAction()
	RestoreTerminalSettings()
}
func testSelector() {
	commonResource.isDebugEnabled = false
	layerAlias1 := "Layer1"
	layerAlias2 := "Layer2"
	InitializeTerminal(80, 40)
	AddLayer(layerAlias1, 0, 0, 80, 40, 1, "")
	AddLayer(layerAlias2, 20, 15, 40, 20, 1, layerAlias1)
	Layer(layerAlias1)
	FillLayer(layerAlias1, "#")
	FillLayer(layerAlias2, "@")
	styleEntry := memory.NewTuiStyleEntry()
	selectionEntry := memory.NewSelectionEntry()
	selectionEntry.Add("1", "OK")
	selectionEntry.Add("2", "CANCEL")
	selectionEntry.Add("3", "EXIT")
	selectionEntry.Add("4", "GOTO")
	selectionEntry.Add("5", "RUN")
	selectionEntry.Add("6", "DELETE")
	menuBarInstance := AddSelector(layerAlias1, "menuBar", styleEntry, selectionEntry, 0, 10, 7, 4, 0)
	LocateLayer(layerAlias1, 3, 3)
	PrintLayer(layerAlias1, menuBarInstance.layerAlias)
	for i := 0; i < 200; i++ {
		UpdateDisplay()
		time.Sleep(50 * time.Millisecond)
		if menuBarInstance.GetSelected() == "6" {
			LocateLayer(layerAlias1, 3, 3)
			PrintLayer(layerAlias1, "HIT")
		}
	}
	DeleteAllLayers()
	RestoreTerminalSettings()
}

func testTextField() {
	commonResource.isDebugEnabled = false
	layerAlias1 := "Layer1"
	layerAlias2 := "Layer2"
	InitializeTerminal(80, 40)
	AddLayer(layerAlias1, 0, 0, 80, 40, 1, "")
	AddLayer(layerAlias2, 20, 15, 40, 20, 1, layerAlias1)
	Layer(layerAlias1)
	FillLayer(layerAlias1, "#")
	FillLayer(layerAlias2, "@")
	styleEntry := memory.NewTuiStyleEntry()
	textFieldInstance := AddTextField(layerAlias1, "textField", styleEntry, 0, 0, 10, 60, false, "Alex Chang is the man" )
	LocateLayer(layerAlias1, 3, 3)
	PrintLayer(layerAlias1, textFieldInstance.GetValue())
	memory.AddTextField(layerAlias1, "textField2", styleEntry, 0, 2, 10, 20, false, "Test" )
	for i := 0; i < 15; i++ {
		UpdateDisplay()
		time.Sleep(1000 * time.Millisecond)
	}
	DeleteAllLayers()
	RestoreTerminalSettings()
}

func testWindowMovement() {
	commonResource.isDebugEnabled = false
	layerAlias1 := "Layer1"
	layerAlias2 := "Layer2"
	layerAlias3 := "Layer3"
	layerAlias4 := "Layer4"
	xLocation := 1
	yLocation := 1
	InitializeTerminal(80, 40)
	AddLayer(layerAlias1, 0, 0, 80, 40, 1, "")
	AddLayer(layerAlias2, 20, 15, 40, 20, 1, layerAlias1)
	AddLayer(layerAlias3, 0, 0, 20, 10, 1, layerAlias2)
	AddLayer(layerAlias4, 15, 10, 20, 10, 2, layerAlias1)
	Layer(layerAlias1)
	FillLayer(layerAlias1, "#")
	FillLayer(layerAlias2, "@")
	Locate(xLocation, yLocation)
	Print("This is a test of the layer system")
	Locate(1,1)
	Print(strconv.Itoa(memory.GetLayer(layerAlias2).ZOrder))
	styleEntry := memory.NewTuiStyleEntry()
	DrawWindow(layerAlias3, styleEntry, 0, 0, 18, 9,true)
	DrawWindow(layerAlias4, styleEntry, 0, 0, 18, 9,true)
	for i := 0; i < 10; i++ {
		UpdateDisplay()
		time.Sleep(1000 * time.Millisecond)
	}
	DeleteAllLayers()
	RestoreTerminalSettings()
}

func testButtonPressAction() {
	commonResource.isDebugEnabled = false
	layerAlias1 := "Layer1"
	xLocation := 0
	yLocation := 0
	InitializeTerminal(80, 40)
	AddLayer(layerAlias1, 0, 0, 80, 40, 1, "")
	Locate(xLocation, yLocation)
	Print("This is a test")
	styleEntry := memory.NewTuiStyleEntry()
	AddButton(layerAlias1, "button1", "CANCEL", styleEntry, 2, 2, 10, 10)
	AddButton(layerAlias1, "button2", "OK", styleEntry, 15, 2, 10, 10)
	layerInformation := memory.GetLayer(layerAlias1)
	drawButtonsOnLayer(*layerInformation)
	UpdateDisplay()
	time.Sleep(15000 * time.Millisecond)
	DeleteAllLayers()
	RestoreTerminalSettings()
}