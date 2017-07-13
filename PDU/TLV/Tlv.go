package TLV

import (
	"fmt"

	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/PDU/Common"
	"github.com/linxGnu/gosmpp/Utils"
)

var DONT_CHECK_LIMIT = -1

type ITLV interface {
	GetTag() int16
	SetTag(tag int16)
	SetData(bb *Utils.ByteBuffer) (err *Exception.Exception)
	GetData() (*Utils.ByteBuffer, *Exception.Exception)
	HasValue() bool
	GetValueData() (*Utils.ByteBuffer, *Exception.Exception)
	SetValueData(buffer *Utils.ByteBuffer) *Exception.Exception
}

type TLV struct {
	Common.ByteData
	Tag        int16
	ValueIsSet bool
	MinLength  int
	MaxLength  int
}

func NewTLV() *TLV {
	a := &TLV{}
	a.Construct()

	return a
}

func NewTLVWithTag(tag int16) *TLV {
	a := NewTLV()
	a.SetTag(tag)

	return a
}

func NewTLVWithTagAndLength(tag int16, minLenght, maxLength int) *TLV {
	a := NewTLVWithTag(tag)
	a.MinLength = minLenght
	a.MaxLength = maxLength

	return a
}

func (c *TLV) Construct() {
	c.ByteData.Construct()
	c.SetRealReference(c)

	c.MinLength = DONT_CHECK_LIMIT
	c.MaxLength = DONT_CHECK_LIMIT
	c.Tag = 0
	c.ValueIsSet = false
}

func (c *TLV) GetTag() int16 {
	return c.Tag
}

func (c *TLV) SetTag(tag int16) {
	c.Tag = tag
}

func (c *TLV) GetValueData() (*Utils.ByteBuffer, *Exception.Exception) {
	return nil, nil
}

func (c *TLV) SetValueData(buffer *Utils.ByteBuffer) *Exception.Exception {
	return nil
}

func (c *TLV) HasValue() bool {
	return c.ValueIsSet
}

func (c *TLV) GetLength() (len int, err *Exception.Exception) {
	if c.HasValue() {
		valueBuf, er := c.GetValueData()
		if er != nil {
			return 0, er
		}

		if valueBuf == nil {
			return 0, nil
		}

		return valueBuf.Len(), nil
	}

	return 0, nil
}

func (c *TLV) SetData(bb *Utils.ByteBuffer) (err *Exception.Exception) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
		}
	}()

	if bb.Len() < Utils.SZ_SHORT<<1 {
		return Exception.NewExceptionFromStr("TLV: Buffer not enough to SetData")
	}

	c.Tag = bb.Read_UnsafeShort()
	length := bb.Read_UnsafeShort()

	tmp, err := bb.Read_Bytes(int(length))
	if err != nil {
		return err
	}

	return c.SetValueData(tmp)
}

func (c *TLV) GetData() (*Utils.ByteBuffer, *Exception.Exception) {
	if c.HasValue() {
		count, err := c.GetLength()
		if err != nil {
			return nil, err
		}

		valueData, err := c.GetValueData()
		if err != nil {
			return nil, err
		}

		tmp := Utils.NewBuffer(make([]byte, 0, Utils.SZ_SHORT*2+valueData.Len()))
		tmp.Write_UnsafeShort(c.Tag)
		tmp.Write_UnsafeShort(Common.EncodeUnsignedFromInt(count))
		err = tmp.Write_Buffer(valueData)

		return tmp, err
	}

	return nil, nil
}

func (c *TLV) MarkValueSet() {
	c.ValueIsSet = true
}

func (c *TLV) CheckLength(length int) bool {
	min := 0
	max := 0
	if c.MinLength != DONT_CHECK_LIMIT {
		min = c.MinLength
	} else {
		min = 0
	}

	if c.MaxLength != DONT_CHECK_LIMIT {
		max = c.MaxLength
	} else {
		max = Common.MaxInt
	}

	return min <= length && length <= max
}

func (c *TLV) CheckLengthBuffer(buffer *Utils.ByteBuffer) bool {
	var len int
	if buffer == nil {
		len = 0
	} else {
		len = buffer.Len()
	}

	return c.CheckLength(len)
}
