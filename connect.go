package gosmpp

import (
	"net"

	"github.com/linxGnu/gosmpp/pdu"
)

// Dialer is connection dialer.
type Dialer func(addr string) (net.Conn, error)

var (
	// NonTLSDialer is non-tls connection dialer.
	NonTLSDialer = func(addr string) (net.Conn, error) {
		return net.Dial("tcp", addr)
	}
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
func ConnectAsReceiver(dialer Dialer, s Auth) (conn *Connection, err error) {
	conn, err = connect(dialer, s.SMSC, newBindRequest(s, pdu.Receiver))
	return
}

// ConnectAsTransmitter connects to SMSC as Transmitter.
func ConnectAsTransmitter(dialer Dialer, s Auth) (conn *Connection, err error) {
	conn, err = connect(dialer, s.SMSC, newBindRequest(s, pdu.Transmitter))
	return
}

// ConnectAsTransceiver connects to SMSC as Transceiver.
func ConnectAsTransceiver(dialer Dialer, s Auth) (conn *Connection, err error) {
	conn, err = connect(dialer, s.SMSC, newBindRequest(s, pdu.Transceiver))
	return
}
