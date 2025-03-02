package consolizer

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/internal/memory"
	"github.com/supercom32/consolizer/types"
	"image"
	"image/color"
	"image/png"
	"strings"
)

/*
UnloadImage allows you to remove an image from memory. Since images can
take up a large amount of space, it is recommended to unload images
any time you are done working with them. However, you may consider
retaining images if they are frequently used and repeatedly loading
them would be less effective. In addition, the following information should
be noted:

- If you pass in an image alias that does not exist, then the delete
operation will be ignored.
*/
func UnloadImage(imageAlias string) {
	memory.DeleteImage(imageAlias)
}

/*
LoadImage allows you to load an image into memory without performing
any ansi conversions ahead of time. This takes up more memory for larger images
but allows you to render those images at arbitrary resolutions. For example,
loading a large image to retain detail and dynamically rendering that
image later depending on the available terminal resolution detected.
*/
func LoadImage(imageFile string) error {
	imageEntry, err := getImageEntryFromFileSystem(imageFile)
	if err != nil {
		return err
	}
	memory.AddImage(imageFile, imageEntry)
	return err
}

/*
LoadImagesInBulk allows you to load multiple images into memory at once.
This is useful since it eliminates the need for error checking over each
image as they are loaded. An example use of this method is as follows:

	// Create a new asset list.
	assetList := dosktop.NewAssetList()
	// Add an image file to our asset list, with a filename of 'MyImageFile'
	// and an image alias of 'MyImageAlias'.
	assetList.AddImage("MyImageFile", "MyImageAlias")
	// Load the list of images into memory.
	err := dosktop.LoadImagesInBulk(assetList)

In addition, the following information should be noted:

- This method works by reading in the provided asset list and then calling
'LoadImage' accordingly each time. For more information about the loading
of images, please see 'LoadImage' for more details.

- In the event an error occurs, it will be returned to the user immediately
and further loading will stop.
*/
func LoadImagesInBulk(assetList types.AssetListType) error {
	var err error
	for _, currentAsset := range assetList.ImageList {
		err = LoadImage(currentAsset.FileName)
		if err != nil {
			return err
		}
	}
	for _, currentAsset := range assetList.PreloadedImageList {
		err = LoadPreRenderedImage(currentAsset.FileName, currentAsset.FileAlias, currentAsset.ImageStyle, currentAsset.WidthInCharacters, currentAsset.HeightInCharacters, currentAsset.BlurSigma)
		if err != nil {
			return err
		}
	}
	return err
}

/*
LoadPreRenderedImage allows you to pre-render an image before loading it into
memory. This enables you to save memory by rendering larger images ahead of
time instead of storing the image data for later use. For example, you can
take a large image and pre-render it at a much lower resolution suitable for
the terminal. In addition, the following information should be noted:

- If you load a pre-rendered image, you are not able to draw them dynamically
at various resolutions. The image can only be drawn with the settings specified
at load time.

- If you specify a value of 0 for ether the width or height, then that
dimension will be automatically calculated to a value that best maintain
the images aspect ratio.

- If you specify a value less than or equal to 0 for both the width and
height, a panic will be generated to fail as fast as possible.

- When pre-rendering an image, it should be noted that each text cell assigned
contains a top and bottom pixel. This is done to provide as much resolution as
possible for your image. That means for a pre-rendered image with a size of
10x10 characters, the actual image being rendered would be 10x20 pixels tall.
If the user wishes to maintain proper aspect ratios, they must manually select
a height that appropriately compensates for this effect, or leave the height
value as 0 to have it done automatically.

- The blur sigma controls how much blurring occurs after your image has been
resized. This allows you to soften your image before it is rendered in ansi
so that hard edges are removed. A value of 0.0 means no blurring will occur,
with higher values increasing the blur factor.
*/
func LoadPreRenderedImage(imageFile string, imageAlias string, imageStyle types.ImageStyleEntryType, widthInCharacters int, heightInCharacters int, blurSigma float64) error {
	imageEntry, err := getImageEntryFromFileSystem(imageFile)
	if err != nil {
		return err
	}
	imageEntry.LayerEntry = getImageLayerAsHighColor(imageEntry.ImageData, imageStyle, widthInCharacters, heightInCharacters, blurSigma)
	imageEntry.ImageData = nil
	memory.AddImage(imageAlias, imageEntry)
	return err
}

