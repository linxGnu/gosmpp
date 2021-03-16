package gosmpp

import (
	"sync/atomic"
	"time"

	"github.com/linxGnu/gosmpp/pdu"
)

// TransceiveSettings is listener for Transceiver.
type TransceiveSettings struct {
	// WriteTimeout is timeout for submitting PDU.
	WriteTimeout time.Duration

	// EnquireLink periodically sends EnquireLink to SMSC.
	// Zero duration means disable auto enquire link.
	EnquireLink time.Duration

	// OnPDU handles received PDU from SMSC.
	//
	// `Responded` flag indicates this pdu is responded automatically,
	// no manual respond needed.
	OnPDU PDUCallback

	// OnSubmitError notifies fail-to-submit PDU with along error.
	OnSubmitError PDUErrorCallback

	// OnReceivingError notifies happened error while reading PDU
	// from SMSC.
	OnReceivingError ErrorCallback

	// OnRebindingError notifies error while rebinding.
	OnRebindingError ErrorCallback

	// OnClosed notifies `closed` event due to State.
	OnClosed ClosedCallback
}

type transceiver struct {
	settings TransceiveSettings
	conn     *Connection
	in       *receiver
	out      *transmitter
	state    int32
}

func newTransceiver(conn *Connection, settings TransceiveSettings) Transceiver {
	t := &transceiver{
		settings: settings,
		conn:     conn,
	}

	t.out = newTransmitter(conn, TransmitSettings{
		Timeout: settings.WriteTimeout,

		EnquireLink: settings.EnquireLink,

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
	})

	t.in = newReceiver(conn, ReceiveSettings{
		OnPDU: settings.OnPDU,

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
			if t.out.Submit(p) != nil { // only happened when transceiver is closed
				_, _ = t.out.write(marshal(p))
			}
		},
	})

	t.out.start(false)
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
		_ = t.out.close(StoppingProcessOnly)
		_ = t.in.close(StoppingProcessOnly)

		// close underlying conn
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
