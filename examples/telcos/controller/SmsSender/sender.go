package SmsSender

import (
	"github.com/linxGnu/gosmpp/examples/telcos/Libs/Utils"
	"github.com/linxGnu/gosmpp/examples/telcos/config"
	"github.com/linxGnu/gosmpp/examples/telcos/controller"
	"github.com/linxGnu/gosmpp/examples/telcos/dao"

	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

// SendSMS send sms handler
func SendSMS(c echo.Context) (erro error) {
	erro = nil
	response := controller.Response{}
	defer func() {
		if erro == nil {
			erro = c.JSON(http.StatusOK, response)
		}
	}()

	alias, content := controller.Sanitize(c.FormValue("alias")), controller.Sanitize(c.FormValue("content"))

	var isdn uint64
	if _isdn, e := strconv.ParseUint(Utils.StandardizePhone(c.FormValue("isdn")), 10, 64); e == nil {
		isdn = _isdn
	} else {
		response.SetCodeMessage(http.StatusBadRequest, "ISDN Invalid!")
		return
	}

	// Save to database
	if e := dao.SendSMSDAO.Save(alias, isdn, content, time.Duration(config.GetConfigurations().SmsGateway.SendSMSExpiredInMinute)*time.Minute); e != nil {
		response.SetCodeMessage(http.StatusInternalServerError, e.Error())
		return
	}

	response.SetCodeMessage(http.StatusOK, "")
	return
}
