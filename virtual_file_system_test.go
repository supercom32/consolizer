package consolizer

import (
	"archive/zip"
	"bytes"
	"embed"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/nwaples/rardecode/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"
)

const (
	BASE_DIRECTORY = "./test_data/virtual_file_systems/"
)

const (
	EMBEDDED_TEST_DIRECTORY = "test_data/virtual_file_systems/embedded"
)

//go:embed test_data/virtual_file_systems/embedded
var embeddedTestFileSystem embed.FS

/*
writeTestZipArchive is a test helper which builds an unencrypted ZIP archive on disk from a map of archive entry paths
to their string contents and returns the full path to the created archive. Any failure while creating the archive is
reported through the supplied testing handle.
*/
func writeTestZipArchive(test *testing.T, directory string, archiveName string, files map[string]string) string {
	test.Helper()
	archivePath := filepath.Join(directory, archiveName)
	archiveFile, err := os.Create(archivePath)
	assert.NoErrorf(test, err, "Failed to create the test archive file.")
	zipWriter := zip.NewWriter(archiveFile)
	for entryName, entryContent := range files {
		entryWriter, err := zipWriter.Create(entryName)
		assert.NoErrorf(test, err, "Failed to create the archive entry '%s'.", entryName)
		_, err = entryWriter.Write([]byte(entryContent))
		assert.NoErrorf(test, err, "Failed to write the archive entry '%s'.", entryName)
	}
	assert.NoErrorf(test, zipWriter.Close(), "Failed to finalize the test archive.")
	assert.NoErrorf(test, archiveFile.Close(), "Failed to close the test archive file.")
	return archivePath
}

/*
TestGetScrambledPassword is a test which verifies that passwords can be correctly scrambled using a base64 encoded XOR
result.

Example:
    Expected Inputs:
        Password string "SamplePassword" and scramble key "SampleScrambleKey".
    Expected Outputs:
        The resulting base64 string matches "awVdQ1tUYVAVR0VXEwE=".
*/
func TestGetScrambledPassword(test *testing.T) {
	obtainedResult := GetScrambledPassword("SamplePassword", "SampleScrambleKey")
	expectedResult := "awVdQ1tUYVAVR0VXEwE="
	assert.Equalf(test, obtainedResult, expectedResult, "The scrambled password did not match what was expected!")
}

/*
TestGetUnscrambledPassword is a test which verifies that scrambled passwords can be correctly unscrambled using the
original key.

Example:
    Expected Inputs:
        Scrambled base64 string "awVdQ1tUYVAVR0VXEwE=" and original scramble key.
    Expected Outputs:
        The decrypted string matches the original "SamplePassword".
*/
func TestGetUnscrambledPassword(test *testing.T) {
	obtainedResult := getUnscrambledPassword("awVdQ1tUYVAVR0VXEwE=", "SampleScrambleKey")
	expectedResult := "SamplePassword"
	assert.Equalf(test, obtainedResult, expectedResult, "The unscrambled password did not match what was expected!")
}

/*
TestMountVirtualFileSystem is a test which verifies mounting and file retrieval from ZIP and RAR virtual file systems,
including error handling for bad passwords, missing files, and unreadable archives. Each scenario fully unmounts before
the next because mounting is now additive rather than replacing.

Example:
    Expected Inputs:
        The archives valid.zip and valid.rar, the scrambled password "TAFDRw==", and the scramble key
        "SampleScrambleKey" which together decode to the password "test".
    Expected Outputs:
        Valid mounts with a correct key read "myFolder-1/myFolder-2/sample3.png" without error. A bad key produces a
        read error for ZIP and a mount error for RAR. A missing path produces a read error. invalid.zip produces a
        mount error.
*/
func TestMountVirtualFileSystem(test *testing.T) {
	defer UnmountVirtualFileSystem()
	scrambledPassword := "TAFDRw=="
	scrambleKey := "SampleScrambleKey"
	badScrambleKey := "SampleScrambleKey_BAD"

	// A valid ZIP archive mounts and its files can be read.
	assert.NoErrorf(test, MountVirtualFileSystem(BASE_DIRECTORY+"valid.zip", scrambledPassword, scrambleKey), "Failed to mount a valid ZIP file system!")
	_, err := getImageFromFileSystem("myFolder-1/myFolder-2/sample3.png")
	assert.NoErrorf(test, err, "Failed to obtain image data from the ZIP file system.")
	UnmountVirtualFileSystem()

	// A valid RAR archive mounts and its files can be read.
	assert.NoErrorf(test, MountVirtualFileSystem(BASE_DIRECTORY+"valid.rar", scrambledPassword, scrambleKey), "Failed to mount a valid RAR file system!")
	_, err = getImageFromFileSystem("myFolder-1/myFolder-2/sample3.png")
	assert.NoErrorf(test, err, "Failed to obtain image data from the RAR file system.")
	UnmountVirtualFileSystem()

	// A bad password lets the ZIP mount succeed but fails when the encrypted entry is read.
	assert.NoErrorf(test, MountVirtualFileSystem(BASE_DIRECTORY+"valid.zip", scrambledPassword, badScrambleKey), "Expected a bad ZIP password to still allow mounting!")
	_, err = getImageFromFileSystem("myFolder-1/myFolder-2/sample3.png")
	assert.Errorf(test, err, "Expected an error reading from the ZIP file system with a bad password.")
	UnmountVirtualFileSystem()

	// A bad password fails the RAR mount outright because the archive uses encrypted headers.
	err = MountVirtualFileSystem(BASE_DIRECTORY+"valid.rar", scrambledPassword, badScrambleKey)
	assert.Errorf(test, err, "Expected mounting a RAR file system with a bad password to fail.")
	UnmountVirtualFileSystem()

	// A request for a path that is not in the archive returns an error for ZIP.
	assert.NoErrorf(test, MountVirtualFileSystem(BASE_DIRECTORY+"valid.zip", scrambledPassword, scrambleKey), "Failed to mount a valid ZIP file system!")
	_, err = getImageFromFileSystem("myFolder-1/myFolder-2/sample3.png_BAD")
	assert.Errorf(test, err, "Expected an error retrieving a file from the ZIP file system that does not exist.")
	UnmountVirtualFileSystem()

	// A request for a path that is not in the archive returns an error for RAR.
	assert.NoErrorf(test, MountVirtualFileSystem(BASE_DIRECTORY+"valid.rar", scrambledPassword, scrambleKey), "Failed to mount a valid RAR file system!")
	_, err = getImageFromFileSystem("myFolder-1/myFolder-2/sample3.png_BAD")
	assert.Errorf(test, err, "Expected an error retrieving a file from the RAR file system that does not exist.")
	UnmountVirtualFileSystem()

	// An archive that is neither a valid ZIP nor a valid RAR fails to mount.
	err = MountVirtualFileSystem(BASE_DIRECTORY+"invalid.zip", scrambledPassword, scrambleKey)
	assert.Errorf(test, err, "Expected mounting an invalid archive to fail.")
	UnmountVirtualFileSystem()
}

