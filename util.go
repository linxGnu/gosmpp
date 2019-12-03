package gosmpp

import (
	"github.com/linxGnu/gosmpp/pdu"
	"github.com/linxGnu/gosmpp/utils"
)

func marshal(p pdu.PDU) []byte {
	buf := utils.NewBuffer(make([]byte, 0, 64))
	p.Marshal(buf)
	return buf.Bytes()
}
