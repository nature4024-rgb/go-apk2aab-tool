package hmessages

import (
   "strings"
   "testing"
)

func TestGetSuccessMessage(t *testing.T) {
   msg := GetSuccessMessage("Done")
   if !strings.Contains(msg, "SUCCESS") || !strings.Contains(msg, "Done") {
      t.Errorf("Unexpected success message: %s", msg)
   }
}

func TestGetErrorMessage(t *testing.T) {
   msg := GetErrorMessage("Failed")
   if !strings.Contains(msg, "ERROR") || !strings.Contains(msg, "Failed") {
      t.Errorf("Unexpected error message: %s", msg)
   }
}

func TestGetInfoMessage(t *testing.T) {
   msg := GetInfoMessage("Loading")
   if !strings.Contains(msg, "INFO") || !strings.Contains(msg, "Loading") {
      t.Errorf("Unexpected info message: %s", msg)
   }
}
