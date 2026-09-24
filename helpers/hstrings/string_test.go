package hstrings

import "testing"

func TestIsEmpty(t *testing.T) {
   if !IsEmpty("") {
      t.Errorf("Expected IsEmpty(\"\") to be true")
   }
   if IsEmpty("hello") {
      t.Errorf("Expected IsEmpty(\"hello\") to be false")
   }
}
