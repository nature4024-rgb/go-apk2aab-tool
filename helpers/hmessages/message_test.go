package hmessages

import (
   "strings"
   "testing"
)

func TestGetStepMessage(t *testing.T) {
   msg := GetStepMessage(1, 7, "📦", "Decompressing input APK...")
   if !strings.Contains(msg, "[1/7 📦]") {
      t.Errorf("Expected step indicator '[1/7 📦]', got: %s", msg)
   }
   if !strings.Contains(msg, "Decompressing input APK...") {
      t.Errorf("Expected message text in step output, got: %s", msg)
   }
}

func TestGetToastNoticeSuccess(t *testing.T) {
   toast := GetToastNotice("Build Complete", "Output created successfully", true)
   if !strings.Contains(toast, "Build Complete") {
      t.Errorf("Expected title 'Build Complete' in toast, got: %s", toast)
   }
   if !strings.Contains(toast, "✨") {
      t.Errorf("Expected success icon '✨' in toast, got: %s", toast)
   }
}

func TestGetToastNoticeError(t *testing.T) {
   toast := GetToastNotice("Build Failed", "An error occurred", false)
   if !strings.Contains(toast, "Build Failed") {
      t.Errorf("Expected title 'Build Failed' in toast, got: %s", toast)
   }
   if !strings.Contains(toast, "⚠️") {
      t.Errorf("Expected error icon '⚠️' in toast, got: %s", toast)
   }
}