/*
TestVirtualFileSystemUnionAndClobber is a test which verifies that multiple archives can be mounted at once, that their
contents are unioned, that a later mount shadows an earlier one for a shared path, and that removing the later mount
makes the shadowed file visible again.

Example:
    Expected Inputs:
        Archive A with {"shared/file.txt": "AAA", "only_a.txt": "a"} and archive B with
        {"shared/file.txt": "BBB", "only_b.txt": "b"}, mounted in that order.
    Expected Outputs:
        "shared/file.txt" reads "BBB" while both are mounted and "only_a.txt" still reads "a". After unmounting B,
        "shared/file.txt" reads "AAA" and "only_b.txt" errors. After unmounting A every read errors.
*/
func TestVirtualFileSystemUnionAndClobber(test *testing.T) {
	defer UnmountVirtualFileSystem()
	temporaryDirectory := test.TempDir()
	archivePathA := writeTestZipArchive(test, temporaryDirectory, "a.zip", map[string]string{
		"shared/file.txt": "AAA",
		"only_a.txt":      "a",
	})
	archivePathB := writeTestZipArchive(test, temporaryDirectory, "b.zip", map[string]string{
		"shared/file.txt": "BBB",
		"only_b.txt":      "b",
	})

	assert.NoErrorf(test, MountVirtualFileSystem(archivePathA, "", ""), "Failed to mount archive A.")
	assert.NoErrorf(test, MountVirtualFileSystem(archivePathB, "", ""), "Failed to mount archive B.")

	sharedContent, err := getTextFromFileSystem("shared/file.txt")
	assert.NoErrorf(test, err, "Failed to read the shared file from the union.")
	assert.Equalf(test, "BBB", sharedContent, "The later mount did not shadow the earlier mount for a shared path.")

	onlyAContent, err := getTextFromFileSystem("only_a.txt")
	assert.NoErrorf(test, err, "Failed to read a file unique to the earlier mount.")
	assert.Equalf(test, "a", onlyAContent, "A file unique to the earlier mount was not visible through the union.")

	assert.NoErrorf(test, UnmountVirtualFileSystemArchive(archivePathB), "Failed to unmount archive B.")

	sharedContent, err = getTextFromFileSystem("shared/file.txt")
	assert.NoErrorf(test, err, "Failed to read the shared file after unmounting the shadowing archive.")
	assert.Equalf(test, "AAA", sharedContent, "Unmounting the later archive did not reveal the shadowed file.")

	_, err = getTextFromFileSystem("only_b.txt")
	assert.Errorf(test, err, "Expected a file unique to the unmounted archive to no longer be found.")

	assert.NoErrorf(test, UnmountVirtualFileSystemArchive(archivePathA), "Failed to unmount archive A.")

	_, err = getTextFromFileSystem("shared/file.txt")
	assert.Errorf(test, err, "Expected every read to fail once all archives are unmounted.")
}

/*
TestUnmountVirtualFileSystemArchive is a test which verifies that a single archive can be removed from the virtual file
system while other mounts stay intact, and that unmounting an archive that was never mounted returns an error.

Example:
    Expected Inputs:
        Archive A with {"only_a.txt": "a"} and archive B with {"only_b.txt": "b"}, both mounted.
    Expected Outputs:
        After unmounting A, "only_b.txt" still reads "b" and "only_a.txt" errors. Unmounting the path "missing.zip"
        returns an error.
*/
func TestUnmountVirtualFileSystemArchive(test *testing.T) {
	defer UnmountVirtualFileSystem()
	temporaryDirectory := test.TempDir()
	archivePathA := writeTestZipArchive(test, temporaryDirectory, "a.zip", map[string]string{"only_a.txt": "a"})
	archivePathB := writeTestZipArchive(test, temporaryDirectory, "b.zip", map[string]string{"only_b.txt": "b"})

	assert.NoErrorf(test, MountVirtualFileSystem(archivePathA, "", ""), "Failed to mount archive A.")
	assert.NoErrorf(test, MountVirtualFileSystem(archivePathB, "", ""), "Failed to mount archive B.")

	assert.NoErrorf(test, UnmountVirtualFileSystemArchive(archivePathA), "Failed to unmount archive A by path.")

	remainingContent, err := getTextFromFileSystem("only_b.txt")
	assert.NoErrorf(test, err, "Expected the archive that was not unmounted to remain readable.")
	assert.Equalf(test, "b", remainingContent, "The remaining archive returned the wrong content.")

	_, err = getTextFromFileSystem("only_a.txt")
	assert.Errorf(test, err, "Expected files from the unmounted archive to no longer be found.")

	err = UnmountVirtualFileSystemArchive(filepath.Join(temporaryDirectory, "missing.zip"))
	assert.Errorf(test, err, "Expected unmounting an archive that was never mounted to return an error.")
}

/*
TestVirtualFileSystemMixedFormatUnion is a test which verifies that a ZIP archive and a RAR archive can be mounted at
the same time and that a single archive can be removed from the union by path regardless of its format.

Example:
    Expected Inputs:
        valid.zip and valid.rar mounted together with the password "test" derived from the scrambled password
        "TAFDRw==" and the scramble key "SampleScrambleKey". Both archives contain
        "myFolder-1/myFolder-2/sample3.png".
    Expected Outputs:
        The shared path reads without error while both are mounted and after the RAR is unmounted by path. After the
        ZIP is also unmounted the read errors.
*/
func TestVirtualFileSystemMixedFormatUnion(test *testing.T) {
	defer UnmountVirtualFileSystem()
	scrambledPassword := "TAFDRw=="
	scrambleKey := "SampleScrambleKey"
	sharedPath := "myFolder-1/myFolder-2/sample3.png"

	assert.NoErrorf(test, MountVirtualFileSystem(BASE_DIRECTORY+"valid.zip", scrambledPassword, scrambleKey), "Failed to mount the ZIP archive.")
	assert.NoErrorf(test, MountVirtualFileSystem(BASE_DIRECTORY+"valid.rar", scrambledPassword, scrambleKey), "Failed to mount the RAR archive.")

	_, err := getImageFromFileSystem(sharedPath)
	assert.NoErrorf(test, err, "Failed to read a file present in both mounted archives.")

	assert.NoErrorf(test, UnmountVirtualFileSystemArchive(BASE_DIRECTORY+"valid.rar"), "Failed to unmount the RAR archive by path.")
	_, err = getImageFromFileSystem(sharedPath)
	assert.NoErrorf(test, err, "Expected the ZIP archive to still serve the shared file after the RAR was unmounted.")

	assert.NoErrorf(test, UnmountVirtualFileSystemArchive(BASE_DIRECTORY+"valid.zip"), "Failed to unmount the ZIP archive by path.")
	_, err = getImageFromFileSystem(sharedPath)
	assert.Errorf(test, err, "Expected the read to fail once both archives are unmounted.")
}

/*
TestGetFileDataFromLocalFileSystem is a test which verifies that file data can be correctly read directly from the
local host file system.

Example:
    Expected Inputs:
        Path to an existing RAR file on the local disk.
    Expected Outputs:
        File bytes are successfully retrieved without error.
*/
func TestGetFileDataFromLocalFileSystem(test *testing.T) {
	_, err := getFileDataFromLocalFileSystem(BASE_DIRECTORY + "valid.rar")
	assert.NoErrorf(test, err, "Did not expect an error reading a file that should exist!")
}

/*
TestGetTextFromFileSystem is a test which verifies that text content can be successfully retrieved from a file on the
file system.

Example:
    Expected Inputs:
        Path to an existing text file "text_file.txt".
    Expected Outputs:
        The string content of the file is correctly returned.
*/
func TestGetTextFromFileSystem(test *testing.T) {
	_, err := getTextFromFileSystem(BASE_DIRECTORY + "text_file.txt")
	assert.NoErrorf(test, err, "Did not expect an error reading a text file that should exist!")
}

/*
TestGetFileData is a test which verifies that the public GetFileData method can read the bytes of an arbitrary file
type from both the local file system and a mounted archive, using the same union and shadowing rules as every other
virtual file system accessor.

Example:
    Expected Inputs:
        The local text file "text_file.txt", and a ZIP archive mounted with the entry "data/save.dat" containing the
        bytes "SAVE-DATA-V1".
    Expected Outputs:
        Both reads succeed and the archive read returns exactly "SAVE-DATA-V1".
*/
func TestGetFileData(test *testing.T) {
	defer UnmountVirtualFileSystem()
	_, err := GetFileData(BASE_DIRECTORY + "text_file.txt")
	assert.NoErrorf(test, err, "Did not expect an error reading an arbitrary file from the local file system.")

	temporaryDirectory := test.TempDir()
	archivePath := writeTestZipArchive(test, temporaryDirectory, "arbitrary.zip", map[string]string{
		"data/save.dat": "SAVE-DATA-V1",
	})
	assert.NoErrorf(test, MountVirtualFileSystem(archivePath, "", ""), "Failed to mount the archive.")

	fileData, err := GetFileData("data/save.dat")
	assert.NoErrorf(test, err, "Failed to read an arbitrary file type from the mounted virtual file system.")
	assert.Equalf(test, "SAVE-DATA-V1", string(fileData), "The arbitrary file data did not match what was expected.")

	_, err = GetFileData("data/missing.dat")
	assert.Errorf(test, err, "Expected an error reading a file that does not exist in the mounted archive.")
}

