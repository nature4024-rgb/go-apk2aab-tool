package hmessages

import (
   "strings"
   "testing"
)

func TestGetSuccessMessage(t *testing.T) {
   msg := GetSuccessMessage("Build completed")
   if !strings.Contains(msg, "SUCCESS") || !strings.Contains(msg, "Build completed") {
      t.Errorf("GetSuccessMessage failed: got %s", msg)
   }
}

func TestGetErrorMessage(t *testing.T) {
   msg := GetErrorMessage("Invalid path")
   if !strings.Contains(msg, "ERROR") || !strings.Contains(msg, "Invalid path") {
      t.Errorf("GetErrorMessage failed: got %s", msg)
   }
}

func TestGetInfoMessage(t *testing.T) {
   msg := GetInfoMessage("Initializing")
   if !strings.Contains(msg, "INFO") || !strings.Contains(msg, "Initializing") {
      t.Errorf("GetInfoMessage failed: got %s", msg)
   }
}

func TestGetProgressMessage(t *testing.T) {
   msg := GetProgressMessage("Processing")
   if !strings.Contains(msg, "WORKING") || !strings.Contains(msg, "Processing") {
      t.Errorf("GetProgressMessage failed: got %s", msg)
   }
}

func TestGetToastBanner(t *testing.T) {
   bannerSuccess := GetToastBanner("Success Title", "Operation completed", true)
   if !strings.Contains(bannerSuccess, "Success Title") || !strings.Contains(bannerSuccess, "Operation completed") {
      t.Errorf("GetToastBanner success failed: got %s", bannerSuccess)
   }

   bannerError := GetToastBanner("Error Title", "Operation failed", false)
   if !strings.Contains(bannerError, "Error Title") || !strings.Contains(bannerError, "Operation failed") {
      t.Errorf("GetToastBanner error failed: got %s", bannerError)
   }
}
