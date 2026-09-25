package hmessages

import (
   "strings"
   "testing"
)

func TestGetSuccessMessage(t *testing.T) {
   msg := GetSuccessMessage("Operation completed")
   if !strings.Contains(msg, "SUCCESS") {
      t.Errorf("Expected msg to contain SUCCESS, got: %s", msg)
   }
   if !strings.Contains(msg, "✔") {
      t.Errorf("Expected msg to contain ✔ icon, got: %s", msg)
   }
}

func TestGetErrorMessage(t *testing.T) {
   msg := GetErrorMessage("Operation failed")
   if !strings.Contains(msg, "ERROR") {
      t.Errorf("Expected msg to contain ERROR, got: %s", msg)
   }
   if !strings.Contains(msg, "✖") {
      t.Errorf("Expected msg to contain ✖ icon, got: %s", msg)
   }
}

func TestGetInfoMessage(t *testing.T) {
   msg := GetInfoMessage("Processing data")
   if !strings.Contains(msg, "INFO") {
      t.Errorf("Expected msg to contain INFO, got: %s", msg)
   }
   if !strings.Contains(msg, "ℹ") {
      t.Errorf("Expected msg to contain ℹ icon, got: %s", msg)
   }
}

func TestGetToastNotification(t *testing.T) {
   toastSuccess := GetToastNotification("Title Success", "Message Detail", true)
   if !strings.Contains(toastSuccess, "Title Success") || !strings.Contains(toastSuccess, "✔") {
      t.Errorf("Expected success toast to contain title and icon, got: %s", toastSuccess)
   }

   toastError := GetToastNotification("Title Error", "Error Detail", false)
   if !strings.Contains(toastError, "Title Error") || !strings.Contains(toastError, "✖") {
      t.Errorf("Expected error toast to contain title and icon, got: %s", toastError)
   }
}
