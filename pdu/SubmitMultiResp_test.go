package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"
)

func TestSubmitMultiResp(t *testing.T) {
	req := NewSubmitMulti().(*SubmitMulti)
	req.SequenceNumber = 13

	v := NewSubmitMultiRespFromReq(req).(*SubmitMultiResp)
	require.False(t, v.CanResponse())
	require.Nil(t, v.GetResponse())

	v.MessageID = "football"

	addr1 := NewUnsuccessSMEWithTonNpi(38, 33, 19)
	require.Nil(t, addr1.SetAddress("Bob1"))
	require.EqualValues(t, 19, addr1.ErrorStatusCode())

	addr2, err := NewUnsuccessSMEWithAddr("Bob2", 20)
	require.Nil(t, err)
	require.EqualValues(t, 20, addr2.ErrorStatusCode())

	v.UnsuccessSMEs.Add(addr1, addr2)
	require.Equal(t, []UnsuccessSME{addr1, addr2}, v.UnsuccessSMEs.Get())

	validate(t,
		v,
		"0000003080000021000000000000000d666f6f7462616c6c00022621426f623100000000130000426f62320000000014",
		data.SUBMIT_MULTI_RESP,
	)
}
