package constants

import (
	"github.com/gdamore/tcell/v2"
)

type ImageStyle int

const (
	OS_WINDOWS = "windows"
	OS_LINUX   = "linux"
	OS_MAC     = "darwin"
	OS_OTHER   = "other"
)

const (
	ImageStyleHighColor ImageStyle = iota
	ImageStyleBraille
	ImageStyleCharacters
	ImageStyleBlockElements
)

type DitheringStyle int

const (
	DitheringStyleBasic = iota
	DitheringStyle2x2BayerMatrix
	DitheringStyle4x4BayerMatrix
	DitheringStyle8x8BayerMatrix
	DitheringStyleErrorDiffusion
)

type EffectStyle int

const (
	EffectNone EffectStyle = iota
	EffectSinWave
	EffectConcentricCircles
	EffectFlagWave
	EffectBlinds
	EffectHorizontalWeaveTransition
	EffectVerticalWeaveTransition
	EffectForwardDiagonalWeaveTransition
	EffectBackwardDiagonalWeaveTransition
	EffectCounterClockwiseSwirlTransition
	EffectClockwiseSwirlTransition
	EffectGrowingCircleTransition
	EffectVerticalCurtainTransition
	EffectHorizontalCurtainTransition
)

type ColorType uint64

var CP437ToUnicode = [256]rune{
	'\u0000', '\u263A', '\u263B', '\u2665', '\u2666', '\u2663', '\u2660', '\u2022',
	'\u25D8', '\u25CB', '\u25D9', '\u2642', '\u2640', '\u266A', '\u266B', '\u263C',
	'\u25BA', '\u25C4', '\u2195', '\u203C', '\u00B6', '\u00A7', '\u25AC', '\u21A8',
	'\u2191', '\u2193', '\u2192', '\u2190', '\u221F', '\u2194', '\u25B2', '\u25BC',
	'\u0020', '\u0021', '\u0022', '\u0023', '\u0024', '\u0025', '\u0026', '\u0027',
	'\u0028', '\u0029', '\u002A', '\u002B', '\u002C', '\u002D', '\u002E', '\u002F',
	'\u0030', '\u0031', '\u0032', '\u0033', '\u0034', '\u0035', '\u0036', '\u0037',
	'\u0038', '\u0039', '\u003A', '\u003B', '\u003C', '\u003D', '\u003E', '\u003F',
	'\u0040', '\u0041', '\u0042', '\u0043', '\u0044', '\u0045', '\u0046', '\u0047',
	'\u0048', '\u0049', '\u004A', '\u004B', '\u004C', '\u004D', '\u004E', '\u004F',
	'\u0050', '\u0051', '\u0052', '\u0053', '\u0054', '\u0055', '\u0056', '\u0057',
	'\u0058', '\u0059', '\u005A', '\u005B', '\u005C', '\u005D', '\u005E', '\u005F',
	'\u0060', '\u0061', '\u0062', '\u0063', '\u0064', '\u0065', '\u0066', '\u0067',
	'\u0068', '\u0069', '\u006A', '\u006B', '\u006C', '\u006D', '\u006E', '\u006F',
	'\u0070', '\u0071', '\u0072', '\u0073', '\u0074', '\u0075', '\u0076', '\u0077',
	'\u0078', '\u0079', '\u007A', '\u007B', '\u007C', '\u007D', '\u007E', '\u2302',
	'\u00C7', '\u00FC', '\u00E9', '\u00E2', '\u00E4', '\u00E0', '\u00E5', '\u00E7',
	'\u00EA', '\u00EB', '\u00E8', '\u00EF', '\u00EE', '\u00EC', '\u00C4', '\u00C5',
	'\u00C9', '\u00E6', '\u00C6', '\u00F4', '\u00F6', '\u00F2', '\u00FB', '\u00F9',
	'\u00FF', '\u00D6', '\u00DC', '\u00A2', '\u00A3', '\u00A5', '\u20A7', '\u0192',
	'\u00E1', '\u00ED', '\u00F3', '\u00FA', '\u00F1', '\u00D1', '\u00AA', '\u00BA',
	'\u00BF', '\u2310', '\u00AC', '\u00BD', '\u00BC', '\u00A1', '\u00AB', '\u00BB',
	'\u2591', '\u2592', '\u2593', '\u2502', '\u2524', '\u2561', '\u2562', '\u2556',
	'\u2555', '\u2563', '\u2551', '\u2557', '\u255D', '\u255C', '\u255B', '\u2510',
	'\u2514', '\u2534', '\u252C', '\u251C', '\u2500', '\u253C', '\u255E', '\u255F',
	'\u255A', '\u2554', '\u2569', '\u2566', '\u2560', '\u2550', '\u256C', '\u2567',
	'\u2568', '\u2564', '\u2565', '\u2559', '\u2558', '\u2552', '\u2553', '\u256B',
	'\u256A', '\u2518', '\u250C', '\u2588', '\u2584', '\u258C', '\u2590', '\u2580',
	'\u03B1', '\u00DF', '\u0393', '\u03C0', '\u03A3', '\u03C3', '\u00B5', '\u03C4',
	'\u03A6', '\u0398', '\u03A9', '\u03B4', '\u221E', '\u03C6', '\u03B5', '\u2229',
	'\u2261', '\u00B1', '\u2265', '\u2264', '\u2320', '\u2321', '\u00F7', '\u2248',
	'\u00B0', '\u2219', '\u00B7', '\u221A', '\u207F', '\u00B2', '\u25A0', '\u00A0',
}

