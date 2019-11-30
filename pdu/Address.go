package pdu

import (
	"fmt"

	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/utils"
)

// Address smpp address of src and dst.
type Address struct {
	Ton              byte
	Npi              byte
	Address          string
	MaxAddressLength int
}

// NewAddress create new address with default max length
func NewAddress() *Address {
	return &Address{Ton: data.GetDefaultTon(), Npi: data.GetDefaultNpi(), MaxAddressLength: data.SM_ADDR_LEN}
}

// NewAddressWithAddr create new address
func NewAddressWithAddr(addr string) (a *Address, err error) {
	a = NewAddress()
	err = a.SetAddress(addr)
	return
}

// NewAddressWithMaxLength create new address, set max length in C of address
func NewAddressWithMaxLength(len int) (a *Address) {
	a = NewAddress()
	a.MaxAddressLength = len
	return
}

// NewAddressWithTonNpiLen create new address with ton, npi, max length.
func NewAddressWithTonNpiLen(ton, npi byte, len int) *Address {
	return &Address{MaxAddressLength: len, Ton: ton, Npi: npi}
}

// Unmarshal from buffer.
func (c *Address) Unmarshal(b *utils.ByteBuffer) (err error) {
	c.Ton, err = b.ReadByte()
	if err == nil {
		c.Npi, err = b.ReadByte()
		if err == nil {
			c.Address, err = b.ReadCString()
		}
	}
	return
}

// Marshal to buffer.
func (c *Address) Marshal(b *utils.ByteBuffer) {
	b.Grow(3 + len(c.Address))
	_ = b.WriteByte(c.Ton)
	_ = b.WriteByte(c.Npi)
	_ = b.WriteCString(c.Address)
}

// SetAddress to pdu.
func (c *Address) SetAddress(addr string) (err error) {
	if c.MaxAddressLength > 0 && len(addr) > c.MaxAddressLength {
		err = fmt.Errorf("Address len exceed limit. (%d > %d)", len(addr), c.MaxAddressLength)
	}
	c.Address = addr
	return
}
