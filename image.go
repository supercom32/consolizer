package consolizer

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/types"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"math"
	"math/rand"
	"runtime"
	"strings"
	"sync"
	"time"
)

type ImageMemoryType struct {
	sync.Mutex
	Entries map[string]*types.ImageEntryType
}

var Image ImageMemoryType

func init() {
	Image.Entries = make(map[string]*types.ImageEntryType)
}

func addImage(imageAlias string, imageEntry types.ImageEntryType) {
	Image.Lock()
	defer func() {
		Image.Unlock()
	}()
	// verify if any errors occurred?
	Image.Entries[imageAlias] = &imageEntry
}

func getImage(imageAlias string) *types.ImageEntryType {
	Image.Lock()
	defer func() {
		Image.Unlock()
	}()
	if Image.Entries[imageAlias] == nil {
		safeSttyPanic(fmt.Sprintf("The requested Image with alias '%s' could not be returned since it does not exist.", imageAlias))
	}
	return Image.Entries[imageAlias]
}

func deleteImage(imageAlias string) {
	Image.Lock()
	defer func() {
		Image.Unlock()
	}()
	delete(Image.Entries, imageAlias)
}

// IsImageExists allows you to check if an image with the given alias exists in memory.
func IsImageExists(imageAlias string) bool {
	Image.Lock()
	defer func() {
		Image.Unlock()
	}()
	if Image.Entries[imageAlias] == nil {
		return false
	}
	return true
}

// ClearAllImages removes all loaded images from memory.
func ClearAllImages() {
	Image.Lock()
	defer func() {
		Image.Unlock()
	}()
	Image.Entries = make(map[string]*types.ImageEntryType)
}

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
	deleteImage(imageAlias)
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
	addImage(imageFile, imageEntry)
	return err
}

/*
LoadImagesInBulk allows you to load multiple images into memory at once.
This is useful since it eliminates the need for error checking over each
image as they are loaded. An example use of this method is as follows:

	// Create a new asset list.
	assetList := dosktop.NewAssetList()
	// AddLayer an image file to our asset list, with a filename of 'MyImageFile'
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
	addImage(imageAlias, imageEntry)
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
	addImage(imageAlias, imageEntry)
	return err
}

/*
LoadPreRenderedBase64Image allows you to pre-render an image before loading it
into memory. This enables you to save memory by rendering larger images ahead of
time instead of storing the image data for later use. For example, you can
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
	addImage(imageAlias, imageEntry)
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
func resizeImage(imageData image.Image, targetWidth, targetHeight uint, isWidthAspectRatioPreserved bool, isHeightAspectRatioPreserved bool) image.Image {
	if !isWidthAspectRatioPreserved && !isHeightAspectRatioPreserved {
		return resize.Resize(targetWidth, targetHeight, imageData, resize.Lanczos3)
	}
	originalBounds := imageData.Bounds()
	originalWidth := originalBounds.Dx()
	originalHeight := originalBounds.Dy()
	scaleWidth := float64(targetWidth) / float64(originalWidth)
	scaleHeight := float64(targetHeight) / float64(originalHeight)
	var scale float64
	switch {
	case isWidthAspectRatioPreserved && !isHeightAspectRatioPreserved:
		scale = scaleWidth
	case isHeightAspectRatioPreserved && !isWidthAspectRatioPreserved:
		scale = scaleHeight
	default:
		scale = scaleWidth
		if scaleHeight < scaleWidth {
			scale = scaleHeight
		}
	}
	newWidth := uint(float64(originalWidth) * scale)
	newHeight := uint(float64(originalHeight) * scale)
	resizedImage := resize.Resize(newWidth, newHeight, imageData, resize.Lanczos3)
	outputImage := image.NewRGBA(image.Rect(0, 0, int(targetWidth), int(targetHeight)))
	draw.Draw(outputImage, image.Rect(0, 0, int(newWidth), int(newHeight)), resizedImage, image.Point{}, draw.Over)
	return outputImage
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
func getImageLayerAsHighColor(
	sourceImageData image.Image,
	imageStyle types.ImageStyleEntryType,
	widthInCharacters int,
	heightInCharacters int,
	blurSigma float64,
	existingLayer ...*types.LayerEntryType,
) types.LayerEntryType {

	if widthInCharacters <= 0 && heightInCharacters <= 0 {
		safeSttyPanic(fmt.Sprintf("The specified width and height of %dx%d for your image is not valid.", widthInCharacters, heightInCharacters))
	}

	// Calculate pixel dimensions
	calculatedPixelWidth := widthInCharacters
	calculatedPixelHeight := heightInCharacters * 2
	if widthInCharacters == 0 {
		calculatedPixelWidth = (heightInCharacters * 2 * sourceImageData.Bounds().Max.X) / sourceImageData.Bounds().Max.Y
	}
	if heightInCharacters == 0 {
		calculatedPixelHeight = (widthInCharacters * sourceImageData.Bounds().Max.Y) / sourceImageData.Bounds().Max.X
	}

	// Resize and optionally blur / grayscale
	processedImageData := resizeImage(sourceImageData, uint(calculatedPixelWidth), uint(calculatedPixelHeight), imageStyle.IsWidthAspectRatioPreserved, imageStyle.IsHeightAspectRatioPreserved)
	if blurSigma > 0 {
		processedImageData = imaging.Blur(processedImageData, blurSigma)
	}
	if imageStyle.IsGrayscale {
		processedImageData = ConvertImageToGrayscale(processedImageData)
	}

	calculatedCharacterWidth := calculatedPixelWidth
	calculatedCharacterHeight := calculatedPixelHeight / 2

	// Create new layer
	layerEntry := types.NewLayerEntry("", "", calculatedCharacterWidth, calculatedCharacterHeight)

	// Prepare underlying layer (optional)
	var underlyingLayer *types.LayerEntryType
	var underlyingWidth, underlyingHeight int
	if len(existingLayer) > 0 && existingLayer[0] != nil {
		underlyingLayer = existingLayer[0]
		underlyingHeight = len(underlyingLayer.CharacterMemory)
		if underlyingHeight > 0 {
			underlyingWidth = len(underlyingLayer.CharacterMemory[0])
		}
	}

	imageBounds := processedImageData.Bounds()
	currentImageY := 0

	for charY := 0; charY < calculatedCharacterHeight; charY++ {
		for charX := 0; charX < calculatedCharacterWidth; charX++ {
			currentChar := layerEntry.CharacterMemory[charY][charX]

			// Underlying cell
			var underlyingCell types.CharacterEntryType
			if underlyingLayer != nil && charY < underlyingHeight && charX < underlyingWidth {
				underlyingCell = underlyingLayer.CharacterMemory[charY][charX]
			} else {
				underlyingCell = types.NewCharacterEntry()
			}

			// Check transparency for upper and lower pixels
			upperTransparent := currentImageY < imageBounds.Min.Y || currentImageY >= imageBounds.Max.Y ||
				isTransparentPixel(processedImageData, charX, currentImageY)
			lowerTransparent := currentImageY+1 < imageBounds.Min.Y || currentImageY+1 >= imageBounds.Max.Y ||
				isTransparentPixel(processedImageData, charX, currentImageY+1)

			switch {
			case upperTransparent && lowerTransparent:
				// Fully transparent → NullRune + transparent background
				currentChar.Character = constants.NullRune
				currentChar.AttributeEntry.IsBackgroundTransparent = true

			case upperTransparent && !lowerTransparent:
				// Only lower visible → lower half block
				currentChar.Character = constants.CharBlockLowerHalf

				lowerPixel := processedImageData.At(charX, currentImageY+1)
				r, g, b, _ := get8BitColorComponents(lowerPixel)
				currentChar.AttributeEntry.ForegroundColor = GetRGBColor(r, g, b)
				currentChar.AttributeEntry.BackgroundColor = underlyingCell.AttributeEntry.BackgroundColor

			case !upperTransparent && lowerTransparent:
				// Only upper visible → upper half block
				currentChar.Character = constants.CharBlockUpperHalf

				upperPixel := processedImageData.At(charX, currentImageY)
				r, g, b, _ := get8BitColorComponents(upperPixel)
				currentChar.AttributeEntry.ForegroundColor = GetRGBColor(r, g, b)
				currentChar.AttributeEntry.BackgroundColor = underlyingCell.AttributeEntry.BackgroundColor

			case !upperTransparent && !lowerTransparent:
				// Both visible → full block (upper = foreground, lower = background)
				currentChar.Character = constants.CharBlockUpperHalf

				upperPixel := processedImageData.At(charX, currentImageY)
				r, g, b, _ := get8BitColorComponents(upperPixel)
				currentChar.AttributeEntry.ForegroundColor = GetRGBColor(r, g, b)

				lowerPixel := processedImageData.At(charX, currentImageY+1)
				r, g, b, _ = get8BitColorComponents(lowerPixel)
				currentChar.AttributeEntry.BackgroundColor = GetRGBColor(r, g, b)
			}

			layerEntry.CharacterMemory[charY][charX] = currentChar
		}
		currentImageY += 2
	}

	return layerEntry
}

func isTransparentPixel(processedImageData image.Image, x, y int) bool {
	// Get the color at the specified pixel
	c := processedImageData.At(x, y)

	// Convert to RGBA to get access to individual channels
	rgba := color.RGBAModel.Convert(c).(color.RGBA)

	// Check if alpha value is 0 (fully transparent)
	return rgba.A == 0
}

// GetRGBComponents is a wrapper for GetRGBColorComponents
func GetRGBComponents(color constants.ColorType) (int32, int32, int32) {
	return GetRGBColorComponents(color)
}

/*
get8BitColorComponents allows you to get red, green, and blue color components
from a specific color.
*/
func get8BitColorComponents(colorEntry color.Color) (int32, int32, int32, uint32) {
	redIndex, greenIndex, blueIndex, alphaIndex := colorEntry.RGBA()
	return int32(redIndex) / 257, int32(greenIndex) / 257, int32(blueIndex) / 257, alphaIndex / 257
}

