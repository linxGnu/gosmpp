package main

import (
	"math/rand"
	"strconv"

	"github.com/linxGnu/gosmpp/pdu"
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
