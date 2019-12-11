package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"
	"github.com/stretchr/testify/require"
)

func TestSubmitMulti(t *testing.T) {
	v := NewSubmitMulti().(*SubmitMulti)
	require.True(t, v.CanResponse())

	validate(t,
		v.GetResponse(),
		"000000128000002100000000000000010000",
		data.SUBMIT_MULTI_RESP,
	)

	v.ServiceType = "abc"
	_ = v.SourceAddr.SetAddress("Alicer")
	v.SourceAddr.SetTon(28)
	v.SourceAddr.SetNpi(29)

	d1, err := NewDestinationAddressFromAddress("Bob1")
	require.Nil(t, err)
	d2, err := NewDestinationAddressFromDistributionList("List1")
	require.Nil(t, err)
	d3 := NewDestinationAddress()
	require.Nil(t, d3.SetDistributionList("List2"))

	v.DestAddrs.Add(d1, d2, d3)
	require.Equal(t, []DestinationAddress{d1, d2, d3}, v.DestAddrs.Get())

	v.EsmClass = 77
	v.ProtocolID = 99
	v.PriorityFlag = 61
	v.RegisteredDelivery = 83
	_ = v.Message.SetMessageWithEncoding("nghắ nghiêng nghiễng ngả", data.UCS2)
	v.Message.message = ""

	validate(t,
		v,
		"0000006e000000210000000000000001616263001c1d416c696365720003010000426f623100024c6973743100024c69737432004d633d00005300080030006e006700681eaf0020006e00670068006900ea006e00670020006e0067006800691ec5006e00670020006e00671ea3",
		data.SUBMIT_MULTI,
	)
}
