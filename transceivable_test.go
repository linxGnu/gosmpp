package gosmpp

import (
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
			t.Logf("%+v\n", pd)

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

func transceivableHandleAllPDU(t *testing.T) func(pdu.PDU) (pdu.PDU, bool) {
	return func(p pdu.PDU) (pdu.PDU, bool) {
		switch pd := p.(type) {
		case *pdu.SubmitSMResp:

			require.EqualValues(t, data.ESME_ROK, pd.CommandStatus)
			require.NotZero(t, len(pd.MessageID))
			atomic.AddInt32(&countSubmitSMResp, 1)
			return nil, false

		case *pdu.GenericNack:
			t.Fatal(pd)
			return nil, false
		case *pdu.DataSM:
			t.Logf("%+v\n", pd)
			return p.GetResponse(), false
		case *pdu.Unbind:
			t.Logf("%+v\n", pd)
			return p.GetResponse(), true

		case *pdu.DeliverSM:
			require.EqualValues(t, data.ESME_ROK, pd.CommandStatus)

			_mess, err := pd.Message.GetMessageWithEncoding(data.UCS2)
			assert.Nil(t, err)
			if mess == _mess {
				atomic.AddInt32(&countDeliverSM, 1)
			}
			return p.GetResponse(), false
		}
		return nil, false
	}
}

func TestTRXSubmitSM(t *testing.T) {
	auth := nextAuth()
	trans, err := NewSession(
		TRXConnector(NonTLSDialer, auth),
		Settings{
			ReadTimeout: 2 * time.Second,

			WriteTimeout: 3 * time.Second,

			EnquireLink: 200 * time.Millisecond,

			OnSubmitError: func(_ pdu.PDU, err error) {
				t.Fatal(err)
			},

			OnReceivingError: func(err error) {
				t.Log(err)
			},

			OnRebindingError: func(err error) {
				t.Log(err)
			},

			OnPDU: handlePDU(t),

			OnClosed: func(state State) {
				t.Log(state)
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
		time.Sleep(50 * time.Millisecond)
	}

	time.Sleep(5 * time.Second)

	// wait response received
	require.True(t, atomic.LoadInt32(&countSubmitSMResp) >= 15)

	// rebind and submit again
	trans.rebind()
	err = trans.Transceiver().Submit(newSubmitSM(auth.SystemID))
	require.Nil(t, err)
	time.Sleep(time.Second)
	require.True(t, atomic.LoadInt32(&countSubmitSMResp) >= 16)
}

func TestTRXSubmitSM_with_OnAllPDU(t *testing.T) {
	auth := nextAuth()
	trans, err := NewSession(
		TRXConnector(NonTLSDialer, auth),
		Settings{
			ReadTimeout: 2 * time.Second,

			WriteTimeout: 3 * time.Second,

			EnquireLink: 200 * time.Millisecond,

			OnSubmitError: func(_ pdu.PDU, err error) {
				t.Fatal(err)
			},

			OnReceivingError: func(err error) {
				t.Log(err)
			},

			OnRebindingError: func(err error) {
				t.Log(err)
			},

			OnAllPDU: transceivableHandleAllPDU(t),

			OnClosed: func(state State) {
				t.Log(state)
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
		time.Sleep(50 * time.Millisecond)
	}

	time.Sleep(5 * time.Second)

	// wait response received
	require.True(t, atomic.LoadInt32(&countSubmitSMResp) >= 15)

	// rebind and submit again
	trans.rebind()
	err = trans.Transceiver().Submit(newSubmitSM(auth.SystemID))
	require.Nil(t, err)
	time.Sleep(time.Second)
	require.True(t, atomic.LoadInt32(&countSubmitSMResp) >= 16)
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

func TestTRXSubmitSM_with_WindowConfig(t *testing.T) {
	auth := nextAuth()
	trans, err := NewSession(
		TRXConnector(NonTLSDialer, auth),
		Settings{
			ReadTimeout: 2 * time.Second,

			WriteTimeout: 3 * time.Second,

			EnquireLink: 200 * time.Millisecond,

			OnSubmitError: func(_ pdu.PDU, err error) {
				t.Fatal(err)
			},

			OnReceivingError: func(err error) {
				t.Log(err)
			},

			OnRebindingError: func(err error) {
				t.Log(err)
			},

			RequestWindowConfig: &RequestWindowConfig{
				OnReceivedPduRequest:  handleReceivedPduRequest(t),
				OnExpectedPduResponse: handleExpectedPduResponse(t),
				OnExpiredPduRequest:   nil,
				PduExpireTimeOut:      30 * time.Second,
				ExpireCheckTimer:      10 * time.Second,
				MaxWindowSize:         30,
				EnableAutoRespond:     false,
				RequestStore:          NewDefaultStore(),
			},

			OnClosed: func(state State) {
				t.Log(state)
			},
		}, 5*time.Second)
	require.Nil(t, err)
	require.NotNil(t, trans)
	defer func() {
		_ = trans.Close()
	}()

	require.Equal(t, "MelroseLabsSMSC", trans.Transceiver().SystemID())

	// sending 20 SMS
	for i := 0; i < 50; i++ {
		err = trans.Transceiver().Submit(newSubmitSM(auth.SystemID))
		require.Nil(t, err)
		time.Sleep(50 * time.Millisecond)
	}

	time.Sleep(5 * time.Second)

	// wait response received
	require.True(t, atomic.LoadInt32(&countSubmitSMResp) >= 15)

	// rebind and submit again
	trans.rebind()
	err = trans.Transceiver().Submit(newSubmitSM(auth.SystemID))
	require.Nil(t, err)
	time.Sleep(time.Second)
	require.True(t, atomic.LoadInt32(&countSubmitSMResp) >= 16)
}

func TestTRXSubmitSM_with_WindowConfig_and_AutoRespond(t *testing.T) {
	auth := nextAuth()
	trans, err := NewSession(
		TRXConnector(NonTLSDialer, auth),
		Settings{
			ReadTimeout: 2 * time.Second,

			WriteTimeout: 3 * time.Second,

			EnquireLink: 1 * time.Second,

			OnSubmitError: func(_ pdu.PDU, err error) {
				t.Fatal(err)
			},

			OnReceivingError: func(err error) {
				t.Log(err)
			},

			OnRebindingError: func(err error) {
				t.Log(err)
			},

			RequestWindowConfig: &RequestWindowConfig{
				OnReceivedPduRequest:  handleReceivedPduRequest(t),
				OnExpectedPduResponse: handleExpectedPduResponse(t),
				OnExpiredPduRequest:   nil,
				PduExpireTimeOut:      30 * time.Second,
				ExpireCheckTimer:      10 * time.Second,
				MaxWindowSize:         30,
				EnableAutoRespond:     true,
				RequestStore:          NewDefaultStore(),
			},

			OnClosed: func(state State) {
				t.Log("rebinded")
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
		time.Sleep(55 * time.Millisecond)
	}

	time.Sleep(5 * time.Second)

	// wait response received
	require.True(t, atomic.LoadInt32(&countSubmitSMResp) >= 15)

	// rebind and submit again
	trans.rebind()
	err = trans.Transceiver().Submit(newSubmitSM(auth.SystemID))
	require.Nil(t, err)
	time.Sleep(time.Second)
	require.True(t, atomic.LoadInt32(&countSubmitSMResp) >= 16)
}

func handleReceivedPduRequest(t *testing.T) func(pdu.PDU) (pdu.PDU, bool) {
	return func(p pdu.PDU) (pdu.PDU, bool) {
		switch pd := p.(type) {
		case *pdu.Unbind:
			return pd.GetResponse(), true

		case *pdu.GenericNack:
			t.Fatal(pd)

		case *pdu.EnquireLink:
			return pd.GetResponse(), false

		case *pdu.DataSM:
			t.Logf("%+v\n", pd)
			return pd.GetResponse(), false

		case *pdu.DeliverSM:
			require.EqualValues(t, data.ESME_ROK, pd.CommandStatus)

			_mess, err := pd.Message.GetMessageWithEncoding(data.UCS2)
			assert.Nil(t, err)
			if mess == _mess {
				atomic.AddInt32(&countDeliverSM, 1)
			}
			return pd.GetResponse(), false
		}
		return nil, false
	}
}

func handleExpectedPduResponse(t *testing.T) func(response Response) {
	return func(response Response) {

		switch pp := response.PDU.(type) {
		case *pdu.UnbindResp:
			//t.Logf("%+v\n", pp)

		case *pdu.SubmitSMResp:
			require.NotZero(t, len(pp.MessageID))
			atomic.AddInt32(&countSubmitSMResp, 1)
			t.Logf("%+v with original %+v\n", pp, response.OriginalRequest.PDU)

		case *pdu.EnquireLinkResp:
			t.Logf("%+v\n", pp)
		}
	}
}

func Test_newTransceivable(t *testing.T) {
	t.Run("always receive a non nil response", func(t *testing.T) {
		trans := newTransceivable(nil, Settings{})
		assert.NotNil(t, trans.in.settings.response)
	})
}
