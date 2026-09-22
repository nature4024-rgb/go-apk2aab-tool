package main

import (
	"apk2aab/helpers/hcolors"
	"apk2aab/helpers/hcompressions"
	"apk2aab/helpers/hfiles"
	"apk2aab/helpers/hmessages"
	"apk2aab/helpers/hstrings"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

type appConfig struct {
	javaBinFilePath       string
	aapt2BinFilePath      string
	apktoolJarFilePath    string
	bundletoolJarFilePath string
	androidJarFilePath    string
	tempFolderPath        string
	errorMessage          string
}

const APP_AUTHOR_NAME = "Ivan Ricart Borges"
const APP_VERSION = "1.1"

const APP_FOLDER_TEMP = "temp"
const APP_FOLDER_TEMP_INPUT = "input"
const APP_FOLDER_TEMP_OUTPUT = "output"
const APP_FILE_TEMP_OUTPUT = "output"

const OS_PLATFORM_WINDOWS = "windows"
const OS_ENVIRONMENT_VAR_JAVA_HOME = "JAVA_HOME"
const OS_ENVIRONMENT_VAR_JAVA_JRE = "JAVA_JRE"
const OS_ENVIRONMENT_VAR_ANDROID_HOME = "ANDROID_HOME"
const OS_ENVIRONMENT_VAR_ANDROID_SDK_ROOT = "ANDROID_SDK_ROOT"

const REGEX_NUMERIC = "^[0-9]+$"

var oAppConfig *appConfig = new(appConfig)
var sSeparatorCharacter string
var sExecutableExtension string

func main() {
	if len(os.Args) == 5 {
		if hfiles.FileExists(os.Args[1]) && strings.ToLower(filepath.Ext(os.Args[1])) == hfiles.FILE_EXTENSION_APK && regexp.MustCompile(REGEX_NUMERIC).MatchString(os.Args[3]) && regexp.MustCompile(REGEX_NUMERIC).MatchString(os.Args[4]) {
			setAppConfig(os.Args[2], os.Args[3], os.Args[4])

			if hstrings.IsEmpty(oAppConfig.errorMessage) {
				fmt.Println(getAppBanner())
				fmt.Println(" " + getLine() + "\r\n")

				steps := []struct {
					name   string
					action func() error
				}{
					{"Clean the environment...", cleanEnvironment},
					{"Decompress input APK package using apktool...", func() error { return decompressInputAPKPackage(os.Args[1]) }},
					{"Compiling input resources using aapt2...", compileInputResources},
					{"Generating output APK base using aapt2...", func() error { return generateOutputAPKBase(os.Args[3], os.Args[4]) }},
					{"Unzipping output APK base...", unzipOutputAPKBase},
					{"Creating output structure...", createOutputStructure},
					{"Zipping output structure...", zipOutoutStructure},
					{"Generating output AAB...", func() error { return generateOutputAAB(os.Args[1]) }},
					{"Validating generated AAB using bundletool...", func() error { return validateOutputAAB(os.Args[1]) }},
				}

				bError := false
				for _, step := range steps {
					if !runStep(step.name, step.action) {
						bError = true
						break
					}
				}

				if bError {
					fmt.Println(" " + hmessages.GetErrorMessage("Conversion failed."))
				} else {
					fmt.Println(" " + hmessages.GetSuccessMessage("AAB generated and validated successfully!"))
				}

				fmt.Println(" " + getLine())

				cleanEnvironment()
			} else {
				fmt.Println(getAppBanner())
				fmt.Println(" " + getLine() + "\r\n")
				fmt.Println(" " + hmessages.GetErrorMessage(oAppConfig.errorMessage))
				fmt.Println(" " + getLine())
			}
		} else {
			fmt.Println(getAppBanner())
			fmt.Println(" " + getLine() + "\r\n")

			if hfiles.FileExists(os.Args[1]) && strings.ToLower(filepath.Ext(os.Args[1])) == hfiles.FILE_EXTENSION_APK {
				fmt.Println(" " + hmessages.GetErrorMessage("Parameters min-sdk-version and target-sdk-version must be numeric"))
			} else {
				if strings.ToLower(filepath.Ext(os.Args[1])) == hfiles.FILE_EXTENSION_APK {
					fmt.Println(" " + hmessages.GetErrorMessage("File "+os.Args[1]+" not exists"))
				} else {
					fmt.Println(" " + hmessages.GetErrorMessage("File "+os.Args[1]+" isn't APK file"))
				}
			}

			fmt.Println(" " + getLine())
		}
	} else {
		fmt.Println(getAppBanner())
		fmt.Println(" " + getLine() + "\r\n")
		fmt.Println(" " + hmessages.GetMessage("Application to transform a file with APK format to AAB", hcolors.Yellow, "INFO   "))
		fmt.Println(" " + hmessages.GetMessage("apk2aab file-apk build-tools-version min-sdk-version target-sdk-version", hcolors.Green, "INPUT  "))
		fmt.Println(" " + hmessages.GetMessage("apk2aab file.apk 34.0.0 21 34", hcolors.Green, "EXAMPLE"))
		fmt.Println(" " + hmessages.GetMessage("file.aab", hcolors.Green, "OUTPUT "))
		fmt.Println(" " + getLine() + "\r\n")
		fmt.Println(" Author: " + APP_AUTHOR_NAME + " | Version: " + APP_VERSION)
	}
}

func runStep(stepName string, action func() error) bool {
	fmt.Print(" " + hmessages.GetInfoMessage(stepName))
	err := action()
	if err != nil {
		fmt.Println(" " + hmessages.GetErrorMessage("FAILED"))
		cleanErrStr := strings.ReplaceAll(err.Error(), "\r", "")
		fmt.Printf("\nError details:\n%s\n\n", cleanErrStr)
		return false
	}
	fmt.Println(" " + hmessages.GetSuccessMessage(hstrings.STRING_EMPTY))
	return true
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		outStr := strings.ReplaceAll(string(output), "\r", "")
		if len(outStr) > 0 {
			return fmt.Errorf("command execution failed (%v):\n%s", err, outStr)
		}
		return fmt.Errorf("command execution failed (%v)", err)
	}
	return nil
}

