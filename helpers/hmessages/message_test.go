package hmessages

import (
   "strings"
   "testing"
)

func TestGetSuccessMessage(t *testing.T) {
   msg := GetSuccessMessage("Operation completed")
   if !strings.Contains(msg, "✔ SUCCESS") {
      t.Errorf("Expected success message to contain '✔ SUCCESS', got %q", msg)
   }
   if !strings.Contains(msg, "Operation completed") {
      t.Errorf("Expected success message to contain payload, got %q", msg)
   }
}

func TestGetErrorMessage(t *testing.T) {
   msg := GetErrorMessage("Operation failed")
   if !strings.Contains(msg, "✖ ERROR") {
      t.Errorf("Expected error message to contain '✖ ERROR', got %q", msg)
   }
   if !strings.Contains(msg, "Operation failed") {
      t.Errorf("Expected error message to contain payload, got %q", msg)
   }
}

func TestGetInfoMessage(t *testing.T) {
   msg := GetInfoMessage("Processing step")
   if !strings.Contains(msg, "ℹ INFO") {
      t.Errorf("Expected info message to contain 'ℹ INFO', got %q", msg)
   }
   if !strings.Contains(msg, "Processing step") {
      t.Errorf("Expected info message to contain payload, got %q", msg)
   }
}

func TestGetProgressMessage(t *testing.T) {
   msg := GetProgressMessage("Working...")
   if !strings.Contains(msg, "⏳ BUSY") {
      t.Errorf("Expected progress message to contain '⏳ BUSY', got %q", msg)
   }
   if !strings.Contains(msg, "Working...") {
      t.Errorf("Expected progress message to contain payload, got %q", msg)
   }
}
