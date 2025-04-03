package consolizer

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math"
	"os"
	"sort"
	"supercom32.net/consolizer/constants"
	"supercom32.net/consolizer/internal/memory"
	"supercom32.net/consolizer/types"
)

type ImageComposerEntryType struct {
	images             map[string]*types.ImageComposerImageEntryType
	imageStyle         types.ImageStyleEntryType
	widthInCharacters  int
	heightInCharacters int
}

var ImageComposer ImageComposerEntryType

func (shared *ImageComposerEntryType) New() ImageComposerEntryType {
	var newImageComposer ImageComposerEntryType
	newImageComposer.images = make(map[string]*types.ImageComposerImageEntryType)
	return newImageComposer
}

// Specifying 0 widht and height means use native image size
func (shared *ImageComposerEntryType) Add(fileName string, xLocation int, yLocation int, width int, height int, imageStyle types.ImageStyleEntryType, effectStyle constants.EffectStyle, effectStep float64, zOrder int, alphaValue float32) *types.ImageComposerImageEntryType {
	imageComposerImage := types.NewImageComposerImageEntry()
	imageComposerImage.ZOrder = zOrder
	imageComposerImage.XLocation = xLocation
	imageComposerImage.YLocation = yLocation
	imageComposerImage.Width = width
	imageComposerImage.Height = height
	imageComposerImage.ImageStyle = imageStyle
	imageComposerImage.EffectStyle = effectStyle
	imageComposerImage.EffectStep = effectStep
	imageComposerImage.IsVisible = true
	imageComposerImage.AlphaValue = alphaValue
	imageEntry, err := getImage(fileName)
	if err != nil {
		// Here we perform this action anyway to trigger a panic, so we don't need to duplicate the panic code.
		memory.GetImage(fileName)
	}
	imageComposerImage.ImageData = imageEntry.ImageData
	shared.images[fileName] = &imageComposerImage
	return &imageComposerImage
}

func (shared *ImageComposerEntryType) Delete(imageAlias string) {
	delete(shared.images, imageAlias)
}

func (shared *ImageComposerEntryType) RenderImage() image.Image {
	sortedImageList := sortImagesByZOrder(shared.images)
	var baseImage image.Image
	for _, imageEntry := range sortedImageList {
		if imageEntry.IsVisible {
			imageBounds := imageEntry.ImageData.Bounds()
			width := imageBounds.Max.X
			height := imageBounds.Max.Y
			transformedImageData := getTransformedImage(imageEntry.ImageData, *imageEntry)
			if imageEntry.Width > 0 && imageEntry.Height > 0 {
				transformedImageData = resizeImage(transformedImageData, uint(imageEntry.Width), uint(imageEntry.Height), imageEntry.ImageStyle.IsWidthAspectRatioPreserved, imageEntry.ImageStyle.IsHeightAspectRatioPreserved)
			}
			if baseImage != nil {
				baseImage = OverlayImageWithAlpha(transformedImageData, 0, 0, width, height, baseImage, imageEntry.XLocation, imageEntry.YLocation, imageEntry.AlphaValue)
			} else {
				baseImage = transformedImageData
			}
		}
	}
	return baseImage
}