/*
LoadBase64Image allows you to load a base64 encoded image into memory without
performing any ansi conversions ahead of time. This takes up more memory for
larger images but allows you to render those images at arbitrary resolutions.
For example, loading a large image to retain detail and dynamically rendering
that image later depending on the available terminal resolution detected.
Since base64 encoded images can be stored in strings, they are ideal for
directly embedding them into applications.
*/
func LoadBase64Image(imageDataAsBase64 string, imageAlias string) error {
	imageEntry := types.NewImageEntry()
	imageData, err := getImageFromBase64String(imageDataAsBase64)
	if err != nil {
		return err
	}
	imageEntry.ImageData = imageData
	memory.AddImage(imageAlias, imageEntry)
	return err
}

/*
LoadPreRenderedBase64Image allows you to pre-render an image before loading it
into memory. This enables you to save memory by rendering larger images ahead
of time instead of storing the image data for later use. For example, you can
take a large image and pre-render it at a much lower resolution suitable for
the terminal. Since base64 encoded images can be stored in strings, they are
ideal for directly embedding them into applications. In addition, the following
information should be noted:

- If you load a pre-rendered image, you are not able to draw them dynamically
at various resolutions. The image can only be drawn with the settings specified
at load time.

- If you specify a value of 0 for ether the width or height, then that
dimension will be automatically calculated to a value that best maintain
the images aspect ratio. This is useful since it removes the need to
calculate this manually.

- If you specify a value less than or equal to 0 for both the width and
height, a panic will be generated to fail as fast as possible.

- When pre-rendering an image, it should be noted that each text cell assigned
contains a top and bottom pixel. This is done to provide as much resolution as
possible for your image. That means for a pre-rendered image with a size of
10x10 characters, the actual image being rendered would be 10x20 pixels tall.
If the user wishes to maintain proper aspect ratios, they must manually select
a height that appropriately compensates for this effect, or leave the height
value as 0 to have it done automatically.

- The blur sigma controls how much blurring occurs after your image has been
resized. This allows you to soften your image before it is rendered in ansi
so that hard edges are removed. A value of 0.0 means no blurring will occur,
with higher values increasing the blur factor.
*/
func LoadPreRenderedBase64Image(imageDataAsBase64 string, imageAlias string, imageStyle types.ImageStyleEntryType, widthInCharacters int, heightInCharacters int, blurSigma float64) error {
	imageEntry := types.NewImageEntry()
	imageData, err := getImageFromBase64String(imageDataAsBase64)
	if err != nil {
		return err
	}
	imageEntry.LayerEntry = getImageLayerAsHighColor(imageData, imageStyle, widthInCharacters, heightInCharacters, blurSigma)
	memory.AddImage(imageAlias, imageEntry)
	return err
}

/*
getBase64PngFromImage allows you to covert raw image data into a base64 encoded
string. This is useful for embedding images directly in applications.
*/
func getBase64PngFromImage(imageToConvert image.Image) (string, error) {
	var imageAsBase64 string
	buffer := new(bytes.Buffer)
	err := png.Encode(buffer, imageToConvert)
	if err != nil {
		return imageAsBase64, err
	}
	imageAsBase64 = base64.StdEncoding.EncodeToString(buffer.Bytes())
	return imageAsBase64, err
}

/*
getImageFromBase64String allows you to obtain raw image data from a base64
encoded string. This is useful for when images are embedded  directly into
applications.
*/
func getImageFromBase64String(imageAsBase64 string) (image.Image, error) {
	fileReader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(imageAsBase64))
	imageData, _, err := image.Decode(fileReader)
	return imageData, err
}

/*
resizeImage allows you to resize an image.
*/
func resizeImage(imageData image.Image, width uint, height uint) image.Image {
	return resize.Resize(width, height, imageData, resize.Lanczos3)
	// return resize.Resize(width, height, imageData, resize.NearestNeighbor)
}

