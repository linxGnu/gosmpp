package pdu

import (
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/utils"
)

// BindResp PDU.
type BindResp struct {
	base
	Request *BindRequest
}

// NewBindResp returns BindResp.
func NewBindResp(req *BindRequest) (c *BindResp) {
	c = &BindResp{
		base:    newBase(),
		Request: req,
	}

	if req != nil {
		switch req.Type {
		case Transceiver:
			c.CommandID = data.BIND_TRANSCEIVER_RESP

		case Receiver:
			c.CommandID = data.BIND_RECEIVER_RESP

		case Transmitter:
			c.CommandID = data.BIND_TRANSMITTER_RESP
		}
	}

	return
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
	c.base.marshal(b, nil)
}

// Unmarshal implements PDU interface.
func (c *BindResp) Unmarshal(b *utils.ByteBuffer) error {
	return c.base.unmarshal(b, nil)
}
