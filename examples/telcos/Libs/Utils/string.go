package Utils

import (
	"errors"
	"regexp"
)

const (
	SOURCE_CHARACTERS      = "ÀÁÂÃÈÉÊÌÍÒÓÔÕÙÚÝàáâãèéêìíòóôõùúýĂăĐđĨĩŨũƠơƯưẠạẢảẤấẦầẨẩẪẫẬậẮắẰằẲẳẴẵẶặẸẹẺẻẼẽẾếỀềỂểỄễỆệỈỉỊịỌọỎỏỐốỒồỔổỖỗỘộỚớỜờỞởỠỡỢợỤụỦủỨứỪừỬửỮữỰự"
	DESTINATION_CHARACTERS = "AAAAEEEIIOOOOUUYaaaaeeeiioooouuyAaDdIiUuOoUuAaAaAaAaAaAaAaAaAaAaAaAaEeEeEeEeEeEeEeEeIiIiOoOoOoOoOoOoOoOoOoOoOoOoUuUuUuUuUuUuUu"
)

//
var source_rune = []rune(SOURCE_CHARACTERS)
var des_rune = []rune(DESTINATION_CHARACTERS)

//
func searchRune(v rune) int {
	left, right := 0, len(source_rune)-1
	if source_rune[right] < v || v < source_rune[left] {
		return -1
	}

	mid := 0
	for left <= right {
		if mid = (left + right) >> 1; mid == left {
			if source_rune[left] == v {
				return left
			} else if source_rune[right] == v {
				return right
			}
			return -1
		}

		if source_rune[mid] == v {
			return mid
		} else if source_rune[mid] < v {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return -1
}

// TransformVietnamese into "khong dau"
func TransformVietnamese(inp string) string {
	tmp := []rune(inp)

	ind := 0
	for i := range tmp {
		if ind = searchRune(tmp[i]); ind != -1 {
			tmp[i] = des_rune[ind]
		}
	}

	return string(tmp)
}

var reg, _ = regexp.Compile("[^a-zA-Z0-9]+")

// StripNonAlphabet ...
func StripNonAlphabet(input string) string {
	return reg.ReplaceAllString(input, "")
}

// IsUnicode check is unicode string
func IsUnicode(input string) bool {
	return len(input) != len([]rune(input))
}

var (
	ErrEmailBadFormat = errors.New("invalid format")
	emailRegexp       = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// ValidateEmailFormat validate email format
func ValidateEmailFormat(email string) error {
	if !emailRegexp.MatchString(email) {
		return ErrEmailBadFormat
	}
	return nil
}
