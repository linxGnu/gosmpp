package gosmpp

import (
	"sync/atomic"
	"time"
)

// ReceiverSession represents session for Receiver.
type ReceiverSession struct {
	dialer Dialer
	auth   Auth

	originalOnClosed func(State)
	settings         ReceiveSettings

	rebindingInterval time.Duration

	r atomic.Value // Receiver

	state     int32
	rebinding int32
}

// NewReceiverSession creates new session for Receiver.
//
// Session will `non-stop`, automatically rebind (create new and authenticate connection with SMSC) when
// unexpected error happened.
//
// `rebindingInterval` indicates duration that Session has to wait before rebinding again.
//
// Setting `rebindingInterval <= 0` will disable `auto-rebind` functionality.
func NewReceiverSession(dialer Dialer, auth Auth, settings ReceiveSettings, rebindingInterval time.Duration) (session *ReceiverSession, err error) {
	if conn, err := ConnectAsReceiver(dialer, auth); err == nil {
		session = &ReceiverSession{
			dialer:            dialer,
			auth:              auth,
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

		// create new receiver
		r := NewReceiver(conn, session.settings)

		// bind to session
		session.r.Store(r)
	}
	return
}

// Receiver returns bound Receiver.
func (s *ReceiverSession) Receiver() (r Receiver) {
	r, _ = s.r.Load().(Receiver)
	return
}

// Close session.
func (s *ReceiverSession) Close() (err error) {
	if atomic.CompareAndSwapInt32(&s.state, 0, 1) {
		// close underlying Receiver
		err = s.close()
	}
	return
}

// close underlying Receiver
func (s *ReceiverSession) close() (err error) {
	if r := s.Receiver(); r != nil {
		err = r.Close()
	}
	return
}

func (s *ReceiverSession) rebind() {
	if atomic.CompareAndSwapInt32(&s.rebinding, 0, 1) {
		// close underlying Receiver
		_ = s.close()

		for atomic.LoadInt32(&s.state) == 0 {
			conn, err := ConnectAsReceiver(s.dialer, s.auth)
			if err != nil {
				if s.settings.OnRebindingError != nil {
					s.settings.OnRebindingError(err)
				}
				time.Sleep(s.rebindingInterval)
			} else {
				r := NewReceiver(conn, s.settings)

				// bind to session
				s.r.Store(r)

				// reset rebinding state
				atomic.StoreInt32(&s.rebinding, 0)

				return
			}
		}
	}
}
