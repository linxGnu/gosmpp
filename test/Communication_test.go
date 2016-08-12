package test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/linxGnu/gosmpp"
	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/PDU"
)

var session *gosmpp.Session

// TestBindingSMSC test binding connection with SMSC
func TestBindingSMSC(t *testing.T) {
	// connection, err := gosmpp.NewTCPIPConnectionWithAddrPort("localhost", 34567)
	connection, err := gosmpp.NewTCPIPConnectionWithAddrPort("localhost", 2775)
	if err != nil {
		t.Fail()
		return
	}

	request := PDU.NewBindTransciever()
	request.SetSystemId("smppclient1")
	request.SetPassword("password")
	request.SetSystemType("CMT")

	session = gosmpp.NewSessionWithConnection(connection)
	session.EnableStateChecking()

	listener := &TestPDUListener{}

	resp, e := session.BindWithListener(request, listener)
	if e != nil || resp.GetCommandStatus() != 0 {
		t.Fail()
		return
	}
	fmt.Println("Binding done!")

	resp, e = session.Unbind()
	if e != nil {
		t.Fail()
		return
	}

	time.Sleep(3 * time.Second)

	resp, e = session.BindWithListener(request, listener)
	if e != nil || resp.GetCommandStatus() != 0 {
		fmt.Println(e)
		t.Fail()
		return
	}
	fmt.Println("ReBinding done!")

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
	_, e = session.Submit(submit)
	time.Sleep(5 * time.Second)

	if e != nil {
		t.Errorf(e.Error.Error())
		t.Fail()
		return
	}
	fmt.Println("Done submit content:", "Biết đâu mà đợi")

}

type TestPDUListener struct {
}

func (c *TestPDUListener) HandleEvent(event *gosmpp.ServerPDUEvent) *Exception.Exception {
	switch event.GetPDU().(type) {
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
