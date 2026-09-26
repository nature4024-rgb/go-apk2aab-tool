package hmessages

import (
   "apk2aab/helpers/hcolors"
   "apk2aab/helpers/hstrings"
   "fmt"
)

func GetSuccessMessage(sMessage string) string {
   if hstrings.IsEmpty(sMessage) {
      return GetMessage("✔ Done", hcolors.Green, "SUCCESS")
   }
   return GetMessage(sMessage, hcolors.Green, "SUCCESS")
}

func GetErrorMessage(sMessage string) string {
   if hstrings.IsEmpty(sMessage) {
      return GetMessage("✖ Failed", hcolors.Red, "ERROR  ")
   }
   return GetMessage(sMessage, hcolors.Red, "ERROR  ")
}

func GetInfoMessage(sMessage string) string {
   return GetMessage(sMessage, hcolors.Yellow, "INFO   ")
}

func GetWarningMessage(sMessage string) string {
   return GetMessage(sMessage, hcolors.Yellow, "WARN   ")
}

func GetStepMessage(iStep int, iTotalSteps int, sMessage string) string {
   sStepTag := fmt.Sprintf("STEP %d/%d", iStep, iTotalSteps)
   return GetMessage(sMessage, hcolors.Cyan, sStepTag)
}

func GetMessage(sMessage string, sExtraParams ...string) string {
   var sFullMessage string = ""

   if len(sExtraParams) > 0 {
      if len(sExtraParams) > 1 {
         var sTypeMessage string = ""

         if !hstrings.IsEmpty(sExtraParams[0]) {
            sTypeMessage = "[ " + sExtraParams[0] + sExtraParams[1] + hcolors.Reset + " ]"
         } else {
            sTypeMessage = "[ " + sExtraParams[1] + " ]"
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