func setAppConfig(sBuildToolsVersion string, sMinSdkVersion string, sTargetSdkVersion string) {
	if runtime.GOOS == OS_PLATFORM_WINDOWS {
		sSeparatorCharacter = "\\"
		sExecutableExtension = hfiles.FILE_EXTENSION_EXE
	} else {
		sSeparatorCharacter = "/"
		sExecutableExtension = hstrings.STRING_EMPTY
	}

	oAppConfig.tempFolderPath = APP_FOLDER_TEMP

	// 1. Check if JAVA is available
	oAppConfig.javaBinFilePath = findJavaBinary()
	if hstrings.IsEmpty(oAppConfig.javaBinFilePath) {
		oAppConfig.errorMessage = "Java isn't installed or couldn't be found in JAVA_HOME, JAVA_JRE, or PATH"
		return
	}

	// 2. Check if apktool is available
	oAppConfig.apktoolJarFilePath = filepath.Join("tools", "apktool"+hfiles.FILE_EXTENSION_JAR)
	if !hfiles.FileExists(oAppConfig.apktoolJarFilePath) {
		oAppConfig.errorMessage = "Apktool isn't available, please download apktool" + hfiles.FILE_EXTENSION_JAR + " and put it inside tools folder"
		return
	}

	// 3. Check if bundletool is available
	oAppConfig.bundletoolJarFilePath = filepath.Join("tools", "bundletool"+hfiles.FILE_EXTENSION_JAR)
	if !hfiles.FileExists(oAppConfig.bundletoolJarFilePath) {
		oAppConfig.errorMessage = "Bundletool isn't available, please download bundletool" + hfiles.FILE_EXTENSION_JAR + " and put it inside tools folder"
		return
	}

	// 4. Check if AAPT2 is available
	oAppConfig.aapt2BinFilePath = findAapt2Binary(sBuildToolsVersion)
	if hstrings.IsEmpty(oAppConfig.aapt2BinFilePath) {
		oAppConfig.errorMessage = "Aapt2 isn't available, please check that build-tools " + sBuildToolsVersion + " is installed in ANDROID_HOME or aapt2 is in PATH"
		return
	}

	// 5. Check if android.jar is available
	oAppConfig.androidJarFilePath = findAndroidJar(sTargetSdkVersion)
	if hstrings.IsEmpty(oAppConfig.androidJarFilePath) {
		oAppConfig.errorMessage = "Android" + hfiles.FILE_EXTENSION_JAR + " isn't available for target API " + sTargetSdkVersion + ", please check ANDROID_HOME platforms folder"
		return
	}
}

