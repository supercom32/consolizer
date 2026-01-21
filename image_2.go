package consolizer

import (
	"github.com/supercom32/consolizer/types"
	"image"
	"math"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/uniseg"
)

// Types of dithering applied to images.
const (
	DitheringNone           = iota // No dithering.
	DitheringFloydSteinberg        // Floyd-Steinberg dithering (the default).
)

// The number of colors supported by true color terminals (R*G*B = 256*256*256).
const TrueColor = 16777216

// This map describes what each block element looks like. A 1 bit represents a
// pixel that is drawn, a 0 bit represents a pixel that is not drawn. The least
// significant bit is the top left pixel, the most significant bit is the bottom
// right pixel, moving row by row from left to right, top to bottom.
var blockElements = map[rune]uint64{
	BlockLowerOneEighthBlock:            0b1111111100000000000000000000000000000000000000000000000000000000,
	BlockLowerOneQuarterBlock:           0b1111111111111111000000000000000000000000000000000000000000000000,
	BlockLowerThreeEighthsBlock:         0b1111111111111111111111110000000000000000000000000000000000000000,
	BlockLowerHalfBlock:                 0b1111111111111111111111111111111100000000000000000000000000000000,
	BlockLowerFiveEighthsBlock:          0b1111111111111111111111111111111111111111000000000000000000000000,
	BlockLowerThreeQuartersBlock:        0b1111111111111111111111111111111111111111111111110000000000000000,
	BlockLowerSevenEighthsBlock:         0b1111111111111111111111111111111111111111111111111111111100000000,
	BlockLeftSevenEighthsBlock:          0b0111111101111111011111110111111101111111011111110111111101111111,
	BlockLeftThreeQuartersBlock:         0b0011111100111111001111110011111100111111001111110011111100111111,
	BlockLeftFiveEighthsBlock:           0b0001111100011111000111110001111100011111000111110001111100011111,
	BlockLeftHalfBlock:                  0b0000111100001111000011110000111100001111000011110000111100001111,
	BlockLeftThreeEighthsBlock:          0b0000011100000111000001110000011100000111000001110000011100000111,
	BlockLeftOneQuarterBlock:            0b0000001100000011000000110000001100000011000000110000001100000011,
	BlockLeftOneEighthBlock:             0b0000000100000001000000010000000100000001000000010000000100000001,
	BlockQuadrantLowerLeft:              0b0000111100001111000011110000111100000000000000000000000000000000,
	BlockQuadrantLowerRight:             0b1111000011110000111100001111000000000000000000000000000000000000,
	BlockQuadrantUpperLeft:              0b0000000000000000000000000000000000001111000011110000111100001111,
	BlockQuadrantUpperRight:             0b0000000000000000000000000000000011110000111100001111000011110000,
	BlockQuadrantUpperLeftAndLowerRight: 0b1111000011110000111100001111000000001111000011110000111100001111,
}

// pixel represents a character on screen used to draw part of an image.
type pixel struct {
	style   tcell.Style
	element rune // The block element.
}

// Image_2 implements a widget that displays one image. The original image
// (specified with [Image_2.SetImage]) is resized according to the specified size
// (see [Image_2.SetSize]), using the specified number of colors (see
// [Image_2.SetColors]), while applying dithering if necessary (see
// [Image_2.SetDithering]).
//
// Images are approximated by graphical characters in the terminal. The
// resolution is therefore limited by the number and type of characters that can
// be drawn in the terminal and the colors available in the terminal. The
// quality of the final image also depends on the terminal's font and spacing
// settings, none of which are under the control of this package. Results may
// vary.
type Image_2 struct {
	*Box

	// The image to be displayed. If nil, the widget will be empty.
	image image.Image

	// The size of the image. If a value is 0, the corresponding size is chosen
	// automatically based on the other size while preserving the image's aspect
	// ratio. If both are 0, the image uses as much space as possible. A
	// negative value represents a percentage, e.g. -50 means 50% of the
	// available space.
	width, height int

	// The number of colors to use. If 0, the number of colors is chosen based
	// on the terminal's capabilities.
	colors int

	// The dithering algorithm to use, one of the constants starting with
	// "ImageDithering".
	dithering int

	// The width of a terminal's cell divided by its height.
	aspectRatio float64

	// Horizontal and vertical alignment, one of the "Align" constants.
	alignHorizontal, alignVertical int

	// The text to be displayed before the image.
	label string

	// The label style.
	labelStyle tcell.Style

	// The screen width of the label area. A value of 0 means use the width of
	// the label text.
	labelWidth int

	// The actual image size (in cells) when it was drawn the last time.
	lastWidth, lastHeight int

	// The actual image (in cells) when it was drawn the last time. The size of
	// this slice is lastWidth * lastHeight, indexed by y*lastWidth + x.
	pixels []pixel

	// A callback function set by the Form class and called when the user leaves
	// this form item.
	finished func(tcell.Key)
}

// NewImage returns a new [Image_2] widget with an empty image (use
// [Image_2.SetImage] to specify the image to be displayed). The image will use
// the widget's entire available space. The default dithering algorithm is set
// to Floyd-Steinberg dithering. The terminal's cell aspect ratio defaults to
// 0.5.
func NewImage() *Image_2 {
	i := &Image_2{
		Box:             NewBox(),
		dithering:       DitheringFloydSteinberg,
		aspectRatio:     0.5,
		alignHorizontal: AlignCenter,
		alignVertical:   AlignCenter,
	}
	i.Box.Primitive = i
	return i
}

// SetImage sets the image to be displayed. If nil, the widget will be empty.
func (i *Image_2) SetImage(image image.Image) *Image_2 {
	i.image = image
	i.lastWidth, i.lastHeight = 0, 0
	return i
}

// SetSize sets the size of the image. Positive values refer to cells in the
// terminal. Negative values refer to a percentage of the available space (e.g.
// -50 means 50%). A value of 0 means that the corresponding size is chosen
// automatically based on the other size while preserving the image's aspect
// ratio. If both are 0, the image uses as much space as possible while still
// preserving the aspect ratio.
func (i *Image_2) SetSize(rows, columns int) *Image_2 {
	i.width = columns
	i.height = rows
	return i
}

// SetColors sets the number of colors to use. This should be the number of
// colors supported by the terminal. If 0, the number of colors is chosen based
// on the TERM environment variable (which may or may not be reliable).
//
// Only the values 0, 2, 8, 256, and 16777216 ([TrueColor]) are supported. Other
// values will be rounded up to the next supported value, to a maximum of
// 16777216.
//
// The effect of using more colors than supported by the terminal is undefined.
func (i *Image_2) SetColors(colors int) *Image_2 {
	i.colors = colors
	i.lastWidth, i.lastHeight = 0, 0
	return i
}

// GetColors returns the number of colors that will be used while drawing the
// image. This is one of the values listed in [Image_2.SetColors], except 0 which
// will be replaced by the actual number of colors used.
func (i *Image_2) GetColors() int {
	switch {
	case i.colors == 0:
		return availableColors
	case i.colors <= 2:
		return 2
	case i.colors <= 8:
		return 8
	case i.colors <= 256:
		return 256
	}
	return TrueColor
}

// SetDithering sets the dithering algorithm to use, one of the constants
// starting with "Dithering", for example [DitheringFloydSteinberg] (the
// default). Dithering is not applied when rendering in true-color.
func (i *Image_2) SetDithering(dithering int) *Image_2 {
	i.dithering = dithering
	i.lastWidth, i.lastHeight = 0, 0
	return i
}

// SetAspectRatio sets the width of a terminal's cell divided by its height.
// You may change the default of 0.5 if your terminal / font has a different
// aspect ratio. This is used to calculate the size of the image if the
// specified width or height is 0. The function will panic if the aspect ratio
// is 0 or less.
func (i *Image_2) SetAspectRatio(aspectRatio float64) *Image_2 {
	if aspectRatio <= 0 {
		panic("aspect ratio must be greater than 0")
	}
	i.aspectRatio = aspectRatio
	i.lastWidth, i.lastHeight = 0, 0
	return i
}

// SetAlign sets the vertical and horizontal alignment of the image within the
// widget's space. The possible values are [AlignTop], [AlignCenter], and
// [AlignBottom] for vertical alignment and [AlignLeft], [AlignCenter], and
// [AlignRight] for horizontal alignment. The default is [AlignCenter] for both
// (or [AlignTop] and [AlignLeft] if the image is part of a [Form]).
func (i *Image_2) SetAlign(vertical, horizontal int) *Image_2 {
	i.alignHorizontal = horizontal
	i.alignVertical = vertical
	return i
}

// SetLabel sets the text to be displayed before the image.
func (i *Image_2) SetLabel(label string) *Image_2 {
	i.label = label
	return i
}

// GetLabel returns the text to be displayed before the image.
func (i *Image_2) GetLabel() string {
	return i.label
}

// SetLabelWidth sets the screen width of the label. A value of 0 will cause the
// primitive to use the width of the label string.
func (i *Image_2) SetLabelWidth(width int) *Image_2 {
	i.labelWidth = width
	return i
}

// GetFieldWidth returns this primitive's field width. This is the image's width
// or, if the width is 0 or less, the proportional width of the image based on
// its height as returned by [Image_2.GetFieldHeight]. If there is no image, 0 is
// returned.
func (i *Image_2) GetFieldWidth() int {
	if i.width <= 0 {
		if i.image == nil {
			return 0
		}
		bounds := i.image.Bounds()
		height := i.GetFieldHeight()
		return bounds.Dx() * height / bounds.Dy()
	}
	return i.width
}

// GetFieldHeight returns this primitive's field height. This is the image's
// height or 8 if the height is 0 or less.
func (i *Image_2) GetFieldHeight() int {
	if i.height <= 0 {
		return 8
	}
	return i.height
}

// SetDisabled sets whether or not the item is disabled / read-only.
func (i *Image_2) SetDisabled(disabled bool) FormItem {
	return i // Images are always read-only.
}

// GetDisabled returns whether or not the item is disabled / read-only.
func (i *Image_2) GetDisabled() bool {
	return true // Images are always read-only.
}

// SetFormAttributes sets attributes shared by all form items.
func (i *Image_2) SetFormAttributes(labelWidth int, labelColor, bgColor, fieldTextColor, fieldBgColor tcell.Color) FormItem {
	i.labelWidth = labelWidth
	i.backgroundColor = bgColor
	i.SetLabelStyle(tcell.StyleDefault.Foreground(labelColor).Background(bgColor))
	i.lastWidth, i.lastHeight = 0, 0
	return i
}

// SetLabelStyle sets the style of the label.
func (i *Image_2) SetLabelStyle(style tcell.Style) *Image_2 {
	i.labelStyle = style
	return i
}

// GetLabelStyle returns the style of the label.
func (i *Image_2) GetLabelStyle() tcell.Style {
	return i.labelStyle
}

// SetFinishedFunc sets a callback invoked when the user leaves this form item.
func (i *Image_2) SetFinishedFunc(handler func(key tcell.Key)) FormItem {
	i.finished = handler
	return i
}

// Focus is called when this primitive receives focus.
func (i *Image_2) Focus(delegate func(p Primitive)) {
	// If we're part of a form, there's nothing the user can do here so we're
	// finished.
	if i.finished != nil {
		i.finished(-1)
		return
	}

	i.Box.Focus(delegate)
}

