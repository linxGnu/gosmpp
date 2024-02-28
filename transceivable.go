package gosmpp

import (
	"github.com/linxGnu/gosmpp/pdu"
	"sync/atomic"
)

type transceivable struct {
	settings Settings

	conn *Connection
	in   *receivable
	out  *transmittable

	aliveState int32
}

func newTransceivable(conn *Connection, settings Settings) *transceivable {

	t := &transceivable{
		settings: settings,
		conn:     conn,
	}

	t.out = newTransmittable(conn, Settings{
		WriteTimeout: settings.WriteTimeout,

		EnquireLink: settings.EnquireLink,

		OnSubmitError: settings.OnSubmitError,

		OnClosed: func(state State) {
			switch state {
			case ConnectionIssue:
				// also close input
				_ = t.in.close(ExplicitClosing)

				if t.settings.OnClosed != nil {
					t.settings.OnClosed(ConnectionIssue)
				}
			default:
				return
			}
		},

		RequestWindowConfig: settings.RequestWindowConfig,
	})

	t.in = newReceivable(conn, Settings{
		ReadTimeout: settings.ReadTimeout,

		OnPDU: settings.OnPDU,

		OnAllPDU: settings.OnAllPDU,

		OnReceivingError: settings.OnReceivingError,

		OnClosed: func(state State) {
			switch state {
			case InvalidStreaming, UnbindClosing:
				// also close output
				_ = t.out.close(ExplicitClosing)

				if t.settings.OnClosed != nil {
					t.settings.OnClosed(state)
				}
			default:
				return
			}
		},

		RequestWindowConfig: settings.RequestWindowConfig,

		response: func(p pdu.PDU) {
			_ = t.Submit(p)
		},
	})
	return t
}

func (t *transceivable) start() {
	t.out.start()
	t.in.start()
}

// SystemID returns tagged SystemID which is attached with bind_resp from SMSC.
func (t *transceivable) SystemID() string {
	return t.conn.systemID
}

// Close transceiver and stop underlying daemons.
func (t *transceivable) Close() (err error) {
	if atomic.CompareAndSwapInt32(&t.aliveState, Alive, Closed) {
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
func (t *transceivable) Submit(p pdu.PDU) error {
	return t.out.Submit(p)
}

func (t *transceivable) GetWindowSize() int {
	return t.out.GetWindowSize()
}
