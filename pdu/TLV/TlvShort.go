package TLV

import (
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/Utils"
)

type TLVShort struct {
	TLV
	Value int16
}

func NewTLVShort() *TLVShort {
	a := &TLVShort{}
	a.Construct()

	return a
}

func NewTLVShortWithTag(tag int16) *TLVShort {
	a := NewTLVShort()
	a.Tag = tag

	return a
}

func NewTLVShortWithTagValue(tag int16, value int16) *TLVShort {
	a := NewTLVShortWithTag(tag)
	a.SetValue(value)

	return a
}

func (c *TLVShort) Construct() {
	c.TLV.Construct()
	c.SetRealReference(c)

	c.MinLength = 2
	c.MaxLength = 2
}

func (c *TLVShort) GetValueData() (b *Utils.ByteBuffer, er *Exception.Exception) {
	val, er := c.GetValue()
	if er != nil {
		return nil, er
	}

	buf := Utils.NewBuffer(make([]byte, 0, 2))
	return buf, buf.Write_Short(val)
}

func (c *TLVShort) SetValueData(buffer *Utils.ByteBuffer) *Exception.Exception {
	if !c.CheckLengthBuffer(buffer) {
		return Exception.NotEnoughDataInByteBufferException
	}

	val, err := buffer.Read_Short()
	if err != nil {
		return err
	}

	c.SetValue(val)

	return nil
}

func (c *TLVShort) SetValue(value int16) *Exception.Exception {
	c.Value = value
	c.MarkValueSet()

	return nil
}

func (c *TLVShort) GetValue() (int16, *Exception.Exception) {
	if c.HasValue() {
		return c.Value, nil
	}

	return 0, Exception.ValueNotSetException
}
