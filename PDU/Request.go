package PDU

import "fmt"

type IRequest interface {
	IPDU
	GetResponse() (IResponse, error)
	CreateResponse() (IResponse, error)
	CanResponse() bool
}

type Request struct {
	PDU
}

func NewRequest() *Request {
	a := &Request{}
	a.Construct()

	return a
}

func NewRequestWithCmdId(cmdId int32) *Request {
	a := NewRequest()
	a.SetCommandId(cmdId)

	return a
}

func (c *Request) Construct() {
	defer c.SetRealReference(c)
	c.PDU.Construct()
}

func (c *Request) GetResponse() (res IResponse, err error) {
	defer func() {
		if errs := recover(); errs != nil {
			err = fmt.Errorf("%v", errs)
		}
	}()

	source := c.This.(IRequest)

	res, err = source.CreateResponse()
	if err != nil {
		return
	}

	res.SetSequenceNumber(c.GetSequenceNumber())
	res.SetOriginalRequest(source)

	return
}

func (c *Request) GetResponseCommandId() (int32, error) {
	res, err := c.This.(IRequest).CreateResponse()
	if err != nil {
		return 0, err
	}

	return res.GetCommandId(), nil
}

func (c *Request) CanResponse() bool {
	return true
}

func (c *Request) IsRequest() bool {
	return true
}

func (c *Request) IsResponse() bool {
	return false
}
