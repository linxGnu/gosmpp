package gosmpp

import (
	"fmt"
	"time"

	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/PDU"
)

type Session struct {
	opened             bool
	bound              bool
	disallowUnknownPDU bool
	state              int32
	stateChecking      bool
	sessionType        int32
	connection         IConnection
	transmitter        *Transmitter
	receiver           *Receiver
	pduListener        ServerPDUEventListener
	isAsync            bool
}

const (
	STATE_NOT_ALLOWED int32 = 0x00
	STATE_CLOSED      int32 = 0x01
	STATE_OPENED      int32 = 0x02
	STATE_TRANSMITTER int32 = 0x04
	STATE_RECEIVER    int32 = 0x08
	STATE_TRANSCEIVER int32 = 0x10
	STATE_ALWAYS      int32 = STATE_OPENED | STATE_TRANSMITTER | STATE_RECEIVER | STATE_TRANSCEIVER
	TYPE_ESME         int32 = 1
	TYPE_MC           int32 = 2
)

var esmeStateMatrix map[int]int = make(map[int]int)
var mcStateMatrix map[int]int = make(map[int]int)
var isMatrixInit = false

func addValidState(m map[int]int, k, v int32) {
	m[int(k)] = int(v)
}

func getStateMatrix(t int) map[int]int {
	if !isMatrixInit {
		addValidState(esmeStateMatrix, Data.BIND_TRANSMITTER, STATE_CLOSED)
		addValidState(esmeStateMatrix, Data.BIND_TRANSMITTER_RESP, STATE_NOT_ALLOWED)
		addValidState(esmeStateMatrix, Data.BIND_RECEIVER, STATE_CLOSED)
		addValidState(esmeStateMatrix, Data.BIND_RECEIVER_RESP, STATE_NOT_ALLOWED)
		addValidState(esmeStateMatrix, Data.BIND_TRANSCEIVER, STATE_CLOSED)
		addValidState(esmeStateMatrix, Data.BIND_TRANSCEIVER_RESP, STATE_NOT_ALLOWED)
		addValidState(esmeStateMatrix, Data.OUTBIND, STATE_NOT_ALLOWED)
		addValidState(esmeStateMatrix, Data.UNBIND, STATE_TRANSMITTER|STATE_RECEIVER|STATE_TRANSCEIVER)
		addValidState(esmeStateMatrix, Data.UNBIND_RESP, STATE_TRANSMITTER|STATE_RECEIVER|STATE_TRANSCEIVER)
		addValidState(esmeStateMatrix, Data.SUBMIT_SM, STATE_TRANSMITTER|STATE_TRANSCEIVER)
		addValidState(esmeStateMatrix, Data.SUBMIT_SM_RESP, STATE_NOT_ALLOWED)
		addValidState(esmeStateMatrix, Data.SUBMIT_MULTI, STATE_TRANSMITTER|STATE_TRANSCEIVER)
		addValidState(esmeStateMatrix, Data.SUBMIT_MULTI_RESP, STATE_NOT_ALLOWED)
		addValidState(esmeStateMatrix, Data.DATA_SM, STATE_TRANSMITTER|STATE_TRANSCEIVER)
		addValidState(esmeStateMatrix, Data.DATA_SM_RESP, STATE_RECEIVER|STATE_TRANSCEIVER)
		addValidState(esmeStateMatrix, Data.DELIVER_SM, STATE_NOT_ALLOWED)
		addValidState(esmeStateMatrix, Data.DELIVER_SM_RESP, STATE_RECEIVER|STATE_TRANSCEIVER)
		addValidState(esmeStateMatrix, Data.QUERY_SM, STATE_TRANSMITTER|STATE_TRANSCEIVER)
		addValidState(esmeStateMatrix, Data.QUERY_SM_RESP, STATE_NOT_ALLOWED)
		addValidState(esmeStateMatrix, Data.CANCEL_SM, STATE_TRANSMITTER|STATE_TRANSCEIVER)
		addValidState(esmeStateMatrix, Data.CANCEL_SM_RESP, STATE_NOT_ALLOWED)
		addValidState(esmeStateMatrix, Data.REPLACE_SM, STATE_TRANSMITTER|STATE_TRANSCEIVER)
		addValidState(esmeStateMatrix, Data.REPLACE_SM_RESP, STATE_NOT_ALLOWED)
		addValidState(esmeStateMatrix, Data.ENQUIRE_LINK, STATE_TRANSMITTER|STATE_RECEIVER|STATE_TRANSCEIVER)
		// STATE_ALWAYS);
		addValidState(esmeStateMatrix, Data.ENQUIRE_LINK_RESP, STATE_TRANSMITTER|STATE_RECEIVER|STATE_TRANSCEIVER)
		// STATE_ALWAYS);
		addValidState(esmeStateMatrix, Data.ALERT_NOTIFICATION, STATE_NOT_ALLOWED)
		addValidState(esmeStateMatrix, Data.GENERIC_NACK, STATE_TRANSMITTER|STATE_RECEIVER|STATE_TRANSCEIVER)
		// STATE_ALWAYS);

		addValidState(mcStateMatrix, Data.BIND_TRANSMITTER, STATE_NOT_ALLOWED)
		addValidState(mcStateMatrix, Data.BIND_TRANSMITTER_RESP, STATE_TRANSMITTER|STATE_RECEIVER|STATE_TRANSCEIVER)
		// STATE_OPENED);
		addValidState(mcStateMatrix, Data.BIND_RECEIVER, STATE_NOT_ALLOWED)
		addValidState(mcStateMatrix, Data.BIND_RECEIVER_RESP, STATE_TRANSMITTER|STATE_RECEIVER|STATE_TRANSCEIVER)
		// STATE_OPENED);
		addValidState(mcStateMatrix, Data.BIND_TRANSCEIVER, STATE_NOT_ALLOWED)
		addValidState(mcStateMatrix, Data.BIND_TRANSCEIVER_RESP, STATE_TRANSMITTER|STATE_RECEIVER|STATE_TRANSCEIVER)
		// STATE_OPENED);
		addValidState(mcStateMatrix, Data.OUTBIND, STATE_TRANSMITTER|STATE_RECEIVER|STATE_TRANSCEIVER)
		// STATE_OPENED);
		addValidState(mcStateMatrix, Data.UNBIND, STATE_TRANSMITTER|STATE_RECEIVER|STATE_TRANSCEIVER)
		addValidState(mcStateMatrix, Data.UNBIND_RESP, STATE_TRANSMITTER|STATE_RECEIVER|STATE_TRANSCEIVER)
		addValidState(mcStateMatrix, Data.SUBMIT_SM, STATE_NOT_ALLOWED)
		addValidState(mcStateMatrix, Data.SUBMIT_SM_RESP, STATE_TRANSMITTER|STATE_TRANSCEIVER)
		addValidState(mcStateMatrix, Data.SUBMIT_MULTI, STATE_NOT_ALLOWED)
		addValidState(mcStateMatrix, Data.SUBMIT_MULTI_RESP, STATE_TRANSMITTER|STATE_TRANSCEIVER)
		addValidState(mcStateMatrix, Data.DATA_SM, STATE_RECEIVER|STATE_TRANSCEIVER)
		addValidState(mcStateMatrix, Data.DATA_SM_RESP, STATE_TRANSMITTER|STATE_TRANSCEIVER)
		addValidState(mcStateMatrix, Data.DELIVER_SM, STATE_RECEIVER|STATE_TRANSCEIVER)
		addValidState(mcStateMatrix, Data.DELIVER_SM_RESP, STATE_NOT_ALLOWED)
		addValidState(mcStateMatrix, Data.QUERY_SM, STATE_NOT_ALLOWED)
		addValidState(mcStateMatrix, Data.QUERY_SM_RESP, STATE_TRANSMITTER|STATE_TRANSCEIVER)
		addValidState(mcStateMatrix, Data.CANCEL_SM, STATE_NOT_ALLOWED)
		addValidState(mcStateMatrix, Data.CANCEL_SM_RESP, STATE_TRANSMITTER|STATE_TRANSCEIVER)
		addValidState(mcStateMatrix, Data.REPLACE_SM, STATE_NOT_ALLOWED)
		addValidState(mcStateMatrix, Data.REPLACE_SM_RESP, STATE_TRANSMITTER|STATE_TRANSCEIVER)
		addValidState(mcStateMatrix, Data.ENQUIRE_LINK, STATE_TRANSMITTER|STATE_RECEIVER|STATE_TRANSCEIVER)
		// STATE_ALWAYS);
		addValidState(mcStateMatrix, Data.ENQUIRE_LINK_RESP, STATE_TRANSMITTER|STATE_RECEIVER|STATE_TRANSCEIVER)
		// STATE_ALWAYS);
		addValidState(mcStateMatrix, Data.ALERT_NOTIFICATION, STATE_RECEIVER|STATE_TRANSCEIVER)
		addValidState(mcStateMatrix, Data.GENERIC_NACK, STATE_TRANSMITTER|STATE_RECEIVER|STATE_TRANSCEIVER)
		isMatrixInit = true
	}

	if t == int(TYPE_ESME) {
		return esmeStateMatrix
	}

	if t == int(TYPE_MC) {
		return mcStateMatrix
	}

	return nil
}

