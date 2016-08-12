package gosmpp

import "github.com/linxGnu/gosmpp/PDU"

type OutbindEvent struct {
	ReceivedPDUEvent
	serialVersionUID int64
}

func NewOutbindEvent(source *OutbindReceiver, con IConnection, pdu *PDU.Outbind) *OutbindEvent {
	a := &OutbindEvent{}
	a.Source = source
	a.serialVersionUID = 1808913846085130877
	a.Connection = con
	a.Pdu = pdu

	return a
}

func (c *OutbindEvent) GetReceiver() *OutbindReceiver {
	return c.Source.(*OutbindReceiver)
}

func (c *OutbindEvent) GetOutbindPDU() *PDU.Outbind {
	return c.Pdu.(*PDU.Outbind)
}
