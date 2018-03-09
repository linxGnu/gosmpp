package dao

import (
	"github.com/linxGnu/gosmpp/examples/telcos/models"

	"fmt"
	"time"
)

const (
	qSendSMSSelectNotSend = "SELECT * FROM " + models.TableSendSMS + " WHERE submit_status = 0 and submit_expired >= ?"
	qSendSMSSave          = "INSERT INTO " + models.TableSendSMS + " (alias, isdn, content, submit_status, submit_expired, created_at) VALUES (?,?,?,?,?,?)"
	qSendSMSUpdate        = "UPDATE " + models.TableSendSMS + " SET submit_status = ?, smsc_message_id = ? WHERE id = ?"
)

// ISendSMSDAO send sms dao interface
type ISendSMSDAO interface {
	SelectNotSend() ([]*models.SendSMS, error)
	Save(alias string, isdn uint64, content string, duration time.Duration) error
	Update(id uint64, submitStatus byte, msgID int64) error
}

type sendSMSDAO struct{}

var errDBNotInit = fmt.Errorf("DB not initialized")

func (c *sendSMSDAO) SelectNotSend() (result []*models.SendSMS, err error) {
	db := getDB()
	if db == nil {
		err = errDBNotInit
		return
	}

	result = []*models.SendSMS{}
	err = db.Select(&result, qSendSMSSelectNotSend, time.Now())

	return
}

func (c *sendSMSDAO) Save(alias string, isdn uint64, content string, duration time.Duration) error {
	db := getDB()
	if db == nil {
		return errDBNotInit
	}

	now := time.Now()

	_, err := db.Exec(qSendSMSSave, alias, isdn, content, models.SUBMIT_STATUS_NOT_SEND, now.Add(duration), now)
	return err
}

func (c *sendSMSDAO) Update(id uint64, submitStatus byte, msgID int64) error {
	db := getDB()
	if db == nil {
		return errDBNotInit
	}

	_, err := db.Exec(qSendSMSUpdate, submitStatus, msgID, id)
	return err
}

// SendSMSDAO instance of send sms dao
var SendSMSDAO ISendSMSDAO = &sendSMSDAO{}