// Function to calculate the brightness of a pixel
func calculateBrightness(r, g, b uint8) float64 {
	// Using a common formula to calculate brightness (perceived luminance)
	// Scale result to be between 0 and 1
	return (0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b)) / 255.0
}

// Function to map brightness to an ASCII character from the bToC_full mapping
func mapBrightnessToCharacter(brightness float64) rune {
	// Mapping brightness (0 to 1) to corresponding ASCII character from bToC_full
	bToC_full := map[float64][]rune{
		0.0:      {'.'},
		0.1:      {'.', '`'},
		0.133333: {'.', '`'},
		0.155556: {'-'},
		0.177778: {'\'', ',', '_'},
		0.266667: {':', '=', '^'},
		0.311111: {'"', '+', '/', '\\'},
		0.333333: {'~'},
		0.355556: {';', '|'},
		0.4:      {'(', ')', '<', '>'},
		0.444444: {'%', '?', 'c', 's', '{', '}'},
		0.488889: {'!', 'I', '[', ']', 'i', 't', 'v', 'x', 'z'},
		0.511111: {'1', 'r'},
		0.533333: {'*', 'a', 'e', 'l', 'o'},
		0.555556: {'n', 'u'},
		0.577778: {'T', 'f', 'w'},
		0.6:      {'3', '7'},
		0.622222: {'J', 'j', 'y'},
		0.644444: {'5'},
		0.666667: {'$', '2', '6', '9', 'C', 'L', 'Y', 'm'},
		0.688889: {'S'},
		0.711111: {'4', 'g', 'k', 'p', 'q'},
		0.733333: {'F', 'P', 'b', 'd', 'h'},
		0.755556: {'G', 'O', 'V', 'X'},
		0.777778: {'E', 'Z'},
		0.8:      {'8', 'A', 'U'},
		0.844444: {'D', 'H', 'K', 'W'},
		0.888889: {'&', '@', 'R'},
		0.911111: {'B', 'Q'},
		0.933333: {'#'},
		1.0:      {'0', 'M', 'N'},
	}

	// Find the appropriate character for the brightness level
	for threshold, characters := range bToC_full {
		if brightness <= threshold {
			// Pick a random character from the list at this threshold
			randomIndex := rand.Intn(len(characters)) // Generates a random index within the range of available characters
			return characters[randomIndex]
		}
	}
	return ' ' // Default to space if no match
}

