package PDU

import "github.com/linxGnu/gosmpp/Data"

type BindTransceiverResp struct {
	BindResponse
}

func NewBindTransceiverResp() *BindTransceiverResp {
	a := &BindTransceiverResp{}
	a.Construct()

	return a
}

func (c *BindTransceiverResp) Construct() {
	defer c.SetRealReference(c)
	c.BindResponse.Construct()

	c.SetCommandId(Data.BIND_TRANSCEIVER_RESP)
}

func (c *BindTransceiverResp) GetInstance() (IPDU, error) {
	return NewBindTransceiverResp(), nil
}
