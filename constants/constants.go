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

const (
	CharDot                                 = '\u2022'
	CharArrowLeft                           = '\u2190'
	CharArrowUp                             = '\u2191'
	CharArrowDown                           = '\u2193'
	CharHLine                               = '\u2500'
	CharSingleLineHorizontal                = '\u2500'
	CharSingleLineVertical                  = '\u2502'
	CharVLine                               = '\u2502'
	CharSingleLineUpLeftCorner              = '\u250C'
	CharULCorner                            = '\u250C'
	CharSingleLineUpRightCorner             = '\u2510'
	CharURCorner                            = '\u2510'
	CharSingleLineTRight                    = '\u251C'
	CharSingleLineTLeft                     = '\u2524'
	CharSingleLineTDown                     = '\u252C'
	CharSingleLineTUp                       = '\u2534'
	CharSingleLineCross                     = '\u253C'
	CharLLCorner                            = '\u2514'
	CharSingleLineLowerLeftCorner           = '\u2514'
	CharLRCorner                            = '\u2518'
	CharSingleLineLowerRightCorner          = '\u2518'
	CharDoubleLineHorizontal                = '\u2550'
	CharDoubleLineVertical                  = '\u2551'
	CharDoubleLineUpLeftCorner              = '\u2554'
	CharDoubleLineUpRightCorner             = '\u2557'
	CharDoubleLineLowerLeftCorner           = '\u255A'
	CharDoubleLineLowerRightCorner          = '\u255D'
	CharSingleLineDoubleRight               = '\u255E'
	CharDoubleLineTSingleRight              = '\u255F'
	CharDoubleLineTRight                    = '\u2560'
	CharSingleLineDoubleLeft                = '\u2561'
	CharDoubleLineTSingleLeft               = '\u2562'
	CharDoubleLineTLeft                     = '\u2563'
	CharDoubleLineTSingleDown               = '\u2564'
	CharSingleLineDoubleDown                = '\u2565'
	CharDoubleLineTDown                     = '\u2566'
	CharDoubleLineTSingleUp                 = '\u2567'
	CharDoubleLineTUp                       = '\u2569'
	CharSingleLineDoubleUp                  = '\u256B'
	CharDoubleLineCross                     = '\u256C'
	CharRoundedULCorner                     = '\u256D'
	CharRoundedURCorner                     = '\u256E'
	CharRoundedLRCorner                     = '\u256F'
	CharRoundedLLCorner                     = '\u2570'
	CharBlockUpperHalf                      = '\u2580'
	CharBlockLowerOneEighth                 = '\u2581'
	CharBlockLowerOneQuarter                = '\u2582'
	CharBlockLowerThreeEighths              = '\u2583'
	CharBlockLowerHalf                      = '\u2584'
	CharBlockLowerFiveEighths               = '\u2585'
	CharBlockLowerThreeQuarters             = '\u2586'
	CharBlockLowerSevenEighths              = '\u2587'
	CharBlockFull                           = '\u2588'
	CharBlockSolid                          = '\u2588'
	CharBlockLeftSevenEighths               = '\u2589'
	CharBlockLeftThreeQuarters              = '\u258A'
	CharBlockLeftFiveEighths                = '\u258B'
	CharBlockLeftHalf                       = '\u258C'
	CharBlockLeftThreeEighths               = '\u258D'
	CharBlockLeftOneQuarter                 = '\u258E'
	CharBlockLeftOneEighth                  = '\u258F'
	CharBlockLightShade                     = '\u2591'
	CharBlockSparce                         = '\u2591'
	CharBlockMedium                         = '\u2592'
	CharBlockMediumShade                    = '\u2592'
	CharBlockDarkShade                      = '\u2593'
	CharBlockDense                          = '\u2593'
	CharBlockQuadrantLowerLeft              = '\u2596'
	CharBlockQuadrantLowerRight             = '\u2597'
	CharBlockQuadrantUpperLeft              = '\u2598'
	CharBlockQuadrantUpperRight             = '\u2599'
	CharBlockQuadrantUpperLeftAndLowerRight = '\u259A'
	CharTriangleUp                          = '\u25B4'
	CharTriangleRight                       = '\u25B8'
	CharTriangleDown                        = '\u25BE'
	CharTriangleLeft                        = '\u25C2'
	CharCheckedRadioButton                  = '\u25CB'
	CharUncheckedRadioButton                = '\u25C9'
	CharCheckedBox                          = '\u2610'
	CharUncheckedBox                        = '\u2611'
	CharFaceWhite                           = '\u263A'
	CharFaceBlack                           = '\u263B'
	CharSpades                              = '\u2660'
	CharClub                                = '\u2663'
	CharHeart                               = '\u2665'
	CharDiamond                             = '\u2666'
)

var BlockElementRunes = []rune{
	CharBlockLowerOneEighth,
	CharBlockLowerOneQuarter,
	CharBlockLowerThreeEighths,
	CharBlockLowerHalf,
	CharBlockLowerFiveEighths,
	CharBlockLowerThreeQuarters,
	CharBlockLowerSevenEighths,
	CharBlockLeftSevenEighths,
	CharBlockLeftThreeQuarters,
	CharBlockLeftFiveEighths,
	CharBlockLeftHalf,
	CharBlockLeftThreeEighths,
	CharBlockLeftOneQuarter,
	CharBlockLeftOneEighth,
	CharBlockQuadrantLowerLeft,
	CharBlockQuadrantLowerRight,
	CharBlockQuadrantUpperLeft,
	CharBlockQuadrantUpperRight,
	CharBlockQuadrantUpperLeftAndLowerRight,
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

const CellControlIdUpScrollArrow = -1
const CellControlIdDownScrollArrow = -2
const CellControlIdScrollbarHandle = -3
const CellControlIdUnchecked = 1
const CellControlIdChecked = 2

const VirtualFileSystemZip = 1
const VirtualFileSystemRar = 2
const EventStateNone = 0
const EventStateDragAndDrop = 1
const EventStateDragAndDropScrollbar = 2
const NullControlType = 0

const NullScrollbarValue = -1

const DefaultTooltipHoverTime = 1000
