package pdu

import (
	"encoding/hex"
	"log"
	"testing"

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

func validate(t *testing.T, p PDU, hexValue string, expectCommandID int32) {
	p.ResetSequenceNumber()
	require.EqualValues(t, 1, p.GetSequenceNumber())

	buf := NewBuffer(nil)
	p.Marshal(buf)
	require.Equal(t, fromHex(hexValue), buf.Bytes())

	expectAfterParse(t, buf, p, expectCommandID)
}

func expectAfterParse(t *testing.T, b *ByteBuffer, expect PDU, expectCommandID int32) {
	c, err := Parse(b)
	require.Nil(t, err)
	require.Equal(t, expect, c)
	require.EqualValues(t, expectCommandID, c.GetHeader().CommandID)
	require.Zero(t, b.Len())
}
