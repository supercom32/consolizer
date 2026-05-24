package consolizer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/supercom32/consolizer/constants"
	"github.com/supercom32/consolizer/types"
	"testing"
)

const (
	IMAGE_TEST_SUITE_NAME = "image"
	IMAGE_COMPLEX         = "./test_data/image/complex_image.png"
	IMAGE_TRANSPARENCY    = "./test_data/image/transparency_circle2.png"
	IMAGE_GEOMETRY        = "./test_data/image/complex_geometry.png"
)

/*
setupTest is a test which allows you to initialize the testing environment by clearing image entries and setting up
standard styles and layers.

Example:

	layer1, layer2, layer3, tuiStyle, imageStyle := setupTest()
*/
func setupTest() (*LayerInstanceType, *LayerInstanceType, *LayerInstanceType, types.TuiStyleEntryType, types.ImageStyleEntryType) {
	ClearAllImages()
	imageStyle := NewImageStyle()
	imageStyle.DrawingStyle = constants.ImageStyleBlockElementsAccurate
	imageStyle.DitheringStyle = constants.DitheringStyle4x4BayerMatrix
	imageStyle.IsHistogramEqualized = true
	imageStyle.IsGrayscale = false
	imageStyle.BlurSigmaIntensity = 0
	imageStyle.TransparentForegroundPenalty = 5000 // Strongly penalize foreground on transparent pixels
	imageStyle.AggressiveErrorThreshold = 0.8      // Only very well-fitting blocks survive
	imageStyle.AggressiveCoverageThreshold = 0.5   // Cells must be at least 50% filled to survive
	layer1, layer2, layer3, tuiStyleEntry := CommonTestSetupImages()
	layer1.Color24Bit(GetRGBColor(255, 0, 0), GetRGBColor(0, 0, 255))
	layer1.FillLayer("$")
	return layer1, layer2, layer3, tuiStyleEntry, imageStyle
}

/*
TestAddAndIsImageExists is a test which allows you to verify that an image can be successfully loaded and its existence
confirmed using the IMAGE_COMPLEX asset.
*/
func TestAddAndIsImageExists(test *testing.T) {
	setupTest()
	// We can't directly add an image anymore, so we need to load it.
	// For this test, we'll create a dummy image file.
	// In a real test, you'd likely have test assets.
	err := LoadImage(IMAGE_COMPLEX)
	assert.Nil(test, err, "Loading image should not produce an error")
	assert.True(test, IsImageExists(IMAGE_COMPLEX), "Image should exist after being loaded")
}

/*
TestDeleteImage is a test which allows you to verify that a loaded image can be successfully unloaded and is
subsequently reported as non-existent.
*/
func TestDeleteImage(t *testing.T) {
	setupTest()
	err := LoadImage(IMAGE_COMPLEX)
	assert.Nil(t, err, "Loading image should not produce an error")
	UnloadImage(IMAGE_COMPLEX)
	assert.False(t, IsImageExists(IMAGE_COMPLEX), "Image should no longer exist after being deleted")
}

/*
TestIsImageExists is a test which allows you to verify the IsImageExists function correctly reports both existent and
non-existent images.
*/
func TestIsImageExists(test *testing.T) {
	setupTest()
	err := LoadImage(IMAGE_COMPLEX)
	assert.Nil(test, err, "Loading image should not produce an error")
	assert.True(test, IsImageExists(IMAGE_COMPLEX), "Image should exist after being added")
	assert.False(test, IsImageExists("nonExistentImage"), "Non-existent image should return false")
}

/*
TestUnloadImage is a test which allows you to verify that unloading an image correctly removes it from memory and
existence checks.
*/
func TestUnloadImage(test *testing.T) {
	setupTest()
	err := LoadImage(IMAGE_COMPLEX)
	assert.Nil(test, err, "Loading image should not produce an error")
	UnloadImage(IMAGE_COMPLEX)
	assert.False(test, IsImageExists(IMAGE_COMPLEX), "Image should no longer exist after being unloaded")
}

