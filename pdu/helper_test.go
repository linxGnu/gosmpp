package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/utils"

	"github.com/stretchr/testify/require"
)

func validate(t *testing.T, p PDU, hexValue string, expectCommandID int32) {
	p.ResetSequenceNumber()

	buf := utils.NewBuffer(nil)
	p.Marshal(buf)
	require.Equal(t, fromHex(hexValue), buf.Bytes())
	expectAfterParse(t, buf, p, expectCommandID)
}

func expectAfterParse(t *testing.T, b *utils.ByteBuffer, expect PDU, expectCommandID int32) {
	c, err := Parse(b)
	require.Nil(t, err)
	require.Equal(t, expect, c)
	require.EqualValues(t, expectCommandID, c.GetHeader().CommandID)
}
