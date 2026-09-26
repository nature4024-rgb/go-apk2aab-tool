package hcolors

import (
   "testing"
)

func TestColorizeFunctions(t *testing.T) {
   s := "test"

   if InRed(s) == "" {
      t.Error("InRed returned empty string")
   }

   if InGreen(s) == "" {
      t.Error("InGreen returned empty string")
   }

   if InYellow(s) == "" {
      t.Error("InYellow returned empty string")
   }

   if InBlue(s) == "" {
      t.Error("InBlue returned empty string")
   }

   if InPurple(s) == "" {
      t.Error("InPurple returned empty string")
   }

   if InCyan(s) == "" {
      t.Error("InCyan returned empty string")
   }

   if InGray(s) == "" {
      t.Error("InGray returned empty string")
   }

   if InWhite(s) == "" {
      t.Error("InWhite returned empty string")
   }

   if InBold(s) == "" {
      t.Error("InBold returned empty string")
   }
}

func TestIzeAlias(t *testing.T) {
   res := Ize(Red, "hello")
   if res == "" {
      t.Error("Ize returned empty string")
   }
}