const (
	CharDot                                              = '\u2022' // •
	CharArrowLeft                                        = '\u2190' // ←
	CharArrowUp                                          = '\u2191' // ↑
	CharArrowDown                                        = '\u2193' // ↓
	CharHLine                                            = '\u2500' // ─
	CharSingleLineHorizontal                             = '\u2500' // ─
	CharSingleLineVertical                               = '\u2502' // │
	CharVLine                                            = '\u2502' // │
	CharSingleLineUpLeftCorner                           = '\u250C' // ┌
	CharULCorner                                         = '\u250C' // ┌
	CharSingleLineUpRightCorner                          = '\u2510' // ┐
	CharURCorner                                         = '\u2510' // ┐
	CharSingleLineTRight                                 = '\u251C' // ├
	CharSingleLineTLeft                                  = '\u2524' // ┤
	CharSingleLineTDown                                  = '\u252C' // ┬
	CharSingleLineTUp                                    = '\u2534' // ┴
	CharSingleLineCross                                  = '\u253C' // ┼
	CharLLCorner                                         = '\u2514' // └
	CharSingleLineLowerLeftCorner                        = '\u2514' // └
	CharLRCorner                                         = '\u2518' // ┘
	CharSingleLineLowerRightCorner                       = '\u2518' // ┘
	CharDoubleLineHorizontal                             = '\u2550' // ═
	CharDoubleLineVertical                               = '\u2551' // ║
	CharDoubleLineUpLeftCorner                           = '\u2554' // ╔
	CharDoubleLineUpRightCorner                          = '\u2557' // ╗
	CharDoubleLineLowerLeftCorner                        = '\u255A' // ╚
	CharDoubleLineLowerRightCorner                       = '\u255D' // ╝
	CharSingleLineDoubleRight                            = '\u255E' // ╞
	CharDoubleLineTSingleRight                           = '\u255F' // ╟
	CharDoubleLineTRight                                 = '\u2560' // ╠
	CharSingleLineDoubleLeft                             = '\u2561' // ╡
	CharDoubleLineTSingleLeft                            = '\u2562' // ╢
	CharDoubleLineTLeft                                  = '\u2563' // ╣
	CharDoubleLineTSingleDown                            = '\u2564' // ╤
	CharSingleLineDoubleDown                             = '\u2565' // ╥
	CharDoubleLineTDown                                  = '\u2566' // ╦
	CharDoubleLineTSingleUp                              = '\u2567' // ╧
	CharDoubleLineTUp                                    = '\u2569' // ╩
	CharSingleLineDoubleUp                               = '\u256B' // ╪
	CharDoubleLineCross                                  = '\u256C' // ╬
	CharRoundedULCorner                                  = '\u256D' // ╭
	CharRoundedURCorner                                  = '\u256E' // ╮
	CharRoundedLRCorner                                  = '\u256F' // ╯
	CharRoundedLLCorner                                  = '\u2570' // ╰
	CharBlockUpperHalf                                   = '\u2580' // ▀
	CharBlockLowerOneEighth                              = '\u2581' //
	CharBlockLowerOneQuarter                             = '\u2582' // ▂
	CharBlockLowerThreeEighths                           = '\u2583' // ▃
	CharBlockLowerHalf                                   = '\u2584' // ▄
	CharBlockLowerFiveEighths                            = '\u2585' // ▅
	CharBlockLowerThreeQuarters                          = '\u2586' // ▆
	CharBlockLowerSevenEighths                           = '\u2587' // ▇
	CharBlockFull                                        = '\u2588' // █
	CharBlockSolid                                       = '\u2588' // █
	CharBlockLeftSevenEighths                            = '\u2589' // ▉
	CharBlockLeftThreeQuarters                           = '\u258A' // ▊
	CharBlockLeftFiveEighths                             = '\u258B' // ▋
	CharBlockLeftHalf                                    = '\u258C' // ▌
	CharBlockLeftThreeEighths                            = '\u258D' // ▍
	CharBlockLeftOneQuarter                              = '\u258E' // ▎
	CharBlockLeftOneEighth                               = '\u258F' // ▏
	CharBlockRightHalf                                   = '\u2590' // ▐
	CharBlockLightShade                                  = '\u2591' // ░
	CharBlockSparce                                      = '\u2591' // ░
	CharBlockMedium                                      = '\u2592' // ▒
	CharBlockMediumShade                                 = '\u2592' // ▒
	CharBlockDarkShade                                   = '\u2593' // ▓
	CharBlockDense                                       = '\u2593' // ▓
	CharBlockUpperOneEighth                              = '\u2594' // ▔
	CharBlockRightOneEighth                              = '\u2595' // ▕
	CharBlockQuadrantLowerLeft                           = '\u2596' // ▖
	CharBlockQuadrantLowerRight                          = '\u2597' // ▗
	CharBlockQuadrantUpperLeft                           = '\u2598' // ▘
	CharBlockQuadrantUpperLeftAndLowerLeftAndLowerRight  = '\u2599' // ▙
	CharBlockQuadrantUpperLeftAndLowerRight              = '\u259A' // ▚
	CharBlockQuadrantUpperLeftAndUpperRightAndLowerLeft  = '\u259B' // ▛
	CharBlockQuadrantUpperLeftAndUpperRightAndLowerRight = '\u259C' // ▜
	CharBlockQuadrantUpperRight                          = '\u259D' // ▝
	CharBlockQuadrantUpperRightAndLowerLeft              = '\u259E' // ▞
	CharBlockQuadrantUpperRightAndLowerLeftAndLowerRight = '\u259F' // ▟
	CharTriangleUp                                       = '\u25B4' // ▴
	CharTriangleRight                                    = '\u25B8' // ▸
	CharTriangleDown                                     = '\u25BE' // ▾
	CharTriangleLeft                                     = '\u25C2' // ◂
	CharCheckedRadioButton                               = '\u25CB' // ○
	CharUncheckedRadioButton                             = '\u25C9' // ◉
	CharCheckedBox                                       = '\u2610' // ☐
	CharUncheckedBox                                     = '\u2611' // ☑
	CharFaceWhite                                        = '\u263A' // ☺
	CharFaceBlack                                        = '\u263B' // ☻
	CharSpades                                           = '\u2660' // ♠
	CharClub                                             = '\u2663' // ♣
	CharHeart                                            = '\u2665' // ♥
	CharDiamond                                          = '\u2666' // ♦
)