/*
getImageLayerAsHighColor allows you to specify an image and convert it into a text layer
suitable for drawing with. In addition, the following information should be
noted:

- If you specify a value of 0 for ether the width or height, then that
dimension will be automatically calculated to a value that best maintain
the images aspect ratio. This is useful since it removes the need to
calculate this manually.

- If you specify a value less than or equal to 0 for both the width and
height, a panic will be generated to fail as fast as possible.

- When pre-rendering an image, it should be noted that each text cell assigned
contains a top and bottom pixel. This is done to provide as much resolution as
possible for your image. That means for a pre-rendered image with a size of
10x10 characters, the actual image being rendered would be 10x20 pixels tall.
If the user wishes to maintain proper aspect ratios, they must manually select
a height that appropriately compensates for this effect, or leave the height
value as 0 to have it done automatically.

- The blur sigma controls how much blurring occurs after your image has been
resized. This allows you to soften your image before it is rendered in ansi
so that hard edges are removed. A value of 0.0 means no blurring will occur,
with higher values increasing the blur factor.
*/
func getImageLayerAsHighColor(sourceImageData image.Image, imageStyle types.ImageStyleEntryType, widthInCharacters int, heightInCharacters int, blurSigma float64) types.LayerEntryType {
	if widthInCharacters <= 0 && heightInCharacters <= 0 {
		panic(fmt.Sprintf("The specified width and height of %dx%d for your image is not valid.", widthInCharacters, heightInCharacters))
	}
	calculatedPixelWidth := widthInCharacters
	calculatedPixelHeight := heightInCharacters * 2
	if widthInCharacters == 0 {
		calculatedPixelWidth = (heightInCharacters * 2 * sourceImageData.Bounds().Max.X) / sourceImageData.Bounds().Max.Y
	}
	if heightInCharacters == 0 {
		calculatedPixelHeight = (widthInCharacters * sourceImageData.Bounds().Max.Y) / sourceImageData.Bounds().Max.X
	}
	processedImageData := resizeImage(sourceImageData, uint(calculatedPixelWidth), uint(calculatedPixelHeight))
	if blurSigma > 0 {
		processedImageData = imaging.Blur(processedImageData, blurSigma)
	}
	if imageStyle.IsGrayscale {
		processedImageData = ConvertImageToGrayscale(processedImageData)
	}
	calculatedCharacterWidth := calculatedPixelWidth
	calculatedCharacterHeight := calculatedPixelHeight / 2
	layerEntry := types.NewLayerEntry("", "", calculatedCharacterWidth, calculatedCharacterHeight)
	currentImageYLocation := 0
	for currentYLocation := 0; currentYLocation < calculatedCharacterHeight; currentYLocation++ {
		for currentXLocation := 0; currentXLocation < calculatedCharacterWidth; currentXLocation++ {
			currentCharacter := layerEntry.CharacterMemory[currentYLocation][currentXLocation]
			currentCharacter.Character = constants.CharBlockUpperHalf
			upperPixel := processedImageData.At(currentXLocation, currentImageYLocation)
			redColorIndex, greenColorIndex, blueColorIndex, firstAlphaIndex := get8BitColorComponents(upperPixel)
			currentCharacter.AttributeEntry.ForegroundColor = GetRGBColor(int32(redColorIndex), int32(greenColorIndex), int32(blueColorIndex))
			if currentImageYLocation < calculatedCharacterHeight*2 {
				lowerPixel := processedImageData.At(currentXLocation, currentImageYLocation+1)
				redColorIndex, greenColorIndex, blueColorIndex, secondAlphaIndex := get8BitColorComponents(lowerPixel)
				currentCharacter.AttributeEntry.BackgroundColor = GetRGBColor(int32(redColorIndex), int32(greenColorIndex), int32(blueColorIndex))
				if firstAlphaIndex <= 150 || secondAlphaIndex <= 150 {
					currentCharacter.Character = constants.NullRune
				}
			}
			layerEntry.CharacterMemory[currentYLocation][currentXLocation] = currentCharacter
		}
		currentImageYLocation += 2
	}
	return layerEntry
}

/*
get8BitColorComponents allows you to get red, green, and blue color components
from a specific color.
*/
func get8BitColorComponents(colorEntry color.Color) (int32, int32, int32, uint32) {
	redIndex, greenIndex, blueIndex, alphaIndex := colorEntry.RGBA()
	return int32(redIndex) / 257, int32(greenIndex) / 257, int32(blueIndex) / 257, alphaIndex / 257
}

