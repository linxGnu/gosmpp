package gosmpp

import (
	"context"
	cmap "github.com/orcaman/concurrent-map/v2"
	"strconv"
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
	window     cmap.ConcurrentMap[string, Request]
}

func newReceivable(conn *Connection, window cmap.ConcurrentMap[string, Request], settings Settings) *receivable {
	r := &receivable{
		settings: settings,
		conn:     conn,
		window:   window,
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
	if t.settings.WindowPDUHandlerConfig != nil && t.settings.PduExpireTimeOut > 0 && t.settings.ExpireCheckTimer > 0 {
		go func() {
			t.loopWithVerifyExpiredPdu()
			t.wg.Done()
		}()
	} else {
		go func() {
			t.loop()
			t.wg.Done()
		}()
	}
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

func (t *receivable) loopWithVerifyExpiredPdu() {
	ticker := time.NewTicker(t.settings.ExpireCheckTimer)
	defer func() {
		ticker.Stop()
	}()
	var err error
	for {
		select {
		case <-t.ctx.Done():
			return
		case <-ticker.C:
			for request := range t.window.IterBuffered() {
				if time.Since(request.Val.TImeSent) > t.settings.PduExpireTimeOut {
					t.window.Remove(request.Key)
					if t.settings.OnExpiredPduRequest != nil {
						t.settings.OnExpiredPduRequest(request.Val.PDU)
					}
				}
			}
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
		if t.settings.WindowPDUHandlerConfig != nil {
			if t.settings.OnExpectedPduResponse != nil {
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
						sequence := strconv.Itoa(int(p.GetSequenceNumber()))
						request, ok := t.window.Get(sequence)
						//request, found := t.conn.window[p.GetSequenceNumber()]

						if ok {
							t.window.Remove(sequence)
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
					if t.settings.EnableAutoRespond && t.settings.response != nil {
						t.settings.response(pp.GetResponse())
					} else {
						if t.settings.OnReceivedPduRequest != nil {
							r, closeBind := t.settings.OnReceivedPduRequest(p)
							t.settings.response(r)
							if closeBind {
								time.Sleep(50 * time.Millisecond)
								closing = true
								t.closing(UnbindClosing)
							}
							return
						}
					}

				case *pdu.Unbind:
					if t.settings.EnableAutoRespond && t.settings.response != nil {
						t.settings.response(pp.GetResponse())

						// wait to send response before closing
						time.Sleep(50 * time.Millisecond)
					} else {
						if t.settings.OnReceivedPduRequest != nil {
							r, closeBind := t.settings.OnReceivedPduRequest(p)
							t.settings.response(r)
							if closeBind {
								time.Sleep(50 * time.Millisecond)
								closing = true
								t.closing(UnbindClosing)
							}
							return
						}
					}

					closing = true
					t.closing(UnbindClosing)
				default:
					if t.settings.OnReceivedPduRequest != nil {
						r, closeBind := t.settings.OnReceivedPduRequest(p)
						t.settings.response(r)
						if closeBind {
							time.Sleep(50 * time.Millisecond)
							closing = true
							t.closing(UnbindClosing)
						}
						return
					}
				}
			}
			return
		}
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
