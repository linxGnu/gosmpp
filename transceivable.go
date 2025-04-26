package gosmpp

import (
	"context"
	"errors"
	"github.com/linxGnu/gosmpp/pdu"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrWindowNotConfigured = errors.New("window settings not configured")
)

type transceivable struct {
	settings Settings

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	conn   *Connection
	in     *receivable
	out    *transmittable

	aliveState   int32
	requestStore RequestStore
}
type TransceivableOption func(session *Session)

func newTransceivable(conn *Connection, settings Settings, requestStore RequestStore) *transceivable {

	trx := &transceivable{
		settings:     settings,
		conn:         conn,
		requestStore: requestStore,
	}
	trx.ctx, trx.cancel = context.WithCancel(context.Background())

	trx.out = newTransmittable(conn, Settings{
		WriteTimeout: settings.WriteTimeout,

		EnquireLink: settings.EnquireLink,

		OnSubmitError: settings.OnSubmitError,

		OnClosed: func(state State) {
			switch state {
			case ConnectionIssue:
				// also close input
				_ = trx.in.close(ExplicitClosing)

				if trx.settings.OnClosed != nil {
					trx.settings.OnClosed(ConnectionIssue)
				}
			default:
				return
			}
		},

		WindowedRequestTracking: settings.WindowedRequestTracking,
	}, requestStore)

	trx.in = newReceivable(conn, Settings{
		ReadTimeout: settings.ReadTimeout,

		OnPDU: settings.OnPDU,

		OnAllPDU: settings.OnAllPDU,

		OnReceivingError: settings.OnReceivingError,

		OnClosed: func(state State) {
			switch state {
			case InvalidStreaming, UnbindClosing:
				// also close output
				_ = trx.out.close(ExplicitClosing)

				if trx.settings.OnClosed != nil {
					trx.settings.OnClosed(state)
				}
			default:
				return
			}
		},

		WindowedRequestTracking: settings.WindowedRequestTracking,

		response: func(p pdu.PDU) {
			_ = trx.Submit(p)
		},
	},
		requestStore,
	)
	return trx
}

func (trx *transceivable) start() {
	if trx.settings.WindowedRequestTracking != nil && trx.settings.ExpireCheckTimer > 0 {
		trx.wg.Add(1)
		go func() {
			defer trx.wg.Done()
			trx.windowCleanup()
		}()

	}
	trx.out.start()
	trx.in.start()
}

// SystemID returns tagged SystemID which is attached with bind_resp from SMSC.
func (trx *transceivable) SystemID() string {
	return trx.conn.systemID
}

// Close transceiver and stop underlying daemons.
func (trx *transceivable) Close() (err error) {
	return trx.closing(ExplicitClosing)
}

// Submit a PDU.
func (trx *transceivable) Submit(p pdu.PDU) error {
	return trx.out.Submit(p)
}

func (trx *transceivable) GetWindowSize() (int, error) {
	if trx.settings.WindowedRequestTracking != nil {
		ctx, cancelFunc := context.WithTimeout(context.Background(), trx.settings.StoreAccessTimeOut*time.Millisecond)
		defer cancelFunc()
		return trx.requestStore.Length(ctx)
	}
	return 0, ErrWindowNotConfigured

}

func (trx *transceivable) windowCleanup() {
	ticker := time.NewTicker(trx.settings.ExpireCheckTimer)
	defer ticker.Stop()
	for {
		select {
		case <-trx.ctx.Done():
			return
		case <-ticker.C:
			ctx, cancelFunc := context.WithTimeout(context.Background(), trx.settings.StoreAccessTimeOut*time.Millisecond)
			for _, request := range trx.requestStore.List(ctx) {
				if time.Since(request.TimeSent) > trx.settings.PduExpireTimeOut {
					_ = trx.requestStore.Delete(ctx, request.GetSequenceNumber())
					if trx.settings.OnExpiredPduRequest != nil {
						if trx.settings.OnExpiredPduRequest(request.PDU) {
							_ = trx.closing(ConnectionIssue)
						}
					}
				}
			}
			cancelFunc() //defer should not be used because we are inside loop
		}
	}
}

func (trx *transceivable) closing(state State) (err error) {
	if atomic.CompareAndSwapInt32(&trx.aliveState, Alive, Closed) {
		trx.cancel()

		// closing input and output
		_ = trx.out.close(StoppingProcessOnly)
		_ = trx.in.close(StoppingProcessOnly)

		// close underlying conn
		err = trx.conn.Close()

		// notify transceiver closed
		if trx.settings.OnClosed != nil {
			trx.settings.OnClosed(state)
		}

		trx.wg.Wait()
	}
	return
}
