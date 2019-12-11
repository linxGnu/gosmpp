package pdu

import (
	"fmt"

	"github.com/linxGnu/gosmpp/data"
)

// DistributionList represents group of contacts.
type DistributionList struct {
	name string
}

// NewDistributionList returns a new DistributionList.
func NewDistributionList(name string) (c DistributionList, err error) {
	err = c.SetName(name)
	return
}

// Unmarshal from buffer.
func (c *DistributionList) Unmarshal(b *ByteBuffer) (err error) {
	c.name, err = b.ReadCString()
	return
}

// Marshal to buffer.
func (c *DistributionList) Marshal(b *ByteBuffer) {
	b.Grow(1 + len(c.name))

	_ = b.WriteCString(c.name)
}

// SetName sets DistributionList name.
func (c *DistributionList) SetName(name string) (err error) {
	if len(name) > data.SM_DL_NAME_LEN {
		err = fmt.Errorf("Distribution List name exceed limit. (%d > %d)", len(name), data.SM_DL_NAME_LEN)
	} else {
		c.name = name
	}
	return
}

// Name returns name of DistributionList
func (c DistributionList) Name() string {
	return c.name
}
