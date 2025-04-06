package consolizer

import (
	"fmt"
	"log"
	"strconv"
	"supercom32.net/consolizer/constants"
	"supercom32.net/consolizer/types"
	"testing"
	"time"
)

func TestMainStub(test *testing.T) {
	// testRadioButtons()
	// testTextboxes()
	// testCheckboxes()
	// testDropdown()
	// testScrollBars()
	// testSelector()
	testTextField()
	//testProgressBar()
	// testWindowMovement()
	// testButtonPressAction()
	//RestoreTerminalSettings()
}
func testProgressBar() {
	commonResource.isDebugEnabled = false
	InitializeTerminal(80, 40)
	layer1 := AddLayer(0, 0, 80, 40, 1, nil)
	layer2 := AddLayer(20, 15, 40, 20, 1, &layer1)
	Layer(layer1)
	ColorRGB(255, 0, 0, 0, 0, 0)
	layer1.FillLayer("#")
	layer2.FillLayer("@")
	Locate(0, 0)
	styleEntry := types.NewTuiStyleEntry()
	styleEntry.SelectorTextAlignment = 0
	styleEntry.ProgressBarUnfilledBackgroundColor = constants.ColorBrightGreen
	layer1.AddProgressBar("Any Label", styleEntry, 21, 5, 40, 3, 10, 20, true)

	for {
		UpdateDisplay(false)
		key := string(Inkey())
		if key == "w" {
			log.Printf("TESTING")
			Locate(0, 0)
		}
		if key == "a" {
			RefreshDisplay()
		}
		if key == "q" {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	DeleteAllLayers()
	RestoreTerminalSettings()
}

func testRadioButtons() {
	commonResource.isDebugEnabled = false
	InitializeTerminal(80, 40)
	layer1 := AddLayer(0, 0, 80, 40, 1, nil)
	layer2 := AddLayer(20, 15, 40, 20, 1, &layer1)
	Layer(layer1)
	ColorRGB(255, 0, 0, 0, 0, 0)
	layer1.FillLayer("#")
	layer2.FillLayer("@")
	Locate(0, 0)

	Print("Enable ☑ Enable ○ ● (U+25CB, U+25CF) ▾☒♪")
	styleEntry := types.NewTuiStyleEntry()
	styleEntry.SelectorTextAlignment = 0
	radioButton := layer1.AddRadioButton("Enable 文字 Feature 1", styleEntry, 2, 2, 0, true)
	layer1.AddRadioButton("Enable 文字 Feature 2", styleEntry, 2, 3, 0, false)
	layer1.AddRadioButton("Enable 文字 Feature 3", styleEntry, 2, 4, 0, false)

	layer1.AddRadioButton("Option 1.1", styleEntry, 2, 6, 1, true)
	layer1.AddRadioButton("Option 1.2", styleEntry, 2, 7, 1, false)
	layer1.AddRadioButton("Option 1.3", styleEntry, 2, 8, 1, false)

	for {
		UpdateDisplay(false)
		selectedButton := radioButton.GetSelectedRadioButton()
		Locate(0, 0)
		Print("                                ")
		Locate(0, 0)
		Print(selectedButton)
		key := string(Inkey())
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
func testTextboxes() {
	commonResource.isDebugEnabled = false
	InitializeTerminal(80, 40)
	layer1 := AddLayer(0, 0, 80, 40, 1, nil)
	layer2 := AddLayer(25, 15, 40, 20, 1, &layer1)
	Layer(layer1)
	ColorRGB(255, 0, 0, 0, 0, 0)
	layer1.FillLayer("#")
	layer2.FillLayer("@")
	Locate(0, 0)
	Print("Enable ☑ Enable ○ ● (U+25CB, U+25CF)")
	styleEntry := types.NewTuiStyleEntry()
	styleEntry.SelectorTextAlignment = 2
	textBox := layer1.AddTextbox(styleEntry, 2, 2, 20, 5, true)
	textBox.setText("This is a test 1\nthis is ☑ second line which is very long and big 1\n李克强宣布中国今年经济增长的目标为 third line. 1")

	textBox2 := layer1.AddTextbox(styleEntry, 40, 2, 20, 5, false)
	textBox2.setText("This is a test\nthis is ☑ second line which is very long and big\nthis is 文字 third line.")
	textBox3 := layer1.AddTextbox(styleEntry, 1, 10, 20, 7, true)
	textBox3.setText("This is a test123456\nThis DDDDtesdfsfsdfsdfsdfsdfsdfsddffdsfdsst123456\nThis is GGGGGst1234\nThis is a ZZZZst123456\nThis is a test123456\nzzzzzzzzz\naaaa\ndddddd\nttttt\n222222\n555555")

	for {
		mouseXLocation, mouseYLocation, _, _ := GetMouseStatus()
		characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
		Locate(0, 0)
		Print(fmt.Sprintf("%d, %d   ", characterEntry.AttributeEntry.CellControlId, characterEntry.AttributeEntry.CellControlLocation))

		UpdateDisplay(false)
		key := string(Inkey())
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

func testCheckboxes() {
	commonResource.isDebugEnabled = false
	InitializeTerminal(80, 40)
	layer1 := AddLayer(0, 0, 80, 40, 1, nil)
	layer2 := AddLayer(20, 15, 40, 20, 1, &layer1)
	Layer(layer1)
	ColorRGB(255, 0, 0, 0, 0, 0)
	layer1.FillLayer("#")
	layer2.FillLayer("@")
	Locate(0, 0)

	Print("Enable ☑ Enable ○ ● (U+25CB, U+25CF) ▾☒♪")
	styleEntry := types.NewTuiStyleEntry()
	styleEntry.SelectorTextAlignment = 0
	layer1.AddCheckbox("Enable 文字 Feature", styleEntry, 2, 2, true, true)
	for {
		UpdateDisplay(false)
		key := string(Inkey())
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
	InitializeTerminal(80, 40)
	layer1 := AddLayer(0, 0, 80, 40, 1, nil)
	layer2 := AddLayer(20, 15, 40, 20, 1, &layer1)
	Layer(layer1)
	ColorRGB(255, 0, 0, 0, 0, 0)
	layer1.FillLayer("#")
	layer2.FillLayer("@")
	Locate(0, 0)
	Print("Enable ☑ Enable ○ ● (U+25CB, U+25CF)")
	styleEntry := types.NewTuiStyleEntry()
	styleEntry.SelectorTextAlignment = 0
	selectionEntry := types.NewSelectionEntry()
	selectionEntry.Add("1", "")
	selectionEntry.Add("2", "Enabl文e ○ ●")
	selectionEntry.Add("3", "☑ Enable ○ ●")
	selectionEntry.Add("4", "GET 文字文字")
	selectionEntry.Add("5", "IE文字文字")
	selectionEntry.Add("6", "DELETE")
	for i := 0; i < 20; i++ {
		selectionEntry.Add(strconv.Itoa(i), strconv.Itoa(i))
	}
	selectionEntry2 := types.NewSelectionEntry()
	selectionEntry2.Add("1", "1")
	selectionEntry2.Add("2", "2")
	selectionEntry2.Add("3", "3")
	layer1.AddDropdown(styleEntry, selectionEntry, 2, 2, 8, 7, 3)
	layer1.AddDropdown(styleEntry, selectionEntry2, 2, 39, 3, 7, 1)

	layer1.AddSelector(styleEntry, selectionEntry, 6, 10, 4, 7, 3, 0, -1, false)
	layer1.AddSelector(styleEntry, selectionEntry, 10, 20, 4, 7, 1, 0, -1, false)

	styleEntry.SelectorTextAlignment = constants.AlignmentNoPadding
	layer1.AddSelector(styleEntry, selectionEntry, 6, 30, 4, 7, 1, 0, -1, false)
	for {
		UpdateDisplay(false)
		key := string(Inkey())
		// fmt.Print(key)
		if key == "w" {
			Locate(0, 0)
			PrintLayer(layer1, "                                                ")
			Locate(0, 0)
			_, _, pressed, _ := GetPreviousMouseStatus()
			PrintLayer(layer1, "***"+strconv.Itoa(int(pressed))+"***")
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
	InitializeTerminal(80, 40)
	layer1 := AddLayer(0, 0, 80, 40, 1, nil)
	layer2 := AddLayer(20, 15, 40, 20, 1, &layer1)
	Layer(layer1)
	layer1.FillLayer("#")
	layer2.FillLayer("@")
	styleEntry := types.NewTuiStyleEntry()
	styleEntry.SelectorTextAlignment = 0
	selectionEntry := types.NewSelectionEntry()
	selectionEntry.Add("1", "OK")
	selectionEntry.Add("2", "CANCEL")
	selectionEntry.Add("3", "EXIT")
	selectionEntry.Add("4", "GOTO")
	selectionEntry.Add("5", "RUN")
	selectionEntry.Add("6", "DELETE")
	s1 := layer1.AddScrollbar(styleEntry, 2, 2, 8, 80, 0, 1, false)
	s2 := layer1.AddScrollbar(styleEntry, 10, 5, 8, 8, 4, 1, true)
	s1.setScrollValue(4)
	s2.setHandlePosition(4)
	for {
		UpdateDisplay(false)
		// x, y, _, _ := memory.GetMouseStatus()
		// a := getCellInformationUnderMouseCursor(x, y)
		LocateLayer(layer1, 10, 0)
		PrintLayer(layer1, "  ")
		LocateLayer(layer1, 10, 0)
		PrintLayer(layer1, strconv.Itoa(s1.getScrollValue()))
		LocateLayer(layer1, 10, 1)
		PrintLayer(layer1, "  ")
		LocateLayer(layer1, 10, 1)
		PrintLayer(layer1, strconv.Itoa(s2.getScrollValue()))
		key := string(Inkey())
		// fmt.Print(key)
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
	InitializeTerminal(80, 40)
	layer1 := AddLayer(0, 0, 80, 40, 1, nil)
	layer2 := AddLayer(20, 15, 40, 20, 1, &layer1)
	Layer(layer1)
	layer1.FillLayer("#")
	layer2.FillLayer("@")
	styleEntry := types.NewTuiStyleEntry()
	styleEntry.SelectorTextAlignment = 0
	selectionEntry := types.NewSelectionEntry()
	selectionEntry.Add("1", "OK")
	selectionEntry.Add("2", "CANCEL")
	selectionEntry.Add("3", "EXIT")
	selectionEntry.Add("4", "GOTO")
	selectionEntry.Add("5", "RUN")
	selectionEntry.Add("6", "DELETE")
	for i := 0; i < 20; i++ {
		selectionEntry.Add(strconv.Itoa(i), strconv.Itoa(i))
	}
	menuBarInstance := layer1.AddSelector(styleEntry, selectionEntry, 2, 10, 4, 7, 3, 0, -1, false)
	// menuBarInstance2 := Selector.AddLayer(layerAlias1, "menuBar2", styleEntry, selectionEntry, 10, 20, 4, 7, 1, 0, -1, false)
	LocateLayer(layer1, 3, 3)
	PrintLayer(layer1, menuBarInstance.layerAlias)
	for {
		UpdateDisplay(false)
		time.Sleep(50 * time.Millisecond)
		LocateLayer(layer1, 3, 3)
		// PrintLayer(layerAlias1, menuBarInstance.GetSelected())
		LocateLayer(layer1, 3, 4)
		// PrintLayer(layerAlias1, menuBarInstance2.GetSelected())
		key := string(Inkey())
		// fmt.Print(key)
		if key == "q" {
			break
		}
	}
	DeleteAllLayers()
	RestoreTerminalSettings()
}

func testTextField() {
	commonResource.isDebugEnabled = false
	InitializeTerminal(80, 40)
	layer1 := AddLayer(0, 0, 80, 40, 1, nil)
	layer2 := AddLayer(20, 15, 40, 20, 1, &layer1)
	Layer(layer1)
	layer1.FillLayer("#")
	layer2.FillLayer("@")
	styleEntry := types.NewTuiStyleEntry()
	textFieldInstance := layer1.AddTextField(styleEntry, 0, 3, 10, 60, true, "Alex Chang is the man", true)
	// Drag a scroll bar, move out of the scroll bar, moving on another control should not crash.
	layer1.AddTextField(styleEntry, 0, 7, 50, 60, false, "यह व्यवस्था की परीक्षा है", true)
	layer1.AddTextField(styleEntry, 0, 9, 50, 60, false, "Zhè shì duì wǒ de yìngyòng chéngxù de cèshì", true)

	LocateLayer(layer1, 3, 20)
	PrintLayer(layer1, textFieldInstance.GetValue())
	layer1.AddTextField(styleEntry, 0, 5, 15, 30, false, "Test 李克强宣布中国今年经济增长的目 acbc1", true)
	layer1.AddScrollbar(styleEntry, 0, 8, 10, 10, 0, 1, false)
	for {
		mouseXLocation, mouseYLocation, _, _ := GetMouseStatus()
		characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
		Locate(0, 0)
		Print(fmt.Sprintf("%d, %d   ", characterEntry.AttributeEntry.CellControlId, characterEntry.AttributeEntry.CellControlLocation))
		UpdateDisplay(false)
		time.Sleep(50 * time.Millisecond)
		key := string(Inkey())
		if key == "q" {
			break
		}
	}
	DeleteAllLayers()
	RestoreTerminalSettings()
}

func testWindowMovement() {
	commonResource.isDebugEnabled = false
	xLocation := 1
	yLocation := 1
	InitializeTerminal(80, 40)
	layer1 := AddLayer(0, 0, 80, 40, 1, nil)
	layer2 := AddLayer(20, 15, 40, 20, 1, &layer1)
	layer3 := AddLayer(0, 0, 20, 10, 1, &layer2)
	layer4 := AddLayer(15, 10, 20, 10, 2, &layer1)
	Layer(layer1)
	layer1.FillLayer("#")
	layer2.FillLayer("@")
	Locate(xLocation, yLocation)
	Print("This is a test of the layer system")
	Locate(1, 1)
	Print(strconv.Itoa(GetLayer(layer2.layerAlias).ZOrder))
	styleEntry := types.NewTuiStyleEntry()
	layer3.DrawWindow(styleEntry, 0, 0, 18, 9, true)
	layer4.DrawWindow(styleEntry, 0, 0, 18, 9, true)
	for i := 0; i < 10; i++ {
		UpdateDisplay(false)
		time.Sleep(1000 * time.Millisecond)
	}
	DeleteAllLayers()
	RestoreTerminalSettings()
}

func testButtonPressAction() {
	commonResource.isDebugEnabled = false
	xLocation := 0
	yLocation := 0
	InitializeTerminal(80, 40)
	layer1 := AddLayer(0, 0, 80, 40, 1, nil)
	Locate(xLocation, yLocation)
	Print("This is a test")
	styleEntry := types.NewTuiStyleEntry()
	layer1.AddButton("CANCEL", styleEntry, 2, 2, 10, 10, true)
	layer1.AddButton("OK", styleEntry, 15, 2, 10, 10, true)
	layerInformation := GetLayer(layer1.layerAlias)
	Button.drawButtonsOnLayer(*layerInformation)
	for {
		mouseXLocation, mouseYLocation, _, _ := GetMouseStatus()
		characterEntry := getCellInformationUnderMouseCursor(mouseXLocation, mouseYLocation)
		Locate(0, 0)
		Print(fmt.Sprintf("%d, %d, %d   ", characterEntry.AttributeEntry.CellType, characterEntry.AttributeEntry.CellControlId, characterEntry.AttributeEntry.CellControlLocation))
		UpdateDisplay(false)
		key := string(Inkey())
		if key == "q" {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	DeleteAllLayers()
	RestoreTerminalSettings()
}