func newSession() *Session {
	a := &Session{}
	a.bound = false
	a.opened = false
	a.disallowUnknownPDU = false
	a.state = STATE_CLOSED
	a.stateChecking = false
	a.sessionType = TYPE_ESME
	a.isAsync = false

	return a
}

func NewSessionWithConnection(conn IConnection) *Session {
	if conn == nil {
		return nil
	}

	a := newSession()
	a.connection = conn

	return a
}

func (c *Session) Open() *Exception.Exception {
	if !c.opened {
		err := c.connection.Open()
		if err != nil {
			return err
		}

		c.opened = true
		c.setState(STATE_OPENED)
	}

	return nil
}

func (c *Session) Close() *Exception.Exception {
	if c.IsOpened() {
		c.connection.Close()
		if c.receiver != nil {
			c.receiver.Stop()
		}

		c.opened = false
	}

	c.bound = false
	c.setState(STATE_CLOSED)

	return nil
}

func (c *Session) IsOpened() bool {
	return c.opened && c.connection.IsOpened()
}

func (c *Session) IsBound() bool {
	return c.bound
}

func (c *Session) setState(state int32) {
	c.state = state
}

func (c *Session) SetType(t int32) {
	c.sessionType = t
}

func (c *Session) GetType() int32 {
	return c.sessionType
}

