package pdu

import (
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/utils"
)

// AlertNotification PDU.
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
func (a *AlertNotification) Marshal(b *utils.ByteBuffer) {
	a.base.marshal(b, func(b *utils.ByteBuffer) {
		a.SourceAddr.Marshal(b)
		a.EsmeAddr.Marshal(b)
	})
}

// Unmarshal implements PDU interface.
func (a *AlertNotification) Unmarshal(b *utils.ByteBuffer) error {
	return a.base.unmarshal(b, func(b *utils.ByteBuffer) (err error) {
		if err = a.SourceAddr.Unmarshal(b); err == nil {
			err = a.EsmeAddr.Unmarshal(b)
		}
		return
	})
}
