package gosmpp

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/linxGnu/gosmpp/pdu"
)

type receivable struct {
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	settings     Settings
	conn         *Connection
	aliveState   int32
	requestStore RequestStore
}

func newReceivable(conn *Connection, settings Settings, requestStore RequestStore) *receivable {
	r := &receivable{
		settings:     settings,
		conn:         conn,
		requestStore: requestStore,
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
		defer t.wg.Done()
		t.loop()
	}()
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
		if err != nil {
			if atomic.LoadInt32(&t.aliveState) == Alive {
				if t.settings.OnReceivingError != nil {
					t.settings.OnReceivingError(err)
				}
				t.closing(InvalidStreaming)
			}
			return
		}

		var closeOnUnbind bool
		if p != nil {
			if t.settings.WindowedRequestTracking != nil && t.settings.OnExpectedPduResponse != nil {
				closeOnUnbind = t.handleWindowPdu(p)
			} else if t.settings.OnAllPDU != nil {
				closeOnUnbind = t.handleAllPdu(p)
			} else {
				closeOnUnbind = t.handleOrClose(p)
			}
			if closeOnUnbind {
				t.closing(UnbindClosing)
			}
		}

	}
}

func (t *receivable) handleWindowPdu(p pdu.PDU) (closing bool) {
	if t.settings.WindowedRequestTracking != nil && t.settings.OnExpectedPduResponse != nil && p != nil {
		// This case must match the same request item list in transmittable write func
		switch pp := p.(type) {
		case *pdu.CancelSMResp,
			*pdu.DataSMResp,
			*pdu.DeliverSMResp,
			*pdu.EnquireLinkResp,
			*pdu.QuerySMResp,
			*pdu.ReplaceSMResp,
			*pdu.SubmitMultiResp,
			*pdu.SubmitSMResp:
			if t.settings.OnExpectedPduResponse != nil {
				ctx, cancelFunc := context.WithTimeout(context.Background(), t.settings.StoreAccessTimeOut*time.Millisecond)
				defer cancelFunc()
				request, ok := t.requestStore.Get(ctx, p.GetSequenceNumber())
				if ok {
					_ = t.requestStore.Delete(ctx, p.GetSequenceNumber())
					response := Response{
						PDU:             p,
						OriginalRequest: request,
					}
					t.settings.OnExpectedPduResponse(response)
				} else if t.settings.OnUnexpectedPduResponse != nil {
					t.settings.OnUnexpectedPduResponse(p)
				}
			}
		case *pdu.EnquireLink:
			if t.settings.EnableAutoRespond {
				t.settings.response(pp.GetResponse())
			} else if t.settings.OnReceivedPduRequest != nil {
				r, _ := t.settings.OnReceivedPduRequest(p)
				t.settings.response(r)

			}
		case *pdu.Unbind:
			if t.settings.EnableAutoRespond {
				t.settings.response(pp.GetResponse())

				// wait to send response before closing
				time.Sleep(50 * time.Millisecond)
				closing = true
			} else if t.settings.OnReceivedPduRequest != nil {
				r, closeBind := t.settings.OnReceivedPduRequest(p)
				t.settings.response(r)
				if closeBind {
					time.Sleep(50 * time.Millisecond)
					closing = true
				}
			}
		default:
			if t.settings.OnReceivedPduRequest != nil {
				r, closeBind := t.settings.OnReceivedPduRequest(p)
				t.settings.response(r)
				if closeBind {
					time.Sleep(50 * time.Millisecond)
					closing = true
				}
			}
		}
	}
	return
}

func (t *receivable) handleAllPdu(p pdu.PDU) (closing bool) {
	if t.settings.OnAllPDU != nil && p != nil {
		r, closeBind := t.settings.OnAllPDU(p)
		t.settings.response(r)
		if closeBind {
			time.Sleep(50 * time.Millisecond)
			closing = true
		}
	}
	return
}

func (t *receivable) handleOrClose(p pdu.PDU) (closing bool) {
	if p != nil {
		switch pp := p.(type) {
		case *pdu.EnquireLink:
			t.settings.response(pp.GetResponse())

		case *pdu.Unbind:
			t.settings.response(pp.GetResponse())
			// wait to send response before closing
			time.Sleep(50 * time.Millisecond)

			closing = true

		default:
			var responded bool
			if p.CanResponse() {
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
