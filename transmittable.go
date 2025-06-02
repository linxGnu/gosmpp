package gosmpp

import (
	"context"
	"errors"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/linxGnu/gosmpp/pdu"
)

var (
	// ErrConnectionClosing indicates transmitter is closing. Can not send any PDU.
	ErrConnectionClosing = errors.New("connection is closing, can not send PDU to SMSC")
	ErrWindowsFull       = errors.New("window full")
)

type transmittable struct {
	settings Settings

	wg    sync.WaitGroup
	input chan pdu.PDU

	conn *Connection

	aliveState   int32
	pendingWrite int32
	requestStore RequestStore
}

func newTransmittable(conn *Connection, settings Settings, requestStore RequestStore) *transmittable {
	tx := &transmittable{
		settings:     settings,
		conn:         conn,
		input:        make(chan pdu.PDU, 1),
		aliveState:   Alive,
		pendingWrite: 0,
		requestStore: requestStore,
	}

	return tx
}

func (tx *transmittable) close() {
	if atomic.CompareAndSwapInt32(&tx.aliveState, Alive, Closed) {
		for atomic.LoadInt32(&tx.pendingWrite) != 0 {
			runtime.Gosched()
		}

		// notify daemon
		close(tx.input)

		// wait daemon
		tx.wg.Wait()

		// try to send unbind
		_, _ = tx.write(pdu.NewUnbind())

		// concurrent-map has no func to verify initialization
		// we need to do the same check in
		if tx.settings.WindowedRequestTracking != nil {
			ctx, cancelFunc := context.WithTimeout(context.Background(), tx.settings.StoreAccessTimeOut)
			defer cancelFunc()
			var size int
			size, err := tx.requestStore.Length(ctx)
			if err != nil {
				return
			}
			if size > 0 {
				for _, request := range tx.requestStore.List(ctx) {
					if tx.settings.OnClosePduRequest != nil {
						tx.settings.OnClosePduRequest(request.PDU)
					}
					_ = tx.requestStore.Delete(ctx, request.GetSequenceNumber())
				}
			}
		}
	}
}

func (tx *transmittable) closing(state State) {
	// notify transceiver of closing
	go func() {
		tx.settings.OnClosed(state)
	}()
}

// Submit a PDU.
func (tx *transmittable) Submit(p pdu.PDU) (err error) {
	atomic.AddInt32(&tx.pendingWrite, 1)

	if atomic.LoadInt32(&tx.aliveState) == Alive {
		tx.input <- p
	} else {
		err = ErrConnectionClosing
	}

	atomic.AddInt32(&tx.pendingWrite, -1)
	return
}

func (tx *transmittable) start() {
	tx.wg.Add(1)
	if tx.settings.EnquireLink > 0 {
		go func() {
			defer tx.wg.Done()
			tx.loopWithEnquireLink()
		}()
	} else {
		go func() {
			defer tx.wg.Done()
			tx.loop()
		}()
	}
}

func (tx *transmittable) drain() {
	for range tx.input {
	}
}

func (tx *transmittable) loop() {
	defer tx.drain()

	for p := range tx.input {
		if p != nil {
			n, err := tx.write(p)
			if tx.check(p, n, err) {
				return
			}
		}
	}
}

func (tx *transmittable) loopWithEnquireLink() {
	ticker := time.NewTicker(tx.settings.EnquireLink)
	defer func() {
		ticker.Stop()
		tx.drain()
	}()

	for {
		select {
		case <-ticker.C:
			eqp := pdu.NewEnquireLink()
			n, err := tx.write(eqp)
			if tx.check(eqp, n, err) {
				return
			}

		case p, ok := <-tx.input:
			if !ok {
				return
			}

			if p != nil {
				n, err := tx.write(p)
				if tx.check(p, n, err) {
					return
				}
			}
		}
	}
}

// check error and do closing if need
func (tx *transmittable) check(p pdu.PDU, n int, err error) (closing bool) {
	if err == nil {
		return
	}

	if tx.settings.OnSubmitError != nil {
		tx.settings.OnSubmitError(p, err)
	}

	if n == 0 {
		if errors.Is(err, ErrWindowsFull) {
			closing = false
		} else if nErr, ok := err.(net.Error); ok {
			closing = nErr.Timeout()
		} else {
			closing = true
		}
	} else {
		closing = true // force closing
	}

	if closing {
		tx.closing(ConnectionIssue) // start closing
	}

	return
}

// low level writing
func (tx *transmittable) write(p pdu.PDU) (n int, err error) {
	if tx.settings.WriteTimeout > 0 {
		err = tx.conn.SetWriteTimeout(tx.settings.WriteTimeout)
	}
	if err != nil {
		return
	}

	if tx.settings.WindowedRequestTracking != nil && tx.settings.MaxWindowSize > 0 && isAllowPDU(p) {
		ctx, cancelFunc := context.WithTimeout(context.Background(), tx.settings.StoreAccessTimeOut)
		defer cancelFunc()
		var length int
		length, err = tx.requestStore.Length(ctx)
		if err != nil {
			return 0, err
		}
		if length < int(tx.settings.MaxWindowSize) {
			n, err = tx.conn.WritePDU(p)
			if err != nil {
				return 0, err
			}
			request := Request{
				PDU:      p,
				TimeSent: time.Now(),
			}
			err = tx.requestStore.Set(ctx, request)
			if err != nil {
				return 0, err
			}
		} else {
			return 0, ErrWindowsFull
		}
	} else {
		n, err = tx.conn.WritePDU(p)
	}

	return
}

func isAllowPDU(p pdu.PDU) bool {
	if p.CanResponse() {
		switch p.(type) {
		case *pdu.BindRequest, *pdu.Unbind:
			return false
		}
		return true
	}
	return false
}
