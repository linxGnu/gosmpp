package gosmpp

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/linxGnu/gosmpp/pdu"
)

var (
	// ErrConnectionClosing indicates transmitter is closing. Can not send any PDU.
	ErrConnectionClosing = fmt.Errorf("connection is closing, can not send PDU to SMSC")
)

type transmittable struct {
	conn *Connection

	wg    sync.WaitGroup
	input chan pdu.PDU

	settings Settings

	lock   sync.RWMutex
	closed bool
}

func newTransmittable(conn *Connection, settings Settings) *transmittable {
	t := &transmittable{
		settings: settings,
		conn:     conn,
		input:    make(chan pdu.PDU, 1),
	}

	return t
}

func (t *transmittable) close(state State) (err error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if !t.closed {
		t.closed = true

		// notify daemon
		close(t.input)

		// wait daemon
		t.wg.Wait()

		// try to send unbind
		_, _ = t.write(pdu.NewUnbind())

		// close connection
		if state != StoppingProcessOnly {
			err = t.conn.Close()
		}

		// notify transmitter closed
		if t.settings.OnClosed != nil {
			t.settings.OnClosed(state)
		}
	}

	return
}

func (t *transmittable) closing(state State) {
	go func() {
		_ = t.close(state)
	}()
}

// Submit a PDU.
func (t *transmittable) Submit(p pdu.PDU) (err error) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	if !t.closed {
		t.input <- p
	} else {
		err = ErrConnectionClosing
	}

	return
}

func (t *transmittable) start() {
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

func (t *transmittable) loop() {
	for p := range t.input {
		if p != nil {
			n, err := t.write(p)
			if t.check(p, n, err) {
				return
			}
		}
	}
}

func (t *transmittable) loopWithEnquireLink() {
	ticker := time.NewTicker(t.settings.EnquireLink)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			eqp := pdu.NewEnquireLink()
			n, err := t.write(eqp)
			if t.check(eqp, n, err) {
				return
			}

		case p, ok := <-t.input:
			if !ok {
				return
			}

			if p != nil {
				n, err := t.write(p)
				if t.check(p, n, err) {
					return
				}
			}
		}
	}
}

// check error and do closing if need
func (t *transmittable) check(p pdu.PDU, n int, err error) (closing bool) {
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
func (t *transmittable) write(p pdu.PDU) (n int, err error) {
	if t.settings.WriteTimeout > 0 {
		err = t.conn.SetWriteTimeout(t.settings.WriteTimeout)
	}

	if err == nil {
		n, err = t.conn.WritePDU(p)
	}

	return
}
