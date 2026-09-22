package hstrings

import "testing"

func TestIsEmpty(t *testing.T) {
	if !IsEmpty("") {
		t.Errorf("IsEmpty(\"\") should be true")
	}
	if IsEmpty("hello") {
		t.Errorf("IsEmpty(\"hello\") should be false")
	}
}