func findJavaBinary() string {
	sOSEnvVarJava := os.Getenv(OS_ENVIRONMENT_VAR_JAVA_HOME)
	if hstrings.IsEmpty(sOSEnvVarJava) {
		sOSEnvVarJava = os.Getenv(OS_ENVIRONMENT_VAR_JAVA_JRE)
	}

	if !hstrings.IsEmpty(sOSEnvVarJava) {
		javaPath := filepath.Join(sOSEnvVarJava, "bin", "java"+sExecutableExtension)
		if hfiles.FileExists(javaPath) {
			return javaPath
		}
	}

	if path, err := exec.LookPath("java" + sExecutableExtension); err == nil {
		return path
	}

	return hstrings.STRING_EMPTY
}

func findAapt2Binary(sBuildToolsVersion string) string {
	sOSEnvVarAndroid := os.Getenv(OS_ENVIRONMENT_VAR_ANDROID_HOME)
	if hstrings.IsEmpty(sOSEnvVarAndroid) {
		sOSEnvVarAndroid = os.Getenv(OS_ENVIRONMENT_VAR_ANDROID_SDK_ROOT)
	}

	if !hstrings.IsEmpty(sOSEnvVarAndroid) {
		// Try requested build-tools version
		aapt2Path := filepath.Join(sOSEnvVarAndroid, "build-tools", sBuildToolsVersion, "aapt2"+sExecutableExtension)
		if hfiles.FileExists(aapt2Path) {
			return aapt2Path
		}

		// Fallback: search in any other available build-tools version
		buildToolsDir := filepath.Join(sOSEnvVarAndroid, "build-tools")
		if entries, err := os.ReadDir(buildToolsDir); err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					candidate := filepath.Join(buildToolsDir, entry.Name(), "aapt2"+sExecutableExtension)
					if hfiles.FileExists(candidate) {
						return candidate
					}
				}
			}
		}
	}

	if path, err := exec.LookPath("aapt2" + sExecutableExtension); err == nil {
		return path
	}

	return hstrings.STRING_EMPTY
}

func findAndroidJar(sTargetSdkVersion string) string {
	sOSEnvVarAndroid := os.Getenv(OS_ENVIRONMENT_VAR_ANDROID_HOME)
	if hstrings.IsEmpty(sOSEnvVarAndroid) {
		sOSEnvVarAndroid = os.Getenv(OS_ENVIRONMENT_VAR_ANDROID_SDK_ROOT)
	}

	if !hstrings.IsEmpty(sOSEnvVarAndroid) {
		// Try target SDK platform
		jarPath := filepath.Join(sOSEnvVarAndroid, "platforms", "android-"+sTargetSdkVersion, "android"+hfiles.FILE_EXTENSION_JAR)
		if hfiles.FileExists(jarPath) {
			return jarPath
		}

		// Fallback: search any available platforms directory
		platformsDir := filepath.Join(sOSEnvVarAndroid, "platforms")
		if entries, err := os.ReadDir(platformsDir); err == nil {
			var fallbackPath string
			for _, entry := range entries {
				if entry.IsDir() && strings.HasPrefix(entry.Name(), "android-") {
					candidate := filepath.Join(platformsDir, entry.Name(), "android"+hfiles.FILE_EXTENSION_JAR)
					if hfiles.FileExists(candidate) {
						fallbackPath = candidate
					}
				}
			}
			if !hstrings.IsEmpty(fallbackPath) {
				return fallbackPath
			}
		}
	}

	return hstrings.STRING_EMPTY
}

func cleanEnvironment() error {
	return os.RemoveAll(APP_FOLDER_TEMP)
}

func decompressInputAPKPackage(sAPKFile string) error {
	inputDir := filepath.Join(oAppConfig.tempFolderPath, APP_FOLDER_TEMP_INPUT)
	return runCommand(oAppConfig.javaBinFilePath, "-jar", oAppConfig.apktoolJarFilePath, "d", sAPKFile, "-s", "-o", inputDir, "-f")
}

func compileInputResources() error {
	resDir := filepath.Join(oAppConfig.tempFolderPath, APP_FOLDER_TEMP_INPUT, "res")
	compiledZip := filepath.Join(oAppConfig.tempFolderPath, "compiled_resources"+hfiles.FILE_EXTENSION_ZIP)
	return runCommand(oAppConfig.aapt2BinFilePath, "compile", "--dir", resDir, "-o", compiledZip)
}