// Function to process the image and convert it to ASCII art
func GetImageLayerAsAsciiColorArt(sourceImageData image.Image, imageStyle types.ImageStyleEntryType, widthInCharacters int, heightInCharacters int, blurSigma float64, existingLayer ...*types.LayerEntryType) types.LayerEntryType {
	if widthInCharacters <= 0 && heightInCharacters <= 0 {
		safeSttyPanic(fmt.Sprintf("The specified width and height of %dx%d for your image is not valid.", widthInCharacters, heightInCharacters))
	}

	// Seed the random number generator for random character selection
	rand.Seed(time.Now().UnixNano())

	calculatedPixelWidth := widthInCharacters
	calculatedPixelHeight := heightInCharacters * 2
	if widthInCharacters == 0 {
		calculatedPixelWidth = (heightInCharacters * 2 * sourceImageData.Bounds().Max.X) / sourceImageData.Bounds().Max.Y
	}
	if heightInCharacters == 0 {
		calculatedPixelHeight = (widthInCharacters * sourceImageData.Bounds().Max.Y) / sourceImageData.Bounds().Max.X
	}

	// Resize the image based on calculated dimensions
	processedImageData := resizeImage(sourceImageData, uint(calculatedPixelWidth), uint(calculatedPixelHeight), imageStyle.IsWidthAspectRatioPreserved, imageStyle.IsHeightAspectRatioPreserved)

	// Apply blur if needed
	if blurSigma > 0 {
		processedImageData = imaging.Blur(processedImageData, blurSigma)
	}

	// Convert to grayscale if specified
	if imageStyle.IsGrayscale {
		processedImageData = ConvertImageToGrayscale(processedImageData)
	}

	// Initialize the layer entry for the image
	calculatedCharacterWidth := calculatedPixelWidth
	calculatedCharacterHeight := calculatedPixelHeight / 2

	// Use existing layer if provided, otherwise create a new one
	var layerEntry types.LayerEntryType
	if len(existingLayer) > 0 && existingLayer[0] != nil {
		layerEntry = *existingLayer[0]
	} else {
		layerEntry = types.NewLayerEntry("", "", calculatedCharacterWidth, calculatedCharacterHeight)
	}

	// Loop through each character position in the grid
	for currentYLocation := 0; currentYLocation < calculatedCharacterHeight; currentYLocation++ {
		for currentXLocation := 0; currentXLocation < calculatedCharacterWidth; currentXLocation++ {
			currentCharacter := layerEntry.CharacterMemory[currentYLocation][currentXLocation]

			// Get the upper pixel's color (as uint8)
			upperPixel := processedImageData.At(currentXLocation, currentYLocation*2) // Upper half of the character
			redColor, greenColor, blueColor, _ := get8BitColorComponents(upperPixel)

			// Calculate brightness based on RGB components
			brightness := calculateBrightness(uint8(redColor), uint8(greenColor), uint8(blueColor))

			// Map the brightness to an ASCII character using bToC_full mapping
			asciiCharacter := mapBrightnessToCharacter(brightness)

			// Set the ASCII character
			currentCharacter.Character = asciiCharacter

			// Set the foreground color based on the pixel color
			currentCharacter.AttributeEntry.ForegroundColor = GetRGBColor(redColor, greenColor, blueColor)

			// Get the lower pixel's color for the background (if applicable)
			lowerPixel := processedImageData.At(currentXLocation, currentYLocation*2+1) // Lower half of the character
			redColor, greenColor, blueColor, _ = get8BitColorComponents(lowerPixel)

			// Set the background color
			currentCharacter.AttributeEntry.BackgroundColor = GetRGBColor(0, 0, 0)

			// If the alpha value is low, set character to null rune
			if redColor <= 150 || greenColor <= 150 || blueColor <= 150 {
				currentCharacter.Character = constants.NullRune
			}

			// Update the layer entry with the character and its color attributes
			layerEntry.CharacterMemory[currentYLocation][currentXLocation] = currentCharacter
		}
	}

	return layerEntry
}