/*
DrawImageToLayer allows you to draw a loaded image to the specified layer.
In addition, the following information should be noted:

- If the location specified falls outside the range for the layer, then
only the visible portion of the image will be drawn.

- If you are drawing an image which has already been pre-rendered, then
your width, height, and blur factor will be ignored.

- If you specify a value of 0 for ether the width or height, then that
dimension will be automatically calculated to a value that best maintain
the images aspect ratio. This is useful since it removes the need to
calculate this manually.

- If you specify a value less than or equal to 0 for both the width and
height, a panic will be generated to fail as fast as possible.

- When pre-rendering an image, it should be noted that each text cell assigned
contains a top and bottom pixel. This is done to provide as much resolution as
possible for your image. That means for a pre-rendered image with a size of
10x10 characters, the actual image being rendered would be 10x20 pixels tall.
If the user wishes to maintain proper aspect ratios, they must manually select
a height that appropriately compensates for this effect, or leave the height
value as 0 to have it done automatically.

- The blur sigma controls how much blurring occurs after your image has been
resized. This allows you to soften your image before it is rendered in ansi
so that hard edges are removed. A value of 0.0 means no blurring will occur,
with higher values increasing the blur factor.
*/
func DrawImageToLayer(layerAlias string, imageAlias string, imageStyle types.ImageStyleEntryType, xLocation int, yLocation int, widthInCharacters int, heightInCharacters int, blurSigma float64) {
	imageEntryType := memory.GetImage(imageAlias)
	imageLayer := imageEntryType.LayerEntry
	if memory.Image.Entries[imageAlias].ImageData != nil {
		imageData := memory.Image.Entries[imageAlias].ImageData
		imageLayer = getImageLayer(imageData, imageStyle, widthInCharacters, heightInCharacters, blurSigma)
	}
	drawImageToLayer(layerAlias, imageLayer, xLocation, yLocation)
}

func getImageLayer(sourceImageData image.Image, imageStyle types.ImageStyleEntryType, widthInCharacters int, heightInCharacters int, blurSigma float64) types.LayerEntryType {
	imageLayer := types.NewLayerEntry("", "", widthInCharacters, heightInCharacters)
	if imageStyle.DrawingStyle == constants.ImageStyleHighColor {
		imageLayer = getImageLayerAsHighColor(sourceImageData, imageStyle, widthInCharacters, heightInCharacters, blurSigma)
	} else {
		imageLayer = getImageLayerAsBraille(sourceImageData, imageStyle, widthInCharacters, heightInCharacters, blurSigma)
	}
	return imageLayer
}

/*
drawImageToLayer allows you to draw a loaded image to the specified layer.
*/
func drawImageToLayer(layerAlias string, imageLayer types.LayerEntryType, xLocation int, yLocation int) {
	layerEntry := memory.GetLayer(layerAlias)
	imageLayer.ScreenXLocation = xLocation
	imageLayer.ScreenYLocation = yLocation
	overlayLayers(&imageLayer, layerEntry)
}

func getImage(fileName string) (*types.ImageEntryType, error) {
	var err error
	if !memory.IsImageExists(fileName) {
		err = LoadImage(fileName)
		if err != nil {
			return nil, err
		}
		defer func() {
			UnloadImage(fileName)
		}()
	}
	imageEntry := memory.GetImage(fileName)
	return imageEntry, err
}

func FloydSteinbergDithering2x2(inputImage image.Image) image.Image {
	bounds := inputImage.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	newImage := image.NewGray(image.Rect(0, 0, width, height))
	bayerMatrix := [][]float64{
		{0, 128},
		{192, 64},
	}
	for yLocation := 0; yLocation < height; yLocation++ {
		for xLocation := 0; xLocation < width; xLocation++ {
			oldPixel := inputImage.At(xLocation, yLocation)
			redColor, _, _, _ := oldPixel.RGBA()
			grayValue := uint8(redColor >> 8)
			threshold := bayerMatrix[yLocation%2][xLocation%2]
			if float64(grayValue) >= threshold {
				grayValue = 255
			} else {
				grayValue = 0
			}
			newColor := color.Gray{Y: grayValue}
			newImage.SetGray(xLocation, yLocation, newColor)
		}
	}
	return newImage
}

