package pdu

import (
	"fmt"

	"github.com/linxGnu/gosmpp/data"
)

// Address smpp address of src and dst.
type Address struct {
	ton     byte
	npi     byte
	address string
}

// NewAddress returns new address with default max length.
func NewAddress() Address {
	return Address{ton: data.GetDefaultTon(), npi: data.GetDefaultNpi()}
}

// NewAddressWithAddr returns new address.
func NewAddressWithAddr(addr string) (a Address, err error) {
	a = NewAddress()
	err = a.SetAddress(addr)
	return
}

// NewAddressWithTonNpi returns new address with ton, npi.
func NewAddressWithTonNpi(ton, npi byte) Address {
	return Address{ton: ton, npi: npi}
}

// NewAddressWithTonNpiAddr returns new address with ton, npi, addr string.
func NewAddressWithTonNpiAddr(ton, npi byte, addr string) (a Address, err error) {
	a = NewAddressWithTonNpi(ton, npi)
	err = a.SetAddress(addr)
	return
}

// Unmarshal from buffer.
func (c *Address) Unmarshal(b *ByteBuffer) (err error) {
	if c.ton, err = b.ReadByte(); err == nil {
		if c.npi, err = b.ReadByte(); err == nil {
			c.address, err = b.ReadCString()
		}
	}
	return
}

// Marshal to buffer.
func (c *Address) Marshal(b *ByteBuffer) {
	b.Grow(3 + len(c.address))

	_ = b.WriteByte(c.ton)
	_ = b.WriteByte(c.npi)
	_ = b.WriteCString(c.address)
}

// SetTon sets ton.
func (c *Address) SetTon(ton byte) {
	c.ton = ton
}

// SetNpi sets npi.
func (c *Address) SetNpi(npi byte) {
	c.npi = npi
}

// SetAddress sets address.
func (c *Address) SetAddress(addr string) (err error) {
	if len(addr) > data.SM_ADDR_LEN {
		err = fmt.Errorf("Address len exceed limit. (%d > %d)", len(addr), data.SM_ADDR_LEN)
	} else {
		c.address = addr
	}
	return
}

// Ton returns assigned ton.
func (c Address) Ton() byte {
	return c.ton
}

// Npi returns assigned npi.
func (c Address) Npi() byte {
	return c.npi
}

// Address returns assigned address (in string).
func (c Address) Address() string {
	return c.address
}

// String returns assigned address
func (c Address) String() string {
	return c.Address()
}
