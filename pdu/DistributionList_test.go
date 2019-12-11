package pdu

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDistributionList(t *testing.T) {
	_, err := NewDistributionList("1234567890123456789012")
	require.NotNil(t, err)
}
