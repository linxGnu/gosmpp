package gosmpp

import (
	"context"
	"fmt"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/linxGnu/gosmpp/pdu"

	"github.com/stretchr/testify/require"
)

func TestTransmitter(t *testing.T) {
	t.Run("binding", func(t *testing.T) {
		auth := nextAuth()
		transmitter, err := NewTransmitterSession(NonTLSDialer, auth, TransmitSettings{
			WriteTimeout: time.Second,
			OnSubmitError: func(p pdu.PDU, err error) {
				t.Fatal(err)
			},
			OnRebindingError: func(err error) {
				t.Fatal(err)
			},
			OnClosed: func(state State) {
				fmt.Println(state)
			},
		}, 5*time.Second)
		require.Nil(t, err)
		require.NotNil(t, transmitter)
		defer func() {
			_ = transmitter.Close()
		}()

		require.Equal(t, "MelroseLabsSMSC", transmitter.Transmitter().SystemID())

		err = transmitter.Transmitter().Submit(newSubmitSM(auth.SystemID))
		require.Nil(t, err)

		time.Sleep(400 * time.Millisecond)

		transmitter.rebind()
		err = transmitter.Transmitter().Submit(newSubmitSM(auth.SystemID))
		require.Nil(t, err)
	})

	errorHandling := func(t *testing.T, trigger func(*transmitter)) {
		conn, err := net.Dial("tcp", "smscsim.melroselabs.com:2775")
		require.Nil(t, err)

		var tr transmitter
		tr.input = make(chan pdu.PDU, 1)

		c := NewConnection(conn)
		defer func() {
			_ = c.Close()

			// write on closed conn?
			n, err := tr.write([]byte{1, 2, 3})
			require.NotNil(t, err)
			require.Zero(t, n)
		}()

		// fake settings
		tr.conn = c
		tr.ctx, tr.cancel = context.WithCancel(context.Background())

		var count int32
		tr.settings.OnClosed = func(state State) {
			atomic.AddInt32(&count, 1)
		}

		tr.settings.OnSubmitError = func(p pdu.PDU, err error) {
			require.NotNil(t, err)
			_, ok := p.(*pdu.CancelSM)
			require.True(t, ok)
		}
		tr.settings.WriteTimeout = 500 * time.Millisecond

		// do trigger
		trigger(&tr)

		time.Sleep(300 * time.Millisecond)
		require.NotZero(t, atomic.LoadInt32(&count))
	}

	t.Run("errorHandling1", func(t *testing.T) {
		errorHandling(t, func(tr *transmitter) {
			var p pdu.CancelSM
			tr.check(&p, 100, fmt.Errorf("fake error"))
		})
	})

	t.Run("errorHandling2", func(t *testing.T) {
		errorHandling(t, func(tr *transmitter) {
			var p pdu.CancelSM
			tr.check(&p, 0, fmt.Errorf("fake error"))
		})
	})

	t.Run("errorHandling3", func(t *testing.T) {
		errorHandling(t, func(tr *transmitter) {
			var p pdu.CancelSM
			tr.check(&p, 0, &net.DNSError{IsTemporary: false})
		})
	})
}
