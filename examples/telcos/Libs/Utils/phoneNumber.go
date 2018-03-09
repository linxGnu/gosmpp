package Utils

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var r, _ = regexp.Compile("^(0|84)(9\\d|16[2-9]|12\\d|86|88|89|186|188|199)(\\d{7})$")
var viettel, _ = regexp.Compile("^(0|84)(9[6-8]|16[2-9]|86)(\\d{7})$")
var vinaphone, _ = regexp.Compile("^(0|84)(9[14]|12[34579]|88)(\\d{7})$")
var mobifone, _ = regexp.Compile("^(0|84)(9[03]|12[01268]|89)(\\d{7})$")
var vietnamobile, _ = regexp.Compile("^(0|84)(92|186|188)(\\d{7})$")
var gmobile, _ = regexp.Compile("^(0|84)(99|199)(\\d{7})$")
var landline, _ = regexp.Compile("^(0|84)\\d{3}\\d{7}$")

// ValidateVNPhoneNumber ...
func ValidateVNPhoneNumber(phone string) bool {
	return r.Match([]byte(phone)) || landline.Match([]byte(phone))
}

// StandardizePhone ...
func StandardizePhone(phone string) string {
	phone = strings.TrimPrefix(phone, "+")
	phone = strings.TrimPrefix(phone, "00")

	if !strings.HasPrefix(phone, "84") {
		if utf8.RuneCountInString(phone) <= 10 && !strings.HasPrefix(phone, "0") {
			phone = "84" + phone
		} else {
			if strings.HasPrefix(phone, "0") {
				phone = strings.TrimPrefix(phone, "0")
				phone = "84" + phone
			}
		}
	}

	return phone
}

// GetVNTelcoByPhone ...
func GetVNTelcoByPhone(phone string) string {
	p := []byte(phone)

	if viettel.Match(p) {
		return "viettel"
	}

	if mobifone.Match(p) {
		return "mobifone"
	}

	if vinaphone.Match(p) {
		return "vinaphone"
	}

	if vietnamobile.Match(p) {
		return "vietnamobile"
	}

	if gmobile.Match(p) {
		return "gmobile"
	}

	if landline.Match(p) {
		return "landline"
	}

	return "other"
}