/*
TestUnloadNonExistentImage is a test which allows you to verify that attempting to unload a non-existent image does not
cause errors and correctly reports non-existence.
*/
func TestUnloadNonExistentImage(test *testing.T) {
	setupTest()
	nonExistentImageAlias := "nonExistentImage"
	UnloadImage(nonExistentImageAlias)
	assert.False(test, IsImageExists(nonExistentImageAlias))
}

/*
TestTransparentImageBlockStyleBackgroundTransparency is a test which allows you to verify that drawing a transparent
image using block style correctly handles background transparency by comparing against a known base64 encoded ANSI
string.
*/
func TestTransparentImageBlockStyleBackgroundTransparency(test *testing.T) {
	layer1, _, _, _, imageStyle := setupTest()
	layer1.DrawImage(IMAGE_TRANSPARENCY, imageStyle, 1, 0, 40, 20, 0)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, IMAGE_TEST_SUITE_NAME, "TestTransparentImageBlockStyleBackgroundTransparency", obtainedValue)
	expectedValue := LoadMasterImage(IMAGE_TEST_SUITE_NAME, "TestTransparentImageBlockStyleBackgroundTransparency")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestTransparentImageBlockStyleForegroundTransparency is a test which allows you to verify that drawing a transparent
image using block style correctly handles foreground transparency by comparing against a known base64 encoded ANSI
string.
*/
func TestTransparentImageBlockStyleForegroundTransparency(test *testing.T) {
	_, layer2, _, _, imageStyle := setupTest()
	layer2.DrawImage(IMAGE_TRANSPARENCY, imageStyle, 1, 0, 40, 20, 0)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, IMAGE_TEST_SUITE_NAME, "TestTransparentImageBlockStyleForegroundTransparency", obtainedValue)
	expectedValue := LoadMasterImage(IMAGE_TEST_SUITE_NAME, "TestTransparentImageBlockStyleForegroundTransparency")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestTransparentImageBlockStyleBlendedTransparency is a test which allows you to verify that drawing a transparent image
using block style correctly handles blended transparency by comparing against a known base64 encoded ANSI string.
*/
func TestTransparentImageBlockStyleBlendedTransparency(test *testing.T) {
	_, layer2, _, _, imageStyle := setupTest()
	layer2.DrawImage(IMAGE_TRANSPARENCY, imageStyle, 1, 0, 40, 20, 0)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, IMAGE_TEST_SUITE_NAME, "TestTransparentImageBlockStyleBlendedTransparency", obtainedValue)
	expectedValue := LoadMasterImage(IMAGE_TEST_SUITE_NAME, "TestTransparentImageBlockStyleBlendedTransparency")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestTransparentImageMediumResolutionStyleBackgroundTransparency is a test which allows you to verify that drawing a
transparent image using medium resolution style correctly handles background transparency by comparing against a known
base64 encoded ANSI string.
*/
func TestTransparentImageMediumResolutionStyleBackgroundTransparency(test *testing.T) {
	layer1, _, _, _, imageStyle := setupTest()
	imageStyle.DrawingStyle = constants.ImageStyleHalfBlock
	layer1.DrawImage(IMAGE_TRANSPARENCY, imageStyle, 1, 0, 40, 20, 0)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, IMAGE_TEST_SUITE_NAME, "TestTransparentImageMediumResolutionStyleBackgroundTransparency", obtainedValue)
	expectedValue := LoadMasterImage(IMAGE_TEST_SUITE_NAME, "TestTransparentImageMediumResolutionStyleBackgroundTransparency")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestTransparentImageMediumResolutionStyleForegroundTransparency is a test which allows you to verify that drawing a
transparent image using medium resolution style correctly handles foreground transparency by comparing against a known
base64 encoded ANSI string.
*/
func TestTransparentImageMediumResolutionStyleForegroundTransparency(test *testing.T) {
	layer1, _, _, _, imageStyle := setupTest()
	imageStyle.DrawingStyle = constants.ImageStyleHalfBlock
	layer1.DrawImage(IMAGE_TRANSPARENCY, imageStyle, 1, 0, 40, 20, 0)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, IMAGE_TEST_SUITE_NAME, "TestTransparentImageMediumResolutionStyleForegroundTransparency", obtainedValue)
	expectedValue := LoadMasterImage(IMAGE_TEST_SUITE_NAME, "TestTransparentImageMediumResolutionStyleForegroundTransparency")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestTransparentImageMediumResolutionStyleBlendedTransparency is a test which allows you to verify that drawing a
