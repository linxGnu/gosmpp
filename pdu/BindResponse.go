package pdu

import (
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/utils"
)

// BindResp PDU.
type BindResp struct {
	base
	Request  BindRequest
	SystemID string
}

// NewBindResp returns BindResp.
func NewBindResp(req BindRequest) (c *BindResp) {
	c = &BindResp{
		base:    newBase(),
		Request: req,
	}

	switch req.BindingType {
	case Transceiver:
		c.CommandID = data.BIND_TRANSCEIVER_RESP

	case Receiver:
		c.CommandID = data.BIND_RECEIVER_RESP

	case Transmitter:
		c.CommandID = data.BIND_TRANSMITTER_RESP
	}

	return
}

// NewBindTransmitterResp returns new bind transmitter resp.
func NewBindTransmitterResp() PDU {
	c := &BindResp{
		base: newBase(),
	}
	c.CommandID = data.BIND_TRANSMITTER_RESP
	return c
}

// NewBindTransceiverResp returns new bind transceiver resp.
func NewBindTransceiverResp() PDU {
	c := &BindResp{
		base: newBase(),
	}
	c.CommandID = data.BIND_TRANSCEIVER_RESP
	return c
}

// NewBindReceiverResp returns new bind receiver resp.
func NewBindReceiverResp() PDU {
	c := &BindResp{
		base: newBase(),
	}
	c.CommandID = data.BIND_RECEIVER_RESP
	return c
}

// CanResponse implements PDU interface.
func (c *BindResp) CanResponse() bool {
	return false
}

// GetResponse implements PDU interface.
func (c *BindResp) GetResponse() PDU {
	return nil
}

// Marshal implements PDU interface.
func (c *BindResp) Marshal(b *utils.ByteBuffer) {
	c.base.marshal(b, func(w *utils.ByteBuffer) {
		w.Grow(len(c.SystemID) + 1)
		_ = w.WriteCString(c.SystemID)
	})
}

// Unmarshal implements PDU interface.
func (c *BindResp) Unmarshal(b *utils.ByteBuffer) error {
	return c.base.unmarshal(b, func(w *utils.ByteBuffer) (err error) {
		if c.CommandID == data.BIND_TRANSCEIVER_RESP || c.CommandStatus == data.ESME_ROK {
			c.SystemID, err = w.ReadCString()
		}
		return
	})
}
