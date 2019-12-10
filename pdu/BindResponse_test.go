package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"
)

func TestBindResponse(t *testing.T) {
	t.Run("receiver", func(t *testing.T) {
		v := NewBindReceiverResp().(*BindResp)
		require.False(t, v.CanResponse())
		require.Nil(t, v.GetResponse())

		v.SystemID = "system_id_fake"

		validate(t,
			v,
			"0000001f80000001000000000000000173797374656d5f69645f66616b6500",
			data.BIND_RECEIVER_RESP,
		)
	})

	t.Run("transmitter", func(t *testing.T) {
		v := NewBindTransmitterResp().(*BindResp)
		require.False(t, v.CanResponse())
		require.Nil(t, v.GetResponse())

		v.SystemID = "system_id_fake"

		validate(t,
			v,
			"0000001f80000002000000000000000173797374656d5f69645f66616b6500",
			data.BIND_TRANSMITTER_RESP,
		)
	})

	t.Run("transceiver", func(t *testing.T) {
		v := NewBindTransceiverResp().(*BindResp)
		require.False(t, v.CanResponse())
		require.Nil(t, v.GetResponse())

		v.SystemID = "system_id_fake"

		validate(t,
			v,
			"0000001f80000009000000000000000173797374656d5f69645f66616b6500",
			data.BIND_TRANSCEIVER_RESP,
		)
	})
}