/*
TestGetTextFileData is a test which verifies that the public GetTextFileData method returns the string contents of an
arbitrary file from both the local file system and a mounted archive.

Example:
    Expected Inputs:
        The local text file "text_file.txt", and a ZIP archive mounted with the entry "config/settings.ini"
        containing "key=value".
    Expected Outputs:
        Both reads succeed and the archive read returns exactly "key=value".
*/
func TestGetTextFileData(test *testing.T) {
	defer UnmountVirtualFileSystem()
	_, err := GetTextFileData(BASE_DIRECTORY + "text_file.txt")
	assert.NoErrorf(test, err, "Did not expect an error reading an arbitrary text file from the local file system.")

	temporaryDirectory := test.TempDir()
	archivePath := writeTestZipArchive(test, temporaryDirectory, "arbitrary.zip", map[string]string{
		"config/settings.ini": "key=value",
	})
	assert.NoErrorf(test, MountVirtualFileSystem(archivePath, "", ""), "Failed to mount the archive.")

	textContent, err := GetTextFileData("config/settings.ini")
	assert.NoErrorf(test, err, "Failed to read arbitrary text content from the mounted virtual file system.")
	assert.Equalf(test, "key=value", textContent, "The arbitrary text content did not match what was expected.")
}

/*
TestGetFileReader is a test which verifies that the public GetFileReader method returns a readable stream of an
arbitrary file's contents from both the local file system and a mounted archive, and that the caller is responsible
for closing it.

Example:
    Expected Inputs:
        The local text file "text_file.txt", and a ZIP archive mounted with the entry "audio/click.raw" containing
        the bytes "RAW-AUDIO-BYTES".
    Expected Outputs:
        Both readers open without error, and reading the archive entry to completion yields exactly
        "RAW-AUDIO-BYTES".
*/
func TestGetFileReader(test *testing.T) {
	defer UnmountVirtualFileSystem()
	localReader, err := GetFileReader(BASE_DIRECTORY + "text_file.txt")
	assert.NoErrorf(test, err, "Did not expect an error opening a reader for a local file.")
	assert.NoErrorf(test, localReader.Close(), "Failed to close the local file reader.")

	temporaryDirectory := test.TempDir()
	archivePath := writeTestZipArchive(test, temporaryDirectory, "arbitrary.zip", map[string]string{
		"audio/click.raw": "RAW-AUDIO-BYTES",
	})
	assert.NoErrorf(test, MountVirtualFileSystem(archivePath, "", ""), "Failed to mount the archive.")

	archiveReader, err := GetFileReader("audio/click.raw")
	assert.NoErrorf(test, err, "Failed to open a reader for an arbitrary file in the mounted virtual file system.")
	streamedData, err := io.ReadAll(archiveReader)
	assert.NoErrorf(test, err, "Failed to read the full contents of the streamed file.")
	assert.NoErrorf(test, archiveReader.Close(), "Failed to close the archive file reader.")
	assert.Equalf(test, "RAW-AUDIO-BYTES", string(streamedData), "The streamed file data did not match what was expected.")
}

/*
readRangeType is a class which allows you to describe one ReadAt call recorded by a recordingReaderAt, as the offset
the read started at and the number of bytes requested.
*/
type readRangeType struct {
	offset int64
	length int64
}

/*
recordingReaderAt is a class which allows you to wrap an io.ReaderAt and record every byte range that is read through
it, so a test can prove exactly which parts of an archive the virtual file system touched. It also implements
io.Closer and counts calls to Close, so a test can prove the virtual file system never closes a reader it was handed.
*/
type recordingReaderAt struct {
	reader     io.ReaderAt
	mutex      sync.Mutex
	ranges     []readRangeType
	closeCount int
}

/*
ReadAt is a method which allows you to read from the wrapped reader while recording the requested range.
*/
func (recorder *recordingReaderAt) ReadAt(buffer []byte, offset int64) (int, error) {
	recorder.mutex.Lock()
	recorder.ranges = append(recorder.ranges, readRangeType{offset: offset, length: int64(len(buffer))})
	recorder.mutex.Unlock()
	return recorder.reader.ReadAt(buffer, offset)
}

/*
Close is a method which allows you to count how many times something attempted to close the recorder.
*/
func (recorder *recordingReaderAt) Close() error {
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	recorder.closeCount++
	return nil
}

/*
takeRanges is a method which allows you to obtain every range recorded since the previous call and clear the record.
*/
func (recorder *recordingReaderAt) takeRanges() []readRangeType {
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	takenRanges := recorder.ranges
	recorder.ranges = nil
	return takenRanges
}

/*
getCloseCount is a method which allows you to obtain how many times Close was called on the recorder.
*/
func (recorder *recordingReaderAt) getCloseCount() int {
	recorder.mutex.Lock()
	defer recorder.mutex.Unlock()
	return recorder.closeCount
}

/*
newRecordingReaderAt is a constructor which allows you to create a recordingReaderAt over a byte slice.
*/
func newRecordingReaderAt(archiveData []byte) *recordingReaderAt {
	return &recordingReaderAt{reader: bytes.NewReader(archiveData)}
}

/*
testZipEntryType is a class which allows you to describe one entry of a ZIP archive built in memory by
buildTestZipArchiveInMemory. An entry whose name ends in a slash is written as a directory entry with no data.
*/
type testZipEntryType struct {
	name string
	data []byte
}

/*
buildTestZipArchiveInMemory is a method which allows you to build an unencrypted, uncompressed ZIP archive in
memory from an ordered list of entries. Entries are stored rather than deflated so their sizes in the archive are
predictable, which the read tracking tests rely on.
*/
func buildTestZipArchiveInMemory(test *testing.T, entries []testZipEntryType) []byte {
	test.Helper()
	var archiveBuffer bytes.Buffer
	zipWriter := zip.NewWriter(&archiveBuffer)
	for _, entry := range entries {
		entryWriter, err := zipWriter.CreateHeader(&zip.FileHeader{Name: entry.name, Method: zip.Store})
		assert.NoErrorf(test, err, "Failed to create the archive entry '%s'.", entry.name)
		if !strings.HasSuffix(entry.name, "/") {
			_, err = entryWriter.Write(entry.data)
			assert.NoErrorf(test, err, "Failed to write the archive entry '%s'.", entry.name)
		}
	}
	assert.NoErrorf(test, zipWriter.Close(), "Failed to finalize the in memory test archive.")
	return archiveBuffer.Bytes()
}

/*
getTestRandomBytes is a method which allows you to obtain a deterministic block of incompressible bytes.
*/
func getTestRandomBytes(seed int64, length int) []byte {
	randomBytes := make([]byte, length)
	rand.New(rand.NewSource(seed)).Read(randomBytes)
	return randomBytes
}

/*
getTestZipRecordRegions is a method which allows you to locate, for every entry of an in memory ZIP archive, the
byte region holding that entry's whole local record: its local header, its data, and any trailing data descriptor.
Each region runs from the entry's local header up to the next entry's local header, or up to the central directory for
the last entry. The offset of the central directory is returned as well.
*/
func getTestZipRecordRegions(test *testing.T, archiveData []byte) (map[string]readRangeType, int64) {
	test.Helper()
	zipReader, err := zip.NewReader(bytes.NewReader(archiveData), int64(len(archiveData)))
	assert.NoErrorf(test, err, "Failed to parse the in memory test archive.")
	endOfCentralDirectory := archiveData[len(archiveData)-22:]
	assert.Equalf(test, uint32(0x06054b50), binary.LittleEndian.Uint32(endOfCentralDirectory), "The test archive must end with a comment free end of central directory record.")
	centralDirectoryOffset := int64(binary.LittleEndian.Uint32(endOfCentralDirectory[16:]))
	type recordStartType struct {
		name  string
		start int64
	}
	var recordStarts []recordStartType
	for _, currentFile := range zipReader.File {
		dataOffset, err := currentFile.DataOffset()
		assert.NoErrorf(test, err, "Failed to locate the data of '%s'.", currentFile.Name)
		headerOffset := dataOffset - 30 - int64(len(currentFile.Name))
		assert.Equalf(test, uint32(0x04034b50), binary.LittleEndian.Uint32(archiveData[headerOffset:]), "The local header of '%s' was not where expected.", currentFile.Name)
		recordStarts = append(recordStarts, recordStartType{name: currentFile.Name, start: headerOffset})
	}
	sort.Slice(recordStarts, func(first int, second int) bool { return recordStarts[first].start < recordStarts[second].start })
	regions := make(map[string]readRangeType)
	for index, recordStart := range recordStarts {
		regionEnd := centralDirectoryOffset
		if index+1 < len(recordStarts) {
			regionEnd = recordStarts[index+1].start
		}
		regions[recordStart.name] = readRangeType{offset: recordStart.start, length: regionEnd - recordStart.start}
	}
	return regions, centralDirectoryOffset
}

