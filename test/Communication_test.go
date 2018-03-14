package test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/tsocial/gosmpp"
	"github.com/tsocial/gosmpp/Data"
	"github.com/tsocial/gosmpp/Exception"
	"github.com/tsocial/gosmpp/PDU"
)

const (
	testSMSCAddr = "localhost"
	testSMSCPort = 34567
)

var session *gosmpp.Session

// TestBindingSMSC test binding connection with SMSC
func TestBindingSMSC(t *testing.T) {
	connection, err := gosmpp.NewTCPIPConnectionWithAddrPort(testSMSCAddr, testSMSCPort)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	request := PDU.NewBindTransceiver()
	request.SetSystemId("smppclient1")
	request.SetPassword("password")
	request.SetSystemType("CMT")

	session = gosmpp.NewSessionWithConnection(connection)
	session.EnableStateChecking()

	listener := &TestPDUListener{}

	resp, e := session.BindWithListener(request, listener)
	if e != nil || resp.GetCommandStatus() != 0 {
		t.Error(e)
		t.Fail()
		return
	}

	resp, e = session.Unbind()
	if e != nil {
		t.Error(e)
		t.Fail()
		return
	}
}

// TestSubmitSMSC test submit to SMSC
func TestSubmitSMSC(t *testing.T) {
	connection, err := gosmpp.NewTCPIPConnectionWithAddrPort(testSMSCAddr, testSMSCPort)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	request := PDU.NewBindTransceiver()
	request.SetSystemId("smppclient1")
	request.SetPassword("password")
	request.SetSystemType("CMT")

	session = gosmpp.NewSessionWithConnection(connection)
	session.EnableStateChecking()

	listener := &TestPDUListener{}

	resp, e := session.BindWithListener(request, listener)
	if e != nil || resp.GetCommandStatus() != 0 {
		t.Error(e)
		t.Fail()
		return
	}

	// Test submit
	submit := PDU.NewSubmitSM()
	sourceAddr, _ := PDU.NewAddressWithAddr("smppclient1")
	sourceAddr.SetTon(5)
	sourceAddr.SetNpi(0)
	desAddr, _ := PDU.NewAddressWithAddr("smppclient2")
	desAddr.SetTon(1)
	desAddr.SetNpi(1)
	submit.SetSourceAddr(sourceAddr)
	submit.SetDestAddr(desAddr)
	submit.SetShortMessageWithEncoding("Biết đâu mà đợi", Data.ENC_UTF16)
	submit.SetDataCoding(8)
	submit.SetProtocolId(0)
	submit.SetRegisteredDelivery(1)
	submit.SetReplaceIfPresentFlag(0)
	submit.SetEsmClass(0)
	submit.SetSequenceNumber(10)

	if _, e = session.Submit(submit); e != nil {
		t.Errorf(e.Error.Error())
		t.Fail()
		return
	}

	fmt.Println("Waiting 15 seconds to receive submitSMResp from SMSC or deliverSM")
	time.Sleep(15 * time.Second)
	fmt.Println("Done")
}

type TestPDUListener struct {
}

func (c *TestPDUListener) HandleEvent(event *gosmpp.ServerPDUEvent) *Exception.Exception {
	switch event.GetPDU().(type) {
	case *PDU.SubmitSMResp:
		t := event.GetPDU().(*PDU.SubmitSMResp)
		fmt.Println("SUBMIT SM RESP", t.GetMessageId())
	case *PDU.DeliverSM:
		t := event.GetPDU().(*PDU.DeliverSM)

		// It's always better to do response without worrying!
		resp, _ := t.GetResponse()
		if resp != nil {
			session.Respond(resp)
		}

		if t.GetEsmClass() == 0 {
			if t.GetDataCoding() == 0 {
				x, er := t.GetShortMessage()
				fmt.Println("From:", t.GetSourceAddr(), "with message:", x, er)
			} else {
				x, er := t.GetShortMessageWithEncoding(Data.ENC_UTF16)
				fmt.Println("From:", t.GetSourceAddr(), "with message:", x, er)
			}
		} else {
			if t.HasReceiptedMessageId() {
				rm, e := t.GetReceiptedMessageId()
				if e != nil {
					fmt.Println("DeliverSM Message ID:", rm, e)
				} else {
					rmid, err := strconv.Atoi(rm)
					fmt.Println("DeliverSM Message ID:", rmid, err)
				}
			}
		}
	}

	return nil
}
