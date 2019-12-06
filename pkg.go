package gosmpp

import (
	"io"

	"github.com/linxGnu/gosmpp/pdu"
)

// Submiter submits PDU to SMSC.
type Submiter interface {
	Submit(pdu.PDU) error
}

// Transceiver interface.
type Transceiver interface {
	io.Closer
	Submiter
}

// Transmitter interface.
type Transmitter interface {
	io.Closer
	Submiter
}

// Receiver interface.
type Receiver interface {
	io.Closer
}
