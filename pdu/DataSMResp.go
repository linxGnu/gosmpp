package pdu

import (
	"github.com/linxGnu/gosmpp/data"
)

// DataSMResp PDU.
type DataSMResp struct {
	base
	MessageID string
}

// NewDataSMResp returns DataSMResp.
func NewDataSMResp() PDU {
	c := &DataSMResp{
		base:      newBase(),
		MessageID: data.DFLT_MSGID,
	}
	c.CommandID = data.DATA_SM_RESP
	return c
}

// NewDataSMRespFromReq returns DataSMResp.
func NewDataSMRespFromReq(req *DataSM) PDU {
	c := NewDataSMResp().(*DataSMResp)
	if req != nil {
		c.SequenceNumber = req.SequenceNumber
	}
	return c
}

// CanResponse implements PDU interface.
func (c *DataSMResp) CanResponse() bool {
	return false
}

// GetResponse implements PDU interface.
func (c *DataSMResp) GetResponse() PDU {
	return nil
}

// Marshal implements PDU interface.
func (c *DataSMResp) Marshal(b *ByteBuffer) {
	c.base.marshal(b, func(b *ByteBuffer) {
		b.Grow(len(c.MessageID) + 1)

		_ = b.WriteCString(c.MessageID)
	})
}

// Unmarshal implements PDU interface.
func (c *DataSMResp) Unmarshal(b *ByteBuffer) error {
	return c.base.unmarshal(b, func(b *ByteBuffer) (err error) {
		c.MessageID, err = b.ReadCString()
		return
	})
}
