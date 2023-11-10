package gosmpp

import "github.com/linxGnu/gosmpp/pdu"

// PDUCallback handles received PDU.
type PDUCallback func(pdu pdu.PDU, responded bool)

// AllPDUCallback handles all received PDU.
//
// This pdu is NOT responded to automatically, manual response/handling is needed
// and the bind can be closed by retuning true on closeBind.
type AllPDUCallback func(pdu pdu.PDU) (responsePdu pdu.PDU, closeBind bool)

// PDUErrorCallback notifies fail-to-submit PDU with along error.
type PDUErrorCallback func(pdu pdu.PDU, err error)

// ErrorCallback notifies happened error while reading PDU.
type ErrorCallback func(error)

// ClosedCallback notifies closed event due to State.
type ClosedCallback func(State)
