package pdu

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/linxGnu/gosmpp/utils"
)

func TestMalformUSME(t *testing.T) {
	t.Run("malformSME", func(t *testing.T) {
		b := utils.NewBuffer(nil)
		var u UnsuccessSME
		require.NotNil(t, u.Unmarshal(b))
	})

	t.Run("malformSMEs", func(t *testing.T) {
		b := utils.NewBuffer(nil)
		_ = b.WriteByte(1)
		var u UnsuccessSMEs
		require.NotNil(t, u.Unmarshal(b))
	})
}
