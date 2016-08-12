package PDU

import "github.com/linxGnu/gosmpp/Data"

type BindTranscieverResp struct {
	BindResponse
}

func NewBindTranscieverResp() *BindTranscieverResp {
	a := &BindTranscieverResp{}
	a.Construct()

	return a
}

func (c *BindTranscieverResp) Construct() {
	defer c.SetRealReference(c)
	c.BindResponse.Construct()

	c.SetCommandId(Data.BIND_TRANSCEIVER_RESP)
}

func (c *BindTranscieverResp) GetInstance() (IPDU, error) {
	return NewBindTranscieverResp(), nil
}
