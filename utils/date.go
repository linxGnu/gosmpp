package utils

import (
	"strconv"
	"time"

	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/errors"
)

// CheckDate checks date string format.
func CheckDate(dateStr string) (err error) {
	strLen := len(dateStr)

	count := strLen + 1
	if count != 1 && count != int(data.SM_DATE_LEN) {
		err = errors.ErrWrongDateFormat
		return
	}

	if count == 1 {
		return
	}

	locTime := string(dateStr[strLen-1])
	if locTime != "-" && locTime != "+" && locTime != "R" {
		err = errors.ErrWrongDateFormat
		return
	}

	formatLen := len("060102150405")
	dateGoStr := dateStr[0:formatLen]

	if _, err = time.Parse("060102150405", dateGoStr); err != nil {
		return
	}

	tenthsOfSecStr := dateStr[formatLen : formatLen+1]
	if _, err = strconv.Atoi(tenthsOfSecStr); err != nil {
		return
	}

	timeDiffStr := dateStr[formatLen+1 : formatLen+3]
	timeDiff, err := strconv.Atoi(timeDiffStr)
	if err != nil {
		return
	}

	if timeDiff < 0 || timeDiff > 48 {
		err = errors.ErrWrongDateFormat
	}

	return
}
