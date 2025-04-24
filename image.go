package consolizer

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/types"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

type imageMemoryType struct {
	sync.Mutex
	Entries map[string]*types.ImageEntryType
}

var Image imageMemoryType

func init() {
	Image.Entries = make(map[string]*types.ImageEntryType)
}

func AddImage(imageAlias string, imageEntry types.ImageEntryType) {
	Image.Lock()
	defer func() {
		Image.Unlock()
	}()
	// verify if any errors occurred?
	Image.Entries[imageAlias] = &imageEntry
}

func GetImage(imageAlias string) *types.ImageEntryType {
	Image.Lock()
	defer func() {
		Image.Unlock()
	}()
	if Image.Entries[imageAlias] == nil {
		panic(fmt.Sprintf("The requested Image with alias '%s' could not be returned since it does not exist.", imageAlias))
	}
	return Image.Entries[imageAlias]
}
func DeleteImage(imageAlias string) {
	Image.Lock()
	defer func() {
		Image.Unlock()
	}()
	delete(Image.Entries, imageAlias)
}
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
	DeleteImage(imageAlias)
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
	AddImage(imageFile, imageEntry)
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
	AddImage(imageAlias, imageEntry)
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
	AddImage(imageAlias, imageEntry)
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
	AddImage(imageAlias, imageEntry)
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
	processedImageData := resizeImage(sourceImageData, uint(calculatedPixelWidth), uint(calculatedPixelHeight), imageStyle.IsWidthAspectRatioPreserved, imageStyle.IsHeightAspectRatioPreserved)
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
			var upperPixel color.Color
			var redColorIndex int32
			var greenColorIndex int32
			var blueColorIndex int32
			imageBounds := processedImageData.Bounds()
			if currentXLocation < imageBounds.Min.X || currentXLocation >= imageBounds.Max.X ||
				currentImageYLocation < imageBounds.Min.Y || currentImageYLocation >= imageBounds.Max.Y {
				// Out of bounds, treat as transparent
				currentCharacter.Character = constants.NullRune
			} else {
				// In bounds, get the actual pixel color
				upperPixel = processedImageData.At(currentXLocation, currentImageYLocation)
				redColorIndex, greenColorIndex, blueColorIndex, _ = get8BitColorComponents(upperPixel)
			}
			currentCharacter.AttributeEntry.ForegroundColor = GetRGBColor(int32(redColorIndex), int32(greenColorIndex), int32(blueColorIndex))
			if currentImageYLocation < calculatedCharacterHeight*2 {
				// Check for null parts of an image.
				if currentXLocation < imageBounds.Min.X || currentXLocation >= imageBounds.Max.X ||
					currentImageYLocation+1 < imageBounds.Min.Y || currentImageYLocation+1 >= imageBounds.Max.Y {
					// Out of bounds, treat as transparent
					currentCharacter.Character = constants.NullRune // For now we blank upper square since black bar may be less desierable.
					// lowerPixel := processedImageData.At(currentXLocation, currentImageYLocation+1)
					// redColorIndex, greenColorIndex, blueColorIndex, _ := get8BitColorComponents(lowerPixel)
					// currentCharacter.AttributeEntry.BackgroundColor = GetRGBColor(0, 0, 0)
				} else {
					// In bounds, get the actual pixel color
					lowerPixel := processedImageData.At(currentXLocation, currentImageYLocation+1)
					redColorIndex, greenColorIndex, blueColorIndex, _ := get8BitColorComponents(lowerPixel)
					currentCharacter.AttributeEntry.BackgroundColor = GetRGBColor(int32(redColorIndex), int32(greenColorIndex), int32(blueColorIndex))
				}
			}
			layerEntry.CharacterMemory[currentYLocation][currentXLocation] = currentCharacter
		}
		currentImageYLocation += 2
	}
	return layerEntry
}

