package TLV

import (
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/Utils"
)

type TLVUByte struct {
	TLV
	Value uint8
}

func NewTLVUByte() *TLVUByte {
	a := &TLVUByte{}
	a.Construct()

	return a
}

func NewTLVUByteWithTag(tag int16) *TLVUByte {
	a := NewTLVUByte()
	a.Tag = tag

	return a
}

func NewTLVUByteWithTagValue(tag int16, value uint8) *TLVUByte {
	a := NewTLVUByteWithTag(tag)
	a.SetValue(value)

	return a
}

func (c *TLVUByte) Construct() {
	c.TLV.Construct()
	c.SetRealReference(c)

	c.MinLength = 1
	c.MaxLength = 1
}

func (c *TLVUByte) GetValueData() (b *Utils.ByteBuffer, er *Exception.Exception) {
	val, er := c.GetValue()
	if er != nil {
		return nil, er
	}

	buf := Utils.NewBuffer(make([]byte, 0, 1))
	return buf, buf.Write_Byte(val)
}

func (c *TLVUByte) SetValueData(buffer *Utils.ByteBuffer) *Exception.Exception {
	if !c.CheckLengthBuffer(buffer) {
		return Exception.NotEnoughDataInByteBufferException
	}

	val, err := buffer.Read_Byte()
	if err != nil {
		return err
	}

	c.SetValue(val)

	return nil
}

func (c *TLVUByte) SetValue(value uint8) *Exception.Exception {
	c.Value = value
	c.MarkValueSet()

	return nil
}

func (c *TLVUByte) GetValue() (uint8, *Exception.Exception) {
	if c.HasValue() {
		return c.Value, nil
	}

	return 0, Exception.ValueNotSetException
}
