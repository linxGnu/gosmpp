package pdu

import (
	"encoding/hex"
	"log"
	"testing"

	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/utils"

	"github.com/stretchr/testify/require"
)

var submitSMPayload []byte

func init() {
	var err error
	submitSMPayload, err = hex.DecodeString("0000003f00000004000000000000000200010034363730313133333131310001013436373039373731333337004000000000000001000803240103747B7374")
	if err != nil {
		log.Fatal(err)
	}
}

func TestSubmitSM(t *testing.T) {
	check := func(s *SubmitSM) {
		require.True(t, s.CanResponse())
		resp := s.GetResponse().(*SubmitSMResp)
		require.Empty(t, resp.MessageID)
		require.EqualValues(t, data.SUBMIT_SM_RESP, resp.CommandID)

		require.EqualValues(t, 63, s.CommandLength)
		require.EqualValues(t, 4, s.CommandID)
		require.EqualValues(t, 0, s.CommandStatus)
		require.EqualValues(t, 2, s.SequenceNumber)
		require.Equal(t, "", s.ServiceType)
		require.EqualValues(t, 1, s.SourceAddr.ton)
		require.EqualValues(t, 0, s.SourceAddr.npi)
		require.EqualValues(t, "46701133111", s.SourceAddr.Address())
		require.EqualValues(t, 1, s.DestAddr.ton)
		require.EqualValues(t, 1, s.DestAddr.npi)
		require.EqualValues(t, "46709771337", s.DestAddr.Address())
		require.EqualValues(t, 64, s.EsmClass)
		require.EqualValues(t, 0, s.ProtocolID)
		require.EqualValues(t, 0, s.PriorityFlag)
		require.Equal(t, "", s.ScheduleDeliveryTime)
		require.Equal(t, "", s.ValidityPeriod)
		require.EqualValues(t, 0, s.RegisteredDelivery)
		require.EqualValues(t, 0, s.ReplaceIfPresentFlag)
		require.EqualValues(t, 0, s.Message.SmDefaultMsgID)
		require.Equal(t, data.ASCII, s.Message.enc)
		message, err := s.Message.GetMessageWithEncoding(data.ASCII)
		require.Nil(t, err)
		require.EqualValues(t, "$t{st", message)

	}

	{
		b := utils.NewBuffer(submitSMPayload)
		pdu, err := Parse(b)
		require.Nil(t, err)
		check(pdu.(*SubmitSM))
	}

	{
		b := utils.NewBuffer(submitSMPayload)
		s := NewSubmitSM().(*SubmitSM)
		require.Nil(t, s.Unmarshal(b))
		check(s)

		// Check marshaling
		vb := utils.NewBuffer(nil)
		s.Marshal(vb)
		require.Equal(t, submitSMPayload, vb.Bytes())
	}
}
