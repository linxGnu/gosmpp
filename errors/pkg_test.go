package errors

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErr(t *testing.T) {
	require.True(t, strings.HasPrefix(ErrInvalidPDU.Error(), "Error happened: ["))
}
