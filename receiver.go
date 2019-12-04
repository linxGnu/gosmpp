package gosmpp

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/linxGnu/gosmpp/pdu"
)

// ReceiveListener is event listener for Receiver.
type ReceiveListener struct {
	// OnRecevingError handles fail to submit pdu with along error.
	OnRecevingError func(error)

	// OnClosed handles receiver `closed` event
	// with processing pdu and the specific error that causes receiver closed.
	OnClosed func()

	// OnPDU handles received pdu from SMSC.
	OnPDU func(pdu.PDU)

	// Response indicates that invoker should response specific provided pdu to SMSC.
	Response func(pdu.PDU)
}

type receiver struct {
	wg       sync.WaitGroup
	listener ReceiveListener
	conn     Connection
	state    int32
}

// Close receiver, close connection and stop underlying daemons.
func (t *receiver) Close() (err error) {
	if atomic.CompareAndSwapInt32(&t.state, 0, 1) {
		// close connection to notify daemons to stop
		if t.conn.Dedicated {
			err = t.conn.Conn.Close()
		}

		// wait daemons
		t.wg.Wait()

		// notify receiver closed
		if t.listener.OnClosed != nil {
			t.listener.OnClosed()
		}
	}
	return
}

func (t *receiver) close() {
	go func() {
		_ = t.Close()
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

	if t.listener.OnRecevingError != nil {
		t.listener.OnRecevingError(err)
	}

	closing = true
	return
}

// pdu loop processing
func (t *receiver) loop() {
	for {
		p, err := pdu.Parse(t.conn.Conn)
		if t.check(err) {
			t.close()
			return
		}
		t.handle(p)
	}
}

func (t *receiver) handle(p pdu.PDU) {
	if p != nil {
		switch pp := p.(type) {
		case *pdu.EnquireLink:
			if t.listener.Response != nil {
				t.listener.Response(pp.GetResponse())
			}
			return

		case *pdu.Unbind:
			if t.listener.Response != nil {
				t.listener.Response(pp.GetResponse())

				// wait to send response before closing
				time.Sleep(100 * time.Millisecond)
			}

			t.close()
			return

		default:
			if t.listener.OnPDU != nil {
				t.listener.OnPDU(p)
			}
		}
	}
	return
}
