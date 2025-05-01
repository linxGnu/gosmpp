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
	rx := &receivable{
		settings:     settings,
		conn:         conn,
		requestStore: requestStore,
	}
	rx.ctx, rx.cancel = context.WithCancel(context.Background())

	return rx
}

func (rx *receivable) close() {
	if atomic.CompareAndSwapInt32(&rx.aliveState, Alive, Closed) {
		// cancel to notify stop
		rx.cancel()

		// set read deadline for current blocking read
		_ = rx.conn.SetReadDeadline(time.Now().Add(200 * time.Millisecond))

		// wait daemons
		rx.wg.Wait()
	}
	return
}

func (rx *receivable) closing(state State) {
	// notify transceiver of closing
	go func() {
		rx.settings.OnClosed(state)
	}()
}

func (rx *receivable) start() {
	rx.wg.Add(1)
	go func() {
		defer rx.wg.Done()
		rx.loop()
	}()
}

func (rx *receivable) loop() {
	var err error
	for {
		select {
		case <-rx.ctx.Done():
			return
		default:
		}

		// read pdu from conn
		var p pdu.PDU
		if err = rx.conn.SetReadTimeout(rx.settings.ReadTimeout); err == nil {
			p, err = pdu.Parse(rx.conn)
		}
		if err != nil {
			if atomic.LoadInt32(&rx.aliveState) == Alive {
				if rx.settings.OnReceivingError != nil {
					rx.settings.OnReceivingError(err)
				}
				rx.closing(InvalidStreaming)
			}
			return
		}

		var closeOnUnbind bool
		if p != nil {
			if rx.settings.WindowedRequestTracking != nil && rx.settings.OnExpectedPduResponse != nil {
				closeOnUnbind = rx.handleWindowPdu(p)
			} else if rx.settings.OnAllPDU != nil {
				closeOnUnbind = rx.handleAllPdu(p)
			} else {
				closeOnUnbind = rx.handleOrClose(p)
			}
			if closeOnUnbind {
				rx.closing(UnbindClosing)
			}
		}

	}
}

func (rx *receivable) handleWindowPdu(p pdu.PDU) (closing bool) {
	if rx.settings.WindowedRequestTracking != nil && rx.settings.OnExpectedPduResponse != nil && p != nil {
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
			if rx.settings.OnExpectedPduResponse != nil {
				ctx, cancelFunc := context.WithTimeout(context.Background(), rx.settings.StoreAccessTimeOut)
				defer cancelFunc()
				request, ok := rx.requestStore.Get(ctx, p.GetSequenceNumber())
				if ok {
					_ = rx.requestStore.Delete(ctx, p.GetSequenceNumber())
					response := Response{
						PDU:             p,
						OriginalRequest: request,
					}
					rx.settings.OnExpectedPduResponse(response)
				} else if rx.settings.OnUnexpectedPduResponse != nil {
					rx.settings.OnUnexpectedPduResponse(p)
				}
			}
		case *pdu.EnquireLink:
			if rx.settings.EnableAutoRespond {
				rx.settings.response(pp.GetResponse())
			} else if rx.settings.OnReceivedPduRequest != nil {
				r, _ := rx.settings.OnReceivedPduRequest(p)
				rx.settings.response(r)

			}
		case *pdu.Unbind:
			if rx.settings.EnableAutoRespond {
				rx.settings.response(pp.GetResponse())

				// wait to send response before closing
				time.Sleep(50 * time.Millisecond)
				closing = true
			} else if rx.settings.OnReceivedPduRequest != nil {
				r, closeBind := rx.settings.OnReceivedPduRequest(p)
				rx.settings.response(r)
				if closeBind {
					time.Sleep(50 * time.Millisecond)
					closing = true
				}
			}
		default:
			if rx.settings.OnReceivedPduRequest != nil {
				r, closeBind := rx.settings.OnReceivedPduRequest(p)
				rx.settings.response(r)
				if closeBind {
					time.Sleep(50 * time.Millisecond)
					closing = true
				}
			}
		}
	}
	return
}

func (rx *receivable) handleAllPdu(p pdu.PDU) (closing bool) {
	if rx.settings.OnAllPDU != nil && p != nil {
		r, closeBind := rx.settings.OnAllPDU(p)
		rx.settings.response(r)
		if closeBind {
			time.Sleep(50 * time.Millisecond)
			closing = true
		}
	}
	return
}

func (rx *receivable) handleOrClose(p pdu.PDU) (closing bool) {
	if p != nil {
		switch pp := p.(type) {
		case *pdu.EnquireLink:
			rx.settings.response(pp.GetResponse())

		case *pdu.Unbind:
			rx.settings.response(pp.GetResponse())
			// wait to send response before closing
			time.Sleep(50 * time.Millisecond)

			closing = true

		default:
			var responded bool
			if p.CanResponse() {
				rx.settings.response(p.GetResponse())
				responded = true
			}

			if rx.settings.OnPDU != nil {
				rx.settings.OnPDU(p, responded)
			}
		}
	}
	return
}
