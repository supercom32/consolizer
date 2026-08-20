package consolizer

import (
	"bytes"
	"crypto/md5"
	"embed"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/nwaples/rardecode"
	"github.com/supercom32/consolizer/constants"
	"github.com/yeka/zip"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

var virtualFileSystemArchive string
var virtualFileSystemPassword string
var virtualFileSystemArchiveType int
var virtualFileSystemEncryptionKey string
var virtualEmbeddedFileSystem embed.FS

/*
GetScrambledPassword is a method which allows you to scramble a password with a simple XOR algorithm. This allows a user
to provide a password for a virtual file system without having to store it in their own program as plaintext. To use
this feature, simply pass in your desired password and scramble key to obtain your encoded password. This password can
then be used to mount a virtual file system provided you use the same scramble key to decode it. In addition, the
following should be noted:

- This method is not designed for cryptographic security. It is simply a method of transposing a password so that it.

- The length and randomness of your password will directly influence the usefulness of your chosen scrambleKey.

Example:
    scrambled := GetScrambledPassword("mySecret", "myKey")
*/
func GetScrambledPassword(password string, scrambleKey string) string {
	scrambledPassword := xorString(password, scrambleKey)
	scrambledPassword = base64.StdEncoding.EncodeToString([]byte(scrambledPassword))
	return scrambledPassword
}

/*
getUnscrambledPassword is a method which allows you to obtain an unscrambled password that was created using the
GetScrambledPassword method. This is used by the virtual file system to decode a password that was previously scrambled
by the user in order to avoid storing passwords in plaintext. In addition, the following should be noted:

- Password scrambling is not designed for cryptographic security. It is simply a method of transposing a password so.

Example:
    plaintext := getUnscrambledPassword(scrambled, "myKey")
*/
func getUnscrambledPassword(password string, scrambleKey string) string {
	decodedString, _ := base64.StdEncoding.DecodeString(password)
	unscrambledPassword := xorString(string(decodedString), scrambleKey)
	return unscrambledPassword
}

/*
xorString is a method which allows you to perform an XOR over a given string using a given scrambleKey. This method is
useful for when you don't want to store something in plaintext. In addition, the following should be noted:

- String scrambling is not designed for cryptographic security. It is simply a method of transposing data so that it.

- While this method will properly XOR your string, it will not guarantee that the obtained result is screen printable.

- The scrambleKey will be MD5 hashed before being used. This ensures that if any part of the scramble key has been.

Example:
    xored := xorString("Secret Data", "myKey")
*/
func xorString(stringToXor string, scrambleKey string) string {
	var xoredString string
	hashedScrambleKey := getMD5Hash(scrambleKey)
	for i := 0; i < len(stringToXor); i++ {
		xoredString += string(stringToXor[i] ^ hashedScrambleKey[i%len(hashedScrambleKey)])
	}
	return xoredString
}

/*
getMD5Hash is a method which allows you to obtain an MD5 hash from a provided text string.

Example:
    hash := getMD5Hash("input string")
*/
func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

/*
mountVirtualFileSystem is a method which allows you to specify a virtual file system to mount. A virtual file system is
a ZIP or RAR archive that contains all files in which you wish to access. This is useful since instead of distributing
multiple files and folders with your application, you can simply package everything inside a virtual file system and
just include that instead. In addition, the following should be noted:

- If the archive you wish to use as your virtual file system is password protected, you must provide it here at mount.

- If password protection is used, your password should not be considered secure. Password protection is designed for.

- If you do not wish to store virtual file system passwords in plaintext, you can provide a scrambled password instead.

- If for some reason the virtual file system was unable to be mounted, an error will be returned so that your.

Example:
    err := mountVirtualFileSystem("assets.zip", "pwd", "")
*/
func mountVirtualFileSystem(archivePath string, password string, scrambleKey string) error {
	err := isArchiveFormatZip(archivePath)
	if err == nil {
		virtualFileSystemArchiveType = constants.VirtualFileSystemZip
		virtualFileSystemArchive = archivePath
		virtualFileSystemPassword = password
		virtualFileSystemEncryptionKey = scrambleKey
		return err
	}
	err = isArchiveFormatRar(archivePath, password)
	if err == nil {
		virtualFileSystemArchiveType = constants.VirtualFileSystemRar
		virtualFileSystemArchive = archivePath
		virtualFileSystemPassword = password
		virtualFileSystemEncryptionKey = scrambleKey
		return err
	}
	err = errors.New(fmt.Sprintf("Failed to open or decode '%s'.", archivePath))
	return err
}

/*
MountEmbeddedFileSystem is a method which allows you to mount an embedded file system as the virtual file system.

Example:
    MountEmbeddedFileSystem(myEmbeddedFS)
*/
func MountEmbeddedFileSystem(fs embed.FS) {
	virtualFileSystemArchiveType = constants.VirtualFileSystemEmbedded
	virtualEmbeddedFileSystem = fs
}

/*
UnmountVirtualFileSystem is a method which allows you to reset the virtual file system to an unmounted state. This is
useful for when you want to access the physical file system directly.

Example:
    UnmountVirtualFileSystem()
*/
func UnmountVirtualFileSystem() {
	virtualFileSystemArchiveType = 0
	virtualFileSystemArchive = ""
	virtualFileSystemPassword = ""
	virtualFileSystemEncryptionKey = ""
	virtualEmbeddedFileSystem = embed.FS{}
}

/*
isArchiveFormatZip is a method which allows you to detect if the provided archive is in ZIP format or not.

Example:
    err := isArchiveFormatZip("assets.zip")
*/
func isArchiveFormatZip(archivePath string) error {
	readCloser, err := zip.OpenReader(archivePath)
	if err == nil {
		_ = readCloser.Close()
	}
	return err
}

/*
isArchiveFormatRar is a method which allows you to detect if the provided archive is in RAR format or not.

Example:
    err := isArchiveFormatRar("assets.rar", "pwd")
*/
func isArchiveFormatRar(archivePath string, password string) error {
	readCloser, err := rardecode.OpenReader(archivePath, password)
	if err == nil {
		_ = readCloser.Close()
	}
	return err
}

/*
getImageFromFileSystem is a method which allows you to obtain image data from a file from the default file system. In
addition, the following should be noted:

- If for some reason the requested image could not be obtained, an error will be returned so that your application can
  handle this case appropriately.

Example:
    img, err := getImageFromFileSystem("logo.png")
*/
func getImageFromFileSystem(imageFile string) (image.Image, error) {
	var imageData image.Image
	fileData, err := getFileDataFromFileSystem(imageFile)
	if err != nil {
		err = errors.New(fmt.Sprintf("Could not get image data from '%s': %s", imageFile, err.Error()))
		return nil, err
	}
	if strings.HasSuffix(strings.ToLower(imageFile), ".jpg") || strings.HasSuffix(strings.ToLower(imageFile), ".jpeg") {
		imageData, err = jpeg.Decode(bytes.NewReader(fileData))
	}
	if strings.HasSuffix(strings.ToLower(imageFile), ".png") {
		imageData, err = png.Decode(bytes.NewReader(fileData))
	}
	if err != nil {
		err = errors.New(fmt.Sprintf("Could not decode the image '%s': %s", imageFile, err.Error()))
		return nil, err
	}
	return imageData, err
}

/*
getTextFromFileSystem is a method which allows you to obtain text data from a file from the default file system. In
addition, the following should be noted:

- If for some reason the requested text data could not be obtained, an error will be returned so that your application
  can handle this case appropriately.

Example:
    content, err := getTextFromFileSystem("readme.txt")
*/
func getTextFromFileSystem(textFile string) (string, error) {
	fileData, err := getFileDataFromFileSystem(textFile)
	dataAsString := string(fileData)
	return dataAsString, err
}

/*
getFileDataFromFileSystem is a method which allows you to get the contents of a file from the default file system. If
you have a virtual file system mounted, then the file will be retrieved from it instead of your local file system. In
addition, the following should be noted:

- If a file is being accessed from a password protected virtual file system, then the password provided at mount time
  will be used to decrypt the file automatically.

Example:
    data, err := getFileDataFromFileSystem("config.json")
*/
func getFileDataFromFileSystem(fileName string) ([]byte, error) {
	var fileData []byte
	var err error
	if virtualFileSystemArchiveType == constants.VirtualFileSystemEmbedded {
		fileData, err = virtualEmbeddedFileSystem.ReadFile(fileName)
		return fileData, err
	} else if virtualFileSystemArchiveType == constants.VirtualFileSystemZip {
		fileData, err = getFileDataFromZipArchive(fileName)
		return fileData, err
	} else if virtualFileSystemArchiveType == constants.VirtualFileSystemRar {
		fileData, err = getFileDataFromRarArchive(fileName)
		return fileData, err
	}
	fileData, err = getFileDataFromLocalFileSystem(fileName)
	if err != nil {
		err = errors.New(fmt.Sprintf("Could not open the file '%s': %s", fileName, err.Error()))
	}
	return fileData, err
}

/*
getFileDataFromLocalFileSystem is a method which allows you to get the contents of a file from the local file system. If
the contents of the file cannot be retrieved, then an error is returned instead.

Example:
    data, err := getFileDataFromLocalFileSystem("local.txt")
*/
func getFileDataFromLocalFileSystem(fileName string) ([]byte, error) {
	var fileData []byte
	fileReadCloser, err := os.Open(fileName)
	if err != nil {
		err = errors.New(fmt.Sprintf("Could not open the file '%s': %s", fileName, err.Error()))
		return fileData, err
	}
	defer fileReadCloser.Close()
	fileData, err = ioutil.ReadAll(fileReadCloser)
	if err != nil {
		err = errors.New(fmt.Sprintf("Could not read data from the file '%s': %s", fileName, err.Error()))
		return fileData, err
	}
	return fileData, err
}

/*
getFileDataFromZipArchive is a method which allows you to get the contents of a file from a ZIP archive. If the contents
of the file cannot be retrieved, then an error is returned instead. In addition, the following should be noted:

- If a file is being accessed from a password protected virtual file system, then the password provided at mount time.

Example:
    data, err := getFileDataFromZipArchive("archive_file.txt")
*/
func getFileDataFromZipArchive(fileName string) ([]byte, error) {
	var err error
	var fileReadCloser io.ReadCloser
	var fileData []byte
	archivePassword := virtualFileSystemPassword
	if virtualFileSystemEncryptionKey != "" {
		archivePassword = getUnscrambledPassword(archivePassword, virtualFileSystemEncryptionKey)
	}
	archiveReadCloser, err := zip.OpenReader(virtualFileSystemArchive)
	if err != nil {
		err = errors.New(fmt.Sprintf("Could not open '%s': %s", virtualFileSystemArchive, err.Error()))
		return fileData, err
	}
	defer archiveReadCloser.Close()
	for _, currentFile := range archiveReadCloser.File {
		if currentFile.Name == fileName {
			if currentFile.IsEncrypted() {
				currentFile.SetPassword(archivePassword)
			}
			fileReadCloser, err = currentFile.Open()
			if err != nil {
				err = errors.New(fmt.Sprintf("Could not open the file '%s' from the virtual file system: %s", fileName, err.Error()))
				return fileData, err
			}
			fileData, err = ioutil.ReadAll(fileReadCloser)
			if err != nil {
				err = errors.New(fmt.Sprintf("Could not read data from the file '%s': %s", fileName, err.Error()))
				return fileData, err
			}
			_ = fileReadCloser.Close()
		}
	}
	if fileData == nil {
		err = errors.New(fmt.Sprintf("Could not find the file '%s' from the virtual file system.", fileName))
		return fileData, err
	}
	return fileData, err
}

/*
getFileDataFromRarArchive is a method which allows you to get the contents of a file from a RAR archive. If the contents
of the file cannot be retrieved, then an error is returned instead. In addition, the following should be noted:

- If a file is being accessed from a password protected virtual file system, then the password provided at mount time.

Example:
    data, err := getFileDataFromRarArchive("archive_file.txt")
*/
func getFileDataFromRarArchive(fileName string) ([]byte, error) {
	var fileData []byte
	archivePassword := virtualFileSystemPassword
	if virtualFileSystemEncryptionKey != "" {
		archivePassword = getUnscrambledPassword(archivePassword, virtualFileSystemEncryptionKey)
	}
	archiveReadCloser, err := rardecode.OpenReader(virtualFileSystemArchive, archivePassword)
	if err != nil {
		err = errors.New(fmt.Sprintf("Could not open '%s': %s", virtualFileSystemArchive, err.Error()))
		return fileData, err
	}
	defer archiveReadCloser.Close()
	for {
		fileHeader, err := archiveReadCloser.Next()
		if err == io.EOF {
			// If EOF then we are done reading the archive.
			if fileData == nil {
				err = errors.New(fmt.Sprintf("Could not find file '%s' in archive '%s': %s", fileName, virtualFileSystemArchive, err.Error()))
			}
			return fileData, err
		}
		if err != nil {
			err = errors.New(fmt.Sprintf("Failed while scanning archive '%s': %s", virtualFileSystemArchive, err.Error()))
			return fileData, err
		}
		if fileHeader.Name == fileName {
			fileData, err = ioutil.ReadAll(archiveReadCloser)
			if err != nil {
				err = errors.New(fmt.Sprintf("Could not read data from the file '%s': %s", fileName, err.Error()))
				return fileData, err
			}
			return fileData, err
		}
	}
	return fileData, err
}

/*
writeFileDataToFileSystem is a method which allows you to write data to a file in the file system. If a virtual file
system is mounted, the file will be written to the local file system since virtual file systems are read-only. In
addition, the following should be noted:

- If the file does not already exist, it will be created with the specified permissions.

- If the file already exists, it will be overwritten.

- If the permissions parameter is 0, the default value of 0644 will be used.

Example:
    err := writeFileDataToFileSystem("output.txt", byte("Hello"), 0644)
*/
func writeFileDataToFileSystem(fileName string, data []byte, permissions int) error {
	// Virtual file systems are read-only, so always write to the local file system
	if permissions == 0 {
		permissions = 0644
	}
	perm := os.FileMode(permissions)
	err := ioutil.WriteFile(fileName, data, perm)
	if err != nil {
		err = errors.New(fmt.Sprintf("Could not write data to the file '%s': %s", fileName, err.Error()))
	}
	return err
}

/*
getFileReaderFromFileSystem is a method which allows you to obtain a file reader for a file in the file system. If a
virtual file system is mounted, the reader will be obtained from it.

Example:
    reader, err := getFileReaderFromFileSystem("data.txt")
*/
func getFileReaderFromFileSystem(fileName string) (io.ReadCloser, error) {
	if virtualFileSystemArchiveType == constants.VirtualFileSystemEmbedded {
		data, err := virtualEmbeddedFileSystem.ReadFile(fileName)
		if err != nil {
			return nil, err
		}
		return io.NopCloser(bytes.NewReader(data)), nil
	} else if virtualFileSystemArchiveType == constants.VirtualFileSystemZip {
		fileData, err := getFileDataFromZipArchive(fileName)
		if err != nil {
			return nil, err
		}
		return io.NopCloser(bytes.NewReader(fileData)), nil
	} else if virtualFileSystemArchiveType == constants.VirtualFileSystemRar {
		fileData, err := getFileDataFromRarArchive(fileName)
		if err != nil {
			return nil, err
		}
		return io.NopCloser(bytes.NewReader(fileData)), nil
	}
	// Fallback: local file system
	return os.Open(fileName)
}
