package gosmpp

import "github.com/linxGnu/gosmpp/PDU"

type ReceivedPDUEvent struct {
	serialVersionUID int64
	Connection       IConnection
	Pdu              PDU.IPDU
	Source           interface{}
}

func NewReceivedPDUEvent(source *ReceiverBase, con IConnection, pdu PDU.IPDU) *ReceivedPDUEvent {
	a := &ReceivedPDUEvent{}
	a.serialVersionUID = 2888578757849035826
	a.Source = source
	a.Connection = con
	a.Pdu = pdu

	return a
}

func (c *ReceivedPDUEvent) GetConnection() IConnection {
	return c.Connection
}

func (c *ReceivedPDUEvent) GetPDU() PDU.IPDU {
	return c.Pdu
}
