// Command cjkdemo renders every interactive consolizer widget populated with Chinese and other wide (double-cell)
// characters so a human can verify alignment, truncation, the caret, hit testing, scrolling and word wrap in a
// real UTF-8 terminal.
//
// Run it from the repository root:
//
//	go run ./examples/cjkdemo
//
// Tab cycles focus through the buttons, checkboxes, radio buttons, list and dropdown; click a text field or a
// textbox to put the caret in it. Click, type, drag-select, use the arrow keys, and drag the scrollbars to
// exercise every widget. One textbox has word wrap on (no horizontal scroll) and one has it off (both
// scrollbars). Esc quits and restores the terminal. Your terminal must be at least 84 columns by 40 rows.
package main

import (
	"fmt"
	"strings"

	"github.com/supercom32/consolizer"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/stringformat"
)

const (
	minWidth  = 84
	minHeight = 40
)

func main() {
	consolizer.InitializeTerminal(0, 0)
	terminalWidth, terminalHeight := consolizer.GetTerminalSize()
	if terminalWidth < minWidth || terminalHeight < minHeight {
		consolizer.RestoreTerminalSettings()
		fmt.Printf("Please resize the terminal to at least %dx%d (currently %dx%d) and run again.\n",
			minWidth, minHeight, terminalWidth, terminalHeight)
		return
	}

	style := consolizer.NewTuiStyleEntry()
	mainLayer := consolizer.AddLayer(0, 0, terminalWidth, terminalHeight, 1, nil)
	mainLayer.Color(15, 4)
	mainLayer.FillLayer(" ")

	mainLayer.AddLabel("  中文宽字符控件测试   Wide-Character Widget Test", style, 0, 0, terminalWidth)

	mainLayer.AddLabel("窄字符不可加宽 (must stay one column each):", style, 2, 2, 42)
	mainLayer.AddLabel("café Ñoño Ελλάδα Привет  ABCabc 123", style, 44, 2, 36)

	mainLayer.AddLabel("截断 (truncation):", style, 2, 4, 20)
	mainLayer.AddLabel("这是一个很长的中文标签会被截断显示", style, 22, 4, 20)

	flatButtonStyle := style
	flatButtonStyle.Button.StyleMode = constants.ButtonStyleFlat
	borderlessButtonStyle := style
	borderlessButtonStyle.Button.StyleMode = constants.ButtonStyleBorderless
	borderlessButtonStyle.Button.BackgroundColor = consolizer.GetColor(4)
	borderlessButtonStyle.Button.ForegroundColor = consolizer.GetColor(15)
	borderlessButtonStyle.Button.PressedForegroundColor = consolizer.GetColor(0)
	borderlessButtonStyle.Button.PressedBackgroundColor = consolizer.GetColor(2)

	mainLayer.AddLabel("按钮样式 Button styles — beveled / flat / borderless:", style, 2, 5, terminalWidth-4)
	okButton := mainLayer.AddButton("确定 OK", style, 2, 6, 8, 3, true)
	cancelButton := mainLayer.AddButton("取消 Cancel", style, 11, 6, 13, 3, true)
	helpButton := mainLayer.AddButton("帮助", style, 25, 6, 6, 3, true)
	flatButton := mainLayer.AddButton("平面 Flat", flatButtonStyle, 32, 6, 11, 3, true)
	borderlessButton := mainLayer.AddButton("无框按钮 Borderless", borderlessButtonStyle, 2, 9, 20, 1, true)
	tinyBorderlessButton := mainLayer.AddButton("BTN", borderlessButtonStyle, 24, 9, 3, 1, true)
	okButton.AddToTabIndex()
	cancelButton.AddToTabIndex()
	helpButton.AddToTabIndex()
	flatButton.AddToTabIndex()
	borderlessButton.AddToTabIndex()
	tinyBorderlessButton.AddToTabIndex()

	mainLayer.AddLabel("城市 City:", style, 44, 7, 10)
	cityEntry := consolizer.NewSelectionEntry()
	cityEntry.Add("beijing", "北京")
	cityEntry.Add("shanghai", "上海")
	cityEntry.Add("guangzhou", "广州")
	cityEntry.Add("shenzhen", "深圳")
	cityEntry.Add("chengdu", "成都市")
	cityDropdown := mainLayer.AddDropdown(style, cityEntry, 55, 7, 5, 12, 0)
	cityDropdown.AddToTabIndex()

	optionCheckbox := mainLayer.AddCheckbox("启用中文选项", style, 2, 10, true, true)
	nightCheckbox := mainLayer.AddCheckbox("夜间模式", style, 24, 10, false, true)
	gridCheckbox := mainLayer.AddCheckbox("显示网格线", style, 42, 10, false, true)
	optionCheckbox.AddToTabIndex()
	nightCheckbox.AddToTabIndex()
	gridCheckbox.AddToTabIndex()

	firstRadio := mainLayer.AddRadioButton("选项一", style, 2, 12, 1, true)
	secondRadio := mainLayer.AddRadioButton("选项二内容较长", style, 14, 12, 1, false)
	thirdRadio := mainLayer.AddRadioButton("选项三", style, 36, 12, 1, false)
	firstRadio.AddToTabIndex()
	secondRadio.AddToTabIndex()
	thirdRadio.AddToTabIndex()

	mainLayer.AddLabel("列表 List:", style, 2, 14, 18)
	listEntry := consolizer.NewSelectionEntry()
	listEntry.Add("row1", "第一项")
	listEntry.Add("row2", "第二项内容较长")
	listEntry.Add("row3", "第三项目")
	listEntry.Add("row4", "第四选项")
	listEntry.Add("row5", "第五项")
	listSelector := mainLayer.AddSelector(style, listEntry, 2, 15, 5, 18, 1, 0, 0, false, true)
	listSelector.AddToTabIndex()

	mainLayer.AddLabel("换行文本框 Textbox, word wrap ON:", style, 28, 13, 44)
	wrappedTextbox := mainLayer.AddTextbox(style, 28, 14, 44, 4, true)
	wrappedTextbox.SetWordWrap(true)
	wrappedTextbox.SetText("中文段落自动换行测试 这一段文字用来检查在指定宽度下 中文和 English words 混排时是否在正确的列位置换行 并且不会把一个全角字符切成两半 用方向键上下移动光标 应当逐个可见行前进 不会卡住")

	mainLayer.AddLabel("不换行文本框 Textbox, word wrap OFF (both scrollbars):", style, 28, 20, 52)
	scrollTextbox := mainLayer.AddTextbox(style, 28, 22, 44, 4, true)
	scrollTextbox.SetText("这一行不会换行 this line does NOT wrap so drag the horizontal scrollbar 结束标记\n" +
		"第二行同样很长 second long line for horizontal scrolling as well 需要水平滚动才能看完\n" +
		"第三行 third line\n第四行 fourth line 内容\n第五行 fifth line\n第六行 sixth line to force the vertical scrollbar 结束")

	mainLayer.AddLabel("文本框 Field:", style, 2, 28, 12)
	textField := mainLayer.AddTextField(style, 14, 28, 26, 80, false, "中文abc测试内容", true)

	mainLayer.AddLabel("密码 Password:", style, 2, 29, 12)
	passwordField := mainLayer.AddTextField(style, 14, 29, 26, 80, true, "中文pw密码", true)

	mainLayer.AddLabel("进度条 Progress bars — top: label bg follows the fill   bottom: fixed label bg", style, 2, 30, terminalWidth-4)
	progressBarTransparent := mainLayer.AddProgressBar("下载进度", style, 2, 31, 52, 1, false, 0, 100, true)
	progressBarOpaque := mainLayer.AddProgressBar("下载进度", style, 2, 32, 52, 1, false, 0, 100, false)

	mainLayer.AddLabel("[Tab] cycles widgets   click / type / arrow-keys / drag scrollbars in a textbox   [Esc] quit",
		style, 2, 38, terminalWidth-4)

	consolizer.UpdateDisplay(true)

	lastButtonAction := "(none)"
	tick := 0
	for {
		for {
			keystroke := consolizer.Inkey()
			if len(keystroke) == 0 {
				break
			}
			name := string(keystroke)
			if name == "esc" || name == "escape" {
				consolizer.RestoreTerminalSettings()
				return
			}
		}

		tick++
		percent := (tick / 2) % 101
		progressLabel := fmt.Sprintf("下载进度 %d%%", percent)
		progressBarTransparent.SetValue(percent)
		progressBarTransparent.SetLabel(progressLabel)
		progressBarOpaque.SetValue(percent)
		progressBarOpaque.SetLabel(progressLabel)

		if _, pressedAlias := okButton.GetPressed(); pressedAlias != "" {
			switch pressedAlias {
			case okButton.GetAlias():
				lastButtonAction = "确定 OK (beveled)"
			case cancelButton.GetAlias():
				lastButtonAction = "取消 Cancel (beveled)"
			case helpButton.GetAlias():
				lastButtonAction = "帮助 Help (beveled)"
			case flatButton.GetAlias():
				lastButtonAction = "平面 Flat"
			case borderlessButton.GetAlias():
				lastButtonAction = "无框按钮 Borderless"
			case tinyBorderlessButton.GetAlias():
				lastButtonAction = "BTN (borderless, text = button)"
			}
		}

		selectedRadio := "选项一"
		if secondRadio.IsSelected() {
			selectedRadio = "选项二内容较长"
		} else if thirdRadio.IsSelected() {
			selectedRadio = "选项三"
		}

		listValue := "(none)"
		if _, listIndex := listSelector.GetSelected(); listIndex >= 0 {
			if _, values := listSelector.GetAllItems(); listIndex < len(values) {
				listValue = values[listIndex]
			}
		}

		drawStatusLine(mainLayer, terminalWidth, 34,
			fmt.Sprintf("文本框内容 Field: [%s]   密码真实值 Password: [%s]", textField.GetValue(), passwordField.GetValue()))
		drawStatusLine(mainLayer, terminalWidth, 35,
			fmt.Sprintf("下拉 Dropdown: %s   列表 List: %s   单选 Radio: %s", cityDropdown.GetValue(), listValue, selectedRadio))
		drawStatusLine(mainLayer, terminalWidth, 36,
			fmt.Sprintf("复选 Checkbox: 中文选项=%s 夜间=%s 网格=%s   最近按钮 Button: %s",
				checkMark(optionCheckbox.IsSelected()), checkMark(nightCheckbox.IsSelected()),
				checkMark(gridCheckbox.IsSelected()), lastButtonAction))

		consolizer.UpdateDisplay(false)
		consolizer.Sleep(33)
	}
}

func checkMark(isChecked bool) string {
	if isChecked {
		return "是"
	}
	return "否"
}

func drawStatusLine(layer *consolizer.LayerInstanceType, terminalWidth int, rowIndex int, text string) {
	availableColumns := terminalWidth - 4
	printedWidth := stringformat.GetWidthOfRunesWhenPrinted([]rune(text))
	if printedWidth < availableColumns {
		text += strings.Repeat(" ", availableColumns-printedWidth)
	}
	layer.Locate(2, rowIndex)
	layer.Print(text)
}
