package gosmpp

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/pdu"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	countSubmitSMResp, countDeliverSM int32
)

func handlePDU(t *testing.T) func(pdu.PDU, bool) {
	return func(p pdu.PDU, responded bool) {
		switch pd := p.(type) {
		case *pdu.SubmitSMResp:
			require.False(t, responded)
			require.EqualValues(t, data.ESME_ROK, pd.CommandStatus)
			require.NotZero(t, len(pd.MessageID))
			atomic.AddInt32(&countSubmitSMResp, 1)

		case *pdu.GenericNack:
			require.False(t, responded)
			t.Fatal(pd)

		case *pdu.DataSM:
			require.True(t, responded)
			fmt.Println(pd.Header)

		case *pdu.DeliverSM:
			require.True(t, responded)
			require.EqualValues(t, data.ESME_ROK, pd.CommandStatus)

			_mess, err := pd.Message.GetMessageWithEncoding(data.UCS2)
			assert.Nil(t, err)
			if mess == _mess {
				atomic.AddInt32(&countDeliverSM, 1)
			}

		}
	}
}

func TestSubmitSM(t *testing.T) {
	auth := nextAuth()
	trans, err := NewTransceiverSession(NonTLSDialer, auth, TransceiveSettings{
		EnquireLink: 200 * time.Millisecond,

		OnSubmitError: func(_ pdu.PDU, err error) {
			t.Fatal(err)
		},

		OnReceivingError: func(err error) {
			fmt.Println(err)
		},

		OnRebindingError: func(err error) {
			fmt.Println(err)
		},

		OnPDU: handlePDU(t),

		OnClosed: func(state State) {
			fmt.Println(state)
		},
	}, 5*time.Second)
	require.Nil(t, err)
	require.NotNil(t, trans)
	defer func() {
		_ = trans.Close()
	}()

	require.Equal(t, "MelroseLabsSMSC", trans.Transceiver().SystemID())

	// sending 20 SMS
	for i := 0; i < 20; i++ {
		err = trans.Transceiver().Submit(newSubmitSM(auth.SystemID))
		require.Nil(t, err)
		time.Sleep(100 * time.Millisecond)
	}

	time.Sleep(5 * time.Second)

	// wait response received
	require.EqualValues(t, 20, atomic.LoadInt32(&countSubmitSMResp))

	// rebind and submit again
	trans.rebind()
	err = trans.Transceiver().Submit(newSubmitSM(auth.SystemID))
	require.Nil(t, err)
	time.Sleep(time.Second)
	require.EqualValues(t, 21, atomic.LoadInt32(&countSubmitSMResp))
}

func newSubmitSM(systemID string) *pdu.SubmitSM {
	// build up submitSM
	srcAddr := pdu.NewAddress()
	srcAddr.SetTon(5)
	srcAddr.SetNpi(0)
	_ = srcAddr.SetAddress(systemID)

	destAddr := pdu.NewAddress()
	destAddr.SetTon(1)
	destAddr.SetNpi(1)
	_ = destAddr.SetAddress("12" + systemID)

	submitSM := pdu.NewSubmitSM().(*pdu.SubmitSM)
	submitSM.SourceAddr = srcAddr
	submitSM.DestAddr = destAddr
	_ = submitSM.Message.SetMessageWithEncoding(mess, data.UCS2)
	submitSM.ProtocolID = 0
	submitSM.RegisteredDelivery = 1
	submitSM.ReplaceIfPresentFlag = 0
	submitSM.EsmClass = 0

	return submitSM
}
