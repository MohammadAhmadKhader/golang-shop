package payloads

import "strings"

// use this function when you want to trim a string that you not sure if its nil or string
func TrimStrPtr(strPtr *string) *string {
	var trimmedStr *string
    if strPtr != nil {
		desc := strings.Trim(*strPtr, " ")
		trimmedStr = &desc
        return trimmedStr
	}
    
    return trimmedStr
}