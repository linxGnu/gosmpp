package utils

import "unicode/utf8"

// GetStringLength returns string length.
func GetStringLength(st string) int {
	return utf8.RuneCountInString(st)
}

// Prefix returns prefix of s with length n.
func Prefix(s string, n int) string {
	return string([]rune(s)[:n])
}
