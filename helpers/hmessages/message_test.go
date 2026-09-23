package hmessages

import (
   "strings"
   "testing"
)

func TestGetSuccessMessage(t *testing.T) {
   msg := GetSuccessMessage("Operation completed")
   if !strings.Contains(msg, "✔ SUCCESS") || !strings.Contains(msg, "Operation completed") {
      t.Errorf("Unexpected success message output: %s", msg)
   }
}

func TestGetErrorMessage(t *testing.T) {
   msg := GetErrorMessage("An error occurred")
   if !strings.Contains(msg, "✖ ERROR") || !strings.Contains(msg, "An error occurred") {
      t.Errorf("Unexpected error message output: %s", msg)
   }
}

func TestGetInfoMessage(t *testing.T) {
   msg := GetInfoMessage("Informational text")
   if !strings.Contains(msg, "ℹ INFO") || !strings.Contains(msg, "Informational text") {
      t.Errorf("Unexpected info message output: %s", msg)
   }
}

func TestGetLoadingMessage(t *testing.T) {
   msg := GetLoadingMessage("Loading process")
   if !strings.Contains(msg, "⏳ BUSY") || !strings.Contains(msg, "Loading process") {
      t.Errorf("Unexpected loading message output: %s", msg)
   }
}