transparent image using medium resolution style correctly handles blended transparency by comparing against a known
base64 encoded ANSI string.
*/
func TestTransparentImageMediumResolutionStyleBlendedTransparency(test *testing.T) {
	layer1, _, _, _, imageStyle := setupTest()
	imageStyle.DrawingStyle = constants.ImageStyleHalfBlock
	layer1.DrawImage(IMAGE_TRANSPARENCY, imageStyle, 1, 0, 40, 20, 0)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, IMAGE_TEST_SUITE_NAME, "TestTransparentImageMediumResolutionStyleBlendedTransparency", obtainedValue)
	expectedValue := LoadMasterImage(IMAGE_TEST_SUITE_NAME, "TestTransparentImageMediumResolutionStyleBlendedTransparency")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestTransparentImageBrailleStyleBackgroundTransparency is a test which allows you to verify that drawing a transparent
image using Braille style correctly handles background transparency by comparing against a known base64 encoded ANSI
string.
*/
func TestTransparentImageBrailleStyleBackgroundTransparency(test *testing.T) {
	layer1, _, _, _, imageStyle := setupTest()
	imageStyle.DrawingStyle = constants.ImageStyleBraille
	layer1.DrawImage(IMAGE_TRANSPARENCY, imageStyle, 1, 0, 40, 20, 0)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, IMAGE_TEST_SUITE_NAME, "TestTransparentImageBrailleStyleBackgroundTransparency", obtainedValue)
	expectedValue := LoadMasterImage(IMAGE_TEST_SUITE_NAME, "TestTransparentImageBrailleStyleBackgroundTransparency")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestTransparentImageBrailleStyleForegroundTransparency is a test which allows you to verify that drawing a transparent
image using Braille style correctly handles foreground transparency by comparing against a known base64 encoded ANSI
string.
*/
func TestTransparentImageBrailleStyleForegroundTransparency(test *testing.T) {
	layer1, _, _, _, imageStyle := setupTest()
	imageStyle.DrawingStyle = constants.ImageStyleBraille
	layer1.DrawImage(IMAGE_TRANSPARENCY, imageStyle, 1, 0, 40, 20, 0)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, IMAGE_TEST_SUITE_NAME, "TestTransparentImageBrailleStyleForegroundTransparency", obtainedValue)
	expectedValue := LoadMasterImage(IMAGE_TEST_SUITE_NAME, "TestTransparentImageBrailleStyleForegroundTransparency")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestTransparentImageBrailleStyleBlendedTransparency is a test which allows you to verify that drawing a transparent
image using Braille style correctly handles blended transparency by comparing against a known base64 encoded ANSI
string.
*/
func TestTransparentImageBrailleStyleBlendedTransparency(test *testing.T) {
	layer1, _, _, _, imageStyle := setupTest()
	imageStyle.DrawingStyle = constants.ImageStyleBraille
	layer1.DrawImage(IMAGE_TRANSPARENCY, imageStyle, 1, 0, 40, 20, 0)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, IMAGE_TEST_SUITE_NAME, "TestTransparentImageBrailleStyleBlendedTransparency", obtainedValue)
	expectedValue := LoadMasterImage(IMAGE_TEST_SUITE_NAME, "TestTransparentImageBrailleStyleBlendedTransparency")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestTransparentImageAsciiStyleBackgroundTransparency is a test which allows you to verify that drawing a transparent
image using ASCII style correctly handles background transparency by comparing against a known base64 encoded ANSI
string.
*/
func TestTransparentImageAsciiStyleBackgroundTransparency(test *testing.T) {
	_, layer2, _, _, imageStyle := setupTest()
	imageStyle.DrawingStyle = constants.ImageStyleCharacters
	imageStyle.RandomSeed = 1
	layer2.DrawImage(IMAGE_TRANSPARENCY, imageStyle, 1, 0, 40, 20, 0)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, IMAGE_TEST_SUITE_NAME, "TestTransparentImageAsciiStyleBackgroundTransparency", obtainedValue)
	expectedValue := LoadMasterImage(IMAGE_TEST_SUITE_NAME, "TestTransparentImageAsciiStyleBackgroundTransparency")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestTransparentImageAsciiStyleForegroundTransparency is a test which allows you to verify that drawing a transparent
