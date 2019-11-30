package pdu

import (
	"fmt"

	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/utils"
)

// AddressRange smpp address range of src and dst.
type AddressRange struct {
	Ton              byte
	Npi              byte
	AddressRange     string
	MaxAddressLength int
}

// NewAddressRange create new AddressRange with default max length.
func NewAddressRange() *AddressRange {
	return &AddressRange{Ton: data.GetDefaultTon(), Npi: data.GetDefaultNpi(), MaxAddressLength: data.SM_ADDR_LEN}
}

// NewAddressRangeWithAddr create new AddressRange.
func NewAddressRangeWithAddr(addr string) (a *AddressRange, err error) {
	a = NewAddressRange()
	err = a.SetAddress(addr)
	return
}

// NewAddressRangeWithMaxLength create new AddressRange, set max length in C of address.
func NewAddressRangeWithMaxLength(len int) (a *AddressRange) {
	a = NewAddressRange()
	a.MaxAddressLength = len
	return
}

// NewAddressRangeWithTonNpiLen create new AddressRange with ton, npi, max length.
func NewAddressRangeWithTonNpiLen(ton, npi byte, len int) *AddressRange {
	return &AddressRange{MaxAddressLength: len, Ton: ton, Npi: npi}
}

// Unmarshal from buffer.
func (c *AddressRange) Unmarshal(b *utils.ByteBuffer) (err error) {
	c.Ton, err = b.ReadByte()
	if err == nil {
		c.Npi, err = b.ReadByte()
		if err == nil {
			c.AddressRange, err = b.ReadCString()
		}
	}
	return
}

// Marshal to buffer.
func (c *AddressRange) Marshal(b *utils.ByteBuffer) {
	b.Grow(3 + len(c.AddressRange))
	_ = b.WriteByte(c.Ton)
	_ = b.WriteByte(c.Npi)
	_ = b.WriteCString(c.AddressRange)
}

// SetAddress to pdu.
func (c *AddressRange) SetAddress(addr string) (err error) {
	if c.MaxAddressLength > 0 && len(addr) > c.MaxAddressLength {
		err = fmt.Errorf("Address len exceed limit. (%d > %d)", len(addr), c.MaxAddressLength)
	}
	c.AddressRange = addr
	return
}
