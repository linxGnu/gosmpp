package Utils

import "unicode/utf8"

func GetStringLength(st string) int {
	return utf8.RuneCountInString(st)
}

func Substring(s string, len int) string {
	if len == 0 {
		return ""
	}

	by := []byte(s)

	if int(by[len-1]) >= 224 {
		len += 2
	} else if int(by[len-1]) >= 192 && int(by[len-1]) < 224 {
		len++
	}

	return s[0:len]
}