func (c *Session) Bind(req PDU.IBindRequest) (PDU.IResponse, *Exception.Exception) {
	err := c.CheckPDUState(req)
	if err != nil {
		return nil, err
	}

	return c.BindWithListener(req, nil)
}

func (c *Session) BindWithListener(req PDU.IBindRequest, pduListener ServerPDUEventListener) (bindResp PDU.IResponse, err *Exception.Exception) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
		}
	}()

	err = c.CheckPDUState(req)
	if err != nil {
		return
	}

	if c.bound {
		return nil, nil
	}

	c.Open()
	c.transmitter = NewTransmitterWithConnection(c.connection)
	c.receiver = NewReceiverWithTransmitterCon(c.transmitter, c.connection)

	resp, err := c.SendWithAsyncStyle(req, false)
	if err != nil {
		return
	}

	if resp != nil {
		bindResp = resp.(PDU.IResponse)
	}

	c.bound = bindResp != nil && bindResp.GetCommandStatus() == Data.ESME_ROK
	if !c.bound {
		c.Close()
	} else {
		c.receiver.Start()
		if req.IsTransmitter() {
			if req.IsReceiver() {
				c.setState(STATE_TRANSCEIVER)
			} else {
				c.setState(STATE_TRANSMITTER)
			}
		} else {
			c.setState(STATE_RECEIVER)
		}
		c.setServerPDUEventListener(pduListener)
	}

	return bindResp, nil
}

func (c *Session) Unbind() (unbindResp *PDU.UnbindResp, err *Exception.Exception) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
		}
	}()

	if c.bound {
		unbindReq := PDU.NewUnbind()
		err = c.CheckPDUState(unbindReq)
		if err != nil {
			return
		}

		origListener := c.getServerPDUEventListener()
		if c.isAsync {
			unbindReq.AssignSequenceNumber()

			unbindListener := NewUnbindServerPDUEventListener(c, origListener, unbindReq)

			c.setServerPDUEventListener(unbindListener)
			c.Send(unbindReq)

			unbindListener.StartWait(Data.UNBIND_RECEIVE_TIMEOUT)
			<-unbindListener.GetWaitChan()
			unbindListener.CloseWaitChan()

			unbindResp = unbindListener.GetUnbindResp()

			// Reset listener
			c.setServerPDUEventListener(origListener)
		} else {
			_unbindResp, _err := c.Send(unbindReq)
			if _err != nil {
				err = _err
				return
			}

			if _unbindResp != nil {
				unbindResp = _unbindResp.(*PDU.UnbindResp)
			}
		}

		c.Close()
	}

	return unbindResp, nil
}

