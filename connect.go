package gosmpp

import (
	"github.com/linxGnu/gosmpp/pdu"
)

// Auth represents basic authentication to SMSC.
type Auth struct {
	// SMSC represents SMSC address.
	SMSC string

	// authentication infos
	SystemID   string
	Password   string
	SystemType string
}

func newBindRequest(s Auth, bindingType pdu.BindingType) (bindReq *pdu.BindRequest) {
	bindReq = pdu.NewBindRequest(bindingType)
	bindReq.SystemID = s.SystemID
	bindReq.Password = s.Password
	bindReq.SystemType = s.SystemType
	return
}

// ConnectAsReceiver connects to SMSC as Receiver.
func ConnectAsReceiver(s Auth) (conn *Connection, err error) {
	conn, err = connect(s.SMSC, newBindRequest(s, pdu.Receiver))
	return
}

// ConnectAsTransmitter connects to SMSC as Transmitter.
func ConnectAsTransmitter(s Auth) (conn *Connection, err error) {
	conn, err = connect(s.SMSC, newBindRequest(s, pdu.Transmitter))
	return
}

// ConnectAsTransceiver connects to SMSC as Transceiver.
func ConnectAsTransceiver(s Auth) (conn *Connection, err error) {
	conn, err = connect(s.SMSC, newBindRequest(s, pdu.Transceiver))
	return
}
