package gosmpp

import (
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/PDU"
)

type Transmitter struct {
	connection IConnection
}

func NewTransmitter() *Transmitter {
	a := &Transmitter{}

	return a
}

func NewTransmitterWithConnection(con IConnection) *Transmitter {
	a := NewTransmitter()
	a.connection = con

	return a
}

func (c *Transmitter) Send(pdu PDU.IPDU) *Exception.Exception {
	if pdu == nil {
		return nil
	}

	if c.connection == nil {
		return Exception.NewExceptionFromStr("Connection not set")
	}

	pdu.AssignSequenceNumber()

	dat, err, _ := pdu.GetData()
	if err != nil {
		return err
	}

	return c.connection.Send(dat)
}
