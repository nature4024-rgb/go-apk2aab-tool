package main

import (
	"archive/zip"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestAppConfigDiscovery(t *testing.T) {
	setAppConfig("34.0.0", "21", "34")
	if oAppConfig.errorMessage != "" {
		t.Fatalf("setAppConfig failed: %s", oAppConfig.errorMessage)
	}

	if oAppConfig.javaBinFilePath == "" {
		t.Errorf("javaBinFilePath should not be empty")
	}
	if oAppConfig.aapt2BinFilePath == "" {
		t.Errorf("aapt2BinFilePath should not be empty")
	}
	if oAppConfig.androidJarFilePath == "" {
		t.Errorf("androidJarFilePath should not be empty")
	}
	if oAppConfig.apktoolJarFilePath == "" {
		t.Errorf("apktoolJarFilePath should not be empty")
	}
	if oAppConfig.bundletoolJarFilePath == "" {
		t.Errorf("bundletoolJarFilePath should not be empty")
	}
}

func TestAPKToAABPipeline(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Prepare minimal Android project files
	resDir := filepath.Join(tempDir, "res", "values")
	if err := os.MkdirAll(resDir, 0755); err != nil {
		t.Fatalf("Failed to create res dir: %v", err)
	}

	manifestContent := `<manifest xmlns:android="http://schemas.android.com/apk/res/android" package="com.example.pipeline.test">
    <application android:label="Pipeline Test">
        <activity android:name=".MainActivity" android:exported="true">
            <intent-filter>
                <action android:name="android.intent.action.MAIN" />
                <category android:name="android.intent.category.LAUNCHER" />
            </intent-filter>
        </activity>
    </application>
</manifest>`

	manifestPath := filepath.Join(tempDir, "AndroidManifest.xml")
	if err := os.WriteFile(manifestPath, []byte(manifestContent), 0644); err != nil {
		t.Fatalf("Failed to write manifest: %v", err)
	}

	stringsXml := `<resources><string name="app_name">Pipeline Test</string></resources>`
	if err := os.WriteFile(filepath.Join(resDir, "strings.xml"), []byte(stringsXml), 0644); err != nil {
		t.Fatalf("Failed to write strings.xml: %v", err)
	}

	setAppConfig("34.0.0", "21", "34")
	if oAppConfig.errorMessage != "" {
		t.Fatalf("Environment config error: %s", oAppConfig.errorMessage)
	}

	// 2. Compile resources and link base APK
	compiledZip := filepath.Join(tempDir, "compiled.zip")
	if err := runCommand(oAppConfig.aapt2BinFilePath, "compile", "--dir", filepath.Join(tempDir, "res"), "-o", compiledZip); err != nil {
		t.Fatalf("aapt2 compile failed: %v", err)
	}

	sampleAPK := filepath.Join(tempDir, "sample_test.apk")
	if err := runCommand(oAppConfig.aapt2BinFilePath, "link",
		"-o", sampleAPK,
		"-I", oAppConfig.androidJarFilePath,
		"--manifest", manifestPath,
		"-R", compiledZip,
		"--auto-add-overlay",
	); err != nil {
		t.Fatalf("aapt2 link failed: %v", err)
	}

	// 3. Inject DEX files, Assets, Native Libs, and Unknown files into test APK
	apkContentsDir := filepath.Join(tempDir, "apk_contents")
	assetsDir := filepath.Join(apkContentsDir, "assets")
	libDir := filepath.Join(apkContentsDir, "lib", "arm64-v8a")
	unknownDir := filepath.Join(apkContentsDir, "unknown")

	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		t.Fatalf("Failed to create assets dir: %v", err)
	}
	if err := os.MkdirAll(libDir, 0755); err != nil {
		t.Fatalf("Failed to create lib dir: %v", err)
	}
	if err := os.MkdirAll(unknownDir, 0755); err != nil {
		t.Fatalf("Failed to create unknown dir: %v", err)
	}

	if err := os.WriteFile(filepath.Join(apkContentsDir, "classes.dex"), []byte("dex1"), 0644); err != nil {
		t.Fatalf("Failed to write classes.dex: %v", err)
	}
	if err := os.WriteFile(filepath.Join(apkContentsDir, "classes2.dex"), []byte("dex2"), 0644); err != nil {
		t.Fatalf("Failed to write classes2.dex: %v", err)
	}
	if err := os.WriteFile(filepath.Join(assetsDir, "sample.txt"), []byte("asset text"), 0644); err != nil {
		t.Fatalf("Failed to write asset file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(libDir, "libtest.so"), []byte("native library"), 0644); err != nil {
		t.Fatalf("Failed to write native lib: %v", err)
	}
	if err := os.WriteFile(filepath.Join(unknownDir, "config.json"), []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to write unknown file: %v", err)
	}

	// Zip injected items into sampleAPK
	cmd := exec.Command("zip", "-u", sampleAPK, "classes.dex", "classes2.dex", "assets/sample.txt", "lib/arm64-v8a/libtest.so", "unknown/config.json")
	cmd.Dir = apkContentsDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to zip additional files into sample APK (%v): %s", err, string(out))
	}

	// 4. Run the full conversion pipeline steps
	steps := []struct {
		name   string
		action func() error
	}{
		{"Clean environment", cleanEnvironment},
		{"Decompress APK", func() error { return decompressInputAPKPackage(sampleAPK) }},
		{"Compile resources", compileInputResources},
		{"Generate base APK", func() error { return generateOutputAPKBase("21", "34") }},
		{"Unzip base APK", unzipOutputAPKBase},
		{"Create output structure", createOutputStructure},
		{"Zip output structure", zipOutoutStructure},
		{"Generate AAB", func() error { return generateOutputAAB(sampleAPK) }},
		{"Validate AAB", func() error { return validateOutputAAB(sampleAPK) }},
	}

	for _, step := range steps {
		if err := step.action(); err != nil {
			t.Fatalf("Pipeline step '%s' failed: %v", step.name, err)
		}
	}

	sampleAAB := getAABPath(sampleAPK)
	if !fileInZipExists(t, sampleAAB, "base/manifest/AndroidManifest.xml") {
		t.Errorf("AAB missing base/manifest/AndroidManifest.xml")
	}
	if !fileInZipExists(t, sampleAAB, "base/dex/classes.dex") {
		t.Errorf("AAB missing base/dex/classes.dex")
	}
	if !fileInZipExists(t, sampleAAB, "base/dex/classes2.dex") {
		t.Errorf("AAB missing base/dex/classes2.dex")
	}
	if !fileInZipExists(t, sampleAAB, "base/assets/sample.txt") {
		t.Errorf("AAB missing base/assets/sample.txt")
	}
	if !fileInZipExists(t, sampleAAB, "base/lib/arm64-v8a/libtest.so") {
		t.Errorf("AAB missing base/lib/arm64-v8a/libtest.so")
	}
	if !fileInZipExists(t, sampleAAB, "base/resources.pb") {
		t.Errorf("AAB missing base/resources.pb")
	}

	cleanEnvironment()
}

func fileInZipExists(t *testing.T, zipPath string, targetFile string) bool {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		t.Fatalf("Failed to open zip file %s: %v", zipPath, err)
		return false
	}
	defer reader.Close()

	for _, file := range reader.File {
		if file.Name == targetFile {
			return true
		}
	}
	return false
}
