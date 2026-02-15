package testutils

import (
	"log"
	"path/filepath"
	"runtime"
)

// GetTestDataFilePath returns the absolute path to a test data file
// It looks for a "testdata" directory relative to the calling test file
func GetTestDataFilePath(file string) string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		log.Fatal("failed to get current file path")
	}
	// Get directory of the caller and join with testdata path
	dir := filepath.Dir(filename)
	path := filepath.Join(dir, "testdata", file)
	return path
}