// resizeImageForBlockElements resizes the source image using an area-averaging
// method to prepare it for block element rendering. It returns the raw pixel
// color data and the weight of each pixel's contribution.
func resizeImageForBlockElements(sourceImageData image.Image, targetWidth, targetHeight int) ([][][4]float64, [][]float64) {
	sourceBounds := sourceImageData.Bounds()
	sourceImageWidth, sourceImageHeight := sourceBounds.Dx(), sourceBounds.Dy()
	coverageWidth := float64(targetWidth) / float64(sourceImageWidth)
	coverageHeight := float64(targetHeight) / float64(sourceImageHeight)

	numberOfWorkers := runtime.NumCPU()
	var waitGroup sync.WaitGroup

	// Slices to hold results from each goroutine
	partialPixelColorInformation := make([][][4]float64, numberOfWorkers)
	partialPixelWeightInformation := make([][]float64, numberOfWorkers)

	rowsPerWorker := (sourceImageHeight + numberOfWorkers - 1) / numberOfWorkers

	for workerIndex := 0; workerIndex < numberOfWorkers; workerIndex++ {
		// Initialize slices for this worker
		partialPixelColorInformation[workerIndex] = make([][4]float64, targetWidth*targetHeight)
		partialPixelWeightInformation[workerIndex] = make([]float64, targetWidth*targetHeight)

		// Calculate row range for this worker
		startSourceYLocation := sourceBounds.Min.Y + workerIndex*rowsPerWorker
		endSourceYLocation := startSourceYLocation + rowsPerWorker
		if endSourceYLocation > sourceBounds.Max.Y {
			endSourceYLocation = sourceBounds.Max.Y
		}

		waitGroup.Add(1)
		go func(currentWorkerIndex int, currentStartSourceYLocation, currentEndSourceYLocation int) {
			defer waitGroup.Done()
			localPixelColorInformation := partialPixelColorInformation[currentWorkerIndex]
			localPixelWeightInformation := partialPixelWeightInformation[currentWorkerIndex]

			for sourceYLocation := currentStartSourceYLocation; sourceYLocation < currentEndSourceYLocation; sourceYLocation++ {
				for sourceXLocation := sourceBounds.Min.X; sourceXLocation < sourceBounds.Max.X; sourceXLocation++ {
					redComponent, greenComponent, blueComponent, alphaComponent := sourceImageData.At(sourceXLocation, sourceYLocation).RGBA()
					redColorValue := float64(redComponent) / 0xffff
					greenColorValue := float64(greenComponent) / 0xffff
					blueColorValue := float64(blueComponent) / 0xffff
					alphaColorValue := float64(alphaComponent) / 0xffff

					targetYLocationStartingPoint := float64(sourceYLocation-sourceBounds.Min.Y) * coverageHeight
					targetYLocationEndingPoint := targetYLocationStartingPoint + coverageHeight
					fromTargetYLocation := int(targetYLocationStartingPoint)
					toTargetYLocation := int(targetYLocationEndingPoint)

					for targetYLocation := fromTargetYLocation; targetYLocation <= toTargetYLocation && targetYLocation < targetHeight; targetYLocation++ {
						yAxisCoverage := 1.0
						if targetYLocation == fromTargetYLocation {
							yAxisCoverage -= math.Mod(targetYLocationStartingPoint, 1.0)
						}
						if targetYLocation == toTargetYLocation {
							yAxisCoverage -= 1.0 - math.Mod(targetYLocationEndingPoint, 1.0)
						}

						targetXLocationStartingPoint := float64(sourceXLocation-sourceBounds.Min.X) * coverageWidth
						targetXLocationEndingPoint := targetXLocationStartingPoint + coverageWidth
						fromTargetXLocation := int(targetXLocationStartingPoint)
						toTargetXLocation := int(targetXLocationEndingPoint)

						for targetXLocation := fromTargetXLocation; targetXLocation <= toTargetXLocation && targetXLocation < targetWidth; targetXLocation++ {
							xAxisCoverage := 1.0
							if targetXLocation == fromTargetXLocation {
								xAxisCoverage -= math.Mod(targetXLocationStartingPoint, 1.0)
							}
							if targetXLocation == toTargetXLocation {
								xAxisCoverage -= 1.0 - math.Mod(targetXLocationEndingPoint, 1.0)
							}

							totalCoverage := xAxisCoverage * yAxisCoverage
							if totalCoverage <= 0 {
								continue
							}

							currentPixelIndex := targetYLocation*targetWidth + targetXLocation
							localPixelColorInformation[currentPixelIndex][0] += redColorValue * totalCoverage
							localPixelColorInformation[currentPixelIndex][1] += greenColorValue * totalCoverage
							localPixelColorInformation[currentPixelIndex][2] += blueColorValue * totalCoverage
							localPixelColorInformation[currentPixelIndex][3] += alphaColorValue * totalCoverage
							localPixelWeightInformation[currentPixelIndex] += totalCoverage
						}
					}
				}
			}
		}(workerIndex, startSourceYLocation, endSourceYLocation)
	}

	waitGroup.Wait()

	return partialPixelColorInformation, partialPixelWeightInformation
}

// mergeAndNormalizeInParallel merges and normalizes pixel data in parallel.
func mergeAndNormalizeInParallel(partialPixelColorInformation [][][4]float64, partialPixelWeightInformation [][]float64, targetWidth, targetHeight int) [][4]float64 {
	numberOfWorkers := runtime.NumCPU()
	var waitGroup sync.WaitGroup

	finalPixelColorInformation := make([][4]float64, targetWidth*targetHeight)
	pixelsPerWorker := (len(finalPixelColorInformation) + numberOfWorkers - 1) / numberOfWorkers

	for workerIndex := 0; workerIndex < numberOfWorkers; workerIndex++ {
		startPixelIndex := workerIndex * pixelsPerWorker
		endPixelIndex := startPixelIndex + pixelsPerWorker
		if endPixelIndex > len(finalPixelColorInformation) {
			endPixelIndex = len(finalPixelColorInformation)
		}

		waitGroup.Add(1)
		go func(start, end int) {
			defer waitGroup.Done()
			for pixelDataIndex := start; pixelDataIndex < end; pixelDataIndex++ {
				var totalRed, totalGreen, totalBlue, totalAlpha, totalWeight float64
				for i := 0; i < numberOfWorkers; i++ {
					totalRed += partialPixelColorInformation[i][pixelDataIndex][0]
					totalGreen += partialPixelColorInformation[i][pixelDataIndex][1]
					totalBlue += partialPixelColorInformation[i][pixelDataIndex][2]
					totalAlpha += partialPixelColorInformation[i][pixelDataIndex][3]
					totalWeight += partialPixelWeightInformation[i][pixelDataIndex]
				}

				if totalWeight > 0 {
					finalPixelColorInformation[pixelDataIndex][0] = totalRed / totalWeight
					finalPixelColorInformation[pixelDataIndex][1] = totalGreen / totalWeight
					finalPixelColorInformation[pixelDataIndex][2] = totalBlue / totalWeight
					finalPixelColorInformation[pixelDataIndex][3] = totalAlpha / totalWeight
				}
			}
		}(startPixelIndex, endPixelIndex)
	}

	waitGroup.Wait()
	return finalPixelColorInformation
}

