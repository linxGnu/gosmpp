package PDU

import "github.com/linxGnu/gosmpp/Data"

type BindTransciever struct {
	BindRequest
}

func NewBindTransciever() *BindTransciever {
	a := &BindTransciever{}
	a.Construct()

	return a
}

func (c *BindTransciever) Construct() {
	defer c.SetRealReference(c)
	c.BindRequest.Construct()

	c.SetCommandId(Data.BIND_TRANSCEIVER)
}

func (c *BindTransciever) GetInstance() (IPDU, error) {
	return NewBindTransciever(), nil
}

func (c *BindTransciever) CreateResponse() (IResponse, error) {
	return NewBindTranscieverResp(), nil
}

func (c *BindTransciever) IsTransmitter() bool {
	return true
}

func (c *BindTransciever) IsReceiver() bool {
	return true
}
