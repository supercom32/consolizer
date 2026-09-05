package consolizer

import (
	"archive/zip"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"path/filepath"
	"testing"
)

const (
	BASE_DIRECTORY = "./test_data/virtual_file_systems/"
)

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