/*
assertNoReadCoversWholeArchive is a method which allows you to assert that none of the recorded reads requested
the entire archive, or even half of it, in a single call.
*/
func assertNoReadCoversWholeArchive(test *testing.T, ranges []readRangeType, archiveSize int64, context string) {
	test.Helper()
	for _, currentRange := range ranges {
		assert.Lessf(test, currentRange.length, archiveSize/2, "%s: a single read of %d bytes at offset %d covers too much of the %d byte archive.", context, currentRange.length, currentRange.offset, archiveSize)
	}
}

/*
getTotalBytesRead is a method which allows you to add up the lengths of every recorded read.
*/
func getTotalBytesRead(ranges []readRangeType) int64 {
	var totalBytes int64
	for _, currentRange := range ranges {
		totalBytes += currentRange.length
	}
	return totalBytes
}

/*
getVirtualFileSystemMountCount is a method which allows you to obtain the number of mounts currently in place.
*/
func getVirtualFileSystemMountCount() int {
	virtualFileSystemMutex.RLock()
	defer virtualFileSystemMutex.RUnlock()
	return len(virtualFileSystemMounts)
}

/*
assertFileExistsMatchesGetFileData is a method which allows you to assert that FileExists returns the expected
answer for a path and that it agrees with GetFileData, meaning FileExists is true exactly when GetFileData succeeds.
*/
func assertFileExistsMatchesGetFileData(test *testing.T, filePath string, expectedResult bool, context string) {
	test.Helper()
	obtainedResult := FileExists(filePath)
	assert.Equalf(test, expectedResult, obtainedResult, "%s: FileExists(%q) returned the wrong answer.", context, filePath)
	_, err := GetFileData(filePath)
	assert.Equalf(test, obtainedResult, err == nil, "%s: FileExists(%q) returned %t but GetFileData returned the error %v.", context, filePath, obtainedResult, err)
}

/*
TestMountVirtualFileSystemFromBytes is a test which allows you to verify that a ZIP archive held in a byte slice can be
mounted without writing it to disk and that every public reader, including the image loader, serves its files.

Example:

	Expected Inputs:
		An in memory ZIP with {"data/save.dat": "SAVE-DATA-V1", "config/settings.ini": "key=value"} mounted as
		"memory-archive", then the bytes of valid.zip mounted as "memory-images" with the scrambled password
		"TAFDRw==" and the scramble key "SampleScrambleKey".
	Expected Outputs:
		GetFileData returns "SAVE-DATA-V1", GetTextFileData returns "key=value", GetFileReader streams
		"SAVE-DATA-V1", and LoadImage loads "myFolder-1/myFolder-2/sample3.png" without error.
*/
func TestMountVirtualFileSystemFromBytes(test *testing.T) {
	defer UnmountVirtualFileSystem()
	archiveData := buildTestZipArchiveInMemory(test, []testZipEntryType{
		{name: "data/save.dat", data: []byte("SAVE-DATA-V1")},
		{name: "config/settings.ini", data: []byte("key=value")},
	})
	assert.NoErrorf(test, MountVirtualFileSystemFromBytes("memory-archive", archiveData, "", ""), "Failed to mount a ZIP archive from bytes.")

	fileData, err := GetFileData("data/save.dat")
	assert.NoErrorf(test, err, "Failed to read a file from an archive mounted from bytes.")
	assert.Equalf(test, "SAVE-DATA-V1", string(fileData), "The file data did not match what was expected.")

	textContent, err := GetTextFileData("config/settings.ini")
	assert.NoErrorf(test, err, "Failed to read text from an archive mounted from bytes.")
	assert.Equalf(test, "key=value", textContent, "The text content did not match what was expected.")

	fileReader, err := GetFileReader("data/save.dat")
	assert.NoErrorf(test, err, "Failed to open a reader for a file in an archive mounted from bytes.")
	streamedData, err := io.ReadAll(fileReader)
	assert.NoErrorf(test, err, "Failed to stream a file from an archive mounted from bytes.")
	assert.NoErrorf(test, fileReader.Close(), "Failed to close the streamed file reader.")
	assert.Equalf(test, "SAVE-DATA-V1", string(streamedData), "The streamed file data did not match what was expected.")

	imageArchiveData, err := os.ReadFile(BASE_DIRECTORY + "valid.zip")
	assert.NoErrorf(test, err, "Failed to read the image test archive.")
	assert.NoErrorf(test, MountVirtualFileSystemFromBytes("memory-images", imageArchiveData, "TAFDRw==", "SampleScrambleKey"), "Failed to mount the image archive from bytes.")
	imagePath := "myFolder-1/myFolder-2/sample3.png"
	assert.NoErrorf(test, LoadImage(imagePath), "Failed to load an image from an archive mounted from bytes.")
	UnloadImage(imagePath)
}

/*
TestMountVirtualFileSystemFromReaderAt is a test which allows you to verify that a ZIP archive behind a caller supplied
io.ReaderAt can be mounted and read, including an image decoded through the image loader.

Example:

	Expected Inputs:
		An in memory ZIP with {"audio/click.raw": "RAW-AUDIO-BYTES"} mounted from a bytes.Reader as "reader-archive",
		then valid.zip opened as an *os.File and mounted from it as "file-reader" with the password "test".
	Expected Outputs:
		GetFileData, GetTextFileData and GetFileReader all return "RAW-AUDIO-BYTES", and the image
		"myFolder-1/sample2.png" decodes without error.
*/
func TestMountVirtualFileSystemFromReaderAt(test *testing.T) {
	defer UnmountVirtualFileSystem()
	archiveData := buildTestZipArchiveInMemory(test, []testZipEntryType{
		{name: "audio/click.raw", data: []byte("RAW-AUDIO-BYTES")},
	})
	assert.NoErrorf(test, MountVirtualFileSystemFromReaderAt("reader-archive", bytes.NewReader(archiveData), int64(len(archiveData)), "", ""), "Failed to mount a ZIP archive from a reader.")

	fileData, err := GetFileData("audio/click.raw")
	assert.NoErrorf(test, err, "Failed to read a file from an archive mounted from a reader.")
	assert.Equalf(test, "RAW-AUDIO-BYTES", string(fileData), "The file data did not match what was expected.")

	textContent, err := GetTextFileData("audio/click.raw")
	assert.NoErrorf(test, err, "Failed to read text from an archive mounted from a reader.")
	assert.Equalf(test, "RAW-AUDIO-BYTES", textContent, "The text content did not match what was expected.")

	fileReader, err := GetFileReader("audio/click.raw")
	assert.NoErrorf(test, err, "Failed to open a reader for a file in an archive mounted from a reader.")
	streamedData, err := io.ReadAll(fileReader)
	assert.NoErrorf(test, err, "Failed to stream a file from an archive mounted from a reader.")
	assert.NoErrorf(test, fileReader.Close(), "Failed to close the streamed file reader.")
	assert.Equalf(test, "RAW-AUDIO-BYTES", string(streamedData), "The streamed file data did not match what was expected.")

	archiveFile, err := os.Open(BASE_DIRECTORY + "valid.zip")
	assert.NoErrorf(test, err, "Failed to open the image test archive.")
	defer archiveFile.Close()
	fileInfo, err := archiveFile.Stat()
	assert.NoErrorf(test, err, "Failed to stat the image test archive.")
	assert.NoErrorf(test, MountVirtualFileSystemFromReaderAt("file-reader", archiveFile, fileInfo.Size(), "test", ""), "Failed to mount the image archive from an *os.File.")
	_, err = getImageFromFileSystem("myFolder-1/sample2.png")
	assert.NoErrorf(test, err, "Failed to decode an image from an archive mounted from a reader.")
}

