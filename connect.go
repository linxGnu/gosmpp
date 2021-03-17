package gosmpp

import (
	"net"

	"github.com/linxGnu/gosmpp/pdu"
)

var (
	// NonTLSDialer is non-tls connection dialer.
	NonTLSDialer = func(addr string) (net.Conn, error) {
		return net.Dial("tcp", addr)
	}
)

// Dialer is connection dialer.
type Dialer func(addr string) (net.Conn, error)

// Auth represents basic authentication to SMSC.
type Auth struct {
	// SMSC is SMSC address.
	SMSC       string
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

// Connector is connection factory interface.
type Connector interface {
	Connect() (conn *Connection, err error)
}

type connector struct {
	dialer      Dialer
	auth        Auth
	bindingType pdu.BindingType
}

func (c *connector) Connect() (conn *Connection, err error) {
	conn, err = connect(c.dialer, c.auth.SMSC, newBindRequest(c.auth, c.bindingType))
	return
}

// TXConnector returns a Transmitter (TX) connector.
func TXConnector(dialer Dialer, auth Auth) Connector {
	return &connector{
		dialer:      dialer,
		auth:        auth,
		bindingType: pdu.Transmitter,
	}
}

// RXConnector returns a Receiver (RX) connector.
func RXConnector(dialer Dialer, auth Auth) Connector {
	return &connector{
		dialer:      dialer,
		auth:        auth,
		bindingType: pdu.Receiver,
	}
}

// TRXConnector returns a Transceiver (TRX) connector.
func TRXConnector(dialer Dialer, auth Auth) Connector {
	return &connector{
		dialer:      dialer,
		auth:        auth,
		bindingType: pdu.Transceiver,
	}
}
