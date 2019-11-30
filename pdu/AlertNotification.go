package PDU

import (
	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/Utils"
)

type AlertNotification struct {
	Request
}

func NewAlertNotification() *AlertNotification {
	a := &AlertNotification{}
	a.Construct()

	return a
}

func (c *AlertNotification) Construct() {
	defer c.SetRealReference(c)
	c.Request.Construct()

	c.SetCommandId(Data.ALERT_NOTIFICATION)
}

func (c *AlertNotification) GetInstance() (IPDU, error) {
	return NewAlertNotification(), nil
}

func (c *AlertNotification) CanResponse() bool {
	return false
}

func (c *AlertNotification) CreateResponse() (IResponse, error) {
	return nil, nil
}

func (c *AlertNotification) GetBody() (*Utils.ByteBuffer, *Exception.Exception, IPDU) {
	return nil, nil, nil
}

func (c *AlertNotification) SetBody(buffer *Utils.ByteBuffer) (*Exception.Exception, IPDU) {
	return nil, nil
}
