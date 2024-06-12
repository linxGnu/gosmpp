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

	t := &transceivable{
		settings:     settings,
		conn:         conn,
		requestStore: requestStore,
	}
	t.ctx, t.cancel = context.WithCancel(context.Background())

	t.out = newTransmittable(conn, Settings{
		WriteTimeout: settings.WriteTimeout,

		EnquireLink: settings.EnquireLink,

		OnSubmitError: settings.OnSubmitError,

		OnClosed: func(state State) {
			switch state {
			case ConnectionIssue:
				// also close input
				_ = t.in.close(ExplicitClosing)

				if t.settings.OnClosed != nil {
					t.settings.OnClosed(ConnectionIssue)
				}
			default:
				return
			}
		},

		WindowedRequestTracking: settings.WindowedRequestTracking,
	}, requestStore)

	t.in = newReceivable(conn, Settings{
		ReadTimeout: settings.ReadTimeout,

		OnPDU: settings.OnPDU,

		OnAllPDU: settings.OnAllPDU,

		OnReceivingError: settings.OnReceivingError,

		OnClosed: func(state State) {
			switch state {
			case InvalidStreaming, UnbindClosing:
				// also close output
				_ = t.out.close(ExplicitClosing)

				if t.settings.OnClosed != nil {
					t.settings.OnClosed(state)
				}
			default:
				return
			}
		},

		WindowedRequestTracking: settings.WindowedRequestTracking,

		response: func(p pdu.PDU) {
			_ = t.Submit(p)
		},
	},
		requestStore,
	)
	return t
}

func (t *transceivable) start() {
	if t.settings.WindowedRequestTracking != nil && t.settings.ExpireCheckTimer > 0 {
		t.wg.Add(1)
		go func() {
			t.windowCleanup()
			defer t.wg.Done()
		}()

	}
	t.out.start()
	t.in.start()
}

// SystemID returns tagged SystemID which is attached with bind_resp from SMSC.
func (t *transceivable) SystemID() string {
	return t.conn.systemID
}

// Close transceiver and stop underlying daemons.
func (t *transceivable) Close() (err error) {
	if atomic.CompareAndSwapInt32(&t.aliveState, Alive, Closed) {
		// closing input and output
		_ = t.out.close(StoppingProcessOnly)
		_ = t.in.close(StoppingProcessOnly)

		// close underlying conn
		err = t.conn.Close()

		// notify transceiver closed
		if t.settings.OnClosed != nil {
			t.settings.OnClosed(ExplicitClosing)
		}
	}
	return
}

// Submit a PDU.
func (t *transceivable) Submit(p pdu.PDU) error {
	return t.out.Submit(p)
}

func (t *transceivable) GetWindowSize() (int, error) {
	if t.settings.WindowedRequestTracking != nil {
		ctx, cancelFunc := context.WithTimeout(context.Background(), t.settings.StoreAccessTimeOut*time.Millisecond)
		defer cancelFunc()
		return t.requestStore.Length(ctx)
	}
	return 0, ErrWindowNotConfigured

}

func (t *transceivable) windowCleanup() {
	ticker := time.NewTicker(t.settings.ExpireCheckTimer)
	defer ticker.Stop()
	for {
		select {
		case <-t.ctx.Done():
			return
		case <-ticker.C:
			ctx, cancelFunc := context.WithTimeout(context.Background(), t.settings.StoreAccessTimeOut*time.Millisecond)
			for _, request := range t.requestStore.List(ctx) {
				if time.Since(request.TimeSent) > t.settings.PduExpireTimeOut {
					_ = t.requestStore.Delete(ctx, request.GetSequenceNumber())
					if t.settings.OnExpiredPduRequest != nil {
						bindClose := t.settings.OnExpiredPduRequest(request.PDU)
						if bindClose {
							_ = t.Close()
						}
					}
				}
			}
			cancelFunc() //defer should not be used because we are inside loop
		}
	}
}
