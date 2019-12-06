package gosmpp

import (
	"sync/atomic"
	"time"
)

// TransmitterSession represents session for Transmitter.
type TransmitterSession struct {
	auth Auth

	originalOnClosed func(State)
	settings         TransmitSettings

	rebindingInterval time.Duration

	r atomic.Value // Transmitter

	state     int32
	rebinding int32
}

// NewTransmitterSession creates new session for Transmitter.
//
// Session will `non-stop`, automatically rebind (create new and authenticate connection with SMSC) when
// unexpected error happened.
//
// `rebindingInterval` indicates duration that Session has to wait before rebinding again.
//
// Setting `rebindingInterval <= 0` will disable `auto-rebind` functionality.
func NewTransmitterSession(auth Auth, settings TransmitSettings, rebindingInterval time.Duration) (session *TransmitterSession, err error) {
	conn, err := ConnectAsTransmitter(auth)
	if err != nil {
		return
	}

	session = &TransmitterSession{
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

	// create new Transmitter
	r := NewTransmitter(conn, session.settings)

	// bind to session
	session.r.Store(r)

	return
}

// Transmitter returns bound Transmitter.
func (s *TransmitterSession) Transmitter() (r Transmitter) {
	r, _ = s.r.Load().(Transmitter)
	return
}

// Close session.
func (s *TransmitterSession) Close() (err error) {
	if atomic.CompareAndSwapInt32(&s.state, 0, 1) {
		// close underlying Transmitter
		err = s.close()
	}
	return
}

// close underlying Transmitter
func (s *TransmitterSession) close() (err error) {
	if r := s.Transmitter(); r != nil {
		err = r.Close()
	}
	return
}

func (s *TransmitterSession) rebind() {
	if atomic.CompareAndSwapInt32(&s.rebinding, 0, 1) {
		// close underlying Transmitter
		_ = s.close()

		for {
			conn, err := ConnectAsTransmitter(s.auth)
			if err != nil {
				if s.settings.OnRebindingError != nil {
					s.settings.OnRebindingError(err)
				}
				time.Sleep(s.rebindingInterval)
			} else {
				r := NewTransmitter(conn, s.settings)

				// bind to session
				s.r.Store(r)

				// reset rebinding state
				atomic.StoreInt32(&s.rebinding, 0)

				return
			}
		}
	}
}
