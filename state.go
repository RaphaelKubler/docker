package docker

import (
	"fmt"
	"sync"
	"time"
)

type State struct {
	Running   bool
	Pid       int
	ExitCode  int
	StartedAt time.Time

	stateChangeLock *sync.Mutex
	stateChangeCond *sync.Cond
}

// String returns a human-readable description of the state
func (s *State) String() string {
	if s.Running {
		return fmt.Sprintf("Up %s", HumanDuration(time.Now().Sub(s.StartedAt)))
	}
	return fmt.Sprintf("Exit %d", s.ExitCode)
}

func (s *State) setRunning(pid int) {
	s.Running = true
	s.ExitCode = 0
	s.Pid = pid
	s.StartedAt = time.Now()
	s.broadcast()
}

func (s *State) setStopped(exitCode int) {
	s.Running = false
	s.Pid = 0
	s.ExitCode = exitCode
	s.broadcast()
}

func (s *State) initLock() {
	if s.stateChangeLock == nil {
		s.stateChangeLock = &sync.Mutex{}
		s.stateChangeCond = sync.NewCond(s.stateChangeLock)
	}
}

func (s *State) broadcast() {
	s.stateChangeLock.Lock()
	s.stateChangeCond.Broadcast()
	s.stateChangeLock.Unlock()
}

func (s *State) wait() {
	s.stateChangeLock.Lock()
	s.stateChangeCond.Wait()
	s.stateChangeLock.Unlock()
}
