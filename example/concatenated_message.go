package main

import (
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/pdu"
)

func NewConcatenatedSubmitSM(phoneSrc, phoneDest, text string) (multiSubmitSM []*pdu.SubmitSM, err error) {

	multiShortMess, err := pdu.NewLongMessage(text)
	if err != nil {
		return
	}

	multiSubmitSM = make([]*pdu.SubmitSM, len(multiShortMess))
	for i, shortMess := range multiShortMess {

		submitSM := pdu.NewSubmitSM().(*pdu.SubmitSM)

		// source_ton set to alpha numeric since we use sms brandname
		submitSM.SourceAddr, err = pdu.NewAddressWithTonNpiAddr(data.GSM_TON_ALPHANUMERIC, data.GSM_NPI_UNKNOWN, phoneSrc)
		if err != nil {
			return nil, err
		}

		submitSM.DestAddr, err = pdu.NewAddressWithTonNpiAddr(data.GSM_TON_INTERNATIONAL, data.GSM_NPI_ISDN, phoneDest)
		if err != nil {
			return nil, err
		}

		submitSM.Message = *shortMess
		multiSubmitSM[i] = submitSM
	}

	return
}
