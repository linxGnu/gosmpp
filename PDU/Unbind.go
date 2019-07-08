package PDU

import (
	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/Utils"
)

type Unbind struct {
	Request
}

func NewUnbind() *Unbind {
	a := &Unbind{}
	a.Construct()

	return a
}

func (c *Unbind) Construct() {
	defer c.SetRealReference(c)
	c.Request.Construct()

	c.SetCommandId(Data.UNBIND)
}

func (c *Unbind) GetInstance() (IPDU, error) {
	return NewUnbind(), nil
}

func (c *Unbind) CreateResponse() (IResponse, error) {
	return NewUnbindResp(), nil
}

func (c *Unbind) SetBody(buf *Utils.ByteBuffer) (err *Exception.Exception, source IPDU) {
	source = c.This.(IPDU)

	return nil, source
}

func (c *Unbind) GetBody() (buf *Utils.ByteBuffer, err *Exception.Exception, source IPDU) {
	source = c.This.(IPDU)

	return nil, nil, source
}
