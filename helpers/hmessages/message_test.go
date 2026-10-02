package hmessages

import (
	"apk2aab/helpers/hcolors"
	"strings"
	"testing"
)

func TestGetSuccessMessage(t *testing.T) {
	msg := GetSuccessMessage("Done")
	if !strings.Contains(msg, "SUCCESS") || !strings.Contains(msg, "Done") {
		t.Errorf("Unexpected GetSuccessMessage output: %q", msg)
	}
}

func TestGetErrorMessage(t *testing.T) {
	msg := GetErrorMessage("Failed")
	if !strings.Contains(msg, "ERROR") || !strings.Contains(msg, "Failed") {
		t.Errorf("Unexpected GetErrorMessage output: %q", msg)
	}
}

func TestGetInfoMessage(t *testing.T) {
	msg := GetInfoMessage("Processing")
	if !strings.Contains(msg, "INFO") || !strings.Contains(msg, "Processing") {
		t.Errorf("Unexpected GetInfoMessage output: %q", msg)
	}
}

func TestGetMessage(t *testing.T) {
	msg := GetMessage("Custom message", hcolors.Yellow, "WARN")
	if !strings.Contains(msg, "WARN") || !strings.Contains(msg, "Custom message") {
		t.Errorf("Unexpected GetMessage output: %q", msg)
	}
}
