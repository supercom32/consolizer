package consolizer

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/types"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math"
	"os"
	"sort"
)

/*
ImageComposerEntryType is a structure which manages multiple image layers for compositing.
*/
type ImageComposerEntryType struct {
	images             map[string]*types.ImageComposerImageEntryType
	imageStyle         types.ImageStyleEntryType
	widthInCharacters  int
	heightInCharacters int
}

/*
ImageComposer is a variable which is the global instance for managing image composition.
*/
var ImageComposer ImageComposerEntryType

/*
New is a method which allows you to create a new instance of an image composer. In addition, the following should
be noted:

- It initializes the internal image map for tracking multiple image layers.

Example:
    composer := consolizer.ImageComposer.New()
*/
func (shared *ImageComposerEntryType) New() ImageComposerEntryType {
	var newImageComposer ImageComposerEntryType
	newImageComposer.images = make(map[string]*types.ImageComposerImageEntryType)
	return newImageComposer
}

/*
Add is a method which allows you to add a new image to the composer with specific dimensions and styles. In addition, the
following should be noted:

- If width or height are set to 0, the native image size will be used.

Example:
    composer.Add("test.png", 0, 0, 100, 100, style, constants.EffectNone, 0, 1, 1.0)
*/
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
	imageEntry, err := loadImageAndGetEntry(fileName)
	if err != nil {
		safeSttyPanic(fmt.Sprintf("Could not load image '%s': %s", fileName, err.Error()))
	}
	imageComposerImage.ImageData = imageEntry.ImageData
	shared.images[fileName] = &imageComposerImage
	return &imageComposerImage
}

/*
Delete is a method which allows you to remove an image from the composer using its alias.

Example:
    composer.Delete("test.png")
*/
func (shared *ImageComposerEntryType) Delete(imageAlias string) {
	delete(shared.images, imageAlias)
}

/*
Clear is a method which allows you to remove all images currently managed by the composer.

Example:
    composer.Clear()
*/
func (shared *ImageComposerEntryType) Clear() {
	shared.images = make(map[string]*types.ImageComposerImageEntryType)
}

/*
RenderImage is a method which allows you to composite all visible images in the composer into a single image object. In
addition, the following should be noted:

- The images are processed according to their Z-order, effects, and alpha values.

Example:
    finalImage := composer.RenderImage()
*/
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

/*
getTransformedImage is a method which allows you to apply visual effects to a source image based on the provided
composer entry settings.

Example:
    transformed := getTransformedImage(img, entry)
*/
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

/*
sortImagesByZOrder is a method which allows you to sort a map of image entries into a slice ordered by their Z-order.

Example:
    sorted := sortImagesByZOrder(composer.images)
*/
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

/*
BrailleFromDots is a method which allows you to convert a set of eight dots into a corresponding Braille rune.

Example:
    r := BrailleFromDots(true, false, true, false, true, false, true, false)
*/
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

