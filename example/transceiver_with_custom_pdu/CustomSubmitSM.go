package main

import (
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/pdu"
	"math/rand"
	"strconv"
)

// CustomSubmitSM by embedding the PDU interface
// and adding messageId as an extra field to SubmitSM
type CustomSubmitSM struct {
	pdu.PDU
	messageId string
}

// newCustomSubmitSM returns CustomSubmitSM PDU.
// Using rand.Int to generate new id for each CustomSubmitSM
func newCustomSubmitSM() CustomSubmitSM {
	return CustomSubmitSM{
		PDU:       newSubmitSM(),
		messageId: strconv.Itoa(rand.Int()),
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
