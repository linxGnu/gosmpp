package gosmpp

import (
	"io"
	"net"

	"github.com/linxGnu/gosmpp/pdu"
)

// Connection wraps over net.Conn along with setting(s).
type Connection struct {
	Conn      net.Conn // underlying connection
	Dedicated bool     // indicates connection is dedicated from invoker. Invoker won't care about its state
}

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