// render re-populates the [Image_2.pixels] slice based on the current settings,
// if [Image_2.lastWidth] and [Image_2.lastHeight] don't match the current image's
// size. It also sets the new image size in these two variables.
func (i *Image_2) render() {
	// If there is no image, there are no pixels.
	if i.image == nil {
		i.pixels = nil
		return
	}

	// Calculate the new (terminal-space) image size.
	bounds := i.image.Bounds()
	imageWidth, imageHeight := bounds.Dx(), bounds.Dy()
	if i.aspectRatio != 1.0 {
		imageWidth = int(float64(imageWidth) / i.aspectRatio)
	}
	width, height := i.width, i.height
	_, _, innerWidth, innerHeight := i.GetInnerRect()
	if i.labelWidth > 0 {
		innerWidth -= i.labelWidth
	} else {
		innerWidth -= TaggedStringWidth(i.label)
	}
	if innerWidth <= 0 {
		i.pixels = nil
		return
	}
	if width == 0 && height == 0 {
		// Use all available space.
		width, height = innerWidth, innerHeight
		if adjustedWidth := imageWidth * height / imageHeight; adjustedWidth < width {
			width = adjustedWidth
		} else {
			height = imageHeight * width / imageWidth
		}
	} else {
		// Turn percentages into absolute values.
		if width < 0 {
			width = innerWidth * -width / 100
		}
		if height < 0 {
			height = innerHeight * -height / 100
		}
		if width == 0 {
			// Adjust the width.
			width = imageWidth * height / imageHeight
		} else if height == 0 {
			// Adjust the height.
			height = imageHeight * width / imageWidth
		}
	}
	if width <= 0 || height <= 0 {
		i.pixels = nil
		return
	}

	// If nothing has changed, we're done.
	if i.lastWidth == width && i.lastHeight == height {
		return
	}
	i.lastWidth, i.lastHeight = width, height // This could still be larger than the available space but that's ok for now.

	// Generate the initial pixels by resizing the image (8x8 per cell).
	pixels := i.resize()

	// Turn them into block elements with background/foreground colors.
	i.stamp(pixels)
}

// resize resizes the image to the current size and returns the result as a
// slice of pixels. It is assumed that [Image_2.lastWidth] (w) and
// [Image_2.lastHeight] (h) are positive, non-zero values, and the slice has a
// size of 64*w*h, with each pixel being represented by 3 float64 values in the
// range of 0-1. The factor of 64 is due to the fact that we calculate 8x8
// pixels per cell.
func (i *Image_2) resize() [][3]float64 {
	// Because most of the time, we will be downsizing the image, we don't even
	// attempt to do any fancy interpolation. For each target pixel, we
	// calculate a weighted average of the source pixels using their coverage
	// area.

	bounds := i.image.Bounds()
	srcWidth, srcHeight := bounds.Dx(), bounds.Dy()
	tgtWidth, tgtHeight := i.lastWidth*8, i.lastHeight*8
	coverageWidth, coverageHeight := float64(tgtWidth)/float64(srcWidth), float64(tgtHeight)/float64(srcHeight)
	pixels := make([][3]float64, tgtWidth*tgtHeight)
	weights := make([]float64, tgtWidth*tgtHeight)
	for srcY := bounds.Min.Y; srcY < bounds.Max.Y; srcY++ {
		for srcX := bounds.Min.X; srcX < bounds.Max.X; srcX++ {
			r32, g32, b32, _ := i.image.At(srcX, srcY).RGBA()
			r, g, b := float64(r32)/0xffff, float64(g32)/0xffff, float64(b32)/0xffff

			// Iterate over all target pixels. Outer loop is Y.
			startY := float64(srcY-bounds.Min.Y) * coverageHeight
			endY := startY + coverageHeight
			fromY, toY := int(startY), int(endY)
			for tgtY := fromY; tgtY <= toY && tgtY < tgtHeight; tgtY++ {
				coverageY := 1.0
				if tgtY == fromY {
					coverageY -= math.Mod(startY, 1.0)
				}
				if tgtY == toY {
					coverageY -= 1.0 - math.Mod(endY, 1.0)
				}

				// Inner loop is X.
				startX := float64(srcX-bounds.Min.X) * coverageWidth
				endX := startX + coverageWidth
				fromX, toX := int(startX), int(endX)
				for tgtX := fromX; tgtX <= toX && tgtX < tgtWidth; tgtX++ {
					coverageX := 1.0
					if tgtX == fromX {
						coverageX -= math.Mod(startX, 1.0)
					}
					if tgtX == toX {
						coverageX -= 1.0 - math.Mod(endX, 1.0)
					}

					// Add a weighted contribution to the target pixel.
					index := tgtY*tgtWidth + tgtX
					coverage := coverageX * coverageY
					pixels[index][0] += r * coverage
					pixels[index][1] += g * coverage
					pixels[index][2] += b * coverage
					weights[index] += coverage
				}
			}
		}
	}

	// Normalize the pixels.
	for index, weight := range weights {
		if weight > 0 {
			pixels[index][0] /= weight
			pixels[index][1] /= weight
			pixels[index][2] /= weight
		}
	}

	return pixels
}