func FloydSteinbergDithering4x4(inputImage image.Image) image.Image {
	bounds := inputImage.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	newImage := image.NewGray(image.Rect(0, 0, width, height))
	bayerMatrix := [][]float64{
		{0, 128, 32, 160},
		{192, 64, 224, 96},
		{48, 176, 16, 144},
		{240, 112, 208, 80},
	}
	for yLocation := 0; yLocation < height; yLocation++ {
		for xLocation := 0; xLocation < width; xLocation++ {
			oldPixel := inputImage.At(xLocation, yLocation)
			redColor, _, _, _ := oldPixel.RGBA()
			grayValue := uint8(redColor >> 8)
			threshold := bayerMatrix[yLocation%4][xLocation%4]
			if float64(grayValue) >= threshold {
				grayValue = 255
			} else {
				grayValue = 0
			}
			newColor := color.Gray{Y: grayValue}
			newImage.SetGray(xLocation, yLocation, newColor)
		}
	}
	return newImage
}

func FloydSteinbergDithering8x8(inputImage image.Image) image.Image {
	bounds := inputImage.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	newImage := image.NewGray(image.Rect(0, 0, width, height))
	bayerMatrix := [][]float64{
		{0, 48, 12, 60, 3, 51, 15, 63},
		{32, 16, 44, 28, 35, 19, 47, 31},
		{8, 56, 4, 52, 11, 59, 7, 55},
		{40, 24, 36, 20, 43, 27, 39, 23},
		{2, 50, 14, 62, 1, 49, 13, 61},
		{34, 18, 46, 30, 33, 17, 45, 29},
		{10, 58, 6, 54, 9, 57, 5, 53},
		{42, 26, 38, 22, 41, 25, 37, 21},
	}
	for yLocation := 0; yLocation < height; yLocation++ {
		for xLocation := 0; xLocation < width; xLocation++ {
			oldPixel := inputImage.At(xLocation, yLocation)
			redColor, _, _, _ := oldPixel.RGBA()
			grayValue := uint8(redColor >> 8)
			threshold := (bayerMatrix[yLocation%8][xLocation%8] + 0.5) * 255 / 64
			if float64(grayValue) >= threshold {
				grayValue = 255
			} else {
				grayValue = 0
			}
			newColor := color.Gray{grayValue}
			newImage.SetGray(xLocation, yLocation, newColor)
		}
	}
	return newImage
}

func FloydSteinbergDitheringBasic(inputImage *image.Gray) image.Image {
	for yLocation := inputImage.Bounds().Min.Y; yLocation < inputImage.Bounds().Max.Y; yLocation++ {
		for xLocation := inputImage.Bounds().Min.X; xLocation < inputImage.Bounds().Max.X; xLocation++ {
			currentPixel := inputImage.GrayAt(xLocation, yLocation).Y
			newPixel := 0
			if currentPixel > 122 {
				newPixel = 255
			}
			inputImage.Set(xLocation, yLocation, color.Gray{Y: uint8(newPixel)})
			quantizationError := int(currentPixel) - newPixel
			temporaryPixel := float64(inputImage.GrayAt(xLocation+1, yLocation).Y)
			temporaryColor := temporaryPixel + float64(quantizationError*7/16)
			inputImage.Set(xLocation+1, yLocation, color.Gray{Y: uint8(temporaryColor)})
			temporaryPixel = float64(inputImage.GrayAt(xLocation-1, yLocation+1).Y)
			temporaryColor = temporaryPixel + float64(quantizationError*3/16)
			inputImage.Set(xLocation-1, yLocation+1, color.Gray{Y: uint8(temporaryColor)})
			temporaryPixel = float64(inputImage.GrayAt(xLocation, yLocation+1).Y)
			temporaryColor = temporaryPixel + float64(quantizationError*5/16)
			inputImage.Set(xLocation, yLocation+1, color.Gray{Y: uint8(temporaryColor)})
			temporaryPixel = float64(inputImage.GrayAt(xLocation+1, yLocation+1).Y)
			temporaryColor = temporaryPixel + float64(quantizationError*1/16)
			inputImage.Set(xLocation+1, yLocation+1, color.Gray{Y: uint8(temporaryColor)})
		}
	}
	return inputImage
}

