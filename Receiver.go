package gosmpp

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/PDU"
	"github.com/linxGnu/gosmpp/Utils"
)

const (
	RECEIVER_THREAD_NAME = "Receiver"
	RECEIVE_CHAN_SIZE    = 10000
)

type IReceiver interface {
	Receive()
	TryReceivePDU(IConnection, PDU.IPDU) (PDU.IPDU, *Exception.Exception)
}

type Receiver struct {
	*receiverBase
	transmitter     *Transmitter
	connection      IConnection
	unprocessed     *Utils.Unprocessed
	pduListener     ServerPDUEventListener
	pduListenerLock sync.RWMutex
	automaticNack   bool // If true then GenericNack messages will be sent automatically if message can't be parsed
}

func NewReceiver(listener ServerPDUEventListener) (r *Receiver) {
	r = &Receiver{pduListener: listener}
	r.receiverBase = newReceiverBase(r)
	r.automaticNack = true
	r.unprocessed = Utils.NewUnprocessed()
	return
}

func NewReceiverWithConnection(listener ServerPDUEventListener, conn IConnection) (r *Receiver) {
	r = NewReceiver(listener)
	r.connection = conn
	return
}

func NewReceiverWithTransmitterCon(listener ServerPDUEventListener, trans *Transmitter, conn IConnection) (r *Receiver) {
	r = NewReceiverWithConnection(listener, conn)
	r.transmitter = trans
	return
}

// Start will reset unprocessed data and start receiving on the background.
func (c *Receiver) Start() {
	c.unprocessed.Reset()
	c.receiverBase.start()
}

// Stop receiver.
func (c *Receiver) Stop() {
	c.receiverBase.stop()
}

// StopByException stops receiver and print err log.
func (c *Receiver) StopByException(e *Exception.Exception) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e.Error)
	}
	c.Stop()
}

/**
 * This method receives a PDU or returns PDU received on background,
 * if there is any. It tries to receive a PDU for the specified timeout.
 * If the receiver is asynchronous, then no attempt to receive a PDU
 * and <code>null</code> is returned.
 * The function calls are nested as follows:<br>
 * <ul>
 *   <li>No background receiver thread<br><code>
 *       Receiver.receive(long)<br>
 *       ReceiverBase.tryReceivePDUWithTimeout(Connection,PDU,long)<br>
 *       Receiver.tryReceivePDU(Connection,PDU)<br>
 *       ReceiverBase.receivePDUFromConnection<br>
 *       Connection.receive()</code>
 *   <li>Has background receiver thread<br><code>
 *       Receiver.receive(long)<br>
 *       ReceiverBase.tryReceivePDUWithTimeout(Connection,PDU,long)<br>
 *       Receiver.tryReceivePDU(Connection,PDU)<br>
 *       Queue.dequeue(PDU)</code><br>
 *       and the ReceiverBase.run() function which actually receives the
 *       PDUs and stores them to a queue looks as follows:<br><code>
 *       ReceiverBase.run()<br>
 *       Receiver.receiveAsync()<br>
 *       ReceiverBase.receivePDUFromConnection<br>
 *       Connection.receive()</code>
 *
 * @param timeout for how long is tried to receive a PDU
 * @return the received PDU or null if none received for the spec. timeout
 *
 * @exception IOException exception during communication
 * @exception PDUException incorrect format of PDU
 * @exception TimeoutException rest of PDU not received for too long time
 * @exception UnknownCommandIdException PDU with unknown id was received
 * @see ReceiverBase#tryReceivePDUWithTimeout(Connection,PDU,long)
 */
func (c *Receiver) ReceiveSyncWTimeout(timeout int64) (PDU.IPDU, *Exception.Exception) {
	return c.tryReceivePDUWithCustomTimeout(c.connection, nil, timeout)
}

/**
 * Called from session to receive a response for previously sent request.
 *
 * @param expectedPDU the template for expected PDU; the PDU returned
 *                    must have the same sequence number
 * @return the received PDU or null if none
 * @see ReceiverBase#tryReceivePDUWithTimeout(Connection,PDU,long)
 */
func (c *Receiver) ReceiveSyncWithExpectedPDU(pdu PDU.IPDU) (PDU.IPDU, *Exception.Exception) {
	return c.tryReceivePDUWithCustomTimeout(c.connection, pdu, c.receiveTimeout)
}

