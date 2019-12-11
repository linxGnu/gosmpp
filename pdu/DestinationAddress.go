package pdu

import (
	"fmt"

	"github.com/linxGnu/gosmpp/data"
)

// DestinationAddress represents Address or Distribution List based on destination flag.
type DestinationAddress struct {
	destFlag byte
	address  Address
	dl       DistributionList
}

// NewDestinationAddress returns new DestinationAddress.
func NewDestinationAddress() (c DestinationAddress) {
	c.destFlag = data.DFLT_DEST_FLAG
	return
}

// NewDestinationAddressFromAddress returns new DestinationAddress from an address.
func NewDestinationAddressFromAddress(addr string) (c DestinationAddress, err error) {
	err = c.SetAddress(addr)
	return
}

// NewDestinationAddressFromDistributionList returns new DestinationAddress from a DistributionList name.
func NewDestinationAddressFromDistributionList(name string) (c DestinationAddress, err error) {
	err = c.SetDistributionList(name)
	return
}

// Unmarshal from buffer.
func (c *DestinationAddress) Unmarshal(b *ByteBuffer) (err error) {
	if c.destFlag, err = b.ReadByte(); err == nil {
		switch c.destFlag {

		case data.SM_DEST_SME_ADDRESS:
			err = c.address.Unmarshal(b)

		case data.SM_DEST_DL_NAME:
			err = c.dl.Unmarshal(b)

		default:
			err = fmt.Errorf("Unrecognize dest_flag %d", c.destFlag)

		}
	}
	return
}

// Marshal to buffer.
func (c *DestinationAddress) Marshal(b *ByteBuffer) {
	switch c.destFlag {
	case data.SM_DEST_DL_NAME:
		_ = b.WriteByte(data.SM_DEST_DL_NAME)
		c.dl.Marshal(b)

	default:
		_ = b.WriteByte(data.SM_DEST_SME_ADDRESS)
		c.address.Marshal(b)
	}
}

// Address returns underlying Address.
func (c *DestinationAddress) Address() Address {
	return c.address
}

// DistributionList returns underlying DistributionList.
func (c *DestinationAddress) DistributionList() DistributionList {
	return c.dl
}

// SetAddress marks DistributionAddress as a SME Address and assign.
func (c *DestinationAddress) SetAddress(addr string) (err error) {
	c.destFlag = data.SM_DEST_SME_ADDRESS
	c.address, err = NewAddressWithAddr(addr)
	return
}

// SetDistributionList marks DistributionAddress as a DistributionList and assign.
func (c *DestinationAddress) SetDistributionList(name string) (err error) {
	c.destFlag = data.SM_DEST_DL_NAME
	c.dl, err = NewDistributionList(name)
	return
}

// HasValue returns true if underlying DistributionList/Address is assigned.
func (c *DestinationAddress) HasValue() bool {
	return c.destFlag != byte(data.DFLT_DEST_FLAG)
}

// IsAddress returns true if DestinationAddress is a SME Address.
func (c *DestinationAddress) IsAddress() bool {
	return c.destFlag == byte(data.SM_DEST_SME_ADDRESS)
}

// IsDistributionList returns true if DestinationAddress is a DistributionList.
func (c *DestinationAddress) IsDistributionList() bool {
	return c.destFlag == byte(data.SM_DEST_DL_NAME)
}

// DestinationAddresses represents list of DestinationAddress.
type DestinationAddresses struct {
	l []DestinationAddress
}

// NewDestinationAddresses returns list of DestinationAddress.
func NewDestinationAddresses() (u DestinationAddresses) {
	u.l = make([]DestinationAddress, 0, 8)
	return
}

// Add to list.
func (c *DestinationAddresses) Add(addresses ...DestinationAddress) {
	c.l = append(c.l, addresses...)
}

// Get list.
func (c *DestinationAddresses) Get() []DestinationAddress {
	return c.l
}

// Unmarshal from buffer.
func (c *DestinationAddresses) Unmarshal(b *ByteBuffer) (err error) {
	var n byte
	if n, err = b.ReadByte(); err == nil {
		c.l = make([]DestinationAddress, n)

		var i byte
		for ; i < n; i++ {
			if err = c.l[i].Unmarshal(b); err != nil {
				return
			}
		}
	}
	return
}

// Marshal to buffer.
func (c *DestinationAddresses) Marshal(b *ByteBuffer) {
	n := byte(len(c.l))
	_ = b.WriteByte(n)

	var i byte
	for ; i < n; i++ {
		c.l[i].Marshal(b)
	}
}