func generateOutputAPKBase(sMinSdkVersion string, sTargetSdkVersion string) error {
	outputAPK := filepath.Join(oAppConfig.tempFolderPath, APP_FILE_TEMP_OUTPUT+hfiles.FILE_EXTENSION_APK)
	manifestPath := filepath.Join(oAppConfig.tempFolderPath, APP_FOLDER_TEMP_INPUT, "AndroidManifest"+hfiles.FILE_EXTENSION_XML)
	compiledZip := filepath.Join(oAppConfig.tempFolderPath, "compiled_resources"+hfiles.FILE_EXTENSION_ZIP)

	return runCommand(oAppConfig.aapt2BinFilePath, "link", "--proto-format",
		"-o", outputAPK,
		"-I", oAppConfig.androidJarFilePath,
		"--min-sdk-version", sMinSdkVersion,
		"--target-sdk-version", sTargetSdkVersion,
		"--version-code", "1",
		"--version-name", "1.0",
		"--manifest", manifestPath,
		"-R", compiledZip,
		"--auto-add-overlay",
	)
}

func unzipOutputAPKBase() error {
	outputAPK := filepath.Join(oAppConfig.tempFolderPath, APP_FILE_TEMP_OUTPUT+hfiles.FILE_EXTENSION_APK)
	outputDir := filepath.Join(oAppConfig.tempFolderPath, APP_FOLDER_TEMP_OUTPUT)
	return hcompressions.ZipDecompression(outputAPK, outputDir)
}

func createOutputStructure() error {
	baseInput := filepath.Join(oAppConfig.tempFolderPath, APP_FOLDER_TEMP_INPUT)
	baseOutput := filepath.Join(oAppConfig.tempFolderPath, APP_FOLDER_TEMP_OUTPUT)

	// 1. Move AndroidManifest.xml file
	manifestSrc := filepath.Join(baseOutput, "AndroidManifest"+hfiles.FILE_EXTENSION_XML)
	if hfiles.FileExists(manifestSrc) {
		manifestDir := filepath.Join(baseOutput, "manifest")
		if err := os.MkdirAll(manifestDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create manifest folder: %w", err)
		}
		manifestDst := filepath.Join(manifestDir, "AndroidManifest"+hfiles.FILE_EXTENSION_XML)
		if err := os.Rename(manifestSrc, manifestDst); err != nil {
			return fmt.Errorf("failed to move AndroidManifest.xml: %w", err)
		}
	}

	// 2. Move assets folder
	assetsSrc := filepath.Join(baseInput, "assets")
	if hfiles.FolderExists(assetsSrc) {
		assetsDst := filepath.Join(baseOutput, "assets")
		if err := os.Rename(assetsSrc, assetsDst); err != nil {
			return fmt.Errorf("failed to move assets folder: %w", err)
		}
	}

	// 3. Move lib folder
	libSrc := filepath.Join(baseInput, "lib")
	if hfiles.FolderExists(libSrc) {
		libDst := filepath.Join(baseOutput, "lib")
		if err := os.Rename(libSrc, libDst); err != nil {
			return fmt.Errorf("failed to move lib folder: %w", err)
		}
	}

	// 4. Ensure root folder exists
	rootDir := filepath.Join(baseOutput, "root")
	if err := os.MkdirAll(rootDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create root folder: %w", err)
	}

	// 5. Move kotlin folder
	kotlinSrc := filepath.Join(baseInput, "kotlin")
	if hfiles.FolderExists(kotlinSrc) {
		kotlinDst := filepath.Join(rootDir, "kotlin")
		if err := os.Rename(kotlinSrc, kotlinDst); err != nil {
			return fmt.Errorf("failed to move kotlin folder: %w", err)
		}
	}

	// 6. Move meta-inf folder
	metaInfSrc := filepath.Join(baseInput, "original", "meta-inf")
	if !hfiles.FolderExists(metaInfSrc) {
		metaInfSrc = filepath.Join(baseInput, "original", "META-INF")
	}
	if hfiles.FolderExists(metaInfSrc) {
		metaInfDst := filepath.Join(rootDir, "meta-inf")
		if err := os.Rename(metaInfSrc, metaInfDst); err != nil {
			return fmt.Errorf("failed to move meta-inf folder: %w", err)
		}
	}

	// 7. Move unknown folder contents to root
	unknownSrc := filepath.Join(baseInput, "unknown")
	if hfiles.FolderExists(unknownSrc) {
		entries, err := os.ReadDir(unknownSrc)
		if err == nil {
			for _, entry := range entries {
				srcPath := filepath.Join(unknownSrc, entry.Name())
				dstPath := filepath.Join(rootDir, entry.Name())
				if err := os.Rename(srcPath, dstPath); err != nil {
					return fmt.Errorf("failed to move unknown item %s: %w", entry.Name(), err)
				}
			}
		}
	}

	// 8. Move all .dex files
	dexDir := filepath.Join(baseOutput, "dex")
	entries, err := os.ReadDir(baseInput)
	if err == nil {
		hasDex := false
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), hfiles.FILE_EXTENSION_DEX) {
				if !hasDex {
					if err := os.MkdirAll(dexDir, os.ModePerm); err != nil {
						return fmt.Errorf("failed to create dex folder: %w", err)
					}
					hasDex = true
				}
				srcPath := filepath.Join(baseInput, entry.Name())
				dstPath := filepath.Join(dexDir, entry.Name())
				if err := os.Rename(srcPath, dstPath); err != nil {
					return fmt.Errorf("failed to move dex file %s: %w", entry.Name(), err)
				}
			}
		}
	}

	return nil
}

