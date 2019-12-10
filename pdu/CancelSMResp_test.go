package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"
)

func TestCancelSMResp(t *testing.T) {
	v := NewCancelSMResp().(*CancelSMResp)
	require.False(t, v.CanResponse())
	require.Nil(t, v.GetResponse())

	validate(t,
		v,
		"00000010800000080000000000000001",
		data.CANCEL_SM_RESP,
	)
}
