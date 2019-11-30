package PDU

import "github.com/linxGnu/gosmpp/Data"

type BindTransceiver struct {
	BindRequest
}

func NewBindTransceiver() *BindTransceiver {
	a := &BindTransceiver{}
	a.Construct()

	return a
}

func (c *BindTransceiver) Construct() {
	defer c.SetRealReference(c)
	c.BindRequest.Construct()

	c.SetCommandId(Data.BIND_TRANSCEIVER)
}

func (c *BindTransceiver) GetInstance() (IPDU, error) {
	return NewBindTransceiver(), nil
}

func (c *BindTransceiver) CreateResponse() (IResponse, error) {
	return NewBindTransceiverResp(), nil
}

func (c *BindTransceiver) IsTransmitter() bool {
	return true
}

func (c *BindTransceiver) IsReceiver() bool {
	return true
}
