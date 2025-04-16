package constants

import (
	"github.com/gdamore/tcell/v2"
)

type ImageStyle int

const (
	ImageStyleHighColor ImageStyle = iota
	ImageStyleBraille
	ImageStyleCharacters
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

// Constants
const (
	CharULCorner                   = '\u250C'
	CharURCorner                   = '\u2510'
	CharLLCorner                   = '\u2514'
	CharLRCorner                   = '\u2518'
	CharHLine                      = '\u2500'
	CharVLine                      = '\u2502'
	CharFaceWhite                  = '\u263A'
	CharFaceBlack                  = '\u263B'
	CharHeart                      = '\u2665'
	CharClub                       = '\u2663'
	CharDiamond                    = '\u2666'
	CharSpades                     = '\u2660'
	CharDot                        = '\u2022'
	CharTriangleUp                 = '\u25B4'
	CharTriangleDown               = '\u25BE'
	CharTriangleLeft               = '\u25C2'
	CharTriangleRight              = '\u25B8'
	CharArrowUp                    = '\u2191'
	CharArrowDown                  = '\u2193'
	CharArrowLeft                  = '\u2190'
	CharBlockSolid                 = '\u2588'
	CharBlockDense                 = '\u2593'
	CharBlockMedium                = '\u2592'
	CharBlockSparce                = '\u2591'
	CharSingleLineHorizontal       = CharHLine
	CharDoubleLineHorizontal       = '\u2550'
	CharSingleLineVertical         = '\u2502'
	CharDoubleLineVertical         = '\u2551'
	CharSingleLineUpLeftCorner     = CharULCorner
	CharDoubleLineUpLeftCorner     = '\u2554'
	CharSingleLineUpRightCorner    = CharURCorner
	CharDoubleLineUpRightCorner    = '\u2557'
	CharSingleLineLowerLeftCorner  = CharLLCorner
	CharDoubleLineLowerLeftCorner  = '\u255A'
	CharSingleLineLowerRightCorner = CharLRCorner
	CharDoubleLineLowerRightCorner = '\u255D'
	CharSingleLineCross            = '\u253C'
	CharDoubleLineCross            = '\u256C'
	CharSingleLineTUp              = '\u2534'
	CharSingleLineTDown            = '\u252C'
	CharSingleLineTLeft            = '\u2524'
	CharSingleLineTRight           = '\u251C'
	CharSingleLineDoubleUp         = '\u256B'
	CharSingleLineDoubleDown       = '\u2565'
	CharSingleLineDoubleLeft       = '\u2561'
	CharSingleLineDoubleRight      = '\u255E'
	CharDoubleLineTUp              = '\u2569'
	CharDoubleLineTDown            = '\u2566'
	CharDoubleLineTLeft            = '\u2563'
	CharDoubleLineTRight           = '\u2560'
	CharDoubleLineTSingleUp        = '\u2567'
	CharDoubleLineTSingleDown      = '\u2564'
	CharDoubleLineTSingleLeft      = '\u2562'
	CharDoubleLineTSingleRight     = '\u255F'
	CharBlockLowerHalf             = '\u2584'
	CharBlockUpperHalf             = '\u2580'
	CharCheckedBox                 = '\u2610'
	CharUncheckedBox               = '\u2611'
	CharCheckedRadioButton         = '\u25CB'
	CharUncheckedRadioButton       = '\u25C9'
)

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
