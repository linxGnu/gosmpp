package daemons

import (
	"sync"
)

var lock sync.RWMutex
var shouldStop bool

var wg sync.WaitGroup

// RunDaemons run all daemons
func RunDaemons() {
	shouldStop = false

	wg.Add(1)
	go SmsSenderDaemon()
}

// StopDaemons stop all daemons
func StopDaemons() {
	lock.Lock()
	shouldStop = true
	lock.Unlock()

	wg.Wait()
}
