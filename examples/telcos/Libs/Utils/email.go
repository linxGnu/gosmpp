package Utils

import (
	"errors"
	"net"
	"regexp"
	"strings"
)

/**
 * Code below forged from: https://github.com/goware/emailx
 */
var (
	ErrInvalidFormat    = errors.New("invalid format")
	ErrUnresolvableHost = errors.New("unresolvable host")

	userRegexp = regexp.MustCompile("^[a-zA-Z0-9!#$%&'*+/=?^_`{|}~.-]+$")
	hostRegexp = regexp.MustCompile("^[^\\s]+\\.[^\\s]+$")

	// As per RFC 5332 secion 3.2.3: https://tools.ietf.org/html/rfc5322#section-3.2.3
	// Dots are not allowed in the beginning, end or in occurances of more than 1 in the email address
	userDotRegexp = regexp.MustCompile("(^[.]{1})|([.]{1}$)|([.]{2,})")
)

// Validate checks format of a given email and resolves its host name.
func Validate(email string) error {
	host, err := ValidateFast(email)
	if err != nil {
		return err
	}

	if _, err := net.LookupMX(host); err != nil {
		if _, err := net.LookupIP(host); err != nil {
			// Only fail if both MX and A records are missing - any of the
			// two is enough for an email to be deliverable
			return ErrUnresolvableHost
		}
	}

	return nil
}

// ValidateFast checks format of a given email.
func ValidateFast(email string) (host string, err error) {
	if len(email) < 6 || len(email) > 254 {
		err = ErrInvalidFormat
		return
	}

	at := strings.LastIndex(email, "@")
	if at <= 0 || at > len(email)-3 {
		err = ErrInvalidFormat
		return
	}

	user := email[:at]
	host = email[at+1:]

	if len(user) > 64 {
		err = ErrInvalidFormat
		return
	}

	if userDotRegexp.MatchString(user) || !userRegexp.MatchString(user) || !hostRegexp.MatchString(host) {
		err = ErrInvalidFormat
		return
	}

	err = nil
	return
}

// Normalize normalizes email address.
func Normalize(email string) string {
	// Trim whitespaces.
	email = strings.TrimSpace(email)

	// Trim extra dot in hostname.
	email = strings.TrimRight(email, ".")

	// Lowercase.
	email = strings.ToLower(email)

	return email
}
