package PDU

import (
	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/Utils"
)

type EnquireLinkResp struct {
	Response
}

func NewEnquireLinkResp() *EnquireLinkResp {
	a := &EnquireLinkResp{}
	a.Construct()

	return a
}

func (c *EnquireLinkResp) Construct() {
	defer c.SetRealReference(c)
	c.Response.Construct()

	c.SetCommandId(Data.ENQUIRE_LINK_RESP)
}

func (c *EnquireLinkResp) GetInstance() (IPDU, error) {
	return NewEnquireLinkResp(), nil
}

func (c *EnquireLinkResp) SetBody(buf *Utils.ByteBuffer) (*Exception.Exception, IPDU) {
	return nil, nil
}

func (c *EnquireLinkResp) GetBody() (*Utils.ByteBuffer, *Exception.Exception, IPDU) {
	return nil, nil, nil
}