var CharBlockBitmasks = map[rune]uint64{
	CharBlockQuadrantUpperLeftAndLowerRight: 0b1111000011110000111100001111000000001111000011110000111100001111,
	CharBlockQuadrantUpperRight:             0b0000000000000000000000000000000011110000111100001111000011110000,
	CharBlockQuadrantUpperLeft:              0b0000000000000000000000000000000000001111000011110000111100001111,
	CharBlockQuadrantLowerRight:             0b1111000011110000111100001111000000000000000000000000000000000000,
	CharBlockQuadrantLowerLeft:              0b0000111100001111000011110000111100000000000000000000000000000000,
	CharBlockLeftOneEighth:                  0b0000000100000001000000010000000100000001000000010000000100000001,
	CharBlockLeftOneQuarter:                 0b0000001100000011000000110000001100000011000000110000001100000011,
	CharBlockLeftThreeEighths:               0b0000011100000111000001110000011100000111000001110000011100000111,
	CharBlockLeftHalf:                       0b0000111100001111000011110000111100001111000011110000111100001111,
	CharBlockLeftFiveEighths:                0b0001111100011111000111110001111100011111000111110001111100011111,
	CharBlockLeftThreeQuarters:              0b0011111100111111001111110011111100111111001111110011111100111111,
	CharBlockLeftSevenEighths:               0b0111111101111111011111110111111101111111011111110111111101111111,
	CharBlockLowerSevenEighths:              0b1111111111111111111111111111111111111111111111111111111100000000,
	CharBlockLowerThreeQuarters:             0b1111111111111111111111111111111111111111111111110000000000000000,
	CharBlockLowerFiveEighths:               0b1111111111111111111111111111111111111111000000000000000000000000,
	CharBlockLowerHalf:                      0b1111111111111111111111111111111100000000000000000000000000000000,
	CharBlockLowerThreeEighths:              0b1111111111111111111111110000000000000000000000000000000000000000,
	CharBlockLowerOneQuarter:                0b1111111111111111000000000000000000000000000000000000000000000000,
	CharBlockLowerOneEighth:                 0b1111111100000000000000000000000000000000000000000000000000000000,
}

