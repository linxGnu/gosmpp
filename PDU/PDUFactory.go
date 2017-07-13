package PDU

import "errors"

var pduList []IPDU = []IPDU{
	NewBindTransmitter(),
	NewBindTransmitterResp(),
	NewBindReceiver(),
	NewBindReceiverResp(),
	NewBindTransceiver(),
	NewBindTransceiverResp(),
	NewUnbind(),
	NewUnbindResp(),
	NewOutbind(),
	NewSubmitSM(),
	NewSubmitSMResp(),
	NewSubmitMultiSM(),
	NewSubmitMultiSMResp(),
	NewDeliverSM(),
	NewDeliverSMResp(),
	NewDataSM(),
	NewDataSMResp(),
	NewQuerySM(),
	NewQuerySMResp(),
	NewCancelSM(),
	NewCancelSMResp(),
	NewReplaceSM(),
	NewReplaceSMResp(),
	NewEnquireLink(),
	NewEnquireLinkResp(),
	NewAlertNotification(),
	NewGenericNack(),
}

func CreatePDUWithCmdId(cmdId int32) (IPDU, error) {
	if pduList == nil {
		return nil, nil
	}

	for _, v := range pduList {
		if v != nil && v.GetCommandId() == cmdId {
			return v.GetInstance()
		}
	}

	return nil, errors.New("Unknown command id")
}
