package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/linxGnu/gosmpp"
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/pdu"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go sendingAndReceiveSMS(&wg)

	wg.Wait()
}

func sendingAndReceiveSMS(wg *sync.WaitGroup) {
	defer wg.Done()

	auth := gosmpp.Auth{
		SMSC:       "localhost:2775",
		SystemID:   "169994",
		Password:   "EDXPJU",
		SystemType: "",
	}

	sendResponse := make(chan pdu.PDU)
	defer close(sendResponse)
	trans, err := gosmpp.NewSession(
		gosmpp.TRXConnector(gosmpp.NonTLSDialer, auth),
		gosmpp.Settings{
			EnquireLink: 5 * time.Second,

			ReadTimeout: 10 * time.Second,

			OnSubmitError: func(_ pdu.PDU, err error) {
				log.Fatal("SubmitPDU error:", err)
			},

			OnReceivingError: func(err error) {
				fmt.Println("Receiving PDU/Network error:", err)
			},

			OnRebindingError: func(err error) {
				fmt.Println("Rebinding but error:", err)
			},

			OnAllPDU: handlePDU(sendResponse),

			OnClosed: func(state gosmpp.State) {
				fmt.Println(state)
			},
		}, 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = trans.Close()
	}()

	go func() {
		for p := range sendResponse {
			if err = trans.Transceiver().Respond(p); err != nil {
				fmt.Println(err)
			} else {
				switch p.(type) {
				case *pdu.UnbindResp:
					fmt.Println("UnbindResp Sent")
					_ = trans.Transceiver().Close()

				case *pdu.EnquireLinkResp:
					fmt.Println("EnquireLinkResp Sent")

				case *pdu.DataSMResp:
					fmt.Println("DataSMResp Sent")

				case *pdu.DeliverSMResp:
					fmt.Println("DeliverSMResp Sent")
				}
			}
		}
	}()

	// sending SMS(s)
	for i := 0; i < 30; i++ {
		if err = trans.Transceiver().Submit(newSubmitSM()); err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Second)
	}

}

func handlePDU(sendResponse chan pdu.PDU) func(pdu.PDU) {
	return func(p pdu.PDU) {
		switch pd := p.(type) {
		case *pdu.Unbind:
			fmt.Printf("Unbind:%+v\n", pd)
			sendResponse <- pd.GetResponse()

		case *pdu.UnbindResp:
			fmt.Println("UnbindResp Received")

		case *pdu.SubmitSMResp:
			fmt.Println("SubmitSMResp Received")

		case *pdu.GenericNack:
			fmt.Println("GenericNack Received")

		case *pdu.EnquireLinkResp:
			fmt.Println("EnquireLinkResp Received")

		case *pdu.EnquireLink:
			fmt.Println("EnquireLinkResp Received")

		case *pdu.DataSM:
			fmt.Println("DataSM receiver")
			sendResponse <- pd.GetResponse()

		case *pdu.DeliverSM:
			fmt.Println("DeliverSM receiver")
			sendResponse <- pd.GetResponse()
		}
	}
}

func newSubmitSM() *pdu.SubmitSM {
	// build up submitSM
	srcAddr := pdu.NewAddress()
	srcAddr.SetTon(5)
	srcAddr.SetNpi(0)
	_ = srcAddr.SetAddress("00" + "522241")

	destAddr := pdu.NewAddress()
	destAddr.SetTon(1)
	destAddr.SetNpi(1)
	_ = destAddr.SetAddress("99" + "522241")

	submitSM := pdu.NewSubmitSM().(*pdu.SubmitSM)
	submitSM.SourceAddr = srcAddr
	submitSM.DestAddr = destAddr
	_ = submitSM.Message.SetMessageWithEncoding("Đừng buồn thế dù ngoài kia vẫn mưa nghiễng rợi tý tỵ", data.UCS2)
	submitSM.ProtocolID = 0
	submitSM.RegisteredDelivery = 1
	submitSM.ReplaceIfPresentFlag = 0
	submitSM.EsmClass = 0

	return submitSM
}
