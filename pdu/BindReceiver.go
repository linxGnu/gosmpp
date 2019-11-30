package PDU

import "github.com/linxGnu/gosmpp/Data"

type BindReceiver struct {
	BindRequest
}

func NewBindReceiver() *BindReceiver {
	a := &BindReceiver{}
	a.Construct()

	return a
}

func (c *BindReceiver) Construct() {
	defer c.SetRealReference(c)
	c.BindRequest.Construct()

	c.SetCommandId(Data.BIND_RECEIVER)
}

func (c *BindReceiver) GetInstance() (IPDU, error) {
	return NewBindReceiver(), nil
}

func (c *BindReceiver) CreateResponse() (IResponse, error) {
	return NewBindReceiverResp(), nil
}

func (c *BindReceiver) IsTransmitter() bool {
	return false
}

func (c *BindReceiver) IsReceiver() bool {
	return true
}
