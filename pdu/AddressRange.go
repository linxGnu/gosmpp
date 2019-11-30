package PDU

import (
	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/PDU/Common"
	"github.com/linxGnu/gosmpp/Utils"
)

type AddressRange struct {
	Common.ByteData
	Ton          byte
	Npi          byte
	AddressRange string
}

func NewAddressRange() *AddressRange {
	a := &AddressRange{}
	a.Construct()

	return a
}

func NewAddressRangeWithAddr(addr string) (*AddressRange, error) {
	a := NewAddressRange()
	a.AddressRange = addr

	return a, nil
}

func NewAddressRangeWithTonNpiAddr(ton, npi byte, addr string) (*AddressRange, error) {
	a, err := NewAddressRangeWithAddr(addr)
	if err != nil {
		return nil, err
	}

	a.Ton = ton
	a.Npi = npi

	return a, nil
}

func (c *AddressRange) Construct() {
	defer c.SetRealReference(c)
	c.ByteData.Construct()

	c.Ton = Data.DFLT_GSM_TON
	c.Npi = Data.DFLT_GSM_NPI
	c.AddressRange = Data.DFLT_ADDR_RANGE
}

func (c *AddressRange) SetAddressRange(addr string) *Exception.Exception {
	err := c.CheckCStringMax(addr, int(Data.SM_ADDR_RANGE_LEN))
	if err != nil {
		return err
	}

	c.AddressRange = addr
	return nil
}

func (c *AddressRange) GetAddressRangeWithEncoding(enc Data.Encoding) (string, *Exception.Exception) {
	bytes, err := enc.Encode(c.AddressRange)
	if err != nil {
		return "", Exception.UnsupportedEncodingException
	}

	res, err := enc.Decode(bytes)
	if err != nil {
		return "", Exception.UnsupportedEncodingException
	}

	return res, nil
}

func (c *AddressRange) SetData(bb *Utils.ByteBuffer) (err *Exception.Exception) {
	if bb == nil || bb.Buffer == nil {
		return Exception.NewExceptionFromStr("AddressRange: buffer is nil")
	}

	if bb.Len() < Utils.SZ_BYTE<<1 {
		err = Exception.NewExceptionFromStr("Address: buffer is not enough for Address")
		return
	}

	c.Ton = bb.Read_UnsafeByte()
	c.Npi = bb.Read_UnsafeByte()
	c.AddressRange, err = bb.Read_CString()

	return
}

func (c *AddressRange) GetData() (result *Utils.ByteBuffer, err *Exception.Exception) {
	bb := Utils.NewBuffer(make([]byte, 0, len(c.AddressRange)<<1+1+(Utils.SZ_BYTE<<1)))

	bb.Write_UnsafeByte(c.Ton)
	bb.Write_UnsafeByte(c.Npi)

	return bb, bb.Write_CString(c.AddressRange)
}

func (c *AddressRange) SetTon(data byte) {
	c.Ton = data
}

func (c *AddressRange) SetNpi(data byte) {
	c.Npi = data
}

func (c *AddressRange) GetTon() byte {
	return c.Ton
}

func (c *AddressRange) GetNpi() byte {
	return c.Npi
}
