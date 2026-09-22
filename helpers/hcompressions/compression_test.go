package hcompressions

import (
	"os"
	"path/filepath"
	"testing"
)

func TestZipCompressionAndDecompression(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()
	zipFile := filepath.Join(t.TempDir(), "archive.zip")

	sampleFile := filepath.Join(srcDir, "sample.txt")
	sampleContent := []byte("Compression test content")
	if err := os.WriteFile(sampleFile, sampleContent, 0644); err != nil {
		t.Fatalf("Failed to create sample file: %v", err)
	}

	subDir := filepath.Join(srcDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	subFile := filepath.Join(subDir, "sub.txt")
	if err := os.WriteFile(subFile, []byte("Sub content"), 0644); err != nil {
		t.Fatalf("Failed to create sub file: %v", err)
	}

	err := ZipCompression(srcDir+string(filepath.Separator), zipFile, false)
	if err != nil {
		t.Fatalf("ZipCompression failed: %v", err)
	}

	err = ZipDecompression(zipFile, dstDir)
	if err != nil {
		t.Fatalf("ZipDecompression failed: %v", err)
	}

	extractedSample := filepath.Join(dstDir, "sample.txt")
	content, err := os.ReadFile(extractedSample)
	if err != nil {
		t.Fatalf("Failed to read extracted file: %v", err)
	}
	if string(content) != string(sampleContent) {
		t.Errorf("Extracted content mismatch. Got: %s, Want: %s", string(content), string(sampleContent))
	}

	extractedSub := filepath.Join(dstDir, "subdir", "sub.txt")
	subContent, err := os.ReadFile(extractedSub)
	if err != nil {
		t.Fatalf("Failed to read extracted sub file: %v", err)
	}
	if string(subContent) != "Sub content" {
		t.Errorf("Extracted sub content mismatch. Got: %s, Want: Sub content", string(subContent))
	}
}
