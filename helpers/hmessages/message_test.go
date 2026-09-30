package hmessages

import (
   "strings"
   "testing"
)

func TestGetSuccessMessage(t *testing.T) {
   sMsg := GetSuccessMessage("Done")
   if !strings.Contains(sMsg, "Done") || !strings.Contains(sMsg, "✔ SUCCESS") {
      t.Errorf("Unexpected success message output: %s", sMsg)
   }
}

func TestGetErrorMessage(t *testing.T) {
   sMsg := GetErrorMessage("Failed")
   if !strings.Contains(sMsg, "Failed") || !strings.Contains(sMsg, "✖ ERROR") {
      t.Errorf("Unexpected error message output: %s", sMsg)
   }
}

func TestGetInfoMessage(t *testing.T) {
   sMsg := GetInfoMessage("Notice")
   if !strings.Contains(sMsg, "Notice") || !strings.Contains(sMsg, "ℹ INFO") {
      t.Errorf("Unexpected info message output: %s", sMsg)
   }
}

func TestGetLoadingMessage(t *testing.T) {
   sMsg := GetLoadingMessage("Processing")
   if !strings.Contains(sMsg, "Processing") || !strings.Contains(sMsg, "⏳ RUNNING") {
      t.Errorf("Unexpected loading message output: %s", sMsg)
   }
}

func TestGetToastSuccess(t *testing.T) {
   sMsg := GetToastSuccess("Title", "Details")
   if !strings.Contains(sMsg, "Title") || !strings.Contains(sMsg, "Details") || !strings.Contains(sMsg, "TOAST") {
      t.Errorf("Unexpected toast success output: %s", sMsg)
   }
}

func TestGetToastError(t *testing.T) {
   sMsg := GetToastError("AlertTitle", "AlertDetails")
   if !strings.Contains(sMsg, "AlertTitle") || !strings.Contains(sMsg, "AlertDetails") || !strings.Contains(sMsg, "ALERT") {
      t.Errorf("Unexpected toast error output: %s", sMsg)
   }
}