image using ASCII style correctly handles foreground transparency by comparing against a known base64 encoded ANSI
string.
*/
func TestTransparentImageAsciiStyleForegroundTransparency(test *testing.T) {
	_, layer2, _, _, imageStyle := setupTest()
	imageStyle.DrawingStyle = constants.ImageStyleCharacters
	imageStyle.RandomSeed = 1
	layer2.DrawImage(IMAGE_TRANSPARENCY, imageStyle, 1, 0, 40, 20, 0)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, IMAGE_TEST_SUITE_NAME, "TestTransparentImageAsciiStyleForegroundTransparency", obtainedValue)
	expectedValue := LoadMasterImage(IMAGE_TEST_SUITE_NAME, "TestTransparentImageAsciiStyleForegroundTransparency")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestTransparentImageAsciiStyleBlendedTransparency is a test which allows you to verify that drawing a transparent image
using ASCII style correctly handles blended transparency by comparing against a known base64 encoded ANSI string.
*/
func TestTransparentImageAsciiStyleBlendedTransparency(test *testing.T) {
	_, layer2, _, _, imageStyle := setupTest()
	imageStyle.DrawingStyle = constants.ImageStyleCharacters
	imageStyle.RandomSeed = 1
	layer2.DrawImage(IMAGE_TRANSPARENCY, imageStyle, 1, 0, 40, 20, 0)
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, IMAGE_TEST_SUITE_NAME, "TestTransparentImageAsciiStyleBlendedTransparency", obtainedValue)
	expectedValue := LoadMasterImage(IMAGE_TEST_SUITE_NAME, "TestTransparentImageAsciiStyleBlendedTransparency")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}

/*
TestComplexGeometryImage is a test which allows you to verify that a complex geometry image can be rendered using the
accurate block element style at a specific size and aspect ratio.

Example:

	Expected Inputs:
	    IMAGE_GEOMETRY asset, 140x50 size, accurate block style.

	Expected Outputs:
	    The rendered image matches the master base64 string stored in the master images directory.
*/
func TestComplexGeometryImage(test *testing.T) {
	layer1, _, _, _ := CommonTestSetupHighResolutionImages()
	ClearAllImages()
	imageStyle := NewImageStyle()
	imageStyle.DrawingStyle = constants.ImageStyleBlockElementsAccurate
	imageStyle.DitheringStyle = constants.DitheringStyle4x4BayerMatrix
	imageStyle.IsHistogramEqualized = true
	imageStyle.IsGrayscale = false
	imageStyle.BlurSigmaIntensity = 0
	imageStyle.TransparentForegroundPenalty = 5000
	imageStyle.AggressiveErrorThreshold = 0.8
	imageStyle.AggressiveCoverageThreshold = 0.5
	err := layer1.DrawImage(IMAGE_GEOMETRY, imageStyle, 0, 0, 140, 50, 0)
	assert.Nil(test, err, "Drawing image should not produce an error")
	UpdateDisplay(false)
	layerEntry := commonResource.screenLayer
	obtainedValue := layerEntry.GetBasicAnsiStringAsBase64()
	UpdateMasterImages(false, IMAGE_TEST_SUITE_NAME, "TestComplexGeometryImage", obtainedValue)
	expectedValue := LoadMasterImage(IMAGE_TEST_SUITE_NAME, "TestComplexGeometryImage")
	obtainedValueBase64 := layerEntry.GetAnsiStringFromBase64(obtainedValue)
	expectedValueBase64 := layerEntry.GetAnsiStringFromBase64(expectedValue)
	if !assert.Equalf(test, expectedValue, obtainedValue, "The updated screen does not match the master original!") {
		fmt.Println("Expected:\n", expectedValueBase64)
		fmt.Println("Obtained:\n", obtainedValueBase64)
	}
}
