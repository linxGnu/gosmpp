package gosmpp

import (
	"io"
	"time"

	"github.com/linxGnu/gosmpp/pdu"
)

// Transceiver interface.
type Transceiver interface {
	io.Closer
	Submit(pdu.PDU) error
	SystemID() string
}

// Transmitter interface.
type Transmitter interface {
	io.Closer
	Submit(pdu.PDU) error
	SystemID() string
}

// Receiver interface.
type Receiver interface {
	io.Closer
	SystemID() string
}

// Settings for TX (transmitter), RX (receiver), TRX (transceiver).
type Settings struct {
	// ReadTimeout is timeout for reading PDU from SMSC.
	// Underlying net.Conn will be stricted with ReadDeadline(now + timeout).
	// This setting is very important to detect connection failure.
	//
	// Must: ReadTimeout > max(0, EnquireLink)
	ReadTimeout time.Duration

	// WriteTimeout is timeout for submitting PDU.
	WriteTimeout time.Duration

	// EnquireLink periodically sends EnquireLink to SMSC.
	// The duration must not be smaller than 1 minute.
	//
	// Zero duration disables auto enquire link.
	EnquireLink time.Duration

	// OnPDU handles received PDU from SMSC.
	//
	// `Responded` flag indicates this pdu is responded automatically,
	// no manual respond needed.
	//
	// Will be ignored if OnAllPDU or WindowedRequestTracking is set
	OnPDU PDUCallback

	// OnAllPDU handles all received PDU from SMSC.
	//
	// This pdu is NOT responded to automatically,
	// manual response/handling is needed
	//
	// User can also decide to close bind by retuning true, default is false
	//
	// Will be ignored if WindowedRequestTracking is set
	OnAllPDU AllPDUCallback

	// OnReceivingError notifies happened error while reading PDU
	// from SMSC.
	OnReceivingError ErrorCallback

	// OnSubmitError notifies fail-to-submit PDU with along error.
	OnSubmitError PDUErrorCallback

	// OnRebindingError notifies error while rebinding.
	OnRebindingError ErrorCallback

	// OnClosed notifies `closed` event due to State.
	OnClosed ClosedCallback

	// OnRebind notifies `rebind` event due to State.
	OnRebind RebindCallback

	// SMPP Bind Window tracking feature config
	*WindowedRequestTracking

	response func(pdu.PDU)
}

// WindowedRequestTracking settings for TX (transmitter) and TRX (transceiver) request store.
type WindowedRequestTracking struct {

	// OnReceivedPduRequest handles received PDU request from SMSC.
	//
	// User can also decide to close bind by retuning true, default is false
	OnReceivedPduRequest AllPDUCallback

	// OnExpectedPduResponse handles expected PDU response from SMSC.
	// Only triggered when the original request is found in the window cache
	//
	// Handle is optional
	// If not set, response will be dropped
	OnExpectedPduResponse func(Response)

	// OnUnexpectedPduResponse handles unexpected PDU response from SMSC.
	// Only triggered if the original request is not found in the window cache
	//
	// Handle is optional
	// If not set, response will be dropped
	OnUnexpectedPduResponse func(pdu.PDU)

	// OnExpiredPduRequest handles expired PDU request with no response received
	//
	// Mandatory: the PduExpireTimeOut must be set
	// Handle is optional
	// If not set, expired PDU will be removed from cache
	// the bind can be closed by retuning true on closeBind.
	OnExpiredPduRequest func(pdu.PDU) (closeBind bool)

	// OnClosePduRequest will return all PDU request found in the store when the bind closes
	OnClosePduRequest func(pdu.PDU)

	// Set the number of second to expire a request sent to the SMSC
	//
	// Zero duration disables pdu expire check and the cache may fill up over time with expired PDU request
	// Recommended: eual or less to the value set in ReadTimeout + EnquireLink
	PduExpireTimeOut time.Duration

	// The time period between each check of the expired PDU in the cache
	//
	// Zero duration disables pdu expire check and the cache may fill up over time with expired PDU request
	// Recommended: Less or half the time set in for PduExpireTimeOut
	// Don't be too aggressive, there is a performance hit if the check is done often
	ExpireCheckTimer time.Duration

	// The maximum number of pending request sent to the SMSC
	//
	// Maximum value is 255
	MaxWindowSize uint8

	// if enabled, EnquireLink and Unbind request will be responded to automatically
	EnableAutoRespond bool

	// Set the number of millisecond to expire a request to store or retrieve data from request store
	//
	// Value must be greater than 0
	// 200 to 1000 is a good starting point
	StoreAccessTimeOut time.Duration
}
