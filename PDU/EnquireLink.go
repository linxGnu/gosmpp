package PDU

import (
	"github.com/tsocial/gosmpp/Data"
	"github.com/tsocial/gosmpp/Exception"
	"github.com/tsocial/gosmpp/Utils"
)

type EnquireLink struct {
	Request
}

func NewEnquireLink() *EnquireLink {
	a := &EnquireLink{}
	a.Construct()

	return a
}

func (c *EnquireLink) Construct() {
	defer c.SetRealReference(c)
	c.Request.Construct()

	c.SetCommandId(Data.ENQUIRE_LINK)
}

func (c *EnquireLink) GetInstance() (IPDU, error) {
	return NewEnquireLink(), nil
}

func (c *EnquireLink) CreateResponse() (IResponse, error) {
	return NewEnquireLinkResp(), nil
}

func (c *EnquireLink) SetBody(buf *Utils.ByteBuffer) (*Exception.Exception, IPDU) {
	return nil, nil
}

func (c *EnquireLink) GetBody() (*Utils.ByteBuffer, *Exception.Exception, IPDU) {
	return nil, nil, nil
}
