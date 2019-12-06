package gosmpp

import (
	"context"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/linxGnu/gosmpp/pdu"
)

// TransmitSettings is listener for transmitter.
type TransmitSettings struct {
	// WriteTimeout is timeout/deadline for submitting PDU.
	WriteTimeout time.Duration

	// EnquireLink periodically sends EnquireLink to SMSC.
	// Zero duration means disable auto enquire link.
	EnquireLink time.Duration

	// OnSubmitError notifies fail-to-submit PDU with along error.
	OnSubmitError func(pdu.PDU, error)

	// OnClosed notifies `closed` event due to State.
	OnClosed func(State)
}

type transmitter struct {
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	settings TransmitSettings
	conn     net.Conn
	input    chan pdu.PDU
	state    int32
}

// NewTransmitter returns new Transmitter.
func NewTransmitter(conn net.Conn, settings TransmitSettings) Transmitter {
	return newTransmitter(conn, settings, true)
}

func newTransmitter(conn net.Conn, settings TransmitSettings, startDaemon bool) *transmitter {
	t := &transmitter{
		settings: settings,
		conn:     conn,
		input:    make(chan pdu.PDU, 1),
	}
	t.ctx, t.cancel = context.WithCancel(context.Background())

	// start transmitter daemon(s)
	if startDaemon {
		t.start()
	}

	return t
}

// Close transmitter and stop underlying daemons.
func (t *transmitter) Close() (err error) {
	return t.close(ExplicitClosing)
}

func (t *transmitter) close(state State) (err error) {
	if atomic.CompareAndSwapInt32(&t.state, 0, 1) {
		// cancel context to notify stop
		t.cancel()

		// wait daemons
		t.wg.Wait()

		// close connection
		err = t.conn.Close()

		// notify transmitter closed
		if t.settings.OnClosed != nil {
			t.settings.OnClosed(state)
		}
	}
	return
}

func (t *transmitter) closing(state State) {
	go func() {
		_ = t.Close()
	}()
}

// Submit a PDU.
func (t *transmitter) Submit(p pdu.PDU) (err error) {
	select {
	case <-t.ctx.Done():
		err = t.ctx.Err()
		return

	case t.input <- p:
		return
	}
}

func (t *transmitter) start() {
	t.wg.Add(1)
	if t.settings.EnquireLink > 0 {
		go func() {
			t.loopWithEnquireLink()
			t.wg.Done()
		}()
	} else {
		go func() {
			t.loop()
			t.wg.Done()
		}()
	}
}

// PDU loop processing
func (t *transmitter) loop() {
	for {
		select {
		case <-t.ctx.Done():
			return

		case p := <-t.input:
			if p != nil {
				n, err := t.write(marshal(p))
				if t.check(p, n, err) {
					return
				}
			}
		}
	}
}

// PDU loop processing with enquire link support
func (t *transmitter) loopWithEnquireLink() {
	ticker := time.NewTicker(t.settings.EnquireLink)

	// enquireLink payload
	eqp := pdu.NewEnquireLink()
	enquireLink := marshal(eqp)

	for {
		select {
		case <-t.ctx.Done():
			return

		case <-ticker.C:
			n, err := t.write(enquireLink)
			if t.check(eqp, n, err) {
				return
			}

		case p := <-t.input:
			if p != nil {
				n, err := t.write(marshal(p))
				if t.check(p, n, err) {
					return
				}
			}
		}
	}
}

// check error and do closing if need
func (t *transmitter) check(p pdu.PDU, n int, err error) (closing bool) {
	if err == nil {
		return
	}

	if t.settings.OnSubmitError != nil {
		t.settings.OnSubmitError(p, err)
	}

	if n == 0 {
		if nErr, ok := err.(net.Error); ok {
			closing = nErr.Timeout() || !nErr.Temporary()
		}
	} else {
		closing = true // force closing
	}

	if closing {
		t.closing(ConnectionIssue) // start closing
	}

	return
}

// low level writing
func (t *transmitter) write(v []byte) (n int, err error) {
	hasTimeout := t.settings.WriteTimeout > 0

	if hasTimeout {
		t.conn.SetWriteDeadline(time.Now().Add(t.settings.WriteTimeout))
	}

	if n, err = t.conn.Write(v); err != nil && n == 0 {
		// retry again with double timeout
		if hasTimeout {
			t.conn.SetWriteDeadline(time.Now().Add(t.settings.WriteTimeout << 1))
		}

		n, err = t.conn.Write(v)
	}

	return
}
