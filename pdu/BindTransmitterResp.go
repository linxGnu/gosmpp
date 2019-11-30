package PDU

import "github.com/linxGnu/gosmpp/Data"

type BindTransmitterResp struct {
	BindResponse
}

func NewBindTransmitterResp() *BindTransmitterResp {
	a := &BindTransmitterResp{}
	a.Construct()

	return a
}

func (c *BindTransmitterResp) Construct() {
	defer c.SetRealReference(c)
	c.BindResponse.Construct()

	c.SetCommandId(Data.BIND_TRANSMITTER_RESP)
}

func (c *BindTransmitterResp) GetInstance() (IPDU, error) {
	return NewBindTransmitterResp(), nil
}
