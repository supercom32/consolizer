package types

/*
ImageListEntryType is a structure which represents a single image entry in an asset list.
*/
type ImageListEntryType struct {
	FileName  string
	FileAlias string
}

/*
PreloadedImageListEntryType is a structure which represents a single preloaded image entry in an asset list.
*/
type PreloadedImageListEntryType struct {
	FileName           string
	FileAlias          string
	ImageStyle         ImageStyleEntryType
	WidthInCharacters  int
	HeightInCharacters int
	BlurSigma          float64
}

/*
AssetListType is a structure which contains a list of images and preloaded images to be loaded by the engine.
*/
type AssetListType struct {
	PreloadedImageList []PreloadedImageListEntryType
	ImageList          []ImageListEntryType
}

/*
NewAssetList is a constructor which creates a new asset list.

Example:
    assetList := NewAssetList()
*/
func NewAssetList() AssetListType {
	var assetList AssetListType
	return assetList
}

/*
AddImage is a method which adds an image file to the asset list.

Example:
    assetList.AddImage("image.png")
*/
func (shared *AssetListType) AddImage(fileName string) {
	var newImageListEntryType ImageListEntryType
	newImageListEntryType.FileName = fileName
	newImageListEntryType.FileAlias = fileName
	shared.ImageList = append(shared.ImageList, newImageListEntryType)
}

/*
AddPreloadedImage is a method which adds a pre-configured image to the asset list.

Example:
    assetList.AddPreloadedImage("image.png", style, 80, 24, 0.0)
*/
func (shared *AssetListType) AddPreloadedImage(fileName string, imageStyle ImageStyleEntryType, widthInCharacters int, heightInCharacters int, blurSigma float64) {
	var preloadedImageListEntryType PreloadedImageListEntryType
	preloadedImageListEntryType.FileName = fileName
	preloadedImageListEntryType.FileAlias = fileName
	preloadedImageListEntryType.WidthInCharacters = widthInCharacters
	preloadedImageListEntryType.HeightInCharacters = heightInCharacters
	preloadedImageListEntryType.BlurSigma = blurSigma
	preloadedImageListEntryType.ImageStyle = imageStyle
	shared.PreloadedImageList = append(shared.PreloadedImageList, preloadedImageListEntryType)
}

/*
Clear is a method which removes all images and preloaded images from the asset list.

Example:
    assetList.Clear()
*/
func (shared *AssetListType) Clear() {
	shared.ImageList = nil
	shared.PreloadedImageList = nil
}
