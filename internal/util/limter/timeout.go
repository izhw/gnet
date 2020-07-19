package limter

import (
	"sync"
	"time"
)

type timeoutLimiter struct {
	c      chan struct{}
	t      time.Duration
	timers sync.Pool
}

func NewTimeoutLimiter(n uint32, timeout time.Duration) Limiter {
	return &timeoutLimiter{
		c: make(chan struct{}, n),
		t: timeout,
		timers: sync.Pool{
			New: func() interface{} {
				t := time.NewTimer(timeout)
				t.Stop()
				return t
			}},
	}
}

// Allow returns true if request is allowed, false if timeout
func (l *timeoutLimiter) Allow() bool {
	select {
	case l.c <- struct{}{}:
		return true
	default:
	}
	timer := l.timers.Get().(*time.Timer)
	timer.Reset(l.t)
	select {
	case l.c <- struct{}{}:
		if !timer.Stop() {
			<-timer.C
		}
		l.timers.Put(timer)
		return true
	case <-timer.C:
		l.timers.Put(timer)
		return false
	}
}

func (l *timeoutLimiter) Revert() {
	<-l.c
}
