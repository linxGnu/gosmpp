package gosmpp

import (
	"context"
	"github.com/linxGnu/gosmpp/pdu"
	"time"
)

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

// RebindCallback notifies rebind event due to State.
type RebindCallback func()

type Request struct {
	pdu.PDU
	TimeSent time.Time
}

type Response struct {
	pdu.PDU
	OriginalRequest Request
}

type RequestWindowStore interface {
	Set(ctx context.Context, request Request)
	Get(ctx context.Context, sequenceNumber int32) (Request, bool)
	List(ctx context.Context) []Request
	Delete(ctx context.Context, sequenceNumber int32)
	Clear(ctx context.Context)
	Length(ctx context.Context) int
}
