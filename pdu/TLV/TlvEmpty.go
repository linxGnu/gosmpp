package TLV

import (
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/Utils"
)

type TLVEmpty struct {
	TLV
	Present bool
}

func NewTLVEmpty() *TLVEmpty {
	a := &TLVEmpty{}
	a.Construct()

	return a
}

func NewTLVEmptyWithTag(tag int16) *TLVEmpty {
	a := NewTLVEmpty()
	a.Tag = tag

	return a
}

func NewTLVEmptyWithTagValue(tag int16, present bool) *TLVEmpty {
	a := NewTLVEmptyWithTag(tag)
	a.Present = present
	a.MarkValueSet()

	return a
}

func (c *TLVEmpty) Construct() {
	c.TLV.Construct()
	c.SetRealReference(c)

	c.MinLength = 0
	c.MaxLength = 0
}

func (c *TLVEmpty) GetValueData() (b *Utils.ByteBuffer, er *Exception.Exception) {
	return nil, nil
}

func (c *TLVEmpty) SetValueData(buffer *Utils.ByteBuffer) *Exception.Exception {
	if !c.CheckLengthBuffer(buffer) {
		return Exception.NewExceptionFromStr("TLVEmpty: Buffer length is not valid")
	}

	c.SetValue(true)

	return nil
}

func (c *TLVEmpty) SetValue(value bool) *Exception.Exception {
	c.Present = value
	c.MarkValueSet()

	return nil
}

func (c *TLVEmpty) GetValue() (bool, *Exception.Exception) {
	if c.HasValue() {
		return c.Present, nil
	}

	return false, Exception.ValueNotSetException
}