// Black
const (
	ColorBlack = iota
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
	ColorBrightBlack
	ColorBrightRed
	ColorBrightGreen
	ColorBrightYellow
	ColorBrightBlue
	ColorBrightMagenta
	ColorBrightCyan
	ColorBrightWhite
)

var AnsiColorByIndex = map[int]ColorType{
	ColorBlack:         ColorType(tcell.NewRGBColor(0, 0, 0)),
	ColorRed:           ColorType(tcell.GetColor("maroon")),
	ColorGreen:         ColorType(tcell.GetColor("green")),
	ColorYellow:        ColorType(tcell.GetColor("olive")),
	ColorBlue:          ColorType(tcell.GetColor("navy")),
	ColorMagenta:       ColorType(tcell.GetColor("purple")),
	ColorCyan:          ColorType(tcell.GetColor("teal")),
	ColorWhite:         ColorType(tcell.GetColor("silver")),
	ColorBrightBlack:   ColorType(tcell.GetColor("gray")),
	ColorBrightRed:     ColorType(tcell.GetColor("red")),
	ColorBrightGreen:   ColorType(tcell.GetColor("lime")),
	ColorBrightYellow:  ColorType(tcell.GetColor("yellow")),
	ColorBrightBlue:    ColorType(tcell.GetColor("blue")),
	ColorBrightMagenta: ColorType(tcell.GetColor("fuchsia")),
	ColorBrightCyan:    ColorType(tcell.GetColor("aqua")),
	ColorBrightWhite:   ColorType(tcell.GetColor("white")),
}

