package main

import (
   "apk2aab/helpers/hcolors"
   "apk2aab/helpers/hcompressions"
   "apk2aab/helpers/hfiles"
   "apk2aab/helpers/hmessages"
   "apk2aab/helpers/hstrings"
   "fmt"
   "io/ioutil"
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
const APP_VERSION = "1.0"

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

func main() {
   if len(os.Args) == 5 {
      if hfiles.FileExists(os.Args[1]) && strings.ToLower(filepath.Ext(os.Args[1])) == hfiles.FILE_EXTENSION_APK && regexp.MustCompile(REGEX_NUMERIC).MatchString(os.Args[3]) && regexp.MustCompile(REGEX_NUMERIC).MatchString(os.Args[4]) {
         setAppConfig(os.Args[2], os.Args[3], os.Args[4])

         if hstrings.IsEmpty(oAppConfig.errorMessage) {
            var bError bool = true

            fmt.Println(getAppBanner())
            fmt.Println(" " + getLine() + "\r\n")

            // Prepare the environment
            fmt.Println(" " + hmessages.GetProgressMessage("Cleaning temporary environment..."))
            if cleanEnvironment() {
               fmt.Println(" " + hmessages.GetSuccessMessage("Environment cleaned"))

               // 1. Decompress input APK package
               fmt.Println(" " + hmessages.GetProgressMessage("Decompressing input APK package using apktool..."))
               if decompressInputAPKPackage(os.Args[1]) {
                  fmt.Println(" " + hmessages.GetSuccessMessage("Decompressed APK successfully"))

                  // 2. Compile input resources
                  fmt.Println(" " + hmessages.GetProgressMessage("Compiling input resources using aapt2..."))
                  if compileInputResources() {
                     fmt.Println(" " + hmessages.GetSuccessMessage("Compiled input resources"))

                     // 3. Generate output APK base
                     fmt.Println(" " + hmessages.GetProgressMessage("Generating output APK base using aapt2..."))
                     if generateOutputAPKBase(os.Args[3], os.Args[4]) {
                        fmt.Println(" " + hmessages.GetSuccessMessage("Generated APK base link"))

                        // 4. Unzip output APK base
                        fmt.Println(" " + hmessages.GetProgressMessage("Unzipping output APK base..."))
                        if unzipOutputAPKBase() {
                           fmt.Println(" " + hmessages.GetSuccessMessage("Unzipped output APK base"))

                           // 5. Create output structure
                           fmt.Println(" " + hmessages.GetProgressMessage("Creating output structure..."))
                           if createOutputStructure() {
                              fmt.Println(" " + hmessages.GetSuccessMessage("Created output structure"))

                              // 6. Compress output structure
                              fmt.Println(" " + hmessages.GetProgressMessage("Zipping output structure..."))
                              if zipOutoutStructure() {
                                 fmt.Println(" " + hmessages.GetSuccessMessage("Zipped output structure"))

                                 // 7. Generate output AAB
                                 fmt.Println(" " + hmessages.GetProgressMessage("Generating output AAB bundle using bundletool..."))
                                 if generateOutputAAB(os.Args[1]) {
                                    fmt.Println(" " + hmessages.GetSuccessMessage("Generated output AAB"))

                                    bError = false
                                 }
                              }
                           }
                        }
                     }
                  }
               }
            }

            fmt.Println(" " + getLine() + "\r\n")
            if bError {
               fmt.Println(hmessages.GetToastBanner("Conversion Failed", "Check parameter values or tool installation and try again.", false))
            } else {
               aabOutput := strings.ReplaceAll(os.Args[1], hfiles.FILE_EXTENSION_APK, hfiles.FILE_EXTENSION_AAB)
               fmt.Println(hmessages.GetToastBanner("Conversion Successful!", "Output bundle generated at: "+aabOutput, true))
            }
            fmt.Println(" " + getLine())

            cleanEnvironment()
         } else {
            fmt.Println(getAppBanner())
            fmt.Println(" " + getLine() + "\r\n")
            fmt.Println(hmessages.GetToastBanner("Environment Error", oAppConfig.errorMessage, false))
            fmt.Println(" " + getLine())
         }
      } else {
         fmt.Println(getAppBanner())
         fmt.Println(" " + getLine() + "\r\n")

         if hfiles.FileExists(os.Args[1]) && strings.ToLower(filepath.Ext(os.Args[1])) == hfiles.FILE_EXTENSION_APK {
            fmt.Println(" " + hmessages.GetErrorMessage("Validation Error: min-sdk-version and target-sdk-version must be numeric integers"))
         } else {
            if strings.ToLower(filepath.Ext(os.Args[1])) == hfiles.FILE_EXTENSION_APK {
               fmt.Println(" " + hmessages.GetErrorMessage("File Error: Input file '"+os.Args[1]+"' does not exist"))
            } else {
               fmt.Println(" " + hmessages.GetErrorMessage("File Error: '"+os.Args[1]+"' is not a valid .apk file"))
            }
         }

         fmt.Println(" " + getLine())
      }
   } else {
      fmt.Println(getAppBanner())
      fmt.Println(" " + getLine() + "\r\n")
      fmt.Println(" " + hmessages.GetMessage("Tool to transform Android APK packages into AAB (Android App Bundle)", hcolors.Yellow, "PURPOSE"))
      fmt.Println(" " + hmessages.GetMessage("apk2aab <apk_file> <build_tools_version> <min_sdk> <target_sdk>", hcolors.Green, "USAGE  "))
      fmt.Println(" " + hmessages.GetMessage("apk2aab app.apk 34.0.0 21 34", hcolors.Green, "EXAMPLE"))
      fmt.Println(" " + hmessages.GetMessage("app.aab (Generated in same folder)", hcolors.Green, "OUTPUT "))
      fmt.Println(" " + getLine() + "\r\n")
      fmt.Println(" 💡 Hint: Ensure ANDROID_HOME and JAVA_HOME environment variables are set correctly.")
      fmt.Println(" Author: " + APP_AUTHOR_NAME + " | Version: " + APP_VERSION)
   }
}

func setAppConfig(sBuildToolsVersion string, sMinSdkVersion string, sTargetSdkVersion string) {
   var sExecutableExtension string

   if runtime.GOOS == OS_PLATFORM_WINDOWS {
      sSeparatorCharacter = "\\"
      sExecutableExtension = hfiles.FILE_EXTENSION_EXE
   } else {
      sSeparatorCharacter = "/"
      sExecutableExtension = hstrings.STRING_EMPTY
   }

   oAppConfig.tempFolderPath = APP_FOLDER_TEMP

   // 1. Check if JAVA is installed
   var sOSEnvVarJava string = os.Getenv(OS_ENVIRONMENT_VAR_JAVA_HOME)

   if hstrings.IsEmpty(sOSEnvVarJava) {
      sOSEnvVarJava = os.Getenv(OS_ENVIRONMENT_VAR_JAVA_JRE)
   }

   if !hstrings.IsEmpty(sOSEnvVarJava) {
      oAppConfig.javaBinFilePath = sOSEnvVarJava + sSeparatorCharacter + "bin" + sSeparatorCharacter + "java"

      // 2. Check if apktool is available
      if hfiles.FileExists("tools" + sSeparatorCharacter + "apktool" + hfiles.FILE_EXTENSION_JAR) {
         oAppConfig.apktoolJarFilePath = "tools" + sSeparatorCharacter + "apktool" + hfiles.FILE_EXTENSION_JAR

         // 3. Check if bundletool is available
         if hfiles.FileExists("tools" + sSeparatorCharacter + "bundletool" + hfiles.FILE_EXTENSION_JAR) {
            oAppConfig.bundletoolJarFilePath = "tools" + sSeparatorCharacter + "bundletool" + hfiles.FILE_EXTENSION_JAR

            // 4. Check if SDK is installed
            var sOSEnvVarAndroid string = os.Getenv(OS_ENVIRONMENT_VAR_ANDROID_HOME)

            if hstrings.IsEmpty(sOSEnvVarAndroid) {
               sOSEnvVarAndroid = os.Getenv(OS_ENVIRONMENT_VAR_ANDROID_SDK_ROOT)
            }

            if !hstrings.IsEmpty(sOSEnvVarAndroid) {
               if hfiles.FileExists(sOSEnvVarAndroid + sSeparatorCharacter + "build-tools" + sSeparatorCharacter + sBuildToolsVersion + sSeparatorCharacter + "aapt2" + sExecutableExtension) {
                  oAppConfig.aapt2BinFilePath = sOSEnvVarAndroid + sSeparatorCharacter + "build-tools" + sSeparatorCharacter + sBuildToolsVersion + sSeparatorCharacter + "aapt2" + sExecutableExtension

                  // 5. Check if Android.jar exists
                  if hfiles.FileExists(sOSEnvVarAndroid + sSeparatorCharacter + "platforms" + sSeparatorCharacter + "android-" + sTargetSdkVersion + sSeparatorCharacter + "android" + hfiles.FILE_EXTENSION_JAR) {
                     oAppConfig.androidJarFilePath = sOSEnvVarAndroid + sSeparatorCharacter + "platforms" + sSeparatorCharacter + "android-" + sTargetSdkVersion + sSeparatorCharacter + "android" + hfiles.FILE_EXTENSION_JAR
                  } else {
                     oAppConfig.errorMessage = "Android" + hfiles.FILE_EXTENSION_JAR + " isn't available, please check that the jar exists in the {ANDROID_SDK}" + sSeparatorCharacter + "platforms" + sSeparatorCharacter + "android-" + sTargetSdkVersion + sSeparatorCharacter + " folder"
                  }
               } else {
                  oAppConfig.errorMessage = "Aapt2 isn't available, please check that the binary exists in the {ANDROID_SDK}" + sSeparatorCharacter + "build-tools" + sSeparatorCharacter + sBuildToolsVersion + " folder"
               }
            } else {
               oAppConfig.errorMessage = "Android Studio isn't installed the ANDROID_HOME or ANDROID_SDK_ROOT environment variables couldn't be detected"
            }
         } else {
            oAppConfig.errorMessage = "Bundletool isn't available, please download bundletool" + hfiles.FILE_EXTENSION_JAR + " and put it inside tools folder"
         }
      } else {
         oAppConfig.errorMessage = "Apktool isn't available, please download apktool" + hfiles.FILE_EXTENSION_JAR + " and put it inside tools folder"
      }
   } else {
      oAppConfig.errorMessage = "Java isn't installed, the JAVA_JDK or JAVA_JRE environment variables couldn't be detected"
   }
}

func cleanEnvironment() bool {
   var oPrepareEnvironment error = os.RemoveAll(APP_FOLDER_TEMP)

   return oPrepareEnvironment == nil
}

func decompressInputAPKPackage(sAPKFile string) bool {
   var oDecompressInputAPKPackage error = exec.Command(oAppConfig.javaBinFilePath, "-jar", oAppConfig.apktoolJarFilePath, "d", sAPKFile, "-s", "-o", oAppConfig.tempFolderPath+sSeparatorCharacter+APP_FOLDER_TEMP_INPUT).Run()

   return oDecompressInputAPKPackage == nil
}

func compileInputResources() bool {
   var oCompileInputResources error = exec.Command(oAppConfig.aapt2BinFilePath, "compile", "--dir", oAppConfig.tempFolderPath+sSeparatorCharacter+APP_FOLDER_TEMP_INPUT+sSeparatorCharacter+"res", "-o", oAppConfig.tempFolderPath+sSeparatorCharacter+"compiled_resources"+hfiles.FILE_EXTENSION_ZIP).Run()

   return oCompileInputResources == nil
}

func generateOutputAPKBase(sMinSdkVersion string, sTargetSdkVersion string) bool {
   var oGenerateOutputAPKBase error = exec.Command(oAppConfig.aapt2BinFilePath, "link", "--proto-format", "-o", oAppConfig.tempFolderPath+sSeparatorCharacter+APP_FILE_TEMP_OUTPUT+hfiles.FILE_EXTENSION_APK, "-I", oAppConfig.androidJarFilePath, "--min-sdk-version", sMinSdkVersion, "--target-sdk-version", sTargetSdkVersion, "--version-code", "1", "--version-name", "1.0", "--manifest", oAppConfig.tempFolderPath+sSeparatorCharacter+APP_FOLDER_TEMP_INPUT+sSeparatorCharacter+"AndroidManifest"+hfiles.FILE_EXTENSION_XML, "-R", oAppConfig.tempFolderPath+sSeparatorCharacter+"compiled_resources"+hfiles.FILE_EXTENSION_ZIP, "--auto-add-overlay").Run()

   return oGenerateOutputAPKBase == nil
}

func unzipOutputAPKBase() bool {
   var oUnzipOutputAPKBase error = hcompressions.ZipDecompression(oAppConfig.tempFolderPath+sSeparatorCharacter+APP_FILE_TEMP_OUTPUT+hfiles.FILE_EXTENSION_APK, oAppConfig.tempFolderPath+sSeparatorCharacter+APP_FOLDER_TEMP_OUTPUT)

   return oUnzipOutputAPKBase == nil
}

func createOutputStructure() bool {
   var oCreateOutputStructure error = nil

   // Move AndroidManifest.xml file
   if hfiles.FileExists(oAppConfig.tempFolderPath + sSeparatorCharacter + APP_FOLDER_TEMP_OUTPUT + sSeparatorCharacter + "AndroidManifest" + hfiles.FILE_EXTENSION_XML) {
      oCreateOutputStructure = os.Mkdir(oAppConfig.tempFolderPath+sSeparatorCharacter+APP_FOLDER_TEMP_OUTPUT+sSeparatorCharacter+"manifest", os.ModePerm)

      if oCreateOutputStructure == nil {
         oCreateOutputStructure = os.Rename(oAppConfig.tempFolderPath+sSeparatorCharacter+APP_FOLDER_TEMP_OUTPUT+sSeparatorCharacter+"AndroidManifest"+hfiles.FILE_EXTENSION_XML, oAppConfig.tempFolderPath+sSeparatorCharacter+APP_FOLDER_TEMP_OUTPUT+sSeparatorCharacter+"manifest"+sSeparatorCharacter+"AndroidManifest"+hfiles.FILE_EXTENSION_XML)
      }
   }

   // Move assets folder
   if oCreateOutputStructure == nil {
      if hfiles.FolderExists(oAppConfig.tempFolderPath + sSeparatorCharacter + APP_FOLDER_TEMP_INPUT + sSeparatorCharacter + "assets") {
         oCreateOutputStructure = os.Rename(oAppConfig.tempFolderPath+sSeparatorCharacter+APP_FOLDER_TEMP_INPUT+sSeparatorCharacter+"assets", oAppConfig.tempFolderPath+sSeparatorCharacter+APP_FOLDER_TEMP_OUTPUT+sSeparatorCharacter+"assets")
      }
   }

   // Move lib folder
   if oCreateOutputStructure == nil {
      if hfiles.FolderExists(oAppConfig.tempFolderPath + sSeparatorCharacter + APP_FOLDER_TEMP_INPUT + sSeparatorCharacter + "lib") {
         oCreateOutputStructure = os.Rename(oAppConfig.tempFolderPath+sSeparatorCharacter+APP_FOLDER_TEMP_INPUT+sSeparatorCharacter+"lib", oAppConfig.tempFolderPath+sSeparatorCharacter+APP_FOLDER_TEMP_OUTPUT+sSeparatorCharacter+"lib")
      }
   }

   // Create root folder
   if !hfiles.FolderExists(oAppConfig.tempFolderPath + sSeparatorCharacter + APP_FOLDER_TEMP_OUTPUT + sSeparatorCharacter + "root") {
      oCreateOutputStructure = os.Mkdir(oAppConfig.tempFolderPath+sSeparatorCharacter+APP_FOLDER_TEMP_OUTPUT+sSeparatorCharacter+"root", os.ModePerm)
   }

   // Move kotlin folder
   if oCreateOutputStructure == nil {
      if hfiles.FolderExists(oAppConfig.tempFolderPath + sSeparatorCharacter + APP_FOLDER_TEMP_INPUT + sSeparatorCharacter + "kotlin") {
         oCreateOutputStructure = os.Rename(oAppConfig.tempFolderPath+sSeparatorCharacter+APP_FOLDER_TEMP_INPUT+sSeparatorCharacter+"kotlin", oAppConfig.tempFolderPath+sSeparatorCharacter+APP_FOLDER_TEMP_OUTPUT+sSeparatorCharacter+"root"+sSeparatorCharacter+"kotlin")
      }
   }

   // Move meta-inf folder
   if oCreateOutputStructure == nil {
      if hfiles.FolderExists(oAppConfig.tempFolderPath + sSeparatorCharacter + APP_FOLDER_TEMP_INPUT + sSeparatorCharacter + "original" + sSeparatorCharacter + "meta-inf") {
         oCreateOutputStructure = os.Rename(oAppConfig.tempFolderPath+sSeparatorCharacter+APP_FOLDER_TEMP_INPUT+sSeparatorCharacter+"original"+sSeparatorCharacter+"meta-inf", oAppConfig.tempFolderPath+sSeparatorCharacter+APP_FOLDER_TEMP_OUTPUT+sSeparatorCharacter+"root"+sSeparatorCharacter+"meta-inf")
      } else if hfiles.FolderExists(oAppConfig.tempFolderPath + sSeparatorCharacter + APP_FOLDER_TEMP_INPUT + sSeparatorCharacter + "original" + sSeparatorCharacter + "META-INF") {
         oCreateOutputStructure = os.Rename(oAppConfig.tempFolderPath+sSeparatorCharacter+APP_FOLDER_TEMP_INPUT+sSeparatorCharacter+"original"+sSeparatorCharacter+"META-INF", oAppConfig.tempFolderPath+sSeparatorCharacter+APP_FOLDER_TEMP_OUTPUT+sSeparatorCharacter+"root"+sSeparatorCharacter+"meta-inf")
      }
   }

   // Move all .dex files
   oFiles, oCreateOutputStructure := ioutil.ReadDir(oAppConfig.tempFolderPath + sSeparatorCharacter + APP_FOLDER_TEMP_INPUT)
   if oCreateOutputStructure == nil {
      var bFirstDexFile = true
      for _, oFiles := range oFiles {
         if strings.Contains(oFiles.Name(), hfiles.FILE_EXTENSION_DEX) {
            if bFirstDexFile {
               oCreateOutputStructure = os.Mkdir(oAppConfig.tempFolderPath+sSeparatorCharacter+APP_FOLDER_TEMP_OUTPUT+sSeparatorCharacter+"dex", os.ModePerm)
            }

            if oCreateOutputStructure == nil {
               oCreateOutputStructure = os.Rename(oAppConfig.tempFolderPath+sSeparatorCharacter+APP_FOLDER_TEMP_INPUT+sSeparatorCharacter+oFiles.Name(), oAppConfig.tempFolderPath+sSeparatorCharacter+APP_FOLDER_TEMP_OUTPUT+sSeparatorCharacter+"dex"+sSeparatorCharacter+oFiles.Name())
            }

            bFirstDexFile = false
         }
      }
   }

   return oCreateOutputStructure == nil
}

func zipOutoutStructure() bool {
   var oZipOutoutStructure error = hcompressions.ZipCompression(oAppConfig.tempFolderPath+sSeparatorCharacter+APP_FOLDER_TEMP_OUTPUT+sSeparatorCharacter, oAppConfig.tempFolderPath+sSeparatorCharacter+APP_FILE_TEMP_OUTPUT+hfiles.FILE_EXTENSION_ZIP, false)

   return oZipOutoutStructure == nil
}

func generateOutputAAB(sAPKFile string) bool {
   var oGenerateOutputAAB error = nil

   if hfiles.FileExists(strings.ReplaceAll(sAPKFile, hfiles.FILE_EXTENSION_APK, hfiles.FILE_EXTENSION_AAB)) {
      oGenerateOutputAAB = os.Remove(strings.ReplaceAll(sAPKFile, hfiles.FILE_EXTENSION_APK, hfiles.FILE_EXTENSION_AAB))
   }

   if oGenerateOutputAAB == nil {
      oGenerateOutputAAB = exec.Command(oAppConfig.javaBinFilePath, "-jar", oAppConfig.bundletoolJarFilePath, "build-bundle", "--modules="+oAppConfig.tempFolderPath+sSeparatorCharacter+APP_FILE_TEMP_OUTPUT+hfiles.FILE_EXTENSION_ZIP, "--output="+strings.ReplaceAll(sAPKFile, hfiles.FILE_EXTENSION_APK, hfiles.FILE_EXTENSION_AAB)).Run()
   }

   return oGenerateOutputAAB == nil
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
   sAppBanner += "    \\ \\__\\ \\__\\ \\__\\    \\ \\__\\\\ \\__\\|\\________\\ \\__\\ \\__\\ \\__\\ \\__\\ \\_______\\\r\n"
   sAppBanner += "     \\|__|\\|__|\\|__|    \\|__| \\\\|__| \\|_______|\\|__|\\|__|\\|__|\\|__|\\|_______|"

   return sAppBanner
}
