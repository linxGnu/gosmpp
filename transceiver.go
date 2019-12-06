package gosmpp

import (
	"sync/atomic"
	"time"

	"github.com/linxGnu/gosmpp/pdu"
)

// TransceiveSettings is listener for Transceiver.
type TransceiveSettings struct {
	// WriteTimeout is timeout/deadline for submitting PDU.
	WriteTimeout time.Duration

	// EnquireLink periodically sends EnquireLink to SMSC.
	// Zero duration means disable auto enquire link.
	EnquireLink time.Duration

	// OnPDU handles received PDU from SMSC.
	OnPDU func(pdu.PDU)

	// OnSubmitError notifies fail-to-submit PDU with along error.
	OnSubmitError func(pdu.PDU, error)

	// OnReceivingError notifies happened error while reading PDU
	// from SMSC.
	OnReceivingError func(error)

	// OnRebindingError notifies error while rebinding.
	OnRebindingError func(error)

	// OnClosed notifies `closed` event due to State.
	OnClosed func(State)
}

type transceiver struct {
	settings TransceiveSettings
	conn     *Connection
	in       *receiver
	out      *transmitter
	state    int32
}

// NewTransceiver creates new Transceiver from bound connection.
func NewTransceiver(conn *Connection, settings TransceiveSettings) Transceiver {
	t := &transceiver{
		settings: settings,
		conn:     conn,
	}

	t.out = newTransmitter(conn, TransmitSettings{
		WriteTimeout:  settings.WriteTimeout,
		EnquireLink:   settings.EnquireLink,
		OnSubmitError: settings.OnSubmitError,
		OnClosed: func(state State) {
			switch state {
			case ExplicitClosing:
				return

			case ConnectionIssue:
				// also close input
				_ = t.in.Close()

				if t.settings.OnClosed != nil {
					t.settings.OnClosed(ConnectionIssue)
				}
			}
		},
	}, false)

	t.in = newReceiver(conn, ReceiveSettings{
		OnPDU:            settings.OnPDU,
		OnReceivingError: settings.OnReceivingError,

		OnClosed: func(state State) {
			switch state {
			case ExplicitClosing:
				return

			case InvalidStreaming, UnbindClosing:
				// also close output
				_ = t.out.Close()

				if t.settings.OnClosed != nil {
					t.settings.OnClosed(state)
				}
			}
		},

		response: func(p pdu.PDU) {
			_ = t.out.Submit(p)
		},
	}, false)

	t.out.start()
	t.in.start()

	return t
}

// SystemID returns tagged SystemID, returned from bind_resp from SMSC.
func (t *transceiver) SystemID() string {
	return t.conn.systemID
}

// Close transceiver and stop underlying daemons.
func (t *transceiver) Close() (err error) {
	if atomic.CompareAndSwapInt32(&t.state, 0, 1) {
		// closing input and output
		_ = t.out.Close()
		_ = t.in.Close()

		// close connection
		err = t.conn.Close()

		// notify transceiver closed
		if t.settings.OnClosed != nil {
			t.settings.OnClosed(ExplicitClosing)
		}
	}
	return
}

// Submit a PDU.
func (t *transceiver) Submit(p pdu.PDU) error {
	return t.out.Submit(p)
}
