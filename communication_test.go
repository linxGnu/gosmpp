package gosmpp

import (
	"testing"
	"time"

	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/pdu"

	"github.com/stretchr/testify/require"
)

const (
	smscAddr   = "smscsim.melroselabs.com:2775"
	systemID   = "422315"
	password   = "GNCKOO"
	systemID2  = "459119"
	password2  = "JRAEPF"
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

func newAuth2() Auth {
	return Auth{
		SMSC:       smscAddr,
		SystemID:   systemID2,
		Password:   password2,
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
	_ = destAddr.SetAddress(systemID2)

	submitSM := pdu.NewSubmitSM().(*pdu.SubmitSM)
	submitSM.SourceAddr = srcAddr
	submitSM.DestAddr = destAddr
	_ = submitSM.Message.SetMessageWithEncoding(mess, data.UTF16)
	submitSM.DataCoding = 8
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
			t.Log(state)
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
	receiver, err := NewReceiverSession(NonTLSDialer, newAuth2(), ReceiveSettings{
		OnReceivingError: func(err error) {
			t.Log(err)
		},
		OnRebindingError: func(err error) {
			t.Log(err)
		},
		OnPDU: func(p pdu.PDU) {
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
}

func TestSubmitSM(t *testing.T) {
	trans, err := NewTransceiverSession(NonTLSDialer, newAuth(), TransceiveSettings{
		OnSubmitError: func(p pdu.PDU, err error) {
			t.Fatal(err)
		},
		OnReceivingError: func(err error) {
			t.Log(err)
		},
		OnRebindingError: func(err error) {
			t.Log(err)
		},
		OnPDU: func(p pdu.PDU) {
			t.Log(p)
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

	err = trans.Transceiver().Submit(newSubmitSM())
	require.Nil(t, err)

	time.Sleep(time.Second)
}

// // TestSubmitSMSC test submit to SMSC
// func TestSubmitSMSC(t *testing.T) {
// 	connection, err := gosmpp.NewTCPIPConnectionWithAddrPort(testSMSCAddr, testSMSCPort)
// 	if err != nil {
// 		t.Error(err)
// 		t.Fail()
// 		return
// 	}

// 	connection2, err := gosmpp.NewTCPIPConnectionWithAddrPort(testSMSCAddr, testSMSCPort)
// 	if err != nil {
// 		t.Error(err)
// 		t.Fail()
// 		return
// 	}

// 	request := PDU.NewBindTransceiver()
// 	request.SetSystemId("smppclient1")
// 	request.SetPassword("password")
// 	request.SetSystemType("CMT")

// 	request2 := PDU.NewBindTransceiver()
// 	request2.SetSystemId("smppclient2")
// 	request2.SetPassword("password")
// 	request2.SetSystemType("CMT")

// 	session := gosmpp.NewSessionWithConnection(connection)
// 	session.EnableStateChecking()

// 	session2 := gosmpp.NewSessionWithConnection(connection2)
// 	session.EnableStateChecking()

// 	listener := &TestPDUListener{id: "smppclient1", session: session}
// 	listener2 := &TestPDUListener{id: "smppclient2", session: session2}

// 	resp, e := session.BindWithListener(request, listener)
// 	if e != nil || resp.GetCommandStatus() != 0 {
// 		t.Error(e)
// 		t.Fail()
// 		return
// 	}
// 	session.GetReceiver().SetReceiveTimeout(-1)

// 	resp, e = session2.BindWithListener(request2, listener2)
// 	if e != nil || resp.GetCommandStatus() != 0 {
// 		t.Error(e)
// 		t.Fail()
// 		return
// 	}
// 	session2.GetReceiver().SetReceiveTimeout(-1)

// 	// Test submit
// 	submit := PDU.NewSubmitSM()
// 	sourceAddr, _ := PDU.NewAddressWithAddr("smppclient1")
// 	sourceAddr.SetTon(5)
// 	sourceAddr.SetNpi(0)
// 	desAddr, _ := PDU.NewAddressWithAddr("smppclient2")
// 	desAddr.SetTon(1)
// 	desAddr.SetNpi(1)
// 	submit.SetSourceAddr(sourceAddr)
// 	submit.SetDestAddr(desAddr)
// 	submit.SetShortMessageWithEncoding("Biết đâu mà đợi", Data.ENC_UTF16)
// 	submit.SetDataCoding(8)
// 	submit.SetProtocolId(0)
// 	submit.SetRegisteredDelivery(1)
// 	submit.SetReplaceIfPresentFlag(0)
// 	submit.SetEsmClass(0)
// 	submit.SetSequenceNumber(10)

// 	for i := 0; i < 100; i++ {
// 		if _, e = session.Submit(submit); e != nil {
// 			t.Errorf(e.Error.Error())
// 			t.Fail()
// 			return
// 		}
// 	}

// 	fmt.Println("Waiting 60 seconds to receive submitSMResp from SMSC or deliverSM")
// 	time.Sleep(60 * time.Second)
// 	session.Unbind()
// 	session2.Unbind()
// }

// type TestPDUListener struct {
// 	id      string
// 	session *gosmpp.Session
// }

// func (c *TestPDUListener) HandleEvent(event *gosmpp.ServerPDUEvent) *Exception.Exception {
// 	switch ev := event.GetPDU().(type) {
// 	case *PDU.DeliverSMResp:
// 		fmt.Println("DeliverSMResp", ev)

// 	case *PDU.SubmitSMResp:
// 		fmt.Println("SubmitSMResp", ev)

// 	case *PDU.DataSM:
// 		fmt.Println("DataSM", ev)
// 		resp, _ := ev.GetResponse()
// 		if resp != nil {
// 			c.session.Respond(resp)
// 		}

// 	case *PDU.DeliverSM:
// 		fmt.Println("DeliverSM", ev, c.id)

// 		// It's always better to do response without worrying!
// 		resp, _ := ev.GetResponse()
// 		if resp != nil {
// 			fmt.Println("Responding", resp, c.session)
// 			c.session.Respond(resp)
// 		}

// 		if ev.GetEsmClass() == 0 {
// 			if ev.GetDataCoding() == 0 {
// 				x, er := ev.GetShortMessage()
// 				fmt.Println("From:", ev.GetSourceAddr(), "with message:", x, er)
// 			} else {
// 				x, er := ev.GetShortMessageWithEncoding(Data.ENC_UTF16)
// 				fmt.Println("From:", ev.GetSourceAddr(), "with message:", x, er)
// 			}
// 		} else {
// 			if ev.HasReceiptedMessageId() {
// 				rm, e := ev.GetReceiptedMessageId()
// 				if e != nil {
// 					fmt.Println("DeliverSM Message ID:", rm, e)
// 				} else {
// 					rmid, err := strconv.Atoi(rm)
// 					fmt.Println("DeliverSM Message ID:", rmid, err)
// 				}
// 			}
// 		}
// 	}

// 	return nil
// }

// func (c *TestPDUListener) HandleException(err *Exception.Exception) {}