/*
getImageLayerAsBraille is a method which allows you to convert an image into a layer of Braille characters based on the
specified dimensions and style. In addition, the following should be noted:

- The blur sigma controls how much blurring occurs after your image has been resized. This allows you to soften your
  image before it is rendered in ansi so that hard edges are removed. A value of 0.0 means no blurring will occur, with
  higher values increasing the blur factor.

Example:
    layer := getImageLayerAsBraille(img, style, 80, 24, 0.5)
*/
func getImageLayerAsBraille(sourceImageData image.Image, imageStyle types.ImageStyleEntryType, widthInCharacters int, heightInCharacters int, blurSigma float64) types.LayerEntryType {
	if widthInCharacters <= 0 && heightInCharacters <= 0 {
		safeSttyPanic(fmt.Sprintf("The specified width and height of %dx%d for your image is not valid.", widthInCharacters, heightInCharacters))
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

	// Get the braille image data
	brailleImageData := getBrailleImageData(processedImageData, imageStyle)

	// Safely copy the braille image data to the layer entry
	// Only copy data that fits within the bounds of the layer entry
	for y := 0; y < len(brailleImageData) && y < len(layerEntry.CharacterMemory); y++ {
		for x := 0; x < len(brailleImageData[y]) && x < len(layerEntry.CharacterMemory[y]); x++ {
			layerEntry.CharacterMemory[y][x] = brailleImageData[y][x]
		}
	}

	return layerEntry
}

/*
ConvertImageToGrayscale is a method which allows you to convert a color image into a grayscale image while preserving
transparency.

Example:
    gray := ConvertImageToGrayscale(img)
*/
func ConvertImageToGrayscale(inputImage image.Image) image.Image {
	bounds := inputImage.Bounds()
	// Create a new RGBA image to support alpha channel.
	grayImage := image.NewRGBA(bounds)
	// Iterate through each pixel in the input image and convert it to grayscale.
	for yLocation := bounds.Min.Y; yLocation < bounds.Max.Y; yLocation++ {
		for xLocation := bounds.Min.X; xLocation < bounds.Max.X; xLocation++ {
			pixel := inputImage.At(xLocation, yLocation)
			_, _, _, a := pixel.RGBA()
			if a == 0 {
				grayImage.Set(xLocation, yLocation, color.Transparent)
			} else {
				grayValue := color.GrayModel.Convert(pixel).(color.Gray)
				newColor := color.RGBA{R: grayValue.Y, G: grayValue.Y, B: grayValue.Y, A: uint8(a >> 8)}
				grayImage.Set(xLocation, yLocation, newColor)
			}
		}
	}
	return grayImage
}

/*
getBrailleImageData is a method which allows you to convert an image into a 2D grid of character entries representing
Braille dots.

Example:
    data := getBrailleImageData(img, style)
*/
func getBrailleImageData(inputImage image.Image, imageStyle types.ImageStyleEntryType) [][]types.CharacterEntryType {
	var monochromeImage image.Image
	var grayscaleImage *image.Gray
	contrastAdjustedImage := inputImage
	if imageStyle.DitheringIntensity != 1 {
		contrastAdjustedImage = adjustContrast(contrastAdjustedImage, imageStyle.DitheringIntensity)
	}
	grayscaleWithAlpha := ConvertImageToGrayscale(contrastAdjustedImage)
	bounds := grayscaleWithAlpha.Bounds()
	grayscaleImage = image.NewGray(bounds)
	draw.Draw(grayscaleImage, bounds, grayscaleWithAlpha, bounds.Min, draw.Src)

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
	/*
	contextForImage := gg.NewContextForImage(monochromeImage) if err := contextForImage.SavePNG("dithered_output.png"); err != nil { safeSttyPanic(err) }
	*/
	bounds = inputImage.Bounds()
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
			isTransparentBlock := false
			for xBlockLocation := 0; xBlockLocation < 2; xBlockLocation++ {
				for yBlockLocation := 0; yBlockLocation < 4; yBlockLocation++ {
					pixelX := xCanvasLocation*2 + xBlockLocation
					pixelY := yCanvasLocation*4 + yBlockLocation
					pixelColor := inputImage.At(pixelX, pixelY)
					_, _, _, a := pixelColor.RGBA()
					if a == 0 {
						isTransparentBlock = true
						break
					}
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
				if isTransparentBlock {
					break
				}
			}

			if isTransparentBlock {
				result[yCanvasLocation][xCanvasLocation] = types.CharacterEntryType{
					Character:      rune(0),
					AttributeEntry: types.NewAttributeEntry(),
				}
				continue
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

/*
adjustContrast is a method which allows you to modify the contrast of an image by a specified factor. In addition, the
following should be noted:

- The contrast factor is a scaling value where 1.0 results in no change, values greater than 1.0 increase contrast, and
  values less than 1.0 decrease it.

Example:
    brighter := adjustContrast(img, 1.2)
*/
func adjustContrast(inputImage image.Image, contrastFactor float64) image.Image {
	bounds := inputImage.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

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

/*
clampColorValue is a method which allows you to ensure a color intensity value remains within the valid 0-255 range.

Example:
    v := clampColorValue(260.0)
*/
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

/*
ReadImage is a method which allows you to load an image from a file path.

Example:
    img, err := ReadImage("path/to/image.png")
*/
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

/*
SaveImageToFile is a method which allows you to save an image object to a file as a JPEG.

Example:
    err := SaveImageToFile("output.jpg", img)
*/
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

/*
getSubImageFromImage is a method which allows you to extract a rectangular sub-region from a source image.

Example:
    sub := getSubImageFromImage(img, 10, 10, 50, 50)
*/
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

/*
OverlayImageWithAlpha is a method which allows you to overlay one image onto another with a specified alpha transparency
level. In addition, the following should be noted:

- The alpha value can range from 0.0 (fully transparent) to 1.0 (fully opaque).

Example:
    result := OverlayImageWithAlpha(fg, 0, 0, 10, 10, bg, 5, 5, 0.5)
*/
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

/*
max is a method which allows you to determine the maximum of two float64 values.

Example:
    m := max(10.5, 20.1)
*/
func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

/*
min is a method which allows you to determine the minimum of two float64 values.

Example:
    m := min(10.5, 20.1)
*/
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

/*
blendColors is a method which allows you to blend two colors based on an alpha value. In addition, the following should
be noted:

- The alpha value should range from 0 (fully transparent) to 255 (fully opaque).

Example:
    c := blendColors(color.Black, color.White, 128)
*/
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

/*
applyWaveEffect is a method which allows you to apply a sine wave distortion effect to an image. In addition, the
following should be noted:

- The amplitude determines the maximum displacement of the wave effect in pixels.

Example:
    wavy := applyWaveEffect(img, 5.0, true)
*/
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

			// GetLayer the color from the original image and set it in the result image.
			color := inputImage.At(newX, newY)
			resultImage.Set(x, y, color)
		}
	}

	return resultImage
}

/*
applyRippleEffect is a method which allows you to apply a concentric ripple effect to an image. In addition, the
following should be noted:

- The amplitude determines the maximum displacement of the ripple effect in pixels.

Example:
    rippled := applyRippleEffect(img, 10.0)
*/
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

			// GetLayer the color from the original image and set it in the result image.
			color := inputImage.At(newX, newY)
			resultImage.Set(x, y, color)
		}
	}

	return resultImage
}

