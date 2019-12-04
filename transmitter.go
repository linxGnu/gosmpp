package gosmpp

import (
	"context"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/linxGnu/gosmpp/pdu"
)

// TransmitListener is listener for transmitter.
type TransmitListener struct {
	// WriteTimeout is timeout/deadline for writting.
	WriteTimeout time.Duration

	// EnquireLink periodically sends EnquireLink to SMSC.
	// 0 duration means disable auto enquire link.
	EnquireLink time.Duration

	// OnError handles fail to submit pdu with along error.
	OnError func(pdu.PDU, error)

	// OnClosed handles transmitter `closed` event
	// with processing pdu and the specific error that causes transmitter closed.
	OnClosed func()
}

type transmitter struct {
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	listener TransmitListener

	dedicated bool
	conn      Connection

	input chan pdu.PDU
	state int32
}

// NewTransmitter returns new transmitter.
func NewTransmitter(conn Connection, listener TransmitListener) Transmitter {
	t := &transmitter{
		listener: listener,
		conn:     conn,
		input:    make(chan pdu.PDU),
	}
	t.ctx, t.cancel = context.WithCancel(context.Background())

	// start transmitter daemon(s)
	t.start()

	return t
}

// Close transmitter and stop underlying daemons.
func (t *transmitter) Close() (err error) {
	if atomic.CompareAndSwapInt32(&t.state, 0, 1) {
		// cancel context to notify stop
		t.cancel()

		// wait daemons
		t.wg.Wait()

		if t.conn.Dedicated {
			err = t.conn.Conn.Close()
		}

		// notify transmitter closed
		if t.listener.OnClosed != nil {
			t.listener.OnClosed()
		}
	}
	return
}

func (t *transmitter) close() {
	go func() {
		_ = t.Close()
	}()
}

// Write submits a pdu.
func (t *transmitter) Write(p pdu.PDU) (err error) {
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
	if t.listener.EnquireLink > 0 {
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

// pdu loop processing
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

// pdu loop processing with enquire link support
func (t *transmitter) loopWithEnquireLink() {
	ticker := time.NewTicker(t.listener.EnquireLink)

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

	if t.listener.OnError != nil {
		t.listener.OnError(p, err)
	}

	if n == 0 {
		if nErr, ok := err.(net.Error); ok {
			closing = nErr.Timeout() || !nErr.Temporary()
		}
	} else {
		closing = true // force closing
	}

	if closing {
		t.close() // start closing
	}

	return
}

// low level writing
func (t *transmitter) write(v []byte) (n int, err error) {
	hasTimeout := t.listener.WriteTimeout > 0

	if hasTimeout {
		t.conn.Conn.SetWriteDeadline(time.Now().Add(t.listener.WriteTimeout))
	}

	if n, err = t.conn.Conn.Write(v); err != nil && n == 0 {
		// retry again with double timeout
		if hasTimeout {
			t.conn.Conn.SetWriteDeadline(time.Now().Add(t.listener.WriteTimeout << 1))
		}

		n, err = t.conn.Conn.Write(v)
	}

	return
}
