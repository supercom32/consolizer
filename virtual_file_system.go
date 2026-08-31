package consolizer

import (
	"bytes"
	"crypto/md5"
	"embed"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/nwaples/rardecode/v2"
	"github.com/supercom32/consolizer/constants"
	"github.com/yeka/zip"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
)

/*
virtualFileSystemMount is a structure which represents a single mounted archive together with the pre-built index of
every file path that archive provides. One instance exists for each successful mount and is retained for the lifetime
of the mount so that reads never have to rescan the archive. In addition, the following should be noted:

- Exactly one of the backend handles is populated depending on mountType. A ZIP mount keeps its zipReader open and
  stores every entry in zipEntries. A RAR mount keeps its rarFileSystem, which records the byte offset of every entry
  so retrieval seeks directly to the file. An embedded mount keeps the embed.FS value itself.

- The paths set holds every normalized entry path the mount can serve and is the only field consulted when rebuilding
  the union index.
*/
type virtualFileSystemMount struct {
	id                 uint64
	source             string
	mountType          int
	paths              map[string]struct{}
	zipReader          *zip.ReadCloser
	zipEntries         map[string]*zip.File
	rarFileSystem      *rardecode.RarFS
	embeddedFileSystem embed.FS
}

var virtualFileSystemMounts []*virtualFileSystemMount
var virtualFileSystemIndex = make(map[string]*virtualFileSystemMount)
var virtualFileSystemMutex sync.RWMutex
var virtualFileSystemNextId uint64

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
normalizeVirtualFileSystemPath is a method which allows you to convert an archive entry name or a lookup path into the
single canonical form used as a key in the union index. Path separators are unified to forward slashes and any leading
"./", embedded "../", duplicate slash, or trailing slash segments are removed. The result never has a leading slash and
is therefore also a valid io/fs path.

:param filePath: The archive entry name or lookup path to normalize.

:return: The canonical slash separated path with no leading slash.

Example:
    key := normalizeVirtualFileSystemPath("./myFolder\\sample.png")
*/
func normalizeVirtualFileSystemPath(filePath string) string {
	slashedPath := strings.ReplaceAll(filePath, "\\", "/")
	cleanedPath := path.Clean("/" + slashedPath)
	return strings.TrimPrefix(cleanedPath, "/")
}

/*
rebuildVirtualFileSystemIndex is a method which allows you to recompute the flat union index from the ordered list of
mounts. Every mount contributes each of its normalized paths to the index in mount order, so a path present in more
than one archive resolves to the most recently mounted archive that provides it. In addition, the following should be
noted:

- The caller must already hold the virtualFileSystemMutex write lock.

Example:
    rebuildVirtualFileSystemIndex()
*/
func rebuildVirtualFileSystemIndex() {
	rebuiltIndex := make(map[string]*virtualFileSystemMount)
	for _, mount := range virtualFileSystemMounts {
		for entryPath := range mount.paths {
			rebuiltIndex[entryPath] = mount
		}
	}
	virtualFileSystemIndex = rebuiltIndex
}

/*
newZipVirtualFileSystemMount is a constructor which allows you to open a ZIP archive and build a fully indexed mount
record for it. Every non directory entry is recorded so that a later read is a direct map lookup followed by a direct
extract with no scanning. In addition, the following should be noted:

- The returned mount keeps the underlying zip.ReadCloser open for its entire lifetime. It must be released with
  UnmountVirtualFileSystem or UnmountVirtualFileSystemArchive.

- If an entry is encrypted, the supplied password is applied to it at index time. Both traditional ZipCrypto and
  WinZip AES encrypted entries are supported for reading.

:param source: The path to the ZIP archive on the local file system.
:param password: The already unscrambled password to apply to encrypted entries, or an empty string.

:return: The indexed mount record, or an error if the archive could not be opened as a ZIP.

Example:
    mount, err := newZipVirtualFileSystemMount("assets.zip", "secret")
*/
func newZipVirtualFileSystemMount(source string, password string) (*virtualFileSystemMount, error) {
	zipReadCloser, err := zip.OpenReader(source)
	if err != nil {
		return nil, err
	}
	mount := &virtualFileSystemMount{
		source:     source,
		mountType:  constants.VirtualFileSystemZip,
		paths:      make(map[string]struct{}),
		zipReader:  zipReadCloser,
		zipEntries: make(map[string]*zip.File),
	}
	for _, currentFile := range zipReadCloser.File {
		if currentFile.FileInfo().IsDir() {
			continue
		}
		entryPath := normalizeVirtualFileSystemPath(currentFile.Name)
		if currentFile.IsEncrypted() {
			currentFile.SetPassword(password)
		}
		mount.zipEntries[entryPath] = currentFile
		mount.paths[entryPath] = struct{}{}
	}
	return mount, nil
}

