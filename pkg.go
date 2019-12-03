package gosmpp

import (
	"io"

	"github.com/linxGnu/gosmpp/pdu"
)

// Writer submits PDU to SMSC.
type Writer interface {
	Write(pdu.PDU) error
}

// Reader handles received PDU from SMSC.
type Reader interface {
	Read(func(pdu.PDU))
}

// Transceiver interface.
type Transceiver interface {
	io.Closer
	Writer
	Reader
}

// Transmitter interface.
type Transmitter interface {
	io.Closer
	Writer
}

// Receiver interface.
type Receiver interface {
	io.Closer
	Reader
}