/*
TestMountVirtualFileSystemFromReaderAtReadsOnlyWhatIsNeeded is a test which allows you to verify, with a reader that
records every byte range read, that a ZIP mount from memory reads only the archive's index and that reading one file
reads only that file's own record.

Example:

	Expected Inputs:
		An in memory stored ZIP of three 64 KiB incompressible entries "first.bin", "second.bin" and "third.bin",
		about 192 KiB in all, mounted from a recordingReaderAt.
	Expected Outputs:
		Every read during mounting starts at or after the earlier of the central directory offset and the last
		1 KiB of the archive, where the end of central directory signature is searched for. Reading "second.bin"
		returns its exact bytes and every read lies inside the record of "second.bin". No read at any stage covers
		half of the archive or more.
*/
func TestMountVirtualFileSystemFromReaderAtReadsOnlyWhatIsNeeded(test *testing.T) {
	defer UnmountVirtualFileSystem()
	entrySize := 64 * 1024
	secondEntryData := getTestRandomBytes(2, entrySize)
	archiveData := buildTestZipArchiveInMemory(test, []testZipEntryType{
		{name: "first.bin", data: getTestRandomBytes(1, entrySize)},
		{name: "second.bin", data: secondEntryData},
		{name: "third.bin", data: getTestRandomBytes(3, entrySize)},
	})
	archiveSize := int64(len(archiveData))
	regions, centralDirectoryOffset := getTestZipRecordRegions(test, archiveData)
	recorder := newRecordingReaderAt(archiveData)

	assert.NoErrorf(test, MountVirtualFileSystemFromReaderAt("tracked", recorder, archiveSize, "", ""), "Failed to mount the tracked archive.")
	mountRanges := recorder.takeRanges()
	assert.NotEmptyf(test, mountRanges, "Expected mounting to read the archive index.")
	indexStart := centralDirectoryOffset
	if archiveSize-1024 < indexStart {
		indexStart = archiveSize - 1024
	}
	for _, currentRange := range mountRanges {
		assert.GreaterOrEqualf(test, currentRange.offset, indexStart, "Mounting read %d bytes at offset %d, before the archive index that starts at %d.", currentRange.length, currentRange.offset, indexStart)
	}
	assertNoReadCoversWholeArchive(test, mountRanges, archiveSize, "Mounting")

	fileData, err := GetFileData("second.bin")
	assert.NoErrorf(test, err, "Failed to read the tracked file.")
	assert.Truef(test, bytes.Equal(secondEntryData, fileData), "The tracked file data did not match what was written.")
	readRanges := recorder.takeRanges()
	assert.NotEmptyf(test, readRanges, "Expected reading a file to read its data.")
	secondRegion := regions["second.bin"]
	for _, currentRange := range readRanges {
		assert.GreaterOrEqualf(test, currentRange.offset, secondRegion.offset, "Reading second.bin read %d bytes at offset %d, before its record.", currentRange.length, currentRange.offset)
		assert.LessOrEqualf(test, currentRange.offset+currentRange.length, secondRegion.offset+secondRegion.length, "Reading second.bin read %d bytes at offset %d, past the end of its record.", currentRange.length, currentRange.offset)
	}
	assertNoReadCoversWholeArchive(test, readRanges, archiveSize, "Reading second.bin")
	assert.LessOrEqualf(test, getTotalBytesRead(readRanges), secondRegion.length, "Reading second.bin read more bytes than its record holds.")
}

/*
TestMountVirtualFileSystemFromReaderAtRarReadsOnlyWhatIsNeeded is a test which allows you to verify that a RAR archive
mounted from memory keeps random access, so mounting reads only block headers and reading one file reads about that
file's bytes rather than the archive from the start.

Example:

	Expected Inputs:
		The bytes of valid.rar, about 237 KiB holding three PNG files with encrypted headers, mounted from a
		recordingReaderAt with the password "test", then "sample.png", whose uncompressed size is 50284 bytes, read.
	Expected Outputs:
		Mounting reads less than half of the archive. Reading "sample.png" returns 50284 bytes and reads no more than
		50284 bytes plus two 4 KiB decoder buffers. No single read covers half of the archive or more.
*/
func TestMountVirtualFileSystemFromReaderAtRarReadsOnlyWhatIsNeeded(test *testing.T) {
	defer UnmountVirtualFileSystem()
	archiveData, err := os.ReadFile(BASE_DIRECTORY + "valid.rar")
	assert.NoErrorf(test, err, "Failed to read the RAR test archive.")
	archiveSize := int64(len(archiveData))
	recorder := newRecordingReaderAt(archiveData)

	assert.NoErrorf(test, MountVirtualFileSystemFromReaderAt("tracked-rar", recorder, archiveSize, "test", ""), "Failed to mount the RAR archive from memory.")
	mountRanges := recorder.takeRanges()
	assert.Lessf(test, getTotalBytesRead(mountRanges), archiveSize/2, "Mounting the RAR archive read too much of it.")
	assertNoReadCoversWholeArchive(test, mountRanges, archiveSize, "Mounting the RAR archive")

	expectedFileSize := int64(50284)
	fileData, err := GetFileData("sample.png")
	assert.NoErrorf(test, err, "Failed to read a file from the RAR archive mounted from memory.")
	assert.Equalf(test, expectedFileSize, int64(len(fileData)), "The RAR entry read from memory had the wrong size.")
	readRanges := recorder.takeRanges()
	assert.LessOrEqualf(test, getTotalBytesRead(readRanges), expectedFileSize+2*4096, "Reading one RAR entry read far more than that entry.")
	assertNoReadCoversWholeArchive(test, readRanges, archiveSize, "Reading sample.png from the RAR archive")
}

/*
TestMountVirtualFileSystemFromBytesPassword is a test which allows you to verify that a password protected ZIP mounted
from memory handles a correct password, a wrong password and a scrambled password exactly like a path mount.

Example:

	Expected Inputs:
		The bytes of the AES encrypted valid.zip, whose password is "test", mounted with "test" in plaintext, with
		"wrong" in plaintext, and with the scrambled password "TAFDRw==" and the scramble key "SampleScrambleKey".
	Expected Outputs:
		The correct and scrambled passwords read "sample.png" without error. The wrong password still mounts, as a
		path mount does, but reading "sample.png" returns an error.
*/
func TestMountVirtualFileSystemFromBytesPassword(test *testing.T) {
	defer UnmountVirtualFileSystem()
	archiveData, err := os.ReadFile(BASE_DIRECTORY + "valid.zip")
	assert.NoErrorf(test, err, "Failed to read the encrypted test archive.")

	assert.NoErrorf(test, MountVirtualFileSystemFromBytes("encrypted", archiveData, "test", ""), "Failed to mount the encrypted archive with the right password.")
	_, err = GetFileData("sample.png")
	assert.NoErrorf(test, err, "Failed to read an encrypted entry with the right password.")
	UnmountVirtualFileSystem()

	assert.NoErrorf(test, MountVirtualFileSystemFromBytes("encrypted", archiveData, "wrong", ""), "Expected a wrong ZIP password to still allow mounting.")
	_, err = GetFileData("sample.png")
	assert.Errorf(test, err, "Expected reading an encrypted entry with the wrong password to fail.")
	UnmountVirtualFileSystem()

	assert.NoErrorf(test, MountVirtualFileSystemFromBytes("encrypted", archiveData, "TAFDRw==", "SampleScrambleKey"), "Failed to mount the encrypted archive with a scrambled password.")
	_, err = GetFileData("sample.png")
	assert.NoErrorf(test, err, "Failed to read an encrypted entry with a scrambled password.")
}