/*
findBestBlockElementForCell analyzes an 8x8 grid of pixels to find the optimal
block element character and corresponding foreground/background colors to
represent that portion of the image.

This function is at the core of the block element rendering style. It takes a
small section of the source image (corresponding to a single character cell)
and determines the best Unicode block element character (like '▀', '▐', '░', etc.)
to approximate it. It does this by trying every available block element and
calculating which one, along with its optimal foreground and background colors,
minimizes the visual error (Sum of Absolute Differences) compared to the
original pixels.

The function also includes logic to handle transparency and to discard cells
that have very little content or are poorly represented, preventing visual noise
in the final output. In addition, the following information should be noted:

- `transparentForegroundPenalty`: This parameter controls how strongly the
algorithm avoids placing foreground parts of a block element over transparent
areas of the original image. A higher value results in a larger penalty,
making the algorithm more aggressively select block elements that do not
have "spikes" or "stray pixels" protruding into transparent regions. This is
useful for cleaning up the edges of sprites against a transparent background.
A typical range is 10.0 to 100.0. A value of 0 disables this feature.

- `aggressiveCoverageThreshold`: This is the minimum percentage of opaque
pixels required within an 8x8 cell to consider it for rendering. If the
coverage is below this threshold (e.g., less than 35% of the 64 pixels are
opaque), the cell may be culled. This helps remove isolated, noisy pixels
or very thin, faint parts of an image that don't render well as block
elements. The value should be between 0.0 (nothing is culled) and 1.0 (everything
is culled). A typical value is around 0.35.

- `aggressiveErrorThreshold`: This sets the maximum allowed error (Sum of
Absolute Differences) for a low-coverage cell to survive culling. Even if a
cell is below `aggressiveCoverageThreshold`, it can be kept if it's still a
very good fit for a block element (i.e., its error is below this threshold).
Lowering this value makes the culling more aggressive, as it requires even
low-coverage cells to be a near-perfect match. A typical range is 1.0 to 5.0.
*/
func findBestBlockElementForCell(
	pixelColorData [][4]float64,
	cellRowLocation, cellColumnLocation,
	characterGridWidth, characterGridHeight int,
	transparentForegroundPenalty float64,
	aggressiveCoverageThreshold float64,
	aggressiveErrorThreshold float64,
) (rune, [3]float64, [3]float64, bool, bool) {

	minimumSAD := math.MaxFloat64
	var bestBlockElement rune
	var bestForegroundColor, bestBackgroundColor [3]float64
	var isForegroundColorTransparent, isBackgroundColorTransparent bool

	// Compute total opaque pixels
	var totalOpaquePixels float64
	for pixelY := 0; pixelY < 8; pixelY++ {
		for pixelX := 0; pixelX < 8; pixelX++ {
			pixelIndex := (cellRowLocation*8+pixelY)*characterGridWidth*8 + (cellColumnLocation*8 + pixelX)
			if pixelColorData[pixelIndex][3] >= 0.1 {
				totalOpaquePixels++
			}
		}
	}
	coverage := totalOpaquePixels / 64.0

	// Try each block element
	for blockElement, bitmask := range constants.CharBlockBitmasks {
		var foregroundColor, backgroundColor [3]float64
		var foregroundAlpha, backgroundAlpha float64
		var setBitCount, unsetBitCount float64
		currentBit := uint64(1)

		// Compute average fg/bg colors
		for pixelY := 0; pixelY < 8; pixelY++ {
			for pixelX := 0; pixelX < 8; pixelX++ {
				pixelIndex := (cellRowLocation*8+pixelY)*characterGridWidth*8 + (cellColumnLocation*8 + pixelX)
				alpha := pixelColorData[pixelIndex][3]

				if alpha < 0.1 {
					currentBit <<= 1
					continue
				}

				if bitmask&currentBit != 0 {
					foregroundColor[0] += pixelColorData[pixelIndex][0]
					foregroundColor[1] += pixelColorData[pixelIndex][1]
					foregroundColor[2] += pixelColorData[pixelIndex][2]
					foregroundAlpha += alpha
					setBitCount++
				} else {
					backgroundColor[0] += pixelColorData[pixelIndex][0]
					backgroundColor[1] += pixelColorData[pixelIndex][1]
					backgroundColor[2] += pixelColorData[pixelIndex][2]
					backgroundAlpha += alpha
					unsetBitCount++
				}
				currentBit <<= 1
			}
		}

		// Normalize colors
		if setBitCount > 0 {
			for c := 0; c < 3; c++ {
				foregroundColor[c] /= setBitCount
				foregroundColor[c] = math.Min(math.Max(foregroundColor[c], 0), 1)
			}
		}
		if unsetBitCount > 0 {
			for c := 0; c < 3; c++ {
				backgroundColor[c] /= unsetBitCount
				backgroundColor[c] = math.Min(math.Max(backgroundColor[c], 0), 1)
			}
		}

		// Compute SAD (Sum of Absolute Differences)
		var sumAbsDiff float64
		currentBit = 1
		for pixelY := 0; pixelY < 8; pixelY++ {
			for pixelX := 0; pixelX < 8; pixelX++ {
				pixelIndex := (cellRowLocation*8+pixelY)*characterGridWidth*8 + (cellColumnLocation*8 + pixelX)
				alpha := pixelColorData[pixelIndex][3]

				if alpha < 0.1 && bitmask&currentBit != 0 {
					// Aggressive penalty for placing foreground over transparent pixel
					sumAbsDiff += transparentForegroundPenalty
					currentBit <<= 1
					continue
				}

				var pixelCol [3]float64
				if bitmask&currentBit != 0 {
					pixelCol = foregroundColor
				} else {
					pixelCol = backgroundColor
				}

				for c := 0; c < 3; c++ {
					sumAbsDiff += math.Abs(pixelColorData[pixelIndex][c] - pixelCol[c])
				}
				currentBit <<= 1
			}
		}

		// Update best block if SAD is lower
		if sumAbsDiff < minimumSAD {
			minimumSAD = sumAbsDiff
			bestBlockElement = blockElement
			bestForegroundColor = foregroundColor
			bestBackgroundColor = backgroundColor
			isForegroundColorTransparent = (setBitCount == 0) || (foregroundAlpha/setBitCount < 0.5)
			isBackgroundColorTransparent = (unsetBitCount == 0) || (backgroundAlpha/unsetBitCount < 0.5)
		}
	}

	// Aggressive culling: remove cells with very low coverage OR poorly-fitting blocks
	if coverage < aggressiveCoverageThreshold && minimumSAD > aggressiveErrorThreshold {
		bestBlockElement = ' '
		isForegroundColorTransparent = true
		isBackgroundColorTransparent = true
	}

	return bestBlockElement, bestForegroundColor, bestBackgroundColor, isForegroundColorTransparent, isBackgroundColorTransparent
}

