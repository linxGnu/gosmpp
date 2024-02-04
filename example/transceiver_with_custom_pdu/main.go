package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/linxGnu/gosmpp"
	"github.com/linxGnu/gosmpp/pdu"
)

// This example uses the WindowPDUHandlerConfig, to show that a custom PDU can be used
// and the expected response will contain that original custom PDU with extra fields
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

			WindowPDUHandlerConfig: &gosmpp.WindowPDUHandlerConfig{
				OnReceivedPduRequest:  handleReceivedPduRequest(),
				OnExpectedPduResponse: handleExpectedPduResponse(),
				OnExpiredPduRequest:   handleExpirePduRequest(),
				PduExpireTimeOut:      30 * time.Second,
				ExpireCheckTimer:      10 * time.Second,
				MaxWindowSize:         30,
				EnableAutoRespond:     false,
			},
		}, 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = trans.Close()
	}()

	// sending SMS(s)
	for i := 0; i < 60; i++ {
		p := newCustomSubmitSM()
		if err = trans.Transceiver().Submit(p); err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Sent CustomSubmitSM with id: %+v\n", p.messageId)
		time.Sleep(time.Second)

	}
	time.Sleep(3 * time.Second)

}

func handleExpirePduRequest() func(pdu.PDU) {
	return func(p pdu.PDU) {
		switch p.(type) {

		case *CustomSubmitSM:
			fmt.Printf("Expired CustomSubmitSM:%+v\n", p)
		}
	}
}

func handleExpectedPduResponse() func(response gosmpp.Response) {
	// for this example, we only care about receiving our CustomSubmitSM response
	return func(response gosmpp.Response) {

		switch response.PDU.(type) {
		case *pdu.SubmitSMResp:
			p, ok := response.OriginalRequest.PDU.(CustomSubmitSM)
			if ok {
				fmt.Printf("SubmitSMResp Received, original CustomSubmitSM with id: %+v\n", p.messageId)
			}
		}
	}
}

func handleReceivedPduRequest() func(pdu.PDU) (pdu.PDU, bool) {
	// for this example, we are ignoring all Received PDU
	return func(p pdu.PDU) (pdu.PDU, bool) {
		return p.GetResponse(), false
	}
}