func (c *Session) Submit(req *PDU.SubmitSM) (resp *PDU.SubmitSMResp, err *Exception.Exception) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewExceptionFromStr(fmt.Sprintf("%v", errs))
		}
	}()

	err = c.CheckPDUState(req)
	if err != nil {
		return nil, err
	}

	tmp, err := c.Send(req)
	if err != nil {
		return nil, err
	}

	if tmp != nil {
		resp = tmp.(*PDU.SubmitSMResp)
	}

	return
}

func (c *Session) SubmitMulti(req *PDU.SubmitMultiSM) (resp *PDU.SubmitMultiSMResp, err *Exception.Exception) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewExceptionFromStr(fmt.Sprintf("%v", errs))
		}
	}()

	err = c.CheckPDUState(req)
	if err != nil {
		return nil, err
	}

	tmp, err := c.Send(req)
	if err != nil {
		return nil, err
	}

	if tmp != nil {
		resp = tmp.(*PDU.SubmitMultiSMResp)
	}

	return
}

func (c *Session) Deliver(req *PDU.DeliverSM) (resp *PDU.DeliverSMResp, err *Exception.Exception) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewExceptionFromStr(fmt.Sprintf("%v", errs))
		}
	}()

	err = c.CheckPDUState(req)
	if err != nil {
		return nil, err
	}

	tmp, err := c.Send(req)
	if err != nil {
		return nil, err
	}

	if tmp != nil {
		resp = tmp.(*PDU.DeliverSMResp)
	}

	return
}

func (c *Session) Data(req *PDU.DataSM) (resp *PDU.DataSMResp, err *Exception.Exception) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewExceptionFromStr(fmt.Sprintf("%v", errs))
		}
	}()

	err = c.CheckPDUState(req)
	if err != nil {
		return nil, err
	}

	tmp, err := c.Send(req)
	if err != nil {
		return nil, err
	}

	if tmp != nil {
		resp = tmp.(*PDU.DataSMResp)
	}

	return
}

func (c *Session) Query(req *PDU.QuerySM) (resp *PDU.QuerySMResp, err *Exception.Exception) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewExceptionFromStr(fmt.Sprintf("%v", errs))
		}
	}()

	err = c.CheckPDUState(req)
	if err != nil {
		return nil, err
	}

	tmp, err := c.Send(req)
	if err != nil {
		return nil, err
	}

	if tmp != nil {
		resp = tmp.(*PDU.QuerySMResp)
	}

	return
}

func (c *Session) Cancel(req *PDU.CancelSM) (resp *PDU.CancelSMResp, err *Exception.Exception) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewExceptionFromStr(fmt.Sprintf("%v", errs))
		}
	}()

	err = c.CheckPDUState(req)
	if err != nil {
		return nil, err
	}

	tmp, err := c.Send(req)
	if err != nil {
		return nil, err
	}

	if tmp != nil {
		resp = tmp.(*PDU.CancelSMResp)
	}

	return
}

func (c *Session) Replace(req *PDU.ReplaceSM) (resp *PDU.ReplaceSMResp, err *Exception.Exception) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewExceptionFromStr(fmt.Sprintf("%v", errs))
		}
	}()

	err = c.CheckPDUState(req)
	if err != nil {
		return nil, err
	}

	tmp, err := c.Send(req)
	if err != nil {
		return nil, err
	}

	if tmp != nil {
		resp = tmp.(*PDU.ReplaceSMResp)
	}

	return
}

func (c *Session) EnquireLink(req *PDU.EnquireLink) (resp *PDU.EnquireLinkResp, err *Exception.Exception) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewExceptionFromStr(fmt.Sprintf("%v", errs))
		}
	}()

	err = c.CheckPDUState(req)
	if err != nil {
		return nil, err
	}

	tmp, err := c.Send(req)
	if err != nil {
		return nil, err
	}

	if tmp != nil {
		resp = tmp.(*PDU.EnquireLinkResp)
	}

	return
}

