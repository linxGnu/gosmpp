package daemons

import (
	"github.com/linxGnu/gosmpp/examples/telcos/controller/SmppSession"
	"github.com/linxGnu/gosmpp/examples/telcos/dao"
	"github.com/linxGnu/gosmpp/examples/telcos/models"

	"time"
)

// SmsSenderDaemon sms sender daemon
func SmsSenderDaemon() {
	defer wg.Done()

	for {
		lock.RLock()
		if shouldStop {
			lock.RUnlock()
			return
		}
		lock.RUnlock()

		// select not send sms from database
		if notSends, err := dao.SendSMSDAO.SelectNotSend(); err == nil && len(notSends) > 0 {
			for _, v := range notSends {
				if SmppSession.Instance != nil && SmppSession.Instance.SendSMS(v) {
					dao.SendSMSDAO.Update(v.Id, models.SUBMIT_STATUS_SUBMITTED, -1)
				}
			}

			time.Sleep(1 * time.Second)
			continue
		}

		// wait 5 second before next checking
		time.Sleep(5 * time.Second)
	}
}
