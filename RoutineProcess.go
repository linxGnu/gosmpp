package gosmpp

import (
	"fmt"
	"os"
	"sync"
)

const (
	PROC_INITIALISING byte = 0
	PROC_RECEIVING    byte = 1
	PROC_FINISHED     byte = 2
)

type IRoutineProcess interface {
	Process()
	GetProcessStatus() byte
}

type RoutineProcess struct {
	KeepProcessing bool
	ProcessStatus  byte
	ProcessLock    sync.RWMutex
	ProcessUnit    IRoutineProcess
}

func (c *RoutineProcess) RegisterProcessUnit(proc IRoutineProcess) {
	c.ProcessUnit = proc
}

func (c *RoutineProcess) GetProcessStatus() byte {
	c.ProcessLock.RLock()
	defer c.ProcessLock.RUnlock()
	return c.ProcessStatus
}

func (c *RoutineProcess) IsKeepProcessing() bool {
	c.ProcessLock.RLock()
	defer c.ProcessLock.RUnlock()

	return c.KeepProcessing
}

func (c *RoutineProcess) StartProcess() {
	c.ProcessLock.Lock()
	defer c.ProcessLock.Unlock()

	if c.ProcessStatus != PROC_RECEIVING {
		c.ProcessStatus = PROC_RECEIVING
		c.KeepProcessing = true
		go c.Run()
	}
}

func (c *RoutineProcess) Run() {
	defer func() {
		if errs := recover(); errs != nil {
			fmt.Fprintf(os.Stderr, "%v", errs)
		}
	}()

	defer func() {
		c.ProcessLock.Lock()
		defer c.ProcessLock.Unlock()

		c.ProcessStatus = PROC_FINISHED
		c.KeepProcessing = false
	}()

	for {
		if !c.IsKeepProcessing() {
			return
		}

		if c.ProcessUnit != nil {
			c.ProcessUnit.Process()
		}
	}
}

func (c *RoutineProcess) StopProcess() {
	c.ProcessLock.Lock()
	defer c.ProcessLock.Unlock()

	c.KeepProcessing = false
}
