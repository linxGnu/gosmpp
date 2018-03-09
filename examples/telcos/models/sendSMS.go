package models

import "time"

const (
	// TableSendSMS ...
	TableSendSMS = "send_sms"

	SUBMIT_STATUS_NOT_SEND  = 0
	SUBMIT_STATUS_SUBMITTED = 1
	SUBMIT_STATUS_SUCC      = 2
	SUBMIT_STATUS_FAIL      = 3
)

// SendSMS ...
type SendSMS struct {
	Id                 uint64
	Alias              string
	Isdn               uint64
	Content            string
	Submit_status      byte // 0: not send, 1: submitted, 2: status success, 3: status fail
	Submit_expired     time.Time
	Submit_status_code int32
	Smsc_message_id    int64
	Created_at         time.Time
}

// SendSMSReq ...
type SendSMSReq struct {
	Id      uint64
	Alias   string
	Content string
}
