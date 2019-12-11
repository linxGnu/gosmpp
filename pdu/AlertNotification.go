package pdu

import (
	"github.com/linxGnu/gosmpp/data"
)

// AlertNotification PDU is sent by the SMSC to the ESME, when the SMSC has detected that
// a particular mobile subscriber has become available and a delivery pending flag had been
// set for that subscriber from a previous data_sm operation.
type AlertNotification struct {
	base
	SourceAddr Address
	EsmeAddr   Address
}

// NewAlertNotification create new alert notification pdu.
func NewAlertNotification() PDU {
	a := &AlertNotification{
		base: newBase(),
	}
	a.CommandID = data.ALERT_NOTIFICATION
	return a
}

// CanResponse implements PDU interface.
func (a *AlertNotification) CanResponse() bool {
	return false
}

// GetResponse implements PDU interface.
func (a *AlertNotification) GetResponse() PDU {
	return nil
}

// Marshal implements PDU interface.
func (a *AlertNotification) Marshal(b *ByteBuffer) {
	a.base.marshal(b, func(b *ByteBuffer) {
		a.SourceAddr.Marshal(b)
		a.EsmeAddr.Marshal(b)
	})
}

// Unmarshal implements PDU interface.
func (a *AlertNotification) Unmarshal(b *ByteBuffer) error {
	return a.base.unmarshal(b, func(b *ByteBuffer) (err error) {
		if err = a.SourceAddr.Unmarshal(b); err == nil {
			err = a.EsmeAddr.Unmarshal(b)
		}
		return
	})
}
