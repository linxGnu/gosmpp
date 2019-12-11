package pdu

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNextSeq(t *testing.T) {
	var v int32 = math.MaxInt32
	require.EqualValues(t, 1, nextSequenceNumber(&v))
}