func (c *Session) DoEnquireLink() (resp *PDU.EnquireLinkResp, err *Exception.Exception) {
	return c.EnquireLink(PDU.NewEnquireLink())
}

func (c *Session) AlertNotification(req *PDU.AlertNotification) (err *Exception.Exception) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewExceptionFromStr(fmt.Sprintf("%v", errs))
		}
	}()

	err = c.CheckPDUState(req)
	if err != nil {
		return
	}

	_, err = c.Send(req)
	return
}

func (c *Session) Receive() (pdu PDU.IPDU, err *Exception.Exception) {
	if !c.isAsync {
		return c.ReceiveWTimeout(Data.RECEIVE_BLOCKING)
	}

	err = Exception.NotSynchronousException
	return
}

func (c *Session) ReceiveWTimeout(timeout int64) (pdu PDU.IPDU, err *Exception.Exception) {
	if !c.isAsync {
		pdu, err = c.receiver.ReceiveSyncWTimeout(timeout)
	} else {
		err = Exception.NotSynchronousException
	}

	return pdu, nil
}

func (c *Session) Respond(resp PDU.IResponse) (err *Exception.Exception) {
	err = c.CheckPDUState(resp)
	if err != nil {
		return
	}

	err = c.transmitter.Send(resp)
	return
}

func (c *Session) GetTransmitter() *Transmitter {
	return c.transmitter
}

func (c *Session) GetReceiver() *Receiver {
	return c.receiver
}

func (c *Session) GetConnection() IConnection {
	return c.connection
}

func (c *Session) Send(req PDU.IRequest) (resp PDU.IResponse, err *Exception.Exception) {
	return c.SendWithAsyncStyle(req, c.isAsync)
}

func (c *Session) SendWithAsyncStyle(req PDU.IRequest, isAsync bool) (resp PDU.IResponse, err *Exception.Exception) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
		}
	}()

	resp = nil

	err = c.transmitter.Send(req)
	if err != nil {
		return
	}

	if !isAsync && req.CanResponse() {
		expected, e := req.GetResponse()
		if e != nil {
			err = Exception.NewException(e)
			return
		}

		var pdu PDU.IPDU

		pdu, err = c.receiver.ReceiveSyncWithExpectedPDU(expected)
		if err == Exception.UnknownCommandIdException {
			c.safeGenericNack(Data.ESME_RINVCMDID, pdu.GetSequenceNumber())
		} else if err == Exception.TerminatingZeroNotFoundException || err == Exception.NotEnoughDataInByteBufferException {
			c.safeGenericNack(Data.ESME_RINVMSGLEN, pdu.GetSequenceNumber())
		} else if err != nil {
			return
		}

		if pdu != nil {
			resp, err = c.checkResponse(pdu, expected)
			return
		}
	}

	return
}

func (c *Session) checkResponse(resp PDU.IPDU, exp PDU.IResponse) (result PDU.IResponse, err *Exception.Exception) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewExceptionFromStr(fmt.Sprintf("%v", errs))
		}
	}()

	if resp.GetCommandId() != exp.GetCommandId() {
		if resp.GetCommandId() == Data.GENERIC_NACK {
			exp.SetCommandId(Data.GENERIC_NACK)
			exp.SetCommandLength(resp.GetCommandLength())
			exp.SetCommandStatus(resp.GetCommandStatus())
			exp.SetSequenceNumber(resp.GetSequenceNumber())
			return exp, nil
		} else {
			err := c.safeGenericNack(Data.ESME_RINVCMDID, resp.GetSequenceNumber())
			if err != nil {
				return nil, err
			}

			return nil, nil
		}
	} else {
		return resp.(PDU.IResponse), nil
	}
}

func (c *Session) safeGenericNack(commandStatus, sequenceNumber int32) *Exception.Exception {
	return c.GenericNackWithCmStatusSeqNum(commandStatus, sequenceNumber)
}

func (c *Session) GenericNack(resp *PDU.GenericNack) *Exception.Exception {
	err := c.CheckPDUState(resp)
	if err != nil {
		return err
	}

	return c.Respond(resp)
}