func getTransformedImage(sourceImageData image.Image, imageComposerImageEntry types.ImageComposerImageEntryType) image.Image {
	transformedImage := sourceImageData
	if imageComposerImageEntry.EffectStyle == constants.EffectSinWave {
		transformedImage = applyWaveEffect(sourceImageData, float64(imageComposerImageEntry.EffectStep), false)
	} else if imageComposerImageEntry.EffectStyle == constants.EffectConcentricCircles {
		transformedImage = applyRippleEffect(sourceImageData, float64(imageComposerImageEntry.EffectStep))
	} else if imageComposerImageEntry.EffectStyle == constants.EffectFlagWave {
		transformedImage = applyFlagWavingEffect(sourceImageData, float64(imageComposerImageEntry.EffectStep), 100)
	} else if imageComposerImageEntry.EffectStyle == constants.EffectBlinds {
		transformedImage = applyBlindsEffect(sourceImageData, int(imageComposerImageEntry.EffectStep))
	} else if imageComposerImageEntry.EffectStyle == constants.EffectHorizontalWeaveTransition {
		transformedImage = applyHorizontalWeavingEffect(sourceImageData, int(imageComposerImageEntry.EffectStep))
	} else if imageComposerImageEntry.EffectStyle == constants.EffectVerticalWeaveTransition {
		transformedImage = applyVerticalWeavingEffect(sourceImageData, int(imageComposerImageEntry.EffectStep))
	} else if imageComposerImageEntry.EffectStyle == constants.EffectForwardDiagonalWeaveTransition {
		transformedImage = applyForwardDiagonalWeavingEffect(sourceImageData, int(imageComposerImageEntry.EffectStep))
	} else if imageComposerImageEntry.EffectStyle == constants.EffectBackwardDiagonalWeaveTransition {
		transformedImage = applyBackwardDiagonalWeavingEffect(sourceImageData, int(imageComposerImageEntry.EffectStep))
	} else if imageComposerImageEntry.EffectStyle == constants.EffectGrowingCircleTransition {
		transformedImage = applyGrowingCircleEffect(sourceImageData, int(imageComposerImageEntry.EffectStep))
	} else if imageComposerImageEntry.EffectStyle == constants.EffectCounterClockwiseSwirlTransition {
		transformedImage = applyCounterClockwiseSwirlEffect(sourceImageData, int(imageComposerImageEntry.EffectStep))
	} else if imageComposerImageEntry.EffectStyle == constants.EffectClockwiseSwirlTransition {
		transformedImage = applyClockwiseSwirlEffect(sourceImageData, int(imageComposerImageEntry.EffectStep))
	} else if imageComposerImageEntry.EffectStyle == constants.EffectVerticalCurtainTransition {
		transformedImage = applyVerticalCurtainEffect(sourceImageData, int(imageComposerImageEntry.EffectStep))
	} else if imageComposerImageEntry.EffectStyle == constants.EffectHorizontalCurtainTransition {
		transformedImage = applyHorizontalCurtainEffect(sourceImageData, int(imageComposerImageEntry.EffectStep))
	} else {
		transformedImage = sourceImageData
	}
	return transformedImage
}

func sortImagesByZOrder(imageList map[string]*types.ImageComposerImageEntryType) []*types.ImageComposerImageEntryType {
	// Create a slice to store the image entries
	var imageSlice []*types.ImageComposerImageEntryType

	// Populate the slice with the image entries
	for _, value := range imageList {
		imageSlice = append(imageSlice, value)
	}

	// Sort the slice based on ZOrder
	sort.Slice(imageSlice, func(i, j int) bool {
		return imageSlice[i].ZOrder < imageSlice[j].ZOrder
	})

	return imageSlice
}

func BrailleFromDots(dot0, dot1, dot2, dot3, dot4, dot5, dot6, dot7 bool) rune {
	blocks := 0
	if dot0 {
		blocks += 1
	}
	if dot1 {
		blocks += 8
	}
	if dot2 {
		blocks += 2
	}
	if dot3 {
		blocks += 16
	}
	if dot4 {
		blocks += 4
	}
	if dot5 {
		blocks += 32
	}
	if dot6 {
		blocks += 64
	}
	if dot7 {
		blocks += 128
	}
	brailleChar := rune(0x2800 + blocks)
	return brailleChar
}

func getImageLayerAsBraille(sourceImageData image.Image, imageStyle types.ImageStyleEntryType, widthInCharacters int, heightInCharacters int, blurSigma float64) types.LayerEntryType {
	if widthInCharacters <= 0 && heightInCharacters <= 0 {
		panic(fmt.Sprintf("The specified width and height of %dx%d for your image is not valid.", widthInCharacters, heightInCharacters))
	}
	calculatedPixelWidth := widthInCharacters * 2
	calculatedPixelHeight := heightInCharacters * 4
	if widthInCharacters == 0 {
		calculatedPixelWidth = (heightInCharacters * 4 * sourceImageData.Bounds().Max.X) / sourceImageData.Bounds().Max.Y
	}
	if heightInCharacters == 0 {
		calculatedPixelHeight = (widthInCharacters * 2 * sourceImageData.Bounds().Max.Y) / sourceImageData.Bounds().Max.X
	}
	processedImageData := resizeImage(sourceImageData, uint(calculatedPixelWidth), uint(calculatedPixelHeight), imageStyle.IsWidthAspectRatioPreserved, imageStyle.IsHeightAspectRatioPreserved)
	if blurSigma > 0 {
		processedImageData = imaging.Blur(processedImageData, blurSigma)
	}
	layerEntry := types.NewLayerEntry("", "", widthInCharacters, heightInCharacters)
	layerEntry.CharacterMemory = getBrailleImageData(processedImageData, imageStyle)
	return layerEntry
}

