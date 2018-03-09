package dao

import (
	"github.com/linxGnu/gosmpp/examples/telcos/models"

	"fmt"
	"time"
)

const (
	qReceiveSMSSave = "INSERT INTO " + models.TableReceiveSMS + " (isdn, content, is_processed, process_expired, created_at) VALUES (?,?,?,?,?)"
)

// IReceiveSMSDAO interface of receive sms dao
type IReceiveSMSDAO interface {
	SaveReceiveSMS(isdn uint64, content string) error
}

type receiveSMSDAO struct{}

func (c *receiveSMSDAO) SaveReceiveSMS(isdn uint64, content string) error {
	db := getDB()
	if db == nil {
		return fmt.Errorf("DB not initialized")
	}

	now := time.Now()
	_, err := db.Exec(qReceiveSMSSave, isdn, content, false, now.Add(2*time.Minute), now)
	return err
}

// ReceiveSMSDAO instance of receive sms dao
var ReceiveSMSDAO IReceiveSMSDAO = &receiveSMSDAO{}
