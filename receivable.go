package gosmpp

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/linxGnu/gosmpp/pdu"
)

type receivable struct {
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	settings   Settings
	conn       *Connection
	aliveState int32
}

func newReceivable(conn *Connection, settings Settings) *receivable {
	r := &receivable{
		settings: settings,
		conn:     conn,
	}
	r.ctx, r.cancel = context.WithCancel(context.Background())

	return r
}

func (t *receivable) close(state State) (err error) {
	if atomic.CompareAndSwapInt32(&t.aliveState, Alive, Closed) {
		// cancel to notify stop
		t.cancel()

		// set read deadline for current blocking read
		_ = t.conn.SetReadDeadline(time.Now().Add(200 * time.Millisecond))

		// wait daemons
		t.wg.Wait()

		// close connection to notify daemons to stop
		if state != StoppingProcessOnly {
			err = t.conn.Close()
		}

		// notify receiver closed
		if t.settings.OnClosed != nil {
			t.settings.OnClosed(state)
		}
	}
	return
}

func (t *receivable) closing(state State) {
	go func() {
		_ = t.close(state)
	}()
}

func (t *receivable) start() {
	t.wg.Add(1)
	go func() {
		t.loop()
		t.wg.Done()
	}()
}

// check error and do closing if need
func (t *receivable) check(err error) (closing bool) {
	if err == nil {
		return
	}

	if t.settings.OnReceivingError != nil {
		t.settings.OnReceivingError(err)
	}

	closing = true
	return
}

func (t *receivable) loop() {
	var err error
	for {
		select {
		case <-t.ctx.Done():
			return
		default:
		}

		// read pdu from conn
		var p pdu.PDU
		if err = t.conn.SetReadTimeout(t.settings.ReadTimeout); err == nil {
			p, err = pdu.Parse(t.conn)
		}

		// check error
		if closeOnError := t.check(err); closeOnError || t.handleOrClose(p) {
			if closeOnError {
				t.closing(InvalidStreaming)
			}
			return
		}
	}
}

func (t *receivable) handleOrClose(p pdu.PDU) (closing bool) {
	if p != nil {
		if t.settings.OnAllPDU != nil {
			r, closeBind := t.settings.OnAllPDU(p)
			t.settings.response(r)
			if closeBind {
				time.Sleep(50 * time.Millisecond)
				closing = true
				t.closing(UnbindClosing)
			}
			return
		}

		switch pp := p.(type) {
		case *pdu.EnquireLink:
			if t.settings.response != nil {
				t.settings.response(pp.GetResponse())
			}

		case *pdu.Unbind:
			if t.settings.response != nil {
				t.settings.response(pp.GetResponse())

				// wait to send response before closing
				time.Sleep(50 * time.Millisecond)
			}

			closing = true
			t.closing(UnbindClosing)

		default:
			var responded bool
			if p.CanResponse() && t.settings.response != nil {
				t.settings.response(p.GetResponse())
				responded = true
			}

			if t.settings.OnPDU != nil {
				t.settings.OnPDU(p, responded)
			}
		}
	}
	return
}
