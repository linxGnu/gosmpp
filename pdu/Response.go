package pdu

type Response struct {
	PDU
	OriginalRequest Request
}
