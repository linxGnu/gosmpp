package SmppSession

import (
	"github.com/linxGnu/gosmpp/examples/telcos/Libs/Utils"
	"github.com/linxGnu/gosmpp/examples/telcos/config"
	"github.com/linxGnu/gosmpp/examples/telcos/consts"
	"github.com/linxGnu/gosmpp/examples/telcos/dao"
	"github.com/linxGnu/gosmpp/examples/telcos/models"

	"bytes"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/linxGnu/gosmpp"
	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/PDU"
	bb "github.com/linxGnu/gosmpp/Utils"
	cache "github.com/patrickmn/go-cache"
)

// ISmppSessionManager smpp session manager interface
type ISmppSessionManager interface {
	Bind() error
	Rebind() error
	Destroy()
	SendSMS(data *models.SendSMS) bool
	Unbind() (unbindResp *PDU.UnbindResp, err *Exception.Exception)
}

// Instance instance of SmppSessionManager
var Instance ISmppSessionManager
var lock sync.RWMutex

func newSmppSessionManager(gw *config.SmsGateway, pduConf *config.PDUConfigs) *SmppSessionManager {
	dur := time.Duration(gw.CacheDurForProcessStatusInSec) * time.Second

	res := &SmppSessionManager{smsGatewayConfigs: gw, pduConfigs: pduConf, cache: cache.New(dur, dur*2), wg: &sync.WaitGroup{}, firstBind: true, seq: gw.SeqSeed}

	// initialize concurrency handler
	if gw.NumConcurrentHandler <= 2 {
		gw.NumConcurrentHandler = 2 // minimum is 2
	}

	return res
}

// InitializeSmppSessionInstance ...
func InitializeSmppSessionInstance(gw *config.SmsGateway, pduConf *config.PDUConfigs) {
	if gw == nil {
		return
	}

	lock.Lock()
	defer lock.Unlock()

	if Instance != nil {
		Instance.Unbind()
	}

	Instance = newSmppSessionManager(gw, pduConf)
}

const (
	STATE_UNBIND = 0
	STATE_BIND   = 1
)

// SmppSessionManager default smpp session manager
type SmppSessionManager struct {
	session     *gosmpp.Session
	sessionLock sync.RWMutex

	smsGatewayConfigs *config.SmsGateway
	pduConfigs        *config.PDUConfigs
	cache             *cache.Cache

	// eventCh smpp server pdu event channel
	eventCh chan *gosmpp.ServerPDUEvent
	wg      *sync.WaitGroup

	// sms channel
	smsCh chan *models.SendSMS

	// enquireLink
	enquireLink     bool
	enquireLinkLock sync.RWMutex

	firstBind bool

	//
	state byte

	//
	seq     uint
	seqLock sync.RWMutex

	seq8     uint
	seq8Lock sync.RWMutex

	seq16     uint
	seq16Lock sync.RWMutex
}

// Unbind unbinding connections
func (c *SmppSessionManager) Unbind() (unbindResp *PDU.UnbindResp, err *Exception.Exception) {
	c.sessionLock.RLock()
	if c.session == nil {
		c.sessionLock.RUnlock()
		return
	}
	c.sessionLock.RUnlock()

	if c.eventCh != nil {
		go func(c *SmppSessionManager) {
			defer func() {
				if e := recover(); e != nil { // avoid close on closed channel
					fmt.Println(e)
				}
			}()

			close(c.eventCh)
		}(c)
	}

	if c.smsCh != nil {
		go func(c *SmppSessionManager) {
			defer func() {
				if e := recover(); e != nil { // avoid close on closed channel
					fmt.Println(e)
				}
			}()

			close(c.smsCh)
		}(c)
	}

	c.wg.Wait()

	// do session lock
	c.sessionLock.Lock()
	unbindResp, err = c.session.Unbind()
	if err != nil { // always log error
		// Logger.GetGlobalLogger().ErrorLog("Unbind: " + err.Error.Error())
	} else if c.smsGatewayConfigs.IsDebugMode {
		// Logger.GetGlobalLogger().InfoLog(fmt.Sprintf("UnbindResp: %+v", unbindResp))
	}
	c.session = nil
	c.state = STATE_UNBIND
	c.sessionLock.Unlock()

	return
}

