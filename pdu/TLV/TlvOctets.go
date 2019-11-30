package TLV

import (
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/Utils"
)

type TLVOctets struct {
	TLV
	Value *Utils.ByteBuffer
}

func NewTLVOctets() *TLVOctets {
	a := &TLVOctets{}
	a.Construct()

	return a
}

func NewTLVOctetsWithTag(tag int16) *TLVOctets {
	a := NewTLVOctets()
	a.Tag = tag

	return a
}

func NewTLVOctetsWithTagLength(tag int16, min, max int) *TLVOctets {
	a := NewTLVOctetsWithTag(tag)
	a.MinLength = min
	a.MaxLength = max

	return a
}

func NewTLVOctetsWithTagValue(tag int16, buf *Utils.ByteBuffer) *TLVOctets {
	a := NewTLVOctetsWithTag(tag)
	a.SetValue(buf)

	return a
}

func NewTLVOctetsWithTagLengthValue(tag int16, min, max int, buf *Utils.ByteBuffer) *TLVOctets {
	a := NewTLVOctetsWithTagLength(tag, min, max)
	a.SetValue(buf)

	return a
}

func (c *TLVOctets) Construct() {
	c.TLV.Construct()
	c.SetRealReference(c)
}

func (c *TLVOctets) SetValueData(buffer *Utils.ByteBuffer) *Exception.Exception {
	if !c.CheckLengthBuffer(buffer) {
		return Exception.NotEnoughDataInByteBufferException
	}

	return c.SetValue(buffer)
}

func (c *TLVOctets) GetValueData() (*Utils.ByteBuffer, *Exception.Exception) {
	val, err := c.GetValue()
	if err != nil {
		return nil, err
	}

	buf := Utils.NewBuffer(make([]byte, 0, 16))
	return buf, buf.Write_Buffer(val)
}

func (c *TLVOctets) SetValue(buffer *Utils.ByteBuffer) *Exception.Exception {
	if c.Value != nil && c.Value.Buffer != nil {
		c.Value = buffer
		c.MarkValueSet()
	} else {
		c.Value = nil
	}

	return nil
}

func (c *TLVOctets) GetValue() (*Utils.ByteBuffer, *Exception.Exception) {
	if c.Value != nil && c.HasValue() {
		return c.Value, nil
	}

	return nil, Exception.ValueNotSetException
}
