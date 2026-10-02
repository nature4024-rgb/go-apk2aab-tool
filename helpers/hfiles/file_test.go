package hfiles

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileExists(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.txt")

	if FileExists(tempFile) {
		t.Errorf("Expected FileExists(%q) to be false before creation", tempFile)
	}

	err := os.WriteFile(tempFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if !FileExists(tempFile) {
		t.Errorf("Expected FileExists(%q) to be true after creation", tempFile)
	}

	if FileExists(tempDir) {
		t.Errorf("Expected FileExists(%q) on directory to be false", tempDir)
	}
}

func TestFolderExists(t *testing.T) {
	tempDir := t.TempDir()
	nonExistentDir := filepath.Join(tempDir, "nonexistent")

	if !FolderExists(tempDir) {
		t.Errorf("Expected FolderExists(%q) to be true", tempDir)
	}

	if FolderExists(nonExistentDir) {
		t.Errorf("Expected FolderExists(%q) to be false", nonExistentDir)
	}
}
