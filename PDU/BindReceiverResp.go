package PDU

import "github.com/linxGnu/gosmpp/Data"

type BindReceiverResp struct {
	BindResponse
}

func NewBindReceiverResp() *BindReceiverResp {
	a := &BindReceiverResp{}
	a.Construct()

	return a
}

func (c *BindReceiverResp) Construct() {
	defer c.SetRealReference(c)
	c.BindResponse.Construct()

	c.SetCommandId(Data.BIND_RECEIVER_RESP)
}

func (c *BindReceiverResp) GetInstance() (IPDU, error) {
	return NewBindReceiverResp(), nil
}