func ConvertImageToGrayscale(inputImage image.Image) *image.Gray {
	bounds := inputImage.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	// Create a new grayscale image with the same dimensions as the input image.
	grayImage := image.NewGray(bounds)
	// Iterate through each pixel in the input image and convert it to grayscale.
	for yLocation := 0; yLocation < height; yLocation++ {
		for xLocation := 0; xLocation < width; xLocation++ {
			pixel := inputImage.At(xLocation, yLocation)
			grayValue := color.GrayModel.Convert(pixel).(color.Gray)
			grayImage.Set(xLocation, yLocation, grayValue)
		}
	}
	return grayImage
}

func getBrailleImageData(inputImage image.Image, imageStyle types.ImageStyleEntryType) [][]types.CharacterEntryType {
	var monochromeImage image.Image
	var grayscaleImage *image.Gray
	contrastAdjustedImage := inputImage
	if imageStyle.DitheringIntensity != 1 {
		contrastAdjustedImage = adjustContrast(contrastAdjustedImage, imageStyle.DitheringIntensity)
	}
	grayscaleImage = ConvertImageToGrayscale(contrastAdjustedImage)
	if imageStyle.IsHistogramEqualized {
		monochromeImage = HistogramEqualization(grayscaleImage)
	}
	if imageStyle.DitheringStyle == constants.DitheringStyle2x2BayerMatrix {
		monochromeImage = FloydSteinbergDithering2x2(contrastAdjustedImage)
	} else if imageStyle.DitheringStyle == constants.DitheringStyle4x4BayerMatrix {
		monochromeImage = FloydSteinbergDithering4x4(contrastAdjustedImage)
	} else if imageStyle.DitheringStyle == constants.DitheringStyle8x8BayerMatrix {
		monochromeImage = FloydSteinbergDithering8x8(contrastAdjustedImage)
	} else if imageStyle.DitheringStyle == constants.DitheringStyleBasic {
		monochromeImage = FloydSteinbergDitheringBasic(grayscaleImage)
	} else if imageStyle.DitheringStyle == constants.DitheringStyleErrorDiffusion {
		monochromeImage = FloydSteinbergDitheringErrorDiffusion(grayscaleImage)
	} else {
		monochromeImage = FloydSteinbergDithering4x4(contrastAdjustedImage)
	}
	/*contextForImage := gg.NewContextForImage(monochromeImage)
	if err := contextForImage.SavePNG("dithered_output.png"); err != nil {
		panic(err)
	}*/
	bounds := inputImage.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	brailleWidth := width / 2
	brailleHeight := height / 4
	result := make([][]types.CharacterEntryType, brailleHeight)
	for currentIndex := range result {
		result[currentIndex] = make([]types.CharacterEntryType, brailleWidth)
	}
	for yCanvasLocation := 0; yCanvasLocation < brailleHeight; yCanvasLocation++ {
		for xCanvasLocation := 0; xCanvasLocation < brailleWidth; xCanvasLocation++ {
			var dots [8]bool
			var redColor, greenColor, blueColor int32
			var totalRed, totalGreen, totalBlue int32
			numPixels := int32(2 * 4)
			for xBlockLocation := 0; xBlockLocation < 2; xBlockLocation++ {
				for yBlockLocation := 0; yBlockLocation < 4; yBlockLocation++ {
					pixelX := xCanvasLocation*2 + xBlockLocation
					pixelY := yCanvasLocation*4 + yBlockLocation
					pixelColor := inputImage.At(pixelX, pixelY)
					redColor, greenColor, blueColor, _ = get8BitColorComponents(pixelColor)
					totalRed += redColor
					totalGreen += greenColor
					totalBlue += blueColor
					monochromeColor := monochromeImage.At(pixelX, pixelY)
					monochromeRed, monochromeGreen, monochromeBlue, _ := get8BitColorComponents(monochromeColor)
					if monochromeRed == 255 && monochromeGreen == 255 && monochromeBlue == 255 {
						dots[xBlockLocation*4+yBlockLocation] = true
					}
				}
			}
			avgRed := totalRed / numPixels
			avgGreen := totalGreen / numPixels
			avgBlue := totalBlue / numPixels
			brailleCharacter := BrailleFromDots(dots[0], dots[1], dots[2], dots[3], dots[4], dots[5], dots[6], dots[7])
			attributeEntry := types.NewAttributeEntry()
			if !imageStyle.IsGrayscale {
				attributeEntry.ForegroundColor = GetRGBColor(avgRed, avgGreen, avgBlue)
			}
			result[yCanvasLocation][xCanvasLocation] = types.CharacterEntryType{
				Character:      brailleCharacter,
				AttributeEntry: attributeEntry,
			}
		}
	}
	return result
}

