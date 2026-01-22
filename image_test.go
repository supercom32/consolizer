package consolizer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supercom32/consolizer"
)

// setupTest is a helper function to clean up the image entries before each test.
func setupTest() {
	consolizer.ClearAllImages()
}

func TestAddAndIsImageExists(t *testing.T) {
	setupTest()
	imageAlias := "test_data/image/test.png"
	// We can't directly add an image anymore, so we need to load it.
	// For this test, we'll create a dummy image file.
	// In a real test, you'd likely have test assets.
	err := consolizer.LoadImage(imageAlias)
	assert.Nil(t, err, "Loading image should not produce an error")
	assert.True(t, consolizer.IsImageExists(imageAlias), "Image should exist after being loaded")
}

func TestDeleteImage(t *testing.T) {
	setupTest()
	imageAlias := "test_data/image/test.png"
	err := consolizer.LoadImage(imageAlias)
	assert.Nil(t, err, "Loading image should not produce an error")
	consolizer.UnloadImage(imageAlias)
	assert.False(t, consolizer.IsImageExists(imageAlias), "Image should no longer exist after being deleted")
}

func TestIsImageExists(t *testing.T) {
	setupTest()
	imageAlias := "test_data/image/test.png"
	err := consolizer.LoadImage(imageAlias)
	assert.Nil(t, err, "Loading image should not produce an error")
	assert.True(t, consolizer.IsImageExists(imageAlias), "Image should exist after being added")
	assert.False(t, consolizer.IsImageExists("nonExistentImage"), "Non-existent image should return false")
}

func TestUnloadImage(t *testing.T) {
	setupTest()
	imageAlias := "test_data/image/test.png"
	err := consolizer.LoadImage(imageAlias)
	assert.Nil(t, err, "Loading image should not produce an error")
	consolizer.UnloadImage(imageAlias)
	assert.False(t, consolizer.IsImageExists(imageAlias), "Image should no longer exist after being unloaded")
}

func TestUnloadNonExistentImage(t *testing.T) {
	setupTest()
	nonExistentImageAlias := "nonExistentImage"
	consolizer.UnloadImage(nonExistentImageAlias)
	assert.False(t, consolizer.IsImageExists(nonExistentImageAlias))
}
