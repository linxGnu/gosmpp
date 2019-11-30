package pdu

import (
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/utils"
)

// AlertNotification PDU.
type AlertNotification struct {
	base
}

// NewAlertNotification create new alert notification pdu.
func NewAlertNotification() (a *AlertNotification) {
	a = &AlertNotification{}
	a.CommandID = data.ALERT_NOTIFICATION
	return
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
	a.base.marshal(b, nil)
}

// Unmarshal implements PDU interface.
func (a *AlertNotification) Unmarshal(b *utils.ByteBuffer) error {
	return a.base.unmarshal(b, nil)
}