func (c *Session) GenericNackWithCmStatusSeqNum(commandStatus, sequenceNumber int32) *Exception.Exception {
	gnack := PDU.NewGenericNackWithCmStatusSeqNum(commandStatus, sequenceNumber)
	err := c.CheckPDUState(gnack)
	if err != nil {
		return err
	}

	return c.GenericNack(gnack)
}

func (c *Session) GetState() int32 {
	return c.state
}

func (c *Session) CheckPDUState(pdu PDU.IPDU) *Exception.Exception {
	if pdu == nil {
		return Exception.ValueNotSetException
	}

	if c.stateChecking {
		pduMatrix := getStateMatrix(int(c.sessionType))
		if pduMatrix == nil {
			if c.disallowUnknownPDU {
				return Exception.WrongSessionStateException
			}
		} else {
			cId := pdu.GetCommandId()
			r, ok := pduMatrix[int(cId)]
			if ok {
				return c.CheckState(int32(r))
			}
		}
	}

	return nil
}

func (c *Session) CheckState(requestedState int32) *Exception.Exception {
	if c.stateChecking {
		if c.state&requestedState == 0 {
			return Exception.WrongSessionStateException
		}
	}

	return nil
}

func (c *Session) EnableStateChecking() {
	c.stateChecking = true
}

func (c *Session) DisableStateChecking() {
	c.stateChecking = false
}

func (c *Session) IsStateAllowed(requestedState int32) bool {
	return c.CheckState(requestedState) == nil
}

func (c *Session) IsPDUAllowed(pdu PDU.IPDU) bool {
	return c.CheckPDUState(pdu) == nil
}

func (c *Session) setServerPDUEventListener(pduListener ServerPDUEventListener) {
	c.pduListener = pduListener
	c.receiver.SetServerPDUEventListener(pduListener)
	c.isAsync = pduListener != nil
}

func (c *Session) getServerPDUEventListener() ServerPDUEventListener {
	return c.pduListener
}

type UnbindServerPDUEventListener struct {
	session      *Session
	origListener ServerPDUEventListener
	unbindReq    *PDU.Unbind
	expectedResp *PDU.UnbindResp
	unbindResp   *PDU.UnbindResp
	waitChan     chan bool
}

func NewUnbindServerPDUEventListener(sess *Session, origListener ServerPDUEventListener, unbindReq *PDU.Unbind) *UnbindServerPDUEventListener {
	a := &UnbindServerPDUEventListener{}
	a.session = sess
	a.origListener = origListener
	a.unbindReq = unbindReq

	tmp, _ := unbindReq.GetResponse()
	a.expectedResp = tmp.(*PDU.UnbindResp)

	a.waitChan = make(chan bool, 10)

	return a
}

func (c *UnbindServerPDUEventListener) HandleEvent(event *ServerPDUEvent) *Exception.Exception {
	defer func() {
		if errs := recover(); errs != nil {
		}
	}()

	pdu := event.GetPDU()
	if pdu == nil {
		return nil
	}

	if pdu.GetSequenceNumber() == c.unbindReq.GetSequenceNumber() {
		resp, err := c.session.checkResponse(pdu, c.expectedResp)
		if err != nil {
			return nil
		}

		if resp != nil {
			c.unbindResp = resp.(*PDU.UnbindResp)
			c.waitChan <- true
		}
	} else {
		if c.origListener != nil {
			c.origListener.HandleEvent(event)
		}
	}

	return nil
}

func (c *UnbindServerPDUEventListener) GetUnbindResp() *PDU.UnbindResp {
	return c.unbindResp
}

func (c *UnbindServerPDUEventListener) GetWaitChan() chan bool {
	return c.waitChan
}

func (c *UnbindServerPDUEventListener) CloseWaitChan() {
	close(c.waitChan)
}

func (c *UnbindServerPDUEventListener) StartWait(miliSecond int64) {
	timer := time.NewTimer(time.Millisecond * time.Duration(miliSecond))
	go func(timer *time.Timer, a *UnbindServerPDUEventListener) {
		defer func() {
			if errs := recover(); errs != nil {
			}

			a.GetWaitChan() <- true
		}()

		<-timer.C
	}(timer, c)
}