type cellJob struct {
	rowLocation    int
	columnLocation int
}

// processCellsInParallel processes all character cells in parallel using a worker pool.
func processCellsInParallel(pixelColorData [][4]float64, characterWidth, characterHeight int, targetLayerEntry *types.LayerEntryType, underlyingLayerEntry *types.LayerEntryType, imageStyle types.ImageStyleEntryType) {
	numberOfWorkers := runtime.NumCPU()
	jobs := make(chan cellJob, characterWidth*characterHeight)
	var waitGroup sync.WaitGroup

	// Get the actual dimensions of the layer
	layerHeight := len(targetLayerEntry.CharacterMemory)
	var layerWidth int
	if layerHeight > 0 {
		layerWidth = len(targetLayerEntry.CharacterMemory[0])
	}

	underlyingLayerHeight := len(underlyingLayerEntry.CharacterMemory)
	var underlyingLayerWidth int
	if underlyingLayerHeight > 0 {
		underlyingLayerWidth = len(underlyingLayerEntry.CharacterMemory[0])
	}

	for workerIndex := 0; workerIndex < numberOfWorkers; workerIndex++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			for job := range jobs {
				// Skip processing if the cell is out of bounds
				if job.rowLocation < 0 || job.rowLocation >= layerHeight ||
					job.columnLocation < 0 || job.columnLocation >= layerWidth {
					continue
				}

				// Get the underlying cell from the destination layer before modification.
				var underlyingCell types.CharacterEntryType
				if job.rowLocation >= 0 && job.rowLocation < underlyingLayerHeight &&
					job.columnLocation >= 0 && job.columnLocation < underlyingLayerWidth {
					underlyingCell = underlyingLayerEntry.CharacterMemory[job.rowLocation][job.columnLocation]
				} else {
					underlyingCell = types.NewCharacterEntry()
				}

				bestBlockElement, bestForegroundColor, bestBackgroundColor, isForegroundColorTransparent, isBackgroundColorTransparent := findBestBlockElementForCell(pixelColorData, job.rowLocation, job.columnLocation, characterWidth, characterHeight, imageStyle.TransparentForegroundPenalty, imageStyle.AggressiveCoverageThreshold, imageStyle.AggressiveErrorThreshold)

				// Convert the calculated best colors for the incoming image cell to ColorType.
				red := int32(math.Min(255, bestForegroundColor[0]*255))
				green := int32(math.Min(255, bestForegroundColor[1]*255))
				blue := int32(math.Min(255, bestForegroundColor[2]*255))
				foregroundColor := GetRGBColor(red, green, blue)

				red = int32(math.Min(255, bestBackgroundColor[0]*255))
				green = int32(math.Min(255, bestBackgroundColor[1]*255))
				blue = int32(math.Min(255, bestBackgroundColor[2]*255))
				backgroundColor := GetRGBColor(red, green, blue)

				// If the foreground part of the incoming cell is transparent, determine its color
				// based on the underlying cell and the transparency mode.
				if isForegroundColorTransparent {
					switch imageStyle.TransparencyMode {
					case constants.TransparencyModeForeground:
						foregroundColor = underlyingCell.AttributeEntry.ForegroundColor
					case constants.TransparencyModeBackground:
						foregroundColor = underlyingCell.AttributeEntry.BackgroundColor
					case constants.TransparencyModeBlended:
						r1, g1, b1 := GetRGBComponents(underlyingCell.AttributeEntry.ForegroundColor)
						r2, g2, b2 := GetRGBComponents(underlyingCell.AttributeEntry.BackgroundColor)
						blendedColor := GetRGBColor((r1+r2)/2, (g1+g2)/2, (b1+b2)/2)
						foregroundColor = blendedColor
					}
				}

				// If the background part of the incoming cell is transparent, determine its color
				// based on the underlying cell and the transparency mode.
				if isBackgroundColorTransparent {
					switch imageStyle.TransparencyMode {
					case constants.TransparencyModeForeground:
						backgroundColor = underlyingCell.AttributeEntry.ForegroundColor
					case constants.TransparencyModeBackground:
						backgroundColor = underlyingCell.AttributeEntry.BackgroundColor
					case constants.TransparencyModeBlended:
						r1, g1, b1 := GetRGBComponents(underlyingCell.AttributeEntry.ForegroundColor)
						r2, g2, b2 := GetRGBComponents(underlyingCell.AttributeEntry.BackgroundColor)
						blendedColor := GetRGBColor((r1+r2)/2, (g1+g2)/2, (b1+b2)/2)
						backgroundColor = blendedColor
					}
				}

				// Update the layer entry with the final character and colors.
				if bestBlockElement == ' ' {
					// A space character indicates the entire 8x8 grid is transparent.
					// We mark it as NullRune so the overlay process skips it, preserving the underlying cell.
					targetLayerEntry.CharacterMemory[job.rowLocation][job.columnLocation].Character = constants.NullRune
					targetLayerEntry.CharacterMemory[job.rowLocation][job.columnLocation].AttributeEntry.IsBackgroundTransparent = true
					targetLayerEntry.CharacterMemory[job.rowLocation][job.columnLocation].AttributeEntry.CellType = constants.CellTypeShadow
					//layerEntry.CharacterMemory[job.rowLocation][job.columnLocation].AttributeEntry.ForegroundColor = underlyingCell.AttributeEntry.ForegroundColor
					//layerEntry.CharacterMemory[job.rowLocation][job.columnLocation].AttributeEntry.ForegroundColor = underlyingCell.AttributeEntry.BackgroundColor
				} else {
					targetLayerEntry.CharacterMemory[job.rowLocation][job.columnLocation].Character = bestBlockElement
					targetLayerEntry.CharacterMemory[job.rowLocation][job.columnLocation].AttributeEntry.ForegroundColor = foregroundColor
					targetLayerEntry.CharacterMemory[job.rowLocation][job.columnLocation].AttributeEntry.BackgroundColor = backgroundColor
				}
			}
		}()
	}

	// Only process cells that are within the bounds of both the image and the layer
	for rowLocation := 0; rowLocation < characterHeight; rowLocation++ {
		for columnLocation := 0; columnLocation < characterWidth; columnLocation++ {
			jobs <- cellJob{rowLocation: rowLocation, columnLocation: columnLocation}
		}
	}
	close(jobs)
	waitGroup.Wait()
}