/*
newRarVirtualFileSystemMount is a constructor which allows you to open a RAR archive and build a fully indexed mount
record for it. The archive is opened once through rardecode, which records the byte offset of every entry, so a later
read seeks straight to the requested file instead of scanning from the start. In addition, the following should be
noted:

- A wrong password fails here at open time when the archive uses encrypted headers, and otherwise fails later when the
  file is read.

- Solid RAR archives do not support direct entry access. Reads of entries in a solid archive will return an error.
  Create the archive without the solid option to allow random access.

:param source: The path to the RAR archive on the local file system.
:param password: The already unscrambled password for the archive, or an empty string.

:return: The indexed mount record, or an error if the archive could not be opened as a RAR.

Example:
    mount, err := newRarVirtualFileSystemMount("assets.rar", "secret")
*/
func newRarVirtualFileSystemMount(source string, password string) (*virtualFileSystemMount, error) {
	rarFileSystem, err := rardecode.OpenFS(source, rardecode.Password(password))
	if err != nil {
		return nil, err
	}
	mount := &virtualFileSystemMount{
		source:        source,
		mountType:     constants.VirtualFileSystemRar,
		paths:         make(map[string]struct{}),
		rarFileSystem: rarFileSystem,
	}
	err = fs.WalkDir(rarFileSystem, ".", func(entryPath string, dirEntry fs.DirEntry, walkError error) error {
		if walkError != nil {
			return walkError
		}
		if entryPath == "." || dirEntry.IsDir() {
			return nil
		}
		mount.paths[normalizeVirtualFileSystemPath(entryPath)] = struct{}{}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return mount, nil
}

/*
MountVirtualFileSystem is a method which allows you to add a ZIP or RAR archive to the virtual file system as an
additional layer. A virtual file system is a ZIP or RAR archive that contains the files you wish to access. Any number
of archives of either format can be mounted at the same time, and their contents are unioned into a single directory
structure. In addition, the following should be noted:

- Mounting does not replace earlier mounts. When a file exists at the same path in more than one mounted archive, the
  most recently mounted archive wins and shadows the earlier ones.

- The archive format is detected automatically by attempting to open it first as a ZIP and then as a RAR. If neither
  succeeds, an error describing both failures is returned and nothing is mounted.

- If the archive is password protected, the password must be supplied here at mount time. If a non empty scrambleKey
  is supplied, the password is treated as a scrambled password and is decoded with that key first.

:param archivePath: The path to the ZIP or RAR archive to mount.
:param password: The archive password, scrambled or plaintext, or an empty string if the archive is not protected.
:param scrambleKey: The key used to unscramble the password, or an empty string to use the password as plaintext.

:return: An error if the archive could not be opened as either a ZIP or a RAR archive.

Example:
    err := MountVirtualFileSystem("assets.zip", "pwd", "")
*/
func MountVirtualFileSystem(archivePath string, password string, scrambleKey string) error {
	if scrambleKey != "" {
		password = getUnscrambledPassword(password, scrambleKey)
	}
	mount, zipError := newZipVirtualFileSystemMount(archivePath, password)
	if zipError != nil {
		var rarError error
		mount, rarError = newRarVirtualFileSystemMount(archivePath, password)
		if rarError != nil {
			return errors.New(fmt.Sprintf("Failed to open or decode '%s' as a ZIP (%s) or RAR (%s) archive.", archivePath, zipError.Error(), rarError.Error()))
		}
	}
	virtualFileSystemMutex.Lock()
	defer virtualFileSystemMutex.Unlock()
	virtualFileSystemNextId++
	mount.id = virtualFileSystemNextId
	virtualFileSystemMounts = append(virtualFileSystemMounts, mount)
	rebuildVirtualFileSystemIndex()
	return nil
}

/*
MountEmbeddedFileSystem is a method which allows you to add an embedded file system as an additional layer of the
virtual file system. The embedded file system is walked once at mount time and every file it contains is indexed, so
its contents take part in the same union and shadowing rules as mounted archives. In addition, the following should be
noted:

- Mounting is additive. An embedded file system mounted after another mount shadows it for any shared file path.

:param embeddedFileSystem: The embedded file system whose contents should be added to the virtual file system.

:return: An error if the embedded file system could not be walked.

Example:
    err := MountEmbeddedFileSystem(myEmbeddedFileSystem)
*/
func MountEmbeddedFileSystem(embeddedFileSystem embed.FS) error {
	mount := &virtualFileSystemMount{
		mountType:          constants.VirtualFileSystemEmbedded,
		paths:              make(map[string]struct{}),
		embeddedFileSystem: embeddedFileSystem,
	}
	err := fs.WalkDir(embeddedFileSystem, ".", func(entryPath string, dirEntry fs.DirEntry, walkError error) error {
		if walkError != nil {
			return walkError
		}
		if entryPath == "." || dirEntry.IsDir() {
			return nil
		}
		mount.paths[normalizeVirtualFileSystemPath(entryPath)] = struct{}{}
		return nil
	})
	if err != nil {
		return err
	}
	virtualFileSystemMutex.Lock()
	defer virtualFileSystemMutex.Unlock()
	virtualFileSystemNextId++
	mount.id = virtualFileSystemNextId
	mount.source = fmt.Sprintf("embedded:%d", mount.id)
	virtualFileSystemMounts = append(virtualFileSystemMounts, mount)
	rebuildVirtualFileSystemIndex()
	return nil
}

/*
closeVirtualFileSystemMount is a method which allows you to release any operating system resources held by a mount. A
ZIP mount holds an open file handle through its zip.ReadCloser and that handle is closed here. RAR and embedded mounts
hold no long lived handle and require no action. In addition, the following should be noted:

- The caller must already hold the virtualFileSystemMutex write lock.

:param mount: The mount whose resources should be released.

Example:
    closeVirtualFileSystemMount(mount)
*/
func closeVirtualFileSystemMount(mount *virtualFileSystemMount) {
	if mount.zipReader != nil {
		_ = mount.zipReader.Close()
		mount.zipReader = nil
	}
}

/*
UnmountVirtualFileSystem is a method which allows you to reset the virtual file system to a completely unmounted state.
Every mounted archive and embedded file system is removed, all open archive handles are closed, and subsequent reads
go directly to the local file system again.

Example:
    UnmountVirtualFileSystem()
*/
func UnmountVirtualFileSystem() {
	virtualFileSystemMutex.Lock()
	defer virtualFileSystemMutex.Unlock()
	for _, mount := range virtualFileSystemMounts {
		closeVirtualFileSystemMount(mount)
	}
	virtualFileSystemMounts = nil
	virtualFileSystemIndex = make(map[string]*virtualFileSystemMount)
	virtualFileSystemNextId = 0
}

/*
UnmountVirtualFileSystemArchive is a method which allows you to remove a single previously mounted archive from the
virtual file system while leaving every other mount in place. The union index is recomputed afterwards so that files
which were being shadowed by the removed archive become visible again. In addition, the following should be noted:

- Archives are matched by the exact path that was passed to the mount call. If the same path was mounted more than
  once, every instance of it is removed.

:param archivePath: The archive path that was originally passed to the mount call.

:return: An error if no mounted archive matches the supplied path.

Example:
    err := UnmountVirtualFileSystemArchive("assets.zip")
*/
func UnmountVirtualFileSystemArchive(archivePath string) error {
	virtualFileSystemMutex.Lock()
	defer virtualFileSystemMutex.Unlock()
	var remainingMounts []*virtualFileSystemMount
	removedCount := 0
	for _, mount := range virtualFileSystemMounts {
		if mount.source == archivePath {
			closeVirtualFileSystemMount(mount)
			removedCount++
			continue
		}
		remainingMounts = append(remainingMounts, mount)
	}
	if removedCount == 0 {
		return errors.New(fmt.Sprintf("No virtual file system is mounted from '%s'.", archivePath))
	}
	virtualFileSystemMounts = remainingMounts
	rebuildVirtualFileSystemIndex()
	return nil
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
getFileDataFromFileSystem is a method which allows you to get the contents of a file from the virtual file system. Every
mounted archive and embedded file system is consulted through the pre-built union index, so the lookup is a single map
access followed by a direct extract from the archive that owns the file. In addition, the following should be noted:

- When one or more mounts are present the virtual file system is authoritative. A path that is not found in any mount
  returns an error and the local file system is not consulted.

- When nothing is mounted the file is read directly from the local file system.

- The password supplied when the owning archive was mounted is applied automatically to decrypt encrypted entries.

Example:
    data, err := getFileDataFromFileSystem("config.json")
*/
func getFileDataFromFileSystem(fileName string) ([]byte, error) {
	normalizedName := normalizeVirtualFileSystemPath(fileName)
	virtualFileSystemMutex.RLock()
	defer virtualFileSystemMutex.RUnlock()
	if mount := virtualFileSystemIndex[normalizedName]; mount != nil {
		fileData, err := readFileFromVirtualFileSystemMount(mount, normalizedName)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Could not read the file '%s' from the virtual file system '%s': %s", fileName, mount.source, err.Error()))
		}
		return fileData, nil
	}
	if len(virtualFileSystemMounts) > 0 {
		return nil, errors.New(fmt.Sprintf("Could not find the file '%s' in any mounted virtual file system.", fileName))
	}
	fileData, err := getFileDataFromLocalFileSystem(fileName)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not open the file '%s': %s", fileName, err.Error()))
	}
	return fileData, nil
}

/*
readFileFromVirtualFileSystemMount is a method which allows you to extract a single file from a specific mount once the
union index has identified which mount owns it. The extraction path is chosen from the mount type and performs no
scanning. In addition, the following should be noted:

- The caller must already hold at least the virtualFileSystemMutex read lock so that the mount is not released while
  the file is being extracted.

:param mount: The mount that owns the requested file.
:param normalizedName: The requested path already reduced to its canonical form.

:return: The file contents, or an error if the file could not be extracted.

Example:
    data, err := readFileFromVirtualFileSystemMount(mount, "myFolder/sample.png")
*/
func readFileFromVirtualFileSystemMount(mount *virtualFileSystemMount, normalizedName string) ([]byte, error) {
	switch mount.mountType {
	case constants.VirtualFileSystemZip:
		zipEntry := mount.zipEntries[normalizedName]
		if zipEntry == nil {
			return nil, errors.New(fmt.Sprintf("The file '%s' could not be located in the ZIP archive.", normalizedName))
		}
		fileReadCloser, err := zipEntry.Open()
		if err != nil {
			return nil, err
		}
		defer fileReadCloser.Close()
		return io.ReadAll(fileReadCloser)
	case constants.VirtualFileSystemRar:
		return mount.rarFileSystem.ReadFile(normalizedName)
	case constants.VirtualFileSystemEmbedded:
		return mount.embeddedFileSystem.ReadFile(normalizedName)
	}
	return nil, errors.New(fmt.Sprintf("Unknown virtual file system mount type %d.", mount.mountType))
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
writeFileDataToFileSystem is a method which allows you to write data to a file in the file system. Virtual file systems
are read only, so the file is always written to the local file system. In addition, the following should be noted:

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
getFileReaderFromFileSystem is a method which allows you to obtain a file reader for a file in the file system. The file
is resolved through the same union index as getFileDataFromFileSystem and the returned reader serves an in memory copy
of the file contents.

Example:
    reader, err := getFileReaderFromFileSystem("data.txt")
*/
func getFileReaderFromFileSystem(fileName string) (io.ReadCloser, error) {
	fileData, err := getFileDataFromFileSystem(fileName)
	if err != nil {
		return nil, err
	}
	return io.NopCloser(bytes.NewReader(fileData)), nil
}