func FloydSteinbergDitheringErrorDiffusion(inputImage *image.Gray) image.Image {
	imageBounds := inputImage.Bounds()
	imageWidth, imageHeight := imageBounds.Max.X, imageBounds.Max.Y
	outputImage := image.NewRGBA(imageBounds)
	errorMatrix := make([][]float32, imageHeight)
	for currentIndex := range errorMatrix {
		errorMatrix[currentIndex] = make([]float32, imageWidth)
	}
	for yLocation := 0; yLocation < imageHeight; yLocation++ {
		for xLocation := 0; xLocation < imageWidth; xLocation++ {
			currentPixel := inputImage.GrayAt(xLocation, yLocation)
			threshold := float32(currentPixel.Y) + errorMatrix[yLocation][xLocation]
			var newPixelValue uint8
			if threshold > 127 {
				newPixelValue = 255
			} else {
				newPixelValue = 0
			}
			newColor := color.RGBA{newPixelValue, newPixelValue, newPixelValue, 255}
			outputImage.Set(xLocation, yLocation, newColor)
			quantizationError := float32(currentPixel.Y) - float32(newPixelValue)
			if xLocation+1 < imageWidth {
				errorMatrix[yLocation][xLocation+1] += quantizationError * 7 / 16
			}
			if yLocation+1 < imageHeight {
				if xLocation-1 >= 0 {
					errorMatrix[yLocation+1][xLocation-1] += quantizationError * 3 / 16
				}
				errorMatrix[yLocation+1][xLocation] += quantizationError * 5 / 16
				if xLocation+1 < imageWidth {
					errorMatrix[yLocation+1][xLocation+1] += quantizationError * 1 / 16
				}
			}
		}
	}
	return outputImage
}

/*
HistogramEqualization performs histogram equalization on a grayscale image to enhance
its contrast and improve the overall image quality. This technique redistributes the
intensity values of the image, resulting in a more balanced and visually appealing
representation. In addition, the following information should be noted:

- An alpha value of 1.0 is equal to 100% visible, while an alpha value of
0.0 is 0% visible. Specifying a value outside this range indicates that
you want to over amplify or under amplify the color transparency effect.

- If the percent change specified is outside the RGB color range (for
example, if you specified 200%), then the color will simply bottom or max
out at RGB(0, 0, 0) or RGB(255, 255, 255) respectively.
*/
func HistogramEqualization(inputImage *image.Gray) *image.Gray {
	pixelCount := 256
	histogram := make([]int, pixelCount)
	for yLocation := 0; yLocation < inputImage.Rect.Max.Y; yLocation++ {
		for xLocation := 0; xLocation < inputImage.Rect.Max.X; xLocation++ {
			grayValue := inputImage.GrayAt(xLocation, yLocation).Y
			histogram[grayValue]++
		}
	}
	cumulativeDistributionFunction := make([]int, pixelCount)
	cumulativeDistributionFunction[0] = histogram[0]
	for currentIndex := 1; currentIndex < pixelCount; currentIndex++ {
		cumulativeDistributionFunction[currentIndex] = cumulativeDistributionFunction[currentIndex-1] + histogram[currentIndex]
	}
	minimumCDF := cumulativeDistributionFunction[0]
	for currentIndex := 0; currentIndex < pixelCount; currentIndex++ {
		if cumulativeDistributionFunction[currentIndex] > 0 {
			minimumCDF = cumulativeDistributionFunction[currentIndex]
			break
		}
	}
	width, height := inputImage.Rect.Max.X, inputImage.Rect.Max.Y
	outputImage := image.NewGray(inputImage.Rect)
	for yLocation := 0; yLocation < height; yLocation++ {
		for xLocation := 0; xLocation < width; xLocation++ {
			grayValue := inputImage.GrayAt(xLocation, yLocation).Y
			newGrayValue := uint8(((cumulativeDistributionFunction[grayValue] - minimumCDF) * 255) / (width*height - minimumCDF))
			outputImage.SetGray(xLocation, yLocation, color.Gray{Y: newGrayValue})
		}
	}
	return outputImage
}
