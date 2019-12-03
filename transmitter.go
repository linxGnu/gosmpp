package gosmpp

import (
	"context"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/linxGnu/gosmpp/pdu"
)

// TransmitterConfig is configuration for transmitter.
type TransmitterConfig struct {
	// WriteTimeout is timeout/deadline for writting.
	WriteTimeout time.Duration

	// EnquireLink periodically sends EnquireLink to SMSC.
	// 0 duration means disable auto enquire link.
	EnquireLink time.Duration

	// OnError handles fail to submit pdu with along error.
	OnError func(pdu.PDU, error)

	// OnClosed handles transmitter `closed` event
	// with processing pdu and the specific error that causes transmitter closed.
	OnClosed func(pdu.PDU, error)
}

type transmitter struct {
	ctx    context.Context
	cancel context.CancelFunc

	wg sync.WaitGroup

	config TransmitterConfig

	shared bool
	conn   net.Conn

	input chan pdu.PDU

	state int32
}

// NewTransmitter returns new transmitter.
func NewTransmitter(ctx context.Context, config TransmitterConfig) Transmitter {
	return &transmitter{}
}

// Close transmitter and stop underlying daemons.
func (t *transmitter) Close() (err error) {
	if atomic.CompareAndSwapInt32(&t.state, 0, 1) {
		// cancel context to notify stop
		t.cancel()

		// wait daemons
		t.wg.Wait()

		if !t.shared {
			err = t.conn.Close()
		}
	}
	return
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

// write pdu daemon
func (t *transmitter) write() {
	if t.config.EnquireLink > 0 {
		t.loopWithEnquireLink()
	} else {
		t.loop()
	}
	t.wg.Done()
}

// pdu loop processing
func (t *transmitter) loop() {
	for {
		select {
		case <-t.ctx.Done():
			return

		case p := <-t.input:
			if p != nil {
				n, err := t._write(marshal(p))
				if t.check(p, n, err) {
					return
				}
			}
		}
	}
}

// pdu loop processing with enquire link support
func (t *transmitter) loopWithEnquireLink() {
	ticker := time.NewTicker(t.config.EnquireLink)

	// enquireLink payload
	eqp := pdu.NewEnquireLink()
	enquireLink := marshal(eqp)

	for {
		select {
		case <-t.ctx.Done():
			return

		case <-ticker.C:
			n, err := t._write(enquireLink)
			if t.check(eqp, n, err) {
				return
			}

		case p := <-t.input:
			if p != nil {
				n, err := t._write(marshal(p))
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

	if t.config.OnError != nil {
		t.config.OnError(p, err)
	}

	if n == 0 {
		if nErr, ok := err.(net.Error); ok {
			closing = nErr.Timeout() || !nErr.Temporary()
		}
	} else {
		closing = true // force closing
	}

	if closing {
		go closeTransmitter(t, p, err)
	}

	return
}

// low level writing
func (t *transmitter) _write(v []byte) (n int, err error) {
	hasTimeout := t.config.WriteTimeout > 0

	if hasTimeout {
		t.conn.SetWriteDeadline(time.Now().Add(t.config.WriteTimeout))
	}

	if n, err = t.conn.Write(v); err != nil && n == 0 {
		// retry again with double timeout
		if hasTimeout {
			t.conn.SetWriteDeadline(time.Now().Add(t.config.WriteTimeout << 1))
		}

		n, err = t.conn.Write(v)
	}

	return
}

func closeTransmitter(t *transmitter, processingPDU pdu.PDU, err error) {
	_ = t.Close()

	if t.config.OnClosed != nil {
		t.config.OnClosed(processingPDU, err)
	}
}
