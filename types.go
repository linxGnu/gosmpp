package gosmpp

import "github.com/linxGnu/gosmpp/pdu"

type PDUCallback func(pdu pdu.PDU, responded bool)

type PDUErrorCallback func(pdu pdu.PDU, err error)

type ErrorCallback func(error)

type ClosedCallback func(State)
