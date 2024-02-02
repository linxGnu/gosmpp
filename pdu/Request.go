package pdu

import "time"

type Request struct {
	PDU
	TImeSent time.Time
}
