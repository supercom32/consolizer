package consolizer

import (
	"fmt"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/internal/memory"
	"strconv"
	"testing"
	"time"
)

func TestMainStub(test *testing.T) {
	testTextboxes()
	//testCheckboxes()
	//testDropdown()
	//testScrollBars()
	//testSelector()
	//testTextField()
	//testWindowMovement()
	//testButtonPressAction()
	RestoreTerminalSettings()
}

func testTextboxes () {
	commonResource.isDebugEnabled = false
	layerAlias1 := "Layer1"
	layerAlias2 := "Layer2"
	InitializeTerminal(80, 40)
	AddLayer(layerAlias1, 0, 0, 80, 40, 1, "")
	AddLayer(layerAlias2, 25, 15, 40, 20, 1, layerAlias1)
	Layer(layerAlias1)
	ColorRGB(255,0,0,0,0,0)
	FillLayer(layerAlias1, "#")
	FillLayer(layerAlias2, "@")
	Locate(0,0)
	Print("Enable ☑ Enable ○ ● (U+25CB, U+25CF)")
	styleEntry := memory.NewTuiStyleEntry()
	styleEntry.SelectorTextAlignment = 2
	textBox := AddTextbox(layerAlias1, "textbox1", styleEntry, 2, 2, 20, 5, true)
	textBox.setText("This is a test 1\nthis is ☑ second line which is very long and big 1\n文字文字文字文字文字文字文字文字 third line. 1")

	textBox2 := AddTextbox(layerAlias1, "textbox2", styleEntry, 40, 2, 20, 5, false)
	textBox2.setText("This is a test\nthis is ☑ second line which is very long and big\nthis is 文字 third line.")
	textBox3 := AddTextbox(layerAlias1, "textbox3", styleEntry, 1, 10, 20, 7, true)
	textBox3.setText("This is a test123456\nThis DDDDtesdfsfsdfsdfsdfsdfsdfsddffdsfdsst123456\nThis is GGGGGst1234\nThis is a ZZZZst123456\nThis is a test123456\nzzzzzzzzz\naaaa\ndddddd\nttttt\n222222\n555555")

	for {
		mouseXLocation, mouseYLocation, _, _ := memory.GetMouseStatus()
		characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
		Locate(0,0)
		Print(fmt.Sprintf("%d, %d   ", characterEntry.AttributeEntry.CellControlId, characterEntry.AttributeEntry.CellControlLocation))

		UpdateDisplay()
		key := Inkey()
		if key == "d" {
		}
		if key == "a" {
		}
		if key == "w" {
		}
		if key == "s" {
		}
		if key == "q" {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	DeleteAllLayers()
	RestoreTerminalSettings()
}

func testCheckboxes () {
	commonResource.isDebugEnabled = false
	layerAlias1 := "Layer1"
	layerAlias2 := "Layer2"
	InitializeTerminal(80, 40)
	AddLayer(layerAlias1, 0, 0, 80, 40, 1, "")
	AddLayer(layerAlias2, 20, 15, 40, 20, 1, layerAlias1)
	Layer(layerAlias1)
	ColorRGB(255,0,0,0,0,0)
	FillLayer(layerAlias1, "#")
	FillLayer(layerAlias2, "@")
	Locate(0,0)

	Print("Enable ☑ Enable ○ ● (U+25CB, U+25CF) ▾☒♪")
	styleEntry := memory.NewTuiStyleEntry()
	styleEntry.SelectorTextAlignment = 0
	AddCheckbox(layerAlias1, "checkbox1", "Enable 文字 Feature", styleEntry, 2, 2, true)
	for {
		UpdateDisplay()
		key := Inkey()
		if key == "w" {
		}
		if key == "q" {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	DeleteAllLayers()
	RestoreTerminalSettings()
}
func testDropdown() {
	commonResource.isDebugEnabled = false
	layerAlias1 := "Layer1"
	layerAlias2 := "Layer2"
	InitializeTerminal(80, 40)
	AddLayer(layerAlias1, 0, 0, 80, 40, 1, "")
	AddLayer(layerAlias2, 20, 15, 40, 20, 1, layerAlias1)
	Layer(layerAlias1)
	ColorRGB(255,0,0,0,0,0)
	FillLayer(layerAlias1, "#")
	FillLayer(layerAlias2, "@")
	Locate(0,0)
	Print("Enable ☑ Enable ○ ● (U+25CB, U+25CF)")
	styleEntry := memory.NewTuiStyleEntry()
	styleEntry.SelectorTextAlignment = 0
	selectionEntry := memory.NewSelectionEntry()
	selectionEntry.Add("1", "")
	selectionEntry.Add("2", "Enabl文e ○ ●")
	selectionEntry.Add("3", "☑ Enable ○ ●")
	selectionEntry.Add("4", "GET 文字文字")
	selectionEntry.Add("5", "IE文字文字")
	selectionEntry.Add("6", "DELETE")
	for i := 0; i < 20; i++ {
		selectionEntry.Add(strconv.Itoa(i), strconv.Itoa(i))
	}
	selectionEntry2 := memory.NewSelectionEntry()
	selectionEntry2.Add("1", "1")
	selectionEntry2.Add("2", "2")
	selectionEntry2.Add("3", "3")
	AddDropdown(layerAlias1, "myDropdown", styleEntry, selectionEntry, 2, 2, 8, 7, 0)
	AddDropdown(layerAlias1, "myDropdown2", styleEntry, selectionEntry2, 18, 9, 3, 7, 0)

	AddSelector(layerAlias1, "menuBar", styleEntry, selectionEntry, 6, 10, 4, 7, 3, 0, -1, false)
	AddSelector(layerAlias1, "menuBar2", styleEntry, selectionEntry, 10, 20, 4, 7, 1, 0, -1, false)

	styleEntry.SelectorTextAlignment = constants.AlignmentNoPadding
	AddSelector(layerAlias1, "menuBar3", styleEntry, selectionEntry, 6, 30, 4, 7, 1, 0, -1, false)
	for {
		UpdateDisplay()
		key := Inkey()
		//fmt.Print(key)
		if key == "w" {
			Locate(0,0)
			PrintLayer(layerAlias1, "                                                ")
			Locate(0,0)
			_, _, pressed, _ := memory.GetPreviousMouseStatus()
			PrintLayer(layerAlias1, "***" + strconv.Itoa(int(pressed)) + "***")
		}
		if key == "q" {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	DeleteAllLayers()
	RestoreTerminalSettings()
}
func testScrollBars() {
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
	styleEntry.SelectorTextAlignment = 0
	selectionEntry := memory.NewSelectionEntry()
	selectionEntry.Add("1", "OK")
	selectionEntry.Add("2", "CANCEL")
	selectionEntry.Add("3", "EXIT")
	selectionEntry.Add("4", "GOTO")
	selectionEntry.Add("5", "RUN")
	selectionEntry.Add("6", "DELETE")
	s1 := AddScrollBar(layerAlias1, "scrollBar1", styleEntry, 2, 2, 8,80,0, 1, false)
	s2 := AddScrollBar(layerAlias1, "scrollBar2", styleEntry, 10, 5, 8,8,4, 1,true)
	s1.setScrollValue(4)
	s2.setHandlePosition(4)
	for {
		UpdateDisplay()
		//x, y, _, _ := memory.GetMouseStatus()
		//a := getCellInformationUnderMouseCursor(x, y)
		LocateLayer(layerAlias1, 10, 0)
		PrintLayer(layerAlias1, "  ")
		LocateLayer(layerAlias1, 10, 0)
		PrintLayer(layerAlias1, strconv.Itoa(s1.getScrollValue()))
		LocateLayer(layerAlias1, 10, 1)
		PrintLayer(layerAlias1, "  ")
		LocateLayer(layerAlias1, 10, 1)
		PrintLayer(layerAlias1, strconv.Itoa(s2.getScrollValue()))
		key := Inkey()
		//fmt.Print(key)
		if key == "q" {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	DeleteAllLayers()
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
	styleEntry.SelectorTextAlignment = 0
	selectionEntry := memory.NewSelectionEntry()
	selectionEntry.Add("1", "OK")
	selectionEntry.Add("2", "CANCEL")
	selectionEntry.Add("3", "EXIT")
	selectionEntry.Add("4", "GOTO")
	selectionEntry.Add("5", "RUN")
	selectionEntry.Add("6", "DELETE")
	for i := 0; i < 20; i++ {
		selectionEntry.Add(strconv.Itoa(i), strconv.Itoa(i))
	}
	menuBarInstance := AddSelector(layerAlias1, "menuBar", styleEntry, selectionEntry, 2, 10, 4, 7, 3, 0, -1, false)
	menuBarInstance2 := AddSelector(layerAlias1, "menuBar2", styleEntry, selectionEntry, 10, 20, 4, 7, 1, 0, -1, false)
	LocateLayer(layerAlias1, 3, 3)
	PrintLayer(layerAlias1, menuBarInstance.layerAlias)
	for {
		UpdateDisplay()
		time.Sleep(50 * time.Millisecond)
		LocateLayer(layerAlias1, 3, 3)
		PrintLayer(layerAlias1, menuBarInstance.GetSelected())
		LocateLayer(layerAlias1, 3, 4)
		PrintLayer(layerAlias1, menuBarInstance2.GetSelected())
		key := Inkey()
		//fmt.Print(key)
		if key == "q" {
			break
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