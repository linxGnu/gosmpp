package main

import (
	"errors"
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

	trans, err := gosmpp.NewSession(
		gosmpp.TRXConnector(gosmpp.NonTLSDialer, auth),
		gosmpp.Settings{
			EnquireLink: 5 * time.Second,

			ReadTimeout: 10 * time.Second,

			OnSubmitError: func(p pdu.PDU, err error) {
				if errors.Is(err, gosmpp.ErrWindowsFull) {
					log.Println("SubmitPDU error:", err)
				} else {
					log.Fatal("SubmitPDU error:", err)
				}

			},

			OnReceivingError: func(err error) {
				fmt.Println("Receiving PDU/Network error:", err)
			},

			OnRebindingError: func(err error) {
				fmt.Println("Rebinding but error:", err)
			},

			OnClosed: func(state gosmpp.State) {
				fmt.Println(state)
			},

			OnExpectedPduResponse: handleExpectedPdu(),

			OnExpiredPduRequest: handleExpirePDU(),

			PduExpireTimeOut: 5 * time.Second,

			MaxWindowSize: 30,
		}, 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = trans.Close()
	}()

	// sending SMS(s)
	for i := 0; i < 60; i++ {
		fmt.Println("Window size: ", trans.Transceiver().GetWindowSize())
		if err = trans.Transceiver().Submit(newSubmitSM()); err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Second)

	}
	time.Sleep(3 * time.Second)

}

func handleExpirePDU() func(pdu.PDU) {
	return func(p pdu.PDU) {
		switch p.(type) {
		case *pdu.Unbind:
			fmt.Printf("Expired Unbind:%+v\n", p)
			fmt.Println("Unbind Expired")

		case *pdu.SubmitSM:
			fmt.Printf("Expired SubmitSM:%+v\n", p)

		case *pdu.EnquireLink:
			fmt.Printf("Expired EnquireLink:%+v\n", p)

		case *pdu.DataSM:
			fmt.Printf("Expired DataSM:%+v\n", p)
		}
	}
}

func handleExpectedPdu() func(response pdu.Response) {
	return func(response pdu.Response) {

		switch response.PDU.(type) {

		case *pdu.UnbindResp:

			fmt.Println("UnbindResp Received")
			fmt.Printf("OriginalSM:%+v\n", response.OriginalRequest.PDU)

		case *pdu.SubmitSMResp:
			fmt.Println("SubmitSMResp Received")
			fmt.Printf("Original SubmitSM:%+v\n", response.OriginalRequest.PDU)

		case *pdu.EnquireLinkResp:
			fmt.Println("EnquireLinkResp Received")
			fmt.Printf("Original EnquireLink:%+v\n", response.OriginalRequest.PDU)

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
