package handler

import (
	"strings"
)

//return true if uri pattern is supported
func IsSupportedPattern(uri string, patterns []string) bool{
	for i:=0; i < len(patterns); i++ {
		//if uri contains patterns[i]
		if strings.Contains(uri, patterns[i]) {
			return true
		}
	}
	return false
}

//return uri prefix pattern
func GetUriPattern(uri string, patterns []string) string{
	for i:=0; i < len(patterns); i++ {
		if strings.Contains(uri, patterns[i]) {
			return patterns[i]
		}
	}
	return ""
}
