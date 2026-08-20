package consolizer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	BASE_DIRECTORY = "./test_data/virtual_file_systems/"
)

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
TestMountVirtualFileSystem is a test which verifies the mounting and file retrieval from virtual ZIP and RAR file
systems, including error handling for invalid passwords and missing files.

Example:
    Expected Inputs:
        Valid and invalid ZIP/RAR archives and password/key combinations.
    Expected Outputs:
        Successful mounting returns nil, and subsequent file reads return correct data or appropriate errors.
*/
func TestMountVirtualFileSystem(test *testing.T) {
	var expectedResult error
	scrambledPassword := "TAFDRw=="
	scrambleKey := "SampleScrambleKey"

	// Verify valid statements.
	obtainedResult := mountVirtualFileSystem(BASE_DIRECTORY+"valid.zip", scrambledPassword, scrambleKey)
	assert.Equalf(test, obtainedResult, expectedResult, "Failed to open a valid ZIP filesystem!")
	_, err := getImageFromFileSystem("myFolder-1/myFolder-2/sample3.png")
	assert.NoErrorf(test, err, "Failed to obtain image data from ZIP file system.")
	obtainedResult = mountVirtualFileSystem(BASE_DIRECTORY+"valid.rar", scrambledPassword, scrambleKey)
	_, err = getImageFromFileSystem("myFolder-1/myFolder-2/sample3.png")
	assert.NoErrorf(test, err, "Failed to obtain image data from RAR file system.")

	// Verify invalid archive passwords generate errors.
	scrambleKey = "SampleScrambleKey_BAD"
	obtainedResult = mountVirtualFileSystem(BASE_DIRECTORY+"valid.zip", scrambledPassword, scrambleKey)
	assert.Equalf(test, obtainedResult, expectedResult, "Failed to open a valid ZIP filesystem!")
	_, err = getImageFromFileSystem("myFolder-1/myFolder-2/sample3.png")
	assert.Errorf(test, err, "Expected an error retrieving a file from ZIP file system with a bad password.")
	obtainedResult = mountVirtualFileSystem(BASE_DIRECTORY+"valid.rar", scrambledPassword, scrambleKey)
	_, err = getImageFromFileSystem("myFolder-1/myFolder-2/sample3.png")
	assert.Errorf(test, err, "Expected an error retrieving a file from RAR file system with a bad password.")

	// Verify invalid file requests generate errors.
	scrambleKey = "SampleScrambleKey"
	obtainedResult = mountVirtualFileSystem(BASE_DIRECTORY+"valid.zip", scrambledPassword, scrambleKey)
	assert.Equalf(test, obtainedResult, expectedResult, "Failed to open a valid ZIP filesystem!")
	_, err = getImageFromFileSystem("myFolder-1/myFolder-2/sample3.png_BAD")
	assert.Errorf(test, err, "Expected an error retrieving a file from ZIP file system that does not exist.")
	obtainedResult = mountVirtualFileSystem(BASE_DIRECTORY+"valid.rar", scrambledPassword, scrambleKey)
	_, err = getImageFromFileSystem("myFolder-1/myFolder-2/sample3.png_BAD")
	assert.Errorf(test, err, "Expected an error retrieving a file from RAR file system that does not exist.")

	// Verify invalid archives.
	assert.Equalf(test, obtainedResult, expectedResult, "Failed to open a valid ZIP filesystem!")
	obtainedResult = mountVirtualFileSystem(BASE_DIRECTORY+"invalid.zip", scrambledPassword, scrambleKey)
	assert.NotNil(test, obtainedResult, "Opening an invalid ZIP file was expected to fail when it didn't!")
	UnmountVirtualFileSystem()
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
