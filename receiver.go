package gosmpp

import (
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/linxGnu/gosmpp/pdu"
)

// ReceiveSettings is event listener for Receiver.
type ReceiveSettings struct {
	// OnPDU handles received PDU from SMSC.
	OnPDU func(pdu.PDU)

	// OnReceivingError notifies happened error while reading PDU
	// from SMSC.
	OnReceivingError func(error)

	// OnClosed notifies `closed` event due to State.
	OnClosed func(State)

	response func(pdu.PDU)
}

type receiver struct {
	wg       sync.WaitGroup
	settings ReceiveSettings
	conn     net.Conn
	state    int32
}

// NewReceiver returns new Receiver.
func NewReceiver(conn net.Conn, settings ReceiveSettings) Receiver {
	return newReceiver(conn, settings, true)
}

func newReceiver(conn net.Conn, settings ReceiveSettings, startDaemon bool) *receiver {
	r := &receiver{
		settings: settings,
		conn:     conn,
	}

	// start receiver daemon(s)
	if startDaemon {
		r.start()
	}

	return r
}

// Close receiver, close connection and stop underlying daemons.
func (t *receiver) Close() (err error) {
	return t.close(ExplicitClosing)
}

func (t *receiver) close(state State) (err error) {
	if atomic.CompareAndSwapInt32(&t.state, 0, 1) {
		// close connection to notify daemons to stop
		err = t.conn.Close()

		// wait daemons
		t.wg.Wait()

		// notify receiver closed
		if t.settings.OnClosed != nil {
			t.settings.OnClosed(state)
		}
	}
	return
}

func (t *receiver) closing(state State) {
	go func() {
		_ = t.close(state)
	}()
}

func (t *receiver) start() {
	t.wg.Add(1)
	go func() {
		t.loop()
		t.wg.Done()
	}()
}

// check error and do closing if need
func (t *receiver) check(err error) (closing bool) {
	if err == nil {
		return
	}

	if t.settings.OnReceivingError != nil {
		t.settings.OnReceivingError(err)
	}

	closing = true
	return
}

// PDU loop processing
func (t *receiver) loop() {
	for {
		p, err := pdu.Parse(t.conn)

		closeOnError := t.check(err)
		if closeOnError || t.handleOrClose(p) {
			if closeOnError {
				t.closing(InvalidStreaming)
			}
			return
		}
	}
}

func (t *receiver) handleOrClose(p pdu.PDU) (closing bool) {
	if p != nil {
		switch pp := p.(type) {
		case *pdu.EnquireLink:
			if t.settings.response != nil {
				t.settings.response(pp.GetResponse())
			}
			return

		case *pdu.Unbind:
			if t.settings.response != nil {
				t.settings.response(pp.GetResponse())

				// wait to send response before closing
				time.Sleep(100 * time.Millisecond)
			}

			closing = true
			t.closing(UnbindClosing)
			return

		default:
			if t.settings.OnPDU != nil {
				t.settings.OnPDU(p)
			}
		}
	}
	return
}