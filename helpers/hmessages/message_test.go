package hmessages

import (
   "strings"
   "testing"
)

func TestGetSuccessMessage(t *testing.T) {
   sResult := GetSuccessMessage("Operation completed")
   if !strings.Contains(sResult, "SUCCESS") || !strings.Contains(sResult, "Operation completed") {
      t.Errorf("Unexpected result for GetSuccessMessage: %s", sResult)
   }
}

func TestGetErrorMessage(t *testing.T) {
   sResult := GetErrorMessage("An error occurred")
   if !strings.Contains(sResult, "ERROR") || !strings.Contains(sResult, "An error occurred") {
      t.Errorf("Unexpected result for GetErrorMessage: %s", sResult)
   }
}

func TestGetInfoMessage(t *testing.T) {
   sResult := GetInfoMessage("Information note")
   if !strings.Contains(sResult, "INFO") || !strings.Contains(sResult, "Information note") {
      t.Errorf("Unexpected result for GetInfoMessage: %s", sResult)
   }
}

func TestGetRunningMessage(t *testing.T) {
   sResult := GetRunningMessage("Task running")
   if !strings.Contains(sResult, "RUNNING") || !strings.Contains(sResult, "Task running") {
      t.Errorf("Unexpected result for GetRunningMessage: %s", sResult)
   }
}