func adjustContrast(inputImage image.Image, contrastFactor float64) image.Image {
	bounds := inputImage.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Create a new RGBA image to store the modified image.
	outputImage := image.NewRGBA(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			originalColor := inputImage.At(x, y)
			r, g, b, a := originalColor.RGBA()

			// Adjust contrast by scaling pixel values based on contrastFactor
			newR := uint8(clampColorValue(float64(r>>8) * contrastFactor))
			newG := uint8(clampColorValue(float64(g>>8) * contrastFactor))
			newB := uint8(clampColorValue(float64(b>>8) * contrastFactor))

			// Create a new color with the adjusted intensity values
			newColor := color.RGBA{newR, newG, newB, uint8(a >> 8)}
			outputImage.Set(x, y, newColor)
		}
	}

	return outputImage
}

func clampColorValue(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}
	return value
}

// ==========================================================================================================
// ==========================================================================================================
// ==========================================================================================================

func ReadImage(filePath string) (image.Image, error) {
	// Open the image file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func SaveImageToFile(filePath string, img image.Image) error {
	// Create a new file for writing
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Encode and save the image as JPEG
	if err := jpeg.Encode(file, img, nil); err != nil {
		return err
	}
	return nil
}

func getSubImageFromImage(sourceImage image.Image, xLocation, yLocation, width, height int) image.Image {
	sourceBounds := sourceImage.Bounds()
	widthToCut := width
	heightToCut := height
	if xLocation >= sourceBounds.Max.X || yLocation >= sourceBounds.Max.Y {
		return nil
	}
	if xLocation < 0 {
		xLocation = int(math.Abs(float64(xLocation)))
	}
	if yLocation < 0 {
		yLocation = int(math.Abs(float64(yLocation)))
	}
	if widthToCut+xLocation > sourceBounds.Max.X {
		widthToCut = sourceBounds.Max.X - xLocation
	}
	if heightToCut+yLocation > sourceBounds.Max.Y {
		heightToCut = sourceBounds.Max.Y - yLocation
	}
	if heightToCut <= 0 || widthToCut <= 0 {
		return nil
	}
	resultImage := image.NewRGBA(image.Rect(0, 0, widthToCut, heightToCut))
	draw.Draw(resultImage, resultImage.Bounds(), sourceImage, image.Point{xLocation, yLocation}, draw.Src)
	return resultImage
}

