package gosmpp

import (
	"io"

	"github.com/linxGnu/gosmpp/pdu"
)

// Transceiver interface.
type Transceiver interface {
	io.Closer
	Submit(pdu.PDU) error
	SystemID() string
}

// Transmitter interface.
type Transmitter interface {
	io.Closer
	Submit(pdu.PDU) error
	SystemID() string
}

// Receiver interface.
type Receiver interface {
	io.Closer
	SystemID() string
}
