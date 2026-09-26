package hmessages

import (
   "strings"
   "testing"
)

func TestGetSuccessMessage(t *testing.T) {
   msg := GetSuccessMessage("All good")
   if !strings.Contains(msg, "SUCCESS") || !strings.Contains(msg, "All good") {
      t.Errorf("Unexpected success message: %s", msg)
   }

   defaultMsg := GetSuccessMessage("")
   if !strings.Contains(defaultMsg, "✔ Done") {
      t.Errorf("Unexpected default success message: %s", defaultMsg)
   }
}

func TestGetErrorMessage(t *testing.T) {
   msg := GetErrorMessage("Invalid input")
   if !strings.Contains(msg, "ERROR") || !strings.Contains(msg, "Invalid input") {
      t.Errorf("Unexpected error message: %s", msg)
   }

   defaultMsg := GetErrorMessage("")
   if !strings.Contains(defaultMsg, "✖ Failed") {
      t.Errorf("Unexpected default error message: %s", defaultMsg)
   }
}

func TestGetInfoMessage(t *testing.T) {
   msg := GetInfoMessage("Processing")
   if !strings.Contains(msg, "INFO") || !strings.Contains(msg, "Processing") {
      t.Errorf("Unexpected info message: %s", msg)
   }
}

func TestGetWarningMessage(t *testing.T) {
   msg := GetWarningMessage("Low disk space")
   if !strings.Contains(msg, "WARN") || !strings.Contains(msg, "Low disk space") {
      t.Errorf("Unexpected warning message: %s", msg)
   }
}

func TestGetStepMessage(t *testing.T) {
   msg := GetStepMessage(2, 5, "Building bundle")
   if !strings.Contains(msg, "STEP 2/5") || !strings.Contains(msg, "Building bundle") {
      t.Errorf("Unexpected step message: %s", msg)
   }
}
