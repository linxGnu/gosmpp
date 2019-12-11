package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"
)

func TestSubmitSMResp(t *testing.T) {
	v := NewSubmitSMResp().(*SubmitSMResp)
	require.False(t, v.CanResponse())
	require.Nil(t, v.GetResponse())

	v.MessageID = "football"

	validate(t,
		v,
		"00000019800000040000000000000001666f6f7462616c6c00",
		data.SUBMIT_SM_RESP,
	)
}