func zipOutoutStructure() error {
	src := filepath.Join(oAppConfig.tempFolderPath, APP_FOLDER_TEMP_OUTPUT) + string(filepath.Separator)
	dst := filepath.Join(oAppConfig.tempFolderPath, APP_FILE_TEMP_OUTPUT+hfiles.FILE_EXTENSION_ZIP)
	return hcompressions.ZipCompression(src, dst, false)
}

func getAABPath(sAPKFile string) string {
	ext := filepath.Ext(sAPKFile)
	return sAPKFile[:len(sAPKFile)-len(ext)] + hfiles.FILE_EXTENSION_AAB
}

func generateOutputAAB(sAPKFile string) error {
	aabPath := getAABPath(sAPKFile)
	if hfiles.FileExists(aabPath) {
		if err := os.Remove(aabPath); err != nil {
			return fmt.Errorf("failed to remove existing AAB file: %w", err)
		}
	}
	moduleZip := filepath.Join(oAppConfig.tempFolderPath, APP_FILE_TEMP_OUTPUT+hfiles.FILE_EXTENSION_ZIP)
	return runCommand(oAppConfig.javaBinFilePath, "-jar", oAppConfig.bundletoolJarFilePath, "build-bundle", "--modules="+moduleZip, "--output="+aabPath)
}

func validateOutputAAB(sAPKFile string) error {
	aabPath := getAABPath(sAPKFile)
	if !hfiles.FileExists(aabPath) {
		return fmt.Errorf("AAB file does not exist: %s", aabPath)
	}
	return runCommand(oAppConfig.javaBinFilePath, "-jar", oAppConfig.bundletoolJarFilePath, "validate", "--bundle="+aabPath)
}

func getLine() string {
	return "____________________________________________________________________________"
}

func getAppBanner() string {
	var sAppBanner string

	sAppBanner = "  ________  ________  ___  __      _______  ________  ________  ________\r\n"
	sAppBanner += " |\\   __  \\|\\   __  \\|\\  \\|\\  \\   /  ___  \\|\\   __  \\|\\   __  \\|\\   __  \\\r\n"
	sAppBanner += " \\ \\  \\|\\  \\ \\  \\|\\  \\ \\  \\/  /|_/__/|_/  /\\ \\  \\|\\  \\ \\  \\|\\  \\ \\  \\|\\ /_\r\n"
	sAppBanner += "  \\ \\   __  \\ \\   ____\\ \\   ___  \\__|//  / /\\ \\   __  \\ \\   __  \\ \\   __  \\\r\n"
	sAppBanner += "   \\ \\  \\ \\  \\ \\  \\___|\\ \\  \\\\ \\  \\  /  /_/__\\ \\  \\ \\  \\ \\  \\ \\  \\ \\  \\|\\  \\\r\n"
	sAppBanner += "    \\ \\__\\ \\__\\ \\__\\    \\ \\__\\\\ \\__\\|\\________\\ \\__\\ \\__\\ \\__\\ \\__\\ \\_______\r\n"
	sAppBanner += "     \\|__|\\|__|\\|__|    \\|__| \\|__| \\|_______|\\|__|\\|__|\\|__|\\|__|\\|_______|"

	return sAppBanner
}