/*
TestMountVirtualFileSystemFromBytesRar is a test which allows you to verify that RAR archives can be mounted from
memory, that a wrong password fails exactly as a path mount does, and that the solid archive limit and its error are
unchanged.

Example:

	Expected Inputs:
		The bytes of valid.rar mounted with the scrambled password "TAFDRw==" and the scramble key
		"SampleScrambleKey", then with a bad scramble key. The solid archive solid.rar, holding "first.txt" and
		"second.txt", mounted both from memory and from its path.
	Expected Outputs:
		The image "myFolder-1/myFolder-2/sample3.png" loads from the RAR in memory. The bad key fails the mount and
		mounts nothing. For the solid archive "first.txt" reads in both cases, while "second.txt" fails in both cases
		with rardecode.ErrSolidOpen.
*/
func TestMountVirtualFileSystemFromBytesRar(test *testing.T) {
	defer UnmountVirtualFileSystem()
	archiveData, err := os.ReadFile(BASE_DIRECTORY + "valid.rar")
	assert.NoErrorf(test, err, "Failed to read the RAR test archive.")

	assert.NoErrorf(test, MountVirtualFileSystemFromBytes("memory-rar", archiveData, "TAFDRw==", "SampleScrambleKey"), "Failed to mount a RAR archive from bytes.")
	_, err = getImageFromFileSystem("myFolder-1/myFolder-2/sample3.png")
	assert.NoErrorf(test, err, "Failed to decode an image from a RAR archive mounted from bytes.")
	UnmountVirtualFileSystem()

	err = MountVirtualFileSystemFromBytes("memory-rar", archiveData, "TAFDRw==", "SampleScrambleKey_BAD")
	assert.Errorf(test, err, "Expected mounting a RAR archive with encrypted headers and a bad password to fail.")
	assert.Equalf(test, 0, getVirtualFileSystemMountCount(), "A failed RAR mount from bytes must not mount anything.")

	solidArchiveData, err := os.ReadFile(BASE_DIRECTORY + "solid.rar")
	assert.NoErrorf(test, err, "Failed to read the solid RAR test archive.")
	mountSolidArchive := []func() error{
		func() error { return MountVirtualFileSystemFromBytes("memory-solid", solidArchiveData, "", "") },
		func() error { return MountVirtualFileSystem(BASE_DIRECTORY+"solid.rar", "", "") },
	}
	for index, mountSolid := range mountSolidArchive {
		context := fmt.Sprintf("Solid archive mount %d", index)
		assert.NoErrorf(test, mountSolid(), "%s: failed to mount the solid archive.", context)
		_, err = GetFileData("first.txt")
		assert.NoErrorf(test, err, "%s: expected the first entry of a solid archive to be readable.", context)
		_, err = GetFileData("second.txt")
		assert.Truef(test, err != nil && strings.Contains(err.Error(), rardecode.ErrSolidOpen.Error()), "%s: expected the solid archive error, got %v.", context, err)
		UnmountVirtualFileSystem()
	}
}

/*
TestMountVirtualFileSystemFromBytesShadowing is a test which allows you to verify that a mount from memory and a path
mount shadow each other in both directions, with the most recent mount winning.

Example:

	Expected Inputs:
		A path archive with {"shared.txt": "PATH", "only_path.txt": "p"} and an in memory archive named "memory" with
		{"shared.txt": "MEMORY", "only_memory.txt": "m"}, mounted path first and then memory first.
	Expected Outputs:
		With the memory mount on top "shared.txt" reads "MEMORY", and after unmounting "memory" it reads "PATH". With
		the path mount on top "shared.txt" reads "PATH", and after unmounting the path it reads "MEMORY". Files unique
		to either mount stay readable while their mount is in place.
*/
func TestMountVirtualFileSystemFromBytesShadowing(test *testing.T) {
	defer UnmountVirtualFileSystem()
	archivePath := writeTestZipArchive(test, test.TempDir(), "path.zip", map[string]string{
		"shared.txt":    "PATH",
		"only_path.txt": "p",
	})
	archiveData := buildTestZipArchiveInMemory(test, []testZipEntryType{
		{name: "shared.txt", data: []byte("MEMORY")},
		{name: "only_memory.txt", data: []byte("m")},
	})
	assertText := func(filePath string, expectedText string, context string) {
		test.Helper()
		textContent, err := GetTextFileData(filePath)
		assert.NoErrorf(test, err, "%s: failed to read '%s'.", context, filePath)
		assert.Equalf(test, expectedText, textContent, "%s: '%s' had the wrong content.", context, filePath)
	}

	assert.NoErrorf(test, MountVirtualFileSystem(archivePath, "", ""), "Failed to mount the path archive.")
	assert.NoErrorf(test, MountVirtualFileSystemFromBytes("memory", archiveData, "", ""), "Failed to mount the memory archive.")
	assertText("shared.txt", "MEMORY", "Memory mounted over path")
	assertText("only_path.txt", "p", "Memory mounted over path")
	assertText("only_memory.txt", "m", "Memory mounted over path")
	assert.NoErrorf(test, UnmountVirtualFileSystemArchive("memory"), "Failed to unmount the memory archive by name.")
	assertText("shared.txt", "PATH", "Memory archive unmounted")
	UnmountVirtualFileSystem()

	assert.NoErrorf(test, MountVirtualFileSystemFromBytes("memory", archiveData, "", ""), "Failed to mount the memory archive.")
	assert.NoErrorf(test, MountVirtualFileSystem(archivePath, "", ""), "Failed to mount the path archive.")
	assertText("shared.txt", "PATH", "Path mounted over memory")
	assertText("only_memory.txt", "m", "Path mounted over memory")
	assert.NoErrorf(test, UnmountVirtualFileSystemArchive(archivePath), "Failed to unmount the path archive.")
	assertText("shared.txt", "MEMORY", "Path archive unmounted")
}

/*
TestMountVirtualFileSystemFromReaderAtUnmount is a test which allows you to verify that a mount from memory is removed
by its name without disturbing other mounts, that mounting a name twice adds two layers which one unmount removes, and
that no unmount ever closes the caller's reader.

Example:

	Expected Inputs:
		Recording readers over in memory archives "alpha" with {"alpha.txt": "a"} and "beta" with {"beta.txt": "b"},
		plus "alpha" mounted a second time with {"alpha.txt": "a2"}.
	Expected Outputs:
		"alpha.txt" reads "a2" while both "alpha" layers are mounted. UnmountVirtualFileSystemArchive("alpha") removes
		both layers so "alpha.txt" errors while "beta.txt" still reads "b". UnmountVirtualFileSystem() removes the
		rest. Close is never called on any reader, and each reader is still usable afterwards.
*/
func TestMountVirtualFileSystemFromReaderAtUnmount(test *testing.T) {
	defer UnmountVirtualFileSystem()
	alphaData := buildTestZipArchiveInMemory(test, []testZipEntryType{{name: "alpha.txt", data: []byte("a")}})
	alphaSecondData := buildTestZipArchiveInMemory(test, []testZipEntryType{{name: "alpha.txt", data: []byte("a2")}})
	betaData := buildTestZipArchiveInMemory(test, []testZipEntryType{{name: "beta.txt", data: []byte("b")}})
	alphaReader := newRecordingReaderAt(alphaData)
	alphaSecondReader := newRecordingReaderAt(alphaSecondData)
	betaReader := newRecordingReaderAt(betaData)

	assert.NoErrorf(test, MountVirtualFileSystemFromReaderAt("alpha", alphaReader, int64(len(alphaData)), "", ""), "Failed to mount alpha.")
	assert.NoErrorf(test, MountVirtualFileSystemFromReaderAt("beta", betaReader, int64(len(betaData)), "", ""), "Failed to mount beta.")
	assert.NoErrorf(test, MountVirtualFileSystemFromReaderAt("alpha", alphaSecondReader, int64(len(alphaSecondData)), "", ""), "Failed to mount alpha a second time.")
	assert.Equalf(test, 3, getVirtualFileSystemMountCount(), "Mounting a name twice should add a second layer.")
	textContent, err := GetTextFileData("alpha.txt")
	assert.NoErrorf(test, err, "Failed to read alpha.txt.")
	assert.Equalf(test, "a2", textContent, "The second alpha layer should shadow the first.")

	assert.NoErrorf(test, UnmountVirtualFileSystemArchive("alpha"), "Failed to unmount alpha by name.")
	assert.Equalf(test, 1, getVirtualFileSystemMountCount(), "Unmounting a name should remove every layer that uses it.")
	_, err = GetFileData("alpha.txt")
	assert.Errorf(test, err, "Expected alpha.txt to be gone after unmounting alpha.")
	textContent, err = GetTextFileData("beta.txt")
	assert.NoErrorf(test, err, "Expected beta.txt to survive unmounting alpha.")
	assert.Equalf(test, "b", textContent, "beta.txt had the wrong content after unmounting alpha.")
	assert.Errorf(test, UnmountVirtualFileSystemArchive("alpha"), "Expected unmounting an already removed name to fail.")

	UnmountVirtualFileSystem()
	assert.Equalf(test, 0, getVirtualFileSystemMountCount(), "UnmountVirtualFileSystem should remove every mount.")
	for readerName, reader := range map[string]*recordingReaderAt{"alpha": alphaReader, "alpha second": alphaSecondReader, "beta": betaReader} {
		assert.Equalf(test, 0, reader.getCloseCount(), "The %s reader must never be closed by the virtual file system.", readerName)
		_, err = reader.ReadAt(make([]byte, 4), 0)
		assert.NoErrorf(test, err, "The %s reader should remain usable after unmounting.", readerName)
	}
}

