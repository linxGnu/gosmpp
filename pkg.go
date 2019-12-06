package gosmpp

import (
	"io"

	"github.com/linxGnu/gosmpp/pdu"
)

// Transceiver interface.
type Transceiver interface {
	io.Closer
	Submit(pdu.PDU) error
}

// Transmitter interface.
type Transmitter interface {
	io.Closer
	Submit(pdu.PDU) error
}

// Receiver interface.
type Receiver interface {
	io.Closer
}
