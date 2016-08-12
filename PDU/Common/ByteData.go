package Common

import (
	"fmt"
	"strconv"
	"time"

	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/Utils"
)

const (
	SMPP_TIME_DATE_FORMAT string = "060102150405"
	MaxUint               uint   = ^uint(0)
	MinUint               uint   = 0
	MaxInt                int    = int(MaxUint >> 1)
	MinInt                int    = -MaxInt - 1
)

type IByteData interface {
	SetData(*Utils.ByteBuffer) *Exception.Exception
	GetData() (*Utils.ByteBuffer, *Exception.Exception)
}

type ByteData struct {
	LibraryCheckDateFormat bool
	This                   interface{}
}

func (c *ByteData) Construct() {
	c.LibraryCheckDateFormat = true
	c.SetRealReference(c)
}

func (c *ByteData) SetRealReference(real interface{}) {
	c.This = real
}

func (c *ByteData) CheckStringMax(st string, max int) *Exception.Exception {
	return c.CheckStringMinMax(st, 0, max)
}

func (c *ByteData) CheckStringMaxEncoding(st string, max int, enc Data.Encoding) *Exception.Exception {
	return c.CheckStringMinMaxEncoding(st, 0, max, enc)
}

func (c *ByteData) CheckStringMinMax(String string, min int, max int) *Exception.Exception {
	count := Utils.GetStringLength(String)
	if count < min || count > max {
		return Exception.WrongLengthOfStringException
	}

	return nil
}

func (c *ByteData) CheckStringMinMaxEncoding(String string, min int, max int, enc Data.Encoding) *Exception.Exception {
	t1, err := enc.Encode(String)
	if err != nil {
		return Exception.UnsupportedEncodingException
	}

	count := len(t1)
	if count < min || count > max {
		return Exception.WrongLengthOfStringException
	}

	return nil
}

func (c *ByteData) CheckCStringMax(st string, max int) *Exception.Exception {
	return c.CheckCStringMinMax(st, 0, max)
}

func (c *ByteData) CheckCStringMinMax(String string, min int, max int) *Exception.Exception {
	count := len(String) + 1

	if count < min || count > max {
		return Exception.WrongLengthOfStringException
	}

	return nil
}

func (c *ByteData) CheckDate(dateStr string) (err *Exception.Exception) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
		}
	}()

	strLen := len(dateStr)
	count := strLen + 1
	if count != 1 && count != int(Data.SM_DATE_LEN) {
		return Exception.WrongDateFormatException
	}

	if count == 1 || !c.LibraryCheckDateFormat {
		return nil
	}

	locTime := string(dateStr[strLen-1])
	if locTime != "-" && locTime != "+" && locTime != "R" {
		return Exception.WrongDateFormatException
	}

	formatLen := len(SMPP_TIME_DATE_FORMAT)
	dateGoStr := dateStr[0:formatLen]
	_, err1 := time.Parse(SMPP_TIME_DATE_FORMAT, dateGoStr)
	if err1 != nil {
		return Exception.NewException(err1)
	}

	tenthsOfSecStr := dateStr[formatLen : formatLen+1]
	_, err1 = strconv.Atoi(tenthsOfSecStr)
	if err1 != nil {
		return Exception.NewException(err1)
	}

	timeDiffStr := dateStr[formatLen+1 : formatLen+3]
	timeDiff, err1 := strconv.Atoi(timeDiffStr)
	if err1 != nil {
		return Exception.NewException(err1)
	}

	if timeDiff < 0 || timeDiff > 48 {
		return Exception.WrongDateFormatException
	}

	return nil
}
