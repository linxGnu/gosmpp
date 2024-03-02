package gosmpp

import (
	"errors"
	"fmt"
	"github.com/linxGnu/gosmpp/pdu"
	"sync/atomic"
	"time"
)

var (
	ErrRequestWindowIsNil     = errors.New("request window is nil")
	ErrWindowSizeEqualZero    = errors.New("request window size cannot be 0")
	ErrExpireCheckTimerNotSet = errors.New("ExpireCheckTimer cannot be 0 if PduExpireTimeOut is set")
)

// Session represents session for TX, RX, TRX.
type Session struct {
	c Connector

	originalOnClosed func(State)
	settings         Settings

	rebindingInterval time.Duration

	trx atomic.Value // transceivable

	state     int32
	rebinding int32
}

// NewSession creates new session for TX, RX, TRX.
//
// Session will `non-stop`, automatically rebind (create new and authenticate connection with SMSC) when
// unexpected error happened.
//
// `rebindingInterval` indicates duration that Session has to wait before rebinding again.
//
// Setting `rebindingInterval <= 0` will disable `auto-rebind` functionality.
func NewSession(c Connector, settings Settings, rebindingInterval time.Duration) (session *Session, err error) {
	if settings.ReadTimeout <= 0 || settings.ReadTimeout <= settings.EnquireLink {
		return nil, fmt.Errorf("invalid settings: ReadTimeout must greater than max(0, EnquireLink)")
	}
	if settings.RequestWindowConfig != nil {
		if settings.RequestStore == nil {
			return nil, ErrRequestWindowIsNil
		}
		if settings.MaxWindowSize == 0 {
			return nil, ErrWindowSizeEqualZero
		}
		if settings.PduExpireTimeOut > 0 && settings.ExpireCheckTimer == 0 {
			return nil, ErrExpireCheckTimerNotSet
		}
	}

	conn, err := c.Connect()
	if err == nil {
		session = &Session{
			c:                 c,
			rebindingInterval: rebindingInterval,
			originalOnClosed:  settings.OnClosed,
		}

		if rebindingInterval > 0 {
			newSettings := settings
			newSettings.OnClosed = func(state State) {
				switch state {
				case ExplicitClosing:
					return

				default:
					if session.originalOnClosed != nil {
						session.originalOnClosed(state)
					}
					session.rebind()
				}
			}
			session.settings = newSettings
		} else {
			session.settings = settings
		}

		// bind to session
		trans := newTransceivable(conn, session.settings)
		trans.start()
		session.trx.Store(trans)
	}
	return
}

func (s *Session) bound() *transceivable {
	r, _ := s.trx.Load().(*transceivable)
	return r
}

// Transmitter returns bound Transmitter.
func (s *Session) Transmitter() Transmitter {
	return s.bound()
}

// Receiver returns bound Receiver.
func (s *Session) Receiver() Receiver {
	return s.bound()
}

// Transceiver returns bound Transceiver.
func (s *Session) Transceiver() Transceiver {
	return s.bound()
}

func (s *Session) GetWindowSize() int {
	if s.c.GetBindType() == pdu.Transmitter || s.c.GetBindType() == pdu.Transceiver {
		return s.bound().GetWindowSize()
	}
	return -1
}

// Close session.
func (s *Session) Close() (err error) {
	if atomic.CompareAndSwapInt32(&s.state, Alive, Closed) {
		err = s.close()
	}
	return
}

func (s *Session) close() (err error) {
	if b := s.bound(); b != nil {
		err = b.Close()
	}
	return
}

func (s *Session) rebind() {
	if atomic.CompareAndSwapInt32(&s.rebinding, 0, 1) {
		_ = s.close()

		for atomic.LoadInt32(&s.state) == Alive {
			conn, err := s.c.Connect()
			if err != nil {
				if s.settings.OnRebindingError != nil {
					s.settings.OnRebindingError(err)
				}
				time.Sleep(s.rebindingInterval)
			} else {
				// bind to session
				trans := newTransceivable(conn, s.settings)
				trans.start()
				s.trx.Store(trans)

				// reset rebinding state
				atomic.StoreInt32(&s.rebinding, 0)
				if s.settings.OnRebind != nil {
					s.settings.OnRebind()
				}

				return
			}
		}
	}
}
