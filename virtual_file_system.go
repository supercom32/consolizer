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
	"time"
)

/*
virtualFileSystemMount is a structure which represents a single mounted archive together with the pre-built index of
every file path that archive provides. One instance exists for each successful mount and is retained for the lifetime
of the mount so that reads never have to rescan the archive. In addition, the following should be noted:

  - Exactly one of the backend handles is populated depending on mountType. A ZIP mount keeps its zipReader, which
    reads entries on demand through an io.ReaderAt, and stores every entry in zipEntries. A RAR mount keeps its
    rarFileSystem, which records the byte offset of every entry so retrieval seeks directly to the file. An embedded
    mount keeps the embed.FS value itself.

  - The closer field is only set when the mount owns the handle its archive is read through, which is the case for a
    ZIP archive mounted from a path. A mount created from a caller supplied io.ReaderAt leaves it nil, because the
    reader belongs to the caller and must never be closed by the virtual file system.

  - The paths set holds every normalized entry path the mount can serve and is the only field consulted when rebuilding
    the union index.
*/
type virtualFileSystemMount struct {
	id                 uint64
	source             string
	mountType          int
	paths              map[string]struct{}
	closer             io.Closer
	zipReader          *zip.Reader
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
virtualFileSystemMemoryRarVolumeName is the volume name handed to rardecode when a RAR archive is mounted from a
caller supplied io.ReaderAt. rardecode opens archive volumes by name through an fs.FS, so the in memory archive is
exposed under this single name and every other volume name reports that it does not exist.
*/
const virtualFileSystemMemoryRarVolumeName = "memory.rar"

/*
virtualFileSystemReaderAtFileSystem is a class which allows you to expose a caller supplied io.ReaderAt as a read only
fs.FS holding exactly one file. It exists so that rardecode, which opens archive volumes by name, can read an archive
held in memory with full random access instead of streaming it from the start. In addition, the following should be
noted:

  - Every Open returns an independent io.SectionReader over the same io.ReaderAt, so concurrent reads never share a
    file position and nothing is copied or read until the decoder asks for a specific byte range.

  - Any name other than the single volume name reports fs.ErrNotExist, which rardecode treats as the end of the
    archive. Multi volume RAR archives therefore cannot be mounted from memory.

Example:

	fileSystem := &virtualFileSystemReaderAtFileSystem{name: "memory.rar", reader: reader, size: size}
*/
type virtualFileSystemReaderAtFileSystem struct {
	name   string
	reader io.ReaderAt
	size   int64
}

/*
virtualFileSystemReaderAtFile is a class which allows you to read the single file exposed by a
virtualFileSystemReaderAtFileSystem. It satisfies fs.File and io.Seeker, which rardecode requires in order to seek
straight to an entry instead of reading the archive front to back. In addition, the following should be noted:

  - Close is a no operation because the underlying io.ReaderAt belongs to the caller and must stay usable for the
    lifetime of the mount.

Example:

	file := &virtualFileSystemReaderAtFile{SectionReader: io.NewSectionReader(reader, 0, size), name: "memory.rar"}
*/
type virtualFileSystemReaderAtFile struct {
	*io.SectionReader
	name string
}

/*
virtualFileSystemReaderAtFileInfo is a class which allows you to describe the single file exposed by a
virtualFileSystemReaderAtFileSystem through the fs.FileInfo interface.

Example:

	info := virtualFileSystemReaderAtFileInfo{name: "memory.rar", size: size}
*/
type virtualFileSystemReaderAtFileInfo struct {
	name string
	size int64
}

/*
Open is a method which allows you to open the single file held by the file system. Opening the volume name returns a
fresh seekable view over the caller's io.ReaderAt, and any other name returns an fs.PathError wrapping fs.ErrNotExist.

Example:

	file, err := fileSystem.Open("memory.rar")
*/
func (fileSystem *virtualFileSystemReaderAtFileSystem) Open(name string) (fs.File, error) {
	if name != fileSystem.name {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}
	return &virtualFileSystemReaderAtFile{
		SectionReader: io.NewSectionReader(fileSystem.reader, 0, fileSystem.size),
		name:          fileSystem.name,
	}, nil
}

/*
Stat is a method which allows you to obtain the name and size of the file without reading any of its contents.

Example:

	info, err := file.Stat()
*/
func (file *virtualFileSystemReaderAtFile) Stat() (fs.FileInfo, error) {
	return virtualFileSystemReaderAtFileInfo{name: file.name, size: file.Size()}, nil
}

/*
Close is a method which allows you to satisfy the fs.File interface. It intentionally does nothing, because the
io.ReaderAt behind the file belongs to the caller.

Example:

	err := file.Close()
*/
func (file *virtualFileSystemReaderAtFile) Close() error {
	return nil
}

/*
Name is a method which allows you to obtain the base name of the file.

Example:

	name := info.Name()
*/
func (info virtualFileSystemReaderAtFileInfo) Name() string {
	return info.name
}

/*
Size is a method which allows you to obtain the length of the file in bytes.

Example:

	size := info.Size()
*/
func (info virtualFileSystemReaderAtFileInfo) Size() int64 {
	return info.size
}

/*
Mode is a method which allows you to obtain the file mode, which always describes a read only regular file.

Example:

	mode := info.Mode()
*/
func (info virtualFileSystemReaderAtFileInfo) Mode() fs.FileMode {
	return 0444
}

/*
ModTime is a method which allows you to obtain the modification time of the file. An in memory archive has no
modification time, so the zero time is returned.

Example:

	modificationTime := info.ModTime()
*/
func (info virtualFileSystemReaderAtFileInfo) ModTime() time.Time {
	return time.Time{}
}

/*
IsDir is a method which allows you to determine whether the file is a directory, which it never is.

Example:

	isDirectory := info.IsDir()
*/
func (info virtualFileSystemReaderAtFileInfo) IsDir() bool {
	return false
}

/*
Sys is a method which allows you to obtain the underlying data source of the file information, which does not exist.

Example:

	underlying := info.Sys()
*/
func (info virtualFileSystemReaderAtFileInfo) Sys() any {
	return nil
}

/*
indexZipVirtualFileSystemEntries is a method which allows you to record every file entry of an opened ZIP archive in
a mount so that later reads are a direct map lookup. Directory entries, and entries whose name normalizes to an empty
path, are skipped because the virtual file system indexes files only. Only the central directory already parsed by
the ZIP reader is consulted, so no entry data is read or decompressed. In addition, the following should be noted:

  - The mount is modified in place. Its zipEntries and paths maps must already be allocated.

  - If an entry is encrypted, the supplied password is applied to it here. Both traditional ZipCrypto and WinZip AES
    encrypted entries are supported for reading, and a wrong password only surfaces when the entry is read.

Example:

	indexZipVirtualFileSystemEntries(mount, zipReader, "secret")
*/
func indexZipVirtualFileSystemEntries(mount *virtualFileSystemMount, zipReader *zip.Reader, password string) {
	for _, currentFile := range zipReader.File {
		if currentFile.FileInfo().IsDir() {
			continue
		}
		entryPath := normalizeVirtualFileSystemPath(currentFile.Name)
		if entryPath == "" {
			continue
		}
		if currentFile.IsEncrypted() {
			currentFile.SetPassword(password)
		}
		mount.zipEntries[entryPath] = currentFile
		mount.paths[entryPath] = struct{}{}
	}
}

/*
newZipVirtualFileSystemMount is a constructor which allows you to open a ZIP archive through an io.ReaderAt and build a
fully indexed mount record for it. Only the end of central directory record and the central directory are read at
mount time, and each entry's compressed bytes are read later only when that entry is requested. In addition, the
following should be noted:

  - The returned mount keeps the supplied reader for its entire lifetime and never closes it. A caller that owns a
    handle which must be released on unmount, such as an open file, sets the mount's closer itself.

Example:

	mount, err := newZipVirtualFileSystemMount("assets.zip", bytes.NewReader(archiveData), int64(len(archiveData)), "")
*/
func newZipVirtualFileSystemMount(source string, reader io.ReaderAt, size int64, password string) (*virtualFileSystemMount, error) {
	zipReader, err := zip.NewReader(reader, size)
	if err != nil {
		return nil, err
	}
	mount := &virtualFileSystemMount{
		source:     source,
		mountType:  constants.VirtualFileSystemZip,
		paths:      make(map[string]struct{}),
		zipReader:  zipReader,
		zipEntries: make(map[string]*zip.File),
	}
	indexZipVirtualFileSystemEntries(mount, zipReader, password)
	return mount, nil
}

/*
newZipVirtualFileSystemMountFromPath is a constructor which allows you to open a ZIP archive on the local file system
and build a fully indexed mount record for it. The file is opened once and served to newZipVirtualFileSystemMount as
its io.ReaderAt, so path mounts and in memory mounts share the same ZIP core. In addition, the following should be
noted:

  - The returned mount keeps the file open for its entire lifetime and records it as the mount's closer, so it is
    closed by UnmountVirtualFileSystem or UnmountVirtualFileSystemArchive.

  - If the file is not a valid ZIP archive, the file is closed again before the error is returned.

Example:

	mount, err := newZipVirtualFileSystemMountFromPath("assets.zip", "secret")
*/
func newZipVirtualFileSystemMountFromPath(archivePath string, password string) (*virtualFileSystemMount, error) {
	archiveFile, err := os.Open(archivePath)
	if err != nil {
		return nil, err
	}
	fileInfo, err := archiveFile.Stat()
	if err != nil {
		_ = archiveFile.Close()
		return nil, err
	}
	mount, err := newZipVirtualFileSystemMount(archivePath, archiveFile, fileInfo.Size(), password)
	if err != nil {
		_ = archiveFile.Close()
		return nil, err
	}
	mount.closer = archiveFile
	return mount, nil
}

/*
newRarVirtualFileSystemMount is a constructor which allows you to open a RAR archive and build a fully indexed mount
record for it. The archive is opened once through rardecode, which walks the block headers and records the byte offset
of every entry, so a later read seeks straight to the requested file instead of scanning from the start. In addition,
the following should be noted:

  - The volumePath is the name rardecode opens. For a path mount it is the archive path itself and the operating
    system file system is used, which also allows multi volume archives. For an in memory mount it is the single
    volume name of a virtualFileSystemReaderAtFileSystem supplied through additionalOptions.

  - A wrong password fails here at open time when the archive uses encrypted headers, and otherwise fails later when
    the file is read.

  - Solid RAR archives do not support direct entry access. Reads of entries in a solid archive will return an error.
    Create the archive without the solid option to allow random access.

Example:

	mount, err := newRarVirtualFileSystemMount("assets.rar", "assets.rar", "secret")
*/
func newRarVirtualFileSystemMount(source string, volumePath string, password string, additionalOptions ...rardecode.Option) (*virtualFileSystemMount, error) {
	options := append([]rardecode.Option{rardecode.Password(password)}, additionalOptions...)
	rarFileSystem, err := rardecode.OpenFS(volumePath, options...)
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
newRarVirtualFileSystemMountFromReaderAt is a constructor which allows you to open a RAR archive held behind an
io.ReaderAt and build a fully indexed mount record for it. The reader is exposed to rardecode as a one file fs.FS whose
file is seekable, so rardecode reads only the block headers at mount time and seeks directly to an entry's packed data
when that entry is read. In addition, the following should be noted:

  - Multi volume RAR archives are not supported from memory, because only the single supplied volume exists.

Example:

	mount, err := newRarVirtualFileSystemMountFromReaderAt("assets", bytes.NewReader(archiveData), size, "")
*/
func newRarVirtualFileSystemMountFromReaderAt(source string, reader io.ReaderAt, size int64, password string) (*virtualFileSystemMount, error) {
	fileSystem := &virtualFileSystemReaderAtFileSystem{
		name:   virtualFileSystemMemoryRarVolumeName,
		reader: reader,
		size:   size,
	}
	return newRarVirtualFileSystemMount(source, virtualFileSystemMemoryRarVolumeName, password, rardecode.FileSystem(fileSystem))
}

/*
addVirtualFileSystemMount is a method which allows you to add a fully indexed mount to the top of the virtual file
system. The mount receives the next mount identifier, is appended after every existing mount so that it shadows them
for shared paths, and the union index is rebuilt. In addition, the following should be noted:

  - This method acquires the virtualFileSystemMutex write lock itself, so the caller must not already hold it.

Example:

	addVirtualFileSystemMount(mount)
*/
func addVirtualFileSystemMount(mount *virtualFileSystemMount) {
	virtualFileSystemMutex.Lock()
	defer virtualFileSystemMutex.Unlock()
	virtualFileSystemNextId++
	mount.id = virtualFileSystemNextId
	virtualFileSystemMounts = append(virtualFileSystemMounts, mount)
	rebuildVirtualFileSystemIndex()
}

/*
MountVirtualFileSystem is a method which allows you to add a ZIP or RAR archive on the local file system to the virtual
file system as an additional layer. Any number of archives of either format can be mounted at the same time, together
with embedded file systems and archives mounted from memory, and their contents are unioned into a single directory
structure that GetFileData, GetTextFileData, GetFileReader, FileExists and the image and font loaders all read from.
In addition, the following should be noted:

  - Mounting does not replace earlier mounts. When a file exists at the same path in more than one mount, the most
    recently mounted one wins and shadows the earlier ones. Mounting the same path twice adds a second layer, and
    UnmountVirtualFileSystemArchive removes every layer mounted from that path.

  - The archive format is detected automatically by attempting to open it first as a ZIP and then as a RAR. If neither
    succeeds, an error describing both failures is returned and nothing is mounted.

  - If the archive is password protected, the password must be supplied here at mount time. If a non empty
    scrambleKey is supplied, the password is treated as a scrambled password and is decoded with that key first.

  - Only the archive's index is read at mount time and each file is extracted when it is requested. A ZIP archive's
    file stays open until it is unmounted, while a RAR archive is reopened for each read.

  - To mount an archive that is already in memory, such as one embedded in the executable with go:embed, use
    MountVirtualFileSystemFromBytes or MountVirtualFileSystemFromReaderAt instead. Neither copies the archive.

Example:

	err := MountVirtualFileSystem("assets.zip", "pwd", "")
*/
func MountVirtualFileSystem(archivePath string, password string, scrambleKey string) error {
	if scrambleKey != "" {
		password = getUnscrambledPassword(password, scrambleKey)
	}
	mount, zipError := newZipVirtualFileSystemMountFromPath(archivePath, password)
	if zipError != nil {
		var rarError error
		mount, rarError = newRarVirtualFileSystemMount(archivePath, archivePath, password)
		if rarError != nil {
			return errors.New(fmt.Sprintf("Failed to open or decode '%s' as a ZIP (%s) or RAR (%s) archive.", archivePath, zipError.Error(), rarError.Error()))
		}
	}
	addVirtualFileSystemMount(mount)
	return nil
}

/*
MountVirtualFileSystemFromReaderAt is a method which allows you to add a ZIP or RAR archive that is held behind an
io.ReaderAt to the virtual file system as an additional layer, without copying it and without reading it in full. Only
the archive's index is read at mount time, and each file's compressed bytes are read from the reader only when that
file is requested. The mount takes part in the same union and shadowing rules as every other mount. In addition, the
following should be noted:

  - The reader is kept, not copied, for the whole life of the mount and is never closed by the virtual file system.
    The caller must keep it valid, and its contents unchanged, until the mount is released with
    UnmountVirtualFileSystemArchive or UnmountVirtualFileSystem. Its ReadAt must be safe for concurrent use, as the
    io.ReaderAt contract requires, because files may be read from several goroutines at once.

  - The name identifies the mount in place of an archive path. Pass it to UnmountVirtualFileSystemArchive to release
    the mount, and it appears in error messages. Mounting a name that is already mounted adds another layer exactly as
    mounting the same path twice does, and unmounting that name removes every layer that uses it. Names share one
    namespace with archive paths.

  - The archive format is detected by trying ZIP and then RAR. If neither succeeds, an error naming the mount and
    describing both failures is returned and nothing is mounted. A RAR archive must be a single volume, and entries
    of a solid RAR archive cannot be read individually, exactly as with a path mount.

  - Password handling matches MountVirtualFileSystem. If a non empty scrambleKey is supplied, the password is
    unscrambled with that key first.

  - An empty name, a nil reader, or a size of zero or less is rejected with an error and nothing is mounted.

  - FileExists answers from the index alone and never reads from the reader.

Example:

	archiveFile, _ := os.Open("assets.zip")
	fileInfo, _ := archiveFile.Stat()
	err := MountVirtualFileSystemFromReaderAt("assets", archiveFile, fileInfo.Size(), "", "")
*/
func MountVirtualFileSystemFromReaderAt(name string, reader io.ReaderAt, size int64, password string, scrambleKey string) error {
	if name == "" {
		return errors.New("Failed to mount a virtual file system from memory: the mount name must not be empty.")
	}
	if reader == nil {
		return errors.New(fmt.Sprintf("Failed to mount the virtual file system '%s' from memory: the reader must not be nil.", name))
	}
	if size <= 0 {
		return errors.New(fmt.Sprintf("Failed to mount the virtual file system '%s' from memory: the archive size must be greater than zero, but %d was supplied.", name, size))
	}
	if scrambleKey != "" {
		password = getUnscrambledPassword(password, scrambleKey)
	}
	mount, zipError := newZipVirtualFileSystemMount(name, reader, size, password)
	if zipError != nil {
		var rarError error
		mount, rarError = newRarVirtualFileSystemMountFromReaderAt(name, reader, size, password)
		if rarError != nil {
			return errors.New(fmt.Sprintf("Failed to open or decode '%s' as a ZIP (%s) or RAR (%s) archive.", name, zipError.Error(), rarError.Error()))
		}
	}
	addVirtualFileSystemMount(mount)
	return nil
}

/*
MountVirtualFileSystemFromBytes is a method which allows you to add a ZIP or RAR archive held in a byte slice to the
virtual file system as an additional layer, without copying the slice. This is the natural way to mount an archive
embedded in the executable with go:embed, whose pages the operating system loads from disk only when they are touched.
It behaves exactly like MountVirtualFileSystemFromReaderAt over bytes.NewReader(archiveData). In addition, the
following should be noted:

  - The slice is referenced, not copied, for the whole life of the mount. The caller must not modify it until the
    mount is released with UnmountVirtualFileSystemArchive or UnmountVirtualFileSystem. An embedded archive is never
    modified, so it is always safe.

  - Only the archive's index is read at mount time, and each file's compressed bytes are read only when that file is
    requested. FileExists answers from the index alone and reads nothing.

  - An empty or nil slice, or an empty name, is rejected with an error and nothing is mounted.

Example:

	//go:embed assets.zip
	var assetsArchive []byte

	err := MountVirtualFileSystemFromBytes("assets", assetsArchive, "", "")
*/
func MountVirtualFileSystemFromBytes(name string, archiveData []byte, password string, scrambleKey string) error {
	return MountVirtualFileSystemFromReaderAt(name, bytes.NewReader(archiveData), int64(len(archiveData)), password, scrambleKey)
}

/*
MountEmbeddedFileSystem is a method which allows you to add an embedded file system as an additional layer of the
virtual file system. The embedded file system is walked once at mount time and every file it contains is indexed, so
its contents take part in the same union and shadowing rules as mounted archives and are visible to GetFileData,
GetTextFileData, GetFileReader and FileExists. An error is returned if the embedded file system could not be walked.
In addition, the following should be noted:

  - Mounting is additive. An embedded file system mounted after another mount shadows it for any shared file path.

  - Only files are indexed. Paths are the full paths inside the embedded file system, including the embedded
    directory name, and a path that only names a directory is not found by reads or by FileExists.

  - To mount a ZIP or RAR archive embedded as a byte slice rather than as a directory tree, use
    MountVirtualFileSystemFromBytes.

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
ZIP archive mounted from a path holds its open file as the mount's closer and that file is closed here. RAR mounts,
embedded mounts, and mounts created from a caller supplied io.ReaderAt hold nothing the virtual file system owns, so
they require no action and the caller's reader is never closed. In addition, the following should be noted:

  - The caller must already hold the virtualFileSystemMutex write lock.

  - The closer is cleared after closing, so calling this method twice for the same mount is harmless.

Example:

	closeVirtualFileSystemMount(mount)
*/
func closeVirtualFileSystemMount(mount *virtualFileSystemMount) {
	if mount.closer != nil {
		_ = mount.closer.Close()
		mount.closer = nil
	}
}

/*
UnmountVirtualFileSystem is a method which allows you to reset the virtual file system to a completely unmounted state.
Every mounted archive, archive mounted from memory, and embedded file system is removed, all archive files opened by
the virtual file system are closed, and subsequent reads and FileExists checks go directly to the local file system
again. In addition, the following should be noted:

  - A reader or byte slice supplied to MountVirtualFileSystemFromReaderAt or MountVirtualFileSystemFromBytes is never
    closed. Once this method returns, the virtual file system no longer references it.

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

  - Archives are matched by the exact path that was passed to MountVirtualFileSystem, or the exact name that was
    passed to MountVirtualFileSystemFromReaderAt or MountVirtualFileSystemFromBytes. If the same path or name was
    mounted more than once, every instance of it is removed. An error is returned if nothing matches.

  - A reader or byte slice supplied to a mount from memory is never closed. Once this method returns, the virtual
    file system no longer references it.

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
GetFileData is a method which allows you to read the raw bytes of any file through the virtual file system. This is
the general purpose entry point for loading file types the virtual file system has no dedicated loader for, such as
audio clips, save data, or configuration files, using the exact same mount and shadowing rules that back images and
fonts. In addition, the following should be noted:

  - When one or more archives or embedded file systems are mounted, only the virtual file system is consulted. A path
    that is not found in any mount returns an error even if it exists on the local file system.

  - When nothing is mounted, the file is read directly from the local file system, so this method behaves identically
    with or without a virtual file system in place.

  - Use FileExists to ask whether a file is available under these same rules without reading it.

Example:

	data, err := GetFileData("sounds/click.wav")
*/
func GetFileData(filePath string) ([]byte, error) {
	return getFileDataFromFileSystem(filePath)
}

/*
GetTextFileData is a method which allows you to read the contents of any file through the virtual file system as a
string. This is useful for loading arbitrary text based file formats, such as configuration files, scripts, or save
data, through the same mount and shadowing rules that back images and fonts.

Example:

	content, err := GetTextFileData("config/settings.json")
*/
func GetTextFileData(filePath string) (string, error) {
	return getTextFromFileSystem(filePath)
}

/*
GetFileReader is a method which allows you to obtain a streaming reader for any file through the virtual file system.
This lets your own loaders consume arbitrary file formats without first deciding whether to read the whole file into
memory. In addition, the following should be noted:

  - The returned reader serves an in memory copy of the file contents rather than a live handle into the archive, so
    it remains valid even after the owning mount is later released with UnmountVirtualFileSystem.

Example:

	reader, err := GetFileReader("data/save.bin")
	defer reader.Close()
*/
func GetFileReader(filePath string) (io.ReadCloser, error) {
	return getFileReaderFromFileSystem(filePath)
}

/*
FileExists is a method which allows you to check whether a file is available through the virtual file system without
reading it. It answers exactly what GetFileData would do with the same path, so it is the replacement for os.Stat when
files may live inside a mounted archive or embedded file system. In addition, the following should be noted:

  - When one or more archives or embedded file systems are mounted, only the virtual file system is consulted. The
    path is normalized exactly as reads normalize it, so spellings such as a leading "./" or backslash separators
    give the same answer, and true is returned only for a file that a mount serves. A file that exists only on the
    local file system returns false.

  - When nothing is mounted, true is returned only when the path names a regular file on the local file system.
    Directories, missing paths, and paths that cannot be inspected return false.

  - A path that only names a directory always returns false, because mounts index files only.

  - The check is a lookup in the union index. No file is read or decompressed and no archive bytes are touched, so it
    costs nothing even for an archive mounted from memory. It is safe to call while mounts are added or removed from
    other goroutines.

  - A file that exists but cannot be read, such as a local file without read permission or an encrypted entry mounted
    with the wrong password, still returns true because FileExists does not attempt the read.

Example:

	hasPortrait := FileExists("portraits/hero.png")
*/
func FileExists(filePath string) bool {
	virtualFileSystemMutex.RLock()
	mount, _, isMounted := lookupVirtualFileSystemMount(filePath)
	virtualFileSystemMutex.RUnlock()
	if isMounted {
		return mount != nil
	}
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return fileInfo.Mode().IsRegular()
}

/*
lookupVirtualFileSystemMount is a method which allows you to resolve a requested path against the union index. It is
the single place where a lookup path is normalized, so file reads and FileExists can never disagree about which
spellings of a path refer to the same file. The mount that serves the file is returned together with the normalized
path, or a nil mount when no mount serves it, and isMounted reports whether anything is mounted at all, which decides
whether the local file system may be consulted. In addition, the following should be noted:

  - The caller must already hold at least the virtualFileSystemMutex read lock.

Example:

	mount, normalizedName, isMounted := lookupVirtualFileSystemMount("./myFolder\\sample.png")
*/
func lookupVirtualFileSystemMount(filePath string) (*virtualFileSystemMount, string, bool) {
	normalizedName := normalizeVirtualFileSystemPath(filePath)
	return virtualFileSystemIndex[normalizedName], normalizedName, len(virtualFileSystemMounts) > 0
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
	virtualFileSystemMutex.RLock()
	defer virtualFileSystemMutex.RUnlock()
	mount, normalizedName, isMounted := lookupVirtualFileSystemMount(fileName)
	if mount != nil {
		fileData, err := readFileFromVirtualFileSystemMount(mount, normalizedName)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Could not read the file '%s' from the virtual file system '%s': %s", fileName, mount.source, err.Error()))
		}
		return fileData, nil
	}
	if isMounted {
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