// Destroy session
func (c *SmppSessionManager) Destroy() {
	c.enquireLinkLock.Lock()
	c.enquireLink = false
	c.enquireLinkLock.Unlock()

	c.Unbind()
}

// Bind safe binding
func (c *SmppSessionManager) Bind() error {
	if c.smsGatewayConfigs == nil {
		return fmt.Errorf("Configurations not initialized!")
	}

	c.sessionLock.Lock()
	defer c.sessionLock.Unlock()

	// initialize new connection
	connection, err := gosmpp.NewTCPIPConnectionWithAddrPort(c.smsGatewayConfigs.Addr, c.smsGatewayConfigs.Port)
	if err != nil {
		return err
	}

	// create session based on this connection
	c.session = gosmpp.NewSessionWithConnection(connection)
	c.session.EnableStateChecking()

	// create bind request
	request := PDU.NewBindTransceiver()
	request.SetSystemId(c.smsGatewayConfigs.SystemID)
	request.SetPassword(c.smsGatewayConfigs.Password)
	request.SetSystemType(c.smsGatewayConfigs.SystemType)

	// try to bind
	if resp, e := c.session.BindWithListener(request, c); e != nil || resp.GetCommandStatus() != 0 {
		// Logger.GetGlobalLogger().ErrorLog("Binding: " + e.Error.Error())

		connection.Close() // try to close connection()

		c.session = nil // set nil

		// return error
		return e.Error
	} else if c.smsGatewayConfigs.IsDebugMode {
		// Logger.GetGlobalLogger().InfoLog(fmt.Sprintf("BindResp: %+v", resp))
	}

	// no timeout
	c.session.GetReceiver().SetReceiveTimeout(-1)
	c.state = STATE_BIND

	// Start pdu event handlers
	c.eventCh = make(chan *gosmpp.ServerPDUEvent, c.smsGatewayConfigs.NumConcurrentHandler)
	for i := 0; i < c.smsGatewayConfigs.NumConcurrentHandler; i++ {
		c.wg.Add(1)
		go c.handleEventWorker()
	}

	// Start sms submit handler
	c.smsCh = make(chan *models.SendSMS, 2000)
	c.wg.Add(1)
	go c.submitWorker()

	// Start enquire link
	if c.firstBind {
		c.firstBind = false

		if c.smsGatewayConfigs.EnquiryIntervalInSec > 0 {
			c.enquireLink = true
			go c.doEnquireLink()
		}
	}

	return nil
}

// SendSMS do send sms
func (c *SmppSessionManager) SendSMS(data *models.SendSMS) (subRes bool) {
	defer func() {
		if e := recover(); e != nil { // prevent send to closed channel
			fmt.Println(e)

			subRes = false
		}
	}()

	if data == nil {
		return true
	}

	c.sessionLock.RLock()
	if c.state == STATE_UNBIND {
		c.sessionLock.RUnlock()
		return false
	}
	c.sessionLock.RUnlock()

	// do send to channel
	c.smsCh <- data

	return true
}

