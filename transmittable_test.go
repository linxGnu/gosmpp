package gosmpp

import (
	"fmt"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/linxGnu/gosmpp/pdu"

	"github.com/stretchr/testify/require"
)

func TestTransmit(t *testing.T) {
	t.Run("Binding", func(t *testing.T) {
		auth := nextAuth()
		transmitter, err := NewSession(
			TXConnector(NonTLSDialer, auth),
			Settings{
				ReadTimeout: 2 * time.Second,

				OnPDU: func(p pdu.PDU, _ bool) {
					t.Logf("%+v\n", p)
				},

				OnSubmitError: func(_ pdu.PDU, err error) {
					t.Fatal(err)
				},

				OnRebindingError: func(err error) {
					t.Fatal(err)
				},

				OnClosed: func(state State) {
					t.Log(state)
				},
			}, -1)
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

	errorHandling := func(t *testing.T, trigger func(*transmittable)) {
		conn, err := net.Dial("tcp", "smscsim.melroselabs.com:2775")
		require.Nil(t, err)

		var tr transmittable
		tr.input = make(chan pdu.PDU, 1)

		c := NewConnection(conn)
		defer func() {
			_ = c.Close()

			// write on closed conn?
			n, err := tr.write(pdu.NewEnquireLink())
			require.NotNil(t, err)
			require.Zero(t, n)
		}()

		// fake settings
		tr.conn = c

		var count int32
		tr.settings.OnClosed = func(State) {
			atomic.AddInt32(&count, 1)
		}

		tr.settings.OnSubmitError = func(p pdu.PDU, err error) {
			require.NotNil(t, err)
			_, ok := p.(*pdu.CancelSM)
			require.True(t, ok)
		}

		// do trigger
		trigger(&tr)

		time.Sleep(300 * time.Millisecond)
		require.NotZero(t, atomic.LoadInt32(&count))
	}

	t.Run("ErrorHandling", func(t *testing.T) {
		errorHandling(t, func(tr *transmittable) {
			var p pdu.CancelSM
			tr.check(&p, 100, fmt.Errorf("fake error"))
		})

		errorHandling(t, func(tr *transmittable) {
			var p pdu.CancelSM
			tr.check(&p, 0, fmt.Errorf("fake error"))
		})

		errorHandling(t, func(tr *transmittable) {
			var p pdu.CancelSM
			tr.check(&p, 0, &net.DNSError{IsTemporary: false})
		})
	})

	t.Run("SubmitErr", func(t *testing.T) {
		var tr transmittable
		tr.input = make(chan pdu.PDU, 1)

		tr.aliveState = 1
		err := tr.Submit(nil)
		require.Error(t, err)

		tr.aliveState = 0
		err = tr.Submit(nil)
		require.NoError(t, err)
	})
}
