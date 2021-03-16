package gosmpp

import (
	"fmt"
	"testing"
	"time"

	"github.com/linxGnu/gosmpp/pdu"

	"github.com/stretchr/testify/require"
)

func TestReceiver(t *testing.T) {
	auth := nextAuth()
	receiver, err := NewReceiverSession(NonTLSDialer, auth, ReceiveSettings{
		Timeout: 2 * time.Second,

		OnReceivingError: func(err error) {
			fmt.Println(err)
		},

		OnRebindingError: func(err error) {
			fmt.Println(err)
		},

		OnPDU: func(p pdu.PDU, _ bool) {
			fmt.Println(p)
		},

		OnClosed: func(state State) {
			fmt.Println(state)
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
