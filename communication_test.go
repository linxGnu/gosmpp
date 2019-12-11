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

const (
	smscAddr   = "smscsim.melroselabs.com:2775"
	systemID   = "062347"
	password   = "ZLMSQS"
	systemType = ""
	mess       = "Thử nghiệm: chuẩn bị nế mễ"
)

func newAuth() Auth {
	return Auth{
		SMSC:       smscAddr,
		SystemID:   systemID,
		Password:   password,
		SystemType: systemType,
	}
}

func newSubmitSM() *pdu.SubmitSM {
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

// TestBindingSMSC test binding connection with SMSC
func TestBindingSMSC(t *testing.T) {
	connection, err := ConnectAsTransceiver(NonTLSDialer, newAuth())
	require.Nil(t, err)
	require.NotNil(t, connection)

	// close connection
	_ = connection.Close()

	connection, err = ConnectAsReceiver(NonTLSDialer, newAuth())
	require.Nil(t, err)
	require.NotNil(t, connection)

	// close connection
	_ = connection.Close()

	connection, err = ConnectAsTransmitter(NonTLSDialer, newAuth())
	require.Nil(t, err)
	require.NotNil(t, connection)

	// close connection
	_ = connection.Close()
}

func TestTransmitter(t *testing.T) {
	transmitter, err := NewTransmitterSession(NonTLSDialer, newAuth(), TransmitSettings{
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

	err = transmitter.Transmitter().Submit(newSubmitSM())
	require.Nil(t, err)

	time.Sleep(time.Second)
}

func TestReceiver(t *testing.T) {
	receiver, err := NewReceiverSession(NonTLSDialer, newAuth(), ReceiveSettings{
		OnReceivingError: func(err error) {
			fmt.Println(err)
		},
		OnRebindingError: func(err error) {
			fmt.Println(err)
		},
		OnPDU: func(p pdu.PDU, responded bool) {
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
}

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

		case *pdu.GenerickNack:
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

	trans, err := NewTransceiverSession(NonTLSDialer, newAuth(), TransceiveSettings{
		EnquireLink: 200 * time.Millisecond,

		OnSubmitError: func(p pdu.PDU, err error) {
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

	// sending 10 SMS
	for i := 0; i < 10; i++ {
		err = trans.Transceiver().Submit(newSubmitSM())
		require.Nil(t, err)
		time.Sleep(100 * time.Millisecond)
	}

	time.Sleep(15 * time.Second)

	// wait response received
	require.True(t, atomic.LoadInt32(&countSubmitSMResp) == 10)
	require.True(t, atomic.LoadInt32(&countDeliverSM) > 0)
}
