package gosmpp

import (
	"testing"
	"time"

	"github.com/linxGnu/gosmpp/pdu"

	"github.com/stretchr/testify/require"
)

func TestReceive(t *testing.T) {
	auth := nextAuth()
	receiver, err := NewSession(
		RXConnector(NonTLSDialer, auth),
		Settings{
			ReadTimeout: 2 * time.Second,

			OnReceivingError: func(err error) {
				t.Log(err)
			},

			OnRebindingError: func(err error) {
				t.Log(err)
			},

			OnPDU: func(p pdu.PDU, _ bool) {
				t.Log(p)
			},

			OnClosed: func(state State) {
				t.Log(state)
			},
		}, 5*time.Second)
	require.Nil(t, err)
	require.NotNil(t, receiver)
	defer func() {
		_ = receiver.Close()
	}()

	require.Equal(t, "MelroseLabsSMSC", receiver.Receiver().SystemID())

	time.Sleep(time.Second)
	receiver.rebind()
}
