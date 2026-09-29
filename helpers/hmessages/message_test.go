package hmessages

import (
   "strings"
   "testing"
)

func TestGetSuccessMessage(t *testing.T) {
   msg := GetSuccessMessage("Operation completed")
   if !strings.Contains(msg, "✔ SUCCESS") {
      t.Errorf("Expected message to contain '✔ SUCCESS', got: %s", msg)
   }
   if !strings.Contains(msg, "Operation completed") {
      t.Errorf("Expected message to contain 'Operation completed', got: %s", msg)
   }
}

func TestGetErrorMessage(t *testing.T) {
   msg := GetErrorMessage("Something went wrong")
   if !strings.Contains(msg, "✖ ERROR") {
      t.Errorf("Expected message to contain '✖ ERROR', got: %s", msg)
   }
   if !strings.Contains(msg, "Something went wrong") {
      t.Errorf("Expected message to contain 'Something went wrong', got: %s", msg)
   }
}

func TestGetInfoMessage(t *testing.T) {
   msg := GetInfoMessage("Informational text")
   if !strings.Contains(msg, "ℹ INFO") {
      t.Errorf("Expected message to contain 'ℹ INFO', got: %s", msg)
   }
   if !strings.Contains(msg, "Informational text") {
      t.Errorf("Expected message to contain 'Informational text', got: %s", msg)
   }
}

func TestGetRunningMessage(t *testing.T) {
   msg := GetRunningMessage("Processing step...")
   if !strings.Contains(msg, "⏳ RUNNING") {
      t.Errorf("Expected message to contain '⏳ RUNNING', got: %s", msg)
   }
   if !strings.Contains(msg, "Processing step...") {
      t.Errorf("Expected message to contain 'Processing step...', got: %s", msg)
   }
}
