package PDU

import (
	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/Utils"
)

type ReplaceSMResp struct {
	Response
}

func NewReplaceSMResp() *ReplaceSMResp {
	a := &ReplaceSMResp{}
	a.Construct()

	return a
}

func (c *ReplaceSMResp) Construct() {
	defer c.SetRealReference(c)
	c.Response.Construct()

	c.SetCommandId(Data.REPLACE_SM_RESP)
}

func (c *ReplaceSMResp) GetInstance() (IPDU, error) {
	return NewReplaceSMResp(), nil
}

func (c *ReplaceSMResp) SetBody(buf *Utils.ByteBuffer) (*Exception.Exception, IPDU) {
	return nil, nil
}

func (c *ReplaceSMResp) GetBody() (*Utils.ByteBuffer, *Exception.Exception, IPDU) {
	return nil, nil, nil
}