func OverlayImageWithAlpha(sourceImage image.Image, sourceXLocation int, sourceYLocation int, sourceWidth int, sourceHeight int, targetImage image.Image, targetXLocation, targetYLocation int, alphaValue float32) image.Image {
	sourceImage = getSubImageFromImage(sourceImage, sourceXLocation, sourceYLocation, sourceWidth, sourceHeight)
	if sourceImage == nil {
		return targetImage
	}
	// Get the bounds of the source and target images.
	sourceBounds := sourceImage.Bounds()
	targetBounds := targetImage.Bounds()

	// Calculate the visible region of the source image within the target bounds.
	x1 := int(max(float64(targetXLocation), float64(targetBounds.Min.X)))
	y1 := int(max(float64(targetYLocation), float64(targetBounds.Min.Y)))
	x2 := int(min(float64(targetXLocation+sourceBounds.Dx()), float64(targetBounds.Max.X)))
	y2 := int(min(float64(targetYLocation+sourceBounds.Dy()), float64(targetBounds.Max.Y)))

	// Create a new RGBA image for the result.
	resultImage := image.NewRGBA(targetBounds)
	draw.Draw(resultImage, targetBounds, targetImage, image.Point{0, 0}, draw.Src)

	for y := y1; y < y2; y++ {
		for x := x1; x < x2; x++ {
			// Calculate the corresponding coordinates in the source image.
			sourceX := x - targetXLocation
			sourceY := y - targetYLocation

			// Retrieve the colors of the source and target pixels.
			sourceColor := sourceImage.At(sourceX, sourceY)
			targetColor := resultImage.At(x, y)

			// Extract the alpha channel of the source pixel.
			_, _, _, sourceAlpha := sourceColor.RGBA()

			// Calculate the final alpha value based on the source alpha and the specified alphaValue.
			finalAlpha := uint8(float32(sourceAlpha>>8) * alphaValue)

			// Blend the source and target colors using the alpha value.
			blendedColor := blendColors(targetColor, sourceColor, finalAlpha)

			// Set the blended color in the result image.
			resultImage.Set(x, y, blendedColor)
		}
	}
	return resultImage
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func blendColors(bg color.Color, fg color.Color, alpha uint8) color.Color {
	bgR, bgG, bgB, bgA := bg.RGBA()
	fgR, fgG, fgB, fgA := fg.RGBA()
	bgFactor := uint32(0xFF - alpha)
	fgFactor := uint32(alpha)
	blendR := (fgR*fgFactor + bgR*bgFactor) / 0xFF
	blendG := (fgG*fgFactor + bgG*bgFactor) / 0xFF
	blendB := (fgB*fgFactor + bgB*bgFactor) / 0xFF
	blendA := (fgA*fgFactor + bgA*bgFactor) / 0xFF
	return color.RGBA{uint8(blendR >> 8), uint8(blendG >> 8), uint8(blendB >> 8), uint8(blendA >> 8)}
}

func applyWaveEffect(inputImage image.Image, amplitude float64, isHorizontal bool) image.Image {
	bounds := inputImage.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	resultImage := image.NewRGBA(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var newX, newY int

			if isHorizontal {
				// Calculate the offset based on the sine wave for vertical direction.
				offset := int(amplitude * math.Sin(2*math.Pi*float64(x)/100.0))
				newX, newY = x, y+offset
			} else {
				// Calculate the offset based on the sine wave for horizontal direction.
				offset := int(amplitude * math.Sin(2*math.Pi*float64(y)/100.0))
				newX, newY = x+offset, y
			}

			// Clamp the new coordinates to the image boundaries.
			if newX < 0 {
				newX = 0
			} else if newX >= width {
				newX = width - 1
			}
			if newY < 0 {
				newY = 0
			} else if newY >= height {
				newY = height - 1
			}

			// Get the color from the original image and set it in the result image.
			color := inputImage.At(newX, newY)
			resultImage.Set(x, y, color)
		}
	}

	return resultImage
}

func applyRippleEffect(inputImage image.Image, amplitude float64) image.Image {
	bounds := inputImage.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	resultImage := image.NewRGBA(bounds)

	centerX, centerY := width/2, height/2

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Calculate the distance from the current pixel to the center of the image.
			dx := float64(x - centerX)
			dy := float64(y - centerY)
			distance := math.Sqrt(dx*dx + dy*dy)

			// Calculate the displacement based on the distance and amplitude.
			displacement := amplitude * math.Sin(distance/20.0)

			// Calculate the new coordinates.
			newX := int(float64(x) + dx*displacement)
			newY := int(float64(y) + dy*displacement)

			// Clamp the new coordinates to the image boundaries.
			if newX < 0 || newX >= width || newY < 0 || newY >= height {
				continue
			}

			// Get the color from the original image and set it in the result image.
			color := inputImage.At(newX, newY)
			resultImage.Set(x, y, color)
		}
	}

	return resultImage
}

func applyFlagWavingEffect(inputImage image.Image, amplitude float64, frequency float64) image.Image {
	bounds := inputImage.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	resultImage := image.NewRGBA(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var newX, newY int

			// Calculate the offset based on the sine wave for vertical direction with back-and-forth oscillation.
			offset := int(amplitude * math.Sin(2*math.Pi*float64(x)/frequency))

			// Reverse the offset's direction on every other iteration.
			if x/int(frequency)%2 == 0 {
				offset = -offset
			}

			newX, newY = x, y+offset

			// Clamp the new coordinates to the image boundaries.
			if newX < 0 {
				newX = 0
			} else if newX >= width {
				newX = width - 1
			}
			if newY < 0 {
				newY = 0
			} else if newY >= height {
				newY = height - 1
			}

			// Get the color from the original image and set it in the result image.
			color := inputImage.At(newX, newY)
			resultImage.Set(x, y, color)
		}
	}

	return resultImage
}

