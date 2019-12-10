package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"
)

func TestBindRequest(t *testing.T) {
	t.Run("receiver", func(t *testing.T) {
		req := NewBindReceiver().(*BindRequest)
		require.True(t, req.CanResponse())

		resp := req.GetResponse()
		require.NotNil(t, resp)
		require.EqualValues(t, data.BIND_RECEIVER_RESP, resp.GetHeader().CommandID)

		req.SystemID = "system_id_fake"
		req.Password = "password"
		req.SystemType = "only"
		req.InterfaceVersion = 44
		require.Nil(t, req.AddressRange.SetAddressRange("emptY"))
		req.AddressRange.SetTon(23)
		req.AddressRange.SetNpi(101)

		validate(t,
			req,
			"0000003600000001000000000000000173797374656d5f69645f66616b650070617373776f7264006f6e6c79002c1765656d70745900",
			data.BIND_RECEIVER,
		)
	})

	t.Run("transmitter", func(t *testing.T) {
		req := NewBindTransmitter().(*BindRequest)
		require.True(t, req.CanResponse())

		resp := req.GetResponse()
		require.NotNil(t, resp)
		require.EqualValues(t, data.BIND_TRANSMITTER_RESP, resp.GetHeader().CommandID)

		req.SystemID = "system_id_fake"
		req.Password = "password"
		req.SystemType = "only"
		req.InterfaceVersion = 44
		require.Nil(t, req.AddressRange.SetAddressRange("emptY"))
		req.AddressRange.SetTon(23)
		req.AddressRange.SetNpi(101)

		validate(t,
			req,
			"0000003600000002000000000000000173797374656d5f69645f66616b650070617373776f7264006f6e6c79002c1765656d70745900",
			data.BIND_TRANSMITTER,
		)
	})

	t.Run("transceiver", func(t *testing.T) {
		req := NewBindTransceiver().(*BindRequest)
		require.True(t, req.CanResponse())

		resp := req.GetResponse()
		require.NotNil(t, resp)
		require.EqualValues(t, data.BIND_TRANSCEIVER_RESP, resp.GetHeader().CommandID)

		req.SystemID = "system_id_fake"
		req.Password = "password"
		req.SystemType = "only"
		req.InterfaceVersion = 44
		require.Nil(t, req.AddressRange.SetAddressRange("emptY"))
		req.AddressRange.SetTon(23)
		req.AddressRange.SetNpi(101)

		validate(t,
			req,
			"0000003600000009000000000000000173797374656d5f69645f66616b650070617373776f7264006f6e6c79002c1765656d70745900",
			data.BIND_TRANSCEIVER,
		)
	})
}
