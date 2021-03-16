package gosmpp

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/linxGnu/gosmpp/pdu"
)

var (
	// ErrTransmitterClosing indicates transmitter is closing. Can not send any PDU.
	ErrTransmitterClosing = fmt.Errorf("Transmitter is closing. Can not send PDU to SMSC")
)

// TransmitSettings is listener for transmitter.
type TransmitSettings struct {
	// Timeout is timeout/deadline for submitting PDU.
	Timeout time.Duration

	// EnquireLink periodically sends EnquireLink to SMSC.
	// The duration must not be smaller than 1 minute.
	//
	// Zero duration disables auto enquire link.
	EnquireLink time.Duration

	// OnPDU handles received PDU from SMSC.
	//
	// `Responded` flag indicates this pdu is responded automatically,
	// no manual respond needed.
	OnPDU PDUCallback

	// OnReceivingError notifies happened error while reading PDU
	// from SMSC.
	OnReceivingError ErrorCallback

	// OnSubmitError notifies fail-to-submit PDU with along error.
	OnSubmitError PDUErrorCallback

	// OnRebindingError notifies error while rebinding.
	OnRebindingError ErrorCallback

	// OnClosed notifies `closed` event due to State.
	OnClosed ClosedCallback
}

type transmitter struct {
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	settings TransmitSettings
	conn     *Connection
	input    chan pdu.PDU
	lock     sync.RWMutex
	state    int32
}

func newTransmitter(conn *Connection, settings TransmitSettings) *transmitter {
	t := &transmitter{
		settings: settings,
		conn:     conn,
		input:    make(chan pdu.PDU, 1),
	}
	t.ctx, t.cancel = context.WithCancel(context.Background())

	return t
}

// SystemID returns tagged SystemID, returned from bind_resp from SMSC.
func (t *transmitter) SystemID() string {
	return t.conn.systemID
}

// Close transmitter and stop underlying daemons.
func (t *transmitter) Close() (err error) {
	return t.close(ExplicitClosing)
}

func (t *transmitter) close(state State) (err error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.state == 0 {
		// don't receive anymore SubmitSM
		t.cancel()

		// notify daemon
		close(t.input)

		// wait daemon
		t.wg.Wait()

		// try to send unbind
		if t.settings.Timeout > 0 {
			_ = t.conn.SetWriteTimeout(t.settings.Timeout)
		}
		_, _ = t.conn.Write(marshal(pdu.NewUnbind()))

		// close connection
		if state != StoppingProcessOnly {
			err = t.conn.Close()
		}

		// notify transmitter closed
		if t.settings.OnClosed != nil {
			t.settings.OnClosed(state)
		}

		t.state = 1
	}

	return
}

func (t *transmitter) closing(state State) {
	go func() {
		_ = t.close(state)
	}()
}

// Submit a PDU.
func (t *transmitter) Submit(p pdu.PDU) (err error) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	if t.state == 0 {
		select {
		case <-t.ctx.Done():
			err = t.ctx.Err()

		case t.input <- p:
		}
	} else {
		err = ErrTransmitterClosing
	}

	return
}

func (t *transmitter) start(receiving bool) {
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

	if receiving {
		go t.loopReceiving()
	}
}

// PDU loop processing
func (t *transmitter) loop() {
	for p := range t.input {
		if p != nil {
			n, err := t.write(marshal(p))
			if t.check(p, n, err) {
				return
			}
		}
	}
}

// PDU loop processing with enquire link support
func (t *transmitter) loopWithEnquireLink() {
	ticker := time.NewTicker(t.settings.EnquireLink)
	defer ticker.Stop()

	// enquireLink payload
	eqp := pdu.NewEnquireLink()
	enquireLink := marshal(eqp)

	for {
		select {
		case <-ticker.C:
			n, err := t.write(enquireLink)
			if t.check(eqp, n, err) {
				return
			}

		case p, ok := <-t.input:
			if !ok {
				return
			}

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
		} else {
			closing = true
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
	hasTimeout := t.settings.Timeout > 0

	if hasTimeout {
		err = t.conn.SetWriteTimeout(t.settings.Timeout)
	}

	if err == nil {
		if n, err = t.conn.Write(v); err != nil &&
			n == 0 &&
			hasTimeout &&
			t.conn.SetWriteTimeout(t.settings.Timeout<<1) == nil {
			// retry again with double timeout
			n, err = t.conn.Write(v)
		}
	}

	return
}

// PDU loop processing
func (t *transmitter) loopReceiving() {
	checkErr := func(err error) (closing bool) {
		if err == nil {
			return
		}

		if t.settings.OnReceivingError != nil {
			t.settings.OnReceivingError(err)
		}

		closing = true
		return
	}

	for {
		select {
		case <-t.ctx.Done():
			return
		default:
		}

		// read pdu from conn
		p, err := pdu.Parse(t.conn)

		// check error
		if closeOnError := checkErr(err); closeOnError || t.handleOrClose(p) {
			if closeOnError {
				t.closing(InvalidStreaming)
			}
			return
		}
	}
}

func (t *transmitter) handleOrClose(p pdu.PDU) (closing bool) {
	if p != nil {
		switch pp := p.(type) {
		case *pdu.EnquireLink:
			_ = t.Submit(pp.GetResponse())

		case *pdu.Unbind:
			_ = t.Submit(pp.GetResponse())

			// wait to send response before closing
			time.Sleep(50 * time.Millisecond)

			closing = true
			t.closing(UnbindClosing)

		default:
			var responded bool
			if p.CanResponse() {
				_ = t.Submit(p.GetResponse())
			}

			if t.settings.OnPDU != nil {
				t.settings.OnPDU(p, responded)
			}
		}
	}
	return
}