/*
TestMountVirtualFileSystemFromReaderAtInvalidInput is a test which allows you to verify that every invalid mount from
memory is rejected with an error and mounts nothing.

Example:

	Expected Inputs:
		Garbage bytes "this is not an archive", an empty slice, a nil slice, a nil reader, a size of zero, a negative
		size, and an empty name with an otherwise valid archive.
	Expected Outputs:
		Every call returns an error and the mount count stays zero. The garbage error names the mount and describes
		both the ZIP and the RAR failure.
*/
func TestMountVirtualFileSystemFromReaderAtInvalidInput(test *testing.T) {
	defer UnmountVirtualFileSystem()
	UnmountVirtualFileSystem()
	validArchiveData := buildTestZipArchiveInMemory(test, []testZipEntryType{{name: "file.txt", data: []byte("x")}})
	garbageData := []byte("this is not an archive")

	err := MountVirtualFileSystemFromBytes("garbage", garbageData, "", "")
	assert.Errorf(test, err, "Expected garbage bytes to fail to mount.")
	if err != nil {
		assert.Containsf(test, err.Error(), "'garbage'", "The error should name the mount.")
		assert.Containsf(test, err.Error(), "ZIP (", "The error should describe the ZIP failure.")
		assert.Containsf(test, err.Error(), "RAR (", "The error should describe the RAR failure.")
	}
	invalidMounts := map[string]func() error{
		"empty slice": func() error { return MountVirtualFileSystemFromBytes("empty", []byte{}, "", "") },
		"nil slice":   func() error { return MountVirtualFileSystemFromBytes("nil-slice", nil, "", "") },
		"nil reader":  func() error { return MountVirtualFileSystemFromReaderAt("nil-reader", nil, 100, "", "") },
		"zero size": func() error {
			return MountVirtualFileSystemFromReaderAt("zero", bytes.NewReader(validArchiveData), 0, "", "")
		},
		"negative size": func() error {
			return MountVirtualFileSystemFromReaderAt("negative", bytes.NewReader(validArchiveData), -1, "", "")
		},
		"empty name bytes": func() error { return MountVirtualFileSystemFromBytes("", validArchiveData, "", "") },
		"empty name reader": func() error {
			return MountVirtualFileSystemFromReaderAt("", bytes.NewReader(validArchiveData), int64(len(validArchiveData)), "", "")
		},
	}
	for caseName, invalidMount := range invalidMounts {
		assert.Errorf(test, invalidMount(), "Expected the %s case to fail to mount.", caseName)
	}
	assert.Equalf(test, 0, getVirtualFileSystemMountCount(), "No invalid mount from memory may mount anything.")
}

/*
TestFileExistsWithNothingMounted is a test which allows you to verify that with nothing mounted FileExists reports only
regular files on the local file system, in agreement with GetFileData.

Example:

	Expected Inputs:
		The local file "text_file.txt", a missing file, the test data directory, and an empty path.
	Expected Outputs:
		Only "text_file.txt" is true. Each answer matches whether GetFileData succeeds.
*/
func TestFileExistsWithNothingMounted(test *testing.T) {
	UnmountVirtualFileSystem()
	assertFileExistsMatchesGetFileData(test, BASE_DIRECTORY+"text_file.txt", true, "Nothing mounted")
	assertFileExistsMatchesGetFileData(test, BASE_DIRECTORY+"missing.txt", false, "Nothing mounted")
	assertFileExistsMatchesGetFileData(test, BASE_DIRECTORY, false, "Nothing mounted")
	assertFileExistsMatchesGetFileData(test, strings.TrimSuffix(BASE_DIRECTORY, "/"), false, "Nothing mounted")
	assertFileExistsMatchesGetFileData(test, "", false, "Nothing mounted")
}

/*
TestFileExistsWithArchiveMountedFromPath is a test which allows you to verify FileExists against a ZIP archive mounted
from a path, including alternate spellings of the same path, a file that exists only on the local disk, and a folder.

Example:

	Expected Inputs:
		A path archive with the directory entry "folder/" and the file "folder/file.txt", then the RAR valid.rar.
	Expected Outputs:
		"folder/file.txt", "./folder/file.txt", "folder\\file.txt" and "/folder/file.txt" are true. The local only
		file "text_file.txt", "folder", "folder/", a missing file and an empty path are false. For the RAR,
		"myFolder-1/sample2.png" is true and "myFolder-1" is false. Each answer matches GetFileData.
*/
func TestFileExistsWithArchiveMountedFromPath(test *testing.T) {
	defer UnmountVirtualFileSystem()
	archivePath := filepath.Join(test.TempDir(), "folders.zip")
	archiveData := buildTestZipArchiveInMemory(test, []testZipEntryType{
		{name: "folder/"},
		{name: "folder/file.txt", data: []byte("inside")},
	})
	assert.NoErrorf(test, os.WriteFile(archivePath, archiveData, 0644), "Failed to write the test archive.")
	assert.NoErrorf(test, MountVirtualFileSystem(archivePath, "", ""), "Failed to mount the test archive.")

	for _, filePath := range []string{"folder/file.txt", "./folder/file.txt", "folder\\file.txt", "/folder/file.txt", "folder//file.txt"} {
		assertFileExistsMatchesGetFileData(test, filePath, true, "Path mount")
	}
	for _, filePath := range []string{BASE_DIRECTORY + "text_file.txt", "folder", "folder/", "folder/missing.txt", ""} {
		assertFileExistsMatchesGetFileData(test, filePath, false, "Path mount")
	}
	UnmountVirtualFileSystem()

	assert.NoErrorf(test, MountVirtualFileSystem(BASE_DIRECTORY+"valid.rar", "test", ""), "Failed to mount the RAR archive.")
	assertFileExistsMatchesGetFileData(test, "myFolder-1/sample2.png", true, "RAR path mount")
	assertFileExistsMatchesGetFileData(test, "myFolder-1", false, "RAR path mount")
	assertFileExistsMatchesGetFileData(test, "myFolder-1/myFolder-2", false, "RAR path mount")
}

/*
TestFileExistsWithEmbeddedFileSystem is a test which allows you to verify FileExists against a mounted embedded file
system.

Example:

	Expected Inputs:
		The embedded test folder holding "embedded_file.txt" and "embedded_folder/nested_file.txt", mounted with
		MountEmbeddedFileSystem.
	Expected Outputs:
		Both files are true, including the spelling with a leading "./". The local only file "text_file.txt", the
		embedded folders, and a missing file are false. Each answer matches GetFileData.
*/
func TestFileExistsWithEmbeddedFileSystem(test *testing.T) {
	defer UnmountVirtualFileSystem()
	assert.NoErrorf(test, MountEmbeddedFileSystem(embeddedTestFileSystem), "Failed to mount the embedded test file system.")

	for _, filePath := range []string{
		EMBEDDED_TEST_DIRECTORY + "/embedded_file.txt",
		"./" + EMBEDDED_TEST_DIRECTORY + "/embedded_file.txt",
		EMBEDDED_TEST_DIRECTORY + "/embedded_folder/nested_file.txt",
	} {
		assertFileExistsMatchesGetFileData(test, filePath, true, "Embedded mount")
	}
	for _, filePath := range []string{
		BASE_DIRECTORY + "text_file.txt",
		EMBEDDED_TEST_DIRECTORY,
		EMBEDDED_TEST_DIRECTORY + "/embedded_folder",
		EMBEDDED_TEST_DIRECTORY + "/missing.txt",
		"",
	} {
		assertFileExistsMatchesGetFileData(test, filePath, false, "Embedded mount")
	}
}

