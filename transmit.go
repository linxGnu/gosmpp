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

type transmitable struct {
	ctx    context.Context
	cancel context.CancelFunc

	conn *Connection

	wg    sync.WaitGroup
	input chan pdu.PDU

	settings Settings

	lock  sync.RWMutex
	state int32
}

func newTransmitable(conn *Connection, settings Settings) *transmitable {
	t := &transmitable{
		settings: settings,
		conn:     conn,
		input:    make(chan pdu.PDU, 1),
	}
	t.ctx, t.cancel = context.WithCancel(context.Background())

	return t
}

func (t *transmitable) close(state State) (err error) {
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
		_, _ = t.write(marshal(pdu.NewUnbind()))

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

func (t *transmitable) closing(state State) {
	go func() {
		_ = t.close(state)
	}()
}

// Submit a PDU.
func (t *transmitable) Submit(p pdu.PDU) (err error) {
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

func (t *transmitable) start() {
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
func (t *transmitable) loop() {
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
func (t *transmitable) loopWithEnquireLink() {
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
func (t *transmitable) check(p pdu.PDU, n int, err error) (closing bool) {
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
func (t *transmitable) write(v []byte) (n int, err error) {
	if t.settings.WriteTimeout > 0 {
		err = t.conn.SetWriteTimeout(t.settings.WriteTimeout)
	}

	if err == nil {
		n, err = t.conn.Write(v)
	}

	return
}
