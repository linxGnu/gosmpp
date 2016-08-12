package gosmpp

import "github.com/linxGnu/gosmpp/PDU"

type ServerPDUEvent struct {
	ReceivedPDUEvent
}

func NewServerPDUEvent(source *Receiver, con IConnection, pdu PDU.IPDU) *ServerPDUEvent {
	a := &ServerPDUEvent{}
	a.serialVersionUID = 8400363453588829420
	a.Source = source
	a.Connection = con
	a.Pdu = pdu

	return a
}

func (c *ServerPDUEvent) GetReceiver() *Receiver {
	return c.Source.(*Receiver)
}