/*
TestFileExistsWithArchiveMountedFromMemory is a test which allows you to verify FileExists against ZIP and RAR archives
mounted from memory, and that answering reads nothing from the caller's reader.

Example:

	Expected Inputs:
		An in memory ZIP with the directory entry "folder/" and the file "folder/file.txt" mounted from a
		recordingReaderAt, then valid.rar mounted from a second recordingReaderAt with the password "test".
	Expected Outputs:
		"folder/file.txt" and "myFolder-1/sample2.png" are true, while the local only file "text_file.txt", the
		folders, and a missing file are false. Each answer matches GetFileData. FileExists records zero reads on
		either reader.
*/
func TestFileExistsWithArchiveMountedFromMemory(test *testing.T) {
	defer UnmountVirtualFileSystem()
	archiveData := buildTestZipArchiveInMemory(test, []testZipEntryType{
		{name: "folder/"},
		{name: "folder/file.txt", data: []byte("inside")},
	})
	zipRecorder := newRecordingReaderAt(archiveData)
	assert.NoErrorf(test, MountVirtualFileSystemFromReaderAt("memory-zip", zipRecorder, int64(len(archiveData)), "", ""), "Failed to mount the ZIP from memory.")
	rarData, err := os.ReadFile(BASE_DIRECTORY + "valid.rar")
	assert.NoErrorf(test, err, "Failed to read the RAR test archive.")
	rarRecorder := newRecordingReaderAt(rarData)
	assert.NoErrorf(test, MountVirtualFileSystemFromReaderAt("memory-rar", rarRecorder, int64(len(rarData)), "test", ""), "Failed to mount the RAR from memory.")
	zipRecorder.takeRanges()
	rarRecorder.takeRanges()

	expectedAnswers := map[string]bool{
		"folder/file.txt":                   true,
		"./folder/file.txt":                 true,
		"folder\\file.txt":                  true,
		"myFolder-1/sample2.png":            true,
		"myFolder-1/myFolder-2/sample3.png": true,
		BASE_DIRECTORY + "text_file.txt":    false,
		"folder":                            false,
		"folder/":                           false,
		"myFolder-1":                        false,
		"myFolder-1/myFolder-2/missing.png": false,
		"":                                  false,
	}
	for filePath, expectedAnswer := range expectedAnswers {
		assert.Equalf(test, expectedAnswer, FileExists(filePath), "Memory mount: FileExists(%q) returned the wrong answer.", filePath)
	}
	assert.Emptyf(test, zipRecorder.takeRanges(), "FileExists must not read from a ZIP mounted from memory.")
	assert.Emptyf(test, rarRecorder.takeRanges(), "FileExists must not read from a RAR mounted from memory.")
	for filePath, expectedAnswer := range expectedAnswers {
		assertFileExistsMatchesGetFileData(test, filePath, expectedAnswer, "Memory mount")
	}
}

/*
TestFileExistsWithSeveralMountsAndUnmount is a test which allows you to verify that FileExists sees a file served by
any of several mounts, follows each unmount, and returns to the local file system rule once everything is unmounted.

Example:

	Expected Inputs:
		A path archive with {"path.txt": "p", "shared.txt": "s"}, an in memory archive "memory" with
		{"memory.txt": "m", "shared.txt": "s"}, and the embedded test folder, all mounted together.
	Expected Outputs:
		"path.txt", "memory.txt", "shared.txt" and the embedded "embedded_file.txt" are true. After unmounting the path
		archive "path.txt" is false and the rest stay true. After unmounting "memory" its files are false and the
		embedded file stays true. After UnmountVirtualFileSystem the local file "text_file.txt" is true again and
		archive only paths are false. Each answer matches GetFileData.
*/
func TestFileExistsWithSeveralMountsAndUnmount(test *testing.T) {
	defer UnmountVirtualFileSystem()
	archivePath := writeTestZipArchive(test, test.TempDir(), "path.zip", map[string]string{"path.txt": "p", "shared.txt": "s"})
	archiveData := buildTestZipArchiveInMemory(test, []testZipEntryType{
		{name: "memory.txt", data: []byte("m")},
		{name: "shared.txt", data: []byte("s")},
	})
	embeddedFile := EMBEDDED_TEST_DIRECTORY + "/embedded_file.txt"
	localFile := BASE_DIRECTORY + "text_file.txt"
	assert.NoErrorf(test, MountVirtualFileSystem(archivePath, "", ""), "Failed to mount the path archive.")
	assert.NoErrorf(test, MountVirtualFileSystemFromBytes("memory", archiveData, "", ""), "Failed to mount the memory archive.")
	assert.NoErrorf(test, MountEmbeddedFileSystem(embeddedTestFileSystem), "Failed to mount the embedded file system.")

	for _, filePath := range []string{"path.txt", "memory.txt", "shared.txt", embeddedFile} {
		assertFileExistsMatchesGetFileData(test, filePath, true, "All mounted")
	}
	assertFileExistsMatchesGetFileData(test, localFile, false, "All mounted")

	assert.NoErrorf(test, UnmountVirtualFileSystemArchive(archivePath), "Failed to unmount the path archive.")
	assertFileExistsMatchesGetFileData(test, "path.txt", false, "Path archive unmounted")
	for _, filePath := range []string{"memory.txt", "shared.txt", embeddedFile} {
		assertFileExistsMatchesGetFileData(test, filePath, true, "Path archive unmounted")
	}

	assert.NoErrorf(test, UnmountVirtualFileSystemArchive("memory"), "Failed to unmount the memory archive.")
	for _, filePath := range []string{"memory.txt", "shared.txt"} {
		assertFileExistsMatchesGetFileData(test, filePath, false, "Memory archive unmounted")
	}
	assertFileExistsMatchesGetFileData(test, embeddedFile, true, "Memory archive unmounted")

	UnmountVirtualFileSystem()
	assertFileExistsMatchesGetFileData(test, localFile, true, "Everything unmounted")
	assertFileExistsMatchesGetFileData(test, "memory.txt", false, "Everything unmounted")
	assertFileExistsMatchesGetFileData(test, embeddedFile, true, "Everything unmounted")
}

/*
TestFileExistsConcurrentWithMountAndUnmount is a test which allows you to verify, under the race detector, that
FileExists can be called while other goroutines mount and unmount archives, and that it never reports a file outside
the set of files some mount could serve.

Example:

	Expected Inputs:
		One goroutine repeatedly mounting and unmounting an in memory archive "churn" with {"churn.txt": "c"}, while
		several goroutines call FileExists on "churn.txt", the always mounted "stable.txt" and "never.txt".
	Expected Outputs:
		No data race is reported, "stable.txt" is always true, "never.txt" is always false, and once the churn stops
		with "churn" unmounted "churn.txt" is false.
*/
func TestFileExistsConcurrentWithMountAndUnmount(test *testing.T) {
	defer UnmountVirtualFileSystem()
	stableData := buildTestZipArchiveInMemory(test, []testZipEntryType{{name: "stable.txt", data: []byte("s")}})
	churnData := buildTestZipArchiveInMemory(test, []testZipEntryType{{name: "churn.txt", data: []byte("c")}})
	assert.NoErrorf(test, MountVirtualFileSystemFromBytes("stable", stableData, "", ""), "Failed to mount the stable archive.")

	var waitGroup sync.WaitGroup
	stopChannel := make(chan struct{})
	errorChannel := make(chan error, 16)
	for checkerIndex := 0; checkerIndex < 4; checkerIndex++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			for {
				select {
				case <-stopChannel:
					return
				default:
				}
				FileExists("churn.txt")
				if !FileExists("stable.txt") {
					errorChannel <- errors.New("stable.txt was reported missing while it was mounted")
					return
				}
				if FileExists("never.txt") {
					errorChannel <- errors.New("never.txt was reported present although no mount serves it")
					return
				}
			}
		}()
	}
	for iteration := 0; iteration < 200; iteration++ {
		assert.NoErrorf(test, MountVirtualFileSystemFromBytes("churn", churnData, "", ""), "Failed to mount the churn archive.")
		assert.NoErrorf(test, UnmountVirtualFileSystemArchive("churn"), "Failed to unmount the churn archive.")
	}
	close(stopChannel)
	waitGroup.Wait()
	close(errorChannel)
	for err := range errorChannel {
		assert.NoErrorf(test, err, "A concurrent FileExists check returned a wrong answer.")
	}
	assertFileExistsMatchesGetFileData(test, "churn.txt", false, "Churn stopped")
	assertFileExistsMatchesGetFileData(test, "stable.txt", true, "Churn stopped")
}
