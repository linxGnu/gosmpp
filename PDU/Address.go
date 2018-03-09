package PDU

import (
	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/PDU/Common"
	"github.com/linxGnu/gosmpp/Utils"
)

var defaultMaxAddressLength = Data.SM_ADDR_LEN

// Address smpp address of src and dst
type Address struct {
	Common.ByteData
	Ton              byte
	Npi              byte
	Address          string
	MaxAddressLength int32
}

// NewAddress create new address with default max length
func NewAddress() *Address {
	res := &Address{}
	res.Construct()

	return res
}

// NewAddressWithAddr create new address
func NewAddressWithAddr(addr string) (*Address, *Exception.Exception) {
	res := NewAddress()

	err := res.SetAddress(addr)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// NewAddressWithMaxLength create new address, set max length in C of address
func NewAddressWithMaxLength(len int32) *Address {
	addr := NewAddress()
	addr.MaxAddressLength = len

	return addr
}

func NewAddressWithTonNpiLen(ton, npi byte, len int32) *Address {
	addr := NewAddress()
	addr.Ton = ton
	addr.Npi = npi
	addr.MaxAddressLength = len

	return addr
}

func NewAddressWithTonNpiAddr(ton, npi byte, addr string) (*Address, *Exception.Exception) {
	a, err := NewAddressWithAddr(addr)
	if err != nil {
		return nil, err
	}
	a.Ton = ton
	a.Npi = npi

	return a, nil
}

func NewAddressWithTonNpiAddrMaxLen(ton, npi byte, addr string, len int32) (*Address, *Exception.Exception) {
	a := NewAddressWithTonNpiLen(ton, npi, len)
	err := a.SetAddress(addr)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (c *Address) Construct() {
	defer c.SetRealReference(c)
	c.ByteData.Construct()

	c.Ton = Data.DFLT_GSM_TON
	c.Npi = Data.DFLT_GSM_NPI
	c.Address = Data.DFLT_ADDR
	c.MaxAddressLength = defaultMaxAddressLength
}

func (c *Address) SetData(bb *Utils.ByteBuffer) (err *Exception.Exception) {
	if bb == nil || bb.Buffer == nil {
		err = Exception.NewExceptionFromStr("Address: buffer is nil")
		return
	}

	if bb.Len() < Utils.SZ_BYTE<<1 {
		err = Exception.NewExceptionFromStr("Address: buffer is not enough for Address")
		return
	}

	c.Ton = bb.Read_UnsafeByte()
	c.Npi = bb.Read_UnsafeByte()
	c.Address, err = bb.Read_CString()

	return
}

func (c *Address) GetData() (result *Utils.ByteBuffer, err *Exception.Exception) {
	bb := Utils.NewBuffer(make([]byte, 0, len(c.Address)<<1+1+(Utils.SZ_BYTE<<1)))

	bb.Write_UnsafeByte(c.Ton)
	bb.Write_UnsafeByte(c.Npi)

	return bb, bb.Write_CString(c.Address)
}

func (c *Address) SetAddress(addr string) *Exception.Exception {
	err := c.CheckCStringMax(addr, int(c.MaxAddressLength))
	if err != nil {
		return err
	}

	c.Address = addr
	return nil
}

func (c *Address) GetAddress() string {
	return c.Address
}

func (c *Address) GetAddressWithEncoding(enc Data.Encoding) (string, *Exception.Exception) {
	bytes, err := enc.Encode(c.Address)
	if err != nil {
		return "", Exception.UnsupportedEncodingException
	}

	res, err := enc.Decode(bytes)
	if err != nil {
		return "", Exception.UnsupportedEncodingException
	}

	return res, nil
}

func (c *Address) GetTon() byte {
	return c.Ton
}

func (c *Address) GetNpi() byte {
	return c.Npi
}

func (c *Address) SetTon(data byte) {
	c.Ton = data
}

func (c *Address) SetNpi(data byte) {
	c.Npi = data
}
