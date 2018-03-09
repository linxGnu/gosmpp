package Utils

import "time"

const (
	maxMonthRange = 12 * 1000
)

// MonthBetween total months between time point
func MonthBetween(t1, t2 time.Time) int {
	if t1.Equal(t2) || t1.After(t2) {
		return 0
	}

	left, right, mid := 0, maxMonthRange, 0
	for left < right {
		if mid = (left + right) >> 1; mid == left {
			return left
		} else if t1.AddDate(0, mid, 0).After(t2) {
			right = mid
		} else {
			left = mid
		}
	}

	return left
}