// getImageLayerAsBlockElements renders an image using block elements.
// It divides each character cell into an 8x8 grid and finds the best block element
// to represent the image in that cell.
func getImageLayerAsBlockElements(sourceImageData image.Image, imageStyle types.ImageStyleEntryType, widthInCharacters int, heightInCharacters int, blurSigma float64, existingLayer ...*types.LayerEntryType) types.LayerEntryType {
	if !imageStyle.IsWidthAspectRatioPreserved && !imageStyle.IsHeightAspectRatioPreserved {
		if widthInCharacters <= 0 || heightInCharacters <= 0 {
			safeSttyPanic(fmt.Sprintf("The specified width and height of %dx%d for your image is not valid when aspect ratio is not preserved.", widthInCharacters, heightInCharacters))
		}
	} else {
		if widthInCharacters <= 0 && heightInCharacters <= 0 {
			safeSttyPanic(fmt.Sprintf("The specified width and height of %dx%d for your image is not valid.", widthInCharacters, heightInCharacters))
		}
	}

	// Apply blur and grayscale if specified
	processedImageData := sourceImageData
	if blurSigma > 0 {
		processedImageData = imaging.Blur(sourceImageData, blurSigma)
	}
	if imageStyle.IsGrayscale {
		processedImageData = ConvertImageToGrayscale(processedImageData)
	}

	// Calculate the dimensions
	sourceBounds := processedImageData.Bounds()
	sourceImageWidth, sourceImageHeight := sourceBounds.Dx(), sourceBounds.Dy()

	// Calculate width and height in characters
	characterWidth, characterHeight := widthInCharacters, heightInCharacters

	// If both width and height are 0, use the image's dimensions
	if characterWidth == 0 && characterHeight == 0 {
		characterWidth, characterHeight = sourceImageWidth, sourceImageHeight
		if adjustedWidth := sourceImageWidth * characterHeight / sourceImageHeight; adjustedWidth < characterWidth {
			characterWidth = adjustedWidth
		} else {
			characterHeight = sourceImageHeight * characterWidth / sourceImageWidth
		}
	} else {
		// If only one dimension is specified, calculate the other to preserve aspect ratio
		if imageStyle.IsWidthAspectRatioPreserved && characterWidth == 0 {
			characterWidth = sourceImageWidth * characterHeight / sourceImageHeight
		} else if imageStyle.IsHeightAspectRatioPreserved && characterHeight == 0 {
			characterHeight = sourceImageHeight * characterWidth / sourceImageWidth
		}
	}

	// Use existing layer if provided, otherwise create a new one
	layerEntry := types.NewLayerEntry("", "", characterWidth, characterHeight)

	// Resize the image to an 8x8 grid per character cell
	targetPixelWidth, targetPixelHeight := characterWidth*8, characterHeight*8
	partialPixelColorInformation, partialPixelWeightInformation := resizeImageForBlockElements(processedImageData, targetPixelWidth, targetPixelHeight)

	// Merge and normalize the resized pixels in parallel
	finalPixelColorInformation := mergeAndNormalizeInParallel(partialPixelColorInformation, partialPixelWeightInformation, targetPixelWidth, targetPixelHeight)

	// Process each cell to find the best block element
	processCellsInParallel(finalPixelColorInformation, characterWidth, characterHeight, &layerEntry, existingLayer[0], imageStyle)

	return layerEntry
}

func getImageLayer(sourceImageData image.Image, imageStyle types.ImageStyleEntryType, widthInCharacters int, heightInCharacters int, blurSigma float64, layerEntry ...*types.LayerEntryType) types.LayerEntryType {
	var imageLayer types.LayerEntryType
	var existingLayer *types.LayerEntryType

	// If a layer entry is provided, use it
	if len(layerEntry) > 0 && layerEntry[0] != nil {
		existingLayer = layerEntry[0]
	}

	if imageStyle.DrawingStyle == constants.ImageStyleHighColor {
		imageLayer = getImageLayerAsHighColor(sourceImageData, imageStyle, widthInCharacters, heightInCharacters, blurSigma, existingLayer)
	} else if imageStyle.DrawingStyle == constants.ImageStyleCharacters {
		imageLayer = GetImageLayerAsAsciiColorArt(sourceImageData, imageStyle, widthInCharacters, heightInCharacters, blurSigma, existingLayer)
	} else if imageStyle.DrawingStyle == constants.ImageStyleBlockElements {
		imageLayer = getImageLayerAsBlockElements(sourceImageData, imageStyle, widthInCharacters, heightInCharacters, blurSigma, existingLayer)
	} else {
		imageLayer = getImageLayerAsBraille(sourceImageData, imageStyle, widthInCharacters, heightInCharacters, blurSigma, existingLayer)
	}
	return imageLayer
}

/*
drawImageToLayer allows you to draw a loaded image to the specified layer.
*/
func drawImageToLayer(layerEntry *types.LayerEntryType, imageLayer types.LayerEntryType, xLocation int, yLocation int) {
	imageLayer.ScreenXLocation = xLocation
	imageLayer.ScreenYLocation = yLocation
	overlayLayers(&imageLayer, layerEntry)
}

