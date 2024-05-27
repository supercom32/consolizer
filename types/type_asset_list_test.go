package types

import (
	"github.com/stretchr/testify/assert"
	"github.com/supercom32/consolizer/internal/recast"
	"testing"
)

func TestAddImage(test *testing.T) {
	assetList := NewAssetList()
	assetList.AddImage("fileName1")
	obtainedValue := recast.GetArrayOfInterfaces(assetList.ImageList[0].FileName, assetList.ImageList[0].FileAlias)
	expectedValue := recast.GetArrayOfInterfaces("fileName1", "fileName1")
	assert.Equalf(test, expectedValue, obtainedValue, "The file entry obtained does not match what was set!")
	assetList.Clear()
	obtainedValue = recast.GetArrayOfInterfaces(len(assetList.ImageList))
	expectedValue = recast.GetArrayOfInterfaces(0)
	assert.Equalf(test, expectedValue, obtainedValue, "The number of file entries does not what was expected!")
}

func TestAddPreloadedImage(test *testing.T) {
	imageStyle := NewImageStyleEntry()
	assetList := NewAssetList()
	assetList.AddPreloadedImage("fileName1", imageStyle, 10, 11, 0.6)
	obtainedValue := recast.GetArrayOfInterfaces(assetList.PreloadedImageList[0].FileName, assetList.PreloadedImageList[0].FileAlias, assetList.PreloadedImageList[0].WidthInCharacters, assetList.PreloadedImageList[0].HeightInCharacters, assetList.PreloadedImageList[0].BlurSigma)
	expectedValue := recast.GetArrayOfInterfaces("fileName1", "fileName1", 10, 11, 0.6)
	assert.Equalf(test, expectedValue, obtainedValue, "The file entry obtained does not match what was set!")
	assetList.Clear()
	obtainedValue = recast.GetArrayOfInterfaces(len(assetList.PreloadedImageList))
	expectedValue = recast.GetArrayOfInterfaces(0)
	assert.Equalf(test, expectedValue, obtainedValue, "The number of file entries does not what was expected!")
}
