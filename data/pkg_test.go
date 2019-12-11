package data

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultNpi(t *testing.T) {
	SetDefaultNpi(13)
	require.EqualValues(t, 13, GetDefaultNpi())
}

func TestDefaultTon(t *testing.T) {
	SetDefaultTon(19)
	require.EqualValues(t, 19, GetDefaultTon())
}
