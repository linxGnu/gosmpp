package Common

import (
	"errors"
	"fmt"

	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/Utils"
)

type IByteDataList interface {
	CreateValue() IByteData
}

type ByteDataList struct {
	ByteData
	Values       []IByteData
	MaxSize      int
	LengthOfSize byte
}

func NewByteDataList() *ByteDataList {
	a := &ByteDataList{}
	a.Construct()

	return a
}

func NewByteDataListWithSize(max int, lengthOfSize int) (*ByteDataList, error) {
	a := NewByteDataList()

	a.MaxSize = max
	if lengthOfSize != Utils.SZ_BYTE && lengthOfSize != Utils.SZ_SHORT && lengthOfSize != Utils.SZ_INT {
		return nil, errors.New("ByteDataList: constructor with size length not valid")
	}
	a.LengthOfSize = byte(lengthOfSize)

	return a, nil
}

func (c *ByteDataList) Construct() {
	c.ByteData.Construct()
	c.SetRealReference(c)

	c.Values = make([]IByteData, 0)
}

func (c *ByteDataList) ResetValues() {
	c.Values = make([]IByteData, 0)
}

func (c *ByteDataList) SetData(buffer *Utils.ByteBuffer) (err *Exception.Exception) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
		}
	}()

	c.ResetValues()

	var nrValues int
	switch int(c.LengthOfSize) {
	case Utils.SZ_BYTE:
		v, err := buffer.Read_Byte()
		if err != nil {
			return err
		}

		nrValues = int(DecodeUnsigned(v))
	case Utils.SZ_SHORT:
		v, err := buffer.Read_Short()
		if err != nil {
			return err
		}

		nrValues = DecodeUnsignedFromInt16(v)
	case Utils.SZ_INT:
		v, err := buffer.Read_Int()
		if err != nil {
			return err
		}

		nrValues = int(v)
	}

	test := c.This.(IByteDataList).CreateValue()
	if test == nil {
		return nil
	}

	c.Values = make([]IByteData, nrValues)
	for i := 0; i < nrValues; i++ {
		c.Values[i] = c.This.(IByteDataList).CreateValue()
		err := c.Values[i].SetData(buffer)
		if err != nil {
			c.ResetValues()
			return err
		}
	}

	return nil
}

func (c *ByteDataList) GetCount() int {
	if c.Values == nil {
		c.Values = make([]IByteData, 0)
	}

	return len(c.Values)
}

func (c *ByteDataList) GetData() (*Utils.ByteBuffer, *Exception.Exception) {
	buf := Utils.NewBuffer([]byte{})

	numberValues := c.GetCount()
	switch int(c.LengthOfSize) {
	case Utils.SZ_BYTE:
		err := buf.Write_Byte(EncodeUnsigned(int16(numberValues)))
		buf.Grow(numberValues * Utils.SZ_BYTE)
		if err != nil {
			return nil, err
		}
	case Utils.SZ_SHORT:
		err := buf.Write_Short(EncodeUnsignedFromInt(int(numberValues)))
		buf.Grow(numberValues * Utils.SZ_SHORT)
		if err != nil {
			return nil, err
		}
	case Utils.SZ_INT:
		err := buf.Write_Int(int32(numberValues))
		buf.Grow(numberValues * Utils.SZ_INT)
		if err != nil {
			return nil, err
		}
	}

	for _, val := range c.Values {
		if val == nil {
			return nil, Exception.ValueNotSetException
		}

		tmp, err := val.GetData()
		if err != nil {
			return nil, err
		}

		err = buf.Write_Buffer(tmp)
		if err != nil {
			return nil, err
		}
	}

	return buf, nil
}

func (c *ByteDataList) AddValue(val IByteData) *Exception.Exception {
	if c.GetCount() > c.MaxSize {
		return Exception.TooManyValuesException
	}

	if val != nil {
		c.Values = append(c.Values, val)
	}

	return nil
}

func (c *ByteDataList) GetValue(index int) IByteData {
	if index < c.GetCount() {
		return c.Values[index]
	}

	return nil
}
