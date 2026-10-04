package main

import (
   "strings"
   "testing"
)

func TestGetLine(t *testing.T) {
   line := getLine()
   if len(line) == 0 {
      t.Errorf("expected non-empty line separator")
   }
}

func TestGetAppBanner(t *testing.T) {
   banner := getAppBanner()
   if !strings.Contains(banner, "________") {
      t.Errorf("expected ASCII art in app banner")
   }
}
