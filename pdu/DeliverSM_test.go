package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"
)

func TestDeliverSM(t *testing.T) {
	v := NewDeliverSM().(*DeliverSM)
	require.True(t, v.CanResponse())
	v.SequenceNumber = 13

	validate(t,
		v.GetResponse(),
		"0000001180000005000000000000000d00",
		data.DELIVER_SM_RESP,
	)

	v.ServiceType = "abc"
	_ = v.SourceAddr.SetAddress("Alicer")
	v.SourceAddr.SetTon(28)
	v.SourceAddr.SetNpi(29)
	_ = v.DestAddr.SetAddress("Bobo")
	v.DestAddr.SetTon(30)
	v.DestAddr.SetNpi(31)
	v.EsmClass = 13
	v.ProtocolID = 99
	v.PriorityFlag = 61
	v.RegisteredDelivery = 83
	_ = v.Message.SetMessageWithEncoding("nghắ nghiêng nghiễng ngả", data.UCS2)
	v.Message.message = ""

	validate(t,
		v,
		"0000005e00000005000000000000000d616263001c1d416c69636572001e1f426f626f000d633d00005300080030006e006700681eaf0020006e00670068006900ea006e00670020006e0067006800691ec5006e00670020006e00671ea3",
		data.DELIVER_SM,
	)
}

func TestDeliverSMwithUDH(t *testing.T) {
	v := NewDeliverSM().(*DeliverSM)
	require.True(t, v.CanResponse())
	v.SequenceNumber = 13

	validate(t,
		v.GetResponse(),
		"0000001180000005000000000000000d00",
		data.DELIVER_SM_RESP,
	)

	v.ServiceType = "abc"
	_ = v.SourceAddr.SetAddress("Alicer")
	v.SourceAddr.SetTon(28)
	v.SourceAddr.SetNpi(29)
	_ = v.DestAddr.SetAddress("Bobo")
	v.DestAddr.SetTon(30)
	v.DestAddr.SetNpi(31)
	v.EsmClass = 77
	v.ProtocolID = 99
	v.PriorityFlag = 61
	v.RegisteredDelivery = 83
	_ = v.Message.SetMessageWithEncoding("nghắ nghiêng nghiễng ngả", data.UCS2)
	v.Message.message = ""
	v.Message.SetUDH(UDH{NewIEConcatMessage(2, 1, 254)})

	validate(t,
		v,
		"0000006400000005000000000000000d616263001c1d416c69636572001e1f426f626f004d633d00005300080036050003fe0201006e006700681eaf0020006e00670068006900ea006e00670020006e0067006800691ec5006e00670020006e00671ea3",
		data.DELIVER_SM,
	)
}
