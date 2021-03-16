package gosmpp

import (
	"sync/atomic"
	"time"
)

// TransceiverSession represents session for Transceiver.
type TransceiverSession struct {
	dialer Dialer
	auth   Auth

	originalOnClosed func(State)
	settings         TransceiveSettings

	rebindingInterval time.Duration

	r atomic.Value // Transceiver

	state     int32
	rebinding int32
}

// NewTransceiverSession creates new session for Transceiver.
//
// Session will `non-stop`, automatically rebind (create new and authenticate connection with SMSC) when
// unexpected error happened.
//
// `rebindingInterval` indicates duration that Session has to wait before rebinding again.
//
// Setting `rebindingInterval <= 0` will disable `auto-rebind` functionality.
func NewTransceiverSession(dialer Dialer, auth Auth, settings TransceiveSettings, rebindingInterval time.Duration) (session *TransceiverSession, err error) {
	conn, err := ConnectAsTransceiver(dialer, auth)
	if err == nil {
		session = &TransceiverSession{
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

		// create new Transceiver
		r := newTransceiver(conn, session.settings)

		// bind to session
		session.r.Store(r)
	}
	return
}

// Transceiver returns bound Transceiver.
func (s *TransceiverSession) Transceiver() (r Transceiver) {
	r, _ = s.r.Load().(Transceiver)
	return
}

// Close session.
func (s *TransceiverSession) Close() (err error) {
	if atomic.CompareAndSwapInt32(&s.state, 0, 1) {
		// close underlying Transceiver
		err = s.close()
	}
	return
}

// close underlying Transceiver
func (s *TransceiverSession) close() (err error) {
	if r := s.Transceiver(); r != nil {
		err = r.Close()
	}
	return
}

func (s *TransceiverSession) rebind() {
	if atomic.CompareAndSwapInt32(&s.rebinding, 0, 1) {
		// close underlying Transceiver
		_ = s.close()

		for atomic.LoadInt32(&s.state) == 0 {
			conn, err := ConnectAsTransceiver(s.dialer, s.auth)
			if err != nil {
				if s.settings.OnRebindingError != nil {
					s.settings.OnRebindingError(err)
				}
				time.Sleep(s.rebindingInterval)
			} else {
				r := newTransceiver(conn, s.settings)

				// bind to session
				s.r.Store(r)

				// reset rebinding state
				atomic.StoreInt32(&s.rebinding, 0)

				return
			}
		}
	}
}
