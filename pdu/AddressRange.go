package pdu

import (
	"fmt"

	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/utils"
)

// AddressRange smpp address range of src and dst.
type AddressRange struct {
	ton          byte
	npi          byte
	addressRange string
}

// NewAddressRange create new AddressRange with default max length.
func NewAddressRange() *AddressRange {
	return &AddressRange{ton: data.GetDefaultTon(), npi: data.GetDefaultNpi()}
}

// NewAddressRangeWithAddr create new AddressRange.
func NewAddressRangeWithAddr(addr string) (a *AddressRange, err error) {
	a = NewAddressRange()
	err = a.SetAddressRange(addr)
	return
}

// NewAddressRangeWithTonNpiLen create new AddressRange with ton, npi, max length.
func NewAddressRangeWithTonNpiLen(ton, npi byte) *AddressRange {
	return &AddressRange{ton: ton, npi: npi}
}

// Unmarshal from buffer.
func (c *AddressRange) Unmarshal(b *utils.ByteBuffer) (err error) {
	if c.ton, err = b.ReadByte(); err == nil {
		if c.npi, err = b.ReadByte(); err == nil {
			c.addressRange, err = b.ReadCString()
		}
	}
	return
}

// Marshal to buffer.
func (c *AddressRange) Marshal(b *utils.ByteBuffer) {
	b.Grow(3 + len(c.addressRange))

	_ = b.WriteByte(c.ton)
	_ = b.WriteByte(c.npi)
	_ = b.WriteCString(c.addressRange)
}

// SetAddressRange sets address range.
func (c *AddressRange) SetAddressRange(addr string) (err error) {
	if len(addr) > data.SM_ADDR_RANGE_LEN {
		err = fmt.Errorf("Address len exceed limit. (%d > %d)", len(addr), data.SM_ADDR_RANGE_LEN)
	}
	c.addressRange = addr
	return
}

// AddressRange returns assigned address range (in string).
func (c *AddressRange) AddressRange() string {
	return c.addressRange
}

// Ton returns assigned ton.
func (c *AddressRange) Ton() byte {
	return c.ton
}

// Npi returns assigned npi.
func (c *AddressRange) Npi() byte {
	return c.npi
}
