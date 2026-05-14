package types

type ImageListEntryType struct {
	FileName  string
	FileAlias string
}

type PreloadedImageListEntryType struct {
	FileName           string
	FileAlias          string
	ImageStyle         ImageStyleEntryType
	WidthInCharacters  int
	HeightInCharacters int
	BlurSigma          float64
}

type AssetListType struct {
	PreloadedImageList []PreloadedImageListEntryType
	ImageList          []ImageListEntryType
}

/*
NewAssetList is a constructor which allows you to newassetlist.

:return: AssetListType

Example:

	NewAssetList()
*/
func NewAssetList() AssetListType {
	var assetList AssetListType
	return assetList
}

/*
AddImage is a method which allows you to addimage.

:param fileName: The fileName parameter.

Example:

	instance.AddImage(fileName)
*/
func (shared *AssetListType) AddImage(fileName string) {
	var newImageListEntryType ImageListEntryType
	newImageListEntryType.FileName = fileName
	newImageListEntryType.FileAlias = fileName
	shared.ImageList = append(shared.ImageList, newImageListEntryType)
}

/*
AddPreloadedImage is a method which allows you to addpreloadedimage.

:param fileName: The fileName parameter.
:param imageStyle: The imageStyle parameter.
:param widthInCharacters: The widthInCharacters parameter.
:param heightInCharacters: The heightInCharacters parameter.
:param blurSigma: The blurSigma parameter.

Example:

	instance.AddPreloadedImage(fileName, imageStyle, widthInCharacters, heightInCharacters, blurSigma)
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
Clear is a method which allows you to clear.

Example:

	instance.Clear()
*/
func (shared *AssetListType) Clear() {
	shared.ImageList = nil
	shared.PreloadedImageList = nil
}
