package pdu

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInvalidCmdID(t *testing.T) {
	v, err := CreatePDUFromCmdID(-12)
	require.Nil(t, v)
	require.NotNil(t, err)
}
