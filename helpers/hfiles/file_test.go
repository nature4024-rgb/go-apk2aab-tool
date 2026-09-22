package hfiles

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileAndFolderExists(t *testing.T) {
	tempDir := t.TempDir()

	testFile := filepath.Join(tempDir, "test.txt")
	testFolder := filepath.Join(tempDir, "subfolder")

	if FileExists(testFile) {
		t.Errorf("FileExists returned true for non-existent file")
	}
	if FolderExists(testFolder) {
		t.Errorf("FolderExists returned true for non-existent folder")
	}

	if err := os.WriteFile(testFile, []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	if err := os.Mkdir(testFolder, 0755); err != nil {
		t.Fatalf("failed to create test folder: %v", err)
	}

	if !FileExists(testFile) {
		t.Errorf("FileExists returned false for existing file")
	}

	if FolderExists(testFile) {
		t.Errorf("FolderExists returned true for a file path")
	}

	if !FolderExists(testFolder) {
		t.Errorf("FolderExists returned false for existing folder")
	}

	if FileExists(testFolder) {
		t.Errorf("FileExists returned true for a folder path")
	}
}
