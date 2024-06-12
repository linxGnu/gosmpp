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
					log.Println("SubmitPDU error: ", err)
				} else {
					log.Fatal("SubmitPDU error: ", err)
				}

			},

			OnReceivingError: func(err error) {
				fmt.Println("Receiving PDU/Network error: ", err)
			},

			OnRebindingError: func(err error) {
				fmt.Println("Rebinding but error: ", err)
			},

			OnClosed: func(state gosmpp.State) {
				fmt.Println("Bind has been closed: ", state.String())
			},

			WindowedRequestTracking: &gosmpp.WindowedRequestTracking{
				OnReceivedPduRequest:  handleReceivedPduRequest(),
				OnExpectedPduResponse: handleExpectedPduResponse(),
				OnExpiredPduRequest:   handleExpirePduRequest(),
				OnClosePduRequest:     handleOnClosePduRequest(),
				PduExpireTimeOut:      30 * time.Second,
				ExpireCheckTimer:      10 * time.Second,
				MaxWindowSize:         30,
				StoreAccessTimeOut:    1 * time.Second,
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
		fmt.Println("Current window size: ", trans.GetWindowSize())
		p := newCustomSubmitSM()
		if err = trans.Transceiver().Submit(p); err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Sent CustomSubmitSM with id: %+v\n", p.messageId)
		time.Sleep(1 * time.Second)
	}
	time.Sleep(1 * time.Second)
}

func handleExpirePduRequest() func(pdu.PDU) bool {
	return func(p pdu.PDU) bool {
		switch p.(type) {
		case *pdu.Unbind:
			fmt.Printf("Expired Unbind :%+v\n", p)
			fmt.Println("Unbind Expired")

		case *pdu.SubmitSM:
			fmt.Printf("Expired SubmitSM :%+v\n", p)

		case *pdu.EnquireLink:
			fmt.Printf("Expired EnquireLink:%+v\n", p)
			return true // if the enquire_link expired, usually means the bind is stale

		case *pdu.DataSM:
			fmt.Printf("Expired DataSM:%+v\n", p)
		}
		return false
	}
}

func handleOnClosePduRequest() func(pdu.PDU) {
	return func(p pdu.PDU) {
		switch p.(type) {
		case *pdu.Unbind:
			fmt.Printf("OnClose Unbind:%+v\n", p)

		case *pdu.SubmitSM:
			fmt.Printf("OnClose SubmitSM:%+v\n", p)

		case *pdu.EnquireLink:
			fmt.Printf("OnClose EnquireLink:%+v\n", p)

		case *pdu.DataSM:
			fmt.Printf("OnClose DataSM:%+v\n", p)
		}
	}
}

func handleExpectedPduResponse() func(response gosmpp.Response) {
	return func(response gosmpp.Response) {

		switch response.PDU.(type) {
		case *pdu.UnbindResp:
			fmt.Println("UnbindResp Received")
			fmt.Printf("OriginalSM id:%+v\n", response.OriginalRequest.PDU)

		case *pdu.SubmitSMResp:
			fmt.Printf("SubmitSMResp Received: %+v\n", response.PDU)
			fmt.Printf("OriginalSM SubmitSM:%+v\n", response.OriginalRequest.PDU)

		case *pdu.EnquireLinkResp:
			fmt.Println("EnquireLinkResp Received")
			fmt.Printf("Original EnquireLink:%+v\n", response.OriginalRequest.PDU)

		}
	}
}

func handleReceivedPduRequest() func(pdu.PDU) (pdu.PDU, bool) {
	return func(p pdu.PDU) (pdu.PDU, bool) {
		switch pd := p.(type) {
		case *pdu.Unbind:
			fmt.Println("Unbind Received")
			return pd.GetResponse(), true

		case *pdu.GenericNack:
			fmt.Println("GenericNack Received")

		case *pdu.EnquireLinkResp:
			fmt.Println("EnquireLinkResp Received")

		case *pdu.EnquireLink:
			fmt.Println("EnquireLink Received")
			return pd.GetResponse(), false

		case *pdu.DataSM:
			fmt.Println("DataSM Received")
			return pd.GetResponse(), false

		case *pdu.DeliverSM:
			fmt.Println("DeliverSM Received")
			return pd.GetResponse(), false
		}
		return nil, false
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
