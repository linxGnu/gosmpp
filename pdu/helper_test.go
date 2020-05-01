package pdu

import (
	"encoding/hex"
	"log"
	"testing"

	"github.com/linxGnu/gosmpp/data"
	"github.com/stretchr/testify/require"
)

func fromHex(h string) (v []byte) {
	var err error
	v, err = hex.DecodeString(h)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func toHex(v []byte) (h string) {
	h = hex.EncodeToString(v)
	return
}

func validate(t *testing.T, p PDU, hexValue string, expectCommandID data.CommandIDType) {
	buf := NewBuffer(nil)
	p.Marshal(buf)
	require.Equal(t, fromHex(hexValue), buf.Bytes())

	expectAfterParse(t, buf, p, expectCommandID)
}

func expectAfterParse(t *testing.T, b *ByteBuffer, expect PDU, expectCommandID data.CommandIDType) {
	c, err := Parse(b)
	require.Nil(t, err)
	require.Equal(t, expect, c)
	require.EqualValues(t, expectCommandID, c.GetHeader().CommandID)
	require.Zero(t, b.Len())
}