// TdfToRgbMap maps the 16-color TDF palette to 24-bit RGB ColorType values.
var TdfToRgbMap = []ColorType{
	ColorType(tcell.NewRGBColor(0, 0, 0)),       // 0: Black
	ColorType(tcell.NewRGBColor(0, 0, 170)),     // 1: Blue
	ColorType(tcell.NewRGBColor(0, 170, 0)),     // 2: Green
	ColorType(tcell.NewRGBColor(0, 170, 170)),   // 3: Cyan
	ColorType(tcell.NewRGBColor(170, 0, 0)),     // 4: Red
	ColorType(tcell.NewRGBColor(170, 0, 170)),   // 5: Magenta
	ColorType(tcell.NewRGBColor(170, 85, 0)),    // 6: Brown
	ColorType(tcell.NewRGBColor(170, 170, 170)), // 7: Light Gray
	ColorType(tcell.NewRGBColor(85, 85, 85)),    // 8: Dark Gray
	ColorType(tcell.NewRGBColor(85, 85, 255)),   // 9: Light Blue
	ColorType(tcell.NewRGBColor(85, 255, 85)),   // 10: Light Green
	ColorType(tcell.NewRGBColor(85, 255, 255)),  // 11: Light Cyan
	ColorType(tcell.NewRGBColor(255, 85, 85)),   // 12: Light Red
	ColorType(tcell.NewRGBColor(255, 85, 255)),  // 13: Light Magenta
	ColorType(tcell.NewRGBColor(255, 255, 85)),  // 14: Yellow
	ColorType(tcell.NewRGBColor(255, 255, 255)), // 15: White
}

const (
	AnsiEsc = '\u001b'
)
const NullRune = '\x00'
const NullColor = -1
const NullDataType = -1
const NullTransformValue = -1
const NullSelectionIndex = -1
const NullItemSelection = -1
const NullCellId = -1
const NullCellType = -1
const NullCellControlLocation = -1
const NullCellControlId = -2
const TransformContrast = 0
const TransformTransparency = 0

const AlignmentLeft = 0
const AlignmentRight = 1
const AlignmentCenter = 2
const AlignmentNoPadding = 3
const FrameStyleNormal = 0
const FrameStyleRaised = 1
const FrameStyleSunken = 2

const CellTypeButton = 1
const CellTypeTextField = 2
const CellTypeFrameTop = 3
const CellTypeSelectorItem = 4
const CellTypeScrollbar = 5
const CellTypeDropdown = 6
const CellTypeCheckbox = 7
const CellTypeTextbox = 8
const CellTypeRadioButton = 9
const CellTypeProgressBar = 10
const CellTypeLabel = 11
const CellTypeTooltip = 12
const CellTypeFileMenuHeading = 13
const CellTypeFileMenuItem = 14
const CellTypeShadow = 15

const CellControlIdUpScrollArrow = -1
const CellControlIdDownScrollArrow = -2
const CellControlIdScrollbarHandle = -3
const CellControlIdUnchecked = 1
const CellControlIdChecked = 2

const VirtualFileSystemZip = 1
const VirtualFileSystemRar = 2
const VirtualFileSystemEmbedded = 3
const EventStateNone = 0
const EventStateDragAndDrop = 1
const EventStateDragAndDropScrollbar = 2
const NullControlType = 0

const NullScrollbarValue = -1

// Control type string constants
const TYPE_BUTTON = "button"
const TYPE_CHECKBOX = "checkbox"
const TYPE_DROPDOWN = "dropdown"
const TYPE_FONT = "font"
const TYPE_LABEL = "label"
const TYPE_PROGRESSBAR = "progressbar"
const TYPE_SCROLLBAR = "scrollbar"
const TYPE_SELECTOR = "selector"
const TYPE_TEXTBOX = "textbox"
const TYPE_TEXTFIELD = "textfield"
const TYPE_TOOLTIP = "tooltip"
const TYPE_RADIOBUTTON = "radiobutton"
const TYPE_VIEWPORT = "viewport"
const TYPE_FILEMENU = "filemenu"

const DefaultTooltipHoverTime = 1000
const SELECTED_NONE = -1

const (
	ButtonStateUnpressed = iota
	ButtonStatePressed
	ButtonStateHovering
)

const (
	_ = iota
	MouseButtonLeft
	MouseButtonMiddle
	MouseButtonRight
)
