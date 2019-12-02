package pdu

import (
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/utils"
)

// QuerySM PDU.
type QuerySM struct {
	base
	MessageID  string
	SourceAddr Address
}

// NewQuerySM returns new QuerySM PDU.
func NewQuerySM() (c *QuerySM) {
	c = &QuerySM{
		SourceAddr: *NewAddress(),
	}
	c.CommandID = data.QUERY_SM
	return
}

// CanResponse implements PDU interface.
func (c *QuerySM) CanResponse() bool {
	return true
}

// GetResponse implements PDU interface.
func (c *QuerySM) GetResponse() PDU {
	return NewQuerySMResp(c)
}

// Marshal implements PDU interface.
func (c *QuerySM) Marshal(b *utils.ByteBuffer) {
	c.base.marshal(b, func(b *utils.ByteBuffer) {
		b.Grow(len(c.MessageID) + 1)

		b.WriteCString(c.MessageID)
		c.SourceAddr.Marshal(b)
	})
}

// Unmarshal implements PDU interface.
func (c *QuerySM) Unmarshal(b *utils.ByteBuffer) error {
	return c.base.unmarshal(b, func(b *utils.ByteBuffer) (err error) {
		if c.MessageID, err = b.ReadCString(); err == nil {
			err = c.SourceAddr.Unmarshal(b)
		}
		return
	})
}
