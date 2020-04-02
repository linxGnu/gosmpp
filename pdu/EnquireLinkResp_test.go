package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"
)

func TestEnquireLinkResp(t *testing.T) {
	req := NewEnquireLink().(*EnquireLink)
	req.SequenceNumber = 13

	v := NewEnquireLinkRespFromReq(req).(*EnquireLinkResp)
	require.False(t, v.CanResponse())
	require.Nil(t, v.GetResponse())

	validate(t,
		v,
		"0000001080000015000000000000000d",
		data.ENQUIRE_LINK_RESP,
	)
}