// submitWorker ...
func (c *SmppSessionManager) submitWorker() {
	defer c.wg.Done()

	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
		}
	}()

	var subSM *PDU.SubmitSM
	var trunkLength int
	s, e, ln, firstSMSInSeq, isUnicode := 0, 0, 0, true, true
	var sq uint

	for v := range c.smsCh {
		if v == nil || v.Submit_expired.Before(time.Now()) {
			continue
		}

		// now do send sms
		isUnicode, ln = Utils.IsUnicode(v.Content), len(v.Content)
		if (isUnicode && ln <= 70) || (!isUnicode && ln <= 160) {
			subSM = c.createSubmitSM(v.Alias, strconv.FormatUint(v.Isdn, 10), v.Content, isUnicode)

			// cache this seq
			seq := int32(c.incSeq())
			c.cache.Set(strconv.FormatUint(uint64(seq), 10), v, time.Duration(c.smsGatewayConfigs.CacheDurForProcessStatusInSec)*time.Second)

			// set sequence number
			subSM.SetSequenceNumber(seq)

			if c.smsGatewayConfigs.IsDebugMode {
				// Logger.GetGlobalLogger().ErrorLog(fmt.Sprintf("SUBMIT_SM: %+v", subSM))
			}

			c.sessionLock.RLock()
			c.session.Submit(subSM)
			c.sessionLock.RUnlock()
		} else {
			s, e, firstSMSInSeq = 0, 0, true
			des := strconv.FormatUint(v.Isdn, 10)

			content := c.getLongText(v.Content)

			if isUnicode {
				trunkLength = 139
			} else {
				trunkLength = 140
			}

			for ln = len(content); e < ln; {
				if s+trunkLength < ln {
					e = s + trunkLength
				} else {
					e = ln
				}

				subSM = c.createSubmitSM(v.Alias, des, "", isUnicode)
				subSM.SetShortMessageData(bb.NewBuffer(content[s:e]))
				subSM.SetEsmClass(64)

				if c.smsGatewayConfigs.IsDebugMode {
					// Logger.GetGlobalLogger().ErrorLog(fmt.Sprintf("SUBMIT_SM: %+v", subSM))
				}

				if firstSMSInSeq {
					sq = c.incSeq()

					c.sessionLock.RLock()
					if _, e := c.session.Submit(subSM); e == nil {
						c.cache.Set(strconv.FormatUint(uint64(sq), 10), v, time.Duration(c.smsGatewayConfigs.CacheDurForProcessStatusInSec)*time.Second)
						firstSMSInSeq = false
					} else {
						// Logger.GetGlobalLogger().ErrorLog(fmt.Sprintf("SUBMIT_FAIL: %+v %+v", e, resp))
					}
					c.sessionLock.RUnlock()
				}

				s = e
			}
		}

		time.Sleep(time.Duration(c.smsGatewayConfigs.SendSMSIntervalInMiliSec) * time.Millisecond)
	}
}

// handleEventWorker ...
func (c *SmppSessionManager) handleEventWorker() {
	defer c.wg.Done()
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
		}
	}()

	var msgID int64
	var commandStatus int32
	var status byte
	var content string

	for event := range c.eventCh {
		if event == nil {
			continue
		}

		t := event.GetPDU()
		if t == nil {
			continue
		}

		commandStatus = t.GetCommandStatus()
		if commandStatus == Data.ESME_RTHROTTLED {
			// Logger.GetGlobalLogger().ErrorLog(fmt.Sprintf("THROTTLED: %+v", t))
		} else if commandStatus == Data.ESME_RMSGQFUL {
			// Logger.GetGlobalLogger().ErrorLog(fmt.Sprintf("QUEUE_FULL: %+v", t))
		}

		switch v := t.(type) {
		case *PDU.DeliverSM:
			if c.smsGatewayConfigs.IsDebugMode {
				// Logger.GetGlobalLogger().InfoLog(fmt.Sprintf("DELIVERSM: %+v", v))
			}

			if response, err := v.GetResponse(); err == nil {
				c.session.Respond(response)
			}

			if v.GetEsmClass() == 0 { // receive customer message
				content = ""
				if v.GetDataCoding() == 0 {
					content, _ = v.GetShortMessage()
				} else {
					content, _ = v.GetShortMessageWithEncoding(Data.ENC_UTF16)
				}

				if v.GetSourceAddr() != nil {
					isdn := Utils.StandardizePhone(v.GetSourceAddr().GetAddress())
					if pn, e := strconv.ParseUint(isdn, 10, 64); e == nil {
						if e = dao.ReceiveSMSDAO.SaveReceiveSMS(pn, content); e != nil { // rarely happend: response sms to customer that we fail to save
							now := time.Now()
							c.SendSMS(&models.SendSMS{
								Alias:          c.smsGatewayConfigs.Alias,
								Content:        consts.InternalSystemError,
								Created_at:     now,
								Isdn:           pn,
								Submit_expired: now.Add(time.Duration(c.smsGatewayConfigs.SendSMSExpiredInMinute) * time.Minute),
								Submit_status:  models.SUBMIT_STATUS_NOT_SEND,
							})
						}
					}
				}
			}
		case *PDU.Unbind:
			// Logger.GetGlobalLogger().InfoLog("UNBIND: received")

			if response, err := v.GetResponse(); err == nil {
				c.session.Respond(response)
			}
		case *PDU.EnquireLink:
			if response, err := v.GetResponse(); err == nil {
				c.session.Respond(response)
			}
		case *PDU.SubmitSMResp:
			seq := strconv.Itoa(int(v.GetSequenceNumber()))
			if tmp, ok := c.cache.Get(seq); ok {
				switch sms := tmp.(type) {
				case *models.SendSMS:
					msgID = -1
					if _msgID, err := strconv.ParseInt(v.GetMessageId(), 16, 64); err != nil {
						msgID = -1
					} else {
						msgID = _msgID
					}

					if v.GetCommandStatus() == Data.ESME_ROK {
						status = models.SUBMIT_STATUS_SUCC
					} else {
						status = models.SUBMIT_STATUS_FAIL
					}

					// try to update to database
					dao.SendSMSDAO.Update(sms.Id, status, msgID)
				}
			}
		}
	}
}

