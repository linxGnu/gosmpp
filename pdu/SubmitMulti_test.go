package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"
)

func TestSubmitMulti(t *testing.T) {
	v := NewSubmitMulti().(*SubmitMulti)
	require.True(t, v.CanResponse())
	v.SequenceNumber = 13

	validate(t,
		v.GetResponse(),
		"0000001280000021000000000000000d0000",
		data.SUBMIT_MULTI_RESP,
	)

	v.ServiceType = "abc"
	_ = v.SourceAddr.SetAddress("Alicer")
	v.SourceAddr.SetTon(28)
	v.SourceAddr.SetNpi(29)

	addr := NewAddress()
	require.Nil(t, addr.SetAddress("Bob1"))
	d1 := NewDestinationAddress()
	d1.SetAddress(addr)

	dl, err := NewDistributionList("List1")
	require.Nil(t, err)
	d2 := NewDestinationAddress()
	d2.SetDistributionList(dl)

	dl, err = NewDistributionList("List2")
	require.Nil(t, err)
	d3 := NewDestinationAddress()
	d3.SetDistributionList(dl)

	v.DestAddrs.Add(d1, d2, d3)
	require.Equal(t, []DestinationAddress{d1, d2, d3}, v.DestAddrs.Get())

	v.EsmClass = 77
	v.ProtocolID = 99
	v.PriorityFlag = 61
	v.RegisteredDelivery = 83

	v.Message, err = NewShortMessageWithEncoding("nghắ nghiêng nghiễng ngả", data.UCS2)
	require.Nil(t, err)
	v.Message.message = ""

	validate(t,
		v,
		"0000006e00000021000000000000000d616263001c1d416c696365720003010000426f623100024c6973743100024c69737432004d633d00005300080030006e006700681eaf0020006e00670068006900ea006e00670020006e0067006800691ec5006e00670020006e00671ea3",
		data.SUBMIT_MULTI,
	)
}
