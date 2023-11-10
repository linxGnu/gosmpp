package gosmpp

import (
	"io"
	"time"

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

// Settings for TX (transmitter), RX (receiver), TRX (transceiver).
type Settings struct {
	// ReadTimeout is timeout for reading PDU from SMSC.
	// Underlying net.Conn will be stricted with ReadDeadline(now + timeout).
	// This setting is very important to detect connection failure.
	//
	// Must: ReadTimeout > max(0, EnquireLink)
	ReadTimeout time.Duration

	// WriteTimeout is timeout for submitting PDU.
	WriteTimeout time.Duration

	// EnquireLink periodically sends EnquireLink to SMSC.
	// The duration must not be smaller than 1 minute.
	//
	// Zero duration disables auto enquire link.
	EnquireLink time.Duration

	// OnPDU handles received PDU from SMSC.
	//
	// `Responded` flag indicates this pdu is responded automatically,
	// no manual respond needed.
	//
	// Will be ignored if OnAllPDU is set
	OnPDU PDUCallback

	// OnAllPDU handles all received PDU from SMSC.
	//
	// This pdu is NOT responded to automatically,
	// manual response/handling is needed
	//
	// User can also decide to close bind by retuning true, default is false
	OnAllPDU AllPDUCallback

	// OnReceivingError notifies happened error while reading PDU
	// from SMSC.
	OnReceivingError ErrorCallback

	// OnSubmitError notifies fail-to-submit PDU with along error.
	OnSubmitError PDUErrorCallback

	// OnRebindingError notifies error while rebinding.
	OnRebindingError ErrorCallback

	// OnClosed notifies `closed` event due to State.
	OnClosed ClosedCallback

	response func(pdu.PDU)
}
