package pdu

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMalformUSME(t *testing.T) {
	t.Run("malformSME", func(t *testing.T) {
		b := NewBuffer(nil)
		var u UnsuccessSME
		require.NotNil(t, u.Unmarshal(b))
	})

	t.Run("malformSMEs", func(t *testing.T) {
		b := NewBuffer(nil)
		_ = b.WriteByte(1)
		var u UnsuccessSMEs
		require.NotNil(t, u.Unmarshal(b))
	})
}
