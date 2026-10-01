package hmessages

import (
   "apk2aab/helpers/hcolors"
   "apk2aab/helpers/hstrings"
   "fmt"
   "strings"
)

func GetSuccessMessage(sMessage string) string {
   if hstrings.IsEmpty(sMessage) {
      return GetMessage("✔ Done", hcolors.Green, "SUCCESS")
   }
   return GetMessage("✔ "+sMessage, hcolors.Green, "SUCCESS")
}

func GetErrorMessage(sMessage string) string {
   if hstrings.IsEmpty(sMessage) {
      return GetMessage("✖ Failed", hcolors.Red, "ERROR  ")
   }
   return GetMessage("✖ "+sMessage, hcolors.Red, "ERROR  ")
}

func GetInfoMessage(sMessage string) string {
   return GetMessage("ℹ "+sMessage, hcolors.Yellow, "INFO   ")
}

func GetProgressMessage(sMessage string) string {
   return GetMessage("⏳ "+sMessage, hcolors.Cyan, "WORKING")
}

func GetToastBanner(sTitle string, sDetails string, isSuccess bool) string {
   var colorCode string = hcolors.Green
   var symbol string = "✔"
   if !isSuccess {
      colorCode = hcolors.Red
      symbol = "✖"
   }

   var border string = strings.Repeat("─", 60)
   var result string = fmt.Sprintf("┌%s┐\r\n", border)
   result += fmt.Sprintf("│ %s%-58s%s │\r\n", colorCode, symbol+"  "+sTitle, hcolors.Reset)
   if !hstrings.IsEmpty(sDetails) {
      result += fmt.Sprintf("│ %-58s │\r\n", sDetails)
   }
   result += fmt.Sprintf("└%s┘", border)
   return result
}

func GetMessage(sMessage string, sExtraParams ...string) string {
   var sFullMessage string = ""

   if len(sExtraParams) > 0 {
      if len(sExtraParams) > 1 {
         var sTypeMessage string = ""

         if !hstrings.IsEmpty(sExtraParams[0]) {
            sTypeMessage = "[" + sExtraParams[0] + " " + sExtraParams[1] + " " + hcolors.Reset + "]"
         } else {
            sTypeMessage = "[" + sExtraParams[1] + "]"
         }

         sFullMessage += sTypeMessage + " "
      }
   }

   if len(sExtraParams) > 0 && !hstrings.IsEmpty(sExtraParams[0]) {
      sFullMessage += sExtraParams[0] + sMessage + hcolors.Reset
   } else {
      sFullMessage += sMessage
   }

   return sFullMessage
}