func (c *Receiver) tryReceivePDUWithCustomTimeout(conn IConnection, expectedPDU PDU.IPDU, timeout int64) (pduResult PDU.IPDU, expc *Exception.Exception) {
	if timeout <= 0 {
		return c.TryReceivePDU(conn, expectedPDU)
	}

	deadLine := time.Now().Add(time.Duration(timeout) * time.Millisecond)

	var pdu PDU.IPDU

	for pdu == nil && time.Now().Before(deadLine) {
		_pdu, err := c.TryReceivePDU(conn, expectedPDU)
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

func (c *Receiver) TryReceivePDU(conn IConnection, expected PDU.IPDU) (pduResult PDU.IPDU, expc *Exception.Exception) {
	var pdu PDU.IPDU
	pdu, err := c.ReceivePDUFromConnection(c.connection, c.unprocessed)
	if err != nil {
		return pdu, err
	}

	if expected != nil && expected.IsEquals(pdu) {
		return pdu, nil
	}

	return nil, nil
}

/**
 * This method receives a PDU from connection and stores it into
 * <code>pduQueue</code>. It's called from the <code>ReceiverBase</code>'s
 * p<code>process</code> method which is called in loop from
 * <code>ProcessingThread</code>'s <code>run</code> method.
 * <p>
 * If an exception occurs during receiving, depending on type
 * of the exception this method either just reports the exception to
 * debug & event objects or stops processing to indicate
 * that it isn't able to process the exception. The function
 * <code>setTermException</code> is then called with the caught exception.
 *
 * @see ReceiverBase#run()
 */
func (c *Receiver) Receive() {
	if c.connection != nil && !c.connection.IsOpened() {
		c.Stop()
		return
	}

	pdu, err := c.ReceivePDUFromConnection(c.connection, c.unprocessed)
	if err != nil {
		if err == Exception.InvalidPDUException {
			var seqNr int32
			if pdu != nil {
				seqNr = pdu.GetSequenceNumber()
			}

			if c.automaticNack {
				c.sendGenericNack(Data.ESME_RINVMSGLEN, seqNr)
			} else {
				pdu = PDU.NewGenericNackWithCmStatusSeqNum(Data.ESME_RINVMSGLEN, seqNr)
			}
		} else if err == Exception.UnknownCommandIdException {
			var seqNr int32
			if pdu != nil {
				seqNr = pdu.GetSequenceNumber()
			}

			if c.automaticNack {
				c.sendGenericNack(Data.ESME_RINVCMDID, seqNr)
			} else {
				pdu = PDU.NewGenericNackWithCmStatusSeqNum(Data.ESME_RINVCMDID, seqNr)
			}
		} else if pdu != nil {
			if c.automaticNack {
				c.sendGenericNack(err.ErrorCode, pdu.GetSequenceNumber())
			} else {
				pdu = PDU.NewGenericNackWithCmStatusSeqNum(err.ErrorCode, pdu.GetSequenceNumber())
			}
		}
	}

	if pdu != nil {
		c.handle(pdu)
	}
}

func (c *Receiver) ReceivePDUFromConnection(conn IConnection, unprocessed *Utils.Unprocessed) (PDU.IPDU, *Exception.Exception) {
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
		} else if unprocessedBuf.Len() > 0 && time.Now().UnixNano() > c.receiveTimeout*int64(1000000)+unprocessed.GetLastTimeReceived() {
			unprocessed.Reset()
			return pdu, Exception.TimeoutException
		}
	}

	return pdu, nil
}

func (c *Receiver) tryGetUnprocessedPDU(unproc *Utils.Unprocessed) (PDU.IPDU, *Exception.Exception) {
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

func (c *Receiver) getListener() (lis ServerPDUEventListener) {
	c.pduListenerLock.RLock()
	lis = c.pduListener
	c.pduListenerLock.RUnlock()
	return
}

func (c *Receiver) setListener(lis ServerPDUEventListener) {
	c.pduListenerLock.Lock()
	c.pduListener = lis
	c.pduListenerLock.Unlock()
}

func (c *Receiver) handle(pdu PDU.IPDU) {
	if pdu == nil {
		return
	}

	if pduListener := c.getListener(); pduListener != nil {
		pduListener.HandleEvent(NewServerPDUEvent(c, c.connection, pdu))
	} else {
		t, _, _ := pdu.GetData()
		if t != nil {
			fmt.Fprintf(os.Stdout, "Receiver doesn't have ServerPDUEventListener, "+"discarding "+t.GetHexDump()+"\n")
		}
	}
}

func (c *Receiver) sendGenericNack(commandStatus, seqNum int32) {
	if c.transmitter != nil {
		gnack := PDU.NewGenericNackWithCmStatusSeqNum(commandStatus, seqNum)
		c.transmitter.Send(gnack)
	}
}
