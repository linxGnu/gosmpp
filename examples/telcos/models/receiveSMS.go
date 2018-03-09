package models

import (
	"database/sql"
	"time"
)

const (
	// TableReceiveSMS ...
	TableReceiveSMS = "receive_sms"
)

// ReceiveSMS ...
type ReceiveSMS struct {
	Id      uint64
	Isdn    uint64
	Content string

	Is_processed    bool
	Process_expired time.Time

	Request_check_errcode sql.NullInt64
	Request_check_message sql.NullString

	Request_charge_errcode sql.NullInt64
	Request_charge_message sql.NullString

	Request_pay_errcode sql.NullInt64
	Request_pay_message sql.NullString

	Created_at time.Time
}