func loadImageAndGetEntry(fileName string) (*types.ImageEntryType, error) {
	var err error
	if !IsImageExists(fileName) {
		err = LoadImage(fileName)
		if err != nil {
			return nil, err
		}
		defer func() {
			UnloadImage(fileName)
		}()
	}
	imageEntry := getImage(fileName)
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
			newColor := color.Gray{Y: grayValue}
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

/*
SaveLayer allows you to save a pre-rendered layer to disk. This is useful for caching
complex layers that take time to render, such as image layers. The layer is saved in a
compressed format to minimize disk space usage. In addition, the following information
should be noted:

- The file extension ".clayer" is automatically appended to the filename if not provided.
- The layer is saved using gzip compression to minimize disk space.
- If the file cannot be written, an error is returned.
*/
func SaveLayer(layerAlias string, filePath string) error {
	validateLayer(layerAlias)
	layerEntry := Layers.Get(layerAlias)

	// Ensure the file has the correct extension
	if len(filePath) < 7 || filePath[len(filePath)-7:] != ".clayer" {
		filePath += ".clayer"
	}

	// Marshal the layer to JSON
	jsonData, err := json.Marshal(layerEntry)
	if err != nil {
		return fmt.Errorf("failed to marshal layer: %w", err)
	}

	// Compress the JSON data
	var compressedData bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedData)
	_, err = gzipWriter.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed to compress layer data: %w", err)
	}

	// Close the gzip writer to flush any remaining data
	if err = gzipWriter.Close(); err != nil {
		return fmt.Errorf("failed to finalize compression: %w", err)
	}

	// Write the compressed data to file system (local only, as virtual file systems are read-only)
	err = writeFileDataToFileSystem(filePath, compressedData.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("failed to write layer to file: %w", err)
	}

	return nil
}

/*
LoadLayer allows you to load a pre-rendered layer from disk. This is useful for quickly
loading complex layers that were previously saved, such as image layers. The layer is
loaded from a compressed format that was created by SaveLayer. In addition, the following
information should be noted:

- The file extension ".clayer" is automatically appended to the filename if not provided.
- If the file cannot be read or is not a valid layer file, an error is returned.
- The loaded layer is returned and can be added to the layer system as needed.
*/
func LoadLayer(filePath string) (types.LayerEntryType, error) {
	var layerEntry types.LayerEntryType

	// Ensure the file has the correct extension
	if len(filePath) < 7 || filePath[len(filePath)-7:] != ".clayer" {
		filePath += ".clayer"
	}

	// Read the compressed data from file system (local or virtual)
	compressedData, err := getFileDataFromFileSystem(filePath)
	if err != nil {
		return layerEntry, fmt.Errorf("failed to read layer file: %w", err)
	}

	// Create a reader for the compressed data
	gzipReader, err := gzip.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return layerEntry, fmt.Errorf("failed to decompress layer data: %w", err)
	}
	defer gzipReader.Close()

	// Read all the decompressed data
	jsonData, err := io.ReadAll(gzipReader)
	if err != nil {
		return layerEntry, fmt.Errorf("failed to read decompressed data: %w", err)
	}

	// Unmarshal the JSON data into a layer entry
	err = json.Unmarshal(jsonData, &layerEntry)
	if err != nil {
		return layerEntry, fmt.Errorf("failed to unmarshal layer data: %w", err)
	}

	return layerEntry, nil
}

/*
LoadLayerFromFile allows you to load a pre-rendered layer from disk and add it to the layer system.
This is useful for quickly loading complex layers that were previously saved, such as image layers.
The layer is loaded from a compressed format that was created by SaveLayer. In addition, the following
information should be noted:

- The file extension ".clayer" is automatically appended to the filename if not provided.
- If the file cannot be read or is not a valid layer file, an error is returned.
- The loaded layer is added to the layer system with the specified alias, position, and z-order.
- The function returns a LayerInstanceType that can be used to manipulate the loaded layer.
*/
func LoadLayerFromFile(layerAlias string, xLocation int, yLocation int, zOrderPriority int, parentAlias string, filePath string) (LayerInstanceType, error) {
	// Load the layer from file
	layerEntry, err := LoadLayer(filePath)
	if err != nil {
		return LayerInstanceType{}, err
	}

	// Add the loaded layer to the layer system
	layer.Add(layerAlias, xLocation, yLocation, layerEntry.Width, layerEntry.Height, zOrderPriority, parentAlias)

	// Get the layer entry from the layer system
	newLayerEntry := Layers.Get(layerAlias)

	// Copy the character memory from the loaded layer to the new layer
	for y := 0; y < layerEntry.Height; y++ {
		for x := 0; x < layerEntry.Width; x++ {
			newLayerEntry.CharacterMemory[y][x] = layerEntry.CharacterMemory[y][x]
			newLayerEntry.CharacterMemory[y][x].LayerAlias = layerAlias
			newLayerEntry.CharacterMemory[y][x].ParentAlias = parentAlias
		}
	}

	// Create and return a LayerInstanceType for the new layer
	layerInstance := LayerInstanceType{
		layerAlias:  layerAlias,
		parentAlias: parentAlias,
		LayerWidth:  layerEntry.Width,
		LayerHeight: layerEntry.Height,
	}

	return layerInstance, nil
}

/*
LoadPreRenderedLayerImage allows you to load a pre-rendered layer image directly into image memory.
This is different from loading an image and pre-rendering it afterwards, as it directly loads
a layer that has already been rendered. This is useful for quickly loading complex images
that have been pre-processed and saved as layers. In addition, the following information
should be noted:

- The file extension ".clayer" is automatically appended to the filename if not provided.
- If the file cannot be read or is not a valid layer file, an error is returned.
- The loaded layer is added to the image system with the specified alias.
*/
func LoadPreRenderedLayerImage(filePath string, imageAlias string) error {
	// Load the layer from file
	layerEntry, err := LoadLayer(filePath)
	if err != nil {
		return err
	}

	// Create a new image entry
	imageEntry := types.NewImageEntry()

	// Store the layer in the image entry
	imageEntry.LayerEntry = layerEntry

	// Set ImageData to nil since we're using a pre-rendered layer
	imageEntry.ImageData = nil

	// Add the image to the image system
	addImage(imageAlias, imageEntry)

	return nil
}
