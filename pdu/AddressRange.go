package pdu

import "github.com/linxGnu/gosmpp/data"

// AddressRange smpp address range of src and dst.
type AddressRange struct {
	Ton          byte
	Npi          byte
	AddressRange string
}

// NewAddressRange create new AddressRange with default max length.
func NewAddressRange() AddressRange {
	return AddressRange{Ton: data.GetDefaultTon(), Npi: data.GetDefaultNpi()}
}

// NewAddressRangeWithAddr create new AddressRange.
func NewAddressRangeWithAddr(addr string) (a AddressRange, err error) {
	a = NewAddressRange()
	a.AddressRange = addr
	return
}

// NewAddressRangeWithTonNpi create new AddressRange with ton, npi.
func NewAddressRangeWithTonNpi(ton, npi byte) AddressRange {
	return AddressRange{Ton: ton, Npi: npi}
}

// NewAddressRangeWithTonNpiAddr returns new address with ton, npi, addr string.
func NewAddressRangeWithTonNpiAddr(ton, npi byte, addr string) (a AddressRange, err error) {
	a = NewAddressRangeWithTonNpi(ton, npi)
	a.AddressRange = addr
	return
}

// Unmarshal from buffer.
func (c *AddressRange) Unmarshal(b *ByteBuffer) (err error) {
	if c.Ton, err = b.ReadByte(); err == nil {
		if c.Npi, err = b.ReadByte(); err == nil {
			c.AddressRange, err = b.ReadCString()
		}
	}
	return
}

// Marshal to buffer.
func (c *AddressRange) Marshal(b *ByteBuffer) {
	b.Grow(3 + len(c.AddressRange))

	_ = b.WriteByte(c.Ton)
	_ = b.WriteByte(c.Npi)
	_ = b.WriteCString(c.AddressRange)
}
