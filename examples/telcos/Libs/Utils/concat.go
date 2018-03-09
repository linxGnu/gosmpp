package Utils

import (
	"bytes"
	"strconv"
)

// ConcatCopyPreAllocate ..
func ConcatCopyPreAllocate(slices [][]byte) []byte {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}
	tmp := make([]byte, totalLen)
	var i int
	for _, s := range slices {
		copy(tmp[i:], s)
		i += len(s)
	}
	return tmp
}

// ConcatArrUint64 ...
func ConcatArrUint64(r []uint64) string {
	if r == nil || len(r) == 0 {
		return ""
	}

	var buffer bytes.Buffer
	n := len(r)
	for i := range r {
		if i < n-1 {
			buffer.WriteString(strconv.FormatUint(r[i], 10) + ", ")
		} else {
			buffer.WriteString(strconv.FormatUint(r[i], 10))
		}
	}

	return buffer.String()
}

// ConcatArrUint32 ...
func ConcatArrUint32(r []uint32) string {
	if r == nil || len(r) == 0 {
		return ""
	}

	var buffer bytes.Buffer
	n := len(r)
	for i := range r {
		if i < n-1 {
			buffer.WriteString(strconv.FormatUint(uint64(r[i]), 10) + ", ")
		} else {
			buffer.WriteString(strconv.FormatUint(uint64(r[i]), 10))
		}
	}

	return buffer.String()
}
