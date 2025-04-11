package consolizer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"supercom32.net/consolizer"
	"supercom32.net/consolizer/types"
)

func setupTest() {
	imageAliases := []string{}
	for alias := range consolizer.Image.Entries {
		imageAliases = append(imageAliases, alias)
	}
	for _, alias := range imageAliases {
		consolizer.DeleteImage(alias)
	}
}

func TestAddAndGetImage(t *testing.T) {
	setupTest()
	imageAlias := "testImage"
	imageEntry := types.ImageEntryType{}
	consolizer.AddImage(imageAlias, imageEntry)
	retrievedImage := consolizer.GetImage(imageAlias)
	assert.NotNil(t, retrievedImage, "Retrieved image should not be nil")
	assert.True(t, consolizer.IsImageExists(imageAlias), "Image should exist after being added")
}

func TestGetNonExistentImage(t *testing.T) {
	setupTest()
	imageAlias := "nonExistentImage"
	assert.PanicsWithValue(t,
		"The requested Image with alias 'nonExistentImage' could not be returned since it does not exist.",
		func() {
			consolizer.GetImage(imageAlias)
		}, "Should panic when getting a non-existent image")
}

func TestDeleteImage(t *testing.T) {
	setupTest()
	imageAlias := "testImageToDelete"
	imageEntry := types.ImageEntryType{}
	consolizer.AddImage(imageAlias, imageEntry)
	consolizer.DeleteImage(imageAlias)
	assert.False(t, consolizer.IsImageExists(imageAlias), "Image should no longer exist after being deleted")
}

func TestIsImageExists(t *testing.T) {
	setupTest()
	imageAlias := "existingImage"
	imageEntry := types.ImageEntryType{}
	consolizer.AddImage(imageAlias, imageEntry)
	assert.True(t, consolizer.IsImageExists(imageAlias), "Image should exist after being added")
	assert.False(t, consolizer.IsImageExists("nonExistentImage"), "Non-existent image should return false")
}

func TestUnloadImage(t *testing.T) {
	setupTest()
	imageAlias := "imageToBeUnloaded"
	imageEntry := types.ImageEntryType{}
	consolizer.AddImage(imageAlias, imageEntry)
	consolizer.UnloadImage(imageAlias)
	assert.False(t, consolizer.IsImageExists(imageAlias), "Image should no longer exist after being unloaded")
}

func TestLoadImage(t *testing.T) {
	t.Skip("Requires mocking of filesystem operations")
}

func TestLoadImagesInBulk(t *testing.T) {
	t.Skip("Requires mocking of LoadImage behavior")
}

func TestLoadPreRenderedImage(t *testing.T) {
	t.Skip("Requires mocking of filesystem operations")
}

func TestDeleteNonExistentImage(t *testing.T) {
	setupTest()
	nonExistentImageAlias := "nonExistentImage"
	consolizer.DeleteImage(nonExistentImageAlias)
	assert.False(t, consolizer.IsImageExists(nonExistentImageAlias))
}

func TestUnloadNonExistentImage(t *testing.T) {
	setupTest()
	nonExistentImageAlias := "nonExistentImage"
	consolizer.UnloadImage(nonExistentImageAlias)
	assert.False(t, consolizer.IsImageExists(nonExistentImageAlias))
}
