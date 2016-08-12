package PDU

import "github.com/linxGnu/gosmpp/Data"

type BindTransmitter struct {
	BindRequest
}

func NewBindTransmitter() *BindTransmitter {
	a := &BindTransmitter{}
	a.Construct()

	return a
}

func (c *BindTransmitter) Construct() {
	defer c.SetRealReference(c)
	c.BindRequest.Construct()

	c.SetCommandId(Data.BIND_TRANSMITTER)
}

func (c *BindTransmitter) GetInstance() (IPDU, error) {
	return NewBindTransmitter(), nil
}

func (c *BindTransmitter) CreateResponse() (IResponse, error) {
	return NewBindTransmitterResp(), nil
}

func (c *BindTransmitter) IsTransmitter() bool {
	return true
}

func (c *BindTransmitter) IsReceiver() bool {
	return false
}
