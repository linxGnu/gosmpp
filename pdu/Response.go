package PDU

type IResponse interface {
	IPDU
	SetOriginalRequest(t IRequest)
	GetOriginalRequest() IRequest
	CanResponse() bool
}

type Response struct {
	PDU
	originalRequest IRequest
}

func NewResponse() *Response {
	a := &Response{}
	a.Construct()

	return a
}

func NewResponseWithCmdId(cmdId int32) *Response {
	a := NewResponse()
	a.SetCommandId(cmdId)

	return a
}

func (c *Response) Construct() {
	defer c.SetRealReference(c)
	c.PDU.Construct()
}

func (c *Response) CanResponse() bool {
	return false
}

func (c *Response) IsRequest() bool {
	return false
}

func (c *Response) IsResponse() bool {
	return true
}

func (c *Response) SetOriginalRequest(t IRequest) {
	c.originalRequest = t
}

func (c *Response) GetOriginalRequest() IRequest {
	return c.originalRequest
}