// Rebind ...
func (c *SmppSessionManager) Rebind() error {
	// Unbind first
	c.Unbind()

	// Do rebind
	return c.Bind()
}

// HandleEvent handle server pdu events
func (c *SmppSessionManager) HandleEvent(event *gosmpp.ServerPDUEvent) (ex *Exception.Exception) {
	defer func() {
		if e := recover(); e != nil {
			// panic caused by sending to closed channel
			// Logger.GetGlobalLogger().ErrorLog(e)

			ex = Exception.NewException(fmt.Errorf("%v", e))
		}
	}()

	c.eventCh <- event
	return
}

func (c *SmppSessionManager) doEnquireLink() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
		}
	}()

	for {
		// check enquire link state
		c.enquireLinkLock.RLock()
		if !c.enquireLink {
			c.enquireLinkLock.RUnlock()
			return
		}
		c.enquireLinkLock.RUnlock()

		// do enquire link
		c.sessionLock.RLock()
		if c.session != nil {
			if _, err := c.session.EnquireLink(PDU.NewEnquireLink()); err != nil {
				c.sessionLock.RUnlock()

				if c.smsGatewayConfigs.IsDebugMode {
					// Logger.GetGlobalLogger().ErrorLog(fmt.Sprintf("TryingEnquire: %+v", err))
				}

				c.Rebind() // rebind now

				goto sleepEnq // sleep before next enquirying
			} else if c.smsGatewayConfigs.ShowHeartBeat && c.smsGatewayConfigs.IsDebugMode {
				// Logger.GetGlobalLogger().InfoLog(fmt.Sprintf("Heartbeat: %+v", resp))
			}
		} else {
			c.sessionLock.RUnlock()

			if c.smsGatewayConfigs.IsDebugMode {
				// Logger.GetGlobalLogger().ErrorLog(fmt.Sprintf("TryingEnquire: Session is nil/closed"))
			}

			time.Sleep(10 * time.Second) // sleep before rebinding
			c.Rebind()                   // rebind now

			goto sleepEnq // sleep before next enquirying
		}
		c.sessionLock.RUnlock()

	sleepEnq:
		time.Sleep(time.Duration(c.smsGatewayConfigs.EnquiryIntervalInSec) * time.Second) // sleep to wait next enquire link
	}
}

