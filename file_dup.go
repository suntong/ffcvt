// Copy or hardlink a file in Go

package main

import (
	"io"
	"os"
)

func mainSample() {
	copyFile("test.txt", "test_copy.txt")
	copyFile("test_copy.txt", "test_link.txt")
}

// hardlink a file
func linkFile(src, dst string) error {
	return os.Link(src, dst)
}

// Copy a file
func copyFile(src, dst string) {
	srcFile, err := os.Open(src)
	check(err)
	defer srcFile.Close()

	destFile, err := os.Create(dst) // creates if file doesn't exist
	check(err)
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile) // check first var for number of bytes copied
	check(err)

	err = destFile.Sync()
	check(err)
}

func check(err error) {
	checkError(err)
}
