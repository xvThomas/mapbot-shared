package testutils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetTestDataFilePath(t *testing.T) {
	t.Run("returns path relative to caller", func(t *testing.T) {
		path := GetTestDataFilePath("sample.txt")

		if !strings.Contains(path, "testdata") {
			t.Errorf("GetTestDataFilePath() = %q, should contain 'testdata'", path)
		}

		if !strings.HasSuffix(path, filepath.Join("testdata", "sample.txt")) {
			t.Errorf("GetTestDataFilePath() = %q, should end with 'testdata/sample.txt'", path)
		}

		if !filepath.IsAbs(path) {
			t.Errorf("GetTestDataFilePath() should return absolute path, got %q", path)
		}
	})

	t.Run("works with subdirectory paths", func(t *testing.T) {
		path := GetTestDataFilePath("subdir/file.json")

		expected := filepath.Join("testdata", "subdir", "file.json")
		if !strings.HasSuffix(path, expected) {
			t.Errorf("GetTestDataFilePath() = %q, should end with %q", path, expected)
		}
	})

	t.Run("returns valid path even if file does not exist", func(t *testing.T) {
		path := GetTestDataFilePath("nonexistent.txt")

		if !strings.Contains(path, "testdata") {
			t.Errorf("GetTestDataFilePath() = %q, should contain 'testdata'", path)
		}
	})

	t.Run("actual file access", func(t *testing.T) {
		testDir := filepath.Join(".", "testdata")
		if err := os.MkdirAll(testDir, 0750); err != nil {
			t.Fatalf("Failed to create testdata directory: %v", err)
		}
		defer func() { _ = os.RemoveAll(testDir) }()

		testFile := filepath.Join(testDir, "test_file.txt")
		testContent := []byte("test content")
		if err := os.WriteFile(testFile, testContent, 0600); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		path := GetTestDataFilePath("test_file.txt")

		// #nosec G304 - Reading test file with dynamic path is intentional in tests
		content, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("GetTestDataFilePath() returned %q, cannot be read: %v", path, err)
		}

		if string(content) != string(testContent) {
			t.Errorf("File content = %q, want %q", string(content), string(testContent))
		}
	})
}
