package gosmpp

import (
	"context"
	"time"

	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/PDU"
	"github.com/linxGnu/gosmpp/Utils"
)

// base struct for receiver.
type receiverBase struct {
	ctx                         context.Context
	cancel                      context.CancelFunc
	receiveTimeout              int64
	messageIncompleteRetryCount byte
	receiver                    IReceiver
}

func newReceiverBase(receiver IReceiver) (r *receiverBase) {
	r = &receiverBase{
		receiver:       receiver,
		receiveTimeout: Data.RECEIVER_TIMEOUT,
	}
	r.ctx, r.cancel = context.WithCancel(context.Background())
	r.init()
	return r
}

func (c *receiverBase) init() {
	c.receiveTimeout = Data.RECEIVER_TIMEOUT
	c.messageIncompleteRetryCount = 0
}

func (c *receiverBase) start() {
	go func() {
		for {
			select {
			case <-c.ctx.Done():
				return

			default:
				if c.receiver != nil {
					c.receiver.Receive()
				}
			}
		}
	}()
}

func (c *receiverBase) stop() {
	c.cancel()
}

func (c *receiverBase) canContinueReceiving(deadLine time.Time, timeout int64) bool {
	if timeout == Data.RECEIVE_BLOCKING {
		return true
	}
	return time.Now().Before(deadLine)
}

func (c *receiverBase) tryReceivePDUWithTimeout(conn IConnection, expected PDU.IPDU) (PDU.IPDU, *Exception.Exception) {
	return c.tryReceivePDUWithCustomTimeout(conn, expected, c.receiveTimeout)
}

func (c *receiverBase) tryReceivePDUWithCustomTimeout(conn IConnection, expectedPDU PDU.IPDU, timeout int64) (pduResult PDU.IPDU, expc *Exception.Exception) {
	if c.receiver == nil {
		return nil, Exception.NewExceptionFromStr("Receiver not initialized")
	}

	if timeout <= 0 {
		return c.receiver.TryReceivePDU(conn, expectedPDU)
	}

	deadLine := time.Now().Add(time.Duration(timeout) * time.Millisecond)

	var pdu PDU.IPDU

	for pdu == nil && c.canContinueReceiving(deadLine, timeout) {
		_pdu, err := c.receiver.TryReceivePDU(conn, expectedPDU)
		if err != nil {
			return nil, err
		}
		pdu = _pdu

		if pdu == nil {
			time.Sleep(50 * time.Millisecond)
		}
	}

	return pdu, nil
}

func (c *receiverBase) ReceivePDUFromConnection(conn IConnection, unprocessed *Utils.Unprocessed) (PDU.IPDU, *Exception.Exception) {
	if unprocessed == nil {
		return nil, nil
	}

	var pdu PDU.IPDU // nil by default
	var unprocessedBuf *Utils.ByteBuffer

	if unprocessed.GetHasUnprocessed() {
		unprocessedBuf = unprocessed.GetUnprocessed()
		_pdu, err := c.tryGetUnprocessedPDU(unprocessed)
		if err != nil {
			return _pdu, err
		}
		pdu = _pdu
	}

	if pdu == nil {
		buffer, err := conn.Receive()
		if err != nil {
			return pdu, err
		}

		unprocessedBuf = unprocessed.GetUnprocessed()
		if buffer.Len() != 0 {
			_, err := unprocessedBuf.Write(buffer.Bytes())
			if err != nil {
				return pdu, Exception.NewException(err)
			}

			unprocessed.SetLastTimeReceivedCurTime()

			_pdu, e := c.tryGetUnprocessedPDU(unprocessed)
			if e != nil {
				if e == Exception.UnknownCommandIdException {
					unprocessed.Reset()
				}

				return pdu, e
			}

			pdu = _pdu
		} else {
			if unprocessedBuf.Len() > 0 && time.Now().UnixNano() > c.receiveTimeout*int64(1000000)+unprocessed.GetLastTimeReceived() {
				unprocessed.Reset()
				return pdu, Exception.TimeoutException
			}
		}
	}

	return pdu, nil
}

func (c *receiverBase) tryGetUnprocessedPDU(unproc *Utils.Unprocessed) (PDU.IPDU, *Exception.Exception) {
	unprocBuffer := unproc.GetUnprocessed()

	pdu, err, header := PDU.CreatePDU(unprocBuffer)
	if err == Exception.HeaderIncompleteException {
		unproc.SetHasUnprocessed(false)
		unproc.SetExpected(Data.PDU_HEADER_SIZE)

		// incomplete message header, will wait for the rest.
		return nil, nil
	} else if err == Exception.MessageIncompleteException {
		if c.messageIncompleteRetryCount > 5 {
			// Giving up on incomplete messages - probably garbage in unprocessed buffer. Flushing unprocessed buffer.
			c.messageIncompleteRetryCount = 0
			unproc.Reset()
		}

		unproc.SetHasUnprocessed(false)
		unproc.SetExpected(Data.PDU_HEADER_SIZE)
		c.messageIncompleteRetryCount++

		// incomplete message, will wait for the rest.
		return nil, nil
	} else if err == Exception.UnknownCommandIdException {
		if int(header.GetCommandLength()) <= unprocBuffer.Len() {
			_, err1 := unprocBuffer.Read_Bytes(int(header.GetCommandId()))
			if err1 != nil {
				return nil, err1
			}

			unproc.Check()
		}

		return pdu, err
	} else if err != nil {
		unproc.Check()
		return pdu, err
	}

	unproc.Check()
	c.messageIncompleteRetryCount = 0

	return pdu, nil
}