func isTransparentPixel(processedImageData image.Image, x, y int) bool {
	// GetLayer the color at the specified pixel
	c := processedImageData.At(x, y)

	// Convert to RGBA to get access to individual channels
	rgba := color.RGBAModel.Convert(c).(color.RGBA)

	// Check if alpha value is 0 (fully transparent)
	return rgba.A == 0
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
func GetImageLayerAsAsciiColorArt(sourceImageData image.Image, imageStyle types.ImageStyleEntryType, widthInCharacters int, heightInCharacters int, blurSigma float64) types.LayerEntryType {
	if widthInCharacters <= 0 && heightInCharacters <= 0 {
		panic(fmt.Sprintf("The specified width and height of %dx%d for your image is not valid.", widthInCharacters, heightInCharacters))
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
	layerEntry := types.NewLayerEntry("", "", calculatedCharacterWidth, calculatedCharacterHeight)

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

// getImageLayerAsBlockElements renders an image using block elements.
// It divides each character cell into an 8x8 grid and finds the best block element
// to represent the image in that cell.
func getImageLayerAsBlockElements(sourceImageData image.Image, imageStyle types.ImageStyleEntryType, widthInCharacters int, heightInCharacters int, blurSigma float64) types.LayerEntryType {
	if widthInCharacters <= 0 && heightInCharacters <= 0 {
		panic(fmt.Sprintf("The specified width and height of %dx%d for your image is not valid.", widthInCharacters, heightInCharacters))
	}

	// Calculate the pixel dimensions
	calculatedPixelWidth := widthInCharacters * 8
	calculatedPixelHeight := heightInCharacters * 8
	if widthInCharacters == 0 {
		calculatedPixelWidth = (heightInCharacters * 8 * sourceImageData.Bounds().Max.X) / sourceImageData.Bounds().Max.Y
	}
	if heightInCharacters == 0 {
		calculatedPixelHeight = (widthInCharacters * 8 * sourceImageData.Bounds().Max.Y) / sourceImageData.Bounds().Max.X
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
	calculatedCharacterWidth := calculatedPixelWidth / 8
	calculatedCharacterHeight := calculatedPixelHeight / 8
	layerEntry := types.NewLayerEntry("", "", calculatedCharacterWidth, calculatedCharacterHeight)

	// Loop through each character position in the grid
	for currentYLocation := 0; currentYLocation < calculatedCharacterHeight; currentYLocation++ {
		for currentXLocation := 0; currentXLocation < calculatedCharacterWidth; currentXLocation++ {
			currentCharacter := layerEntry.CharacterMemory[currentYLocation][currentXLocation]

			// For each 8x8 pixel block, we find the best block element to represent it,
			// given the available colors.
			var (
				minMSE float64 = math.MaxFloat64 // Mean squared error.
				bestElement rune
				bestFg, bestBg constants.ColorType
			)

			// Try each block element
			for _, element := range blockElementRunes {
				// Get the bit pattern for this element
				bits := getBlockElementPattern(element)

				// Calculate the average color for the pixels covered by the set
				// bits and unset bits.
				var (
					fg, bg           [3]float64
					setBits          float64
					bit              uint64 = 1
				)
				for y := 0; y < 8; y++ {
					for x := 0; x < 8; x++ {
						pixelX := currentXLocation*8 + x
						pixelY := currentYLocation*8 + y

						// Get the pixel color
						pixel := processedImageData.At(pixelX, pixelY)
						r, g, b, _ := get8BitColorComponents(pixel)

						if bits&bit != 0 {
							fg[0] += float64(r) / 255.0
							fg[1] += float64(g) / 255.0
							fg[2] += float64(b) / 255.0
							setBits++
						} else {
							bg[0] += float64(r) / 255.0
							bg[1] += float64(g) / 255.0
							bg[2] += float64(b) / 255.0
						}
						bit <<= 1
					}
				}

				// Normalize the colors
				for ch := 0; ch < 3; ch++ {
					if setBits > 0 {
						fg[ch] /= setBits
					}
					if (64 - setBits) > 0 {
						bg[ch] /= (64 - setBits)
					}
				}

				// Calculate the error
				var mse float64
				bit = 1
				for y := 0; y < 8; y++ {
					for x := 0; x < 8; x++ {
						pixelX := currentXLocation*8 + x
						pixelY := currentYLocation*8 + y

						// Get the pixel color
						pixel := processedImageData.At(pixelX, pixelY)
						r, g, b, _ := get8BitColorComponents(pixel)

						// Calculate the error
						var targetColor [3]float64
						if bits&bit != 0 {
							targetColor = fg
						} else {
							targetColor = bg
						}

						err := math.Pow(float64(r)/255.0-targetColor[0], 2) +
							math.Pow(float64(g)/255.0-targetColor[1], 2) +
							math.Pow(float64(b)/255.0-targetColor[2], 2)
						mse += err

						bit <<= 1
					}
				}

				// Normalize the error
				mse /= 64

				// Check if this is the best match so far
				if mse < minMSE {
					minMSE = mse
					bestElement = element
					bestFg = GetRGBColor(int32(fg[0]*255), int32(fg[1]*255), int32(fg[2]*255))
					bestBg = GetRGBColor(int32(bg[0]*255), int32(bg[1]*255), int32(bg[2]*255))
				}
			}

			// Set the character and colors
			currentCharacter.Character = bestElement
			currentCharacter.AttributeEntry.ForegroundColor = bestFg
			currentCharacter.AttributeEntry.BackgroundColor = bestBg

			// Update the layer entry with the character and its color attributes
			layerEntry.CharacterMemory[currentYLocation][currentXLocation] = currentCharacter
		}
	}

	return layerEntry
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
	imageEntryType := GetImage(imageAlias)
	imageLayer := imageEntryType.LayerEntry
	if Image.Entries[imageAlias].ImageData != nil {
		imageData := Image.Entries[imageAlias].ImageData
		imageLayer = getImageLayer(imageData, imageStyle, widthInCharacters, heightInCharacters, blurSigma)
	}
	drawImageToLayer(layerAlias, imageLayer, xLocation, yLocation)
}

// blockElementRunes is a slice of all block element runes used for image rendering.
// This is more memory-efficient than storing a map with all bit patterns.
var blockElementRunes = []rune{
	constants.BlockLowerOneEighthBlock,
	constants.BlockLowerOneQuarterBlock,
	constants.BlockLowerThreeEighthsBlock,
	constants.BlockLowerHalfBlock,
	constants.BlockLowerFiveEighthsBlock,
	constants.BlockLowerThreeQuartersBlock,
	constants.BlockLowerSevenEighthsBlock,
	constants.BlockLeftSevenEighthsBlock,
	constants.BlockLeftThreeQuartersBlock,
	constants.BlockLeftFiveEighthsBlock,
	constants.BlockLeftHalfBlock,
	constants.BlockLeftThreeEighthsBlock,
	constants.BlockLeftOneQuarterBlock,
	constants.BlockLeftOneEighthBlock,
	constants.BlockQuadrantLowerLeft,
	constants.BlockQuadrantLowerRight,
	constants.BlockQuadrantUpperLeft,
	constants.BlockQuadrantUpperRight,
	constants.BlockQuadrantUpperLeftAndLowerRight,
}

// originalBlockElements contains the original hardcoded bit patterns for verification purposes.
// This map is only used for testing and validation, not in the actual rendering code.
var originalBlockElements = map[rune]uint64{
	constants.BlockLowerOneEighthBlock:            0b1111111100000000000000000000000000000000000000000000000000000000,
	constants.BlockLowerOneQuarterBlock:           0b1111111111111111000000000000000000000000000000000000000000000000,
	constants.BlockLowerThreeEighthsBlock:         0b1111111111111111111111110000000000000000000000000000000000000000,
	constants.BlockLowerHalfBlock:                 0b1111111111111111111111111111111100000000000000000000000000000000,
	constants.BlockLowerFiveEighthsBlock:          0b1111111111111111111111111111111111111111000000000000000000000000,
	constants.BlockLowerThreeQuartersBlock:        0b1111111111111111111111111111111111111111111111110000000000000000,
	constants.BlockLowerSevenEighthsBlock:         0b1111111111111111111111111111111111111111111111111111111100000000,
	constants.BlockLeftSevenEighthsBlock:          0b0111111101111111011111110111111101111111011111110111111101111111,
	constants.BlockLeftThreeQuartersBlock:         0b0011111100111111001111110011111100111111001111110011111100111111,
	constants.BlockLeftFiveEighthsBlock:           0b0001111100011111000111110001111100011111000111110001111100011111,
	constants.BlockLeftHalfBlock:                  0b0000111100001111000011110000111100001111000011110000111100001111,
	constants.BlockLeftThreeEighthsBlock:          0b0000011100000111000001110000011100000111000001110000011100000111,
	constants.BlockLeftOneQuarterBlock:            0b0000001100000011000000110000001100000011000000110000001100000011,
	constants.BlockLeftOneEighthBlock:             0b0000000100000001000000010000000100000001000000010000000100000001,
	constants.BlockQuadrantLowerLeft:              0b0000111100001111000011110000111100000000000000000000000000000000,
	constants.BlockQuadrantLowerRight:             0b1111000011110000111100001111000000000000000000000000000000000000,
	constants.BlockQuadrantUpperLeft:              0b0000000000000000000000000000000000001111000011110000111100001111,
	constants.BlockQuadrantUpperRight:             0b0000000000000000000000000000000011110000111100001111000011110000,
	constants.BlockQuadrantUpperLeftAndLowerRight: 0b1111000011110000111100001111000000001111000011110000111100001111,
}

// verifyBlockElementPatterns checks that the programmatically generated patterns
// match the original hardcoded values. This ensures that our optimization
// doesn't change the behavior of the code.
func init() {
	// Only run verification in debug/test mode
	if os.Getenv("CONSOLIZER_DEBUG") == "1" {
		for r, originalPattern := range originalBlockElements {
			generatedPattern := getBlockElementPattern(r)
			if generatedPattern != originalPattern {
				panic(fmt.Sprintf("Block element pattern mismatch for rune %U: original=%064b, generated=%064b",
					r, originalPattern, generatedPattern))
			}
		}
	}
}

// getBlockElementPattern returns the bit pattern for a given block element rune.
// A 1 bit represents a pixel that is drawn, a 0 bit represents a pixel that is not drawn.
// The least significant bit is the top left pixel, the most significant bit is the bottom
// right pixel, moving row by row from left to right, top to bottom.
func getBlockElementPattern(r rune) uint64 {
	// Generate patterns programmatically based on the block element type
	switch {
	// Lower blocks (horizontal)
	case r >= constants.BlockLowerOneEighthBlock && r <= constants.BlockLowerSevenEighthsBlock:
		// Calculate how many rows of 8 pixels should be filled (1-7)
		rows := int(r - constants.BlockLowerOneEighthBlock + 1)
		var pattern uint64
		// Fill the appropriate number of rows from the bottom
		for row := 0; row < rows; row++ {
			// Each row is 8 bits, starting from the bottom (highest bits)
			startBit := 56 - (row * 8) // 56 is the bit position of the first bit in the bottom row
			for bit := 0; bit < 8; bit++ {
				pattern |= 1 << (startBit + bit)
			}
		}
		return pattern

	// Left blocks (vertical)
	case r >= constants.BlockLeftOneEighthBlock && r <= constants.BlockLeftSevenEighthsBlock:
		// Calculate how many columns of 8 pixels should be filled (1-7)
		// Note: BlockLeftSevenEighthsBlock is the smallest value, BlockLeftOneEighthBlock is the largest
		cols := 8 - int(r - constants.BlockLeftOneEighthBlock)
		var pattern uint64
		// Fill the appropriate number of columns from the left
		for row := 0; row < 8; row++ {
			for col := 0; col < cols; col++ {
				// Calculate bit position: row * 8 + col
				bitPos := row*8 + col
				pattern |= 1 << bitPos
			}
		}
		return pattern

	// Quadrant blocks
	case r == constants.BlockQuadrantLowerLeft:
		// Lower left quadrant (bottom 4 rows, left 4 columns)
		var pattern uint64
		for row := 4; row < 8; row++ {
			for col := 0; col < 4; col++ {
				bitPos := row*8 + col
				pattern |= 1 << bitPos
			}
		}
		return pattern

	case r == constants.BlockQuadrantLowerRight:
		// Lower right quadrant (bottom 4 rows, right 4 columns)
		var pattern uint64
		for row := 4; row < 8; row++ {
			for col := 4; col < 8; col++ {
				bitPos := row*8 + col
				pattern |= 1 << bitPos
			}
		}
		return pattern

	case r == constants.BlockQuadrantUpperLeft:
		// Upper left quadrant (top 4 rows, left 4 columns)
		var pattern uint64
		for row := 0; row < 4; row++ {
			for col := 0; col < 4; col++ {
				bitPos := row*8 + col
				pattern |= 1 << bitPos
			}
		}
		return pattern

	case r == constants.BlockQuadrantUpperRight:
		// Upper right quadrant (top 4 rows, right 4 columns)
		var pattern uint64
		for row := 0; row < 4; row++ {
			for col := 4; col < 8; col++ {
				bitPos := row*8 + col
				pattern |= 1 << bitPos
			}
		}
		return pattern

	case r == constants.BlockQuadrantUpperLeftAndLowerRight:
		// Upper left and lower right quadrants
		var pattern uint64
		// Upper left
		for row := 0; row < 4; row++ {
			for col := 0; col < 4; col++ {
				bitPos := row*8 + col
				pattern |= 1 << bitPos
			}
		}
		// Lower right
		for row := 4; row < 8; row++ {
			for col := 4; col < 8; col++ {
				bitPos := row*8 + col
				pattern |= 1 << bitPos
			}
		}
		return pattern

	default:
		return 0
	}
}

func getImageLayer(sourceImageData image.Image, imageStyle types.ImageStyleEntryType, widthInCharacters int, heightInCharacters int, blurSigma float64) types.LayerEntryType {
	imageLayer := types.NewLayerEntry("", "", widthInCharacters, heightInCharacters)
	if imageStyle.DrawingStyle == constants.ImageStyleHighColor {
		imageLayer = getImageLayerAsHighColor(sourceImageData, imageStyle, widthInCharacters, heightInCharacters, blurSigma)
	} else if imageStyle.DrawingStyle == constants.ImageStyleCharacters {
		imageLayer = GetImageLayerAsAsciiColorArt(sourceImageData, imageStyle, widthInCharacters, heightInCharacters, blurSigma)
	} else if imageStyle.DrawingStyle == constants.ImageStyleBlockElements {
		imageLayer = getImageLayerAsBlockElements(sourceImageData, imageStyle, widthInCharacters, heightInCharacters, blurSigma)
	} else {
		imageLayer = getImageLayerAsBraille(sourceImageData, imageStyle, widthInCharacters, heightInCharacters, blurSigma)
	}
	return imageLayer
}

/*
drawImageToLayer allows you to draw a loaded image to the specified layer.
*/
func drawImageToLayer(layerAlias string, imageLayer types.LayerEntryType, xLocation int, yLocation int) {
	layerEntry := Layers.Get(layerAlias)
	imageLayer.ScreenXLocation = xLocation
	imageLayer.ScreenYLocation = yLocation
	overlayLayers(&imageLayer, layerEntry)
}

func getImage(fileName string) (*types.ImageEntryType, error) {
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
	imageEntry := GetImage(fileName)
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
