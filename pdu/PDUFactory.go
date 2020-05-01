package pdu

import (
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/errors"
)

type pduGenerator func() PDU

var pduMap = map[data.CommandIDType]pduGenerator{
	data.BIND_TRANSMITTER:      NewBindTransmitter,
	data.BIND_TRANSMITTER_RESP: NewBindTransmitterResp,
	data.BIND_TRANSCEIVER:      NewBindTransceiver,
	data.BIND_TRANSCEIVER_RESP: NewBindTransceiverResp,
	data.BIND_RECEIVER:         NewBindReceiver,
	data.BIND_RECEIVER_RESP:    NewBindReceiverResp,
	data.UNBIND:                NewUnbind,
	data.UNBIND_RESP:           NewUnbindResp,
	data.OUTBIND:               NewOutbind,
	data.SUBMIT_SM:             NewSubmitSM,
	data.SUBMIT_SM_RESP:        NewSubmitSMResp,
	data.SUBMIT_MULTI:          NewSubmitMulti,
	data.SUBMIT_MULTI_RESP:     NewSubmitMultiResp,
	data.DELIVER_SM:            NewDeliverSM,
	data.DELIVER_SM_RESP:       NewDeliverSMResp,
	data.DATA_SM:               NewDataSM,
	data.DATA_SM_RESP:          NewDataSMResp,
	data.QUERY_SM:              NewQuerySM,
	data.QUERY_SM_RESP:         NewQuerySMResp,
	data.CANCEL_SM:             NewCancelSM,
	data.CANCEL_SM_RESP:        NewCancelSMResp,
	data.REPLACE_SM:            NewReplaceSM,
	data.REPLACE_SM_RESP:       NewReplaceSMResp,
	data.ENQUIRE_LINK:          NewEnquireLink,
	data.ENQUIRE_LINK_RESP:     NewEnquireLinkResp,
	data.ALERT_NOTIFICATION:    NewAlertNotification,
	data.GENERIC_NACK:          NewGenerickNack,
}

// CreatePDUFromCmdID creates PDU from cmd id.
func CreatePDUFromCmdID(cmdID data.CommandIDType) (PDU, error) {
	if g, ok := pduMap[cmdID]; ok {
		return g(), nil
	}
	return nil, errors.ErrUnknownCommandID
}
