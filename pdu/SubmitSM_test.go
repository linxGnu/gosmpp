package pdu

import (
	"fmt"
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"
)

func TestSubmitSM(t *testing.T) {
	v := NewSubmitSM().(*SubmitSM)
	require.True(t, v.CanResponse())
	v.SequenceNumber = 13

	validate(t,
		v.GetResponse(),
		"0000001180000004000000000000000d00",
		data.SUBMIT_SM_RESP,
	)

	v.ServiceType = "abc"
	_ = v.SourceAddr.SetAddress("Alicer")
	v.SourceAddr.SetTon(28)
	v.SourceAddr.SetNpi(29)

	_ = v.DestAddr.SetAddress("Bob")
	v.DestAddr.SetTon(79)
	v.DestAddr.SetNpi(80)

	v.EsmClass = 77
	v.ProtocolID = 99
	v.PriorityFlag = 61
	v.RegisteredDelivery = 83
	_ = v.Message.SetMessageWithEncoding("nghắ nghiêng nghiễng ngả", data.UCS2)
	v.Message.message = ""
	fmt.Println(v.Message.udHeader)

	validate(t,
		v,
		"0000005d00000004000000000000000d616263001c1d416c69636572004f50426f62004d633d00005300080030006e006700681eaf0020006e00670068006900ea006e00670020006e0067006800691ec5006e00670020006e00671ea3",
		data.SUBMIT_SM,
	)
}