// stamp takes the pixels generated by [Image_2.resize] and populates the
// [Image_2.pixels] slice accordingly.
func (i *Image_2) stamp(resized [][3]float64) {
	// For each 8x8 pixel block, we find the best block element to represent it,
	// given the available colors.
	i.pixels = make([]pixel, i.lastWidth*i.lastHeight)
	colors := i.GetColors()
	for row := 0; row < i.lastHeight; row++ {
		for col := 0; col < i.lastWidth; col++ {
			// Calculate an error for each potential block element + color. Keep
			// the one with the lowest error.

			// Note that the values in "resize" may lie outside [0, 1] due to
			// the error distribution during dithering.

			minMSE := math.MaxFloat64 // Mean squared error.
			var final [64][3]float64  // The final pixel values.
			for element, bits := range blockElements {
				// Calculate the average color for the pixels covered by the set
				// bits and unset bits.
				var (
					bg, fg  [3]float64
					setBits float64
					bit     uint64 = 1
				)
				for y := 0; y < 8; y++ {
					for x := 0; x < 8; x++ {
						index := (row*8+y)*i.lastWidth*8 + (col*8 + x)
						if bits&bit != 0 {
							fg[0] += resized[index][0]
							fg[1] += resized[index][1]
							fg[2] += resized[index][2]
							setBits++
						} else {
							bg[0] += resized[index][0]
							bg[1] += resized[index][1]
							bg[2] += resized[index][2]
						}
						bit <<= 1
					}
				}
				for ch := 0; ch < 3; ch++ {
					fg[ch] /= setBits
					if fg[ch] < 0 {
						fg[ch] = 0
					} else if fg[ch] > 1 {
						fg[ch] = 1
					}
					bg[ch] /= 64 - setBits
					if bg[ch] < 0 {
						bg[ch] = 0
					} else if bg[ch] > 1 {
						bg[ch] = 1
					}
				}

				// Quantize to the nearest acceptable color.
				for _, color := range []*[3]float64{&fg, &bg} {
					if colors == 2 {
						// Monochrome. The following weights correspond better
						// to human perception than the arithmetic mean.
						gray := 0.299*color[0] + 0.587*color[1] + 0.114*color[2]
						if gray < 0.5 {
							*color = [3]float64{0, 0, 0}
						} else {
							*color = [3]float64{1, 1, 1}
						}
					} else {
						for index, ch := range color {
							switch {
							case colors == 8:
								// Colors vary wildly for each terminal. Expect
								// suboptimal results.
								if ch < 0.5 {
									color[index] = 0
								} else {
									color[index] = 1
								}
							case colors == 256:
								color[index] = math.Round(ch*6) / 6
							}
						}
					}
				}

				// Calculate the error (and the final pixel values).
				var (
					mse         float64
					values      [64][3]float64
					valuesIndex int
				)
				bit = 1
				for y := 0; y < 8; y++ {
					for x := 0; x < 8; x++ {
						if bits&bit != 0 {
							values[valuesIndex] = fg
						} else {
							values[valuesIndex] = bg
						}
						index := (row*8+y)*i.lastWidth*8 + (col*8 + x)
						for ch := 0; ch < 3; ch++ {
							err := resized[index][ch] - values[valuesIndex][ch]
							mse += err * err
						}
						bit <<= 1
						valuesIndex++
					}
				}

				// Do we have a better match?
				if mse < minMSE {
					// Yes. Save it.
					minMSE = mse
					final = values
					index := row*i.lastWidth + col
					i.pixels[index].element = element
					i.pixels[index].style = tcell.StyleDefault.
						Foreground(tcell.NewRGBColor(int32(math.Min(255, fg[0]*255)), int32(math.Min(255, fg[1]*255)), int32(math.Min(255, fg[2]*255)))).
						Background(tcell.NewRGBColor(int32(math.Min(255, bg[0]*255)), int32(math.Min(255, bg[1]*255)), int32(math.Min(255, bg[2]*255))))
				}
			}

			// Check if there is a shade block which results in a smaller error.

			// What's the overall average color?
			var avg [3]float64
			for y := 0; y < 8; y++ {
				for x := 0; x < 8; x++ {
					index := (row*8+y)*i.lastWidth*8 + (col*8 + x)
					for ch := 0; ch < 3; ch++ {
						avg[ch] += resized[index][ch] / 64
					}
				}
			}
			for ch := 0; ch < 3; ch++ {
				if avg[ch] < 0 {
					avg[ch] = 0
				} else if avg[ch] > 1 {
					avg[ch] = 1
				}
			}

			// Quantize and choose shade element.
			element := BlockFullBlock
			var fg, bg tcell.Color
			shades := []rune{' ', BlockLightShade, BlockMediumShade, BlockDarkShade, BlockFullBlock}
			if colors == 2 {
				// Monochrome.
				gray := 0.299*avg[0] + 0.587*avg[1] + 0.114*avg[2] // See above for details.
				shade := int(math.Round(gray * 4))
				element = shades[shade]
				for ch := 0; ch < 3; ch++ {
					avg[ch] = float64(shade) / 4
				}
				bg = tcell.ColorBlack
				fg = tcell.ColorWhite
			} else if colors == TrueColor {
				// True color.
				fg = tcell.NewRGBColor(int32(math.Min(255, avg[0]*255)), int32(math.Min(255, avg[1]*255)), int32(math.Min(255, avg[2]*255)))
				bg = fg
			} else {
				// 8 or 256 colors.
				steps := 1.0
				if colors == 256 {
					steps = 6.0
				}
				var (
					lo, hi, pos [3]float64
					shade       float64
				)
				for ch := 0; ch < 3; ch++ {
					lo[ch] = math.Floor(avg[ch]*steps) / steps
					hi[ch] = math.Ceil(avg[ch]*steps) / steps
					if r := hi[ch] - lo[ch]; r > 0 {
						pos[ch] = (avg[ch] - lo[ch]) / r
						if math.Abs(pos[ch]-0.5) < math.Abs(shade-0.5) {
							shade = pos[ch]
						}
					}
				}
				shade = math.Round(shade * 4)
				element = shades[int(shade)]
				shade /= 4
				for ch := 0; ch < 3; ch++ { // Find the closest channel value.
					best := math.Abs(avg[ch] - (lo[ch] + (hi[ch]-lo[ch])*shade)) // Start shade from lo to hi.
					if value := math.Abs(avg[ch] - (hi[ch] - (hi[ch]-lo[ch])*shade)); value < best {
						best = value // Swap lo and hi.
						lo[ch], hi[ch] = hi[ch], lo[ch]
					}
					if value := math.Abs(avg[ch] - lo[ch]); value < best {
						best = value // Use lo.
						hi[ch] = lo[ch]
					}
					if value := math.Abs(avg[ch] - hi[ch]); value < best {
						lo[ch] = hi[ch] // Use hi.
					}
					avg[ch] = lo[ch] + (hi[ch]-lo[ch])*shade // Quantize.
				}
				bg = tcell.NewRGBColor(int32(math.Min(255, lo[0]*255)), int32(math.Min(255, lo[1]*255)), int32(math.Min(255, lo[2]*255)))
				fg = tcell.NewRGBColor(int32(math.Min(255, hi[0]*255)), int32(math.Min(255, hi[1]*255)), int32(math.Min(255, hi[2]*255)))
			}

			// Calculate the error (and the final pixel values).
			var (
				mse         float64
				values      [64][3]float64
				valuesIndex int
			)
			for y := 0; y < 8; y++ {
				for x := 0; x < 8; x++ {
					index := (row*8+y)*i.lastWidth*8 + (col*8 + x)
					for ch := 0; ch < 3; ch++ {
						err := resized[index][ch] - avg[ch]
						mse += err * err
					}
					values[valuesIndex] = avg
					valuesIndex++
				}
			}

			// Is this shade element better than the block element?
			if mse < minMSE {
				// Yes. Save it.
				final = values
				index := row*i.lastWidth + col
				i.pixels[index].element = element
				i.pixels[index].style = tcell.StyleDefault.Foreground(fg).Background(bg)
			}

			// Apply dithering.
			if colors < TrueColor && i.dithering == DitheringFloydSteinberg {
				// The dithering mask determines how the error is distributed.
				// Each element has three values: dx, dy, and weight (in 16th).
				var mask = [4][3]int{
					{1, 0, 7},
					{-1, 1, 3},
					{0, 1, 5},
					{1, 1, 1},
				}

				// We dither the 8x8 block as a 2x2 block, transferring errors
				// to its 2x2 neighbors.
				for ch := 0; ch < 3; ch++ {
					for y := 0; y < 2; y++ {
						for x := 0; x < 2; x++ {
							// What's the error for this 4x4 block?
							var err float64
							for dy := 0; dy < 4; dy++ {
								for dx := 0; dx < 4; dx++ {
									err += (final[(y*4+dy)*8+(x*4+dx)][ch] - resized[(row*8+(y*4+dy))*i.lastWidth*8+(col*8+(x*4+dx))][ch]) / 16
								}
							}

							// Distribute it to the 2x2 neighbors.
							for _, dist := range mask {
								for dy := 0; dy < 4; dy++ {
									for dx := 0; dx < 4; dx++ {
										targetX, targetY := (x+dist[0])*4+dx, (y+dist[1])*4+dy
										if targetX < 0 || col*8+targetX >= i.lastWidth*8 || targetY < 0 || row*8+targetY >= i.lastHeight*8 {
											continue
										}
										resized[(row*8+targetY)*i.lastWidth*8+(col*8+targetX)][ch] -= err * float64(dist[2]) / 16
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

// Draw draws this primitive onto the screen.
func (i *Image_2) Draw(screen tcell.Screen) {
	i.DrawForSubclass(screen, i)

	// Regenerate image if necessary.
	i.render()

	// Draw label.
	viewX, viewY, viewWidth, viewHeight := i.GetInnerRect()
	_, labelBg, _ := i.labelStyle.Decompose()
	if i.labelWidth > 0 {
		labelWidth := i.labelWidth
		if labelWidth > viewWidth {
			labelWidth = viewWidth
		}
		printWithStyle(screen, i.label, viewX, viewY, 0, labelWidth, AlignLeft, i.labelStyle, labelBg == tcell.ColorDefault)
		viewX += labelWidth
		viewWidth -= labelWidth
	} else {
		_, _, drawnWidth := printWithStyle(screen, i.label, viewX, viewY, 0, viewWidth, AlignLeft, i.labelStyle, labelBg == tcell.ColorDefault)
		viewX += drawnWidth
		viewWidth -= drawnWidth
	}

	// Determine image placement.
	x, y, width, height := viewX, viewY, i.lastWidth, i.lastHeight
	if i.alignHorizontal == AlignCenter {
		x += (viewWidth - width) / 2
	} else if i.alignHorizontal == AlignRight {
		x += viewWidth - width
	}
	if i.alignVertical == AlignCenter {
		y += (viewHeight - height) / 2
	} else if i.alignVertical == AlignBottom {
		y += viewHeight - height
	}

	// Draw the image.
	for row := 0; row < height; row++ {
		if y+row < viewY || y+row >= viewY+viewHeight {
			continue
		}
		for col := 0; col < width; col++ {
			if x+col < viewX || x+col >= viewX+viewWidth {
				continue
			}

			index := row*width + col
			screen.SetContent(x+col, y+row, i.pixels[index].element, nil, i.pixels[index].style)
		}
	}
}

// Primitive is the top-most interface for all graphical primitives.
type Primitive interface {
	// Draw draws this primitive onto the screen. Implementers can call the
	// screen's ShowCursor() function but should only do so when they have focus.
	// (They will need to keep track of this themselves.)
	Draw(screen tcell.Screen)

	// GetRect returns the current position of the primitive, x, y, width, and
	// height.
	GetRect() (int, int, int, int)

	// SetRect sets a new position of the primitive.
	SetRect(x, y, width, height int)

	// InputHandler returns a handler which receives key events when it has focus.
	// It is called by the Application class.
	//
	// A value of nil may also be returned, in which case this primitive cannot
	// receive focus and will not process any key events.
	//
	// The handler will receive the key event and a function that allows it to
	// set the focus to a different primitive, so that future key events are sent
	// to that primitive.
	//
	// The Application's Draw() function will be called automatically after the
	// handler returns.
	//
	// The Box class provides functionality to intercept keyboard input. If you
	// subclass from Box, it is recommended that you wrap your handler using
	// Box.WrapInputHandler() so you inherit that functionality.
	InputHandler() func(event *tcell.EventKey, setFocus func(p Primitive))

	// Focus is called by the application when the primitive receives focus.
	// Implementers may call delegate() to pass the focus on to another
	// primitive which is usually a child primitive. This is not called on
	// parents of the primitive that receives focus.
	Focus(delegate func(p Primitive))

	// HasFocus determines if the primitive (or any of its child primitives) has
	// focus.
	HasFocus() bool

	// Blur is called by the application when the primitive loses focus. This is
	// not called on parents of the primitive that loses focus.
	Blur()

	// MouseHandler returns a handler which receives mouse events.
	// It is called by the Application class.
	//
	// A value of nil may also be returned to stop the downward propagation of
	// mouse events.
	//
	// The Box class provides functionality to intercept mouse events. If you
	// subclass from Box, it is recommended that you wrap your handler using
	// Box.WrapMouseHandler() so you inherit that functionality.
	MouseHandler() func(action MouseAction, event *tcell.EventMouse, setFocus func(p Primitive)) (consumed bool, capture Primitive)

	// PasteHandler returns a handler which receives pasted text.
	// It is called by the Application class.
	//
	// A value of nil may also be returned to stop the downward propagation of
	// paste events.
	//
	// The Box class may provide functionality to intercept paste events in the
	// future. If you subclass from [Box], it is recommended that you wrap your
	// handler using Box.WrapPasteHandler() so you inherit that functionality.
	PasteHandler() func(text string, setFocus func(p Primitive))

	// focusChain adds the chain of primitives that have focus to the given
	// slice, starting with the bottom-most primitive that has focus and ending
	// with this box. If this box or none of its descendents has focus, the
	// slice is not modified. If chain is nil, no chain is added. Returns
	// whether or not this box or one of its descendents has focus.
	focusChain(chain *[]Primitive) bool

	// focused is called when the current input focus changes. It is called on
	// the primitive which newly received focus as well as on all of its
	// ancestors (in no defined order). The default implementation in [Box]
	// invokes the callback set with [Box.SetFocusFunc]. This can also happen
	// when the focus is set to the primitive that already has focus.
	focused()

	// blurred is called when the current input focus changes. It is called on
	// the primitive which lost focus as well as on all of its ancestors (in no
	// defined order). The default implementation in [Box] invokes the callback
	// set with [Box.SetBlurFunc]. This can also happen when the focus is set to
	// the primitive that already has focus.
	blurred()
}

// Box implements the Primitive interface with an empty background and optional
// elements such as a border and a title. Box itself does not hold any content
// but serves as the superclass of all other primitives. Subclasses add their
// own content, typically (but not necessarily) keeping their content within the
// box's rectangle.
//
// Box provides a number of utility functions available to all primitives.
//
// See https://github.com/rivo/tview/wiki/Box for an example.
type Box struct {
	// Points to the implementing primitive at the bottom of the hierarchy.
	Primitive

	// The position of the rect.
	x, y, width, height int

	// The inner rect reserved for the box's content. If innerX is negative,
	// the rect is undefined and must be calculated.
	innerX, innerY, innerWidth, innerHeight int

	// Border padding.
	paddingTop, paddingBottom, paddingLeft, paddingRight int

	// The box's background color.
	backgroundColor tcell.Color

	// If set to true, the background of this box is not cleared while drawing.
	dontClear bool

	// Whether or not a border is drawn, reducing the box's space for content by
	// two in width and height.
	border bool

	// The border style.
	borderStyle tcell.Style

	// The title. Only visible if there is a border, too.
	title string

	// The color of the title.
	titleColor tcell.Color

	// The alignment of the title.
	titleAlign int

	// Whether or not this box has focus. At any time, this must be true only
	// for one primitive in the entire application. Such a primitive is usually
	// a visible and enabled widget but may also be a container primitive (if
	// no contained primitive has focus) or a primitive inaccessible to the user
	// (e.g. a child primitive of a widget to which interaction is delegated).
	hasFocus bool

	// Optional callback functions invoked when the primitive receives or loses
	// focus.
	focus, blur func()

	// Callback function invoked when the box itself is resized, nil if not set.
	boxResize func()

	// Callback function invoked when the box's inner content area is resized,
	// nil if not set.
	contentResize func()

	// An optional capture function which receives a key event and returns the
	// event to be forwarded to the primitive's default input handler (nil if
	// nothing should be forwarded).
	inputCapture func(event *tcell.EventKey) *tcell.EventKey

	// An optional function which is called before the box is drawn.
	draw func(screen tcell.Screen, x, y, width, height int) (int, int, int, int)

	// An optional capture function which receives a mouse event and returns the
	// event to be forwarded to the primitive's default mouse event handler (at
	// least one nil if nothing should be forwarded).
	mouseCapture func(action MouseAction, event *tcell.EventMouse) (MouseAction, *tcell.EventMouse)
}

// NewBox returns a [Box] without a border.
func NewBox() *Box {
	b := &Box{
		width:           15,
		height:          10,
		innerX:          -1, // Mark as uninitialized.
		backgroundColor: Styles.PrimitiveBackgroundColor,
		borderStyle:     tcell.StyleDefault.Foreground(Styles.BorderColor).Background(Styles.PrimitiveBackgroundColor),
		titleColor:      Styles.TitleColor,
		titleAlign:      AlignCenter,
	}
	b.Primitive = b
	return b
}

// SetBorderPadding sets the size of the borders around the box content.
func (b *Box) SetBorderPadding(top, bottom, left, right int) *Box {
	b.paddingTop, b.paddingBottom, b.paddingLeft, b.paddingRight = top, bottom, left, right
	return b
}

// GetRect returns the current position of the rectangle, x, y, width, and
// height.
func (b *Box) GetRect() (int, int, int, int) {
	return b.x, b.y, b.width, b.height
}

// GetInnerRect returns the position of the inner rectangle (x, y, width,
// height), without the border and without any padding. Width and height values
// will clamp to 0 and thus never be negative.
func (b *Box) GetInnerRect() (int, int, int, int) {
	if b.innerX >= 0 {
		return b.innerX, b.innerY, b.innerWidth, b.innerHeight
	}
	x, y, width, height := b.GetRect()
	if b.border {
		x++
		y++
		width -= 2
		height -= 2
	}
	x, y, width, height = x+b.paddingLeft,
		y+b.paddingTop,
		width-b.paddingLeft-b.paddingRight,
		height-b.paddingTop-b.paddingBottom
	if width < 0 {
		width = 0
	}
	if height < 0 {
		height = 0
	}
	return x, y, width, height
}

// SetRect sets a new position of the primitive. Note that this has no effect
// if this primitive is part of a layout (e.g. Flex, Grid) or if it was added
// like this:
//
//	application.SetRoot(p, true)
func (b *Box) SetRect(x, y, width, height int) {
	b.x = x
	b.y = y
	b.width, width = width, b.width
	b.height, height = height, b.height
	if b.width != width || b.height != height {
		if b.boxResize != nil {
			b.boxResize()
		}
		if b.contentResize != nil {
			b.contentResize()
		}
	}
	b.innerX = -1 // Mark inner rect as uninitialized.
}

// SetDrawFunc sets a callback function which is invoked after the box primitive
// has been drawn. This allows you to add a more individual style to the box
// (and all primitives which extend it).
//
// The function is provided with the box's dimensions (set via SetRect()). It
// must return the box's inner dimensions (x, y, width, height) which will be
// returned by GetInnerRect(), used by descendent primitives to draw their own
// content.
func (b *Box) SetDrawFunc(handler func(screen tcell.Screen, x, y, width, height int) (int, int, int, int)) *Box {
	b.draw = handler
	return b
}

// GetDrawFunc returns the callback function which was installed with
// SetDrawFunc() or nil if no such function has been installed.
func (b *Box) GetDrawFunc() func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
	return b.draw
}

// SetBoxResizeFunc sets a callback function which is invoked when the size of
// the box itself changes. Note that this is not called when the box is moved
// (i.e. when only x and y change). Set to nil to remove the callback function.
func (b *Box) SetBoxResizeFunc(handler func()) *Box {
	b.boxResize = handler
	return b
}

// SetContentResizeFunc sets a callback function which is invoked when the size
// of the box's inner content area changes. Note that this is not called when
// the area is moved (i.e. when only x and y change). Set to nil to remove the
// callback function.
func (b *Box) SetContentResizeFunc(handler func()) *Box {
	b.contentResize = handler
	return b
}

// WrapInputHandler wraps an input handler (see [Box.InputHandler]) with the
// functionality to capture input (see [Box.SetInputCapture]) before passing it
// on to the provided (default) input handler.
//
// This is only meant to be used by subclassing primitives.
func (b *Box) WrapInputHandler(inputHandler func(*tcell.EventKey, func(p Primitive))) func(*tcell.EventKey, func(p Primitive)) {
	return func(event *tcell.EventKey, setFocus func(p Primitive)) {
		if b.inputCapture != nil {
			event = b.inputCapture(event)
		}
		if event != nil && inputHandler != nil {
			inputHandler(event, setFocus)
		}
	}
}

// InputHandler returns nil. Box has no default input handling.
func (b *Box) InputHandler() func(event *tcell.EventKey, setFocus func(p Primitive)) {
	return b.WrapInputHandler(nil)
}

// WrapPasteHandler wraps a paste handler (see [Box.PasteHandler]).
func (b *Box) WrapPasteHandler(pasteHandler func(string, func(p Primitive))) func(string, func(p Primitive)) {
	return func(text string, setFocus func(p Primitive)) {
		if pasteHandler != nil {
			pasteHandler(text, setFocus)
		}
	}
}

// PasteHandler returns nil. Box has no default paste handling.
func (b *Box) PasteHandler() func(pastedText string, setFocus func(p Primitive)) {
	return b.WrapPasteHandler(nil)
}

// SetInputCapture installs a function which captures key events before they are
// forwarded to the primitive's default key event handler. This function can
// then choose to forward that key event (or a different one) to the default
// handler by returning it. If nil is returned, the default handler will not
// be called.
//
// Providing a nil handler will remove a previously existing handler.
//
// This function can also be used on container primitives (like Flex, Grid, or
// Form) as keyboard events will be handed down until they are handled.
//
// Pasted key events are not forwarded to the input capture function if pasting
// is enabled (see [Application.EnablePaste]).
func (b *Box) SetInputCapture(capture func(event *tcell.EventKey) *tcell.EventKey) *Box {
	b.inputCapture = capture
	return b
}

// GetInputCapture returns the function installed with SetInputCapture() or nil
// if no such function has been installed.
func (b *Box) GetInputCapture() func(event *tcell.EventKey) *tcell.EventKey {
	return b.inputCapture
}

// WrapMouseHandler wraps a mouse event handler (see [Box.MouseHandler]) with the
// functionality to capture mouse events (see [Box.SetMouseCapture]) before passing
// them on to the provided (default) event handler.
//
// This is only meant to be used by subclassing primitives.
func (b *Box) WrapMouseHandler(mouseHandler func(MouseAction, *tcell.EventMouse, func(p Primitive)) (bool, Primitive)) func(action MouseAction, event *tcell.EventMouse, setFocus func(p Primitive)) (consumed bool, capture Primitive) {
	return func(action MouseAction, event *tcell.EventMouse, setFocus func(p Primitive)) (consumed bool, capture Primitive) {
		if b.mouseCapture != nil {
			action, event = b.mouseCapture(action, event)
		}
		if event == nil {
			if action == MouseConsumed {
				consumed = true
			}
		} else if mouseHandler != nil {
			consumed, capture = mouseHandler(action, event, setFocus)
		}
		return
	}
}

// MouseHandler returns nil. Box has no default mouse handling.
func (b *Box) MouseHandler() func(action MouseAction, event *tcell.EventMouse, setFocus func(p Primitive)) (consumed bool, capture Primitive) {
	return b.WrapMouseHandler(func(action MouseAction, event *tcell.EventMouse, setFocus func(p Primitive)) (consumed bool, capture Primitive) {
		if action == MouseLeftDown && b.InRect(event.Position()) {
			setFocus(b)
			consumed = true
		}
		return
	})
}

// SetMouseCapture sets a function which captures mouse events (consisting of
// the original tcell mouse event and the semantic mouse action) before they are
// forwarded to the primitive's default mouse event handler. This function can
// then choose to forward that event (or a different one) by returning it or
// returning a nil mouse event, in which case the default handler will not be
// called.
//
// When a nil event is returned, the returned mouse action value may be set to
// [MouseConsumed] to indicate that the event was consumed and the screen should
// be redrawn. Any other value will not cause a redraw.
//
// Providing a nil handler will remove a previously existing handler.
//
// Note that mouse events are ignored completely if the application has not been
// enabled for mouse events (see [Application.EnableMouse]), which is the
// default.
func (b *Box) SetMouseCapture(capture func(action MouseAction, event *tcell.EventMouse) (MouseAction, *tcell.EventMouse)) *Box {
	b.mouseCapture = capture
	return b
}

// InRect returns true if the given coordinate is within the bounds of the box's
// rectangle.
func (b *Box) InRect(x, y int) bool {
	rectX, rectY, width, height := b.GetRect()
	return x >= rectX && x < rectX+width && y >= rectY && y < rectY+height
}

// InInnerRect returns true if the given coordinate is within the bounds of the
// box's inner rectangle (within the border and padding).
func (b *Box) InInnerRect(x, y int) bool {
	rectX, rectY, width, height := b.GetInnerRect()
	return x >= rectX && x < rectX+width && y >= rectY && y < rectY+height
}

// GetMouseCapture returns the function installed with SetMouseCapture() or nil
// if no such function has been installed.
func (b *Box) GetMouseCapture() func(action MouseAction, event *tcell.EventMouse) (MouseAction, *tcell.EventMouse) {
	return b.mouseCapture
}

// SetBackgroundColor sets the box's background color.
func (b *Box) SetBackgroundColor(color tcell.Color) *Box {
	b.backgroundColor = color
	b.borderStyle = b.borderStyle.Background(color)
	return b
}

// SetBorder sets the flag indicating whether or not the box should have a
// border.
func (b *Box) SetBorder(show bool) *Box {
	b.border, show = show, b.border
	if b.border != show {
		if b.contentResize != nil {
			b.contentResize()
		}
	}
	return b
}

// SetBorderStyle sets the box's border style.
func (b *Box) SetBorderStyle(style tcell.Style) *Box {
	b.borderStyle = style
	return b
}

// SetBorderColor sets the box's border color.
func (b *Box) SetBorderColor(color tcell.Color) *Box {
	b.borderStyle = b.borderStyle.Foreground(color)
	return b
}

// SetBorderAttributes sets the border's style attributes. You can combine
// different attributes using bitmask operations:
//
//	box.SetBorderAttributes(tcell.AttrItalic | tcell.AttrBold)
func (b *Box) SetBorderAttributes(attr tcell.AttrMask) *Box {
	b.borderStyle = b.borderStyle.Attributes(attr)
	return b
}

// GetBorderAttributes returns the border's style attributes.
func (b *Box) GetBorderAttributes() tcell.AttrMask {
	_, _, attr := b.borderStyle.Decompose()
	return attr
}

// GetBorderColor returns the box's border color.
func (b *Box) GetBorderColor() tcell.Color {
	color, _, _ := b.borderStyle.Decompose()
	return color
}

// GetBackgroundColor returns the box's background color.
func (b *Box) GetBackgroundColor() tcell.Color {
	return b.backgroundColor
}

// SetTitle sets the box's title.
func (b *Box) SetTitle(title string) *Box {
	b.title = title
	return b
}

// GetTitle returns the box's current title.
func (b *Box) GetTitle() string {
	return b.title
}

// SetTitleColor sets the box's title color.
func (b *Box) SetTitleColor(color tcell.Color) *Box {
	b.titleColor = color
	return b
}

// SetTitleAlign sets the alignment of the title, one of AlignLeft, AlignCenter,
// or AlignRight.
func (b *Box) SetTitleAlign(align int) *Box {
	b.titleAlign = align
	return b
}

// Draw draws this primitive onto the screen.
func (b *Box) Draw(screen tcell.Screen) {
	b.DrawForSubclass(screen, b)
}

// DrawForSubclass draws this box under the assumption that primitive p is a
// subclass of this box. This is needed e.g. to draw proper box frames which
// depend on the subclass's focus.
//
// Only call this function from your own custom primitives. It is not needed in
// applications that have no custom primitives.
func (b *Box) DrawForSubclass(screen tcell.Screen, p Primitive) {
	// Don't draw anything if there is no space.
	if b.width <= 0 || b.height <= 0 {
		return
	}

	// Fill background.
	background := tcell.StyleDefault.Background(b.backgroundColor)
	if !b.dontClear {
		for y := b.y; y < b.y+b.height; y++ {
			for x := b.x; x < b.x+b.width; x++ {
				screen.SetContent(x, y, ' ', nil, background)
			}
		}
	}

	// Draw border.
	if b.border && b.width >= 2 && b.height >= 2 {
		var vertical, horizontal, topLeft, topRight, bottomLeft, bottomRight rune
		if p.HasFocus() {
			horizontal = Borders.HorizontalFocus
			vertical = Borders.VerticalFocus
			topLeft = Borders.TopLeftFocus
			topRight = Borders.TopRightFocus
			bottomLeft = Borders.BottomLeftFocus
			bottomRight = Borders.BottomRightFocus
		} else {
			horizontal = Borders.Horizontal
			vertical = Borders.Vertical
			topLeft = Borders.TopLeft
			topRight = Borders.TopRight
			bottomLeft = Borders.BottomLeft
			bottomRight = Borders.BottomRight
		}
		for x := b.x + 1; x < b.x+b.width-1; x++ {
			screen.SetContent(x, b.y, horizontal, nil, b.borderStyle)
			screen.SetContent(x, b.y+b.height-1, horizontal, nil, b.borderStyle)
		}
		for y := b.y + 1; y < b.y+b.height-1; y++ {
			screen.SetContent(b.x, y, vertical, nil, b.borderStyle)
			screen.SetContent(b.x+b.width-1, y, vertical, nil, b.borderStyle)
		}
		screen.SetContent(b.x, b.y, topLeft, nil, b.borderStyle)
		screen.SetContent(b.x+b.width-1, b.y, topRight, nil, b.borderStyle)
		screen.SetContent(b.x, b.y+b.height-1, bottomLeft, nil, b.borderStyle)
		screen.SetContent(b.x+b.width-1, b.y+b.height-1, bottomRight, nil, b.borderStyle)

		// Draw title.
		if b.title != "" && b.width >= 4 {
			printed, _ := Print_2(screen, b.title, b.x+1, b.y, b.width-2, b.titleAlign, b.titleColor)
			if len(b.title)-printed > 0 && printed > 0 {
				xEllipsis := b.x + b.width - 2
				if b.titleAlign == AlignRight {
					xEllipsis = b.x + 1
				}
				_, _, style, _ := screen.GetContent(xEllipsis, b.y)
				fg, _, _ := style.Decompose()
				Print_2(screen, string(SemigraphicsHorizontalEllipsis), xEllipsis, b.y, 1, AlignLeft, fg)
			}
		}
	}

	// Call custom draw function.
	if b.draw != nil {
		b.innerX, b.innerY, b.innerWidth, b.innerHeight = b.draw(screen, b.x, b.y, b.width, b.height)
	} else {
		// Remember the inner rect.
		b.innerX = -1
		b.innerX, b.innerY, b.innerWidth, b.innerHeight = b.GetInnerRect()
	}
}

// SetFocusFunc sets a callback function which is invoked when this primitive
// receives focus. Container primitives such as [Flex] or [Grid] will also be
// notified if one of their descendents receive focus directly. Note that this
// may result in a blur notification, immediately followed by a focus
// notification, when the focus is set to a different descendent of the
// container primitive.
//
// At this point, the order in which the focus callbacks are invoked during one
// draw cycle, is not defined. However, the blur callbacks are always invoked
// before the focus callbacks.
//
// Set to nil to remove the callback function.
func (b *Box) SetFocusFunc(callback func()) *Box {
	b.focus = callback
	return b
}

// SetBlurFunc sets a callback function which is invoked when this primitive
// loses focus. Container primitives such as [Flex] or [Grid] will also be
// notified if one of their descendents lose focus. Note that this may result in
// a blur notification, immediately followed by a focus notification, when the
// focus is set to a different different descendent of the container primitive.
//
// At this point, the order in which the blur callbacks are invoked during one
// draw cycle, is not defined. However, the blur callbacks are always invoked
// before the focus callbacks.
//
// Set to nil to remove the callback function.
func (b *Box) SetBlurFunc(callback func()) *Box {
	b.blur = callback
	return b
}

// Focus is called when this primitive directly receives focus.
func (b *Box) Focus(delegate func(p Primitive)) {
	b.hasFocus = true
}

// focused is called when this primitive or one of its descendents receives
// focus.
func (b *Box) focused() {
	if b.focus != nil {
		b.focus()
	}
}

// Blur is called when this primitive directly loses focus.
func (b *Box) Blur() {
	b.hasFocus = false
}

// blurred is called when this primitive or one of its descendents loses focus.
func (b *Box) blurred() {
	if b.blur != nil {
		b.blur()
	}
}

// HasFocus returns whether or not this primitive has focus.
func (b *Box) HasFocus() bool {
	return b.Primitive.focusChain(nil)
}

// focusChain implements the [Primitive]'s focusChain method.
func (b *Box) focusChain(chain *[]Primitive) bool {
	if !b.hasFocus {
		return false
	}
	if chain != nil {
		*chain = append(*chain, b.Primitive)
	}
	return true
}

// FormItem is the interface all form items must implement to be able to be
// included in a form.
type FormItem interface {
	Primitive

	// GetLabel returns the item's label text.
	GetLabel() string

	// SetFormAttributes sets a number of item attributes at once.
	SetFormAttributes(labelWidth int, labelColor, bgColor, fieldTextColor, fieldBgColor tcell.Color) FormItem

	// GetFieldWidth returns the width of the form item's field (the area which
	// is manipulated by the user) in number of screen cells. A value of 0
	// indicates the field width is flexible and may use as much space as
	// required.
	GetFieldWidth() int

	// GetFieldHeight returns the height of the form item's field (the area which
	// is manipulated by the user). This value must be greater than 0.
	GetFieldHeight() int

	// SetFinishedFunc sets the handler function for when the user finished
	// entering data into the item. The handler may receive events for the
	// Enter key (we're done), the Escape key (cancel input), the Tab key (move
	// to next field), the Backtab key (move to previous field), or a negative
	// value, indicating that the action for the last known key should be
	// repeated.
	SetFinishedFunc(handler func(key tcell.Key)) FormItem

	// SetDisabled sets whether or not the item is disabled / read-only. A form
	// must have at least one item that is not disabled.
	SetDisabled(disabled bool) FormItem

	// GetDisabled returns whether or not the item is disabled / read-only.
	GetDisabled() bool
}

// escapedTagPattern matches an escaped tag, e.g. "[red[]", at the beginning of
// a string.
var escapedTagPattern = regexp.MustCompile(`^\[[^\[\]]+\[+\]`)

// stepOptions is a bit field of options for [step]. A value of 0 results in
// [step] having the same behavior as uniseg.Step, i.e. no tview-related parsing
// is performed.
type stepOptions int

// Bit fields for [stepOptions].
const (
	stepOptionsNone   stepOptions = 0
	stepOptionsStyle  stepOptions = 1 << iota // Parse style tags.
	stepOptionsRegion                         // Parse region tags.
)

// stepState represents the current state of the parser implemented in [step].
type stepState struct {
	unisegState     int         // The state of the uniseg parser.
	boundaries      int         // Information about boundaries, as returned by uniseg.Step.
	style           tcell.Style // The style of the returned grapheme cluster.
	region          string      // The region of the returned grapheme cluster.
	escapedTagState int         // States for parsing escaped tags (defined in [step]).
	grossLength     int         // The length of the cluster, including any tags not returned.

	// The styles for the initial call to [step].
	initialForeground tcell.Color
	initialBackground tcell.Color
	initialAttributes tcell.AttrMask
}

// IsWordBoundary returns true if the boundary between the returned grapheme
// cluster and the one following it is a word boundary.
func (s *stepState) IsWordBoundary() bool {
	return s.boundaries&uniseg.MaskWord != 0
}

// IsSentenceBoundary returns true if the boundary between the returned grapheme
// cluster and the one following it is a sentence boundary.
func (s *stepState) IsSentenceBoundary() bool {
	return s.boundaries&uniseg.MaskSentence != 0
}

// LineBreak returns whether the string can be broken into the next line after
// the returned grapheme cluster. If optional is true, the line break is
// optional. If false, the line break is mandatory, e.g. after a newline
// character.
func (s *stepState) LineBreak() (lineBreak, optional bool) {
	switch s.boundaries & uniseg.MaskLine {
	case uniseg.LineCanBreak:
		return true, true
	case uniseg.LineMustBreak:
		return true, false
	}
	return false, false // uniseg.LineDontBreak.
}

// Width returns the grapheme cluster's width in cells.
func (s *stepState) Width() int {
	return s.boundaries >> uniseg.ShiftWidth
}

// GrossLength returns the grapheme cluster's length in bytes, including any
// tags that were parsed but not explicitly returned.
func (s *stepState) GrossLength() int {
	return s.grossLength
}

// Style returns the style for the grapheme cluster.
func (s *stepState) Style() tcell.Style {
	return s.style
}

// step uses uniseg.Step to iterate over the grapheme clusters of a string but
// (optionally) also parses the string for style or region tags.
//
// This function can be called consecutively to extract all grapheme clusters
// from str, without returning any contained (parsed) tags. The return values
// are the first grapheme cluster, the remaining string, and the new state. Pass
// the remaining string and the returned state to the next call. If the rest
// string is empty, parsing is complete. Call the returned state's methods for
// boundary and cluster width information.
//
// The returned cluster may be empty if the given string consists of only
// (parsed) tags. The boundary and width information will be meaningless in
// this case but the style will describe the style at the end of the string.
//
// Pass nil for state on the first call. This will assume an initial style with
// [Styles.PrimitiveBackgroundColor] as the background color and
// [Styles.PrimaryTextColor] as the text color, no current region. If you want
// to start with a different style or region, you can set the state accordingly
// but you must then set [state.unisegState] to -1.
//
// There is no need to call uniseg.HasTrailingLineBreakInString on the last
// non-empty cluster as this function will do this for you and adjust the
// returned boundaries accordingly.
func step(str string, state *stepState, opts stepOptions) (cluster, rest string, newState *stepState) {
	// Set up initial state.
	if state == nil {
		state = &stepState{
			unisegState: -1,
			style:       tcell.StyleDefault.Background(Styles.PrimitiveBackgroundColor).Foreground(Styles.PrimaryTextColor),
		}
	}
	if state.unisegState < 0 {
		state.initialForeground, state.initialBackground, state.initialAttributes = state.style.Decompose()
	}
	if len(str) == 0 {
		newState = state
		return
	}

	// Get a grapheme cluster.
	preState := state.unisegState
	cluster, rest, state.boundaries, state.unisegState = uniseg.StepString(str, preState)
	state.grossLength = len(cluster)
	if rest == "" {
		if !uniseg.HasTrailingLineBreakInString(cluster) {
			state.boundaries &^= uniseg.MaskLine
		}
	}

	// Parse tags.
	if opts != stepOptionsNone {
		const (
			etNone int = iota
			etStart
			etChar
			etClosing
		)

		// Finite state machine for escaped tags.
		switch state.escapedTagState {
		case etStart:
			if cluster[0] == '[' || cluster[0] == ']' { // Invalid escaped tag.
				state.escapedTagState = etNone
			} else { // Other characters are allowed.
				state.escapedTagState = etChar
			}
		case etChar:
			if cluster[0] == ']' { // In theory, this should not happen.
				state.escapedTagState = etNone
			} else if cluster[0] == '[' { // Starting closing sequence.
				// Swallow the first one.
				cluster, rest, state.boundaries, state.unisegState = uniseg.StepString(rest, preState)
				state.grossLength += len(cluster)
				if cluster[0] == ']' {
					state.escapedTagState = etNone
				} else {
					state.escapedTagState = etClosing
				}
			} // More characters. Remain in etChar.
		case etClosing:
			if cluster[0] != '[' {
				state.escapedTagState = etNone
			}
		}

		// Regular tags.
		if state.escapedTagState == etNone {
			if cluster[0] == '[' {
				// We've already opened a tag. Parse it.
				length, style, region := parseTag(str, state, opts)
				if length > 0 {
					state.style = style
					state.region = region
					cluster, rest, state.boundaries, state.unisegState = uniseg.StepString(str[length:], preState)
					state.grossLength = len(cluster) + length
					if rest == "" {
						if !uniseg.HasTrailingLineBreakInString(cluster) {
							state.boundaries &^= uniseg.MaskLine
						}
					}
				}
				// Is this an escaped tag?
				if escapedTagPattern.MatchString(str[length:]) {
					state.escapedTagState = etStart
				}
			}
			if len(rest) > 0 && rest[0] == '[' {
				// A tag might follow the cluster. If so, we need to fix the state
				// for the boundaries to be correct.
				if length, _, _ := parseTag(rest, state, opts); length > 0 {
					if len(rest) > length {
						_, l := utf8.DecodeRuneInString(rest[length:])
						cluster += rest[length : length+l]
					}
					var taglessRest string
					cluster, taglessRest, state.boundaries, state.unisegState = uniseg.StepString(cluster, preState)
					if taglessRest == "" {
						if !uniseg.HasTrailingLineBreakInString(cluster) {
							state.boundaries &^= uniseg.MaskLine
						}
					}
				}
			}
		}
	}

	newState = state
	return
}

// parseTag parses str for consecutive style and/or region tags, assuming that
// str starts with the opening bracket for the first tag. It returns the string
// length of all valid tags (0 if the first tag is not valid) and the updated
// style and region for valid tags (based on the provided state).
func parseTag(str string, state *stepState, opts stepOptions) (length int, style tcell.Style, region string) {
	if opts == stepOptionsNone {
		return // No tags to parse.
	}

	// Automata states for parsing tags.
	const (
		tagStateNone = iota
		tagStateDoneTag
		tagStateStart
		tagStateRegionStart
		tagStateEndForeground
		tagStateStartBackground
		tagStateNumericForeground
		tagStateNameForeground
		tagStateEndBackground
		tagStateStartAttributes
		tagStateNumericBackground
		tagStateNameBackground
		tagStateAttributes
		tagStateRegionEnd
		tagStateRegionName
		tagStateEndAttributes
		tagStateStartURL
		tagStateEndURL
		tagStateURL
	)

	// Helper function which checks if the given byte is one of a list of
	// characters, including letters and digits.
	isOneOf := func(b byte, chars string) bool {
		if b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' || b >= '0' && b <= '9' {
			return true
		}
		return strings.IndexByte(chars, b) >= 0
	}

	// Attribute map.
	attrs := map[byte]tcell.AttrMask{
		'B': tcell.AttrBold,
		'I': tcell.AttrItalic,
		'L': tcell.AttrBlink,
		'D': tcell.AttrDim,
		'S': tcell.AttrStrikeThrough,
		'R': tcell.AttrReverse,
	}

	var (
		tagState, tagLength int
		tempStr             strings.Builder
	)
	tStyle := state.style
	tRegion := state.region

	// Process state transitions.
	for len(str) > 0 {
		ch := str[0]
		str = str[1:]
		tagLength++

		// Transition.
		switch tagState {
		case tagStateNone:
			if ch == '[' { // Start of a tag.
				tagState = tagStateStart
			} else { // Not a tag. We're done.
				return
			}
		case tagStateStart:
			if ch == '"' && opts&stepOptionsRegion == 0 {
				return // Region tags are not allowed.
			} else if ch != '"' && opts&stepOptionsStyle == 0 {
				return // Style tags are not allowed.
			}
			switch {
			case ch == '"': // Start of a region tag.
				tempStr.Reset()
				tagState = tagStateRegionStart
			case !isOneOf(ch, "#:-"): // Invalid style tag.
				return
			case ch == '-': // Reset foreground color.
				tStyle = tStyle.Foreground(state.initialForeground)
				tagState = tagStateEndForeground
			case ch == ':': // No foreground color.
				tagState = tagStateStartBackground
			default:
				tempStr.Reset()
				tempStr.WriteByte(ch)
				if ch == '#' { // Numeric foreground color.
					tagState = tagStateNumericForeground
				} else { // Letters or numbers.
					tagState = tagStateNameForeground
				}
			}
		case tagStateEndForeground:
			switch ch {
			case ']': // End of tag.
				tagState = tagStateDoneTag
			case ':':
				tagState = tagStateStartBackground
			default: // Invalid tag.
				return
			}
		case tagStateNumericForeground:
			if ch == ']' || ch == ':' {
				if tempStr.Len() != 7 { // Must be #rrggbb.
					return
				}
				tStyle = tStyle.Foreground(tcell.GetColor(tempStr.String()))
			}
			switch {
			case ch == ']': // End of tag.
				tagState = tagStateDoneTag
			case ch == ':': // Start of background color.
				tagState = tagStateStartBackground
			case strings.IndexByte("0123456789abcdefABCDEF", ch) >= 0: // Hex digit.
				tempStr.WriteByte(ch)
				tagState = tagStateNumericForeground
			default: // Invalid tag.
				return
			}
		case tagStateNameForeground:
			if ch == ']' || ch == ':' {
				name := tempStr.String()
				if name[0] >= '0' && name[0] <= '9' { // Must not start with a digit.
					return
				}
				tStyle = tStyle.Foreground(tcell.ColorNames[name])
			}
			switch {
			case !isOneOf(ch, "]:"): // Invalid tag.
				return
			case ch == ']': // End of tag.
				tagState = tagStateDoneTag
			case ch == ':': // Start of background color.
				tagState = tagStateStartBackground
			default: // Letters or numbers.
				tempStr.WriteByte(ch)
			}
		case tagStateStartBackground:
			switch {
			case !isOneOf(ch, "#:-]"): // Invalid style tag.
				return
			case ch == ']': // End of tag.
				tagState = tagStateDoneTag
			case ch == '-': // Reset background color.
				tStyle = tStyle.Background(state.initialBackground)
				tagState = tagStateEndBackground
			case ch == ':': // No background color.
				tagState = tagStateStartAttributes
			default:
				tempStr.Reset()
				tempStr.WriteByte(ch)
				if ch == '#' { // Numeric background color.
					tagState = tagStateNumericBackground
				} else { // Letters or numbers.
					tagState = tagStateNameBackground
				}
			}
		case tagStateEndBackground:
			switch ch {
			case ']': // End of tag.
				tagState = tagStateDoneTag
			case ':': // Start of attributes.
				tagState = tagStateStartAttributes
			default: // Invalid tag.
				return
			}
		case tagStateNumericBackground:
			if ch == ']' || ch == ':' {
				if tempStr.Len() != 7 { // Must be #rrggbb.
					return
				}
				tStyle = tStyle.Background(tcell.GetColor(tempStr.String()))
			}
			if ch == ']' { // End of tag.
				tagState = tagStateDoneTag
			} else if ch == ':' { // Start of attributes.
				tagState = tagStateStartAttributes
			} else if strings.IndexByte("0123456789abcdefABCDEF", ch) >= 0 { // Hex digit.
				tempStr.WriteByte(ch)
				tagState = tagStateNumericBackground
			} else { // Invalid tag.
				return
			}
		case tagStateNameBackground:
			if ch == ']' || ch == ':' {
				name := tempStr.String()
				if name[0] >= '0' && name[0] <= '9' { // Must not start with a digit.
					return
				}
				tStyle = tStyle.Background(tcell.ColorNames[name])
			}
			switch {
			case !isOneOf(ch, "]:"): // Invalid tag.
				return
			case ch == ']': // End of tag.
				tagState = tagStateDoneTag
			case ch == ':': // Start of background color.
				tagState = tagStateStartAttributes
			default: // Letters or numbers.
				tempStr.WriteByte(ch)
			}
		case tagStateStartAttributes:
			switch {
			case ch == ']': // End of tag.
				tagState = tagStateDoneTag
			case ch == '-': // Reset attributes.
				tStyle = tStyle.Attributes(state.initialAttributes)
				tagState = tagStateEndAttributes
			case ch == ':': // Start of URL.
				tagState = tagStateStartURL
			case strings.IndexByte("buildsrBUILDSR", ch) >= 0: // Attribute tag.
				tempStr.Reset()
				tempStr.WriteByte(ch)
				tagState = tagStateAttributes
			default: // Invalid tag.
				return
			}
		case tagStateAttributes:
			if ch == ']' || ch == ':' {
				flags := tempStr.String()
				_, _, a := tStyle.Decompose()
				for index := 0; index < len(flags); index++ {
					ch := flags[index]
					switch {
					case ch == 'u':
						tStyle = tStyle.Underline(true)
					case ch == 'U':
						tStyle = tStyle.Underline(false)
					case ch >= 'a' && ch <= 'z':
						a |= attrs[ch-('a'-'A')]
					default:
						a &^= attrs[ch]
					}
				}
				tStyle = tStyle.Attributes(a)
			}
			switch {
			case ch == ']': // End of tag.
				tagState = tagStateDoneTag
			case ch == ':': // Start of URL.
				tagState = tagStateStartURL
			case strings.IndexByte("buildsrBUILDSR", ch) >= 0: // Attribute tag.
				tempStr.WriteByte(ch)
			default: // Invalid tag.
				return
			}
		case tagStateEndAttributes:
			switch ch {
			case ']': // End of tag.
				tagState = tagStateDoneTag
			case ':': // Start of URL.
				tagState = tagStateStartURL
			default: // Invalid tag.
				return
			}
		case tagStateStartURL:
			switch ch {
			case ']': // End of tag.
				tagState = tagStateDoneTag
			case '-': // Reset URL.
				tStyle = tStyle.Url("").UrlId("")
				tagState = tagStateEndURL
			default: // URL character.
				tempStr.Reset()
				tempStr.WriteByte(ch)
				tStyle = tStyle.UrlId(strconv.Itoa(int(rand.Uint32()))) // Generate a unique ID for this URL.
				tagState = tagStateURL
			}
		case tagStateEndURL:
			if ch == ']' { // End of tag.
				tagState = tagStateDoneTag
			} else { // Invalid tag.
				return
			}
		case tagStateURL:
			if ch == ']' { // End of tag.
				tStyle = tStyle.Url(tempStr.String())
				tagState = tagStateDoneTag
			} else { // URL character.
				tempStr.WriteByte(ch)
			}
		case tagStateRegionStart:
			switch {
			case ch == '"': // End of region tag.
				tagState = tagStateRegionEnd
			case isOneOf(ch, "_,;: -."): // Region name.
				tempStr.WriteByte(ch)
				tagState = tagStateRegionName
			default: // Invalid tag.
				return
			}
		case tagStateRegionEnd:
			if ch == ']' { // End of tag.
				tRegion = tempStr.String()
				tagState = tagStateDoneTag
			} else { // Invalid tag.
				return
			}
		case tagStateRegionName:
			switch {
			case ch == '"': // End of region tag.
				tagState = tagStateRegionEnd
			case isOneOf(ch, "_,;: -."): // Region name.
				tempStr.WriteByte(ch)
			default: // Invalid tag.
				return
			}
		}

		// The last transition led to a tag end. Make the tag permanent.
		if tagState == tagStateDoneTag {
			length, style, region = tagLength, tStyle, tRegion
			tagState = tagStateNone // Reset state.
		}
	}

	return
}

// TaggedStringWidth returns the width of the given string needed to print it on
// screen. The text may contain style tags which are not counted.
func TaggedStringWidth(text string) (width int) {
	var state *stepState
	for len(text) > 0 {
		_, text, state = step(text, state, stepOptionsStyle)
		width += state.Width()
	}
	return
}

// WordWrap splits a text such that each resulting line does not exceed the
// given screen width. Split points are determined using the algorithm described
// in [Unicode Standard Annex #14].
//
// This function considers style tags to have no width.
//
// [Unicode Standard Annex #14]: https://www.unicode.org/reports/tr14/
func WordWrap(text string, width int) (lines []string) {
	if width <= 0 {
		return
	}

	var (
		state                                              *stepState
		lineWidth, lineLength, lastOption, lastOptionWidth int
	)
	str := text
	for len(str) > 0 {
		// Parse the next character.
		_, str, state = step(str, state, stepOptionsStyle)
		cWidth := state.Width()

		// Would it exceed the line width?
		if lineWidth+cWidth > width {
			if lastOptionWidth == 0 {
				// No split point so far. Just split at the current position.
				lines = append(lines, text[:lineLength])
				text = text[lineLength:]
				lineWidth, lineLength, lastOption, lastOptionWidth = 0, 0, 0, 0
			} else {
				// Split at the last split point.
				lines = append(lines, text[:lastOption])
				text = text[lastOption:]
				lineWidth -= lastOptionWidth
				lineLength -= lastOption
				lastOption, lastOptionWidth = 0, 0
			}
		}

		// Move ahead.
		lineWidth += cWidth
		lineLength += state.GrossLength()

		// Check for split points.
		if lineBreak, optional := state.LineBreak(); lineBreak {
			if optional {
				// Remember this split point.
				lastOption = lineLength
				lastOptionWidth = lineWidth
			} else {
				// We must split here.
				lines = append(lines, strings.TrimRight(text[:lineLength], "\n\r"))
				text = text[lineLength:]
				lineWidth, lineLength, lastOption, lastOptionWidth = 0, 0, 0, 0
			}
		}
	}
	lines = append(lines, text)

	return
}

// Escape escapes the given text such that color and/or region tags are not
// recognized and substituted by the print functions of this package. For
// example, to include a tag-like string in a box title or in a TextView:
//
//	box.SetTitle(tview.Escape("[squarebrackets]"))
//	fmt.Fprint(textView, tview.Escape(`["quoted"]`))
func Escape(text string) string {
	return escapePattern.ReplaceAllString(text, "$1[]")
}

// Unescape unescapes text previously escaped with [Escape].
func Unescape(text string) string {
	return unescapePattern.ReplaceAllString(text, "$1]")
}

// stripTags strips style tags from the given string. (Region tags are not
// stripped.)
func stripTags(text string) string {
	var (
		str   strings.Builder
		state *stepState
	)
	for len(text) > 0 {
		var c string
		c, text, state = step(text, state, stepOptionsStyle)
		str.WriteString(c)
	}
	return str.String()
}

// Text alignment within a box. Also used to align images.
const (
	AlignLeft = iota
	AlignCenter
	AlignRight
	AlignTop    = 0
	AlignBottom = 2
)

var (
	// Regular expression used to escape style/region tags.
	escapePattern = regexp.MustCompile(`(\[[a-zA-Z0-9_,;: \-\."#]+\[*)\]`)

	// Regular expression used to unescape escaped style/region tags.
	unescapePattern = regexp.MustCompile(`(\[[a-zA-Z0-9_,;: \-\."#]+\[*)\[\]`)

	// The number of colors available in the terminal.
	availableColors = 256
)

// Package initialization.
func init() {
	// Determine the number of colors available in the terminal.
	info, err := tcell.LookupTerminfo(os.Getenv("TERM"))
	if err == nil {
		availableColors = info.Colors
	}
}

// Print_2 prints text onto the screen into the given box at (x,y,maxWidth,1),
// not exceeding that box. "align" is one of AlignLeft, AlignCenter, or
// AlignRight. The screen's background color will not be changed.
//
// You can change the colors and text styles mid-text by inserting a style tag.
// See the package description for details.
//
// Returns the number of actual bytes of the text printed (including style tags)
// and the actual width used for the printed runes.
func Print_2(screen tcell.Screen, text string, x, y, maxWidth, align int, color tcell.Color) (int, int) {
	start, end, width := printWithStyle(screen, text, x, y, 0, maxWidth, align, tcell.StyleDefault.Foreground(color), true)
	return end - start, width
}

// printWithStyle works like [Print_2] but it takes a style instead of just a
// foreground color. The skipWidth parameter specifies the number of cells
// skipped at the beginning of the text. It returns the start index, end index
// (exclusively), and screen width of the text actually printed. If
// maintainBackground is "true", the existing screen background is not changed
// (i.e. the style's background color is ignored).
func printWithStyle(screen tcell.Screen, text string, x, y, skipWidth, maxWidth, align int, style tcell.Style, maintainBackground bool) (start, end, printedWidth int) {
	totalWidth, totalHeight := screen.Size()
	if maxWidth <= 0 || len(text) == 0 || y < 0 || y >= totalHeight {
		return 0, 0, 0
	}

	// If we don't overwrite the background, we use the default color.
	if maintainBackground {
		style = style.Background(tcell.ColorDefault)
	}

	// Skip beginning and measure width.
	var textWidth int
	state := &stepState{
		unisegState: -1,
		style:       style,
	}
	newState := *state
	str := text
	for len(str) > 0 {
		_, str, state = step(str, state, stepOptionsStyle)
		if skipWidth > 0 {
			skipWidth -= state.Width()
			text = str
			newState = *state
			start += state.GrossLength()
		} else {
			textWidth += state.Width()
		}
	}
	state = &newState

	// Reduce all alignments to AlignLeft.
	if align == AlignRight {
		// Chop off characters on the left until it fits.
		for len(text) > 0 && textWidth > maxWidth {
			_, text, state = step(text, state, stepOptionsStyle)
			textWidth -= state.Width()
			start += state.GrossLength()
		}
		x, maxWidth = x+maxWidth-textWidth, textWidth
	} else if align == AlignCenter {
		// Chop off characters on the left until it fits.
		subtracted := (textWidth - maxWidth) / 2
		for len(text) > 0 && subtracted > 0 {
			_, text, state = step(text, state, stepOptionsStyle)
			subtracted -= state.Width()
			textWidth -= state.Width()
			start += state.GrossLength()
		}
		if textWidth < maxWidth {
			x, maxWidth = x+maxWidth/2-textWidth/2, textWidth
		}
	}

	// Draw left-aligned text.
	end = start
	rightBorder := x + maxWidth
	for len(text) > 0 && x < rightBorder && x < totalWidth {
		var c string
		c, text, state = step(text, state, stepOptionsStyle)
		if c == "" {
			break // We don't care about the style at the end.
		}
		width := state.Width()

		if width > 0 {
			finalStyle := state.Style()
			if maintainBackground {
				_, backgroundColor, _ := finalStyle.Decompose()
				if backgroundColor == tcell.ColorDefault {
					_, _, existingStyle, _ := screen.GetContent(x, y)
					_, background, _ := existingStyle.Decompose()
					finalStyle = finalStyle.Background(background)
				}
			}
			for offset := width - 1; offset >= 0; offset-- {
				// To avoid undesired effects, we populate all cells.
				runes := []rune(c)
				if offset == 0 {
					screen.SetContent(x+offset, y, runes[0], runes[1:], finalStyle)
				} else {
					screen.SetContent(x+offset, y, ' ', nil, finalStyle)
				}
			}
		}

		x += width
		end += state.GrossLength()
		printedWidth += width
	}

	return
}

// PrintSimple prints white text to the screen at the given position.
func PrintSimple(screen tcell.Screen, text string, x, y int) {
	Print_2(screen, text, x, y, math.MaxInt32, AlignLeft, Styles.PrimaryTextColor)
}

// Theme defines the colors used when primitives are initialized.
type Theme struct {
	PrimitiveBackgroundColor    tcell.Color // Main background color for primitives.
	ContrastBackgroundColor     tcell.Color // Background color for contrasting elements.
	MoreContrastBackgroundColor tcell.Color // Background color for even more contrasting elements.
	BorderColor                 tcell.Color // Box borders.
	TitleColor                  tcell.Color // Box titles.
	GraphicsColor               tcell.Color // Graphics.
	PrimaryTextColor            tcell.Color // Primary text.
	SecondaryTextColor          tcell.Color // Secondary text (e.g. labels).
	TertiaryTextColor           tcell.Color // Tertiary text (e.g. subtitles, notes).
	InverseTextColor            tcell.Color // Text on primary-colored backgrounds.
	ContrastSecondaryTextColor  tcell.Color // Secondary text on ContrastBackgroundColor-colored backgrounds.
}

// Styles defines the theme for applications. The default is for a black
// background and some basic colors: black, white, yellow, green, cyan, and
// blue.
var Styles = Theme{
	PrimitiveBackgroundColor:    tcell.ColorBlack,
	ContrastBackgroundColor:     tcell.ColorBlue,
	MoreContrastBackgroundColor: tcell.ColorGreen,
	BorderColor:                 tcell.ColorWhite,
	TitleColor:                  tcell.ColorWhite,
	GraphicsColor:               tcell.ColorWhite,
	PrimaryTextColor:            tcell.ColorWhite,
	SecondaryTextColor:          tcell.ColorYellow,
	TertiaryTextColor:           tcell.ColorGreen,
	InverseTextColor:            tcell.ColorBlue,
	ContrastSecondaryTextColor:  tcell.ColorNavy,
}

// Borders defines various borders used when primitives are drawn.
// These may be changed to accommodate a different look and feel.
var Borders = struct {
	Horizontal  rune
	Vertical    rune
	TopLeft     rune
	TopRight    rune
	BottomLeft  rune
	BottomRight rune

	LeftT   rune
	RightT  rune
	TopT    rune
	BottomT rune
	Cross   rune

	HorizontalFocus  rune
	VerticalFocus    rune
	TopLeftFocus     rune
	TopRightFocus    rune
	BottomLeftFocus  rune
	BottomRightFocus rune
}{
	Horizontal:  BoxDrawingsLightHorizontal,
	Vertical:    BoxDrawingsLightVertical,
	TopLeft:     BoxDrawingsLightDownAndRight,
	TopRight:    BoxDrawingsLightDownAndLeft,
	BottomLeft:  BoxDrawingsLightUpAndRight,
	BottomRight: BoxDrawingsLightUpAndLeft,

	LeftT:   BoxDrawingsLightVerticalAndRight,
	RightT:  BoxDrawingsLightVerticalAndLeft,
	TopT:    BoxDrawingsLightDownAndHorizontal,
	BottomT: BoxDrawingsLightUpAndHorizontal,
	Cross:   BoxDrawingsLightVerticalAndHorizontal,

	HorizontalFocus:  BoxDrawingsDoubleHorizontal,
	VerticalFocus:    BoxDrawingsDoubleVertical,
	TopLeftFocus:     BoxDrawingsDoubleDownAndRight,
	TopRightFocus:    BoxDrawingsDoubleDownAndLeft,
	BottomLeftFocus:  BoxDrawingsDoubleUpAndRight,
	BottomRightFocus: BoxDrawingsDoubleUpAndLeft,
}

// Semigraphics provides an easy way to access unicode characters for drawing.
//
// Named like the unicode characters, 'Semigraphics'-prefix used if unicode block
// isn't prefixed itself.
const (
	// Block: General Punctuation U+2000-U+206F (http://unicode.org/charts/PDF/U2000.pdf)
	SemigraphicsHorizontalEllipsis rune = '\u2026' // …

	// Block: Box Drawing U+2500-U+257F (http://unicode.org/charts/PDF/U2500.pdf)
	BoxDrawingsLightHorizontal                    rune = '\u2500' // ─
	BoxDrawingsHeavyHorizontal                    rune = '\u2501' // ━
	BoxDrawingsLightVertical                      rune = '\u2502' // │
	BoxDrawingsHeavyVertical                      rune = '\u2503' // ┃
	BoxDrawingsLightTripleDashHorizontal          rune = '\u2504' // ┄
	BoxDrawingsHeavyTripleDashHorizontal          rune = '\u2505' // ┅
	BoxDrawingsLightTripleDashVertical            rune = '\u2506' // ┆
	BoxDrawingsHeavyTripleDashVertical            rune = '\u2507' // ┇
	BoxDrawingsLightQuadrupleDashHorizontal       rune = '\u2508' // ┈
	BoxDrawingsHeavyQuadrupleDashHorizontal       rune = '\u2509' // ┉
	BoxDrawingsLightQuadrupleDashVertical         rune = '\u250a' // ┊
	BoxDrawingsHeavyQuadrupleDashVertical         rune = '\u250b' // ┋
	BoxDrawingsLightDownAndRight                  rune = '\u250c' // ┌
	BoxDrawingsDownLightAndRightHeavy             rune = '\u250d' // ┍
	BoxDrawingsDownHeavyAndRightLight             rune = '\u250e' // ┎
	BoxDrawingsHeavyDownAndRight                  rune = '\u250f' // ┏
	BoxDrawingsLightDownAndLeft                   rune = '\u2510' // ┐
	BoxDrawingsDownLightAndLeftHeavy              rune = '\u2511' // ┑
	BoxDrawingsDownHeavyAndLeftLight              rune = '\u2512' // ┒
	BoxDrawingsHeavyDownAndLeft                   rune = '\u2513' // ┓
	BoxDrawingsLightUpAndRight                    rune = '\u2514' // └
	BoxDrawingsUpLightAndRightHeavy               rune = '\u2515' // ┕
	BoxDrawingsUpHeavyAndRightLight               rune = '\u2516' // ┖
	BoxDrawingsHeavyUpAndRight                    rune = '\u2517' // ┗
	BoxDrawingsLightUpAndLeft                     rune = '\u2518' // ┘
	BoxDrawingsUpLightAndLeftHeavy                rune = '\u2519' // ┙
	BoxDrawingsUpHeavyAndLeftLight                rune = '\u251a' // ┚
	BoxDrawingsHeavyUpAndLeft                     rune = '\u251b' // ┛
	BoxDrawingsLightVerticalAndRight              rune = '\u251c' // ├
	BoxDrawingsVerticalLightAndRightHeavy         rune = '\u251d' // ┝
	BoxDrawingsUpHeavyAndRightDownLight           rune = '\u251e' // ┞
	BoxDrawingsDownHeavyAndRightUpLight           rune = '\u251f' // ┟
	BoxDrawingsVerticalHeavyAndRightLight         rune = '\u2520' // ┠
	BoxDrawingsDownLightAndRightUpHeavy           rune = '\u2521' // ┡
	BoxDrawingsUpLightAndRightDownHeavy           rune = '\u2522' // ┢
	BoxDrawingsHeavyVerticalAndRight              rune = '\u2523' // ┣
	BoxDrawingsLightVerticalAndLeft               rune = '\u2524' // ┤
	BoxDrawingsVerticalLightAndLeftHeavy          rune = '\u2525' // ┥
	BoxDrawingsUpHeavyAndLeftDownLight            rune = '\u2526' // ┦
	BoxDrawingsDownHeavyAndLeftUpLight            rune = '\u2527' // ┧
	BoxDrawingsVerticalHeavyAndLeftLight          rune = '\u2528' // ┨
	BoxDrawingsDownLightAndLeftUpHeavy            rune = '\u2529' // ┨
	BoxDrawingsUpLightAndLeftDownHeavy            rune = '\u252a' // ┪
	BoxDrawingsHeavyVerticalAndLeft               rune = '\u252b' // ┫
	BoxDrawingsLightDownAndHorizontal             rune = '\u252c' // ┬
	BoxDrawingsLeftHeavyAndRightDownLight         rune = '\u252d' // ┭
	BoxDrawingsRightHeavyAndLeftDownLight         rune = '\u252e' // ┮
	BoxDrawingsDownLightAndHorizontalHeavy        rune = '\u252f' // ┯
	BoxDrawingsDownHeavyAndHorizontalLight        rune = '\u2530' // ┰
	BoxDrawingsRightLightAndLeftDownHeavy         rune = '\u2531' // ┱
	BoxDrawingsLeftLightAndRightDownHeavy         rune = '\u2532' // ┲
	BoxDrawingsHeavyDownAndHorizontal             rune = '\u2533' // ┳
	BoxDrawingsLightUpAndHorizontal               rune = '\u2534' // ┴
	BoxDrawingsLeftHeavyAndRightUpLight           rune = '\u2535' // ┵
	BoxDrawingsRightHeavyAndLeftUpLight           rune = '\u2536' // ┶
	BoxDrawingsUpLightAndHorizontalHeavy          rune = '\u2537' // ┷
	BoxDrawingsUpHeavyAndHorizontalLight          rune = '\u2538' // ┸
	BoxDrawingsRightLightAndLeftUpHeavy           rune = '\u2539' // ┹
	BoxDrawingsLeftLightAndRightUpHeavy           rune = '\u253a' // ┺
	BoxDrawingsHeavyUpAndHorizontal               rune = '\u253b' // ┻
	BoxDrawingsLightVerticalAndHorizontal         rune = '\u253c' // ┼
	BoxDrawingsLeftHeavyAndRightVerticalLight     rune = '\u253d' // ┽
	BoxDrawingsRightHeavyAndLeftVerticalLight     rune = '\u253e' // ┾
	BoxDrawingsVerticalLightAndHorizontalHeavy    rune = '\u253f' // ┿
	BoxDrawingsUpHeavyAndDownHorizontalLight      rune = '\u2540' // ╀
	BoxDrawingsDownHeavyAndUpHorizontalLight      rune = '\u2541' // ╁
	BoxDrawingsVerticalHeavyAndHorizontalLight    rune = '\u2542' // ╂
	BoxDrawingsLeftUpHeavyAndRightDownLight       rune = '\u2543' // ╃
	BoxDrawingsRightUpHeavyAndLeftDownLight       rune = '\u2544' // ╄
	BoxDrawingsLeftDownHeavyAndRightUpLight       rune = '\u2545' // ╅
	BoxDrawingsRightDownHeavyAndLeftUpLight       rune = '\u2546' // ╆
	BoxDrawingsDownLightAndUpHorizontalHeavy      rune = '\u2547' // ╇
	BoxDrawingsUpLightAndDownHorizontalHeavy      rune = '\u2548' // ╈
	BoxDrawingsRightLightAndLeftVerticalHeavy     rune = '\u2549' // ╉
	BoxDrawingsLeftLightAndRightVerticalHeavy     rune = '\u254a' // ╊
	BoxDrawingsHeavyVerticalAndHorizontal         rune = '\u254b' // ╋
	BoxDrawingsLightDoubleDashHorizontal          rune = '\u254c' // ╌
	BoxDrawingsHeavyDoubleDashHorizontal          rune = '\u254d' // ╍
	BoxDrawingsLightDoubleDashVertical            rune = '\u254e' // ╎
	BoxDrawingsHeavyDoubleDashVertical            rune = '\u254f' // ╏
	BoxDrawingsDoubleHorizontal                   rune = '\u2550' // ═
	BoxDrawingsDoubleVertical                     rune = '\u2551' // ║
	BoxDrawingsDownSingleAndRightDouble           rune = '\u2552' // ╒
	BoxDrawingsDownDoubleAndRightSingle           rune = '\u2553' // ╓
	BoxDrawingsDoubleDownAndRight                 rune = '\u2554' // ╔
	BoxDrawingsDownSingleAndLeftDouble            rune = '\u2555' // ╕
	BoxDrawingsDownDoubleAndLeftSingle            rune = '\u2556' // ╖
	BoxDrawingsDoubleDownAndLeft                  rune = '\u2557' // ╗
	BoxDrawingsUpSingleAndRightDouble             rune = '\u2558' // ╘
	BoxDrawingsUpDoubleAndRightSingle             rune = '\u2559' // ╙
	BoxDrawingsDoubleUpAndRight                   rune = '\u255a' // ╚
	BoxDrawingsUpSingleAndLeftDouble              rune = '\u255b' // ╛
	BoxDrawingsUpDoubleAndLeftSingle              rune = '\u255c' // ╜
	BoxDrawingsDoubleUpAndLeft                    rune = '\u255d' // ╝
	BoxDrawingsVerticalSingleAndRightDouble       rune = '\u255e' // ╞
	BoxDrawingsVerticalDoubleAndRightSingle       rune = '\u255f' // ╟
	BoxDrawingsDoubleVerticalAndRight             rune = '\u2560' // ╠
	BoxDrawingsVerticalSingleAndLeftDouble        rune = '\u2561' // ╡
	BoxDrawingsVerticalDoubleAndLeftSingle        rune = '\u2562' // ╢
	BoxDrawingsDoubleVerticalAndLeft              rune = '\u2563' // ╣
	BoxDrawingsDownSingleAndHorizontalDouble      rune = '\u2564' // ╤
	BoxDrawingsDownDoubleAndHorizontalSingle      rune = '\u2565' // ╥
	BoxDrawingsDoubleDownAndHorizontal            rune = '\u2566' // ╦
	BoxDrawingsUpSingleAndHorizontalDouble        rune = '\u2567' // ╧
	BoxDrawingsUpDoubleAndHorizontalSingle        rune = '\u2568' // ╨
	BoxDrawingsDoubleUpAndHorizontal              rune = '\u2569' // ╩
	BoxDrawingsVerticalSingleAndHorizontalDouble  rune = '\u256a' // ╪
	BoxDrawingsVerticalDoubleAndHorizontalSingle  rune = '\u256b' // ╫
	BoxDrawingsDoubleVerticalAndHorizontal        rune = '\u256c' // ╬
	BoxDrawingsLightArcDownAndRight               rune = '\u256d' // ╭
	BoxDrawingsLightArcDownAndLeft                rune = '\u256e' // ╮
	BoxDrawingsLightArcUpAndLeft                  rune = '\u256f' // ╯
	BoxDrawingsLightArcUpAndRight                 rune = '\u2570' // ╰
	BoxDrawingsLightDiagonalUpperRightToLowerLeft rune = '\u2571' // ╱
	BoxDrawingsLightDiagonalUpperLeftToLowerRight rune = '\u2572' // ╲
	BoxDrawingsLightDiagonalCross                 rune = '\u2573' // ╳
	BoxDrawingsLightLeft                          rune = '\u2574' // ╴
	BoxDrawingsLightUp                            rune = '\u2575' // ╵
	BoxDrawingsLightRight                         rune = '\u2576' // ╶
	BoxDrawingsLightDown                          rune = '\u2577' // ╷
	BoxDrawingsHeavyLeft                          rune = '\u2578' // ╸
	BoxDrawingsHeavyUp                            rune = '\u2579' // ╹
	BoxDrawingsHeavyRight                         rune = '\u257a' // ╺
	BoxDrawingsHeavyDown                          rune = '\u257b' // ╻
	BoxDrawingsLightLeftAndHeavyRight             rune = '\u257c' // ╼
	BoxDrawingsLightUpAndHeavyDown                rune = '\u257d' // ╽
	BoxDrawingsHeavyLeftAndLightRight             rune = '\u257e' // ╾
	BoxDrawingsHeavyUpAndLightDown                rune = '\u257f' // ╿

	// Block Elements.
	BlockUpperHalfBlock                              rune = '\u2580' // ▀
	BlockLowerOneEighthBlock                         rune = '\u2581' // ▁
	BlockLowerOneQuarterBlock                        rune = '\u2582' // ▂
	BlockLowerThreeEighthsBlock                      rune = '\u2583' // ▃
	BlockLowerHalfBlock                              rune = '\u2584' // ▄
	BlockLowerFiveEighthsBlock                       rune = '\u2585' // ▅
	BlockLowerThreeQuartersBlock                     rune = '\u2586' // ▆
	BlockLowerSevenEighthsBlock                      rune = '\u2587' // ▇
	BlockFullBlock                                   rune = '\u2588' // █
	BlockLeftSevenEighthsBlock                       rune = '\u2589' // ▉
	BlockLeftThreeQuartersBlock                      rune = '\u258A' // ▊
	BlockLeftFiveEighthsBlock                        rune = '\u258B' // ▋
	BlockLeftHalfBlock                               rune = '\u258C' // ▌
	BlockLeftThreeEighthsBlock                       rune = '\u258D' // ▍
	BlockLeftOneQuarterBlock                         rune = '\u258E' // ▎
	BlockLeftOneEighthBlock                          rune = '\u258F' // ▏
	BlockRightHalfBlock                              rune = '\u2590' // ▐
	BlockLightShade                                  rune = '\u2591' // ░
	BlockMediumShade                                 rune = '\u2592' // ▒
	BlockDarkShade                                   rune = '\u2593' // ▓
	BlockUpperOneEighthBlock                         rune = '\u2594' // ▔
	BlockRightOneEighthBlock                         rune = '\u2595' // ▕
	BlockQuadrantLowerLeft                           rune = '\u2596' // ▖
	BlockQuadrantLowerRight                          rune = '\u2597' // ▗
	BlockQuadrantUpperLeft                           rune = '\u2598' // ▘
	BlockQuadrantUpperLeftAndLowerLeftAndLowerRight  rune = '\u2599' // ▙
	BlockQuadrantUpperLeftAndLowerRight              rune = '\u259A' // ▚
	BlockQuadrantUpperLeftAndUpperRightAndLowerLeft  rune = '\u259B' // ▛
	BlockQuadrantUpperLeftAndUpperRightAndLowerRight rune = '\u259C' // ▜
	BlockQuadrantUpperRight                          rune = '\u259D' // ▝
	BlockQuadrantUpperRightAndLowerLeft              rune = '\u259E' // ▞
	BlockQuadrantUpperRightAndLowerLeftAndLowerRight rune = '\u259F' // ▟
)

// MouseAction indicates one of the actions the mouse is logically doing.
type MouseAction int16

// Available mouse actions.
const (
	MouseMove MouseAction = iota
	MouseLeftDown
	MouseLeftUp
	MouseLeftClick
	MouseLeftDoubleClick
	MouseMiddleDown
	MouseMiddleUp
	MouseMiddleClick
	MouseMiddleDoubleClick
	MouseRightDown
	MouseRightUp
	MouseRightClick
	MouseRightDoubleClick
	MouseScrollUp
	MouseScrollDown
	MouseScrollLeft
	MouseScrollRight

	// The following special value will not be provided as a mouse action but
	// indicate that an overridden mouse event was consumed. See
	// [Box.SetMouseCapture] for details.
	MouseConsumed
)

func getImageLayerAsBlockElements2(sourceImageData image.Image, imageStyle types.ImageStyleEntryType, widthInCharacters int, heightInCharacters int, blurSigma float64) types.LayerEntryType {
	img := NewImage()
	img.SetImage(sourceImageData)
	img.SetSize(heightInCharacters, widthInCharacters)
	img.SetColors(TrueColor)
	img.render()

	layerEntry := types.NewLayerEntry("", "", img.lastWidth, img.lastHeight)
	for y := 0; y < img.lastHeight; y++ {
		for x := 0; x < img.lastWidth; x++ {
			pixel := img.pixels[y*img.lastWidth+x]
			fg, bg, _ := pixel.style.Decompose()
			r, g, b := fg.RGB()
			layerEntry.CharacterMemory[y][x].AttributeEntry.ForegroundColor = GetRGBColor(r, g, b)
			r, g, b = bg.RGB()
			layerEntry.CharacterMemory[y][x].AttributeEntry.BackgroundColor = GetRGBColor(r, g, b)
			layerEntry.CharacterMemory[y][x].Character = pixel.element
		}
	}
	return layerEntry
}
