package pdu

import (
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/utils"
)

// QuerySM PDU is issued by the ESME to query the status of a previously submitted short message.
// The matching mechanism is based on the SMSC assigned message_id and source address. Where the
// original submit_sm, data_sm or submit_multi ‘source address’ was defaulted to NULL, then the
// source address in the query_sm command should also be set to NULL.
type QuerySM struct {
	base
	MessageID  string
	SourceAddr Address
}

// NewQuerySM returns new QuerySM PDU.
func NewQuerySM() PDU {
	c := &QuerySM{
		SourceAddr: NewAddress(),
	}
	c.CommandID = data.QUERY_SM
	return c
}

// CanResponse implements PDU interface.
func (c *QuerySM) CanResponse() bool {
	return true
}

// GetResponse implements PDU interface.
func (c *QuerySM) GetResponse() PDU {
	return NewQuerySMResp()
}

// Marshal implements PDU interface.
func (c *QuerySM) Marshal(b *utils.ByteBuffer) {
	c.base.marshal(b, func(b *utils.ByteBuffer) {
		b.Grow(len(c.MessageID) + 1)

		_ = b.WriteCString(c.MessageID)
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
