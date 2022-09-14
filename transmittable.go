package gosmpp

import (
	"fmt"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/linxGnu/gosmpp/pdu"
)

var (
	// ErrConnectionClosing indicates transmitter is closing. Can not send any PDU.
	ErrConnectionClosing = fmt.Errorf("connection is closing, can not send PDU to SMSC")
)

type transmittable struct {
	settings Settings

	wg    sync.WaitGroup
	input chan pdu.PDU

	conn *Connection

	aliveState   int32
	pendingWrite int32
}

func newTransmittable(conn *Connection, settings Settings) *transmittable {
	t := &transmittable{
		settings:     settings,
		conn:         conn,
		input:        make(chan pdu.PDU, 1),
		aliveState:   Alive,
		pendingWrite: 0,
	}

	return t
}

func (t *transmittable) close(state State) (err error) {
	if atomic.CompareAndSwapInt32(&t.aliveState, Alive, Closed) {
		for atomic.LoadInt32(&t.pendingWrite) != 0 {
			runtime.Gosched()
		}

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
	atomic.AddInt32(&t.pendingWrite, 1)

	if atomic.LoadInt32(&t.aliveState) == Alive {
		t.input <- p
	} else {
		err = ErrConnectionClosing
	}

	atomic.AddInt32(&t.pendingWrite, -1)
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

func (t *transmittable) drain() {
	for range t.input {
	}
}

func (t *transmittable) loop() {
	defer t.drain()

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
	defer func() {
		ticker.Stop()
		t.drain()
	}()

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
			closing = nErr.Timeout()
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
