package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"
)

func TestDataSMResp(t *testing.T) {
	v := NewDataSMResp().(*DataSMResp)
	require.False(t, v.CanResponse())
	require.Nil(t, v.GetResponse())

	v.MessageID = "testMID"

	validate(t,
		v,
		"00000018800001030000000000000001746573744d494400",
		data.DATA_SM_RESP,
	)
}