func applyBlindsEffect(inputImage image.Image, step int) image.Image {
	bounds := inputImage.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	resultImage := image.NewRGBA(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var newX int

			if step != 0 {
				if x%(2*step) < step {
					// Calculate the new X position for lines moving from left to right.
					newX = x - step
				} else {
					// Calculate the new X position for lines moving from right to left.
					newX = x + step
				}
			} else {
				// When step is 0, keep the X position the same.
				newX = x
			}

			// Adjust the new X position to avoid the repeating pattern on the left side.
			if newX < 0 {
				newX = -newX
			} else if newX >= width {
				newX = 2*width - newX - 2
			}

			// Clamp the new X position to image boundaries.
			if newX < 0 {
				newX = 0
			} else if newX >= width {
				newX = width - 1
			}

			// Get the color from the original image and set it in the result image.
			color := inputImage.At(newX, y)
			resultImage.Set(x, y, color)
		}
	}

	return resultImage
}

func applyHorizontalWeavingEffect(inputImage image.Image, step int) image.Image {
	bounds := inputImage.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	resultImage := image.NewRGBA(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var newX int

			// Calculate the new X position for even and odd lines.
			if y%2 == 0 {
				// Even lines come from the left.
				newX = x - step
			} else {
				// Odd lines come from the right.
				newX = x + step
			}

			// Ensure the new X position is within the image boundaries.
			if newX < 0 || newX >= width {
				// Set missing areas to black.
				resultImage.Set(x, y, color.Black)
			} else {
				// Get the color from the original image and set it in the result image.
				color := inputImage.At(newX, y)
				resultImage.Set(x, y, color)
			}
		}
	}

	return resultImage
}

func applyVerticalWeavingEffect(inputImage image.Image, step int) image.Image {
	bounds := inputImage.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	resultImage := image.NewRGBA(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var newY int

			// Calculate the new Y position for even and odd lines.
			if x%2 == 0 {
				// Even lines come from the top.
				newY = y - step
			} else {
				// Odd lines come from the bottom.
				newY = y + step
			}

			// Ensure the new Y position is within the image boundaries.
			if newY < 0 || newY >= height {
				// Set missing areas to black.
				resultImage.Set(x, y, color.Black)
			} else {
				// Get the color from the original image and set it in the result image.
				color := inputImage.At(x, newY)
				resultImage.Set(x, y, color)
			}
		}
	}

	return resultImage
}

func applyForwardDiagonalWeavingEffect(inputImage image.Image, step int) image.Image {
	bounds := inputImage.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	resultImage := image.NewRGBA(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var newX, newY int

			// Calculate the new X and Y positions for even and odd diagonals.
			if x%2 == 0 {
				// Even diagonals come from the top-right.
				newX = x + step
				newY = y - step
			} else {
				// Odd diagonals come from the bottom-left.
				newX = x - step
				newY = y + step
			}

			// Ensure the new X and Y positions are within the image boundaries.
			if newX < 0 || newX >= width || newY < 0 || newY >= height {
				// Set missing areas to black.
				resultImage.Set(x, y, color.Black)
			} else {
				// Get the color from the original image and set it in the result image.
				color := inputImage.At(newX, newY)
				resultImage.Set(x, y, color)
			}
		}
	}

	return resultImage
}

func applyBackwardDiagonalWeavingEffect(inputImage image.Image, step int) image.Image {
	bounds := inputImage.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	resultImage := image.NewRGBA(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var newX, newY int

			// Calculate the new X and Y positions for even and odd diagonals.
			if x%2 == 0 {
				// Even diagonals come from the top-left.
				newX = x - step
				newY = y - step
			} else {
				// Odd diagonals come from the bottom-right.
				newX = x + step
				newY = y + step
			}

			// Ensure the new X and Y positions are within the image boundaries.
			if newX < 0 || newX >= width || newY < 0 || newY >= height {
				// Set missing areas to black.
				resultImage.Set(x, y, color.Black)
			} else {
				// Get the color from the original image and set it in the result image.
				color := inputImage.At(newX, newY)
				resultImage.Set(x, y, color)
			}
		}
	}

	return resultImage
}

