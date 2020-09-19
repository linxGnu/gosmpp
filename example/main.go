package main

import (
	"log"
	"strings"
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
		// see https://melroselabs.com/services/smsc-simulator/#smsc-simulator-try
		SMSC:       "smscsim.melroselabs.com:2775",
		SystemID:   "your test system id",
		Password:   "your test password",
		SystemType: "",
	}

	trans, err := gosmpp.NewTransceiverSession(gosmpp.NonTLSDialer, auth, gosmpp.TransceiveSettings{
		EnquireLink:  5 * time.Second,
		WriteTimeout: time.Second,
		// this setting is very important to detect broken conn.
		// After timeout, there is no read packet, then we decide it's connection broken.
		ReadTimeout: 10 * time.Second,

		OnSubmitError:    func(p pdu.PDU, err error) { log.Fatal("SubmitPDU error:", err) },
		OnReceivingError: func(err error) { log.Println("Receiving PDU/Network error:", err) },
		OnRebindingError: func(err error) { log.Println("Rebinding but error:", err) },
		OnPDU:            handlePDU(),
		OnClosed:         func(state gosmpp.State) { log.Println(state) },
	}, 5*time.Second)
	if err != nil {
		log.Println(err)
	}
	defer func() {
		_ = trans.Close()
	}()

	// sending SMS(s)
	for i := 0; i < 1800; i++ {
		if err = trans.Transceiver().Submit(newSubmitSM()); err != nil {
			log.Println(err)
		}
		time.Sleep(time.Second)
	}
}

func handlePDU() gosmpp.PDUCallback {
	concatenated := map[uint16][]string{}
	return func(p pdu.PDU, responded bool) {
		switch pd := p.(type) {
		case *pdu.SubmitSMResp:
			log.Printf("SubmitSMResp:%+v\n", pd)

		case *pdu.GenericNack:
			log.Println("GenericNack Received")

		case *pdu.EnquireLinkResp:
			log.Println("EnquireLinkResp Received")

		case *pdu.DataSM:
			log.Printf("DataSM:%+v\n", pd)

		case *pdu.DeliverSM:
			log.Printf("DeliverSM:%+v\n", pd)
			log.Println(pd.Message.GetMessage())
			// region concatenated sms (sample code)
			message, err := pd.Message.GetMessage()
			if err != nil {
				log.Fatal(err)
			}
			totalParts, sequence, reference, found := pd.Message.UDH().GetConcatInfo()
			if found {
				if _, ok := concatenated[reference]; !ok {
					concatenated[reference] = make([]string, totalParts)
				}
				concatenated[reference][sequence-1] = message
			}
			if !found {
				log.Println(message)
			} else if parts, ok := concatenated[reference]; ok && isConcatenatedDone(parts, totalParts) {
				log.Println(strings.Join(parts, ""))
				delete(concatenated, reference)
			}
			// endregion
		}
	}
}

func isConcatenatedDone(parts []string, total byte) bool {
	for _, part := range parts {
		if part != "" {
			total--
		}
	}
	return total == 0
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
