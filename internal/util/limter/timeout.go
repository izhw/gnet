// Copyright (c) 2020 izhw
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