func (c *SmppSessionManager) createSubmitSM(sender, receiver, content string, isUnicode bool) *PDU.SubmitSM {
	res := PDU.NewSubmitSM()

	srcAddr := PDU.NewAddress()
	srcAddr.SetAddress(sender)
	srcAddr.SetTon(c.pduConfigs.SourceAddr.Ton)
	srcAddr.SetNpi(c.pduConfigs.SourceAddr.Npi)
	res.SetSourceAddr(srcAddr)

	desAddr := PDU.NewAddress()
	desAddr.SetAddress(receiver)
	desAddr.SetTon(c.pduConfigs.DestinationAddr.Ton)
	desAddr.SetNpi(c.pduConfigs.DestinationAddr.Npi)
	res.SetDestAddr(desAddr)

	res.SetProtocolId(c.pduConfigs.ProtocolID)
	res.SetRegisteredDelivery(c.pduConfigs.RegisteredDelivery)
	res.SetReplaceIfPresentFlag(c.pduConfigs.ReplaceIfPresentFlag)
	res.SetEsmClass(c.pduConfigs.EsmClass)

	if len(content) > 0 {
		res.SetShortMessage(content)
	}

	if isUnicode {
		res.SetDataCoding(consts.DataCodingUCS2)
	} else {
		res.SetDataCoding(consts.DataCodingASCII)
	}

	return res
}

func (c *SmppSessionManager) getLongText(data string) []byte {
	if isUnicode := Utils.IsUnicode(data); isUnicode {
		return c.getLongText16bit(data, int(c.incSeq16()))
	}

	return c.getLongText8bit([]byte(data), int(c.incSeq8()))
}

// IncSeq increase sequence
func (c *SmppSessionManager) incSeq() uint {
	c.seqLock.Lock()
	defer c.seqLock.Unlock()

	c.seq += c.smsGatewayConfigs.SeqMod
	if c.seq > 2000000000 {
		c.seq = c.smsGatewayConfigs.SeqSeed
	}

	return c.seq
}

// IncSeq8 ...
func (c *SmppSessionManager) incSeq8() uint {
	c.seq8Lock.Lock()
	defer c.seq8Lock.Unlock()

	if c.seq8++; c.seq8 > 127 {
		c.seq8 = 1
	}

	return c.seq8
}

// IncSeq16 ...
func (c *SmppSessionManager) incSeq16() uint {
	c.seq16Lock.Lock()
	defer c.seq16Lock.Unlock()

	if c.seq16++; c.seq16 > 32767 {
		c.seq16 = 1
	}

	return c.seq16
}

func (c *SmppSessionManager) getLongText8bit(data []byte, ref int) []byte {
	if data == nil {
		return nil
	}

	// Sent ref number 8-bit
	msgLength := len(data)

	// number of messages
	nNumberOfMsg := (msgLength-1)/134 + 1

	// create bytes buffer
	buf := bytes.NewBuffer(make([]byte, msgLength<<1))

	if ref > 127 {
		ref = 1
	}

	index, byteData := 0, 0
	for i := 0; i < nNumberOfMsg; i++ {
		buf.Write([]byte{5, 0, 3, byte(ref), byte(nNumberOfMsg), byte(i + 1)})
		if byteData = msgLength - index; byteData > 134 {
			byteData = 134
		}
		buf.Write(data[index : index+byteData])
		index += byteData
	}

	return buf.Bytes()
}

func (c *SmppSessionManager) getLongText16bit(_data string, ref int) []byte {
	data, err := Data.ENC_UTF16.Encode(_data)
	if err != nil {
		return nil
	}

	msgLength := len(data)

	// number of messages
	nNumberOfMsg := (msgLength-1)/65 + 1

	// create bytes buffer
	buf := bytes.NewBuffer(make([]byte, msgLength<<1))

	if ref > 32768 {
		ref = 1
	}

	index, byteData := 2, 0
	for i := 0; i < nNumberOfMsg; i++ {
		buf.Write([]byte{6, 8, 4, byte((ref >> 8) & 255), byte(ref & 255), byte(nNumberOfMsg), byte(i + 1)})
		if byteData = msgLength - index; byteData > 132 {
			byteData = 132
		}
		buf.Write(data[index : index+byteData])
		index += byteData
	}

	return buf.Bytes()
}
