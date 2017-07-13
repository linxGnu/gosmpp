package TLV

import (
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/Utils"
)

type TLVString struct {
	TLV
	Value string
}

func NewTLVString() *TLVString {
	a := &TLVString{}
	a.Construct()

	return a
}

func NewTLVStringWithTag(tag int16) *TLVString {
	a := NewTLVString()
	a.Tag = tag

	return a
}

func NewTLVStringWithTagValue(tag int16, value string) *TLVString {
	a := NewTLVStringWithTag(tag)
	a.SetValue(value)

	return a
}

func NewTLVStringWithTagLength(tag int16, min, max int) *TLVString {
	a := NewTLVStringWithTag(tag)
	a.MinLength = min
	a.MaxLength = max

	return a
}

func NewTLVStringWithTagLengthValue(tag int16, min, max int, value string) *TLVString {
	a := NewTLVStringWithTagLength(tag, min, max)
	a.SetValue(value)

	return a
}

func (c *TLVString) Construct() {
	c.TLV.Construct()
	c.SetRealReference(c)
}

func (c *TLVString) GetValueData() (b *Utils.ByteBuffer, er *Exception.Exception) {
	val, er := c.GetValue()
	if er != nil {
		return nil, er
	}

	buf := Utils.NewBuffer(make([]byte, 0, len(val)<<1+1))
	return buf, buf.Write_CString(val)
}

func (c *TLVString) SetValueData(buffer *Utils.ByteBuffer) *Exception.Exception {
	if !c.CheckLengthBuffer(buffer) {
		return Exception.WrongLengthException
	}

	val, err := buffer.Read_CString()
	if err != nil {
		return err
	}

	c.SetValue(val)

	return nil
}

func (c *TLVString) SetValue(value string) *Exception.Exception {
	if !c.CheckLength(len(value) + 1) {
		return Exception.WrongLengthOfStringException
	}

	c.Value = value
	c.MarkValueSet()

	return nil
}

func (c *TLVString) GetValue() (string, *Exception.Exception) {
	if c.HasValue() {
		return c.Value, nil
	}

	return "", Exception.ValueNotSetException
}