/*
applyFlagWavingEffect is a method which allows you to apply a flag waving effect to an image. In addition, the
following should be noted:

- The amplitude determines the maximum displacement of the flag waving effect in pixels.

- The frequency determines the length of each wave cycle in pixels.

Example:
    flag := applyFlagWavingEffect(img, 5.0, 100.0)
*/
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

			// GetLayer the color from the original image and set it in the result image.
			color := inputImage.At(newX, newY)
			resultImage.Set(x, y, color)
		}
	}

	return resultImage
}

/*
applyBlindsEffect is a method which allows you to apply a vertical blinds transition effect to an image.

Example:
    blinds := applyBlindsEffect(img, 10)
*/
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

			// GetLayer the color from the original image and set it in the result image.
			color := inputImage.At(newX, y)
			resultImage.Set(x, y, color)
		}
	}

	return resultImage
}

/*
applyHorizontalWeavingEffect is a method which allows you to apply a horizontal weaving transition effect to an image.

Example:
    weaved := applyHorizontalWeavingEffect(img, 5)
*/
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
				// GetLayer the color from the original image and set it in the result image.
				color := inputImage.At(newX, y)
				resultImage.Set(x, y, color)
			}
		}
	}

	return resultImage
}

/*
applyVerticalWeavingEffect is a method which allows you to apply a vertical weaving transition effect to an image.

Example:
    weaved := applyVerticalWeavingEffect(img, 5)
*/
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
				// GetLayer the color from the original image and set it in the result image.
				color := inputImage.At(x, newY)
				resultImage.Set(x, y, color)
			}
		}
	}

	return resultImage
}

/*
applyForwardDiagonalWeavingEffect is a method which allows you to apply a forward diagonal weaving transition effect to
an image.

Example:
    weaved := applyForwardDiagonalWeavingEffect(img, 5)
*/
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
				// GetLayer the color from the original image and set it in the result image.
				color := inputImage.At(newX, newY)
				resultImage.Set(x, y, color)
			}
		}
	}

	return resultImage
}

/*
applyBackwardDiagonalWeavingEffect is a method which allows you to apply a backward diagonal weaving transition effect
to an image.

Example:
    weaved := applyBackwardDiagonalWeavingEffect(img, 5)
*/
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
				// GetLayer the color from the original image and set it in the result image.
				color := inputImage.At(newX, newY)
				resultImage.Set(x, y, color)
			}
		}
	}

	return resultImage
}

/*
applyCounterClockwiseSwirlEffect is a method which allows you to apply a counter-clockwise swirl effect to an image.

Example:
    swirled := applyCounterClockwiseSwirlEffect(img, 5)
*/
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
				// GetLayer the color from the original image and set it in the result image.
				color := inputImage.At(newX, newY)
				resultImage.Set(x, y, color)
			}
		}
	}

	return resultImage
}

/*
applyClockwiseSwirlEffect is a method which allows you to apply a clockwise swirl effect to an image.

Example:
    swirled := applyClockwiseSwirlEffect(img, 5)
*/
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
				// GetLayer the color from the original image and set it in the result image.
				color := inputImage.At(newX, newY)
				resultImage.Set(x, y, color)
			}
		}
	}

	return resultImage
}

/*
applyGrowingCircleEffect is a method which allows you to apply a growing circle transition effect to an image.

Example:
    circle := applyGrowingCircleEffect(img, 20)
*/
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

/*
applyVerticalCurtainEffect is a method which allows you to apply a vertical curtain transition effect to an image.

Example:
    curtain := applyVerticalCurtainEffect(img, 30)
*/
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

/*
applyHorizontalCurtainEffect is a method which allows you to apply a horizontal curtain transition effect to an image.

Example:
    curtain := applyHorizontalCurtainEffect(img, 30)
*/
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