func applyCounterClockwiseSwirlEffect(inputImage image.Image, step int) image.Image {
	bounds := inputImage.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	resultImage := image.NewRGBA(bounds)

	centerX, centerY := float64(width)/2, float64(height)/2

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dx := float64(x) - centerX
			dy := float64(y) - centerY

			distance := math.Sqrt(dx*dx + dy*dy)
			angle := math.Atan2(dy, dx)

			// Calculate the swirl effect by modifying the angle based on the step.
			swirlFactor := float64(step) / 10.0
			newAngle := angle + swirlFactor*(distance/100)

			// Calculate the new coordinates after the swirl effect.
			newX := int(centerX + distance*math.Cos(newAngle))
			newY := int(centerY + distance*math.Sin(newAngle))

			if newX >= 0 && newX < width && newY >= 0 && newY < height {
				// Get the color from the original image and set it in the result image.
				color := inputImage.At(newX, newY)
				resultImage.Set(x, y, color)
			}
		}
	}

	return resultImage
}

func applyClockwiseSwirlEffect(inputImage image.Image, step int) image.Image {
	bounds := inputImage.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	resultImage := image.NewRGBA(bounds)

	centerX, centerY := float64(width)/2, float64(height)/2

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dx := float64(x) - centerX
			dy := float64(y) - centerY

			distance := math.Sqrt(dx*dx + dy*dy)
			angle := math.Atan2(dy, dx)

			// Calculate the clockwise swirl effect by modifying the angle based on the step.
			swirlFactor := float64(step) / 10.0
			newAngle := angle - swirlFactor*(distance/100)

			// Calculate the new coordinates after the swirl effect.
			newX := int(centerX + distance*math.Cos(newAngle))
			newY := int(centerY + distance*math.Sin(newAngle))

			if newX >= 0 && newX < width && newY >= 0 && newY < height {
				// Get the color from the original image and set it in the result image.
				color := inputImage.At(newX, newY)
				resultImage.Set(x, y, color)
			}
		}
	}

	return resultImage
}

func applyGrowingCircleEffect(inputImage image.Image, step int) image.Image {
	bounds := inputImage.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	resultImage := image.NewRGBA(bounds)

	centerX, centerY := width/2, height/2

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Calculate the distance from the center of the image.
			dx := float64(x - centerX)
			dy := float64(y - centerY)
			distance := math.Sqrt(dx*dx + dy*dy)

			// Check if the distance is less than the step value (radius of the circle).
			if distance < float64(step) {
				// If within the circle, set the pixel in the result image to be transparent.
				resultImage.Set(x, y, color.RGBA{0, 0, 0, 0})
			} else {
				// If outside the circle, get the color from the original image.
				color := inputImage.At(x, y)
				resultImage.Set(x, y, color)
			}
		}
	}

	return resultImage
}

func applyVerticalCurtainEffect(inputImage image.Image, step int) image.Image {
	bounds := inputImage.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	resultImage := image.NewRGBA(bounds)

	// Calculate the position of the transition center (midpoint).
	transitionCenterX := width / 2

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Calculate the distance from the transition center.
			dx := math.Abs(float64(x - transitionCenterX))

			// Check if the distance is less than the step value (width of the transition).
			if dx < float64(step) {
				// If within the transition, get the color from the original image.
				color := inputImage.At(x, y)
				resultImage.Set(x, y, color)
			} else {
				// If outside the transition, set the pixel in the result image to be transparent.
				resultImage.Set(x, y, color.RGBA{0, 0, 0, 0})
			}
		}
	}

	return resultImage
}

func applyHorizontalCurtainEffect(inputImage image.Image, step int) image.Image {
	bounds := inputImage.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	resultImage := image.NewRGBA(bounds)

	// Calculate the position of the transition center (midpoint).
	transitionCenterY := height / 2
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Calculate the distance from the transition center along the y-axis.
			dy := float64(y - transitionCenterY)

			// Check if the absolute distance is less than or equal to the step value.
			if math.Abs(dy) <= float64(step) {
				// If within the transition, get the color from the original image.
				color := inputImage.At(x, y)
				resultImage.Set(x, y, color)
			} else {
				// If outside the transition, set the pixel in the result image to be transparent.
				resultImage.Set(x, y, color.RGBA{0, 0, 0, 0})
			}
		}
	}

	return resultImage
}
